# QuantatomAI Data Tiering: Hot, Warm, Cold

## Purpose
This document defines the operational storage lifecycle for QuantatomAI atoms or molecules across hot, warm, and cold tiers.

It exists to answer four database-hardening questions:
- what data lives in each tier
- when data moves between tiers
- how recovery and audit interact with tiering
- how tenancy, privacy, and cost are enforced during lifecycle transitions

## Tier Model
### Hot Tier
- Primary role: low-latency interaction path for active planning, writeback, and interactive recalculation
- Candidate engines: Redis Cluster, DragonflyDB, or Scylla for latency-sensitive access
- Target data shape: active cell windows, current working set, lock state, session-local projections, near-term recalculation queues
- Persistence stance: cache-plus-journal, never sole source of truth

### Warm Tier
- Primary role: authoritative analytical substrate for active and recently active planning models
- Candidate engines: ClickHouse, DuckDB for local workflows, Iceberg-backed partition sets, Postgres metadata sidecar
- Target data shape: recent atoms, reusable aggregates, active scenario partitions, lineage-enriched read models
- Persistence stance: primary query tier for enterprise interactive analytics

### Cold Tier
- Primary role: durable historical archive, regulatory retention, replay, and AI training source subject to policy
- Candidate engines: Iceberg on object storage, Parquet or Arrow bundles, append-only audit archives
- Target data shape: historical snapshots, closed periods, archived scenarios, audit exports, replayable event history
- Persistence stance: lowest-cost durable system of record for historical states and recovery artifacts

## Placement Rules
### Hot placement
Place data in hot tier when all conditions are true:
- active user interaction is expected within minutes
- p95 latency target is sub-second or better
- data participates in live recalculation, locking, or conflict resolution

Examples:
- visible grid windows
- current spread and allocation worksets
- active approval-state overlays
- near-real-time connector delta buffers

### Warm placement
Place data in warm tier when any of these are true:
- data is needed for interactive or near-interactive drill, variance, and reporting
- data is a recent but not ultra-hot scenario or version
- data is required for reusable aggregates, explainability, or governed AI feature computation

Examples:
- open plan versions for current cycle
- current-year actuals with full drill support
- lineage-aware report projections
- recently closed monthly snapshots still under heavy analysis

### Cold placement
Place data in cold tier when any of these are true:
- period is closed and outside the active planning window
- legal retention or audit preservation is required
- data is large-volume historical material used mostly for replay, investigation, or model training

Examples:
- prior-year closed periods
- archived scenarios and snapshots
- immutable audit ledger exports
- historical connector landing artifacts

## Promotion And Demotion Policies
### Promotion to hot
Trigger promotion when:
- a user opens or reopens a view with predicted high interaction likelihood
- a workflow transition changes a node from archived or review to active editing
- predictive recalculation or anomaly investigation requests rapid access

### Demotion from hot to warm
Trigger demotion when:
- session activity drops below threshold
- recalculation or approval work completes
- lock state and edit queues have drained

### Demotion from warm to cold
Trigger demotion when:
- period close is finalized
- retention policy tags data as archive-eligible
- scenario lifecycle changes to archived or audit-only

Demotion must preserve:
- immutable lineage
- tenant identity
- encryption context
- replay metadata

## Recovery And Durability Model
- Hot tier failure must not lose authoritative writes because all durable mutations are anchored in event log plus warm or cold persistence.
- Warm tier failure must support rebuild from event log, cold snapshots, and metadata checkpoints.
- Cold tier is the long-retention recovery anchor for replay and legal evidence.
- Audit records must never depend solely on hot-tier retention.

## Retention Model
| Artifact | Hot | Warm | Cold |
| --- | --- | --- | --- |
| Active planning cells | minutes to hours | days to months | not primary |
| Current-year actuals | optional cache | primary | archive copy |
| Closed periods | no | short analytical window | primary |
| Audit records | no | short query acceleration | primary |
| AI feature history | optional cache | primary recent window | governed archive |

## Tenant And Privacy Controls
- Every tier transition must retain tenant id, policy tags, and encryption metadata.
- Cross-tenant co-residency is allowed only where logical isolation, encryption boundaries, and performance fairness are explicit.
- Cold-tier training or retrieval use must respect tenant-safe AI boundary rules.

## Operational SLO Guidance
- Hot tier: optimize for interaction latency and lock-state coordination
- Warm tier: optimize for scan efficiency, partition pruning, and recent historical analysis
- Cold tier: optimize for durability, replay fidelity, retention cost, and evidence exportability

## Open Phase 2 And Phase 4 Dependencies
- Tenant-aware partition keys still need formalization in Phase 2.
- Benchmark and recovery proof under this lifecycle model still need completion in Phase 4.
