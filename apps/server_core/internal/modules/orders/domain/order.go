package domain

import "time"

type LinkQuality string

const (
	LinkQualityResolved   LinkQuality = "resolved"
	LinkQualityRejected   LinkQuality = "rejected"
	LinkQualityConflict   LinkQuality = "conflict"
	LinkQualityUnresolved LinkQuality = "unresolved"
	LinkQualityMissing    LinkQuality = "missing"
)

type ListingIdentity struct {
	InstallationID      string `json:"installation_id"`
	ProviderItemID      string `json:"provider_item_id"`
	ProviderVariationID string `json:"provider_variation_id,omitempty"`
}

type MarketplaceOrder struct {
	InstallationID      string                  `json:"installation_id"`
	ProviderCode        string                  `json:"provider_code"`
	ProviderOrderID     string                  `json:"provider_order_id"`
	ProviderStatus      string                  `json:"provider_status,omitempty"`
	ProviderStatusDetail string                 `json:"provider_status_detail,omitempty"`
	ProviderCreatedAt   *time.Time              `json:"provider_created_at,omitempty"`
	ProviderClosedAt    *time.Time              `json:"provider_closed_at,omitempty"`
	ProviderUpdatedAt   *time.Time              `json:"provider_updated_at,omitempty"`
	FetchedAt           time.Time               `json:"fetched_at"`
	ShippingID          string                  `json:"shipping_id,omitempty"`
	CancellationDetail  string                  `json:"cancellation_detail,omitempty"`
	Tags                []string                `json:"tags,omitempty"`
	RawProviderRef      string                  `json:"raw_provider_ref,omitempty"`
	Items               []MarketplaceOrderItem  `json:"items"`
	Payments            []MarketplaceOrderPayment `json:"payments"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

type MarketplaceOrderItem struct {
	ProviderItemID      string      `json:"provider_item_id"`
	ProviderVariationID string      `json:"provider_variation_id,omitempty"`
	SellerSKU           string      `json:"seller_sku,omitempty"`
	Title               string      `json:"title,omitempty"`
	Quantity            int         `json:"quantity"`
	UnitPrice           *float64    `json:"unit_price,omitempty"`
	SaleFeeAmount       *float64    `json:"sale_fee_amount,omitempty"`
	LinkQuality         LinkQuality `json:"link_quality"`
	InternalProductID   *int        `json:"internal_product_id,omitempty"`
}

type MarketplaceOrderPayment struct {
	ProviderPaymentID string   `json:"provider_payment_id"`
	ProviderStatus    string   `json:"provider_status,omitempty"`
	TransactionAmount *float64 `json:"transaction_amount,omitempty"`
	TotalPaidAmount   *float64 `json:"total_paid_amount,omitempty"`
}

type ImportResult struct {
	InstallationID string             `json:"installation_id"`
	ImportedCount  int                `json:"imported_count"`
	SkippedCount   int                `json:"skipped_count"`
	Items          []MarketplaceOrder `json:"items"`
}

type ItemLink struct {
	Identity          ListingIdentity
	Quality           LinkQuality
	InternalProductID *int
}
