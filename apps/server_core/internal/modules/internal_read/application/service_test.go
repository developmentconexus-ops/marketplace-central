package application

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestServiceDelegatesSellableStock(t *testing.T) {
	reader := fake.NewReader(fake.Fixtures{
		Stocks: map[int]domain.SellableStock{
			42664: {
				ProductID: 42664,
				Quantity:  float64ptr(3),
				Policy:    domain.DefaultSellableStockPolicy(),
				Source: domain.SourceMetadata{
					System:    "oracle",
					FetchedAt: time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	})
	svc := NewService(reader)

	got, err := svc.GetSellableStock(context.Background(), ports.SellableStockInput{
		ProductID: 42664,
		Policy:    domain.DefaultSellableStockPolicy(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity == nil || *got.Quantity != 3 {
		t.Fatalf("expected 3, got %v", got.Quantity)
	}
}

func float64ptr(v float64) *float64 {
	return &v
}
