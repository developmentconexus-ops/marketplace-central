package main

import (
	"context"
	"fmt"
	"log"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
	canonical "marketplace-central/apps/server_core/migrations"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := pgdb.LoadConfig()
	if err != nil {
		return err
	}
	pool, err := pgdb.NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	applied, err := migrate.Run(ctx, pool, canonical.Source())
	if err != nil {
		return err
	}
	fmt.Printf("applied %d migration(s)\n", applied)
	return nil
}
