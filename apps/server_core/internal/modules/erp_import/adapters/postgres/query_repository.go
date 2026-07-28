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
			-- erp_import_products.codprod keeps the raw spreadsheet string (a human
			-- types it, and IsValidCodprod accepts leading zeros), while
			-- product_links.internal_product_id is a ParseInt'd integer column. A
			-- raw '00101' therefore linked as 101 and compared '101' = '00101' →
			-- false, silently undercounting an operator-facing number. Strip the
			-- padding on BOTH sides so the comparison is one canonical form.
			-- Text-to-text on purpose: codprod is not guaranteed numeric, and a
			-- ::integer cast on junk would raise 22P02 at query time.
			SELECT DISTINCT products.codprod AS codprod
			FROM import_products AS products
			JOIN product_links AS links
			  ON ltrim(links.internal_product_id::text, '0') = ltrim(products.codprod, '0')
			WHERE links.tenant_id = $1
			  AND links.state = 'resolved'
		),
		queued_products AS (
			-- importados / vinculados / enfileirados are read on one screen as a
			-- decomposition of a single population, so the two joined counters have
			-- to agree on what makes two CODPRODs the same product. resolved_products
			-- already answers that canonically; leaving this side raw would let one
			-- padded CODPROD be linked but not queued, and the operator would read
			-- the gap between two numbers as a stalled queue that never existed. The
			-- counted key stays the raw codprod, which is what importados counts —
			-- only the identity test is canonicalized. Text-to-text for the same
			-- reason as above: codprod is not guaranteed numeric.
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
			  ON ltrim(products.codprod, '0') = ltrim(pending.codprod, '0')
			WHERE state.tenant_id = $1
			  AND state.entity = 'market'
		)
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
