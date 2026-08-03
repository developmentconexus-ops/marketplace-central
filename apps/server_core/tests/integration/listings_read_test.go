//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	listingsintegrations "marketplace-central/apps/server_core/internal/modules/listings/adapters/integrations"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	listingstransport "marketplace-central/apps/server_core/internal/modules/listings/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deterministicCosts struct {
	costs map[int64]*ports.CostFact
}

// deterministicCosts is the sole substitution: an explicit fake for the external
// Oracle cost port. It is deterministic contract evidence, never live-Oracle evidence.
func (f deterministicCosts) GetCostFactsByIDs(_ context.Context, ids []int64) (map[int64]*ports.CostFact, error) {
	out := make(map[int64]*ports.CostFact, len(ids))
	for _, id := range ids {
		out[id] = f.costs[id]
	}
	return out, nil
}
func (deterministicCosts) GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*ports.ICMSCeiling, error) {
	rate := 22.0
	return map[int64]*ports.ICMSCeiling{8: {Percent: &rate}}, nil
}

// root.go intentionally composes an unavailable policy source. This is the same
// explicit policy contract used by read_service_test; no new policy contract is invented.
type testPricingPolicy struct{}

func (testPricingPolicy) GetPricingPolicyForInstallation(context.Context, string) (ports.PricingPolicy, bool, error) {
	return ports.PricingPolicy{MinMarginPercent: 10}, true, nil
}

type listingsHarness struct {
	pool                        *pgxpool.Pool
	mux                         http.Handler
	tenant, other, installation string
}

func newListingsHarness(t *testing.T, tag string, costs map[int64]*ports.CostFact) listingsHarness {
	t.Helper()
	pool, _ := testpostgres.OpenPool(t, tag)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	h := listingsHarness{pool: pool, tenant: "f02-a-" + token, other: "f02-b-" + token, installation: "f02-inst-" + token}
	seedInstallation(t, pool, h.tenant, h.installation)
	seedInstallation(t, pool, h.other, h.installation+"-b")
	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(pool, h.tenant), h.tenant)
	svc := listingsapp.NewReadService(listingspostgres.NewRepository(pool, h.tenant), deterministicCosts{costs: costs}, testPricingPolicy{}, listingsintegrations.NewInstallationReader(installationSvc), time.Now)
	mux := httpx.NewRouteClassMux()
	listingstransport.NewReadHandler(svc).Register(mux)
	h.mux = mux
	t.Cleanup(func() {
		ctx := context.Background()
		for _, table := range []string{"listing_sync_events", "product_links", "product_link_listing_snapshots", "listings"} {
			_, _ = pool.Exec(ctx, "DELETE FROM "+table+" WHERE tenant_id = ANY($1)", []string{h.tenant, h.other})
		}
		_, _ = pool.Exec(ctx, `DELETE FROM integration_installations WHERE tenant_id=ANY($1)`, []string{h.tenant, h.other})
	})
	return h
}

func seedInstallation(t *testing.T, pool *pgxpool.Pool, tenant, installation string) {
	t.Helper()
	ctx := context.Background()
	_, err := pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES('mercadolivre','marketplace','Mercado Livre','oauth2','interactive') ON CONFLICT(provider_code) DO NOTHING`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO integration_installations(installation_id,tenant_id,provider_code,family,display_name,status,health_status) VALUES($1,$2,'mercadolivre','marketplace','F02','connected','healthy')`, installation, tenant)
	if err != nil {
		t.Fatal(err)
	}
}

type listingSeed struct {
	id, variation, title, status, sync, listingType, sku, linkState, productName string
	productID                                                                    any
	price                                                                        any
	syncError                                                                    bool
}

