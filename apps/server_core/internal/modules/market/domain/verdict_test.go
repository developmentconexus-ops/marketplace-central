package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestAssembleVerdictPrecedenceAndInputs(t *testing.T) {
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	money := &Money{Amount: "10", Currency: "BRL"}
	agg := func(status MarketAggregateStatus) *MarketAggregate {
		v, err := NewMarketAggregate(MarketAggregateInput{ProductID: "p1", Median: money, MinValid: nil, NOffers: 6, NSellers: 5, Source: MarketPriceSourceMLCatalogOffers, FetchedAt: now, ComputedAt: now, Status: status})
		if err != nil {
			t.Fatal(err)
		}
		return &v
	}
	ok, insufficient, empty := agg(MarketAggregateStatusOK), agg(MarketAggregateStatusInsufficientMarket), agg(MarketAggregateStatusNoPriceEvidence)
	tests := []struct {
		name         string
		status       MatchStatus
		aggregate    *MarketAggregate
		cost         *CostRef
		wantEvidence PriceEvidenceStatus
		wantBlock    *BlockingState
		wantRange    bool
	}{
		{"review", MatchStatusReview, ok, &CostRef{AsOf: now}, PriceEvidenceStatusNoPriceEvidence, block(BlockingStateNoCandidate), false},
		{"reject", MatchStatusReject, ok, &CostRef{AsOf: now}, PriceEvidenceStatusNoPriceEvidence, block(BlockingStateNoCandidate), false},
		{"no candidate", MatchStatusNoCandidate, nil, nil, PriceEvidenceStatusNoPriceEvidence, block(BlockingStateNoCandidate), false},
		{"no aggregate", MatchStatusAccept, nil, nil, PriceEvidenceStatusNoPriceEvidence, block(BlockingStateNoPriceEvidence), false},
		{"empty aggregate", MatchStatusAccept, empty, nil, PriceEvidenceStatusNoPriceEvidence, block(BlockingStateNoPriceEvidence), false},
		{"insufficient", MatchStatusAccept, insufficient, nil, PriceEvidenceStatusInsufficientMarket, block(BlockingStateInsufficientMarket), true},
		{"sem custo", MatchStatusAccept, ok, nil, PriceEvidenceStatusOK, block(BlockingStateSemCusto), true},
		{"complete", MatchStatusAccept, ok, &CostRef{AsOf: now}, PriceEvidenceStatusOK, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: tt.status}, Aggregate: tt.aggregate, Cost: tt.cost}
			got, err := AssembleVerdict(input)
			if err != nil {
				t.Fatal(err)
			}
			if got.PriceEvidenceStatus != tt.wantEvidence || !reflect.DeepEqual(got.BlockingState, tt.wantBlock) || (got.MarketRange != nil) != tt.wantRange {
				t.Fatalf("got %#v", got)
			}
			if got.VerdictLabel != nil {
				t.Fatal("verdict label must be nil")
			}
		})
	}

	input := VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: ok, OurPrice: &SnapshotRef{Source: MarketPriceSourceMLSalePrice, AsOf: now}, TargetPrice: &SnapshotRef{Source: MarketPriceSourceMLPriceToWin, AsOf: now.Add(time.Second)}, Cost: &CostRef{AsOf: now.Add(2 * time.Second)}}
	got, err := AssembleVerdict(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.InputsUsed) != 4 {
		t.Fatalf("inputs %#v", got.InputsUsed)
	}
	if got.MarketRange == nil || got.MarketRange.MinValid != nil || got.MarketRange.Median != money || got.MarketRange.Currency != "BRL" || got.MarketRange.NOffers != 6 || got.MarketRange.NSellers != 5 {
		t.Fatalf("range %#v", got.MarketRange)
	}
	got2, err := AssembleVerdict(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, got2) {
		t.Fatalf("non-deterministic: %#v != %#v", got, got2)
	}
	invalidAggregate := *ok
	invalidAggregate.Status = MarketAggregateStatus("INVALID")
	if _, err := AssembleVerdict(VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusReview}, Aggregate: &invalidAggregate}); err != nil {
		t.Fatalf("identity blocker consulted aggregate: %v", err)
	}

	badSource := *ok
	badSource.Source = MarketPriceSourceMLSalePrice
	errCases := []struct {
		name  string
		input VerdictInput
	}{
		{"empty product id", VerdictInput{ProductID: " ", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: ok}},
		{"invalid decision status", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatus("BOGUS")}, Aggregate: ok}},
		{"accept bad aggregate status", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: &invalidAggregate}},
		{"accept wrong aggregate source", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: &badSource}},
		{"our price wrong source", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: ok, OurPrice: &SnapshotRef{Source: MarketPriceSourceMLCatalogOffers, AsOf: now}}},
		{"target price zero asof", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: ok, TargetPrice: &SnapshotRef{Source: MarketPriceSourceMLPriceToWin}}},
		{"cost zero asof", VerdictInput{ProductID: "p1", Decision: MatchDecision{Status: MatchStatusAccept}, Aggregate: ok, Cost: &CostRef{}}},
	}
	for _, tc := range errCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := AssembleVerdict(tc.input); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func block(v BlockingState) *BlockingState { return &v }
