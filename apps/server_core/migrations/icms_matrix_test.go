//go:build integration

package migrations

// Task 1 (P2.b): live-DB proof that 0093_icms_matrix.sql's partial unique
// indexes are the actual guarantee against two open versions of the same
// cell — not writer discipline. Every other *_test.go in this package
// asserts on the migration's SQL TEXT; this file is deliberately different
// because a partial unique index's enforcement cannot be observed by
// parsing text, only by provoking Postgres into rejecting a real insert.

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestICMSMatrixMirrorVigenteUniqueIndexRejectsSecondOpenRow proves the
// icms_matrix_mirror_vigente partial unique index rejects a second row with
// vigente_ate IS NULL for the same (tenant_id, uf_origem, uf_destino,
// grupo_icms) cell. Before 0093 exists this test fails because the table is
// absent (proves the red before the migration is written); after 0093 the
// first insert succeeds and the second is rejected with SQLSTATE 23505 on
// exactly this index.
func TestICMSMatrixMirrorVigenteUniqueIndexRejectsSecondOpenRow(t *testing.T) {
	const tenantID = "tenant_icms_matrix_vigente"
	pool, _ := testpostgres.OpenPool(t, tenantID)
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM icms_matrix_mirror WHERE tenant_id = $1`, tenantID)
	})

	const insert = `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota,
			 red_base, perc_fcp, linhas_candidatas, ambiguo, vigente_desde, synced_at)
		VALUES ($1, 'SP', 'RJ', 1, NULL, 18.0, 0, 0, 1, false, $2, now())
	`
	now := time.Now().UTC()
	if _, err := pool.Exec(ctx, insert, tenantID, now); err != nil {
		t.Fatalf("first open row must insert cleanly: %v", err)
	}

	// Same cell (tenant_id, uf_origem, uf_destino, grupo_icms), a later
	// vigente_desde, still open (vigente_ate IS NULL) — the exact shape the
	// index exists to forbid: two open versions of one cell.
	_, err := pool.Exec(ctx, insert, tenantID, now.Add(time.Hour))
	if err == nil {
		t.Fatal("expected a unique violation inserting a second open row for the same cell, got none")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected a *pgconn.PgError, got %T: %v", err, err)
	}
	if pgErr.Code != "23505" {
		t.Fatalf("expected SQLSTATE 23505 (unique_violation), got %q: %v", pgErr.Code, pgErr)
	}
	if pgErr.ConstraintName != "icms_matrix_mirror_vigente" {
		t.Fatalf("expected the violation to name icms_matrix_mirror_vigente, got %q — a different constraint caught this, not the one-open-row guarantee", pgErr.ConstraintName)
	}
}

// TestICMSMatrixMirrorVigenteUniqueIndexAllowsAClosedAndAnOpenVersion is the
// must-fail's positive twin: closing the first version (vigente_ate NOT
// NULL) before opening a new one is the ordinary write pattern the index
// must permit, so a naive "no two rows ever" reading of the index would be
// wrong just as much as "no index at all" would be.
func TestICMSMatrixMirrorVigenteUniqueIndexAllowsAClosedAndAnOpenVersion(t *testing.T) {
	const tenantID = "tenant_icms_matrix_vigente_reopen"
	pool, _ := testpostgres.OpenPool(t, tenantID)
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM icms_matrix_mirror WHERE tenant_id = $1`, tenantID)
	})

	now := time.Now().UTC()
	const insertClosed = `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota,
			 red_base, perc_fcp, linhas_candidatas, ambiguo, vigente_desde, vigente_ate, synced_at)
		VALUES ($1, 'SP', 'RJ', 1, NULL, 18.0, 0, 0, 1, false, $2, $3, now())
	`
	if _, err := pool.Exec(ctx, insertClosed, tenantID, now, now.Add(time.Hour)); err != nil {
		t.Fatalf("closed version must insert cleanly: %v", err)
	}

	const insertOpen = `
		INSERT INTO icms_matrix_mirror
			(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota,
			 red_base, perc_fcp, linhas_candidatas, ambiguo, vigente_desde, synced_at)
		VALUES ($1, 'SP', 'RJ', 1, NULL, 18.0, 0, 0, 1, false, $2, now())
	`
	if _, err := pool.Exec(ctx, insertOpen, tenantID, now.Add(time.Hour)); err != nil {
		t.Fatalf("a fresh open version after the prior one closed must insert cleanly: %v", err)
	}
}

// TestProductsMirrorGrupoICMSAcceptsNull proves grupo_icms is a bare
// nullable column on products_mirror (ADR-17): a row written without it
// stays NULL, it never becomes 0. Task 2 (sync) is what populates it; this
// task only owns the column accepting the honest-unknown state.
func TestProductsMirrorGrupoICMSAcceptsNull(t *testing.T) {
	const tenantID = "tenant_icms_grupo_null"
	const codigoProduto = "PROD-ICMS-NULL-1"
	pool, _ := testpostgres.OpenPool(t, tenantID)
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM products_mirror WHERE tenant_id = $1`, tenantID)
	})

	const insert = `
		INSERT INTO products_mirror (tenant_id, source, codigo_produto)
		VALUES ($1, 'xlsx', $2)
	`
	if _, err := pool.Exec(ctx, insert, tenantID, codigoProduto); err != nil {
		t.Fatalf("insert with grupo_icms omitted must succeed (column must be nullable): %v", err)
	}

	var grupoICMS *int32
	row := pool.QueryRow(ctx, `SELECT grupo_icms FROM products_mirror WHERE tenant_id = $1 AND codigo_produto = $2`, tenantID, codigoProduto)
	if err := row.Scan(&grupoICMS); err != nil {
		t.Fatalf("select grupo_icms: %v", err)
	}
	if grupoICMS != nil {
		t.Fatalf("expected grupo_icms to be honest-unknown NULL for an unread product, got %d", *grupoICMS)
	}
}
