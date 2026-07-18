package domain

import "strings"

// OrderBucket is the workflow bucket a marketplace order is derived into for
// the pedidos Lista status tabs, Kanban columns, and KPI row (F-02; ruling
// D-57). It is a display/aggregation concept distinct from the RAW
// provider_status stored on OrderReadModel.Status.
type OrderBucket string

const (
	BucketNovo      OrderBucket = "novo"
	BucketFaturar   OrderBucket = "faturar"
	BucketEnviar    OrderBucket = "enviar"
	BucketEnviado   OrderBucket = "enviado"
	BucketCancelado OrderBucket = "cancelado"
)

// DeriveOrderBucket maps the RAW provider status + shipment presence to a
// workflow bucket. It is the single authority for bucket membership across
// the API and the pedidos UI (F01-A).
//
// hasShipment is a PROXY for the design's NF-emitted state: orders storage
// carries no true ERP NF-emission signal, only the shipping_id column. The
// faturar/enviar split therefore approximates design semantics using
// available marketplace data, degrading honestly per ruling D-57(6).
func DeriveOrderBucket(providerStatus string, hasShipment bool) OrderBucket {
	status := strings.ToLower(strings.TrimSpace(providerStatus))
	switch status {
	case "cancelled", "invalid":
		return BucketCancelado
	case "delivered", "shipped":
		return BucketEnviado
	case "paid", "confirmed":
		if hasShipment {
			return BucketEnviar
		}
		return BucketFaturar
	default:
		return BucketNovo
	}
}
