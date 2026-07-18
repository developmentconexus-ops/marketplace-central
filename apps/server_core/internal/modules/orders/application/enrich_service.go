package application

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// ItemCost is the per-item ERP unit cost resolved for a single order item.
// UnitCost is nil when the cost could not be honestly resolved (ADR-17:
// unknown != zero) — the item's identifier then also appears in
// EnrichedOrder.ComponentesDesconhecidos.
type ItemCost struct {
	ItemIdentifier string
	UnitCost       *float64
}

// EnrichedOrder wraps a canonical read model with derived, read-time-only
// enrichment facts. It is an internal application value — it is never
// marshaled directly onto the HTTP transport contract or the SDK.
type EnrichedOrder struct {
	Order                    domain.OrderReadModel
	Buyer                    domain.MaskedBuyer
	VinculoStatus            domain.VinculoStatus
	ItemCosts                []ItemCost
	ComponentesDesconhecidos []string
}

// EnrichService adds read-time enrichment (masked buyer identity, vinculo
// status, per-item ERP cost) to canonical order read models. Buyer/vinculo
// facts already exist on the read model or its items and require no
// external call. Cost is sourced via CostReader (the internal_read cost
// fact); tenant scoping remains the repo layer's responsibility —
// installationID is accepted for future enrichment concerns but drives no
// query here.
type EnrichService struct {
	cost ports.CostReader
}

func NewEnrichService(cost ports.CostReader) EnrichService {
	return EnrichService{cost: cost}
}

func (s EnrichService) Enrich(ctx context.Context, installationID string, orders []domain.OrderReadModel) []EnrichedOrder {
	enriched := make([]EnrichedOrder, 0, len(orders))
	for _, order := range orders {
		var nickname string
		if order.BuyerNickname != nil {
			nickname = *order.BuyerNickname
		}
		itemCosts, unknown := s.resolveItemCosts(ctx, order)
		enriched = append(enriched, EnrichedOrder{
			Order:                    order,
			Buyer:                    domain.MaskBuyer(nickname),
			VinculoStatus:            domain.DeriveVinculoStatus(order.Items),
			ItemCosts:                itemCosts,
			ComponentesDesconhecidos: unknown,
		})
	}
	return enriched
}

// resolveItemCosts looks up the per-unit ERP cost for every item on order.
// An item is degraded to unknown (nil UnitCost, identifier appended to the
// unknown list) — never a fabricated 0 — whenever it has no linked internal
// product, the order carries no honest effectiveAt date, or the reader
// itself returns an error or a nil Amount. A per-item miss never aborts the
// order.
func (s EnrichService) resolveItemCosts(ctx context.Context, order domain.OrderReadModel) ([]ItemCost, []string) {
	itemCosts := make([]ItemCost, 0, len(order.Items))
	var unknown []string

	effectiveAt, hasEffectiveAt := orderEffectiveAt(order)
	for _, item := range order.Items {
		identifier := itemIdentifier(item)

		if item.InternalProductID == nil || !hasEffectiveAt {
			itemCosts = append(itemCosts, ItemCost{ItemIdentifier: identifier})
			unknown = append(unknown, identifier)
			continue
		}

		cost, err := s.cost.GetCostAsOf(ctx, *item.InternalProductID, effectiveAt)
		if err != nil || cost.Amount == nil {
			itemCosts = append(itemCosts, ItemCost{ItemIdentifier: identifier})
			unknown = append(unknown, identifier)
			continue
		}

		amount := *cost.Amount
		itemCosts = append(itemCosts, ItemCost{ItemIdentifier: identifier, UnitCost: &amount})
	}
	return itemCosts, unknown
}

// orderEffectiveAt derives the honest "as of" date for a cost lookup: the
// first non-nil of ProviderClosedAt, ProviderCreatedAt, CreatedAt. All nil
// means there is no honest date to ask "cost as of when" — callers must not
// query with a zero time.
func orderEffectiveAt(order domain.OrderReadModel) (time.Time, bool) {
	if order.ProviderClosedAt != nil {
		return *order.ProviderClosedAt, true
	}
	if order.ProviderCreatedAt != nil {
		return *order.ProviderCreatedAt, true
	}
	if order.CreatedAt != nil {
		return *order.CreatedAt, true
	}
	return time.Time{}, false
}

// itemIdentifier picks the one stable, non-PII field that lets the UI point
// at a specific unresolved item: SellerSKU when present, else
// ProviderItemID.
func itemIdentifier(item domain.MarketplaceOrderItem) string {
	if item.SellerSKU != "" {
		return item.SellerSKU
	}
	return item.ProviderItemID
}
