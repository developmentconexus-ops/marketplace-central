package ports

import (
	"context"
	"errors"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// ShipmentReader is the orders-local, module-scoped consumer port for the
// shipment fact. Since F-03 (read-path switch), the real adapter
// (orders/adapters/postgres.ShipmentReader) resolves installationID to a
// tenant-scoped Postgres read of order_shipments (migration 0088) — no
// connectors call. installationID stays part of the signature for port
// stability even though the Postgres adapter does not use it (order_shipments
// carries no installation_id column).
type ShipmentReader interface {
	GetShipment(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentInfo, error)
}

// ErrShipmentNotFound is the honest-unknown sentinel a ShipmentReader returns
// when no shipment row/fact exists yet for the given key (e.g. an order
// ingested before F-02, or one with no shipping_id) — never fabricated as a
// zero-value ShipmentInfo with a nil error. EnrichService.fetchShipment
// recognizes this sentinel (via errors.Is) and degrades to a nil
// ShipmentEnrichment WITHOUT a warn log, distinct from a real read failure
// (which still warns): a not-yet-ingested order is an expected, common case.
var ErrShipmentNotFound = errors.New("orders: shipment not found")
