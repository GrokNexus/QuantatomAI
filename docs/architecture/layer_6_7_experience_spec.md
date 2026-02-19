# QuantatomAI Layer 6 & 7 Specification: The Holographic Experience

## Layer 6: Domain Orchestration (The "Traffic Controller")
This layer translates business intent (pivot, filter, consolidate) into technical queries.

### 6.1 The Orchestrator (Go)
*   **Technology:** **Go 1.22+**.
*   **Framework:** **Connect-Go** (High-performance gRPC).
*   **Role:**
    1.  Receives user request: "Pivot by Region."
    2.  Resolves Permissions (ACL).
    3.  Dispatches calc jobs to Layer 5 (Rust).
    4.  Consolidates results and streams to UI.

### 6.2 The Dependency Injector
*   **Technology:** **Fx** (Uber's DI framework).
*   **Why:** Modularizes "Auth Service", "Metadata Service", and "Grid Service" for testability.

### 6.3 Observability
*   **Technology:** **OpenTelemetry**.
*   **Tracing:** 100% Sampling for end-to-end request tracing.

---

## Layer 7: The Experience Layer (The "Projector")
This layer is the **Grid UI**. It is a dumb projector of the data lattice.

### 7.1 The Rendering Core (WebGPU)
*   **Technology:** **WebGPU** (via `wgpu` or `Three.js`).
*   **Why:** Traditional DOM (HTML Table) lags at >5k visible cells. GPU renders 10M cells at 120 FPS.
*   **Canvas:** We draw the grid lines, text, and conditional formatting as **Textures** on a quad.

### 7.2 The State Atom (Zustand)
*   **Technology:** **Zustand** + **Immer**.
*   **Philosophy:** Application state is immutable.
*   **Sync:** Using `persist` middleware to sync state to IndexedDB for offline support.

### 7.3 The Data Stream (Connect-Web)
*   **Technology:** **Connect-Web** (gRPC-Web).
*   **Format:** **Binary Protobuf**. No JSON.
*   **Streaming:** The grid loads incrementally. You see the top-left corner instantly while the rest streams in.

### 7.4 The Editor (Monaco)
*   **Technology:** **Monaco Editor**.
*   **Role:** The "Formula Bar" for editing AtomScript.
*   **Integration:** Custom Language Server Protocol (LSP) for IntelliSense on dimension members.
