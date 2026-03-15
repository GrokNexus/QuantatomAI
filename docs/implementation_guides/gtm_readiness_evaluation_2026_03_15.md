# QuantatomAI GTM Readiness Evaluation
**Date:** 2026-03-15  
**Evaluator roles applied:** Chief Product Architect · Chief Database Technology Architect · Founder · AI Intelligence Architect  
**Scope:** DB layer (Phases 1–4 output) + AI capabilities (Cortex / oracle / formula intelligence)  
**Reference commit:** `380ea79` (HEAD, main)

---

## 1. Executive Verdict

| Dimension | Score | Status |
|---|---|---|
| DB control plane (schema, tenancy, governance) | 6 / 10 | Partial – strong foundations, service-layer enforcement missing |
| Storage lifecycle (hot/warm/cold automation) | 3 / 10 | Designed – no automated tier promotion in code |
| Performance evidence under enterprise load | 4 / 10 | Early – benchmark harness exists, DB-backed results pending CI |
| AI inference (Cortex) | 2 / 10 | Stub – FastAPI scaffolding only, no live model connections |
| Formula intelligence (AtomScript) | 1 / 10 | Concept – no parser, compiler, or Monaco integration |
| Audit / lineage retrieval | 3 / 10 | Interface only – ClickHouse stub returns hardcoded records |
| Connector fabric | 2 / 10 | WASM host stub – no UI wizard, no staging airlock UX |
| Multi-tenant isolation evidence | 5 / 10 | Schema + validation SQL exist, service-layer enforcement unproven |
| Privacy / field-level encryption (PEV) | 2 / 10 | Spec filled, no KMS/Sodium wiring in code |
| Consolidation and close domain pack | 1 / 10 | Architectural concept only |
| **Overall GTM readiness (enterprise finance)** | **3 / 10** | **Not ready for enterprise sales** |

**One-line verdict:** The architecture is enterprise-grade in ambition and the control-plane schema is a real foundation. Everything above the schema — AI inference, formula execution, connector UX, audit retrieval, storage tiering, close workflows — is either a stub or a design document. A Fortune 500 CFO cannot buy stubs.

---

## 2. DB Layer Deep Dive

### 2.1 What is genuinely implemented

| Artifact | Evidence | Notes |
|---|---|---|
| Tenant control plane migration | `services/grid-service/sql/schema/07_tenant_control_plane.sql` | Row-level tenancy, per-tenant quota rows, residency columns |
| Tenant validation SQL | `sql/validation/phase2_tenant_control_plane_checks.sql` | 10+ checks for isolation invariants |
| Migration integration test | `sql/schema/migrate_integration_test.go` | Runs against real Postgres, proves Phase 2 schema lands clean |
| Governance migration | `sql/schema/08_audit_workflow_governance.sql` | Immutable audit events, workflow state machine, metadata promotion, connector staging |
| Governance validation SQL | `sql/validation/phase3_audit_workflow_governance_checks.sql` | Post-migration assertion suite |
| Benchmark profiles A–D | `tools/load-testing/phase4-profiles.json` | Four workload profiles with thresholds |
| Benchmark harness | `tools/load-testing/run-phase4-profile.ps1` | Evidence bundle gen with auto-extracted metrics + threshold eval |
| DB-backed execution for B/C/D | `tools/load-testing/run-grid-service-phase4.ps1` | Docker psql commands wired per profile |
| Fixture smoke validation SQL B/C/D | `sql/validation/phase4_*.sql` | Post-fixture invariant checks |
| CI matrix (B, C, D) with artifact upload | `.github/workflows/ci-grid-service.yml` | Each run publishes `run-manifest.json` + `evidence-summary.md` |
| Hot/warm/cold architecture spec | `docs/architecture/data-tiering-hot-warm-cold.md` | Non-empty spec since Phase 1 |
| Delta-branching (metadata git-flow) | `sql/schema/02_git_metadata.sql`, impl guide | Sparse override pattern implemented in schema |
| CRDT merge benchmarks | `pkg/sync/crdt_benchmark_test.go` | Go-level benchmark, not yet DB-backed |
| Audit lineage interface | `pkg/audit/lineage.go` | Go interface + ClickHouse stub |

