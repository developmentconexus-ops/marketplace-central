// Package composition wires the listings module's daily backfill/sweep
// scheduler for the application root — the ADR-030 "segunda instância" of the
// sync scheduler pattern sync/composition.NewProductsScheduler established.
//
// It cannot reuse that pattern verbatim: products data is tenant-wide (one
// active ERP source per tenant, sync/composition.InstallationScopeERP), but a
// listings backfill run is scoped to ONE specific Mercado Livre
// InstallationAccount, and a tenant can have several ML installations.
// sync/application.Scheduler is bound to a single installationID at
// construction and RegisterJob allows only one registration per entity per
// instance, so a single shared scheduler cannot serve N installations of the
// same entity — this package builds one Scheduler per active ML
// installation instead.
package composition

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	listingsconnectors "marketplace-central/apps/server_core/internal/modules/listings/adapters/connectors"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	synccomposition "marketplace-central/apps/server_core/internal/modules/sync/composition"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// mercadoLivreProviderCode is the only provider this scheduler runs the
// listings backfill/sweep job for today. No named constant exists for this
// literal elsewhere in the codebase — root.go's own provider writers
// (priceWriter/listingWriter) use the same bare string.
const mercadoLivreProviderCode = "mercado_livre"

// installationLister enumerates the tenant's installations — same seam
// product_links/composition.LinkCandidateRefresher uses for its own
// per-installation fan-out, implemented by
// *integrationsapp.InstallationService.
type installationLister interface {
	List(ctx context.Context) ([]integrationsdomain.Installation, error)
}

// Schedulers is one *syncapp.Scheduler per active Mercado Livre installation.
type Schedulers []*syncapp.Scheduler

// StartAll launches every scheduler's ticker loop in its own goroutine and
// returns immediately.
func (s Schedulers) StartAll(ctx context.Context) {
	for _, scheduler := range s {
		scheduler := scheduler
		go scheduler.Start(ctx)
	}
}

// Group captures the F-04 listings scheduler composition inputs. Building it
// via NewListingsSchedulers does NOT touch the database — installations are
// resolved lazily, inside StartAll's goroutine, never synchronously during
// composition. This matches every other composed unit in root.go: DB access
// happens on a request or a scheduler tick, never at NewRootRouter
// construction time (proved by
// TestNewRootRouterBuildsWithoutLegacyConnectorCredentials, which builds the
// whole router with a nil pool and must not panic or error).
type Group struct {
	pool          *pgxpool.Pool
	interval      time.Duration
	installations installationLister
	capabilities  *connectorsapp.MarketplaceCapabilityService
	observer      listingsconnectors.SnapshotObserver
	store         listingsports.CompletedPullStore
}

// NewListingsSchedulers captures the composition inputs needed to build one
// daily listings backfill/sweep scheduler — F-03's NewListingsJob registered
// on domain.EntityListings — per active Mercado Livre installation. Nothing
// is resolved or queried here; call StartAll to actually list installations
// and start the per-installation schedulers.
func NewListingsSchedulers(
	pool *pgxpool.Pool,
	interval time.Duration,
	installations installationLister,
	capabilities *connectorsapp.MarketplaceCapabilityService,
	observer listingsconnectors.SnapshotObserver,
	store listingsports.CompletedPullStore,
) Group {
	return Group{
		pool:          pool,
		interval:      interval,
		installations: installations,
		capabilities:  capabilities,
		observer:      observer,
		store:         store,
	}
}

// StartAll resolves the tenant's installations (via the installationLister
// seam) and starts one daily listings backfill/sweep scheduler per active
// Mercado Livre installation found, all inside a background goroutine so it
// never blocks or fails server boot (root.go's single composition call is
// NewListingsSchedulers(...).StartAll(ctx)).
//
// A List failure, or any panic surfaced while resolving it (e.g. an
// unconfigured pool), is logged and yields no schedulers rather than
// crashing the process (same "log, don't crash boot" posture
// sync/composition.NewProductsScheduler's own RegisterJob call takes: its
// error is intentionally not surfaced from wiring). This is a deliberate
// fail-open: an operator restart re-attempts the list; nothing silently
// fabricates a schedule in the meantime.
//
// installations resolved here are a ONE-TIME snapshot taken at this call —
// an ML installation connected AFTER boot will not get its own scheduler
// until the next restart. A live-reconciling variant would need its own
// ticker to re-list periodically, which is out of this feature's scope.
func (g Group) StartAll(ctx context.Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("listings backfill scheduler: panic resolving installations", "panic", r)
			}
		}()
		for _, scheduler := range resolveListingsSchedulers(ctx, g) {
			scheduler := scheduler
			go scheduler.Start(ctx)
		}
	}()
}

// resolveListingsSchedulers is the synchronous fan-out logic StartAll's
// goroutine runs: list installations, filter to Mercado Livre, build one
// Scheduler + registered EntityListings job per match. Split out from
// StartAll so tests can exercise the fan-out directly with a fake
// installationLister (which never touches a real database and so is safe to
// call synchronously in a test).
func resolveListingsSchedulers(ctx context.Context, g Group) Schedulers {
	all, err := g.installations.List(ctx)
	if err != nil {
		slog.Error("listings backfill scheduler: list installations failed", "err", err)
		return nil
	}

	enumerator := listingsconnectors.NewBackfillSource(g.capabilities)
	hydrator := listingsconnectors.NewMultigetHydrator(g.capabilities, g.observer, time.Now)

	var schedulers Schedulers
	for _, installation := range all {
		if installation.ProviderCode != mercadoLivreProviderCode {
			continue
		}
		account := listingsports.InstallationAccount{
			TenantID:          installation.TenantID,
			InstallationID:    installation.InstallationID,
			ProviderCode:      installation.ProviderCode,
			ProviderAccountID: installation.ExternalAccountID,
		}
		scheduler := synccomposition.NewInstallationScheduler(g.pool, installation.TenantID, installation.InstallationID, g.interval)
		// installation.InstallationID is unique per scheduler instance and
		// EntityListings is the only (first) registration on it, so this
		// cannot fail — the failure modes RegisterJob has (unknown entity,
		// duplicate registration on the SAME instance) are proven generically
		// for this entity in scheduler_test.go, not exercised here.
		_ = scheduler.RegisterJob(domain.EntityListings, listingsapp.NewListingsJob(enumerator, hydrator, g.store, account, time.Now))
		schedulers = append(schedulers, scheduler)
	}
	return schedulers
}
