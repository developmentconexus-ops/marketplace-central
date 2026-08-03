package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

var _ ports.ProductFiscalReader = (*ProductFiscalReader)(nil)

// ProductFiscalReader reads a product's fiscal snapshot (grupo_icms,
// origprod, st_retido_entrada, restituicao_unit, fiscal_dt_ref) from
// products_mirror, scoped to the tenant's CURRENT active source (D-121
// mirror-first routing) — only the Sankhya sync path populates these 5
// columns today, so resolving against the wrong source would honestly (and
// silently) return NULL facts for a product the Sankhya mirror actually
// knows about.
type ProductFiscalReader struct {
	pool         *pgxpool.Pool
	activeSource tenant_config.ActiveSourceLookup
}

// NewProductFiscalReader builds a ProductFiscalReader over pool, resolving
// active source per call via activeSource.
func NewProductFiscalReader(pool *pgxpool.Pool, activeSource tenant_config.ActiveSourceLookup) *ProductFiscalReader {
	return &ProductFiscalReader{pool: pool, activeSource: activeSource}
}

// FactsFor resolves tenantID's active source, then reads the fiscal snapshot
// for codigoProduto within that source. Found=false when the tenant has no
// active-source row (ErrUnknownActiveSource) or no products_mirror row for
// (tenantID, source, codigoProduto) — neither case fabricates a fact.
func (r *ProductFiscalReader) FactsFor(ctx context.Context, tenantID, codigoProduto string) (ports.ProductFiscalFacts, error) {
	cfg, err := r.activeSource.Get(ctx, tenantID)
	if err == tenant_config.ErrUnknownActiveSource {
		return ports.ProductFiscalFacts{Found: false}, nil
	}
	if err != nil {
		return ports.ProductFiscalFacts{}, err
	}

	var grupoICMS *int
	var origprod *int
	var stRetidoEntrada *string
	var restituicaoUnit *string
	var fiscalDtRef *string
	err = r.pool.QueryRow(ctx, `
		SELECT grupo_icms, origprod, st_retido_entrada::text, restituicao_unit::text, fiscal_dt_ref::text
		FROM products_mirror
		WHERE tenant_id = $1 AND source = $2 AND codigo_produto = $3
	`, tenantID, string(cfg.Source), codigoProduto).Scan(&grupoICMS, &origprod, &stRetidoEntrada, &restituicaoUnit, &fiscalDtRef)
	if err == pgx.ErrNoRows {
		return ports.ProductFiscalFacts{Found: false}, nil
	}
	if err != nil {
		return ports.ProductFiscalFacts{}, err
	}

	return ports.ProductFiscalFacts{
		Found:           true,
		GrupoICMS:       grupoICMS,
		Origprod:        origprod,
		StRetidoEntrada: stRetidoEntrada,
		RestituicaoUnit: restituicaoUnit,
		FiscalDtRef:     fiscalDtRef,
	}, nil
}
