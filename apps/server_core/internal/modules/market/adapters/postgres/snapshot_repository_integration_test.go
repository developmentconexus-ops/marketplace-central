//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestSnapshotLatestValidSurvivesLaterFailureAndSeriesShowsBoth(t *testing.T) {
	ctx, repo, _, tenant, token := snapshotHarness(t, "latest")
	now := time.Now().UTC().Truncate(time.Microsecond)
	valid := validSnapshot("listing-"+token, "101.25", now, now.Add(time.Hour))
	reason := "provider unavailable"
	failed := domain.MarketPriceSnapshot{Scope: valid.Scope, RefID: valid.RefID, Source: valid.Source, ObservedAt: now.Add(time.Minute), FetchedAt: now.Add(time.Minute), ExpiresAt: now.Add(time.Hour), Status: domain.MarketPriceSnapshotStatusFailed, FailureReason: &reason, RequestID: "failed-" + token}
	if err := repo.AppendSnapshot(ctx, valid); err != nil {
		t.Fatal(err)
	}
	if err := repo.AppendSnapshot(ctx, failed); err != nil {
		t.Fatal(err)
	}

	got, err := repo.LatestValidSnapshot(ctx, valid.Scope, valid.RefID, valid.Source, now.Add(2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || !got.FetchedAt.Equal(valid.FetchedAt) || got.Price == nil || got.Price.Amount != valid.Price.Amount {
		t.Fatalf("latest valid = %+v, want original %+v", got, valid)
	}
	series, err := repo.SnapshotSeries(ctx, valid.Scope, valid.RefID, valid.Source, now.Add(2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if len(series) != 2 {
		t.Fatalf("series rows = %d, want 2: %+v", len(series), series)
	}
	if series[0].Status != domain.MarketPriceSnapshotStatusFailed || series[0].Price != nil || series[0].FailureReason == nil || *series[0].FailureReason != reason {
		t.Fatalf("series[0] = %+v, want FAILED with nil price and reason %q", series[0], reason)
	}
	_ = tenant
}

func TestLatestValidSnapshotProjectsExpiry(t *testing.T) {
	ctx, repo, _, _, token := snapshotHarness(t, "latestexpiry")
	now := time.Now().UTC().Truncate(time.Microsecond)
	snapshot := validSnapshot("listing-"+token, "77.00", now.Add(-time.Hour), now.Add(-time.Minute))
	if err := repo.AppendSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	got, err := repo.LatestValidSnapshot(ctx, snapshot.Scope, snapshot.RefID, snapshot.Source, now)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Status != domain.MarketPriceSnapshotStatusExpired {
		t.Fatalf("latest valid = %+v, want projected EXPIRED", got)
	}
}

func TestSnapshotDuplicateAppendIsNoOp(t *testing.T) {
	ctx, repo, pool, tenant, token := snapshotHarness(t, "duplicate")
	now := time.Now().UTC().Truncate(time.Microsecond)
	snapshot := validSnapshot("listing-"+token, "22.50", now, now.Add(time.Hour))
	if err := repo.AppendSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	if err := repo.AppendSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM market_price_snapshots WHERE tenant_id = $1 AND scope = $2 AND ref_id = $3 AND source = $4 AND fetched_at = $5`, tenant, snapshot.Scope, snapshot.RefID, snapshot.Source, snapshot.FetchedAt).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("row count = %d, want 1", count)
	}
}

func TestSnapshotRepositoryCannotCrossTenant(t *testing.T) {
	ctx, repoA, pool, _, token := snapshotHarness(t, "tenant")
	repoB := marketpostgres.NewRepository(pool, "snapshot-b-"+token)
	now := time.Now().UTC().Truncate(time.Microsecond)
	snapshot := validSnapshot("listing-"+token, "30.00", now, now.Add(time.Hour))
	if err := repoA.AppendSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	got, err := repoB.SnapshotSeries(ctx, snapshot.Scope, snapshot.RefID, snapshot.Source, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("tenant B rows = %+v, want none", got)
	}
}

func TestSnapshotExpiryIsProjectedWithoutMutation(t *testing.T) {
	ctx, repo, pool, tenant, token := snapshotHarness(t, "expiry")
	now := time.Now().UTC().Truncate(time.Microsecond)
	snapshot := validSnapshot("listing-"+token, "44.00", now.Add(-time.Hour), now.Add(-time.Minute))
	if err := repo.AppendSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	series, err := repo.SnapshotSeries(ctx, snapshot.Scope, snapshot.RefID, snapshot.Source, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(series) != 1 || series[0].Status != domain.MarketPriceSnapshotStatusExpired {
		t.Fatalf("series = %+v, want projected EXPIRED", series)
	}
	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM market_price_snapshots WHERE tenant_id = $1 AND scope = $2 AND ref_id = $3 AND source = $4 AND fetched_at = $5`, tenant, snapshot.Scope, snapshot.RefID, snapshot.Source, snapshot.FetchedAt).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(domain.MarketPriceSnapshotStatusValid) {
		t.Fatalf("stored status = %q, want VALID", status)
	}
}

func snapshotHarness(t *testing.T, name string) (context.Context, *marketpostgres.Repository, *pgxpool.Pool, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_snapshot_"+name)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant := "snapshot-a-" + token
	return ctx, marketpostgres.NewRepository(pool, tenant), pool, tenant, token
}

func validSnapshot(refID, amount string, fetchedAt, expiresAt time.Time) domain.MarketPriceSnapshot {
	return domain.MarketPriceSnapshot{Scope: domain.MarketPriceScopeListing, RefID: refID, Source: domain.MarketPriceSourceMLSalePrice, Price: &domain.Money{Amount: amount, Currency: "BRL"}, ObservedAt: fetchedAt, FetchedAt: fetchedAt, ExpiresAt: expiresAt, Status: domain.MarketPriceSnapshotStatusValid, RequestID: "request-" + refID}
}
