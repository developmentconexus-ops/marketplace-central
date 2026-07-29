package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type fakeCatalogPageReader struct {
	listPages       map[int64]ports.CatalogFactPage
	searchPage      ports.CatalogFactPage
	listCursors     []ports.Cursor
	searchQueries   []string
	searchCursors   []ports.Cursor
	searchLimits    []int
	freshnessPolicy []internalreaddomain.FreshnessPolicy
	activeSources   []erpdomain.ImportSource
	sourcePresent   []bool
	idAsks          [][]int64
	idPage          ports.CatalogFactPage
	err             error
	policies        []*ports.SellableAssortmentPolicy
}

func (f *fakeCatalogPageReader) captureSource(ctx context.Context) {
	source, present := erpinternalread.ActiveSourceFromContext(ctx)
	f.activeSources = append(f.activeSources, source)
	f.sourcePresent = append(f.sourcePresent, present)
}

func (f *fakeCatalogPageReader) ListCatalogProductFacts(ctx context.Context, cursor ports.Cursor, _ int) (ports.CatalogFactPage, error) {
	f.listCursors = append(f.listCursors, cursor)
	f.captureSource(ctx)
	if policy, ok := FreshnessPolicyFromContext(ctx); ok {
		f.freshnessPolicy = append(f.freshnessPolicy, policy)
	}
	if f.err != nil {
		return ports.CatalogFactPage{}, f.err
	}
	return f.listPages[cursor.InternalProductID], nil
}

func (f *fakeCatalogPageReader) CatalogProductFactsByIDs(ctx context.Context, ids []int64) (ports.CatalogFactPage, error) {
	f.idAsks = append(f.idAsks, ids)
	f.captureSource(ctx)
	if f.err != nil {
		return ports.CatalogFactPage{}, f.err
	}
	return f.idPage, nil
}
func (f *fakeCatalogPageReader) SearchCatalogProductFacts(ctx context.Context, q string, cursor ports.Cursor, limit int) (ports.CatalogFactPage, error) {
	f.searchQueries = append(f.searchQueries, q)
	f.searchCursors = append(f.searchCursors, cursor)
	f.searchLimits = append(f.searchLimits, limit)
	f.captureSource(ctx)
	if policy, ok := FreshnessPolicyFromContext(ctx); ok {
		f.freshnessPolicy = append(f.freshnessPolicy, policy)
	}
	if f.err != nil {
		return ports.CatalogFactPage{}, f.err
	}
	return f.searchPage, nil
}

func (f *fakeCatalogPageReader) ListCatalogProductFactsWithPolicy(ctx context.Context, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	f.policies = append(f.policies, policy)
	return f.ListCatalogProductFacts(ctx, cursor, limit)
}

func (f *fakeCatalogPageReader) SearchCatalogProductFactsWithPolicy(ctx context.Context, q string, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	f.policies = append(f.policies, policy)
	return f.SearchCatalogProductFacts(ctx, q, cursor, limit)
}

func (f *fakeCatalogPageReader) GetCatalogAssortmentCounts(context.Context, *ports.SellableAssortmentPolicy) (ports.CatalogAssortmentCounts, error) {
	return ports.CatalogAssortmentCounts{}, f.err
}

func TestHandlerCatalogDefaultsFilteredAndIncludeAllIsRequestLocal(t *testing.T) {
	fake := &fakeCatalogPageReader{listPages: map[int64]ports.CatalogFactPage{
		0: {Items: []ports.CatalogProductFact{}, AsOf: time.Now().UTC()},
	}}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)

	for _, path := range []string{"/catalog/products", "/catalog/products?include_all=true", "/catalog/products"} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d, body = %s", path, recorder.Code, recorder.Body.String())
		}
	}
	// nil is the request "apply whatever rule this tenant configured", which only
	// the routing seam can answer; a named policy crosses as itself. The handler
	// must never author a third thing.
	all := ports.AllProductsAssortment()
	want := []*ports.SellableAssortmentPolicy{nil, &all, nil}
	if !reflect.DeepEqual(fake.policies, want) {
		t.Fatalf("policies read = %+v, want %+v", fake.policies, want)
	}
}

