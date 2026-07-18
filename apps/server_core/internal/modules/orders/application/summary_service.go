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
	store       ports.OrderSummaryStore
	bucketStore ports.OrderBucketStore
}

func NewSummaryService(store ports.OrderSummaryStore) SummaryService {
	return SummaryService{store: store}
}

// NewSummaryServiceWithBuckets is additive to NewSummaryService: it wires an
// optional OrderBucketStore for BucketCounts (F01-A) without widening
// OrderSummaryStore or changing NewSummaryService's signature, so the sole
// production caller (composition/root.go) keeps compiling unchanged.
func NewSummaryServiceWithBuckets(store ports.OrderSummaryStore, bucketStore ports.OrderBucketStore) SummaryService {
	return SummaryService{store: store, bucketStore: bucketStore}
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

// BucketCounts mirrors Summary's honest-failure guards: a missing
// installationID or an unconfigured bucket store never yields fabricated
// zero counts (ADR-17), only an honest error.
func (s SummaryService) BucketCounts(ctx context.Context, installationID string) (ports.OrderBucketCounts, error) {
	if s.bucketStore == nil {
		return ports.OrderBucketCounts{}, ErrSummaryStoreNotConfigured
	}
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return ports.OrderBucketCounts{}, ErrSummaryInstallationRequired
	}
	return s.bucketStore.GetOrderBucketCounts(ctx, installationID)
}
