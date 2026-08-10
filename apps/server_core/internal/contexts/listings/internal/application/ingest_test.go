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
		Key:        key,
		Evidence:   ev,
		RawPayload: []byte(`{"id":"` + listingID + `"}`),
		State: contracts.ListingState{
			Title:             unknown("test fixture omitted title"),
			Status:            unknown("test fixture omitted status"),
			ListingType:       unknown("test fixture omitted listing_type"),
			Price:             unknownMoney("test fixture omitted price"),
			AvailableQuantity: unknownInt("test fixture omitted available_quantity"),
			SellerSKU:         unknown("test fixture omitted seller_sku"),
			GTIN:              unknown("test fixture omitted gtin"),
		},
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

// TestChangedPayloadMintsNewVersion used to prove a changed hash alone minted
// a new version. That is no longer the rule (ingest.go now folds on
// State.Equal, not the hash — see TestIdenticalFactsDifferentEvidenceStaysIdempotent
// for the case where the hash changes and nothing else does), so the second
// observation here now also carries a changed FACT: the assertions are
// unchanged, but the fixture that makes them true is a genuine fact change,
// which is what "changed v2" actually names under the new fold.
func TestChangedPayloadMintsNewVersion(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	second := obs(t, "MLB1", "h2")
	known, err := fact.NewKnown("paused", second.Evidence)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	second.State.Status = known
	got, err := svc.Ingest(context.Background(), second)
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionChanged || got.Version != 2 {
		t.Fatalf("got %+v, want changed v2", got)
	}
}

// TestSameEvidenceCorrectedFactMintsNewVersion is the defect that started this
// change: a corrected mapper produces the SAME evidence (the channel's bytes
// never moved, so the payload hash is identical) but a DIFFERENT fact --
// SellerSKU goes from Unknown to Known. Under the old payload-hash fold this
// read as idempotent, which is exactly why re-running the live ingest after
// today's mapper fix healed only 14 of 34 rows: the other 20 listings' bytes
// had not changed at Mercado Livre, so the fold never looked at the fact that
// did change.
func TestSameEvidenceCorrectedFactMintsNewVersion(t *testing.T) {
	store := newMemStore()
	svc := application.NewService(store)
	first := obs(t, "MLB1", "h1")
	if _, err := svc.Ingest(context.Background(), first); err != nil {
		t.Fatalf("first ingest: %v", err)
	}

	second := first
	known, err := fact.NewKnown("SKU-MLB1", first.Evidence)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	second.State.SellerSKU = known

	got, err := svc.Ingest(context.Background(), second)
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionChanged || got.Version != 2 {
		t.Fatalf("got %+v, want changed v2 -- a corrected fact must mint a new version even though the payload hash did not move", got)
	}
	sku, ok := store.lastObservation.State.SellerSKU.Value()
	if !ok || sku != "SKU-MLB1" {
		t.Fatalf("stored seller_sku = %q known=%v, want the corrected value persisted", sku, ok)
	}
}

// TestIdenticalFactsDifferentEvidenceStaysIdempotent is the converse of the
// healing property above: re-polling an unchanged channel must stay free even
// when the evidence itself differs (a later timestamp, a different payload
// hash) as long as the facts derived from it are the same. A fold that keyed
// on the hash alone would mint a version here for no reason a reader can see.
func TestIdenticalFactsDifferentEvidenceStaysIdempotent(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h2"))
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionIdempotent || got.Version != 1 {
		t.Fatalf("got %+v, want idempotent v1 -- same facts, different evidence, must not mint a version", got)
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
