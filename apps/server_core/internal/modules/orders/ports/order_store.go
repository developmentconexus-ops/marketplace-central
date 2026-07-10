package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type OrderStore interface {
	UpsertOrders(ctx context.Context, orders []domain.MarketplaceOrder) (importedCount int, skippedCount int, err error)
	ListOrders(ctx context.Context, installationID string, limit int) ([]domain.MarketplaceOrder, error)
}
