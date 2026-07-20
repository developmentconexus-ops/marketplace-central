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

	dashboardapp "marketplace-central/apps/server_core/internal/modules/dashboard/application"
	dashboardtransport "marketplace-central/apps/server_core/internal/modules/dashboard/transport"
	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
	erpapp "marketplace-central/apps/server_core/internal/modules/erp_import/application"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationstransport "marketplace-central/apps/server_core/internal/modules/integrations/transport"
	listingsintegrations "marketplace-central/apps/server_core/internal/modules/listings/adapters/integrations"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	orderstransport "marketplace-central/apps/server_core/internal/modules/orders/transport"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

const aggregateProviderMarker = "F01_RAW_PROVIDER_SECRET_MARKER"

type aggregateSyncHarness struct {
	pool          *pgxpool.Pool
	handler       http.Handler
	tenantA       string
	tenantB       string
	installationA string
	installationB string
	now           time.Time
}

type aggregateCosts struct{}

func (aggregateCosts) GetCostFactsByIDs(_ context.Context, ids []int64) (map[int64]*listingsports.CostFact, error) {
	amount := 10.0
	result := make(map[int64]*listingsports.CostFact, len(ids))
	for _, id := range ids {
		result[id] = &listingsports.CostFact{Amount: &amount}
	}
	return result, nil
}

func (aggregateCosts) GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*listingsports.ICMSCeiling, error) {
	percent := 0.0
	return map[int64]*listingsports.ICMSCeiling{13: {Percent: &percent}}, nil
}

type aggregatePolicy struct{}

func (aggregatePolicy) GetPricingPolicyForInstallation(context.Context, string) (listingsports.PricingPolicy, bool, error) {
	return listingsports.PricingPolicy{MinMarginPercent: 10}, true, nil
}

