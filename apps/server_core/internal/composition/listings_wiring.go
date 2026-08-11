package composition

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
)

// ListingsWiring is the assembled listings slice. Like CatalogWiring, this
// file cannot name anything under internal/ — the compiler enforces it.
type ListingsWiring struct {
	Module *listings.Module
	Feed   mercadolivre.Bundle
}

// MLBaseURL and MLProviderCode are the resolved vendor identity for the
// listings slice's one supported marketplace. RuleVendorToken
// (internal/arch/scan.go) only permits a vendor name to appear outside
// adapters/ inside a composition root — that is the whole point of a
// composition root, and cmd/listingsingest is not itself exempt the way
// cmd/mlprobe is. So the literals live here, once, and the CLI references
// them by identifier instead of repeating "mercado_livre" and the API host
// as string literals of its own.
const (
	MLBaseURL      = "https://api.mercadolibre.com"
	MLProviderCode = "mercado_livre"
)

// ListingsReprocessWiring is the credential-free half of the listings slice:
// the context plus the channel's offline mapper, and nothing that can reach
// Mercado Livre. A reprocess assembled from WireListings would hold a live
// client it must never use; this one cannot hold it at all.
type ListingsReprocessWiring struct {
	Module *listings.Module
	Mapper port.ListingMapper
}

// WireListingsReprocess assembles that half. It takes no token source and no
// base URL because it makes no call — the only inputs a reprocess has are the
// database and the mapping already compiled into this binary.
func WireListingsReprocess(pool *pgxpool.Pool) ListingsReprocessWiring {
	return ListingsReprocessWiring{Module: listings.New(pool), Mapper: mercadolivre.NewListingMapper()}
}

// WireListings assembles the slice. The token source and account identity are
// the root's decision, passed in from the operator entry point.
func WireListings(pool *pgxpool.Pool, mlBaseURL, mlUserID, accountID string, token func(context.Context) (string, error)) (ListingsWiring, error) {
	bundle, err := mercadolivre.New(mercadolivre.Config{
		BaseURL:   mlBaseURL,
		UserID:    mlUserID,
		Channel:   "mercado_livre",
		AccountID: accountID,
		Token:     token,
	})
	if err != nil {
		return ListingsWiring{}, err
	}
	return ListingsWiring{Module: listings.New(pool), Feed: bundle}, nil
}
