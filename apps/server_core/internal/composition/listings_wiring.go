package composition

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings"
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
