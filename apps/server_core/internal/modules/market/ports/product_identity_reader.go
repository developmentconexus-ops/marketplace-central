package ports

import (
	"context"
	"errors"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// ProductIdentityReader is the market module's read-only view over ERP
// product identity and its linked marketplace listings.
type ProductIdentityReader interface {
	GetLocalIdentity(ctx context.Context, codprod string) (domain.LocalProductIdentity, error)
	LinkedListings(ctx context.Context, codprod string) ([]string, error)
}

// ErrProductNotFound signals the requested codprod has no local identity row.
var ErrProductNotFound = errors.New("PRODUCT_NOT_FOUND")
