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

// AliquotaInternaFor returns the vigente icms_aliquota_interna row for uf as
// (alíquota headline, fcp_embutido), or nil when no vigente row exists
// (D-37). icms_aliquota_interna is a GLOBAL table (a lei não é por tenant) —
// this is the one query in this reader with no tenant_id predicate, by
// design.
//
// The old version of this comment said fcp_embutido was informational/
// audit-only and "not read here" because it was already folded into the
// headline aliquota — that claim is REFUTED by measurement (Task A2):
// against 6 real RJ invoice items (roundV.txt V2), the ERP abates ICMS+DIFAL
// (which use the FCP-included headline rate) from the PIS/COFINS base but
// leaves the FCP portion IN that base. Reproducing that split needs BOTH
// numbers from the same vigente row, so fcp_embutido is read here now.
func (r *MatrixReader) AliquotaInternaFor(ctx context.Context, uf string) (*ports.AliquotaInterna, error) {
	var aliquota, fcpEmbutido string
	err := r.pool.QueryRow(ctx, `
		SELECT aliquota::text, fcp_embutido::text
		FROM icms_aliquota_interna
		WHERE uf = $1 AND vigente_ate IS NULL
	`, uf).Scan(&aliquota, &fcpEmbutido)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ports.AliquotaInterna{AliquotaPct: aliquota, FcpEmbutidoPct: fcpEmbutido}, nil
}
