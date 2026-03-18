> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# Wave QoS

## Purpose

Define Quality of Service expectations for eventing, projection, and interactive grid operations.

## QoS Principles

- Durable event capture precedes asynchronous propagation.
- Projection paths prioritize consistency with authoritative sources.
- Replay compatibility and idempotency are mandatory.
- QoS routing must preserve tenant and policy constraints.

## Performance Baseline

Target p95 behavior:

- Grid open hot path: sub-second.
- Single cell edit: sub-second.
- Medium spread operations: low-second.
- Heavy allocations: asynchronous, with progress and audit trace.

## Eventing Guarantees

- Event schema is stable and versioned.
- Replay is idempotent.
- Projection rebuild is possible from durable history plus snapshots.
- Critical business state never exists only in ephemeral cache.

## Operational Safeguards

- Emit auditable provenance for mutating and restricted operations.
- Fail closed for policy-restricted operations when policy service is unavailable.
- Preserve residency and encryption context across inter-region movement.

