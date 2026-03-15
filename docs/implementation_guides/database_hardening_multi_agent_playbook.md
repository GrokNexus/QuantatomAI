# QuantatomAI Database Hardening Engineering Playbook

## Document Role
This is the engineering-level execution guide.

Board-facing summary lives in:
- [docs/implementation_guides/database_hardening_board_brief.md](docs/implementation_guides/database_hardening_board_brief.md)

## Purpose
This file is the single reference for evaluating, hardening, and implementing the QuantatomAI database layer in repeated future sessions.

Use this document when you want an agent to:
- assess current database readiness without rereading the entire repo
- understand what is already implemented versus merely designed
- execute one implementation phase at a time with a consistent multi-agent workflow
- carry forward the rationale for each hardening step

## How To Use In A Future Session
Use one of these prompts:

- Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and summarize current database hardening status.
- Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 1.
- Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and run the Phase 3 multi-agent implementation prompt.

## Canonical Multi-Agent Roles
Future agents should simulate and synthesize these perspectives for every phase:

- Chief Product Architect: GTM readiness, enterprise completeness, admin UX, buyer trust
- Chief Database Technology Architect: storage engines, tenancy, partitioning, durability, concurrency, governance
- Founder: moat, long-term extensibility, protocol strategy, innovation pressure
- AI Intelligence Architect: intersection intelligence, metadata intelligence, explainability, anomaly and prediction layers

## Current Baseline Assessment
This is the current assessment of the repo as of 2026-03-14.

### Overall State
- Architecture quality: strong
- Enterprise database hardening: incomplete
- GTM readiness for enterprise finance database layer: not yet ready
- Best characterization: advanced architecture plus partial implementation plus inconsistent readiness evidence

### Why This Assessment Exists
The repository contains both:
- strong target-state architecture documents
- explicit red-team documents admitting critical enterprise gaps

The playbook exists to stop future work from drifting between aspirational design and actual delivery status.

## Primary Evidence Files
These files are the minimum evidence set for database-layer work.

### Core architecture
- [docs/architecture/quantatomai-architecture.md](docs/architecture/quantatomai-architecture.md)
- [docs/architecture/quantatomai_data_structure_innovation.md](docs/architecture/quantatomai_data_structure_innovation.md)
- [docs/architecture/quantatomai-master-schema.md](docs/architecture/quantatomai-master-schema.md)
- [docs/architecture/quantatomai-grid-engine.md](docs/architecture/quantatomai-grid-engine.md)
- [docs/architecture/quantatomai_calculation_engine.md](docs/architecture/quantatomai_calculation_engine.md)
- [docs/architecture/quantatomai_ai_architecture.md](docs/architecture/quantatomai_ai_architecture.md)

### Gap and execution truth sources
- [docs/architecture/quantatomai_gap_analysis.md](docs/architecture/quantatomai_gap_analysis.md)
- [docs/project/quantatomai_implementation_dashboard.md](docs/project/quantatomai_implementation_dashboard.md)
- [docs/architecture/grid-operations-spec-sheet.md](docs/architecture/grid-operations-spec-sheet.md)
- [docs/architecture/GeminiQuantAnalysis.md](docs/architecture/GeminiQuantAnalysis.md)

### Existing implementation guides that matter
- [docs/implementation_guides/layer_2_1_metadata_schema_impl.md](docs/implementation_guides/layer_2_1_metadata_schema_impl.md)
- [docs/implementation_guides/layer_2_2_mdf_impl.md](docs/implementation_guides/layer_2_2_mdf_impl.md)
- [docs/implementation_guides/layer_2_3_audit_impl.md](docs/implementation_guides/layer_2_3_audit_impl.md)
- [docs/implementation_guides/layer_3_1_event_backbone_impl.md](docs/implementation_guides/layer_3_1_event_backbone_impl.md)
- [docs/implementation_guides/layer_3_2_ipc_impl.md](docs/implementation_guides/layer_3_2_ipc_impl.md)
- [docs/implementation_guides/layer_5_compute_impl.md](docs/implementation_guides/layer_5_compute_impl.md)
- [docs/implementation_guides/layer_8_2_git_flow_metadata_impl.md](docs/implementation_guides/layer_8_2_git_flow_metadata_impl.md)
- [docs/implementation_guides/layer_8_4_red_team_alm_hardening_impl.md](docs/implementation_guides/layer_8_4_red_team_alm_hardening_impl.md)

