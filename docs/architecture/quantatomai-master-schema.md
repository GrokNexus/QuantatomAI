# QuantatomAI Master Schema & Resonance Specification

================================================================================
                 MASTER SCHEMA FOR QUANTATOMAI (ALL 7 LAYERS)
================================================================================

Legend:
- [Table/Node]: Entity (e.g., Table in Postgres, Node in Neo4j)
- ────> : One-to-Many / Foreign Key / Relationship
- <──── : Sync / Event Flow (e.g., Kafka Trigger)
- ~~> : Asynchronous Wave / Propagation (Resonance-specific)
- (Index): Performance Hint (e.g., GIN for JSONB)
- {Dynamic}: Expandable Field (JSONB/Map for flexibility)
- <E3 Tier>: Cost Tier Tag (hot/warm/cold from architecture)

## Core Resonance Innovations

- **Resonance Aggregation Bridges (RAB)**: Particles that auto-roll granular Actuals (10k accounts) to high-level Plan members (500 accounts) with formula inheritance.
- **Echo Commentary Layers (ECL)**: Attached commentary particles (notes/adjusts) entangled to any level of the lattice.
- **Bidirectional Resonance Flows (BRF)**: Top-down targets as "guidance echoes" propagating down; bottoms-up actuals flow up as variance echoes.
- **Shard-by-Planning-Type (SPT)**: Auto-sharding the lattice by application type (e.g., `sales_shard`, `finance_shard`) for linear scalability.

---

### [Layer 1 & 2: Infrastructure & Data - Postgres]
*6 Tables: Transactional, Master Metadata, and Small-Scale Real-Time Atoms.*

  +-------------------+
  | viz_apps (1)      |
  | - id (UUID PK)    |
  | - name (VARCHAR)  |
  | - planning_type (ENUM: 'corporate', 'supply_chain', 'headcount') | ← SPT Shard Key
  | - created_at (TIMESTAMP) |
  | (Index: planning_type) |
  +-------------------+
             |
             v
  +-------------------+
  | dimensions (2)    |
  | - id (UUID PK)    |
  | - viz_app_id (FK) |
  | - name (VARCHAR)  |
  | - type (ENUM)     |
  | - properties (JSONB) | ← {formula: 'sum(children)', alignment: 8}
  | - huc_gates (JSONB)  | ← Governance/Access Control
  | (GIN Index on properties) |
  +-------------------+
             |
             v
  +-------------------+
  | dimension_members (3) |
  | - id (UUID PK)       |
  | - dimension_id (FK)  |
  | - parent_id (FK self)| ← Hierarchy Support
  | - code (VARCHAR)     |
  | - name (VARCHAR)     |
  | - attributes (JSONB) | ← Custom weights/properties
  | (Index: dimension_id, parent_id) |
  +-------------------+
             |
             v
  +-------------------+
  | data_atoms (4)    |
  | - id (UUID PK)    |
  | - value (NUMERIC/NaN-Boxed) |
  | - scenario_id (UUID)|
  | - coordinates (JSONB) | ← {dim_id: member_id}
  | - bridge_id (UUID FK) | ← RAB for roll-up paths
  | - target_id (UUID FK) | ← BRF for top-down targeting
  | -- MOAT ENGINEERING -- |
  | - causal_clock (BIGINT[]) | ← Lamport Vector for distributed consistency
  | - bridge_vector (BYTEA)   | ← Roaring Bitmap for SIMD propagation
  | - security_mask (BIGINT)  | ← CPU-level ACL bitmask
  | (GIN Index on coordinates) |
  +-------------------+
             |
             v
  +-------------------+
  | aggregation_bridges (5) | ← RAB (Actuals to Plan Bridge)
  | - id (UUID PK)         |
  | - source_atom_id (FK)  | ← Granular (Actuals)
  | - target_atom_id (FK)  | ← High-level (Plan)
  | - weight (FLOAT)       | ← Proportional allocation
  | - variance_echo (JSONB)| ← Delta analysis
  +-------------------+
             |
             v
  +-------------------+
  | commentary_echo (6)     | ← ECL (Adjustment/Notes)
  | - id (UUID PK)         |
  | - atom_id (FK)         | ← Link to any lattice coord
  | - adjust_value (NUMERIC)| 
  | - note (TEXT)          |
  | - created_by (UUID)    |
  +-------------------+

---

### [Layer 3 & 4: Data & Lattice - ClickHouse]
*3 Tables + Views: Columnar Analytics for Massive Datasets (Synchronized via Kafka).*

  +-------------------+
  | atom_analytics (1)| ← Flattened Lattice Projection
  | - id (UUID)       |
  | - value (Float64) |
  | - scenario_id (UUID)|
  | - dim_map (Map)   | ← Dynamic Dimensions
  | - bridge_id (UUID)| ← RAB
  | - target_id (UUID)| ← BRF
  | (MinMax Index on scenario_id) |
  +-------------------+

  +-------------------+
  | agg_cells (2)     | ← Pre-computed Aggregates (Layer 5 Fuel)
  | - node_id (UUID)  |
  | - sum_val (AggregateFunction) |
  | - avg_val (AggregateFunction) |
  | (AggregatingMergeTree Engine) |
  +-------------------+

