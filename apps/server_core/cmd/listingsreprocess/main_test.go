package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

// This command's whole point is that it never talks to Mercado Livre, so its
// only inputs are the run's own parameters: which tenant to reprocess and how
// big a stored-observations page to request. Both travel as flags, mirroring
// cmd/listingsingest's parseOptions exactly — same fail-closed shape, just
// without the credential concerns a live poll needs and this command does
// not.

func TestParseOptionsRequiresTenant(t *testing.T) {
	if _, err := parseOptions(nil, io.Discard); err == nil {
		t.Fatal("parseOptions(nil) error = nil, want error naming -tenant")
	} else if !strings.Contains(err.Error(), "-tenant") {
		t.Fatalf("parseOptions(nil) error = %q, want it to name -tenant", err.Error())
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
		if raw == "0" || raw == "-1" {
			if _, err := parseOptions([]string{"-tenant", "t1", "-page-size", raw}, io.Discard); err != nil && !strings.Contains(err.Error(), "page-size") {
				t.Fatalf("parseOptions(-page-size %q) error = %q, want it to name -page-size", raw, err.Error())
			}
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

func TestParseOptionsUnexpectedPositionalArgumentNamesIt(t *testing.T) {
	_, err := parseOptions([]string{"-tenant", "t1", "leftover"}, io.Discard)
	if err == nil {
		t.Fatal("parseOptions error = nil, want error naming the unexpected argument")
	}
	if !strings.Contains(err.Error(), "leftover") {
		t.Fatalf("parseOptions error = %q, want it to name the unexpected argument %q", err.Error(), "leftover")
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