### Empty or weak specs that must be treated as open work
- [docs/architecture/data-tiering-hot-warm-cold.md](docs/architecture/data-tiering-hot-warm-cold.md)
- [docs/architecture/multi-cloud-federation.md](docs/architecture/multi-cloud-federation.md)
- [docs/architecture/offline-ich.md](docs/architecture/offline-ich.md)
- [docs/architecture/privacy-pev.md](docs/architecture/privacy-pev.md)

## Readiness Matrix
| Area | Current state | Notes |
| --- | --- | --- |
| Atom or molecule protocol | Strong design, partial implementation evidence | Good moat candidate |
| Metadata store | Partial implementation evidence | Needs stronger versioning and semantic governance |
| Audit and lineage | Partial implementation evidence | Gap analysis still treats it as insufficient for enterprise trust |
| Multi-tenant control plane | Weak | Isolation model not fully specified |
| Hot or warm or cold lifecycle | Weak | Critical doc is empty |
| Privacy and residency | Weak | Critical doc is empty |
| Offline conflict model | Weak | Critical doc is empty |
| Consolidation and reporting pack | Partial concept | Needs auditable domain package |
| AI-native database intelligence | Conceptually strong | Needs governed operationalization |
| GTM readiness | Not ready | Governance, audit, ALM, connectors, tenancy still need hardening |

## Non-Negotiable Design Principles
- Never trade away immutable auditability for low-latency convenience.
- Never claim enterprise multi-tenancy without explicit tenant keys, isolation tests, and AI boundary rules.
- Never treat AI explanations as valid unless they are reproducible from persisted lineage and model provenance.
- Never let implementation dashboards outrun red-team evidence.
- Preserve the molecule or atom protocol moat while hardening operational trust.

## Cross-Phase Status Model
Use this status language when updating the file:
- Designed: described in architecture, little implementation evidence
- Partial: some code or implementation guides exist, but enterprise proof is incomplete
- Hardened: implemented, benchmarked, and verified against success criteria

## Phase Overview
| Phase | Name | Current status |
| --- | --- | --- |
| 1 | Database Truth Baseline | Partial |
| 2 | Multi-Tenant Control Plane | In Progress |
| 3 | Audit, Lineage, and Workflow Governance | In Progress |
| 4 | Storage and Performance Hardening | In Progress |
| 5 | Metadata Intelligence and Visualization Plane | In Progress |
| 6 | Consolidation and External Reporting Domain Pack | In Progress |
| 7 | AI-Native Operationalization | In Progress |

## Execution Log
### 2026-03-14
- Split documentation into board and engineering views.
- Created board brief at [docs/implementation_guides/database_hardening_board_brief.md](docs/implementation_guides/database_hardening_board_brief.md).
- Started Phase 1 execution by filling previously empty architecture specs for storage tiering, multi-cloud federation, offline ICH, and privacy PEV.
- Started Phase 2 execution with a production-grade tenant control-plane migration, validation checks, an executable migration integration test, a dedicated implementation guide, and a database-layer mind map.
- Started Phase 3 execution with a production-grade governance migration for immutable metadata audit events, workflow state controls, metadata promotion governance, and connector staging governance plus Phase 3 validation checks and implementation guide.
- Started Phase 4 execution with a concrete storage and performance hardening implementation guide, benchmark profile taxonomy, governance-on measurement contract, and initial acceptance thresholds.
- Added Phase 4 runnable artifacts under `tools/load-testing` plus a benchmark evidence template so performance claims can be tied to repeatable evidence bundles.
- Wired the Phase 4 runner to concrete `grid-service` benchmark commands for the grid query serialization path and CRDT merge path.
- Expanded `grid-service` fixture preparation from migrate-only and compat seeding into tenant-aware governance fixtures plus replay-and-recovery seed scripts aligned to the real Phase 2 and Phase 3 schema.
- Phases 5 through 7 remain planned but not yet executed in code or detailed subsystem docs.

