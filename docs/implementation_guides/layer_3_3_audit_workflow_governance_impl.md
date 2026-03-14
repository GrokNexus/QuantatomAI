# Layer 3.3 Implementation Guide: Audit, Workflow, and Ingestion Governance

## Status
In progress

## Locations
- Schema migration: [services/grid-service/sql/schema/08_audit_workflow_governance.sql](services/grid-service/sql/schema/08_audit_workflow_governance.sql)
- Validation checks: [services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql](services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql)
- Related phase baseline: [services/grid-service/sql/schema/07_tenant_control_plane.sql](services/grid-service/sql/schema/07_tenant_control_plane.sql)

## Executive Summary
Phase 3 introduces finance-trust controls in the metadata plane by adding:
- immutable metadata audit events
- workflow node state controls and transition rules
- metadata promotion governance records
- connector staging and rejection governance records

This extends the Phase 2 tenant control plane into trust controls needed for enterprise approvals, traceability, and controlled promotions.

## Why This Was Implemented
Phase 2 made tenancy explicit, but enterprise finance buyers also require:
- append-only change history for governance evidence
- controlled workflow state transitions
- explicit promotion records for metadata ALM
- ingestion batch and rejection traceability

These controls provide the first concrete Phase 3 implementation substrate.

## Architecture Decisions
### 1. Audit Is Append-Only
`metadata_audit_events` is an append-only log table. Row-level triggers emit immutable events for `dimensions`, `dimension_members`, `security_policies`, and `branches`.

### 2. Workflow Is Rule-Enforced In SQL
`workflow_state_transitions` is validated by a trigger that rejects invalid transitions. Another trigger applies accepted transitions to `workflow_nodes` lock state and current status.

### 3. Promotion Is Governed As A First-Class Record
`metadata_promotion_requests` records who requested and approved metadata promotion, source branch, status, risk, and timing metadata.

### 4. Connector Staging Is Governed
`connector_ingest_batches` and `connector_ingest_rejections` establish explicit ingestion governance artifacts for batch-level accountability and quarantine traceability.

## Production-Grade Concerns Addressed
### Finance auditability
- append-only change evidence
- actor hooks and operation types
- tenant and app alignment in audit records

### Workflow trust
- explicit transition records
- invalid transition rejection
- lock state updates tied to workflow state

### ALM and release safety
- promotion requests have explicit status and risk model
- branch linkage supports governance checks

### Connector governance
- ingestion records are traceable by tenant and app
- rejected records can be investigated by batch

## Validation Strategy
Run the checks in [services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql](services/grid-service/sql/validation/phase3_audit_workflow_governance_checks.sql).

Expected outcomes:
- invalid workflow node tenant rows are zero
- invalid promotion alignment rows are zero
- orphan connector rejection rows are zero
- invalid metadata audit rows are zero
- invalid audit operation rows are zero

## Open Risks
- Actor and source channel attribution is currently system-default; service-layer propagation is still required.
- Workflow transition policy is intentionally strict and should be expanded with domain-specific transitions over time.
- Cell-value audit trail is still outside this metadata-focused implementation and should be integrated with data-plane writeback path.

## Next Dependencies
- Add service-layer request context propagation into audit event payloads.
- Connect workflow enforcement to writeback and approval APIs.
- Extend promotion governance with diff materialization and approval policy checks.