package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type MarginInputStore interface {
	ReplaceInputs(ctx context.Context, installationID string, orders []string, inputs []domain.MarginInput) error
	ListInputs(ctx context.Context, installationID string, limit int) ([]domain.MarginInput, error)
}

type ManualAdjustmentStore interface {
	CreateOrGetAdjustment(ctx context.Context, adjustment domain.ManualAdjustment) (domain.ManualAdjustment, error)
	ListAdjustments(ctx context.Context, installationID string, limit int) ([]domain.ManualAdjustment, error)
}

type ProfitSnapshotStore interface {
	ReplaceSnapshots(ctx context.Context, installationID string, orders []string, snapshots []domain.ProfitSnapshot) error
	ListSnapshots(ctx context.Context, installationID string, limit int) ([]domain.ProfitSnapshot, error)
}
