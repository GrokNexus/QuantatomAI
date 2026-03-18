# QuantatomAI Single Source Of Truth

Status: Canonical

Version: 1.0

Last Updated: 2026-03-17

Owner: Architecture Council

## 1. Purpose

This document is the authoritative architecture and data-journey reference for QuantatomAI.

It consolidates architecture, storage lifecycle, compute flow, API contracts, privacy, offline behavior, and multi-cloud constraints into one normative source.

If any document conflicts with this file, this file takes precedence.

## 2. Scope

In scope:

- 7-layer platform architecture
- end-to-end data journey
- storage and authority model across hot, warm, cold tiers
- eventing and compute flow
- API contract guardrails
- privacy, residency, and audit controls
- offline and conflict handling

Out of scope:

- line-by-line service implementation details
- vendor procurement decisions

## 3. Precedence Rules

Precedence order for design decisions:

1. This file
2. Service-level ADRs approved after this file date
3. Supporting architecture docs

Supporting docs must align to this file.

## 4. Canonical Architecture

### Layer 7: Experience

- Channels: Web, Excel add-in, external BI
- Requirement: responsive grid interaction and progressive data loading
- Requirement: tenant-safe UX for AI suggestions and explanations

### Layer 6: Domain Services

- Services: modeling, planning, actuals, reconciliation, ALM, connectors, metadata and grid APIs
- Responsibility: orchestrate policy checks, metadata resolution, compute dispatch, and response shaping

### Layer 5: Compute And AI

- Core compute: HelioCalc and orchestration execution paths
- Responsibility: dependency-aware recalculation, aggregation, writeback evaluation, scenario operations
- AI posture: AI-augmented but core compute remains AI-independent

### Layer 4: Data And Lattice

- Data model: atom or molecule-based sparse multidimensional substrate
- Stores: hot, warm, cold plus metadata and vector retrieval stores
- Responsibility: lifecycle tiering, lineage continuity, replay compatibility

### Layer 3: Eventing And Sync

- Backbone: append-only event flow and projection pipeline
- Responsibility: durable intent capture, projection materialization, replayability, QoS routing

### Layer 2: Offline And Integrations

- Responsibility: connectors, legacy migration rails, offline snapshots and delta logs, conflict handling

### Layer 1: Platform And Multi-Cloud

- Responsibility: Kubernetes platform, security, observability, policy enforcement integration, regional deployment controls

## 5. Canonical Data Journey

### 5.1 Inbound Journey

1. Source intent arrives from UI writeback, connector ingest, or offline sync replay.
2. Request is authenticated and bound to tenant context.
3. Policy checks execute before restricted reads, writes, exports, or AI retrieval.
4. A durable event or transaction record is written before asynchronous propagation.

### 5.2 Event And Projection Journey

1. Durable write emits append-only event records with tenant and lineage metadata.
2. Stream processors enrich events and materialize read models.
3. Projections update hot and warm access paths.
4. Audit evidence receives append-only provenance entries.

### 5.3 Compute Journey

1. Domain service issues compute request with model, scenario, and dimensional context.
2. Compute resolves dependencies, executes transforms, and returns result deltas.
3. Accepted deltas are persisted and emitted as events.
4. Downstream views and caches refresh from projections, not from untracked side effects.

### 5.4 Read Journey

1. API receives query with tenant-scoped dimensions and filters.
2. Metadata and policy filters are applied.
3. Serving path prefers hot cache, then warm query tier, then cold replay or archive path when required.
4. Response includes only policy-eligible fields and dimensions.

### 5.5 Offline Journey

1. Client snapshot is scoped by tenant, role, model, and lineage anchor.
2. Every offline edit is stored as append-only local intent.
3. On reconnect, intents replay against current authority state.
4. Intelligent Conflict Handler classifies and resolves value, workflow, metadata, and policy conflicts.
5. Resolution outcomes are always auditable.

## 6. Data Authority Model

