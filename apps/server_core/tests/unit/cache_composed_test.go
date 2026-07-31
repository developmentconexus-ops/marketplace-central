package unit

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	catalogtransport "marketplace-central/apps/server_core/internal/modules/catalog/transport"
	cacheadapter "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache"
	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/routing"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadobservability "marketplace-central/apps/server_core/internal/modules/internal_read/observability"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type composedFakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *composedFakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *composedFakeClock) Advance(delta time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(delta)
	c.mu.Unlock()
}

// composedCatalogOracleReader stands in for the Oracle source at the bottom of
// the live chain. It is a whole internalreadports.Source because that is what
// the chain takes; the per-product reads panic rather than answer, since no
// composed test drives them and a silent zero would be a lie.
//
// The assortment cut is SQL-side in the real adapter, so this stub cannot
// evaluate the predicate — it models the OUTCOME: product 101 survives every
// cut, product 202 only survives the uncut read. It refuses an unresolved
// policy for the same reason the adapter does, so a seat that drops the policy
// on its way down fails here instead of quietly paging the whole catalog.
type composedCatalogOracleReader struct {
	clock       *composedFakeClock
	listCalls   atomic.Int64
	searchCalls atomic.Int64
}

func (r *composedCatalogOracleReader) page(policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	resolved, err := internalreadports.RequireAssortmentPolicy(policy)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	items := []internalreadports.CatalogProductFact{{InternalProductID: 101, Active: true}}
	if !resolved.OnlyRevenda && !resolved.OnlyEmEstoque && !resolved.OnlyEcommerceEligible {
		items = append(items, internalreadports.CatalogProductFact{InternalProductID: 202, Active: true})
	}
	return internalreadports.CatalogFactPage{Items: items, AsOf: r.clock.Now()}, nil
}

func (r *composedCatalogOracleReader) ListCatalogProductFacts(_ context.Context, _ internalreadports.Cursor, _ int, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	r.listCalls.Add(1)
	return r.page(policy)
}

func (r *composedCatalogOracleReader) CatalogProductFactsByIDs(context.Context, []int64) (internalreadports.CatalogFactPage, error) {
	return internalreadports.CatalogFactPage{AsOf: r.clock.Now()}, nil
}

func (r *composedCatalogOracleReader) SearchCatalogProductFacts(_ context.Context, _ string, _ internalreadports.Cursor, _ int, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	r.searchCalls.Add(1)
	page, err := r.page(policy)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	page.Items = nil
	return page, nil
}

func (r *composedCatalogOracleReader) GetCatalogAssortmentCounts(_ context.Context, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogAssortmentCounts, error) {
	page, err := r.page(policy)
	if err != nil {
		return internalreadports.CatalogAssortmentCounts{}, err
	}
	return internalreadports.CatalogAssortmentCounts{SellableCount: len(page.Items), TotalCount: 2}, nil
}

func (r *composedCatalogOracleReader) FindProductsForLinking(context.Context, internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	panic("composed catalog test does not drive the per-product reads")
}

func (r *composedCatalogOracleReader) GetSellableStock(context.Context, internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	panic("composed catalog test does not drive the per-product reads")
}

func (r *composedCatalogOracleReader) GetCurrentPrice(context.Context, internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	panic("composed catalog test does not drive the per-product reads")
}

func (r *composedCatalogOracleReader) GetCostAsOf(context.Context, internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	panic("composed catalog test does not drive the per-product reads")
}

func (r *composedCatalogOracleReader) GetSalesHistory(context.Context, internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	panic("composed catalog test does not drive the per-product reads")
}

func (r *composedCatalogOracleReader) GetTaxInputs(context.Context, internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	panic("composed catalog test does not drive the per-product reads")
}

// composedTenantLookup is the seat that owns the tenant's stored assortment
// rule. Only the routing reader reads it, which is the point: every layer above
// asks with a nil policy and gets the rule resolved here, exactly once.
type composedTenantLookup struct {
	cfg tenant_config.Config
}

func (l composedTenantLookup) Get(context.Context, string) (tenant_config.Config, error) {
	return l.cfg, nil
}

