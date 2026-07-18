package application

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

func strPtr(s string) *string { return &s }

func TestEnrichServiceEnrich(t *testing.T) {
	cases := []struct {
		name        string
		order       domain.OrderReadModel
		wantVinculo domain.VinculoStatus
		wantDisplay string
	}{
		{
			name: "resolved item yields OK",
			order: domain.OrderReadModel{
				ProviderOrderID: "ord-1",
				BuyerNickname:   strPtr("Maria Souza"),
				Items: []domain.MarketplaceOrderItem{
					{LinkQuality: domain.LinkQualityUnresolved},
					{LinkQuality: domain.LinkQualityResolved},
				},
			},
			wantVinculo: domain.VinculoStatusOK,
			wantDisplay: "Maria S.",
		},
		{
			name: "no resolved item yields SEM_VINCULO",
			order: domain.OrderReadModel{
				ProviderOrderID: "ord-2",
				BuyerNickname:   strPtr("PEDRO"),
				Items: []domain.MarketplaceOrderItem{
					{LinkQuality: domain.LinkQualityMissing},
				},
			},
			wantVinculo: domain.VinculoStatusSemVinculo,
			wantDisplay: "PEDRO",
		},
		{
			name: "nil buyer nickname yields empty masked buyer",
			order: domain.OrderReadModel{
				ProviderOrderID: "ord-3",
				BuyerNickname:   nil,
				Items:           nil,
			},
			wantVinculo: domain.VinculoStatusSemVinculo,
			wantDisplay: "",
		},
	}

	svc := NewEnrichService()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{tc.order})
			if len(got) != 1 {
				t.Fatalf("len(got) = %d, want 1", len(got))
			}
			result := got[0]
			if result.VinculoStatus != tc.wantVinculo {
				t.Fatalf("VinculoStatus = %q, want %q", result.VinculoStatus, tc.wantVinculo)
			}
			if result.Buyer.Display != tc.wantDisplay {
				t.Fatalf("Buyer.Display = %q, want %q", result.Buyer.Display, tc.wantDisplay)
			}
			if result.Buyer.City != nil || result.Buyer.UF != nil {
				t.Fatalf("Buyer carries unexpected city/uf: %+v", result.Buyer)
			}
			if result.Order.ProviderOrderID != tc.order.ProviderOrderID {
				t.Fatalf("Order.ProviderOrderID = %q, want %q", result.Order.ProviderOrderID, tc.order.ProviderOrderID)
			}
		})
	}
}
