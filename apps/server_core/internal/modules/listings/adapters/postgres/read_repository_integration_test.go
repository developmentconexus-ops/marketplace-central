//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

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
