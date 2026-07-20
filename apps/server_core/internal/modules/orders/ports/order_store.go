package ports

import (
	"context"
	"errors"
	"time"

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

type OrderSummary struct {
	Today     int64
	SevenDays int64
}

type OrderSummaryStore interface {
	GetOrderSummary(ctx context.Context, installationID string, referenceTime time.Time) (OrderSummary, error)
}

// OrderBucketCounts is the by=status response shape for GET /orders/summary
// (F01-A). It carries the four KPI/Lista-tab/Kanban-column buckets;
// domain.BucketCancelado is intentionally excluded (see domain/order_bucket.go).
type OrderBucketCounts struct {
	Novo    int64
	Faturar int64
	Enviar  int64
	Enviado int64
}

// OrderBucketStore is additive to OrderSummaryStore so existing fakes/stores
// that only implement OrderSummaryStore keep compiling unchanged.
type OrderBucketStore interface {
	GetOrderBucketCounts(ctx context.Context, installationID string) (OrderBucketCounts, error)
}

// OrderFaturadoStore records the operator's "faturado" (invoiced) mark. Our-DB
// only; idempotent; tenant+installation scoped.
type OrderFaturadoStore interface {
	MarkOrderFaturado(ctx context.Context, installationID, providerOrderID string) error
}
