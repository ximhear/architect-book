package com.onepick.outbox.adapter

import com.onepick.outbox.port.OutboxRepository
import org.slf4j.LoggerFactory
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional
import java.time.Duration
import java.time.Instant

/**
 * 1초마다 미발행 outbox 를 읽어 메시지 브로커에 발행.
 *
 * 학습용으로 stdout 로그로 시뮬레이션. 운영은 KafkaTemplate / SnsClient 등으로 교체.
 *
 * 운영 시 추가 고려:
 * 1. 다중 인스턴스 동시 poll → SELECT FOR UPDATE SKIP LOCKED (Go 예제 참조)
 * 2. Debezium-CDC 기반 outbox 로 본 poller 자체 제거 가능 (3장 4.2)
 * 3. 발행 실패 시 재시도 정책 + DLQ
 * 4. 메트릭 — outbox_lag_seconds, batch_size, error_count (3장 8.6)
 */
@Component
class OutboxPublisher(
    private val outboxRepo: OutboxRepository,
) {
    private val log = LoggerFactory.getLogger(javaClass)

    @Scheduled(fixedDelay = 1000)
    @Transactional
    fun publish() {
        val pending = outboxRepo.findUnpublished()
        if (pending.isEmpty()) return

        val now = Instant.now()
        for (event in pending) {
            // TODO: 실제 발행 — kafkaTemplate.send(event.topic, event.partitionKey, event.payload)
            val lag = Duration.between(event.createdAt, now).toMillis()
            log.info(
                "PUBLISH topic={} partition={} lag_ms={} trace={} payload={}",
                event.topic,
                event.partitionKey,
                lag,
                event.traceId ?: "-",
                event.payload,
            )
            event.publishedAt = now
        }
        log.info("PUBLISHED batch_size={}", pending.size)
    }
}
