package composition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cataloginternalread "marketplace-central/apps/server_core/internal/modules/catalog/adapters/internalread"
	catalogpostgres "marketplace-central/apps/server_core/internal/modules/catalog/adapters/postgres"
	catalogapp "marketplace-central/apps/server_core/internal/modules/catalog/application"
	catalogports "marketplace-central/apps/server_core/internal/modules/catalog/ports"
	catalogtransport "marketplace-central/apps/server_core/internal/modules/catalog/transport"
	classpostgres "marketplace-central/apps/server_core/internal/modules/classifications/adapters/postgres"
	classapp "marketplace-central/apps/server_core/internal/modules/classifications/application"
	classtransport "marketplace-central/apps/server_core/internal/modules/classifications/transport"
	melhorenvio "marketplace-central/apps/server_core/internal/modules/connectors/adapters/melhorenvio"
	mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorstransport "marketplace-central/apps/server_core/internal/modules/connectors/transport"
	dashboardapp "marketplace-central/apps/server_core/internal/modules/dashboard/application"
	dashboardtransport "marketplace-central/apps/server_core/internal/modules/dashboard/transport"
	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
	erpproductlinks "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/productlinks"
	erpsync "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/sync"
	erpxlsx "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/xlsx"
	erpapp "marketplace-central/apps/server_core/internal/modules/erp_import/application"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	erptransport "marketplace-central/apps/server_core/internal/modules/erp_import/transport"
	integrationscrypto "marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto"
	integrationsfeesync "marketplace-central/apps/server_core/internal/modules/integrations/adapters/feesync"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsbg "marketplace-central/apps/server_core/internal/modules/integrations/background"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	integrationstransport "marketplace-central/apps/server_core/internal/modules/integrations/transport"
	internalreadcache "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache"
	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	internalreadoracle "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	routing "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/routing"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadobservability "marketplace-central/apps/server_core/internal/modules/internal_read/observability"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	inventoryintegrations "marketplace-central/apps/server_core/internal/modules/inventory/adapters/integrations"
	inventoryinternalread "marketplace-central/apps/server_core/internal/modules/inventory/adapters/internalread"
	inventorymutations "marketplace-central/apps/server_core/internal/modules/inventory/adapters/mutations"
	inventorypostgres "marketplace-central/apps/server_core/internal/modules/inventory/adapters/postgres"
	inventoryproductlinks "marketplace-central/apps/server_core/internal/modules/inventory/adapters/productlinks"
	inventoryapp "marketplace-central/apps/server_core/internal/modules/inventory/application"
	inventoryports "marketplace-central/apps/server_core/internal/modules/inventory/ports"
	inventorytransport "marketplace-central/apps/server_core/internal/modules/inventory/transport"
	listingsconnectors "marketplace-central/apps/server_core/internal/modules/listings/adapters/connectors"
	listingsintegrations "marketplace-central/apps/server_core/internal/modules/listings/adapters/integrations"
	listingsinternalread "marketplace-central/apps/server_core/internal/modules/listings/adapters/internalread"
	listingsmarketplaces "marketplace-central/apps/server_core/internal/modules/listings/adapters/marketplaces"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingscomposition "marketplace-central/apps/server_core/internal/modules/listings/composition"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	listingstransport "marketplace-central/apps/server_core/internal/modules/listings/transport"
	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	marketapp "marketplace-central/apps/server_core/internal/modules/market/application"
	markettransport "marketplace-central/apps/server_core/internal/modules/market/transport"
	marketplacespostgres "marketplace-central/apps/server_core/internal/modules/marketplaces/adapters/postgres"
	marketplacesapp "marketplace-central/apps/server_core/internal/modules/marketplaces/application"
	marketplacesregistry "marketplace-central/apps/server_core/internal/modules/marketplaces/registry"
	marketplacestransport "marketplace-central/apps/server_core/internal/modules/marketplaces/transport"
	mutationsconnectors "marketplace-central/apps/server_core/internal/modules/mutations/adapters/connectors"
	mutationslistings "marketplace-central/apps/server_core/internal/modules/mutations/adapters/listings"
	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	mutationsproductlinks "marketplace-central/apps/server_core/internal/modules/mutations/adapters/productlinks"
	mutationsstub "marketplace-central/apps/server_core/internal/modules/mutations/adapters/stub"
	mutationsapp "marketplace-central/apps/server_core/internal/modules/mutations/application"
	mutationsbg "marketplace-central/apps/server_core/internal/modules/mutations/background"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
	mutationstransport "marketplace-central/apps/server_core/internal/modules/mutations/transport"
	ordersintegrations "marketplace-central/apps/server_core/internal/modules/orders/adapters/integrations"
	ordersinternalread "marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread"
	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	orderspricingtax "marketplace-central/apps/server_core/internal/modules/orders/adapters/pricingtax"
	ordersproductlinks "marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	orderstransport "marketplace-central/apps/server_core/internal/modules/orders/transport"
	pricingcatalog "marketplace-central/apps/server_core/internal/modules/pricing/adapters/catalog"
	pricingcostread "marketplace-central/apps/server_core/internal/modules/pricing/adapters/costread"
	pricingmarket "marketplace-central/apps/server_core/internal/modules/pricing/adapters/marketplace"
	pricingpostgres "marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres"
	pricingtariffcomposite "marketplace-central/apps/server_core/internal/modules/pricing/adapters/tariffcomposite"
	pricingtariffdefaults "marketplace-central/apps/server_core/internal/modules/pricing/adapters/tariffdefaults"
	pricingtarifflive "marketplace-central/apps/server_core/internal/modules/pricing/adapters/tarifflive"
	pricingapp "marketplace-central/apps/server_core/internal/modules/pricing/application"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
	pricingtransport "marketplace-central/apps/server_core/internal/modules/pricing/transport"
	productlinksconnectors "marketplace-central/apps/server_core/internal/modules/product_links/adapters/connectors"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinkscomposition "marketplace-central/apps/server_core/internal/modules/product_links/composition"
	productlinkstransport "marketplace-central/apps/server_core/internal/modules/product_links/transport"
	profitabilityinternalread "marketplace-central/apps/server_core/internal/modules/profitability/adapters/internalread"
	profitabilityorders "marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders"
	profitabilitypostgres "marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres"
	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitytransport "marketplace-central/apps/server_core/internal/modules/profitability/transport"
	syncpg "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	synccomposition "marketplace-central/apps/server_core/internal/modules/sync/composition"
	synctransport "marketplace-central/apps/server_core/internal/modules/sync/transport"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
	tenantconfigtransport "marketplace-central/apps/server_core/internal/modules/tenant_config/transport"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	"marketplace-central/apps/server_core/internal/platform/pgdb"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authFlowFacade struct {
	flow        *integrationsapp.AuthFlowService
	feeSync     *integrationsapp.FeeSyncService
	operations  *integrationsapp.OperationService
	providerOps *integrationsapp.ProviderOperationService
}