### 2026-03-15 (commits 85c2f3b → 3b63244 → 380ea79)
- Hardened Phase 4 fixture orchestration: fixed migrate-only root-relative path (`go -C services/grid-service`), added DATABASE_URL env token to all migrate steps in `grid-service-phase4-fixtures.json`.
- Added docker psql fallback with auto `docker cp` for Windows environments where `psql` is not natively installed.
- Created `services/grid-service/sql/validation/phase4_fixture_smoke_checks.sql` with alignment invariant selects and minimum footprint enforcement for profiles C and D.
- Refactored CI workflow: separated `integration-docker` job from new `phase4-fixture-smoke` matrix job; matrix covers profiles C and D (B added 2026-03-15 session end).
- Extended `run-phase4-profile.ps1` with `Get-AutoExtractedMetrics` (ns/op parser, tenant fairness ratio, replay_invalid_rows) and `Get-ThresholdEvaluation` producing pass/fail/not_evaluated per threshold field in the evidence bundle.
- Added `-DatabaseBacked` switch to `run-grid-service-phase4.ps1` wiring profiles C and D to docker exec psql validation commands; added `Get-LastExitCodeOrZero` for strict-mode safety.
- Fixed `PSNativeCommandUseErrorActionPreference` false-failure from docker NOTICE stderr; wrapped command capture with pre/post toggle.
- Added CI steps to generate Phase 4 evidence bundle (`run-grid-service-phase4.ps1 -DatabaseBacked`) and upload `run-manifest.json` + `evidence-summary.md` as GitHub Actions artifacts per matrix profile and SHA (`phase4-evidence-B-<sha>`, `phase4-evidence-C-<sha>`, `phase4-evidence-D-<sha>`).
- Added `tools/load-testing/results/.gitignore` to prevent timestamped run artifacts from entering git.
- Created `services/grid-service/sql/validation/phase4_planning_workload_smoke_checks.sql` for Profile B (validates Phase 2 tenant + Phase 3 workflow governance footprint and alignment invariants).
- Added Profile B to `databaseCommandMap` in `run-grid-service-phase4.ps1` and expanded CI matrix to `["B", "C", "D"]`.
- Local validation confirmed: profiles C and D smoke checks PASSED (tenants=3, apps=3, promotions=9, ingest_batches=9, rejections=9, replay_audits=3).
- GTM evaluation performed 2026-03-15 — see [docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md](docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md) for full gap analysis, capability scoring, and phase plans 5–8.
- Began Phase 5 implementation with tenant-safe metadata graph API scaffolding in `grid-service`: added `RequireTenantHeader` middleware, `/api/v1/metadata/graph` endpoint, a deterministic graph store stub, and unit tests for tenant-header and query validation behavior.
- Began Phase 6 implementation with `09_consolidation_domain.sql` introducing tenant-aligned close-calendar, intercompany ownership, journals, FX translation policies, elimination rules, and disclosure mappings plus `phase6_consolidation_domain_checks.sql` validation checks.
- Began Phase 7 implementation with `10_ai_inference_governance.sql` introducing `ai_inference_log` provenance and override governance plus `phase7_ai_inference_governance_checks.sql`.
- Extended `cortex-service` with tenant-aware `/api/v1/vector/similarity`, optional live LLM synthesis path in `rag.py` (via `litellm` + `CORTEX_LLM_API_KEY`), and model/tenant metadata in narrative responses.
- Began Phase 8 MVP implementation with a new `compute/heliocalc` Rust crate: AtomScript parser + interpreter for `SUM(...)` with `WHERE` predicates and 10 passing unit tests (`cargo test`).

---

## Phase 1: Database Truth Baseline
### Objective
Convert the current database-layer story into one validated baseline with no empty critical specs and no contradictory status claims.

### Why Implement This
Without this phase, every later phase will be built on mixed assumptions. This is the phase that separates implemented reality from architecture rhetoric.

### Current State
- Strong architecture docs exist.
- Critical database-adjacent specs are empty.
- Implementation dashboard claims full completion, while red-team gap analysis still identifies enterprise blockers.

### Evidence To Review
- [docs/architecture/quantatomai-architecture.md](docs/architecture/quantatomai-architecture.md)
- [docs/architecture/quantatomai_gap_analysis.md](docs/architecture/quantatomai_gap_analysis.md)
- [docs/project/quantatomai_implementation_dashboard.md](docs/project/quantatomai_implementation_dashboard.md)
- [docs/architecture/data-tiering-hot-warm-cold.md](docs/architecture/data-tiering-hot-warm-cold.md)
- [docs/architecture/multi-cloud-federation.md](docs/architecture/multi-cloud-federation.md)
- [docs/architecture/offline-ich.md](docs/architecture/offline-ich.md)
- [docs/architecture/privacy-pev.md](docs/architecture/privacy-pev.md)

### Subsystem Hardening Checklist
- Architecture truth reconciliation
  - Define authoritative current-state versus target-state markers
  - Remove or annotate contradictory progress claims
  - Record open database risks explicitly
