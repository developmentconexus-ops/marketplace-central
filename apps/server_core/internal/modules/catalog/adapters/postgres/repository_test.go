package postgres

import "testing"

func TestPositiveProductIDRejectsLegacyAndNonCanonicalValues(t *testing.T) {
	for _, raw := range []string{"", "REF-1", "0", "-1", "01"} {
		if _, err := positiveProductID(raw); err == nil {
			t.Fatalf("positiveProductID(%q) unexpectedly succeeded", raw)
		}
	}
	if got, err := positiveProductID("1"); err != nil || got != 1 {
		t.Fatalf("positiveProductID(1) = %d, %v", got, err)
	}
}
