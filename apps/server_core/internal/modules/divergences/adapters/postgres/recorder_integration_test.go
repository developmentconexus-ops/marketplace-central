//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/divergences/domain"
	"marketplace-central/apps/server_core/internal/platform/migrate"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
	canonical "marketplace-central/apps/server_core/migrations"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

// TestEvaluateOpensThenAutoResolvesOnConvergence is the feature.md-required
// divergence integration proof, both directions in one test: create → open
// row; converge → resolved_at set, no new row (2 evaluations, not 2 rows).
func TestEvaluateOpensThenAutoResolvesOnConvergence(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_divergence")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM divergences WHERE tenant_id = $1`, "tenant_f02_divergence")
	})

	t1 := time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 7, 31, 9, 5, 0, 0, time.UTC)
	rec := NewRecorder(pool, WithClock(fixedClock(t1)))

	expectedObs := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observedObs := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	baseInput := domain.EvaluateInput{
		TenantID:           "tenant_f02_divergence",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB123",
		Kind:               domain.KindEstoque,
		Linked:             true,
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expectedObs,
		ObservedObservedAt: &observedObs,
	}

	// 1) Divergent: opens a new row.
	divergent := baseInput
	divergent.ExpectedValue = "12"
	divergent.ObservedValue = "9"
	outcome, err := rec.Evaluate(ctx, divergent)
	if err != nil {
		t.Fatalf("first Evaluate (divergent): %v", err)
	}
	if !outcome.Diverged || !outcome.Opened {
		t.Fatalf("expected Diverged=true Opened=true, got %+v", outcome)
	}

	var openCount int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4 AND resolved_at IS NULL
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&openCount); err != nil {
		t.Fatalf("count open rows: %v", err)
	}
	if openCount != 1 {
		t.Fatalf("expected exactly 1 open row after the divergent evaluation, got %d", openCount)
	}

	var detectedAt time.Time
	if err := pool.QueryRow(ctx, `
		SELECT detected_at FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4 AND resolved_at IS NULL
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&detectedAt); err != nil {
		t.Fatalf("select detected_at: %v", err)
	}
	if !detectedAt.Truncate(time.Microsecond).Equal(t1.Truncate(time.Microsecond)) {
		t.Fatalf("expected detected_at %v, got %v", t1, detectedAt)
	}

	// 2) Still divergent, later clock: refreshes the SAME open row.
	// detected_at must stay pinned to t1, never move to t2 (IC-02
	// Persistence Expectations — "detected_at NÃO muda").
	rec2 := NewRecorder(pool, WithClock(fixedClock(t2)))
	stillDivergent := baseInput
	stillDivergent.ExpectedValue = "12"
	stillDivergent.ObservedValue = "8"
	outcome, err = rec2.Evaluate(ctx, stillDivergent)
	if err != nil {
		t.Fatalf("second Evaluate (still divergent): %v", err)
	}
	if !outcome.Diverged || !outcome.Updated || outcome.Opened {
		t.Fatalf("expected Diverged=true Updated=true Opened=false on re-evaluation of an open row, got %+v", outcome)
	}

	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4 AND resolved_at IS NULL
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&openCount); err != nil {
		t.Fatalf("count open rows: %v", err)
	}
	if openCount != 1 {
		t.Fatalf("expected STILL exactly 1 open row (refreshed, not duplicated), got %d", openCount)
	}

	var detectedAtAfterRefresh time.Time
	if err := pool.QueryRow(ctx, `
		SELECT detected_at FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4 AND resolved_at IS NULL
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&detectedAtAfterRefresh); err != nil {
		t.Fatalf("select detected_at after refresh: %v", err)
	}
	if !detectedAtAfterRefresh.Truncate(time.Microsecond).Equal(t1.Truncate(time.Microsecond)) {
		t.Fatalf("MUST-FAIL guard: detected_at must remain immutable at %v, got %v (moved to a later re-evaluation clock)", t1, detectedAtAfterRefresh)
	}

	// 3) Converges: auto-resolves, no new row.
	t3 := time.Date(2026, 7, 31, 9, 10, 0, 0, time.UTC)
	rec3 := NewRecorder(pool, WithClock(fixedClock(t3)))
	converged := baseInput
	converged.ExpectedValue = "12"
	converged.ObservedValue = "12"
	outcome, err = rec3.Evaluate(ctx, converged)
	if err != nil {
		t.Fatalf("third Evaluate (converged): %v", err)
	}
	if outcome.Diverged {
		t.Fatalf("expected Diverged=false on convergence, got %+v", outcome)
	}
	if !outcome.Resolved {
		t.Fatalf("expected Resolved=true (auto-resolve of the open row), got %+v", outcome)
	}

	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4 AND resolved_at IS NULL
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&openCount); err != nil {
		t.Fatalf("count open rows: %v", err)
	}
	if openCount != 0 {
		t.Fatalf("expected ZERO open rows after convergence (auto-resolved), got %d", openCount)
	}

	var totalRows int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM divergences
		WHERE tenant_id=$1 AND entity_type=$2 AND entity_id=$3 AND kind=$4
	`, "tenant_f02_divergence", domain.EntityTypeListing, "MLB123", domain.KindEstoque).Scan(&totalRows); err != nil {
		t.Fatalf("count total rows: %v", err)
	}
	if totalRows != 1 {
		t.Fatalf("expected exactly 1 row total (resolved, kept as history — no delete, no duplicate), got %d", totalRows)
	}
}

// TestEvaluateConvergentWithNoOpenRowIsNoOp proves the quiet no-op branch:
// a convergent evaluation with nothing currently open does nothing (no row
// created, no error) — it is not an error to evaluate a healthy entity.
func TestEvaluateConvergentWithNoOpenRowIsNoOp(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_divergence_noop")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	expectedObs := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observedObs := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	rec := NewRecorder(pool)
	outcome, err := rec.Evaluate(ctx, domain.EvaluateInput{
		TenantID:           "tenant_f02_divergence_noop",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB_HEALTHY",
		Kind:               domain.KindEstoque,
		Linked:             true,
		ExpectedValue:      "12",
		ObservedValue:      "12",
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expectedObs,
		ObservedObservedAt: &observedObs,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if outcome.Resolved || outcome.Opened || outcome.Updated {
		t.Fatalf("expected a pure no-op, got %+v", outcome)
	}

	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM divergences WHERE tenant_id=$1`, "tenant_f02_divergence_noop").Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected zero rows written, got %d", count)
	}
}

