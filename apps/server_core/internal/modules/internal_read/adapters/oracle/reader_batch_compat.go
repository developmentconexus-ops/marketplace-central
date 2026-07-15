package oracle

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

// The production root injects one shared semaphore into BatchReader. These
// forwarding methods keep Reader compatible with the batch port for bounded
// adapter tests and callers that still construct the legacy reader directly.
var compatibilityBatchSemaphore = oraclebatch.NewSemaphore(4)

func (r *Reader) GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.CostAsOf, error) {
	if r == nil {
		return NewBatchReader(nil, compatibilityBatchSemaphore).GetCostFactsByIDs(ctx, ids)
	}
	return NewBatchReader(r.db, compatibilityBatchSemaphore).GetCostFactsByIDs(ctx, ids)
}

func (r *Reader) GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.TaxInputs, error) {
	if r == nil {
		return NewBatchReader(nil, compatibilityBatchSemaphore).GetTaxFactsByIDs(ctx, ids)
	}
	return NewBatchReader(r.db, compatibilityBatchSemaphore).GetTaxFactsByIDs(ctx, ids)
}

func (r *Reader) GetICMSCeilingByOrigin(ctx context.Context, originUF domain.UF) (map[domain.UF]*domain.ICMSCeiling, error) {
	if r == nil {
		return NewBatchReader(nil, compatibilityBatchSemaphore).GetICMSCeilingByOrigin(ctx, originUF)
	}
	return NewBatchReader(r.db, compatibilityBatchSemaphore).GetICMSCeilingByOrigin(ctx, originUF)
}
