# 📖 Layer 2.3 Implementation Guide: The Entropy Ledger (Audit)

**Status:** ✅ Implemented
**Location:**
*   **Schema (ClickHouse):** `services/grid-service/sql/clickhouse/01_init_audit.sql`
*   **Logger (Go):** `services/grid-service/pkg/audit/logger.go`
**Technology:** ClickHouse (MergeTree) + Go Channels

---

## 1. The "Why": The Black Box Recorder
In financial planning, knowing *who* changed a number is as important as the number itself.
**The Problem:** Storing 1 billion audit logs in Postgres kills performance.
**The Solution:** We use **ClickHouse**.
*   It is a column-store optimized for append-only logs.
*   It can ingest 1M rows/sec on a single node.
*   It compresses log data by 10x-20x.

## 2. The Implementation (Code Deep Dive)

### A. The Schema (`audit_log`)
We use the `MergeTree` engine partitioned by month (`toYYYYMM(timestamp)`).
*   **`action_type` (Enum8):** Stores actions as 1-byte integers for extreme efficiency.
*   **`TTL`:** Automatically deletes logs older than 7 years (SOX compliance).

### B. The Async Logger (Go)
The `AuditLogger` is designed to have **Zero Latency Impact** on the user.
*   **Buffered Channel:** Calls to `Log()` simply push to a Go channel (`chan *AuditEvent`). It takes nanoseconds.
*   **Worker Pool:** A background goroutine reads the channel and batches writes to ClickHouse every 1 second or 100 events.
*   **Safety:** If the buffer fills up (e.g., ClickHouse is down), we drop events rather than crashing the application (Availability > Consistency for Logs).

## 3. Usage Instructions

### Initialization
Run the SQL script against your ClickHouse instance:
```bash
clickhouse-client --host localhost --queries-file services/grid-service/sql/clickhouse/01_init_audit.sql
```

### Logging an Event (Go)
```go
logger := audit.NewAsyncLogger()
// ... inside a handler ...
logger.Log(ctx, tenantID, userID, audit.TypeWriteCell, `{"cell": "A1", "val": 100}`)
```

## 4. Key Design Decisions related to "The Moat"
1.  **Materialized Views:** We create a `anomaly_detection_mv` that pre-aggregates user activity. The **AI Cortex (Layer 8)** reads this view to detect "Mass Deletions" instantly without scanning the raw logs.
2.  **Immutable:** ClickHouse makes it very hard to update/delete rows. This guarantees the integrity of the audit trail.
