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
	CreatedAt            *time.Time                `json:"created_at"`
	ProviderCreatedAt    *time.Time                `json:"provider_created_at"`
	ProviderClosedAt     *time.Time                `json:"provider_closed_at"`
	ProviderUpdatedAt    *time.Time                `json:"provider_updated_at"`
	Items                []MarketplaceOrderItem    `json:"items"`
	Payments             []MarketplaceOrderPayment `json:"payments"`
}
