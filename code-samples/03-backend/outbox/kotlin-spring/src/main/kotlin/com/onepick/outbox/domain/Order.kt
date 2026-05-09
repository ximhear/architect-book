package com.onepick.outbox.domain

import jakarta.persistence.*
import java.time.Instant

@Entity
@Table(name = "orders")
class Order(
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    val id: Long? = null,

    @Column(nullable = false)
    val amount: Long,

    @Column(nullable = false)
    var status: String = "CREATED",

    @Column(nullable = false)
    val createdAt: Instant = Instant.now(),
)
