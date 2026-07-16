package application_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

type fakeCollector struct {
	observations []domain.MarketObservation
	references   []domain.MarketReference
	err          error
}

func (f fakeCollector) CollectObservations(context.Context, string, []string) ([]domain.MarketObservation, error) {
	return f.observations, f.err
}

func (f fakeCollector) CollectReferences(context.Context, []string) ([]domain.MarketReference, error) {
	return f.references, f.err
}

type fakeObservationStore struct {
	appended []domain.MarketObservation
	err      error
}

func (f *fakeObservationStore) AppendObservations(_ context.Context, observations []domain.MarketObservation) error {
	f.appended = append(f.appended, observations...)
	return f.err
}

func (f *fakeObservationStore) LatestObservations(context.Context, []string) ([]domain.MarketObservation, error) {
	return nil, nil
}

type fakeReferenceStore struct {
	appended []domain.MarketReference
	err      error
}

func (f *fakeReferenceStore) AppendReferences(_ context.Context, references []domain.MarketReference) error {
	f.appended = append(f.appended, references...)
	return f.err
}

func (f *fakeReferenceStore) LatestReferences(context.Context, []string) ([]domain.MarketReference, error) {
	return nil, nil
}

func TestCollectorPortRoundTripPersistsMandatoryProvenance(t *testing.T) {
	capturedAt := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	source := domain.SourceOfficialAPI
	collector := fakeCollector{
		observations: []domain.MarketObservation{{
			ListingID:     "listing-1",
			EvidenceState: domain.EvidenceStateObserved,
			CapturedAt:    &capturedAt,
			Source:        &source,
		}},
		references: []domain.MarketReference{{
			ProductID:     "product-1",
			MatchState:    domain.MatchStateAccept,
			EvidenceState: domain.EvidenceStateObserved,
			CapturedAt:    &capturedAt,
			Source:        &source,
		}},
	}
	observations := &fakeObservationStore{}
	references := &fakeReferenceStore{}
	service := application.NewCollectionService(collector, observations, references)

	if err := service.CollectObservations(context.Background(), "installation-1", []string{"listing-1"}); err != nil {
		t.Fatalf("CollectObservations() error = %v", err)
	}
	if err := service.CollectReferences(context.Background(), []string{"product-1"}); err != nil {
		t.Fatalf("CollectReferences() error = %v", err)
	}

	if len(observations.appended) != 1 || observations.appended[0].Source == nil ||
		*observations.appended[0].Source != source || observations.appended[0].CapturedAt == nil ||
		!observations.appended[0].CapturedAt.Equal(capturedAt) {
		t.Fatalf("stored observations = %+v, want provenance intact", observations.appended)
	}
	if len(references.appended) != 1 || references.appended[0].Source == nil ||
		*references.appended[0].Source != source || references.appended[0].CapturedAt == nil ||
		!references.appended[0].CapturedAt.Equal(capturedAt) {
		t.Fatalf("stored references = %+v, want provenance intact", references.appended)
	}
}

