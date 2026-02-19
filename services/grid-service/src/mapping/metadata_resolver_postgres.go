package mapping

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"quantatomai/grid-service/planner"
)

// -----------------------------
// Core Postgres Resolver
// -----------------------------

// PostgresMetadataResolver is a production-ready implementation of MetadataResolver
// backed by Postgres metadata tables.
type PostgresMetadataResolver struct {
	db      *sql.DB
	modelID string
	timeout time.Duration

	resolveMembersStmt  *sql.Stmt
	resolveMeasuresStmt *sql.Stmt
	resolveScenStmt     *sql.Stmt
}

// NewPostgresMetadataResolver constructs a new resolver with prepared statements.
func NewPostgresMetadataResolver(db *sql.DB, modelID string, timeout time.Duration) (*PostgresMetadataResolver, error) {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	r := &PostgresMetadataResolver{
		db:      db,
		modelID: modelID,
		timeout: timeout,
	}

	// Prepared statements using ANY($n) for arrays and soft-delete/effective dating.
	var err error

	r.resolveMembersStmt, err = db.Prepare(`
        SELECT m.id, m.code, m.name
        FROM members m
        JOIN dimensions d ON m.dimension_id = d.id
        WHERE d.model_id = $1
          AND d.name = $2
          AND m.code = ANY($3)
          AND m.is_active = TRUE
          AND (m.effective_start <= NOW())
          AND (m.effective_end IS NULL OR m.effective_end >= NOW())
        ORDER BY m.sequence ASC, m.code ASC
    `)
	if err != nil {
		return nil, fmt.Errorf("prepare resolveMembersStmt: %w", err)
	}

	r.resolveMeasuresStmt, err = db.Prepare(`
        SELECT id
        FROM measures
        WHERE model_id = $1
          AND code = ANY($2)
          AND is_active = TRUE
          AND (effective_start <= NOW())
          AND (effective_end IS NULL OR effective_end >= NOW())
        ORDER BY code ASC
    `)
	if err != nil {
		return nil, fmt.Errorf("prepare resolveMeasuresStmt: %w", err)
	}

	r.resolveScenStmt, err = db.Prepare(`
        SELECT id
        FROM scenarios
        WHERE model_id = $1
          AND code = ANY($2)
          AND is_active = TRUE
          AND (effective_start <= NOW())
          AND (effective_end IS NULL OR effective_end >= NOW())
        ORDER BY code ASC
    `)
	if err != nil {
		return nil, fmt.Errorf("prepare resolveScenStmt: %w", err)
	}

	return r, nil
}

// Close releases the prepared statements.
func (r *PostgresMetadataResolver) Close() error {
	var errs []string
	if r.resolveMembersStmt != nil {
		if err := r.resolveMembersStmt.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if r.resolveMeasuresStmt != nil {
		if err := r.resolveMeasuresStmt.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if r.resolveScenStmt != nil {
		if err := r.resolveScenStmt.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing statements: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (r *PostgresMetadataResolver) ResolveMembers(
	ctx context.Context,
	dim string,
	codes []string,
) ([]planner.MemberInfo, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.resolveMembersStmt.QueryContext(ctx, r.modelID, dim, pqStringArray(codes))
	if err != nil {
		return nil, fmt.Errorf("ResolveMembers query failed: %w", err)
	}
	defer rows.Close()

	var result []planner.MemberInfo
	for rows.Next() {
		var mi planner.MemberInfo
		if err := rows.Scan(&mi.ID, &mi.Code, &mi.Name); err != nil {
			return nil, fmt.Errorf("ResolveMembers scan failed: %w", err)
		}
		result = append(result, mi)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ResolveMembers rows error: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no members resolved for dimension %s", dim)
	}

	return result, nil
}

func (r *PostgresMetadataResolver) ResolveMeasureIDs(
	ctx context.Context,
	measures []string,
) ([]int64, error) {
	if len(measures) == 0 {
		return nil, fmt.Errorf("no measures provided")
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.resolveMeasuresStmt.QueryContext(ctx, r.modelID, pqStringArray(measures))
	if err != nil {
		return nil, fmt.Errorf("ResolveMeasureIDs query failed: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("ResolveMeasureIDs scan failed: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ResolveMeasureIDs rows error: %w", err)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no measures resolved")
	}

	return ids, nil
}

func (r *PostgresMetadataResolver) ResolveScenarioIDs(
	ctx context.Context,
	scenarios []string,
) ([]int64, error) {
	if len(scenarios) == 0 {
		return nil, fmt.Errorf("no scenarios provided")
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.resolveScenStmt.QueryContext(ctx, r.modelID, pqStringArray(scenarios))
	if err != nil {
		return nil, fmt.Errorf("ResolveScenarioIDs query failed: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("ResolveScenarioIDs scan failed: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ResolveScenarioIDs rows error: %w", err)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no scenarios resolved")
	}

	return ids, nil
}

// GetCurrencyCode returns the currency code for a given dimension/member pair.
func (r *PostgresMetadataResolver) GetCurrencyCode(ctx context.Context, dimensionID, memberID int64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var code sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT currency_code 
		FROM members 
		WHERE id = $1 AND dimension_id = $2
	`, memberID, dimensionID).Scan(&code)
	
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return code.String, nil
}

// BulkGetCurrencyCodes returns currency codes for many members at once.
func (r *PostgresMetadataResolver) BulkGetCurrencyCodes(ctx context.Context, dimensionID int64, memberIDs []int64) (map[int64]string, error) {
	if len(memberIDs) == 0 {
		return make(map[int64]string), nil
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, currency_code 
		FROM members 
		WHERE dimension_id = $1 AND id = ANY($2)
	`, dimensionID, pq.Array(memberIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[int64]string)
	for rows.Next() {
		var id int64
		var code sql.NullString
		if err := rows.Scan(&id, &code); err != nil {
			return nil, err
		}
		results[id] = code.String
	}
	return results, nil
}

// pqStringArray is a tiny helper to allow using []string with postgres ANY() operator.
type pqStringArray []string

func (a pqStringArray) Value() (driver.Value, error) {
	return pq.Array([]string(a)).Value()
}
