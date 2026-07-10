package domain

import "time"

type ListingSnapshot struct {
	InstallationID      string     `json:"installation_id"`
	ProviderCode        string     `json:"provider_code"`
	ProviderItemID      string     `json:"provider_item_id"`
	ProviderVariationID string     `json:"provider_variation_id,omitempty"`
	ProviderStatus      string     `json:"provider_status,omitempty"`
	SellerSKU           string     `json:"seller_sku,omitempty"`
	EAN                 string     `json:"ean,omitempty"`
	Title               string     `json:"title,omitempty"`
	AvailableQuantity   *int       `json:"available_quantity,omitempty"`
	SourceUpdatedAt     *time.Time `json:"source_updated_at,omitempty"`
	FetchedAt           time.Time  `json:"fetched_at"`
	CreatedAt           time.Time  `json:"created_at,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at,omitempty"`
}

type ListingSnapshotImportResult struct {
	InstallationID string            `json:"installation_id"`
	ImportedCount  int               `json:"imported_count"`
	Items          []ListingSnapshot `json:"items"`
}
