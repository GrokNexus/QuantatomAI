> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# ASCII Diagram

```text
Layer 7  Experience
	Web / Excel Add-in / External BI
						|
						v
Layer 6  Domain Services
	Modeling | Planning | Actuals | Reconciliation | ALM | Connectors | Grid API
						|
						v
Layer 5  Compute and AI
	HelioCalc + orchestration + policy-aware AI augmentation
						|
						v
Layer 4  Data and Lattice
	Atom/Molecule lattice | Hot cache | Warm analytical store | Cold archive/replay
						|
						v
Layer 3  Eventing and Sync
	Append-only events | Stream processors | Projections | QoS routing
						|
						v
Layer 2  Offline and Integrations
	Offline snapshots/delta logs | ICH conflict resolution | Connector ingestion
						|
						v
Layer 1  Platform and Multi-Cloud
	Kubernetes | Security | Observability | Policy enforcement | Regional controls

Cross-cutting authorities:
	Metadata Registry (dimensions/hierarchy/policy tags/tenancy)
	Audit Ledger (immutable evidence trail)
```

## Data Journey Summary

1. Intent enters via UI writeback, connector ingest, or offline replay.
2. Authentication, tenant binding, and policy checks execute.
3. Durable write occurs before asynchronous propagation.
4. Append-only events drive projection updates for serving tiers.
5. Reads prefer hot, then warm, then cold/replay path as required.

