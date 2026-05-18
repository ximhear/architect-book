-- TimescaleDB IoT 스키마 — 2권 5장 시계열 hypertable
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- 스마트미터 raw — 15분 주기 (책 케이스 2.2)
CREATE TABLE IF NOT EXISTS meter_readings (
    household_id  VARCHAR(64) NOT NULL,
    reading_at    TIMESTAMPTZ NOT NULL,
    kwh           DOUBLE PRECISION NOT NULL,
    source        VARCHAR(32) NOT NULL DEFAULT 'gateway',
    PRIMARY KEY (household_id, reading_at)
);

-- hypertable 로 변환 (자동 chunk — 1일 단위)
SELECT create_hypertable(
    'meter_readings',
    'reading_at',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- 가구 + 시각 복합 인덱스 (가구별 최근 N분 조회 최적화)
CREATE INDEX IF NOT EXISTS idx_meter_household_time
    ON meter_readings (household_id, reading_at DESC);

-- 압축 정책 — 7일 지나면 자동 압축
ALTER TABLE meter_readings SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'household_id'
);
SELECT add_compression_policy('meter_readings', INTERVAL '7 days', if_not_exists => TRUE);

-- Continuous Aggregate — 1시간 평균 (책 5장 3절)
CREATE MATERIALIZED VIEW IF NOT EXISTS meter_hourly
WITH (timescaledb.continuous) AS
SELECT
    household_id,
    time_bucket('1 hour', reading_at) AS hour,
    AVG(kwh)  AS avg_kwh,
    SUM(kwh)  AS sum_kwh,
    COUNT(*)  AS sample_count
FROM meter_readings
GROUP BY household_id, hour;

-- 자동 갱신 정책 — 1시간 윈도우 7일 지난 거 다시 안 봄
SELECT add_continuous_aggregate_policy('meter_hourly',
    start_offset      => INTERVAL '7 days',
    end_offset        => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists     => TRUE
);

-- 보존 정책 — raw 90일, 1시간 집계 1년, 1일 집계 5년 (책 5장 3절)
SELECT add_retention_policy('meter_readings', INTERVAL '90 days', if_not_exists => TRUE);
