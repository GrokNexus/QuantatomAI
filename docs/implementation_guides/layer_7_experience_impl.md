# 📖 Layer 7 Implementation Guide: The Formula Experience

**Status:** ✅ Implemented Skeleton
**Location:**`frontend/src/components/FormulaBar.tsx`
**Technology:** Monaco Editor + React

---

## 1. The Formula Bar (Layer 7.3)
We use the **Monaco Editor** (the core of VS Code) to provide a rich code editing experience for AtomScript.
*   **Why Monaco?**
    *   Syntax Highlighting (Keywords like `Sum`, `Avg` in Neon Red).
    *   IntelliSense (Autocomplete for Dimensions `[Region]`).
    *   Minimap/Folding (Disabled for the single-line bar).
*   **Theme:** `quantatom-dark` (Matches the application aesthetic).

## 2. Hierarchy Intelligence (Layer 7.4)
For details on `@Children`, `@Descendants` and the Compile-Time Expansion engine, see:
👉 [Layer 7.4 Implementation Guide](layer_7_4_hierarchy_impl.md)

## 2. Integration
The Formula Bar is mounted at the top of `App.tsx` with `z-index: 100`.
*   It controls the `formula` state.
*   In the future, pressing `Enter` will dispatch a `QueryGridRequest` to the Orchestrator.

## 3. Future Work (Layer 7.5: Visual Intelligence)
*   **Charts:** We will use `Recharts` or `D3` (via WebGPU) to render charts overlaid on the grid.
*   **Context Menu:** Right-click to pivot/filter.

## 4. Usage
Type a formula: `Sum([Revenue])`
(Currently just updates local state, does not trigger calc yet).
