package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type StockActionStore interface {
	FindByID(ctx context.Context, actionID string) (domain.StockAction, bool, error)
	Save(ctx context.Context, action domain.StockAction) error
}

type StockWriter interface {
	UpdateAvailableQuantity(ctx context.Context, request domain.StockWriteRequest) (domain.StockWriteResult, error)
}
