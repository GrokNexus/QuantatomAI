# 🧠 QuantatomAI Formula Intelligence & Grid Integration

**Status:** 💎 **Diamond Specification**
**Layer:** 7.3 (Editor) <-> 5.0 (Compute Kernel)

This document answers: *How does the User's Formula connect to the Iron Engine?*

---

## 1. The AtomScript Function Library
The Monaco Editor will support a curated library of high-performance financial functions. These are not just text; they are JIT-compiled into SIMD machine code.

### A. Aggregation (The "Big Iron")
These run on the **LatticeArena** using AVX-512 instructions.
*   `SUM( [Dimension] )`: Aggregates values along a hierarchy.
*   `AVG( [Dimension] )`: Weighted averages.
*   `MIN / MAX`: Peak detection.
*   `COUNT`: Sparse density analysis.

### B. Time Intelligence (The "Quantum Leap")
These leverage the **Dimension Awareness** of the metadata layer.
*   `PREV( [Time] )`: Previous period value (t-1).
*   `NEXT`: Forecast/Budget chaining.
*   `YTD`: Year-to-Date accumulation.
*   `CAGR`: Compound Annual Growth Rate over a range.

### C. Logic & Control
*   `IF( condition, true, false )`: Branchless select (no CPU branch misprediction).
*   `SWITCH`: Multi-path routing.
*   `ISBLANK / ISERROR`: Data quality guards.

### E. Hierarchy Intelligence (The "Family Tree")
You asked: *Where are @Children and @Descendants?*
These are **Meta-Functions** that resolve against the Layer 2.1 Metadata Graph before calculation begins.

*   `@Children( [Dimension], [Member] )`: Returns the immediate children.
    *   *Example:* `@Children([Region], [North America])` -> `[USA, Canada, Mexico]`
*   `@Descendants( ... )`: Returns the full recursive subtree (Closure Table).
*   `@Parent( ... )`: Returns the immediate parent.
*   `@Siblings( ... )`: Returns members at the same level.
*   `@Level( [Dimension], 0 )`: Returns all members at the root level.

**Implementation Plan:**
1.  **Parser:** Detects `@` token.
2.  **Resolver:** Queries the `dimension_hierarchy` table (Layer 2.1).
3.  **Expansion:** The Compiler "unrolls" the function into a standard aggregation.
    *   *Input:* `SUM(@Children([Region], [North America]))`
    *   *Compiled:* `SUM([USA]) + SUM([Canada]) + SUM([Mexico])` (Zero Runtime Recursion).

### F. Hyper-Fast Lookups (The "O(1) Moat")
You asked: *What about VLOOKUP and XLOOKUP?*
In Excel, `VLOOKUP` is **O(N)** (Scanning). In QuantatomAI, it is **O(1)** (Hashing).

*   `LOOKUP( [TargetDimension], [TargetValue] )`: Fetches data from a peer cell.
    *   *Mechanism:* The engine takes the current `CoordinateHash`, replaces the 16 bytes of the `TargetDimension` with the hash of `TargetValue`, and reads memory directly.
    *   *Performance:* 1 nanosecond (Pointer jump) vs 1 millisecond (Table scan).
*   `XLOOKUP`: Supported as `LOOKUP` with a fallback value.

### G. "Moat" Innovations
1.  **The "Time-Travel" Operator (`->`)**:
    *   `[Revenue] -> [Prev Year]`
    *   This is syntactic sugar for a coordinate shift optimization.
2.  **The "Context" Function**:
    *   `@Drive( [Driver] )`: Automatically finds the driver for the current intersection based on metadata attributes (e.g., finding the correct Tax Rate for a specific Region without large IF statements).
How does `[Region]` become a memory address?

1.  **The User Types:** `Sum([Revenue])` in Monaco.
2.  **LSP (Language Server):**
    *   Queries **Layer 2.1 (Metadata Service)**.
    *   Confirms `Revenue` is a valid member of the `Account` dimension.
    *   *Moat Feature:* Auto-suggests `[Revenue - US]`, `[Revenue - EU]` based on user context.
3.  **The Map:**
    *   The **Orchestrator** retrieves the `CoordinateHash` for `Revenue`.
    *   It passes this Hash to the **Rust Engine**.

---

## 3. Utilizing the Data Grid Engine (Layer 5 & 6)
You asked: *What is the plan to use the data grid engine we built?*

### The Execution Pipeline "Ferrari"
When the user hits **Enter**, this pipeline fires in < 50ms:

1.  **Parse (Layer 5.2):**
    *   The `AtomScript` string is tokenized by `Logos`.
    *   Converted to an AST (Abstract Syntax Tree).

2.  **Dependency Resolution (Layer 5.3):**
    *   The **Graph Resolver** checks: *Does Revenue depend on anything else?*
    *   It sorts the calculation order (DAG).

3.  **JIT Compilation (Layer 5.2):**
    *   The AST is compiled into **Stack Bytecode** (`OP_PUSH`, `OP_ADD`).
    *   *Optimization:* Constant Folding (`1+2 -> 3`).

4.  **Massive Parallel Execution (Layer 5.1):**
    *   The **LatticeArena** (Sharded) is engaged.
    *   **SIMD:** The engine loads 8 doubles (64-bit floats) at a time into CPU registers.
    *   It executes the logic across 1 billion cells in parallel shards.

5.  **Zero-Copy Streaming (Layer 6):**
    *   The results are written to an **Arrow RecordBatch**.
    *   Streamed via gRPC to the Frontend.
    *   **WebGPU (Layer 6.2)** renders the new numbers instantly.

### Summary
The **Monaco Editor** is the steering wheel. The **AtomEngine** is the V12 Motor. They are connected by the **Metadata Drive Shaft** (Layer 2.1).
