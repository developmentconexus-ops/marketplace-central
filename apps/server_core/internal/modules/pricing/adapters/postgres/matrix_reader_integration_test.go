//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestMatrixReaderCellForReadsOnlyVigenteRow proves CellFor reads the
// vigente_ate IS NULL row and never a closed (historical) version — a
// versioned cell with two rows for the same key (one closed, one open) must
// return the OPEN row's codtrib/ambiguo, not the closed one's.
func TestMatrixReaderCellForReadsOnlyVigenteRow(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_matrix_reader")
	tenant := "matrix-reader-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icms_matrix_mirror WHERE tenant_id = $1", tenant)
	})

	closedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	// Closed (historical) row: codtrib=60 (ST) — must NOT be what CellFor sees.
	if _, err := pool.Exec(ctx, `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota, red_base, perc_fcp,
			 linhas_candidatas, ambiguo, vigente_desde, vigente_ate, synced_at)
		VALUES ($1, 'MG', 'SP', 7, 60, 18.0, 0, 0, 1, false, $2, $3, $2)
	`, tenant, closedAt.Add(-time.Hour), closedAt); err != nil {
		t.Fatalf("insert closed row: %v", err)
	}
	// Open (vigente) row: codtrib=0 (not ST).
	if _, err := pool.Exec(ctx, `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota, red_base, perc_fcp,
			 linhas_candidatas, ambiguo, vigente_desde, vigente_ate, synced_at)
		VALUES ($1, 'MG', 'SP', 7, 0, 18.0, 0, 0, 1, false, $2, NULL, $2)
	`, tenant, closedAt); err != nil {
		t.Fatalf("insert open row: %v", err)
	}

	r := postgres.NewMatrixReader(pool)
	cell, err := r.CellFor(ctx, tenant, "MG", "SP", 7)
	if err != nil {
		t.Fatalf("CellFor: %v", err)
	}
	if !cell.Found {
		t.Fatalf("cell.Found = false, want true")
	}
	if cell.CodTrib == nil || *cell.CodTrib != 0 {
		t.Fatalf("cell.CodTrib = %v, want 0 (the VIGENTE row, not the closed codtrib=60 one)", cell.CodTrib)
	}
	if cell.Ambiguo {
		t.Fatalf("cell.Ambiguo = true, want false")
	}
}

// TestMatrixReaderCellForNoVigenteRowIsNotFound proves a key with no vigente
// row (or no row at all) returns Found=false — never a fabricated codtrib.
func TestMatrixReaderCellForNoVigenteRowIsNotFound(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_matrix_reader")
	tenant := "matrix-reader-absent-" + time.Now().UTC().Format("150405.000000000")

	r := postgres.NewMatrixReader(pool)
	cell, err := r.CellFor(ctx, tenant, "MG", "AC", 9)
	if err != nil {
		t.Fatalf("CellFor: %v", err)
	}
	if cell.Found {
		t.Fatalf("cell.Found = true for a key with no row, want false")
	}
	if cell.CodTrib != nil {
		t.Fatalf("cell.CodTrib = %v for Found=false, want nil", cell.CodTrib)
	}
}

// TestMatrixReaderCellForAmbiguoCarriesNoCodTrib proves an ambíguo cell
// (ResolveCell's "never chooses blind" contract, written by Task 3's writer
// with CodTrib=nil) round-trips as Found=true, Ambiguo=true, CodTrib=nil —
// the reader must not confuse "cell absent" with "cell ambíguo".
func TestMatrixReaderCellForAmbiguoCarriesNoCodTrib(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_matrix_reader")
	tenant := "matrix-reader-ambiguo-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icms_matrix_mirror WHERE tenant_id = $1", tenant)
	})

	now := time.Now().UTC()
	if _, err := pool.Exec(ctx, `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota, red_base, perc_fcp,
			 linhas_candidatas, ambiguo, vigente_desde, vigente_ate, synced_at)
		VALUES ($1, 'MG', 'PR', 311, NULL, NULL, NULL, NULL, 2, true, $2, NULL, $2)
	`, tenant, now); err != nil {
		t.Fatalf("insert ambiguo row: %v", err)
	}

	r := postgres.NewMatrixReader(pool)
	cell, err := r.CellFor(ctx, tenant, "MG", "PR", 311)
	if err != nil {
		t.Fatalf("CellFor: %v", err)
	}
	if !cell.Found {
		t.Fatalf("cell.Found = false, want true (ambíguo IS a resolved-but-unresolvable finding)")
	}
	if !cell.Ambiguo {
		t.Fatalf("cell.Ambiguo = false, want true")
	}
	if cell.CodTrib != nil {
		t.Fatalf("cell.CodTrib = %v, want nil (ambíguo never carries a chosen codtrib)", *cell.CodTrib)
	}
}