func seedListing(t *testing.T, h listingsHarness, tenant, installation string, s listingSeed) {
	t.Helper()
	ctx := context.Background()
	var problem any
	if s.syncError {
		problem = []byte(`{"code":"provider_rejected","message":"rejected"}`)
	}
	var currency any = "BRL"
	if s.price == nil {
		currency = nil
	}
	_, err := h.pool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,listing_type_code,status,price_amount,price_currency,published_quantity,sync_state,sync_error,quality_score,sales_30d,fetched_at) VALUES($1,$2,'mercadolivre',$3,$4,$5,$6,$7,$8,$9,NULL,$10,$11,NULL,NULL,now())`, tenant, installation, s.id, s.variation, s.title, s.listingType, s.status, s.price, currency, s.sync, problem)
	if err != nil {
		t.Fatal(err)
	}
	identity := s.variation
	if identity == "-" {
		identity = ""
	}
	if s.linkState != "" {
		_, err = h.pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id,internal_product_name) VALUES($1,$2,'mercadolivre',$3,$4,$5,$6,$7)`, tenant, installation, s.id, identity, s.linkState, s.productID, s.productName)
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err = h.pool.Exec(ctx, `INSERT INTO product_link_listing_snapshots(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,seller_sku,fetched_at) VALUES($1,$2,'mercadolivre',$3,$4,$5,now())`, tenant, installation, s.id, identity, s.sku)
	if err != nil {
		t.Fatal(err)
	}
}

func requestJSON(t *testing.T, handler http.Handler, path string) (int, map[string]any) {
	t.Helper()
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s: status=%d body=%q: %v", path, rr.Code, rr.Body.String(), err)
	}
	return rr.Code, body
}
func items(body map[string]any) []any {
	if v, ok := body["items"].([]any); ok {
		return v
	}
	return nil
}
func ids(rows []any) []string {
	out := make([]string, 0, len(rows))
	for _, raw := range rows {
		out = append(out, raw.(map[string]any)["provider_listing_id"].(string))
	}
	return out
}

func seedContractRows(t *testing.T, h listingsHarness, costs map[int64]*ports.CostFact) {
	rows := []listingSeed{
		{"MLBTEST0001", "-", "Alpha", "active", "synced", "gold_pro", "SKU-ALPHA", "resolved", "Mesmo", int64(101), 90, false},
		{"MLBTEST0002", "V2", "Beta", "paused", "stale", "gold_special", "SKU-BETA", "resolved", "Mesmo", int64(102), 100, false},
		{"MLBTEST0003", "-", "Gamma", "active", "error", "free", "SKU-ERROR", "conflict", "", nil, nil, true},
		{"MLBTEST0004", "-", "Same", "closed", "synced", "gold_pro", "SKU-SAME-A", "resolved", "Tie", int64(104), 70, false},
		{"MLBTEST0005", "-", "Same", "active", "synced", "gold_pro", "SKU-SAME-B", "resolved", "Tie", int64(105), 70, false},
		{"MLBTEST0006", "-", "Zulu unlinked", "active", "synced", "gold_pro", "SKU-UNLINKED", "", "", nil, 50, false},
	}
	for _, r := range rows {
		seedListing(t, h, h.tenant, h.installation, r)
	}
	// Identical provider identity in Tenant B, with data that must never leak.
	seedListing(t, h, h.other, h.installation, listingSeed{"MLBTEST0001", "-", "Tenant B secret", "active", "synced", "gold_pro", "LEAK", "resolved", "Leak", int64(999), 1, false})
	for i := 0; i < 12; i++ {
		_, err := h.pool.Exec(context.Background(), `INSERT INTO listing_sync_events(tenant_id,installation_id,provider_listing_id,variation_id,event_id,at,kind,message_pt) VALUES($1,$2,'MLBTEST0001','-',$3,$4,'synced',$5)`, h.tenant, h.installation, fmt.Sprintf("event-%02d", i), time.Date(2026, 1, 1, 0, i/2, i%2, 0, time.UTC), fmt.Sprintf("message-%02d", i))
		if err != nil {
			t.Fatal(err)
		}
	}
	_ = costs
}

