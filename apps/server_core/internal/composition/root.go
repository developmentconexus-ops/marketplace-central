package composition

import (
	"context"
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
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu"
	melhorenvio "marketplace-central/apps/server_core/internal/modules/connectors/adapters/melhorenvio"
	mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee"
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorstransport "marketplace-central/apps/server_core/internal/modules/connectors/transport"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/amazon"
	integrationscrypto "marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto"
	integrationsfeesync "marketplace-central/apps/server_core/internal/modules/integrations/adapters/feesync"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/leroymerlin"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/madeiramadeira"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsbg "marketplace-central/apps/server_core/internal/modules/integrations/background"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	integrationstransport "marketplace-central/apps/server_core/internal/modules/integrations/transport"
	internalreadoracle "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
	internalreadcache "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache"
	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadobservability "marketplace-central/apps/server_core/internal/modules/internal_read/observability"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	inventoryconnectors "marketplace-central/apps/server_core/internal/modules/inventory/adapters/connectors"
	inventoryintegrations "marketplace-central/apps/server_core/internal/modules/inventory/adapters/integrations"
	inventoryinternalread "marketplace-central/apps/server_core/internal/modules/inventory/adapters/internalread"
	inventorypostgres "marketplace-central/apps/server_core/internal/modules/inventory/adapters/postgres"
	inventoryproductlinks "marketplace-central/apps/server_core/internal/modules/inventory/adapters/productlinks"
	inventoryapp "marketplace-central/apps/server_core/internal/modules/inventory/application"
	inventoryports "marketplace-central/apps/server_core/internal/modules/inventory/ports"
	inventorytransport "marketplace-central/apps/server_core/internal/modules/inventory/transport"
	marketplacespostgres "marketplace-central/apps/server_core/internal/modules/marketplaces/adapters/postgres"
	marketplacesapp "marketplace-central/apps/server_core/internal/modules/marketplaces/application"
	marketplacesregistry "marketplace-central/apps/server_core/internal/modules/marketplaces/registry"
	marketplacestransport "marketplace-central/apps/server_core/internal/modules/marketplaces/transport"
	ordersintegrations "marketplace-central/apps/server_core/internal/modules/orders/adapters/integrations"
	ordersinternalread "marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread"
	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	ordersproductlinks "marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	orderstransport "marketplace-central/apps/server_core/internal/modules/orders/transport"
	pricingcatalog "marketplace-central/apps/server_core/internal/modules/pricing/adapters/catalog"
	pricingfee "marketplace-central/apps/server_core/internal/modules/pricing/adapters/feeschedule"
	pricingmarket "marketplace-central/apps/server_core/internal/modules/pricing/adapters/marketplace"
	pricingpostgres "marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres"
	pricingapp "marketplace-central/apps/server_core/internal/modules/pricing/application"
	pricingtransport "marketplace-central/apps/server_core/internal/modules/pricing/transport"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinkstransport "marketplace-central/apps/server_core/internal/modules/product_links/transport"
	profitabilityinternalread "marketplace-central/apps/server_core/internal/modules/profitability/adapters/internalread"
	profitabilityorders "marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders"
	profitabilitypostgres "marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres"
	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitytransport "marketplace-central/apps/server_core/internal/modules/profitability/transport"
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

type RootRuntime struct {
	Handler   http.Handler
	PoolStats *internalreadobservability.PoolStatsLoop
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
	credentialRepo := integrationspostgres.NewCredentialRepository(pool, cfg.DefaultTenantID)
	credentialSvc := integrationsapp.NewCredentialService(credentialRepo, cfg.DefaultTenantID)
	authSessionRepo := integrationspostgres.NewAuthSessionRepository(pool, cfg.DefaultTenantID)
	authSvc := integrationsapp.NewAuthService(authSessionRepo, cfg.DefaultTenantID)
	capabilityStateRepo := integrationspostgres.NewCapabilityStateRepository(pool, cfg.DefaultTenantID)
	capabilitySvc := integrationsapp.NewCapabilityService(capabilityStateRepo, cfg.DefaultTenantID)
	operationRunRepo := integrationspostgres.NewOperationRunRepository(pool, cfg.DefaultTenantID)
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
	})
	marketplaceCapabilities := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{
		mercadoLivreCapabilities.ProviderCapabilitySet(),
	})
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
	})

	flowFacade := authFlowFacade{
		flow:        authFlowSvc,
		feeSync:     feeSyncSvc,
		operations:  operationSvc,
		providerOps: providerOperationSvc,
	}

	integrationstransport.NewHandler(providerSvc, installationSvc).Register(mux)
	integrationstransport.NewAuthHandler(flowFacade).Register(mux)

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
	oracleBatchSemaphore := oraclebatch.NewSemaphore(oracleBatchPermits(os.Getenv))
	var internalReadSvc internalreadapp.Service
	var internalReadAvailable bool
	var oracleDB internalreadoracle.Database
	var poolStats *internalreadobservability.PoolStatsLoop
	if oracleCfg, err := internalreadoracle.LoadConfigFromEnv(os.Getenv); err != nil {
		slog.Warn("product links oracle reader unavailable", "err", err)
	} else if db, err := internalreadoracle.OpenDB(context.Background(), oracleCfg); err != nil {
		slog.Warn("product links oracle connection failed", "err", err)
	} else {
		oracleDB = db
		cachedReader := internalreadcache.NewReader(
			internalreadoracle.NewReader(oracleDB),
			freshnessCache,
		)
		reader := internalreadobservability.NewTimingReader(
			cachedReader,
			slog.Default(),
			observabilityCfg.SlowQueryThreshold,
		)
		internalReadSvc = internalreadapp.NewService(reader)
		poolStats = internalreadobservability.NewPoolStatsLoop(oracleDB, slog.Default(), observabilityCfg.PoolStatsInterval)
		internalReadAvailable = true
		productMatcher = internalReadSvc
		inventoryStockReader = internalreadcache.NewStockBatchReader(
			internalreadoracle.NewBatchStockReader(oracleDB, oracleBatchSemaphore),
			freshnessCache,
		)
	}
	var canonicalCatalogReader catalogports.CanonicalProductReader = cataloginternalread.UnavailableReader{Err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle catalog reader is unavailable", nil)}
	if internalReadAvailable {
		canonicalCatalogReader = cataloginternalread.NewReader(internalReadSvc)
	}
	catalogEnrichments := catalogpostgres.NewEnrichmentRepository(pool, cfg.DefaultTenantID)
	legacyReader := canonicalCatalogReader.(catalogports.ProductReader)
	catalogSvc := catalogapp.NewServiceWithInvalidator(legacyReader, catalogEnrichments, cfg.DefaultTenantID, freshnessCache)
	var catalogPageReader internalreadports.CatalogPageReader
	if internalReadAvailable {
		catalogPageReader = internalReadSvc
	}
	catalogtransport.Handler{Service: catalogapp.NewCanonicalService(canonicalCatalogReader), CompatibilityService: catalogSvc, PageReader: catalogPageReader}.Register(mux)

	productLinkGenerationSvc := productlinksapp.NewGenerationService(productlinksapp.GenerationServiceConfig{
		Snapshots: productLinkSnapshotRepo,
		Matcher:   productMatcher,
		Store:     productLinkCandidateRepo,
	})
	productLinkResolutionSvc := productlinksapp.NewResolutionService(productlinksapp.ResolutionServiceConfig{
		Candidates: productLinkCandidateRepo,
		Workflows:  productLinkCandidateRepo,
	})
	productlinkstransport.NewHandler(productLinkImportSvc, productLinkGenerationSvc, productLinkGenerationSvc, productLinkResolutionSvc, productLinkResolutionSvc).Register(mux)

	inventoryActionRepo := inventorypostgres.NewStockActionRepository(pool, cfg.DefaultTenantID)
	inventoryRiskSvc := inventoryapp.NewStockRiskService(
		inventoryproductlinks.NewSnapshotReader(productLinkSnapshotRepo),
		inventoryproductlinks.NewLinkReader(productLinkCandidateRepo, productLinkCandidateRepo),
		inventoryStockReader,
	)
	inventoryActionSvc := inventoryapp.NewStockActionServiceWithInvalidator(
		inventoryActionRepo,
		inventoryconnectors.NewStockWriter(marketplaceCapabilities),
		freshnessCache,
	)
	inventoryManualAction := inventoryapp.NewManualActionFacade(
		inventoryRiskSvc,
		inventoryActionSvc,
		inventoryintegrations.NewInstallationReader(installationSvc),
		inventoryproductlinks.NewSnapshotStore(productLinkSnapshotRepo),
	)
	inventorytransport.NewHandler(inventoryRiskSvc, inventoryManualAction).Register(mux)

	ordersRepo := orderspostgres.NewOrderRepository(pool, cfg.DefaultTenantID)
	ordersImportSvc := ordersapp.NewImportService(ordersapp.ImportServiceConfig{
		Source: ordersintegrations.NewOrderSource(providerOperationSvc),
		Links:  ordersproductlinks.NewLinkReader(productLinkCandidateRepo, productLinkCandidateRepo),
		Store:  ordersRepo,
	})
	ordersListSvc := ordersapp.NewListService(ordersRepo)
	orderstransport.NewHandler(ordersImportSvc, ordersListSvc).Register(mux)

	linkageRepo := orderspostgres.NewSankhyaLinkageRepository(pool, cfg.DefaultTenantID)
	var assistedLinkageApp orderstransport.AssistedSankhyaLinkageApplication
	var assistedLinkageService *ordersapp.AssistedSankhyaLinkageService
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
	orderstransport.NewSankhyaLinkageHandler(assistedLinkageApp, cfg.DefaultTenantID).Register(mux)

	profitabilityStore := profitabilitypostgres.NewStore(pool, cfg.DefaultTenantID)
	profitabilityCfg := profitabilityapp.ServiceConfig{
		Orders:      profitabilityorders.NewOrderReader(ordersListSvc),
		Inputs:      profitabilityStore,
		Adjustments: profitabilityStore,
		Snapshots:   profitabilityStore,
	}
	if internalReadAvailable {
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

	marketRepo := marketplacespostgres.NewRepository(pool, cfg.DefaultTenantID)
	marketSvc := marketplacesapp.NewService(marketRepo, cfg.DefaultTenantID)

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
	feeAdapter := pricingfee.NewAdapter(feeSvc)
	prodReader := pricingcatalog.NewReader(catalogSvc)
	polReader := pricingmarket.NewReader(marketSvc)
	batchOrch := pricingapp.NewBatchOrchestrator(prodReader, polReader, meClient, feeAdapter, cfg.DefaultTenantID)
	pricingtransport.NewHandler(pricingSvc, batchOrch).Register(mux)

	// Connectors (Melhor Envio auth + fee seeding foundations)
	connectorstransport.NewHandler(meOAuth).Register(mux)

	return &RootRuntime{Handler: httpx.CORSMiddleware(mux), PoolStats: poolStats}, nil
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
