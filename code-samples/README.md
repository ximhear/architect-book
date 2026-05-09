# 코드 예제

> 책의 핵심 패턴을 동작하는 최소 코드로 보여주는 예제 모음.
> 라이선스: [MIT](LICENSE) — 영리 목적 포함 자유 활용.

## 디렉토리 구조

```
code-samples/
├── LICENSE                        # MIT
├── README.md                      # 이 파일
└── 03-backend/                    # 챕터별
    └── outbox/                    # 패턴별
        ├── README.md              # 예제 안내
        ├── kotlin-spring/         # Spring Boot (Kotlin) 구현
        ├── go/                    # Go 구현 (별도 poller)
        └── sql/                   # 스키마
```

## 기본 언어 ([ADR-0005](../adr/0005-코드-예제-기본-언어.md))

| 영역 | 언어 | 용도 |
|------|------|------|
| 도메인 / 트랜잭션 / 결제 | **Spring Boot (Kotlin)** | 무거운 비즈니스 |
| 게이트웨이 / Worker / Poller / Sidecar | **Go** | 가볍고 빠른 인프라 서비스 |
| 프론트엔드 (TBD) | TypeScript / Next.js | 향후 추가 |
| 데이터 / ML (TBD) | Python / FastAPI | 향후 추가 |

## 현재 예제

| 챕터 | 패턴 | 위치 | 언어 |
|------|------|------|------|
| 3장 4.2 | Outbox 패턴 (DB-이벤트 원자성) | [03-backend/outbox/](03-backend/outbox/) | Spring Boot + Go |

## 추가 예정

- Saga + 보상 트랜잭션 (3장 4.2)
- Idempotency 키 처리 (3장 4.3)
- PG 라우터 (1차 실패 → 2차 fallback) (3장 8.4 / 6장 4.1)
- BFF GraphQL (Spring + Apollo) (4장 4.3)
- ISR + on-demand revalidate (Next.js) (4장 2.2)
- dbt + Iceberg 변환 모델 (5장 2.3)
- LLM 호출 + 캐시 + AI Gateway (5장 6.3)

## 실행 환경 가정

- JDK 21+, Kotlin 2.0+, Gradle (Kotlin DSL)
- Go 1.22+
- PostgreSQL 또는 MySQL (예제는 H2 / PostgreSQL)
- Docker (선택 — DB / Kafka 컨테이너)

## 한국 환경 메모

- 운영에서는 H2 → MySQL (Aurora) / PostgreSQL (Aurora) 사용
- 메시지 브로커 — Kafka (MSK) / NCP Cloud Data Streaming Service / SQS
- 시크릿 — 환경변수가 아닌 AWS Secrets Manager / Vault 권장 (6장 9.1)
