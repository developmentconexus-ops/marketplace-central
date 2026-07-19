package application

import (
	"context"
	"log/slog"
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

// ShipmentEnrichment carries the per-order shipment facts sourced via
// ShipmentReader. It is nil on EnrichedOrder whenever the shipment could not
// be honestly read (no shipping_id, reader error) — ADR-17: unknown != zero,
// never a fabricated date/status/UF. Rastreio is represented by ShipmentID
// plus Status only; the provider exposes no tracking-URL/code field.
type ShipmentEnrichment struct {
	ShipmentID    string
	Status        string
	Substatus     string
	SLADue        *time.Time
	Delayed       *bool
	DestinationUF *string
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
	Shipment                 *ShipmentEnrichment
	Profitability            domain.OrderProfitability
}

// EnrichService adds read-time enrichment (masked buyer identity, vinculo
// status, per-item ERP cost, per-order shipment facts) to canonical order
// read models. Buyer/vinculo facts already exist on the read model or its
// items and require no external call. Cost is sourced via CostReader (the
// internal_read cost fact); shipment is sourced via ShipmentReader keyed by
// installationID + the order's ShippingID. Tenant scoping remains the repo
// layer's responsibility.
type EnrichService struct {
	cost       ports.CostReader
	shipment   ports.ShipmentReader
	decomposer ports.Decomposer
	logger     *slog.Logger
}

// NewEnrichService keeps constructing an EnrichService with a nil decomposer
// (delegating to NewEnrichServiceWithDecomposer) so it yields the honest-empty
// profitability path — this is the C1 contract seam; the hub wires a real
// decomposer via NewEnrichServiceWithDecomposer in C2 post-merge.
func NewEnrichService(cost ports.CostReader, shipment ports.ShipmentReader, logger *slog.Logger) EnrichService {
	return NewEnrichServiceWithDecomposer(cost, shipment, nil, logger)
}

// NewEnrichServiceWithDecomposer is the full constructor: decomposer may be
// nil (honest-unknown profitability, never a panic) — see resolveProfitability.
func NewEnrichServiceWithDecomposer(cost ports.CostReader, shipment ports.ShipmentReader, decomposer ports.Decomposer, logger *slog.Logger) EnrichService {
	if logger == nil {
		logger = slog.Default()
	}
	return EnrichService{cost: cost, shipment: shipment, decomposer: decomposer, logger: logger}
}

func (s EnrichService) Enrich(ctx context.Context, installationID string, orders []domain.OrderReadModel) []EnrichedOrder {
	enriched := make([]EnrichedOrder, 0, len(orders))
	for _, order := range orders {
		var nickname string
		if order.BuyerNickname != nil {
			nickname = *order.BuyerNickname
		}
		itemCosts, unknown := s.resolveItemCosts(ctx, order)
		buyer := domain.MaskBuyer(nickname)
		shipmentInfo := s.resolveShipment(ctx, installationID, order, &buyer)
		profitability := s.resolveProfitability(ctx, order)
		enriched = append(enriched, EnrichedOrder{
			Order:                    order,
			Buyer:                    buyer,
			VinculoStatus:            domain.DeriveVinculoStatus(order.Items),
			ItemCosts:                itemCosts,
			ComponentesDesconhecidos: unknown,
			Shipment:                 shipmentInfo,
			Profitability:            profitability,
		})
	}
	return enriched
}

// resolveShipment looks up the order's shipment via ShipmentReader and, when
// found, fills buyer.UF from the shipment's honest DestinationUF. It never
// fabricates a value: an empty ShippingID skips the reader entirely (not an
// error), and a reader error degrades to a nil ShipmentEnrichment plus a
// structured warn — the order itself is still returned. A nil ShipmentReader
// (real adapter not wired yet) is honest-unknown too: the lookup is skipped
// exactly as an empty ShippingID would be, never a panic.
func (s EnrichService) resolveShipment(ctx context.Context, installationID string, order domain.OrderReadModel, buyer *domain.MaskedBuyer) *ShipmentEnrichment {
	if s.shipment == nil || order.ShippingID == "" {
		return nil
	}
	info, err := s.shipment.GetShipment(ctx, installationID, order.ShippingID)
	if err != nil {
		s.logger.Warn("orders: shipment lookup failed",
			"installation_id", installationID,
			"provider_order_id", order.ProviderOrderID,
			"shipping_id", order.ShippingID,
			"error", err,
		)
		return nil
	}
	if info.DestinationUF != nil {
		uf := *info.DestinationUF
		buyer.UF = &uf
	}
	return &ShipmentEnrichment{
		ShipmentID:    info.ID,
		Status:        info.Status,
		Substatus:     info.Substatus,
		SLADue:        info.SLADue,
		Delayed:       info.Delayed,
		DestinationUF: info.DestinationUF,
	}
}

// resolveProfitability looks up the order's decomposição/DIFAL/retorno via
// Decomposer. A nil Decomposer (real adapter not wired yet — C1) or a
// not-ok result both degrade to domain.UnknownOrderProfitability(), the
// honest-empty value (ADR-17: unknown != zero, never fabricated) — never a
// panic. C1 does not populate Difal.UFRoute from destino: a route needs
// both origin (seller UF, a tenant/engine fact not on the orders base) and
// destino, so a decomposer-less path leaves UFRoute nil for C2 (which has
// the engine's tenant origin) to fill.
func (s EnrichService) resolveProfitability(ctx context.Context, order domain.OrderReadModel) domain.OrderProfitability {
	if s.decomposer == nil {
		return domain.UnknownOrderProfitability()
	}
	if p, ok := s.decomposer.Decompose(ctx, order); ok {
		return p
	}
	return domain.UnknownOrderProfitability()
}

// resolveItemCosts looks up the per-unit ERP cost for every item on order.
// An item is degraded to unknown (nil UnitCost, identifier appended to the
// unknown list) — never a fabricated 0 — whenever it has no linked internal
// product, the order carries no honest effectiveAt date, the CostReader is
// not wired yet (nil, real adapter is a later slice), or the reader itself
// returns an error or a nil Amount. A per-item miss never aborts the order.
func (s EnrichService) resolveItemCosts(ctx context.Context, order domain.OrderReadModel) ([]ItemCost, []string) {
	itemCosts := make([]ItemCost, 0, len(order.Items))
	var unknown []string

	effectiveAt, hasEffectiveAt := orderEffectiveAt(order)
	for _, item := range order.Items {
		identifier := itemIdentifier(item)

		if s.cost == nil || item.InternalProductID == nil || !hasEffectiveAt {
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
