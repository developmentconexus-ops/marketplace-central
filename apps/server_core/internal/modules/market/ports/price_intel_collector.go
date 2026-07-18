package ports

import (
	"context"
	"errors"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// PriceIntelCollector is the market module's read-only view over marketplace
// price intelligence (catalog identity, catalog offers, own-item pricing,
// price-to-win). Real adapters are wired at the composition root; this
// module depends only on the interface.
type PriceIntelCollector interface {
	SearchCatalogByEAN(ctx context.Context, ean string) (CatalogIdentityResult, error)
	GetCatalogProduct(ctx context.Context, catalogProductID string) (CatalogProductDetail, error)
	ListCatalogOffers(ctx context.Context, catalogProductID string) ([]domain.ValidatedOffer, error)
	GetOwnItemPricing(ctx context.Context, listingID string) (OwnItemPricing, error)
	GetPriceToWin(ctx context.Context, listingID string) (PriceToWinSignal, error)
}

type CatalogIdentityResult struct {
	Candidates []domain.IdentityCandidate
	FetchedAt  time.Time
}

type CatalogProductDetail struct {
	CatalogProductID  string
	BuyBoxWinnerPrice *domain.Money
	BuyBoxSellerID    string
	FetchedAt         time.Time
}

type OwnItemPricing struct {
	ListingID   string
	OurPrice    *domain.Money
	WinnerPrice *domain.Money
	FetchedAt   time.Time
}

type PriceToWinSignal struct {
	ListingID   string
	TargetPrice *domain.Money
	Position    *domain.MarketPosition
	FetchedAt   time.Time
}

var (
	// ErrCatalogOffersDisabled signals the catalog-offers read is unavailable
	// (feature flag OFF or the route is forbidden); the pipeline maps it to
	// an explicit NO_PRICE_EVIDENCE aggregate, never a 500.
	ErrCatalogOffersDisabled = errors.New("CATALOG_OFFERS_DISABLED")
	// ErrRateLimited signals a provider backoff condition; the pipeline stops
	// remaining provider work for the codprod rather than retrying.
	ErrRateLimited = errors.New("PROVIDER_RATE_LIMITED")
)

// ProviderStatusError classifies a non-rate-limited provider failure by HTTP
// status so the collection pipeline can map it to a PROVIDER_4XX/PROVIDER_5XX
// cause without depending on the connectors module (import cycle).
type ProviderStatusError struct {
	StatusCode int
}

func (e *ProviderStatusError) Error() string {
	return fmt.Sprintf("provider responded with status %d", e.StatusCode)
}
