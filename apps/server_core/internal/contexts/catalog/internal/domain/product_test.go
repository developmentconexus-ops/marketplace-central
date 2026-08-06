package domain_test

import (
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
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

// obs builds an observation with the given description and payload hash.
func obs(t *testing.T, description, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	ean, err := contracts.NewIdentifier(contracts.IdentifierEAN, "7891234567890")
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{
		Key:         key,
		Description: desc,
		Identifiers: []contracts.Identifier{ean},
		Evidence:    e,
	}
}

func pid(t *testing.T) domain.ProductID {
	t.Helper()
	id, err := domain.NewProductID("prd_0198f65fc4477d4f9e9710c2b8ff922a")
	if err != nil {
		t.Fatalf("domain.NewProductID: %v", err)
	}
	return id
}

func TestNewProductIDRejectsBlank(t *testing.T) {
	if _, err := domain.NewProductID("  "); err == nil {
		t.Fatal("NewProductID accepted a blank id")
	}
}

// The canonical id must carry no source semantics. A caller that can pass the
// ERP code as the product id has already lost the property this type exists for.
func TestNewProductIDRejectsABareSourceCode(t *testing.T) {
	if _, err := domain.NewProductID("10529"); err == nil {
		t.Fatal("NewProductID accepted a bare numeric source code as a canonical id")
	}
}

func TestNewProductStartsAtVersionOne(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	if p.Version() != 1 {
		t.Fatalf("Version() = %d, want 1", p.Version())
	}
	if p.ID().String() != "prd_0198f65fc4477d4f9e9710c2b8ff922a" {
		t.Fatalf("ID() = %q", p.ID().String())
	}
	if p.Tenant().String() != "tnt_7f3b2" {
		t.Fatalf("Tenant() = %q, want tnt_7f3b2", p.Tenant().String())
	}
	if p.LastPayloadHash() != "sha256:ab91" {
		t.Fatalf("LastPayloadHash() = %q", p.LastPayloadHash())
	}
	if len(p.SourceKeys()) != 1 {
		t.Fatalf("SourceKeys() has %d entries, want 1", len(p.SourceKeys()))
	}
}

func TestNewProductRejectsAnInvalidObservation(t *testing.T) {
	bad := obs(t, "x", "sha256:ab91")
	bad.Evidence = provenance.Evidence{}
	if _, err := domain.NewProduct(pid(t), bad); err == nil {
		t.Fatal("NewProduct accepted an observation with no evidence")
	}
}

// Re-polling must be free. Same payload hash, same everything: no new version.
func TestApplySamePayloadIsIdempotent(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	next, disp, err := p.Apply(obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if disp != contracts.DispositionIdempotent {
		t.Fatalf("Disposition = %q, want %q", disp, contracts.DispositionIdempotent)
	}
	if next.Version() != 1 {
		t.Fatalf("Version() = %d after an idempotent apply, want 1", next.Version())
	}
}

func TestApplyChangedPayloadBumpsVersion(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	next, disp, err := p.Apply(obs(t, "Cafeteira Eletrica Inox 30 Xicaras 220 V", "sha256:cf02"))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if disp != contracts.DispositionChanged {
		t.Fatalf("Disposition = %q, want %q", disp, contracts.DispositionChanged)
	}
	if next.Version() != 2 {
		t.Fatalf("Version() = %d, want 2", next.Version())
	}
	desc, ok := next.Description().Value()
	if !ok || !strings.Contains(desc, "Inox") {
		t.Fatalf("Description() = %q, ok=%v; want the new description", desc, ok)
	}
}

// Apply returns a new value. The receiver must not move, or a caller that
// persists the old one after a failed transaction writes a version that the
// database never saw.
func TestApplyDoesNotMutateTheReceiver(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	if _, _, err := p.Apply(obs(t, "changed", "sha256:cf02")); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if p.Version() != 1 {
		t.Fatalf("receiver Version() = %d after Apply, want 1", p.Version())
	}
	desc, _ := p.Description().Value()
	if desc != "original" {
		t.Fatalf("receiver Description() = %q, want %q", desc, "original")
	}
}

// A second source key for the same canonical product is added, not replaced.
func TestApplyAccumulatesSourceKeys(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	second := obs(t, "original", "sha256:cf02")
	second.Key, err = contracts.NewSourceProductKey(tid(t), "spreadsheet", "import-2026-08", "product", "ROW-42")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	next, _, err := p.Apply(second)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(next.SourceKeys()) != 2 {
		t.Fatalf("SourceKeys() has %d entries, want 2", len(next.SourceKeys()))
	}
}

// An observation belonging to another tenant is not a version, it is a bug.
func TestApplyRejectsAForeignTenant(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	other, err := tenant.Parse("tnt_other")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	foreign := obs(t, "original", "sha256:cf02")
	foreign.Key, err = contracts.NewSourceProductKey(other, "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	if _, _, err := p.Apply(foreign); err == nil {
		t.Fatal("Apply accepted an observation from another tenant")
	}
}
