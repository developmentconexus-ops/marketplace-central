package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

// The class fix these tests pin: invocation parameters travel as arguments,
// not as process environment. "Which tenant" and "how big a page" are things
// the caller decides per run; MC_DATABASE_URL and MPC_ENCRYPTION_KEY are
// things the deployment decides once. Only the latter stay in the
// environment, read by the one registered typed loader (pgdb.LoadConfig).
//
// The old shape read MC_DEFAULT_TENANT_ID directly here so it could tell
// "unset" from "legitimately set to tenant_default" — pgdb.LoadConfig
// substitutes the latter for the former. That guard was a symptom of the
// wrong channel, not a feature: a required flag has no default to be confused
// with, so fail-closed is by construction and there is nothing left to guard.

func TestParseOptionsRequiresTenant(t *testing.T) {
	if _, err := parseOptions(nil, io.Discard); err == nil {
		t.Fatal("parseOptions(nil) error = nil, want error naming -tenant")
	} else if !strings.Contains(err.Error(), "-tenant") {
		t.Fatalf("parseOptions(nil) error = %q, want it to name -tenant", err.Error())
	}
}

// The next two are the anti-regression assertions for the whole refactor:
// with a retired variable exported, the command must still refuse to take its
// value from there. Each asserts the error names its own flag rather than
// merely that some error came back, so a different check firing first cannot
// mask the regression it is looking for.
func TestParseOptionsIgnoresEnvironmentTenant(t *testing.T) {
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_from_environment")

	_, err := parseOptions(nil, io.Discard)
	if err == nil {
		t.Fatal("error = nil with MC_DEFAULT_TENANT_ID exported; the environment must not supply the tenant")
	}
	if !strings.Contains(err.Error(), "-tenant") {
		t.Fatalf("error = %q, want it to name -tenant", err.Error())
	}
}

// The page size variable is gone rather than replaced, so nothing else would
// notice it coming back: an exported value must not move the default.
func TestParseOptionsIgnoresEnvironmentPageSize(t *testing.T) {
	t.Setenv("MPC_LISTINGS_INGEST_PAGE_SIZE", "9999")

	opts, err := parseOptions([]string{"-tenant", "t1"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.pageSize != defaultPageSize {
		t.Fatalf("pageSize = %d with MPC_LISTINGS_INGEST_PAGE_SIZE=9999 exported, want the flag default %d", opts.pageSize, defaultPageSize)
	}
}

// -h is a request, not a failure: parseOptions returns flag.ErrHelp (which
// main turns into a silent exit 0) and writes the usage text exactly once, to
// the writer it was handed rather than to the FlagSet's default stderr.
func TestParseOptionsHelpIsNotAFailure(t *testing.T) {
	var usage bytes.Buffer

	_, err := parseOptions([]string{"-h"}, &usage)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("parseOptions(-h) error = %v, want flag.ErrHelp", err)
	}
	if n := strings.Count(usage.String(), "-tenant"); n != 1 {
		t.Fatalf("usage names -tenant %d times, want exactly 1 (the FlagSet's own printer must stay silenced)", n)
	}
}

func TestParseOptionsRejectsBlankTenant(t *testing.T) {
	for _, arg := range []string{"", "   "} {
		if _, err := parseOptions([]string{"-tenant", arg}, io.Discard); err == nil {
			t.Fatalf("parseOptions(-tenant %q) error = nil, want error", arg)
		}
	}
}

func TestParseOptionsDefaultsPageSize(t *testing.T) {
	opts, err := parseOptions([]string{"-tenant", "tenant_default"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.tenant != "tenant_default" {
		t.Fatalf("tenant = %q, want tenant_default", opts.tenant)
	}
	if opts.pageSize != defaultPageSize {
		t.Fatalf("pageSize = %d, want %d", opts.pageSize, defaultPageSize)
	}
}

func TestParseOptionsAcceptsPageSize(t *testing.T) {
	opts, err := parseOptions([]string{"-tenant", "t1", "-page-size", "13"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.pageSize != 13 {
		t.Fatalf("pageSize = %d, want 13", opts.pageSize)
	}
}

// A page size of zero or below is not a smaller batch, it is a walk that
// never advances — the caller has to be told, not silently corrected.
func TestParseOptionsRejectsNonPositivePageSize(t *testing.T) {
	for _, raw := range []string{"0", "-1", "abc"} {
		if _, err := parseOptions([]string{"-tenant", "t1", "-page-size", raw}, io.Discard); err == nil {
			t.Fatalf("parseOptions(-page-size %q) error = nil, want error", raw)
		}
	}
}

func TestParseOptionsRejectsUnknownFlagAndPositionalArgs(t *testing.T) {
	for _, args := range [][]string{
		{"-tenant", "t1", "-pagesize", "10"},
		{"-tenant", "t1", "leftover"},
	} {
		if _, err := parseOptions(args, io.Discard); err == nil {
			t.Fatalf("parseOptions(%q) error = nil, want error", args)
		}
	}
}
