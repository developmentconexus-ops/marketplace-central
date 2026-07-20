package domain

import "time"

// OrderReadModel is the canonical, provider-independent order projection.
// Fields unavailable from current storage remain nullable so reads do not
// manufacture facts that the system cannot currently prove.
type OrderReadModel struct {
	ProviderOrderID      string                    `json:"provider_order_id"`
	ProviderCode         string                    `json:"provider_code"`
	Status               string                    `json:"status"`
	ProviderStatusDetail string                    `json:"provider_status_detail"`
	BuyerNickname        *string                   `json:"buyer_nickname"`
	Total                *float64                  `json:"total"`
	Currency             *string                   `json:"currency"`
	Fulfillment          *string                   `json:"fulfillment"`
	NFState              *string                   `json:"nf_state"`
	ShippingID           string                    `json:"shipping_id,omitempty"`
	CreatedAt            *time.Time                `json:"created_at"`
	ProviderCreatedAt    *time.Time                `json:"provider_created_at"`
	ProviderClosedAt     *time.Time                `json:"provider_closed_at"`
	ProviderUpdatedAt    *time.Time                `json:"provider_updated_at"`
	// FaturadoAt is when the operator recorded the order as invoiced
	// ("Foi faturado"), our-DB only. nil = not yet marked invoiced, an
	// honest unknown (ADR-17) -- never inferred from shipment presence.
	// It drives DeriveOrderBucket's Faturar<->Enviar gate and is exposed on the
	// OrderRead contract as an optional field (hub-granted seam shape:
	// faturado_at?: string). omitempty keeps it an honest ADR-17 absence: nil =>
	// key absent (not yet invoiced), a set timestamp => RFC3339 string. The FE
	// drives its bucket UI off the derived `bucket`, but the raw fact is
	// documented in OpenAPI + sdk-runtime so wire == contract == SDK.
	FaturadoAt           *time.Time                `json:"faturado_at,omitempty"`
	Items                []MarketplaceOrderItem    `json:"items"`
	Payments             []MarketplaceOrderPayment `json:"payments"`
	// Tags carries the provider order tags (e.g. "delivered", "paid"). It is a
	// bucket-derivation signal only (DeriveOrderBucket) — Mercado Livre reports
	// delivered on the order tags, not on provider_status — and is not exposed
	// on the transport DTO. Empty/absent means no honest tag signal (ADR-17).
	Tags []string `json:"tags,omitempty"`
}
