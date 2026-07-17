package application_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

type fakeMarketObservationReadStore struct {
	rows      []domain.MarketObservation
	called    int
	requested []string
}

func (f *fakeMarketObservationReadStore) AppendObservations(context.Context, []domain.MarketObservation) error {
	return nil
}

func (f *fakeMarketObservationReadStore) LatestObservations(_ context.Context, listingIDs []string) ([]domain.MarketObservation, error) {
	f.called++
	f.requested = append([]string(nil), listingIDs...)
	return f.rows, nil
}

type fakeMarketReferenceReadStore struct {
	rows      []domain.MarketReference
	called    int
	requested []string
}

func (f *fakeMarketReferenceReadStore) AppendReferences(context.Context, []domain.MarketReference) error {
	return nil
}

func (f *fakeMarketReferenceReadStore) LatestReferences(_ context.Context, productIDs []string) ([]domain.MarketReference, error) {
	f.called++
	f.requested = append([]string(nil), productIDs...)
	return f.rows, nil
}

func newMarketReadService(observations *fakeMarketObservationReadStore, references *fakeMarketReferenceReadStore) *application.ReadService {
	return application.NewReadService(observations, references, func() time.Time {
		return time.Date(2026, 7, 16, 18, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	})
}

func TestVirginObservationReadReturnsNoPriceEvidenceWithNullSignals(t *testing.T) {
	store := &fakeMarketObservationReadStore{}
	service := newMarketReadService(store, &fakeMarketReferenceReadStore{})

	got, err := service.ListObservations(context.Background(), []string{"listing-b", "listing-a"})
	if err != nil {
		t.Fatalf("ListObservations() error = %v", err)
	}
	if want := time.Date(2026, 7, 16, 21, 0, 0, 0, time.UTC); !got.AsOf.Equal(want) {
		t.Fatalf("AsOf = %v, want %v", got.AsOf, want)
	}
	if len(got.Items) != 2 || got.Items[0].ListingID != "listing-b" || got.Items[1].ListingID != "listing-a" {
		t.Fatalf("items = %+v, want input order", got.Items)
	}
	for _, item := range got.Items {
		if item.EvidenceState != domain.EvidenceStateNoPriceEvidence {
			t.Errorf("evidence state = %q, want %q", item.EvidenceState, domain.EvidenceStateNoPriceEvidence)
		}
		if item.OurSalePrice != nil || item.CompetitiveStatus != nil || item.WinnerPrice != nil ||
			item.CompetitiveTarget != nil || item.CatalogOfferPrice != nil || item.CatalogStats != nil ||
			item.CapturedAt != nil || item.Source != nil {
			t.Errorf("synthetic item = %+v, want all signals and provenance nil", item)
		}
	}
}

func TestMalformedListingIDReturnsItemLevelNoEvidence(t *testing.T) {
	store := &fakeMarketObservationReadStore{}
	service := newMarketReadService(store, &fakeMarketReferenceReadStore{})

	got, err := service.ListObservations(context.Background(), []string{"not a listing id", "listing-valid"})
	if err != nil {
		t.Fatalf("ListObservations() error = %v, want nil", err)
	}
	if got.Items[0].ListingID != "not a listing id" || got.Items[0].EvidenceState != domain.EvidenceStateNoPriceEvidence {
		t.Fatalf("malformed item = %+v, want item-level no evidence", got.Items[0])
	}
}

func TestUnknownProductIDReturnsNoEvidence(t *testing.T) {
	store := &fakeMarketReferenceReadStore{}
	service := newMarketReadService(&fakeMarketObservationReadStore{}, store)

	got, err := service.ListReferences(context.Background(), []string{"unknown-product"})
	if err != nil {
		t.Fatalf("ListReferences() error = %v, want nil", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(got.Items))
	}
	item := got.Items[0]
	if item.ProductID != "unknown-product" || item.MatchState != domain.MatchStateNoCandidate || item.EvidenceState != domain.EvidenceStateNoPriceEvidence {
		t.Fatalf("synthetic reference = %+v, want no-candidate/no-evidence", item)
	}
	if item.CatalogProductID != nil || item.MatchMethod != nil || item.CatalogStats != nil || item.CapturedAt != nil || item.Source != nil {
		t.Fatalf("synthetic reference = %+v, want nullable fields nil", item)
	}
}

func TestStoredSignalsRoundTripWithoutMerge(t *testing.T) {
	capturedAt := time.Date(2026, 7, 16, 17, 30, 0, 0, time.UTC)
	source := domain.SourceManual
	status := domain.CompetitiveStatusCompeting
	stored := domain.MarketObservation{
		ListingID:         "listing-stored",
		OurSalePrice:      &domain.Money{Amount: "101.01", Currency: "BRL"},
		CompetitiveStatus: &status,
		WinnerPrice:       &domain.Money{Amount: "202.02", Currency: "BRL"},
		CompetitiveTarget: &domain.Money{Amount: "303.03", Currency: "BRL"},
		CatalogOfferPrice: &domain.Money{Amount: "404.04", Currency: "BRL"},
		CatalogStats: &domain.CatalogStats{
			Min: "401.01", P25: "402.02", Median: "403.03", TrimmedMean: "404.04", P75: "405.05", Max: "406.06",
			NOffers: 8, NSellers: 5,
		},
		EvidenceState: domain.EvidenceStateObserved,
		CapturedAt:    &capturedAt,
		Source:        &source,
	}
	store := &fakeMarketObservationReadStore{rows: []domain.MarketObservation{stored}}
	service := newMarketReadService(store, &fakeMarketReferenceReadStore{})

	got, err := service.ListObservations(context.Background(), []string{"listing-stored"})
	if err != nil {
		t.Fatalf("ListObservations() error = %v", err)
	}
	if len(got.Items) != 1 || !reflect.DeepEqual(got.Items[0], stored) {
		t.Fatalf("item = %+v, want stored observation unchanged", got.Items)
	}
}

func TestTwoHundredOneIDsReturnsTooManyIDs(t *testing.T) {
	observations := &fakeMarketObservationReadStore{}
	service := newMarketReadService(observations, &fakeMarketReferenceReadStore{})

	ids := make([]string, 200)
	if _, err := service.ListObservations(context.Background(), ids); err != nil {
		t.Fatalf("ListObservations(200 ids) error = %v, want nil", err)
	}
	if observations.called != 1 {
		t.Fatalf("store calls after 200 ids = %d, want 1", observations.called)
	}

	ids = append(ids, "id-201")
	_, err := service.ListObservations(context.Background(), ids)
	var tooMany *application.TooManyIDsError
	if !errors.As(err, &tooMany) || tooMany.Code() != "too_many_ids" {
		t.Fatalf("error = %v, want typed too_many_ids", err)
	}
	if observations.called != 1 {
		t.Fatalf("store calls after 201 ids = %d, want unchanged at 1", observations.called)
	}
}
