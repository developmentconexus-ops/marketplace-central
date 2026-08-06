package arch_test

import (
	"os"
	"path/filepath"
	"testing"

	"marketplace-central/apps/server_core/internal/arch"
)

// contextsRoot is internal/contexts. Until Task 7 it does not exist, and a test
// that passes because it looked at nothing must say so out loud.
func contextsRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "contexts")
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		t.Skip("internal/contexts does not exist yet; nothing was scanned")
	}
	return root
}

func report(t *testing.T, got arch.Findings, what string) {
	t.Helper()
	for _, f := range got {
		t.Errorf("%s:%d %s: %s", f.File, f.Line, f.Rule, f.Detail)
	}
	if len(got) != 0 {
		t.Fatalf("%d %s", len(got), what)
	}
}

func TestNoContextImportsAnotherContextInternal(t *testing.T) {
	got, err := arch.ScanCrossContextInternal(contextsRoot(t))
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "cross-context internal imports")
}

func TestNoFloatInAnyContract(t *testing.T) {
	got, err := arch.ScanFloatInContracts(contextsRoot(t))
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "float fields in published contracts")
}

func TestNoVendorTokenInsideContexts(t *testing.T) {
	got, err := arch.ScanVendorTokens(contextsRoot(t), arch.VendorTokens)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "vendor tokens inside contexts/")
}

// The kernel exists from Task 1, so this one never skips.
func TestNoVendorTokenInKernel(t *testing.T) {
	got, err := arch.ScanVendorTokens(filepath.Join("..", "kernel"), arch.VendorTokens)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "vendor tokens in the kernel")
}