### 2.2 What is designed but not implemented

| Gap | Risk level | Notes |
|---|---|---|
| Automated hot→warm→cold tier promotion | HIGH | No promotion daemon, no retention policy triggers, no compaction jobs |
| ClickHouse entropy ledger (live) | HIGH | `lineage.go` returns hardcoded simulation records; no real ClickHouse connection wired |
| Field-level encryption (PEV / AWS KMS / Sodium) | HIGH | Spec exists (`privacy-pev.md`), zero code in service layer |
| Service-layer tenant enforcement (RLS, middleware) | HIGH | Schema has tenant columns; no HTTP/gRPC middleware enforcing tenant context on every request |
| Connector staging airlock (UI + backend validation) | HIGH | `connectors-service/` directory exists; `wasm_host.go` stub; no UI wizard or mapping layer |
| Causal vector clocks on atoms | MEDIUM | Designed in `quantatomai_competitive_analysis.md`; not in schema |
| Bit-packed atom coordinate key (u128/u256) | MEDIUM | Designed for 25-dim scale; still using UUID columns |
| Write buffer (LSM memtable for 5k+ concurrent writes) | MEDIUM | Designed as Redis Scylla-style memtable; not implemented |
| Online schema evolution + dimension add without downtime | MEDIUM | Metadata virtualization claimed; not proven in test |
| OLAP-grade read paths via ClickHouse | MEDIUM | ClickHouse schema folder exists (`sql/clickhouse/`); no wired query path |
| Multi-cloud federation (Crossplane SAIC) | LOW-MEDIUM | Crossplane manifests may exist but multi-cloud path not exercised |
| GDPR / data-residency routing in request path | MEDIUM | Columns exist; no policy evaluation in runtime |

### 2.3 DB layer critical path to GTM

To close the minimum enterprise DB gap the following three must land before any serious sales motion:

1. **Service-layer tenant enforcement** — every gRPC/REST handler must propagate `tenant_id` through a middleware and reject cross-tenant operations at runtime, not just at schema level.
2. **Live ClickHouse audit retrieval** — the "right-click > show history" experience must return real records within 500ms. The hardcoded stub disqualifies the audit story entirely.
3. **Hot/warm/cold tier promotion daemon** — at minimum a background job that enforces retention windows and promotes warm atoms to cold storage; without it the cost and performance claims in the architecture are unverifiable.

---

## 3. AI Capabilities Deep Dive

### 3.1 What is genuinely implemented

| Artifact | Evidence | Notes |
|---|---|---|
| Cortex FastAPI service | `services/cortex-service/src/main.py` | Health, `/forecast/auto-baseline` stub, `/narrative/variance` (wired to RAG) |
| Variance narrative RAG | `services/cortex-service/src/rag.py` | Deterministic template narrative; confidence score hardcoded |
| Cortex K8s deployment | `infra/k8s/base/cortex-service/` | Deployment manifest registered in kustomization |
| `MDFVectorReader` | `services/cortex-service/src/data/` | Reads Parquet via pyarrow; fallback simulation logic |
| `pgvector` in schema | Phase 2 schema | `embedding_vector` column type exists on metadata tables |
| Audit lineage interface | `pkg/audit/lineage.go` | Go interface ready for ClickHouse backing |
| Oracle adapters directory | `ai/oracle/adapters/` | Scaffold for LLM adapters |

### 3.2 What is designed but not implemented

