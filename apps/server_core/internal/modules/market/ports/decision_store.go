package ports

import (
	"context"

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
}
