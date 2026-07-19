package domain

import (
	"math/big"
	"time"
)

// SignalStatus is a listings-owned COMPOSITE over IC-03 states (NOT a market
// enum): SEM_VINCULO (no codprod link) != NO_PRICE_EVIDENCE (linked, no market
// evidence) != STALE (evidence older than signalStaleTTL) — distinct honest
// states per ADR-17, never collapsed, never 0/null-silent.
type SignalStatus string

const (
	SignalStatusOK              SignalStatus = "OK"
	SignalStatusSemVinculo      SignalStatus = "SEM_VINCULO"
	SignalStatusNoPriceEvidence SignalStatus = "NO_PRICE_EVIDENCE"
	SignalStatusStale           SignalStatus = "STALE"
)

// signalStaleTTL mirrors the market module's snapshotTTL
// (internal/modules/market/application/collection_pipeline_service.go:19,
// ExpiresAt = fetchedAt.Add(snapshotTTL)) — the IC-03 EXPIRED horizon.
// listings owns the composite signal_status per IC-03 line 29; the reader
// port exposes no expires_at, so age is derived here from FetchedAt.
const signalStaleTTL = time.Hour

// SignalPosition mirrors IC-03 CompetitiveSignal.position {rank, total}.
type SignalPosition struct {
	Rank  int `json:"rank"`
	Total int `json:"total"`
}

// SignalEvidence carries the P1c evidence obligation (source, age, freshness)
// that must accompany every market-price display.
type SignalEvidence struct {
	Source    string    `json:"source"`
	FetchedAt time.Time `json:"fetched_at"`
	Freshness string    `json:"freshness"`
}

// MarketSignal is the listings-owned mirror of IC-03 CompetitiveSignal, enriched
// with the evidence obligation fields. It is nil whenever SignalStatus is
// SEM_VINCULO or NO_PRICE_EVIDENCE (no signal to show).
type MarketSignal struct {
	Status      SignalStatus    `json:"status"`
	Position    *SignalPosition `json:"position"`
	PriceToWin  *Money          `json:"price_to_win"`
	DeltaPct    *string         `json:"delta_pct"`
	MatchStatus string          `json:"match_status"`
	NOffers     int             `json:"n_offers"`
	NSellers    int             `json:"n_sellers"`
	Evidence    SignalEvidence  `json:"evidence"`
}

// DeriveSignalStatus is pure: no port dependency. signal is the listings-owned
// per-listing MarketSignal (nil when none was mapped). now is injected for
// deterministic staleness checks.
func DeriveSignalStatus(link ListingLink, signal *MarketSignal, now time.Time) SignalStatus {
	if link.ProductID == nil {
		return SignalStatusSemVinculo
	}
	if signal == nil {
		return SignalStatusNoPriceEvidence
	}
	if now.Sub(signal.Evidence.FetchedAt) > signalStaleTTL {
		return SignalStatusStale
	}
	return SignalStatusOK
}

// belowCost implements the abaixo_custo rule: cost known AND target < cost.
// Cost nil (custo desconhecido) is EXCLUDED per ADR-17 — unknown never counts
// as below cost. Reuses the big.Rat decimal-string comparison pattern from
// belowMargin (internal/modules/listings/application/read_service.go:532).
func belowCost(cost *Money, target *Money) bool {
	if cost == nil || target == nil {
		return false
	}
	costRat, ok := new(big.Rat).SetString(cost.Amount)
	if !ok {
		return false
	}
	targetRat, ok := new(big.Rat).SetString(target.Amount)
	if !ok {
		return false
	}
	return targetRat.Cmp(costRat) < 0
}

// BelowCost is the exported form of belowCost (above) for the F01-S4
// consumers outside this package (application.ReadService's summary tally
// and List/ByProduct exception-filter scan): same ADR-17 rule, same
// unknown-cost exclusion, not re-implemented at the call site.
func BelowCost(cost *Money, target *Money) bool {
	return belowCost(cost, target)
}