func TestAggregateSyncReadContract(t *testing.T) {
	h := newAggregateSyncHarness(t)
	seedAggregateSyncFacts(t, h)
	base := "?installation_id=" + url.QueryEscape(h.installationA)

	t.Run("orders cursor walk preserves tie boundary and legacy limit", func(t *testing.T) {
		path := "/orders" + base + "&limit=1"
		var got []map[string]any
		for {
			status, body, _ := aggregateGET(t, h, path)
			if status != http.StatusOK {
				t.Fatalf("status=%d body=%v", status, body)
			}
			got = append(got, aggregateItems(body)...)
			cursor, ok := body["next_cursor"].(string)
			if !ok {
				break
			}
			path = "/orders" + base + "&limit=1&cursor=" + url.QueryEscape(cursor)
		}
		if gotIDs := aggregateIDs(got); fmt.Sprint(gotIDs) != "[order-tie-b order-tie-a order-mid order-week order-old]" {
			t.Fatalf("order ids=%v", gotIDs)
		}
		seen := map[string]bool{}
		for _, row := range got {
			id := row["provider_order_id"].(string)
			if seen[id] {
				t.Fatalf("duplicate order %s", id)
			}
			seen[id] = true
		}
		if got[0]["nf_state"] != nil || got[1]["nf_state"] != "linked" {
			t.Fatalf("nf states=%v,%v; unknown evidence must stay null and exact evidence must link", got[0]["nf_state"], got[1]["nf_state"])
		}

		status, body, _ := aggregateGET(t, h, "/orders"+base+"&limit=2")
		if status != http.StatusOK || len(aggregateItems(body)) != 2 {
			t.Fatalf("legacy limit-only request status=%d body=%v", status, body)
		}
	})

	t.Run("orders filters and tenant isolation", func(t *testing.T) {
		status, body, _ := aggregateGET(t, h, "/orders"+base+"&status=cancelled")
		if status != http.StatusOK || fmt.Sprint(aggregateIDs(aggregateItems(body))) != "[order-mid]" {
			t.Fatalf("status filter status=%d body=%v", status, body)
		}

		status, body, _ = aggregateGET(t, h, "/orders"+base+"&date_from=2026-07-16T00:00:00Z&date_to=2026-07-16T23:59:59Z")
		if status != http.StatusOK || fmt.Sprint(aggregateIDs(aggregateItems(body))) != "[order-tie-b order-tie-a order-mid]" {
			t.Fatalf("date filter status=%d body=%v", status, body)
		}

		status, body, _ = aggregateGET(t, h, "/orders"+base+"&q=needle")
		if status != http.StatusOK || fmt.Sprint(aggregateIDs(aggregateItems(body))) != "[order-tie-a]" {
			t.Fatalf("q filter status=%d body=%v", status, body)
		}

		status, body, _ = aggregateGET(t, h, "/orders"+base+"&q=tenant-b-secret")
		if status != http.StatusOK || len(aggregateItems(body)) != 0 || strings.Contains(string(mustJSON(body)), h.tenantB) {
			t.Fatalf("tenant B leaked through orders list status=%d body=%v", status, body)
		}

		status, body, _ = aggregateGET(t, h, "/orders/order-missing?installation_id="+url.QueryEscape(h.installationA))
		if status != http.StatusNotFound || aggregateErrorCode(body) != "order_not_found" {
			t.Fatalf("unknown order status=%d body=%v", status, body)
		}
		status, body, _ = aggregateGET(t, h, "/orders/order-b-secret?installation_id="+url.QueryEscape(h.installationA))
		if status != http.StatusNotFound || aggregateErrorCode(body) != "order_not_found" {
			t.Fatalf("cross-tenant order status=%d body=%v", status, body)
		}
	})

	t.Run("order detail returns canonical facts only", func(t *testing.T) {
		status, body, _ := aggregateGET(t, h, "/orders/order-tie-a?installation_id="+url.QueryEscape(h.installationA))
		if status != http.StatusOK {
			t.Fatalf("status=%d body=%v", status, body)
		}
		if body["nf_state"] != "linked" || len(body["items"].([]any)) != 1 || len(body["payments"].([]any)) != 2 {
			t.Fatalf("detail=%v", body)
		}
		if total, ok := body["total"].(float64); !ok || total != 45 {
			t.Fatalf("total=%v, want payment sum 45", body["total"])
		}
		for _, key := range []string{"buyer_nickname", "currency", "fulfillment"} {
			if body[key] != nil {
				t.Fatalf("unsupported fact %s=%v", key, body[key])
			}
		}
	})

	t.Run("sync runs cursor window null and filters", func(t *testing.T) {
		path := "/sync/runs" + base + "&limit=1"
		var got []map[string]any
		var rawPages []string
		for {
			status, body, raw := aggregateGET(t, h, path)
			if status != http.StatusOK {
				t.Fatalf("status=%d body=%v", status, body)
			}
			got = append(got, aggregateItems(body)...)
			rawPages = append(rawPages, raw)
			cursor, ok := body["next_cursor"].(string)
			if !ok {
				break
			}
			path = "/sync/runs" + base + "&limit=1&cursor=" + url.QueryEscape(cursor)
		}
		if gotIDs := aggregateIDsFromKey(got, "operation_run_id"); fmt.Sprint(gotIDs) != "[run-orders-tie run-orders run-listings run-linkage]" {
			t.Fatalf("run ids=%v", gotIDs)
		}
		for _, raw := range rawPages {
			if strings.Contains(raw, `"operation_run_id":"run-orders"`) && !strings.Contains(raw, `"finished_at":null`) {
				t.Fatalf("running run did not emit raw JSON null: %s", raw)
			}
		}
		for _, row := range got {
			if row["operation_run_id"] == "run-old" || row["installation_id"] == h.installationB {
				t.Fatalf("old or tenant B run leaked: %v", row)
			}
		}

		status, body, _ := aggregateGET(t, h, "/sync/runs"+base+"&filter.module=pricing_fee_sync&filter.status=failed")
		if status != http.StatusOK || fmt.Sprint(aggregateIDsFromKey(aggregateItems(body), "operation_run_id")) != "[run-linkage]" {
			t.Fatalf("module/status filter status=%d body=%v", status, body)
		}
		status, body, _ = aggregateGET(t, h, "/sync/runs"+base+"&filter.module=tenant-b-secret")
		if status != http.StatusOK || len(aggregateItems(body)) != 0 {
			t.Fatalf("tenant B module leaked status=%d body=%v", status, body)
		}
	})

	t.Run("dashboard summary composes real tenant scoped counts", func(t *testing.T) {
		status, body, _ := aggregateGET(t, h, "/dashboard/summary"+base)
		if status != http.StatusOK {
			t.Fatalf("status=%d body=%v", status, body)
		}
		for key, want := range map[string]float64{
			"orders_today":  3,
			"orders_7d":     4,
			"pending_links": 2,
			"missing_gtin":  2,
			"sync_errors":   1,
		} {
			if got, ok := body[key].(float64); !ok || got != want {
				t.Fatalf("%s=%v, want %v; body=%v", key, body[key], want, body)
			}
		}
		lastSync, ok := body["last_sync_at"].(map[string]any)
		if !ok {
			t.Fatalf("last_sync_at=%v", body["last_sync_at"])
		}
		for _, module := range []string{"orders_read", "listings_refresh", "pricing_fee_sync"} {
			if value, ok := lastSync[module].(string); !ok || value == "" {
				t.Fatalf("last_sync_at[%s]=%v", module, lastSync[module])
			}
		}
		if strings.Contains(string(mustJSON(body)), h.tenantB) || strings.Contains(string(mustJSON(body)), "TENANT_B_SECRET") {
			t.Fatalf("tenant B leaked through dashboard summary: %s", mustJSON(body))
		}
	})
}

