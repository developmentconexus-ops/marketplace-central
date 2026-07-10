package integrations

import (
	"context"
	"reflect"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type stubProviderOrderReader struct {
	items []connectorsdomain.OrderSnapshot
}

func (s stubProviderOrderReader) ListOrders(context.Context, string, int) ([]connectorsdomain.OrderSnapshot, error) {
	return s.items, nil
}

func TestOrderSourceTranslatesConnectorSnapshot(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	closedAt := now.Add(time.Hour)
	unitPrice := 125.92
	saleFee := 11.07
	amount := 130.0
	transactionAmount := 125.92
	totalPaidAmount := 125.92
	source := NewOrderSource(stubProviderOrderReader{items: []connectorsdomain.OrderSnapshot{{
		ProviderCode:         "mercado_livre",
		ProviderOrderID:      "2001",
		ProviderStatus:       "paid",
		ProviderStatusDetail: "accredited",
		ProviderCreatedAt:    &now,
		ProviderClosedAt:     &closedAt,
		ProviderUpdatedAt:    &now,
		FetchedAt:            now,
		RawProviderRef:       map[string]any{"buyer": "not translated"},
		SaleFeeAmount:        &saleFee,
		ShippingID:           "ship-1",
		CancellationDetail:   "none",
		Tags:                 []string{"paid", "delivered"},
		Items: []connectorsdomain.OrderItemSnapshot{{
			ProviderItemID:      "MLB123",
			ProviderVariationID: "VAR1",
			SellerSKU:           "SKU-1",
			EAN:                 "7890000000012",
			Title:               "Produto",
			Quantity:            2,
			UnitPrice:           &unitPrice,
			SaleFeeAmount:       &saleFee,
		}},
		Payments: []connectorsdomain.OrderPaymentSnapshot{{
			PaymentID:         "pay-1",
			Status:            "approved",
			Amount:            &amount,
			TransactionAmount: &transactionAmount,
			TotalPaidAmount:   &totalPaidAmount,
		}},
	}}})

	snapshots, err := source.ListOrders(context.Background(), "inst-1", 1)
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("snapshots = %d, want 1", len(snapshots))
	}

	got := snapshots[0]
	if got.ProviderCode != "mercado_livre" || got.ProviderOrderID != "2001" || got.ProviderStatus != "paid" || got.ProviderStatusDetail != "accredited" {
		t.Fatalf("provider fields = %#v", got)
	}
	if got.ProviderCreatedAt != &now || got.ProviderClosedAt != &closedAt || got.ProviderUpdatedAt != &now || !got.FetchedAt.Equal(now) {
		t.Fatalf("timestamps were not translated: %#v", got)
	}
	if got.SaleFeeAmount != &saleFee || got.ShippingID != "ship-1" || got.CancellationDetail != "none" || !reflect.DeepEqual(got.Tags, []string{"paid", "delivered"}) {
		t.Fatalf("order operational fields = %#v", got)
	}
	item := got.Items[0]
	if item.ProviderItemID != "MLB123" || item.ProviderVariationID != "VAR1" || item.SellerSKU != "SKU-1" || item.EAN != "7890000000012" || item.Title != "Produto" || item.Quantity != 2 || item.UnitPrice != &unitPrice || item.SaleFeeAmount != &saleFee {
		t.Fatalf("item fields = %#v", item)
	}
	payment := got.Payments[0]
	if payment.PaymentID != "pay-1" || payment.Status != "approved" || payment.Amount != &amount || payment.TransactionAmount != &transactionAmount || payment.TotalPaidAmount != &totalPaidAmount {
		t.Fatalf("payment fields = %#v", payment)
	}
}

func TestOrderIngestionSnapshotDoesNotExposeRawProviderReference(t *testing.T) {
	if _, found := reflect.TypeOf(ordersdomain.OrderIngestionSnapshot{}).FieldByName("RawProviderRef"); found {
		t.Fatal("OrderIngestionSnapshot must not expose RawProviderRef")
	}
}
