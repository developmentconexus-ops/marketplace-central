package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// IngestServiceConfig wires every per-order fetch IC-06 names for IngestOrder: OrderDetail
// (required — the order fetch itself), ShipmentDetail (only called when the order carries a
// shipping.id) and BuyerFiscal (called on EVERY order — see the doc comment on IngestOrder for
// why this differs from EnrichService.Enrich's list-path skip), plus the existing LinkReader
// (SKU resolution, same port ImportService used) and the new single-transaction OrderIngestStore.
type IngestServiceConfig struct {
	OrderDetail    ports.OrderDetailReader
	ShipmentDetail ports.ShipmentDetailReader
	BuyerFiscal    ports.BuyerFiscalReader
	Links          ports.LinkReader
	Store          ports.OrderIngestStore
	Now            func() time.Time
	Logger         *slog.Logger
}

type IngestService struct {
	orderDetail    ports.OrderDetailReader
	shipmentDetail ports.ShipmentDetailReader
	buyerFiscal    ports.BuyerFiscalReader
	links          ports.LinkReader
	store          ports.OrderIngestStore
	now            func() time.Time
	logger         *slog.Logger
}

var _ ports.OrderIngestor = (*IngestService)(nil)

func NewIngestService(cfg IngestServiceConfig) *IngestService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &IngestService{
		orderDetail:    cfg.OrderDetail,
		shipmentDetail: cfg.ShipmentDetail,
		buyerFiscal:    cfg.BuyerFiscal,
		links:          cfg.Links,
		store:          cfg.Store,
		now:            now,
		logger:         logger,
	}
}

// IngestOrder is IC-06's single write path (ADR-04). One order, one atomic transaction, no
// partial writes: every fetch below runs BEFORE the repository call, so any real error from any
// of them returns without touching the store at all.
//
// A 403 (third-party order) or 404 on the ORDER fetch is the one case the brief calls out by
// name (feature.md Negative Scenarios) and gets a typed classification (domain.ErrOrderUnavailable)
// so ImportService can count it as a skip instead of failing the whole run. Every other real
// error (shipment auth/transient, buyer-fiscal auth/transient, link resolution, the DB write
// itself) is NOT separately classified: it propagates as-is and also results in zero writes for
// this order, because nothing is written until every fetch has already succeeded — the sequencing
// itself is what guarantees atomicity here, no special-case handling needed per failure site.
//
// BuyerFiscal is fetched for EVERY order, unlike EnrichService.Enrich's list path (which skips it
// per FINDING-M08-LIST-TIMEOUT — two sequential ML calls per order blowing the list request's
// deadline). IngestOrder runs off the interactive request path entirely (batch import today;
// backfill/incremental sync/webhook worker later, per milestone.md and IC-06's own "roda fora de
// request interativo" framing) so that constraint does not apply, and IC-03's Operations table
// lists billing as a standard per-order GET during ingest. F-03 retires the live fiscal read site
// (EnrichService) entirely, so these columns are the durable source of that fact going forward —
// leaving it conditional here would silently orphan every order enriched only via the list path.
func (s *IngestService) IngestOrder(ctx context.Context, installationID, providerOrderID string) error {
	if s.orderDetail == nil || s.shipmentDetail == nil || s.buyerFiscal == nil || s.links == nil || s.store == nil {
		return errors.New("ORDERS_INGEST_NOT_CONFIGURED")
	}
	installationID = strings.TrimSpace(installationID)
	providerOrderID = strings.TrimSpace(providerOrderID)
	if installationID == "" || providerOrderID == "" {
		return errors.New("ORDERS_INGEST_IDENTIFIER_REQUIRED")
	}

	detail, err := s.orderDetail.GetOrderDetail(ctx, installationID, providerOrderID)
	if err != nil {
		if errors.Is(err, connectorsdomain.ErrUnauthorized) || errors.Is(err, connectorsdomain.ErrNotFound) {
			return fmt.Errorf("%w: installation=%s provider_order_id=%s: %w", domain.ErrOrderUnavailable, installationID, providerOrderID, err)
		}
		return err
	}

	// providerShipmentID is set whenever the order carries a shipping.id, independent of
	// whether the shipment fetch below finds anything — feature.md Negative Scenarios:
	// "shipment 404 -> order ainda persiste, provider_shipment_id preenchido p/ retry futuro".
	providerShipmentID := detail.ShippingID
	var shipment *domain.OrderShipment
	var shipmentStatus string
	if detail.ShippingID != nil {
		sd, err := s.shipmentDetail.GetShipmentDetail(ctx, installationID, *detail.ShippingID)
		if err != nil {
			return err
		}
		if sd.Found {
			shipmentStatus = sd.Status
			shipment = mapOrderShipment(sd, detail.ProviderCode, *detail.ShippingID, providerOrderID, s.now())
		}
		// sd.Found == false is honest-absence (404/410 on the primary shipment call) — no
		// order_shipments row this round, order persists regardless (see below).
	}

	fiscal, err := s.buyerFiscal.GetBuyerFiscal(ctx, installationID, providerOrderID)
	if err != nil {
		return err
	}

	links, err := s.links.ResolveLinks(ctx, installationID, identitiesFromDetail(installationID, detail))
	if err != nil {
		return err
	}

	// faturado_at is a local (Sankhya/operator) fact IngestOrder never writes — only reads, to
	// classify the bucket the same way DeriveOrderBucket does everywhere else in this module.
	faturadoAt, err := s.store.GetFaturadoAt(ctx, installationID, providerOrderID)
	if err != nil {
		return err
	}

	order := mapMarketplaceOrder(installationID, detail, links, fiscal, providerShipmentID, shipmentStatus, faturadoAt, s.now())
	return s.store.IngestOrder(ctx, order, shipment)
}