// newComposedCatalogHandler builds the LIVE chain, in the live order:
// source → cache → timing → routing → application service → HTTP handler
// (internal/composition/root.go:450-481, 527-531). Wiring the cache straight to
// the handler — as this test used to — skips the routing seam, so the handler's
// nil policy reaches the cache unresolved and every read fails; the composition
// this test exists to prove is precisely the one that resolves it.
func newComposedCatalogHandler(oracle internalreadports.Source, c *cacheadapter.Cache, assortment tenant_config.SellableAssortment) http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	live := internalreadobservability.NewTimingReader(cacheadapter.NewReader(oracle, c), logger, time.Second)
	lookup := composedTenantLookup{cfg: tenant_config.Config{
		TenantID:           "tenant-composed",
		Source:             tenant_config.SourceSankhya,
		SellableAssortment: assortment,
	}}
	routed := routing.NewReader(nil, live, lookup, "tenant-composed")

	mux := httpx.NewRouteClassMux()
	(catalogtransport.Handler{
		PageReader: internalreadapp.NewService(routed),
	}).Register(mux)
	return mux
}

// sellableAssortment is the cut the tenant stores: the three toggles on.
func sellableAssortment() tenant_config.SellableAssortment {
	return tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}
}

type composedCatalogResponse struct {
	AsOf  time.Time `json:"as_of"`
	Items []struct {
		InternalProductID int `json:"internal_product_id"`
	} `json:"items"`
}

func (r composedCatalogResponse) ids() []int {
	ids := make([]int, 0, len(r.Items))
	for _, item := range r.Items {
		ids = append(ids, item.InternalProductID)
	}
	return ids
}

func getComposedCatalog(t *testing.T, handler http.Handler, noCache bool) composedCatalogResponse {
	t.Helper()
	return getComposedCatalogURL(t, handler, "/catalog/products", noCache)
}

func getComposedCatalogURL(t *testing.T, handler http.Handler, url string, noCache bool) composedCatalogResponse {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, url, nil)
	if noCache {
		request.Header.Set("Cache-Control", "no-cache")
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /catalog/products status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var response composedCatalogResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode catalog response: %v; body=%s", err, recorder.Body.String())
	}
	if response.AsOf.IsZero() {
		t.Fatal("catalog response has zero as_of")
	}
	return response
}

func newComposedCache(clock *composedFakeClock) *cacheadapter.Cache {
	return cacheadapter.New(cacheadapter.Config{
		Clock: clock,
		Policies: map[string]internalreaddomain.FreshnessPolicy{
			cacheadapter.ClassCatalog:   {MaxAge: 5 * time.Minute},
			cacheadapter.ClassPriceCost: {MaxAge: 2 * time.Minute},
		},
	})
}

func TestComposedCatalogHTTPNoCacheBypassesAndRepopulates(t *testing.T) {
	clock := &composedFakeClock{now: time.Date(2026, 7, 14, 16, 0, 0, 0, time.UTC)}
	oracle := &composedCatalogOracleReader{clock: clock}
	handler := newComposedCatalogHandler(oracle, newComposedCache(clock), sellableAssortment())

	warm := getComposedCatalog(t, handler, false)
	warmAgain := getComposedCatalog(t, handler, false)
	if got := oracle.listCalls.Load(); got != 1 {
		t.Fatalf("warm Oracle calls=%d, want 1", got)
	}
	if !warm.AsOf.Equal(warmAgain.AsOf) {
		t.Fatalf("warm as_of values=%s/%s, want identical", warm.AsOf, warmAgain.AsOf)
	}

	clock.Advance(time.Second)
	bypassed := getComposedCatalog(t, handler, true)
	if got := oracle.listCalls.Load(); got != 2 {
		t.Fatalf("bypass Oracle calls=%d, want 2", got)
	}
	if !bypassed.AsOf.After(warm.AsOf) {
		t.Fatalf("bypass as_of=%s, want strictly newer than warm %s", bypassed.AsOf, warm.AsOf)
	}

	repopulated := getComposedCatalog(t, handler, false)
	if got := oracle.listCalls.Load(); got != 2 {
		t.Fatalf("post-bypass Oracle calls=%d, want 2", got)
	}
	if !repopulated.AsOf.Equal(bypassed.AsOf) {
		t.Fatalf("repopulated as_of=%s, want bypass as_of %s", repopulated.AsOf, bypassed.AsOf)
	}
}

