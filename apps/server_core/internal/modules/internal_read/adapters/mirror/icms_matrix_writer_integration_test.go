//go:build integration

package mirror_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func ptrF64(v float64) *float64 { return &v }

// TestICMSMatrixWriterVersionsOnChangeAndTouchesOnIdentical is the P2.b Task 3
// mandatory writer scenario: round 1 writes one open row; round 2 with a different
// aliquota closes the round-1 row and opens a second one (two rows total, one open);
// round 3 with the SAME value as round 2 must NOT create a third row — only
// synced_at advances on the still-open row. codtrib wrong is the difference between
// 0% and 24% ICMS, so the version history itself must be exact, not approximate.
func TestICMSMatrixWriterVersionsOnChangeAndTouchesOnIdentical(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mirror_icms_matrix")
	tenant := "mirror-icms-matrix-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icms_matrix_mirror WHERE tenant_id = $1", tenant)
	})

	w := mirror.NewICMSMatrixWriter(pool)
	t0 := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	t2 := t1.Add(time.Hour)

	cell := func(codTrib int, aliquota float64) domain.ICMSMatrixCell {
		return domain.ICMSMatrixCell{
			UFOrigem: "MG", UFDestino: "SP", GrupoICMS: 7,
			CodTrib: ptrI(codTrib), Aliquota: ptrF64(aliquota),
			LinhasCandidatas: 1, Ambiguo: false,
		}
	}

	countRows := func() int {
		var n int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id=$1`, tenant).Scan(&n); err != nil {
			t.Fatalf("count rows: %v", err)
		}
		return n
	}
	openRow := func() (codTrib *int, aliquota *float64, syncedAt time.Time, vigenteDesde time.Time) {
		err := pool.QueryRow(ctx, `SELECT codtrib, aliquota, synced_at, vigente_desde FROM icms_matrix_mirror
			WHERE tenant_id=$1 AND vigente_ate IS NULL`, tenant).Scan(&codTrib, &aliquota, &syncedAt, &vigenteDesde)
		if err != nil {
			t.Fatalf("read open row: %v", err)
		}
		return
	}

	// Round 1: insert.
	w.Now = func() time.Time { return t0 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{cell(0, 12.0)}); err != nil {
		t.Fatalf("round1 ApplyCells: %v", err)
	}
	if got := countRows(); got != 1 {
		t.Fatalf("round1 row count = %d, want 1", got)
	}
	if codTrib, aliquota, _, vigenteDesde := openRow(); codTrib == nil || *codTrib != 0 || aliquota == nil || *aliquota != 12.0 || !vigenteDesde.Equal(t0) {
		t.Fatalf("round1 open row = codtrib=%v aliquota=%v vigente_desde=%v, want 0/12.0/%v", codTrib, aliquota, vigenteDesde, t0)
	}

	// Round 2: different aliquota → close round-1 row, open a new one.
	w.Now = func() time.Time { return t1 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{cell(0, 18.0)}); err != nil {
		t.Fatalf("round2 ApplyCells: %v", err)
	}
	if got := countRows(); got != 2 {
		t.Fatalf("round2 row count = %d, want 2 (value changed → new version)", got)
	}
	if codTrib, aliquota, syncedAt, vigenteDesde := openRow(); aliquota == nil || *aliquota != 18.0 || !vigenteDesde.Equal(t1) || !syncedAt.Equal(t1) {
		t.Fatalf("round2 open row = codtrib=%v aliquota=%v synced_at=%v vigente_desde=%v, want */18.0/%v/%v", codTrib, aliquota, syncedAt, vigenteDesde, t1, t1)
	}
	var closedCount int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id=$1 AND vigente_ate = $2`, tenant, t1).Scan(&closedCount); err != nil {
		t.Fatalf("count closed: %v", err)
	}
	if closedCount != 1 {
		t.Fatalf("closed-at-t1 count = %d, want 1 (round-1 row must close exactly at round-2's now)", closedCount)
	}

	// Round 3: SAME value as round 2 → must NOT create a third row, only synced_at advances.
	w.Now = func() time.Time { return t2 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{cell(0, 18.0)}); err != nil {
		t.Fatalf("round3 ApplyCells: %v", err)
	}
	if got := countRows(); got != 2 {
		t.Fatalf("round3 row count = %d, want 2 — identical value must only bump synced_at, never create a third row", got)
	}
	if _, aliquota, syncedAt, vigenteDesde := openRow(); aliquota == nil || *aliquota != 18.0 || !vigenteDesde.Equal(t1) || !syncedAt.Equal(t2) {
		t.Fatalf("round3 open row = aliquota=%v synced_at=%v vigente_desde=%v, want 18.0/%v/%v (same version, synced_at advanced)", aliquota, syncedAt, vigenteDesde, t2, t1)
	}
}

