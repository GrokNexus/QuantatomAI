> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# QuantatomAI Architecture

## 1. Overview
QuantatomAI is an AI-native, multi-cloud, and atom-based planning substrate.

## 2. Core substrate: Atoms in SAS
- Atoms as the fundamental unit of data (substrate).
- Sparse Atom Storage (SAS) optimizes for multidimensionality and sparsity.

## 3. The 7-Layer Model
- **Layer 7 â€“ Experiences:** Web, Excel, BI interfaces.
- **Layer 6 â€“ Domain Services:** Modeling, Planning, Actuals (DFN).
- **Layer 5 â€“ Compute & AI:** HelioCalc (Rust), Oracle 2.0, EEG/GRB.
- **Layer 4 â€“ Data & Lattice:** Atom stores (Hot/Warm/Cold), PEV, APWO.
- **Layer 3 â€“ Eventing & Sync:** AODL, WRM, Projectors.
- **Layer 2 â€“ Offline & Integrations:** ICH, LRI, connectors.
- **Layer 1 â€“ Platform:** K8s, Istio (AFM), Crossplane (SAIC).

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
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                           LAYER 7 â€” EXPERIENCES                              â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  Web UI (React / Next.js, React Query, Zustand, TanStack Table)              â”‚
â”‚    â€¢ AUH: AI UX Harmonizer (LlamaIndex / local LLMs for suggestions)         â”‚
â”‚    â€¢ Offline Snapshot + Delta Log (IndexedDB)                                â”‚
â”‚                                                                              â”‚
â”‚  Excel Addâ€‘in (OfficeJS, Custom Functions, Task Pane APIs)                   â”‚
â”‚                                                                              â”‚
â”‚  External BI (Power BI, Tableau, Looker, Oracle Analytics)                   â”‚
â”‚    â€¢ Semantic Layer (GraphQL Federation / Cube.js)                           â”‚
â”‚    â€¢ AUH â€œExplain / Suggest Next Viewâ€                                       â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                     LAYER 6 â€” DOMAIN MICROSERVICES (DFN)                     â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  Microservices (Kubernetes Pods, gRPC/REST):                                 â”‚
â”‚    â€¢ Modeling Service (Node.js / Go / Rust)                                  â”‚
â”‚    â€¢ Planning Service (Go / Rust)                                             â”‚
â”‚    â€¢ Actuals Service (Python + Pandas / DBT / Airbyte)                       â”‚
â”‚    â€¢ Reconciliation Service (Go / Rust)                                      â”‚
â”‚    â€¢ ALM Service (Go / Node.js)                                              â”‚
â”‚    â€¢ Connectors Service (Python, Airbyte, Fivetran)                          â”‚
â”‚                                                                              â”‚
â”‚  DFN: Decentralized Flow Nexus                                               â”‚
â”‚    â€¢ Microâ€‘oracles (ONNX Runtime / TinyLlama / Distilled Models)             â”‚
â”‚    â€¢ Event-driven coordination (AODL + WRM)                                  â”‚
â”‚                                                                              â”‚
â”‚  Infra: K8s HPA/VPA, Envoy Sidecars, Istio mTLS                              â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                          LAYER 5 â€” COMPUTE & AI                              â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  HelioCalc Engine Pool                                                       â”‚
â”‚    â€¢ Rust + SIMD + Rayon                                                     â”‚
â”‚    â€¢ GPU-optional (CUDA / ROCm)                                              â”‚
â”‚    â€¢ Vectorized aggregations, FX, eliminations                               â”‚
â”‚                                                                              â”‚
â”‚  EWA + WRM: Elastic Wave Accelerator + Wave Resilience Mesh                  â”‚
â”‚    â€¢ Istio Mesh + Envoy                                                      â”‚
â”‚    â€¢ OpenTelemetry (Tracing)                                                 â”‚
â”‚    â€¢ Priority Queues (Redis Streams / NATS JetStream)                        â”‚
â”‚                                                                              â”‚
â”‚  Global Oracle 2.0 (Async AI)                                                â”‚
â”‚    â€¢ Azure OpenAI / AWS Bedrock / Vertex AI / OCI AI                         â”‚
â”‚    â€¢ Local fallback models (GGUF + llama.cpp)                                â”‚
â”‚    â€¢ MLflow for model lifecycle                                              â”‚
â”‚                                                                              â”‚
â”‚  EEG + GRB: Ethics & Green Resonance Balancer                                â”‚
â”‚    â€¢ OPA (Open Policy Agent) for PII/HUC                                     â”‚
â”‚    â€¢ Carbon Intensity API + Scheduler                                        â”‚
â”‚                                                                              â”‚
â”‚  VRA / RLE + EVW: Monetization & Loyalty                                     â”‚
â”‚    â€¢ Mixpanel (optâ€‘in analytics)                                             â”‚
â”‚    â€¢ Consent Manager (Open Source CMP)                                       â”‚
â”‚    â€¢ Stripe / Chargebee for billing                                          â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                          LAYER 4 â€” DATA & LATTICE                            â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  Hot Store: HelioHot                                                         â”‚
â”‚    â€¢ Redis Cluster / ScyllaDB                                                â”‚
â”‚    â€¢ LRU + LFU eviction policies                                             â”‚
â”‚                                                                              â”‚
â”‚  Warm Store: Sparse Atom Store (SAS)                                         â”‚
â”‚    â€¢ ClickHouse (MergeTree) / DuckDB                                         â”‚
â”‚    â€¢ Iceberg-backed partitions                                                â”‚
â”‚    â€¢ Write buffers (Kafka Connect / Debezium)                                â”‚
â”‚                                                                              â”‚
â”‚  Cold Store: Iceberg on Object Storage                                       â”‚
â”‚    â€¢ S3 / GCS / Azure Blob / OCI Object Storage                              â”‚
â”‚                                                                              â”‚
â”‚  Metadata: Postgres (Citus optional)                                         â”‚
â”‚    â€¢ RLS, JSONB, Partitioning                                                â”‚
â”‚                                                                              â”‚
â”‚  Vector Store: Qdrant / Weaviate / pgvector                                  â”‚
â”‚                                                                              â”‚
â”‚  PEV: Privacy Echo Veil                                                      â”‚
â”‚    â€¢ Encryption: AWS KMS / Azure KeyVault / HashiCorp Vault                  â”‚
â”‚    â€¢ Sodium (libsodium) for field-level crypto                               â”‚
â”‚                                                                              â”‚
â”‚  APWO: Auto-Part Wave Optimizer                                              â”‚
â”‚    â€¢ Ray (distributed compute)                                               â”‚
â”‚    â€¢ MLflow (pattern models)                                                 â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                         LAYER 3 â€” EVENTING & SYNC                            â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  AODL: Append-Only Delta Log                                                 â”‚
â”‚    â€¢ AWS Kinesis / Azure EventHub / GCP PubSub / OCI Streaming               â”‚
â”‚    â€¢ Protobuf events                                                         â”‚
â”‚                                                                              â”‚
â”‚  WRM: Wave Resilience Mesh                                                   â”‚
â”‚    â€¢ Istio + Envoy + OpenTelemetry                                           â”‚
â”‚    â€¢ Retry/backoff, wave routing, QoS                                        â”‚
â”‚                                                                              â”‚
â”‚  Projectors                                                                  â”‚
â”‚    â€¢ Kafka Connect / Flink / Spark Structured Streaming                      â”‚
â”‚    â€¢ Materialize Hot/Warm/Read Models                                        â”‚
â”‚                                                                              â”‚
â”‚  MCGF: Multi-Cloud Governance Fabric                                         â”‚
â”‚    â€¢ Fluentd / FluentBit                                                     â”‚
â”‚    â€¢ Elasticsearch / OpenSearch                                              â”‚
â”‚    â€¢ Holographic audit trails                                                â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                     LAYER 2 â€” OFFLINE & INTEGRATIONS                         â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  Offline Web & Excel                                                         â”‚
â”‚    â€¢ IndexedDB, Service Workers                                              â”‚
â”‚    â€¢ ICH + AUH (ML-assisted conflict resolution)                             â”‚
â”‚                                                                              â”‚
â”‚  Connectors                                                                  â”‚
â”‚    â€¢ Airbyte / Fivetran / DBT                                                â”‚
â”‚    â€¢ ERP/GL/HR/CRM (Oracle Fusion, SAP, Workday, Netsuite, Dynamics)         â”‚
â”‚    â€¢ Data Warehouses (Snowflake, BigQuery, Redshift, Synapse)                â”‚
â”‚                                                                              â”‚
â”‚  LRI: Legacy Resonance Infuser                                               â”‚
â”‚    â€¢ Pandas / Polars                                                         â”‚
â”‚    â€¢ Scikit-learn (vector similarity)                                        â”‚
â”‚    â€¢ Dual-hierarchy staging                                                  â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜


â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                     LAYER 1 â€” PLATFORM & MULTI-CLOUD                         â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚  Kubernetes: AKS / EKS / GKE / OKE                                           â”‚
â”‚  Istio + AFM: Auto-Federation Mesh                                           â”‚
â”‚  Crossplane: Multi-cloud infra orchestration                                 â”‚
â”‚                                                                              â”‚
â”‚  StorageProvider: S3 / GCS / Azure Blob / OCI                                â”‚
â”‚  IdentityProvider: Azure AD / AWS IAM / Google Identity / Oracle IDCS        â”‚
â”‚  AIProvider: Azure OpenAI / AWS Bedrock / Vertex AI / OCI AI                 â”‚
â”‚                                                                              â”‚
â”‚  Observability: OpenTelemetry, Prometheus, Grafana                           â”‚
â”‚  Security: mTLS, RLS, TLS, PEV encryption                                    â”‚
â”‚  SAIC: Self-Auditing Integrity Core (scheduled jobs)                         â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
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