| Gap | Risk level | Notes |
|---|---|---|
| Live LLM integration (OpenAI / Bedrock / Llama) | CRITICAL | `rag.py` calls no external model; `litellm` commented out as placeholder |
| Auto-baseline forecasting (TimeGPT / Lag-Llama) | CRITICAL | `/forecast/auto-baseline` returns `"status": "queued"` with no model execution |
| Anomaly detection (Isolation Forest on ingest stream) | HIGH | Designed in `quantatomai_ai_architecture.md`; no Redpanda/Kafka consumer or model |
| AtomScript parser (Rust / pest / nom) | CRITICAL | Entire formula language is a design doc; Monaco editor not integrated |
| LLVM JIT compilation for AtomScript | CRITICAL | Described as moat differentiator; no code |
| Scenario generation via NLP → AtomScript | HIGH | LLM-to-DSL pipeline designed; no code |
| Embedding pipeline (atoms → pgvector) | HIGH | Column exists; no embedding computation or ingestion pipeline |
| Semantic similarity / scenario retrieval | HIGH | Depends on embedding pipeline |
| AUH (AI UX Harmonizer — LlamaIndex suggestions) | HIGH | Listed in Layer 7 architecture; no code |
| Explainable variance "Narrative Card" in UI | HIGH | RAG endpoint exists; no UI plumbing |
| Model provenance / confidence storage | HIGH | Designed in Phase 7; no persistence layer |
| Tenant-safe retrieval boundaries | HIGH | Must ensure pgvector queries cannot return cross-tenant embeddings |
| Predictive recalculation priority scoring | MEDIUM | Graph heat scores designed; not implemented |
| Drift detection on metadata | MEDIUM | Phase 5 target; no code |
| Ethics / carbon intensity (EEG) | LOW-MEDIUM | OPA policy exists conceptually; no integration |

### 3.3 AI critical path to GTM

The minimum AI story for enterprise sales requires exactly three things working end-to-end:

1. **Live variance narrative** — call a real LLM (even GPT-4o-mini) from the `/narrative/variance` endpoint with the correct driver data model; the hardcoded template creates demo-ware risk in any serious proof-of-concept.
2. **One real anomaly detection signal** — an Isolation Forest running on the ingest stream that emits at least a flag on a known anomaly; this closes the "watchtower" story for CFO trust.
3. **AtomScript MVP** — at minimum a parser that accepts the formula DSL and produces an AST; without this, the "50× faster than Anaplan" claim has no verifiable foundation.

---

## 4. GTM Readiness Matrix (Full)

| Capability area | Enterprise buyer question | Current answer | Gap to close |
|---|---|---|---|
| Tenant isolation | "Can tenant A see tenant B's data?" | Schema says no; runtime unproven | Service-layer middleware + integration test |
| Audit trail ("who changed what, when?") | Right-click cell > show history | Hardcoded stub | Real ClickHouse entropy ledger query |
| Workflow governance ("can I lock a budget node?") | Yes/No with state machine | Schema exists; service-layer enforcement not wired | Workflow state machine in planning-service |
| Connector onboarding ("can my users upload a CSV?") | No | WASM stub | Connector UI wizard + staging airlock |
| Formula language ("can I write IF/SUM/XREF?") | No | Design doc only | AtomScript parser → executor |
| Performance ("sub-second on 50M cells?") | Claimed; unproven | Benchmark harness runs; DB-backed profiles CI pending | Full DB-backed benchmark results with CI evidence |
| AI variance explanation ("why is revenue down?") | Template narrative | No live LLM | Connect litellm to real model |
| Anomaly detection ("flag data entry errors?") | No | Designed | Isolation Forest on ingest stream |
| FX and eliminations ("can I consolidate 50 entities?") | No | Concept in architecture | Consolidation domain pack (Phase 6) |
| ESG/statutory reporting | No | Not mentioned in code | Phase 6+ |
| Scenario generation via natural language | No | Design doc | AtomScript + LLM pipeline |
| Time-travel debugging (infinite undo) | Designed; AODL exists | AODL conceptual; no replay UI | ClickHouse + UI AODL slider |
| Offline conflict resolution (ICH) | Designed | Spec exists; no implementation | Phase 5+ |
| GDPR / data residency | Columns exist | No runtime enforcement | PEV implementation (Phase 5+) |

---

## 5. Gap Summary (Prioritised)

### Tier 1 — Enterprise sales blockers (must close before first contract)

