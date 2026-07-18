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
	// InsertBatch persists the module-owned batch audit row for a
	// completed ApplyBatch run (S3). Additive: existing callers of the
	// other methods on this interface are unaffected.
	InsertBatch(ctx context.Context, batch domain.ProductLinkBatch) error
	// GetAuditEntry loads a single audit row by its AuditID (S4 undo target
	// lookup). found=false (nil error) when no such row exists for this
	// tenant.
	GetAuditEntry(ctx context.Context, auditID string) (domain.ProductLinkAuditEntry, bool, error)
	// LatestAuditForIdentity returns the most recent audit row for a given
	// listing identity (S4 undo ordering: decides SUPERSEDED/ALREADY_UNDONE
	// deterministically from the audit chain, never from time/random).
	LatestAuditForIdentity(ctx context.Context, identity domain.ListingIdentity) (domain.ProductLinkAuditEntry, bool, error)
	// ListAuditByBatch returns every audit row tagged with the given
	// batch_id, newest first (S4 batch-undo fan-out).
	ListAuditByBatch(ctx context.Context, batchID string) ([]domain.ProductLinkAuditEntry, error)
}
