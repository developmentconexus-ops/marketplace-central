package domain_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

type fakeCollector struct{}

func (*fakeCollector) CollectObservations(context.Context, string, []string) ([]domain.MarketObservation, error) {
	return nil, nil
}

func (*fakeCollector) CollectReferences(context.Context, []string) ([]domain.MarketReference, error) {
	return nil, nil
}

var _ ports.CollectorPort = (*fakeCollector)(nil)

func TestObservationKeepsAllSignalsSeparate(t *testing.T) {
	stats := &domain.CatalogStats{
		Min:         "80.00",
		P25:         "90.00",
		Median:      "100.00",
		TrimmedMean: "101.00",
		P75:         "110.00",
		Max:         "120.00",
		NOffers:     10,
		NSellers:    5,
	}
	input := validObservationInput()
	input.OurSalePrice = &domain.Money{Amount: "100.00", Currency: "BRL"}
	input.WinnerPrice = &domain.Money{Amount: "95.00", Currency: "BRL"}
	input.CompetitiveTarget = &domain.Money{Amount: "94.00", Currency: "BRL"}
	input.CatalogOfferPrice = &domain.Money{Amount: "90.00", Currency: "BRL"}
	input.CatalogStats = stats

	observation, err := domain.NewStoredObservation(input)
	if err != nil {
		t.Fatalf("NewStoredObservation() error = %v", err)
	}

	if observation.OurSalePrice != input.OurSalePrice || observation.WinnerPrice != input.WinnerPrice ||
		observation.CompetitiveTarget != input.CompetitiveTarget || observation.CatalogOfferPrice != input.CatalogOfferPrice ||
		observation.CatalogStats != stats {
		t.Fatal("observation did not preserve each independent market signal")
	}
	for i := 0; i < reflect.TypeOf(observation).NumField(); i++ {
		name := reflect.TypeOf(observation).Field(i).Name
		if strings.Contains(strings.ToLower(name), "market_price") {
			t.Fatalf("observation contains forbidden aggregate field %q", name)
		}
	}
}

func TestCatalogStatsRequireAtLeastFiveSellers(t *testing.T) {
	input := validObservationInput()
	input.CatalogStats = &domain.CatalogStats{NSellers: 4}
	if _, err := domain.NewStoredObservation(input); err == nil {
		t.Fatal("expected catalog stats with four sellers to be rejected")
	}

	input.CatalogStats.NSellers = 5
	if _, err := domain.NewStoredObservation(input); err != nil {
		t.Fatalf("catalog stats with five sellers rejected: %v", err)
	}

	input.CatalogStats = nil
	if _, err := domain.NewStoredObservation(input); err != nil {
		t.Fatalf("nil catalog stats rejected: %v", err)
	}
}

func TestStoredObservationRequiresSourceAndCapturedAt(t *testing.T) {
	input := validObservationInput()
	input.Source = nil
	if _, err := domain.NewStoredObservation(input); err == nil || !strings.Contains(err.Error(), "source") {
		t.Fatalf("missing source error = %v, want source", err)
	}

	input = validObservationInput()
	input.CapturedAt = nil
	if _, err := domain.NewStoredObservation(input); err == nil || !strings.Contains(err.Error(), "captured_at") {
		t.Fatalf("missing captured_at error = %v, want captured_at", err)
	}
}

func TestMarketEnumValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*domain.MarketObservationInput)
		want   string
	}{
		{name: "evidence_state", mutate: func(input *domain.MarketObservationInput) { input.EvidenceState = "unknown" }, want: "evidence_state"},
		{name: "source", mutate: func(input *domain.MarketObservationInput) { bad := domain.Source("scraper"); input.Source = &bad }, want: "source"},
		{name: "competitive_status", mutate: func(input *domain.MarketObservationInput) {
			bad := domain.CompetitiveStatus("unknown")
			input.CompetitiveStatus = &bad
		}, want: "competitive_status"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validObservationInput()
			tt.mutate(&input)
			if _, err := domain.NewStoredObservation(input); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want field %q", err, tt.want)
			}
		})
	}

	reference := validReferenceInput()
	reference.MatchState = "unknown"
	if _, err := domain.NewStoredReference(reference); err == nil || !strings.Contains(err.Error(), "match_state") {
		t.Fatalf("reference error = %v, want match_state", err)
	}
}

func TestMoneyPairValidation(t *testing.T) {
	tests := []struct {
		name  string
		money *domain.Money
	}{
		{name: "currency", money: &domain.Money{Amount: "10.00", Currency: "USD"}},
		{name: "amount", money: &domain.Money{Amount: "", Currency: "BRL"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validObservationInput()
			input.OurSalePrice = tt.money
			if _, err := domain.NewStoredObservation(input); err == nil || !strings.Contains(err.Error(), "our_sale_price") {
				t.Fatalf("error = %v, want our_sale_price", err)
			}
		})
	}
}

func TestStoredReferenceRequiresProvenanceAndValidMatchState(t *testing.T) {
	input := validReferenceInput()
	input.Source = nil
	if _, err := domain.NewStoredReference(input); err == nil || !strings.Contains(err.Error(), "source") {
		t.Fatalf("missing source error = %v, want source", err)
	}

	input = validReferenceInput()
	input.CapturedAt = nil
	if _, err := domain.NewStoredReference(input); err == nil || !strings.Contains(err.Error(), "captured_at") {
		t.Fatalf("missing captured_at error = %v, want captured_at", err)
	}

	input = validReferenceInput()
	input.MatchState = "bad"
	if _, err := domain.NewStoredReference(input); err == nil || !strings.Contains(err.Error(), "match_state") {
		t.Fatalf("invalid match state error = %v, want match_state", err)
	}
}

func validObservationInput() domain.MarketObservationInput {
	source := domain.SourceOfficialAPI
	competitiveStatus := domain.CompetitiveStatusWinning
	capturedAt := time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)
	return domain.MarketObservationInput{
		ListingID:         "inst_1~MLB3456790~-",
		CompetitiveStatus: &competitiveStatus,
		EvidenceState:     domain.EvidenceStateObserved,
		CapturedAt:        &capturedAt,
		Source:            &source,
	}
}

func validReferenceInput() domain.MarketReferenceInput {
	source := domain.SourceVendor
	capturedAt := time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)
	return domain.MarketReferenceInput{
		ProductID:     "12345",
		MatchState:    domain.MatchStateReview,
		EvidenceState: domain.EvidenceStateInsufficientMarket,
		CapturedAt:    &capturedAt,
		Source:        &source,
	}
}
