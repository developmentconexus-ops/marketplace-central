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
				Codprod:         42664,
				Quantity:        3,
				Scope:           domain.DefaultSellableStockScope(),
				SourceFetchedAt: time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
			},
		},
	})
	svc := NewService(reader)

	got, err := svc.GetSellableStock(context.Background(), ports.SellableStockInput{Codprod: 42664})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity != 3 {
		t.Fatalf("expected 3, got %v", got.Quantity)
	}
}
