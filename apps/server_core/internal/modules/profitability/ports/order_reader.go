package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type OrderReader interface {
	ListOrders(ctx context.Context, installationID string, limit int) ([]domain.OrderFact, error)
}
