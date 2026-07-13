package postgres

import "testing"

func TestPricingProductIDMustBeCanonicalCODPROD(t *testing.T) {
	for _, raw := range []string{"legacy-ref", "0", "-1", "01"} {
		if _, err := positiveProductID(raw); err == nil {
			t.Fatalf("positiveProductID(%q) unexpectedly succeeded", raw)
		}
	}
	if got, err := positiveProductID("42"); err != nil || got != 42 {
		t.Fatalf("positiveProductID(42) = %d, %v", got, err)
	}
}
