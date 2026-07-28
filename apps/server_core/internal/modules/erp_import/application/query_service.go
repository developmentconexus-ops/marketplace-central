package application

import (
	"context"
	"errors"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

// ErrImportChainUnavailable is returned when the repository handed to
// NewQueryService cannot read import chains. GetImportChain is an optional
// capability discovered by type assertion, so a decorator that wraps the
// repository without re-exposing ImportChainRepository silently erases it. This
// error makes that erasure say so instead of panicking on a nil interface.
var ErrImportChainUnavailable = errors.New("erp import repository does not read import chains")

type QueryService struct {
	repo      ports.ImportRepository
	chainRepo ports.ImportChainRepository
	tenantID  string
}

func NewQueryService(repo ports.ImportRepository, tenantID string) *QueryService {
	chainRepo, _ := repo.(ports.ImportChainRepository)
	return &QueryService{repo: repo, chainRepo: chainRepo, tenantID: tenantID}
}

func (q *QueryService) ListImports(ctx context.Context) ([]domain.ImportReport, error) {
	return q.repo.ListImports(ctx, q.tenantID)
}

func (q *QueryService) GetImport(ctx context.Context, id domain.ImportID) (domain.ImportReport, error) {
	return q.repo.GetImport(ctx, q.tenantID, id)
}

func (q *QueryService) GetImportChain(ctx context.Context, id domain.ImportID) (domain.ImportChain, error) {
	if q.chainRepo == nil {
		return domain.ImportChain{}, ErrImportChainUnavailable
	}
	return q.chainRepo.GetImportChain(ctx, q.tenantID, id)
}
