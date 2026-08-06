package tenant_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func TestParseRejectsEmpty(t *testing.T) {
	if _, err := tenant.Parse(""); err == nil {
		t.Fatal("Parse(\"\") returned no error; empty tenant must be rejected")
	}
}

func TestParseRejectsWhitespaceOnly(t *testing.T) {
	if _, err := tenant.Parse("   "); err == nil {
		t.Fatal("Parse(\"   \") returned no error; blank tenant must be rejected")
	}
}

func TestParseRoundTrips(t *testing.T) {
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := id.String(); got != "tnt_7f3b2" {
		t.Fatalf("String() = %q, want %q", got, "tnt_7f3b2")
	}
	if id.IsZero() {
		t.Fatal("IsZero() = true for a parsed id")
	}
}

func TestZeroValueIsZero(t *testing.T) {
	var id tenant.ID
	if !id.IsZero() {
		t.Fatal("zero ID.IsZero() = false")
	}
	if id.String() != "" {
		t.Fatalf("zero ID.String() = %q, want empty", id.String())
	}
}
