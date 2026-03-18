> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# QuantatomAI Offline ICH Specification

## Purpose
This document defines offline database behavior and Intelligent Conflict Handler semantics for QuantatomAI web and Excel experiences.

## Design Goals
- enable meaningful offline work without breaking lattice integrity
- preserve every local change as durable intent
- resolve collisions with lineage, policy, and ownership awareness
- keep conflict resolution explainable to end users and administrators

## Offline Data Model
- local snapshot: bounded working-set projection of atoms or molecules plus metadata required for editing
- local delta log: append-only record of user intents and local recalculation effects
- local lineage anchors: references to server-side authoritative versions used when snapshot was taken

## Snapshot Rules
- snapshots are scoped by tenant, model, role, and view
- snapshots must carry metadata version id, scenario id, approval state, and policy tags
- snapshots are read-optimized, but every write is recorded as intent rather than silent truth

## Delta Log Rules
Every local edit must persist:
- actor
- time
- original visible value
- intended new value
- local formula or spread context if relevant
- snapshot lineage id
- lock or approval state observed at edit time

## Conflict Classes
### Value conflict
Two parties changed the same intersection or overlapping intersections.

### Workflow conflict
User edited locally, but server-side node became submitted, approved, or locked.

### Metadata conflict
Local snapshot references hierarchy or formula metadata that changed before sync.

### Policy conflict
A privacy, residency, or access policy changed while the client was offline.

## Resolution Policy
- server authority wins for workflow locks and access revocation
- metadata conflicts require remapping or user review, not silent overwrite
- value conflicts use a merge strategy chosen by policy and context
- every automated suggestion must produce a human-readable explanation

## ICH Resolution Modes
- auto-merge safe: only for non-overlapping or provably commutative changes
- guided merge: default for overlapping business edits
- admin escalation: required for policy, approval, or high-risk financial conflicts

## Explainability Requirements
ICH must be able to explain:
- what changed locally
- what changed on the server
- why the conflict was classified the way it was
- what merge options exist
- what audit record will be created after resolution

## Sync Pipeline
1. authenticate and re-evaluate access policy
2. validate metadata version compatibility
3. replay local delta intents against current authoritative state
4. classify conflicts
5. auto-merge safe intents
6. route unresolved conflicts to guided or escalated workflow
7. emit audit records for accepted, rejected, or transformed intents

## Hard Constraints
- offline support must never bypass approval locks
- conflict resolution must never discard local intent without audit trace
- local AI suggestions must not exceed tenant and policy boundaries

## Dependencies
- Phase 3 for workflow and audit trust
- Phase 4 for replay and recovery proof
- Phase 7 for confidence-aware conflict suggestions
