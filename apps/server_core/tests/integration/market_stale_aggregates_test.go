//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Proves against real Postgres that the enumeration returns only what went
// past the TTL, oldest first, respecting the ceiling. No fake knows ORDER BY
// or LIMIT, and it is exactly the order that decides what the job collects
// first when the ceiling cuts.
func TestStaleAggregateProductIDsReturnsOldestFirstUnderTheCeiling(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "market_stale_"+token)
	tenant := "stale-" + token
	repo := marketpostgres.NewRepository(pool, tenant)

	now := time.Now().UTC().Truncate(time.Microsecond)
	productOld := "stale-old-" + token
	productMid := "stale-mid-" + token
	productFresh := "stale-fresh-" + token

	computedOld := now.Add(-72 * time.Hour)
	computedMid := now.Add(-2 * time.Hour)
	computedFresh := now.Add(-5 * time.Minute)

	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM market_aggregates WHERE tenant_id = $1`, tenant); err != nil {
			t.Errorf("cleanup stale aggregate fixtures: %v", err)
		}
	})

	aggregates := []domain.MarketAggregate{
		{ProductID: productOld, NOffers: 1, NSellers: 1, Source: domain.MarketPriceSourceMLSalePrice, FetchedAt: computedOld, ComputedAt: computedOld, Status: domain.MarketAggregateStatusOK},
		{ProductID: productMid, NOffers: 1, NSellers: 1, Source: domain.MarketPriceSourceMLSalePrice, FetchedAt: computedMid, ComputedAt: computedMid, Status: domain.MarketAggregateStatusOK},
		{ProductID: productFresh, NOffers: 1, NSellers: 1, Source: domain.MarketPriceSourceMLSalePrice, FetchedAt: computedFresh, ComputedAt: computedFresh, Status: domain.MarketAggregateStatusOK},
	}
	if err := repo.AppendMarketAggregates(ctx, aggregates); err != nil {
		t.Fatal(err)
	}

	// 1. TTL 1h, high ceiling -> the 3-days and 2-hours products come back,
	// oldest first, and the 5-minutes product (negative control) never leaks in.
	stale, err := repo.StaleAggregateProductIDs(ctx, time.Hour, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(stale) != 2 || stale[0] != productOld || stale[1] != productMid {
		t.Fatalf("stale under high ceiling = %#v, want [%s %s]", stale, productOld, productMid)
	}
	for _, id := range stale {
		if id == productFresh {
			t.Fatalf("stale under high ceiling leaked the fresh (5-minute) product: %#v", stale)
		}
	}

	// 2. Same data, ceiling 1 -> only the oldest (3-days) product comes back.
	capped, err := repo.StaleAggregateProductIDs(ctx, time.Hour, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(capped) != 1 || capped[0] != productOld {
		t.Fatalf("stale under ceiling=1 = %#v, want [%s]", capped, productOld)
	}

	// 3. TTL greater than all ages -> empty, non-nil slice, not an error.
	none, err := repo.StaleAggregateProductIDs(ctx, 200*time.Hour, 10)
	if err != nil {
		t.Fatal(err)
	}
	if none == nil || len(none) != 0 {
		t.Fatalf("stale under TTL greater than all ages = %#v, want non-nil empty", none)
	}
}
