// P2.b Task 3: icms_matrix_mirror is versioned (bitemporal), unlike products_mirror's
// keep-absent-in-place pattern. A cell's fiscal value changes when TGFICM's vigência
// changes (a new rate takes effect) — we must be able to say "as of date X, this cell
// was Y", so a change closes the old row (vigente_ate=now) and opens a new one, rather
// than overwriting in place. An identical value across syncs only bumps synced_at (no
// version bloat from re-syncing the same fact). A cell key that disappears from the
// current resolution round (e.g. Oracle no longer yields any whitelist survivor for
// that uf_destino×grupo) closes its open row and is NOT reopened by a future sync
// finding the same key again — that would misrepresent a gap as continuity.
package mirror

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

// ErrEmptyICMSMatrix is returned by ApplyCells when cells is empty. Applying it would
// close every currently-open cell for the tenant (nothing survived this round), which
// for a genuinely empty resolution is almost always an upstream failure (Oracle
// unreachable / query error) rather than "the matrix has no cells". Fail-closed,
// mirroring ErrEmptySnapshot's rationale for products_mirror.
var ErrEmptyICMSMatrix = errors.New("mirror: refusing to apply empty icms matrix")

// ICMSMatrixWriter applies domain.ICMSMatrixCell resolutions to icms_matrix_mirror
// with versioned write semantics.
type ICMSMatrixWriter struct {
	pool *pgxpool.Pool
	// Now is exported for test override (deterministic vigente_desde/vigente_ate);
	// defaults to time.Now.
	Now func() time.Time
}

// NewICMSMatrixWriter builds an ICMSMatrixWriter over pool.
func NewICMSMatrixWriter(pool *pgxpool.Pool) *ICMSMatrixWriter {
	return &ICMSMatrixWriter{pool: pool, Now: time.Now}
}

type icmsMatrixCellKey struct {
	UFOrigem  string
	UFDestino string
	GrupoICMS int
}

type icmsMatrixExistingCell struct {
	CodTrib          *int
	Aliquota         *float64
	RedBase          *float64
	PercFcp          *float64
	LinhasCandidatas int
	Ambiguo          bool
}

func (e icmsMatrixExistingCell) equals(c domain.ICMSMatrixCell) bool {
	return intPtrEqual(e.CodTrib, c.CodTrib) &&
		floatPtrEqual(e.Aliquota, c.Aliquota) &&
		floatPtrEqual(e.RedBase, c.RedBase) &&
		floatPtrEqual(e.PercFcp, c.PercFcp) &&
		e.LinhasCandidatas == c.LinhasCandidatas &&
		e.Ambiguo == c.Ambiguo
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func floatPtrEqual(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

const selectOpenCellsSQL = `
SELECT uf_origem, uf_destino, grupo_icms, codtrib, aliquota, red_base, perc_fcp,
       linhas_candidatas, ambiguo
FROM icms_matrix_mirror
WHERE tenant_id = $1 AND vigente_ate IS NULL`

const closeCellSQL = `
UPDATE icms_matrix_mirror
SET vigente_ate = $5
WHERE tenant_id = $1 AND uf_origem = $2 AND uf_destino = $3 AND grupo_icms = $4
  AND vigente_ate IS NULL`

const touchCellSyncedAtSQL = `
UPDATE icms_matrix_mirror
SET synced_at = $5
WHERE tenant_id = $1 AND uf_origem = $2 AND uf_destino = $3 AND grupo_icms = $4
  AND vigente_ate IS NULL`

const insertCellSQL = `
INSERT INTO icms_matrix_mirror
	(tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, aliquota, red_base, perc_fcp,
	 linhas_candidatas, ambiguo, vigente_desde, vigente_ate, synced_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NULL, $11)`

// ApplyCells resolves each cell against the currently open row (if any) for its key:
//   - no open row: insert a new open row (vigente_desde=now).
//   - open row with an identical fiscal value: bump synced_at only, no new version.
//   - open row with a different fiscal value: close it (vigente_ate=now) and insert a
//     new open row.
//
// Any open row whose key is absent from cells this round is closed and NOT reopened.
// Returns the number of cells applied (inserted or matched); the close-only sweep
// count is not included.
func (w *ICMSMatrixWriter) ApplyCells(ctx context.Context, tenantID string, cells []domain.ICMSMatrixCell) (int, error) {
	if tenantID == "" {
		return 0, errors.New("mirror: empty tenantID")
	}
	if len(cells) == 0 {
		return 0, ErrEmptyICMSMatrix
	}

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("mirror: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	existing := make(map[icmsMatrixCellKey]icmsMatrixExistingCell)
	rows, err := tx.Query(ctx, selectOpenCellsSQL, tenantID)
	if err != nil {
		return 0, fmt.Errorf("mirror: select open icms matrix cells: %w", err)
	}
	for rows.Next() {
		var k icmsMatrixCellKey
		var e icmsMatrixExistingCell
		if err := rows.Scan(&k.UFOrigem, &k.UFDestino, &k.GrupoICMS, &e.CodTrib, &e.Aliquota,
			&e.RedBase, &e.PercFcp, &e.LinhasCandidatas, &e.Ambiguo); err != nil {
			rows.Close()
			return 0, fmt.Errorf("mirror: scan open icms matrix cell: %w", err)
		}
		existing[k] = e
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("mirror: read open icms matrix cells: %w", err)
	}

	now := w.Now()
	seen := make(map[icmsMatrixCellKey]bool, len(cells))

	for _, c := range cells {
		key := icmsMatrixCellKey{UFOrigem: c.UFOrigem, UFDestino: c.UFDestino, GrupoICMS: c.GrupoICMS}
		seen[key] = true

		if prev, ok := existing[key]; ok && prev.equals(c) {
			if _, err := tx.Exec(ctx, touchCellSyncedAtSQL, tenantID, c.UFOrigem, c.UFDestino, c.GrupoICMS, now); err != nil {
				return 0, fmt.Errorf("mirror: touch icms matrix cell synced_at: %w", err)
			}
			continue
		}

		if _, ok := existing[key]; ok {
			if _, err := tx.Exec(ctx, closeCellSQL, tenantID, c.UFOrigem, c.UFDestino, c.GrupoICMS, now); err != nil {
				return 0, fmt.Errorf("mirror: close icms matrix cell: %w", err)
			}
		}
		if _, err := tx.Exec(ctx, insertCellSQL, tenantID, c.UFOrigem, c.UFDestino, c.GrupoICMS,
			c.CodTrib, c.Aliquota, c.RedBase, c.PercFcp, c.LinhasCandidatas, c.Ambiguo, now); err != nil {
			return 0, fmt.Errorf("mirror: insert icms matrix cell: %w", err)
		}
	}

	for key := range existing {
		if seen[key] {
			continue
		}
		if _, err := tx.Exec(ctx, closeCellSQL, tenantID, key.UFOrigem, key.UFDestino, key.GrupoICMS, now); err != nil {
			return 0, fmt.Errorf("mirror: close disappeared icms matrix cell: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("mirror: commit: %w", err)
	}
	return len(cells), nil
}