- Storage lifecycle
  - Fill hot or warm or cold storage tiering rules
  - Define retention, compaction, archival, recovery boundaries
- Privacy and federation
  - Fill privacy enforcement model
  - Fill data residency and cross-cloud movement rules
- Offline sync
  - Define conflict semantics, merge policies, and replay guarantees

### Success Criteria
- All critical database-adjacent architecture files are non-empty and concrete.
- There is one explicit current-state assessment of the database layer.
- Contradictions between dashboard and red-team docs are reconciled or annotated.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 1. Act as Chief Product Architect, Chief Database Technology Architect, Founder, and AI Intelligence Architect. Reconcile current-state versus target-state database architecture, replace empty critical specs with concrete technical design, and update status language so the repository has one truthful baseline. Deliver updated docs, a resolved gap summary, and a short evidence table showing what is implemented, what is partial, and what is still designed only.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [docs/architecture/data-tiering-hot-warm-cold.md](docs/architecture/data-tiering-hot-warm-cold.md), [docs/architecture/multi-cloud-federation.md](docs/architecture/multi-cloud-federation.md), [docs/architecture/offline-ich.md](docs/architecture/offline-ich.md), [docs/architecture/privacy-pev.md](docs/architecture/privacy-pev.md), [docs/implementation_guides/database_hardening_board_brief.md](docs/implementation_guides/database_hardening_board_brief.md)
- Why implemented: Establish a truthful and non-empty baseline for the database control plane before further hardening work.
- Remaining risks: Dashboard and gap-analysis reconciliation still needs explicit update, and Phase 2 tenant hardening is still only designed.
- Verification evidence: Critical Phase 1 architecture files are no longer empty and the engineering/board documentation split is in place.

---

## Phase 2: Multi-Tenant Control Plane
### Objective
Define and implement tenant isolation, tenant-aware partitioning, cost boundaries, and tenant-safe AI boundaries.

### Why Implement This
The database layer cannot be called enterprise-grade until tenant isolation is explicit, testable, and auditable.

### Current State
- Architecture mentions RLS, security masks, and shard-by-planning-type.
- Full tenant keying, isolation proofs, quota models, and AI learning boundaries are not concretely specified.

### Evidence To Review
- [docs/architecture/quantatomai-master-schema.md](docs/architecture/quantatomai-master-schema.md)
- [docs/architecture/quantatomai-architecture.md](docs/architecture/quantatomai-architecture.md)
- [docs/architecture/privacy-pev.md](docs/architecture/privacy-pev.md)
- [docs/architecture/multi-cloud-federation.md](docs/architecture/multi-cloud-federation.md)

### Subsystem Hardening Checklist
- Tenant identity model
  - Tenant id propagation across metadata, atoms, events, audit, vector store
  - Tenant-aware compound keys and indexes
- Isolation model
  - Row isolation
  - Storage partition isolation
  - Cache isolation
  - Event topic and consumer isolation
- Security and privacy
  - Per-tenant key hierarchy
  - Encryption domain design
  - Tenant-aware policy evaluation
- Performance and cost
  - Noisy-neighbor controls
  - Quotas and budgets
  - Chargeback or showback model
- AI boundaries
  - No cross-tenant retrieval leakage
  - Explicit federated learning policy if any
  - Tenant-scoped embedding and inference rules

### Success Criteria
- Tenant control-plane spec exists.
- Isolation threat model exists.
- Tenant-aware partitioning and keying are implemented or designed in executable detail.
- Tenant-safe AI boundary rules are explicit.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 2. Use the four canonical roles to design and, where possible, implement the QuantatomAI multi-tenant control plane. Produce tenant key strategy, partitioning scheme, isolation policy matrix, quota model, encryption boundary, and tenant-safe AI data-access rules. Update code or docs as needed and include explicit tests or validation criteria for tenant isolation and noisy-neighbor behavior.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [services/grid-service/sql/schema/07_tenant_control_plane.sql](services/grid-service/sql/schema/07_tenant_control_plane.sql), [services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql](services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql), [services/grid-service/sql/schema/migrate_integration_test.go](services/grid-service/sql/schema/migrate_integration_test.go), [docs/implementation_guides/layer_2_4_multi_tenant_control_plane_impl.md](docs/implementation_guides/layer_2_4_multi_tenant_control_plane_impl.md), [docs/architecture/database-layer-mind-map.md](docs/architecture/database-layer-mind-map.md)
- Why implemented: Make tenant isolation, quota control, residency, key domains, app partitioning, and AI boundaries explicit in the database control plane.
- Remaining risks: Service-layer enforcement, runtime throttling, and AI-stack namespace enforcement still need follow-through in later phases.
- Verification evidence: Migration and validation artifacts exist, tenant context is propagated directly into core metadata tables, an executable integration test now verifies control-plane behavior against a real Postgres instance, and the database-layer mind map is documented.

