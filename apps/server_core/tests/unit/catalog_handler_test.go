package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	catalogapp "marketplace-central/apps/server_core/internal/modules/catalog/application"
	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	catalogtransport "marketplace-central/apps/server_core/internal/modules/catalog/transport"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type catalogHandlerReaderStub struct {
	products []catalogdomain.Product
}

func (r *catalogHandlerReaderStub) GetProduct(_ context.Context, id string) (catalogdomain.Product, error) {
	for _, p := range r.products {
		if p.ProductID == id {
			return p, nil
		}
	}
	return catalogdomain.Product{}, nil
}

func (r *catalogHandlerReaderStub) ListTaxonomyNodes(_ context.Context) ([]catalogdomain.TaxonomyNode, error) {
	return nil, nil
}

func (r *catalogHandlerReaderStub) ListProductsByIDs(_ context.Context, productIDs []string) ([]catalogdomain.Product, error) {
	result := make([]catalogdomain.Product, 0, len(productIDs))
	seen := make(map[string]struct{}, len(productIDs))
	for _, id := range productIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		for _, p := range r.products {
			if p.ProductID == id {
				result = append(result, p)
				break
			}
		}
	}
	return result, nil
}

func (r *catalogHandlerReaderStub) GetCanonicalProduct(_ context.Context, id catalogdomain.InternalProductID) (catalogdomain.CanonicalProduct, error) {
	for _, p := range canonicalCatalogProducts(r.products) {
		if p.InternalProductID == id {
			return p, nil
		}
	}
	return catalogdomain.CanonicalProduct{}, nil
}
func canonicalCatalogProducts(products []catalogdomain.Product) []catalogdomain.CanonicalProduct {
	zero := 0.0
	now := time.Now().UTC()
	fact, _ := catalogdomain.NewNumericSourceFact("fake", &zero, catalogdomain.SourceFactQualityCurrent, &now, nil)
	result := make([]catalogdomain.CanonicalProduct, 0, len(products))
	for i, p := range products {
		result = append(result, catalogdomain.CanonicalProduct{InternalProductID: catalogdomain.InternalProductID(i + 1), Name: p.Name, CostAmount: fact, PriceAmount: fact, StockQuantity: fact})
	}
	return result
}

type catalogHandlerEnrichmentStub struct{}

func (s catalogHandlerEnrichmentStub) GetEnrichment(_ context.Context, productID string) (catalogdomain.ProductEnrichment, error) {
	return catalogdomain.ProductEnrichment{ProductID: productID}, nil
}

func (s catalogHandlerEnrichmentStub) UpsertEnrichment(_ context.Context, _ catalogdomain.ProductEnrichment) error {
	return nil
}

func (s catalogHandlerEnrichmentStub) ListEnrichments(_ context.Context, _ []string) (map[string]catalogdomain.ProductEnrichment, error) {
	return make(map[string]catalogdomain.ProductEnrichment), nil
}

// catalogHandlerPageReader stands for everything below transport: the routing
// seam that answers a nil policy with the tenant's stored rule, and the source
// under it. The handler asks with nil on every read that is not include_all, so
// a page reader that refused nil would be testing the wrong seat — resolving it
// is what the layer below transport is for. The cut itself is proved end-to-end,
// through the real chain, in cache_composed_test.go.
type catalogHandlerPageReader struct {
	products []catalogdomain.Product
}

func (r catalogHandlerPageReader) facts(match string) readports.CatalogFactPage {
	items := make([]readports.CatalogProductFact, 0, len(r.products))
	for i, p := range r.products {
		if match != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(match)) {
			continue
		}
		name := p.Name
		items = append(items, readports.CatalogProductFact{
			InternalProductID: int64(i + 1),
			Description:       &name,
			Active:            true,
		})
	}
	return readports.CatalogFactPage{Items: items, AsOf: time.Now().UTC()}
}

func (r catalogHandlerPageReader) ListCatalogProductFacts(_ context.Context, _ readports.Cursor, _ int, _ *readports.SellableAssortmentPolicy) (readports.CatalogFactPage, error) {
	return r.facts(""), nil
}

func (r catalogHandlerPageReader) SearchCatalogProductFacts(_ context.Context, q string, _ readports.Cursor, _ int, _ *readports.SellableAssortmentPolicy) (readports.CatalogFactPage, error) {
	return r.facts(q), nil
}

