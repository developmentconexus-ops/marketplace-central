package application

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
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

// fakeShipmentReader is a table-driven test double for ports.ShipmentReader.
// It never touches connectors; it lets callers script a ShipmentInfo or an
// error per shipmentID and records how many times it was invoked.
type fakeShipmentReader struct {
	calls int
	infos map[string]connectorsdomain.ShipmentInfo
	errs  map[string]error
}

func newFakeShipmentReader() *fakeShipmentReader {
	return &fakeShipmentReader{infos: map[string]connectorsdomain.ShipmentInfo{}, errs: map[string]error{}}
}

func (f *fakeShipmentReader) GetShipment(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentInfo, error) {
	f.calls++
	if err := f.errs[shipmentID]; err != nil {
		return connectorsdomain.ShipmentInfo{}, err
	}
	return f.infos[shipmentID], nil
}

func testLogger() *slog.Logger { return slog.Default() }

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

	svc := NewEnrichService(newFakeCostReader(), newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
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
		svc := NewEnrichService(reader, newFakeShipmentReader(), testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if len(result.ComponentesDesconhecidos) != 1 || result.ComponentesDesconhecidos[0] != "prov-7" {
			t.Fatalf("ComponentesDesconhecidos = %v, want [prov-7]", result.ComponentesDesconhecidos)
		}
	})
}

func TestEnrichServiceEnrich_NilReaders(t *testing.T) {
	closedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	t.Run("nil cost reader marks every item unknown without panicking", func(t *testing.T) {
		shipmentReader := newFakeShipmentReader()
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-nilcost-1",
			ProviderClosedAt: timePtr(closedAt),
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-1", InternalProductID: intPtr(1)},
				{SellerSKU: "sku-2", InternalProductID: intPtr(2)},
			},
		}
		svc := NewEnrichService(nil, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		result := got[0]
		if len(result.ItemCosts) != 2 {
			t.Fatalf("ItemCosts = %+v, want 2 entries", result.ItemCosts)
		}
		for _, ic := range result.ItemCosts {
			if ic.UnitCost != nil {
				t.Fatalf("ItemCosts = %+v, want all UnitCost nil (cost reader is nil)", result.ItemCosts)
			}
		}
		if len(result.ComponentesDesconhecidos) != 2 {
			t.Fatalf("ComponentesDesconhecidos = %v, want [sku-1 sku-2]", result.ComponentesDesconhecidos)
		}
	})

	t.Run("nil shipment reader skips shipment lookup and leaves Buyer.UF nil", func(t *testing.T) {
		costReader := newFakeCostReader()
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-nilship-1",
			BuyerNickname:   strPtr("Rita Alves"),
			ShippingID:      "SHIP-9",
		}
		svc := NewEnrichService(costReader, nil, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		result := got[0]
		if result.Shipment != nil {
			t.Fatalf("Shipment = %+v, want nil (shipment reader is nil)", result.Shipment)
		}
		if result.Buyer.UF != nil {
			t.Fatalf("Buyer.UF = %v, want nil", *result.Buyer.UF)
		}
		if result.Buyer.Display != "Rita A." {
			t.Fatalf("Buyer.Display = %q, want %q", result.Buyer.Display, "Rita A.")
		}
	})

	t.Run("both readers nil produces a fully honest-null enrichment without panicking", func(t *testing.T) {
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-nilboth-1",
			ShippingID:      "SHIP-10",
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-3", InternalProductID: intPtr(3)},
			},
		}
		svc := NewEnrichService(nil, nil, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		result := got[0]
		if result.Shipment != nil {
			t.Fatalf("Shipment = %+v, want nil", result.Shipment)
		}
		if result.ItemCosts[0].UnitCost != nil {
			t.Fatalf("UnitCost = %v, want nil", *result.ItemCosts[0].UnitCost)
		}
	})
}

