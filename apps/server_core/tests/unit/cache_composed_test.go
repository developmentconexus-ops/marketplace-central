package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	catalogtransport "marketplace-central/apps/server_core/internal/modules/catalog/transport"
	cacheadapter "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
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

type composedCatalogOracleReader struct {
	clock       *composedFakeClock
	listCalls   atomic.Int64
	searchCalls atomic.Int64
}

func (r *composedCatalogOracleReader) ListCatalogProductFacts(context.Context, internalreadports.Cursor, int) (internalreadports.CatalogFactPage, error) {
	r.listCalls.Add(1)
	return internalreadports.CatalogFactPage{
		Items: []internalreadports.CatalogProductFact{{InternalProductID: 101, Active: true}},
		AsOf:  r.clock.Now(),
	}, nil
}

func (r *composedCatalogOracleReader) SearchCatalogProductFacts(context.Context, string, int) (internalreadports.CatalogFactPage, error) {
	r.searchCalls.Add(1)
	return internalreadports.CatalogFactPage{AsOf: r.clock.Now()}, nil
}

func newComposedCatalogHandler(oracle internalreadports.CatalogPageReader, c *cacheadapter.Cache) http.Handler {
	mux := httpx.NewRouteClassMux()
	(catalogtransport.Handler{
		PageReader: cacheadapter.NewCatalogPageReader(oracle, c),
	}).Register(mux)
	return mux
}

type composedCatalogResponse struct {
	AsOf time.Time `json:"as_of"`
}

func getComposedCatalog(t *testing.T, handler http.Handler, noCache bool) composedCatalogResponse {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/catalog/products", nil)
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
	handler := newComposedCatalogHandler(oracle, newComposedCache(clock))

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
	handler := newComposedCatalogHandler(catalogOracle, c)
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

var _ internalreadports.BatchReader = (*composedBatchOracleReader)(nil)
var _ ordersports.OrderLookup = composedOrderLookup{}
var _ ordersports.SankhyaLinkageReader = composedSankhyaReader{}
var _ ordersports.SankhyaLinkageRepository = composedLinkageRepository{}
