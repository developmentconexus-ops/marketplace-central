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
