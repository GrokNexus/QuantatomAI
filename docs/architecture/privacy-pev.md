> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# QuantatomAI Privacy Echo Veil

## Purpose
Privacy Echo Veil defines how sensitive data is classified, protected, masked, and audited across the QuantatomAI database layer.

## Design Goals
- protect finance-sensitive and person-sensitive data at the atom or molecule level
- preserve usability for planning while enforcing policy-aware masking
- support tenant, region, and role-aware privacy decisions
- make privacy behavior explainable and auditable

## Protection Model
PEV operates across four levels:
- classification
- encryption
- masking and minimization
- audit and policy enforcement

## Data Classification Classes
- public operational metadata
- tenant-confidential planning data
- regulated personal or workforce data
- restricted financial close and disclosure data

Each atom or molecule should carry:
- tenant id
- classification tag
- policy tag set
- residency tag
- encryption context reference

## Encryption Strategy
- field-level encryption for high-sensitivity attributes
- per-tenant key hierarchy
- optional per-region subkeys for residency-sensitive data
- rotation support without loss of lineage continuity

## Masking Strategy
- masking is policy-aware, not purely role-aware
- masked views may differ by workflow state, geography, and purpose of use
- AI retrieval must only see the least data necessary for the approved task

Examples:
- workforce planning user can see salary ranges but not raw person-level identifiers
- regional planner can see regional totals but not restricted legal entity detail outside scope

## Policy Evaluation
Policy decisions should consider:
- tenant
- role
- workflow state
- data class
- geography and residency
- access channel such as UI, API, export, or AI retrieval

## Export And AI Boundaries
- exports inherit the highest applicable privacy classification in scope
- embeddings for restricted data must remain tenant-scoped
- no cross-tenant retrieval corpus is allowed by default
- explanation artifacts must not leak masked source content

## Audit Requirements
Every access to restricted or masked data should be attributable by:
- actor
- purpose
- channel
- policy decision
- data class

## Residency Interaction
Privacy and residency policies are linked. If residency blocks movement, privacy cannot be overridden for convenience. Cross-cloud replication of sensitive data requires explicit policy approval and audit.

## Operational Guardrails
- fail closed when policy service is unavailable for restricted data
- prevent debugging and logs from storing decrypted sensitive fields
- require redaction in audit displays shown to non-privileged users

## Phase Dependencies
- Phase 2 must formalize tenant-aware key domains and policy boundaries
- Phase 3 must provide auditable access lineage
- Phase 7 must ensure AI explanations and anomaly flows remain privacy-safe