---

### [Layer 5: Compute & AI - Neo4j]
*Graph Database for Entanglement & Dependency Resolution.*

  [Node: :Member] ──── (:PARENT_OF) ────> [Node: :Member]
  [Node: :Atom] ────── (:BELONGS_TO) ─────> [Node: :Member]
  [Node: :Atom] ────── (:BRIDGE) ─────────> [Node: :Atom] ← RAB Path
  [Node: :Atom] ────── (:TARGET_FLOW) ────> [Node: :Atom] ← BRF Path

---

## Data Flows Across the 7-Layer Mesh

1.  **Write Path (Input Ceremony)**:
    - User (L7) ────> Grid API (L6) ────> Write to Postgres (L2)
    - Postgres Trigger ────> Kafka (L3) ────> Neo4j Entanglement (L5)
    - Neo4j ────> Trigger Recalculation (L5 Rust Engine)

## 4. Realization Stack: The "Absolute Best" Toolkit

To realize this architecture with **"Moat-Grade" Performance**, we reject generic "Enterprise Java" stacks. We use a **Bilingual High-Performance Core (Rust + Go)** wrapped in a **Tensor-Based AI Lattice**.

### Layer 7: Holographic Experiences (The "Projector")
*   **Language:** **TypeScript 5.x** (Strict Mode)
*   **Rendering Core:** **WebGPU** (via `wgpu` or `Three.js`) for the Canvas Grid. DOM is only for chrome/menus.
*   **State Management:** **Zustand** (Atomic state) + **TanStack Query** (Async).
*   **Transport:** **Connect-Web** (gRPC-Web) with **Binary Protobufs** (no JSON on the wire).

### Layer 6: Domain Orchestration (The "Traffic Controller")
*   **Language:** **Go 1.22+** (PGO Enabled)
*   **Framework:** **Connect-Go** (High-perf gRPC) + **Fx** (Dependency Injection).
*   **Concurrency:** Heavy use of **Goroutines** for fan-out/fan-in query planning.
*   **Observability:** **OpenTelemetry** (Auto-instrumentation) with 100% sampling on traces.

### Layer 5: The AtomEngine Kernel (The "Brain")
*   **Language:** **Rust 1.75+** (Nightly for AVX-512 features).
*   **SIMD Framework:** **`portable-simd`** or **`aggregates`** crate for explicit vectorization.
*   **Parallelism:** **Rayon** for work-stealing parallelism across lattice slices.
*   **Graph Logic:** **Petgraph** (Rust) for in-memory dependency resolution (beating Neo4j for hot-path calcs).
*   **JIT Compiler:** **Inkwell** (LLVM wrapper) to compile user formulas into machine code at runtime.

### Layer 4: The Lattice Spine (Data Sovereignty)
*   **Warm Store:** **ClickHouse** (MergeTree engine) for columnar scans.
*   **Hot Store:** **Redis 7.2** (Cluster Mode) or **DragonflyDB** (Thread-per-core architecture).
*   **Metadata:** **Postgres 16** with **pgvector** installed.
*   **IPC:** **Apache Arrow Flight** for zero-copy data movement between Rust Engine and ClickHouse.

### Layer 3: Resonance Eventing (The "Nervous System")
*   **Log Backbone:** **Redpanda** (C++ Kafka API) for 10x lower latency than JVM Kafka.
*   **Stream Processing:** **Rust** consumers (via `rdkafka`) for sub-ms event handling.
*   **Format:** **FlatBuffers v2** (Zero-parsing overhead) for all internal events.

### Layer 1: Infrastructure (The "Bedrock")
*   **Orchestration:** **Kubernetes 1.29** (Cilium CNI for eBPF networking).
*   **Control Plane:** **Crossplane** (Go) to manage cloud resources as code.
*   **Mesh:** **Istio Ambient Mesh** (No sidecar overhead).

---

## 5. Why This Stack Beats Anaplan/Pigment?

| Component | Competitor Stack (Legacy) | QuantatomAI Stack (Modern) | The "Moat" Result |
| :--- | :--- | :--- | :--- |
| **Compute** | Java / proprietary scripts | **Rust + LLVM + SIMD** | **100x Faster** math using CPU vector instructions. No GC pauses. |
| **Grid UI** | HTML DOM Tables | **WebGPU Canvas** | **120 FPS** scrolling on 10M cell grids. No DOM lag. |
| **Messaging** | RabbitMQ / JMS | **Redpanda + FlatBuffers** | **Zero-Copy** event propagation. Sub-ms variance updates. |
| **Storage** | Proprietary Blob / Cube | **ClickHouse + Arrow** | **Open Columnar** format. Instant analysis of 100B rows. |