// TestEvaluateFlappingReopensNewRowWithNewDetectedAt proves IC-02's flap
// rule: after a resolve, a fresh divergence opens a BRAND NEW row (new
// detected_at) — the partial unique index permits exactly this because the
// prior row is no longer open.
func TestEvaluateFlappingReopensNewRowWithNewDetectedAt(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_flap")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM divergences WHERE tenant_id = $1`, "tenant_f02_flap")
	})

	expectedObs := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observedObs := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	base := domain.EvaluateInput{
		TenantID:           "tenant_f02_flap",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB_FLAP",
		Kind:               domain.KindEstoque,
		Linked:             true,
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expectedObs,
		ObservedObservedAt: &observedObs,
	}

	t1 := time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC)
	open1 := base
	open1.ExpectedValue, open1.ObservedValue = "12", "9"
	if _, err := NewRecorder(pool, WithClock(fixedClock(t1))).Evaluate(ctx, open1); err != nil {
		t.Fatalf("open first divergence: %v", err)
	}

	t2 := time.Date(2026, 7, 31, 9, 5, 0, 0, time.UTC)
	resolve := base
	resolve.ExpectedValue, resolve.ObservedValue = "12", "12"
	if _, err := NewRecorder(pool, WithClock(fixedClock(t2))).Evaluate(ctx, resolve); err != nil {
		t.Fatalf("resolve first divergence: %v", err)
	}

	t3 := time.Date(2026, 7, 31, 9, 10, 0, 0, time.UTC)
	open2 := base
	open2.ExpectedValue, open2.ObservedValue = "12", "7"
	outcome, err := NewRecorder(pool, WithClock(fixedClock(t3))).Evaluate(ctx, open2)
	if err != nil {
		t.Fatalf("reopen after flap: %v", err)
	}
	if !outcome.Opened {
		t.Fatalf("expected a brand-new open row on flap-reopen (partial unique index allows it), got %+v", outcome)
	}

	var totalRows, openRows int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM divergences WHERE tenant_id=$1`, "tenant_f02_flap").Scan(&totalRows); err != nil {
		t.Fatalf("count total: %v", err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM divergences WHERE tenant_id=$1 AND resolved_at IS NULL`, "tenant_f02_flap").Scan(&openRows); err != nil {
		t.Fatalf("count open: %v", err)
	}
	if totalRows != 2 {
		t.Fatalf("expected 2 total rows (1 resolved history + 1 new open), got %d", totalRows)
	}
	if openRows != 1 {
		t.Fatalf("expected exactly 1 open row after the flap-reopen, got %d", openRows)
	}

	var newDetectedAt time.Time
	if err := pool.QueryRow(ctx, `SELECT detected_at FROM divergences WHERE tenant_id=$1 AND resolved_at IS NULL`, "tenant_f02_flap").Scan(&newDetectedAt); err != nil {
		t.Fatalf("select new detected_at: %v", err)
	}
	if !newDetectedAt.Truncate(time.Microsecond).Equal(t3.Truncate(time.Microsecond)) {
		t.Fatalf("expected the reopened row's detected_at to be the NEW detection time %v, got %v", t3, newDetectedAt)
	}
}

// TestEvaluateSkipsWithoutLinkageWritesNothing is the integration-level
// half of the skip rule: a real DB is present, yet Evaluate must not touch
// it at all when Linked is false.
func TestEvaluateSkipsWithoutLinkageWritesNothing(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_no_linkage")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	expectedObs := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observedObs := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	outcome, err := NewRecorder(pool).Evaluate(ctx, domain.EvaluateInput{
		TenantID:           "tenant_f02_no_linkage",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB_UNLINKED",
		Kind:               domain.KindEstoque,
		Linked:             false,
		ExpectedValue:      "12",
		ObservedValue:      "9",
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expectedObs,
		ObservedObservedAt: &observedObs,
	})
	if err != nil {
		t.Fatalf("expected no error for a skip, got %v", err)
	}
	if !outcome.Skipped {
		t.Fatalf("expected Skipped=true, got %+v", outcome)
	}

	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM divergences WHERE tenant_id=$1`, "tenant_f02_no_linkage").Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected zero rows written for an unlinked entity, got %d", count)
	}
}