func newAggregateSyncHarness(t *testing.T) *aggregateSyncHarness {
	t.Helper()
	pool, cfg := testpostgres.OpenPool(t, "aggregate_sync_contract")
	token := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	h := &aggregateSyncHarness{
		pool:          pool,
		tenantA:       "f01-i1-a-" + token,
		tenantB:       "f01-i1-b-" + token,
		installationA: "f01-i1-inst-a-" + token,
		installationB: "f01-i1-inst-b-" + token,
		now:           time.Date(2026, 7, 16, 15, 0, 0, 0, time.UTC),
	}
	seedAggregateInstallations(t, h, cfg.DefaultTenantID)

	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(pool, h.tenantA), h.tenantA)
	operationSvc := integrationsapp.NewOperationService(integrationspostgres.NewOperationRunRepository(pool, h.tenantA), h.tenantA)
	listingSvc := listingsapp.NewReadService(
		listingspostgres.NewRepository(pool, h.tenantA),
		aggregateCosts{},
		aggregatePolicy{},
		listingsintegrations.NewInstallationReader(installationSvc),
		func() time.Time { return h.now },
	)
	dashboardSvc := dashboardapp.NewService(
		installationSvc,
		listingSvc,
		productlinksapp.NewSummaryService(productlinkspostgres.NewSummaryReader(pool, h.tenantA)),
		ordersapp.NewSummaryService(orderspostgres.NewOrderRepository(pool, h.tenantA)),
		operationSvc,
		erpapp.NewQueryService(erppostgres.NewRepository(pool), h.tenantA),
		func() time.Time { return h.now },
	)

	mux := httpx.NewRouteClassMux()
	orderstransport.NewHandlerWithReader(nil, ordersapp.NewReadService(orderspostgres.NewOrderReadRepository(pool, h.tenantA))).Register(mux)
	integrationstransport.NewRunReadHandler(operationSvc).Register(mux)
	dashboardtransport.NewHandler(dashboardSvc).Register(mux)
	h.handler = mux

	// The Sankhya linkage ledger is append-only (DELETE is rejected by a trigger), and unique-per-run tenant IDs isolate these rows, so no cleanup is performed.
	return h
}

