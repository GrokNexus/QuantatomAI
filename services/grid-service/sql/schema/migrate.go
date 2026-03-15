package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed *.sql
var schemaFS embed.FS

// Run applies embedded SQL migrations in filename order, ensuring each runs once.
func Run(ctx context.Context, db *sql.DB) error {
	if err := ensureTable(ctx, db); err != nil {
		return err
	}

	entries, err := schemaFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read embedded schemas: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, entry := range entries {
		name := entry.Name()
		if applied, err := alreadyApplied(ctx, db, name); err != nil {
			return err
		} else if applied {
			continue
		}

		contents, err := schemaFS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read schema %s: %w", name, err)
		}

		if err := applyOne(ctx, db, name, string(contents)); err != nil {
			return fmt.Errorf("apply schema %s: %w", name, err)
		}
	}

	return nil
}

func ensureTable(ctx context.Context, db *sql.DB) error {
	const stmt = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    name TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	return nil
}

func alreadyApplied(ctx context.Context, db *sql.DB, name string) (bool, error) {
	const q = `SELECT 1 FROM schema_migrations WHERE name = $1`
	var one int
	err := db.QueryRowContext(ctx, q, name).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", name, err)
	}
	return true, nil
}

func applyOne(ctx context.Context, db *sql.DB, name, sqlBody string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if _, err := tx.ExecContext(ctx, sqlBody); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("exec migration %s: %w (rollback failed: %v)", name, err, rbErr)
		}
		return fmt.Errorf("exec migration %s: %w", name, err)
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(name) VALUES($1)`, name); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("record migration %s: %w (rollback failed: %v)", name, err, rbErr)
		}
		return fmt.Errorf("record migration %s: %w", name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", name, err)
	}

	return nil
}

// NormalizeWhitespace can be used by future callers if we need to trim heredoc spacing.
func NormalizeWhitespace(s string) string {
	return strings.TrimSpace(s)
}
