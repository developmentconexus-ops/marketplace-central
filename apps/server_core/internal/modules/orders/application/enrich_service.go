package application

import (
	"context"
	"log/slog"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
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
// never a fabricated date/status/UF. Destination (city/zip/receiver), carrier
// (name/tracking URL) and Costs are all pointer/nil-able honest-absence facts.
type ShipmentEnrichment struct {
	ShipmentID      string
	Status          string
	Substatus       string
	SLADue          *time.Time
	Delayed         *bool
	DestinationUF   *string
	DestinationCity *string
	DestinationZip  *string
	ReceiverName    *string
	CarrierName     *string
	TrackingURL     *string
	Costs           *connectorsdomain.ShipmentCosts
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
	// BuyerFiscal is the buyer's fiscal identity (name, opaque document,
	// billing address) for ERP registration, resolved read-time via
	// BuyerFiscalReader. It is nil whenever the buyer has no billing data
	// (honest absence) or the lookup could not run — never a fabricated block.
	BuyerFiscal *connectorsdomain.BuyerFiscalInfo
}

// EnrichService adds read-time enrichment (masked buyer identity, vinculo
// status, per-item ERP cost, per-order shipment facts) to canonical order
// read models. Buyer/vinculo facts already exist on the read model or its
// items and require no external call. Cost is sourced via CostReader (the
// internal_read cost fact); shipment is sourced via ShipmentReader keyed by
// installationID + the order's ShippingID. Tenant scoping remains the repo
// layer's responsibility.
type EnrichService struct {
	cost        ports.CostReader
	shipment    ports.ShipmentReader
	decomposer  ports.Decomposer
	buyerFiscal ports.BuyerFiscalReader
	logger      *slog.Logger
}

// NewEnrichService keeps constructing an EnrichService with a nil decomposer
// and a nil buyer-fiscal reader (delegating through the full constructor) so it
// yields the honest-empty profitability path and a nil BuyerFiscal — this is
// the C1 contract seam; the hub wires the real readers post-merge.
func NewEnrichService(cost ports.CostReader, shipment ports.ShipmentReader, logger *slog.Logger) EnrichService {
	return NewEnrichServiceWithReaders(cost, shipment, nil, nil, logger)
}

// NewEnrichServiceWithDecomposer keeps the decomposer-only constructor byte
// stable (the planned C2 root.go swap targets it): a nil buyer-fiscal reader,
// honest-unknown profitability when the decomposer is nil.
func NewEnrichServiceWithDecomposer(cost ports.CostReader, shipment ports.ShipmentReader, decomposer ports.Decomposer, logger *slog.Logger) EnrichService {
	return NewEnrichServiceWithReaders(cost, shipment, decomposer, nil, logger)
}

// NewEnrichServiceWithReaders is the full constructor: every dependency may be
// nil (honest-unknown, never a panic) — see resolveProfitability /
// resolveBuyerFiscal. The composition layer wires the real buyer-fiscal adapter
// here.
func NewEnrichServiceWithReaders(cost ports.CostReader, shipment ports.ShipmentReader, decomposer ports.Decomposer, buyerFiscal ports.BuyerFiscalReader, logger *slog.Logger) EnrichService {
	if logger == nil {
		logger = slog.Default()
	}
	return EnrichService{cost: cost, shipment: shipment, decomposer: decomposer, buyerFiscal: buyerFiscal, logger: logger}
}

// Enrich is the LIST path: it fills every read-time fact that is cheap at
// scale (masked buyer, vinculo, per-item cost, per-order shipment,
// profitability) for N orders. It deliberately does NOT resolve buyer fiscal
// identity — that is drawer-only data whose two-step provider lookup (+2
// sequential ML calls per order) would blow the request deadline across a full
// list page (FINDING-M08-LIST-TIMEOUT). Every EnrichedOrder here carries a nil
// BuyerFiscal; use EnrichOne on the detail path to resolve it.
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

// EnrichOne is the DETAIL path: it runs the same base enrichment as Enrich for
// a single order and additionally resolves the buyer's fiscal identity via
// BuyerFiscalReader. Bounding the fiscal lookup to the one order the drawer is
// showing keeps it off the list hot path (FINDING-M08-LIST-TIMEOUT) while still
// surfacing comprador_fiscal on GET /orders/{id}. Degrade semantics are the
// reader's (honest absence / warn-once) — see resolveBuyerFiscal.
func (s EnrichService) EnrichOne(ctx context.Context, installationID string, order domain.OrderReadModel) EnrichedOrder {
	enriched := s.Enrich(ctx, installationID, []domain.OrderReadModel{order})[0]
	enriched.BuyerFiscal = s.resolveBuyerFiscal(ctx, installationID, order)
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
	if info.DestinationCity != nil {
		city := *info.DestinationCity
		buyer.City = &city
	}
	return &ShipmentEnrichment{
		ShipmentID:      info.ID,
		Status:          info.Status,
		Substatus:       info.Substatus,
		SLADue:          info.SLADue,
		Delayed:         info.Delayed,
		DestinationUF:   info.DestinationUF,
		DestinationCity: info.DestinationCity,
		DestinationZip:  info.DestinationZip,
		ReceiverName:    info.ReceiverName,
		CarrierName:     info.CarrierName,
		TrackingURL:     info.TrackingURL,
		Costs:           info.Costs,
	}
}

// resolveBuyerFiscal looks up the buyer's fiscal identity via BuyerFiscalReader
// (keyed by installationID + the order's own ProviderOrderID — the read model
// carries no billing_info.id, so the reader re-fetches the order internally).
// A nil reader (real adapter not wired yet) or an empty ProviderOrderID skips
// the lookup entirely — honest unknown, never a panic. The reader already maps
// a buyer-without-billing (a 404 / undecodable payload) to an empty
// honest-absence value, so an error here is a REAL failure (auth/transient):
// it degrades to nil AND warns once (never silently swallow a provider fault).
// An empty (HasData()==false) value is honest absence: nil, no warn.
func (s EnrichService) resolveBuyerFiscal(ctx context.Context, installationID string, order domain.OrderReadModel) *connectorsdomain.BuyerFiscalInfo {
	if s.buyerFiscal == nil || order.ProviderOrderID == "" {
		return nil
	}
	info, err := s.buyerFiscal.GetBuyerFiscal(ctx, installationID, order.ProviderOrderID)
	if err != nil {
		s.logger.Warn("orders: buyer fiscal lookup failed",
			"installation_id", installationID,
			"provider_order_id", order.ProviderOrderID,
			"error", err,
		)
		return nil
	}
	if !info.HasData() {
		return nil
	}
	return &info
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
