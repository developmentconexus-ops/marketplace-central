package ports

import (
	"context"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

// CostReader is the orders-local, module-scoped consumer port for the
// internal_read cost fact. It reuses the exact signature established by
// profitability's InternalFactReader.GetCostAsOf — see
// apps/server_core/internal/modules/profitability/ports/internal_read.go:10-13.
type CostReader interface {
	GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error)
}
