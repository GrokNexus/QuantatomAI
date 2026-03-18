> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# QuantatomAI Multi-Cloud Federation

## Purpose
This document defines how QuantatomAI operates across cloud boundaries without breaking tenant isolation, privacy policy, or database consistency guarantees.

## Design Goals
- preserve tenant and data-residency controls across clouds
- keep metadata, audit, and policy behavior consistent across providers
- support regional deployment without forcing a single cloud vendor
- avoid uncontrolled cross-cloud data drift

## Federation Principles
- Control-plane consistency is more important than perfect data symmetry.
- Metadata, policy, and audit semantics must be uniform even if underlying engines differ by cloud.
- Cross-cloud movement is policy-driven, not convenience-driven.
- Tenant residency policy overrides optimization policy.

## Federation Model
### Global Control Plane
- global metadata contract
- global policy vocabulary
- global tenant registry
- global audit taxonomy

### Regional Data Planes
Each region or cloud deployment may host its own:
- hot tier instance
- warm analytical tier
- cold archive bucket
- eventing cluster or bridge
- privacy and key-management integration

### Portable Data Protocol
QuantatomAI molecules or atoms are the portability layer.
This enables consistent data semantics across:
- AWS
- Azure
- GCP
- OCI
- sovereign or on-prem extensions where required

## What Can Be Federated
- metadata definitions
- model packages and approved ALM promotions
- policy definitions
- AI model routing metadata
- audit taxonomy and reporting shape

## What Must Remain Residency-Bound Unless Policy Allows
- tenant operational data
- regulated personal data
- customer-specific embeddings
- raw audit evidence tied to regional control requirements

## Cross-Cloud Movement Rules
- no data movement without explicit policy evaluation
- all movement must preserve tenant id, lineage, encryption context, and classification tags
- warm and cold data replication must be class-aware and residency-aware
- audit exports must record cloud origin, destination, reason, and policy decision

## Identity And Key Management
- each cloud may use its native KMS or secret system
- QuantatomAI must maintain a normalized key-domain abstraction per tenant and per region
- encryption keys must not be silently shared across clouds without policy approval

## Event Federation
- event schemas remain globally consistent
- event routing may be regional with explicit bridge points
- replay boundaries must be documented by region and tenant
- cross-cloud event replication must be idempotent and auditable

## Consistency Model
- metadata and policy: strongly governed, promotion-based consistency
- operational planning data: region-local authority with governed replication where required
- audit: append-only and provenance-preserving, with global reporting possible through federated views

## Failure Model
- loss of one regional hot plane must not compromise the authoritative audit and metadata control planes
- cross-cloud network partition must degrade to region-local safe mode, not data corruption mode
- policy engine failure must fail closed for restricted movement

## Board-Level Constraint
Multi-cloud is not a feature until residency, audit, and tenant-safety behaviors are verifiable. This document provides the operating model; Phase 2 and Phase 4 must provide the proof.
