# 📖 Layer 7.4 Implementation Guide: Hierarchy Intelligence

**Status:** ✅ Implemented (Parser/Compiler)
**Location:**`services/atom-engine/src/lattice/metadata.rs`, `src/atom_script/compiler.rs`
**Technology:** Rust (Traits, AST Transformation)

---

## 1. Overview
Hierarchy Intelligence allows the Calculation Engine to understand the semantic relationships between data points (e.g., "North America contains USA, Canada, Mexico").
Unlike traditional engines that resolve this at runtime (slow recursive queries), QuantatomAI resolves this at **Compile Time**.

## 2. Architecture: The Metadata Bridge
We introduced a semantic bridge between the **Compute Kernel** (which knows math) and the **Metadata Store** (which knows relationships).

### A. The `HierarchyResolver` Trait
Located in `src/lattice/metadata.rs`.
```rust
pub trait HierarchyResolver {
    fn get_children(&self, dimension: &str, member: &str) -> Vec<String>;
    fn get_parent(&self, dimension: &str, member: &str) -> Option<String>;
}
```
*   **Design Choice:** By using a Trait, we decouple the engine from the actual storage (Postgres/Redis).
*   **Mock Implementation:** We use `MockHierarchyResolver` for unit testing the compiler without a DB connection.

## 3. Implementation Details

### B. Parser Upgrades (`parser.rs`)
*   **New Token:** Added `@` support in `lexer.rs`.
*   **AST Node:** Added `Expr::HierarchyCall { name, args }`.
*   **Logic:** The parser distinguishes between standard functions (`SUM`) and hierarchy macros (`@Children`).

### C. Compile-Time Expansion (`compiler.rs`)
This is the **Moat Innovation**.
Instead of emitting a loop, the compiler "unrolls" the hierarchy.

*   **Input:** `SUM(@Children([Region], [North America]))` (where N.A. has 3 children)
*   **Process:**
    1.  Compiler calls `resolver.get_children("Region", "North America")`.
    2.  Returns `["USA", "Canada", "Mexico"]`.
    3.  Compiler emits load instructions for each child.
    4.  Compiler emits `OpCode::Sum(3)`.
*   **Output Bytecode:**
    ```
    LOAD_CONST "USA"
    LOAD_CONST "Canada"
    LOAD_CONST "Mexico"
    SUM 3
    ```
*   **Performance:** 0 overhead at runtime. The VM just sees a summation of 3 numbers.

## 4. Verification
*   **Unit Tests:** `src/atom_script/tests.rs` verifies that `SUM(@Children)` produces the correct OpCodes.
