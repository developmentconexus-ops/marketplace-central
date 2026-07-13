package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type OrderLookup interface {
	FindExactOrder(ctx context.Context, scope domain.LinkageScope) (domain.MarketplaceOrder, bool, error)
}
