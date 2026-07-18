package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// CostReader is the market module's minimal consumer view of ERP cost (IC-02).
// The production adapter over internal_read.GetCostAsOf is wired at the composition root by the hub post-merge.
type CostReader interface {
	GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (*domain.Money, time.Time, error)
}
