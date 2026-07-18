package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

type SnapshotStore interface {
	AppendSnapshot(context.Context, domain.MarketPriceSnapshot) error
	LatestValidSnapshot(context.Context, domain.MarketPriceScope, string, domain.MarketPriceSource, time.Time) (*domain.MarketPriceSnapshot, error)
	SnapshotSeries(context.Context, domain.MarketPriceScope, string, domain.MarketPriceSource, time.Time) ([]domain.MarketPriceSnapshot, error)
}

type ValidatedOfferStore interface {
	AppendValidatedOffers(context.Context, []domain.ValidatedOffer) error
	ValidatedOffersByCatalogProduct(context.Context, string) ([]domain.ValidatedOffer, error)
}
