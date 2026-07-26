package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

const (
	CapabilityAccountProbe        = "account_probe"
	CapabilityListingRead         = "listing_read"
	CapabilityFeeQuoteRead        = "fee_quote_read"
	CapabilityStockRead           = "stock_read"
	CapabilityStockWrite          = "stock_write"
	CapabilityPriceWrite          = "price_write"
	CapabilityListingWrite        = "listing_write"
	CapabilityOrderRead           = "order_read"
	CapabilityShipmentPlaceholder = "shipment_read"
	CapabilityQuestionPlaceholder = "question_read"
)

type IdentityAnchor string

const (
	IdentityAnchorSellerSKU IdentityAnchor = "seller_sku"
	IdentityAnchorEAN       IdentityAnchor = "ean"
	IdentityAnchorTitle     IdentityAnchor = "title"
	IdentityAnchorMarca     IdentityAnchor = "marca"
	IdentityAnchorRefforn   IdentityAnchor = "refforn"
)

var knownIdentityAnchors = []IdentityAnchor{
	IdentityAnchorSellerSKU,
	IdentityAnchorEAN,
	IdentityAnchorTitle,
	IdentityAnchorMarca,
	IdentityAnchorRefforn,
}

func KnownIdentityAnchors() []IdentityAnchor {
	return append([]IdentityAnchor(nil), knownIdentityAnchors...)
}

type AccountProber interface {
	ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error)
}

type ListingReader interface {
	ListListings(ctx context.Context, input domain.ListListingsInput) ([]domain.ListingSnapshot, error)
	ReadListing(ctx context.Context, ref domain.ProviderListingRef) (domain.ListingSnapshot, error)
}

type StockReader interface {
	ReadStock(ctx context.Context, ref domain.ProviderListingRef) (domain.StockSnapshot, error)
}

type StockWriter interface {
	UpdateAvailableQuantity(ctx context.Context, request domain.StockWriteRequest) (domain.StockWriteResult, error)
}

// PriceWriter is consumed by the mutations adapter for absolute price writes.
type PriceWriter interface {
	UpdatePrice(ctx context.Context, request domain.PriceWriteRequest) (domain.PriceWriteResult, error)
}

// ListingWriter is consumed by the mutations adapter for pause and edit writes.
type ListingWriter interface {
	UpdateListing(ctx context.Context, request domain.ListingWriteRequest) (domain.ListingWriteResult, error)
}

type OrderReader interface {
	ListOrders(ctx context.Context, input domain.ListOrdersInput) ([]domain.OrderSnapshot, error)
	ReadOrder(ctx context.Context, ref domain.ProviderOrderRef) (domain.OrderSnapshot, error)
}

type FeeQuoteReader interface {
	ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error)
}
