package integrations

import (
	"context"
	"reflect"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type stubProviderOrderReader struct {
	items []connectorsdomain.OrderSearchHit
	got   connectorsdomain.ListOrdersInput
}

func (s *stubProviderOrderReader) ListOrderRefs(_ context.Context, input connectorsdomain.ListOrdersInput, _ string) ([]connectorsdomain.OrderSearchHit, error) {
	s.got = input
	return s.items, nil
}

// TestOrderSourceTranslatesSearchHits is F-00 Task 3's translation criterion: OrderSource now
// asks the connectors boundary for id-only hits (ListOrderRefs) instead of full snapshots, and
// maps ProviderOrderID/ProviderUpdatedAt straight through — there is no item/payment mapping left
// here because ImportService, the only caller, never consumed those fields off the old snapshot.
func TestOrderSourceTranslatesSearchHits(t *testing.T) {
	updatedAt := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	source := NewOrderSource(&stubProviderOrderReader{items: []connectorsdomain.OrderSearchHit{{
		ProviderOrderID:   "2001",
		ProviderUpdatedAt: &updatedAt,
	}}})

	refs, err := source.ListOrders(context.Background(), ordersports.ListOrdersInput{InstallationID: "inst-1", Limit: 1})
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("refs = %d, want 1", len(refs))
	}

	got := refs[0]
	if got.ProviderOrderID != "2001" {
		t.Fatalf("provider_order_id = %q, want \"2001\"", got.ProviderOrderID)
	}
	if got.ProviderUpdatedAt != &updatedAt {
		t.Fatalf("provider_updated_at was not translated: %#v", got.ProviderUpdatedAt)
	}
}

func TestOrderIngestionSnapshotDoesNotExposeRawProviderReference(t *testing.T) {
	if _, found := reflect.TypeOf(ordersdomain.OrderIngestionSnapshot{}).FieldByName("RawProviderRef"); found {
		t.Fatal("OrderIngestionSnapshot must not expose RawProviderRef")
	}
}
