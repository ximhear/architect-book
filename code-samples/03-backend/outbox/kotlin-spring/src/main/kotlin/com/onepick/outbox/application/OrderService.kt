package com.onepick.outbox.application

import com.fasterxml.jackson.databind.ObjectMapper
import com.onepick.outbox.domain.Order
import com.onepick.outbox.domain.OutboxEvent
import com.onepick.outbox.port.OrderRepository
import com.onepick.outbox.port.OutboxRepository
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class OrderService(
    private val orderRepo: OrderRepository,
    private val outboxRepo: OutboxRepository,
    private val mapper: ObjectMapper,
) {

    /**
     * Order 생성 + Outbox INSERT 가 한 트랜잭션 안에 있음.
     *
     * 핵심 (3장 4.2): 둘 중 하나만 성공하는 경우가 없음 — DB-이벤트 원자성.
     * 트랜잭션이 commit 되어야 outbox 도 보임 → poller 가 그때부터 발행 가능.
     */
    @Transactional
    fun createOrder(amount: Long, traceId: String? = null): Order {
        val order = orderRepo.save(Order(amount = amount))

        val payload = mapper.writeValueAsString(
            mapOf(
                "orderId" to order.id,
                "amount" to order.amount,
                "status" to order.status,
            )
        )

        outboxRepo.save(
            OutboxEvent(
                aggregateId = order.id.toString(),
                aggregateType = "Order",
                topic = "OrderCreated",
                partitionKey = order.id.toString(),
                payload = payload,
                traceId = traceId,
            )
        )

        return order
    }
}
