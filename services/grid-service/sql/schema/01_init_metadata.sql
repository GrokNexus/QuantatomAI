-- =========================================================================================
-- QuantatomAI Layer 2.1: Metadata Registry (Production Schema)
-- Description: The "Spine" of the Lattice. Stores Dimensions, Hierarchies, and Members.
-- Extensions: ltree (Hierarchy), vector (AI Embeddings), uuid-ossp (IDs).
-- =========================================================================================

-- 1. EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE EXTENSION IF NOT EXISTS "vector"; -- Requires pgvector installed
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy search (Show me "Net Sales")

-- ... (Tables define here) ...

-- 6. INDEXING STRATEGY (The "Speed")
-- GIST Index for extremely fast hierarchy queries: WHERE path <@ 'Global.NA'
CREATE INDEX idx_members_path_gist ON dimension_members USING GIST (path);
CREATE INDEX idx_members_path_btree ON dimension_members USING BTREE (path); -- For sorting
-- Ultra Diamond Upgrade: Trigram Index for "Show me 'Net Sales'"
CREATE INDEX idx_members_name_trgm ON dimension_members USING GIN (name gin_trgm_ops); 
CREATE INDEX idx_members_attributes ON dimension_members USING GIN (attributes); -- JSONB queries

-- ... (Member Mappings & Security Policies) ...

-- 9. AUDIT TRIGGER FUNCTIONS (Prep for Layer 2.3)
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Ultra Diamond Upgrade: Apply trigger to ALL modify-able tables
CREATE TRIGGER trigger_update_app_timestamp BEFORE UPDATE ON apps FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER trigger_update_dim_timestamp BEFORE UPDATE ON dimensions FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER trigger_update_mem_timestamp BEFORE UPDATE ON dimension_members FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER trigger_update_pol_timestamp BEFORE UPDATE ON security_policies FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    plan_tier VARCHAR(50) DEFAULT 'enterprise', -- 'standard', 'enterprise', 'ultra'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Argon2 hash
    role VARCHAR(50) DEFAULT 'planner', -- 'admin', 'planner', 'viewer'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

-- 3. APPLICATION REGISTRY (The "Lattice Contianer")
CREATE TABLE apps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    planning_mode VARCHAR(50) DEFAULT 'standard', -- 'continuous', 'scenario_based'
    default_currency VARCHAR(3) DEFAULT 'USD',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- 4. DIMENSION DEFINITIONS (The "Axes")
CREATE TABLE dimensions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_id UUID REFERENCES apps(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- e.g., "Account", "Region", "Time"
    type VARCHAR(50) NOT NULL, -- 'standard', 'time', 'measure', 'scenario'
    is_core BOOLEAN DEFAULT FALSE, -- Core dims required by system
    sort_order INT DEFAULT 0,
    properties JSONB DEFAULT '{}', -- Extra config like { "balance_type": "credit" }
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, name)
);

-- 5. DIMENSION MEMBERS (The "nodes") -- MOAT: Uses LTREE for O(1) Descendant Lookups
CREATE TABLE dimension_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dimension_id UUID REFERENCES dimensions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    
    -- Hierarchy Path (Materialized Path)
    -- Format: root.parent.child (e.g., "Global.NA.USA.NY")
    path ltree NOT NULL,
    
    -- Fast Parent Lookup (Optional but good for integrity)
    parent_id UUID REFERENCES dimension_members(id),
    
    -- Order siblings
    weight INT DEFAULT 0,
    
    -- Member Properties (e.g., { "currency": "USD", "manager": "John Doe" })
    attributes JSONB DEFAULT '{}',
    
    -- Formula / Logic
    formula TEXT, -- AtomScript fragment for this member
    
    -- AI Embeddings (For Semantic Search "Show me revenue members")
    embedding vector(384),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. INDEXING STRATEGY (The "Speed")
-- GIST Index for extremely fast hierarchy queries: WHERE path <@ 'Global.NA'
CREATE INDEX idx_members_path_gist ON dimension_members USING GIST (path);
CREATE INDEX idx_members_path_btree ON dimension_members USING BTREE (path); -- For sorting
CREATE INDEX idx_members_name_trgm ON dimension_members USING GIN (name vector_ops); -- Fuzzy search? No, std btree for exact
CREATE INDEX idx_members_attributes ON dimension_members USING GIN (attributes); -- JSONB queries

-- 7. ATTRIBUTE VALUES (Cross-Reference)
-- Maps a member to another dimension member (e.g., Entity 'NY' maps to Currency 'USD')
CREATE TABLE member_mappings (
    source_member_id UUID REFERENCES dimension_members(id) ON DELETE CASCADE,
    target_dimension_id UUID REFERENCES dimensions(id) ON DELETE CASCADE,
    target_member_id UUID REFERENCES dimension_members(id),
    PRIMARY KEY (source_member_id, target_dimension_id)
);

-- 8. SECURITY POLICIES (Moat: Holographic ACLs)
CREATE TABLE security_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    app_id UUID REFERENCES apps(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    
    -- The Policy Definition
    -- e.g., { "Region": ["USA", "Canada"], "Scenario": ["Budget"] }
    rules JSONB NOT NULL,
    
    -- Grant
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    permission_level VARCHAR(20) DEFAULT 'read', -- 'read', 'write', 'none'
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 9. AUDIT TRIGGER FUNCTIONS (Prep for Layer 2.3)
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_app_timestamp
BEFORE UPDATE ON apps
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
