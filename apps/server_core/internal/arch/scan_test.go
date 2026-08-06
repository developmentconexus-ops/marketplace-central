package arch_test

import (
	"path/filepath"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/arch"
)

const fixtureSuffix = ".go.txt"

func violations() string { return filepath.Join("testdata", "violations") }
func clean() string      { return filepath.Join("testdata", "clean") }

func TestCrossContextInternalFiresOnFixture(t *testing.T) {
	got, err := arch.ScanCrossContextInternalSuffix(violations(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1: %v", len(got), got)
	}
	if got[0].Rule != arch.RuleCrossContextInternal {
		t.Fatalf("rule = %q, want %q", got[0].Rule, arch.RuleCrossContextInternal)
	}
	if got[0].Line != 4 {
		t.Fatalf("line = %d, want 4 (the beta/internal/domain import)", got[0].Line)
	}
}

func TestCrossContextInternalIsSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanCrossContextInternalSuffix(clean(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}

func TestFloatInContractsFiresOnFixture(t *testing.T) {
	got, err := arch.ScanFloatInContractsSuffix(violations(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("found %d findings, want 2 (float64 and float32): %v", len(got), got)
	}
}

func TestFloatInContractsIsSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanFloatInContractsSuffix(clean(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}

// The detector must see the token in a string literal, not only in identifiers.
// The measured defect IS a literal: orders/application/ingest_service.go:334
// compares the string "mercado_livre" in the application layer. A detector that
// only reads identifiers walks straight past it.
func TestVendorTokenFiresOnStringLiteral(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(violations(), arch.VendorTokens, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1: %v", len(got), got)
	}
	if got[0].Rule != arch.RuleVendorToken {
		t.Fatalf("rule = %q, want %q", got[0].Rule, arch.RuleVendorToken)
	}
	if !strings.Contains(got[0].Detail, "string literal") {
		t.Fatalf("detail = %q, want it to name the string literal", got[0].Detail)
	}
}

func TestVendorTokenFiresOnIdentifier(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(violations(), []string{"channelisdefault"}, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1 identifier hit: %v", len(got), got)
	}
	if !strings.Contains(got[0].Detail, "identifier") {
		t.Fatalf("detail = %q, want it to name the identifier", got[0].Detail)
	}
}

func TestVendorTokensAreSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(clean(), arch.VendorTokens, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}

// TestVendorTokensIgnoresAdaptersAndOwnList pins the two halves of Regra 2.3:
// a vendor name inside adapters/ is the design, and the detector's own token
// list is not a violation of the rule it implements.
func TestVendorTokensIgnoresAdaptersAndOwnList(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix("testdata", arch.VendorTokens, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	for _, f := range got {
		if strings.Contains(f.File, "/adapters/") {
			t.Fatalf("test=TestVendorTokensIgnoresAdaptersAndOwnList: adapters/ must be exempt, got %s:%d", f.File, f.Line)
		}
	}
	// Positive control: the fixture outside adapters/ must still be caught.
	var caught bool
	for _, f := range got {
		if strings.HasSuffix(f.File, "vendor_in_context.go.txt") {
			caught = true
		}
	}
	if !caught {
		t.Fatalf("test=TestVendorTokensIgnoresAdaptersAndOwnList: positive control not caught; got %+v", got)
	}
}

// TestCrossContextInternalSeesOutsideContexts is the positive control for the
// blindness that shipped: the composition root imported catalog/internal/postgres
// and the detector reported zero findings, because it skipped every file that was
// not itself inside contexts/.
func TestCrossContextInternalSeesOutsideContexts(t *testing.T) {
	got, err := arch.ScanCrossContextInternalSuffix("testdata", ".go.txt")
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	var hit *arch.Finding
	for i := range got {
		if strings.HasSuffix(got[i].File, "outside_context.go.txt") {
			hit = &got[i]
			break
		}
	}
	if hit == nil {
		t.Fatalf("test=TestCrossContextInternalSeesOutsideContexts: no finding for a file outside contexts/ importing catalog/internal; got %d finding(s): %+v", len(got), got)
	}
	if hit.Rule != arch.RuleCrossContextInternal {
		t.Fatalf("rule = %q, want %q", hit.Rule, arch.RuleCrossContextInternal)
	}
}

// TestFactValueDiscardIsCaught is the level-2 instrument for Regra 4.2. The
// package that exists to stop "unknown becomes zero" shipped exactly one call
// site doing `v, _ := f.Value()`, and nothing looked.
func TestFactValueDiscardIsCaught(t *testing.T) {
	got, err := arch.ScanFactValueDiscardSuffix("testdata", ".go.txt")
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("test=TestFactValueDiscardIsCaught: want exactly 1 finding, got %d: %+v", len(got), got)
	}
	if !strings.HasSuffix(got[0].File, "discarded_value.go.txt") || got[0].Rule != arch.RuleFactValueDiscard {
		t.Fatalf("finding = %+v", got[0])
	}
}
