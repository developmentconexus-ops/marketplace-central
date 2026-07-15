//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

type queryCounter struct{ count atomic.Int64 }

func (q *queryCounter) TraceQueryStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData) context.Context {
	q.count.Add(1)
	return ctx
}
func (*queryCounter) TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData) {}

func TestGetListingsSummaryUsesOneTenantScopedAggregate(t *testing.T) {
	ctx := context.Background()
	seedPool, cfg := testpostgres.OpenPool(t, "tenant_harness_listings_summary")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, other, installation := "summary-a-"+token, "summary-b-"+token, "summary-inst-"+token
	t.Cleanup(func() {
		for _, table := range []string{"product_links", "listings"} {
			_, _ = seedPool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1) AND installation_id=$2", []string{tenant, other}, installation)
		}
	})
	seed := func(owner, id, status, sync string, syncError bool, price any) {
		var problem any
		if syncError {
			problem = []byte(`{"code":"x","message":"x"}`)
		}
		var currency any = "BRL"
		if price == nil {
			currency = nil // listings_price_pair_check: amount+currency both-null or both-set
		}
		if _, err := seedPool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,price_amount,price_currency,sync_state,sync_error) VALUES($1,$2,'mercadolivre',$3,'-',$3,$4,$5,$6,$7,$8)`, owner, installation, id, status, price, currency, sync, problem); err != nil {
			t.Fatal(err)
		}
	}
	for _, owner := range []string{tenant, other} {
		seed(owner, "active", "active", "synced", false, 90)
		seed(owner, "paused", "paused", "stale", false, 100)
		seed(owner, "closed", "closed", "error", true, nil)
		for _, row := range []struct {
			id      string
			product int64
		}{{"active", 101}, {"paused", 102}, {"closed", 103}} {
			if _, err := seedPool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id) VALUES($1,$2,'mercadolivre',$3,'','resolved',$4)`, owner, installation, row.id, row.product); err != nil {
				t.Fatal(err)
			}
		}
		seed(owner, "unlinked", "active", "synced", false, 90)
	}
	counter := &queryCounter{}
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	poolConfig.ConnConfig.Tracer = counter
	countedPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(countedPool.Close)
	got, err := listingspostgres.NewRepository(countedPool, tenant).GetListingsSummary(ctx, ports.SummaryQuery{InstallationID: installation})
	if err != nil {
		t.Fatal(err)
	}
	if counter.count.Load() != 1 {
		t.Fatalf("queries=%d", counter.count.Load())
	}
	if got.Total != 4 || got.Active != 2 || got.Paused != 1 || got.SyncError != 1 || got.Stale != 1 || got.Unlinked != 1 || len(got.Linked) != 3 {
		t.Fatalf("summary=%+v", got)
	}
	for _, linked := range got.Linked {
		if linked.CostID == 103 && linked.Price != nil {
			t.Fatalf("null price defaulted: %+v", linked)
		}
	}
	low, high := 65.6597, 10.0
	service := application.NewReadService(
		listingspostgres.NewRepository(seedPool, tenant),
		integrationFacts{costs: map[int64]*ports.CostFact{101: {Amount: &low}, 102: {Amount: &high}, 103: nil}},
		integrationPolicy{}, integrationInstallation{}, time.Now,
	)
	summary, err := service.Summary(ctx, ports.SummaryQuery{InstallationID: installation})
	if err != nil {
		t.Fatal(err)
	}
	if summary.BelowMarginWorstCase == nil || *summary.BelowMarginWorstCase != 1 || summary.MarginUnknown == nil || *summary.MarginUnknown != 1 {
		t.Fatalf("margin counters=%+v", summary)
	}
}

