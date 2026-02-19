# 📖 Layer 2.1 Implementation Guide: The Metadata Schema (Postgres)

**Status:** ✅ Implemented
**Location:** `services/grid-service/sql/schema/01_init_metadata.sql`
**Technology:** PostgreSQL 16 + `ltree` + `pgvector`

---

## 1. The "Why": Solving the Hierarchy Problem
In legacy systems (TM1/Essbase), hierarchies are stored as parent-child adjacency lists (e.g., `ParentID` -> `ChildID`).
**The Problem:** To find all descendants of "Total Company", you need recursive CTEs (Common Table Expressions) that are O(N) slow.

**The Solution:** We use the PostgreSQL `ltree` extension.
*   We store the full path: `Global.NA.USA.NY`.
*   To find all descendants of North America (`NA`), we run: `WHERE path <@ 'Global.NA'`.
*   **Performance:** This is an O(1) B-Tree index lookup. It takes <1ms even with 1M members.

## 2. The Tables (Schema Deep Dive)

### A. `dimensions` & `dimension_members`
This is the core graph.
*   **`path` (ltree):** The materialized path for fast aggregation.
*   **`embedding` (vector):** A 384-dimensional vector slot for the QuantatomAI Cortex. This allows "Semantic Search" (e.g., "Show me all revenue-related accounts" -> finds `Net Sales`, `Gross Income`).

### B. `security_policies` (Holographic ACLs)
We do not use standard SQL GRANTS involved with Row Level Security (RLS) policies directly on data tables because it's too slow for meaningful aggregations.
Instead, we store the "Policy Definition" here as JSONB `rules`.
*   The Application Layer (Go/Rust) reads this policy on boot.
*   It compiles it into a **Bitmask** (`security_mask`).
*   The query engine filters using bitwise operations (`data.mask & user.mask > 0`).

## 3. Usage Instructions

### Initialization
Run the SQL script against your Postgres instance:
```bash
psql -h localhost -U postgres -d quantatomai -f services/grid-service/sql/schema/01_init_metadata.sql
```

### Verification Query (Test the Moat)
To verify the `ltree` performance, run:
```sql
EXPLAIN ANALYZE SELECT * FROM dimension_members WHERE path <@ 'Global.NA';
```
*Target: Execution time < 5ms.*

## 4. Key Design Decisions related to "The Moat"
1.  **UUIDs Everywhere:** We use `uuid_generate_v4()` for all IDs to allow purely distributed, conflict-free writes from multiple regions.
2.  **Separate "App" Registry:** Allows multi-tenancy and multiple applications (e.g., "Finance App", "HR App") within the same cluster.
3.  **JSONB Attributes:** We typically avoid EAV (Entity-Attribute-Value) patterns. We use Postgres Binary JSON (`JSONB`) to store sparse attributes like `currency`, `manager`, or `start_date` without schema migrations.
