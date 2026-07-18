package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

type QueryService struct {
	repo     ports.ImportRepository
	tenantID string
}

func NewQueryService(repo ports.ImportRepository, tenantID string) *QueryService {
	return &QueryService{repo: repo, tenantID: tenantID}
}

func (q *QueryService) ListImports(ctx context.Context) ([]domain.ImportReport, error) {
	return q.repo.ListImports(ctx, q.tenantID)
}

func (q *QueryService) GetImport(ctx context.Context, id domain.ImportID) (domain.ImportReport, error) {
	return q.repo.GetImport(ctx, q.tenantID, id)
}
