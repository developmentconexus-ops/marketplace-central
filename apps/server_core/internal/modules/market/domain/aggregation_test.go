package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestComputeMarketAggregate(t *testing.T) {
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	offer := func(seller, amount string) ValidatedOffer {
		return ValidatedOffer{SellerID: seller, Price: &Money{Amount: amount, Currency: "BRL"}, Condition: "new", MatchStatus: MatchStatusAccept, FetchedAt: now}
	}
	tests := []struct {
		name     string
		offers   []ValidatedOffer
		status   MarketAggregateStatus
		median   string
		min      string
		max      string
		nOffers  int
		nSellers int
	}{
		{name: "empty", status: MarketAggregateStatusNoPriceEvidence},
		{name: "one seller", offers: []ValidatedOffer{offer("s1", "12.30")}, status: MarketAggregateStatusInsufficientMarket, median: "12.30", min: "12.30", max: "12.30", nOffers: 1, nSellers: 1},
		{name: "four sellers", offers: []ValidatedOffer{offer("s1", "10"), offer("s2", "20"), offer("s3", "30"), offer("s4", "40")}, status: MarketAggregateStatusInsufficientMarket, median: "25.00", min: "10", max: "40", nOffers: 4, nSellers: 4},
		{name: "five sellers", offers: []ValidatedOffer{offer("s1", "5"), offer("s2", "4"), offer("s3", "3"), offer("s4", "2"), offer("s5", "1")}, status: MarketAggregateStatusOK, median: "3", min: "1", max: "5", nOffers: 5, nSellers: 5},
		{name: "even median", offers: []ValidatedOffer{offer("s1", "10.00"), offer("s2", "11.00")}, status: MarketAggregateStatusInsufficientMarket, median: "10.50", min: "10.00", max: "11.00", nOffers: 2, nSellers: 2},
		{name: "odd median exact", offers: []ValidatedOffer{offer("s1", "3.0"), offer("s2", "1"), offer("s3", "2.000")}, status: MarketAggregateStatusInsufficientMarket, median: "2.000", min: "1", max: "3.0", nOffers: 3, nSellers: 3},
		{name: "dedupe lowest", offers: []ValidatedOffer{offer("s1", "12"), offer("s1", "9"), offer("s2", "20")}, status: MarketAggregateStatusInsufficientMarket, median: "14.50", min: "9", max: "20", nOffers: 3, nSellers: 2},
		{name: "even median half-up tie", offers: []ValidatedOffer{offer("s1", "10.00"), offer("s2", "10.01")}, status: MarketAggregateStatusInsufficientMarket, median: "10.01", min: "10.00", max: "10.01", nOffers: 2, nSellers: 2},
		{name: "drops case variant condition", offers: []ValidatedOffer{offer("ok", "7"), {SellerID: "cap", Price: &Money{Amount: "5", Currency: "BRL"}, Condition: "New", MatchStatus: MatchStatusAccept}}, status: MarketAggregateStatusInsufficientMarket, median: "7", min: "7", max: "7", nOffers: 1, nSellers: 1},
		{name: "drops non qualifying", offers: []ValidatedOffer{offer("ok", "7"), {SellerID: "usd", Price: &Money{Amount: "1", Currency: "USD"}, Condition: "new", MatchStatus: MatchStatusAccept}, {SellerID: "used", Price: &Money{Amount: "2", Currency: "BRL"}, Condition: "used", MatchStatus: MatchStatusAccept}, {SellerID: "review", Price: &Money{Amount: "3", Currency: "BRL"}, Condition: "new", MatchStatus: MatchStatusReview}, {SellerID: "nil", Condition: "new", MatchStatus: MatchStatusAccept}}, status: MarketAggregateStatusInsufficientMarket, median: "7", min: "7", max: "7", nOffers: 1, nSellers: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeMarketAggregate("p1", tt.offers, MarketPriceSourceMLCatalogOffers, now, now.Add(time.Second))
			if err != nil {
				t.Fatal(err)
			}
			if got.Status != tt.status || got.NOffers != tt.nOffers || got.NSellers != tt.nSellers {
				t.Fatalf("got status/counts %#v", got)
			}
			if amount(got.Median) != tt.median || amount(got.MinValid) != tt.min || amount(got.MaxValid) != tt.max {
				t.Fatalf("got median=%q min=%q max=%q", amount(got.Median), amount(got.MinValid), amount(got.MaxValid))
			}
		})
	}
	got1, err := ComputeMarketAggregate("p1", tests[2].offers, MarketPriceSourceMLCatalogOffers, now, now)
	if err != nil {
		t.Fatal(err)
	}
	got2, err := ComputeMarketAggregate("p1", tests[2].offers, MarketPriceSourceMLCatalogOffers, now, now)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("non-deterministic: %#v != %#v", got1, got2)
	}
	malformed := []ValidatedOffer{{SellerID: "bad", Price: &Money{Amount: "abc", Currency: "BRL"}, Condition: "new", MatchStatus: MatchStatusAccept}}
	if _, err := ComputeMarketAggregate("p1", malformed, MarketPriceSourceMLCatalogOffers, now, now); err == nil {
		t.Fatal("malformed decimal amount must error, not panic or silently drop")
	}
}

func amount(m *Money) string {
	if m == nil {
		return ""
	}
	return m.Amount
}