1. **Live LLM variance narrative** — demo-ware risk is disqualifying in any enterprise POC.
2. **Service-layer tenant enforcement** — without runtime isolation proof, enterprise security reviews will block.
3. **Connector UX (CSV / Salesforce upload)** — identified in original red-team as Gap 3; still unaddressed; every enterprise finance team needs it.
4. **Audit history retrieval from real ClickHouse** — "right-click > show history" must work.
5. **One end-to-end workflow lock/unlock from UI to DB** — Gap 4 from the original red-team; state machine schema exists but no API or UI path.

### Tier 2 — Phase 2 sales cycle (needed for renewal and upsell)

6. AtomScript MVP parser and executor.
7. Hot/warm/cold tier promotion daemon.
8. Anomaly detection on ingest stream.
9. Embedding pipeline (atoms → pgvector → similarity retrieval).
10. Field-level encryption (PEV wiring).

### Tier 3 — Full enterprise platform maturity

11. Consolidation, FX, and external reporting domain pack.
12. ESG/statutory reporting path.
13. Bit-packed atom coordinate keys for 25-dim scale.
14. Causal vector clocks for distributed consistency.
15. WASM edge-calc for browser-local formula execution.

---

## 6. Phase Plans 5–8 (Future Agent Reference)

---

### Phase 5: Metadata Intelligence and Visualization Plane

**Objective:** Make metadata first-class, versioned, explorable, and AI-exploitable. Bridge the gap between schema design and runtime intelligence.

**Scope:**
- Semantic identity layer (stable business concept IDs, alias resolution, schema evolution policy)
- Metadata graph query API and lineage graph endpoints
- Drift detection on dimension members and formula dependencies
- Impact analysis view (change scope estimator)
- Hierarchy suggestion workflow with governance approval gate
- PEV field-level encryption wiring (AWS KMS / Sodium) tied to metadata sensitivity labels
- ICH (Intelligent Conflict Handler) — offline CRDT merge policy implementation
- Embedding pipeline: compute and store `embedding_vector` for metadata members via Cortex

**Success criteria:**
- `GET /api/v1/apps/:appId/metadata/graph` returns a lineage graph for any dimension
- Drift detection emits an event when a dimension member changes semantic category
- Hierarchy suggestion is governable through the same metadata promotion approval flow as Phase 3
- Field-level encryption is on by default for columns tagged `sensitivity = high`
- ICH resolves a synthetic offline conflict in an integration test

**Implementation prompt (use in a new agent session):**

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md` and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`. Execute Phase 5. Act as Chief Product Architect, Chief Database Technology Architect, Founder, and AI Intelligence Architect. Deliver: (1) a metadata graph API endpoint in grid-service returning ancestor/descendant/impact chains; (2) a drift detection event when dimension member semantic tags change; (3) an embedding computation job in cortex-service that populates `embedding_vector` columns via pgvector; (4) field-level encryption wiring for `sensitivity = high` columns using the PEV spec in `docs/architecture/privacy-pev.md`; (5) update this playbook with implemented artifacts and remaining risks. Do not over-engineer — build each item to the minimum that makes it testable.

---

### Phase 6: Consolidation and External Reporting Domain Pack

**Objective:** Package the DB and compute layers into a finance-grade consolidation workflow covering multi-entity close, FX translation, intercompany eliminations, and external disclosure mapping.

**Scope:**
- Entity close calendar (periods, deadlines, owner assignment)
- Intercompany ownership metadata and elimination rule definitions
- Journal and adjustment model (governed write-back with immutable audit)
- FX translation policy layers (average rate, closing rate, historical rate selection)
- Disclosure mapping (atom → line item → GAAP/IFRS schedule)
- Drillback from any reported figure to source atoms
- ESG and statutory extension scaffold (emissions per atom, assurance path flag)

**Success criteria:**
- One complete close cycle can be simulated: open → assign → submit → approve → lock → report
- FX translation is a governed calculation with provenance, not a formula hack
- Any disclosed figure can be traced back to its constituent atoms within 2 clicks
- ESG data flows through the same audit path as financial data