// TestICMSMatrixWriterClosesDisappearedCellWithoutReopening proves the sweep half of
// the versioned contract: a cell key present in round 1 but absent from round 2 gets
// its open row closed, and a later round bringing the SAME key back must open a THIRD
// row rather than un-closing the first — a gap in Oracle's resolution must never be
// silently smoothed over as continuity.
func TestICMSMatrixWriterClosesDisappearedCellWithoutReopening(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mirror_icms_matrix_sweep")
	tenant := "mirror-icms-matrix-sweep-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icms_matrix_mirror WHERE tenant_id = $1", tenant)
	})

	w := mirror.NewICMSMatrixWriter(pool)
	t0 := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	t2 := t1.Add(time.Hour)

	cell := domain.ICMSMatrixCell{
		UFOrigem: "MG", UFDestino: "RJ", GrupoICMS: 3,
		CodTrib: ptrI(0), Aliquota: ptrF64(12.0), LinhasCandidatas: 1, Ambiguo: false,
	}

	w.Now = func() time.Time { return t0 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{cell}); err != nil {
		t.Fatalf("round1: %v", err)
	}

	// Round 2: cell key absent (e.g. no whitelist survivor this round) — close, don't delete.
	otherCell := domain.ICMSMatrixCell{
		UFOrigem: "MG", UFDestino: "BA", GrupoICMS: 9,
		CodTrib: ptrI(0), Aliquota: ptrF64(7.0), LinhasCandidatas: 1, Ambiguo: false,
	}
	w.Now = func() time.Time { return t1 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{otherCell}); err != nil {
		t.Fatalf("round2: %v", err)
	}
	var openCount, closedCount int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id=$1 AND uf_destino='RJ' AND vigente_ate IS NULL`, tenant).Scan(&openCount); err != nil {
		t.Fatalf("count RJ open: %v", err)
	}
	if openCount != 0 {
		t.Fatalf("RJ open row count after disappearing = %d, want 0", openCount)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id=$1 AND uf_destino='RJ' AND vigente_ate = $2`, tenant, t1).Scan(&closedCount); err != nil {
		t.Fatalf("count RJ closed: %v", err)
	}
	if closedCount != 1 {
		t.Fatalf("RJ closed-at-t1 count = %d, want 1", closedCount)
	}

	// Round 3: RJ key reappears — must open a NEW (third overall, second for RJ) row,
	// not resurrect/reopen the round-1 row.
	w.Now = func() time.Time { return t2 }
	if _, err := w.ApplyCells(ctx, tenant, []domain.ICMSMatrixCell{cell, otherCell}); err != nil {
		t.Fatalf("round3: %v", err)
	}
	var rjRows int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id=$1 AND uf_destino='RJ'`, tenant).Scan(&rjRows); err != nil {
		t.Fatalf("count RJ total: %v", err)
	}
	if rjRows != 2 {
		t.Fatalf("RJ total row count after reappearing = %d, want 2 (round1 closed + round3 new, never reopened)", rjRows)
	}
	var rjOpenVigenteDesde time.Time
	if err := pool.QueryRow(ctx, `SELECT vigente_desde FROM icms_matrix_mirror WHERE tenant_id=$1 AND uf_destino='RJ' AND vigente_ate IS NULL`, tenant).Scan(&rjOpenVigenteDesde); err != nil {
		t.Fatalf("read RJ reopened row: %v", err)
	}
	if !rjOpenVigenteDesde.Equal(t2) {
		t.Fatalf("RJ reopened row vigente_desde = %v, want %v (a fresh row, not the round-1 one un-closed)", rjOpenVigenteDesde, t2)
	}
}
