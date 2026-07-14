package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/platform/config"
	"marketplace-central/apps/server_core/internal/platform/logging"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	logger := logging.New()

	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		log.Fatalf("db config: %v", err)
	}
	pool, err := pgdb.NewPool(ctx, dbCfg)
	if err != nil {
		log.Fatalf("db pool: %v", err)
	}
	logger.Printf("server starting on %s", cfg.Addr)
	runtime, err := composition.NewRootRuntime(pool, dbCfg)
	if err != nil {
		log.Fatalf("root router: %v", err)
	}
	if runtime.PoolStats != nil {
		runtime.PoolStats.Start(ctx)
	}
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           runtime.Handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
	}
	serveErr := server.ListenAndServe()
	if runtime.PoolStats != nil {
		runtime.PoolStats.Stop()
	}
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
