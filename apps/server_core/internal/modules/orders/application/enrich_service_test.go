package application

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
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
	calls atomic.Int64
	infos map[string]connectorsdomain.ShipmentInfo
	errs  map[string]error
}

func newFakeShipmentReader() *fakeShipmentReader {
	return &fakeShipmentReader{infos: map[string]connectorsdomain.ShipmentInfo{}, errs: map[string]error{}}
}

func (f *fakeShipmentReader) GetShipment(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentInfo, error) {
	f.calls.Add(1)
	if err := f.errs[shipmentID]; err != nil {
		return connectorsdomain.ShipmentInfo{}, err
	}
	return f.infos[shipmentID], nil
}

// fakeBuyerFiscalReader is a table-driven test double for ports.BuyerFiscalReader.
// It never touches connectors; it lets callers script a BuyerFiscalInfo or an
// error per providerOrderID and records how many times it was invoked.
type fakeBuyerFiscalReader struct {
	calls int
	infos map[string]connectorsdomain.BuyerFiscalInfo
	errs  map[string]error
}

func newFakeBuyerFiscalReader() *fakeBuyerFiscalReader {
	return &fakeBuyerFiscalReader{infos: map[string]connectorsdomain.BuyerFiscalInfo{}, errs: map[string]error{}}
}

func (f *fakeBuyerFiscalReader) GetBuyerFiscal(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.BuyerFiscalInfo, error) {
	f.calls++
	if err := f.errs[providerOrderID]; err != nil {
		return connectorsdomain.BuyerFiscalInfo{}, err
	}
	return f.infos[providerOrderID], nil
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

// assertHonestEmptyProfitability asserts got equals
// domain.UnknownOrderProfitability(): every pointer nil, and
// ComponentesDesconhecidos names exactly the 7 known-unknown components.
func assertHonestEmptyProfitability(t *testing.T, got domain.OrderProfitability) {
	t.Helper()
	want := domain.UnknownOrderProfitability()
	if got.RetornoLiquido != nil {
		t.Fatalf("RetornoLiquido = %v, want nil", *got.RetornoLiquido)
	}
	if got.MargemPct != nil {
		t.Fatalf("MargemPct = %v, want nil", *got.MargemPct)
	}
	if got.Difal.Amount != nil || got.Difal.UFRoute != nil || got.Difal.DueDate != nil || got.Difal.Paid != nil {
		t.Fatalf("Difal = %+v, want all nil", got.Difal)
	}
	d := got.Decomposition
	if d.Comissao != nil || d.TaxaFixa != nil || d.Frete != nil || d.Imposto != nil ||
		d.Difal != nil || d.TarifaFull != nil || d.Custo != nil || d.MargemValor != nil || d.MargemPct != nil {
		t.Fatalf("Decomposition = %+v, want all pointers nil", d)
	}
	if len(d.ComponentesDesconhecidos) != len(want.Decomposition.ComponentesDesconhecidos) {
		t.Fatalf("ComponentesDesconhecidos = %v, want %v", d.ComponentesDesconhecidos, want.Decomposition.ComponentesDesconhecidos)
	}
	for i, name := range want.Decomposition.ComponentesDesconhecidos {
		if d.ComponentesDesconhecidos[i] != name {
			t.Fatalf("ComponentesDesconhecidos[%d] = %q, want %q", i, d.ComponentesDesconhecidos[i], name)
		}
	}
}

func TestEnrichServiceEnrich_NilDecomposerYieldsHonestEmptyProfitability(t *testing.T) {
	orders := []domain.OrderReadModel{
		{ProviderOrderID: "ord-prof-1"},
		{
			ProviderOrderID: "ord-prof-2",
			Items:           []domain.MarketplaceOrderItem{{SellerSKU: "sku-1"}},
		},
	}
	svc := NewEnrichService(newFakeCostReader(), newFakeShipmentReader(), testLogger())
	got := svc.Enrich(context.Background(), "install-1", orders)
	if len(got) != len(orders) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(orders))
	}
	for _, result := range got {
		assertHonestEmptyProfitability(t, result.Profitability)
	}
}

