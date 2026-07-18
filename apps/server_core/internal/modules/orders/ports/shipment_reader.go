package ports

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// ShipmentReader is the orders-local, module-scoped consumer port for the
// connectors shipment fact. It mirrors the simplified-key pattern established
// by CostReader (cost_reader.go): the real adapter (a later, grant-gated
// slice) resolves installationID to a connectors domain.ProviderAccountRef
// and calls connectors ports.ShipmentReader.GetShipmentInfo — that
// resolution is not this module's concern.
type ShipmentReader interface {
	GetShipment(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentInfo, error)
}
