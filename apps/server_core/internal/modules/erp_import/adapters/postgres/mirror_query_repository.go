package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

const mirrorProductColumns = `codigo_produto,descricao,referencia,ean,marca,grupo_codigo,grupo_descricao,ncm,custo::text,preco_venda::text,estoque_total::text,updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMirrorProduct(row rowScanner) (domain.MirrorProduct, error) {
	var product domain.MirrorProduct
	var descricao, referencia, ean, marca, grupoCodigo, grupoDescricao, ncm sql.NullString
	var custo, precoVenda, estoqueTotal sql.NullString
	if err := row.Scan(
		&product.CodigoProduto, &descricao, &referencia, &ean, &marca, &grupoCodigo,
		&grupoDescricao, &ncm, &custo, &precoVenda, &estoqueTotal, &product.UpdatedAt,
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
	return product, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func (r *Repository) MirrorRows(ctx context.Context, tenantID string, source domain.ImportSource) ([]domain.MirrorProduct, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+mirrorProductColumns+` FROM products_mirror WHERE tenant_id=$1 AND source=$2 AND absent_in_last_snapshot=false ORDER BY codigo_produto`, tenantID, source)
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
	product, err := scanMirrorProduct(r.pool.QueryRow(ctx, `SELECT `+mirrorProductColumns+` FROM products_mirror WHERE tenant_id=$1 AND source=$2 AND codigo_produto=$3 AND absent_in_last_snapshot=false`, tenantID, source, codigo))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.MirrorProduct{}, false, nil
	}
	if err != nil {
		return domain.MirrorProduct{}, false, fmt.Errorf("query current mirror product by code: %w", err)
	}
	return product, true, nil
}

// MirrorCatalogPage returns up to limit+1 rows when limit is positive so callers can detect a following page.
func (r *Repository) MirrorCatalogPage(ctx context.Context, tenantID string, source domain.ImportSource, query string, afterInternalID int64, limit int) ([]domain.MirrorProduct, error) {
	rows, err := r.pool.Query(ctx, `SELECT codigo_produto,descricao,referencia,ean,marca,grupo_codigo,grupo_descricao,ncm,custo::text,preco_venda::text,estoque_total::text,updated_at
FROM products_mirror
WHERE tenant_id=$1
  AND source=$2
  AND absent_in_last_snapshot=false
  AND codigo_produto ~ '^[0-9]{1,18}$'
  AND codigo_produto::bigint > $3
  AND ($4='' OR descricao ILIKE '%'||$4||'%')
ORDER BY codigo_produto::bigint ASC
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
