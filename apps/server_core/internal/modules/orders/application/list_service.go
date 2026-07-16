package application

import (
	"context"
	"errors"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type ListOrdersInput struct {
	InstallationID string
	Limit          int
}

type ListService struct {
	store ports.OrderStore
}

func NewListService(store ports.OrderStore) ListService {
	return ListService{store: store}
}

func (s ListService) List(ctx context.Context, input ListOrdersInput) ([]domain.MarketplaceOrder, error) {
	if s.store == nil {
		return nil, errors.New("ORDERS_LIST_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return nil, errors.New("ORDERS_INSTALLATION_REQUIRED")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	return s.store.ListOrders(ctx, installationID, limit)
}

// ReadService delegates canonical order reads to the tenant-scoped read port.
// The legacy ListService above remains available to existing consumers.
type ReadService struct {
	store ports.OrderReadStore
}

var ErrOrderReadStoreNotConfigured = errors.New("ORDERS_READ_STORE_NOT_CONFIGURED")

func NewReadService(store ports.OrderReadStore) ReadService {
	return ReadService{store: store}
}

func (s ReadService) List(ctx context.Context, query ports.OrderListQuery) (ports.OrderPage, error) {
	if s.store == nil {
		return ports.OrderPage{}, ErrOrderReadStoreNotConfigured
	}
	return s.store.ListOrders(ctx, query)
}

func (s ReadService) Get(ctx context.Context, installationID, providerOrderID string) (domain.OrderReadModel, error) {
	if s.store == nil {
		return domain.OrderReadModel{}, ErrOrderReadStoreNotConfigured
	}
	return s.store.GetOrder(ctx, installationID, providerOrderID)
}
