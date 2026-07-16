//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksapplication "marketplace-central/apps/server_core/internal/modules/product_links/application"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestLinkageSummaryCountsUnresolvedAndEmptyEAN(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_product_links_summary")
	token := time.Now().UTC().Format("150405.000000000")
	tenant := "summary-a-" + token
	installation := "summary-inst-" + token
	t.Cleanup(func() {
		for _, table := range []string{"product_link_candidates", "product_links", "product_link_listing_snapshots"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id=$1 AND installation_id=$2", tenant, installation)
		}
	})

	for _, row := range []struct {
		itemID string
		state  string
	}{{"unresolved-link", "unresolved"}, {"conflict-link", "conflict"}, {"resolved-link", "resolved"}} {
		_, err := pool.Exec(ctx, `
			INSERT INTO product_links(tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id, state)
			VALUES($1, $2, 'mercadolivre', $3, '', $4)
		`, tenant, installation, row.itemID, row.state)
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO product_link_candidates(
			tenant_id, candidate_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			state, match_input
		) VALUES
			($1, 'candidate-unresolved', $2, 'mercadolivre', 'candidate-link', '', 'unresolved', 'none'),
			($1, 'candidate-conflict', $2, 'mercadolivre', 'candidate-link', '', 'conflict', 'none'),
			($1, 'candidate-exact', $2, 'mercadolivre', 'ignored-link', '', 'exact_sku', 'seller_sku')
	`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		itemID string
		ean    string
	}{{"missing-one", ""}, {"missing-two", "   "}, {"known", "7890123456789"}} {
		_, err := pool.Exec(ctx, `
			INSERT INTO product_link_listing_snapshots(
				tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id, ean, fetched_at
			) VALUES($1, $2, 'mercadolivre', $3, '', $4, now())
		`, tenant, installation, row.itemID, row.ean)
		if err != nil {
			t.Fatal(err)
		}
	}

	service := productlinksapplication.NewSummaryService(productlinkspostgres.NewSummaryReader(pool, tenant))
	got, err := service.Summary(ctx, installation)
	if err != nil {
		t.Fatal(err)
	}
	if got.PendingLinks != 3 || got.MissingGTIN != 2 {
		t.Fatalf("summary=%+v, want pending_links=3 missing_gtin=2", got)
	}
}

func TestLinkageSummaryScopesTenantAndInstallation(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_product_links_scope")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, otherTenant := "summary-scope-a-"+token, "summary-scope-b-"+token
	installation, otherInstallation := "summary-scope-inst-"+token, "summary-scope-other-"+token
	t.Cleanup(func() {
		for _, table := range []string{"product_link_candidates", "product_links", "product_link_listing_snapshots"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1)", []string{tenant, otherTenant})
		}
	})

	seed := func(owner, install, item, ean string) {
		t.Helper()
		if _, err := pool.Exec(ctx, `
			INSERT INTO product_links(tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id, state)
			VALUES($1, $2, 'mercadolivre', $3, '', 'unresolved')
		`, owner, install, item); err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, `
			INSERT INTO product_link_listing_snapshots(
				tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id, ean, fetched_at
			) VALUES($1, $2, 'mercadolivre', $3, '', $4, now())
		`, owner, install, item, ean); err != nil {
			t.Fatal(err)
		}
	}
	seed(tenant, installation, "owned", "")
	seed(tenant, otherInstallation, "other-installation", "")
	seed(otherTenant, installation, "other-tenant", "")

	service := productlinksapplication.NewSummaryService(productlinkspostgres.NewSummaryReader(pool, tenant))
	got, err := service.Summary(ctx, installation)
	if err != nil {
		t.Fatal(err)
	}
	if got.PendingLinks != 1 || got.MissingGTIN != 1 {
		t.Fatalf("summary=%+v, want tenant and installation scoped counts of 1/1", got)
	}
	got, err = service.Summary(ctx, otherInstallation)
	if err != nil {
		t.Fatal(err)
	}
	if got.PendingLinks != 1 || got.MissingGTIN != 1 {
		t.Fatalf("other installation summary=%+v, want 1/1", got)
	}
}
