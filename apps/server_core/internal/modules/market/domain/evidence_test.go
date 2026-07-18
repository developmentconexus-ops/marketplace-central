package domain_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

func TestMarketPriceSnapshotRejectsInvalidMoneyAndFailedWithoutReason(t *testing.T) {
	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
	price := &domain.Money{Amount: "10.50", Currency: "BRL"}

	valid, err := domain.NewMarketPriceSnapshot(domain.MarketPriceSnapshotInput{
		Scope: domain.MarketPriceScopeListing, RefID: "MLB-1",
		Source: domain.MarketPriceSourceMLSalePrice, Price: price,
		ObservedAt: now, FetchedAt: now, ExpiresAt: now.Add(time.Hour),
		Status: domain.MarketPriceSnapshotStatusValid, RequestID: "req-1",
	})
	if err != nil || valid.Price != price {
		t.Fatalf("valid snapshot = %#v, error = %v", valid, err)
	}

	failed := domain.MarketPriceSnapshotInput{
		Scope: domain.MarketPriceScopeListing, RefID: "MLB-1",
		Source: domain.MarketPriceSourceMLSalePrice, ObservedAt: now,
		FetchedAt: now, ExpiresAt: now.Add(time.Hour),
		Status: domain.MarketPriceSnapshotStatusFailed, RequestID: "req-2",
	}
	if _, err := domain.NewMarketPriceSnapshot(failed); err == nil || !strings.Contains(err.Error(), "failure_reason") {
		t.Fatalf("failed snapshot error = %v, want failure_reason", err)
	}

	invalid := domain.MarketPriceSnapshotInput{
		Scope: domain.MarketPriceScopeListing, RefID: "MLB-1",
		Source:     domain.MarketPriceSourceMLSalePrice,
		Price:      &domain.Money{Amount: "", Currency: "BRL"},
		ObservedAt: now, FetchedAt: now, ExpiresAt: now.Add(time.Hour),
		Status: domain.MarketPriceSnapshotStatusValid, RequestID: "req-3",
	}
	if _, err := domain.NewMarketPriceSnapshot(invalid); err == nil || !strings.Contains(err.Error(), "price") {
		t.Fatalf("invalid money error = %v, want price", err)
	}
	invalid.Price = &domain.Money{Amount: "0", Currency: "BRL"}
	if _, err := domain.NewMarketPriceSnapshot(invalid); err == nil || !strings.Contains(err.Error(), "price") {
		t.Fatalf("zero money error = %v, want price", err)
	}

	unknown, err := domain.NewMarketPriceSnapshot(domain.MarketPriceSnapshotInput{
		Scope: domain.MarketPriceScopeListing, RefID: "MLB-1",
		Source: domain.MarketPriceSourceMLSalePrice, Price: nil,
		ObservedAt: now, FetchedAt: now, ExpiresAt: now.Add(time.Hour),
		Status: domain.MarketPriceSnapshotStatusFailed, FailureReason: stringPtr("provider unavailable"), RequestID: "req-4",
	})
	if err != nil || unknown.Price != nil {
		t.Fatalf("unknown price = %#v, error = %v; want nil price", unknown.Price, err)
	}
}

