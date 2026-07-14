package internalread

import (
	"context"
	"errors"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type UnavailableStockReader struct {
	Err error
}

func (r UnavailableStockReader) GetSellableStock(context.Context, int, inventorydomain.StockPolicy) (inventorydomain.InternalStockEvidence, inventorydomain.ProductEvidence, error) {
	if r.Err != nil {
		return inventorydomain.InternalStockEvidence{}, inventorydomain.ProductEvidence{}, r.Err
	}
	return inventorydomain.InternalStockEvidence{}, inventorydomain.ProductEvidence{}, errors.New("INVENTORY_INTERNAL_READ_UNAVAILABLE")
}

type UnavailableStockBatchReader struct {
	Err error
}

func (r UnavailableStockBatchReader) GetStockFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.StockFact, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return nil, errors.New("INVENTORY_INTERNAL_READ_UNAVAILABLE")
}
