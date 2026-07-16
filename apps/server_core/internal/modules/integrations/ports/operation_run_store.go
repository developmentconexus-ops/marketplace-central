package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type OperationRunStore interface {
	BeginExclusive(ctx context.Context, run domain.OperationRun) (active domain.OperationRun, alreadyActive bool, err error)
	SaveOperationRun(ctx context.Context, run domain.OperationRun) error
	ListByInstallation(ctx context.Context, installationID string) ([]domain.OperationRun, error)
}

type OperationRunReadStore interface {
	ListRuns(ctx context.Context, query RunListQuery) (RunPage, error)
}
