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
//
// Which tenant to ingest for, which Sankhya instance to read, and how big a
// page to request are parameters of the run, so they arrive as flags:
//
//	go run ./cmd/catalogingest -tenant tenant_default -instance <name> [-page-size 200]
//
// The environment carries only deployment configuration — the Postgres URL
// and encryption key via pgdb.LoadConfig, the Oracle connection via
// oracleconfig.LoadConfigFromEnv — each read once by its typed loader.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	oracleconfig "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

const (
	commandName     = "catalogingest"
	defaultPageSize = 200
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx, os.Args[1:]); err != nil {
		// -h is a request, not a failure: the usage text is already on stderr
		// and this process has nothing to report.
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		// Never a fabricated success: any failure here is printed verbatim to
		// stderr and the process exits non-zero. A "0 processed, exit 0" on a
		// real error would be the invisible failure this command exists to end.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// options are the parameters of one run — what this invocation is for, as
// opposed to how the deployment is configured.
type options struct {
	tenant   string
	instance string
	pageSize int
}

// parseOptions reads the run's parameters from args and nothing else. There
// is deliberately no environment fallback for any of them: a required flag
// has no default that a forgotten export could be mistaken for, so the
// fail-closed behaviour a live write path needs is a property of the channel
// rather than a guard that has to be remembered (global constraint 9 —
// unknown never becomes a plausible default). This is what the retired
// requireTenantConfigured guard was compensating for: pgdb.LoadConfig
// substitutes "tenant_default" for an absent variable, so a tenant arriving
// through the environment could not be told apart from a forgotten export.
//
// usageOut receives the usage text; the FlagSet's own printer is silenced so
// that every failure leaves this process through main's single error path
// instead of being written twice.
func parseOptions(args []string, usageOut io.Writer) (options, error) {
	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tenantFlag := fs.String("tenant", "", "tenant id this run ingests for (required)")
	instanceFlag := fs.String("instance", "", "Sankhya instance the catalogue is read from (required)")
	pageSizeFlag := fs.Int("page-size", defaultPageSize, "catalogue rows requested per page")

	if err := fs.Parse(args); err != nil {
		fs.SetOutput(usageOut)
		fs.Usage()
		return options{}, err
	}
	if fs.NArg() > 0 {
		return options{}, fmt.Errorf("unexpected argument %q (this command takes flags only)", fs.Arg(0))
	}
	tenantID := strings.TrimSpace(*tenantFlag)
	if tenantID == "" {
		return options{}, errors.New("-tenant is required: the tenant to ingest for is a parameter of this run, not a deployment setting")
	}
	instance := strings.TrimSpace(*instanceFlag)
	if instance == "" {
		return options{}, errors.New("-instance is required: the Sankhya instance this run reads is a parameter of this run, not a deployment setting")
	}
	if *pageSizeFlag <= 0 {
		return options{}, fmt.Errorf("-page-size must be a positive integer, got %d", *pageSizeFlag)
	}
	return options{tenant: tenantID, instance: instance, pageSize: *pageSizeFlag}, nil
}

func run(ctx context.Context, args []string) error {
	opts, err := parseOptions(args, os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return err
		}
		return fmt.Errorf("catalogingest: %w", err)
	}
	tenantID, err := tenant.Parse(opts.tenant)
	if err != nil {
		return fmt.Errorf("catalogingest: tenant: %w", err)
	}

	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		return fmt.Errorf("catalogingest: database config: %w", err)
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

	wiring, err := composition.WireCatalog(pool, oracleDB, opts.instance)
	if err != nil {
		return fmt.Errorf("catalogingest: wire catalog: %w", err)
	}

	report, err := composition.RunCatalogIngest(ctx, wiring.Module, wiring.Feed.CatalogFeed, tenantID, opts.pageSize)
	if err != nil {
		return fmt.Errorf("catalogingest: run ingest: %w", err)
	}

	fmt.Printf("catalog ingest report: pages=%d observed=%d created=%d changed=%d idempotent=%d conflicts=%d\n",
		report.Pages, report.Observed, report.Created, report.Changed, report.Idempotent, len(report.Conflicts))
	return nil
}
