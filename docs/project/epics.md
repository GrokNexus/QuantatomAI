# QuantatomAI – v1 Platform Epics

This document tracks the high-level roadmap and architectural milestones for the QuantatomAI platform.

## Epic 1 – Core Platform & Infra

- **Milestones:**
  - **M1.1:** K8s + Istio base cluster
  - **M1.2:** AODL + WRM + Observability
- **Tasks:**
  - Set up base K8s manifests (`infra/k8s/base`)
  - Configure Istio gateway/virtual services
  - Deploy OpenTelemetry, Prometheus, Grafana
  - Define AODL event schemas (Protobuf)
  - Implement WRM routing + QoS policies

## Epic 2 – Data Lattice & Atoms

- **Milestones:**
  - **M2.1:** Atom schema + metadata
  - **M2.2:** Hot/Warm/Cold integration
- **Tasks:**
  - Define atom schema (`data/atoms/schema`)
  - Implement codecs + lineage tracking
  - Integrate Redis/Scylla (Hot), ClickHouse/DuckDB (Warm), Iceberg (Cold)
  - Implement APWO prototype

## Epic 3 – HelioCalc & Compute

- **Milestones:**
  - **M3.1:** HelioCalc MVP
  - **M3.2:** Wave QoS integration
- **Tasks:**
  - Implement core Rust engine (aggregations, variance, V%, FX)
  - Add job sizing + time budgets
  - Integrate with AODL and WRM
  - Benchmarks + golden tests

## Epic 4 – Domain Services

- **Milestones:**
  - **M4.1:** Modeling + Planning services
  - **M4.2:** Actuals + Reconciliation + ALM
- **Tasks:**
  - Scaffold services from template
  - Implement core APIs and events
  - Wire to data lattice and compute

## Epic 5 – AI & UX Intelligence

- **Milestones:**
  - **M5.1:** Oracle 2.0 + AIProvider abstraction
  - **M5.2:** AUH, ICH, LRI, EVW, GRB
- **Tasks:**
  - Implement AIProvider abstraction (multi-cloud)
  - Build AUH suggestions for BI and conflicts
  - Implement LRI mapping/reconciliation
  - Implement EVW consent-aware upsell

## Epic 6 – Experiences

- **Milestones:**
  - **M6.1:** Web planner
  - **M6.2:** Excel add-in + BI integration
- **Tasks:**
  - Build core grids, forms, scenario UX
  - Implement offline snapshot + delta log
  - Integrate AUH guidance
  - Semantic layer for BI tools