// identitiesFromDetail mirrors the dedup rules of the retired collectIdentities (import_service.go,
// pre-F-02): same key shape, same skip-empty-ProviderItemID rule, now built from
// connectorsdomain.OrderDetail.Items instead of domain.OrderIngestionItem.
func identitiesFromDetail(installationID string, detail connectorsdomain.OrderDetail) []domain.ListingIdentity {
	seen := map[string]struct{}{}
	identities := make([]domain.ListingIdentity, 0, len(detail.Items))
	for _, item := range detail.Items {
		identity := domain.ListingIdentity{
			InstallationID:      installationID,
			ProviderItemID:      strings.TrimSpace(item.ProviderItemID),
			ProviderVariationID: strings.TrimSpace(item.ProviderVariationID),
		}
		if identity.ProviderItemID == "" {
			continue
		}
		key := keyOf(identity)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		identities = append(identities, identity)
	}
	return identities
}

// mapMarketplaceOrder builds the module-owned write shape from the connectors-layer fetches.
// Bucket is computed here, once, via domain.DeriveOrderBucket verbatim — never reimplemented.
func mapMarketplaceOrder(
	installationID string,
	detail connectorsdomain.OrderDetail,
	links map[string]domain.ItemLink,
	fiscal connectorsdomain.BuyerFiscalInfo,
	providerShipmentID *string,
	shipmentStatus string,
	faturadoAt *time.Time,
	now time.Time,
) domain.MarketplaceOrder {
	order := domain.MarketplaceOrder{
		InstallationID:       installationID,
		ProviderCode:         strings.TrimSpace(detail.ProviderCode),
		ProviderOrderID:      strings.TrimSpace(detail.ProviderOrderID),
		ProviderStatus:       strings.TrimSpace(detail.ProviderStatus),
		ProviderStatusDetail: strings.TrimSpace(detail.ProviderStatusDetail),
		ProviderCreatedAt:    detail.ProviderCreatedAt,
		ProviderClosedAt:     detail.ProviderClosedAt,
		ProviderUpdatedAt:    detail.ProviderUpdatedAt,
		FetchedAt:            detail.FetchedAt,
		ShippingID:           derefString(detail.ShippingID),
		BuyerNickname:        derefString(detail.BuyerNickname),
		CancellationDetail:   strings.TrimSpace(detail.CancellationDetail),
		Tags:                 detail.Tags,
		RawProviderRef:       safeOrderProviderReference(detail.ProviderCode, detail.ProviderOrderID),
		CreatedAt:            now,
		UpdatedAt:            now,

		PackID:             detail.PackID,
		ProviderShipmentID: providerShipmentID,
		Bucket:             domain.DeriveOrderBucket(detail.ProviderStatus, shipmentStatus, detail.Tags, faturadoAt != nil),
		DateLastUpdatedML:  detail.DateLastUpdatedML,

		BuyerName:      fiscal.Name,
		BuyerDocType:   fiscal.DocType,
		BuyerDocNumber: fiscal.DocNumber,
	}
	if fiscal.Address != nil {
		order.BuyerAddressStreet = fiscal.Address.StreetName
		order.BuyerAddressNumber = fiscal.Address.StreetNumber
		order.BuyerAddressCity = fiscal.Address.City
		order.BuyerAddressStateCode = fiscal.Address.StateCode
		order.BuyerAddressZip = fiscal.Address.ZipCode
		order.BuyerAddressCountry = fiscal.Address.CountryID
	}

	for _, payment := range detail.Payments {
		order.Payments = append(order.Payments, domain.MarketplaceOrderPayment{
			ProviderPaymentID: strings.TrimSpace(payment.PaymentID),
			ProviderStatus:    strings.TrimSpace(payment.Status),
			TransactionAmount: payment.TransactionAmount,
			TotalPaidAmount:   payment.TotalPaidAmount,
		})
	}
	for _, item := range detail.Items {
		identity := domain.ListingIdentity{
			InstallationID:      installationID,
			ProviderItemID:      strings.TrimSpace(item.ProviderItemID),
			ProviderVariationID: strings.TrimSpace(item.ProviderVariationID),
		}
		link := links[keyOf(identity)]
		if link.Quality == "" {
			link.Quality = domain.LinkQualityMissing
		}
		order.Items = append(order.Items, domain.MarketplaceOrderItem{
			ProviderItemID:      identity.ProviderItemID,
			ProviderVariationID: identity.ProviderVariationID,
			SellerSKU:           strings.TrimSpace(item.SellerSKU),
			Title:               strings.TrimSpace(item.Title),
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			// SaleFeeAmount stores whatever the fetch handed us, unchanged — SaleFeeUnit is
			// documented PER UNIT (order_detail.go T2 contract), and no downstream consumer in
			// this feature's scope asked for a *Quantity total; that transform, if ever wanted,
			// is a read-side/M-06 decision, not F-02's to make (ADR-17: never fabricate beyond
			// the source fact).
			SaleFeeAmount:     item.SaleFeeUnit,
			LinkQuality:       link.Quality,
			InternalProductID: link.InternalProductID,
		})
	}
	return order
}

