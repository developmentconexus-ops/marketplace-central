package ports

import (
	"context"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type InternalFactReader interface {
	GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error)
	GetTaxInputs(ctx context.Context, productID int, effectiveAt time.Time, source internalreaddomain.TaxSourceIdentity) (internalreaddomain.TaxInputs, error)
}

type BatchInternalFactReader interface {
	InternalFactReader
	GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.CostAsOf, error)
	GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.TaxInputs, error)
}
