package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

// These replace the three TestRequireTenantConfigured_* tests that pinned the
// old shape. That guard existed only because MC_DEFAULT_TENANT_ID arrived
// through the environment, where pgdb.LoadConfig substitutes "tenant_default"
// for an absent value and the two are indistinguishable downstream. Tenant
// and Sankhya instance are invocation parameters — what this run is for — so
// they arrive as required flags, and a required flag has no default to be
// mistaken for a real value. The old tests are not deleted silently: their
// subject (fail-closed on an unnamed tenant) is carried by
// TestParseOptionsRequiresTenant and TestParseOptionsIgnoresEnvironmentTenant
// below, which assert the same property against the new channel.

func TestParseOptionsRequiresTenant(t *testing.T) {
	if _, err := parseOptions([]string{"-instance", "prod"}, io.Discard); err == nil {
		t.Fatal("parseOptions without -tenant: error = nil, want error naming -tenant")
	} else if !strings.Contains(err.Error(), "-tenant") {
		t.Fatalf("error = %q, want it to name -tenant", err.Error())
	}
}

func TestParseOptionsRequiresInstance(t *testing.T) {
	if _, err := parseOptions([]string{"-tenant", "t1"}, io.Discard); err == nil {
		t.Fatal("parseOptions without -instance: error = nil, want error naming -instance")
	} else if !strings.Contains(err.Error(), "-instance") {
		t.Fatalf("error = %q, want it to name -instance", err.Error())
	}
}

// The next three are the anti-regression assertions for the class fix: with a
// retired variable exported, the command must still refuse to take its value
// from there. A live catalogue write under a tenant nobody named is exactly
// the failure D-39 recorded.
//
// One test per field, and each supplies the OTHER required flag correctly.
// Exporting both variables and passing no flags at all would be vacuous: a
// reintroduced tenant fallback would be masked by the still-missing
// -instance, and the test would stay green through the regression it claims
// to catch. For the same reason each asserts the error names its own flag,
// not merely that some error came back.

func TestParseOptionsIgnoresEnvironmentTenant(t *testing.T) {
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_from_environment")

	_, err := parseOptions([]string{"-instance", "prod"}, io.Discard)
	if err == nil {
		t.Fatal("error = nil with MC_DEFAULT_TENANT_ID exported; the environment must not supply the tenant")
	}
	if !strings.Contains(err.Error(), "-tenant") {
		t.Fatalf("error = %q, want it to name -tenant: another check fired instead, so this test would not see a tenant fallback", err.Error())
	}
}

func TestParseOptionsIgnoresEnvironmentInstance(t *testing.T) {
	t.Setenv("MPC_SANKHYA_INSTANCE", "instance_from_environment")

	_, err := parseOptions([]string{"-tenant", "t1"}, io.Discard)
	if err == nil {
		t.Fatal("error = nil with MPC_SANKHYA_INSTANCE exported; the environment must not supply the instance")
	}
	if !strings.Contains(err.Error(), "-instance") {
		t.Fatalf("error = %q, want it to name -instance: another check fired instead, so this test would not see an instance fallback", err.Error())
	}
}

// The page size variable is gone rather than replaced, so nothing else would
// notice it coming back: an exported value must not move the default.
func TestParseOptionsIgnoresEnvironmentPageSize(t *testing.T) {
	t.Setenv("MPC_CATALOG_INGEST_PAGE_SIZE", "9999")

	opts, err := parseOptions([]string{"-tenant", "t1", "-instance", "prod"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.pageSize != defaultPageSize {
		t.Fatalf("pageSize = %d with MPC_CATALOG_INGEST_PAGE_SIZE=9999 exported, want the flag default %d", opts.pageSize, defaultPageSize)
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
	if !strings.Contains(usage.String(), "-tenant") {
		t.Fatalf("usage = %q, want it to describe -tenant", usage.String())
	}
	if n := strings.Count(usage.String(), "-instance"); n != 1 {
		t.Fatalf("usage names -instance %d times, want exactly 1 (the FlagSet's own printer must stay silenced)", n)
	}
}

func TestParseOptionsRejectsBlankValues(t *testing.T) {
	for _, args := range [][]string{
		{"-tenant", "   ", "-instance", "prod"},
		{"-tenant", "t1", "-instance", ""},
	} {
		if _, err := parseOptions(args, io.Discard); err == nil {
			t.Fatalf("parseOptions(%q) error = nil, want error", args)
		}
	}
}

func TestParseOptionsDefaultsPageSize(t *testing.T) {
	opts, err := parseOptions([]string{"-tenant", "tenant_default", "-instance", "prod"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.tenant != "tenant_default" {
		t.Fatalf("tenant = %q, want tenant_default", opts.tenant)
	}
	if opts.instance != "prod" {
		t.Fatalf("instance = %q, want prod", opts.instance)
	}
	if opts.pageSize != defaultPageSize {
		t.Fatalf("pageSize = %d, want %d", opts.pageSize, defaultPageSize)
	}
}

func TestParseOptionsAcceptsPageSize(t *testing.T) {
	opts, err := parseOptions([]string{"-tenant", "t1", "-instance", "prod", "-page-size", "500"}, io.Discard)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.pageSize != 500 {
		t.Fatalf("pageSize = %d, want 500", opts.pageSize)
	}
}

func TestParseOptionsRejectsNonPositivePageSize(t *testing.T) {
	for _, raw := range []string{"0", "-1", "abc"} {
		if _, err := parseOptions([]string{"-tenant", "t1", "-instance", "p", "-page-size", raw}, io.Discard); err == nil {
			t.Fatalf("parseOptions(-page-size %q) error = nil, want error", raw)
		}
	}
}

func TestParseOptionsRejectsUnknownFlagAndPositionalArgs(t *testing.T) {
	for _, args := range [][]string{
		{"-tenant", "t1", "-instance", "p", "-pagesize", "10"},
		{"-tenant", "t1", "-instance", "p", "leftover"},
	} {
		if _, err := parseOptions(args, io.Discard); err == nil {
			t.Fatalf("parseOptions(%q) error = nil, want error", args)
		}
	}
}
