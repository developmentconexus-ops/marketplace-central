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

// DeriveOrderBucket maps the RAW provider status + shipment facts to a
// workflow bucket. It is the single authority for bucket membership across
// the API and the pedidos UI (F01-A).
//
// Why more than provider_status: Mercado Livre keeps order.provider_status at
// "paid" for the WHOLE post-payment lifecycle — shipped/delivered live on the
// SHIPMENT object (status+substatus) and on the order TAGS (e.g. "delivered"),
// never on provider_status (ML docs: gerenciamento-de-envios status table;
// pedidos-e-opinioes order tags). Deriving from provider_status alone therefore
// pins every paid+shipped order in "enviar" forever (the "A ENVIAR 24 /
// ENVIADOS 0" defect). Two extra signals lift it:
//
//   - tags: the order-level "delivered" tag is the most robust delivered signal
//     — it is persisted with the order (tags_json) and survives even a total
//     shipment-lookup failure, so a delivered order leaves "enviar" regardless
//     of the fragile live /shipments read.
//   - shipmentStatus: the live shipment status distinguishes shipped/delivered
//     (→ enviado) from ready_to_ship/handling/pending (→ enviar) when the tag
//     is absent.
//
// hasShipment stays the pre-shipment faturar/enviar proxy (shipping_id present
// but no live status/tag yet), degrading honestly per ruling D-57(6).
// providerStatus keeps sole authority over novo/faturar/cancelado.
func DeriveOrderBucket(providerStatus, shipmentStatus string, tags []string, hasShipment bool) OrderBucket {
	status := strings.ToLower(strings.TrimSpace(providerStatus))

	// Order-level cancel wins first: a cancelled order is cancelado even if a
	// shipment was created before the cancellation.
	if status == "cancelled" || status == "invalid" {
		return BucketCancelado
	}

	// Order tag "delivered" is authoritative and shipment-lookup-independent.
	if hasTag(tags, "delivered") {
		return BucketEnviado
	}

	// Live shipment status is authoritative for the ML dispatch lifecycle that
	// provider_status does not track. Every "shipped" substatus (out_for_delivery,
	// receiver_absent, buyer_rescheduled, …) stays under status "shipped" → enviado.
	switch strings.ToLower(strings.TrimSpace(shipmentStatus)) {
	case "shipped", "delivered":
		return BucketEnviado
	case "ready_to_ship", "handling", "pending":
		return BucketEnviar
	}

	// No delivered tag and no (or an unrecognised) live shipment status: fall
	// back to provider_status. This is also the shipment-blind SQL summary path.
	switch status {
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

// hasTag reports whether tags contains target (case-insensitive, trimmed).
func hasTag(tags []string, target string) bool {
	for _, tag := range tags {
		if strings.EqualFold(strings.TrimSpace(tag), target) {
			return true
		}
	}
	return false
}
