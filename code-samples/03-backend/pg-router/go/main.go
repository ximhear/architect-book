// PG 라우터 — 한국 결제 PG 3사 (토스 / KCP / NICE) 자동 fallback.
// 책 3장 8.5 + 6장 4.1 / 4.2 의 다중화 패턴 구현.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ChargeRequest — 결제 요청 한 건.
type ChargeRequest struct {
	OrderID        int64  `json:"orderId"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotencyKey"`
}

// ChargeResult — 결제 성공 한 건.
type ChargeResult struct {
	PaymentID string `json:"paymentId"`
	PG        string `json:"pg"`
}

// PG 인터페이스 — 토스 / KCP / NICE 모두 같은 시그니처.
type PG interface {
	Name() string
	Charge(ctx context.Context, req ChargeRequest) (ChargeResult, error)
}

// circuitBreaker — 단순 구현. 운영은 gobreaker 등 라이브러리 권장.
type circuitBreaker struct {
	mu               sync.Mutex
	consecutiveFails int
	openUntil        time.Time
	threshold        int
	openDuration     time.Duration
}

func newBreaker(threshold int, openDuration time.Duration) *circuitBreaker {
	return &circuitBreaker{threshold: threshold, openDuration: openDuration}
}

func (b *circuitBreaker) Healthy() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return time.Now().After(b.openUntil)
}

func (b *circuitBreaker) RecordSuccess() {
	b.mu.Lock()
	b.consecutiveFails = 0
	b.mu.Unlock()
}

func (b *circuitBreaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.consecutiveFails++
	if b.consecutiveFails >= b.threshold {
		b.openUntil = time.Now().Add(b.openDuration)
		b.consecutiveFails = 0
	}
}

// fakePG — 50% 확률로 실패하는 학습용 PG. 실제는 PG 사 SDK / HTTP 호출.
type fakePG struct {
	name      string
	successOn float64 // 0~1
	breaker   *circuitBreaker
	tries     atomic.Uint64
	successes atomic.Uint64
}

func newFakePG(name string, successOn float64) *fakePG {
	return &fakePG{
		name:      name,
		successOn: successOn,
		breaker:   newBreaker(3, 30*time.Second), // 연속 3회 실패 → 30초 open
	}
}

func (p *fakePG) Name() string { return p.name }

func (p *fakePG) Charge(ctx context.Context, req ChargeRequest) (ChargeResult, error) {
	if !p.breaker.Healthy() {
		return ChargeResult{}, fmt.Errorf("%s breaker open", p.name)
	}
	p.tries.Add(1)
	// 시뮬: PG 호출 100~300ms
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	if rand.Float64() > p.successOn {
		p.breaker.RecordFailure()
		return ChargeResult{}, fmt.Errorf("%s charge failed (simulated)", p.name)
	}
	p.breaker.RecordSuccess()
	p.successes.Add(1)
	return ChargeResult{
		PaymentID: fmt.Sprintf("%s_%d_%d", p.name, req.OrderID, time.Now().UnixNano()),
		PG:        p.name,
	}, nil
}

// PGRouter — 순차 fallback.
type PGRouter struct {
	pgs []PG
}

func (r *PGRouter) Charge(ctx context.Context, req ChargeRequest) (ChargeResult, error) {
	var lastErr error
	for _, pg := range r.pgs {
		res, err := pg.Charge(ctx, req)
		if err == nil {
			slog.Info("charge success", "pg", pg.Name(), "order", req.OrderID, "payment", res.PaymentID)
			return res, nil
		}
		slog.Warn("charge failed — falling back", "pg", pg.Name(), "order", req.OrderID, "err", err)
		lastErr = err
	}
	return ChargeResult{}, fmt.Errorf("all PGs failed: %w", lastErr)
}

// idempotency 캐시 — 같은 키면 같은 결과 반환.
type idemCache struct {
	mu sync.RWMutex
	m  map[string]ChargeResult
}

func newIdemCache() *idemCache { return &idemCache{m: make(map[string]ChargeResult)} }

func (c *idemCache) Get(key string) (ChargeResult, bool) {
	c.mu.RLock()
	r, ok := c.m[key]
	c.mu.RUnlock()
	return r, ok
}

func (c *idemCache) Put(key string, r ChargeResult) {
	c.mu.Lock()
	c.m[key] = r
	c.mu.Unlock()
}

func main() {
	// 토스 95% / KCP 90% / NICE 85% — 시뮬 성공률 (실제 PG 가용성 아님)
	router := &PGRouter{pgs: []PG{
		newFakePG("toss", 0.95),
		newFakePG("kcp", 0.90),
		newFakePG("nice", 0.85),
	}}
	cache := newIdemCache()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /charge", func(w http.ResponseWriter, r *http.Request) {
		var req ChargeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if k := r.Header.Get("Idempotency-Key"); k != "" {
			req.IdempotencyKey = k
			if cached, ok := cache.Get(k); ok {
				w.Header().Set("X-Idempotent-Replayed", "true")
				_ = json.NewEncoder(w).Encode(cached)
				return
			}
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		res, err := router.Charge(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if req.IdempotencyKey != "" {
			cache.Put(req.IdempotencyKey, res)
		}
		_ = json.NewEncoder(w).Encode(res)
	})

	addr := ":8082"
	slog.Info("PG router listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server", "err", err)
	}
}

// 사용되지 않는 errors import 회피용 (단순화 목적)
var _ = errors.New