func TestEnrichServiceEnrich_Shipment(t *testing.T) {
	slaDue := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)

	t.Run("shipment present populates SLA/Delayed/DestinationUF/ShipmentID+Status and Buyer.UF", func(t *testing.T) {
		costReader := newFakeCostReader()
		shipmentReader := newFakeShipmentReader()
		delayed := true
		uf := "SP"
		shipmentReader.infos["SHIP-1"] = connectorsdomain.ShipmentInfo{
			ID:            "SHIP-1",
			Status:        "shipped",
			SLADue:        &slaDue,
			Delayed:       &delayed,
			DestinationUF: &uf,
		}
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-ship-1",
			BuyerNickname:   strPtr("Joao Silva"),
			ShippingID:      "SHIP-1",
		}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if shipmentReader.calls != 1 {
			t.Fatalf("shipment calls = %d, want 1", shipmentReader.calls)
		}
		if result.Shipment == nil {
			t.Fatalf("Shipment = nil, want populated")
		}
		if result.Shipment.ShipmentID != "SHIP-1" || result.Shipment.Status != "shipped" {
			t.Fatalf("Shipment = %+v, want ShipmentID=SHIP-1 Status=shipped", result.Shipment)
		}
		if result.Shipment.SLADue == nil || !result.Shipment.SLADue.Equal(slaDue) {
			t.Fatalf("Shipment.SLADue = %v, want %v", result.Shipment.SLADue, slaDue)
		}
		if result.Shipment.Delayed == nil || !*result.Shipment.Delayed {
			t.Fatalf("Shipment.Delayed = %v, want true", result.Shipment.Delayed)
		}
		if result.Shipment.DestinationUF == nil || *result.Shipment.DestinationUF != "SP" {
			t.Fatalf("Shipment.DestinationUF = %v, want SP", result.Shipment.DestinationUF)
		}
		if result.Buyer.UF == nil || *result.Buyer.UF != "SP" {
			t.Fatalf("Buyer.UF = %v, want SP", result.Buyer.UF)
		}
	})

	t.Run("reader error yields nil shipment, order still present, Buyer.UF stays nil", func(t *testing.T) {
		costReader := newFakeCostReader()
		shipmentReader := newFakeShipmentReader()
		shipmentReader.errs["SHIP-2"] = errors.New("boom")
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-ship-2",
			ShippingID:      "SHIP-2",
		}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1 (order must not fail on shipment miss)", len(got))
		}
		result := got[0]
		if shipmentReader.calls != 1 {
			t.Fatalf("shipment calls = %d, want 1 (reader was called)", shipmentReader.calls)
		}
		if result.Shipment != nil {
			t.Fatalf("Shipment = %+v, want nil", result.Shipment)
		}
		if result.Buyer.UF != nil {
			t.Fatalf("Buyer.UF = %v, want nil", *result.Buyer.UF)
		}
	})

	t.Run("empty shipping_id skips the reader and yields nil shipment", func(t *testing.T) {
		costReader := newFakeCostReader()
		shipmentReader := newFakeShipmentReader()
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-ship-3",
			ShippingID:      "",
		}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if shipmentReader.calls != 0 {
			t.Fatalf("shipment calls = %d, want 0 (reader must not be called)", shipmentReader.calls)
		}
		if result.Shipment != nil {
			t.Fatalf("Shipment = %+v, want nil", result.Shipment)
		}
	})

	t.Run("provider nil fields pass through as nil, never fabricated, Buyer.UF stays nil", func(t *testing.T) {
		costReader := newFakeCostReader()
		shipmentReader := newFakeShipmentReader()
		shipmentReader.infos["SHIP-4"] = connectorsdomain.ShipmentInfo{
			ID:     "SHIP-4",
			Status: "pending",
		}
		order := domain.OrderReadModel{
			ProviderOrderID: "ord-ship-4",
			ShippingID:      "SHIP-4",
		}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if result.Shipment == nil {
			t.Fatalf("Shipment = nil, want populated (with nil sub-fields)")
		}
		if result.Shipment.SLADue != nil {
			t.Fatalf("Shipment.SLADue = %v, want nil", result.Shipment.SLADue)
		}
		if result.Shipment.Delayed != nil {
			t.Fatalf("Shipment.Delayed = %v, want nil", result.Shipment.Delayed)
		}
		if result.Shipment.DestinationUF != nil {
			t.Fatalf("Shipment.DestinationUF = %v, want nil", result.Shipment.DestinationUF)
		}
		if result.Buyer.UF != nil {
			t.Fatalf("Buyer.UF = %v, want nil", *result.Buyer.UF)
		}
	})

	t.Run("S1 buyer-mask + vinculo + S2 cost stay correct alongside shipment", func(t *testing.T) {
		costReader := newFakeCostReader()
		costReader.amounts[20] = floatPtr(15.0)
		shipmentReader := newFakeShipmentReader()
		uf := "RJ"
		shipmentReader.infos["SHIP-5"] = connectorsdomain.ShipmentInfo{
			ID:            "SHIP-5",
			Status:        "delivered",
			DestinationUF: &uf,
		}
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-ship-5",
			BuyerNickname:    strPtr("Carlos Souza"),
			ProviderClosedAt: timePtr(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)),
			ShippingID:       "SHIP-5",
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-20", InternalProductID: intPtr(20), LinkQuality: domain.LinkQualityResolved},
			},
		}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
		result := got[0]
		if result.Buyer.Display != "Carlos S." {
			t.Fatalf("Buyer.Display = %q, want %q", result.Buyer.Display, "Carlos S.")
		}
		if result.Buyer.UF == nil || *result.Buyer.UF != "RJ" {
			t.Fatalf("Buyer.UF = %v, want RJ", result.Buyer.UF)
		}
		if result.VinculoStatus != domain.VinculoStatusOK {
			t.Fatalf("VinculoStatus = %q, want OK", result.VinculoStatus)
		}
		if result.ItemCosts[0].UnitCost == nil || *result.ItemCosts[0].UnitCost != 15.0 {
			t.Fatalf("ItemCosts[0] = %+v, want unit cost 15.0", result.ItemCosts[0])
		}
		if result.Shipment == nil || result.Shipment.ShipmentID != "SHIP-5" {
			t.Fatalf("Shipment = %+v, want ShipmentID=SHIP-5", result.Shipment)
		}
	})
}
