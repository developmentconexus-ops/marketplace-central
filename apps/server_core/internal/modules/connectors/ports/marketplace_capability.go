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
)

// knownIdentityAnchors is the identity vocabulary THIS file governs: the
// anchors a connector declares for the product_links candidate generator. D-A
// (D-122) removed `refforn` from it — it is the supplier reference INSIDE the
// ERP (`ZP1704.1.`), no connector declares it, and keeping it here minted one
// permanent never-changing reason per candidate row. It stays a field on the
// ERP side (`erp_import_products.refforn`); what left is this list, not the
// datum.
//
// Deliberately NOT a claim that no marketplace datum ever compares against
// refforn. The catalog matcher in `market/domain/identity_resolver.go` scores
// the ERP `RefForn` against a candidate's `MODEL` attribute under the anchor
// name "refforn" today, and that resolver is wired; it runs on the market
// module's own vocabulary and never reads this list. A marketplace field that
// one day belongs HERE enters as a NEW anchor with its own name, not by
// reusing a term that means something else on the ERP side.
var knownIdentityAnchors = []IdentityAnchor{
	IdentityAnchorSellerSKU,
	IdentityAnchorEAN,
	IdentityAnchorTitle,
	IdentityAnchorMarca,
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
	// ListOrderRefs enumera identidades sem hidratar cada pedido. Existe
	// separado de ListOrders porque o consumidor em lote (orders.ImportService)
	// descarta tudo do snapshot exceto o id e chama o caminho de escrita único
	// (IngestOrder), que refaz a leitura completa — ou seja, a hidratação da
	// enumeração era uma chamada de provider por pedido, por ciclo, jogada fora.
	ListOrderRefs(ctx context.Context, input domain.ListOrdersInput) ([]domain.OrderSearchHit, error)
	ListOrders(ctx context.Context, input domain.ListOrdersInput) ([]domain.OrderSnapshot, error)
	ReadOrder(ctx context.Context, ref domain.ProviderOrderRef) (domain.OrderSnapshot, error)
}

type FeeQuoteReader interface {
	ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error)
}