func TestMarketPriceSnapshotCrossFieldInvariants(t *testing.T) {
	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
	price := &domain.Money{Amount: "10.50", Currency: "BRL"}
	reason := stringPtr("provider unavailable")

	base := domain.MarketPriceSnapshotInput{
		Scope: domain.MarketPriceScopeListing, RefID: "MLB-1",
		Source:     domain.MarketPriceSourceMLSalePrice,
		ObservedAt: now, FetchedAt: now, ExpiresAt: now.Add(time.Hour),
		RequestID: "req-1",
	}

	validNilPrice := base
	validNilPrice.Status = domain.MarketPriceSnapshotStatusValid
	validNilPrice.Price = nil
	if _, err := domain.NewMarketPriceSnapshot(validNilPrice); err == nil || !strings.Contains(err.Error(), "price is required for VALID/EXPIRED snapshot") {
		t.Fatalf("VALID with nil price error = %v, want price required", err)
	}

	expiredNilPrice := base
	expiredNilPrice.Status = domain.MarketPriceSnapshotStatusExpired
	expiredNilPrice.Price = nil
	if _, err := domain.NewMarketPriceSnapshot(expiredNilPrice); err == nil || !strings.Contains(err.Error(), "price is required for VALID/EXPIRED snapshot") {
		t.Fatalf("EXPIRED with nil price error = %v, want price required", err)
	}

	failedWithPrice := base
	failedWithPrice.Status = domain.MarketPriceSnapshotStatusFailed
	failedWithPrice.Price = price
	failedWithPrice.FailureReason = reason
	if _, err := domain.NewMarketPriceSnapshot(failedWithPrice); err == nil || !strings.Contains(err.Error(), "price must be absent for FAILED snapshot") {
		t.Fatalf("FAILED with non-nil price error = %v, want price absent", err)
	}

	validWithReason := base
	validWithReason.Status = domain.MarketPriceSnapshotStatusValid
	validWithReason.Price = price
	validWithReason.FailureReason = reason
	if _, err := domain.NewMarketPriceSnapshot(validWithReason); err == nil || !strings.Contains(err.Error(), "failure_reason only valid for FAILED snapshot") {
		t.Fatalf("VALID with failure_reason error = %v, want failure_reason only valid for FAILED", err)
	}

	validOK := base
	validOK.Status = domain.MarketPriceSnapshotStatusValid
	validOK.Price = price
	got, err := domain.NewMarketPriceSnapshot(validOK)
	if err != nil || got.Price != price {
		t.Fatalf("VALID with valid price = %#v, error = %v; want success", got, err)
	}

	failedOK := base
	failedOK.Status = domain.MarketPriceSnapshotStatusFailed
	failedOK.Price = nil
	failedOK.FailureReason = reason
	gotFailed, err := domain.NewMarketPriceSnapshot(failedOK)
	if err != nil || gotFailed.Price != nil {
		t.Fatalf("FAILED with reason and nil price = %#v, error = %v; want success", gotFailed, err)
	}
}

func TestEvidenceEnumsAreExactAndVerdictStatesAreDistinct(t *testing.T) {
	if domain.MatchStatusAccept != "ACCEPT" || domain.PriceEvidenceStatusOK != "OK" || domain.BlockingStateNoCandidate != "NO_CANDIDATE" {
		t.Fatal("verdict enum values are not exact uppercase contract values")
	}
	if reflect.TypeOf(domain.MatchStatusAccept) == reflect.TypeOf(domain.PriceEvidenceStatusOK) ||
		reflect.TypeOf(domain.PriceEvidenceStatusOK) == reflect.TypeOf(domain.BlockingStateNoCandidate) ||
		reflect.TypeOf(domain.MatchStatusAccept) == reflect.TypeOf(domain.BlockingStateNoCandidate) {
		t.Fatal("verdict-state enums must be distinct Go types")
	}
	if domain.PriceEvidenceStatus("ok").IsValid() || domain.BlockingState("SEM_CUSTO ").IsValid() || domain.MatchStatus("ACCEPT ").IsValid() {
		t.Fatal("invalid enum strings were accepted")
	}

	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
	decision := domain.MatchDecisionInput{
		ProductID: "P-1", CatalogProductID: stringPtr("MLB-CAT-1"),
		Status: domain.MatchStatusAccept, Anchors: []string{"ean", "marca"},
		Contradictions: []string{}, ConfidenceBand: domain.ConfidenceBandAlta, DecidedAt: now,
	}
	if _, err := domain.NewMatchDecision(decision); err != nil {
		t.Fatalf("valid match decision rejected: %v", err)
	}
	decision.Status = domain.MatchStatus("accept")
	if _, err := domain.NewMatchDecision(decision); err == nil || !strings.Contains(err.Error(), "match_status") {
		t.Fatalf("invalid match status error = %v", err)
	}
}

