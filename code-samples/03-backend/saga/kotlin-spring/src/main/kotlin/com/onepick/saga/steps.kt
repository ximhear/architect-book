package com.onepick.saga

/**
 * 도메인 모델 — 학습 단순화를 위해 데이터 클래스로.
 */
data class Item(val sku: String, val quantity: Int)

data class Order(
    val id: Long,
    val items: List<Item>,
    val amount: Long,
    val idempotencyKey: String,
)

data class InventoryReservation(val id: String, val orderId: Long)
data class PaymentResult(val paymentId: String)

data class OrderResult(val ok: Boolean, val orderId: Long? = null, val paymentId: String? = null, val error: String? = null) {
    companion object {
        fun success(orderId: Long, paymentId: String) = OrderResult(true, orderId, paymentId)
        fun failed(error: String) = OrderResult(false, error = error)
    }
}

class PaymentFailedException(msg: String) : RuntimeException(msg)

/**
 * 단계 인터페이스 — port. 어댑터(JPA, RestTemplate 등) 는 별도.
 */
interface InventoryStep {
    fun reserve(orderId: Long, items: List<Item>): InventoryReservation
    fun release(reservation: InventoryReservation, reason: String)
}

interface PaymentStep {
    fun charge(orderId: Long, amount: Long, idempotencyKey: String): PaymentResult
    fun refund(paymentId: String, reason: String)
}

interface OrderStore {
    fun markPaid(orderId: Long, paymentId: String)
    fun markCancelled(orderId: Long, reason: String)
}

interface DeadLetterQueue {
    fun enqueue(msg: DlqMessage)
}

data class DlqMessage(val type: String, val payload: Map<String, Any?>)
