package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type ShipmentReader interface {
	GetShipmentInfo(ctx context.Context, accountRef domain.ProviderAccountRef, shipmentID string) (domain.ShipmentInfo, error)
	GetFreeShippingCost(ctx context.Context, accountRef domain.ProviderAccountRef, query domain.FreeShippingQuery) (domain.FreeShippingCost, error)
}
