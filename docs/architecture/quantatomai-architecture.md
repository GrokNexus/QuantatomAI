# QuantatomAI Architecture

## 1. Overview
QuantatomAI is an AI-native, multi-cloud, and atom-based planning substrate.

## 2. Core substrate: Atoms in SAS
- Atoms as the fundamental unit of data (substrate).
- Sparse Atom Storage (SAS) optimizes for multidimensionality and sparsity.

## 3. The 7-Layer Model
- **Layer 7 – Experiences:** Web, Excel, BI interfaces.
- **Layer 6 – Domain Services:** Modeling, Planning, Actuals (DFN).
- **Layer 5 – Compute & AI:** HelioCalc (Rust), Oracle 2.0, EEG/GRB.
- **Layer 4 – Data & Lattice:** Atom stores (Hot/Warm/Cold), PEV, APWO.
- **Layer 3 – Eventing & Sync:** AODL, WRM, Projectors.
- **Layer 2 – Offline & Integrations:** ICH, LRI, connectors.
- **Layer 1 – Platform:** K8s, Istio (AFM), Crossplane (SAIC).

## 4. Architectural Components (Acronyms)
- **WRM:** Wave Resilience Mesh (Istio/Envoy/OTel)
- **PEV:** Privacy Echo Veil (AWS KMS/Sodium field-level crypto)
- **APWO:** Auto-Part Wave Optimizer (Ray/MLflow partitioning)
- **DFN:** Decentralized Flow Nexus (Micro-oracles/ONNX)
- **AUH:** AI UX Harmonizer (LlamaIndex/Local LLM suggestions)
- **EWA:** Elastic Wave Accelerator (SIMD/GPU-optional)
- **EEG:** Ethics & Green Resonance Balancer (OPA/Carbon Intensity)
- **GRB:** Governance & Reliability Bridge (Prometheus/Grafana)
- **VRA / RLE:** Valor Resilience Engine (Monetization/Billing)
- **EVW:** Ethical Value Wave (Consent/Mixpanel)
- **AODL:** Append-Only Delta Log (Kinesis/EventHub/PubSub)
- **MCGF:** Multi-Cloud Governance Fabric (Elasticsearch/Audit trails)
- **LRI:** Legacy Resonance Infuser (Polars/Scikit-learn migration)
- **ICH:** Intelligent Conflict Handler (Offline resolution)
- **AFM:** Auto-Federation Mesh (Istio multi-cluster)
- **SAIC:** Source-Agnostic Infrastructure Control (Integrity jobs)