func TestHandlerCatalogCountsAndMalformedIncludeAll(t *testing.T) {
	fake := &fakeCatalogPageReader{}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)

	counts := httptest.NewRecorder()
	mux.ServeHTTP(counts, httptest.NewRequest(http.MethodGet, "/catalog/products/counts", nil))
	if counts.Code != http.StatusOK || trimJSON(counts.Body.String()) != `{"sellable_count":0,"total_count":0}` {
		t.Fatalf("counts status/body = %d/%s", counts.Code, trimJSON(counts.Body.String()))
	}
	for _, path := range []string{"/catalog/products?include_all=TRUE", "/catalog/products?include_all=1", "/catalog/products/search?include_all=yes"} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("GET %s status = %d, body = %s", path, recorder.Code, recorder.Body.String())
		}
	}
}

// A screen that already knows which products it needs asks for exactly those.
// The ids it asks for reach the reader verbatim, and an id the active source
// does not carry is simply missing from the answer — the caller can tell which
// products the source knows nothing about.
func TestCatalogProductsByIDsAsksForExactlyTheRequestedProducts(t *testing.T) {
	fake := &fakeCatalogPageReader{idPage: ports.CatalogFactPage{
		Items: []ports.CatalogProductFact{catalogFact(17413), catalogFact(90034)},
		AsOf:  time.Date(2026, 7, 25, 3, 0, 0, 0, time.UTC),
	}}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/catalog/products?ids=17413,90034,24864", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !reflect.DeepEqual(fake.idAsks, [][]int64{{17413, 90034, 24864}}) {
		t.Fatalf("id asks = %+v", fake.idAsks)
	}
	if len(fake.listCursors) != 0 {
		t.Fatalf("an id ask must not fall through to the keyset page: %+v", fake.listCursors)
	}

	var response catalogPageEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Items) != 2 || response.Items[0].InternalProductID != 17413 {
		t.Fatalf("items = %+v", response.Items)
	}
	// The answer is a set, not a page of a longer sequence.
	if response.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want null", *response.NextCursor)
	}
}

func TestCatalogProductsByIDsRejectsAnUnusableIDList(t *testing.T) {
	for _, path := range []string{
		"/catalog/products?ids=",
		"/catalog/products?ids=abc",
		"/catalog/products?ids=0",
		"/catalog/products?ids=-3",
		"/catalog/products?ids=1&ids=2",
	} {
		fake := &fakeCatalogPageReader{}
		mux := httpx.NewRouteClassMux()
		(Handler{PageReader: fake}).Register(mux)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("GET %s status = %d, body = %s", path, recorder.Code, recorder.Body.String())
		}
		if len(fake.idAsks) != 0 {
			t.Fatalf("GET %s reached the port with %+v", path, fake.idAsks)
		}
	}
}

