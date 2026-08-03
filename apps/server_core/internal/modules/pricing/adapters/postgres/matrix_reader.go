package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

var _ ports.TaxMatrixReader = (*MatrixReader)(nil)

// MatrixReader reads the VIGENTE icms_matrix_mirror cell (Task 3's writer
// output) and the VIGENTE icms_aliquota_interna row (Task 1's legal seed) —
// both filtered by vigente_ate IS NULL, never as-of (P2.b reads only the
// current version; as-of is P2.c). Numeric columns come back via ::text so
// values never round-trip through float64.
type MatrixReader struct {
	pool *pgxpool.Pool
}

// NewMatrixReader builds a MatrixReader over pool.
func NewMatrixReader(pool *pgxpool.Pool) *MatrixReader {
	return &MatrixReader{pool: pool}
}

// CellFor returns the vigente icms_matrix_mirror cell for (tenantID,
// ufOrigem, ufDestino, grupoICMS). No vigente row ⇒ MatrixCell{Found: false}
// — D-37, never a fabricated codtrib.
func (r *MatrixReader) CellFor(ctx context.Context, tenantID, ufOrigem, ufDestino string, grupoICMS int) (ports.MatrixCell, error) {
	var codTrib *int
	var ambiguo bool
	err := r.pool.QueryRow(ctx, `
		SELECT codtrib, ambiguo
		FROM icms_matrix_mirror
		WHERE tenant_id = $1 AND uf_origem = $2 AND uf_destino = $3 AND grupo_icms = $4
		  AND vigente_ate IS NULL
	`, tenantID, ufOrigem, ufDestino, grupoICMS).Scan(&codTrib, &ambiguo)
	if err == pgx.ErrNoRows {
		return ports.MatrixCell{Found: false}, nil
	}
	if err != nil {
		return ports.MatrixCell{}, err
	}
	return ports.MatrixCell{Found: true, CodTrib: codTrib, Ambiguo: ambiguo}, nil
}

// AliquotaInternaFor returns the vigente icms_aliquota_interna.aliquota for
// uf as a decimal-percent string, or nil when no vigente row exists (D-37).
// fcp_embutido is not read here — D-43: it is informational/audit-only,
// already folded into the headline aliquota column where general.
func (r *MatrixReader) AliquotaInternaFor(ctx context.Context, uf string) (*string, error) {
	var aliquota string
	err := r.pool.QueryRow(ctx, `
		SELECT aliquota::text
		FROM icms_aliquota_interna
		WHERE uf = $1 AND vigente_ate IS NULL
	`, uf).Scan(&aliquota)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &aliquota, nil
}
