package ports

import "context"

// OrderIngestor is IC-06's single write-path port (ADR-024): every trigger that mutates an order
// — today's import enumeration, and per milestone.md, backfill/incremental sync/webhook worker
// later — converges on this ONE call per order. IngestOrder fetches OrderDetail (and, when
// present, ShipmentDetail + BuyerFiscal), derives the bucket, and persists header+items+
// payments+shipment atomically; it never partially writes.
type OrderIngestor interface {
	IngestOrder(ctx context.Context, installationID, providerOrderID string) error
}