// TestMatrixReaderAliquotaInternaForReadsVigenteAndFoldsFCP proves
// AliquotaInternaFor reads the vigente row and returns BOTH the headline
// aliquota (FCP already folded in per D-43) AND fcp_embutido from that same
// row — Task A2: the old claim that fcp_embutido was informational-only and
// not read here was refuted by measurement (roundV.txt V2), so this test now
// pins both fields instead of just the headline one.
func TestMatrixReaderAliquotaInternaForReadsVigenteAndFoldsFCP(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_matrix_reader")
	uf := "ZZ" // synthetic UF, never collides with the real 27-row seed
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM icms_aliquota_interna WHERE uf = $1", uf)
	})

	// Closed (historical) row: 18.0 — must NOT be what AliquotaInternaFor sees.
	if _, err := pool.Exec(ctx, `
		INSERT INTO icms_aliquota_interna (uf, aliquota, fcp_embutido, fonte, lei, vigente_desde, vigente_ate)
		VALUES ($1, 18.0, 0, 'teste', 'lei-velha', '2020-01-01', '2024-01-01')
	`, uf); err != nil {
		t.Fatalf("insert closed row: %v", err)
	}
	// Open (vigente) row: 22.0 with fcp_embutido=2.0 already folded into the headline value.
	if _, err := pool.Exec(ctx, `
		INSERT INTO icms_aliquota_interna (uf, aliquota, fcp_embutido, fonte, lei, vigente_desde, vigente_ate)
		VALUES ($1, 22.0, 2.0, 'teste', 'lei-nova', '2024-01-01', NULL)
	`, uf); err != nil {
		t.Fatalf("insert open row: %v", err)
	}

	r := postgres.NewMatrixReader(pool)
	got, err := r.AliquotaInternaFor(ctx, uf)
	if err != nil {
		t.Fatalf("AliquotaInternaFor: %v", err)
	}
	if got == nil {
		t.Fatalf("AliquotaInternaFor(%q) = nil, want the vigente row", uf)
	}
	if got.AliquotaPct != "22.000" {
		t.Fatalf("AliquotaInternaFor(%q).AliquotaPct = %v, want \"22.000\" (a vigente, nao a fechada 18.0)", uf, got.AliquotaPct)
	}
	if got.FcpEmbutidoPct != "2.000" {
		t.Fatalf("AliquotaInternaFor(%q).FcpEmbutidoPct = %v, want \"2.000\" (fcp_embutido da mesma linha vigente, agora lido — D-43 revisado)", uf, got.FcpEmbutidoPct)
	}
}

// TestMatrixReaderAliquotaInternaForNoVigenteRowIsNil proves a UF with no
// vigente row returns nil, never a default/approximation (D-37).
func TestMatrixReaderAliquotaInternaForNoVigenteRowIsNil(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_matrix_reader")

	r := postgres.NewMatrixReader(pool)
	got, err := r.AliquotaInternaFor(ctx, "QQ")
	if err != nil {
		t.Fatalf("AliquotaInternaFor: %v", err)
	}
	if got != nil {
		t.Fatalf("AliquotaInternaFor(\"QQ\") = %v, want nil (no vigente row for a UF outside the 27-row seed)", *got)
	}
}