type unavailableProductMatcher struct {
	err error
}

func (m unavailableProductMatcher) FindProductsForLinking(context.Context, internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	return nil, m.err
}

func (f authFlowFacade) StartAuthorize(ctx context.Context, input integrationsapp.StartAuthorizeInput) (integrationsapp.AuthorizeStart, error) {
	return f.flow.StartAuthorize(ctx, input)
}

func (f authFlowFacade) HandleCallback(ctx context.Context, input integrationsapp.HandleCallbackInput) (integrationsapp.AuthStatus, error) {
	return f.flow.HandleCallback(ctx, input)
}

func (f authFlowFacade) SubmitAPIKey(ctx context.Context, input integrationsapp.SubmitAPIKeyInput) (integrationsapp.AuthStatus, error) {
	return f.flow.SubmitAPIKey(ctx, input)
}

func (f authFlowFacade) RefreshCredential(ctx context.Context, input integrationsapp.RefreshCredentialInput) (integrationsapp.AuthStatus, error) {
	return f.flow.RefreshCredential(ctx, input)
}

func (f authFlowFacade) Disconnect(ctx context.Context, input integrationsapp.DisconnectInput) (integrationsapp.AuthStatus, error) {
	return f.flow.Disconnect(ctx, input)
}

func (f authFlowFacade) StartReauth(ctx context.Context, input integrationsapp.StartReauthInput) (integrationsapp.AuthorizeStart, error) {
	return f.flow.StartReauth(ctx, input)
}

func (f authFlowFacade) GetAuthStatus(ctx context.Context, input integrationsapp.GetAuthStatusInput) (integrationsapp.AuthStatus, error) {
	return f.flow.GetAuthStatus(ctx, input)
}

func (f authFlowFacade) StartSync(ctx context.Context, input integrationsapp.StartFeeSyncInput) (integrationsapp.FeeSyncAccepted, error) {
	return f.feeSync.StartSync(ctx, input)
}

func (f authFlowFacade) ListOperationRuns(ctx context.Context, installationID string) ([]integrationsdomain.OperationRun, error) {
	return f.operations.ListByInstallation(ctx, installationID)
}

func (f authFlowFacade) ProbeAccount(ctx context.Context, installationID string) (connectorsdomain.AccountSnapshot, error) {
	return f.providerOps.ProbeAccount(ctx, installationID)
}

func (f authFlowFacade) ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error) {
	return f.providerOps.ListListings(ctx, installationID, limit)
}

func (f authFlowFacade) ListOrders(ctx context.Context, installationID string, limit int) ([]connectorsdomain.OrderSnapshot, error) {
	return f.providerOps.ListOrders(ctx, installationID, limit)
}

func (f authFlowFacade) ReadFeeQuote(ctx context.Context, installationID string, input connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	return f.providerOps.ReadFeeQuote(ctx, installationID, input)
}

func (f authFlowFacade) ReadStock(ctx context.Context, installationID string, ref connectorsdomain.ProviderListingRef) (connectorsdomain.StockSnapshot, error) {
	return f.providerOps.ReadStock(ctx, installationID, ref)
}

func (f authFlowFacade) ProbeCatalogMatch(ctx context.Context, installationID string, input integrationsapp.CatalogMatchProbeInput) (integrationsapp.CatalogMatchProbeResult, error) {
	return f.providerOps.ProbeCatalogMatch(ctx, installationID, input)
}

type RootRuntime struct {
	Handler        http.Handler
	CatalogReader  catalogtransport.CatalogReader
	PoolStats      *internalreadobservability.PoolStatsLoop
	MutationPoller *mutationsbg.Poller
	mutationLane   mutationLane
}

type mutationLane struct {
	real, price, stock, listing, linkage, resync bool
	envelope                                     inventoryports.StockMutationEnvelope
	writer                                       mutationsports.WriterPort
}

type installationStockWriter struct {
	capabilities  *connectorsapp.MarketplaceCapabilityService
	installations *integrationsapp.InstallationService
	tenantID      string
}

func (w installationStockWriter) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	installation, found, err := w.installations.Get(ctx, item.InstallationID)
	if err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	if !found {
		return mutationsports.WriteOutcome{}, fmt.Errorf("mutation installation %q not found", item.InstallationID)
	}
	return mutationsconnectors.NewStockWriter(w.capabilities, w.tenantID, installation.ProviderCode, installation.ExternalAccountID).Apply(ctx, item)
}

type installationResyncWriter struct {
	hydrator      listingsports.BatchHydrator
	store         listingsports.CompletedPullStore
	installations *integrationsapp.InstallationService
	tenantID      string
}

func (w installationResyncWriter) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	installation, found, err := w.installations.Get(ctx, item.InstallationID)
	if err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	if !found {
		return mutationsports.WriteOutcome{}, fmt.Errorf("mutation installation %q not found", item.InstallationID)
	}
	account := listingsports.InstallationAccount{TenantID: w.tenantID, InstallationID: installation.InstallationID, ProviderCode: installation.ProviderCode, ProviderAccountID: installation.ExternalAccountID}
	ingestor := listingsapp.NewListingIngestor(w.hydrator, w.store, account, time.Now)
	return mutationslistings.NewResyncWriter(ingestor, account).Apply(ctx, item)
}

func NewRootRouter(pool *pgxpool.Pool, cfg pgdb.Config) (http.Handler, error) {
	runtime, err := NewRootRuntime(pool, cfg)
	if err != nil {
		return nil, err
	}
	return runtime.Handler, nil
}

func registerBatchRoutes(mux *httpx.RouteClassMux) {
	for _, pattern := range []string{
		"/profitability/margin-inputs/import",
		"/profitability/profit-snapshots/calculate",
		"/orders/import",
		"/product-links/listing-snapshots/imports",
		"/product-links/link-candidates/generations",
		"/pricing/simulations/batch",
		"/admin/fee-schedules/sync",
		"/admin/fee-schedules/seed",
	} {
		mux.RegisterRouteClass(pattern, httpx.BatchRouteClass)
	}
}

