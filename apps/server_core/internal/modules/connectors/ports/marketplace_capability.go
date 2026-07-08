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
	CapabilityOrderRead           = "order_read"
	CapabilityShipmentPlaceholder = "shipment_read"
	CapabilityQuestionPlaceholder = "question_read"
)

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

type OrderReader interface {
	ListOrders(ctx context.Context, input domain.ListOrdersInput) ([]domain.OrderSnapshot, error)
	ReadOrder(ctx context.Context, ref domain.ProviderOrderRef) (domain.OrderSnapshot, error)
}

type FeeQuoteReader interface {
	ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error)
}