---

## Phase 3: Audit, Lineage, and Workflow Governance
### Objective
Make the database trustworthy for finance: immutable history, explainable lineage, metadata promotion, approvals, locking, and controlled write states.

### Why Implement This
This phase closes the biggest buyer-confidence gap for CFO, controllership, and enterprise platform teams.

### Current State
- Audit and governance are partially described and partially implemented.
- Gap analysis still treats audit, ALM, integration staging, and workflow governance as critical gaps.

### Evidence To Review
- [docs/architecture/quantatomai_gap_analysis.md](docs/architecture/quantatomai_gap_analysis.md)
- [docs/implementation_guides/layer_2_3_audit_impl.md](docs/implementation_guides/layer_2_3_audit_impl.md)
- [docs/implementation_guides/layer_8_2_git_flow_metadata_impl.md](docs/implementation_guides/layer_8_2_git_flow_metadata_impl.md)
- [docs/implementation_guides/layer_8_4_red_team_alm_hardening_impl.md](docs/implementation_guides/layer_8_4_red_team_alm_hardening_impl.md)

### Subsystem Hardening Checklist
- Immutable audit
  - Append-only cell history
  - Old and new values
  - actor, session, machine, API, formula, source lineage
  - efficient cell history query path
- Lineage
  - Atom-to-grid-to-report drill path
  - formula dependency provenance
  - connector source provenance
- Workflow governance
  - node state machine
  - owners and approvers
  - lock propagation and override rules
- Metadata ALM
  - snapshots
  - branching
  - diff and merge
  - promotion gates
- Ingestion governance
  - staging airlock
  - mapping validation
  - reject and quarantine path

### Success Criteria
- Immutable audit trail design and retrieval path are concrete.
- Workflow state machine and ownership model exist.
- Metadata promotion flow is documented or implemented.
- Connector staging and validation rules exist.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 3. As the four canonical roles, design and implement immutable audit, lineage, metadata ALM, workflow governance, and connector staging. Prioritize finance trust, low-latency write safety, and clear admin experience. Deliver updated docs or code plus a verification checklist for cell history, promotion safety, locking, approvals, and staging validation.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [services/grid-service/sql/schema/08_audit_workflow_governance.sql](services/grid-service/sql/schema/08_audit_workflow_governance.sql), [services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql](services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql), [docs/implementation_guides/layer_3_3_audit_workflow_governance_impl.md](docs/implementation_guides/layer_3_3_audit_workflow_governance_impl.md)
- Why implemented: Establish immutable metadata audit trails, workflow transition governance, metadata promotion controls, and connector staging governance as the first concrete Phase 3 trust-control foundation.
- Remaining risks: Request-context attribution into audit events, domain-specific workflow policies, and writeback-path enforcement are still pending service-layer integration.
- Verification evidence: Migration and validation artifacts now exist with SQL-level governance constraints and trigger-based policy enforcement.

---

## Phase 4: Storage and Performance Hardening
### Objective
Prove the database layer under realistic enterprise concurrency, durability, audit, and governance load.

### Why Implement This
Performance claims only matter if they survive durability, audit overhead, tenant mix, and governance logic.

### Current State
- Performance goals are well documented.
- Storage and benchmark proof under enterprise control load is still incomplete.

### Evidence To Review
- [docs/architecture/quantatomai-architecture.md](docs/architecture/quantatomai-architecture.md)
- [docs/architecture/grid-operations-spec-sheet.md](docs/architecture/grid-operations-spec-sheet.md)
- [docs/architecture/GeminiQuantAnalysis.md](docs/architecture/GeminiQuantAnalysis.md)
- [docs/implementation_guides/layer_3_1_event_backbone_impl.md](docs/implementation_guides/layer_3_1_event_backbone_impl.md)
- [docs/implementation_guides/layer_3_2_ipc_impl.md](docs/implementation_guides/layer_3_2_ipc_impl.md)
- [docs/implementation_guides/layer_5_compute_impl.md](docs/implementation_guides/layer_5_compute_impl.md)