**Implementation prompt:**

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md` and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`. Execute Phase 6. As the four canonical roles, design and implement the consolidation and external-reporting domain pack. Deliver: (1) a `09_consolidation_domain.sql` migration adding entity_close_calendar, intercompany_ownership, journal_entries, fx_translation_policies, and disclosure_mappings tables with full tenant_id propagation; (2) a Phase 6 validation SQL file; (3) a close-cycle simulation seed script; (4) a drillback prototype: given a disclosure line item ID, return the contributing atom IDs; (5) update the playbook. Finance correctness and audit-trail completeness are non-negotiable. Do not add UI unless the API contract is unambiguous.

---

### Phase 7: AI-Native Operationalization

**Objective:** Make the AI layer production-grade with live inference, tenant-safe retrieval, model governance, and first-class confidence capture. Close the gap between designed and operational AI.

**Scope:**
- Connect `rag.py` / `litellm` to a real LLM (GPT-4o-mini as default, configurable)
- Auto-baseline endpoint wired to a real time-series model (TimeGPT API or Lag-Llama local)
- Isolation Forest anomaly detection running on the connector ingest event stream
- Semantic similarity retrieval using pgvector on populated embeddings (from Phase 5)
- Variance narrative promoted from template to RAG with real LLM synthesis and grounded citation
- Confidence scoring persisted with each AI output (`ai_inference_log` table)
- Tenant-safe retrieval: pgvector `WHERE tenant_id = $1` on every vector query
- Model provenance: model ID, version, provider, temperature persisted per inference
- Override capture: human corrections journalled to `ai_inference_log` for fine-tuning feedback
- AUH (AI UX Harmonizer) — first suggestion prompt in the UI for dimension naming and hierarchy suggestions

**Success criteria:**
- `/narrative/variance` returns a real LLM-generated narrative for a synthetic variance scenario
- The anomaly detection flag fires on a known outlier in an integration test
- pgvector similarity query never returns cross-tenant results (isolation integration test)
- Every AI inference has a row in `ai_inference_log` with `model_id`, `confidence`, `grounding_atoms[]`
- Human override of an AI suggestion persists to the journal

**Implementation prompt:**

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md`, `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`, and `docs/implementation_guides/layer_8_1_cortex_inference_engine_impl.md`. Execute Phase 7. As the four canonical roles, operationalize the AI layer. Deliver: (1) wire `rag.py` to `litellm` with the OpenAI provider and a model configurable via env var; (2) add a `10_ai_inference_governance.sql` migration for `ai_inference_log` with tenant_id, model_id, confidence, grounding_atoms, human_override; (3) add an Isolation Forest anomaly check on the connector ingest path (can be an in-process Python call from cortex-service synchronously for the POC); (4) add a pgvector similarity search endpoint that enforces tenant_id scoping; (5) add a confidence score and model_id to every Cortex API response; (6) update the playbook. AI safety (no cross-tenant data leakage) is the highest priority constraint.

---

### Phase 8: AtomScript and Formula Intelligence Engine

**Objective:** Build the formula language that converts QuantatomAI from a data platform into an enterprise planning engine. This is the product moat differentiator.

**Scope:**
- AtomScript grammar definition (EBNF or pest grammar file)
- Rust parser producing a typed AST (using `pest` or `nom`)
- AST interpreter in Rust for the core function library (SUM, WHERE, ANCESTOR, DESCENDANTS, YTD, PARALLELPERIOD, XREF)
- Formula dependency DAG construction and topological sort
- Formula execution against the Go grid-service data layer via gRPC
- Monaco editor integration in web UI with basic IntelliSense for dimension members
- Step-through replay debugger (record formula eval steps, expose via API)
- Formula provenance in audit trail (formula_id in metadata_audit_events)
- LLVM JIT compilation as a phase 8.2 stretch goal after interpreter is proven

**Success criteria:**
- A formula `SUM([Revenue]) WHERE [Region] != "Intercompany"` parses without error
- The interpreter produces a correct numeric result against a seeded dataset
- Formula dependency graph is queryable (which cells does this formula depend on?)
- Monaco editor autocompletes dimension member names from live metadata
- Formula change creates an audit event with the old and new formula AST hash

**Implementation prompt:**

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md`, `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`, and `docs/architecture/quantatomai_calculation_engine.md`. Execute Phase 8 (AtomScript MVP). As the four canonical roles, build the minimum viable AtomScript formula engine. Deliver: (1) an EBNF grammar file at `compute/heliocalc/src/atomscript.pest` (or equivalent) covering SUM, WHERE, ANCESTOR, DESCENDANTS, XREF, literal values, dimension references, and filter predicates; (2) a Rust parser that reads a formula string and emits a typed AST (tests for 10 formula examples); (3) a Rust interpreter that evaluates the AST against an in-memory data fixture (unit tests proving SUM and WHERE work correctly); (4) a gRPC endpoint in grid-service that accepts a formula string and returns the evaluated result; (5) Monaco editor integration stub with a single autocomplete source (dimension member names from the metadata API); (6) update the playbook. Do not implement LLVM JIT in this phase. Correctness beats speed until the interpreter is stable.

