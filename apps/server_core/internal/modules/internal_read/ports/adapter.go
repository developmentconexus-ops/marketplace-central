package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

// SyncResult reports the outcome of a single Sync() pass. It is deliberately
// minimal here — M-02 defines the port shape; M-03 (xlsx) and M-04 (sankhya)
// populate it from real sync passes. Concrete adapters may treat these counts
// as the number of mirror rows upserted (Processed) and rows that failed to
// merge (Errors).
type SyncResult struct {
	// Processed is the count of product rows the pass upserted into the mirror.
	Processed int
	// Errors is the count of rows the pass could not merge.
	Errors int
}

// ProductSourceAdapter unifies the read side (the existing Reader, unchanged)
// with the write/sync side that feeds products_mirror. It is the formal name
// for the shape already satisfied by erp_import's internalread reader and
// internal_read's oracle reader (Eport.1).
//
// The embedded Reader preserves every existing signature verbatim, so read-side
// consumers (pricing, linking, market) do not recompile with behavior changes.
// The new methods are additive:
//
//   - Sync feeds products_mirror. For upload sources it materializes the most
//     recently completed snapshot; for live sources it reads through to the
//     mirror. Concrete implementations land in M-03/M-04 — M-02 only declares
//     the signature and the target table.
//   - Kind reports whether this adapter is upload-snapshot or live-read-through
//     shaped, so the composition can wire caching/scheduling per shape.
type ProductSourceAdapter interface {
	Reader
	Sync(ctx context.Context) (SyncResult, error)
	Kind() sourcekind.SourceKind
}
