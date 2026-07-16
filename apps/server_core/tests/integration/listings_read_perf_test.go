//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	listingsintegrations "marketplace-central/apps/server_core/internal/modules/listings/adapters/integrations"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	listingstransport "marketplace-central/apps/server_core/internal/modules/listings/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

type sqlRecorder struct {
	mu         sync.Mutex
	statements []string
}

func (r *sqlRecorder) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	r.mu.Lock()
	r.statements = append(r.statements, data.SQL)
	r.mu.Unlock()
	return ctx
}
func (*sqlRecorder) TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData) {}
func (r *sqlRecorder) reset()                                                        { r.mu.Lock(); r.statements = nil; r.mu.Unlock() }
func (r *sqlRecorder) summaryCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, s := range r.statements {
		if strings.Contains(s, "count(*) FILTER") && strings.Contains(s, "l.tenant_id=$1") {
			n++
		}
	}
	return n
}

func TestListingsReadPerformance2000(t *testing.T) {
	seedPool, cfg := testpostgres.OpenPool(t, "tenant_harness_f02_read_perf")
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant, installation := "f02-perf-"+token, "f02-perf-inst-"+token
	seedInstallation(t, seedPool, tenant, installation)
	t.Cleanup(func() {
		ctx := context.Background()
		for _, table := range []string{"product_links", "product_link_listing_snapshots", "listings"} {
			_, _ = seedPool.Exec(ctx, "DELETE FROM "+table+" WHERE tenant_id=$1 AND installation_id=$2", tenant, installation)
		}
		_, _ = seedPool.Exec(ctx, `DELETE FROM integration_installations WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	ctx := context.Background()
	batch := &pgx.Batch{}
	one := 1.0
	costs := make(map[int64]*ports.CostFact, 2000)
	for i := 0; i < 2000; i++ {
		id := fmt.Sprintf("MLBPERF%04d", i)
		product := int64(200000 + i)
		costs[product] = &ports.CostFact{Amount: &one}
		batch.Queue(`INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,listing_type_code,status,price_amount,price_currency,sync_state,fetched_at) VALUES($1,$2,'mercadolivre',$3,'-',$4,'gold_pro','active',90,'BRL','synced',now())`, tenant, installation, id, fmt.Sprintf("Title %04d", i))
		batch.Queue(`INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id,internal_product_name) VALUES($1,$2,'mercadolivre',$3,'','resolved',$4,$5)`, tenant, installation, id, product, fmt.Sprintf("Product %04d", i))
		batch.Queue(`INSERT INTO product_link_listing_snapshots(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,seller_sku,fetched_at) VALUES($1,$2,'mercadolivre',$3,'',$4,now())`, tenant, installation, id, "SKU-"+id)
	}
	br := seedPool.SendBatch(ctx, batch)
	for i := 0; i < 6000; i++ {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			t.Fatal(err)
		}
	}
	if err := br.Close(); err != nil {
		t.Fatal(err)
	}
	// Prime planner statistics: a freshly bulk-loaded table has zero stats until
	// autovacuum fires (async, not guaranteed in a short-lived ephemeral DB), so the
	// planner could pick a Seq Scan and fail the keyset-index EXPLAIN assertion for a
	// reason unrelated to the index under test. ANALYZE makes the plan deterministic.
	if _, err := seedPool.Exec(ctx, "ANALYZE listings"); err != nil {
		t.Fatal(err)
	}

	recorder := &sqlRecorder{}
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	poolCfg.ConnConfig.Tracer = recorder
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	installationSvc := integrationsapp.NewInstallationService(integrationspostgres.NewInstallationRepository(pool, tenant), tenant)
	svc := listingsapp.NewReadService(listingspostgres.NewRepository(pool, tenant), deterministicCosts{costs: costs}, testPricingPolicy{}, listingsintegrations.NewInstallationReader(installationSvc), time.Now)
	mux := httpx.NewRouteClassMux()
	listingstransport.NewReadHandler(svc).Register(mux)
	path := "/listings?installation_id=" + installation + "&limit=50"
	call := func() time.Duration {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		start := time.Now()
		mux.ServeHTTP(rr, req)
		elapsed := time.Since(start)
		if rr.Code != 200 {
			t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
		}
		return elapsed
	}
	for i := 0; i < 5; i++ {
		call()
	}
	samples := make([]time.Duration, 100)
	for i := range samples {
		samples[i] = call()
	}
	sorted := append([]time.Duration(nil), samples...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	rank := int(math.Ceil(.95 * float64(len(sorted))))
	p95 := sorted[rank-1]
	t.Logf("environment: go=%s os=%s/%s sha=%s", runtime.Version(), runtime.GOOS, runtime.GOARCH, os.Getenv("GITHUB_SHA"))
	t.Logf("samples=%v", samples)
	t.Logf("nearest-rank p95=%s", p95)
	if p95 >= 500*time.Millisecond {
		t.Fatalf("p95=%s >= 500ms", p95)
	}

	rows, err := seedPool.Query(ctx, `EXPLAIN SELECT l.provider_listing_id FROM listings l WHERE l.tenant_id=$1 AND l.installation_id=$2 AND (l.title,l.provider_listing_id,l.variation_id)>($3,$4,$5) ORDER BY l.title,l.provider_listing_id,l.variation_id LIMIT 51`, tenant, installation, "Title 1000", "MLBPERF1000", "-")
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			rows.Close()
			t.Fatal(err)
		}
		lines = append(lines, line)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	plan := strings.Join(lines, "\n")
	t.Logf("EXPLAIN plan:\n%s", plan)
	if !strings.Contains(plan, "idx_listings_f02_title_key") || !strings.Contains(plan, "Index") || strings.Contains(plan, "Seq Scan") {
		t.Fatalf("keyset index not used:\n%s", plan)
	}

	recorder.reset()
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/listings/summary?installation_id="+installation, nil))
	if rr.Code != 200 {
		t.Fatalf("summary status=%d body=%s", rr.Code, rr.Body.String())
	}
	count := recorder.summaryCount()
	t.Logf("summary conditional-aggregate query count=%d", count)
	if count != 1 {
		t.Fatalf("summary aggregate queries=%d, want 1", count)
	}
}