// mapOrderShipment converts the connectors-layer fact into the module-owned persistence shape,
// parsing the decimal-string cost fields into plain float64 (nullableFloat8 at the repository,
// same convention as every other money column in this module). A cost string that fails to
// parse degrades to nil (ADR-17: never a fabricated amount) with a warn — it does not fail the
// whole ingest, since cost is supplementary to the shipment fact, not load-bearing for it.
func mapOrderShipment(sd connectorsdomain.ShipmentDetail, providerCode, providerShipmentID, providerOrderID string, now time.Time) *domain.OrderShipment {
	return &domain.OrderShipment{
		ProviderCode:       strings.TrimSpace(providerCode),
		ProviderShipmentID: strings.TrimSpace(providerShipmentID),
		ProviderOrderID:    strings.TrimSpace(providerOrderID),

		Status:    strings.TrimSpace(sd.Status),
		Substatus: sd.Substatus,

		LogisticType:   sd.LogisticType,
		TrackingMethod: sd.TrackingMethod,
		TrackingNumber: sd.TrackingNumber,

		SLAStatus:  sd.SLAStatus,
		SLALimitAt: sd.SLALimitAt,

		CostGross:  parseDecimalPtr(sd.CostGross),
		CostSeller: parseDecimalPtr(sd.CostSeller),
		Currency:   sd.Currency,

		ReceiverName: sd.ReceiverName,
		DestCity:     sd.DestCity,
		DestState:    sd.DestState,
		DestZip:      sd.DestZip,

		SourceTime: sd.SourceTime,
		FetchedAt:  sd.FetchedAt,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func parseDecimalPtr(value *string) *float64 {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// safeOrderProviderReference is the retired normalizeOrders' helper (import_service.go,
// pre-F-02), moved here verbatim: it is IngestOrder's call site now, the only one left.
func safeOrderProviderReference(providerCode, providerOrderID string) string {
	if strings.TrimSpace(providerCode) != "mercado_livre" {
		return ""
	}
	providerOrderID = strings.TrimSpace(providerOrderID)
	if providerOrderID == "" {
		return ""
	}
	return "/orders/" + url.PathEscape(providerOrderID)
}

// keyOf is the retired collectIdentities/normalizeOrders' helper (import_service.go, pre-F-02),
// moved here verbatim: identity resolution now happens entirely inside IngestOrder.
func keyOf(identity domain.ListingIdentity) string {
	return identity.InstallationID + "|" + identity.ProviderItemID + "|" + identity.ProviderVariationID
}
