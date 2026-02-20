-- =========================================================================================
-- QuantatomAI Layer 8.1: Enterprise Git-Flow for Metadata
-- Description: Unlocks isolated branching and committing for structural models (Hierarchies).
-- Pattern: Sparse Override Delta-Branching.
-- =========================================================================================

-- 1. BRANCHES (The Workspaces)
-- Represents an isolated sandbox for structural changes.
CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_id UUID REFERENCES apps(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL, -- e.g., 'main', 'reorg_europe_q3'
    base_branch_id UUID REFERENCES branches(id), -- Self-referencing (The parent branch)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, name)
);

-- Note: In application logic, every App created automatically gets a 'main' branch.

-- 2. COMMITS (The Snapshots)
-- Represents a finalized group of metadata changes within a branch.
CREATE TABLE commits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    author_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. DIMENSION_MEMBERS UPGRADE (Delta-Branching Structure)
-- IMPORTANT: This table is the core of the Moat. We are augmenting the existing table.

-- A) Add Branching Coordinates
ALTER TABLE dimension_members 
ADD COLUMN branch_id UUID REFERENCES branches(id) ON DELETE CASCADE;

ALTER TABLE dimension_members 
ADD COLUMN commit_id UUID REFERENCES commits(id); -- Null until committed

-- B) Add Tombstone Marker (Since we cannot truly DELETE a main row in a child branch)
ALTER TABLE dimension_members 
ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

-- C) Destroy the Old Global Unique Constraint & Create Branch-Aware Unique Constraint
-- Previously, a member name was unique per dimension globally.
-- Now, a member name is unique per dimension _within a specific branch_.
-- Note: Assuming previous constraint name was 'dimension_members_dimension_id_name_key', 
-- you may need to drop it dynamically in a real migration, but we define the new state here.
ALTER TABLE dimension_members 
ADD CONSTRAINT dimension_members_branch_uniq UNIQUE (dimension_id, name, branch_id);

-- 4. DELTA OVERRIDE VIEW (The Magician's Trick)
-- This view allows the backend to query "What does the world look like in Branch X?"
-- It seamlessly overlays the new/modified/deleted rows in Branch X on top of the 'main' branch.
-- To query this: SELECT * FROM branch_view_members WHERE query_branch_id = 'your-branch-uuid';

CREATE OR REPLACE VIEW branch_view_members AS
SELECT 
    m.*, 
    b.id as query_branch_id
FROM branches b
JOIN dimension_members m ON 
    -- The member belongs directly to the queried branch
    m.branch_id = b.id 
    OR 
    -- OR, the member belongs to the parent branch (e.g., 'main'), 
    -- BUT we ensure it hasn't been overridden or tombstoned in the child branch.
    (
      m.branch_id = b.base_branch_id 
      AND NOT EXISTS (
          SELECT 1 FROM dimension_members child_override
          WHERE child_override.branch_id = b.id
          AND child_override.dimension_id = m.dimension_id
          AND child_override.name = m.name 
          -- Name is the business key linking the main row to the branch override row
      )
    );

-- Note: In a deep tree of branches, this view would need recursive CTEs. 
-- For MVP Layer 8.1, we support a depth of 1 (Main -> Sandbox).
