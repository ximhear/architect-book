// TimescaleDB IoT Ingester — 책 2권 5장 구현.
//
// Kafka (NCP CDSS) → batch buffer → TimescaleDB hypertable.
// Batch 윈도우: 1,000 row 또는 100ms 중 빠른 것.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

// MeterReading — 스마트미터 한 record. 책 2권 케이스 2.2 의 컬럼 구성.
type MeterReading struct {
	HouseholdID string    `json:"household_id"`
	ReadingAt   time.Time `json:"reading_at"`
	Kwh         float64   `json:"kwh"`
	Source      string    `json:"source,omitempty"`
}

// Ingester — batch buffer + flush.
type Ingester struct {
	pool         *pgxpool.Pool
	buffer       []MeterReading
	batchSize    int
	flushEvery   time.Duration
	totalIn      atomic.Uint64
	totalWritten atomic.Uint64
	totalDropped atomic.Uint64
}

func NewIngester(pool *pgxpool.Pool) *Ingester {
	return &Ingester{
		pool:       pool,
		buffer:     make([]MeterReading, 0, 1024),
		batchSize:  1000,
		flushEvery: 100 * time.Millisecond,
	}
}

func (i *Ingester) Add(r MeterReading) {
	i.totalIn.Add(1)
	i.buffer = append(i.buffer, r)
}

// Flush — buffer 를 TimescaleDB 로 batch INSERT.
// ON CONFLICT DO NOTHING — at-least-once 컨슈머의 중복 방어.
func (i *Ingester) Flush(ctx context.Context) error {
	if len(i.buffer) == 0 {
		return nil
	}
	batch := i.buffer
	i.buffer = make([]MeterReading, 0, 1024)

	rows := make([][]any, 0, len(batch))
	for _, r := range batch {
		if r.Source == "" {
			r.Source = "gateway"
		}
		rows = append(rows, []any{r.HouseholdID, r.ReadingAt, r.Kwh, r.Source})
	}

	tag, err := i.pool.CopyFrom(ctx,
		[]string{"meter_readings"},
		[]string{"household_id", "reading_at", "kwh", "source"},
		pgxRows(rows),
	)
	if err != nil {
		// 운영: 일부 partition / 제약 위반은 row 단위 fallback. 본 예제는 batch 전체 실패 시 dropped.
		i.totalDropped.Add(uint64(len(batch)))
		return err
	}
	i.totalWritten.Add(uint64(tag))
	return nil
}

// metrics — stdout. 운영은 OTel exporter.
func (i *Ingester) LogMetrics() {
	slog.Info("ingest metrics",
		"received", i.totalIn.Load(),
		"written", i.totalWritten.Load(),
		"dropped", i.totalDropped.Load(),
		"buffer_size", len(i.buffer),
	)
}

// pgxRows — [][]any 를 pgx.CopyFromSource 어댑터로.
type pgxRows [][]any

func (r pgxRows) Next() bool      { return len(r) > 0 }
func (r pgxRows) Err() error      { return nil }
func (r pgxRows) Values() ([]any, error) {
	row := r[0]
	r = r[1:]
	return row, nil
}

// stdin 모드 — Kafka 없이 학습 시 한 줄씩 JSON 입력.
func runStdin(ctx context.Context, ing *Ingester) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var r MeterReading
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			slog.Warn("invalid json", "line", scanner.Text())
			continue
		}
		ing.Add(r)
		if len(ing.buffer) >= ing.batchSize {
			if err := ing.Flush(ctx); err != nil {
				slog.Error("flush failed", "err", err)
			}
		}
	}
	return ing.Flush(ctx)
}

// Kafka 모드 — 운영.
func runKafka(ctx context.Context, ing *Ingester, brokers []string, topic string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "ingester",
		MinBytes:       1,
		MaxBytes:       10 << 20,
		CommitInterval: time.Second,
	})
	defer r.Close()

	timer := time.NewTimer(ing.flushEvery)
	for {
		select {
		case <-ctx.Done():
			return ing.Flush(ctx)
		case <-timer.C:
			if err := ing.Flush(ctx); err != nil {
				slog.Error("flush failed", "err", err)
			}
			timer.Reset(ing.flushEvery)
		default:
		}

		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return ing.Flush(ctx)
			}
			slog.Error("kafka read", "err", err)
			continue
		}
		var rr MeterReading
		if err := json.Unmarshal(m.Value, &rr); err != nil {
			slog.Warn("invalid json", "topic", topic, "offset", m.Offset)
			continue
		}
		ing.Add(rr)
		if len(ing.buffer) >= ing.batchSize {
			if err := ing.Flush(ctx); err != nil {
				slog.Error("flush failed", "err", err)
			}
			timer.Reset(ing.flushEvery)
		}
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	dsn := env("DATABASE_URL", "postgres://postgres:pass@localhost/postgres?sslmode=disable")
	stdinMode := os.Getenv("STDIN_MODE") == "1"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("db pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	ing := NewIngester(pool)

	// 메트릭 로깅 — 10초 주기
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				ing.LogMetrics()
			}
		}
	}()

	if stdinMode {
		slog.Info("stdin mode — paste JSON readings")
		if err := runStdin(ctx, ing); err != nil {
			slog.Error("stdin", "err", err)
		}
		return
	}

	brokers := []string{env("KAFKA_BROKERS", "localhost:9092")}
	topic := env("KAFKA_TOPIC", "meter.readings")
	slog.Info("kafka mode", "brokers", brokers, "topic", topic)
	if err := runKafka(ctx, ing, brokers, topic); err != nil {
		slog.Error("kafka", "err", err)
	}
	ing.LogMetrics()
}
