package contracts_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func tid(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

func ev(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), "sha256:ab91")
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	return e
}

func TestNewIdentifierRejectsBlank(t *testing.T) {
	if _, err := contracts.NewIdentifier(contracts.IdentifierEAN, "   "); err == nil {
		t.Fatal("NewIdentifier with blank value returned no error")
	}
}

func TestNewIdentifierRejectsUnknownKind(t *testing.T) {
	if _, err := contracts.NewIdentifier(contracts.IdentifierKind("barcode"), "789"); err == nil {
		t.Fatal("NewIdentifier accepted a kind outside the closed set")
	}
}

func TestNewSourceProductKeyRequiresEveryPart(t *testing.T) {
	id := tid(t)
	cases := []struct{ name, system, instance, kind, key string }{
		{"no system", "", "sankhya-prod-01", "product", "10529"},
		{"no instance", "sankhya", "", "product", "10529"},
		{"no kind", "sankhya", "sankhya-prod-01", "", "10529"},
		{"no key", "sankhya", "sankhya-prod-01", "product", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := contracts.NewSourceProductKey(id, c.system, c.instance, c.kind, c.key); err == nil {
				t.Fatalf("NewSourceProductKey accepted %s", c.name)
			}
		})
	}
	var noTenant tenant.ID
	if _, err := contracts.NewSourceProductKey(noTenant, "sankhya", "sankhya-prod-01", "product", "10529"); err == nil {
		t.Fatal("NewSourceProductKey accepted a zero tenant")
	}
}

// The ERP code lives inside the source key and nowhere else. This is the test
// that stops CODPROD from becoming the canonical identifier a second time.
func TestSourceProductKeyCarriesTheErpCodeAndNothingElseDoes(t *testing.T) {
	k, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	want := "tnt_7f3b2/sankhya/sankhya-prod-01/product/10529"
	if got := k.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if k.ExternalKey() != "10529" {
		t.Fatalf("ExternalKey() = %q, want %q", k.ExternalKey(), "10529")
	}
}

// Two Sankhya installations are two sources even for the same product code.
func TestSourceInstanceDiscriminates(t *testing.T) {
	a, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	b, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-02", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	if a.String() == b.String() {
		t.Fatal("two installations produced the same key; the instance is not discriminating")
	}
}

func TestObservationValidateRequiresKeyAndEvidence(t *testing.T) {
	desc, err := fact.NewKnown("Cafeteira Eletrica 30 Xicaras 220 V", ev(t))
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	ean, err := contracts.NewIdentifier(contracts.IdentifierEAN, "7891234567890")
	if err != nil {
		t.Fatalf("NewIdentifier: %v", err)
	}

	obs := contracts.ProductObservation{
		Key:         key,
		Description: desc,
		Identifiers: []contracts.Identifier{ean},
		Evidence:    ev(t),
	}
	if err := obs.Validate(); err != nil {
		t.Fatalf("Validate on a complete observation: %v", err)
	}

	noEvidence := obs
	noEvidence.Evidence = provenance.Evidence{}
	if err := noEvidence.Validate(); err == nil {
		t.Fatal("Validate accepted an observation with no evidence")
	}

	noKey := obs
	noKey.Key = contracts.SourceProductKey{}
	if err := noKey.Validate(); err == nil {
		t.Fatal("Validate accepted an observation with no source key")
	}
}

// A source that returns an empty description has not told us the product has no
// name. It has told us nothing, and those are different — which is why the
// field is a Fact and not a string.
func TestObservationAcceptsAnUnknownDescription(t *testing.T) {
	desc, err := fact.NewUnknown[string]("source returned no description", ev(t))
	if err != nil {
		t.Fatalf("fact.NewUnknown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	obs := contracts.ProductObservation{Key: key, Description: desc, Evidence: ev(t)}
	if err := obs.Validate(); err != nil {
		t.Fatalf("Validate rejected a legitimately unknown description: %v", err)
	}
}
