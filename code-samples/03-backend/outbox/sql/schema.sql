-- Outbox 패턴 스키마 — 책 3장 4.2 컬럼 구성
-- PostgreSQL 기준. MySQL 은 UUID 컬럼·JSONB 대체 (BINARY(16) / JSON) 필요.

CREATE TABLE IF NOT EXISTS orders (
    id           BIGSERIAL PRIMARY KEY,
    amount       BIGINT NOT NULL,
    status       VARCHAR(32) NOT NULL DEFAULT 'CREATED',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outbox (
    id              UUID PRIMARY KEY,
    aggregate_id    VARCHAR(64) NOT NULL,
    aggregate_type  VARCHAR(64) NOT NULL,
    topic           VARCHAR(128) NOT NULL,
    partition_key   VARCHAR(128) NOT NULL,
    payload         JSONB NOT NULL,
    trace_id        VARCHAR(128),
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    published_at    TIMESTAMP WITH TIME ZONE
);

-- 미발행 outbox 만 빠르게 스캔 (partial index)
CREATE INDEX IF NOT EXISTS idx_outbox_unpublished
    ON outbox (created_at)
    WHERE published_at IS NULL;

-- aggregate 단위 추적용
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate
    ON outbox (aggregate_type, aggregate_id);