func (r catalogHandlerPageReader) CatalogProductFactsByIDs(context.Context, []int64) (readports.CatalogFactPage, error) {
	return readports.CatalogFactPage{AsOf: time.Now().UTC()}, nil
}

func (r catalogHandlerPageReader) GetCatalogAssortmentCounts(context.Context, *readports.SellableAssortmentPolicy) (readports.CatalogAssortmentCounts, error) {
	return readports.CatalogAssortmentCounts{SellableCount: len(r.products), TotalCount: len(r.products)}, nil
}

func newCatalogHandler(products []catalogdomain.Product) catalogtransport.Handler {
	reader := &catalogHandlerReaderStub{products: products}
	svc := catalogapp.NewService(reader, catalogHandlerEnrichmentStub{}, "tnt_test")
	return catalogtransport.Handler{
		Service:              catalogapp.NewCanonicalService(reader),
		CompatibilityService: svc,
		PageReader:           catalogHandlerPageReader{products: products},
	}
}

func TestCatalogHandlerGetReturnsProducts(t *testing.T) {
	seeded := []catalogdomain.Product{
		{ProductID: "p-1", SKU: "SKU-1", Name: "Widget", Status: "active", CostAmount: 10.0},
	}
	mux := http.NewServeMux()
	newCatalogHandler(seeded).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/catalog/products", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	items, ok := result["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("expected non-empty items, got %v", result["items"])
	}
}

func TestCatalogHandlerRejectsPost(t *testing.T) {
	mux := http.NewServeMux()
	newCatalogHandler(nil).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/catalog/products", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// The route is registered as "GET /catalog/products", so net/http rejects the
	// verb before the handler runs and writes its own plain-text body. That is
	// what production has always answered here: the JSON 405 this used to assert
	// came from the legacy handler, which production never registered.
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	// net/http serves HEAD automatically for a "GET /..." pattern, so the
	// router's Allow header lists both verbs it actually answers.
	if allow := rec.Header().Get("Allow"); allow != "GET, HEAD" {
		t.Fatalf("expected Allow: GET, HEAD on the 405, got %q", allow)
	}
}

func TestCatalogHandlerRejectsPut(t *testing.T) {
	mux := http.NewServeMux()
	newCatalogHandler(nil).Register(mux)

	req := httptest.NewRequest(http.MethodPut, "/catalog/products", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestCatalogSearchEndpoint(t *testing.T) {
	seeded := []catalogdomain.Product{
		{ProductID: "prd_1", SKU: "SKU-001", Name: "Cuba Inox"},
		{ProductID: "prd_2", SKU: "SKU-002", Name: "Torneira Gourmet"},
	}
	mux := http.NewServeMux()
	newCatalogHandler(seeded).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/catalog/products/search?q=cuba", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	items, ok := result["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected 1 search result, got %v", result["items"])
	}
}

func TestCatalogSearchEndpointRequiresQuery(t *testing.T) {
	mux := http.NewServeMux()
	newCatalogHandler(nil).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/catalog/products/search", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// q is required by the OpenAPI contract. The guard used to live only in the
	// deleted legacy handler, so the paged route accepted a blank q and answered
	// with the whole catalog dressed as search results.
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode 400 body: %v", err)
	}
	if body["error"] != "invalid_q" {
		t.Fatalf("expected error=invalid_q, got %v", body)
	}
}

func TestCatalogGetProductEndpoint(t *testing.T) {
	seeded := []catalogdomain.Product{
		{ProductID: "prd_1", SKU: "SKU-001", Name: "Cuba Inox"},
	}
	mux := http.NewServeMux()
	newCatalogHandler(seeded).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/catalog/products/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if result["internal_product_id"] != float64(1) {
		t.Fatalf("expected canonical product id 1, got %v", result["internal_product_id"])
	}
}

func TestCatalogUpsertEnrichmentEndpoint(t *testing.T) {
	seeded := []catalogdomain.Product{
		{ProductID: "prd_1", SKU: "SKU-001", Name: "Cuba Inox"},
	}
	mux := http.NewServeMux()
	newCatalogHandler(seeded).Register(mux)

	body := `{"height_cm": 30.0, "suggested_price_amount": 199.99}`
	req := httptest.NewRequest(http.MethodPut, "/catalog/products/prd_1/enrichment", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}
