package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type OrderSource interface {
	ListOrders(ctx context.Context, installationID string, limit int) ([]domain.OrderIngestionSnapshot, error)
}