func TestCatalogPageRoutesFollowThreePageCursorChain(t *testing.T) {
	page1Cursor := ports.Cursor{InternalProductID: 1}
	page2Cursor := ports.Cursor{InternalProductID: 2}
	fake := &fakeCatalogPageReader{listPages: map[int64]ports.CatalogFactPage{
		0: {Items: []ports.CatalogProductFact{catalogFact(1)}, NextCursor: &page1Cursor, AsOf: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)},
		1: {Items: []ports.CatalogProductFact{catalogFact(2)}, NextCursor: &page2Cursor, AsOf: time.Date(2026, 7, 14, 12, 1, 0, 0, time.UTC)},
		2: {Items: []ports.CatalogProductFact{catalogFact(3)}, AsOf: time.Date(2026, 7, 14, 12, 2, 0, 0, time.UTC)},
	}}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)

	var nextCursor *string
	for pageNumber := 0; pageNumber < 3; pageNumber++ {
		path := "/catalog/products?limit=1"
		if nextCursor != nil {
			path += "&cursor=" + *nextCursor
		}
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("page %d status = %d, body = %s", pageNumber+1, recorder.Code, recorder.Body.String())
		}
		var response catalogPageEnvelope
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("page %d decode: %v", pageNumber+1, err)
		}
		t.Logf("GET %s -> %d %s", path, recorder.Code, trimJSON(recorder.Body.String()))
		if response.PageSize != 1 || len(response.Items) != 1 || response.Items[0].InternalProductID != pageNumber+1 {
			t.Fatalf("page %d response = %+v", pageNumber+1, response)
		}
		nextCursor = response.NextCursor
	}
	if nextCursor != nil {
		t.Fatalf("last page next_cursor = %q, want null", *nextCursor)
	}
	if !reflect.DeepEqual(fake.listCursors, []ports.Cursor{{}, {InternalProductID: 1}, {InternalProductID: 2}}) {
		t.Fatalf("list cursors = %+v", fake.listCursors)
	}
}

