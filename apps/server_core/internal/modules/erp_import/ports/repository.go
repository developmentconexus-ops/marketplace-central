package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

type ImportChainRepository interface {
	GetImportChain(ctx context.Context, tenantID string, importID domain.ImportID) (domain.ImportChain, error)
}

type ImportRepository interface {
	PersistSnapshotAtomically(ctx context.Context, tenantID string, snapshot domain.ImportSnapshot) error
	FindByFileSHA256(ctx context.Context, tenantID string, fileSHA256 domain.FileSHA256) (*domain.ImportReport, error)
	ListImports(ctx context.Context, tenantID string) ([]domain.ImportReport, error)
	GetImport(ctx context.Context, tenantID string, importID domain.ImportID) (domain.ImportReport, error)
	LatestCompletedSnapshot(ctx context.Context, tenantID string, source domain.ImportSource) (domain.ImportSnapshot, error)
	SyncLatestCompletedSnapshot(ctx context.Context, tenantID string, source domain.ImportSource) (processed int, err error)
	MirrorRows(ctx context.Context, tenantID string, source domain.ImportSource) ([]domain.MirrorProduct, error)
	MirrorProductByCode(ctx context.Context, tenantID string, source domain.ImportSource, codigo string) (domain.MirrorProduct, bool, error)
	// MirrorCatalogPage returns up to limit+1 rows when limit is positive so callers can detect a following page.
	MirrorCatalogPage(ctx context.Context, tenantID string, source domain.ImportSource, query string, afterInternalID int64, limit int) ([]domain.MirrorProduct, error)
	MirrorEANCollisionCounts(ctx context.Context, tenantID string, source domain.ImportSource) (map[string]int, error)
	// MirrorProductsByCodes answers one batch stock/price question per page of
	// listings; codes with no mirror row are simply absent from the result.
	MirrorProductsByCodes(ctx context.Context, tenantID string, source domain.ImportSource, codigos []string) ([]domain.MirrorProduct, error)
}
