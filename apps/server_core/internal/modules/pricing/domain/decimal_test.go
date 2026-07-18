package domain

import "testing"

func TestParseRatAcceptsPlainDecimal(t *testing.T) {
	value, err := ParseRat("79.00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := FormatRatHalfUp(value, 2); got != "79.00" {
		t.Fatalf("got %q, want 79.00", got)
	}
}

func TestParseRatRejectsNonDecimalStrings(t *testing.T) {
	for _, in := range []string{"", "abc", "12.34.56", "12,34", "1e10", "  79.00  ", "79."} {
		if _, err := ParseRat(in); err == nil {
			t.Fatalf("ParseRat(%q) = nil error, want error", in)
		}
	}
}

func TestFormatRatHalfUpMoneyTwoPlaces(t *testing.T) {
	cases := map[string]string{
		"79":     "79.00",
		"79.005": "79.01", // IC-04 golden: x.xx5 rounds up
		"2.675":  "2.68",
		"1.005":  "1.01",
		"0.004":  "0.00",
		"0.005":  "0.01",
	}
	for in, want := range cases {
		value, err := ParseRat(in)
		if err != nil {
			t.Fatalf("ParseRat(%q): unexpected error: %v", in, err)
		}
		if got := FormatRatHalfUp(value, 2); got != want {
			t.Fatalf("FormatRatHalfUp(%q, 2) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatRatHalfUpPctPlaces(t *testing.T) {
	cases := map[string]string{
		"15":         "15.0000",
		"0.049":      "0.0490",
		"0.04995":    "0.0500", // golden boundary: x.xxx5 rounds up at pct precision
		"12.345649":  "12.3456",
		"12.3456500": "12.3457",
	}
	for in, want := range cases {
		value, err := ParseRat(in)
		if err != nil {
			t.Fatalf("ParseRat(%q): unexpected error: %v", in, err)
		}
		if got := FormatRatHalfUp(value, 4); got != want {
			t.Fatalf("FormatRatHalfUp(%q, 4) = %q, want %q", in, got, want)
		}
	}
}

func TestNilMoneyIsUnknownNeverZero(t *testing.T) {
	var m *Money
	if !m.IsUnknown() {
		t.Fatalf("nil *Money must be unknown (ADR-17)")
	}

	known := &Money{Amount: "0.00", Currency: "BRL"}
	if known.IsUnknown() {
		t.Fatalf("a known zero-valued Money must not be reported as unknown")
	}
}
