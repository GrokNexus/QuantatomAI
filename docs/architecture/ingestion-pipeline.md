> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# Ingestion Pipeline

## Purpose

Define the canonical inbound ingestion path for connector feeds, UI writebacks, and offline replay intents.

## Canonical Flow

1. Input intent arrives from connector ingestion, UI writeback, or offline replay.
2. Identity and tenant context are bound to the request.
3. Policy checks run before restricted reads, writes, exports, and AI retrieval.
4. A durable write is recorded before asynchronous fan-out.
5. Append-only events are emitted with lineage and actor metadata.
6. Stream processors enrich and project data into serving tiers.

## Required Metadata

Every ingest operation must include:

- tenant_id
- lineage_id
- classification_tags
- policy_decision_id
- residency_domain
- idempotency_key for mutating operations
- actor and channel metadata

## Authority And Replay Semantics

- Hot tier is acceleration only and not authoritative.
- Warm tier is authoritative for active analytical/planning windows.
- Cold tier is authoritative for retained historical states and replay archives.
- Replay must be idempotent and reconstruct projections from durable history.

## Failure Handling

- Fail closed for restricted operations when policy evaluation is unavailable.
- Preserve append-only audit evidence for every restricted or mutating action.
- Avoid hidden state outside durable events, metadata authority, and audit ledger.

