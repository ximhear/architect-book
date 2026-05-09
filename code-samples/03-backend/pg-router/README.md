# PG 라우터 — 1차 → 2차 → 3차 자동 fallback (3장 8.5 / 6장 4.1)

> 결제 가용성 99.99% ([케이스 6](../../../manuscript/case-study/README.md#6-운영-sla)) 를 위한 PG 다중화.
> Go 단일 프로세스 — 한국 PG 3사 (토스 / KCP / NICE) 라우팅 시뮬레이션.

## 무엇을 보여주는가

- **순차 fallback** — 1차 PG 실패 → 2차 시도 → 3차 시도 → 모두 실패 시 결제 실패
- **per-PG circuit breaker** — 한 PG 가 연속 실패하면 일정 시간 그 PG 건너뜀
- **idempotency 키 전파** — 같은 결제 요청이 여러 PG 로 가도 중복 청구 안 되게
- **메트릭** — PG 별 시도 / 성공 / 실패 카운트 (3장 8.6)

## 디렉토리

```
pg-router/
└── go/
    ├── go.mod
    └── main.go
```

## 핵심 코드 흐름

```go
for _, pg := range []PG{toss, kcp, nice} {
    if !pg.Healthy() { continue }       // circuit breaker open
    res, err := pg.Charge(ctx, req)
    if err == nil { return res, nil }   // 성공 — 종료
    pg.RecordFailure(err)               // 실패 카운트 증가
}
return nil, errors.New("all PGs failed")
```

## 책 인용

- 3장 8.5 — PG 라우터 fallback 시퀀스 다이어그램
- 6장 4.1 / 4.2 — 결제 가용성 / PG 3사 다중화
- 케이스 5.2 — 결제 다중화 정책

## 실행

```bash
cd go
go run .
# 다른 터미널:
curl -X POST -H 'Idempotency-Key: order-123' \
  -d '{"orderId":123,"amount":50000}' \
  http://localhost:8082/charge
```

## 운영 / 학습 차이

| 영역 | 본 예제 | 운영 |
|------|--------|------|
| PG 호출 | 50% 확률 시뮬 | 실제 HTTP / SDK |
| Circuit breaker | 연속 3회 실패 → 30초 open | gobreaker / hystrix-go 등 라이브러리 |
| idempotency | 메모리 map | Redis SETNX |
| 메트릭 | stdout | Prometheus + Grafana Cloud |

## 라이선스

[MIT](../../LICENSE)
