package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

const (
	CapabilityListingRead         = "listing_read"
	CapabilityStockRead           = "stock_read"
	CapabilityStockWrite          = "stock_write"
	CapabilityOrderRead           = "order_read"
	CapabilityShipmentPlaceholder = "shipment_read"
	CapabilityQuestionPlaceholder = "question_read"
)

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