func TestValidatedOfferPreservesUnknownPriceAndAcceptsAnyNonEmptyCondition(t *testing.T) {
	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
	offer, err := domain.NewValidatedOffer(domain.ValidatedOfferInput{
		CatalogProductID: "MLB-CAT-1", SellerID: "seller-1", Price: nil,
		Condition: "new", MatchStatus: domain.MatchStatusReview, FetchedAt: now,
	})
	if err != nil || offer.Price != nil {
		t.Fatalf("unknown offer price = %#v, error = %v; want nil price", offer.Price, err)
	}

	used, err := domain.NewValidatedOffer(domain.ValidatedOfferInput{
		CatalogProductID: "MLB-CAT-1", SellerID: "seller-1",
		Price:     &domain.Money{Amount: "10.00", Currency: "BRL"},
		Condition: "used", MatchStatus: domain.MatchStatusReview, FetchedAt: now,
	})
	if err != nil || used.Condition != "used" {
		t.Fatalf("used-condition offer = %#v, error = %v; want condition preserved, no error", used, err)
	}

	if _, err := domain.NewValidatedOffer(domain.ValidatedOfferInput{
		CatalogProductID: "MLB-CAT-1", SellerID: "seller-1",
		Condition: "", MatchStatus: domain.MatchStatusReview, FetchedAt: now,
	}); err == nil || !strings.Contains(err.Error(), "condition") {
		t.Fatalf("empty condition error = %v, want condition required", err)
	}
}

func TestCompetitiveSignalRequiresOurPrice(t *testing.T) {
	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)

	if _, err := domain.NewCompetitiveSignal(domain.CompetitiveSignalInput{
		ListingID: "MLB-1", OurPrice: nil, WinnerPrice: nil, TargetPrice: nil,
		FetchedAt: now,
	}); err == nil || !strings.Contains(err.Error(), "our_price is required") {
		t.Fatalf("nil our_price error = %v, want our_price is required", err)
	}

	ourPrice := &domain.Money{Amount: "99.90", Currency: "BRL"}
	signal, err := domain.NewCompetitiveSignal(domain.CompetitiveSignalInput{
		ListingID: "MLB-1", OurPrice: ourPrice, WinnerPrice: nil, TargetPrice: nil,
		FetchedAt: now,
	})
	if err != nil || signal.OurPrice != ourPrice || signal.WinnerPrice != nil || signal.TargetPrice != nil {
		t.Fatalf("valid signal = %#v, error = %v; want our_price set, others nil", signal, err)
	}
}

func TestMarketAggregateRequiresBRLAndExplicitCounts(t *testing.T) {
	now := time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)
	median := &domain.Money{Amount: "100.00", Currency: "BRL"}
	aggregate, err := domain.NewMarketAggregate(domain.MarketAggregateInput{
		ProductID: "P-1", Median: median, MinValid: median,
		NOffers: 7, NSellers: 5, ComputedAt: now, FetchedAt: now,
		Source: domain.MarketPriceSourceMLCatalogOffers, Status: domain.MarketAggregateStatusOK,
	})
	if err != nil || aggregate.NOffers != 7 || aggregate.NSellers != 5 {
		t.Fatalf("aggregate = %#v, error = %v; want explicit counts", aggregate, err)
	}

	aggregateInput := domain.MarketAggregateInput{
		ProductID: "P-1", Median: &domain.Money{Amount: "100.00", Currency: "USD"},
		NOffers: 7, NSellers: 5, ComputedAt: now, FetchedAt: now,
		Source: domain.MarketPriceSourceMLCatalogOffers, Status: domain.MarketAggregateStatusOK,
	}
	if _, err := domain.NewMarketAggregate(aggregateInput); err == nil || !strings.Contains(err.Error(), "BRL") {
		t.Fatalf("non-BRL aggregate error = %v, want BRL", err)
	}
}

func stringPtr(value string) *string { return &value }
