package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	catalogapp "marketplace-central/apps/server_core/internal/modules/catalog/application"
	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	catalogtransport "marketplace-central/apps/server_core/internal/modules/catalog/transport"
	classapp "marketplace-central/apps/server_core/internal/modules/classifications/application"
	classdomain "marketplace-central/apps/server_core/internal/modules/classifications/domain"
	classtransport "marketplace-central/apps/server_core/internal/modules/classifications/transport"
	connectorstransport "marketplace-central/apps/server_core/internal/modules/connectors/transport"
	marketplacesapp "marketplace-central/apps/server_core/internal/modules/marketplaces/application"
	marketplacesdomain "marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
	marketplacestransport "marketplace-central/apps/server_core/internal/modules/marketplaces/transport"
	pricingapp "marketplace-central/apps/server_core/internal/modules/pricing/application"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingtransport "marketplace-central/apps/server_core/internal/modules/pricing/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// stubCatalogReader satisfies catalog ports.ProductReader with in-memory no-ops.
type stubCatalogReader struct{}

func (r stubCatalogReader) ListCanonicalProducts(_ context.Context) ([]catalogdomain.CanonicalProduct, error) {
	return []catalogdomain.CanonicalProduct{}, nil
}
func (r stubCatalogReader) GetCanonicalProduct(_ context.Context, _ catalogdomain.InternalProductID) (catalogdomain.CanonicalProduct, error) {
	return catalogdomain.CanonicalProduct{}, nil
}
func (r stubCatalogReader) SearchCanonicalProducts(_ context.Context, _ string) ([]catalogdomain.CanonicalProduct, error) {
	return []catalogdomain.CanonicalProduct{}, nil
}

func (r stubCatalogReader) ListProducts(_ context.Context) ([]catalogdomain.Product, error) {
	return nil, nil
}

func (r stubCatalogReader) GetProduct(_ context.Context, _ string) (catalogdomain.Product, error) {
	return catalogdomain.Product{}, nil
}

func (r stubCatalogReader) SearchProducts(_ context.Context, _ string) ([]catalogdomain.Product, error) {
	return nil, nil
}

func (r stubCatalogReader) ListTaxonomyNodes(_ context.Context) ([]catalogdomain.TaxonomyNode, error) {
	return nil, nil
}

func (r stubCatalogReader) ListProductsByIDs(_ context.Context, productIDs []string) ([]catalogdomain.Product, error) {
	result := make([]catalogdomain.Product, 0, len(productIDs))
	seen := make(map[string]struct{}, len(productIDs))
	for _, id := range productIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
	}
	return result, nil
}

// stubCatalogEnrichments satisfies catalog ports.EnrichmentStore with in-memory no-ops.
type stubCatalogEnrichments struct{}

func (r stubCatalogEnrichments) GetEnrichment(_ context.Context, productID string) (catalogdomain.ProductEnrichment, error) {
	return catalogdomain.ProductEnrichment{ProductID: productID}, nil
}

func (r stubCatalogEnrichments) UpsertEnrichment(_ context.Context, _ catalogdomain.ProductEnrichment) error {
	return nil
}

func (r stubCatalogEnrichments) ListEnrichments(_ context.Context, _ []string) (map[string]catalogdomain.ProductEnrichment, error) {
	return make(map[string]catalogdomain.ProductEnrichment), nil
}

// stubMarketplacesRepo satisfies marketplaces ports.Repository with in-memory no-ops.
type stubMarketplacesRepo struct{}

func (r stubMarketplacesRepo) SaveAccount(_ context.Context, _ marketplacesdomain.Account) error {
	return nil
}
func (r stubMarketplacesRepo) SavePolicy(_ context.Context, _ marketplacesdomain.Policy) error {
	return nil
}
func (r stubMarketplacesRepo) ListAccounts(_ context.Context) ([]marketplacesdomain.Account, error) {
	return nil, nil
}
func (r stubMarketplacesRepo) ListPolicies(_ context.Context) ([]marketplacesdomain.Policy, error) {
	return nil, nil
}

func (r stubMarketplacesRepo) ListPoliciesByIDs(_ context.Context, policyIDs []string) ([]marketplacesdomain.Policy, error) {
	result := make([]marketplacesdomain.Policy, 0, len(policyIDs))
	seen := make(map[string]struct{}, len(policyIDs))
	for _, id := range policyIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
	}
	return result, nil
}

// stubPricingRepo satisfies pricing ports.Repository with in-memory no-ops.
type stubPricingRepo struct{}

func (r stubPricingRepo) SaveSimulation(_ context.Context, _ pricingdomain.Simulation) error {
	return nil
}
func (r stubPricingRepo) ListSimulations(_ context.Context) ([]pricingdomain.Simulation, error) {
	return nil, nil
}

// stubClassificationsRepo satisfies classifications ports.Repository with in-memory no-ops.
type stubClassificationsRepo struct{}

func (r stubClassificationsRepo) List(_ context.Context) ([]classdomain.Classification, error) {
	return nil, nil
}
func (r stubClassificationsRepo) GetByID(_ context.Context, _ string) (classdomain.Classification, error) {
	return classdomain.Classification{}, nil
}
func (r stubClassificationsRepo) Create(_ context.Context, _ classdomain.Classification) error {
	return nil
}
func (r stubClassificationsRepo) Update(_ context.Context, _ classdomain.Classification) error {
	return nil
}
func (r stubClassificationsRepo) Delete(_ context.Context, _ string) error {
	return nil
}

// TestRouterRegistersAllFoundationEndpoints verifies that every expected route
// is registered and returns a non-404 response. It builds a minimal mux with
// stub repository adapters so that no real database connection is required.
func TestRouterRegistersAllFoundationEndpoints(t *testing.T) {
	t.Setenv("ME_CLIENT_ID", "test-client")
	t.Setenv("ME_CLIENT_SECRET", "test-secret")

	mux := http.NewServeMux()

	// /healthz
	base := httpx.NewRouter()
	mux.Handle("/healthz", base)

	// /catalog/products
	catalogSvc := catalogapp.NewService(stubCatalogReader{}, stubCatalogEnrichments{}, "tenant_default")
	catalogtransport.Handler{Service: catalogapp.NewCanonicalService(stubCatalogReader{}), CompatibilityService: catalogSvc}.Register(mux)

	// /classifications
	classSvc := classapp.NewService(stubClassificationsRepo{}, "tenant_default")
	classtransport.NewHandler(classSvc).Register(mux)

	// /marketplaces/accounts, /marketplaces/policies
	marketSvc := marketplacesapp.NewService(stubMarketplacesRepo{}, "tenant_default")
	marketplacestransport.NewHandler(marketSvc, nil, nil).Register(mux)

	// /pricing/simulations
	pricingSvc := pricingapp.NewService(stubPricingRepo{}, "tenant_default")
	pricingtransport.NewHandler(pricingSvc, nil).Register(mux)

	connectorstransport.NewHandler(nil).Register(mux)

	cases := []string{
		"/healthz",
		"/catalog/products",
		"/classifications",
		"/marketplaces/accounts",
		"/marketplaces/policies",
		"/pricing/simulations",
		"/pricing/simulations/batch",
		"/connectors/melhor-envio/auth/start",
		"/connectors/melhor-envio/auth/callback",
		"/connectors/melhor-envio/status",
	}

	for _, path := range cases {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Fatalf("expected route %s to be registered (got 404)", path)
		}
	}
}
