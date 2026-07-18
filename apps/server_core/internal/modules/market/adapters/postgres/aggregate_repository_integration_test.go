//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestMarketAggregateRoundTripOrderIsolationAndIdempotency(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "market_aggregate_"+token)
	tenantA := "aggregate-a-" + token
	repoA := marketpostgres.NewRepository(pool, tenantA)
	repoB := marketpostgres.NewRepository(pool, "aggregate-b-"+token)
	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	computedAt := fetchedAt.Add(time.Second)
	aggregates := []domain.MarketAggregate{
		{ProductID: "product-a-" + token, NOffers: 2, NSellers: 1, Source: domain.MarketPriceSourceMLSalePrice, FetchedAt: fetchedAt, ComputedAt: computedAt, Status: domain.MarketAggregateStatusInsufficientMarket},
		{ProductID: "product-b-" + token, NOffers: 7, NSellers: 6, Source: domain.MarketPriceSourceMLCatalogOffers, FetchedAt: fetchedAt, ComputedAt: computedAt, Status: domain.MarketAggregateStatusNoPriceEvidence},
	}
	if err := repoA.AppendMarketAggregates(ctx, aggregates); err != nil {
		t.Fatal(err)
	}
	if err := repoA.AppendMarketAggregates(ctx, aggregates[:1]); err != nil {
		t.Fatal(err)
	}
	got, err := repoA.LatestMarketAggregates(ctx, []string{aggregates[1].ProductID, "missing-" + token, aggregates[0].ProductID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ProductID != aggregates[1].ProductID || got[1].ProductID != aggregates[0].ProductID {
		t.Fatalf("ordered read = %#v", got)
	}
	for i, aggregate := range got {
		want := aggregates[1-i]
		if aggregate.Median != nil || aggregate.MinValid != nil {
			t.Fatalf("row %d nullable prices = %#v/%#v, want nil", i, aggregate.Median, aggregate.MinValid)
		}
		if aggregate.NOffers != want.NOffers || aggregate.NSellers != want.NSellers || aggregate.Source != want.Source || aggregate.Status != want.Status || !aggregate.FetchedAt.Equal(fetchedAt) || !aggregate.ComputedAt.Equal(computedAt) {
			t.Fatalf("row %d = %#v, want %#v", i, aggregate, want)
		}
	}
	otherTenant, err := repoB.LatestMarketAggregates(ctx, []string{aggregates[0].ProductID})
	if err != nil || len(otherTenant) != 0 {
		t.Fatalf("tenant B read = %#v, err = %v", otherTenant, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM market_aggregates WHERE tenant_id = $1 AND product_id = $2`, tenantA, aggregates[0].ProductID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate row count = %d, want 1", count)
	}
}

func TestMarketAggregatePopulatedPricesRoundTrip(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "market_aggregate_full_"+token)
	repo := marketpostgres.NewRepository(pool, "aggregate-full-"+token)
	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	want := domain.MarketAggregate{
		ProductID: "product-" + token,
		Median:    &domain.Money{Amount: "88.8800", Currency: "BRL"},
		MinValid:  &domain.Money{Amount: "72.3300", Currency: "BRL"},
		NOffers:   9, NSellers: 6, Source: domain.MarketPriceSourceMLSalePrice,
		FetchedAt: fetchedAt, ComputedAt: fetchedAt.Add(time.Second), Status: domain.MarketAggregateStatusOK,
	}
	if err := repo.AppendMarketAggregates(ctx, []domain.MarketAggregate{want}); err != nil {
		t.Fatal(err)
	}
	got, err := repo.LatestMarketAggregates(ctx, []string{want.ProductID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("rows = %d, want 1", len(got))
	}
	a := got[0]
	// Distinct median vs min_valid catch an amount/currency column-order swap.
	if a.Median == nil || a.Median.Amount != "88.8800" || a.MinValid == nil || a.MinValid.Amount != "72.3300" {
		t.Fatalf("median/min_valid = %#v/%#v, want 88.8800/72.3300", a.Median, a.MinValid)
	}
	if a.NOffers != 9 || a.NSellers != 6 || a.Status != domain.MarketAggregateStatusOK {
		t.Fatalf("counts/status = %d/%d/%q, want 9/6/OK", a.NOffers, a.NSellers, a.Status)
	}
}

func TestMarketAggregateLatestBySourceOrderIsolationAndEmptyInput(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "market_aggregate_source_"+token)
	tenantA := "aggregate-source-a-" + token
	repoA := marketpostgres.NewRepository(pool, tenantA)
	repoB := marketpostgres.NewRepository(pool, "aggregate-source-b-"+token)
	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	productA := "source-product-a-" + token
	productB := "source-product-b-" + token
	missing := "source-product-missing-" + token
	if err := repoA.AppendMarketAggregates(ctx, []domain.MarketAggregate{
		{ProductID: productA, Median: &domain.Money{Amount: "90.00", Currency: "BRL"}, Source: domain.MarketPriceSourceMLSalePrice, FetchedAt: fetchedAt, ComputedAt: fetchedAt.Add(time.Second), Status: domain.MarketAggregateStatusOK},
		{ProductID: productA, Median: &domain.Money{Amount: "80.00", Currency: "BRL"}, Source: domain.MarketPriceSourceMLCatalogOffers, FetchedAt: fetchedAt, ComputedAt: fetchedAt.Add(2 * time.Second), Status: domain.MarketAggregateStatusOK},
		{ProductID: productB, Median: &domain.Money{Amount: "70.00", Currency: "BRL"}, Source: domain.MarketPriceSourceMLCatalogOffers, FetchedAt: fetchedAt, ComputedAt: fetchedAt.Add(3 * time.Second), Status: domain.MarketAggregateStatusOK},
	}); err != nil {
		t.Fatal(err)
	}
	if err := repoB.AppendMarketAggregates(ctx, []domain.MarketAggregate{
		{ProductID: productA, Median: &domain.Money{Amount: "1.00", Currency: "BRL"}, Source: domain.MarketPriceSourceMLCatalogOffers, FetchedAt: fetchedAt, ComputedAt: fetchedAt.Add(4 * time.Second), Status: domain.MarketAggregateStatusOK},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := repoA.LatestMarketAggregatesBySource(ctx, []string{productB, missing, productA}, domain.MarketPriceSourceMLCatalogOffers)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ProductID != productB || got[1].ProductID != productA {
		t.Fatalf("ordered source read = %#v", got)
	}
	if got[0].Median == nil || got[0].Median.Amount != "70.00" || got[1].Median == nil || got[1].Median.Amount != "80.00" {
		t.Fatalf("source-selected prices = %#v, want 70.00 and 80.00", got)
	}

	// Discriminating case: sale_price is NOT productA's overall-latest row (catalog_offers @ +2s is newer
	// than sale_price @ +1s), so a dropped source=$3 predicate would wrongly return the 80.00 catalog row here.
	saleOnly, err := repoA.LatestMarketAggregatesBySource(ctx, []string{productA}, domain.MarketPriceSourceMLSalePrice)
	if err != nil {
		t.Fatal(err)
	}
	if len(saleOnly) != 1 || saleOnly[0].Median == nil || saleOnly[0].Median.Amount != "90.00" {
		t.Fatalf("sale-source read = %#v, want exactly 90.00 (not the newer catalog 80.00)", saleOnly)
	}

	otherTenant, err := repoB.LatestMarketAggregatesBySource(ctx, []string{productA}, domain.MarketPriceSourceMLCatalogOffers)
	if err != nil {
		t.Fatal(err)
	}
	if len(otherTenant) != 1 || otherTenant[0].Median == nil || otherTenant[0].Median.Amount != "1.00" {
		t.Fatalf("tenant B source read = %#v", otherTenant)
	}

	empty, err := repoA.LatestMarketAggregatesBySource(ctx, nil, domain.MarketPriceSourceMLCatalogOffers)
	if err != nil {
		t.Fatal(err)
	}
	if empty == nil || len(empty) != 0 {
		t.Fatalf("empty source read = %#v, want non-nil empty", empty)
	}
}
