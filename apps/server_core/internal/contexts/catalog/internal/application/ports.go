// Package application holds Catalog's use cases. It depends on interfaces
// declared here and implemented by internal/postgres, so the decision logic is
// testable without a database and the database is replaceable without touching
// a decision.
package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Store persists products. It is the ONLY writer of the catalog schema.
type Store interface {
	// BySourceKey resolves a source address to the product it already maps to.
	BySourceKey(ctx context.Context, k contracts.SourceProductKey) (domain.Product, bool, error)

	// ByIdentifier returns every product of this tenant carrying the identifier.
	// Scoped by tenant in the signature and not only in the query, because a
	// cross-tenant answer here is a data leak, not a bug in a WHERE clause.
	ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error)

	// Insert writes a new product, its identifiers, its source keys and the
	// observation, atomically.
	Insert(ctx context.Context, p domain.Product) error

	// Update writes a new version of an existing product, atomically.
	Update(ctx context.Context, p domain.Product) error
}

// IDFactory mints canonical product identifiers. It is an interface so a test
// can make them predictable; nothing about the value may depend on the source.
type IDFactory interface {
	NewProductID() (domain.ProductID, error)
}
