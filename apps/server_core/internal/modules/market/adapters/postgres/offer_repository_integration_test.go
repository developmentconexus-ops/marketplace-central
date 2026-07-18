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

func TestValidatedOffersAreTenantScopedAndPreserveUnknownPrice(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_offers")
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	repoA := marketpostgres.NewRepository(pool, "offer-a-"+token)
	repoB := marketpostgres.NewRepository(pool, "offer-b-"+token)
	now := time.Now().UTC().Truncate(time.Microsecond)
	catalogID := "catalog-" + token
	want := []domain.ValidatedOffer{
		{CatalogProductID: catalogID, SellerID: "seller-priced", Price: &domain.Money{Amount: "88.90", Currency: "BRL"}, Condition: "new", MatchStatus: domain.MatchStatusAccept, FetchedAt: now},
		{CatalogProductID: catalogID, SellerID: "seller-unknown", Price: nil, Condition: "new", MatchStatus: domain.MatchStatusReview, FetchedAt: now.Add(time.Microsecond)},
	}
	if err := repoA.AppendValidatedOffers(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := repoA.ValidatedOffersByCatalogProduct(ctx, catalogID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].SellerID != want[1].SellerID || got[0].Price != nil || got[1].Price == nil || got[1].Price.Amount != "88.90" {
		t.Fatalf("tenant A offers = %+v, want newest first with nil preserved", got)
	}
	other, err := repoB.ValidatedOffersByCatalogProduct(ctx, catalogID)
	if err != nil {
		t.Fatal(err)
	}
	if len(other) != 0 {
		t.Fatalf("tenant B offers = %+v, want none", other)
	}
}

func TestValidatedOffersDuplicateAppendIsNoOp(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_offers_dup")
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant := "offer-dup-" + token
	repo := marketpostgres.NewRepository(pool, tenant)
	now := time.Now().UTC().Truncate(time.Microsecond)
	catalogID := "catalog-" + token
	offers := []domain.ValidatedOffer{
		{CatalogProductID: catalogID, SellerID: "seller-dup", Price: &domain.Money{Amount: "12.00", Currency: "BRL"}, Condition: "new", MatchStatus: domain.MatchStatusAccept, FetchedAt: now},
	}
	if err := repo.AppendValidatedOffers(ctx, offers); err != nil {
		t.Fatal(err)
	}
	if err := repo.AppendValidatedOffers(ctx, offers); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM market_validated_offers WHERE tenant_id = $1 AND catalog_product_id = $2 AND seller_id = $3 AND fetched_at = $4`, tenant, catalogID, "seller-dup", now).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("row count = %d, want 1", count)
	}
}
