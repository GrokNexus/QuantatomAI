-- Extend compatibility metadata with hierarchy fields for tree rendering
ALTER TABLE members_compat
    ADD COLUMN IF NOT EXISTS path TEXT,
    ADD COLUMN IF NOT EXISTS parent_code TEXT;

CREATE INDEX IF NOT EXISTS idx_members_compat_dim_parent
    ON members_compat(dimension_id, parent_code);

-- Refresh branch view to include hierarchy fields
DROP VIEW IF EXISTS branch_view_members;
CREATE OR REPLACE VIEW branch_view_members AS
SELECT
    m.id,
    m.code,
    m.name,
    m.dimension_id,
    m.sequence,
    m.is_active,
    m.is_deleted,
    m.effective_start,
    m.effective_end,
    m.path,
    m.parent_code,
    'main'::text AS query_branch_id
FROM members_compat m;
