// Command catalogingest is the operator entry point that finally makes the
// catalog ingest path run for real. It is the composition root's job to wire
// a live Sankhya-Oracle feed to Catalog and drive it to exhaustion; nothing
// else in the platform is allowed to import both the Oracle driver and the
// catalog context in the same breath.
//
// The Oracle driver itself is registered in oracle_cgo.go, not here, and only
// under a cgo build — the same split internal/modules/internal_read/adapters/oracle
// already uses. That keeps this file, and `go build ./...` on a host with no C
// toolchain, green; only a cgo build with the Oracle Instant Client actually
// talks to Oracle.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	oracleconfig "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

const defaultPageSize = 200

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		// Never a fabricated success: any failure here is printed verbatim to
		// stderr and the process exits non-zero. A "0 processed, exit 0" on a
		// real error would be the invisible failure this command exists to end.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		return fmt.Errorf("catalogingest: database config: %w", err)
	}
	if err := requireTenantConfigured(os.Getenv); err != nil {
		return err
	}
	tenantID, err := tenant.Parse(dbCfg.DefaultTenantID)
	if err != nil {
		return fmt.Errorf("catalogingest: tenant: %w", err)
	}
	instance := strings.TrimSpace(os.Getenv("MPC_SANKHYA_INSTANCE"))
	if instance == "" {
		return errors.New("catalogingest: MPC_SANKHYA_INSTANCE is required")
	}
	pageSize, err := loadPageSize(os.Getenv)
	if err != nil {
		return fmt.Errorf("catalogingest: %w", err)
	}

	pool, err := pgdb.NewPool(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("catalogingest: open postgres pool: %w", err)
	}
	defer pool.Close()

	oracleCfg, err := oracleconfig.LoadConfigFromEnv(os.Getenv)
	if err != nil {
		return fmt.Errorf("catalogingest: oracle config: %w", err)
	}
	oracleDB, err := openOracleDB(ctx, oracleCfg)
	if err != nil {
		// The error can carry connection text from the driver; it is never a
		// credential by itself (the driver does not echo the password back),
		// but nothing beyond the driver's own message is added here.
		return fmt.Errorf("catalogingest: open oracle: %w", err)
	}
	defer oracleDB.Close()

	wiring, err := composition.WireCatalog(pool, oracleDB, instance)
	if err != nil {
		return fmt.Errorf("catalogingest: wire catalog: %w", err)
	}

	report, err := composition.RunCatalogIngest(ctx, wiring.Module, wiring.Feed.CatalogFeed, tenantID, pageSize)
	if err != nil {
		return fmt.Errorf("catalogingest: run ingest: %w", err)
	}

	fmt.Printf("catalog ingest report: pages=%d observed=%d created=%d changed=%d idempotent=%d conflicts=%d\n",
		report.Pages, report.Observed, report.Created, report.Changed, report.Idempotent, len(report.Conflicts))
	return nil
}

// requireTenantConfigured fails closed when MC_DEFAULT_TENANT_ID is unset or
// empty. pgdb.LoadConfig (internal/platform/pgdb/config.go) silently
// substitutes "tenant_default" for every caller when the variable is absent
// — correct for its many read-oriented callers, but wrong here: this command
// performs live writes across the entire catalogue (D-39,
// .mnfs/HARNESS-DEBTS.md). An operator who forgets to export the variable
// must see a failure naming it, not have thousands of rows land under an
// invented tenant with a silent exit code 0 (global constraint 9: unknown
// never becomes a plausible default).
//
// getenv is read directly here, exactly mirroring how pgdb.LoadConfig reads
// the same variable (os.Getenv("MC_DEFAULT_TENANT_ID"), no trimming), so
// this guard can never disagree with what LoadConfig saw. It keys on the
// variable being absent/empty, never on the resolved value: a deployment
// that legitimately names its tenant "tenant_default" still passes, and a
// typo'd variable name still fails.
//
// pgdb.LoadConfig's defaulting itself is not changed — it is shared
// infrastructure with other callers (cmd/server) that may have distinct
// reasons to tolerate the default today.
func requireTenantConfigured(getenv func(string) string) error {
	if getenv("MC_DEFAULT_TENANT_ID") == "" {
		return errors.New("catalogingest: MC_DEFAULT_TENANT_ID is required (refusing to silently default to \"tenant_default\" for a live write path)")
	}
	return nil
}

// loadPageSize reads the operator-tunable page size. Unlike a domain fact, a
// paging batch size has no "unknown" state to preserve — it is an operational
// default, not a fabricated business value — so an unset variable falls back
// to defaultPageSize rather than failing closed.
func loadPageSize(getenv func(string) string) (int, error) {
	raw := strings.TrimSpace(getenv("MPC_CATALOG_INGEST_PAGE_SIZE"))
	if raw == "" {
		return defaultPageSize, nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("MPC_CATALOG_INGEST_PAGE_SIZE must be a positive integer, got %q", raw)
	}
	return parsed, nil
}
