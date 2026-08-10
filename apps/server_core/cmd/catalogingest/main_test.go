package main

import (
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

// TestParseOptionsIgnoresEnvironmentTenant is the anti-regression assertion
// for the class fix: with the retired variables exported and the flags
// absent, this command must still refuse to run. A live catalogue write under
// a tenant nobody named is exactly the failure D-39 recorded.
func TestParseOptionsIgnoresEnvironmentTenant(t *testing.T) {
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_from_environment")
	t.Setenv("MPC_SANKHYA_INSTANCE", "instance_from_environment")

	if _, err := parseOptions(nil, io.Discard); err == nil {
		t.Fatal("parseOptions(nil) error = nil with the retired variables exported; the environment must not supply invocation parameters")
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
