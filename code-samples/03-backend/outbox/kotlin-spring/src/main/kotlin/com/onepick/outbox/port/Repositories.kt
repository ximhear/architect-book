package com.onepick.outbox.port

import com.onepick.outbox.domain.Order
import com.onepick.outbox.domain.OutboxEvent
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query
import org.springframework.data.repository.query.Param
import java.util.UUID

interface OrderRepository : JpaRepository<Order, Long>

interface OutboxRepository : JpaRepository<OutboxEvent, UUID> {

    /**
     * 미발행 outbox 를 createdAt 순으로 N개 가져옴.
     *
     * 운영의 다중 인스턴스에서는 SELECT FOR UPDATE SKIP LOCKED 권장 — Go 예제 참조.
     * 본 학습 예제는 단일 노드 가정.
     */
    @Query("SELECT o FROM OutboxEvent o WHERE o.publishedAt IS NULL ORDER BY o.createdAt")
    fun findUnpublished(): List<OutboxEvent>
}