func TestGetListingRowAndTimelineUseCompositeTenantKey(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings_get")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, other, installation := "get-a-"+token, "get-b-"+token, "get-inst-"+token
	t.Cleanup(func() {
		for _, table := range []string{"listing_sync_events", "product_links", "product_link_listing_snapshots", "listings"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1) AND installation_id=$2", []string{tenant, other}, installation)
		}
	})
	for _, row := range []struct{ id, variation, title string }{{"42664", "-", "Base"}, {"42664", "v1", "Variation"}} {
		_, err := pool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,price_amount,price_currency,sync_state) VALUES($1,$2,'mercadolivre',$3,$4,$5,'active',90,'BRL','synced')`, tenant, installation, row.id, row.variation, row.title)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 1; i <= 12; i++ {
		_, err := pool.Exec(ctx, `INSERT INTO listing_sync_events(tenant_id,installation_id,provider_listing_id,variation_id,event_id,at,kind,message_pt) VALUES($1,$2,'42664','-', $3, $4, 'synced', $3)`, tenant, installation, fmt.Sprintf("e%02d", i), time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err := pool.Exec(ctx, `INSERT INTO listing_sync_events(tenant_id,installation_id,provider_listing_id,variation_id,event_id,at,kind,message_pt) VALUES($1,$2,'42664','-', 'e99', $3, 'synced', 'other tenant')`, other, installation, time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	repo := listingspostgres.NewRepository(pool, tenant)
	page, err := repo.ListListingRows(ctx, ports.ListingQuery{InstallationID: installation, Limit: 10})
	if err != nil || len(page.Items) != 2 {
		t.Fatalf("list items=%+v err=%v", page.Items, err)
	}
	for _, id := range []domain.ListingID{{InstallationID: installation, ProviderListingID: "42664", VariationID: "-"}, {InstallationID: installation, ProviderListingID: "42664", VariationID: "v1"}} {
		got, found, err := repo.GetListingRow(ctx, domain.ListingKey{TenantID: other, InstallationID: id.InstallationID, ProviderListingID: id.ProviderListingID, VariationID: id.VariationID})
		if err != nil || !found {
			t.Fatalf("get %s found=%v err=%v", id.String(), found, err)
		}
		var listed domain.ListingReadModel
		for _, item := range page.Items {
			if item.ListingID == id.String() {
				listed = item
			}
		}
		if !reflect.DeepEqual(got, listed) {
			t.Fatalf("get=%+v list=%+v", got, listed)
		}
	}
	_, found, err := repo.GetListingRow(ctx, domain.ListingKey{InstallationID: installation, ProviderListingID: "unknown", VariationID: "-"})
	if err != nil || found {
		t.Fatalf("unknown found=%v err=%v", found, err)
	}
	events, err := repo.ListListingTimeline(ctx, domain.ListingKey{TenantID: other, InstallationID: installation, ProviderListingID: "42664", VariationID: "-"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 10 || events[0].MessagePT != "e12" || events[9].MessagePT != "e03" {
		t.Fatalf("events=%+v", events)
	}
	empty, err := repo.ListListingTimeline(ctx, domain.ListingKey{InstallationID: installation, ProviderListingID: "42664", VariationID: "v1"}, 10)
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("empty=%#v err=%v", empty, err)
	}
}

type integrationFacts struct{ costs map[int64]*ports.CostFact }

func (f integrationFacts) GetCostFactsByIDs(_ context.Context, ids []int64) (map[int64]*ports.CostFact, error) {
	out := map[int64]*ports.CostFact{}
	for _, id := range ids {
		out[id] = f.costs[id]
	}
	return out, nil
}
func (integrationFacts) GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*ports.ICMSCeiling, error) {
	rate := 22.0
	return map[int64]*ports.ICMSCeiling{8: {Percent: &rate}}, nil
}

type integrationPolicy struct{}
type integrationInstallation struct{}

func (integrationInstallation) InstallationExists(context.Context, string) (bool, error) {
	return true, nil
}

func (integrationPolicy) GetPricingPolicyForInstallation(context.Context, string) (ports.PricingPolicy, bool, error) {
	return ports.PricingPolicy{MinMarginPercent: 10}, true, nil
}

func TestReadRepositoryKeysetTenantFiltersAndLiteralQuery(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings_read")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, other, installation := "read-a-"+token, "read-b-"+token, "read-inst-"+token
	t.Cleanup(func() {
		for _, table := range []string{"product_links", "product_link_listing_snapshots", "listings"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1) AND installation_id=$2", []string{tenant, other}, installation)
		}
	})
	seedRead := func(owner, id, title, status, sync string) {
		_, err := pool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,listing_type_code,status,price_amount,price_currency,sync_state,fetched_at) VALUES($1,$2,'mercadolivre',$3,'-',$4,'gold_pro',$5,90,'BRL',$6,now())`, owner, installation, id, title, status, sync)
		if err != nil {
			t.Fatal(err)
		}
	}
	seedRead(tenant, "a", "100% literal", "active", "synced")
	seedRead(tenant, "b", "100_ literal", "paused", "stale")
	seedRead(tenant, "c", "Zulu", "active", "error")
	seedRead(other, "leak", "100% literal", "active", "synced")
	_, err := pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id) VALUES($1,$2,'mercadolivre','a','','resolved',42664),($1,$2,'mercadolivre','b','','conflict',NULL)`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO product_link_listing_snapshots(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,seller_sku,fetched_at) VALUES($1,$2,'mercadolivre','a','','SKU_100%',now())`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	repo := listingspostgres.NewRepository(pool, tenant)
	q := ports.ListingQuery{InstallationID: installation, Limit: 2}
	first, err := repo.ListListingRows(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Items) != 2 || first.NextCursor == nil {
		t.Fatalf("first=%+v", first)
	}
	q.Cursor = *first.NextCursor
	last, err := repo.ListListingRows(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	if len(last.Items) != 1 || last.NextCursor != nil || last.Items[0].ProviderListingID != "c" {
		t.Fatalf("last=%+v", last)
	}
	for _, tc := range []struct {
		name   string
		filter domain.ListingFilter
		want   string
	}{{"status", domain.ListingFilter{Status: domain.ListingStatusPaused}, "b"}, {"sync", domain.ListingFilter{SyncState: domain.ListingSyncStateError}, "c"}, {"link", domain.ListingFilter{LinkState: domain.LinkStateConflict}, "b"}, {"type", domain.ListingFilter{ListingTypeCode: "gold_pro"}, "a"}, {"product", domain.ListingFilter{ProductID: "42664"}, "a"}, {"exception", domain.ListingFilter{Exception: domain.ListingExceptionStale}, "b"}} {
		t.Run(tc.name, func(t *testing.T) {
			p, e := repo.ListListingRows(ctx, ports.ListingQuery{InstallationID: installation, Limit: 1, Filter: tc.filter})
			if e != nil {
				t.Fatal(e)
			}
			if len(p.Items) != 1 || p.Items[0].ProviderListingID != tc.want {
				t.Fatalf("items=%+v", p.Items)
			}
		})
	}
	for _, literal := range []string{"100%", "100_", "SKU_100%"} {
		p, e := repo.ListListingRows(ctx, ports.ListingQuery{InstallationID: installation, Limit: 10, Q: literal})
		if e != nil {
			t.Fatal(e)
		}
		if len(p.Items) != 1 {
			t.Fatalf("q=%q items=%+v", literal, p.Items)
		}
	}
}