// TestEnrichServiceEnrich_ListNeverResolvesBuyerFiscal locks the M-08 list
// performance contract (FINDING-M08-LIST-TIMEOUT): buyer fiscal is DRAWER-only
// data. The list path (Enrich over N orders) must NEVER invoke the two-step
// BuyerFiscalReader — even when the reader is wired and has data — because that
// is +2 sequential ML calls per order and blows the request deadline at scale.
// Every EnrichedOrder from Enrich carries a nil BuyerFiscal; the detail path
// (EnrichOne) is the only place fiscal is resolved.
func TestEnrichServiceEnrich_ListNeverResolvesBuyerFiscal(t *testing.T) {
	fetchedAt := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	reader := newFakeBuyerFiscalReader()
	reader.infos["ord-list-1"] = connectorsdomain.BuyerFiscalInfo{Name: strPtr("Maria Souza"), DocType: strPtr("CPF"), FetchedAt: fetchedAt}
	reader.infos["ord-list-2"] = connectorsdomain.BuyerFiscalInfo{Name: strPtr("Ana Lima"), DocType: strPtr("CNPJ"), FetchedAt: fetchedAt}
	orders := []domain.OrderReadModel{
		{ProviderOrderID: "ord-list-1"},
		{ProviderOrderID: "ord-list-2"},
	}
	svc := NewEnrichServiceWithReaders(newFakeCostReader(), newFakeShipmentReader(), nil, reader, testLogger())
	got := svc.Enrich(context.Background(), "install-1", orders)
	if reader.calls != 0 {
		t.Fatalf("buyer fiscal calls = %d, want 0 (list must never resolve fiscal)", reader.calls)
	}
	for i, e := range got {
		if e.BuyerFiscal != nil {
			t.Fatalf("got[%d].BuyerFiscal = %+v, want nil (list path never fills fiscal)", i, e.BuyerFiscal)
		}
	}
}

