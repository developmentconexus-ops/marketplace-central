package ports

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// OrderDetailReader is the orders-local, module-scoped consumer port for the connectors F-01
// ingest fact (GET /orders/{id}). It mirrors the simplified-key pattern of BuyerFiscalReader
// (buyer_fiscal_reader.go): the real adapter (composition-owned, orders_ingest_adapters.go)
// resolves installationID to a connectors domain.ProviderAccountRef and calls
// CapabilityAdapter.GetOrderDetail — that resolution is not this module's concern.
//
// A third-party order (403) or a not-found order (404) surfaces here as an error
// (connectorsdomain.ErrUnauthorized/ErrNotFound, via mapPricingReaderError) — IngestOrder is
// what classifies that into domain.ErrOrderUnavailable and stops before any write.
type OrderDetailReader interface {
	GetOrderDetail(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.OrderDetail, error)
}
