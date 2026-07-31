package ports

import (
	"context"
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
)

// CanonicalProductReader is the Oracle/internal_read catalog boundary.
//
// It reads one product. The whole-catalog reads that used to sit beside it
// (ListCanonicalProducts, SearchCanonicalProducts) walked every page with no
// assortment policy and existed only for the deleted legacy routes; the paged
// reads on internal_read's CatalogPageReader are what answers "which products",
// and they carry the cut.
type CanonicalProductReader interface {
	GetCanonicalProduct(context.Context, domain.InternalProductID) (domain.CanonicalProduct, error)
}
