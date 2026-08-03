package ports

import "context"

// ProductFiscalFacts is the per-product fiscal snapshot TaxesForItem needs
// beyond the matrix cell — grupo_icms selects the matrix cell itself
// (caller's job, before calling FactsFor's sibling CellFor), the rest feed
// directly into pricing/domain.ICMSCell. Found=false means no products_mirror
// row for (tenantID, active source, codigoProduto) — D-37/ADR-17: the caller
// must propagate that as a fully unknown ICMSCell, never a fabricated 0/nil
// mix.
type ProductFiscalFacts struct {
	Found           bool
	GrupoICMS       *int
	Origprod        *int
	StRetidoEntrada *string
	RestituicaoUnit *string
	FiscalDtRef     *string
}

// ProductFiscalReader resolves a product's fiscal snapshot from
// products_mirror, honoring the tenant's active-source routing (D-121
// mirror-first: only the Sankhya sync path populates these 5 columns today —
// a tenant on xlsx/catalogo_cliente active source must see Found=false, not a
// stale row from an inactive source).
type ProductFiscalReader interface {
	// FactsFor returns the fiscal snapshot for (tenantID, codigoProduto)
	// resolved against the tenant's current active source. Found=false when
	// no row exists for that (tenant, active source, codigo_produto) key.
	FactsFor(ctx context.Context, tenantID, codigoProduto string) (ProductFiscalFacts, error)
}
