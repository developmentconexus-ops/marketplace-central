package domain

import "time"

type StockRiskActionability string

const (
	StockRiskActionabilityActionable StockRiskActionability = "actionable"
	StockRiskActionabilityBlocked    StockRiskActionability = "blocked"
)

type ListingIdentity struct {
	InstallationID      string `json:"installation_id"`
	ProviderItemID      string `json:"provider_item_id"`
	ProviderVariationID string `json:"provider_variation_id,omitempty"`
}

type ListingSnapshot struct {
	Identity          ListingIdentity
	ProviderCode      string
	ProviderStatus    string
	SellerSKU         string
	EAN               string
	Title             string
	AvailableQuantity *int
	ObservedAt        *time.Time
}

type LinkRecord struct {
	Identity              ListingIdentity
	ProviderCode          string
	State                 LinkState
	InternalProductID     *int
	InternalProductName   string
	InternalReferenceCode string
}

type StockRiskListItem struct {
	Identity              ListingIdentity        `json:"identity"`
	ProviderCode          string                 `json:"provider_code"`
	ProviderStatus        string                 `json:"provider_status,omitempty"`
	SellerSKU             string                 `json:"seller_sku,omitempty"`
	EAN                   string                 `json:"ean,omitempty"`
	Title                 string                 `json:"title,omitempty"`
	LinkState             LinkState              `json:"link_state"`
	InternalProductID     *int                   `json:"internal_product_id,omitempty"`
	InternalProductName   string                 `json:"internal_product_name,omitempty"`
	InternalReferenceCode string                 `json:"internal_reference_code,omitempty"`
	State                 StockRiskState         `json:"state"`
	Actionability         StockRiskActionability `json:"actionability"`
	Actionable            bool                   `json:"actionable"`
	InternalQuantity      *int                   `json:"internal_quantity,omitempty"`
	ProviderQuantity      *int                   `json:"provider_quantity,omitempty"`
	RecommendedQuantity   *int                   `json:"recommended_quantity,omitempty"`
	PolicyID              string                 `json:"policy_id"`
	InternalObservedAt    *time.Time             `json:"internal_observed_at,omitempty"`
	ProviderObservedAt    *time.Time             `json:"provider_observed_at,omitempty"`
	BlockingReason        BlockingReason         `json:"blocking_reason,omitempty"`
}

func (r StockRiskRow) IsActionable() bool {
	return r.State == StockRiskOversell || r.State == StockRiskUndersell
}
