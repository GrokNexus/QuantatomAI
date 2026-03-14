package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestRun_Phase2TenantControlPlane(t *testing.T) {
	testDatabaseURL := os.Getenv("DATABASE_URL")
	if testDatabaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	adminDB, targetDatabaseURL, databaseName := createEphemeralDatabase(t, ctx, testDatabaseURL)
	defer dropEphemeralDatabase(t, adminDB, databaseName)

	db, err := sql.Open("postgres", targetDatabaseURL)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, db.PingContext(ctx))
	require.NoError(t, Run(ctx, db))
	require.NoError(t, Run(ctx, db))

	var applied bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM schema_migrations
			WHERE name = '07_tenant_control_plane.sql'
		)
	`).Scan(&applied)
	require.NoError(t, err)
	require.True(t, applied)

	tenantID := insertTenantFixture(t, ctx, db)
	userID := insertUserFixture(t, ctx, db, tenantID)
	appID := insertAppFixture(t, ctx, db, tenantID, userID)

	insertTenantRegionFixture(t, ctx, db, tenantID, "us-east-1", true)
	insertTenantKeyDomainFixture(t, ctx, db, tenantID, "us-east-1", "app-data")
	insertTenantQuotaFixture(t, ctx, db, tenantID)
	insertTenantAIPolicyFixture(t, ctx, db, tenantID)

	dimensionID := insertDimensionFixture(t, ctx, db, appID)
	insertDimensionMemberFixture(t, ctx, db, dimensionID)
	insertSecurityPolicyFixture(t, ctx, db, appID, userID)
	insertBranchFixture(t, ctx, db, appID)
	insertAppPartitionFixture(t, ctx, db, appID, "us-east-1")

	assertTenantPropagation(t, ctx, db, tenantID)
	assertTenantAIDefaults(t, ctx, db, tenantID)
	assertSingleWriteRegionEnforced(t, ctx, db, tenantID)
	assertTenantKeyDomainRequiresKnownRegion(t, ctx, db, tenantID)
	assertValidationQueriesPass(t, ctx, db)
}

func createEphemeralDatabase(t *testing.T, ctx context.Context, databaseURL string) (*sql.DB, string, string) {
	t.Helper()

	parsedURL, err := url.Parse(databaseURL)
	require.NoError(t, err)

	baseName := strings.TrimPrefix(parsedURL.Path, "/")
	require.NotEmpty(t, baseName)

	databaseName := fmt.Sprintf("%s_phase2_%d", sanitizeDatabaseName(baseName), time.Now().UnixNano())

	adminURL := *parsedURL
	adminURL.Path = "/postgres"

	adminDB, err := sql.Open("postgres", adminURL.String())
	require.NoError(t, err)
	require.NoError(t, adminDB.PingContext(ctx))

	_, err = adminDB.ExecContext(ctx, "CREATE DATABASE "+pq.QuoteIdentifier(databaseName))
	require.NoError(t, err)

	targetURL := *parsedURL
	targetURL.Path = "/" + databaseName

	return adminDB, targetURL.String(), databaseName
}

func dropEphemeralDatabase(t *testing.T, adminDB *sql.DB, databaseName string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, _ = adminDB.ExecContext(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1 AND pid <> pg_backend_pid()
	`, databaseName)
	_, _ = adminDB.ExecContext(ctx, "DROP DATABASE IF EXISTS "+pq.QuoteIdentifier(databaseName))
	_ = adminDB.Close()
}

func sanitizeDatabaseName(name string) string {
	clean := strings.ToLower(name)
	clean = strings.ReplaceAll(clean, "-", "_")
	clean = strings.ReplaceAll(clean, ".", "_")
	return clean
}

func insertTenantFixture(t *testing.T, ctx context.Context, db *sql.DB) string {
	t.Helper()

	var tenantID string
	err := db.QueryRowContext(ctx, `
		INSERT INTO tenants (name)
		VALUES ('Phase 2 Tenant')
		RETURNING id
	`).Scan(&tenantID)
	require.NoError(t, err)
	return tenantID
}

func insertUserFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) string {
	t.Helper()

	var userID string
	err := db.QueryRowContext(ctx, `
		INSERT INTO users (tenant_id, email, password_hash, role)
		VALUES ($1, 'phase2@example.com', 'hash', 'admin')
		RETURNING id
	`, tenantID).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func insertAppFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID, userID string) string {
	t.Helper()

	var appID string
	err := db.QueryRowContext(ctx, `
		INSERT INTO apps (tenant_id, name, created_by)
		VALUES ($1, 'Phase 2 App', $2)
		RETURNING id
	`, tenantID, userID).Scan(&appID)
	require.NoError(t, err)
	return appID
}

func insertTenantRegionFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID, regionCode string, isWriteRegion bool) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_regions (tenant_id, region_code, region_role, is_write_region)
		VALUES ($1, $2, 'primary', $3)
	`, tenantID, regionCode, isWriteRegion)
	require.NoError(t, err)
}

func insertTenantKeyDomainFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID, regionCode, purpose string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_key_domains (tenant_id, region_code, purpose, kms_provider, key_uri)
		VALUES ($1, $2, $3, 'aws-kms', 'arn:aws:kms:us-east-1:123456789012:key/phase2')
	`, tenantID, regionCode, purpose)
	require.NoError(t, err)
}

func insertTenantQuotaFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_quota_policies (tenant_id)
		VALUES ($1)
	`, tenantID)
	require.NoError(t, err)
}

func insertTenantAIPolicyFixture(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_ai_policies (tenant_id)
		VALUES ($1)
	`, tenantID)
	require.NoError(t, err)
}

func insertDimensionFixture(t *testing.T, ctx context.Context, db *sql.DB, appID string) string {
	t.Helper()

	var dimensionID string
	err := db.QueryRowContext(ctx, `
		INSERT INTO dimensions (app_id, name, type)
		VALUES ($1, 'Account', 'standard')
		RETURNING id
	`, appID).Scan(&dimensionID)
	require.NoError(t, err)
	return dimensionID
}

func insertDimensionMemberFixture(t *testing.T, ctx context.Context, db *sql.DB, dimensionID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO dimension_members (dimension_id, name, path)
		VALUES ($1, 'Revenue', 'global.revenue')
	`, dimensionID)
	require.NoError(t, err)
}

func insertSecurityPolicyFixture(t *testing.T, ctx context.Context, db *sql.DB, appID, userID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO security_policies (app_id, name, rules, user_id)
		VALUES ($1, 'Planner Policy', '{"Region": ["NA"]}', $2)
	`, appID, userID)
	require.NoError(t, err)
}

func insertBranchFixture(t *testing.T, ctx context.Context, db *sql.DB, appID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO branches (app_id, name)
		VALUES ($1, 'main')
	`, appID)
	require.NoError(t, err)
}

func insertAppPartitionFixture(t *testing.T, ctx context.Context, db *sql.DB, appID, regionCode string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO app_partitions (
			app_id,
			write_region,
			hot_namespace,
			warm_partition_template,
			cold_object_prefix,
			event_topic_prefix,
			cache_namespace
		) VALUES (
			$1,
			$2,
			'tenant-cache-hot',
			'tenant={{tenant_id}}/app={{app_id}}',
			's3://quantatomai/archive/tenant/{{tenant_id}}/app/{{app_id}}',
			'tenant.phase2.events',
			'tenant:phase2:cache'
		)
	`, appID, regionCode)
	require.NoError(t, err)
}

