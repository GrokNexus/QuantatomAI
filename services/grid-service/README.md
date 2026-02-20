# Grid Service

The Grid Service is responsible for transforming a user's multidimensional grid selection
(dimension axes, members, filters, scenarios, versions, time) into an executable Atom Query Plan.

It provides:

- Grid Query Planner
- Atom retrieval orchestration (Hot/Warm/Cold)
- HelioCalc compute orchestration
- Writeback handling (cell edits, spreads, allocations)
- Query caching
- Integration with AODL/WRM for updates

## Responsibilities

### 1. Grid Query Planner
- Validates incoming GridQuery payloads
- Resolves dimensions, members, hierarchies
- Builds Atom Query Plan:
  - Required dimensions
  - Required atom blocks
  - Required aggregations
  - Hot vs Warm routing
  - Compute graph for HelioCalc

### 2. Atom Retrieval
- Fetches atom blocks from:
  - Hot Store (Redis/Scylla)
  - Warm Store (ClickHouse/DuckDB)
  - Cold Store (Iceberg)
- Applies APWO partition hints

### 3. Compute Orchestration
- Calls HelioCalc for:
  - Aggregations
  - Variances
  - V%
  - FX translation
  - Eliminations
  - Custom formulas

### 4. Writeback
- Validates edits
- Writes to Hot Store
- Emits `CellUpdated` events to AODL
- Triggers incremental recalculation

### 5. Caching
- Query result cache
- Atom block cache
- Member/hierarchy cache

## API

See `docs/api/grid-api.md`.

## Directory Structure

