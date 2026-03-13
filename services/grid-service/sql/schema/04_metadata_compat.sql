-- Compatibility metadata schema to align with current PostgresMetadataResolver queries
-- Model-level metadata (text model_id for flexibility during bootstrap)
CREATE TABLE IF NOT EXISTS dimensions_compat (
    id BIGSERIAL PRIMARY KEY,
    model_id TEXT NOT NULL,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_dimensions_compat_model_name
    ON dimensions_compat(model_id, name);

CREATE TABLE IF NOT EXISTS members_compat (
    id BIGSERIAL PRIMARY KEY,
    dimension_id BIGINT NOT NULL REFERENCES dimensions_compat(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    sequence INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    effective_start TIMESTAMPTZ DEFAULT NOW(),
    effective_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_members_compat_dim_code
    ON members_compat(dimension_id, code);
CREATE INDEX IF NOT EXISTS idx_members_compat_dim_sequence
    ON members_compat(dimension_id, sequence);

CREATE TABLE IF NOT EXISTS measures (
    id BIGSERIAL PRIMARY KEY,
    model_id TEXT NOT NULL,
    code TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    effective_start TIMESTAMPTZ DEFAULT NOW(),
    effective_end TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_measures_model_code
    ON measures(model_id, code);

CREATE TABLE IF NOT EXISTS scenarios (
    id BIGSERIAL PRIMARY KEY,
    model_id TEXT NOT NULL,
    code TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    effective_start TIMESTAMPTZ DEFAULT NOW(),
    effective_end TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_scenarios_model_code
    ON scenarios(model_id, code);

-- Compatibility view used by resolver (branch-aware stub)
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
    'main'::text AS query_branch_id
FROM members_compat m;
