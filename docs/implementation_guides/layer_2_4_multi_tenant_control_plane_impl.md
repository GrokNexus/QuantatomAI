# Layer 2.4 Implementation Guide: Multi-Tenant Control Plane

## Status
In progress

## Locations
- Schema migration: [services/grid-service/sql/schema/07_tenant_control_plane.sql](services/grid-service/sql/schema/07_tenant_control_plane.sql)
- Validation checks: [services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql](services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql)
- Automated verification: [services/grid-service/sql/schema/migrate_integration_test.go](services/grid-service/sql/schema/migrate_integration_test.go)
- Related architecture: [docs/architecture/privacy-pev.md](docs/architecture/privacy-pev.md), [docs/architecture/multi-cloud-federation.md](docs/architecture/multi-cloud-federation.md), [docs/architecture/database-layer-mind-map.md](docs/architecture/database-layer-mind-map.md)

## Executive Summary
Phase 2 hardens the database layer so tenancy is explicit in control-plane metadata rather than inferred inconsistently across services.

This implementation adds:
- tenant region registry
- tenant key-domain registry
- tenant quota and chargeback policy
- tenant AI boundary policy
- app partition registry
- direct tenant propagation into core metadata tables
- synchronization triggers that keep tenant and app context aligned

## Why This Was Implemented
The base schema already had `tenants`, `users`, and `apps`, but core metadata rows still relied too heavily on indirect tenancy via joins. That is not sufficient for production-grade control over:
- residency
- isolation
- quotas
- partitioning
- AI data boundaries
- cache and event naming

The new control-plane layer closes that gap.

## Architecture Decisions
### 1. Tenant remains the top isolation boundary
Tenant is the primary security, quota, encryption, and AI boundary. Applications, dimensions, members, branches, and policies inherit from that boundary.

### 2. Region registry is explicit
Multi-region behavior must be policy-driven. A tenant can have multiple regions, but only one write region by default. This enables clean failover and residency controls.

### 3. Key domains are purpose-bound
Keys are not generic tenant secrets. They are bound by region and purpose:
- app-data
- audit
- embedding
- export
- backup

This supports future finance-grade key rotation and export isolation.

### 4. AI policy is tenant-scoped by default
The schema defaults to tenant-only retrieval and disallows cross-tenant learning. This is the safe default for an enterprise finance platform.

### 5. App partitions make physical placement explicit
Each app now has a registry entry for:
- write region
- hot namespace
- warm partition template
- cold object prefix
- event topic prefix
- cache namespace

This creates a bridge between logical model ownership and actual data-plane placement.

### 6. Core metadata tables carry tenant context directly
The migration adds tenant context directly to:
- dimensions
- dimension_members
- security_policies
- branches

This reduces ambiguous joins in enforcement code and enables tenant-aware indexing.

## Production-Grade Concerns Addressed
### Tenant isolation
- tenant context is materialized on key metadata tables
- synchronization triggers enforce consistent derivation
- app partition registry keeps cache and event topology tenant-bounded

### Privacy and residency
- tenant regions and key domains bind storage behavior to region-aware policy
- tenant AI policy defaults prevent cross-tenant retrieval leakage

### Performance and cost
- quota policy supports noisy-neighbor mitigation design
- app partitions enable cache namespace and event topic isolation

### Future enterprise finance needs
- consolidation and reporting workloads can be region-aware and tenant-bounded
- audit, export, and AI keys can diverge by purpose without redesigning the schema

## Validation Strategy
Run the validation queries in [services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql](services/grid-service/sql/validation/phase2_tenant_control_plane_checks.sql).

Run the integration test in [services/grid-service/sql/schema/migrate_integration_test.go](services/grid-service/sql/schema/migrate_integration_test.go) against the local grid-service Postgres stack. The workspace task labels are:
- `grid_service_phase2_start_deps`
- `grid_service_phase2_test`

Expected outcomes:
- all invalid row counts are zero
- no tenant has more than one write region
- no AI policy permits cross-tenant learning by default
- every app partition points to a registered tenant write region
- automated migration test proves trigger-based tenant propagation, AI policy defaults, write-region uniqueness, and region-bound key-domain enforcement

## Open Risks
- This phase establishes the control-plane schema, but service-layer enforcement still needs to consume it consistently.
- Quotas are modeled here, but active throttling and fairness enforcement belong in later service and benchmarking phases.
- Tenant-aware vector namespace enforcement still needs application-layer implementation in the AI stack.

## Next Dependencies
- Phase 3 should extend this with workflow- and audit-bound trust controls.
- Phase 4 should prove fairness, replay, and performance under tenant-mixed workloads.
- Phase 7 should bind AI inference and retrieval to `tenant_ai_policies` and `tenant_key_domains`.