package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

type MatchDecisionStore interface {
	AppendMatchDecision(context.Context, domain.MatchDecision) error
	LatestMatchDecisions(context.Context, []string) ([]domain.MatchDecision, error)
}

type CompetitiveSignalStore interface {
	AppendCompetitiveSignals(context.Context, []domain.CompetitiveSignal) error
	LatestCompetitiveSignals(context.Context, []string) ([]domain.CompetitiveSignal, error)
}

type MarketAggregateStore interface {
	AppendMarketAggregates(context.Context, []domain.MarketAggregate) error
	LatestMarketAggregates(context.Context, []string) ([]domain.MarketAggregate, error)
	LatestMarketAggregatesBySource(context.Context, []string, domain.MarketPriceSource) ([]domain.MarketAggregate, error)

	// StaleAggregateProductIDs enumerates the product IDs whose latest
	// aggregate is older than olderThan, oldest first, capped at limit. The
	// other three read methods above all require the caller to already know
	// which product IDs to ask about; this is the only one that answers "who
	// went stale", which is the question a periodic collection job needs
	// answered before it can decide what to re-collect.
	StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error)
}