func TestCollectorSignalsAreNeverMerged(t *testing.T) {
	capturedAt := time.Date(2026, 7, 16, 13, 0, 0, 0, time.UTC)
	source := domain.SourceVendor
	observation := domain.MarketObservation{
		ListingID:         "listing-2",
		OurSalePrice:      &domain.Money{Amount: "101.01", Currency: "BRL"},
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
	store := &fakeObservationStore{}
	service := application.NewCollectionService(fakeCollector{observations: []domain.MarketObservation{observation}}, store, &fakeReferenceStore{})

	if err := service.CollectObservations(context.Background(), "installation-1", []string{"listing-2"}); err != nil {
		t.Fatalf("CollectObservations() error = %v", err)
	}
	if len(store.appended) != 1 {
		t.Fatalf("stored observations = %d, want 1", len(store.appended))
	}
	got := store.appended[0]
	if !reflect.DeepEqual(got.OurSalePrice, observation.OurSalePrice) ||
		!reflect.DeepEqual(got.WinnerPrice, observation.WinnerPrice) ||
		!reflect.DeepEqual(got.CompetitiveTarget, observation.CompetitiveTarget) ||
		!reflect.DeepEqual(got.CatalogOfferPrice, observation.CatalogOfferPrice) ||
		!reflect.DeepEqual(got.CatalogStats, observation.CatalogStats) {
		t.Fatalf("stored signals = %+v, want each input signal unchanged", got)
	}
}

func TestCollectionRejectsInvalidObservationBatchWithoutPersisting(t *testing.T) {
	capturedAt := time.Date(2026, 7, 16, 14, 0, 0, 0, time.UTC)
	source := domain.SourceManual
	valid := domain.MarketObservation{
		ListingID: "listing-valid", EvidenceState: domain.EvidenceStateObserved,
		CapturedAt: &capturedAt, Source: &source,
	}
	invalid := valid
	invalid.ListingID = "listing-invalid"
	invalid.Source = nil
	store := &fakeObservationStore{}
	service := application.NewCollectionService(fakeCollector{observations: []domain.MarketObservation{valid, invalid}}, store, &fakeReferenceStore{})

	err := service.CollectObservations(context.Background(), "installation-1", []string{"listing-valid", "listing-invalid"})
	if !errors.Is(err, domain.ErrSourceRequired) {
		t.Fatalf("CollectObservations() error = %v, want %v", err, domain.ErrSourceRequired)
	}
	if len(store.appended) != 0 {
		t.Fatalf("stored observations = %d, want none", len(store.appended))
	}
}

func TestCollectionRejectsInvalidReferenceBatchWithoutPersisting(t *testing.T) {
	capturedAt := time.Date(2026, 7, 16, 15, 0, 0, 0, time.UTC)
	source := domain.SourceManual
	valid := domain.MarketReference{
		ProductID: "product-valid", MatchState: domain.MatchStateReview,
		EvidenceState: domain.EvidenceStateObserved, CapturedAt: &capturedAt, Source: &source,
	}
	invalid := valid
	invalid.ProductID = "product-invalid"
	invalid.CapturedAt = nil
	store := &fakeReferenceStore{}
	service := application.NewCollectionService(fakeCollector{references: []domain.MarketReference{valid, invalid}}, &fakeObservationStore{}, store)

	err := service.CollectReferences(context.Background(), []string{"product-valid", "product-invalid"})
	if !errors.Is(err, domain.ErrCapturedAtRequired) {
		t.Fatalf("CollectReferences() error = %v, want %v", err, domain.ErrCapturedAtRequired)
	}
	if len(store.appended) != 0 {
		t.Fatalf("stored references = %d, want none", len(store.appended))
	}
}

func TestCollectionPropagatesCollectorAndStoreErrors(t *testing.T) {
	want := errors.New("dependency failed")
	collector := &fakeCollector{err: want}
	observationStore := &fakeObservationStore{err: want}
	referenceStore := &fakeReferenceStore{err: want}
	service := application.NewCollectionService(collector, observationStore, referenceStore)

	if err := service.CollectObservations(context.Background(), "installation-1", nil); !errors.Is(err, want) {
		t.Fatalf("collector error = %v, want %v", err, want)
	}
	collector.err = nil
	collector.observations = []domain.MarketObservation{{ListingID: "listing-1", EvidenceState: domain.EvidenceStateObserved}}
	if err := service.CollectObservations(context.Background(), "installation-1", nil); !errors.Is(err, domain.ErrCapturedAtRequired) {
		t.Fatalf("validation error = %v, want %v", err, domain.ErrCapturedAtRequired)
	}

	collector.observations = nil
	referenceSource := domain.SourceManual
	collector.references = []domain.MarketReference{{
		ProductID: "product-1", MatchState: domain.MatchStateReview,
		EvidenceState: domain.EvidenceStateObserved,
		CapturedAt:    &time.Time{},
		Source:        &referenceSource,
	}}
	if err := service.CollectReferences(context.Background(), nil); !errors.Is(err, want) {
		t.Fatalf("reference store error = %v, want %v", err, want)
	}
}

var _ ports.CollectorPort = fakeCollector{}
var _ ports.ObservationStore = (*fakeObservationStore)(nil)
var _ ports.ReferenceStore = (*fakeReferenceStore)(nil)
