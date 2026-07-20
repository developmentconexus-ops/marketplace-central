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
		faturado       bool
		want           OrderBucket
	}{
		// provider_status-only fallback (no shipment status, no tags) — the
		// shipment-blind SQL summary path and pre-shipment ML orders.
		{"cancelled lowercase", "cancelled", "", nil, false, BucketCancelado},
		{"cancelled uppercase", "CANCELLED", "", nil, false, BucketCancelado},
		{"cancelled with faturado", "cancelled", "", nil, true, BucketCancelado},
		{"invalid", "invalid", "", nil, false, BucketCancelado},
		{"invalid with faturado", "invalid", "", nil, true, BucketCancelado},
		{"delivered provider status", "delivered", "", nil, false, BucketEnviado},
		{"shipped provider status", "shipped", "", nil, false, BucketEnviado},
		{"shipped mixed case", "Shipped", "", nil, true, BucketEnviado},
		{"paid + faturado", "paid", "", nil, true, BucketEnviar},
		{"paid without faturado", "paid", "", nil, false, BucketFaturar},
		{"confirmed + faturado", "confirmed", "", nil, true, BucketEnviar},
		{"confirmed without faturado", "confirmed", "", nil, false, BucketFaturar},
		{"confirmed uppercase + faturado", "CONFIRMED", "", nil, true, BucketEnviar},
		{"pending", "pending", "", nil, false, BucketNovo},
		{"pending with faturado", "pending", "", nil, true, BucketNovo},
		{"payment_required", "payment_required", "", nil, false, BucketNovo},
		{"payment_in_process", "payment_in_process", "", nil, false, BucketNovo},
		{"new", "new", "", nil, false, BucketNovo},
		{"empty status", "", "", nil, false, BucketNovo},
		{"empty status with faturado", "", "", nil, true, BucketNovo},
		{"unknown status", "some_unknown_status", "", nil, false, BucketNovo},
		{"whitespace padded paid faturado", "  paid  ", "", nil, true, BucketEnviar},

		// Order tag "delivered" — the robust order-level delivered signal. ML
		// keeps provider_status "paid" after delivery; the tag lifts it to
		// enviado even with no live shipment status and regardless of faturado.
		{"paid + delivered tag", "paid", "", []string{"paid", "delivered"}, true, BucketEnviado},
		{"paid + delivered tag no faturado", "paid", "", []string{"delivered"}, false, BucketEnviado},
		{"delivered tag mixed case", "paid", "", []string{"Delivered"}, true, BucketEnviado},
		{"delivered tag whitespace", "paid", "", []string{" delivered "}, true, BucketEnviado},
		{"non-delivered tags stay faturado-gated", "paid", "", []string{"paid", "claim_opened"}, true, BucketEnviar},
		{"non-delivered tags stay faturar without faturado", "paid", "", []string{"paid", "claim_opened"}, false, BucketFaturar},
		{"delivered tag overrides ready_to_ship", "paid", "ready_to_ship", []string{"delivered"}, false, BucketEnviado},
		{"cancelled wins over delivered tag", "cancelled", "", []string{"delivered"}, true, BucketCancelado},

		// Live shipment status — ML dispatch lifecycle that provider_status
		// ("paid") does not track. Actually-dispatched wins over the faturado
		// gate: shipped/delivered forces enviado regardless of faturado.
		{"shipment shipped -> enviado", "paid", "shipped", nil, false, BucketEnviado},
		{"shipment delivered -> enviado", "paid", "delivered", nil, false, BucketEnviado},
		{"shipment shipped mixed case", "paid", "Shipped", nil, false, BucketEnviado},
		{"shipment shipped with faturado true too -> enviado", "paid", "shipped", nil, true, BucketEnviado},
		{"shipment ready_to_ship falls back to faturado gate: faturado", "paid", "ready_to_ship", nil, true, BucketEnviar},
		{"shipment ready_to_ship falls back to faturado gate: not faturado", "paid", "ready_to_ship", nil, false, BucketFaturar},
		{"shipment handling falls back to faturado gate", "paid", "handling", nil, true, BucketEnviar},
		{"shipment pending falls back to faturado gate", "paid", "pending", nil, false, BucketFaturar},
		{"shipment not_delivered falls back to faturado gate", "paid", "not_delivered", nil, true, BucketEnviar},
		{"shipment unknown status falls back to faturado gate", "paid", "in_hub", nil, false, BucketFaturar},
		{"confirmed with ready_to_ship shipment and faturado", "confirmed", "ready_to_ship", nil, true, BucketEnviar},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DeriveOrderBucket(tc.providerStatus, tc.shipmentStatus, tc.tags, tc.faturado)
			if got != tc.want {
				t.Fatalf("DeriveOrderBucket(%q, %q, %v, %v) = %q, want %q", tc.providerStatus, tc.shipmentStatus, tc.tags, tc.faturado, got, tc.want)
			}
		})
	}
}
