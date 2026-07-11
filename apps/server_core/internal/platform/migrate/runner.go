package migrate

import (
	"context"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Filenames returns the sorted, full filenames in a migration source.
func Filenames(source fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(source, ".")
	if err != nil {
		return nil, fmt.Errorf("read migration source: %w", err)
	}

	filenames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".sql" {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)
	return filenames, nil
}

// Run applies all pending SQL migrations from source to the database.
// It returns the number of migrations applied and any error encountered.
// Migrations are applied in lexicographic filename order and tracked in
// schema_migrations to guarantee idempotency.
func Run(ctx context.Context, pool *pgxpool.Pool, source fs.FS) (int, error) {
	// Ensure schema_migrations table exists
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename   TEXT        PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return 0, fmt.Errorf("ensure schema_migrations: %w", err)
	}

	filenames, err := Filenames(source)
	if err != nil {
		return 0, err
	}

	// Query already-applied filenames
	rows, err := pool.Query(ctx, `SELECT filename FROM schema_migrations`)
	if err != nil {
		return 0, fmt.Errorf("query schema_migrations: %w", err)
	}
	applied := map[string]bool{}
	for rows.Next() {
		var fn string
		if err := rows.Scan(&fn); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan schema_migrations row: %w", err)
		}
		applied[fn] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate schema_migrations: %w", err)
	}

	count := 0
	for _, filename := range filenames {
		if applied[filename] {
			continue
		}

		content, err := fs.ReadFile(source, filename)
		if err != nil {
			return count, fmt.Errorf("read migration file %s: %w", filename, err)
		}
		if err := applyFile(ctx, pool, filename, content); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func applyFile(ctx context.Context, pool *pgxpool.Pool, filename string, content []byte) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", filename, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("migration %s: %w", filename, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, filename); err != nil {
		return fmt.Errorf("record migration %s: %w", filename, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", filename, err)
	}
	return nil
}
