package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

var (
	ErrSummaryInstallationRequired = errors.New("ORDERS_SUMMARY_INSTALLATION_REQUIRED")
	ErrSummaryStoreNotConfigured   = errors.New("ORDERS_SUMMARY_STORE_NOT_CONFIGURED")
)

type SummaryService struct {
	store ports.OrderSummaryStore
}

func NewSummaryService(store ports.OrderSummaryStore) SummaryService {
	return SummaryService{store: store}
}

func (s SummaryService) Summary(ctx context.Context, installationID string, referenceTime time.Time) (ports.OrderSummary, error) {
	if s.store == nil {
		return ports.OrderSummary{}, ErrSummaryStoreNotConfigured
	}
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return ports.OrderSummary{}, ErrSummaryInstallationRequired
	}
	return s.store.GetOrderSummary(ctx, installationID, referenceTime)
}