func assertTenantPropagation(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	var dimensionTenantID string
	err := db.QueryRowContext(ctx, `SELECT tenant_id::text FROM dimensions LIMIT 1`).Scan(&dimensionTenantID)
	require.NoError(t, err)
	require.Equal(t, tenantID, dimensionTenantID)

	var memberTenantID, memberAppID string
	err = db.QueryRowContext(ctx, `SELECT tenant_id::text, app_id::text FROM dimension_members LIMIT 1`).Scan(&memberTenantID, &memberAppID)
	require.NoError(t, err)
	require.Equal(t, tenantID, memberTenantID)
	require.NotEmpty(t, memberAppID)

	var policyTenantID string
	err = db.QueryRowContext(ctx, `SELECT tenant_id::text FROM security_policies LIMIT 1`).Scan(&policyTenantID)
	require.NoError(t, err)
	require.Equal(t, tenantID, policyTenantID)

	var branchTenantID string
	err = db.QueryRowContext(ctx, `SELECT tenant_id::text FROM branches LIMIT 1`).Scan(&branchTenantID)
	require.NoError(t, err)
	require.Equal(t, tenantID, branchTenantID)

	var appPartitionTenantID string
	err = db.QueryRowContext(ctx, `SELECT tenant_id::text FROM app_partitions LIMIT 1`).Scan(&appPartitionTenantID)
	require.NoError(t, err)
	require.Equal(t, tenantID, appPartitionTenantID)
}

func assertTenantAIDefaults(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	var retrievalScope string
	var allowCrossTenantLearning bool
	err := db.QueryRowContext(ctx, `
		SELECT retrieval_scope, allow_cross_tenant_learning
		FROM tenant_ai_policies
		WHERE tenant_id = $1
	`, tenantID).Scan(&retrievalScope, &allowCrossTenantLearning)
	require.NoError(t, err)
	require.Equal(t, "tenant-only", retrievalScope)
	require.False(t, allowCrossTenantLearning)
}

func assertSingleWriteRegionEnforced(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_regions (tenant_id, region_code, region_role, is_write_region)
		VALUES ($1, 'us-west-2', 'secondary', TRUE)
	`, tenantID)
	require.Error(t, err)
}

func assertTenantKeyDomainRequiresKnownRegion(t *testing.T, ctx context.Context, db *sql.DB, tenantID string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_key_domains (tenant_id, region_code, purpose, kms_provider, key_uri)
		VALUES ($1, 'eu-west-1', 'audit', 'aws-kms', 'arn:aws:kms:eu-west-1:123456789012:key/phase2')
	`, tenantID)
	require.Error(t, err)
}

func assertValidationQueriesPass(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	checks := []string{
		`SELECT COUNT(*) FROM dimensions d JOIN apps a ON a.id = d.app_id WHERE d.tenant_id IS DISTINCT FROM a.tenant_id`,
		`SELECT COUNT(*) FROM dimension_members dm JOIN dimensions d ON d.id = dm.dimension_id WHERE dm.app_id IS DISTINCT FROM d.app_id OR dm.tenant_id IS DISTINCT FROM d.tenant_id`,
		`SELECT COUNT(*) FROM security_policies sp JOIN apps a ON a.id = sp.app_id WHERE sp.tenant_id IS DISTINCT FROM a.tenant_id`,
		`SELECT COUNT(*) FROM branches b JOIN apps a ON a.id = b.app_id WHERE b.tenant_id IS DISTINCT FROM a.tenant_id`,
		`SELECT COUNT(*) FROM tenant_ai_policies WHERE retrieval_scope <> 'tenant-only' OR allow_cross_tenant_learning = TRUE`,
		`SELECT COUNT(*) FROM app_partitions ap LEFT JOIN tenant_regions tr ON tr.tenant_id = ap.tenant_id AND tr.region_code = ap.write_region WHERE tr.tenant_id IS NULL`,
	}

	for _, check := range checks {
		var count int
		err := db.QueryRowContext(ctx, check).Scan(&count)
		require.NoError(t, err)
		require.Zero(t, count, check)
	}
}
