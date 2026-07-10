package orders

import (
	"context"
	"reflect"
	"testing"
	"time"

	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type stubOrderLister struct {
	orders []ordersdomain.MarketplaceOrder
}

func (s stubOrderLister) List(context.Context, ordersapp.ListOrdersInput) ([]ordersdomain.MarketplaceOrder, error) {
	return s.orders, nil
}

func TestOrderReaderMapsOrderFacts(t *testing.T) {
	createdAt := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)
	closedAt := createdAt.Add(2 * time.Hour)
	updatedAt := closedAt.Add(5 * time.Minute)
	fetchedAt := updatedAt.Add(time.Minute)
	price := 125.50
	fee := 18.75
	productID := 42

	reader := NewOrderReader(stubOrderLister{orders: []ordersdomain.MarketplaceOrder{{
		InstallationID:    "installation-1",
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   "order-1",
		ProviderStatus:    "paid",
		ProviderCreatedAt: &createdAt,
		ProviderClosedAt:  &closedAt,
		ProviderUpdatedAt: &updatedAt,
		FetchedAt:         fetchedAt,
		Items: []ordersdomain.MarketplaceOrderItem{{
			ProviderItemID:      "MLB123",
			ProviderVariationID: "V-1",
			Quantity:            2,
			UnitPrice:           &price,
			SaleFeeAmount:       &fee,
			LinkQuality:         ordersdomain.LinkQualityResolved,
			InternalProductID:   &productID,
		}},
	}}})

	facts, err := reader.ListOrders(context.Background(), "installation-1", 10)
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	want := []profitabilitydomain.OrderFact{{
		InstallationID:    "installation-1",
		ProviderOrderID:   "order-1",
		RealizationState:  profitabilitydomain.OrderRealizationRealized,
		ProviderCreatedAt: &createdAt,
		ProviderClosedAt:  &closedAt,
		ProviderUpdatedAt: &updatedAt,
		FetchedAt:         fetchedAt,
		Items: []profitabilitydomain.OrderItemFact{{
			ProviderItemID:      "MLB123",
			ProviderVariationID: "V-1",
			Quantity:            2,
			UnitPrice:           &price,
			SaleFeeAmount:       &fee,
			LinkQuality:         profitabilitydomain.OrderLinkResolved,
			InternalProductID:   &productID,
		}},
	}}
	if !reflect.DeepEqual(facts, want) {
		t.Fatalf("ListOrders() = %#v, want %#v", facts, want)
	}
}

func TestOrderReaderMapsRealizationState(t *testing.T) {
	tests := []struct {
		name         string
		providerCode string
		status       string
		want         profitabilitydomain.OrderRealizationState
	}{
		{name: "mercado livre paid", providerCode: "mercado_livre", status: "paid", want: profitabilitydomain.OrderRealizationRealized},
		{name: "mercado livre cancelled", providerCode: "mercado_livre", status: "cancelled", want: profitabilitydomain.OrderRealizationNotRealized},
		{name: "mercado livre canceled", providerCode: "mercado_livre", status: "canceled", want: profitabilitydomain.OrderRealizationNotRealized},
		{name: "mercado livre blank", providerCode: "mercado_livre", status: "", want: profitabilitydomain.OrderRealizationUnknown},
		{name: "mercado livre unsupported", providerCode: "mercado_livre", status: "handling", want: profitabilitydomain.OrderRealizationUnknown},
		{name: "other provider paid", providerCode: "other", status: "paid", want: profitabilitydomain.OrderRealizationUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewOrderReader(stubOrderLister{orders: []ordersdomain.MarketplaceOrder{{
				ProviderCode:   tt.providerCode,
				ProviderStatus: tt.status,
			}}})

			facts, err := reader.ListOrders(context.Background(), "installation-1", 1)
			if err != nil {
				t.Fatalf("ListOrders() error = %v", err)
			}
			if len(facts) != 1 {
				t.Fatalf("ListOrders() returned %d facts, want 1", len(facts))
			}
			if facts[0].RealizationState != tt.want {
				t.Fatalf("RealizationState = %q, want %q", facts[0].RealizationState, tt.want)
			}
		})
	}
}
