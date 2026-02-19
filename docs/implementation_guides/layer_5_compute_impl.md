# 📖 Layer 5 Implementation Guide: The AtomEngine Kernel (Compute)

**Status:** ✅ Implemented
**Location:** `services/atom-engine/src/`
**Technology:** Rust (Rayon + Logos + SIMD)

---

## 1. The "Why": Speed of Light Math
In standard planning tools (Anaplan/TM1), calculations are slow because:
1.  **Objects:** Every cell is a Java Object (pointer chasing).
2.  **Interpretation:** Formulas are interpreted at runtime (slow dispatch).
**The Solution:** We built the **AtomEngine**.
*   **LatticeArena:** Stores 1 billion numbers in a flat `Vec<f64>`. CPU Prefetcher loves this.
*   **AtomScript VM:** Compiles `Sum(Revenue)` into efficient Bytecode (`OP_ADD`, `OP_RETURN`) that runs in the L1 Cache.

## 2. The Implementation (Code Deep Dive)

### A. The Core (`lattice/arena.rs`)
*   **Structure of Arrays (SoA):** We separate values (`Vec<f64>`) from metadata.
*   **`parking_lot`:** We use ultra-fast spinlocks instead of OS mutexes for high-concurrency cell updates.

### B. The Vector Engine (`compute/simd.rs`)
*   **Rayon:** We use `par_iter()` to split math ops across all CPU cores.
*   **Auto-Vectorization:** The loops are designed so LLVM (Rustc) emits AVX-512 instructions automatically.

### C. AtomScript (`atom_script/`)
*   **Lexer (`logos`):** Scans formula strings at GB/s speeds.
*   **Parser (Pratt):** Handles precedence (`*` before `+`) correctly.
*   **VM (Stack Machine):** Executes the compiled chunk. It fits entirely in CPU registers.

### D. The Graph Resolver (`compute/graph.rs`)
*   **Petgraph:** We model dependencies (e.g., `Net Income` -> `Revenue`) as a DAG.
*   **Topological Sort:** Determines the correct execution order. If `Net Income` depends on `Revenue`, we calculate `Revenue` first.
*   **Cycle Detection:** Automatically detects circular references (e.g., A -> B -> A) and returns an error.

## 3. Usage Instructions

### Running a Calculation
```rust
// 1. Setup Arena
let arena = LatticeArena::new(1000);
let idx_a = arena.set_cell(hash("Revenue"), 100.0);
let idx_b = arena.set_cell(hash("Cost"), 60.0);

// 2. Compile Formula "Revenue - Cost"
let mut compiler = Compiler::new();
compiler.compile(&Expr::Binary { ... });

// 3. Execute
let mut vm = VM::new(compiler.chunk);
let result = vm.run(); // Returns 40.0
```

## 4. Key Design Decisions related to "The Moat"
1.  **Zero GC:** Unlike Java, we never pause. The Arena grows only when needed and drops all at once when the View is closed.
2.  **Thread Safety:** The `LatticeArena` is safe for 1000 concurrent writers (e.g., streaming market data updates) while the Calc Engine reads from it.
