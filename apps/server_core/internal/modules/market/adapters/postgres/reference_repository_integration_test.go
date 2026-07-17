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

func TestReferenceLatestPerProductPreservesInputOrder(t *testing.T) {
	ctx, repo, _, _, token := referenceHarness(t, "latest")
	older := time.Date(2026, time.July, 16, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	for _, productID := range []string{"product-a-" + token, "product-b-" + token, "product-c-" + token} {
		if err := repo.AppendReferences(ctx, []domain.MarketReference{
			storedReference(productID, domain.MatchStateReview, older),
			storedReference(productID, domain.MatchStateAccept, newer),
		}); err != nil {
			t.Fatal(err)
		}
	}

	wantOrder := []string{"product-c-" + token, "product-a-" + token, "product-b-" + token}
	got, err := repo.LatestReferences(ctx, []string{wantOrder[0], "missing-" + token, wantOrder[1], wantOrder[2]})
	if err != nil {
		t.Fatalf("LatestReferences() error = %v", err)
	}
	if len(got) != len(wantOrder) {
		t.Fatalf("LatestReferences() rows = %d, want %d: %+v", len(got), len(wantOrder), got)
	}
	for i, reference := range got {
		if reference.ProductID != wantOrder[i] {
			t.Fatalf("row %d product_id = %q, want %q", i, reference.ProductID, wantOrder[i])
		}
		if reference.CapturedAt == nil || !reference.CapturedAt.Equal(newer) {
			t.Fatalf("row %d captured_at = %v, want %v", i, reference.CapturedAt, newer)
		}
		if reference.MatchState != domain.MatchStateAccept {
			t.Fatalf("row %d match_state = %q, want newest", i, reference.MatchState)
		}
		want := storedReference(wantOrder[i], domain.MatchStateAccept, newer)
		if !reflect.DeepEqual(reference, want) {
			t.Fatalf("row %d round trip = %#v, want %#v", i, reference, want)
		}
	}
}

func TestReferenceRepositoryCannotCrossTenant(t *testing.T) {
	ctx, repoA, pool, _, token := referenceHarness(t, "tenant")
	repoB := marketpostgres.NewRepository(pool, "reference-b-"+token)
	when := time.Date(2026, time.July, 16, 11, 0, 0, 0, time.UTC)
	productA, productB := "product-a-"+token, "product-b-"+token
	if err := repoA.AppendReferences(ctx, []domain.MarketReference{storedReference(productA, domain.MatchStateAccept, when)}); err != nil {
		t.Fatal(err)
	}
	if err := repoB.AppendReferences(ctx, []domain.MarketReference{storedReference(productB, domain.MatchStateReject, when)}); err != nil {
		t.Fatal(err)
	}

	gotA, err := repoA.LatestReferences(ctx, []string{productB, productA})
	if err != nil || len(gotA) != 1 || gotA[0].ProductID != productA {
		t.Fatalf("tenant A read = %+v, err = %v", gotA, err)
	}
	gotB, err := repoB.LatestReferences(ctx, []string{productA, productB})
	if err != nil || len(gotB) != 1 || gotB[0].ProductID != productB {
		t.Fatalf("tenant B read = %+v, err = %v", gotB, err)
	}
}

func TestReferenceProvenanceIsMandatory(t *testing.T) {
	ctx, _, pool, tenant, token := referenceHarness(t, "source")
	_, err := pool.Exec(ctx, `
		INSERT INTO market_references (
			tenant_id, product_id, match_state, evidence_state, captured_at, source
		) VALUES ($1, $2, 'accept', 'observed', $3, $4)
	`, tenant, "product-"+token, time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC), nil)
	if err == nil {
		t.Fatal("INSERT with null source succeeded, want database rejection")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23502" {
		t.Fatalf("INSERT error = %v, want SQLSTATE 23502", err)
	}
}

func referenceHarness(t *testing.T, name string) (context.Context, *marketpostgres.Repository, *pgxpool.Pool, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_reference_"+name)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant := "reference-a-" + token
	return ctx, marketpostgres.NewRepository(pool, tenant), pool, tenant, token
}

func storedReference(productID string, matchState domain.MatchState, capturedAt time.Time) domain.MarketReference {
	catalogProductID := "catalog-" + productID
	matchMethod := "gtin"
	source := domain.SourceOfficialAPI
	stats := referenceStats()
	return domain.MarketReference{
		ProductID:        productID,
		CatalogProductID: &catalogProductID,
		MatchState:       matchState,
		MatchMethod:      &matchMethod,
		CatalogStats:     &stats,
		EvidenceState:    domain.EvidenceStateObserved,
		CapturedAt:       &capturedAt,
		Source:           &source,
	}
}

func referenceStats() domain.CatalogStats {
	return domain.CatalogStats{
		Min: "1.01", P25: "2.02", Median: "3.03", TrimmedMean: "4.04", P75: "5.05", Max: "6.06", NOffers: 7, NSellers: 6,
	}
}