func TestCatalogPageRoutesValidateBeforePortCall(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		body string
	}{
		{name: "garbage cursor", path: "/catalog/products?cursor=%25%25%25garbage", body: `{"error":"invalid_cursor"}`},
		{name: "zero limit", path: "/catalog/products?limit=0", body: `{"error":"invalid_limit","allowed_range":"1..100"}`},
		{name: "over limit", path: "/catalog/products?limit=101", body: `{"error":"invalid_limit","allowed_range":"1..100"}`},
		{name: "search over limit", path: "/catalog/products/search?q=PARAFUSO&limit=51", body: `{"error":"invalid_limit","allowed_range":"1..50"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeCatalogPageReader{}
			mux := httpx.NewRouteClassMux()
			(Handler{PageReader: fake}).Register(mux)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if recorder.Code != http.StatusBadRequest || trimJSON(recorder.Body.String()) != trimJSON(tc.body) {
				t.Fatalf("status/body = %d/%q, want 400/%q", recorder.Code, trimJSON(recorder.Body.String()), tc.body)
			}
			t.Logf("GET %s -> %d %s", tc.path, recorder.Code, trimJSON(recorder.Body.String()))
			if len(fake.listCursors) != 0 || len(fake.searchQueries) != 0 {
				t.Fatal("invalid request reached page port")
			}
		})
	}
}

func TestCatalogSearchPageEnvelopeAndNoCachePolicy(t *testing.T) {
	fake := &fakeCatalogPageReader{searchPage: ports.CatalogFactPage{
		Items:      []ports.CatalogProductFact{catalogFact(10)},
		NextCursor: &ports.Cursor{InternalProductID: 11},
		AsOf:       time.Date(2026, 7, 14, 13, 0, 0, 0, time.UTC),
	}}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/catalog/products/search?q=PARAFUSO&limit=50", nil)
	request.Header.Set("Cache-Control", "max-age=0, no-cache")
	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	t.Logf("GET %s -> %d %s", request.URL.String(), recorder.Code, trimJSON(recorder.Body.String()))
	var response catalogPageEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	// Search is a page of a longer sequence. Withholding its cursor made a
	// truncated result byte-identical to a complete one, and left the caller no
	// way to reach the matches beyond the page.
	if response.NextCursor == nil || response.PageSize != 1 || response.AsOf.IsZero() {
		t.Fatalf("response = %+v", response)
	}
	if !reflect.DeepEqual(fake.searchQueries, []string{"PARAFUSO"}) || !reflect.DeepEqual(fake.searchLimits, []int{50}) {
		t.Fatalf("search calls = %v/%v", fake.searchQueries, fake.searchLimits)
	}
	if !reflect.DeepEqual(fake.searchCursors, []ports.Cursor{{}}) {
		t.Fatalf("search cursors = %+v, want the first page", fake.searchCursors)
	}
	if !reflect.DeepEqual(fake.freshnessPolicy, []internalreaddomain.FreshnessPolicy{{MaxAge: 0}}) {
		t.Fatalf("freshness policy = %+v", fake.freshnessPolicy)
	}

	// The cursor the response handed out must come back as the next page's ask.
	next := httptest.NewRecorder()
	mux.ServeHTTP(next, httptest.NewRequest(http.MethodGet, "/catalog/products/search?q=PARAFUSO&limit=50&cursor="+*response.NextCursor, nil))
	if next.Code != http.StatusOK {
		t.Fatalf("second page status = %d, body = %s", next.Code, next.Body.String())
	}
	if !reflect.DeepEqual(fake.searchCursors, []ports.Cursor{{}, {InternalProductID: 11}}) {
		t.Fatalf("search cursors = %+v, want the round-tripped cursor", fake.searchCursors)
	}

	bad := httptest.NewRecorder()
	mux.ServeHTTP(bad, httptest.NewRequest(http.MethodGet, "/catalog/products/search?q=PARAFUSO&cursor=not-base64", nil))
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("invalid cursor status = %d, body = %s", bad.Code, bad.Body.String())
	}
}

func TestCatalogIdentityJSONNullNotEmptyAndReferenceAlias(t *testing.T) {
	t.Run("identity and collision flags", func(t *testing.T) {
		ean, reference, brand, ncm := "4006381333931", "MF-10", "Marca", "12345678"
		fact := catalogFact(10)
		fact.EAN, fact.Reference, fact.ManufacturerReference, fact.BrandName, fact.NCM = &ean, &reference, &reference, &brand, &ncm
		fact.QualityFlags = []string{"complete", "ean_collision"}
		assertCatalogIdentityJSON(t, fact, `"ean":"4006381333931"`, `"reference":"MF-10"`, `"manufacturer_reference":"MF-10"`, `"brand_name":"Marca"`, `"ncm":"12345678"`, `"quality_flags":["complete","ean_collision"]`)
	})
	t.Run("invalid_ean unknowns are null not empty", func(t *testing.T) {
		fact := catalogFact(11)
		fact.Reference, fact.ManufacturerReference, fact.EAN, fact.BrandName, fact.NCM = nil, nil, nil, nil, nil
		fact.QualityFlags = []string{"invalid_ean"}
		assertCatalogIdentityJSON(t, fact, `"ean":null`, `"reference":null`, `"manufacturer_reference":null`, `"brand_name":null`, `"ncm":null`, `"quality_flags":["invalid_ean"]`)
	})
}

func assertCatalogIdentityJSON(t *testing.T, fact ports.CatalogProductFact, fragments ...string) {
	t.Helper()
	fake := &fakeCatalogPageReader{listPages: map[int64]ports.CatalogFactPage{0: {Items: []ports.CatalogProductFact{fact}, AsOf: time.Now().UTC()}}}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/catalog/products", nil))
	for _, fragment := range fragments {
		if !strings.Contains(recorder.Body.String(), fragment) {
			t.Fatalf("missing %s in %s", fragment, recorder.Body.String())
		}
	}
	if strings.Contains(recorder.Body.String(), `:""`) {
		t.Fatalf("unknown identity emitted empty string: %s", recorder.Body.String())
	}
}

func TestCatalogPageRoutesMapSourceAndDeadlineErrors(t *testing.T) {
	for _, tc := range []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "source unavailable", err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "driver: password=secret", errors.New("driver detail")), wantStatus: http.StatusServiceUnavailable, wantBody: `{"error":"source_unavailable"}`},
		{name: "deadline exceeded", err: context.DeadlineExceeded, wantStatus: http.StatusGatewayTimeout, wantBody: `{"error":"deadline_exceeded"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeCatalogPageReader{err: tc.err}
			mux := httpx.NewRouteClassMux()
			(Handler{PageReader: fake}).Register(mux)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/catalog/products", nil))
			if recorder.Code != tc.wantStatus || trimJSON(recorder.Body.String()) != tc.wantBody {
				t.Fatalf("status/body = %d/%q, want %d/%q", recorder.Code, trimJSON(recorder.Body.String()), tc.wantStatus, tc.wantBody)
			}
			t.Logf("GET /catalog/products -> %d %s", recorder.Code, trimJSON(recorder.Body.String()))
			if strings.Contains(recorder.Body.String(), "secret") || strings.Contains(recorder.Body.String(), "driver") {
				t.Fatal("driver detail leaked into response")
			}
		})
	}
}

