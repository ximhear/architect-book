# Idempotency 키 처리 (3장 4.3)

> 같은 요청이 여러 번 와도 한 번만 실행되도록.
> Spring Boot (Kotlin) HandlerInterceptor + Go net/http middleware 페어.

## 무엇을 보여주는가

- **idempotency_key 헤더** — 클라이언트가 UUID 등 보냄
- **첫 요청** — 정상 처리 + 응답을 키 → 응답 매핑으로 보존
- **재요청 (같은 키)** — 보존된 응답 그대로 반환 (실제 처리 안 함)
- **TTL 도메인별 차등** (3장 4.3) — 일반 24h, 결제 90일+, 알림 1h

## 디렉토리

```
idempotency/
├── README.md
├── kotlin-spring/
│   └── src/main/kotlin/com/onepick/idem/
│       ├── IdempotencyApplication.kt
│       ├── IdempotencyInterceptor.kt
│       ├── IdempotencyStore.kt
│       └── PaymentController.kt
└── go/
    ├── go.mod
    └── middleware.go
```

## 핵심 — Spring HandlerInterceptor

```kotlin
override fun preHandle(req: HttpServletRequest, res: HttpServletResponse, handler: Any): Boolean {
    val key = req.getHeader("Idempotency-Key") ?: return true
    val cached = store.get(key)
    if (cached != null) {
        res.status = cached.status
        res.contentType = "application/json"
        res.writer.write(cached.body)
        return false  // 컨트롤러 호출 안 함
    }
    return true  // 컨트롤러 호출 + postHandle 에서 응답 캐시
}
```

## 책 인용

- 3장 4.3 — Idempotency 키 + TTL 도메인별 차등 (결제 90일+)
- 6장 4.1 — webhook 재처리 안전

## 라이선스

[MIT](../../LICENSE)
