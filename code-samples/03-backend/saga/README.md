# Saga + 보상 트랜잭션 — 주문/재고/결제 (3장 4.2)

> Order → Inventory 예약 → Payment 청구 → 실패 시 보상.
> Spring Boot (Kotlin) Orchestrator + Go 일 단위 reconciler 페어.

## 무엇을 보여주는가

- **Orchestration Saga** — 중앙 OrderSagaOrchestrator 가 단계 호출 + 보상 트리거
- **보상 트랜잭션** — Inventory 예약 후 Payment 실패 시 Inventory 해제
- **보상의 보상 (3장 4.2)** — 보상 자체 실패 시 DLQ + Go reconciler 가 야간에 정합 검증

## 디렉토리

```
saga/
├── README.md
├── kotlin-spring/                   # Saga Orchestrator
│   └── src/main/kotlin/com/onepick/saga/
│       ├── SagaApplication.kt
│       ├── OrderSagaOrchestrator.kt
│       ├── steps.kt                 # ReserveInventory / ChargePayment
│       └── compensations.kt
└── go/                              # 일 단위 reconciler (보상의 보상)
    ├── go.mod
    └── reconciler.go
```

## 핵심 패턴

```kotlin
@Transactional(propagation = REQUIRES_NEW)  // 단계마다 별도 트랜잭션
fun reserveInventory(orderId: Long, items: List<Item>) { ... }

// Orchestrator
fun execute(order: Order) {
    val reservedInventory = reserveInventory(order.id, order.items)
    try {
        chargePayment(order.id, order.amount)
        markOrderPaid(order.id)
    } catch (e: PaymentFailedException) {
        compensateInventory(reservedInventory)  // 보상
        markOrderCancelled(order.id, reason = e.message)
    }
}
```

## 책 인용

- 3장 4.2 — Saga 패턴 (Choreography vs Orchestration)
- 3장 4.2 — 보상의 보상 3계층 (DLQ + 수동 + 일 단위 reconciler)
- 3장 8.5 — 실패 시나리오 표

## 라이선스

[MIT](../../LICENSE)
