-- QuantatomAI Phase 2 Validation Checks
-- Run after applying 07_tenant_control_plane.sql to verify tenant-control integrity.

-- 1. Dimensions must carry tenant context aligned to apps.
SELECT COUNT(*) AS invalid_dimension_tenant_rows
FROM dimensions d
JOIN apps a ON a.id = d.app_id
WHERE d.tenant_id IS DISTINCT FROM a.tenant_id;

-- 2. Dimension members must carry app and tenant context aligned to dimensions.
SELECT COUNT(*) AS invalid_member_context_rows
FROM dimension_members dm
JOIN dimensions d ON d.id = dm.dimension_id
WHERE dm.app_id IS DISTINCT FROM d.app_id
   OR dm.tenant_id IS DISTINCT FROM d.tenant_id;

-- 3. Security policies must remain tenant-aligned with apps.
SELECT COUNT(*) AS invalid_policy_tenant_rows
FROM security_policies sp
JOIN apps a ON a.id = sp.app_id
WHERE sp.tenant_id IS DISTINCT FROM a.tenant_id;

-- 4. Branches must remain tenant-aligned with apps.
SELECT COUNT(*) AS invalid_branch_tenant_rows
FROM branches b
JOIN apps a ON a.id = b.app_id
WHERE b.tenant_id IS DISTINCT FROM a.tenant_id;

-- 5. Each tenant should have at most one write region.
SELECT tenant_id, COUNT(*) AS write_region_count
FROM tenant_regions
WHERE is_write_region = TRUE
GROUP BY tenant_id
HAVING COUNT(*) > 1;

-- 6. AI policies must default to tenant-bounded retrieval.
SELECT COUNT(*) AS invalid_ai_policy_rows
FROM tenant_ai_policies
WHERE retrieval_scope <> 'tenant-only'
   OR allow_cross_tenant_learning = TRUE;

-- 7. App partition records must point to a registered tenant write region.
SELECT COUNT(*) AS invalid_app_partition_rows
FROM app_partitions ap
LEFT JOIN tenant_regions tr
  ON tr.tenant_id = ap.tenant_id
 AND tr.region_code = ap.write_region
WHERE tr.tenant_id IS NULL;