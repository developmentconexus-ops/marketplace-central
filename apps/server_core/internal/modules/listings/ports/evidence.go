package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

// EvidenceSignal is listings' own mirror of the market module's
// CompetitiveSignal (IC-03, per anúncio próprio). Fields are the subset
// ReadService.enrichSignals consumes to build a domain.MarketSignal.
type EvidenceSignal struct {
	ListingID   string
	OurPrice    *domain.Money
	TargetPrice *domain.Money
	Position    *domain.SignalPosition
	FetchedAt   time.Time
}

// EvidenceAggregate is listings' own mirror of the market module's
// MarketAggregate (IC-03, por codprod) — the counters plus the competitor
// price range (min — median — max) ReadService surfaces as the "faixa de
// mercado". The range Money pointers are nil when the market has no valid
// competitor stats (ADR-17: absent renders "—", never a fabricated zero).
type EvidenceAggregate struct {
	ProductID string
	NOffers   int
	NSellers  int
	Median    *domain.Money
	MinValid  *domain.Money
	MaxValid  *domain.Money
}

// EvidenceVerdict is listings' own mirror of the market module's Verdict
// (IC-03, por codprod) — only match_status is consumed here.
type EvidenceVerdict struct {
	ProductID   string
	MatchStatus string
}

// EvidenceReader is listings' consumer-owned port mirroring the three
// methods of market.EvidenceReader
// (internal/modules/market/ports/evidence_reader.go:12) that ReadService
// consumes. listings depends on its OWN port (hexagonal boundary) so this
// module never imports modules/market; the D-21 composition adapter
// (F01-S5) implements this interface over the real market reader at the
// composition root. Named second consumer: ReadService.enrichSignals below.
type EvidenceReader interface {
	Signals(ctx context.Context, listingIDs []string) ([]EvidenceSignal, error)
	Aggregates(ctx context.Context, codprods []string) ([]EvidenceAggregate, error)
	Verdicts(ctx context.Context, codprods []string) ([]EvidenceVerdict, error)
}
