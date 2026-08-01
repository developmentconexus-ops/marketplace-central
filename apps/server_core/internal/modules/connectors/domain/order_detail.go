package domain

import (
	"encoding/json"
	"time"
)

// OrderDetail is the FULL /orders/{id} ingest fact (MIS-007/M-03 F-01, IC-03) — a superset of
// OrderSnapshot (capability.go), which stays scoped to the existing read-path list mapping.
// Do not repurpose OrderSnapshot for the ingest write path; this is a deliberately separate
// type so the two call sites (list-read vs. ingest-write) can evolve independently.
//
// Every optional/unknown field is a nil pointer (ADR-17): provider absence is never a
// fabricated zero or empty string.
type OrderDetail struct {
	ProviderCode         string
	ProviderOrderID      string
	ProviderStatus       string
	ProviderStatusDetail string
	ProviderCreatedAt    *time.Time
	ProviderClosedAt     *time.Time
	// ProviderUpdatedAt mirrors the existing baseline (import_service.go normalizeOrders):
	// last_updated, falling back to date_last_updated, when the provider omits last_updated.
	ProviderUpdatedAt *time.Time
	FetchedAt         time.Time

	// ShippingID is nil when the order payload carries no shipping.id — the ingest caller
	// (F-02) must NOT call GetShipmentDetail in that case (feature.md Negative Scenarios).
	ShippingID *string

	// BuyerNickname is the provider's public buyer handle; nil when the payload carries no
	// buyer block.
	BuyerNickname *string

	// CancellationDetail mirrors the existing baseline convention (mapOrder,
	// capability_adapter.go): it duplicates ProviderStatusDetail verbatim. Kept for parity
	// with normalizeOrders' expectations (import_service.go:108-162), not redefined here.
	CancellationDetail string
	Tags               []string

	// PackID groups orders placed together in one cart. Verified live (fact T2/#17,
	// research/external-ml-api-facts.md — "pack_id agrupa pedidos (carrinho)"); nil when the
	// order was not part of a pack.
	PackID *string

	// DateLastUpdatedML is the provider's own `date_last_updated` value, verbatim — never
	// derived, never defaulted to ProviderUpdatedAt's last_updated fallback. This is the
	// IC-06 incremental-sync watermark (orders_marketplace_orders.date_last_updated_ml, 0089).
	DateLastUpdatedML *time.Time

	Items    []OrderDetailItem
	Payments []OrderDetailPayment

	// RawOrder is the /orders/{id} response body MINUS buyer.billing_info (ADR-03,
	// feature.md Negative/Constraints): held in memory for the ingest layer only. Neither
	// `orders_marketplace_orders` (0089) nor `order_shipments` (0088) carry a raw column —
	// F-02 must NOT persist this field, it exists solely so a single GET can serve both the
	// typed DTO and a raw capture without a second round trip.
	RawOrder json.RawMessage
}

// OrderDetailItem is one order_items[] line. SaleFeeUnit is named to carry the T2 contract
// verbatim: the provider's sale_fee is PER UNIT, never a line total — multiply by Quantity
// downstream if a total is needed. Never rename this field SaleFeeTotal.
type OrderDetailItem struct {
	ProviderItemID      string
	ProviderVariationID string
	SellerSKU           string
	EAN                 string
	Title               string
	Quantity            int
	UnitPrice           *float64
	SaleFeeUnit         *float64
}

type OrderDetailPayment struct {
	PaymentID         string
	Status            string
	TransactionAmount *float64
	TotalPaidAmount   *float64
}
