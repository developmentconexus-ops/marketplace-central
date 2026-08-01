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

// TestListOpenByEntityReturnsOpenRowsNewestFirst proves the Reader side of
// IC-02 Operations: only resolved_at IS NULL rows come back, ordered by
// detected_at DESC, and a resolved row must not leak through.
func TestListOpenByEntityReturnsOpenRowsNewestFirst(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_divergence_reader")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM divergences WHERE tenant_id = $1`, "tenant_f02_divergence_reader")
	})

	expectedObs := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observedObs := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	base := domain.EvaluateInput{
		TenantID:           "tenant_f02_divergence_reader",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB_READER",
		Linked:             true,
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expectedObs,
		ObservedObservedAt: &observedObs,
	}

	// Open an estoque divergence at t1 (older), a tarifa divergence at t2
	// (newer) — same entity, different kind, both currently open.
	t1 := time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC)
	estoque := base
	estoque.Kind = domain.KindEstoque
	estoque.ExpectedValue, estoque.ObservedValue = "12", "9"
	if _, err := NewRecorder(pool, WithClock(fixedClock(t1))).Evaluate(ctx, estoque); err != nil {
		t.Fatalf("open estoque divergence: %v", err)
	}

	t2 := time.Date(2026, 7, 31, 9, 5, 0, 0, time.UTC)
	tarifa := base
	tarifa.Kind = domain.KindTarifa
	tarifa.ExpectedValue, tarifa.ObservedValue = "15.99", "16.50"
	if _, err := NewRecorder(pool, WithClock(fixedClock(t2))).Evaluate(ctx, tarifa); err != nil {
		t.Fatalf("open tarifa divergence: %v", err)
	}

	// A third kind on the SAME entity is opened then resolved — must not
	// appear in ListOpenByEntity at all.
	t3 := time.Date(2026, 7, 31, 9, 10, 0, 0, time.UTC)
	resolvedKind := base
	resolvedKind.Kind = domain.KindEstoque
	resolvedKind.EntityID = "MLB_READER_OTHER"
	resolvedKind.ExpectedValue, resolvedKind.ObservedValue = "5", "3"
	if _, err := NewRecorder(pool, WithClock(fixedClock(t3))).Evaluate(ctx, resolvedKind); err != nil {
		t.Fatalf("open other-entity divergence: %v", err)
	}
	pool.Exec(ctx, `DELETE FROM divergences WHERE tenant_id=$1 AND entity_id=$2`, "tenant_f02_divergence_reader", "MLB_READER_OTHER")

	reader := NewReader(pool)
	open, err := reader.ListOpenByEntity(ctx, "tenant_f02_divergence_reader", "mercado_livre", domain.EntityTypeListing, "MLB_READER")
	if err != nil {
		t.Fatalf("ListOpenByEntity: %v", err)
	}
	if len(open) != 2 {
		t.Fatalf("expected exactly 2 open divergences for MLB_READER, got %d: %+v", len(open), open)
	}
	// Newest first: tarifa (t2) before estoque (t1).
	if open[0].Kind != domain.KindTarifa {
		t.Fatalf("expected tarifa (detected at t2, newest) first, got %+v", open[0])
	}
	if open[1].Kind != domain.KindEstoque {
		t.Fatalf("expected estoque (detected at t1, oldest) second, got %+v", open[1])
	}
	if open[0].ExpectedValue != "15.9900" || open[0].ObservedValue != "16.5000" { // numeric(14,4) round-trips zero-padded
		t.Fatalf("expected tarifa field values preserved verbatim, got %+v", open[0])
	}
	if !open[0].DetectedAt.Truncate(time.Microsecond).Equal(t2.Truncate(time.Microsecond)) {
		t.Fatalf("expected tarifa detected_at %v, got %v", t2, open[0].DetectedAt)
	}
	if !open[1].DetectedAt.Truncate(time.Microsecond).Equal(t1.Truncate(time.Microsecond)) {
		t.Fatalf("expected estoque detected_at %v, got %v", t1, open[1].DetectedAt)
	}
}

// TestListOpenByEntityEmptyWhenNoneOpen proves the "no rows" branch returns
// an empty (never nil) slice with no error — a healthy entity is not an
// error condition.
func TestListOpenByEntityEmptyWhenNoneOpen(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_divergence_reader_empty")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	reader := NewReader(pool)
	open, err := reader.ListOpenByEntity(ctx, "tenant_f02_divergence_reader_empty", "mercado_livre", domain.EntityTypeListing, "MLB_NONE")
	if err != nil {
		t.Fatalf("ListOpenByEntity: %v", err)
	}
	if open == nil {
		t.Fatal("expected an empty non-nil slice, got nil")
	}
	if len(open) != 0 {
		t.Fatalf("expected zero open divergences, got %d", len(open))
	}
}
