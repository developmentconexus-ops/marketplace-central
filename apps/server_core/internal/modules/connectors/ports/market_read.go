package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type MarketReader interface {
	GetOwnItemPricing(ctx context.Context, accountRef domain.ProviderAccountRef, itemID string) (domain.OwnItemPricing, error)
	GetPriceToWin(ctx context.Context, accountRef domain.ProviderAccountRef, itemID string) (domain.PriceToWin, error)
}
