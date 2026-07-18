package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// CostReader resolves a product's custo_erp as decimal Money for the engine.
// A nil Money means the cost is absent — ADR-17: unknown is never zero, so a
// decomposition with a nil custo propagates an unknown component rather than
// computing against 0. Implemented by the costread adapter over IC-02
// internal_read.
type CostReader interface {
	CostForProduct(ctx context.Context, productID int, asOf time.Time) (*domain.Money, error)
}