// TestComposedCatalogCutTravelsWireToSource drives the two answers the wire can
// give — the tenant's stored rule (no parameter) and the explicit whole catalog
// (include_all=true) — through the live chain to the source, and back.
//
// It is the end-to-end form of the 503 this fold came from: the handler asks
// with a nil policy, and only the routing seam knows the tenant's rule. If any
// layer between them drops the policy, the source refuses; if any layer forges
// one, the two reads return the same page. Both failures are visible here in
// the response body, not in a log line.
func TestComposedCatalogCutTravelsWireToSource(t *testing.T) {
	clock := &composedFakeClock{now: time.Date(2026, 7, 14, 16, 0, 0, 0, time.UTC)}
	oracle := &composedCatalogOracleReader{clock: clock}
	handler := newComposedCatalogHandler(oracle, newComposedCache(clock), sellableAssortment())

	cut := getComposedCatalog(t, handler, false)
	if got := cut.ids(); len(got) != 1 || got[0] != 101 {
		t.Fatalf("stored-rule read returned %v, want the sellable cut [101]", got)
	}

	uncut := getComposedCatalogURL(t, handler, "/catalog/products?include_all=true", false)
	if got := uncut.ids(); len(got) != 2 {
		t.Fatalf("include_all read returned %v, want the whole catalog [101 202]", got)
	}
	if got := oracle.listCalls.Load(); got != 2 {
		t.Fatalf("source list calls=%d, want 2 — two cuts are two answers, so they cannot share a cache entry", got)
	}

	// And the entries are stable: re-reading each cut serves its own cached page.
	if got := getComposedCatalog(t, handler, false).ids(); len(got) != 1 || got[0] != 101 {
		t.Fatalf("second stored-rule read returned %v, want [101]", got)
	}
	if got := getComposedCatalogURL(t, handler, "/catalog/products?include_all=true", false).ids(); len(got) != 2 {
		t.Fatalf("second include_all read returned %v, want two items", got)
	}
	if got := oracle.listCalls.Load(); got != 2 {
		t.Fatalf("source list calls after the warm re-reads=%d, want 2", got)
	}
}

type composedBatchOracleReader struct {
	clock     *composedFakeClock
	costCalls atomic.Int64
}

func (r *composedBatchOracleReader) GetCostFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.CostAsOf, error) {
	r.costCalls.Add(1)
	amount := 12.50
	return map[int64]*internalreaddomain.CostAsOf{
		7: {ProductID: 7, Amount: &amount, Source: internalreaddomain.SourceMetadata{System: "fake-oracle"}},
	}, nil
}

func (r *composedBatchOracleReader) GetTaxFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.TaxInputs, error) {
	return map[int64]*internalreaddomain.TaxInputs{}, nil
}

func (r *composedBatchOracleReader) GetICMSCeilingByOrigin(context.Context, internalreaddomain.UF) (map[internalreaddomain.UF]*internalreaddomain.ICMSCeiling, error) {
	return map[internalreaddomain.UF]*internalreaddomain.ICMSCeiling{}, nil
}

type composedOrderLookup struct {
	order ordersdomain.MarketplaceOrder
}

func (r composedOrderLookup) FindExactOrder(context.Context, ordersdomain.LinkageScope) (ordersdomain.MarketplaceOrder, bool, error) {
	return r.order, true, nil
}

type composedSankhyaReader struct {
	candidate ordersdomain.AssistedSankhyaCandidate
}

func (r composedSankhyaReader) ValidateConfiguration(context.Context) error { return nil }
func (r composedSankhyaReader) ConfigurationRevision() string               { return "cfg-composed-1" }
func (r composedSankhyaReader) EvidenceReference() string                   { return "evidence-composed-1" }
func (r composedSankhyaReader) FindCandidates(context.Context, string) ([]ordersdomain.AssistedSankhyaCandidate, error) {
	return []ordersdomain.AssistedSankhyaCandidate{r.candidate}, nil
}
func (r composedSankhyaReader) ListDescendants(context.Context, ordersdomain.InternalDocumentLineIdentity, *float64) (ordersdomain.AssistedSankhyaLineage, error) {
	return ordersdomain.AssistedSankhyaLineage{State: ordersdomain.AssistedSankhyaLineageNone}, nil
}

type composedLinkageRepository struct{}

