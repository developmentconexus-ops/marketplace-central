package domain

import (
	"testing"
	"time"
)

func productID(v string) *string { return &v }

func TestDeriveSignalStatusUnlinkedIsSemVinculoWithNilSignal(t *testing.T) {
	link := ListingLink{ProductID: nil}
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)

	got := DeriveSignalStatus(link, nil, now)

	if got != SignalStatusSemVinculo {
		t.Errorf("DeriveSignalStatus() = %q; want SEM_VINCULO", got)
	}
}

func TestDeriveSignalStatusLinkedNoEvidenceIsNoPriceEvidence(t *testing.T) {
	link := ListingLink{ProductID: productID("42664")}
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)

	got := DeriveSignalStatus(link, nil, now)

	if got != SignalStatusNoPriceEvidence {
		t.Errorf("DeriveSignalStatus() = %q; want NO_PRICE_EVIDENCE", got)
	}
}

func TestDeriveSignalStatusOldSignalIsStaleButValueRetained(t *testing.T) {
	link := ListingLink{ProductID: productID("42664")}
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	signal := &MarketSignal{
		Status:      SignalStatusOK,
		MatchStatus: "ACCEPT",
		Evidence: SignalEvidence{
			Source:    "ml_price_to_win",
			FetchedAt: now.Add(-2 * time.Hour),
			Freshness: "stale",
		},
	}

	got := DeriveSignalStatus(link, signal, now)

	if got != SignalStatusStale {
		t.Errorf("DeriveSignalStatus() = %q; want STALE", got)
	}
	if signal.MatchStatus != "ACCEPT" {
		t.Fatal("STALE derivation must retain the signal value, not discard it")
	}
}

func TestDeriveSignalStatusFreshSignalIsOK(t *testing.T) {
	link := ListingLink{ProductID: productID("42664")}
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	signal := &MarketSignal{
		Status:      SignalStatusOK,
		MatchStatus: "ACCEPT",
		Evidence: SignalEvidence{
			Source:    "ml_price_to_win",
			FetchedAt: now.Add(-10 * time.Minute),
			Freshness: "fresh",
		},
	}

	got := DeriveSignalStatus(link, signal, now)

	if got != SignalStatusOK {
		t.Errorf("DeriveSignalStatus() = %q; want OK", got)
	}
}

func TestBelowCostTargetUnderKnownCostIsTrue(t *testing.T) {
	cost := &Money{Amount: "100.00", Currency: "BRL"}
	target := &Money{Amount: "90.00", Currency: "BRL"}

	if !belowCost(cost, target) {
		t.Error("belowCost() = false; want true when target < cost")
	}
}

func TestBelowCostTargetAtOrAboveKnownCostIsFalse(t *testing.T) {
	cost := &Money{Amount: "100.00", Currency: "BRL"}
	target := &Money{Amount: "100.00", Currency: "BRL"}

	if belowCost(cost, target) {
		t.Error("belowCost() = true; want false when target == cost")
	}
}

func TestBelowCostUnknownCostIsFalseNeverZeroDefault(t *testing.T) {
	target := &Money{Amount: "10.00", Currency: "BRL"}

	if belowCost(nil, target) {
		t.Error("belowCost() = true with nil cost; ADR-17 requires custo desconhecido excluded, never counted below cost")
	}
}
