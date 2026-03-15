# QuantatomAI Phase 4 Load Testing

This directory contains the first concrete Phase 4 benchmark and replay scaffolding.

## Files
- `phase4-profiles.json`: canonical profile definitions for benchmark profiles A through D
- `run-phase4-profile.ps1`: creates a timestamped evidence bundle for one profile run
- `run-grid-service-phase4.ps1`: maps Phase 4 profiles to concrete `grid-service` benchmark commands
- `grid-service-phase4-fixtures.json`: fixture strategy catalog for `grid-service` Phase 4 profiles
- `prepare-grid-service-phase4-fixtures.ps1`: executes or dry-runs fixture setup for `grid-service` profiles
- `services/grid-service/sql/seed/phase4_tenant_governance_seed.sql`: seeds multi-tenant control-plane, app partitions, security policies, workflow nodes, promotions, and ingest batches
- `services/grid-service/sql/seed/phase4_replay_recovery_seed.sql`: seeds recovery branches, branch overrides, quarantined replay batches, and replay audit checkpoints

## Profile Mapping
- `A`: interactive write path with lock and policy checks
- `B`: mixed read and write workload across multiple tenants
- `C`: connector ingest plus reconciliation read pressure
- `D`: replay and recovery flow with governance checks enabled

## Usage
Dry-run profile A:

```powershell
pwsh -File tools/load-testing/run-phase4-profile.ps1 -Profile A -DryRun
```

Dry-run the grid-service profile mapping for profile B:

```powershell
pwsh -File tools/load-testing/run-grid-service-phase4.ps1 -Profile B -DryRun
```

Dry-run profile C with fixture preparation:

```powershell
pwsh -File tools/load-testing/run-grid-service-phase4.ps1 -Profile C -PrepareFixtures -DryRun
```

Run profile D in database-backed mode (fixture smoke SQL path):

```powershell
pwsh -File tools/load-testing/run-grid-service-phase4.ps1 -Profile D -DatabaseBacked -PrepareFixtures
```

Prepare replay fixtures against a non-default database target:

```powershell
pwsh -File tools/load-testing/prepare-grid-service-phase4-fixtures.ps1 \
  -Profile D \
  -DatabaseHost localhost \
  -DatabaseUser quantatomai \
  -DatabaseName quantatomai \
  -DryRun
```

Run profile B and capture an external command result:

```powershell
pwsh -File tools/load-testing/run-phase4-profile.ps1 \
  -Profile B \
  -Command "go test -run TestBenchmarkMixedTenantLoad ./..."
```

## Output
Each run creates a folder under `tools/load-testing/results/` with:
- `run-manifest.json`: profile metadata, environment hints, and execution outcome
- `evidence-summary.md`: benchmark summary template for manual or automated completion
- `command-output.txt`: captured command output when `-Command` is provided

`run-manifest.json` now also includes:
- `autoExtractedMetrics`: best-effort parsed metrics from command output (for example p95 latency derived from `ns/op` and replay invalid rows)
- `thresholdEvaluation`: per-threshold `pass`/`fail`/`not_evaluated` status based on extracted metrics

## Working Rule
Do not publish benchmark claims from ad hoc commands alone. Record the run through `run-phase4-profile.ps1` so each evidence bundle includes:
- git hash
- profile id
- governance assumptions
- threshold targets
- command outcome

When testing `grid-service`, prefer `run-grid-service-phase4.ps1` so the benchmark profile stays aligned to concrete service packages.
When profiles require seeded metadata or migrated schema, add `-PrepareFixtures` so setup steps are declared and captured consistently.
Profiles B through D now seed tenant-aware governance fixtures first, then layer compat or replay datasets on top so benchmark evidence can exercise Phase 2 and Phase 3 controls instead of migrate-only baselines.
If the target database requires a non-default password, set `PGPASSWORD` in the shell before running the fixture preparer or grid-service wrapper.