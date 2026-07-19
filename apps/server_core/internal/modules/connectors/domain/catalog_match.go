package domain

import "time"

// CatalogMatchInput drives a read-only tier-3 catalog match probe: resolve a
// seller product (by EAN and/or a free-text title query) against the provider
// catalog, its buy-box, and category predictions. No provider writes.
type CatalogMatchInput struct {
	AccountRef ProviderAccountRef
	EAN        string
	Query      string
}

// CatalogMatchSnapshot is the normalized (IC-06) catalog-match observation.
// Provider payloads die in the adapter; unknown facts are null/empty, never
// zero-defaulted (ADR-17). BuyBox is nil when the provider exposes no buy-box
// winner for the resolved product.
type CatalogMatchSnapshot struct {
	CatalogHits     []CatalogHit       `json:"catalog_hits"`
	BuyBox          *BuyBoxSnapshot    `json:"buy_box"`
	DomainDiscovery []DomainPrediction `json:"domain_discovery"`
	FetchedAt       time.Time          `json:"fetched_at"`
	RawProviderRef  any                `json:"raw_provider_ref,omitempty"`
}

// CatalogHit is one provider catalog product matched from the identifier search.
// The search result carries domain_id, not category_id.
type CatalogHit struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

// BuyBoxSnapshot captures the resolved product's buy-box winner. Every field is
// nullable/empty when the provider omits it — the winner block itself is often
// absent for catalog products (observed live 2026-07-13).
type BuyBoxSnapshot struct {
	CategoryID  string   `json:"category_id,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ListingType string   `json:"listing_type,omitempty"`
}

// DomainPrediction is one ranked category prediction from domain discovery.
// The provider returns no confidence score; order is the only ranking signal.
type DomainPrediction struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"`
	DomainID     string `json:"domain_id,omitempty"`
}
