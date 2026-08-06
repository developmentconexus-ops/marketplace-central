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

	mercadolivre "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
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

// F-04 retires the old page-based stubListingReader{ListListings/ReadListing}
// + listingsapp.NewIngestion(listingsconnectors.NewSource(caps), ...) wiring
// this file used to exercise. The refresh HTTP path now runs through the
// SAME BackfillRunner (ids-only enumeration + multiget hydration) F-03's
// scheduled job uses (ADR-024: one writer, no coexisting path), so this stub
// implements the two optional capabilities that path type-asserts for:
// idScanReader (ListListingsScanIDs) and multigetReader (GetItemsMultiget).
type stubScanPage struct {
	ids  []string
	next string
	err  error
}

type stubListingReader struct {
	mu       sync.Mutex
	pages    map[string]stubScanPage
	items    map[string]mercadolivre.ItemMultigetDTO
	gate     <-chan struct{}
	gateOnce sync.Once
}

// ListListings/ReadListing satisfy the base connectorsports.ListingReader
// capability this stub is registered under; neither is exercised by the
// backfill/refresh path anymore, so both are harmless no-ops.
func (*stubListingReader) ListListings(context.Context, connectorsdomain.ListListingsInput) ([]connectorsdomain.ListingSnapshot, error) {
	return nil, nil
}

func (*stubListingReader) ReadListing(context.Context, connectorsdomain.ProviderListingRef) (connectorsdomain.ListingSnapshot, error) {
	return connectorsdomain.ListingSnapshot{}, nil
}

