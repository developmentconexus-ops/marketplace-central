// Command listingsreprocess is the offline counterpart to cmd/listingsingest:
// it re-derives listing facts from payloads already stored in Postgres and
// folds them back in. It makes NO calls to Mercado Livre or any other
// marketplace — it only re-maps bytes this binary already holds, using the
// mapper compiled into it now, which may have been corrected since those
// bytes were first observed. Each observation's original observed_at is kept
// as stored, because nothing is observed by a reprocess; only the payload and
// the mapping change.
//
// There is accordingly no token, no credential, and no network client
// anywhere in this command — WireListingsReprocess takes a database pool and
// nothing else, and RunListingsReprocess only ever reads stored rows and
// folds them back through the module.
//
// Which tenant to reprocess and how big a stored-observations page to
// request are parameters of the run, so they arrive as flags:
//
//	go run ./cmd/listingsreprocess -tenant tenant_default [-page-size 50]
//
// The environment carries only deployment configuration — the database URL —
// read once by pgdb.LoadConfig. There is no encryption key to load here: a
// reprocess never decrypts a credential, because it never resolves one.
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
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

const (
	commandName     = "listingsreprocess"
	defaultPageSize = 50
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
		// stderr and the process exits non-zero.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// options are the parameters of one run — what this invocation is for, as
// opposed to how the deployment is configured.
type options struct {
	tenant   string
	pageSize int
}

// parseOptions reads the run's parameters from args and nothing else. There
// is deliberately no environment fallback for either of them: a required flag
// has no default that a forgotten export could be mistaken for, so the
// fail-closed behaviour is a property of the channel rather than a guard that
// has to be remembered.
//
// usageOut receives the usage text; the FlagSet's own printer is silenced so
// that every failure leaves this process through main's single error path
// instead of being written twice.
func parseOptions(args []string, usageOut io.Writer) (options, error) {
	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tenantFlag := fs.String("tenant", "", "tenant id this run reprocesses stored listings for (required)")
	pageSizeFlag := fs.Int("page-size", defaultPageSize, "stored observations requested per page")

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
		return options{}, errors.New("-tenant is required: the tenant to reprocess is a parameter of this run, not a deployment setting")
	}
	if *pageSizeFlag <= 0 {
		return options{}, fmt.Errorf("-page-size must be a positive integer, got %d", *pageSizeFlag)
	}
	return options{tenant: tenantID, pageSize: *pageSizeFlag}, nil
}

func run(ctx context.Context, args []string) error {
	opts, err := parseOptions(args, os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return err
		}
		return fmt.Errorf("listingsreprocess: %w", err)
	}
	tenantID, err := tenant.Parse(opts.tenant)
	if err != nil {
		return fmt.Errorf("listingsreprocess: tenant: %w", err)
	}

	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		return fmt.Errorf("listingsreprocess: database config: %w", err)
	}

	pool, err := pgdb.NewPool(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("listingsreprocess: open postgres pool: %w", err)
	}
	defer pool.Close()

	wiring := composition.WireListingsReprocess(pool)

	report, err := composition.RunListingsReprocess(ctx, wiring.Module, wiring.Mapper, tenantID, opts.pageSize)
	if err != nil {
		return fmt.Errorf("listingsreprocess: run reprocess: %w", err)
	}

	fmt.Printf("listings reprocess report: pages=%d read=%d created=%d changed=%d idempotent=%d\n",
		report.Pages, report.Read, report.Created, report.Changed, report.Idempotent)
	return nil
}