func TestListingsReadContractEndToEnd(t *testing.T) {
	high, low := 100.0, 1.0
	costs := map[int64]*ports.CostFact{101: {Amount: &high}, 102: nil, 104: {Amount: &low}, 105: {Amount: &low}}
	h := newListingsHarness(t, "tenant_harness_f02_read_contract", costs)
	seedContractRows(t, h, costs)
	base := "?installation_id=" + url.QueryEscape(h.installation)

	t.Run("small page cursor walk and JSON null contract", func(t *testing.T) {
		path := "/listings" + base + "&limit=2"
		var got []any
		for {
			status, body := requestJSON(t, h.mux, path)
			if status != 200 {
				t.Fatalf("status=%d body=%v", status, body)
			}
			for _, k := range []string{"items", "next_cursor", "page_size", "as_of"} {
				if _, ok := body[k]; !ok {
					t.Fatalf("missing envelope key %s", k)
				}
			}
			got = append(got, items(body)...)
			cursor, ok := body["next_cursor"].(string)
			if !ok {
				break
			}
			path = "/listings" + base + "&limit=2&cursor=" + url.QueryEscape(cursor)
		}
		want := []string{"MLBTEST0001", "MLBTEST0002", "MLBTEST0003", "MLBTEST0004", "MLBTEST0005", "MLBTEST0006"}
		if fmt.Sprint(ids(got)) != fmt.Sprint(want) {
			t.Fatalf("ids=%v", ids(got))
		}
		seen := map[string]bool{}
		for _, raw := range got {
			row := raw.(map[string]any)
			id := row["listing_id"].(string)
			if seen[id] {
				t.Fatalf("duplicate %s", id)
			}
			seen[id] = true
			for _, k := range []string{"published_quantity", "sync_error", "cost", "below_margin_worst_case", "fetched_at"} {
				if _, ok := row[k]; !ok {
					t.Fatalf("%s missing %s", id, k)
				}
			}
			// quality_score/sales_30d left the HTTP read contract in Task 8 (ADR-C3):
			// no producer wires them into the response anymore.
			for _, k := range []string{"quality_score", "sales_30d"} {
				if _, ok := row[k]; ok {
					t.Fatalf("%s stale %s still on the wire", id, k)
				}
			}
		}
	})

	t.Run("all filter keys", func(t *testing.T) {
		cases := map[string]string{"status=paused": "MLBTEST0002", "sync_state=error": "MLBTEST0003", "link_state=conflict": "MLBTEST0003", "exception=sync_error": "MLBTEST0003", "listing_type_code=free": "MLBTEST0003", "product_id=101": "MLBTEST0001"}
		for filter, want := range cases {
			status, b := requestJSON(t, h.mux, "/listings"+base+"&filter."+filter)
			if status != 200 || fmt.Sprint(ids(items(b))) != "["+want+"]" {
				t.Fatalf("%s: status=%d ids=%v", filter, status, ids(items(b)))
			}
		}
		status, b := requestJSON(t, h.mux, "/listings"+base+"&filter.has_exception=false")
		if status != 200 || len(items(b)) == 0 {
			t.Fatalf("has_exception status=%d body=%v", status, b)
		}
	})

	t.Run("q title provider id and seller sku", func(t *testing.T) {
		for _, q := range []string{"Alpha", "MLBTEST0002", "SKU-ERROR"} {
			status, b := requestJSON(t, h.mux, "/listings"+base+"&q="+url.QueryEscape(q))
			if status != 200 || len(items(b)) != 1 {
				t.Fatalf("q=%s status=%d ids=%v", q, status, ids(items(b)))
			}
		}
	})

	t.Run("by product cursor walk tie order and null last", func(t *testing.T) {
		path := "/listings/by-product" + base + "&limit=1"
		var groups []any
		for {
			status, b := requestJSON(t, h.mux, path)
			if status != 200 {
				t.Fatalf("status=%d body=%v", status, b)
			}
			groups = append(groups, b["groups"].([]any)...)
			c, ok := b["next_cursor"].(string)
			if !ok {
				break
			}
			path = "/listings/by-product" + base + "&limit=1&cursor=" + url.QueryEscape(c)
		}
		if len(groups) != 5 {
			t.Fatalf("groups=%v", groups)
		}
		last := groups[len(groups)-1].(map[string]any)
		if last["product_id"] != nil {
			t.Fatalf("unlinked not last: %v", groups)
		}
		for _, raw := range groups {
			g := raw.(map[string]any)
			if int(g["listing_count"].(float64)) != len(g["listings"].([]any)) {
				t.Fatalf("count mismatch %v", g)
			}
			if _, ok := g["group_state"]; !ok {
				t.Fatalf("missing group_state")
			}
		}
		tieA, tieB := groups[2].(map[string]any), groups[3].(map[string]any)
		if tieA["product_title"] != "Tie" || tieB["product_title"] != "Tie" || tieA["product_id"] != "104" || tieB["product_id"] != "105" {
			t.Fatalf("equal-title tie order=%v/%v", tieA, tieB)
		}
	})

	t.Run("detail variation timeline and not found", func(t *testing.T) {
		status, b := requestJSON(t, h.mux, "/listings/"+url.PathEscape(h.installation+"~MLBTEST0001~-"))
		if status != 200 {
			t.Fatalf("status=%d body=%v", status, b)
		}
		if b["listing_id"] != h.installation+"~MLBTEST0001~-" {
			t.Fatalf("id=%v", b["listing_id"])
		}
		timeline := b["timeline"].([]any)
		if len(timeline) != 10 {
			t.Fatalf("timeline=%d", len(timeline))
		}
		got := []string{}
		for _, e := range timeline {
			got = append(got, e.(map[string]any)["message_pt"].(string))
		}
		want := []string{"message-11", "message-10", "message-09", "message-08", "message-07", "message-06", "message-05", "message-04", "message-03", "message-02"}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("timeline=%v", got)
		}
		for _, id := range []string{"malformed", h.installation + "~absent~-"} {
			s, x := requestJSON(t, h.mux, "/listings/"+url.PathEscape(id))
			if s != 404 || x["error"].(map[string]any)["code"] != "listing_not_found" {
				t.Fatalf("id=%s status=%d body=%v", id, s, x)
			}
		}
		status, b = requestJSON(t, h.mux, "/listings/"+url.PathEscape(h.installation+"~MLBTEST0002~V2"))
		if status != 200 || b["provider_listing_id"] != "MLBTEST0002" {
			t.Fatalf("variation status=%d body=%v", status, b)
		}
	})

	t.Run("error matrix", func(t *testing.T) {
		cases := []struct {
			path   string
			status int
			code   string
		}{{"/listings", 400, "installation_required"}, {"/listings" + base + "&filter.nope=x", 400, "invalid_filter"}, {"/listings" + base + "&filter.status=nope", 400, "invalid_filter"}, {"/listings" + base + "&cursor=nope", 400, "invalid_cursor"}, {"/listings/absent", 404, "listing_not_found"}}
		for _, tc := range cases {
			s, b := requestJSON(t, h.mux, tc.path)
			code := b["error"].(map[string]any)["code"]
			if s != tc.status || code != tc.code {
				t.Fatalf("%s: status=%d code=%v", tc.path, s, code)
			}
		}
		// TODO(F-01 Slice 4): refresh-conflict belongs to the later write slice.
	})

	t.Run("null cost honesty known margin and summary", func(t *testing.T) {
		_, b := requestJSON(t, h.mux, "/listings"+base+"&filter.product_id=102")
		row := items(b)[0].(map[string]any)
		if row["cost"] != nil || row["below_margin_worst_case"] != nil {
			t.Fatalf("unknown cost defaulted: %v", row)
		}
		_, b = requestJSON(t, h.mux, "/listings"+base+"&filter.product_id=101")
		row = items(b)[0].(map[string]any)
		if row["below_margin_worst_case"] != true {
			t.Fatalf("known margin=%v", row)
		}
		status, b := requestJSON(t, h.mux, "/listings/summary"+base)
		if status != 200 {
			t.Fatalf("summary=%v", b)
		}
		ex := b["exceptions"].(map[string]any)
		if ex["below_margin_worst_case"].(float64) != 1 || ex["margin_unknown"].(float64) != 1 {
			t.Fatalf("exceptions=%v", ex)
		}
	})

	t.Run("tenant isolation all read paths and cursors", func(t *testing.T) {
		paths := []string{"/listings" + base, "/listings" + base + "&q=secret", "/listings" + base + "&filter.product_id=999", "/listings/by-product" + base, "/listings/summary" + base, "/listings/" + url.PathEscape(h.installation+"~MLBTEST0001~-")}
		for _, p := range paths {
			s, b := requestJSON(t, h.mux, p)
			if s != 200 {
				t.Fatalf("%s status=%d body=%v", p, s, b)
			}
			encoded, _ := json.Marshal(b)
			if strings.Contains(string(encoded), "Tenant B secret") || strings.Contains(string(encoded), "LEAK") || strings.Contains(string(encoded), "999") {
				t.Fatalf("tenant leak path=%s body=%s", p, encoded)
			}
		}
		_, b := requestJSON(t, h.mux, "/listings"+base+"&limit=2")
		cursor := b["next_cursor"].(string)
		_, b = requestJSON(t, h.mux, "/listings"+base+"&limit=2&cursor="+url.QueryEscape(cursor))
		for _, id := range ids(items(b)) {
			if strings.Contains(id, "LEAK") {
				t.Fatalf("cursor leak")
			}
		}
	})

}