func (s *stubListingReader) ListListingsScanIDs(_ context.Context, input connectorsdomain.ListListingsInput) ([]string, string, error) {
	s.gateOnce.Do(func() {
		if s.gate != nil {
			<-s.gate
		}
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	page := s.pages[input.Cursor]
	return append([]string(nil), page.ids...), page.next, page.err
}

func (s *stubListingReader) GetItemsMultiget(_ context.Context, _ connectorsdomain.ProviderAccountRef, ids []string) ([]mercadolivre.ItemMultigetDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]mercadolivre.ItemMultigetDTO, 0, len(ids))
	for _, id := range ids {
		item, ok := s.items[id]
		if !ok {
			item = mercadolivre.ItemMultigetDTO{ProviderItemID: id, Err: fmt.Errorf("test fixture: unknown item %s", id)}
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *stubListingReader) setPages(pages map[string]stubScanPage) {
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

func newRefreshHarness(t *testing.T, tag string, stub *stubListingReader) refreshHarness {
	t.Helper()
	pool, cfg := testpostgres.OpenPool(t, tag)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant, installation := cfg.DefaultTenantID+"-f01-"+token, "f01-inst-"+token
	seedConnectedInstallation(t, pool, tenant, installation)
	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(pool, tenant), tenant)
	operationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(pool, tenant), tenant)
	caps := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercadolivre", Listings: stub}})
	listingRepo := listingspostgres.NewRepository(pool, tenant)
	runner := listingsapp.NewBackfillRunner(
		listingsconnectors.NewBackfillSource(caps),
		listingsconnectors.NewMultigetHydrator(caps, nil, time.Now),
		listingRepo, time.Now,
	)
	gateway := listingsintegrations.NewGateway(installationSvc, operationSvc)
	reports := make(chan error, 1)
	refreshSvc := listingsapp.NewRefreshService(gateway, runner, func(task func()) { go task() }, time.Now,
		func(err error) { reports <- err },
		func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })
	mux := httpx.NewRouteClassMux()
	listingstransport.NewRefreshHandler(refreshSvc).Register(mux)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(ctx, `DELETE FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
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

// listingsFixtureIDs is the deterministic id order every test below scans in.
var listingsFixtureIDs = []string{"MLBTEST0001", "MLBTEST0002", "MLBTEST0003", "MLBTEST0004", "MLBTEST0005", "MLBTEST0006"}

// listingsFixtureItems is the multiget DTO fixture set. Status is kept
// VERBATIM by the new mapper (MapMultigetItemToListing, multiget_mapper.go —
// F-02/F-03, not something F-04 changes): MLBTEST0004's out-of-vocabulary
// "provider-new-state" is persisted exactly as given, unlike the OLD
// snapshot-based mapper this file used to exercise, which canonicalized any
// unrecognized status to "unknown". MLBTEST0005 carries one variation
// (V-5): the mapper emits ONE row per item regardless (VariationID
// sentinel "-"), with variation facts (price/qty/seller_sku) attached as a
// CHILD listing_variations row, not a second top-level listings row — a
// real cardinality/shape difference from the old flattened-per-variation
// mapper this file's assertions must respect.
func listingsFixtureItems() map[string]mercadolivre.ItemMultigetDTO {
	qty := 7
	price := "129.90"
	return map[string]mercadolivre.ItemMultigetDTO{
		"MLBTEST0001": {ProviderItemID: "MLBTEST0001", Title: "Alpha", Status: "active", AvailableQuantity: &qty},
		"MLBTEST0002": {ProviderItemID: "MLBTEST0002", Title: "Beta", Status: "paused"},
		"MLBTEST0003": {ProviderItemID: "MLBTEST0003", Title: "Gamma", Status: "closed"},
		"MLBTEST0004": {ProviderItemID: "MLBTEST0004", Title: "Delta", Status: "provider-new-state"},
		"MLBTEST0005": {ProviderItemID: "MLBTEST0005", Title: "Epsilon", Status: "active", Variations: []mercadolivre.ItemMultigetVariationDTO{
			{ProviderVariationID: "V-5", Price: &price, AvailableQuantity: &qty, SellerSKU: "SKU-5"},
		}},
		"MLBTEST0006": {ProviderItemID: "MLBTEST0006", Title: "Zeta", Status: "paused"},
	}
}

func TestListingsRefreshSeedsIC02RowsAndClosesMissing(t *testing.T) {
	stub := &stubListingReader{
		pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}},
		items: listingsFixtureItems(),
	}
	h := newRefreshHarness(t, "f01-refresh-complete", stub)
	status, body := postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("status = %d, body = %#v", status, body)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("terminal status = %s", got)
	}
	var count, nullSales, dashVariation, newStateCount int
	err := h.pool.QueryRow(context.Background(), `SELECT count(*),count(*) FILTER (WHERE sales_30d IS NULL),count(*) FILTER (WHERE variation_id='-'),count(*) FILTER (WHERE status='provider-new-state') FROM listings WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, h.installation).Scan(&count, &nullSales, &dashVariation, &newStateCount)
	if err != nil {
		t.Fatal(err)
	}
	// Every item maps to exactly ONE parent row (variation_id sentinel "-"),
	// even MLBTEST0005 which carries a variation — that variation lives in
	// listing_variations, not as a second top-level row.
	if count != 6 || nullSales != 6 || dashVariation != 6 || newStateCount != 1 {
		t.Fatalf("count/null/dash/newState = %d/%d/%d/%d", count, nullSales, dashVariation, newStateCount)
	}
	var distinct int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(DISTINCT (installation_id,provider_listing_id,variation_id)) FROM listings WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, h.installation).Scan(&distinct); err != nil || distinct != 6 {
		t.Fatalf("distinct composite keys = %d, err = %v", distinct, err)
	}
	var variationCount int
	var variationPrice, variationSKU string
	var variationQty int
	if err := h.pool.QueryRow(context.Background(), `SELECT count(*) FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2`, h.tenant, h.installation).Scan(&variationCount); err != nil || variationCount != 1 {
		t.Fatalf("listing_variations rows = %d, err = %v", variationCount, err)
	}
	if err := h.pool.QueryRow(context.Background(), `SELECT price::text,seller_sku,available_quantity FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id='MLBTEST0005' AND variation_id='V-5'`, h.tenant, h.installation).Scan(&variationPrice, &variationSKU, &variationQty); err != nil {
		t.Fatal(err)
	}
	if variationPrice != "129.90" || variationSKU != "SKU-5" || variationQty != 7 {
		t.Fatalf("MLBTEST0005/V-5 variation = price=%q sku=%q qty=%d", variationPrice, variationSKU, variationQty)
	}

	stub.setPages(map[string]stubScanPage{"": {ids: listingsFixtureIDs[:5]}})
	status, body = postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("second status = %d, body = %#v", status, body)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("terminal status = %s", got)
	}
	// MASS-CLOSURE retired (F-02): a row the second page never mentions keeps its
	// own status (keep-absent semantics) and instead gets absent_since stamped —
	// it is no longer force-flipped to "closed". MLBTEST0006 ("Zeta") is dropped
	// from the second page's ids[:5] and must stay "paused" with absent_since set.
	rows, err := h.pool.Query(context.Background(), `SELECT provider_listing_id,status,absent_since FROM listings WHERE tenant_id=$1 AND installation_id=$2 ORDER BY provider_listing_id`, h.tenant, h.installation)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	seen := 0
	wantStatus := map[string]string{
		"MLBTEST0001": "active", "MLBTEST0002": "paused", "MLBTEST0003": "closed",
		"MLBTEST0004": "provider-new-state", "MLBTEST0005": "active", "MLBTEST0006": "paused",
	}
	for rows.Next() {
		var id, state string
		var absentSince *time.Time
		if err := rows.Scan(&id, &state, &absentSince); err != nil {
			t.Fatal(err)
		}
		seen++
		if state != wantStatus[id] {
			t.Fatalf("%s status = %s, want %s", id, state, wantStatus[id])
		}
		if id == "MLBTEST0006" && absentSince == nil {
			t.Fatalf("MLBTEST0006 dropped from page but absent_since not stamped (mass-closure regression: status should stay, absent_since should mark it)")
		}
		if id != "MLBTEST0006" && absentSince != nil {
			t.Fatalf("%s unexpectedly marked absent_since = %v", id, *absentSince)
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
	stub := &stubListingReader{pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}}, items: listingsFixtureItems(), gate: gate}
	h := newRefreshHarness(t, "f01-refresh-concurrent", stub)
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

