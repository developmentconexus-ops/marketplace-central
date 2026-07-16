//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	listingsconnectors "marketplace-central/apps/server_core/internal/modules/listings/adapters/connectors"
	listingsintegrations "marketplace-central/apps/server_core/internal/modules/listings/adapters/integrations"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingstransport "marketplace-central/apps/server_core/internal/modules/listings/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type stubListingPage struct {
	items []connectorsdomain.ListingSnapshot
	err   error
}

type stubListingReader struct {
	mu       sync.Mutex
	pages    map[string]stubListingPage
	gate     <-chan struct{}
	gateOnce sync.Once
}

func (s *stubListingReader) ListListings(_ context.Context, input connectorsdomain.ListListingsInput) ([]connectorsdomain.ListingSnapshot, error) {
	s.gateOnce.Do(func() {
		if s.gate != nil {
			<-s.gate
		}
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	page := s.pages[input.Cursor]
	return append([]connectorsdomain.ListingSnapshot(nil), page.items...), page.err
}

func (*stubListingReader) ReadListing(context.Context, connectorsdomain.ProviderListingRef) (connectorsdomain.ListingSnapshot, error) {
	return connectorsdomain.ListingSnapshot{}, nil
}

func (s *stubListingReader) setPages(pages map[string]stubListingPage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages = pages
}

type refreshHarness struct {
	pool         *pgxpool.Pool
	mux          http.Handler
	operations   *integrationsapp.OperationService
	tenant       string
	installation string
	reports      chan error
}

func newRefreshHarness(t *testing.T, tag string, pageSize int, stub *stubListingReader) refreshHarness {
	t.Helper()
	pool, cfg := testpostgres.OpenPool(t, tag)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant, installation := cfg.DefaultTenantID+"-f01-"+token, "f01-inst-"+token
	seedConnectedInstallation(t, pool, tenant, installation)
	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(pool, tenant), tenant)
	operationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(pool, tenant), tenant)
	caps := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercadolivre", Listings: stub}})
	listingRepo := listingspostgres.NewRepository(pool, tenant)
	ingestion := listingsapp.NewIngestion(listingsconnectors.NewSource(caps), listingRepo, pageSize, time.Now)
	gateway := listingsintegrations.NewGateway(installationSvc, operationSvc)
	reports := make(chan error, 1)
	refreshSvc := listingsapp.NewRefreshService(gateway, ingestion, func(task func()) { go task() }, time.Now,
		func(err error) { reports <- err },
		func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })
	mux := httpx.NewRouteClassMux()
	listingstransport.NewRefreshHandler(refreshSvc).Register(mux)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(ctx, `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(ctx, `DELETE FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(ctx, `DELETE FROM integration_installations WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})
	return refreshHarness{pool: pool, mux: mux, operations: operationSvc, tenant: tenant, installation: installation, reports: reports}
}

func seedConnectedInstallation(t *testing.T, pool *pgxpool.Pool, tenant, installation string) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES('mercadolivre','marketplace','Mercado Livre','oauth2','interactive') ON CONFLICT(provider_code) DO NOTHING`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO integration_installations(installation_id,tenant_id,provider_code,family,display_name,status,health_status) VALUES($1,$2,'mercadolivre','marketplace','F01','connected','healthy')`, installation, tenant); err != nil {
		t.Fatal(err)
	}
}

func postRefresh(t *testing.T, mux http.Handler, installation string) (int, map[string]any) {
	t.Helper()
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/listings/refresh", strings.NewReader(fmt.Sprintf(`{"installation_id":%q}`, installation))))
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response %d %q: %v", rr.Code, rr.Body.String(), err)
	}
	return rr.Code, body
}

func operationID(t *testing.T, body map[string]any) string {
	t.Helper()
	id, _ := body["operation_run_id"].(string)
	if id == "" {
		t.Fatalf("operation_run_id missing from %#v", body)
	}
	return id
}

func waitTerminal(t *testing.T, svc *integrationsapp.OperationService, installation, id string) string {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		runs, err := svc.ListByInstallation(context.Background(), installation)
		if err != nil {
			t.Fatal(err)
		}
		for _, run := range runs {
			if run.OperationRunID == id && (run.Status == "succeeded" || run.Status == "failed") {
				return string(run.Status)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("operation %s did not reach a terminal state within 2s", id)
	return ""
}

func listingSnapshots() []connectorsdomain.ListingSnapshot {
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	price := "129.90"
	quantity := 7
	return []connectorsdomain.ListingSnapshot{
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0001", ProviderStatus: "active", Title: "Alpha", PriceAmount: &price, PriceCurrency: "BRL", AvailableQuantity: &quantity, FetchedAt: now},
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0002", ProviderStatus: "paused", Title: "Beta", FetchedAt: now},
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0003", ProviderStatus: "closed", Title: "Gamma", PriceAmount: &price, PriceCurrency: "BRL", FetchedAt: now},
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0004", ProviderStatus: "provider-new-state", Title: "Delta", FetchedAt: now},
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0005", ProviderStatus: "active", Title: "Epsilon", FetchedAt: now, Variations: []connectorsdomain.ListingVariationSnapshot{{ProviderVariationID: "V-5", AvailableQuantity: &quantity}}},
		{ProviderCode: "mercadolivre", ProviderItemID: "MLBTEST0006", ProviderStatus: "paused", Title: "Zeta", FetchedAt: now},
	}
}

func TestListingsRefreshSeedsIC02RowsAndClosesMissing(t *testing.T) {
	items := listingSnapshots()
	stub := &stubListingReader{pages: map[string]stubListingPage{"": {items: items}}}
	h := newRefreshHarness(t, "f01-refresh-complete", 100, stub)
	status, body := postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("status = %d, body = %#v", status, body)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("terminal status = %s", got)
	}
	var count, nullSales, dashVariation, unknown int
	err := h.pool.QueryRow(context.Background(), `SELECT count(*),count(*) FILTER (WHERE sales_30d IS NULL),count(*) FILTER (WHERE variation_id='-'),count(*) FILTER (WHERE status='unknown') FROM listings WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, h.installation).Scan(&count, &nullSales, &dashVariation, &unknown)
	if err != nil {
		t.Fatal(err)
	}
	if count != 6 || nullSales != 6 || dashVariation != 5 || unknown != 1 {
		t.Fatalf("count/null/dash/unknown = %d/%d/%d/%d", count, nullSales, dashVariation, unknown)
	}
	var distinct int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(DISTINCT (installation_id,provider_listing_id,variation_id)) FROM listings WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, h.installation).Scan(&distinct); err != nil || distinct != 6 {
		t.Fatalf("distinct composite keys = %d, err = %v", distinct, err)
	}

	stub.setPages(map[string]stubListingPage{"": {items: items[:5]}})
	status, body = postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("second status = %d, body = %#v", status, body)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("terminal status = %s", got)
	}
	rows, err := h.pool.Query(context.Background(), `SELECT provider_listing_id,status,price_amount::text,published_quantity FROM listings WHERE tenant_id=$1 AND installation_id=$2 ORDER BY provider_listing_id`, h.tenant, h.installation)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	seen := 0
	wantStatus := map[string]string{
		"MLBTEST0001": "active", "MLBTEST0002": "paused", "MLBTEST0003": "closed",
		"MLBTEST0004": "unknown", "MLBTEST0005": "active", "MLBTEST0006": "closed",
	}
	for rows.Next() {
		var id, state string
		var price *string
		var qty *int
		if err := rows.Scan(&id, &state, &price, &qty); err != nil {
			t.Fatal(err)
		}
		seen++
		if state != wantStatus[id] {
			t.Fatalf("%s status = %s, want %s", id, state, wantStatus[id])
		}
		if id == "MLBTEST0001" && (state != "active" || price == nil || *price != "129.90" || qty == nil || *qty != 7) {
			t.Fatalf("MLBTEST0001 facts changed: %s %#v %#v", state, price, qty)
		}
		if id != "MLBTEST0001" && id != "MLBTEST0003" && price != nil {
			t.Fatalf("%s nullable price changed to %q", id, *price)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if seen != 6 {
		t.Fatalf("rows after close = %d", seen)
	}
}

func TestListingsRefreshRejectsConcurrentRunWithActiveID(t *testing.T) {
	gate := make(chan struct{})
	stub := &stubListingReader{pages: map[string]stubListingPage{"": {items: listingSnapshots()}}, gate: gate}
	h := newRefreshHarness(t, "f01-refresh-concurrent", 100, stub)
	status, first := postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("first status = %d", status)
	}
	id := operationID(t, first)
	status, second := postRefresh(t, h.mux, h.installation)
	if status != http.StatusConflict {
		t.Fatalf("second status = %d, body = %#v", status, second)
	}
	errBody := second["error"].(map[string]any)
	details := errBody["details"].(map[string]any)
	if errBody["code"] != "refresh_in_progress" || details["operation_run_id"] != id {
		t.Fatalf("conflict body = %#v", second)
	}
	var active int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2 AND operation_type='listings_refresh' AND status IN ('queued','running')`, h.tenant, h.installation).Scan(&active); err != nil || active != 1 {
		t.Fatalf("active runs = %d, err = %v", active, err)
	}
	close(gate)
	if got := waitTerminal(t, h.operations, h.installation, id); got != "succeeded" {
		t.Fatalf("terminal status = %s", got)
	}
}

func TestListingsRefreshCapabilityErrorMidPullLeavesRowsUnchanged(t *testing.T) {
	items := listingSnapshots()
	stub := &stubListingReader{pages: map[string]stubListingPage{"": {items: items}}}
	h := newRefreshHarness(t, "f01-refresh-atomic", 100, stub)
	status, body := postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("seed status = %d", status)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("seed terminal = %s", got)
	}
	before := listingStateJSON(t, h.pool, h.tenant, h.installation)
	stub.setPages(map[string]stubListingPage{
		"":  {items: items[:3]},
		"3": {err: connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderAuth, "raw provider body must not persist")},
	})
	// Recompose with page size 3 while retaining the same persisted installation and stub.
	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(h.pool, h.tenant), h.tenant)
	operationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(h.pool, h.tenant), h.tenant)
	caps := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercadolivre", Listings: stub}})
	ingestion := listingsapp.NewIngestion(listingsconnectors.NewSource(caps), listingspostgres.NewRepository(h.pool, h.tenant), 3, time.Now)
	reports := make(chan error, 1)
	svc := listingsapp.NewRefreshService(listingsintegrations.NewGateway(installationSvc, operationSvc), ingestion, func(task func()) { go task() }, time.Now, func(err error) { reports <- err }, func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })
	mux := httpx.NewRouteClassMux()
	listingstransport.NewRefreshHandler(svc).Register(mux)
	status, body = postRefresh(t, mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("failure status = %d", status)
	}
	id := operationID(t, body)
	if got := waitTerminal(t, operationSvc, h.installation, id); got != "failed" {
		t.Fatalf("terminal = %s", got)
	}
	after := listingStateJSON(t, h.pool, h.tenant, h.installation)
	if before != after {
		t.Fatalf("listing rows changed across incomplete pull\nbefore=%s\nafter=%s", before, after)
	}
	var translated string
	var evidence []byte
	if err := h.pool.QueryRow(context.Background(), `SELECT translated_error_code,provider_evidence_json FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2 AND operation_run_id=$3 AND status='failed'`, h.tenant, h.installation, id).Scan(&translated, &evidence); err != nil {
		t.Fatal(err)
	}
	if translated != "CONNECTORS_PROVIDER_AUTH" {
		t.Fatalf("translated_error_code = %q", translated)
	}
	if strings.Contains(string(evidence), "raw provider body") {
		t.Fatalf("provider evidence leaked response: %s", evidence)
	}
	select {
	case err := <-reports:
		t.Fatalf("report callback called for pull failure: %v", err)
	default:
	}
}

func listingStateJSON(t *testing.T, pool *pgxpool.Pool, tenant, installation string) string {
	t.Helper()
	var state []byte
	err := pool.QueryRow(context.Background(), `SELECT COALESCE(jsonb_agg(to_jsonb(x) ORDER BY provider_listing_id,variation_id),'[]'::jsonb) FROM (SELECT provider_listing_id,variation_id,title,status,price_amount::text,price_currency,published_quantity,sync_state,sales_30d,fetched_at FROM listings WHERE tenant_id=$1 AND installation_id=$2) x`, tenant, installation).Scan(&state)
	if err != nil {
		t.Fatal(err)
	}
	return string(state)
}

func TestListingsRefreshIsTenantIsolated(t *testing.T) {
	items := listingSnapshots()
	stub := &stubListingReader{pages: map[string]stubListingPage{"": {items: items}}}
	h := newRefreshHarness(t, "f01-refresh-tenant", 100, stub)
	// integration_installations PK is (installation_id) — globally unique — so tenant B
	// cannot reuse tenant A's installation_id. It gets its own; isolation is proven by
	// tenant A's refresh/close leaving tenant B's tenant-scoped rows untouched.
	other := h.tenant + "-b"
	otherInstall := h.installation + "-b"
	seedConnectedInstallation(t, h.pool, other, otherInstall)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = h.pool.Exec(ctx, `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM integration_installations WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
	})
	otherRepo := listingspostgres.NewRepository(h.pool, other)
	otherCaps := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercadolivre", Listings: &stubListingReader{pages: map[string]stubListingPage{"": {items: items}}}}})
	otherInstallationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(h.pool, other), other)
	otherOperationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(h.pool, other), other)
	otherSvc := listingsapp.NewRefreshService(listingsintegrations.NewGateway(otherInstallationSvc, otherOperationSvc), listingsapp.NewIngestion(listingsconnectors.NewSource(otherCaps), otherRepo, 100, time.Now), func(task func()) { go task() }, time.Now, nil, func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })
	otherMux := httpx.NewRouteClassMux()
	listingstransport.NewRefreshHandler(otherSvc).Register(otherMux)
	status, body := postRefresh(t, otherMux, otherInstall)
	if status != http.StatusAccepted {
		t.Fatalf("tenant B seed status = %d", status)
	}
	if got := waitTerminal(t, otherOperationSvc, otherInstall, operationID(t, body)); got != "succeeded" {
		t.Fatalf("tenant B terminal = %s", got)
	}
	var beforeRuns int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall).Scan(&beforeRuns); err != nil {
		t.Fatal(err)
	}
	stub.setPages(map[string]stubListingPage{"": {items: items[:5]}})
	status, body = postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("tenant A status = %d", status)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("tenant A terminal = %s", got)
	}
	var otherRows, otherClosed, afterRuns int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*),count(*) FILTER (WHERE status='closed') FROM listings WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall).Scan(&otherRows, &otherClosed); err != nil {
		t.Fatal(err)
	}
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall).Scan(&afterRuns); err != nil {
		t.Fatal(err)
	}
	if otherRows != 6 || otherClosed != 1 || afterRuns != beforeRuns {
		t.Fatalf("tenant B rows/closed/runs = %d/%d/%d, runs before=%d", otherRows, otherClosed, afterRuns, beforeRuns)
	}
}

func TestListingsRefreshUnknownInstallation(t *testing.T) {
	stub := &stubListingReader{pages: map[string]stubListingPage{"": {items: listingSnapshots()}}}
	h := newRefreshHarness(t, "f01-refresh-missing", 100, stub)
	missing := h.installation + "-missing"
	status, body := postRefresh(t, h.mux, missing)
	if status != http.StatusNotFound {
		t.Fatalf("status = %d, body = %#v", status, body)
	}
	errBody := body["error"].(map[string]any)
	if errBody["code"] != "installation_not_found" {
		t.Fatalf("body = %#v", body)
	}
	var runs, listings int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, missing).Scan(&runs); err != nil {
		t.Fatal(err)
	}
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM listings WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, missing).Scan(&listings); err != nil {
		t.Fatal(err)
	}
	if runs != 0 || listings != 0 {
		t.Fatalf("missing installation created runs/listings = %d/%d", runs, listings)
	}
}