### Subsystem Hardening Checklist
- Tiered storage behavior
  - promotion and demotion rules
  - retention windows
  - compaction strategy
  - disaster recovery path
- Concurrency
  - mixed read and write bursts
  - tenant fairness
  - workflow-lock interactions
  - audit overhead interactions
- Eventing and recovery
  - replay semantics
  - exactly-once or effective-once guarantees
  - backpressure behavior
- Benchmarks
  - include governance overhead
  - include audit writes
  - include connector ingest
  - include tenant mix

### Success Criteria
- Benchmark suite is reproducible.
- Recovery and replay behavior are documented.
- Tier transition rules are explicit.
- Performance claims are tied to realistic load profiles.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 4. Use the four canonical roles to harden storage lifecycle, concurrency, event recovery, and benchmark evidence for the database layer. Produce or update storage-tier rules, mixed-workload test plans, recovery semantics, and success thresholds that include governance and audit overhead rather than idealized raw compute only.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [docs/implementation_guides/layer_4_1_storage_performance_hardening_impl.md](docs/implementation_guides/layer_4_1_storage_performance_hardening_impl.md), [tools/load-testing/run-phase4-profile.ps1](tools/load-testing/run-phase4-profile.ps1), [tools/load-testing/run-grid-service-phase4.ps1](tools/load-testing/run-grid-service-phase4.ps1), [tools/load-testing/prepare-grid-service-phase4-fixtures.ps1](tools/load-testing/prepare-grid-service-phase4-fixtures.ps1), [tools/load-testing/grid-service-phase4-fixtures.json](tools/load-testing/grid-service-phase4-fixtures.json), [tools/load-testing/phase4-profiles.json](tools/load-testing/phase4-profiles.json), [docs/implementation_guides/layer_4_2_benchmark_evidence_template.md](docs/implementation_guides/layer_4_2_benchmark_evidence_template.md), [services/grid-service/pkg/orchestration/grid_query_service_benchmark_test.go](services/grid-service/pkg/orchestration/grid_query_service_benchmark_test.go), [services/grid-service/pkg/sync/crdt_benchmark_test.go](services/grid-service/pkg/sync/crdt_benchmark_test.go)
- Why implemented: Establish a concrete benchmark and recovery hardening framework that includes Phase 2 and Phase 3 governance overhead in performance evidence.
- Remaining risks: Docker-dependent local execution and database-backed benchmark coverage still need stabilization even after the richer tenant/replay fixture layer was added.
- Verification evidence: Phase 4 now has explicit benchmark profiles, a runnable evidence-bundle generator, a service-specific benchmark wrapper, a PowerShell-safe fixture runner, tenant-aware governance seed scripts, replay-and-recovery seed scripts, benchmark test entry points in `grid-service`, a benchmark evidence template, measurement contract, recovery contract, and acceptance thresholds documented for repeatable execution.

---

## Phase 5: Metadata Intelligence and Visualization Plane
### Objective
Turn metadata into a governed, versioned, explorable intelligence layer that supports both modeling and AI.

### Why Implement This
The platform claims metadata-driven planning and AI guidance, but metadata must become first-class and inspectable to support that claim.

### Current State
- Metadata is recognized as central.
- Semantic identity, drift detection, graph visualization, and governance-assisted suggestions are still underdeveloped.

### Evidence To Review
- [docs/architecture/quantatomai-master-schema.md](docs/architecture/quantatomai-master-schema.md)
- [docs/architecture/quantatomai_calculation_engine.md](docs/architecture/quantatomai_calculation_engine.md)
- [docs/implementation_guides/layer_2_1_metadata_schema_impl.md](docs/implementation_guides/layer_2_1_metadata_schema_impl.md)
- [docs/implementation_guides/layer_7_4_hierarchy_impl.md](docs/implementation_guides/layer_7_4_hierarchy_impl.md)

### Subsystem Hardening Checklist
- Semantic identity
  - stable business concept ids
  - alias handling
  - schema evolution policy
- Visualization plane
  - metadata graph query API
  - lineage graph endpoints
  - impact analysis view
- Metadata intelligence
  - hierarchy suggestions
  - semantic drift alerts
  - duplicate concept detection
- Governance
  - suggestion approval workflow
  - version tags and promotion provenance

