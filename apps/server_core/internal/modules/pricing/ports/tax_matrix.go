package ports

import "context"

// MatrixCell is the VIGENTE icms_matrix_mirror row (vigente_ate IS NULL) for
// one (uf_origem, uf_destino, grupo_icms) key. Found=false means no vigente
// row for that key — D-37: UF/grupo sem linha na matriz vira pendência
// explícita, nunca uma alíquota aproximada. CodTrib nil with Found=true means
// the cell was resolved as ambíguo (ResolveCell never picks a candidate's
// codtrib blind) — Ambiguo carries that distinction explicitly rather than
// overloading CodTrib's nil-ness with two different meanings.
type MatrixCell struct {
	Found   bool
	CodTrib *int
	Ambiguo bool
}

// TaxMatrixReader resolves the two facts pricing/domain.TaxesForItem needs
// beyond the product's own fiscal columns (origprod, st_retido_entrada,
// restituicao_unit — read by the caller, Task 5's per-item adapter):
//   - the vigente ICMS matrix cell (is this product ST at this destino?)
//   - the vigente alíquota interna legal for a UF (D-42: our own versioned
//     table, never TGFICM — the ERP's copy is up to two years stale in
//     RJ/PR/DF)
type TaxMatrixReader interface {
	// CellFor returns the vigente cell for (tenantID, ufOrigem, ufDestino,
	// grupoICMS). MatrixCell{Found: false} when no vigente row exists — the
	// caller must propagate that as an unknown ICMSCell.CodTrib, never a
	// fabricated codtrib.
	CellFor(ctx context.Context, tenantID, ufOrigem, ufDestino string, grupoICMS int) (MatrixCell, error)

	// AliquotaInternaFor returns icms_aliquota_interna.aliquota (FCP já
	// embutido onde geral, D-43) for uf as a decimal-percent string, or nil
	// when no vigente row exists for that UF (D-37 — never falls back to a
	// neighboring UF or a default).
	AliquotaInternaFor(ctx context.Context, uf string) (*string, error)
}
