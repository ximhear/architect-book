package com.onepick.idem

import org.springframework.stereotype.Component
import java.time.Duration
import java.time.Instant
import java.util.concurrent.ConcurrentHashMap

data class CachedResponse(
    val status: Int,
    val body: String,
    val expiresAt: Instant,
)

/**
 * 학습용 인메모리 store. 운영에서는 Redis 또는 DB.
 *
 * TTL 정책 (책 3장 4.3 — 도메인별 차등):
 * - 일반 API: 24h
 * - 결제 / 환불: 90일+ (chargeback 윈도우)
 * - 알림: 1h
 *
 * 키 충돌 / lock — 같은 키 동시 요청 시 첫 요청이 끝나기 전 두 번째가 들어오면 race.
 * 운영 권장: Redis SETNX (set-if-not-exists) 로 처리 중 표시 + 짧은 lock TTL.
 */
@Component
class IdempotencyStore {
    private val store = ConcurrentHashMap<String, CachedResponse>()

    fun get(key: String): CachedResponse? {
        val cached = store[key] ?: return null
        if (cached.expiresAt.isBefore(Instant.now())) {
            store.remove(key)
            return null
        }
        return cached
    }

    fun put(key: String, status: Int, body: String, ttl: Duration) {
        store[key] = CachedResponse(status, body, Instant.now().plus(ttl))
    }
}