func seedAggregateInstallations(t *testing.T, h *aggregateSyncHarness, _ string) {
	t.Helper()
	ctx := context.Background()
	if _, err := h.pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES('mercado_livre','marketplace','Mercado Livre','none','manual') ON CONFLICT(provider_code) DO NOTHING`); err != nil {
		t.Fatal(err)
	}
	for tenant, installation := range map[string]string{h.tenantA: h.installationA, h.tenantB: h.installationB} {
		if _, err := h.pool.Exec(ctx, `INSERT INTO integration_installations(installation_id,tenant_id,provider_code,family,display_name,status,health_status) VALUES($1,$2,'mercado_livre','marketplace',$3,'connected','healthy')`, installation, tenant, tenant); err != nil {
			t.Fatal(err)
		}
	}
}

func seedAggregateSyncFacts(t *testing.T, h *aggregateSyncHarness) {
	t.Helper()
	ctx := context.Background()
	today := h.now
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-tie-b", "paid", timePtr(today.Add(-2*time.Hour)), false)
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-tie-a", "paid", timePtr(today.Add(-2*time.Hour)), true)
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-mid", "cancelled", timePtr(today.Add(-3*time.Hour)), false)
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-week", "paid", timePtr(today.AddDate(0, 0, -6)), false)
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-old", "paid", timePtr(today.AddDate(0, 0, -8)), false)
	seedAggregateOrder(t, h, h.tenantA, h.installationA, "order-unknown", "paid", nil, false)
	seedAggregateOrder(t, h, h.tenantB, h.installationB, "order-b-secret", "paid", timePtr(today.Add(-time.Hour)), false)

	if _, err := h.pool.Exec(ctx, `INSERT INTO orders_marketplace_order_items(tenant_id,installation_id,provider_order_id,line_no,provider_item_id,provider_variation_id,seller_sku,title,quantity,unit_price,sale_fee_amount,link_quality,internal_product_id,created_at,updated_at,mpc_line_id,reconciliation_state) VALUES($1,$2,'order-tie-a',1,'item-a','','SKU-NEEDLE','Needle product',2,20,1,'resolved',101,$3,$3,'mpl_11111111111111111111111111111111','stable')`, h.tenantA, h.installationA, h.now); err != nil {
		t.Fatal(err)
	}
	for _, payment := range []struct {
		line        int
		transaction float64
		paid        *float64
	}{{1, 40, float64Ptr(42)}, {2, 3, nil}} {
		if _, err := h.pool.Exec(ctx, `INSERT INTO orders_marketplace_order_payments(tenant_id,installation_id,provider_order_id,line_no,provider_payment_id,provider_status,transaction_amount,total_paid_amount,created_at,updated_at) VALUES($1,$2,'order-tie-a',$3,$4,'approved',$5,$6,$7,$7)`, h.tenantA, h.installationA, payment.line, fmt.Sprintf("payment-%d", payment.line), payment.transaction, payment.paid, h.now); err != nil {
			t.Fatal(err)
		}
	}
	seedAggregateLinkageEvent(t, h, "order-tie-a", "event-exact", "exact")
	seedAggregateLinkageEvent(t, h, "order-tie-b", "event-unknown", "unknown")

	for _, row := range []struct {
		item      string
		ean       string
		state     string
		productID any
	}{{"listing-pending", "", "unresolved", nil}, {"listing-conflict", "789", "conflict", nil}, {"listing-linked", "", "resolved", 101}} {
		if _, err := h.pool.Exec(ctx, `INSERT INTO product_link_listing_snapshots(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,seller_sku,ean,title,fetched_at) VALUES($1,$2,'mercado_livre',$3,'',$4,$5,$3,$6)`, h.tenantA, h.installationA, row.item, "SKU-"+row.item, row.ean, h.now); err != nil {
			t.Fatal(err)
		}
		if _, err := h.pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id) VALUES($1,$2,'mercado_livre',$3,'',$4,$5)`, h.tenantA, h.installationA, row.item, row.state, row.productID); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := h.pool.Exec(ctx, `INSERT INTO product_link_listing_snapshots(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,seller_sku,ean,title,fetched_at) VALUES($1,$2,'mercado_livre','tenant-b-listing','',$3,'999','TENANT_B_SECRET',$4)`, h.tenantB, h.installationB, "SKU-B", h.now); err != nil {
		t.Fatal(err)
	}

	seedAggregateListing(t, h, h.tenantA, h.installationA, "listing-error", "error", true)
	seedAggregateListing(t, h, h.tenantA, h.installationA, "listing-pending", "stale", false)
	seedAggregateListing(t, h, h.tenantA, h.installationA, "listing-conflict", "synced", false)
	seedAggregateListing(t, h, h.tenantA, h.installationA, "listing-linked", "synced", false)
	seedAggregateListing(t, h, h.tenantB, h.installationB, "tenant-b-listing", "synced", false)

	seedAggregateRun(t, h, h.tenantA, h.installationA, "run-orders", "orders_read", "running", h.now.Add(-time.Hour), nil)
	tieStarted := h.now.Add(-time.Hour)
	seedAggregateRun(t, h, h.tenantA, h.installationA, "run-orders-tie", "orders_read", "succeeded", tieStarted, timePtr(tieStarted.Add(5*time.Minute)))
	seedAggregateRun(t, h, h.tenantA, h.installationA, "run-listings", "listings_refresh", "succeeded", h.now.Add(-2*time.Hour), timePtr(h.now.Add(-90*time.Minute)))
	seedAggregateRun(t, h, h.tenantA, h.installationA, "run-linkage", "pricing_fee_sync", "failed", h.now.Add(-3*time.Hour), timePtr(h.now.Add(-170*time.Minute)))
	old := time.Now().UTC().Add(-91 * 24 * time.Hour)
	seedAggregateRun(t, h, h.tenantA, h.installationA, "run-old", "old_module", "succeeded", old, timePtr(old.Add(time.Minute)))
	seedAggregateRun(t, h, h.tenantB, h.installationB, "run-b-secret", "tenant-b-secret", "succeeded", time.Now().UTC().Add(-time.Hour), timePtr(time.Now().UTC()))
}

