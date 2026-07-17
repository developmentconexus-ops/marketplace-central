package domain

import "testing"

func TestSetFilterValueGrammar(t *testing.T) {
	cases := []struct {
		key, value string
		valid      bool
	}{
		{"status", "active", true}, {"status", "banana", false},
		{"sync_state", "paused_sync", true}, {"link_state", "resolved", true},
		{"exception", "below_margin", true}, {"has_exception", "false", true},
		{"has_exception", "1", false}, {"listing_type_code", "gold_pro", true},
		{"product_id", "42664", true}, {"product_id", "042664", false},
		{"unknown", "x", false}, {"status", "gte:active", false},
	}
	for _, tc := range cases {
		var filter ListingFilter
		err := SetFilterValue(&filter, tc.key, tc.value)
		if (err == nil) != tc.valid {
			t.Errorf("SetFilterValue(%q, %q) error = %v; valid=%v", tc.key, tc.value, err, tc.valid)
		}
	}
}
