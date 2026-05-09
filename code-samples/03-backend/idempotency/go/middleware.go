// Idempotency middleware — net/http 표준.
//
// 같은 Idempotency-Key 로 들어온 요청은 첫 응답을 그대로 재반환.
// 운영에서는 Redis SETNX 로 처리 중 표시 + 짧은 lock TTL 권장.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// CachedResponse 는 보존되는 응답 한 건.
type CachedResponse struct {
	Status    int
	Body      []byte
	ExpiresAt time.Time
}

// Store 는 학습용 인메모리. 운영은 Redis.
type Store struct {
	mu sync.RWMutex
	m  map[string]CachedResponse
}

func NewStore() *Store { return &Store{m: make(map[string]CachedResponse)} }

func (s *Store) Get(key string) (CachedResponse, bool) {
	s.mu.RLock()
	c, ok := s.m[key]
	s.mu.RUnlock()
	if !ok || time.Now().After(c.ExpiresAt) {
		return CachedResponse{}, false
	}
	return c, true
}

func (s *Store) Put(key string, c CachedResponse) {
	s.mu.Lock()
	s.m[key] = c
	s.mu.Unlock()
}

// captureWriter 는 응답을 캡처하면서 실제 ResponseWriter 에도 쓴다.
type captureWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (c *captureWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

func (c *captureWriter) Write(b []byte) (int, error) {
	c.buf.Write(b)
	return c.ResponseWriter.Write(b)
}

// IdempotencyMiddleware — 결제 경로 등에 적용.
// TTL 은 도메인별 차등 (책 3장 4.3) — 본 예제는 결제 가정으로 90일.
func IdempotencyMiddleware(store *Store, ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			if cached, ok := store.Get(key); ok {
				slog.Info("idempotency hit", "key", key, "status", cached.Status)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Idempotent-Replayed", "true")
				w.WriteHeader(cached.Status)
				_, _ = w.Write(cached.Body)
				return
			}

			cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(cw, r)

			store.Put(key, CachedResponse{
				Status:    cw.status,
				Body:      cw.buf.Bytes(),
				ExpiresAt: time.Now().Add(ttl),
			})
			slog.Info("idempotency cached", "key", key, "status", cw.status, "ttl", ttl)
		})
	}
}

// 데모 핸들러 — 매번 다른 paymentId 를 발급.
func chargeHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"paymentId": fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		"status":    "SUCCESS",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	store := NewStore()
	mw := IdempotencyMiddleware(store, 90*24*time.Hour) // 결제 — 90일

	mux := http.NewServeMux()
	mux.Handle("POST /payments", mw(http.HandlerFunc(chargeHandler)))

	addr := ":8081"
	slog.Info("listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server", "err", err)
	}
}
