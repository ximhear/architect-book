# Outbox 패턴 예제 — Spring Boot (Kotlin) + Go

> 본 책 [3장 4.2](../../../manuscript/03-백엔드/README.md) 의 Outbox + Event 패턴을 두 언어로 구현한 학습용 최소 예제.

## 무엇을 보여주는가

핵심 보장: **DB 변경과 메시지 발행의 원자성** — 둘 중 하나만 성공하는 경우가 없게.

| 컴포넌트 | 역할 | 언어 |
|---------|------|------|
| **생산자** | `Order` INSERT 와 `outbox` INSERT 를 같은 트랜잭션 안에서 | Spring Boot (Kotlin) |
| **Poller** | 미발행 outbox 를 읽어 메시지 브로커에 발행 (예제는 stdout) | Go |
| **DB** | 두 컴포넌트가 공유하는 outbox 테이블 | PostgreSQL (Go) / H2 (Spring 학습용) |

## 책의 어느 절을 구현하는가

- **3장 4.2** Outbox + Event 패턴의 SQL 컬럼 구성 (`aggregate_id` / `aggregate_type` / `partition_key` / `trace_id` 등)
- **3장 4.2** "보상의 보상" 의 1계층 (DLQ — 본 예제는 트랜잭션 롤백으로 대체)
- **3장 8.6** 백엔드 관측성 — `outbox_lag_seconds` / 처리 batch count 로깅

## 디렉토리 구조

```
outbox/
├── README.md                    # 이 파일
├── kotlin-spring/               # 생산자
│   ├── build.gradle.kts
│   ├── settings.gradle.kts
│   └── src/main/
│       ├── kotlin/com/onepick/outbox/
│       │   ├── OutboxApplication.kt
│       │   ├── domain/
│       │   │   ├── Order.kt
│       │   │   └── OutboxEvent.kt
│       │   ├── application/
│       │   │   └── OrderService.kt
│       │   ├── port/
│       │   │   └── Repositories.kt
│       │   └── adapter/
│       │       └── OrderController.kt
│       └── resources/
│           └── application.yml
├── go/                          # Poller (별도 프로세스)
│   ├── go.mod
│   └── main.go
└── sql/
    └── schema.sql               # PostgreSQL 스키마 (Go 측 전제)
```

## 실행

### Spring Boot (생산자) — H2 메모리 DB 로 빠르게

```bash
cd kotlin-spring
./gradlew bootRun
```

다른 터미널에서:

```bash
curl -X POST -H 'Content-Type: application/json' \
  -d '{"amount": 50000}' \
  http://localhost:8080/orders
```

H2 콘솔 (http://localhost:8080/h2-console, JDBC URL `jdbc:h2:mem:outbox`) 에서 `OUTBOX` 테이블 확인. 동일 프로세스의 `OutboxPublisher` 가 1초마다 발행 후 `published_at` 갱신.

### Go (별도 Poller) — PostgreSQL

```bash
cd go
psql -d onepick -f ../sql/schema.sql
DATABASE_URL='postgres://localhost/onepick?sslmode=disable' go run .
```

이 시나리오에서는 같은 `outbox` 테이블에 다른 프로세스가 INSERT 했다고 가정. 직접 `INSERT` 후 poller 가 발행하는지 확인.

```sql
INSERT INTO outbox (id, aggregate_id, aggregate_type, topic, partition_key, payload, trace_id)
VALUES (gen_random_uuid(), '12345', 'Order', 'OrderCreated', '12345',
        '{"orderId":12345,"amount":50000}', 'trace-abc');
```

## 운영 / 학습용 차이

| 영역 | 본 예제 (학습용) | 운영 |
|------|----------------|------|
| DB | H2 메모리 (Spring) / 단일 PostgreSQL (Go) | Aurora MySQL / PostgreSQL |
| 메시지 브로커 | stdout 로그 | Kafka (MSK) / SQS / NCP CDSS |
| Poller 동시성 | 단일 프로세스 / Spring 스케줄 | 다중 인스턴스 + `SELECT FOR UPDATE SKIP LOCKED` (Go 예제 적용) |
| 발행 후 처리 | `published_at` 갱신 | 별도 `outbox_archive` 테이블로 이동 또는 TTL 파기 |
| CDC 대안 | 본 예제는 polling | Debezium 으로 outbox 직접 캡처 (poller 불필요) |
| 관측성 | 로그만 | OpenTelemetry trace + `outbox_lag_seconds` 메트릭 (3장 8.6) |

## 한국 환경 메모

- **시크릿** — DATABASE_URL 등은 환경변수가 아닌 AWS Secrets Manager / Vault 가 운영 표준 (6장 9.1)
- **Kafka 대체** — NCP Cloud Data Streaming Service 도 동일 프로토콜 호환
- **다중 인스턴스 Poller** — 한국 클라우드의 일부 RDB 가 `SKIP LOCKED` 미지원 / 제한이 있을 수 있어 사전 검증

## 함정

> - **트랜잭션 분리** — `OrderService.createOrder` 가 `@Transactional` 없으면 Order INSERT 만 commit + outbox 누락. 본 예제는 의도적으로 명시.
> - **단일 인스턴스 가정** — Spring 측 `OutboxPublisher` 는 단일 노드 가정. 다중 노드 운영 시 잠금 필요 (Go 측처럼).
> - **payload 큰 경우** — Kafka 1MB 기본 제한 초과. 큰 페이로드는 S3 포인터로.
> - **published_at 만 표시 후 보존** — 무한 누적. 운영에서는 일 단위 archive / 파기.

## 라이선스

[MIT](../../LICENSE)
