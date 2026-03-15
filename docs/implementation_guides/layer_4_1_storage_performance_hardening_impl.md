# Layer 4.1 Implementation Guide: Storage and Performance Hardening

## Status
In progress

## Locations
- Primary architecture baseline: [docs/architecture/data-tiering-hot-warm-cold.md](docs/architecture/data-tiering-hot-warm-cold.md)
- Program source of truth: [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md)
- Phase 2 control-plane baseline: [services/grid-service/sql/schema/07_tenant_control_plane.sql](services/grid-service/sql/schema/07_tenant_control_plane.sql)
- Phase 3 governance baseline: [services/grid-service/sql/schema/08_audit_workflow_governance.sql](services/grid-service/sql/schema/08_audit_workflow_governance.sql)
- Load-testing runner: [tools/load-testing/run-phase4-profile.ps1](tools/load-testing/run-phase4-profile.ps1)
- Grid-service wrapper: [tools/load-testing/run-grid-service-phase4.ps1](tools/load-testing/run-grid-service-phase4.ps1)
- Grid-service fixture preparer: [tools/load-testing/prepare-grid-service-phase4-fixtures.ps1](tools/load-testing/prepare-grid-service-phase4-fixtures.ps1)
- Profile catalog: [tools/load-testing/phase4-profiles.json](tools/load-testing/phase4-profiles.json)
- Evidence template: [docs/implementation_guides/layer_4_2_benchmark_evidence_template.md](docs/implementation_guides/layer_4_2_benchmark_evidence_template.md)

## Executive Summary
Phase 4 converts architecture-level performance claims into measurable, repeatable evidence under realistic enterprise controls.

This phase defines the first concrete benchmark and recovery hardening frame for:
- hot, warm, and cold tier transitions
- tenant fairness under mixed workload pressure
- governance overhead from audit and workflow controls
- event replay and recovery acceptance thresholds

## Why This Was Implemented
Phases 2 and 3 introduced production-grade control-plane and governance constraints. Raw throughput numbers are not useful unless they include those constraints.

Phase 4 is the first step to ensure benchmark claims reflect real operating conditions rather than idealized runs.

## Scope For This Increment
### 1. Benchmark profile taxonomy
- profile A: low-latency interactive edits with lock and policy checks
- profile B: mixed read and write planning workload with tenant mix
- profile C: connector ingest plus reconciliation read pressure
- profile D: replay and recovery flow with governance checks enabled

### 2. Measurement contract
All benchmark runs must capture:
- p50, p95, p99 latency
- throughput by operation class
- queue depth and backpressure signal
- audit-write amplification
- tenant fairness spread (max/min throughput and latency by tenant)

### 3. Recovery contract
Recovery validation for this phase includes:
- warm-tier rebuild from event and snapshot inputs
- replay idempotency checks for repeated event windows
- consistency checks for workflow-locked nodes during replay

### 4. Governance-on benchmark requirement
Every benchmark report must include at least one run with:
- Phase 2 tenant controls active
- Phase 3 metadata audit triggers active
- workflow transition and lock behavior active

## Initial Acceptance Thresholds
These thresholds are intentionally conservative and should be tightened after baseline evidence is collected.

- interactive write path p95 under governance: <= 500 ms
- mixed workload read p95 under governance: <= 1500 ms
- replay correctness: zero invalid tenant alignment rows after replay checks
- tenant fairness ratio on critical paths: <= 3.0x p95 spread across active tenants

## Verification Plan
1. Run schema and metadata baseline migrations through Phase 3.
2. Seed representative multi-tenant and workflow-state fixtures.
3. Execute profiles A through D with governance enabled through [tools/load-testing/run-phase4-profile.ps1](tools/load-testing/run-phase4-profile.ps1).
4. Store benchmark and replay output with timestamps and git hash.
5. Record threshold pass/fail results with [docs/implementation_guides/layer_4_2_benchmark_evidence_template.md](docs/implementation_guides/layer_4_2_benchmark_evidence_template.md) and then summarize in the program playbook.

## Open Risks
- Current local environment dependency on Docker can delay reproducible benchmark execution.
- Current benchmark wiring reaches the grid query serialization path and CRDT merge path, but not yet the full service endpoint and database-backed workloads.
- Replay workload generation now has tenant-aware and recovery-aware seed scripts, but automated benchmark assertions over those datasets still need to be wired into live database-backed commands.

## Next Dependencies
- Extend the grid-service profile wrapper to database-backed and connector-ingest benchmark commands once stable local infrastructure is available.
- Add automated pass/fail summary output aligned to the acceptance thresholds above.
- Extend phase output with documented benchmark evidence runs per release candidate.