### Success Criteria
- Metadata graph and semantic identity model are explicit.
- Drift detection and impact analysis paths exist.
- Metadata suggestions are governable rather than freeform AI actions.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 5. As the four canonical roles, build the metadata intelligence and visualization plane for QuantatomAI. Produce semantic identity rules, metadata graph APIs, drift detection design, hierarchy suggestion workflows, and governance controls for AI-assisted modeling. Tie every recommendation back to versioned metadata and explainable impact analysis.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [services/grid-service/pkg/orchestration/tenant_middleware.go](services/grid-service/pkg/orchestration/tenant_middleware.go), [services/grid-service/pkg/orchestration/metadata_graph_handler.go](services/grid-service/pkg/orchestration/metadata_graph_handler.go), [services/grid-service/pkg/orchestration/tenant_middleware_test.go](services/grid-service/pkg/orchestration/tenant_middleware_test.go), [services/grid-service/pkg/orchestration/metadata_graph_handler_test.go](services/grid-service/pkg/orchestration/metadata_graph_handler_test.go), [services/grid-service/cmd/grid-service/main.go](services/grid-service/cmd/grid-service/main.go)
- Why implemented: Establish a tenant-safe metadata graph API contract that can be queried by the UI and reused by future lineage and impact-analysis features.
- Remaining risks: Current graph store is a deterministic stub and still needs Postgres-backed retrieval from `dimension_members` and branch overlays.
- Verification evidence: `go test ./pkg/orchestration` runs with no diagnostics errors and the new handler tests validate tenant header enforcement and query parameter validation.

---

## Phase 6: Consolidation and External Reporting Domain Pack
### Objective
Package the database layer into auditable enterprise finance workflows for consolidation, external reporting, and ESG reporting.

### Why Implement This
This is where the database becomes a finance platform rather than a generic planning engine.

### Current State
- FX, eliminations, and reporting concepts exist in architecture.
- There is not yet enough evidence of a full auditable domain package.

### Evidence To Review
- [docs/architecture/grid-operations-spec-sheet.md](docs/architecture/grid-operations-spec-sheet.md)
- [docs/architecture/GeminiQuantAnalysis.md](docs/architecture/GeminiQuantAnalysis.md)
- [docs/architecture/quantatomai-grid-engine.md](docs/architecture/quantatomai-grid-engine.md)

### Subsystem Hardening Checklist
- Consolidation domain model
  - entity close calendar
  - ownership and intercompany metadata
  - journals and adjustments
  - FX translation policy layers
  - eliminations and reclassification rules
- External reporting
  - disclosure mapping
  - report package lineage
  - disclosure-to-atom drillback
- ESG and statutory extensions
  - emission and sustainability data mappings
  - separate assurance path
  - policy-aware calculation provenance

### Success Criteria
- One full close-flow reference design exists.
- One external reporting drillback path exists.
- FX and eliminations are modeled as governed artifacts, not just formulas.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 6. Use the four canonical roles to build a consolidation and external-reporting domain pack on top of the atom lattice. Define entity close controls, journals, FX translation, eliminations, disclosure mapping, and end-to-end audit drillback. Deliver domain models, workflow touchpoints, and acceptance criteria for finance-grade reporting trust.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [services/grid-service/sql/schema/09_consolidation_domain.sql](services/grid-service/sql/schema/09_consolidation_domain.sql), [services/grid-service/sql/validation/phase6_consolidation_domain_checks.sql](services/grid-service/sql/validation/phase6_consolidation_domain_checks.sql)
- Why implemented: Introduce a concrete finance domain schema foundation for close management, journals, FX governance, eliminations, and disclosure mapping with strict tenant alignment constraints.
- Remaining risks: Service endpoints, close-cycle orchestration, drillback APIs, and seeded test data for end-to-end close flow are still pending.
- Verification evidence: Schema files compile under migration ordering rules (`*.sql` embedded in migrate runner) and include tenant-alignment triggers plus targeted validation queries.

---

## Phase 7: AI-Native Operationalization
### Objective
Operationalize AI at intersection, graph, data-quality, and metadata layers with governance, confidence, and tenant safety.

### Why Implement This
This is the actual AI-native moat. Until this is operationalized, the AI layer remains more vision than enterprise product.

### Current State
- AI direction is strong.
- Production governance for AI data access, explainability, drift, and confidence is incomplete.

