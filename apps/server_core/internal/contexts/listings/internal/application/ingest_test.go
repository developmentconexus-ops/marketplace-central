package application_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func obs(t *testing.T, listingID, payloadHash string) contracts.ListingObservation {
	t.Helper()
	tid, _ := tenant.Parse("tenant_default")
	code, _ := channel.ParseCode("mercado_livre")
	account, _ := channel.NewAccountRef(code, "179571326")
	key, err := contracts.NewSourceListingKey(tid, account, listingID)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	ev, err := provenance.NewEvidence("mercado_livre", "item", listingID, time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC), payloadHash)
	if err != nil {
		t.Fatalf("evidence: %v", err)
	}
	// This suite tests Ingest's disposition/version bookkeeping, not what the
	// channel reported, so every fact is a deliberate Unknown-with-reason
	// rather than a zero value — Validate now rejects a zero-value fact
	// (contracts/observation.go), same as the database's CHECK constraints do.
	unknown := func(reason string) fact.Fact[string] {
		f, ferr := fact.NewUnknown[string](reason, ev)
		if ferr != nil {
			t.Fatalf("fact.NewUnknown: %v", ferr)
		}
		return f
	}
	unknownInt := func(reason string) fact.Fact[int] {
		f, ferr := fact.NewUnknown[int](reason, ev)
		if ferr != nil {
			t.Fatalf("fact.NewUnknown: %v", ferr)
		}
		return f
	}
	unknownMoney := func(reason string) fact.Fact[exact.Money] {
		f, ferr := fact.NewUnknown[exact.Money](reason, ev)
		if ferr != nil {
			t.Fatalf("fact.NewUnknown: %v", ferr)
		}
		return f
	}
	return contracts.ListingObservation{
		Key:               key,
		Evidence:          ev,
		RawPayload:        []byte(`{"id":"` + listingID + `"}`),
		Title:             unknown("test fixture omitted title"),
		Status:            unknown("test fixture omitted status"),
		ListingType:       unknown("test fixture omitted listing_type"),
		Price:             unknownMoney("test fixture omitted price"),
		AvailableQuantity: unknownInt("test fixture omitted available_quantity"),
		SellerSKU:         unknown("test fixture omitted seller_sku"),
		GTIN:              unknown("test fixture omitted gtin"),
	}
}

func TestFirstObservationCreatesVersionOne(t *testing.T) {
	svc := application.NewService(newMemStore())
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1"))
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionCreated || got.Version != 1 {
		t.Fatalf("got %+v, want created v1", got)
	}
}

func TestSamePayloadHashIsIdempotent(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1"))
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionIdempotent || got.Version != 1 {
		t.Fatalf("got %+v, want idempotent v1", got)
	}
}

func TestChangedPayloadMintsNewVersion(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h2"))
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionChanged || got.Version != 2 {
		t.Fatalf("got %+v, want changed v2", got)
	}
}

func TestInvalidObservationNeverReachesTheStore(t *testing.T) {
	store := newMemStore()
	svc := application.NewService(store)
	bad := contracts.ListingObservation{} // zero key, zero evidence
	if _, err := svc.Ingest(context.Background(), bad); err == nil {
		t.Fatal("invalid observation accepted")
	}
	if store.saves != 0 {
		t.Fatalf("store touched %d times by an invalid observation", store.saves)
	}
}
