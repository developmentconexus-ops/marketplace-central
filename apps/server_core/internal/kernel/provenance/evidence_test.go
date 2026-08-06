package provenance_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

func at() time.Time { return time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC) }

func TestNewEvidenceRejectsMissingParts(t *testing.T) {
	cases := []struct {
		name     string
		system   string
		kind     string
		key      string
		observed time.Time
		hash     string
	}{
		{"no system", "", "product", "10529", at(), "sha256:ab91"},
		{"no kind", "sankhya", "", "10529", at(), "sha256:ab91"},
		{"no key", "sankhya", "product", "", at(), "sha256:ab91"},
		{"zero time", "sankhya", "product", "10529", time.Time{}, "sha256:ab91"},
		{"no hash", "sankhya", "product", "10529", at(), ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := provenance.NewEvidence(c.system, c.kind, c.key, c.observed, c.hash); err == nil {
				t.Fatalf("NewEvidence accepted %s", c.name)
			}
		})
	}
}

func TestEvidenceRefIsStable(t *testing.T) {
	e, err := provenance.NewEvidence("sankhya", "product", "10529", at(), "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	if got := e.Ref(); got != "sankhya/product:10529" {
		t.Fatalf("Ref() = %q, want %q", got, "sankhya/product:10529")
	}
	if !e.ObservedAt().Equal(at()) {
		t.Fatalf("ObservedAt() = %v, want %v", e.ObservedAt(), at())
	}
}

func TestObservedAtIsStoredInUTC(t *testing.T) {
	saoPaulo := time.FixedZone("-03", -3*60*60)
	local := time.Date(2026, 8, 6, 9, 15, 0, 0, saoPaulo)
	e, err := provenance.NewEvidence("sankhya", "product", "10529", local, "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	if e.ObservedAt().Location() != time.UTC {
		t.Fatalf("ObservedAt().Location() = %v, want UTC", e.ObservedAt().Location())
	}
	if !e.ObservedAt().Equal(at()) {
		t.Fatalf("ObservedAt() = %v, want the same instant as %v", e.ObservedAt(), at())
	}
}

func TestZeroEvidenceIsZero(t *testing.T) {
	var e provenance.Evidence
	if !e.IsZero() {
		t.Fatal("zero Evidence.IsZero() = false")
	}
}
