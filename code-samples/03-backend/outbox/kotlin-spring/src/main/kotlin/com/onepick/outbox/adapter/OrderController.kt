package com.onepick.outbox.adapter

import com.onepick.outbox.application.OrderService
import org.springframework.web.bind.annotation.*

@RestController
@RequestMapping("/orders")
class OrderController(
    private val orderService: OrderService,
) {

    data class CreateOrderRequest(val amount: Long)
    data class OrderResponse(val id: Long?, val amount: Long, val status: String)

    @PostMapping
    fun create(
        @RequestBody req: CreateOrderRequest,
        @RequestHeader("X-Trace-Id", required = false) traceId: String?,
    ): OrderResponse {
        val order = orderService.createOrder(req.amount, traceId)
        return OrderResponse(order.id, order.amount, order.status)
    }
}
