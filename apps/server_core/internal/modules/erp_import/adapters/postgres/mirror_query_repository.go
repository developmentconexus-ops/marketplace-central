package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

const mirrorProductColumns = `m.codigo_produto,m.descricao,m.referencia,m.ean,m.marca,m.grupo_codigo,m.grupo_descricao,m.ncm,m.custo::text,m.preco_venda::text,m.estoque_total::text,m.updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMirrorProduct(row rowScanner) (domain.MirrorProduct, error) {
	var product domain.MirrorProduct
	var descricao, referencia, ean, marca, grupoCodigo, grupoDescricao, ncm sql.NullString
	var custo, precoVenda, estoqueTotal sql.NullString
	var importedAt sql.NullTime
	if err := row.Scan(
		&product.CodigoProduto, &descricao, &referencia, &ean, &marca, &grupoCodigo,
		&grupoDescricao, &ncm, &custo, &precoVenda, &estoqueTotal, &product.UpdatedAt, &importedAt,
	); err != nil {
		return domain.MirrorProduct{}, err
	}
	product.Descricao = nullStringPointer(descricao)
	product.Referencia = nullStringPointer(referencia)
	product.EAN = nullStringPointer(ean)
	product.Marca = nullStringPointer(marca)
	product.GrupoCodigo = nullStringPointer(grupoCodigo)
	product.GrupoDescricao = nullStringPointer(grupoDescricao)
	product.NCM = nullStringPointer(ncm)
	product.Custo = nullStringPointer(custo)
	product.PrecoVenda = nullStringPointer(precoVenda)
	product.EstoqueTotal = nullStringPointer(estoqueTotal)
	if importedAt.Valid {
		product.ImportedAt = &importedAt.Time
	}
	return product, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func (r *Repository) MirrorRows(ctx context.Context, tenantID string, source domain.ImportSource) ([]domain.MirrorProduct, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+mirrorProductColumns+`,p.imported_at
FROM products_mirror m
LEFT JOIN erp_import_protocols p ON p.tenant_id = m.tenant_id AND p.id = m.protocol_id
WHERE m.tenant_id=$1 AND m.source=$2 AND m.absent_in_last_snapshot=false
ORDER BY m.codigo_produto`, tenantID, source)
	if err != nil {
		return nil, fmt.Errorf("query current mirror rows: %w", err)
	}
	defer rows.Close()

	products := make([]domain.MirrorProduct, 0)
	for rows.Next() {
		product, err := scanMirrorProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan current mirror row: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate current mirror rows: %w", err)
	}
	return products, nil
}

func (r *Repository) MirrorProductByCode(ctx context.Context, tenantID string, source domain.ImportSource, codigo string) (domain.MirrorProduct, bool, error) {
	product, err := scanMirrorProduct(r.pool.QueryRow(ctx, `SELECT `+mirrorProductColumns+`,p.imported_at
FROM products_mirror m
LEFT JOIN erp_import_protocols p ON p.tenant_id = m.tenant_id AND p.id = m.protocol_id
WHERE m.tenant_id=$1 AND m.source=$2 AND m.codigo_produto=$3 AND m.absent_in_last_snapshot=false`, tenantID, source, codigo))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.MirrorProduct{}, false, nil
	}
	if err != nil {
		return domain.MirrorProduct{}, false, fmt.Errorf("query current mirror product by code: %w", err)
	}
	return product, true, nil
}

// MirrorProductsByCodes reads the mirror rows for a set of product codes in a
// single round-trip. A code with no row is absent from the result — the caller
// decides whether that is an unknown product or an honest "not in this source",
// and never reads it as a zero (ADR-17).
func (r *Repository) MirrorProductsByCodes(ctx context.Context, tenantID string, source domain.ImportSource, codigos []string) ([]domain.MirrorProduct, error) {
	if len(codigos) == 0 {
		return []domain.MirrorProduct{}, nil
	}
	rows, err := r.pool.Query(ctx, `SELECT `+mirrorProductColumns+`,p.imported_at
FROM products_mirror m
LEFT JOIN erp_import_protocols p ON p.tenant_id = m.tenant_id AND p.id = m.protocol_id
WHERE m.tenant_id=$1 AND m.source=$2 AND m.absent_in_last_snapshot=false
  AND m.codigo_produto = ANY($3)`, tenantID, source, codigos)
	if err != nil {
		return nil, fmt.Errorf("query mirror products by codes: %w", err)
	}
	defer rows.Close()

	products := make([]domain.MirrorProduct, 0, len(codigos))
	for rows.Next() {
		product, err := scanMirrorProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan mirror product by code: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate mirror products by codes: %w", err)
	}
	return products, nil
}

// MirrorCatalogPage returns up to limit+1 rows when limit is positive so callers can detect a following page.
//
// The search term matches description, reference, product code and EAN — the
// four identifiers an operator actually has at hand — and both sides of the
// text comparison are unaccented (migration 0080), because the ERP stores
// "VÁLVULA" while the operator types "valvula".
func (r *Repository) MirrorCatalogPage(ctx context.Context, tenantID string, source domain.ImportSource, query string, afterInternalID int64, limit int) ([]domain.MirrorProduct, error) {
	rows, err := r.pool.Query(ctx, `SELECT DISTINCT ON (m.codigo_produto::bigint)
  `+mirrorProductColumns+`,p.imported_at
FROM products_mirror m
LEFT JOIN erp_import_protocols p ON p.tenant_id = m.tenant_id AND p.id = m.protocol_id
WHERE m.tenant_id=$1 AND m.source=$2 AND m.absent_in_last_snapshot=false
  AND m.codigo_produto ~ '^[0-9]{1,18}$'
  AND m.codigo_produto::bigint > $3
  AND ($4='' OR unaccent(m.descricao) ILIKE '%'||unaccent($4)||'%'
       OR unaccent(coalesce(m.referencia,'')) ILIKE '%'||unaccent($4)||'%'
       OR m.codigo_produto ILIKE '%'||$4||'%'
       OR coalesce(m.ean,'') ILIKE '%'||$4||'%')
ORDER BY m.codigo_produto::bigint ASC, m.updated_at DESC
LIMIT CASE WHEN $5 > 0 THEN $5 + 1 END`, tenantID, source, afterInternalID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query mirror catalog page: %w", err)
	}
	defer rows.Close()

	products := make([]domain.MirrorProduct, 0)
	for rows.Next() {
		product, err := scanMirrorProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan mirror catalog row: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate mirror catalog page: %w", err)
	}
	return products, nil
}

func (r *Repository) MirrorEANCollisionCounts(ctx context.Context, tenantID string, source domain.ImportSource) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT ean,count(*) FROM products_mirror WHERE tenant_id=$1 AND source=$2 AND absent_in_last_snapshot=false AND ean IS NOT NULL GROUP BY ean HAVING count(*) >= 2`, tenantID, source)
	if err != nil {
		return nil, fmt.Errorf("query mirror EAN collision counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var ean string
		var count int
		if err := rows.Scan(&ean, &count); err != nil {
			return nil, fmt.Errorf("scan mirror EAN collision count: %w", err)
		}
		counts[ean] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate mirror EAN collision counts: %w", err)
	}
	return counts, nil
}
