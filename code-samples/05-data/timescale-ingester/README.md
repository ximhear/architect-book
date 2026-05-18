# TimescaleDB IoT Ingester — 2권 5장

> 책 2권 5장 (TimescaleDB / 시계열 다운샘플링) 의 동작하는 최소 예제.
> Go 단일 프로세스 — Kafka (NCP CDSS 호환) 컨슈머 → TimescaleDB hypertable INSERT.

## 무엇을 보여주는가

- **Kafka 컨슈머 (Go)** — IoT 게이트웨이 → CDSS → 우리 ingester
- **TimescaleDB hypertable** — 시간·가구 인덱스 + 자동 chunk 분할
- **Batch INSERT 최적화** — 1,000 row / 100ms 윈도우, ON CONFLICT DO NOTHING
- **Continuous Aggregate** — 15분 raw → 1시간 평균 자동 갱신
- **손실률 모니터링** — 처리 카운트 + lag 메트릭

## 책 인용

- 2권 3장 2절 — IoT ingest 부하 특성 (일 2,700만 record)
- 2권 5장 3절 — Continuous Aggregate 패턴
- 2권 [케이스 11.4](../../../manuscript-2-energy/case-study/README.md#114-데이터--ai) — TimescaleDB + ClickHouse 통합

## 디렉토리

```
timescale-ingester/
├── README.md
├── go/
│   ├── go.mod
│   └── main.go
└── sql/
    └── schema.sql
```

## 실행

```bash
# 1. TimescaleDB 준비 (Docker)
docker run -d --name ts -p 5432:5432 -e POSTGRES_PASSWORD=pass timescale/timescaledb-ha:pg16

# 2. 스키마
psql -h localhost -U postgres -d postgres -f sql/schema.sql

# 3. ingester 실행
cd go
DATABASE_URL='postgres://postgres:pass@localhost/postgres?sslmode=disable' \
KAFKA_BROKERS='localhost:9092' \
KAFKA_TOPIC='meter.readings' \
go run .

# 4. 테스트 메시지 (Kafka 없을 때 stdin 모드)
echo '{"household_id":"H001","reading_at":"2026-05-18T10:15:00Z","kwh":0.42}' | \
  STDIN_MODE=1 go run .
```

## 운영 / 학습 차이

| 영역 | 본 예제 | 운영 |
|------|--------|------|
| Kafka 클라이언트 | 단순 (segmentio/kafka-go) | NCP CDSS SDK + ACL |
| 인증 | 없음 | mTLS (게이트웨이) + SASL (Kafka) |
| 다중 인스턴스 | 단일 | partition 별 consumer group |
| 손실률 측정 | 로그 | OTel 메트릭 → Grafana Cloud |
| ClickHouse 연계 | 없음 | dbt + 일 단위 익스포트 |

## 라이선스

[MIT](../../LICENSE)