Authoritative source rules:

- Hot tier is acceleration only, never sole authority.
- Warm tier is interactive analytical authority for active and recent planning windows.
- Cold tier is retention, replay, and regulatory archive authority for historical states.
- Metadata registry is authority for dimensions, hierarchy, policy tags, and tenancy metadata.
- Audit ledger is authority for immutable evidence trail.

### 6.1 Data Authority Matrix

| Artifact Class | Hot Tier | Warm Tier | Cold Tier | Metadata Registry | Audit Ledger |
| --- | --- | --- | --- | --- | --- |
| Active planning cells | cache or projection | authoritative for active cycles | optional archive copy | context only | write events and approvals |
| Current-year actuals | optional cache | authoritative analytical surface | archival copy | context only | ingestion and adjustments |
| Closed periods | none | short analytical window | authoritative long retention | context only | authoritative evidence |
| Hierarchies and dimensions | projection only | projection only | snapshot copy | authoritative | change evidence |
| Offline deltas | local client cache | replay target | historical replay archive | lineage anchor | resolution evidence |
| AI embeddings | optional cache | tenant-scoped serving | governed archive optional | policy tags and scopes | retrieval and export evidence |

## 7. Required Contract Fields

All write, event, replay, export, and AI retrieval contracts must include:

- tenant_id
- lineage_id
- classification_tags
- policy_decision_id
- residency_domain
- idempotency_key for mutating operations
- actor and channel metadata for audit

## 8. Privacy, Residency, And Governance

### 8.1 Privacy Echo Veil Baseline

- classify data by sensitivity class
- enforce field-level protection for high-sensitivity attributes
- apply policy-aware masking based on role, purpose, workflow state, and geography
- fail closed for restricted decisions when policy service is unavailable

### 8.2 Residency Baseline

- tenant residency policy overrides optimization policy
- cross-cloud movement requires explicit policy evaluation and audit record
- encryption context must remain intact across movement

### 8.3 Audit Baseline

All restricted access and all mutating operations must be attributable by actor, purpose, channel, and policy decision.

## 9. Eventing And Replay Rules

- event schema must be stable and versioned
- replay must be idempotent
- projection rebuild must be possible from durable event history and snapshots
- no critical business state can exist only in ephemeral cache

## 10. Performance And SLO Baseline

Target behavior:

- grid open hot path: sub-second p95
- single cell edit: sub-second p95
- medium spread operations: low-second p95
- heavy allocations: asynchronous with progress and audit trace

SLO tuning can evolve, but authority and audit guarantees in this file are invariant.

## 11. Security Baseline

- tenant isolation by default
- least privilege access across services and data planes
- encrypted transport and encrypted sensitive data at rest
- no plaintext sensitive data in logs

## 12. Multi-Cloud Operating Model

- global control-plane semantics for policy vocabulary, metadata contract, and audit taxonomy
- regional data-plane autonomy with policy-governed federation
- safe regional degrade mode during partition events

## 13. Implementation Conformance Checklist

A service or workflow is conformant only if:

- it implements required contract fields
- it preserves authority model semantics
- it emits auditable provenance for mutating and restricted operations
- it passes policy checks before restricted access and movement
- it supports replay and recovery without hidden state

## 14. Migration Plan For Existing Docs

The following docs must be updated to reference this canonical source and remove conflicting language:

- docs/architecture/quantatomai-architecture.md
- docs/architecture/layer_1_2_foundation_spec.md
- docs/architecture/layer_3_4_spine_spec.md
- docs/architecture/layer_5_compute_spec.md
- docs/architecture/layer_6_7_experience_spec.md
- docs/architecture/quantatomai-master-schema.md
- docs/architecture/ingestion-pipeline.md
- docs/architecture/atoms-lattice.md
- docs/architecture/wave-qos.md
- docs/api/rest-endpoints.md

## 15. Changelog

- 1.0: initial canonical consolidation of architecture and data journey