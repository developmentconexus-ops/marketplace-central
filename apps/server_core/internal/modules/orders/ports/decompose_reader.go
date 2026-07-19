package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

// Decomposer is the orders-local, module-scoped consumer port for the
// decomposição/DIFAL/retorno profitability fact. It mirrors the nil-port
// idiom established by CostReader (cost_reader.go) and ShipmentReader
// (shipment_reader.go): the real adapter (hub C2, post-merge) resolves the
// pricing IC-04 decomposition engine — that resolution is not this module's
// concern. The bool is ok: false means the caller must fall back to
// domain.UnknownOrderProfitability(), the honest-unknown value; a nil
// Decomposer is honest-unknown too, never a panic.
type Decomposer interface {
	Decompose(ctx context.Context, order domain.OrderReadModel) (domain.OrderProfitability, bool)
}
