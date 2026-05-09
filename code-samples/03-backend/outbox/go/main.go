// Outbox poller — 책 3장 4.2 의 별도 프로세스 구현.
//
// Spring Boot 측 OutboxPublisher 가 단일 노드 가정인 반면,
// 본 Go 구현은 다중 인스턴스 동시 실행 안전 (SELECT FOR UPDATE SKIP LOCKED).
//
// 메시지 브로커는 stdout 로 시뮬레이션 — 실제는 Kafka / SQS / NCP CDSS 로 교체.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

// OutboxEvent 는 outbox 테이블의 한 행.
// 책 3장 4.2 의 컬럼 구성과 일치.
type OutboxEvent struct {
	ID            string
	AggregateID   string
	AggregateType string
	Topic         string
	PartitionKey  string
	Payload       json.RawMessage
	TraceID       sql.NullString
	CreatedAt     time.Time
}

// Publisher 는 메시지 브로커 추상.
type Publisher interface {
	Publish(ctx context.Context, ev OutboxEvent) error
}

// StdoutPublisher 는 학습용 stdout 출력. 실제는 Kafka / SQS 등 구현체.
type StdoutPublisher struct{}

func (s *StdoutPublisher) Publish(ctx context.Context, ev OutboxEvent) error {
	lag := time.Since(ev.CreatedAt).Milliseconds()
	traceID := "-"
	if ev.TraceID.Valid {
		traceID = ev.TraceID.String
	}
	slog.Info("PUBLISH",
		"topic", ev.Topic,
		"partition", ev.PartitionKey,
		"lag_ms", lag,
		"trace", traceID,
		"payload", string(ev.Payload),
	)
	return nil
}

// pollAndPublish 는 미발행 outbox 를 다중 인스턴스 안전하게 잠금 후 발행.
// 핵심: SELECT FOR UPDATE SKIP LOCKED 로 다른 인스턴스가 잠근 행은 건너뛰고 처리.
func pollAndPublish(ctx context.Context, db *sql.DB, pub Publisher, batchSize int) (int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	// rollback 은 commit 후 에러 안 남.
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id, aggregate_id, aggregate_type, topic, partition_key, payload, trace_id, created_at
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, batchSize)
	if err != nil {
		return 0, fmt.Errorf("select unpublished: %w", err)
	}

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(
			&e.ID, &e.AggregateID, &e.AggregateType,
			&e.Topic, &e.PartitionKey, &e.Payload, &e.TraceID, &e.CreatedAt,
		); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan: %w", err)
		}
		events = append(events, e)
	}
	rows.Close()

	now := time.Now().UTC()
	for _, ev := range events {
		if err := pub.Publish(ctx, ev); err != nil {
			// 운영: 발행 실패 시 retry / DLQ 정책. 여기서는 트랜잭션 롤백으로 다음 batch 에서 재시도.
			return 0, fmt.Errorf("publish event %s: %w", ev.ID, err)
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE outbox SET published_at = $1 WHERE id = $2`, now, ev.ID); err != nil {
			return 0, fmt.Errorf("update published_at: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return len(events), nil
}

func mustEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func main() {
	dsn := mustEnv("DATABASE_URL", "postgres://localhost/onepick?sslmode=disable")
	interval, _ := time.ParseDuration(mustEnv("POLL_INTERVAL", "1s"))
	batchSize := 100

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("db open failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("db ping failed", "err", err)
		os.Exit(1)
	}

	pub := &StdoutPublisher{}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("outbox poller started",
		"interval", interval,
		"batch_size", batchSize,
	)
	defer slog.Info("outbox poller stopped (graceful)")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := pollAndPublish(ctx, db, pub, batchSize)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				slog.Error("poll failed", "err", err)
				continue
			}
			if n > 0 {
				slog.Info("batch published", "count", n)
			}
		}
	}
}