func NewRootRuntime(pool *pgxpool.Pool, cfg pgdb.Config) (*RootRuntime, error) {
	observabilityCfg, err := internalreadobservability.LoadConfig(os.Getenv)
	if err != nil {
		return nil, fmt.Errorf("observability config: %w", err)
	}
	freshnessCache := internalreadcache.New(internalreadcache.ConfigFromEnv(os.Getenv))

	mux := httpx.NewRouteClassMux()
	registerBatchRoutes(mux)

	base := httpx.NewRouter()
	mux.Handle("/healthz", base)

	erpRepo := erppostgres.NewRepository(pool)
	erpImportSvc := erpapp.NewImportService(erpxlsx.NewParser(), erpRepo, cfg.DefaultTenantID)
	erpQuerySvc := erpapp.NewQueryService(erpRepo, cfg.DefaultTenantID)
	erptransport.NewHandler(erpImportSvc, erpQuerySvc).Register(mux)

	classRepo := classpostgres.NewRepository(pool, cfg.DefaultTenantID)
	classSvc := classapp.NewService(classRepo, cfg.DefaultTenantID)
	classtransport.NewHandler(classSvc).Register(mux)

	providerRepo := integrationspostgres.NewProviderDefinitionRepository(pool)
	providerSvc := integrationsapp.NewProviderService(providerRepo)
	providerRegistry := integrationsproviders.NewRegistry()
	if err := providerRegistry.Validate(); err != nil {
		return nil, fmt.Errorf("integration provider registry: %w", err)
	}

	if pool != nil {
		if err := providerSvc.SeedProviderDefinitions(context.Background(), providerRegistry.All()); err != nil {
			slog.Warn("integration provider definitions seed failed", "err", err)
		}
	}

	installationRepo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	installationSvc := integrationsapp.NewInstallationService(installationRepo, cfg.DefaultTenantID)
	mutationRepo := mutationspostgres.NewRepository(pool, cfg.DefaultTenantID)
	mutationEnvelope := inventorymutations.NewEnvelope(mutationRepo)
	credentialRepo := integrationspostgres.NewCredentialRepository(pool, cfg.DefaultTenantID)
	credentialSvc := integrationsapp.NewCredentialService(credentialRepo, cfg.DefaultTenantID)
	authSessionRepo := integrationspostgres.NewAuthSessionRepository(pool, cfg.DefaultTenantID)
	authSvc := integrationsapp.NewAuthService(authSessionRepo, cfg.DefaultTenantID)
	capabilityStateRepo := integrationspostgres.NewCapabilityStateRepository(pool, cfg.DefaultTenantID)
	capabilitySvc := integrationsapp.NewCapabilityService(capabilityStateRepo, cfg.DefaultTenantID)
	operationRunRepo := integrationspostgres.NewOperationRunRepository(pool, cfg.DefaultTenantID)
	// Runs execute in this process, so anything still queued/running at boot died
	// with the previous one. Closing them out here is what keeps a restart in the
	// middle of a refresh from permanently blocking that operation: BeginExclusive
	// would otherwise keep answering "já em andamento" against a run nobody owns.
	if pool != nil {
		reapCtx, cancelReap := context.WithTimeout(context.Background(), 10*time.Second)
		reaped, reapErr := operationRunRepo.FailInterruptedRuns(reapCtx, time.Now().UTC())
		cancelReap()
		if reapErr != nil {
			slog.Default().Error("integrations: could not close interrupted operation runs", "error", reapErr)
		} else if reaped > 0 {
			slog.Default().Warn("integrations: closed operation runs interrupted by a restart", "runs", reaped)
		}
	}
	operationSvc := integrationsapp.NewOperationService(operationRunRepo, cfg.DefaultTenantID)
	oauthStateRepo := integrationspostgres.NewOAuthStateRepository(pool, cfg.DefaultTenantID)
	feeRepo := marketplacespostgres.NewFeeScheduleRepository(pool)
	registeredFeeSyncers := integrationsproviders.BuildFeeSyncers()
	feeSyncExecutor := integrationsfeesync.NewMarketplaceExecutor(feeRepo, registeredFeeSyncers)
	feeSyncSvc := integrationsapp.NewFeeSyncService(integrationsapp.FeeSyncServiceConfig{
		Installations: installationSvc,
		Providers:     providerSvc,
		Operations:    operationSvc,
		Capabilities:  capabilitySvc,
		Executor:      feeSyncExecutor,
	})

	encryptionSvc, err := integrationscrypto.NewLocalKeyService(cfg.EncryptionKey, "local-key-v1")
	if err != nil {
		return nil, fmt.Errorf("encryption service: %w", err)
	}
	oauthStateCodec, err := integrationsapp.NewHMACOAuthStateCodec(cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("oauth state codec: %w", err)
	}

	authFlowSvc, err := integrationsapp.NewAuthFlowService(integrationsapp.AuthFlowConfig{
		TenantID:        cfg.DefaultTenantID,
		Installations:   installationSvc,
		Credentials:     credentialSvc,
		AuthSessions:    authSvc,
		OAuthStates:     oauthStateRepo,
		OAuthStateCodec: oauthStateCodec,
		Encryptor:       encryptionSvc,
		Adapters:        integrationsproviders.BuildAuthAdapters(),
	})
	if err != nil {
		return nil, fmt.Errorf("auth flow service: %w", err)
	}

	credentialResolver := integrationsapp.NewCredentialResolver(credentialSvc, encryptionSvc)
	mercadoLivreCapabilities := mercadolivreconnector.NewCapabilityAdapter(mercadolivreconnector.CapabilityAdapterConfig{
		AccessTokenResolver: func(ctx context.Context, ref connectorsdomain.ProviderAccountRef) (string, error) {
			resolved, err := credentialResolver.ResolveAccessToken(ctx, integrationsapp.CredentialResolutionRef{
				TenantID:          ref.TenantID,
				InstallationID:    ref.InstallationID,
				ProviderCode:      ref.ProviderCode,
				ProviderAccountID: ref.ProviderAccountID,
			})
			if err != nil {
				return "", err
			}
			return resolved.AccessToken, nil
		},
		CatalogOffersEnabled: mlCatalogOffersEnabled(os.Getenv),
	})
	marketplaceCapabilities := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{
		mercadoLivreCapabilities.ProviderCapabilitySet(),
	})
	productLinkIdentityAnchorReader := productlinksconnectors.NewIdentityAnchorAdapter(marketplaceCapabilities)
	installationSvc.SetStateReconciler(integrationsapp.NewInstallationStateReconciler(integrationsapp.InstallationStateReconcilerConfig{
		TenantID:      cfg.DefaultTenantID,
		Installations: installationSvc,
		AuthSessions:  authSvc,
		Credentials:   credentialSvc,
		Decryptor:     encryptionSvc,
		Capabilities:  marketplaceCapabilities,
	}))
	installationSvc.SetRuntimeCapabilitySource(integrationsapp.NewRuntimeCapabilityProjector(marketplaceCapabilities))
	installationSvc.SetCapabilityStateResolver(capabilitySvc)
	providerOperationSvc := integrationsapp.NewProviderOperationService(integrationsapp.ProviderOperationServiceConfig{
		TenantID:         cfg.DefaultTenantID,
		Installations:    installationSvc,
		Capabilities:     marketplaceCapabilities,
		CapabilityStates: capabilitySvc,
		Operations:       operationSvc,
		CatalogMatch:     mercadoLivreCapabilities,
	})

	flowFacade := authFlowFacade{
		flow:        authFlowSvc,
		feeSync:     feeSyncSvc,
		operations:  operationSvc,
		providerOps: providerOperationSvc,
	}

	integrationstransport.NewHandler(providerSvc, installationSvc).Register(mux)
	integrationstransport.NewAuthHandler(flowFacade).Register(mux)
	integrationstransport.NewRunReadHandler(operationSvc).Register(mux)

	productLinkSnapshotRepo := productlinkspostgres.NewListingSnapshotRepository(pool, cfg.DefaultTenantID)
	productLinkCandidateRepo := productlinkspostgres.NewLinkCandidateRepository(pool, cfg.DefaultTenantID)
	productLinkImportSvc := productlinksapp.NewImportService(productlinksapp.ImportServiceConfig{
		Source: providerOperationSvc,
		Store:  productLinkSnapshotRepo,
	})

	productMatcher := productlinksapp.ProductMatcher(unavailableProductMatcher{
		err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle product linking reader is unavailable", nil),
	})
	var inventoryStockReader inventoryports.InternalStockBatchReader = inventoryinternalread.UnavailableStockBatchReader{
		Err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle inventory stock reader is unavailable", nil),
	}
	// Kept apart from inventoryStockReader so the routing reader below can offer
	// it as the live-source path while still serving the upload sources.
	var liveStockBatchReader internalreadports.StockBatchReader
	oracleBatchSemaphore := oraclebatch.NewSemaphore(oracleBatchPermits(os.Getenv))
	var internalReadSvc internalreadapp.Service
	var internalReadAvailable bool
	var oracleDB internalreadoracle.Database
	var poolStats *internalreadobservability.PoolStatsLoop
	var liveReader internalreadports.Source
	if oracleCfg, err := internalreadoracle.LoadConfigFromEnv(os.Getenv); err != nil {
		slog.Warn("product links oracle reader unavailable", "err", err)
	} else if db, err := internalreadoracle.OpenDB(context.Background(), oracleCfg); err != nil {
		// ReadError.Error() only renders its message, so the wrapped cause — the
		// scrubbed Oracle error code, which is the whole diagnosis — never reaches
		// the log unless it is unwrapped here.
		slog.Warn("product links oracle connection failed", "err", err, "cause", errors.Unwrap(err))
	} else {
		oracleDB = db
		cachedReader := internalreadcache.NewReader(
			internalreadoracle.NewReader(oracleDB),
			freshnessCache,
		)
		liveReader = internalreadobservability.NewTimingReader(
			cachedReader,
			slog.Default(),
			observabilityCfg.SlowQueryThreshold,
		)
		poolStats = internalreadobservability.NewPoolStatsLoop(oracleDB, slog.Default(), observabilityCfg.PoolStatsInterval)
		liveStockBatchReader = internalreadcache.NewStockBatchReader(
			internalreadoracle.NewBatchStockReader(oracleDB, oracleBatchSemaphore),
			freshnessCache,
		)
		inventoryStockReader = liveStockBatchReader
	}
	// Product sync wiring: the scheduler below needs the active-source lookup and
	// one adapter per configured source. Both are only available once the pool
	// exists, so they are declared here and filled in the pool-scoped block.
	var activeSourceLookup tenant_config.ActiveSourceLookup
	productSyncAdapters := map[tenant_config.ActiveSource]internalreadports.ProductSourceAdapter{}
	if pool != nil {
		xlsxReader := erpinternalread.NewReader(erpRepo, cfg.DefaultTenantID)
		cachedUploadReader := internalreadcache.NewReader(xlsxReader, freshnessCache)
		uploadReader := internalreadobservability.NewTimingReader(
			cachedUploadReader,
			slog.Default(),
			observabilityCfg.SlowQueryThreshold,
		)
		activeSourceRepo := tenant_config.NewRepository(pool)
		routingReader := routing.NewReader(uploadReader, liveReader, activeSourceRepo, cfg.DefaultTenantID)
		internalReadSvc = internalreadapp.NewService(routingReader)
		tenantconfigtransport.NewHandler(activeSourceRepo, cfg.DefaultTenantID).Register(mux)
		internalReadAvailable = true
		// Linking reads the mirror for EVERY source, including the live ones. The
		// routing reader is right for a single product question but wrong for
		// candidate generation: that asks about every listing in an account, and one
		// live ERP round-trip per listing cannot finish inside a request. The mirror
		// is the same data the live sync writes, so the answer is unchanged — only
		// the number of round-trips is. Its own reader instance keeps the linking
		// snapshot out of the catalog reader's caches.
		productMatcher = routing.NewMirrorMatcher(
			erpinternalread.NewReader(erpRepo, cfg.DefaultTenantID),
			activeSourceRepo,
			cfg.DefaultTenantID,
		)

		// Stock Seguro asks for the stock of every linked listing at once. Only the
		// oracle path had a batch reader, so upload-source tenants saw "o ERP não
		// informou o estoque" on every row while the mirror held the quantity; the
		// mirror reader closes that gap for the sources it owns.
		inventoryStockReader = routing.NewStockBatchReader(
			erpinternalread.NewBatchStockReader(erpinternalread.NewReader(erpRepo, cfg.DefaultTenantID)),
			liveStockBatchReader,
			activeSourceRepo,
			cfg.DefaultTenantID,
		)

		activeSourceLookup = activeSourceRepo
		// Both upload sources mirror their own protocol history; the reader they
		// carry is the upload reader because that is the source they represent.
		productSyncAdapters[tenant_config.SourceXLSX] = erpxlsx.NewUploadAdapter(uploadReader, erpRepo, cfg.DefaultTenantID, erpdomain.SourceXLSX)
		productSyncAdapters[tenant_config.SourceCatalogoCliente] = erpxlsx.NewUploadAdapter(uploadReader, erpRepo, cfg.DefaultTenantID, erpdomain.SourceCatalogoCliente)
		if oracleDB != nil {
			// Registered only when the ERP connection actually exists: an absent
			// adapter makes the scheduler fail closed with a named error instead of
			// pretending a Sankhya sync ran.
			productSyncAdapters[tenant_config.SourceSankhya] = internalreadoracle.NewSankhyaAdapter(oracleDB, mirror.NewPgWriter(pool), cfg.DefaultTenantID)
		}
	}
	var canonicalCatalogReader catalogports.CanonicalProductReader = cataloginternalread.UnavailableReader{Err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle catalog reader is unavailable", nil)}
	if internalReadAvailable {
		canonicalCatalogReader = cataloginternalread.NewReader(internalReadSvc)
	}
	catalogEnrichments := catalogpostgres.NewEnrichmentRepository(pool, cfg.DefaultTenantID)
	legacyReader := canonicalCatalogReader.(catalogports.ProductReader)
	catalogSvc := catalogapp.NewServiceWithInvalidator(legacyReader, catalogEnrichments, cfg.DefaultTenantID, freshnessCache)
	var catalogPageReader catalogtransport.CatalogReader
	if internalReadAvailable {
		catalogPageReader = internalReadSvc
	}
	catalogtransport.Handler{Service: catalogapp.NewCanonicalService(canonicalCatalogReader), CompatibilityService: catalogSvc, PageReader: catalogPageReader}.Register(mux)

	// Resolution is built first because generation hands it the corroborated
	// candidates to auto-approve (M-05 F-02): the single automatic path runs
	// through the same transition/audit machine as an operator approval.
	productLinkResolutionSvc := productlinksapp.NewResolutionService(productlinksapp.ResolutionServiceConfig{
		Candidates: productLinkCandidateRepo,
		Workflows:  productLinkCandidateRepo,
	})
	productLinkGenerationSvc := productlinksapp.NewGenerationService(productlinksapp.GenerationServiceConfig{
		Snapshots:       productLinkSnapshotRepo,
		Matcher:         productMatcher,
		Store:           productLinkCandidateRepo,
		IdentityAnchors: productLinkIdentityAnchorReader,
		AutoApprover:    productLinkResolutionSvc,
	})
	productLinkSummarySvc := productlinksapp.NewSummaryService(productlinkspostgres.NewSummaryReader(pool, cfg.DefaultTenantID))
	productLinkBatchSvc := productlinksapp.NewBatchService(productlinksapp.BatchServiceConfig{
		Candidates: productLinkCandidateRepo,
		Workflows:  productLinkCandidateRepo,
		Approver:   productLinkResolutionSvc,
		Now:        time.Now,
	})
	productlinkstransport.NewHandlerWithBatchPreview(productLinkImportSvc, productLinkGenerationSvc, productLinkGenerationSvc, productLinkResolutionSvc, productLinkResolutionSvc, productLinkBatchSvc).Register(mux)
	erpImportSvc.SetPostImportHooks(
		erpsync.NewMarketEnqueuer(installationSvc, syncpg.NewSyncStateRepository(pool, cfg.DefaultTenantID)),
		erpproductlinks.NewLinkCandidateGenerator(productLinkGenerationSvc),
	)

	inventoryActionRepo := inventorypostgres.NewStockActionRepository(pool, cfg.DefaultTenantID)
	inventoryRiskSvc := inventoryapp.NewStockRiskService(
		inventoryproductlinks.NewSnapshotReader(productLinkSnapshotRepo),
		inventoryproductlinks.NewLinkReader(productLinkCandidateRepo, productLinkCandidateRepo),
		inventoryStockReader,
	)
	inventoryActionSvc := inventoryapp.NewStockActionServiceWithEnvelope(inventoryActionRepo, mutationEnvelope)
	inventoryManualAction := inventoryapp.NewManualActionFacade(
		inventoryRiskSvc,
		inventoryActionSvc,
		inventoryintegrations.NewInstallationReader(installationSvc),
		inventoryproductlinks.NewSnapshotStore(productLinkSnapshotRepo),
	)
	inventorytransport.NewHandler(inventoryRiskSvc, inventoryManualAction).Register(mux)

	ordersRepo := orderspostgres.NewOrderRepository(pool, cfg.DefaultTenantID)
	// One link reader for BOTH the import (write-time snapshot) and the enrich
	// read path (refresh): an order imported before its anúncio was linked
	// picks the link up on the next read instead of waiting for a re-import.
	ordersLinkReader := ordersproductlinks.NewLinkReader(productLinkCandidateRepo, productLinkCandidateRepo)
	ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
	// F-02: IngestOrder is the single write path (ADR-04) — Import now only enumerates
	// provider_order_ids from ListOrders and delegates each one here.
	ordersOrderDetailReader := newOrdersOrderDetailReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
	ordersShipmentDetailReader := newOrdersShipmentDetailReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
	ordersIngestSvc := ordersapp.NewIngestService(ordersapp.IngestServiceConfig{
		OrderDetail:    ordersOrderDetailReader,
		ShipmentDetail: ordersShipmentDetailReader,
		BuyerFiscal:    ordersBuyerFiscalReader,
		Links:          ordersLinkReader,
		Store:          ordersRepo,
	})
	ordersImportSvc := ordersapp.NewImportService(ordersapp.ImportServiceConfig{
		Source:   ordersintegrations.NewOrderSource(providerOperationSvc),
		Ingestor: ordersIngestSvc,
	})
	ordersListSvc := ordersapp.NewListService(ordersRepo)
	ordersReadRepo := orderspostgres.NewOrderReadRepository(pool, cfg.DefaultTenantID)
	ordersReadSvc := ordersapp.NewReadService(ordersReadRepo)
	ordersSummarySvc := ordersapp.NewSummaryServiceWithBuckets(ordersRepo, ordersRepo)
	ordersCostReader := newOrdersCostReaderAdapter(internalReadSvc, internalReadAvailable)
	// F-03 (read-path switch): EnrichService's shipment/buyer-fiscal reads come
	// from the Postgres readers below (order_shipments/0088,
	// orders_marketplace_orders buyer_* columns/0089 — both populated by F-02's
	// IngestOrder), not a live Mercado Livre call. ordersBuyerFiscalReader above
	// stays live-ML-backed on purpose: it feeds ordersIngestSvc.BuyerFiscal for
	// batch ingest (a distinct call path from the interactive GET /orders/{id}
	// enrich path wired here).
	ordersEnrichShipmentReader := orderspostgres.NewShipmentReader(pool, cfg.DefaultTenantID)
	ordersEnrichBuyerFiscalReader := orderspostgres.NewBuyerFiscalReader(pool, cfg.DefaultTenantID)
	// ICMSSaida/Difal/PisCofins/RestituicaoST for a sold order come from the
	// same D-41 per-item matrix calculation the Simulador prices with
	// (pricingpostgres.MatrixReader + ProductFiscalReader), so the two
	// screens can never quote different tax on the same money. T5 (P2.b):
	// replaces the ProfileSource-based imposto/DIFAL calculation end to end.
	ordersTaxReader := orderspricingtax.NewReader(
		pricingpostgres.NewMatrixReader(pool),
		pricingpostgres.NewProductFiscalReader(pool, activeSourceLookup),
		cfg.DefaultTenantID,
	)
	ordersEnrichSvc := ordersapp.NewEnrichServiceWithReaders(ordersCostReader, ordersEnrichShipmentReader, nil, ordersEnrichBuyerFiscalReader, slog.Default()).
		WithLinkRefresh(ordersLinkReader).
		WithTaxes(ordersTaxReader)
	ordersFaturadoSvc := ordersapp.NewFaturadoService(ordersRepo)
	orderstransport.NewHandlerWithSummary(ordersImportSvc, ordersReadSvc, &ordersEnrichSvc, ordersSummarySvc).WithFaturador(ordersFaturadoSvc).Register(mux)

	linkageRepo := orderspostgres.NewSankhyaLinkageRepository(pool, cfg.DefaultTenantID)
	var assistedLinkageApp orderstransport.AssistedSankhyaLinkageApplication
	var assistedLinkageService *ordersapp.AssistedSankhyaLinkageService
	if oracleDB != nil {
		if runtimeConfig, err := internalreadoracle.LoadSankhyaLinkageRuntimeConfigFromEnv(os.Getenv); err != nil {
			slog.Warn("assisted Sankhya linkage unavailable", "err", err)
		} else if !internalReadAvailable || oracleDB == nil {
			slog.Warn("assisted Sankhya linkage unavailable", "err", "Oracle reader is unavailable")
		} else {
			source := internalreadobservability.NewTimingSankhyaLinkageReader(
				internalreadoracle.NewSankhyaLinkageReader(oracleDB, runtimeConfig.ReaderConfig),
				slog.Default(),
				observabilityCfg.SlowQueryThreshold,
			)
			if err := source.ValidateConfiguration(context.Background()); err != nil {
				_ = oracleDB.Close()
				return nil, fmt.Errorf("sankhya linkage startup validation: %w", err)
			}
			reader, err := ordersinternalread.NewSankhyaLinkageReader(
				source,
				runtimeConfig.ReaderConfig.ConfigurationRevision,
				runtimeConfig.EvidenceReference,
			)
			if err != nil {
				slog.Warn("assisted Sankhya linkage unavailable", "err", err)
			} else {
				assistedLinkageService = ordersapp.NewAssistedSankhyaLinkageService(ordersapp.AssistedSankhyaLinkageServiceConfig{
					Orders: ordersRepo, Reader: reader, Linkages: linkageRepo, Invalidator: freshnessCache,
				})
				assistedLinkageApp = assistedLinkageService
			}
		}
	}
	orderstransport.NewSankhyaLinkageHandler(assistedLinkageApp, cfg.DefaultTenantID).Register(mux)

	profitabilityStore := profitabilitypostgres.NewStore(pool, cfg.DefaultTenantID)
	profitabilityCfg := profitabilityapp.ServiceConfig{
		Orders:      profitabilityorders.NewOrderReader(ordersListSvc),
		Inputs:      profitabilityStore,
		Adjustments: profitabilityStore,
		Snapshots:   profitabilityStore,
	}
	if oracleDB != nil && internalReadAvailable {
		profitabilityCfg.Internal = profitabilityinternalread.NewFactReader(
			internalReadSvc,
			internalreadcache.NewBatchReader(
				internalreadoracle.NewBatchReader(oracleDB, oracleBatchSemaphore),
				freshnessCache,
			),
		)
	}
	profitabilityCfg.Invalidator = freshnessCache
	if assistedLinkageService != nil {
		profitabilityCfg.Lineage = profitabilityorders.NewSankhyaLineageReader(assistedLinkageService, cfg.DefaultTenantID)
	}
	profitabilitySvc := profitabilityapp.NewService(profitabilityCfg)
	profitabilitytransport.NewHandler(profitabilitySvc).Register(mux)

	go integrationsbg.NewRefreshTicker(authSessionRepo, authFlowSvc, 5*time.Minute).Start(context.Background())
	go integrationsbg.NewStateCleanup(oauthStateRepo, time.Hour).Start(context.Background())
	go integrationsbg.NewFeeSyncScheduler(installationSvc, providerSvc, feeSyncSvc, 15*time.Minute).Start(context.Background())
	// MIS-006: products sync. M-01 built the scheduler seam and M-03/M-04 built the
	// adapters, but nothing connected them, so the scheduler ticked a NO-OP job and
	// product data never refreshed on its own. The real job resolves the tenant's
	// active source each cycle and runs that source's adapter.
	if activeSourceLookup != nil {
		// Every sync that lands new product data re-asks the linking question
		// (M-05 F-01): a codprod that arrives today is what makes yesterday's
		// unlinked anúncio linkable.
		go synccomposition.NewProductsScheduler(
			pool, cfg.DefaultTenantID, 15*time.Minute, activeSourceLookup, productSyncAdapters,
			synccomposition.WithLinkCandidateRefresh(
				productlinkscomposition.NewLinkCandidateRefresher(installationSvc, productLinkGenerationSvc),
			),
		).Start(context.Background())
	}

	marketModuleRepo := marketpostgres.NewRepository(pool, cfg.DefaultTenantID)
	marketReadSvc := marketapp.NewReadService(marketModuleRepo, marketModuleRepo, time.Now)
	marketCostReader := newMarketCostReaderAdapter(internalReadSvc, internalReadAvailable)
	marketIdentityReader := newMarketProductIdentityReaderAdapter(internalReadSvc, internalReadAvailable, installationSvc, productLinkCandidateRepo)
	marketPriceIntelReader := newMarketPriceIntelCollectorAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
	marketCollectionSvc := marketapp.NewCollectionPipelineService(
		marketIdentityReader,
		marketPriceIntelReader,
		marketModuleRepo,
		marketModuleRepo,
		marketModuleRepo,
		marketModuleRepo,
		marketModuleRepo,
		time.Now,
	)
	marketEvidenceSvc := marketapp.NewEvidenceReadService(marketModuleRepo, marketModuleRepo, marketModuleRepo, marketCostReader, time.Now)
	markettransport.NewHandlerWithCollections(marketReadSvc, marketCollectionSvc, marketEvidenceSvc).Register(mux)

	marketRepo := marketplacespostgres.NewRepository(pool, cfg.DefaultTenantID)
	marketSvc := marketplacesapp.NewService(marketRepo, cfg.DefaultTenantID)

	listingRepo := listingspostgres.NewRepository(pool, cfg.DefaultTenantID)
	listingCostReader := listingsinternalread.NewCostReader(internalreadoracle.NewBatchReader(oracleDB, oracleBatchSemaphore))
	listingPolicyReader := listingsmarketplaces.NewPolicyReader(marketSvc)
	listingInstallationReader := listingsintegrations.NewInstallationReader(installationSvc)
	listingEvidenceReader := newListingsEvidenceAdapter(marketEvidenceSvc)
	listingSvc := listingsapp.NewReadServiceWithEvidence(
		listingRepo,
		listingCostReader,
		listingPolicyReader,
		listingInstallationReader,
		time.Now,
		listingEvidenceReader,
	)
	listingstransport.NewReadHandler(listingSvc).Register(mux)
	dashboardSvc := dashboardapp.NewService(
		installationSvc,
		listingSvc,
		productLinkSummarySvc,
		ordersSummarySvc,
		operationSvc,
		erpQuerySvc,
		time.Now,
	)
	dashboardtransport.NewHandler(dashboardSvc).Register(mux)
	// The listings sync is the ONLY path that walks the seller's whole catalog,
	// so it is also where the link matcher gets its snapshots: the observer
	// upserts the provider EAN/seller-SKU the canonical listing row cannot hold.
	// listingBackfillHydrator and listingRepo (below) are shared by BOTH
	// producers that need hydrate+persist — the manual API refresh's
	// whole-catalog BackfillRunner right below, and the admin resync writer's
	// per-item ListingIngestor (wired further down as installationResyncWriter)
	// — so a per-item resync and a whole-catalog backfill agree on hydrate
	// and persist behavior without a resync fanning out into N whole-catalog
	// pulls (ADR-04: one writer per shared seam; F-04 retires the old
	// page-based Source/Ingestion).
	listingBackfillHydrator := listingsconnectors.NewMultigetHydrator(marketplaceCapabilities, productLinkImportSvc, time.Now)
	listingBackfillRunner := listingsapp.NewBackfillRunner(
		listingsconnectors.NewBackfillSource(marketplaceCapabilities),
		listingBackfillHydrator,
		listingRepo, time.Now,
	)
	listingRefreshGateway := listingsintegrations.NewGateway(installationSvc, operationSvc)
	listingRefreshSvc := listingsapp.NewRefreshService(
		listingRefreshGateway, listingBackfillRunner, func(task func()) { go task() }, time.Now,
		func(err error) { slog.Error("listings refresh lifecycle persistence failed", "err", err) },
		func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) },
	)
	listingstransport.NewRefreshHandler(listingRefreshSvc).Register(mux)
	// Daily listings backfill/sweep scheduler (ADR-08 second Scheduler
	// instance, mirroring synccomposition.NewProductsScheduler): one instance
	// per active Mercado Livre installation, since (unlike products) a
	// listings backfill run is scoped to one specific InstallationAccount.
	listingscomposition.NewListingsSchedulers(
		pool, 24*time.Hour, installationSvc, marketplaceCapabilities, productLinkImportSvc, listingRepo,
	).StartAll(context.Background())

	mutationLane := mutationLane{envelope: mutationEnvelope}
	mutationWriter := mutationsports.WriterPort(mutationsstub.NewWriter(nil))
	if providerWritesEnabled(os.Getenv) {
		priceWriter := mutationsconnectors.NewPriceWriter(marketplaceCapabilities, cfg.DefaultTenantID, "mercado_livre")
		stockWriter := installationStockWriter{capabilities: marketplaceCapabilities, installations: installationSvc, tenantID: cfg.DefaultTenantID}
		listingWriter := mutationsconnectors.NewListingWriter(marketplaceCapabilities, cfg.DefaultTenantID, "mercado_livre")
		linkageWriter := mutationsproductlinks.NewWriter(productLinkResolutionSvc)
		resyncWriter := installationResyncWriter{hydrator: listingBackfillHydrator, store: listingRepo, installations: installationSvc, tenantID: cfg.DefaultTenantID}
		policyPresent := func(ctx context.Context, listingID string) (bool, error) {
			parsed, err := listingsdomain.ParseListingID(listingID)
			if err != nil {
				return false, err
			}
			_, found, err := marketSvc.GetPricingPolicyForInstallation(ctx, parsed.InstallationID)
			return found, err
		}
		mutationWriter = mutationsapp.NewWriterRouter(priceWriter, stockWriter, listingWriter, resyncWriter, linkageWriter, linkageWriter.HasResolvedLink, policyPresent)
		mutationLane.real, mutationLane.price, mutationLane.stock = true, true, true
		mutationLane.listing, mutationLane.linkage, mutationLane.resync = true, true, true
	}
	mutationLane.writer = mutationWriter
	mutationService := mutationsapp.NewService(mutationRepo, mutationslistings.NewSelectionResolver(listingSvc), listingSvc, time.Now)
	mutationApproval := mutationsapp.NewApprovalService(mutationRepo, time.Now)
	mutationRetry := mutationsapp.NewRetryService(mutationRepo, time.Now)
	mutationHandler := mutationstransport.NewHandler(mutationService, mutationApproval, mutationRetry)
	mutationHandler.Register(mux)
	mux.HandleFunc("POST /mutations/{id}/retry-failures", mutationHandler.HandleRetry)
	mutationstransport.NewQueryHandler(mutationsapp.NewReadService(mutationRepo)).Register(mux)
	mutationPoller := mutationsapp.NewPoller(mutationRepo, mutationWriter, time.Now)
	mutationRunner := mutationsbg.NewPoller(mutationPoller, func(ctx context.Context) ([]string, error) {
		installations, err := installationSvc.List(ctx)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(installations))
		for i := range installations {
			ids[i] = installations[i].InstallationID
		}
		return ids, nil
	}, slog.Default(), 2*time.Second, nil)

	feeSvc := marketplacesapp.NewFeeScheduleService(feeRepo)

	if pool != nil {
		if err := feeSvc.SeedDefinitions(context.Background()); err != nil {
			slog.Warn("marketplace definitions sync failed", "err", err)
		}
	}

	connectorsFeeSyncSvc := connectorsapp.NewFeeSyncService(feeRepo, registeredFeeSyncers...)

	if pool != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			connectorsFeeSyncSvc.SeedAll(ctx)
		}()
	}

	// Seed stub fee rows for channels without a dedicated FeeSyncer (Amazon, Leroy, Madeira).
	if pool != nil {
		go func() {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			slog.Info("registry.SeedAll started", "action", "seed_stub_fees")
			marketplacesregistry.SeedAll(ctx, pool)
			slog.Info("registry.SeedAll completed", "action", "seed_stub_fees", "result", "ok", "duration_ms", time.Since(start).Milliseconds())
		}()
	}

	marketplacestransport.NewHandler(marketSvc, feeSvc, connectorsFeeSyncSvc).Register(mux)

	pricingRepo := pricingpostgres.NewRepository(pool, cfg.DefaultTenantID)
	pricingSvc := pricingapp.NewService(pricingRepo, cfg.DefaultTenantID)

	// Melhor Envio
	meTokenStore := melhorenvio.NewTokenStore(pool, cfg.DefaultTenantID)
	meClient := melhorenvio.NewClient(meTokenStore)
	meOAuth := melhorenvio.NewOAuthHandlerFromEnv(meTokenStore) // nil if ME_CLIENT_ID unset

	// Pricing batch orchestrator
	prodReader := pricingcatalog.NewReader(catalogSvc)
	polReader := pricingmarket.NewReader(marketSvc)
	batchOrch := pricingapp.NewBatchOrchestrator(prodReader, polReader, meClient, cfg.DefaultTenantID)

	// IC-04 calculator (M-07): profile, DIFAL seed+override, scenarios, decompose,
	// solve. Wired only when internal_read (cost facts) is available — mirrors the
	// profitability module's conditional Internal wiring. When unavailable the calc
	// routes are simply absent; W1 /pricing/simulations is unaffected.
	pricingHandler := pricingtransport.NewHandler(pricingSvc, batchOrch)
	if internalReadAvailable {
		calcRepo := pricingpostgres.NewCalcRepository(pool)
		calcCost := pricingcostread.NewReader(profitabilityinternalread.NewFactReader(internalReadSvc))
		calcProducts := pricingcatalog.NewProductChecker(catalogSvc)
		// CHIP-T1: tariff defaults store + degrau-4 resolver (installation "" =
		// default-installation sentinel for the single-installation demo). calcRepo
		// implements ports.TariffDefaultsStore; the resolver reads config-seeded
		// comissão/frete so an empty comissao_pct resolves and frete honors policy.
		calcTariffResolver := pricingtariffdefaults.NewResolver(calcRepo, cfg.DefaultTenantID, "")
		// CHIP-T2-MIN degrau 3 (live COTAÇÃO): compose over the degrau-4 base
		// when the mercado_livre fee-quote reader is available, so a decompose/
		// solve upgrades the commission to a live quote (category via catalog
		// match, commission via listing_prices) and silently falls to degrau 4
		// on any miss. If the provider has no fee-quote capability wired, the
		// base degrau-4 resolver stands alone.
		var tariffResolver pricingports.TariffResolver = calcTariffResolver
		if feeReader, ferr := marketplaceCapabilities.FeeQuoteReader(mercadoLivreProviderCode); ferr == nil {
			identityReader := pricingcatalog.NewProductIdentityReader(catalogSvc)
			categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
			commissionQuoter := newPricingCommissionQuoterAdapter(feeReader, installationSvc, cfg.DefaultTenantID)
			liveResolver := pricingtarifflive.NewResolver(identityReader, categoryResolver, commissionQuoter)
			tariffResolver = pricingtariffcomposite.NewResolver(calcTariffResolver, liveResolver)
		}
		batchOrch.WithTariffResolver(tariffResolver)
		calcSvc := pricingapp.NewCalcService(calcRepo, calcCost, calcProducts, cfg.DefaultTenantID).
			WithTariffStore(calcRepo).
			WithTariffResolver(tariffResolver)
		pricingHandler = pricingHandler.WithCalc(calcSvc)
	}
	pricingHandler.Register(mux)

	// Connectors (Melhor Envio auth + fee seeding foundations)
	connectorstransport.NewHandler(meOAuth).Register(mux)

	// M-09 F-01: GET /sync/health — per-entity sync_state aggregate (GREATEST
	// of last_full_sync_at/last_incremental_at, tolerant cursor.phase read) +
	// a webhook-stats block behind the WebhookStatsReader seam. The webhook
	// reader stays the M-09 default (canonical zero) until a future milestone
	// calls syncHealthSvc.WithWebhookStatsReader from its own region — that
	// call mutates syncHealthSvc in place, so it takes effect on this same
	// already-registered route without touching this region.
	syncHealthReader := syncpg.NewHealthReader(pool)
	syncHealthSvc := syncapp.NewHealthService(syncHealthReader, cfg.DefaultTenantID)
	synctransport.NewHealthHandler(syncHealthSvc).Register(mux)

	return &RootRuntime{Handler: httpx.CORSMiddleware(apierror.Recover(mux)), CatalogReader: catalogPageReader, PoolStats: poolStats, MutationPoller: mutationRunner, mutationLane: mutationLane}, nil
}

func providerWritesEnabled(getenv func(string) string) bool {
	return strings.EqualFold(strings.TrimSpace(getenv("MPC_PROVIDER_WRITES_ENABLED")), "true")
}

// mlCatalogOffersEnabled gates the mercado_livre catalog-offers READ route
// (ADR-05: flag defaults OFF). This is a read-only provider flag, unrelated
// to MPC_PROVIDER_WRITES_ENABLED and the mutations dispatcher.
func mlCatalogOffersEnabled(getenv func(string) string) bool {
	return strings.EqualFold(strings.TrimSpace(getenv("MPC_ML_CATALOG_OFFERS_ENABLED")), "true")
}

func oracleBatchPermits(getenv func(string) string) int {
	const defaultPermits = 4
	raw := strings.TrimSpace(getenv("MPC_ORACLE_BATCH_PERMITS"))
	if raw == "" {
		return defaultPermits
	}
	permits, err := strconv.Atoi(raw)
	if err != nil || permits <= 0 {
		slog.Warn("invalid Oracle batch permits; using default", "value", raw, "default", defaultPermits)
		return defaultPermits
	}
	return permits
}
