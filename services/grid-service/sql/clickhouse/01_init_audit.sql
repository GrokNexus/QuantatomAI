-- =========================================================================================
-- QuantatomAI Layer 2.3: Entropy Ledger (Audit Log)
-- Database: ClickHouse
-- Engine: MergeTree (Optimized for Time-Series Appends)
-- =========================================================================================

-- 1. THE MAIN LOG TABLE
CREATE TABLE IF NOT EXISTS audit_log (
    event_id UUID,
    timestamp DateTime64(3), -- Millisecond precision
    tenant_id UUID,
    user_id UUID,
    
    -- Context
    trace_id String, -- OpenTelemetry Trace ID
    ip_address IPv6,
    user_agent String,
    
    -- The Action
    action_type Enum8(
        'LOGIN' = 1,
        'WRITE_CELL' = 2,
        'UPDATE_METADATA' = 3,
        'EXECUTE_CALC' = 4,
        'EXPORT_DATA' = 5
    ),
    resource_target String, -- e.g., "Cell: [Region=NA, Account=Sales]"
    
    -- The Change (Delta)
    old_value String,
    new_value String,
    
    -- Status
    status Enum8('SUCCESS' = 1, 'FAILURE' = 0),
    error_message String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp) -- Partition by Month
ORDER BY (tenant_id, timestamp) -- Optimized for "Show me all logs for this tenant"
TTL timestamp + INTERVAL 7 YEAR; -- SOX Compliance Retention

-- 2. MATERIALIZED VIEW FOR ANOMALY DETECTION (For Cortex Layer 8)
-- Pre-aggregates write volume per user per minute to detect "Mass Deletions"
CREATE MATERIALIZED VIEW IF NOT EXISTS anomaly_detection_mv
ENGINE = SummingMergeTree()
ORDER BY (tenant_id, user_id, toStartOfMinute(timestamp))
AS SELECT
    tenant_id,
    user_id,
    toStartOfMinute(timestamp) as time_bucket,
    count() as operation_count
FROM audit_log
WHERE action_type = 'WRITE_CELL'
GROUP BY tenant_id, user_id, time_bucket;
