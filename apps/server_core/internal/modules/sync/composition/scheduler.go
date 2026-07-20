// Package composition wires the sync module for the application root. It exists
// so composition/root.go depends on a single sync entrypoint (one import, one
// ticker line) instead of reaching into the adapter and application packages
// directly — keeping the M-01 additive-lock edit in root.go minimal.
package composition

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	syncpg "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// InstallationScopeERP is the installation_id used for ERP-sourced entities
// (products via xlsx/sankhya) that are not tied to a marketplace installation.
// It keeps sync_state's (tenant, installation, entity) key well-formed without
// fabricating an ML installation id.
const InstallationScopeERP = "erp"

// NewProductsSkeletonScheduler builds the M-01 sync scheduler with the single
// products entity registered to a NO-OP placeholder job. The placeholder makes
// NO external ML/Oracle call — it returns the cursor unchanged so a cycle records
// last_full_sync_at without heavy work. M-03 (xlsx) / M-04 (sankhya) replace the
// placeholder with the real product-sync body via the same RegisterJob seam,
// without touching the loop.
//
// Cadence-agnostic (D6): interval is injected by the caller; nothing here
// hardcodes a business cadence.
func NewProductsSkeletonScheduler(pool *pgxpool.Pool, tenantID string, interval time.Duration) *syncapp.Scheduler {
	repo := syncpg.NewSyncStateRepository(pool, tenantID)
	scheduler := syncapp.NewScheduler(repo, InstallationScopeERP, interval, time.Now)
	// Deferral (2026-07-20, MIS-006): this NO-OP placeholder job is the only
	// stub on the live composition path. The real product-sync body lands in
	// M-03 (xlsx adapter) / M-04 (sankhya adapter), which re-register products
	// via the same RegisterJob seam. products is a valid entity and the first
	// registration, so RegisterJob cannot fail here; the error is intentionally
	// not surfaced from wiring.
	_ = scheduler.RegisterJob(domain.EntityProducts, func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		return cursor, nil
	})
	return scheduler
}