func (composedLinkageRepository) LoadCurrent(context.Context, ordersdomain.LinkageScope) (ordersdomain.SankhyaLinkage, bool, error) {
	return ordersdomain.SankhyaLinkage{}, false, nil
}
func (composedLinkageRepository) AppendConfirmation(_ context.Context, linkage ordersdomain.SankhyaLinkage) (ordersdomain.SankhyaLinkage, error) {
	return linkage, nil
}

func TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost(t *testing.T) {
	clock := &composedFakeClock{now: time.Date(2026, 7, 14, 16, 0, 0, 0, time.UTC)}
	catalogOracle := &composedCatalogOracleReader{clock: clock}
	c := newComposedCache(clock)
	handler := newComposedCatalogHandler(catalogOracle, c, sellableAssortment())
	pricecostOracle := &composedBatchOracleReader{clock: clock}
	pricecost := cacheadapter.NewBatchReader(pricecostOracle, c)

	warm := getComposedCatalog(t, handler, false)
	if got := catalogOracle.listCalls.Load(); got != 1 {
		t.Fatalf("warm catalog Oracle calls=%d, want 1", got)
	}
	if _, err := pricecost.GetCostFactsByIDs(context.Background(), []int64{7}); err != nil {
		t.Fatalf("warm pricecost: %v", err)
	}

	lineOne := ordersdomain.MPCLineID("mpl_0123456789abcdef0123456789abcdef")
	origin := ordersdomain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 1}
	candidate := ordersdomain.AssistedSankhyaCandidate{
		Header:        ordersdomain.InternalDocumentIdentity{DocumentID: 31301},
		OperationCode: 313,
		Lines:         []ordersdomain.AssistedSankhyaCandidateLine{{Identity: origin}},
	}
	service := ordersapp.NewAssistedSankhyaLinkageService(ordersapp.AssistedSankhyaLinkageServiceConfig{
		Orders: composedOrderLookup{order: ordersdomain.MarketplaceOrder{
			InstallationID: "install-composed", ProviderCode: "mercado_livre", ProviderOrderID: "order-composed",
			Items: []ordersdomain.MarketplaceOrderItem{{MPCLineID: lineOne, ReconciliationState: ordersdomain.LineReconciliationStable}},
		}},
		Reader:      composedSankhyaReader{candidate: candidate},
		Linkages:    composedLinkageRepository{},
		Now:         func() time.Time { return clock.Now() },
		NewEventID:  func() (string, error) { return "event-composed-1", nil },
		Invalidator: c,
	})
	clock.Advance(time.Second)
	if _, err := service.Confirm(context.Background(), ordersapp.ConfirmAssistedSankhyaLinkageInput{
		TenantID: "tenant-composed", InstallationID: "install-composed", ProviderOrderID: "order-composed",
		SelectedDocumentID: 31301,
		Selections:         []ordersdomain.AssistedSankhyaLineSelection{{MPCLineID: lineOne, CandidateLine: origin}},
		ActorID:            "operator-composed", Reason: "confirmed", IdempotencyKey: "idem-composed",
		SourceAt: time.Date(2026, 7, 14, 15, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("successful Confirm: %v", err)
	}

	if _, err := pricecost.GetCostFactsByIDs(context.Background(), []int64{7}); err != nil {
		t.Fatalf("post-confirm pricecost: %v", err)
	}
	if got := pricecostOracle.costCalls.Load(); got != 1 {
		t.Fatalf("pricecost Oracle calls=%d, want warm entry preserved at 1", got)
	}

	refreshed := getComposedCatalog(t, handler, false)
	if got := catalogOracle.listCalls.Load(); got != 2 {
		t.Fatalf("post-confirm catalog Oracle calls=%d, want 2", got)
	}
	if !refreshed.AsOf.After(warm.AsOf) {
		t.Fatalf("post-confirm as_of=%s, want strictly newer than warm %s", refreshed.AsOf, warm.AsOf)
	}
}

var _ internalreadports.Source = (*composedCatalogOracleReader)(nil)
var _ internalreadports.BatchReader = (*composedBatchOracleReader)(nil)
var _ ordersports.OrderLookup = composedOrderLookup{}
var _ ordersports.SankhyaLinkageReader = composedSankhyaReader{}
var _ ordersports.SankhyaLinkageRepository = composedLinkageRepository{}
