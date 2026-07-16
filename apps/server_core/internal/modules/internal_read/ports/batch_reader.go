package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

// BatchReader is the bounded, product-keyed Oracle read surface used by batch
// application jobs. A nil map value is an explicit unknown fact.
type BatchReader interface {
	GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.CostAsOf, error)
	GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.TaxInputs, error)
	GetICMSCeilingByOrigin(ctx context.Context, originUF domain.UF) (map[domain.UF]*domain.ICMSCeiling, error)
}
