# QuantatomAI: Planning & Reporting Platform

QuantatomAI is a moat-grade, AI-native planning and reporting platform built on the HelioGrid vFinal architecture:

- Event-driven, wave-based compute (WRM + EWA)
- Atom-based sparse lattice (Hot/Warm/Cold tiers)
- AI-augmented but AI-independent core (Oracle 2.0 + micro-oracles)
- Multi-cloud, ESG-aware, privacy-first (PEV, GRB, MCGF)
- Offline-resilient, migration-focused (ICH, LRI)

## High-level architecture

- **Layer 7 – Experiences:** Web, Excel add-in, external BI, AUH-guided UX
- **Layer 6 – Domain Services:** Modeling, Planning, Actuals, Reconciliation, ALM, Connectors (DFN)
- **Layer 5 – Compute & AI:** HelioCalc (Rust), EWA/WRM, Oracle 2.0, EEG/GRB, VRA/RLE/EVW
- **Layer 4 – Data & Lattice:** Atoms in Hot/Warm/Cold, Metadata, Vector Store, PEV, APWO
- **Layer 3 – Eventing & Sync:** AODL, WRM, Projectors, MCGF
- **Layer 2 – Offline & Integrations:** Offline delta logs, ICH, LRI, connectors
- **Layer 1 – Platform & Multi-cloud:** K8s, Istio/AFM, Crossplane/Terraform, observability, security, SAIC

## Repo layout

- `docs/` – architecture, APIs, dev guides, product docs  
- `infra/` – Kubernetes, Istio, Crossplane, Terraform, CI/CD  
- `services/` – domain microservices (Modeling, Planning, Actuals, etc.)  
- `compute/` – HelioCalc, job scheduler, wave QoS engine  
- `ai/` – Oracle, micro-oracles, AUH, ICH, LRI, EVW, GRB  
- `data/` – atoms schema, hot/warm/cold stores, metadata, vector store  
- `eventing/` – AODL, WRM, projectors, MCGF  
- `ui/` – web app, Excel add-in, design system  
- `tools/` – load testing, migration tools, data generators, admin CLI  
- `tests/` – integration, e2e, performance, chaos

## SLOs (at scale)

- **Grid open (hot slice):** p50 150–300 ms, p95 400–700 ms
- **Grid open (cold/complex):** p50 400–800 ms, p95 800–1,500 ms
- **Single cell edit:** p50 150–300 ms, p95 300–600 ms
- **Medium spread (≤500 cells):** p95 ≤1,200 ms
- **Heavy allocation:** 2–10 s (job with progress)
- **Actuals load (2–3M rows):** available in 2–10 min

## Contracts

### Events (AODL)
`CellUpdated`, `StructureChanged`, `ScenarioCreated`, etc.

### APIs
Documented in `docs/api/rest-endpoints.md` and per-service OpenAPI. See `docs/architecture/quantatomai-architecture.md` for full details.

## Getting started

1. **Clone repo**

   ```bash
   git clone https://github.com/<org>/QuantatomAI.git
   cd QuantatomAI
   ```
