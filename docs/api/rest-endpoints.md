> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# REST Endpoints

## Purpose

Define SSOT-aligned REST contract guardrails for grid operations, ingestion, replay, and governance-sensitive actions.

## Endpoint Groups

- Grid query and writeback endpoints.
- Planning/modeling scenario endpoints.
- Connector ingestion and replay endpoints.
- Metadata and hierarchy endpoints.
- Audit and lineage retrieval endpoints.

## Required Request Contract Fields

Mutating, replay, export, and AI-retrieval requests must include:

- tenant_id
- lineage_id
- classification_tags
- policy_decision_id
- residency_domain
- idempotency_key
- actor and channel metadata

## Required Response Behaviors

- Responses are tenant-scoped and policy-filtered.
- Serving path follows hot to warm to cold/replay preference.
- Restricted fields are removed or masked based on policy decision.
- Mutating responses include lineage and audit correlation identifiers.

## Error And Governance Semantics

- Use fail-closed behavior for restricted operations when policy service is unavailable.
- Enforce idempotency on mutating endpoints.
- Preserve full audit attribution for restricted reads and all writes.

## Notes

Concrete endpoint paths and payload schemas must remain consistent with:

- docs/api/grid-api.md
- docs/api/graphql-schema.md
- docs/architecture/quantatomai-single-source-of-truth.md