// TestListingsRefreshMidPullFailureLeavesEarlierPagesPersistedButLaterPagesUntouched
// replaces the old TestListingsRefreshCapabilityErrorMidPullLeavesRowsUnchanged.
//
// OPEN RISK CALLED OUT (see F-04 final report): the old Ingestion
// accumulated a whole run in memory and upserted once at the end, so a
// mid-pull failure left EVERY row unchanged (before == after, byte for
// byte). BackfillRunner instead persists per page via the SAME ingestBatch
// helper NewListingsJob uses (F-03's already-landed CompletedPullStore split
// into UpsertPulledRows + MarkRunComplete) — a page that succeeds before a
// LATER page fails DOES get persisted. This test proves the new, honest
// contract instead: rows from the successful first page are refreshed
// (their fetched_at moves forward), rows from the page that was never
// reached are byte-for-byte unchanged, nothing from the failing page's raw
// error leaks into evidence, and the run is reported failed with the
// correct translated_error_code.
func TestListingsRefreshMidPullFailureLeavesEarlierPagesPersistedButLaterPagesUntouched(t *testing.T) {
	stub := &stubListingReader{pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}}, items: listingsFixtureItems()}
	h := newRefreshHarness(t, "f01-refresh-atomic", stub)
	status, body := postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("seed status = %d", status)
	}
	if got := waitTerminal(t, h.operations, h.installation, operationID(t, body)); got != "succeeded" {
		t.Fatalf("seed terminal = %s", got)
	}
	firstPageIDs, laterPageIDs := listingsFixtureIDs[:3], listingsFixtureIDs[3:]
	beforeFirst := listingRowsState(t, h.pool, h.tenant, h.installation, firstPageIDs)
	beforeLater := listingRowsState(t, h.pool, h.tenant, h.installation, laterPageIDs)

	stub.setPages(map[string]stubScanPage{
		"":      {ids: firstPageIDs, next: "page2"},
		"page2": {err: connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderAuth, "raw provider body must not persist")},
	})
	operationSvc := h.operations
	status, body = postRefresh(t, h.mux, h.installation)
	if status != http.StatusAccepted {
		t.Fatalf("failure status = %d", status)
	}
	id := operationID(t, body)
	if got := waitTerminal(t, operationSvc, h.installation, id); got != "failed" {
		t.Fatalf("terminal = %s", got)
	}

	afterFirst := listingRowsState(t, h.pool, h.tenant, h.installation, firstPageIDs)
	afterLater := listingRowsState(t, h.pool, h.tenant, h.installation, laterPageIDs)
	if afterFirst == beforeFirst {
		t.Fatalf("expected the first (successful) page's rows to be refreshed (new fetched_at), got byte-identical state\n%s", afterFirst)
	}
	if afterLater != beforeLater {
		t.Fatalf("expected the never-reached page's rows to stay byte-identical\nbefore=%s\nafter=%s", beforeLater, afterLater)
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
	case err := <-h.reports:
		t.Fatalf("report callback called for pull failure: %v", err)
	default:
	}
}

