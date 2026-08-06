// Package port carries what other contexts may ASK Catalog. A question, never a
// table: a consumer asks "which product is this?" and does not join to
// catalog.products. That is the difference that lets Catalog change its storage
// without breaking anybody, and the reason there is no foreign key across
// context schemas.
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Summary answers "what is this product?", flattened for a consumer that has no
// business knowing Catalog's internal model.
type Summary struct {
	ProductID   string
	Description string
	Identifiers []contracts.Identifier
	Version     int
}

// Reader answers identity questions about products.
type Reader interface {
	// ByProductID returns the summary, and false when no such product exists
	// for this tenant. Not existing is an answer, not an error.
	ByProductID(ctx context.Context, t tenant.ID, productID string) (Summary, bool, error)

	// ByIdentifier returns every product carrying this identifier.
	//
	// The slice is the point. More than one is a real and expected answer,
	// because identifiers are evidence and not keys. A signature returning a
	// single product would force this method to pick, and picking silently is
	// how a duplicated EAN becomes a wrong link nobody can see.
	ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]Summary, error)
}