func TestReadServiceBelowMarginScanCapAndCursorResumption(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings_scan")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "scan-"+token, "scan-inst-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})
	costs := map[int64]*ports.CostFact{}
	one := 1.0
	for n := 1; n <= 101; n++ {
		id := fmt.Sprintf("scan-%03d", n)
		productID := int64(100000 + n)
		costs[productID] = &ports.CostFact{Amount: &one}
		_, err := pool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,price_amount,price_currency,sync_state) VALUES($1,$2,'mercadolivre',$3,'-',$3,'active',90,'BRL','synced')`, tenant, installation, id)
		if err != nil {
			t.Fatal(err)
		}
		_, err = pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id) VALUES($1,$2,'mercadolivre',$3,'','resolved',$4)`, tenant, installation, id, productID)
		if err != nil {
			t.Fatal(err)
		}
	}
	service := application.NewReadService(listingspostgres.NewRepository(pool, tenant), integrationFacts{costs: costs}, integrationPolicy{}, integrationInstallation{}, time.Now)
	query := ports.ListingQuery{InstallationID: installation, Limit: 2, Filter: domain.ListingFilter{Exception: domain.ListingExceptionBelowMargin}}
	first, err := service.List(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Items) != 0 || first.NextCursor == nil {
		t.Fatalf("cap page=%+v", first)
	}
	query.Cursor = *first.NextCursor
	resumed, err := service.List(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if len(resumed.Items) != 0 || resumed.NextCursor != nil {
		t.Fatalf("resumed page=%+v", resumed)
	}
}

