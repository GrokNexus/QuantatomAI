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

(omitted)

## Running locally

1) Start dependencies (Postgres + Redis):

```
cd services/grid-service
docker-compose up -d
```

Postgres runs on host port 55432 and includes pgvector via a local image build.

2) Apply migrations (idempotent, embedded):

```
go run ./src -migrate-only
```

Phase 2 tenant-control verification:

```
DATABASE_URL=postgres://quantatomai:quantatomai@localhost:55432/quantatomai?sslmode=disable \
go test -run TestRun_Phase2TenantControlPlane -v ./sql/schema
```

This integration test creates an ephemeral database, applies the embedded migrations through `07_tenant_control_plane.sql`, and verifies tenant propagation triggers, AI policy defaults, write-region uniqueness, and key-domain region enforcement.

3) Run the service:

```
DATABASE_URL=postgres://quantatomai:quantatomai@localhost:55432/quantatomai?sslmode=disable \
REDIS_URL=localhost:6379 \
CORS_ORIGINS=http://localhost:3000 \
go run ./src
```

4) Health check:

```
curl http://localhost:8080/health
```

5) Grid query example:

```
curl -X POST http://localhost:8080/grid/query \
  -H "Content-Type: application/json" \
  -d '{
    "dimensions": {"rows": ["Entity", "Product"], "columns": ["Time"], "pages": ["Scenario"], "filters": {"Region": ["NA"]}},
    "members": {"Entity": ["E100", "E200"], "Product": ["P100", "P200"], "Time": ["2025M01", "2025M02"], "Scenario": ["Working"]}
  }'
```

6) Writeback example:

```
curl -X POST http://localhost:8080/grid/writeback \
  -H "Content-Type: application/json" \
  -d '{"cellEdits":[{"dims":{"Entity":"E100","Product":"P200","Time":"2025M01"},"measure":"Revenue","scenario":"Working","value":12345.67}]}'
```

Environment defaults match `docker-compose.yml` (Postgres user/password/db `quantatomai`, Redis `localhost:6379`). Migrations run automatically on startup; the `-migrate-only` flag seeds the database without launching the HTTP server.

## Seeding synthetic metadata (15 dimensions, ~100k members)

After migrations, load synthetic metadata for UI/testing:

```
make seed
```

This inserts 15 dimensions (`Dim_1`..`Dim_15`) and ~100,000 members into the compat tables. Default model is `default_model`; to override, set a PG setting before running the seed, e.g.:

```
PGPASSWORD=quantatomai psql -h localhost -U quantatomai -d quantatomai \
  -c "SELECT set_config('gridseed.model_id','your_model',false);" \
  -f sql/seed/seed_compat.sql
```

