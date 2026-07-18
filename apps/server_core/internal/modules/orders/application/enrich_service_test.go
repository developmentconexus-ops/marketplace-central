package application

import (
	"context"
	"errors"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

func strPtr(s string) *string        { return &s }
func timePtr(t time.Time) *time.Time { return &t }

// fakeCostReader is a table-driven test double for ports.CostReader. It never
// touches internal_read; it lets callers script an Amount or an error per
// productID and records how many times it was invoked.
type fakeCostReader struct {
	calls   int
	amounts map[int]*float64
	errs    map[int]error
}

func newFakeCostReader() *fakeCostReader {
	return &fakeCostReader{amounts: map[int]*float64{}, errs: map[int]error{}}
}

func (f *fakeCostReader) GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error) {
	f.calls++
	if err := f.errs[productID]; err != nil {
		return internalreaddomain.CostAsOf{}, err
	}
	return internalreaddomain.CostAsOf{
		ProductID:   productID,
		Amount:      f.amounts[productID],
		AmountScope: internalreaddomain.CostAmountScopePerUnit,
	}, nil
}

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

	svc := NewEnrichService(newFakeCostReader())
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

func TestEnrichServiceEnrich_Cost(t *testing.T) {
	closedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	t.Run("known amount populates per-unit cost and excludes item", func(t *testing.T) {
		reader := newFakeCostReader()
		reader.amounts[10] = floatPtr(42.5)
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-cost-1",
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-1", InternalProductID: intPtr(10)},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if reader.calls != 1 {
			t.Fatalf("calls = %d, want 1", reader.calls)
		}
		if len(result.ItemCosts) != 1 || result.ItemCosts[0].UnitCost == nil || *result.ItemCosts[0].UnitCost != 42.5 {
			t.Fatalf("ItemCosts = %+v, want unit cost 42.5", result.ItemCosts)
		}
		if len(result.ComponentesDesconhecidos) != 0 {
			t.Fatalf("ComponentesDesconhecidos = %v, want empty", result.ComponentesDesconhecidos)
		}
	})

	t.Run("nil amount marks item unknown", func(t *testing.T) {
		reader := newFakeCostReader()
		reader.amounts[11] = nil
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-cost-2",
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-2", InternalProductID: intPtr(11)},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if reader.calls != 1 {
			t.Fatalf("calls = %d, want 1", reader.calls)
		}
		if result.ItemCosts[0].UnitCost != nil {
			t.Fatalf("UnitCost = %v, want nil", *result.ItemCosts[0].UnitCost)
		}
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "sku-2" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [sku-2]", result.ComponentesDesconhecidos)
		}
	})

	t.Run("reader error marks item unknown and order still returns", func(t *testing.T) {
		reader := newFakeCostReader()
		reader.errs[12] = errors.New("boom")
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-cost-3",
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-3", InternalProductID: intPtr(12)},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1 (order must not fail on cost miss)", len(got))
		}
		result := got[0]
		if result.ItemCosts[0].UnitCost != nil {
			t.Fatalf("UnitCost = %v, want nil", *result.ItemCosts[0].UnitCost)
		}
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "sku-3" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [sku-3]", result.ComponentesDesconhecidos)
		}
	})

	t.Run("nil InternalProductID skips the reader", func(t *testing.T) {
		reader := newFakeCostReader()
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-cost-4",
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-4"},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if reader.calls != 0 {
			t.Fatalf("calls = %d, want 0", reader.calls)
		}
		result := got[0]
		if result.ItemCosts[0].UnitCost != nil {
			t.Fatalf("UnitCost = %v, want nil", *result.ItemCosts[0].UnitCost)
		}
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "sku-4" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [sku-4]", result.ComponentesDesconhecidos)
		}
	})

	t.Run("all-nil order dates skip the reader", func(t *testing.T) {
		reader := newFakeCostReader()
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-cost-5",
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-5", InternalProductID: intPtr(13)},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if reader.calls != 0 {
			t.Fatalf("calls = %d, want 0", reader.calls)
		}
		result := got[0]
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "sku-5" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [sku-5]", result.ComponentesDesconhecidos)
		}
	})

	t.Run("buyer mask and vinculo status stay correct alongside cost", func(t *testing.T) {
		reader := newFakeCostReader()
		reader.amounts[14] = floatPtr(9.99)
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-cost-6",
			BuyerNickname:    strPtr("Ana Lima"),
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-6", InternalProductID: intPtr(14), LinkQuality: domain.LinkQualityResolved},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if result.Buyer.Display != "Ana L." {
			t.Fatalf("Buyer.Display = %q, want %q", result.Buyer.Display, "Ana L.")
		}
		if result.VinculoStatus != domain.VinculoStatusOK {
			t.Fatalf("VinculoStatus = %q, want OK", result.VinculoStatus)
		}
		if result.ItemCosts[0].UnitCost == nil || *result.ItemCosts[0].UnitCost != 9.99 {
			t.Fatalf("ItemCosts[0] = %+v, want unit cost 9.99", result.ItemCosts[0])
		}
	})

	t.Run("fallback to provider item id when seller sku is empty", func(t *testing.T) {
		reader := newFakeCostReader()
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-cost-7",
			Items: []domain.MarketplaceOrderItem{
				{ProviderItemID: "prov-7"},
			},
		}
		svc := NewEnrichService(reader)
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "prov-7" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [prov-7]", result.ComponentesDesconhecidos)
		}
	})
}
