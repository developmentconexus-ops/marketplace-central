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

func TestCompetitiveSignalRoundTripOrderIsolationAndIdempotency(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "competitive_signal_"+token)
	tenantA := "signal-a-" + token
	repoA := marketpostgres.NewRepository(pool, tenantA)
	repoB := marketpostgres.NewRepository(pool, "signal-b-"+token)
	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	signals := []domain.CompetitiveSignal{
		{ListingID: "listing-a-" + token, OurPrice: &domain.Money{Amount: "101.2300", Currency: "BRL"}, FetchedAt: fetchedAt},
		{ListingID: "listing-b-" + token, OurPrice: &domain.Money{Amount: "202.4500", Currency: "BRL"}, FetchedAt: fetchedAt},
	}
	if err := repoA.AppendCompetitiveSignals(ctx, signals); err != nil {
		t.Fatal(err)
	}
	if err := repoA.AppendCompetitiveSignals(ctx, signals[:1]); err != nil {
		t.Fatal(err)
	}
	got, err := repoA.LatestCompetitiveSignals(ctx, []string{signals[1].ListingID, "missing-" + token, signals[0].ListingID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ListingID != signals[1].ListingID || got[1].ListingID != signals[0].ListingID {
		t.Fatalf("ordered read = %#v", got)
	}
	for i, signal := range got {
		wantAmount := []string{"202.4500", "101.2300"}[i]
		if signal.OurPrice == nil || signal.OurPrice.Amount != wantAmount || signal.OurPrice.Currency != "BRL" {
			t.Fatalf("row %d our_price = %#v, want %s BRL", i, signal.OurPrice, wantAmount)
		}
		if signal.WinnerPrice != nil || signal.TargetPrice != nil || signal.Position != nil {
			t.Fatalf("row %d nullable fields = %#v/%#v/%#v, want nil", i, signal.WinnerPrice, signal.TargetPrice, signal.Position)
		}
		if !signal.FetchedAt.Equal(fetchedAt) {
			t.Fatalf("row %d fetched_at = %v, want %v", i, signal.FetchedAt, fetchedAt)
		}
	}
	otherTenant, err := repoB.LatestCompetitiveSignals(ctx, []string{signals[0].ListingID})
	if err != nil || len(otherTenant) != 0 {
		t.Fatalf("tenant B read = %#v, err = %v", otherTenant, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM market_competitive_signals WHERE tenant_id = $1 AND listing_id = $2`, tenantA, signals[0].ListingID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate row count = %d, want 1", count)
	}
}

func TestCompetitiveSignalPopulatedNullableRoundTrip(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "competitive_signal_full_"+token)
	repo := marketpostgres.NewRepository(pool, "signal-full-"+token)
	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	want := domain.CompetitiveSignal{
		ListingID:   "listing-" + token,
		OurPrice:    &domain.Money{Amount: "150.0000", Currency: "BRL"},
		WinnerPrice: &domain.Money{Amount: "140.5500", Currency: "BRL"},
		TargetPrice: &domain.Money{Amount: "139.9900", Currency: "BRL"},
		Position:    &domain.MarketPosition{Rank: 3, Total: 8},
		FetchedAt:   fetchedAt,
	}
	if err := repo.AppendCompetitiveSignals(ctx, []domain.CompetitiveSignal{want}); err != nil {
		t.Fatal(err)
	}
	got, err := repo.LatestCompetitiveSignals(ctx, []string{want.ListingID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("rows = %d, want 1", len(got))
	}
	s := got[0]
	// Distinct values per nullable column catch any INSERT/SELECT/Scan column-order swap.
	if s.WinnerPrice == nil || s.WinnerPrice.Amount != "140.5500" || s.TargetPrice == nil || s.TargetPrice.Amount != "139.9900" {
		t.Fatalf("winner/target = %#v/%#v, want 140.5500/139.9900", s.WinnerPrice, s.TargetPrice)
	}
	if s.OurPrice == nil || s.OurPrice.Amount != "150.0000" {
		t.Fatalf("our_price = %#v, want 150.0000", s.OurPrice)
	}
	if s.Position == nil || s.Position.Rank != 3 || s.Position.Total != 8 {
		t.Fatalf("position = %#v, want rank 3 total 8", s.Position)
	}
}
