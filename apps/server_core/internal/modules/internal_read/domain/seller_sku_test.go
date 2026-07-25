package domain

import "testing"

func TestIsValidCodprod(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{"real codprod from the operator's ERP", "100002", true},
		{"single digit", "7", true},
		{"surrounding whitespace is not part of the code", "  100002  ", true},
		{"legacy REFFORN is not a CODPROD", "ZP1704.1.", false},
		{"legacy REFFORN with dots and digits", "L.87.22", false},
		{"empty", "", false},
		{"whitespace only", "   ", false},
		{"zero names no product", "0", false},
		{"padded zero names no product", "0000", false},
		{"longer than an internal product id can hold", "1234567890", false},
		{"negative is not a CODPROD literal", "-10", false},
		{"decimal is not a CODPROD literal", "10.5", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidCodprod(tc.value); got != tc.want {
				t.Fatalf("IsValidCodprod(%q) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}
