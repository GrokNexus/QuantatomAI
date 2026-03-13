# QuantatomAI Grid Engine Specification

## Overview
The QuantatomAI Grid Engine is a high-performance, reactive engine designed for multi-dimensional planning and reporting. It manages the lifecycle of grids, from query planning to writeback.

## Core Components
- **Core:** The heart of the engine, managing state and coordination.
- **Query:** Optimized query planning and execution against the Heat/Warm/Cold stores.
- **Writeback:** Efficiently handling cell updates and propagating changes.
- **Dimensions:** Managing multi-dimensional structures and hierarchies.
- **Offline:** Support for local operations and synchronization.
- **Utils:** Common utilities for calculation and data transformation.

# QuantatomAI Data Grid Engine

The QuantatomAI Data Grid Engine is the primary interaction surface for planning and reporting. It is designed to match and surpass the grid capabilities of Anaplan, Pigment, Oracle PBCS, Adaptive, Jedox, and other EPM tools, while being deeply integrated with the atom lattice and HelioCalc.

## Design principles

- **Multidimensional by default:** Any dimension on any axis, any combination.
- **Atom-native:** Every cell is backed by atoms, not cube blocks or spreadsheet cells.
- **Wave-aware:** Interactive grid operations are first-class waves with QoS.
- **AI-native:** The grid is guided and explained by AUH and Oracle 2.0.
- **Offline-resilient:** Full offline editing with intelligent conflict harmonization.
- **Programmable:** Grids are API-addressable, versioned, and composable.

## Core capabilities

### Multidimensional layout

- Dimension-on-any-axis (rows, columns, pages, filters).
- Drag-and-drop pivoting between axes.
- Multiple hierarchies per dimension (primary + alternates).
- Deep hierarchy expand/collapse (up to 15 levels, ragged supported).
- Attribute-based filtering (e.g., Region, Product Type).
- Scenario, version, and time layering (side-by-side views).

### Data retrieval and navigation

- Virtualized rendering (10k–100k visible cells).
- Lazy loading with window-based prefetching.
- Drill-down from aggregates to contributing members.
- Drill-through to transactional detail via connectors.
- Linked grids where one grid’s selection filters another.

### Planning operations

- Single and multi-cell edits.
- Spreads: even, proportional, driver-based.
- Allocations: top-down, rule-based.
- Copy/paste within grid, across grids, and to/from Excel.
- Undo/redo with atom-delta based history.
- Data validation with hard/soft rules and warnings.

### Calculations

- Real-time recalculation for local changes.
- Incremental recalculation on atom blocks.
- Variance and V% across scenarios/versions.
- FX translation and eliminations.
- Custom formulas compiled into HelioCalc graphs.

### UX and collaboration

- Freeze panes, sorting, filtering, grouping.
- Conditional formatting.
- Cell comments and attachments.
- Keyboard-first navigation (Excel muscle memory).

### Offline and conflicts

- Offline snapshot of grid data (atom-based).
- Local delta log of edits (atom deltas).
- Conflict detection for overlapping edits.
- ICH-guided conflict resolution with AUH explanations.

### Reuse in VizBoards

- Grids as data sources for charts, KPIs, and tables.
- Saved views as named, shareable grid definitions.
- VizBoards bind to grid queries, not raw tables.

## Integration with the architecture

- **UI Layer:** `ui/web/grid-engine` implements the grid renderer, query builder, writeback, offline, and VizBoard binding.
- **Grid Service:** `services/grid-service` implements the Grid Query Planner, atom retrieval, HelioCalc orchestration, and writeback handling.
- **Data Layer:** `data/atoms`, `data/hot-store`, `data/warm-store`, `data/cold-store`, and `data/metadata` provide the atom lattice and hierarchies.
- **Compute Layer:** `compute/heliocalc` executes aggregations, variances, FX, eliminations, and custom formulas on atom blocks.
- **Eventing Layer:** `eventing/aodl` and `eventing/wrm` stream atom deltas and ensure resilient grid updates.
- **AI Layer:** `ai/auh`, `ai/ich`, and `ai/oracle` provide AI-guided layouts, conflict resolution, and narratives.

## Moat characteristics

- Atom-native lattice behind the grid.
- Rust SIMD compute (HelioCalc) for scaling.
- Wave QoS for interactive vs heavy operations.
- Offline + ICH conflict harmonization.
- AI-native UX (AUH + Oracle 2.0).
- ESG-aware and privacy-first (GRB + PEV).
- Multi-cloud consistent behavior (AFM + MCGF).
- Grids as programmable, versioned, API objects.
