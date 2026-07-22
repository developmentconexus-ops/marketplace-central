// Package mirror is M-04's own products_mirror writer, owned inside internal_read
// (hub ruling R-2 (i)): fully disjoint from M-03's erp_import writer. It applies a
// full snapshot (ADR: full-snapshot v1) to products_mirror + the child
// products_mirror_stock_locations table with keep-absent merge semantics (ADR-04)
// and honest-NULL fields (ADR-17 — a missing value stays NULL, never a fabricated 0).
//
// Tenant scoping (AC-01): every statement is keyed on tenant_id. The Oracle read
// side that feeds Row is single-tenant by construction (one Oracle connection = one
// tenant's ERP; the Sankhya tables carry no tenant column), so tenant_id is stamped
// here at the mirror write, which is the only place the concept exists.
package mirror

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrEmptySnapshot is returned by ApplySnapshot when the snapshot carries zero
// product rows. Applying it would mark the entire mirror absent (keep-absent
// sweep), which for an empty read is almost always an upstream failure (Oracle
// unreachable / query error) rather than "the ERP has no products". We refuse it
// fail-closed instead of wiping the mirror's presence flags.
var ErrEmptySnapshot = errors.New("mirror: refusing to apply empty snapshot")

// Row is one products_mirror row (grain = 1 per CODPROD, as-of today). All
// optional facts are pointers so nil round-trips to SQL NULL — never 0/"".
type Row struct {
	CodigoProduto  string
	Descricao      *string
	Referencia     *string
	EAN            *string
	Marca          *string
	GrupoCodigo    *string
	GrupoDescricao *string
	NCM            *string
	Custo          *float64
	PrecoVenda     *float64
	EstoqueTotal   *float64
}

// StockLocation is one products_mirror_stock_locations row (per-CODLOCAL stock).
type StockLocation struct {
	CodigoProduto  string
	LocalCodigo    string
	LocalDescricao *string
	Quantidade     *float64
}

// Writer applies a snapshot to the mirror. Defined as an interface so the
// SankhyaAdapter depends on the behavior, not the pgx-backed implementation.
type Writer interface {
	ApplySnapshot(ctx context.Context, tenantID string, rows []Row, locs []StockLocation) (int, error)
}

// PgWriter is the Postgres-backed Writer.
type PgWriter struct {
	pool *pgxpool.Pool
}

// NewPgWriter builds a PgWriter over pool.
func NewPgWriter(pool *pgxpool.Pool) *PgWriter {
	return &PgWriter{pool: pool}
}

var _ Writer = (*PgWriter)(nil)

const upsertSQL = `
INSERT INTO products_mirror
	(tenant_id, source, codigo_produto, descricao, referencia, ean, marca,
	 grupo_codigo, grupo_descricao, ncm, custo, preco_venda, estoque_total,
	 absent_in_last_snapshot, stale_since, updated_at)
VALUES ($1, 'sankhya', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, NULL, now())
ON CONFLICT (tenant_id, codigo_produto) DO UPDATE SET
	source = 'sankhya',
	descricao = EXCLUDED.descricao,
	referencia = EXCLUDED.referencia,
	ean = EXCLUDED.ean,
	marca = EXCLUDED.marca,
	grupo_codigo = EXCLUDED.grupo_codigo,
	grupo_descricao = EXCLUDED.grupo_descricao,
	ncm = EXCLUDED.ncm,
	custo = EXCLUDED.custo,
	preco_venda = EXCLUDED.preco_venda,
	estoque_total = EXCLUDED.estoque_total,
	absent_in_last_snapshot = false,
	stale_since = NULL,
	updated_at = now()`

// keepAbsentSQL implements ADR-04: rows present in the mirror for this tenant but
// absent from the current snapshot are flagged (never physically deleted, and their
// last-known values + source are preserved). stale_since is stamped once on the
// present→absent transition (COALESCE keeps the original). Scope is the whole
// tenant regardless of prior source (hub ruling): the mirror reflects the active
// source's complete current view.
const keepAbsentSQL = `
UPDATE products_mirror
SET absent_in_last_snapshot = true,
	stale_since = COALESCE(stale_since, now()),
	updated_at = now()
WHERE tenant_id = $1
	AND absent_in_last_snapshot = false
	AND codigo_produto <> ALL($2::text[])`

const deleteLocationsSQL = `
DELETE FROM products_mirror_stock_locations
WHERE tenant_id = $1 AND codigo_produto = ANY($2::text[])`

const insertLocationSQL = `
INSERT INTO products_mirror_stock_locations
	(tenant_id, codigo_produto, local_codigo, local_descricao, quantidade)
VALUES ($1, $2, $3, $4, $5)`

// ApplySnapshot upserts every row, replaces the stock-location children for the
// products in the snapshot, and runs the keep-absent sweep — all in one
// transaction. Returns the number of product rows processed. All-or-nothing: any
// failure rolls back and returns the error (v1 has no per-row error accounting).
func (w *PgWriter) ApplySnapshot(ctx context.Context, tenantID string, rows []Row, locs []StockLocation) (int, error) {
	if tenantID == "" {
		return 0, errors.New("mirror: empty tenantID")
	}
	if len(rows) == 0 {
		return 0, ErrEmptySnapshot
	}

	codes := make([]string, len(rows))
	for i, r := range rows {
		codes[i] = r.CodigoProduto
	}

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("mirror: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Upsert product rows.
	batch := &pgx.Batch{}
	for _, r := range rows {
		batch.Queue(upsertSQL, tenantID, r.CodigoProduto, r.Descricao, r.Referencia,
			r.EAN, r.Marca, r.GrupoCodigo, r.GrupoDescricao, r.NCM, r.Custo,
			r.PrecoVenda, r.EstoqueTotal)
	}
	// Replace stock-location children only for products in this snapshot.
	batch.Queue(deleteLocationsSQL, tenantID, codes)
	for _, l := range locs {
		batch.Queue(insertLocationSQL, tenantID, l.CodigoProduto, l.LocalCodigo,
			l.LocalDescricao, l.Quantidade)
	}
	// Keep-absent sweep.
	batch.Queue(keepAbsentSQL, tenantID, codes)

	br := tx.SendBatch(ctx, batch)
	total := batch.Len()
	for i := 0; i < total; i++ {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return 0, fmt.Errorf("mirror: apply snapshot batch (op %d/%d): %w", i+1, total, err)
		}
	}
	if err := br.Close(); err != nil {
		return 0, fmt.Errorf("mirror: close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("mirror: commit: %w", err)
	}
	return len(rows), nil
}
