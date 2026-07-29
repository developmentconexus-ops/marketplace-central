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

// GetImportChain is discovered by type assertion at composition time, so a
// signature drift here would surface only as a missing capability at request
// time. Assert it at build time instead.
var _ ports.ImportChainRepository = (*Repository)(nil)

func (r *Repository) FindByFileSHA256(ctx context.Context, tenantID string, hash domain.FileSHA256) (*domain.ImportReport, error) {
	report, err := r.report(ctx, tenantID, `SELECT id,protocol,file_sha256,source,imported_at,status,accepted_count,rejected_count,warning_count FROM erp_import_protocols WHERE tenant_id=$1 AND file_sha256=$2`, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Repository) ListImports(ctx context.Context, tenantID string) ([]domain.ImportReport, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,protocol,file_sha256,source,imported_at,status,accepted_count,rejected_count,warning_count FROM erp_import_protocols WHERE tenant_id=$1 ORDER BY imported_at DESC, id DESC`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list ERP imports: %w", err)
	}
	defer rows.Close()
	out := make([]domain.ImportReport, 0)
	for rows.Next() {
		var v domain.ImportReport
		if err := scanReport(rows, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *Repository) GetImport(ctx context.Context, tenantID string, importID domain.ImportID) (domain.ImportReport, error) {
	report, err := r.report(ctx, tenantID, `SELECT id,protocol,file_sha256,source,imported_at,status,accepted_count,rejected_count,warning_count FROM erp_import_protocols WHERE tenant_id=$1 AND id=$2`, importID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ImportReport{}, ports.ErrImportNotFound
	}
	if err != nil {
		return domain.ImportReport{}, err
	}
	rows, err := r.pool.Query(ctx, `SELECT row_number,column_name,kind,code,detail,offending_value FROM erp_import_issues WHERE tenant_id=$1 AND protocol_id=$2 ORDER BY row_number`, tenantID, importID)
	if err != nil {
		return domain.ImportReport{}, err
	}
	defer rows.Close()
	report.Issues = make([]domain.Issue, 0)
	for rows.Next() {
		var i domain.Issue
		if err := rows.Scan(&i.Row, &i.Column, &i.Kind, &i.Code, &i.Detail, &i.OffendingValue); err != nil {
			return domain.ImportReport{}, err
		}
		report.Issues = append(report.Issues, i)
	}
	return report, rows.Err()
}

func (r *Repository) GetImportChain(ctx context.Context, tenantID string, importID domain.ImportID) (domain.ImportChain, error) {
	const query = `
		WITH import_target AS (
			SELECT eip.id AS import_id, eip.protocol AS protocol
			FROM erp_import_protocols AS eip
			WHERE eip.tenant_id = $1
			  AND eip.id = $2
		),
		import_products AS (
			SELECT products.codprod AS codprod
			FROM erp_import_products AS products
			JOIN import_target AS target ON target.import_id = products.protocol_id
			WHERE products.tenant_id = $1
		),
		resolved_products AS (
			-- The two sides reach this join through different pipelines:
			-- erp_import_products.codprod keeps the raw spreadsheet string (TEXT,
			-- 0046, and its primary key lets '101' and '00101' coexist as two
			-- rows), while product_links.internal_product_id is a ParseInt'd
			-- integer column (0025) — the only identity a link carries for an ERP
			-- product, since internal_reference_code holds refforn and not codprod
			-- (product_links/application.newCandidate). A raw '00101' therefore
			-- linked as 101 and compared '101' = '00101' → false.
			--
			-- ltrim on BOTH sides fixed that by making '101' and '00101' one key,
			-- which is a LOSING canonicalization: an import carrying both rows then
			-- counted the single linked product TWICE. Exact numeric comparison
			-- instead — a codprod names internal product N when, read as a number,
			-- it IS N — with the counted key moved to links.internal_product_id,
			-- because two import rows naming the same internal product are one
			-- linked product, and DISTINCT on the codprod text would still yield 2.
			--
			-- CASE, not "codprod ~ '^[0-9]+$' AND codprod::bigint = id": a join
			-- predicate carries no evaluation-order guarantee, so the planner is
			-- free to cast a non-numeric codprod before the filter rejects it and
			-- raise 22P02 at query time. CASE evaluates WHEN before THEN by
			-- definition of the language. The digit bound keeps an oversized
			-- codprod out of the cast (22003) and is measured after the padding is
			-- stripped, so only significant digits count; 18 covers every value an
			-- integer column can hold.
			SELECT DISTINCT links.internal_product_id AS internal_product_id
			FROM import_products AS products
			JOIN product_links AS links
			  ON links.internal_product_id = CASE
			         WHEN products.codprod ~ '^[0-9]+$'
			          AND length(ltrim(products.codprod, '0')) <= 18
			         THEN products.codprod::bigint
			     END
			WHERE links.tenant_id = $1
			  AND links.state = 'resolved'
		),
		queued_products AS (
			-- Both sides of this join are the SAME raw string out of the same
			-- slice, established by reading the producer chain: import_service
			-- builds its codes slice from the accepted rows' Codprod and passes it to
			-- MarketEnqueuer.EnqueueMarketProducts, which hands it unchanged to
			-- AppendPendingCodigos → cursor->'pending'; those same accepted rows
			-- are what CopyFrom writes as erp_import_products.codprod. No other
			-- writer exists — AppendPendingCodigo (singular) has zero callers and
			-- ProductsCursor has no pending field. Raw equality was therefore
			-- already the EXACT comparison, and ltrim on both sides only bought
			-- collisions: 'ABC' and '00ABC' are two import rows, and a single
			-- pending "ABC" counted both. The counted key is the raw codprod,
			-- which is the row identity importados counts.
			--
			-- COALESCE only defends against NULL. A cursor whose 'pending' is an
			-- object or a scalar makes jsonb_array_elements_text raise at query
			-- time and takes the whole endpoint down for the tenant, so the type
			-- is checked before expansion. CASE (not AND) because a join/filter
			-- predicate gives no evaluation-order guarantee.
			SELECT DISTINCT products.codprod AS codprod
			FROM sync_state AS state
			CROSS JOIN LATERAL jsonb_array_elements_text(
				CASE
					WHEN jsonb_typeof(state.cursor -> 'pending') = 'array'
						THEN state.cursor -> 'pending'
					ELSE '[]'::jsonb
				END
			) AS pending(codprod)
			JOIN import_products AS products
			  ON products.codprod = pending.codprod
			WHERE state.tenant_id = $1
			  AND state.entity = 'market'
		)
		-- importados and enfileirados count import ROWS; vinculados counts linked
		-- internal PRODUCTS. All three differ without any duplicate: a row may
		-- resolve to no product, and only part of the import stays pending.
		SELECT target.protocol AS protocol,
		       (SELECT count(*) FROM import_products) AS importados,
		       (SELECT count(*) FROM resolved_products) AS vinculados,
		       (SELECT count(*) FROM queued_products) AS enfileirados,
		       statement_timestamp() AS queue_read_at
		FROM import_target AS target
	`

	var chain domain.ImportChain
	err := r.pool.QueryRow(ctx, query, tenantID, importID).Scan(
		&chain.Protocol,
		&chain.Importados,
		&chain.Vinculados,
		&chain.Enfileirados,
		&chain.QueueReadAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ImportChain{}, ports.ErrImportNotFound
	}
	if err != nil {
		return domain.ImportChain{}, fmt.Errorf("get ERP import chain: %w", err)
	}
	return chain, nil
}

func (r *Repository) LatestCompletedSnapshot(ctx context.Context, tenantID string, source domain.ImportSource) (domain.ImportSnapshot, error) {
	var s domain.ImportSnapshot
	// Filter by the active dataset source (two-source toggle): xlsx = the Sankhya
	// ERP snapshot, catalogo_cliente = a lenient prospect catalog. Without this
	// filter the newest imported_at blindly wins, so importing a prospect catalog
	// would hijack the active dataset and displace the real ERP snapshot.
	err := r.pool.QueryRow(ctx, `SELECT id,protocol,file_sha256,source,imported_at,status FROM erp_import_protocols WHERE tenant_id=$1 AND source=$2 AND status='COMPLETED' ORDER BY imported_at DESC, id DESC LIMIT 1`, tenantID, source).Scan(&s.ID, &s.Protocol, &s.FileSHA256, &s.Source, &s.ImportedAt, &s.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ImportSnapshot{}, ports.ErrImportNotFound
	}
	if err != nil {
		return domain.ImportSnapshot{}, err
	}
	rows, err := r.pool.Query(ctx, `SELECT codprod,descrprod,custo::text,preco_venda::text,stock_physical,stock_reserved,ean,refforn,marca,ncm,grupo,descrgrupo FROM erp_import_products WHERE tenant_id=$1 AND protocol_id=$2 ORDER BY codprod`, tenantID, s.ID)
	if err != nil {
		return domain.ImportSnapshot{}, err
	}
	defer rows.Close()
	s.AcceptedRows = make([]domain.NormalizedRow, 0)
	for rows.Next() {
		var p domain.NormalizedRow
		// custo / preco_venda / stock_physical are nullable (lenient client-catalog
		// imports omit them, and the sale price is optional on every path, ADR-17).
		// Scan through nullable intermediates and leave the textual fields empty
		// when NULL — empty is the honest-unknown the readers expect.
		var custo, precoVenda, stockPhysical sql.NullString
		if err := rows.Scan(&p.Codprod, &p.Descrprod, &custo, &precoVenda, &stockPhysical, &p.StockReserved, &p.EAN, &p.Refforn, &p.Marca, &p.NCM, &p.Grupo, &p.DescrGrupo); err != nil {
			return domain.ImportSnapshot{}, err
		}
		if custo.Valid {
			p.Custo = domain.Decimal(custo.String)
		}
		if precoVenda.Valid {
			p.PrecoVenda = domain.Decimal(precoVenda.String)
		}
		if stockPhysical.Valid {
			p.StockPhysical = stockPhysical.String
		}
		s.AcceptedRows = append(s.AcceptedRows, p)
	}
	if err := rows.Err(); err != nil {
		return domain.ImportSnapshot{}, err
	}
	rows.Close()
	issueRows, err := r.pool.Query(ctx, `SELECT row_number,column_name,kind,code,detail,offending_value FROM erp_import_issues WHERE tenant_id=$1 AND protocol_id=$2 ORDER BY row_number`, tenantID, s.ID)
	if err != nil {
		return domain.ImportSnapshot{}, err
	}
	defer issueRows.Close()
	s.Issues = make([]domain.Issue, 0)
	for issueRows.Next() {
		var issue domain.Issue
		if err := issueRows.Scan(&issue.Row, &issue.Column, &issue.Kind, &issue.Code, &issue.Detail, &issue.OffendingValue); err != nil {
			return domain.ImportSnapshot{}, err
		}
		s.Issues = append(s.Issues, issue)
	}
	return s, issueRows.Err()
}

func (r *Repository) report(ctx context.Context, tenantID, query string, arg any) (domain.ImportReport, error) {
	var v domain.ImportReport
	err := scanReport(r.pool.QueryRow(ctx, query, tenantID, arg), &v)
	return v, err
}
func scanReport(row interface{ Scan(...any) error }, v *domain.ImportReport) error {
	return row.Scan(&v.ID, &v.Protocol, &v.FileSHA256, &v.Source, &v.ImportedAt, &v.Status, &v.AcceptedCount, &v.RejectedCount, &v.WarningCount)
}
