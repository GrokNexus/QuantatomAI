> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

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

# QuantatomAI Grid Operations â€“ Spec Sheet

This document enumerates the core operations of the QuantatomAI Data Grid Engine, compares them to major EPM vendors, and highlights where QuantatomAI goes beyond.

## Legend

- **Parity:** Capability exists in Anaplan / Pigment / PBCS / Jedox.
- **Beyond:** Differentiated or deeper implementation in QuantatomAI.

## Operations table

| # | Operation / Capability | Description | Parity | QuantatomAI Beyond |
|---|------------------------|-------------|--------|--------------------|
| 1 | Dimension-on-any-axis | Any dimension on rows, columns, pages, filters | âœ… | â€“ |
| 2 | Drag-and-drop pivoting | Move dimensions between axes interactively | âœ… | â€“ |
| 3 | Multi-hierarchy support | Primary + alternate hierarchies per dimension | âœ… (partial) | Hierarchies are API-visible, versioned, and composable |
| 4 | Deep hierarchy expand/collapse | Up to 15 levels, ragged hierarchies | âœ… | Atom-aware pruning for performance |
| 5 | Attribute-based filtering | Filter members by attributes | âœ… | Attribute filters compiled into atom query plans |
| 6 | Scenario layering | Side-by-side scenarios | âœ… | Unlimited scenarios with wave QoS prioritization |
| 7 | Versioning | Working/Approved/Snapshot | âœ… | Version as first-class dimension with lineage |
| 8 | Time layouts | Months/quarters/years/rolling windows | âœ… | Rolling windows as reusable templates |
| 9 | Virtualized rendering | 10kâ€“100k visible cells | âœ… | Atom-window aware prefetching |
| 10 | Lazy loading | Fetch visible window + buffer | âœ… | Wave QoS for interactive vs background loads |
| 11 | Drill-down | From aggregate to contributing members | âœ… | Drill path encoded as atom lineage |
| 12 | Drill-through | To transactional detail | âœ… | Connectors + LRI for cross-system drill-through |
| 13 | Linked grids | One grid filters another | âœ… (partial) | Grid queries as composable, shareable objects |
| 14 | Single cell edit | Direct input | âœ… | â€“ |
| 15 | Range edit | Multi-cell edit | âœ… | â€“ |
| 16 | Spreads | Even, proportional | âœ… | Spread compiled into atom jobs with QoS |
| 17 | Driver-based spreads | Driver-based allocation | âœ… | Drivers as explicit, inspectable dependencies |
| 18 | Allocations | Top-down allocations | âœ… | Allocation plans as reusable, versioned artifacts |
| 19 | Copy/paste | Within grid / Excel | âœ… | Paste interpreted as structured atom deltas |
| 20 | Undo/redo | Per-user, per-session | âœ… | Atom-delta based, auditable |
| 21 | Data validation | Rules, warnings | âœ… | Rules can be AI-suggested (AUH/Oracle) |
| 22 | Conditional formatting | Color rules, thresholds | âœ… | AI-suggested formats for anomalies |
| 23 | Cell comments | Notes on cells | âœ… | Comments linked to atom lineage + scenarios |
| 24 | Attachments | Files linked to cells/rows | âœ… (partial) | Attachments tied to atom context, not just grid |
| 25 | Keyboard navigation | Excel-like shortcuts | âœ… | â€“ |
| 26 | Freeze panes | Lock rows/columns | âœ… | â€“ |
| 27 | Sorting/filtering/grouping | Standard grid ops | âœ… | Grouping aware of hierarchies + attributes |
| 28 | Real-time recalculation | Immediate updates | âœ… | Rust SIMD + atom blocks for better scaling |
| 29 | Incremental recalculation | Only affected cells recomputed | âœ… | Atom-level dependency graph, not cell-level |
| 30 | Variance & V% | Across scenarios/versions | âœ… | Variance templates reusable across grids |
| 31 | FX translation | Multi-currency | âœ… | FX rules as separate, inspectable layer |
| 32 | Eliminations | Intercompany, consolidation | âœ… | Elim logic as atom transformations |
| 33 | Custom formulas | User-defined calc logic | âœ… | Formulas compiled to HelioCalc graphs |
| 34 | Offline snapshot | Full grid offline | âŒ (most) | Atom-based snapshot as first-class feature |
| 35 | Offline delta log | Local edit log | âŒ | Atom-delta log with replay |
| 36 | Conflict detection | Overlapping edits | âŒ | ICH-powered conflict detection |
| 37 | Conflict resolution | Guided merge UX | âŒ | AUH + ICH suggestions and explanations |
| 38 | Grid as data source | Reuse in charts/KPIs | âœ… | Grid queries as primary semantic objects |
| 39 | Saved views | Named, shareable grids | âœ… | Views are versioned, API-addressable |
| 40 | AI-suggested layouts | â€œBest viewâ€ suggestions | âŒ | AUH suggests pivots, filters, groupings |
| 41 | AI-suggested KPIs | KPIs based on context | âŒ | Oracle 2.0 proposes KPIs + narratives |
| 42 | ESG-aware compute | Carbon-aware heavy jobs | âŒ | GRB integration for grid-heavy jobs |
| 43 | Privacy-aware views | Masking, encryption | âŒ | PEV-enforced masking at atom level |
| 44 | Multi-cloud consistency | Same grid across clouds | âŒ | AFM + MCGF ensure consistent behavior |
| 45 | Programmable grids | Grids as API objects | âœ… (partial) | Grids as first-class programmable substrate |

## Summary

QuantatomAI matches the core grid capabilities of leading EPM vendors and goes beyond them with:

- Atom-native storage and compute.
- Wave QoS for interactive vs heavy operations.
- Offline-first behavior with intelligent conflict harmonization.
- AI-native UX for layouts, KPIs, and explanations.
- ESG-aware and privacy-first execution.
- Multi-cloud consistency and programmability.
