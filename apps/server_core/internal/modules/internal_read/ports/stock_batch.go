package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type StockBatchReader interface {
	GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.StockFact, error)
}