// TestEnrichServiceEnrichOne_BuyerFiscal covers the DETAIL path: EnrichOne is
// the only entry that resolves the buyer's fiscal identity (single order = at
// most one two-step lookup), with the honest-absence / degrade semantics
// unchanged from the reader contract.
func TestEnrichServiceEnrichOne_BuyerFiscal(t *testing.T) {
	fetchedAt := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)

	t.Run("present populates BuyerFiscal keyed by provider order id", func(t *testing.T) {
		reader := newFakeBuyerFiscalReader()
		reader.infos["ord-bf-1"] = connectorsdomain.BuyerFiscalInfo{
			Name:      strPtr("Maria Souza"),
			DocType:   strPtr("CPF"),
			DocNumber: strPtr("12345678909"),
			Address: &connectorsdomain.BuyerFiscalAddress{
				City:      strPtr("São Paulo"),
				StateCode: strPtr("SP"),
			},
			FetchedAt: fetchedAt,
		}
		order := domain.OrderReadModel{ProviderOrderID: "ord-bf-1"}
		svc := NewEnrichServiceWithReaders(newFakeCostReader(), newFakeShipmentReader(), nil, reader, testLogger())
		result := svc.EnrichOne(context.Background(), "install-1", order)
		if reader.calls != 1 {
			t.Fatalf("buyer fiscal calls = %d, want 1", reader.calls)
		}
		if result.BuyerFiscal == nil {
			t.Fatalf("BuyerFiscal = nil, want populated")
		}
		if result.BuyerFiscal.Name == nil || *result.BuyerFiscal.Name != "Maria Souza" {
			t.Fatalf("BuyerFiscal.Name = %v, want Maria Souza", result.BuyerFiscal.Name)
		}
		if result.BuyerFiscal.DocType == nil || *result.BuyerFiscal.DocType != "CPF" {
			t.Fatalf("BuyerFiscal.DocType = %v, want CPF", result.BuyerFiscal.DocType)
		}
		if result.BuyerFiscal.Address == nil || result.BuyerFiscal.Address.StateCode == nil || *result.BuyerFiscal.Address.StateCode != "SP" {
			t.Fatalf("BuyerFiscal.Address = %+v, want StateCode=SP", result.BuyerFiscal.Address)
		}
	})

	t.Run("no-data info yields nil BuyerFiscal (honest absence, order still returned)", func(t *testing.T) {
		reader := newFakeBuyerFiscalReader()
		reader.infos["ord-bf-2"] = connectorsdomain.BuyerFiscalInfo{FetchedAt: fetchedAt}
		order := domain.OrderReadModel{ProviderOrderID: "ord-bf-2"}
		svc := NewEnrichServiceWithReaders(newFakeCostReader(), newFakeShipmentReader(), nil, reader, testLogger())
		got := svc.EnrichOne(context.Background(), "install-1", order)
		if reader.calls != 1 {
			t.Fatalf("buyer fiscal calls = %d, want 1", reader.calls)
		}
		if got.BuyerFiscal != nil {
			t.Fatalf("BuyerFiscal = %+v, want nil (HasData()==false)", got.BuyerFiscal)
		}
	})

	t.Run("reader error yields nil BuyerFiscal, order still returned", func(t *testing.T) {
		reader := newFakeBuyerFiscalReader()
		reader.errs["ord-bf-3"] = errors.New("boom")
		order := domain.OrderReadModel{ProviderOrderID: "ord-bf-3"}
		svc := NewEnrichServiceWithReaders(newFakeCostReader(), newFakeShipmentReader(), nil, reader, testLogger())
		got := svc.EnrichOne(context.Background(), "install-1", order)
		if reader.calls != 1 {
			t.Fatalf("buyer fiscal calls = %d, want 1", reader.calls)
		}
		if got.BuyerFiscal != nil {
			t.Fatalf("BuyerFiscal = %+v, want nil", got.BuyerFiscal)
		}
	})

	t.Run("nil reader skips lookup and yields nil BuyerFiscal", func(t *testing.T) {
		order := domain.OrderReadModel{ProviderOrderID: "ord-bf-4"}
		// NewEnrichService (3-arg) wires a nil buyer fiscal reader.
		svc := NewEnrichService(newFakeCostReader(), newFakeShipmentReader(), testLogger())
		result := svc.EnrichOne(context.Background(), "install-1", order)
		if result.BuyerFiscal != nil {
			t.Fatalf("BuyerFiscal = %+v, want nil (reader is nil)", result.BuyerFiscal)
		}
	})

	t.Run("empty provider order id skips the reader", func(t *testing.T) {
		reader := newFakeBuyerFiscalReader()
		order := domain.OrderReadModel{ProviderOrderID: ""}
		svc := NewEnrichServiceWithReaders(newFakeCostReader(), newFakeShipmentReader(), nil, reader, testLogger())
		svc.EnrichOne(context.Background(), "install-1", order)
		if reader.calls != 0 {
			t.Fatalf("buyer fiscal calls = %d, want 0 (empty provider order id must skip)", reader.calls)
		}
	})

	t.Run("stays correct alongside buyer mask, cost and shipment", func(t *testing.T) {
		costReader := newFakeCostReader()
		costReader.amounts[30] = floatPtr(7.5)
		shipmentReader := newFakeShipmentReader()
		uf := "RJ"
		shipmentReader.infos["SHIP-30"] = connectorsdomain.ShipmentInfo{ID: "SHIP-30", Status: "shipped", DestinationUF: &uf}
		bfReader := newFakeBuyerFiscalReader()
		bfReader.infos["ord-bf-5"] = connectorsdomain.BuyerFiscalInfo{DocType: strPtr("CNPJ"), DocNumber: strPtr("11222333000181"), FetchedAt: fetchedAt}
		order := domain.OrderReadModel{
			ProviderOrderID:  "ord-bf-5",
			BuyerNickname:    strPtr("Ana Lima"),
			ProviderClosedAt: timePtr(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)),
			ShippingID:       "SHIP-30",
			Items: []domain.MarketplaceOrderItem{
				{SellerSKU: "sku-30", InternalProductID: intPtr(30), LinkQuality: domain.LinkQualityResolved},
			},
		}
		svc := NewEnrichServiceWithReaders(costReader, shipmentReader, nil, bfReader, testLogger())
		result := svc.EnrichOne(context.Background(), "install-1", order)
		if result.Buyer.Display != "Ana L." {
			t.Fatalf("Buyer.Display = %q, want Ana L.", result.Buyer.Display)
		}
		if result.Buyer.UF == nil || *result.Buyer.UF != "RJ" {
			t.Fatalf("Buyer.UF = %v, want RJ", result.Buyer.UF)
		}
		if result.VinculoStatus != domain.VinculoStatusOK {
			t.Fatalf("VinculoStatus = %q, want OK", result.VinculoStatus)
		}
		if result.ItemCosts[0].UnitCost == nil || *result.ItemCosts[0].UnitCost != 7.5 {
			t.Fatalf("ItemCosts[0] = %+v, want unit cost 7.5", result.ItemCosts[0])
		}
		if result.BuyerFiscal == nil || result.BuyerFiscal.DocType == nil || *result.BuyerFiscal.DocType != "CNPJ" {
			t.Fatalf("BuyerFiscal = %+v, want DocType=CNPJ", result.BuyerFiscal)
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
		if shipmentReader.calls.Load() != 1 {
			t.Fatalf("shipment calls = %d, want 1", shipmentReader.calls.Load())
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

	t.Run("destination city/zip/receiver, carrier and costs propagate to enrichment and Buyer.City", func(t *testing.T) {
		costReader := newFakeCostReader()
		shipmentReader := newFakeShipmentReader()
		uf := "RJ"
		city := "Rio de Janeiro"
		zip := "20040-002"
		receiver := "João Silva"
		carrier := "Total Express"
		tracking := "http://tracking.totalexpress.com.br/poupup_track.php?reid=3"
		shipmentReader.infos["SHIP-C"] = connectorsdomain.ShipmentInfo{
			ID:              "SHIP-C",
			Status:          "shipped",
			DestinationUF:   &uf,
			DestinationCity: &city,
			DestinationZip:  &zip,
			ReceiverName:    &receiver,
			CarrierName:     &carrier,
			TrackingURL:     &tracking,
			Costs: &connectorsdomain.ShipmentCosts{
				GrossAmount: &connectorsdomain.Money{Amount: "23.45", Currency: "BRL"},
			},
		}
		order := domain.OrderReadModel{ProviderOrderID: "ord-ship-c", ShippingID: "SHIP-C"}
		svc := NewEnrichService(costReader, shipmentReader, testLogger())
		result := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})[0]
		if result.Shipment == nil {
			t.Fatalf("Shipment = nil, want populated")
		}
		s := result.Shipment
		if s.DestinationCity == nil || *s.DestinationCity != city {
			t.Fatalf("DestinationCity = %v, want %q", s.DestinationCity, city)
		}
		if s.DestinationZip == nil || *s.DestinationZip != zip {
			t.Fatalf("DestinationZip = %v, want %q", s.DestinationZip, zip)
		}
		if s.ReceiverName == nil || *s.ReceiverName != receiver {
			t.Fatalf("ReceiverName = %v, want %q", s.ReceiverName, receiver)
		}
		if s.CarrierName == nil || *s.CarrierName != carrier {
			t.Fatalf("CarrierName = %v, want %q", s.CarrierName, carrier)
		}
		if s.TrackingURL == nil || *s.TrackingURL != tracking {
			t.Fatalf("TrackingURL = %v, want %q", s.TrackingURL, tracking)
		}
		if s.Costs == nil || s.Costs.GrossAmount == nil || s.Costs.GrossAmount.Amount != "23.45" {
			t.Fatalf("Costs.GrossAmount = %+v, want 23.45", s.Costs)
		}
		if result.Buyer.City == nil || *result.Buyer.City != city {
			t.Fatalf("Buyer.City = %v, want %q (blessed source = shipment)", result.Buyer.City, city)
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
		if shipmentReader.calls.Load() != 1 {
			t.Fatalf("shipment calls = %d, want 1 (reader was called)", shipmentReader.calls.Load())
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
		if shipmentReader.calls.Load() != 0 {
			t.Fatalf("shipment calls = %d, want 0 (reader must not be called)", shipmentReader.calls.Load())
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

// TestEnrichParallelShipmentsPreservesOrderAndDegrades locks the perf fix:
// Enrich resolves shipments concurrently (bounded by shipmentConcurrency) but
// the returned slice stays in input order, and a per-order shipment error
// degrades that one order to a nil Shipment without aborting the batch or
// disturbing any other order's result.
func TestEnrichParallelShipmentsPreservesOrderAndDegrades(t *testing.T) {
	costReader := newFakeCostReader()
	shipmentReader := newFakeShipmentReader()
	uf1, city1 := "SP", "Sao Paulo"
	uf2, city2 := "RJ", "Rio de Janeiro"
	shipmentReader.infos["SHIP-A"] = connectorsdomain.ShipmentInfo{ID: "SHIP-A", Status: "shipped", DestinationUF: &uf1, DestinationCity: &city1}
	shipmentReader.infos["SHIP-D"] = connectorsdomain.ShipmentInfo{ID: "SHIP-D", Status: "shipped", DestinationUF: &uf2, DestinationCity: &city2}
	shipmentReader.errs["SHIP-C"] = errors.New("boom")

	orders := []domain.OrderReadModel{
		{ProviderOrderID: "ord-a", ShippingID: "SHIP-A"},
		{ProviderOrderID: "ord-b", ShippingID: ""},
		{ProviderOrderID: "ord-c", ShippingID: "SHIP-C"},
		{ProviderOrderID: "ord-d", ShippingID: "SHIP-D"},
	}
	svc := NewEnrichService(costReader, shipmentReader, testLogger())
	got := svc.Enrich(context.Background(), "install-1", orders)

	if len(got) != len(orders) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(orders))
	}
	for i, order := range orders {
		if got[i].Order.ProviderOrderID != order.ProviderOrderID {
			t.Fatalf("got[%d].Order.ProviderOrderID = %q, want %q (order preservation broken)", i, got[i].Order.ProviderOrderID, order.ProviderOrderID)
		}
	}

	// ord-a: shipment present, buyer UF/City filled.
	if got[0].Shipment == nil || got[0].Shipment.ShipmentID != "SHIP-A" {
		t.Fatalf("got[0].Shipment = %+v, want ShipmentID=SHIP-A", got[0].Shipment)
	}
	if got[0].Buyer.UF == nil || *got[0].Buyer.UF != "SP" {
		t.Fatalf("got[0].Buyer.UF = %v, want SP", got[0].Buyer.UF)
	}
	if got[0].Buyer.City == nil || *got[0].Buyer.City != "Sao Paulo" {
		t.Fatalf("got[0].Buyer.City = %v, want Sao Paulo", got[0].Buyer.City)
	}

	// ord-b: empty shipping id, reader skipped.
	if got[1].Shipment != nil {
		t.Fatalf("got[1].Shipment = %+v, want nil (empty shipping id)", got[1].Shipment)
	}

	// ord-c: reader error, honest degrade, order still present.
	if got[2].Shipment != nil {
		t.Fatalf("got[2].Shipment = %+v, want nil (reader error degrades)", got[2].Shipment)
	}

	// ord-d: shipment present, buyer UF/City filled.
	if got[3].Shipment == nil || got[3].Shipment.ShipmentID != "SHIP-D" {
		t.Fatalf("got[3].Shipment = %+v, want ShipmentID=SHIP-D", got[3].Shipment)
	}
	if got[3].Buyer.UF == nil || *got[3].Buyer.UF != "RJ" {
		t.Fatalf("got[3].Buyer.UF = %v, want RJ", got[3].Buyer.UF)
	}
}

// TestEnrichDecompositionFromRealData locks Goal C: with a nil Decomposer,
// resolveProfitability surfaces the real, already-resolved comissão (sum of
// item SaleFeeAmount) and frete (parsed sender cost) instead of an always-—
// UnknownOrderProfitability.
func TestEnrichDecompositionFromRealData(t *testing.T) {
	costReader := newFakeCostReader()
	shipmentReader := newFakeShipmentReader()
	senderCost := &connectorsdomain.Money{Amount: "12.34", Currency: "BRL"}
	shipmentReader.infos["SHIP-DEC"] = connectorsdomain.ShipmentInfo{
		ID:     "SHIP-DEC",
		Status: "shipped",
		Costs:  &connectorsdomain.ShipmentCosts{SenderCost: senderCost},
	}
	saleFee := 22.95
	order := domain.OrderReadModel{
		ProviderOrderID: "ord-dec-1",
		Total:           floatPtr(229.5),
		ShippingID:      "SHIP-DEC",
		Items: []domain.MarketplaceOrderItem{
			{SellerSKU: "sku-dec-1", SaleFeeAmount: &saleFee, Quantity: 1},
		},
	}
	svc := NewEnrichService(costReader, shipmentReader, testLogger())
	got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})
	result := got[0]

	if result.Profitability.Decomposition.Comissao == nil || *result.Profitability.Decomposition.Comissao != 22.95 {
		t.Fatalf("Decomposition.Comissao = %v, want 22.95", result.Profitability.Decomposition.Comissao)
	}
	if result.Profitability.Decomposition.Frete == nil || *result.Profitability.Decomposition.Frete != 12.34 {
		t.Fatalf("Decomposition.Frete = %v, want 12.34", result.Profitability.Decomposition.Frete)
	}
}
