//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestObservationLatestPerListingPreservesInputOrder(t *testing.T) {
	ctx, repo, _, _, token := observationHarness(t, "latest")
	older := time.Date(2026, time.July, 15, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	for _, listingID := range []string{"listing-a-" + token, "listing-b-" + token, "listing-c-" + token} {
		if err := repo.AppendObservations(ctx, []domain.MarketObservation{
			storedObservation(listingID, "10.00", older),
			storedObservation(listingID, "20.00", newer),
		}); err != nil {
			t.Fatal(err)
		}
	}

	wantOrder := []string{"listing-c-" + token, "listing-a-" + token, "listing-b-" + token}
	got, err := repo.LatestObservations(ctx, []string{wantOrder[0], "missing-" + token, wantOrder[1], wantOrder[2]})
	if err != nil {
		t.Fatalf("LatestObservations() error = %v", err)
	}
	if len(got) != len(wantOrder) {
		t.Fatalf("LatestObservations() rows = %d, want %d: %+v", len(got), len(wantOrder), got)
	}
	for i, observation := range got {
		if observation.ListingID != wantOrder[i] {
			t.Fatalf("row %d listing_id = %q, want %q", i, observation.ListingID, wantOrder[i])
		}
		if observation.CapturedAt == nil || !observation.CapturedAt.Equal(newer) {
			t.Fatalf("row %d captured_at = %v, want %v", i, observation.CapturedAt, newer)
		}
		if observation.OurSalePrice == nil || observation.OurSalePrice.Amount != "20.00" {
			t.Fatalf("row %d price = %+v, want newest", i, observation.OurSalePrice)
		}
	}
}

func TestObservationRepositoryCannotCrossTenant(t *testing.T) {
	ctx, repoA, pool, tenantA, token := observationHarness(t, "tenant")
	tenantB := "observation-b-" + token
	repoB := marketpostgres.NewRepository(pool, tenantB)
	when := time.Date(2026, time.July, 15, 11, 0, 0, 0, time.UTC)
	listingA, listingB := "listing-a-"+token, "listing-b-"+token
	if err := repoA.AppendObservations(ctx, []domain.MarketObservation{storedObservation(listingA, "11.00", when)}); err != nil {
		t.Fatal(err)
	}
	if err := repoB.AppendObservations(ctx, []domain.MarketObservation{storedObservation(listingB, "22.00", when)}); err != nil {
		t.Fatal(err)
	}

	gotA, err := repoA.LatestObservations(ctx, []string{listingB, listingA})
	if err != nil || len(gotA) != 1 || gotA[0].ListingID != listingA {
		t.Fatalf("tenant A read = %+v, err = %v", gotA, err)
	}
	gotB, err := repoB.LatestObservations(ctx, []string{listingA, listingB})
	if err != nil || len(gotB) != 1 || gotB[0].ListingID != listingB {
		t.Fatalf("tenant B read = %+v, err = %v", gotB, err)
	}
	_ = tenantA
}

func TestObservationWithoutSourceRejectedByDatabase(t *testing.T) {
	ctx, _, pool, tenant, token := observationHarness(t, "source")
	_, err := pool.Exec(ctx, `
		INSERT INTO market_observations (
			tenant_id, listing_id, evidence_state, captured_at, source
		) VALUES ($1, $2, 'observed', $3, $4)
	`, tenant, "listing-"+token, time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC), nil)
	if err == nil {
		t.Fatal("INSERT with null source succeeded, want database rejection")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23502" {
		t.Fatalf("INSERT error = %v, want SQLSTATE 23502", err)
	}
}

func TestObservationSignalsRoundTripSeparately(t *testing.T) {
	ctx, repo, _, _, token := observationHarness(t, "signals")
	when := time.Now().UTC().Truncate(time.Microsecond)
	status := domain.CompetitiveStatusCompeting
	source := domain.SourceVendor
	want := domain.MarketObservation{
		ListingID:         "listing-" + token,
		OurSalePrice:      &domain.Money{Amount: "101.01", Currency: "BRL"},
		CompetitiveStatus: &status,
		WinnerPrice:       &domain.Money{Amount: "202.02", Currency: "BRL"},
		CompetitiveTarget: &domain.Money{Amount: "303.03", Currency: "BRL"},
		CatalogOfferPrice: &domain.Money{Amount: "404.04", Currency: "BRL"},
		CatalogStats: &domain.CatalogStats{
			Min: "1.01", P25: "2.02", Median: "3.03", TrimmedMean: "4.04", P75: "5.05", Max: "6.06", NOffers: 7, NSellers: 6,
		},
		EvidenceState: domain.EvidenceStateObserved,
		CapturedAt:    &when,
		Source:        &source,
	}
	if err := repo.AppendObservations(ctx, []domain.MarketObservation{want}); err != nil {
		t.Fatal(err)
	}
	got, err := repo.LatestObservations(ctx, []string{want.ListingID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !reflect.DeepEqual(got[0], want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}

func observationHarness(t *testing.T, name string) (context.Context, *marketpostgres.Repository, *pgxpool.Pool, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_observation_"+name)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant := "observation-a-" + token
	return ctx, marketpostgres.NewRepository(pool, tenant), pool, tenant, token
}

func storedObservation(listingID, amount string, capturedAt time.Time) domain.MarketObservation {
	source := domain.SourceOfficialAPI
	return domain.MarketObservation{
		ListingID:     listingID,
		OurSalePrice:  &domain.Money{Amount: amount, Currency: "BRL"},
		EvidenceState: domain.EvidenceStateObserved,
		CapturedAt:    &capturedAt,
		Source:        &source,
	}
}
