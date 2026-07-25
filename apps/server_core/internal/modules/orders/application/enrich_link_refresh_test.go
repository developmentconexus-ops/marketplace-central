package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

// fakeLinkReader is a scripted ports.LinkReader for the read-time link refresh.
type fakeLinkReader struct {
	links map[string]domain.ItemLink
	err   error
	calls int
}

func (r *fakeLinkReader) ResolveLinks(_ context.Context, _ string, _ []domain.ListingIdentity) (map[string]domain.ItemLink, error) {
	r.calls++
	if r.err != nil {
		return nil, r.err
	}
	return r.links, nil
}

// An order imported BEFORE its anúncio was linked must pick the link up on the
// next read — otherwise its custo (and with it the whole retorno) stays "—"
// forever even though /vinculos already resolved that anúncio.
func TestEnrichRefreshesAStaleItemLinkAndResolvesItsCost(t *testing.T) {
	closedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	costReader := newFakeCostReader()
	costReader.amounts[90033] = floatPtr(234.43)
	productID := 90033
	links := &fakeLinkReader{links: map[string]domain.ItemLink{
		"install-1|MLB6896003832|": {
			Identity:          domain.ListingIdentity{InstallationID: "install-1", ProviderItemID: "MLB6896003832"},
			Quality:           domain.LinkQualityResolved,
			InternalProductID: &productID,
		},
	}}
	order := domain.OrderReadModel{
		ProviderOrderID:  "ord-stale-link",
		ProviderClosedAt: timePtr(closedAt),
		Items: []domain.MarketplaceOrderItem{
			{ProviderItemID: "MLB6896003832", Quantity: 1, LinkQuality: domain.LinkQualityUnresolved},
		},
	}
	svc := NewEnrichService(costReader, newFakeShipmentReader(), testLogger()).WithLinkRefresh(links)

	got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})[0]

	if got.VinculoStatus != domain.VinculoStatusOK {
		t.Fatalf("VinculoStatus = %q, want OK", got.VinculoStatus)
	}
	if len(got.ItemCosts) != 1 || got.ItemCosts[0].UnitCost == nil || *got.ItemCosts[0].UnitCost != 234.43 {
		t.Fatalf("ItemCosts = %+v, want the refreshed link unit cost", got.ItemCosts)
	}
	if got.Profitability.Decomposition.Custo == nil || *got.Profitability.Decomposition.Custo != 234.43 {
		t.Fatalf("Decomposition.Custo = %v, want 234.43", got.Profitability.Decomposition.Custo)
	}
	// The caller read model is untouched — the refresh works on copies.
	if order.Items[0].InternalProductID != nil {
		t.Fatalf("the input order was mutated: %+v", order.Items[0])
	}
	if links.calls != 1 {
		t.Fatalf("ResolveLinks calls = %d, want exactly 1 batched lookup", links.calls)
	}
}

// A refresh that cannot run leaves the stored snapshot exactly as it was — it
// never blanks a link the order already carries.
func TestEnrichKeepsTheStoredLinkWhenTheRefreshFails(t *testing.T) {
	closedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	costReader := newFakeCostReader()
	costReader.amounts[10] = floatPtr(42.5)
	order := domain.OrderReadModel{
		ProviderOrderID:  "ord-refresh-fail",
		ProviderClosedAt: timePtr(closedAt),
		Items: []domain.MarketplaceOrderItem{
			{ProviderItemID: "MLB1", InternalProductID: intPtr(10), Quantity: 1, LinkQuality: domain.LinkQualityResolved},
			{ProviderItemID: "MLB2", Quantity: 1, LinkQuality: domain.LinkQualityUnresolved},
		},
	}
	links := &fakeLinkReader{err: errors.New("links unavailable")}
	svc := NewEnrichService(costReader, newFakeShipmentReader(), testLogger()).WithLinkRefresh(links)

	got := svc.Enrich(context.Background(), "install-1", []domain.OrderReadModel{order})[0]

	if got.Order.Items[0].InternalProductID == nil || *got.Order.Items[0].InternalProductID != 10 {
		t.Fatalf("stored link was lost: %+v", got.Order.Items[0])
	}
	if got.ItemCosts[0].UnitCost == nil || *got.ItemCosts[0].UnitCost != 42.5 {
		t.Fatalf("ItemCosts[0] = %+v, want the stored link cost", got.ItemCosts[0])
	}
	if got.Order.Items[1].InternalProductID != nil {
		t.Fatalf("unlinked item gained a product on a failed refresh: %+v", got.Order.Items[1])
	}
}
