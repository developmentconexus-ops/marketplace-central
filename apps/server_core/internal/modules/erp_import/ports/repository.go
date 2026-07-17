package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

type ImportRepository interface {
	PersistSnapshotAtomically(ctx context.Context, tenantID string, snapshot domain.ImportSnapshot) error
	FindByFileSHA256(ctx context.Context, tenantID string, fileSHA256 domain.FileSHA256) (*domain.ImportReport, error)
	ListImports(ctx context.Context, tenantID string) ([]domain.ImportReport, error)
	GetImport(ctx context.Context, tenantID string, importID domain.ImportID) (domain.ImportReport, error)
	LatestCompletedSnapshot(ctx context.Context, tenantID string) (domain.ImportSnapshot, error)
}
