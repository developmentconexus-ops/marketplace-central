package ports

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// BuyerFiscalReader is the orders-local, module-scoped consumer port for the
// connectors buyer-fiscal fact. It mirrors the simplified-key pattern of
// ShipmentReader (shipment_reader.go): the real adapter (a grant-gated
// composition slice) resolves installationID to a connectors
// domain.ProviderAccountRef and calls connectors ports.BuyerFiscalReader —
// that resolution is not this module's concern.
//
// A buyer without billing data degrades, inside the adapter, to an empty
// honest-absence BuyerFiscalInfo (HasData() == false), never an error; only a
// real read failure (auth/transient) surfaces as an error here.
type BuyerFiscalReader interface {
	GetBuyerFiscal(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.BuyerFiscalInfo, error)
}