func TestByProductWalkUsesProductIDTieBreakAndOneNullLastBucket(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings_groups")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, other, installation := "groups-a-"+token, "groups-b-"+token, "groups-inst-"+token
	t.Cleanup(func() {
		for _, table := range []string{"product_links", "listings"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1) AND installation_id=$2", []string{tenant, other}, installation)
		}
	})
	seed := func(owner, id string) {
		if _, err := pool.Exec(ctx, `INSERT INTO listings(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,price_amount,price_currency,sync_state) VALUES($1,$2,'mercadolivre',$3,'-',$3,'active',90,'BRL','synced')`, owner, installation, id); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"a1", "a2", "b1", "conflict", "rejected", "absent"} {
		seed(tenant, id)
	}
	seed(other, "leak")
	for _, row := range []struct {
		id, state string
		product   any
	}{{"a1", "resolved", 101}, {"a2", "resolved", 101}, {"b1", "resolved", 102}, {"conflict", "conflict", nil}, {"rejected", "rejected", nil}} {
		name := ""
		if row.state == "resolved" {
			name = "Mesmo"
		}
		if _, err := pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id,internal_product_name) VALUES($1,$2,'mercadolivre',$3,'',$4,$5,$6)`, tenant, installation, row.id, row.state, row.product, name); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `INSERT INTO product_links(tenant_id,installation_id,provider_code,provider_item_id,provider_variation_id,state,internal_product_id,internal_product_name) VALUES($1,$2,'mercadolivre','leak','','resolved',1,'Aardvark')`, other, installation); err != nil {
		t.Fatal(err)
	}
	one := 1.0
	service := application.NewReadService(listingspostgres.NewRepository(pool, tenant), integrationFacts{costs: map[int64]*ports.CostFact{101: {Amount: &one}, 102: {Amount: &one}}}, integrationPolicy{}, integrationInstallation{}, time.Now)
	query := ports.ListingGroupQuery{InstallationID: installation, Limit: 1}
	var got []domain.ListingGroup
	for {
		page, err := service.ByProduct(ctx, query)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, page.Groups...)
		if page.NextCursor == nil {
			break
		}
		query.Cursor = *page.NextCursor
	}
	if len(got) != 3 || got[0].ProductID == nil || *got[0].ProductID != "101" || got[1].ProductID == nil || *got[1].ProductID != "102" || got[2].ProductID != nil || got[2].ProductTitle == nil || *got[2].ProductTitle != "sem produto" {
		t.Fatalf("groups=%+v", got)
	}
	if len(got[0].Listings) != 2 || len(got[2].Listings) != 3 {
		t.Fatalf("children=%d/%d groups=%+v", len(got[0].Listings), len(got[2].Listings), got)
	}
	for _, group := range got {
		if group.ListingCount != len(group.Listings) {
			t.Fatalf("count invariant group=%+v", group)
		}
	}
}
