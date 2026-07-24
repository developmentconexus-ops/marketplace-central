package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

func (r *Repository) SyncLatestCompletedSnapshot(ctx context.Context, tenantID string, source domain.ImportSource) (processed int, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin latest ERP mirror sync: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var protocolID domain.ImportID
	err = tx.QueryRow(ctx, `SELECT id FROM erp_import_protocols WHERE tenant_id=$1 AND source=$2 AND status='COMPLETED' ORDER BY imported_at DESC, id DESC LIMIT 1`, tenantID, source).Scan(&protocolID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrImportNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("load latest completed ERP protocol: %w", err)
	}

	productRows, err := tx.Query(ctx, `SELECT codprod,descrprod,custo::text,stock_physical,stock_reserved,ean,refforn,marca,ncm,grupo,descrgrupo FROM erp_import_products WHERE tenant_id=$1 AND protocol_id=$2 ORDER BY codprod`, tenantID, protocolID)
	if err != nil {
		return 0, fmt.Errorf("load latest completed ERP products: %w", err)
	}
	rows := make([]domain.NormalizedRow, 0)
	for productRows.Next() {
		var row domain.NormalizedRow
		var custo, stockPhysical sql.NullString
		if err := productRows.Scan(&row.Codprod, &row.Descrprod, &custo, &stockPhysical, &row.StockReserved, &row.EAN, &row.Refforn, &row.Marca, &row.NCM, &row.Grupo, &row.DescrGrupo); err != nil {
			productRows.Close()
			return 0, fmt.Errorf("scan latest completed ERP product: %w", err)
		}
		if custo.Valid {
			row.Custo = domain.Decimal(custo.String)
		}
		if stockPhysical.Valid {
			row.StockPhysical = stockPhysical.String
		}
		rows = append(rows, row)
	}
	if err := productRows.Err(); err != nil {
		productRows.Close()
		return 0, fmt.Errorf("iterate latest completed ERP products: %w", err)
	}
	productRows.Close()

	processed, err = r.mergeSnapshotTx(ctx, tx, tenantID, source, protocolID, rows)
	if err != nil {
		return 0, fmt.Errorf("merge latest completed ERP snapshot: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit latest ERP mirror sync: %w", err)
	}
	return processed, nil
}

func (r *Repository) mergeSnapshotTx(ctx context.Context, tx pgx.Tx, tenantID string, source domain.ImportSource, protocolID domain.ImportID, rows []domain.NormalizedRow) (int, error) {
	_, err := tx.Exec(ctx, `CREATE TEMP TABLE erp_import_mirror_stage (
		codigo_produto text NOT NULL,
		descricao text,
		referencia text,
		ean text,
		marca text,
		grupo_codigo text,
		grupo_descricao text,
		ncm text,
		custo text,
		preco_venda text,
		stock_physical text,
		stock_reserved text,
		local_codigo text
	) ON COMMIT DROP`)
	if err != nil {
		return 0, err
	}

	copyRows := make([][]any, 0, len(rows))
	for _, row := range rows {
		copyRows = append(copyRows, []any{
			row.Codprod, row.Descrprod, row.Refforn, row.EAN, row.Marca, row.Grupo,
			row.DescrGrupo, row.NCM, nullableString(string(row.Custo)),
			nullableString(string(row.PrecoVenda)), nullableString(row.StockPhysical),
			row.StockReserved, row.Local,
		})
	}
	_, err = tx.CopyFrom(ctx, pgx.Identifier{"erp_import_mirror_stage"}, []string{
		"codigo_produto", "descricao", "referencia", "ean", "marca", "grupo_codigo",
		"grupo_descricao", "ncm", "custo", "preco_venda", "stock_physical",
		"stock_reserved", "local_codigo",
	}, pgx.CopyFromRows(copyRows))
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO products_mirror (
		tenant_id,source,codigo_produto,descricao,referencia,ean,marca,grupo_codigo,
		grupo_descricao,ncm,custo,preco_venda,estoque_total,protocol_id
	)
	SELECT $1,$2,codigo_produto,descricao,referencia,ean,marca,grupo_codigo,
		grupo_descricao,ncm,custo::numeric,preco_venda::numeric,
		CASE WHEN stock_physical IS NOT NULL AND stock_reserved IS NOT NULL
			THEN stock_physical::numeric-stock_reserved::numeric END,$3
	FROM erp_import_mirror_stage
	ON CONFLICT (tenant_id,codigo_produto) DO UPDATE SET
		source=EXCLUDED.source,descricao=EXCLUDED.descricao,referencia=EXCLUDED.referencia,
		ean=EXCLUDED.ean,marca=EXCLUDED.marca,grupo_codigo=EXCLUDED.grupo_codigo,
		grupo_descricao=EXCLUDED.grupo_descricao,ncm=EXCLUDED.ncm,custo=EXCLUDED.custo,
		preco_venda=EXCLUDED.preco_venda,estoque_total=EXCLUDED.estoque_total,
		protocol_id=EXCLUDED.protocol_id,absent_in_last_snapshot=false,stale_since=NULL,
		updated_at=now()`, tenantID, source, protocolID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `UPDATE products_mirror AS mirror
	SET absent_in_last_snapshot=true,stale_since=COALESCE(stale_since,now()),updated_at=now()
	WHERE mirror.tenant_id=$1 AND mirror.source=$2 AND mirror.absent_in_last_snapshot=false
		AND NOT EXISTS (SELECT 1 FROM erp_import_mirror_stage AS stage WHERE stage.codigo_produto=mirror.codigo_produto)`, tenantID, source)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `DELETE FROM products_mirror_stock_locations AS locations
	WHERE locations.tenant_id=$1 AND EXISTS (
		SELECT 1 FROM erp_import_mirror_stage AS stage WHERE stage.codigo_produto=locations.codigo_produto
	)`, tenantID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO products_mirror_stock_locations (
		tenant_id,codigo_produto,local_codigo,local_descricao,quantidade
	)
	SELECT $1,codigo_produto,local_codigo,NULL,
		CASE WHEN stock_physical IS NOT NULL AND stock_reserved IS NOT NULL
			THEN stock_physical::numeric-stock_reserved::numeric END
	FROM erp_import_mirror_stage WHERE local_codigo IS NOT NULL AND local_codigo<>''
	ON CONFLICT (tenant_id,codigo_produto,local_codigo) DO UPDATE SET
		local_descricao=EXCLUDED.local_descricao,quantidade=EXCLUDED.quantidade`, tenantID)
	if err != nil {
		return 0, err
	}
	return len(rows), nil
}
