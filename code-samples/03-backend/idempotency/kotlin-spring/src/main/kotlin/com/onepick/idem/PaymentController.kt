package com.onepick.idem

import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
@RequestMapping("/payments")
class PaymentController {

    data class ChargeRequest(val orderId: Long, val amount: Long)
    data class ChargeResponse(val paymentId: String, val orderId: Long, val amount: Long)

    /**
     * curl -X POST -H 'Content-Type: application/json' \
     *   -H 'Idempotency-Key: 11111111-1111-1111-1111-111111111111' \
     *   -d '{"orderId": 1, "amount": 50000}' \
     *   http://localhost:8080/payments
     *
     * 같은 Idempotency-Key 로 두 번 호출 → 첫 응답이 그대로 반환됨.
     */
    @PostMapping
    fun charge(@RequestBody req: ChargeRequest): ChargeResponse {
        // 운영: PG 라우터 호출 (pg-router 예제 참조)
        val paymentId = "pay_${UUID.randomUUID()}"
        return ChargeResponse(paymentId, req.orderId, req.amount)
    }
}
