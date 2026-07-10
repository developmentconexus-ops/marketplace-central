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