---

## 7. Recommended Session Order and Cross-Phase Dependencies

```
Phase 5 (Metadata Intelligence + PEV)
   ↓
Phase 7 (AI Operationalization)          ← depends on Phase 5 embedding pipeline
   ↓
Phase 6 (Consolidation Domain Pack)      ← depends on Phase 3 governance + Phase 5 lineage
   ↓
Phase 8 (AtomScript)                     ← can start concurrently with Phase 6
```

**Note:** Phase 5 service-layer tenant enforcement work must be the very first PR in Phase 5 — it is a blocking prerequisite for Phase 7 AI tenant-safe retrieval and for Phase 6 consolidation trust model.

---

## 8. Quick-Start Prompts for Future Sessions

### Resume and assess current state

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md` and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md` to understand current status. Summarize: what phases are complete, what are the top three GTM blockers today, and what is the recommended next concrete step.

### Execute Phase 5

> Read `docs/implementation_guides/database_hardening_multi_agent_playbook.md` and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`. Execute Phase 5 as described in the Phase 5 implementation prompt in the GTM evaluation doc. Begin with service-layer tenant enforcement, then the metadata graph API, then the embedding pipeline.

### Fix the top GTM blocker (live LLM inference)

> Read `services/cortex-service/src/rag.py` and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`. The variance narrative endpoint currently returns deterministic template responses. Wire it to a real LLM using `litellm` with the OpenAI provider (API key via CORTEX_LLM_API_KEY env var). Add a `10_ai_inference_governance.sql` migration to persist inference metadata. Ensure tenant_id is propagated through every LLM call and returned in the response. Update the playbook.

### Wire service-layer tenant enforcement

> Read `services/grid-service/pkg/orchestration/grid_query_service.go`, `docs/implementation_guides/layer_2_4_multi_tenant_control_plane_impl.md`, and `docs/implementation_guides/gtm_readiness_evaluation_2026_03_15.md`. Add tenant enforcement middleware to the grid-service HTTP/gRPC handlers. Every request must carry a `X-Tenant-ID` header or JWT claim. Middleware must validate it matches the row-level tenant_id being queried. Add an integration test that proves a cross-tenant query returns 403 or empty rather than leaking data.

---

## 9. Design Principles That Must Not Be Dropped Under Schedule Pressure

These were established in the playbook and are repeated here as a GTM guard:

- **Never ship AI output without persisted provenance.** Confidence score, model ID, and grounding atoms must be stored before any AI endpoint goes to a customer demo.
- **Never demo the inline audit story without real ClickHouse backing.** The hardcoded stub disqualifies the product in any enterprise security review.
- **Never ship a connector UI without the staging airlock.** Raw CSV data landing directly in the atom lattice bypasses every governance control built in Phases 2 and 3.
- **Never claim multi-tenant isolation without a service-layer enforcement integration test.** Schema constraints are necessary but not sufficient.
- **The AtomScript moat is the product's long-term defensibility.** Prioritise correctness and composability over speed of the JIT compiler. The interpreter must be fully unit-tested before LLVM compilation is considered.
