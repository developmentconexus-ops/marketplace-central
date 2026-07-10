package internalread

import (
	"context"
	"time"

	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type StockReader struct {
	service internalreadapp.Service
}

func NewStockReader(service internalreadapp.Service) StockReader {
	return StockReader{service: service}
}

func (r StockReader) GetSellableStock(ctx context.Context, productID int, policy inventorydomain.StockPolicy) (inventorydomain.InternalStockEvidence, inventorydomain.ProductEvidence, error) {
	result, err := r.service.GetSellableStock(ctx, internalreadports.SellableStockInput{
		ProductID: productID,
		Policy:    MapStockPolicy(policy),
		Freshness: MapFreshnessPolicy(policy.Freshness),
	})
	if err != nil {
		return inventorydomain.InternalStockEvidence{}, inventorydomain.ProductEvidence{}, err
	}
	return inventorydomain.InternalStockEvidence{
			Quantity: floatQuantityToInt(result.Quantity),
			Source: inventorydomain.SourceEvidence{
				ObservedAt: observedAt(result.Source),
			},
		}, inventorydomain.ProductEvidence{
			ProductID: productID,
		}, nil
}

func floatQuantityToInt(value *float64) *int {
	if value == nil {
		return nil
	}
	v := int(*value)
	return &v
}

func observedAt(source internalreaddomain.SourceMetadata) *time.Time {
	if source.ObservedAt != nil {
		t := source.ObservedAt.UTC()
		return &t
	}
	t := source.FetchedAt.UTC()
	return &t
}
