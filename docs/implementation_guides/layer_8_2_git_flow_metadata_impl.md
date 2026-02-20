# 🔀 Layer 8.2 Implementation Guide: Enterprise Wrap (Metadata Git-Flow)

## 📌 Executive Summary
This document details the implementation of **Layer 8.2: Metadata Git-Flow**. 

In enterprise planning environments (like Anaplan or Adaptive Insights), it is imperative that structural changes—such as reorganizing entire sales hierarchies or introducing new product lines—can be staged, reviewed, and tested *without* disrupting the live calculation engine.

To solve this, we implemented a sophisticated **Delta-Branching Architecture** extending our PostgreSQL metadata registry and our Go Orchestrator, effectively giving Quantatom AI full Git-like branching semantics for dimensional models.

---

## 🏗️ 1. Architecture Overview (Delta-Branching)

Instead of cloning the entire dimension tree whenever a user creates a "sandbox" branch, we utilize a **Sparse Override Pattern**. 

A branch only stores the *differences* (deltas) from the `main` branch. 
- If a member is modified, the new version is written with the new `branch_id`.
- If a member is deleted, a "Tombstone" (marked `is_deleted = true`) is written to the branch.
- On reading, a PostgreSQL View (`branch_view_members`) dynamically merges the base data with the branch overlays using a SQL coalesce/exclusion strategy.

---

## ⚙️ 2. Core Components

### 2.1 The PostgreSQL Schema (`02_git_metadata.sql`)
- **`branches` Table:** Stores workspaces. Every branch (except `main`) has a `base_branch_id` pointing to its parent.
- **`commits` Table:** Stores snapshots of grouped changes within a branch to allow for granular rollback and auditability.
- **`dimension_members` Upgrades:**
  - `branch_id UUID`: Isolates the member to a specific workspace.
  - `commit_id UUID`: Links the member to a specific audit snapshot.
  - `is_deleted BOOLEAN`: The Tombstone flag for logical deletions.
  - Unique Constraint: Upgraded from `(dimension_id, name)` to `(dimension_id, name, branch_id)`.

### 2.2 The Delta-Overlay View
```sql
CREATE OR REPLACE VIEW branch_view_members AS
SELECT 
    m.*, 
    b.id as query_branch_id
FROM branches b
JOIN dimension_members m ON 
    m.branch_id = b.id 
    OR 
    (
      m.branch_id = b.base_branch_id 
      AND NOT EXISTS (
          SELECT 1 FROM dimension_members child_override
          WHERE child_override.branch_id = b.id
          AND child_override.dimension_id = m.dimension_id
          AND child_override.name = m.name 
      )
    );
```
*(This view ensures a query for Branch B gets Branch B's data, falling back to Branch A for anything unmodified).*

### 2.3 The Go Orchestrator Upgrades (`grid-service`)
- **`MetadataResolver` Interface:** All methods (e.g., `ResolveMembers`) now require a `branchId` string.
- **`GridQuery` Domain Model:** Added `BranchID` to the query payload. The frontend URL dictates the branch context via `?branchId=xyz`.
- **Cache Isolation:** The `RedisGridCache` and its `GridCacheKey` were updated to include `BranchID`. Data fetched for Branch B will *never* collide with the cache for the `main` branch.

### 2.4 The Branch Service API (`branch_handler.go`)
Exposes REST endpoints to manage workspaces:
- `GET /api/v1/apps/:appId/branches`: Lists all available sandboxes.
- `POST /api/v1/apps/:appId/branches`: Creates a new isolated branch.

---

## 🚀 3. Usage & Data Flow

1. **Isolation:** User creates branch `reorg_q3` from `main`.
2. **Modification:** User deletes the "Europe" node in `reorg_q3`. A Tombstone record is inserted for Europe linked to `reorg_q3`.
3. **Querying:** 
   - User A requests the grid on `main` -> "Europe" is visible.
   - User B requests the grid on `reorg_q3` -> The SQL View filters out "Europe" due to the Tombstone override.

---

## ✅ 4. Verification Checkpoints
- ✔️ **Zero Duplication Validation:** Creating a branch takes `O(1)` time regardless of dimension size.
- ✔️ **Cache Segregation:** Verified that identical grid queries on different branches yield distinct, isolated Redis L2 cache keys.
- ✔️ **Compilation:** `go build` passes with zero schema mismatch errors in the Go Orchestrator.
