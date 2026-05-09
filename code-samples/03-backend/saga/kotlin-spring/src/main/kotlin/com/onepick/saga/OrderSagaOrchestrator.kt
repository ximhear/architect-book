package com.onepick.saga

import org.slf4j.LoggerFactory
import org.springframework.stereotype.Service

/**
 * 주문 결제 Orchestration Saga (3장 4.2).
 *
 * 단계:
 * 1. ReserveInventory — 재고 예약 (별도 트랜잭션)
 * 2. ChargePayment    — PG 결제 (외부 호출)
 * 3. MarkOrderPaid    — 주문 완료 표시
 *
 * 실패 시 역순 보상:
 * - Payment 실패 → Inventory 해제 (보상)
 * - MarkOrderPaid 실패 (네트워크) → 결제 환불 + 재고 해제 (보상)
 *
 * 보상 자체가 실패하면? → 3장 4.2 의 "보상의 보상" — DLQ + 일 단위 reconciler.
 */
@Service
class OrderSagaOrchestrator(
    private val inventory: InventoryStep,
    private val payment: PaymentStep,
    private val orderStore: OrderStore,
    private val dlq: DeadLetterQueue,
) {
    private val log = LoggerFactory.getLogger(javaClass)

    fun execute(order: Order): OrderResult {
        log.info("saga start order={}", order.id)

        // 1. 재고 예약
        val reservation = inventory.reserve(order.id, order.items)

        // 2. 결제 (외부 PG 호출 — 별도 트랜잭션)
        val paymentResult = try {
            payment.charge(order.id, order.amount, order.idempotencyKey)
        } catch (e: PaymentFailedException) {
            log.warn("payment failed order={} reason={}", order.id, e.message)
            compensateInventory(reservation, "payment_failed")
            orderStore.markCancelled(order.id, reason = e.message ?: "payment failed")
            return OrderResult.failed(e.message ?: "payment failed")
        }

        // 3. 주문 완료 표시
        try {
            orderStore.markPaid(order.id, paymentResult.paymentId)
        } catch (e: Exception) {
            // 결제는 됐는데 우리 DB 가 죽음 — 환불 + 재고 해제 + 외부 알림
            log.error("markPaid failed order={} payment={} — compensating", order.id, paymentResult.paymentId, e)
            try {
                payment.refund(paymentResult.paymentId, "internal_error")
            } catch (re: Exception) {
                // 보상의 보상 — DLQ 로
                dlq.enqueue(
                    DlqMessage(
                        type = "PAYMENT_REFUND_FAILED",
                        payload = mapOf(
                            "orderId" to order.id,
                            "paymentId" to paymentResult.paymentId,
                            "reason" to (re.message ?: "unknown"),
                        ),
                    )
                )
            }
            compensateInventory(reservation, "internal_error")
            return OrderResult.failed("internal error after payment — compensated")
        }

        log.info("saga success order={} payment={}", order.id, paymentResult.paymentId)
        return OrderResult.success(order.id, paymentResult.paymentId)
    }

    private fun compensateInventory(reservation: InventoryReservation, reason: String) {
        try {
            inventory.release(reservation, reason)
        } catch (e: Exception) {
            log.error("inventory compensation failed reservation={} reason={}", reservation.id, reason, e)
            // 보상의 보상 — DLQ. 일 단위 reconciler 가 정합 잡음.
            dlq.enqueue(
                DlqMessage(
                    type = "INVENTORY_RELEASE_FAILED",
                    payload = mapOf(
                        "reservationId" to reservation.id,
                        "reason" to reason,
                        "error" to (e.message ?: "unknown"),
                    ),
                )
            )
        }
    }
}
