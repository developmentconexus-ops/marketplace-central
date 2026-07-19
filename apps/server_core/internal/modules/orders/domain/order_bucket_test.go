package domain

import "testing"

// TestDeriveOrderBucket is the exhaustive truth table for DeriveOrderBucket
// (F01-A). It is the bucket authority: any change to the rules must be
// reflected here first.
func TestDeriveOrderBucket(t *testing.T) {
	cases := []struct {
		name           string
		providerStatus string
		hasShipment    bool
		want           OrderBucket
	}{
		{"cancelled lowercase", "cancelled", false, BucketCancelado},
		{"cancelled uppercase", "CANCELLED", false, BucketCancelado},
		{"cancelled with shipment", "cancelled", true, BucketCancelado},
		{"invalid", "invalid", false, BucketCancelado},
		{"invalid with shipment", "invalid", true, BucketCancelado},
		{"delivered", "delivered", false, BucketEnviado},
		{"delivered with shipment", "delivered", true, BucketEnviado},
		{"shipped", "shipped", false, BucketEnviado},
		{"shipped mixed case", "Shipped", true, BucketEnviado},
		{"paid with shipment", "paid", true, BucketEnviar},
		{"paid without shipment", "paid", false, BucketFaturar},
		{"confirmed with shipment", "confirmed", true, BucketEnviar},
		{"confirmed without shipment", "confirmed", false, BucketFaturar},
		{"confirmed uppercase with shipment", "CONFIRMED", true, BucketEnviar},
		{"pending", "pending", false, BucketNovo},
		{"pending with shipment", "pending", true, BucketNovo},
		{"payment_required", "payment_required", false, BucketNovo},
		{"payment_in_process", "payment_in_process", false, BucketNovo},
		{"new", "new", false, BucketNovo},
		{"empty status", "", false, BucketNovo},
		{"empty status with shipment", "", true, BucketNovo},
		{"unknown status", "some_unknown_status", false, BucketNovo},
		{"whitespace padded paid", "  paid  ", true, BucketEnviar},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DeriveOrderBucket(tc.providerStatus, tc.hasShipment)
			if got != tc.want {
				t.Fatalf("DeriveOrderBucket(%q, %v) = %q, want %q", tc.providerStatus, tc.hasShipment, got, tc.want)
			}
		})
	}
}
