package composition

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	syncpg "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

// ProductsCursor is the cursor persisted by the products sync job. It records
// what the last cycle actually did, so an operator can tell a working scheduler
// from one that is merely ticking.
type ProductsCursor struct {
	Source      string    `json:"source"`
	Processed   int       `json:"processed"`
	CompletedAt time.Time `json:"completed_at"`
}

// NewProductsScheduler builds the products sync scheduler with the REAL sync
// body registered: each cycle resolves the tenant's active source and runs that
// source's adapter.
//
// This closes the seam M-01 opened and M-03/M-04 filled: the adapters existed
// and were wired to nothing, so no product data ever refreshed on its own. The
// scheduler stays cadence-agnostic (interval injected) and source-agnostic (the
// adapter is chosen by configuration, not by a hardcoded branch here).
//
// Fail-closed, in line with the read path: an unknown active source, a source
// with no registered adapter, or an adapter error all return an error, which
// the scheduler persists to sync_state.last_error. Nothing silently succeeds,
// and no source is assumed by default.
func NewProductsScheduler(
	pool *pgxpool.Pool,
	tenantID string,
	interval time.Duration,
	lookup tenant_config.ActiveSourceLookup,
	adapters map[tenant_config.ActiveSource]readports.ProductSourceAdapter,
	opts ...ProductsJobOption,
) *syncapp.Scheduler {
	repo := syncpg.NewSyncStateRepository(pool, tenantID)
	scheduler := syncapp.NewScheduler(repo, InstallationScopeERP, interval, time.Now)

	// products is a valid entity and the first registration, so RegisterJob
	// cannot fail here; the error is intentionally not surfaced from wiring.
	_ = scheduler.RegisterJob(domain.EntityProducts, NewProductsJob(tenantID, lookup, adapters, time.Now, opts...))
	return scheduler
}

// LinkCandidateRefresher regenerates link candidates after product data
// changed. Product data changing is exactly when the answer to "which product
// is this anúncio?" can change, so the sync that brought the data in is what
// asks the question again (M-05 F-01).
type LinkCandidateRefresher interface {
	RefreshLinkCandidates(ctx context.Context) error
}

type ProductsJobOption func(*productsJobConfig)

type productsJobConfig struct {
	refresher LinkCandidateRefresher
}

// WithLinkCandidateRefresh registers the post-sync link candidate regeneration.
// Absent, the job syncs products and nothing else — the same behavior it had
// before M-05.
func WithLinkCandidateRefresh(refresher LinkCandidateRefresher) ProductsJobOption {
	return func(cfg *productsJobConfig) {
		cfg.refresher = refresher
	}
}

// NewProductsJob builds the products sync body. It is separate from the
// scheduler so the decision it makes each cycle — which adapter runs, and what
// happens when there is none — is testable without a database.
func NewProductsJob(
	tenantID string,
	lookup tenant_config.ActiveSourceLookup,
	adapters map[tenant_config.ActiveSource]readports.ProductSourceAdapter,
	now func() time.Time,
	opts ...ProductsJobOption,
) syncapp.JobFunc {
	if now == nil {
		now = time.Now
	}
	jobCfg := productsJobConfig{}
	for _, opt := range opts {
		opt(&jobCfg)
	}
	return func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		cfg, err := lookup.Get(ctx, tenantID)
		if err != nil {
			return cursor, fmt.Errorf("products sync: resolve active source: %w", err)
		}
		adapter, ok := adapters[cfg.Source]
		if !ok || adapter == nil {
			return cursor, fmt.Errorf("products sync: no adapter registered for active source %q", cfg.Source)
		}
		result, err := adapter.Sync(ctx)
		if err != nil {
			return cursor, fmt.Errorf("products sync (%s): %w", cfg.Source, err)
		}
		// The sync already landed. A refresh that fails is logged and left for
		// the next cycle (or the generation endpoint) — reporting the cycle as
		// failed would claim the product data did not arrive, which is false
		// (M05-C11).
		if jobCfg.refresher != nil {
			if err := jobCfg.refresher.RefreshLinkCandidates(ctx); err != nil {
				slog.ErrorContext(ctx, "link candidate refresh after products sync failed",
					"source", string(cfg.Source),
					"processed", result.Processed,
					"error", err)
			}
		}
		next, err := json.Marshal(ProductsCursor{
			Source:      string(cfg.Source),
			Processed:   result.Processed,
			CompletedAt: now().UTC(),
		})
		if err != nil {
			// The sync itself succeeded; report the cycle as successful and keep
			// the previous cursor rather than inventing one.
			return cursor, nil
		}
		return next, nil
	}
}
