package ports

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// ShipmentDetailReader is the orders-local, module-scoped consumer port for the connectors F-01
// order_shipments ingest fact (GET /shipments/{id}). Named distinctly from the existing
// ShipmentReader (shipment_reader.go, GetShipmentInfo/connectorsdomain.ShipmentInfo — a
// narrower, still-LIVE EnrichService call site) so the two never collide: this port returns the
// FULL connectorsdomain.ShipmentDetail (0088 shape), a different type for a different call site.
//
// A 404/410 on the primary provider call is honest-absence: the adapter returns
// connectorsdomain.ShipmentDetail{Found: false} with a nil error (feature.md Negative
// Scenarios) — IngestOrder must persist the order without an order_shipments row in that case,
// never treat Found==false as a Go error.
type ShipmentDetailReader interface {
	GetShipmentDetail(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentDetail, error)
}