func seedAggregateOrder(t *testing.T, h *aggregateSyncHarness, tenant, installation, id, status string, createdAt *time.Time, raw bool) {
	t.Helper()
	fetchedAt := h.now
	if createdAt != nil {
		fetchedAt = *createdAt
	}
	var rawRef any
	if raw {
		rawRef = fmt.Sprintf(`{"marker":%q}`, aggregateProviderMarker)
	}
	if _, err := h.pool.Exec(context.Background(), `INSERT INTO orders_marketplace_orders(tenant_id,installation_id,provider_code,provider_order_id,provider_status,provider_created_at,fetched_at,raw_provider_ref,created_at,updated_at) VALUES($1,$2,'mercado_livre',$3,$4,$5,$6,$7::jsonb,$6,$6)`, tenant, installation, id, status, createdAt, fetchedAt, rawRef); err != nil {
		t.Fatal(err)
	}
}

func seedAggregateLinkageEvent(t *testing.T, h *aggregateSyncHarness, orderID, eventID, evidence string) {
	t.Helper()
	_, err := h.pool.Exec(context.Background(), `INSERT INTO orders_sankhya_linkage_events(tenant_id,installation_id,event_id,event_type,provider_order_id,external_order_key,internal_document_id,actor_type,actor_id,reason,idempotency_key,source_at,recorded_at,configuration_revision,evidence_state) VALUES($1,$2,$3,'confirmed',$4,'ml:v1:'||$2||':'||$4,$5,'operator','integration-test','seed',$3,now(),now(),'f01-i1',$6)`, h.tenantA, h.installationA, eventID, orderID, int64(len(eventID)+1000), evidence)
	if err != nil {
		t.Fatal(err)
	}
}

func seedAggregateListing(t *testing.T, h *aggregateSyncHarness, tenant, installation, item, syncState string, syncError bool) {
	t.Helper()
	var problem any
	if syncError {
		problem = `{"code":"provider_rejected"}`
	}
	if _, err := h.pool.Exec(context.Background(), `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,listing_type_code,status,price_amount,price_currency,published_quantity,sync_state,sync_error,fetched_at) VALUES($1,$2,'mercado_livre',$3,'-',$3,'gold','active',100,'BRL',1,$4,$5::jsonb,$6)`, tenant, installation, item, syncState, problem, h.now); err != nil {
		t.Fatal(err)
	}
}

func seedAggregateRun(t *testing.T, h *aggregateSyncHarness, tenant, installation, id, module, status string, startedAt time.Time, finishedAt *time.Time) {
	t.Helper()
	if _, err := h.pool.Exec(context.Background(), `INSERT INTO integration_operation_runs(tenant_id,operation_run_id,installation_id,operation_type,status,result_code,failure_code,attempt_count,duration_ms,started_at,completed_at) VALUES($1,$2,$3,$4,$5,'ok','provider_failed',1,250,$6,$7)`, tenant, id, installation, module, status, startedAt, finishedAt); err != nil {
		t.Fatal(err)
	}
}

func aggregateGET(t *testing.T, h *aggregateSyncHarness, path string) (int, map[string]any, string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	h.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	raw := recorder.Body.String()
	if strings.Contains(raw, aggregateProviderMarker) {
		t.Fatalf("provider payload leaked on %s: %s", path, raw)
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s: status=%d body=%q: %v", path, recorder.Code, raw, err)
	}
	return recorder.Code, body, raw
}

func aggregateItems(body map[string]any) []map[string]any {
	raw, ok := body["items"].([]any)
	if !ok {
		return nil
	}
	items := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		items = append(items, item.(map[string]any))
	}
	return items
}

func aggregateIDs(items []map[string]any) []string {
	return aggregateIDsFromKey(items, "provider_order_id")
}

func aggregateIDsFromKey(items []map[string]any, key string) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item[key].(string))
	}
	return ids
}

func aggregateErrorCode(body map[string]any) string {
	err, ok := body["error"].(map[string]any)
	if !ok {
		return ""
	}
	code, _ := err["code"].(string)
	return code
}

func mustJSON(value any) []byte {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return encoded
}

func float64Ptr(value float64) *float64 { return &value }

func timePtr(value time.Time) *time.Time { return &value }