func catalogFact(id int64) ports.CatalogProductFact {
	quantity := 41.0
	amount := "12.90"
	reference := "ABC-123"
	description := "PARAFUSO SEXTAVADO M8"
	return ports.CatalogProductFact{
		InternalProductID: id,
		Reference:         &reference,
		Description:       &description,
		Active:            true,
		SellableStock:     ports.CatalogQuantityFact{Quantity: &quantity, Quality: []string{}},
		CurrentPrice:      ports.CatalogMoneyFact{Amount: &amount, Currency: "BRL", Quality: []string{}},
		Cost:              ports.CatalogMoneyFact{Amount: nil, Currency: "BRL", Quality: []string{"missing_cost"}},
	}
}

func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

// GOAL C — erp_source active-source toggle on the catalog read seam.

func TestCatalogPageThreadsActiveSourceToggle(t *testing.T) {
	for _, tc := range []struct {
		name        string
		path        string
		wantSource  erpdomain.ImportSource
		wantPresent bool
	}{
		{name: "list absent = default (byte-stable)", path: "/catalog/products?limit=1", wantPresent: false},
		{name: "list catalogo_cliente", path: "/catalog/products?limit=1&erp_source=catalogo_cliente", wantSource: erpdomain.SourceCatalogoCliente, wantPresent: true},
		{name: "list xlsx explicit", path: "/catalog/products?limit=1&erp_source=xlsx", wantSource: erpdomain.SourceXLSX, wantPresent: true},
		{name: "search catalogo_cliente", path: "/catalog/products/search?q=PARAFUSO&limit=1&erp_source=catalogo_cliente", wantSource: erpdomain.SourceCatalogoCliente, wantPresent: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeCatalogPageReader{}
			mux := httpx.NewRouteClassMux()
			(Handler{PageReader: fake}).Register(mux)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
			}
			t.Logf("GET %s -> %d", tc.path, recorder.Code)
			if len(fake.sourcePresent) != 1 || fake.sourcePresent[0] != tc.wantPresent {
				t.Fatalf("source present = %v, want [%v]", fake.sourcePresent, tc.wantPresent)
			}
			if tc.wantPresent && fake.activeSources[0] != tc.wantSource {
				t.Fatalf("active source = %q, want %q", fake.activeSources[0], tc.wantSource)
			}
		})
	}
}

func TestCatalogPageRejectsUnknownActiveSource(t *testing.T) {
	for _, path := range []string{
		"/catalog/products?limit=1&erp_source=bogus",
		"/catalog/products/search?q=X&limit=1&erp_source=sankhya",
	} {
		fake := &fakeCatalogPageReader{}
		mux := httpx.NewRouteClassMux()
		(Handler{PageReader: fake}).Register(mux)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("GET %s status = %d, want 400 (body %s)", path, recorder.Code, recorder.Body.String())
		}
		if trimJSON(recorder.Body.String()) != trimJSON(`{"error":"invalid_erp_source"}`) {
			t.Fatalf("GET %s body = %s, want invalid_erp_source", path, trimJSON(recorder.Body.String()))
		}
		// Unknown value must never silently fall through to a dataset read.
		if len(fake.listCursors) != 0 || len(fake.searchQueries) != 0 {
			t.Fatalf("GET %s reached the page port on an invalid erp_source", path)
		}
	}
}
