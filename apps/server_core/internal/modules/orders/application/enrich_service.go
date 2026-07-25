package application

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// shipmentConcurrency bounds how many per-order shipment lookups (a live
// network GET each) run in parallel during Enrich. Bounded, not unbounded,
// so a large list page cannot fan out into an unbounded burst against the
// provider.
const shipmentConcurrency = 8

// ItemCost is the per-item ERP unit cost resolved for a single order item.
// UnitCost is nil when the cost could not be honestly resolved (ADR-17:
// unknown != zero) — the item's identifier then also appears in
// EnrichedOrder.ComponentesDesconhecidos.
type ItemCost struct {
	ItemIdentifier string
	UnitCost       *float64
	// ObservedAt is the instant the source actually observed this cost, and
	// Approximate says that instant is NOT the order date: the ERP source holds
	// a single snapshot, not a cost history, so the amount is the closest honest
	// observation rather than an exact as-of answer. Both stay nil/false when the
	// source answered exactly at the order date.
	ObservedAt  *time.Time
	Approximate bool
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
	links       ports.LinkReader
	logger      *slog.Logger
}

// WithLinkRefresh wires the read-time link refresher. The item↔produto link on
// a stored order is a SNAPSHOT of what was known when the order was imported:
// an order imported before its anúncio was linked keeps link_quality=unresolved
// forever, so its custo — and therefore the whole retorno — stayed "—" even
// after the operator resolved that anúncio in /vinculos. Re-reading the current
// links at enrich time makes an order pick up a link the moment it exists,
// with no re-import. One batched lookup per Enrich call, not one per item.
// A nil reader keeps the stored snapshot exactly as before.
func (s EnrichService) WithLinkRefresh(links ports.LinkReader) EnrichService {
	s.links = links
	return s
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
	orders = s.refreshLinks(ctx, installationID, orders)
	shipments := make([]*ShipmentEnrichment, len(orders))
	var g errgroup.Group
	g.SetLimit(shipmentConcurrency)
	for index, order := range orders {
		i, order := index, order
		g.Go(func() error {
			shipments[i] = s.fetchShipment(ctx, installationID, order)
			return nil
		})
	}
	_ = g.Wait()

	enriched := make([]EnrichedOrder, 0, len(orders))
	for i, order := range orders {
		var nickname string
		if order.BuyerNickname != nil {
			nickname = *order.BuyerNickname
		}
		itemCosts, unknown := s.resolveItemCosts(ctx, order)
		buyer := domain.MaskBuyer(nickname)
		shipmentInfo := shipments[i]
		applyShipmentToBuyer(shipmentInfo, &buyer)
		profitability := s.resolveProfitability(ctx, order, itemCosts, shipmentInfo)
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

// refreshLinks re-resolves each item's produto interno against the CURRENT
// links and returns copies of the orders carrying whatever is newer. It only
// ever UPGRADES an item — a resolved link that the stored snapshot missed —
// and never downgrades one on a lookup miss: an identity absent from the
// refresh map keeps the stored value rather than being blanked. A reader
// error degrades to the stored snapshot plus a warn; the orders still render
// (their custo simply stays honest-unknown). Copies are made per order so the
// caller's read models — and any cache behind them — stay untouched.
func (s EnrichService) refreshLinks(ctx context.Context, installationID string, orders []domain.OrderReadModel) []domain.OrderReadModel {
	if s.links == nil || len(orders) == 0 {
		return orders
	}
	identities := make([]domain.ListingIdentity, 0)
	seen := map[string]struct{}{}
	for _, order := range orders {
		for _, item := range order.Items {
			if item.ProviderItemID == "" || item.InternalProductID != nil {
				continue
			}
			identity := domain.ListingIdentity{
				InstallationID:      installationID,
				ProviderItemID:      item.ProviderItemID,
				ProviderVariationID: item.ProviderVariationID,
			}
			key := linkRefreshKey(identity)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			identities = append(identities, identity)
		}
	}
	if len(identities) == 0 {
		return orders
	}
	links, err := s.links.ResolveLinks(ctx, installationID, identities)
	if err != nil {
		s.logger.Warn("orders: link refresh failed",
			"installation_id", installationID,
			"identities", len(identities),
			"error", err,
		)
		return orders
	}

	refreshed := make([]domain.OrderReadModel, len(orders))
	copy(refreshed, orders)
	for i := range refreshed {
		var items []domain.MarketplaceOrderItem
		for j, item := range refreshed[i].Items {
			if item.ProviderItemID == "" || item.InternalProductID != nil {
				continue
			}
			link, ok := links[linkRefreshKey(domain.ListingIdentity{
				InstallationID:      installationID,
				ProviderItemID:      item.ProviderItemID,
				ProviderVariationID: item.ProviderVariationID,
			})]
			if !ok || link.InternalProductID == nil {
				continue
			}
			if items == nil {
				items = make([]domain.MarketplaceOrderItem, len(refreshed[i].Items))
				copy(items, refreshed[i].Items)
			}
			items[j].InternalProductID = link.InternalProductID
			items[j].LinkQuality = link.Quality
		}
		if items != nil {
			refreshed[i].Items = items
		}
	}
	return refreshed
}

func linkRefreshKey(identity domain.ListingIdentity) string {
	return identity.InstallationID + "|" + identity.ProviderItemID + "|" + identity.ProviderVariationID
}

// fetchShipment looks up the order's shipment via ShipmentReader. It never
// fabricates a value: an empty ShippingID skips the reader entirely (not an
// error), and a reader error degrades to a nil ShipmentEnrichment plus a
// structured warn — the order itself is still returned. A nil ShipmentReader
// (real adapter not wired yet) is honest-unknown too: the lookup is skipped
// exactly as an empty ShippingID would be, never a panic. It is pure: it never
// mutates a buyer — see applyShipmentToBuyer — so it is safe to call
// concurrently across orders as long as each call writes only its own result.
func (s EnrichService) fetchShipment(ctx context.Context, installationID string, order domain.OrderReadModel) *ShipmentEnrichment {
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

// applyShipmentToBuyer fills buyer.UF/City from the shipment's honest
// DestinationUF/DestinationCity when present. A nil info leaves buyer
// untouched (honest absence, never fabricated).
func applyShipmentToBuyer(info *ShipmentEnrichment, buyer *domain.MaskedBuyer) {
	if info == nil {
		return
	}
	if info.DestinationUF != nil {
		uf := *info.DestinationUF
		buyer.UF = &uf
	}
	if info.DestinationCity != nil {
		city := *info.DestinationCity
		buyer.City = &city
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
// Decomposer when one is wired. When it is nil (real pricing-engine adapter
// not wired yet), it falls back to domain.BuildProfitability assembled from
// the real facts already resolved in this same enrich pass (comissão from
// the order's own item sale fees, custo from the ERP unit costs, frete from
// the shipment's sender cost) — taxa_fixa/difal/tarifa_full/imposto stay
// honest-unknown until an engine is wired (ADR-17: unknown != zero, never
// fabricated). C1 does not populate Difal.UFRoute from destino: a route
// needs both origin (seller UF, a tenant/engine fact not on the orders base)
// and destino, so a decomposer-less path leaves UFRoute nil for a future
// slice (which has the engine's tenant origin) to fill.
func (s EnrichService) resolveProfitability(ctx context.Context, order domain.OrderReadModel, itemCosts []ItemCost, shipment *ShipmentEnrichment) domain.OrderProfitability {
	if s.decomposer != nil {
		if p, ok := s.decomposer.Decompose(ctx, order); ok {
			return p
		}
	}
	return domain.BuildProfitability(domain.ProfitabilityInputs{
		Total:    order.Total,
		Comissao: sumSaleFee(order.Items),
		Custo:    sumItemCosts(itemCosts, order.Items),
		Frete:    senderFreight(shipment),
		Imposto:  nil, // not persisted on the order today — honest unknown (ADR-17); do NOT compute
	})
}

// sumSaleFee sums the order's per-line SaleFeeAmount (already the per-line
// fee, not per-unit — the connector mapper stores it raw without ×qty) into
// the order comissão. Any nil-fee line or an empty item list makes the whole
// sum honest-unknown (ADR-17: never a partial fabricated total).
func sumSaleFee(items []domain.MarketplaceOrderItem) *float64 {
	if len(items) == 0 {
		return nil
	}
	var sum float64
	for _, item := range items {
		if item.SaleFeeAmount == nil {
			return nil
		}
		sum += *item.SaleFeeAmount
	}
	return &sum
}

// sumItemCosts sums the ERP per-unit cost (positionally zipped with items via
// itemCosts, built in the same order) × quantity into the order custo. Any
// unresolved unit cost or an empty item list makes the whole sum
// honest-unknown (ADR-17: never a partial fabricated total).
func sumItemCosts(itemCosts []ItemCost, items []domain.MarketplaceOrderItem) *float64 {
	if len(items) == 0 {
		return nil
	}
	var sum float64
	for i := range items {
		if i >= len(itemCosts) || itemCosts[i].UnitCost == nil {
			return nil
		}
		sum += *itemCosts[i].UnitCost * float64(items[i].Quantity)
	}
	return &sum
}

// senderFreight extracts the seller's actual freight cost — the shipment's
// SenderCost — as the order frete. Any missing shipment/costs/sender-cost or
// an unparseable amount is honest-unknown (ADR-17), never a fabricated 0.
func senderFreight(info *ShipmentEnrichment) *float64 {
	if info == nil || info.Costs == nil || info.Costs.SenderCost == nil {
		return nil
	}
	v, err := strconv.ParseFloat(info.Costs.SenderCost.Amount, 64)
	if err != nil {
		return nil
	}
	return &v
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
		resolved := ItemCost{ItemIdentifier: identifier, UnitCost: &amount}
		// A stale_source answer is a real observation taken at another instant —
		// carry the caveat instead of dropping the amount or hiding the gap.
		if internalreaddomain.HasQualityFlag(cost.QualityFlags, internalreaddomain.QualityStaleSource) {
			resolved.Approximate = true
		}
		if cost.Source.ObservedAt != nil {
			observedAt := *cost.Source.ObservedAt
			resolved.ObservedAt = &observedAt
		}
		itemCosts = append(itemCosts, resolved)
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
