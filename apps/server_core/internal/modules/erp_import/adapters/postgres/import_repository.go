package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) PersistSnapshotAtomically(ctx context.Context, tenantID string, snapshot domain.ImportSnapshot) error {
	if snapshot.Status != domain.ImportStatusCompleted && snapshot.Status != domain.ImportStatusRejected {
		return fmt.Errorf("persist ERP import: unsupported status %q", snapshot.Status)
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin ERP import: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var locked bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock(hashtextextended('erp_import:'||$1, 0)) WHERE $1 = $1`, tenantID).Scan(&locked); err != nil {
		return fmt.Errorf("lock ERP import tenant: %w", err)
	}
	if !locked {
		return ports.ErrImportInProgress
	}

	var existingID domain.ImportID
	var existingProtocol domain.Protocol
	err = tx.QueryRow(ctx, `SELECT id, protocol FROM erp_import_protocols WHERE tenant_id=$1 AND file_sha256=$2`, tenantID, snapshot.FileSHA256).Scan(&existingID, &existingProtocol)
	if err == nil {
		return &ports.DuplicateFileError{ExistingID: existingID, ExistingProtocol: existingProtocol}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check duplicate ERP import: %w", err)
	}

	rejected, warnings := issueCounts(snapshot.Issues)
	// An unset source defaults to the ERP (xlsx) dataset: the strict M-01 path and
	// the pre-source integration fixtures leave it empty, and the source CHECK
	// constraint rejects the empty string. The lenient path sets catalogo_cliente.
	protocolSource := snapshot.Source
	if protocolSource == "" {
		protocolSource = domain.SourceXLSX
	}
	_, err = tx.Exec(ctx, `INSERT INTO erp_import_protocols (id,tenant_id,file_sha256,protocol,source,imported_at,status,accepted_count,rejected_count,warning_count) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, snapshot.ID, tenantID, snapshot.FileSHA256, snapshot.Protocol, protocolSource, snapshot.ImportedAt.UTC(), snapshot.Status, len(snapshot.AcceptedRows), rejected, warnings)
	if err != nil {
		return fmt.Errorf("insert ERP import protocol: %w", err)
	}
	for _, product := range snapshot.AcceptedRows {
		// Honest-unknown custo/stock_physical (client-catalog lenient path, ADR-17)
		// must land as SQL NULL, never a fabricated zero-valued string.
		custo := nullableString(string(product.Custo))
		precoVenda := nullableString(string(product.PrecoVenda))
		stockPhysical := nullableString(product.StockPhysical)
		_, err = tx.Exec(ctx, `INSERT INTO erp_import_products (tenant_id,protocol_id,codprod,descrprod,custo,preco_venda,stock_physical,stock_reserved,ean,refforn,marca,ncm,grupo,descrgrupo) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`, tenantID, snapshot.ID, product.Codprod, product.Descrprod, custo, precoVenda, stockPhysical, product.StockReserved, product.EAN, product.Refforn, product.Marca, product.NCM, product.Grupo, product.DescrGrupo)
		if err != nil {
			return fmt.Errorf("insert ERP import product: %w", err)
		}
	}
	if snapshot.Status == domain.ImportStatusCompleted {
		if _, err := r.mergeSnapshotTx(ctx, tx, tenantID, protocolSource, snapshot.ID, snapshot.AcceptedRows); err != nil {
			return fmt.Errorf("merge ERP import mirror: %w", err)
		}
	}
	for _, issue := range snapshot.Issues {
		_, err = tx.Exec(ctx, `INSERT INTO erp_import_issues (tenant_id,protocol_id,row_number,column_name,kind,code,detail,offending_value) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, tenantID, snapshot.ID, issue.Row, issue.Column, issue.Kind, issue.Code, issue.Detail, issue.OffendingValue)
		if err != nil {
			return fmt.Errorf("insert ERP import issue: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit ERP import: %w", err)
	}
	return nil
}

// nullableString maps an empty (honest-unknown) value to SQL NULL rather
// than persisting a fabricated empty string into a numeric column (ADR-17).
func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// issueCounts reports the number of distinct rejected rows (a row may carry
// several rejection issues) and the number of warning issues.
func issueCounts(issues []domain.Issue) (rejected, warnings int) {
	rejectedRows := make(map[int]struct{})
	for _, issue := range issues {
		switch issue.Kind {
		case domain.Rejection:
			rejectedRows[issue.Row] = struct{}{}
		case domain.Warning:
			warnings++
		}
	}
	return len(rejectedRows), warnings
}
