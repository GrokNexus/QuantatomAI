# QuantatomAI Layer 5 Specification: The AtomEngine Kernel

## Layer 5: The Computation Engine (The "Brain")
This layer performs all mathematical calculations, dependency resolution, and aggregation.

### 5.1 The Rust Core (AtomEngine)
*   **Technology:** **Rust (Nightly)**.
*   **Why Rust?** Manual memory management (No GC pauses) and explicit SIMD control.
*   **Parallelism:** **Rayon**. Uses work-stealing to process "North America" and "Europe" on different CPU cores simultaneously.

### 5.2 The JIT Compiler (AtomScript Runtime)
*   **Technology:** **LLVM (Inkwell Crate)**.
*   **Process:**
    1.  User writes `Sum(Revenue)`.
    2.  Parser emits AST.
    3.  Compiler emits **AVX-512 Machine Code** at runtime.
    4.  **Result:** User formula runs at native speed.

### 5.3 The Graph Resolver (Neo4j Integration)
*   **Technology:** **Neo4j 5.x**.
*   **Role:** Resolves complex dependencies ("Net Income" depends on "Tax" depends on "Revenue").
*   **Optimization:** We replicate the *Hot Path* of the dependency graph into a **Rust `Petgraph` structure** in RAM for sub-microsecond traversal.

### 5.4 The Vector Engine (SIMD)
*   **Library:** **`portable-simd`** (Rust std).
*   **Operation:** Loads 8 floating-point numbers into a single CPU register (`_mm512_add_pd`) and adds them in one clock cycle.
*   **Performance:** 8-16x faster than scalar Java/C# loops.

---

## Why This Beats Anaplan
| Feature | Anaplan (Java) | QuantatomAI (Rust) |
| :--- | :--- | :--- |
| **Math** | Scalar Loops | SIMD Vectors |
| **Memory** | GC Pauses (Stop-the-world) | Arena Allocation (Zero Pause) |
| **Formula** | Interpreted Script | JIT Compiled Assembly |