## 5. 7-Layer Architecture Blueprint

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                           LAYER 7 — EXPERIENCES                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  Web UI (React / Next.js, React Query, Zustand, TanStack Table)              │
│    • AUH: AI UX Harmonizer (LlamaIndex / local LLMs for suggestions)         │
│    • Offline Snapshot + Delta Log (IndexedDB)                                │
│                                                                              │
│  Excel Add‑in (OfficeJS, Custom Functions, Task Pane APIs)                   │
│                                                                              │
│  External BI (Power BI, Tableau, Looker, Oracle Analytics)                   │
│    • Semantic Layer (GraphQL Federation / Cube.js)                           │
│    • AUH “Explain / Suggest Next View”                                       │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                     LAYER 6 — DOMAIN MICROSERVICES (DFN)                     │
├──────────────────────────────────────────────────────────────────────────────┤
│  Microservices (Kubernetes Pods, gRPC/REST):                                 │
│    • Modeling Service (Node.js / Go / Rust)                                  │
│    • Planning Service (Go / Rust)                                             │
│    • Actuals Service (Python + Pandas / DBT / Airbyte)                       │
│    • Reconciliation Service (Go / Rust)                                      │
│    • ALM Service (Go / Node.js)                                              │
│    • Connectors Service (Python, Airbyte, Fivetran)                          │
│                                                                              │
│  DFN: Decentralized Flow Nexus                                               │
│    • Micro‑oracles (ONNX Runtime / TinyLlama / Distilled Models)             │
│    • Event-driven coordination (AODL + WRM)                                  │
│                                                                              │
│  Infra: K8s HPA/VPA, Envoy Sidecars, Istio mTLS                              │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                          LAYER 5 — COMPUTE & AI                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  HelioCalc Engine Pool                                                       │
│    • Rust + SIMD + Rayon                                                     │
│    • GPU-optional (CUDA / ROCm)                                              │
│    • Vectorized aggregations, FX, eliminations                               │
│                                                                              │
│  EWA + WRM: Elastic Wave Accelerator + Wave Resilience Mesh                  │
│    • Istio Mesh + Envoy                                                      │
│    • OpenTelemetry (Tracing)                                                 │
│    • Priority Queues (Redis Streams / NATS JetStream)                        │
│                                                                              │
│  Global Oracle 2.0 (Async AI)                                                │
│    • Azure OpenAI / AWS Bedrock / Vertex AI / OCI AI                         │
│    • Local fallback models (GGUF + llama.cpp)                                │
│    • MLflow for model lifecycle                                              │
│                                                                              │
│  EEG + GRB: Ethics & Green Resonance Balancer                                │
│    • OPA (Open Policy Agent) for PII/HUC                                     │
│    • Carbon Intensity API + Scheduler                                        │
│                                                                              │
│  VRA / RLE + EVW: Monetization & Loyalty                                     │
│    • Mixpanel (opt‑in analytics)                                             │
│    • Consent Manager (Open Source CMP)                                       │
│    • Stripe / Chargebee for billing                                          │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                          LAYER 4 — DATA & LATTICE                            │
├──────────────────────────────────────────────────────────────────────────────┤
│  Hot Store: HelioHot                                                         │
│    • Redis Cluster / ScyllaDB                                                │
│    • LRU + LFU eviction policies                                             │
│                                                                              │
│  Warm Store: Sparse Atom Store (SAS)                                         │
│    • ClickHouse (MergeTree) / DuckDB                                         │
│    • Iceberg-backed partitions                                                │
│    • Write buffers (Kafka Connect / Debezium)                                │
│                                                                              │
│  Cold Store: Iceberg on Object Storage                                       │
│    • S3 / GCS / Azure Blob / OCI Object Storage                              │
│                                                                              │
│  Metadata: Postgres (Citus optional)                                         │
│    • RLS, JSONB, Partitioning                                                │
│                                                                              │
│  Vector Store: Qdrant / Weaviate / pgvector                                  │
│                                                                              │
│  PEV: Privacy Echo Veil                                                      │
│    • Encryption: AWS KMS / Azure KeyVault / HashiCorp Vault                  │
│    • Sodium (libsodium) for field-level crypto                               │
│                                                                              │
│  APWO: Auto-Part Wave Optimizer                                              │
│    • Ray (distributed compute)                                               │
│    • MLflow (pattern models)                                                 │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                         LAYER 3 — EVENTING & SYNC                            │
├──────────────────────────────────────────────────────────────────────────────┤
│  AODL: Append-Only Delta Log                                                 │
│    • AWS Kinesis / Azure EventHub / GCP PubSub / OCI Streaming               │
│    • Protobuf events                                                         │
│                                                                              │
│  WRM: Wave Resilience Mesh                                                   │
│    • Istio + Envoy + OpenTelemetry                                           │
│    • Retry/backoff, wave routing, QoS                                        │
│                                                                              │
│  Projectors                                                                  │
│    • Kafka Connect / Flink / Spark Structured Streaming                      │
│    • Materialize Hot/Warm/Read Models                                        │
│                                                                              │
│  MCGF: Multi-Cloud Governance Fabric                                         │
│    • Fluentd / FluentBit                                                     │
│    • Elasticsearch / OpenSearch                                              │
│    • Holographic audit trails                                                │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                     LAYER 2 — OFFLINE & INTEGRATIONS                         │
├──────────────────────────────────────────────────────────────────────────────┤
│  Offline Web & Excel                                                         │
│    • IndexedDB, Service Workers                                              │
│    • ICH + AUH (ML-assisted conflict resolution)                             │
│                                                                              │
│  Connectors                                                                  │
│    • Airbyte / Fivetran / DBT                                                │
│    • ERP/GL/HR/CRM (Oracle Fusion, SAP, Workday, Netsuite, Dynamics)         │
│    • Data Warehouses (Snowflake, BigQuery, Redshift, Synapse)                │
│                                                                              │
│  LRI: Legacy Resonance Infuser                                               │
│    • Pandas / Polars                                                         │
│    • Scikit-learn (vector similarity)                                        │
│    • Dual-hierarchy staging                                                  │
└──────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                     LAYER 1 — PLATFORM & MULTI-CLOUD                         │
├──────────────────────────────────────────────────────────────────────────────┤
│  Kubernetes: AKS / EKS / GKE / OKE                                           │
│  Istio + AFM: Auto-Federation Mesh                                           │
│  Crossplane: Multi-cloud infra orchestration                                 │
│                                                                              │
│  StorageProvider: S3 / GCS / Azure Blob / OCI                                │
│  IdentityProvider: Azure AD / AWS IAM / Google Identity / Oracle IDCS        │
│  AIProvider: Azure OpenAI / AWS Bedrock / Vertex AI / OCI AI                 │
│                                                                              │
│  Observability: OpenTelemetry, Prometheus, Grafana                           │
│  Security: mTLS, RLS, TLS, PEV encryption                                    │
│  SAIC: Self-Auditing Integrity Core (scheduled jobs)                         │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 6. SLO Sheet
| Metric | SLO (p50) | SLO (p95) | Notes |
|--------|-----------|-----------|-------|
| Grid open (hot) | 150-300ms | 400-700ms | Sub-second interactive response |
| Grid open (cold) | 400-800ms | 800-1,500ms | Background fetching active |
| Single cell edit | 150-300ms | 300-600ms | Atomic update replication |
| Medium spread | N/A | < 1,200ms | Parallel computation wave |
| Heavy allocation | 2-10s | N/A | Significant compute job |
| Actuals load | 2-10 min | N/A | Large scale ingestion |
