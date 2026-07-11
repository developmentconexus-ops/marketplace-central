package main

import (
	"context"
	"fmt"
	"os"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
	canonical "marketplace-central/apps/server_core/migrations"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) != 1 || args[0] != "migrate" {
		return fmt.Errorf("HPG_TESTDB_ACTION_INVALID")
	}
	cfg, err := testpostgres.LoadConfig(os.Getenv, "tenant_harness", "harness-test-key")
	if err != nil {
		return err
	}
	pool, err := pgdb.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("HPG_TARGET_INVALID: %w", err)
	}
	defer pool.Close()

	applied, err := migrate.Run(ctx, pool, canonical.Source())
	if err != nil {
		return fmt.Errorf("HPG_MIGRATION_FAILED: %w", err)
	}
	fmt.Printf("applied %d migration(s)\n", applied)
	return nil
}
