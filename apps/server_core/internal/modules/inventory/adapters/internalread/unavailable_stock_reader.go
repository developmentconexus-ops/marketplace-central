package internalread

import (
	"context"
	"errors"

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