// listingRowsState returns a stable JSON snapshot of exactly the given
// provider_listing_ids' rows, ordered deterministically, for before/after
// comparison.
func listingRowsState(t *testing.T, pool *pgxpool.Pool, tenant, installation string, providerListingIDs []string) string {
	t.Helper()
	var state []byte
	err := pool.QueryRow(context.Background(), `SELECT COALESCE(jsonb_agg(to_jsonb(x) ORDER BY provider_listing_id,variation_id),'[]'::jsonb) FROM (SELECT provider_listing_id,variation_id,title,status,available_quantity,sync_state,sales_30d,fetched_at,absent_since FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id = ANY($3)) x`, tenant, installation, providerListingIDs).Scan(&state)
	if err != nil {
		t.Fatal(err)
	}
	return string(state)
}

func TestListingsRefreshIsTenantIsolated(t *testing.T) {
	stub := &stubListingReader{pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}}, items: listingsFixtureItems()}
	h := newRefreshHarness(t, "f01-refresh-tenant", stub)
	// integration_installations PK is (installation_id) — globally unique — so tenant B
	// cannot reuse tenant A's installation_id. It gets its own; isolation is proven by
	// tenant A's refresh/close leaving tenant B's tenant-scoped rows untouched.
	other := h.tenant + "-b"
	otherInstall := h.installation + "-b"
	seedConnectedInstallation(t, h.pool, other, otherInstall)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = h.pool.Exec(ctx, `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
		_, _ = h.pool.Exec(ctx, `DELETE FROM integration_installations WHERE tenant_id=$1 AND installation_id=$2`, other, otherInstall)
	})
	otherRepo := listingspostgres.NewRepository(h.pool, other)
	otherStub := &stubListingReader{pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}}, items: listingsFixtureItems()}
	otherCaps := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercadolivre", Listings: otherStub}})
	otherInstallationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(h.pool, other), other)
	otherOperationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(h.pool, other), other)
	otherRunner := listingsapp.NewBackfillRunner(listingsconnectors.NewBackfillSource(otherCaps), listingsconnectors.NewMultigetHydrator(otherCaps, nil, time.Now), otherRepo, time.Now)
	otherSvc := listingsapp.NewRefreshService(listingsintegrations.NewGateway(otherInstallationSvc, otherOperationSvc), otherRunner, func(task func()) { go task() }, time.Now, nil, func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })
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
	stub.setPages(map[string]stubScanPage{"": {ids: listingsFixtureIDs[:5]}})
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
	stub := &stubListingReader{pages: map[string]stubScanPage{"": {ids: listingsFixtureIDs}}, items: listingsFixtureItems()}
	h := newRefreshHarness(t, "f01-refresh-missing", stub)
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
