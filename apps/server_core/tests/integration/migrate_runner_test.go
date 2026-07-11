//go:build integration

package integration

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
	canonical "marketplace-central/apps/server_core/migrations"
)

func TestMigrationRunnerIsIdempotent(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_migrate")

	// First run: apply pending migrations
	applied, err := migrate.Run(context.Background(), pool, canonical.Source())
	if err != nil {
		t.Fatalf("first run error: %v", err)
	}
	t.Logf("first run applied %d migrations", applied)

	// Second run: must apply zero (idempotent)
	applied2, err := migrate.Run(context.Background(), pool, canonical.Source())
	if err != nil {
		t.Fatalf("second run error: %v", err)
	}
	if applied2 != 0 {
		t.Fatalf("expected 0 migrations on second run, got %d", applied2)
	}

	// Verify schema_migrations has entries
	var count int
	err = pool.QueryRow(context.Background(), `SELECT count(*) FROM schema_migrations`).Scan(&count)
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if count == 0 {
		t.Fatal("expected applied migrations in schema_migrations")
	}
}
