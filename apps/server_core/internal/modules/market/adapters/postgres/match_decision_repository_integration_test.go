//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestMatchDecisionRoundTripIsolationAndIdempotency(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "match_decision_"+token)
	repoA := marketpostgres.NewRepository(pool, "decision-a-"+token)
	repoB := marketpostgres.NewRepository(pool, "decision-b-"+token)
	decidedAt := time.Now().UTC().Truncate(time.Microsecond)
	want := domain.MatchDecision{
		ProductID: "product-" + token, Status: domain.MatchStatusReview,
		Anchors:        []string{"ean:789000000001", "brand:model"},
		Contradictions: []string{"color:red!=blue"},
		ConfidenceBand: domain.ConfidenceBandMedia, DecidedAt: decidedAt,
	}

	if err := repoA.AppendMatchDecision(ctx, want); err != nil {
		t.Fatal(err)
	}
	if err := repoA.AppendMatchDecision(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := repoA.LatestMatchDecisions(ctx, []string{want.ProductID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("round trip rows = %d, want 1: %#v", len(got), got)
	}
	d := got[0]
	if d.ProductID != want.ProductID || d.Status != want.Status || d.ConfidenceBand != want.ConfidenceBand ||
		!reflect.DeepEqual(d.Anchors, want.Anchors) || !reflect.DeepEqual(d.Contradictions, want.Contradictions) ||
		!d.DecidedAt.Equal(want.DecidedAt) {
		t.Fatalf("round trip = %#v, want %#v", d, want)
	}
	if d.CatalogProductID != nil {
		t.Fatalf("catalog_product_id = %q, want nil", *d.CatalogProductID)
	}
	otherTenant, err := repoB.LatestMatchDecisions(ctx, []string{want.ProductID})
	if err != nil || len(otherTenant) != 0 {
		t.Fatalf("tenant B read = %#v, err = %v", otherTenant, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM market_match_decisions WHERE tenant_id = $1 AND product_id = $2`, "decision-a-"+token, want.ProductID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate row count = %d, want 1", count)
	}
}

func TestMatchDecisionEmptyArraysAndCatalogProductRoundTrip(t *testing.T) {
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "match_decision_empty_"+token)
	repo := marketpostgres.NewRepository(pool, "decision-empty-"+token)
	catalogID := "catalog-" + token
	want := domain.MatchDecision{
		ProductID: "product-" + token, CatalogProductID: &catalogID, Status: domain.MatchStatusNoCandidate,
		Anchors: nil, Contradictions: nil, ConfidenceBand: domain.ConfidenceBandBaixa,
		DecidedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
	if err := repo.AppendMatchDecision(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := repo.LatestMatchDecisions(ctx, []string{want.ProductID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("rows = %d, want 1", len(got))
	}
	d := got[0]
	// nil anchors/contradictions persist as TEXT[] '{}' and read back as empty, non-nil.
	if d.Anchors == nil || len(d.Anchors) != 0 || d.Contradictions == nil || len(d.Contradictions) != 0 {
		t.Fatalf("empty arrays = %#v/%#v, want non-nil empty", d.Anchors, d.Contradictions)
	}
	if d.CatalogProductID == nil || *d.CatalogProductID != catalogID {
		t.Fatalf("catalog_product_id = %#v, want %q", d.CatalogProductID, catalogID)
	}
}
