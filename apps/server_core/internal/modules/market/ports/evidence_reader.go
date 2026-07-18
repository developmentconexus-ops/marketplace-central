package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// EvidenceReader is the market module's exported read-only facade for
// price-intel evidence (F-04). Its signature is frozen at this feature's
// close; M-05 F-01 consumes it read-only.
type EvidenceReader interface {
	Signals(ctx context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error)
	Aggregates(ctx context.Context, codprods []string) ([]domain.MarketAggregate, error)
	Verdicts(ctx context.Context, codprods []string) ([]domain.Verdict, error)
}
