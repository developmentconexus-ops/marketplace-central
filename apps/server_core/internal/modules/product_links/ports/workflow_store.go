package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type ProductLinkWorkflowStore interface {
	GetProductLink(ctx context.Context, identity domain.ListingIdentity) (domain.ProductLink, bool, error)
	ApplyProductLinkTransition(ctx context.Context, transition domain.ProductLinkTransition) error
	ListProductLinks(ctx context.Context, installationID string, limit int) ([]domain.ProductLink, error)
	ListProductLinkAuditEntries(ctx context.Context, installationID string, limit int) ([]domain.ProductLinkAuditEntry, error)
}
