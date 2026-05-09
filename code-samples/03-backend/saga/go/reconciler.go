// 일 단위 reconciler — 책 3장 4.2 의 "보상의 보상" 3계층 중 3계층.
//
// 야간에 실행되어 "주문 vs PG 결제 vs 재고" 3자 정합을 검증.
// 어긋난 건 운영 티켓 + 알람 (실제로는 PagerDuty / Slack 으로).
//
// 본 예제는 stdout 로그로 시뮬레이션.
package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Mismatch 는 정합 불일치 한 건.
type Mismatch struct {
	OrderID         int64
	OrderStatus     string
	PaymentStatus   sql.NullString
	InventoryStatus sql.NullString
	Reason          string
}

// findMismatches 는 다음 케이스를 잡는다:
//
//  1. 주문 PAID 인데 결제 PG 응답 없음 (결제 누락)
//  2. 결제 SUCCESS 인데 주문 CANCELLED (보상 실패)
//  3. 주문 CREATED 인데 24시간 이상 진행 없음 (분실)
//  4. 재고 reserved 인데 주문 CANCELLED 24시간 이상 (보상 누락)
func findMismatches(ctx context.Context, db *sql.DB, since time.Time) ([]Mismatch, error) {
	// 본 예제는 SQL 한 개로 단순화. 실제는 4가지 케이스를 별도 쿼리로.
	rows, err := db.QueryContext(ctx, `
		SELECT
			o.id,
			o.status,
			p.status,
			r.status,
			CASE
				WHEN o.status = 'PAID'      AND p.status IS NULL                  THEN '주문 PAID 인데 결제 기록 없음'
				WHEN o.status = 'CANCELLED' AND p.status = 'SUCCESS'              THEN '결제 SUCCESS 인데 주문 CANCELLED (환불 누락 의심)'
				WHEN o.status = 'CREATED'   AND o.created_at < $1                 THEN '주문 CREATED 24h+ 진행 없음 (Saga 분실)'
				WHEN o.status = 'CANCELLED' AND r.status = 'reserved'             THEN '주문 CANCELLED 인데 재고 reserved (보상 실패)'
			END AS reason
		FROM orders o
		LEFT JOIN payments p   ON p.order_id = o.id
		LEFT JOIN reservations r ON r.order_id = o.id
		WHERE
			(o.status = 'PAID'      AND p.status IS NULL)
		 OR (o.status = 'CANCELLED' AND p.status = 'SUCCESS')
		 OR (o.status = 'CREATED'   AND o.created_at < $1)
		 OR (o.status = 'CANCELLED' AND r.status = 'reserved')
	`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Mismatch
	for rows.Next() {
		var m Mismatch
		if err := rows.Scan(&m.OrderID, &m.OrderStatus, &m.PaymentStatus, &m.InventoryStatus, &m.Reason); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func main() {
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	hours := flag.Int("hours", 24, "검사 기준 시간 (시간 단위)")
	flag.Parse()

	if *dsn == "" {
		*dsn = "postgres://localhost/onepick?sslmode=disable"
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		slog.Error("db open", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	since := time.Now().UTC().Add(-time.Duration(*hours) * time.Hour)
	slog.Info("reconciler started", "since", since)

	mismatches, err := findMismatches(context.Background(), db, since)
	if err != nil {
		slog.Error("query failed", "err", err)
		os.Exit(1)
	}

	if len(mismatches) == 0 {
		slog.Info("no mismatches")
		return
	}

	for _, m := range mismatches {
		// 운영: PagerDuty / Slack / Jira 자동 티켓
		slog.Warn("MISMATCH",
			"order", m.OrderID,
			"order_status", m.OrderStatus,
			"payment_status", m.PaymentStatus.String,
			"inventory_status", m.InventoryStatus.String,
			"reason", m.Reason,
		)
	}
	slog.Info("reconciler finished", "mismatches", len(mismatches))
}
