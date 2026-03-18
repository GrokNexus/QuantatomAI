> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# Atoms Lattice

## Purpose

Specify the canonical atom or molecule-based sparse multidimensional lattice used for planning and reporting.

## Canonical Model

- Atom: smallest addressable unit at the intersection of dimensions and measures.
- Molecule: grouped logical unit composed of related atoms for scenario and workflow operations.
- Lattice: sparse multidimensional substrate that supports projection, replay, and lineage continuity.

## Dimension And Hierarchy Authority

- Metadata registry is authoritative for dimensions, hierarchies, policy tags, and tenancy metadata.
- Serving tiers may project metadata but do not replace metadata authority.

## Lifecycle Semantics

- Active-cycle analytical authority resides in warm tier.
- Hot tier provides acceleration cache/projections only.
- Cold tier retains historical states for replay and regulation.

## Contract Requirements

All lattice writes, events, and exports must carry:

- tenant_id
- lineage_id
- classification_tags
- policy_decision_id
- residency_domain
- idempotency_key for mutating operations
- actor and channel metadata

## Conformance Rules

- Mutations emit append-only events.
- Projection rebuild from durable event history must remain possible.
- No critical state is allowed to exist only in ephemeral cache.