### Evidence To Review
- [docs/architecture/quantatomai_ai_architecture.md](docs/architecture/quantatomai_ai_architecture.md)
- [docs/architecture/quantatomai_calculation_engine.md](docs/architecture/quantatomai_calculation_engine.md)
- [docs/implementation_guides/layer_8_1_cortex_inference_engine_impl.md](docs/implementation_guides/layer_8_1_cortex_inference_engine_impl.md)

### Subsystem Hardening Checklist
- Intersection intelligence
  - persistent feature vectors
  - cell anomaly detection
  - variance attribution storage
  - scenario similarity retrieval
- Graph intelligence
  - dependency heat scores
  - predictive recalculation priorities
  - graph fragility and impact scoring
- Data quality intelligence
  - hybrid rule plus ML validation
  - missing data inference with confidence
  - human override journaling
- Metadata intelligence
  - semantic drift detection
  - model suggestion approval path
- Explainability and governance
  - confidence scoring
  - model provenance
  - persisted explanation artifacts
  - tenant-safe retrieval and inference

### Success Criteria
- AI explanations are reproducible from stored evidence.
- AI data access is tenant-safe.
- Confidence and override capture are first-class.
- Predictive recalculation is policy-aware and measurable.

### Implementation Prompt
Read [docs/implementation_guides/database_hardening_multi_agent_playbook.md](docs/implementation_guides/database_hardening_multi_agent_playbook.md) and execute Phase 7. As the four canonical roles, operationalize AI-native database intelligence for intersections, calculation graphs, data quality, and metadata. Implement or design persistent feature storage, anomaly and drift workflows, confidence-aware explanations, tenant-safe retrieval boundaries, and policy-aware predictive recalculation. Deliver concrete artifacts and proof requirements, not just conceptual AI features.

### Update This File After Completion
- Status: In Progress
- Implemented artifacts: [services/grid-service/sql/schema/10_ai_inference_governance.sql](services/grid-service/sql/schema/10_ai_inference_governance.sql), [services/grid-service/sql/validation/phase7_ai_inference_governance_checks.sql](services/grid-service/sql/validation/phase7_ai_inference_governance_checks.sql), [services/cortex-service/src/rag.py](services/cortex-service/src/rag.py), [services/cortex-service/src/main.py](services/cortex-service/src/main.py), [services/cortex-service/requirements.txt](services/cortex-service/requirements.txt)
- Why implemented: Move AI from static narrative templates toward governed inference by persisting provenance metadata and adding tenant-scoped vector retrieval plus optional live LLM synthesis.
- Remaining risks: AI inference log persistence is not yet wired from runtime calls, anomaly-detection stream processing is not yet implemented, and similarity endpoint still returns deterministic samples pending pgvector integration.
- Verification evidence: Python diagnostics report no syntax errors in updated Cortex files; Phase 7 SQL constraints enforce confidence bounds, override reason requirements, and tenant alignment.

---

## Cross-Subsystem Hardening Reference
Use this matrix when a phase touches multiple subsystems.

| Subsystem | Key hardening questions |
| --- | --- |
| Metadata | Is semantic identity stable, versioned, and auditable |
| Molecule or atom protocol | Are sparse coordinates portable, durable, and lineage-carrying |
| Hot store | Are latency, eviction, tenant fairness, and replay behavior explicit |
| Warm store | Are partitioning, pruning, and analytical access patterns defined |
| Cold store | Are retention, archival, legal hold, and restoration defined |
| Eventing | Are replay, idempotency, backpressure, and failure semantics explicit |
| Audit ledger | Can every value change be reconstructed and queried quickly |
| Workflow | Are approval, locking, ownership, and override rules explicit |
| Connectors | Is data staged, validated, and quarantined before entering the lattice |
| Privacy | Are encryption, masking, residency, and policy boundaries explicit |
| Multi-tenancy | Are keys, quotas, isolation, and AI boundaries enforceable |
| AI | Are confidence, provenance, tenant safety, and override capture persisted |

## Rules For Updating This Playbook
Whenever a phase is worked on, update this file in the same change set.

Minimum update requirements:
- change phase status if warranted
- add implemented artifacts
- explain why the work was done
- record remaining risks
- add validation evidence

If a future session creates a new implementation guide, add it to the relevant phase section here.

## Recommended Future Session Order
1. Phase 1
2. Phase 2
3. Phase 3
4. Phase 4
5. Phase 5
6. Phase 6
7. Phase 7

## Final Working Rule
If there is ever a conflict between a marketing-style status claim and a red-team gap document, treat the red-team gap as the more truthful input until evidence says otherwise.
