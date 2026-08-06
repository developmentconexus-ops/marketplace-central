package application_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
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

// observation builds one product observation.
func observation(t *testing.T, externalKey, description, ean, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", externalKey,
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", externalKey)
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{
		Key: key, Description: desc,
		Identifiers: []contracts.Identifier{id}, Evidence: e,
	}
}

func newService() (*application.Service, *memStore) {
	store := newMemStore()
	return application.NewService(store, &seqIDs{}), store
}

func TestIngestFirstObservationCreates(t *testing.T) {
	svc, _ := newService()
	got, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica 30 Xicaras 220 V", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("Ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionCreated {
		t.Fatalf("Disposition = %q, want %q", got.Disposition, contracts.DispositionCreated)
	}
	if got.Version != 1 {
		t.Fatalf("Version = %d, want 1", got.Version)
	}
	if got.ProductID == "10529" {
		t.Fatal("the canonical id is the ERP code; identity leaked from the source")
	}
	if len(got.DuplicateIdentifiers) != 0 {
		t.Fatalf("DuplicateIdentifiers = %v, want none", got.DuplicateIdentifiers)
	}
}

func TestIngestSamePayloadIsIdempotent(t *testing.T) {
	svc, _ := newService()
	obs := observation(t, "10529", "Cafeteira", "7891234567890", "sha256:ab91")
	first, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionIdempotent {
		t.Fatalf("Disposition = %q, want %q", second.Disposition, contracts.DispositionIdempotent)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("ProductID moved between polls: %q then %q", first.ProductID, second.ProductID)
	}
	if second.Version != 1 {
		t.Fatalf("Version = %d after an idempotent poll, want 1", second.Version)
	}
}

func TestIngestChangedPayloadBumpsVersionAndKeepsIdentity(t *testing.T) {
	svc, _ := newService()
	first, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica 30 Xicaras 220 V", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica Inox 30 Xicaras 220 V", "7891234567890", "sha256:cf02"))
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionChanged {
		t.Fatalf("Disposition = %q, want %q", second.Disposition, contracts.DispositionChanged)
	}
	if second.Version != 2 {
		t.Fatalf("Version = %d, want 2", second.Version)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("a changed description minted a new identity: %q then %q", first.ProductID, second.ProductID)
	}
}

// The load-bearing test. ERP code 10530 reports the SAME EAN as 10529 because
// the master data is bad. Two source keys are two products. Merging them
// silently is irreversible; reporting the duplicate is not.
func TestIngestDuplicateEanDoesNotMerge(t *testing.T) {
	svc, _ := newService()
	first, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(),
		observation(t, "10530", "Cafeteira B", "7891234567890", "sha256:dd77"))
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionCreated {
		t.Fatalf("Disposition = %q, want %q; a duplicate EAN must not fold into the existing product",
			second.Disposition, contracts.DispositionCreated)
	}
	if second.ProductID == first.ProductID {
		t.Fatal("two ERP codes sharing a bad EAN were merged into one product")
	}
	if len(second.DuplicateIdentifiers) != 1 {
		t.Fatalf("DuplicateIdentifiers = %v, want exactly the shared EAN", second.DuplicateIdentifiers)
	}
	if second.DuplicateIdentifiers[0].Value() != "7891234567890" {
		t.Fatalf("DuplicateIdentifiers[0] = %q", second.DuplicateIdentifiers[0].Value())
	}
}

// Another tenant reporting the same EAN is not a duplicate. It is a different
// company's catalogue, and a conflict reported across tenants is a data leak.
func TestIngestDoesNotReportDuplicatesAcrossTenants(t *testing.T) {
	svc, _ := newService()
	if _, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ab91")); err != nil {
		t.Fatalf("first Ingest: %v", err)
	}

	other, err := tenant.Parse("tnt_other")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	obs := observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ee88")
	obs.Key, err = contracts.NewSourceProductKey(other, "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}

	got, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if len(got.DuplicateIdentifiers) != 0 {
		t.Fatalf("DuplicateIdentifiers = %v across tenants; that is a leak", got.DuplicateIdentifiers)
	}
}

func TestIngestRejectsAnInvalidObservation(t *testing.T) {
	svc, _ := newService()
	bad := observation(t, "10529", "Cafeteira", "7891234567890", "sha256:ab91")
	bad.Evidence = provenance.Evidence{}
	if _, err := svc.Ingest(context.Background(), bad); err == nil {
		t.Fatal("Ingest accepted an observation with no evidence")
	}
}
