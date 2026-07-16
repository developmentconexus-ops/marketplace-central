package ports

import (
	"context"
	"errors"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

var ErrOrderNotFound = errors.New("order_not_found")

type OrderNotFoundError struct {
	InstallationID  string
	ProviderOrderID string
}

func (e *OrderNotFoundError) Error() string { return ErrOrderNotFound.Error() }
func (e *OrderNotFoundError) Unwrap() error { return ErrOrderNotFound }

type OrderPage struct {
	Items      []domain.OrderReadModel
	NextCursor *OrderCursor
}

type OrderReadStore interface {
	ListOrders(context.Context, OrderListQuery) (OrderPage, error)
	GetOrder(context.Context, string, string) (domain.OrderReadModel, error)
}

type OrderStore interface {
	UpsertOrders(ctx context.Context, orders []domain.MarketplaceOrder) (importedCount int, skippedCount int, err error)
	ListOrders(ctx context.Context, installationID string, limit int) ([]domain.MarketplaceOrder, error)
}
