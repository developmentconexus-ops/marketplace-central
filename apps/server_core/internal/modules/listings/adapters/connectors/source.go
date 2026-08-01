package connectors

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// SnapshotObserver receives every provider snapshot page/batch this module's
// readers process, BEFORE it is mapped down to canonical listing rows. The
// canonical row keeps only what the listings table stores, so provider
// identity fields the linker needs (EAN, seller SKU) are dropped at that
// mapping; an observer sees them. Implemented by the product_links import
// service — one provider fetch feeds both the listings catalog and the link
// matcher.
//
// This interface's declaration lives here (rather than moving into
// backfill.go, its only remaining implementer/consumer in this package)
// deliberately: MultigetHydrator (backfill.go) depends on this exact type,
// and keeping the interface's own file stable avoids churn on an unrelated
// file for a type that outlived its original implementer (the now-deleted
// page-based Source/NewSourceWithObserver — F-04, ADR-04 single-writer).
type SnapshotObserver interface {
	AbsorbProviderSnapshots(ctx context.Context, installationID string, snapshots []connectorsdomain.ListingSnapshot) error
}
