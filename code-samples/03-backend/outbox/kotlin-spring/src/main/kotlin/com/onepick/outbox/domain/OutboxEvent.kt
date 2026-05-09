package com.onepick.outbox.domain

import jakarta.persistence.*
import java.time.Instant
import java.util.UUID

/**
 * 책 3장 4.2 Outbox SQL 컬럼 구성과 일치.
 *
 * - aggregate_id / aggregate_type — 같은 aggregate 의 순서 보장 / 도메인 추적
 * - partition_key — Kafka partition 라우팅
 * - trace_id — 분산 추적 (3장 8.6) — 컨슈머까지 전파
 */
@Entity
@Table(
    name = "outbox",
    indexes = [
        Index(name = "idx_outbox_unpublished", columnList = "createdAt"),
    ],
)
class OutboxEvent(
    @Id
    val id: UUID = UUID.randomUUID(),

    @Column(nullable = false, length = 64)
    val aggregateId: String,

    @Column(nullable = false, length = 64)
    val aggregateType: String,

    @Column(nullable = false, length = 128)
    val topic: String,

    @Column(nullable = false, length = 128)
    val partitionKey: String,

    @Lob
    @Column(nullable = false)
    val payload: String,

    @Column(length = 128)
    val traceId: String? = null,

    @Column(nullable = false)
    val createdAt: Instant = Instant.now(),

    @Column
    var publishedAt: Instant? = null,
)
