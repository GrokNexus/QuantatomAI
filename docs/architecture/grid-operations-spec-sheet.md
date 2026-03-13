# Grid Operations Spec Sheet

## Parity & Beyond
| Operation | Current SLO (p95) | Target SLO (p95) | Notes |
|-----------|-------------------|------------------|-------|
| Grid Open (Hot) | 700ms | 400ms | Cache warming in progress. |
| Grid Open (Cold) | 1.5s | 800ms | Background fetching enabled. |
| Single Cell Edit | 600ms | 300ms | Optimistic UI updates. |
| Medium Spread | 1.2s | 1.0s | Parallel calculation. |
| Heavy Allocation | 10s | 5s | Job sizing optimization. |
| Actuals Load | 10 min | 2 min | Ingestion pipeline tuning. |

# QuantatomAI Grid Operations – Spec Sheet

This document enumerates the core operations of the QuantatomAI Data Grid Engine, compares them to major EPM vendors, and highlights where QuantatomAI goes beyond.

## Legend

- **Parity:** Capability exists in Anaplan / Pigment / PBCS / Jedox.
- **Beyond:** Differentiated or deeper implementation in QuantatomAI.

## Operations table

| # | Operation / Capability | Description | Parity | QuantatomAI Beyond |
|---|------------------------|-------------|--------|--------------------|
| 1 | Dimension-on-any-axis | Any dimension on rows, columns, pages, filters | ✅ | – |
| 2 | Drag-and-drop pivoting | Move dimensions between axes interactively | ✅ | – |
| 3 | Multi-hierarchy support | Primary + alternate hierarchies per dimension | ✅ (partial) | Hierarchies are API-visible, versioned, and composable |
| 4 | Deep hierarchy expand/collapse | Up to 15 levels, ragged hierarchies | ✅ | Atom-aware pruning for performance |
| 5 | Attribute-based filtering | Filter members by attributes | ✅ | Attribute filters compiled into atom query plans |
| 6 | Scenario layering | Side-by-side scenarios | ✅ | Unlimited scenarios with wave QoS prioritization |
| 7 | Versioning | Working/Approved/Snapshot | ✅ | Version as first-class dimension with lineage |
| 8 | Time layouts | Months/quarters/years/rolling windows | ✅ | Rolling windows as reusable templates |
| 9 | Virtualized rendering | 10k–100k visible cells | ✅ | Atom-window aware prefetching |
| 10 | Lazy loading | Fetch visible window + buffer | ✅ | Wave QoS for interactive vs background loads |
| 11 | Drill-down | From aggregate to contributing members | ✅ | Drill path encoded as atom lineage |
| 12 | Drill-through | To transactional detail | ✅ | Connectors + LRI for cross-system drill-through |
| 13 | Linked grids | One grid filters another | ✅ (partial) | Grid queries as composable, shareable objects |
| 14 | Single cell edit | Direct input | ✅ | – |
| 15 | Range edit | Multi-cell edit | ✅ | – |
| 16 | Spreads | Even, proportional | ✅ | Spread compiled into atom jobs with QoS |
| 17 | Driver-based spreads | Driver-based allocation | ✅ | Drivers as explicit, inspectable dependencies |
| 18 | Allocations | Top-down allocations | ✅ | Allocation plans as reusable, versioned artifacts |
| 19 | Copy/paste | Within grid / Excel | ✅ | Paste interpreted as structured atom deltas |
| 20 | Undo/redo | Per-user, per-session | ✅ | Atom-delta based, auditable |
| 21 | Data validation | Rules, warnings | ✅ | Rules can be AI-suggested (AUH/Oracle) |
| 22 | Conditional formatting | Color rules, thresholds | ✅ | AI-suggested formats for anomalies |
| 23 | Cell comments | Notes on cells | ✅ | Comments linked to atom lineage + scenarios |
| 24 | Attachments | Files linked to cells/rows | ✅ (partial) | Attachments tied to atom context, not just grid |
| 25 | Keyboard navigation | Excel-like shortcuts | ✅ | – |
| 26 | Freeze panes | Lock rows/columns | ✅ | – |
| 27 | Sorting/filtering/grouping | Standard grid ops | ✅ | Grouping aware of hierarchies + attributes |
| 28 | Real-time recalculation | Immediate updates | ✅ | Rust SIMD + atom blocks for better scaling |
| 29 | Incremental recalculation | Only affected cells recomputed | ✅ | Atom-level dependency graph, not cell-level |
| 30 | Variance & V% | Across scenarios/versions | ✅ | Variance templates reusable across grids |
| 31 | FX translation | Multi-currency | ✅ | FX rules as separate, inspectable layer |
| 32 | Eliminations | Intercompany, consolidation | ✅ | Elim logic as atom transformations |
| 33 | Custom formulas | User-defined calc logic | ✅ | Formulas compiled to HelioCalc graphs |
| 34 | Offline snapshot | Full grid offline | ❌ (most) | Atom-based snapshot as first-class feature |
| 35 | Offline delta log | Local edit log | ❌ | Atom-delta log with replay |
| 36 | Conflict detection | Overlapping edits | ❌ | ICH-powered conflict detection |
| 37 | Conflict resolution | Guided merge UX | ❌ | AUH + ICH suggestions and explanations |
| 38 | Grid as data source | Reuse in charts/KPIs | ✅ | Grid queries as primary semantic objects |
| 39 | Saved views | Named, shareable grids | ✅ | Views are versioned, API-addressable |
| 40 | AI-suggested layouts | “Best view” suggestions | ❌ | AUH suggests pivots, filters, groupings |
| 41 | AI-suggested KPIs | KPIs based on context | ❌ | Oracle 2.0 proposes KPIs + narratives |
| 42 | ESG-aware compute | Carbon-aware heavy jobs | ❌ | GRB integration for grid-heavy jobs |
| 43 | Privacy-aware views | Masking, encryption | ❌ | PEV-enforced masking at atom level |
| 44 | Multi-cloud consistency | Same grid across clouds | ❌ | AFM + MCGF ensure consistent behavior |
| 45 | Programmable grids | Grids as API objects | ✅ (partial) | Grids as first-class programmable substrate |

## Summary

QuantatomAI matches the core grid capabilities of leading EPM vendors and goes beyond them with:

- Atom-native storage and compute.
- Wave QoS for interactive vs heavy operations.
- Offline-first behavior with intelligent conflict harmonization.
- AI-native UX for layouts, KPIs, and explanations.
- ESG-aware and privacy-first execution.
- Multi-cloud consistency and programmability.
