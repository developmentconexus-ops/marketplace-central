package domain

import "testing"

// TestDeriveOrderBucket is the exhaustive truth table for DeriveOrderBucket
// (F01-A). It is the bucket authority: any change to the rules must be
// reflected here first.
func TestDeriveOrderBucket(t *testing.T) {
	cases := []struct {
		name           string
		providerStatus string
		shipmentStatus string
		tags           []string
		hasShipment    bool
		want           OrderBucket
	}{
		// provider_status-only fallback (no shipment status, no tags) — the
		// shipment-blind SQL summary path and pre-shipment ML orders.
		{"cancelled lowercase", "cancelled", "", nil, false, BucketCancelado},
		{"cancelled uppercase", "CANCELLED", "", nil, false, BucketCancelado},
		{"cancelled with shipment", "cancelled", "", nil, true, BucketCancelado},
		{"invalid", "invalid", "", nil, false, BucketCancelado},
		{"invalid with shipment", "invalid", "", nil, true, BucketCancelado},
		{"delivered provider status", "delivered", "", nil, false, BucketEnviado},
		{"shipped provider status", "shipped", "", nil, false, BucketEnviado},
		{"shipped mixed case", "Shipped", "", nil, true, BucketEnviado},
		{"paid with shipment", "paid", "", nil, true, BucketEnviar},
		{"paid without shipment", "paid", "", nil, false, BucketFaturar},
		{"confirmed with shipment", "confirmed", "", nil, true, BucketEnviar},
		{"confirmed without shipment", "confirmed", "", nil, false, BucketFaturar},
		{"confirmed uppercase with shipment", "CONFIRMED", "", nil, true, BucketEnviar},
		{"pending", "pending", "", nil, false, BucketNovo},
		{"pending with shipment", "pending", "", nil, true, BucketNovo},
		{"payment_required", "payment_required", "", nil, false, BucketNovo},
		{"payment_in_process", "payment_in_process", "", nil, false, BucketNovo},
		{"new", "new", "", nil, false, BucketNovo},
		{"empty status", "", "", nil, false, BucketNovo},
		{"empty status with shipment", "", "", nil, true, BucketNovo},
		{"unknown status", "some_unknown_status", "", nil, false, BucketNovo},
		{"whitespace padded paid", "  paid  ", "", nil, true, BucketEnviar},

		// Order tag "delivered" — the robust order-level delivered signal. ML
		// keeps provider_status "paid" after delivery; the tag lifts it to enviado
		// even with no live shipment status (shipment lookup may have degraded).
		{"paid + delivered tag", "paid", "", []string{"paid", "delivered"}, true, BucketEnviado},
		{"paid + delivered tag no shipment status", "paid", "", []string{"delivered"}, true, BucketEnviado},
		{"delivered tag mixed case", "paid", "", []string{"Delivered"}, true, BucketEnviado},
		{"delivered tag whitespace", "paid", "", []string{" delivered "}, true, BucketEnviado},
		{"non-delivered tags stay enviar", "paid", "", []string{"paid", "claim_opened"}, true, BucketEnviar},
		{"delivered tag overrides ready_to_ship", "paid", "ready_to_ship", []string{"delivered"}, true, BucketEnviado},
		{"cancelled wins over delivered tag", "cancelled", "", []string{"delivered"}, true, BucketCancelado},

		// Live shipment status — ML dispatch lifecycle that provider_status
		// ("paid") does not track. Authoritative when the delivered tag is absent.
		{"shipment shipped -> enviado", "paid", "shipped", nil, true, BucketEnviado},
		{"shipment delivered -> enviado", "paid", "delivered", nil, true, BucketEnviado},
		{"shipment shipped mixed case", "paid", "Shipped", nil, true, BucketEnviado},
		{"shipment ready_to_ship -> enviar", "paid", "ready_to_ship", nil, true, BucketEnviar},
		{"shipment handling -> enviar", "paid", "handling", nil, true, BucketEnviar},
		{"shipment pending -> enviar", "paid", "pending", nil, true, BucketEnviar},
		{"shipment not_delivered falls back to provider", "paid", "not_delivered", nil, true, BucketEnviar},
		{"shipment unknown status falls back to provider", "paid", "in_hub", nil, true, BucketEnviar},
		{"shipment shipped but faturar provider no shipment flag", "paid", "shipped", nil, false, BucketEnviado},
		{"shipment ready_to_ship with confirmed provider", "confirmed", "ready_to_ship", nil, true, BucketEnviar},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DeriveOrderBucket(tc.providerStatus, tc.shipmentStatus, tc.tags, tc.hasShipment)
			if got != tc.want {
				t.Fatalf("DeriveOrderBucket(%q, %q, %v, %v) = %q, want %q", tc.providerStatus, tc.shipmentStatus, tc.tags, tc.hasShipment, got, tc.want)
			}
		})
	}
}
