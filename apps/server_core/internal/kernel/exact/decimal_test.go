package exact_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/exact"
)

func TestParseDecimalRejectsGarbage(t *testing.T) {
	for _, in := range []string{"", " ", "abc", "1.2.3", "1,50", "NaN", "Inf"} {
		if _, err := exact.ParseDecimal(in); err == nil {
			t.Fatalf("ParseDecimal(%q) returned no error", in)
		}
	}
}

func TestParseDecimalRoundTrips(t *testing.T) {
	d, err := exact.ParseDecimal("82.45")
	if err != nil {
		t.Fatalf("ParseDecimal: %v", err)
	}
	if got := d.StringFixed(2); got != "82.45" {
		t.Fatalf("StringFixed(2) = %q, want %q", got, "82.45")
	}
}

// The whole reason this type exists: 0.1 + 0.2 is exactly 0.3, and a division
// that does not terminate in binary keeps its exact value until it is printed.
func TestArithmeticIsExact(t *testing.T) {
	a, _ := exact.ParseDecimal("0.1")
	b, _ := exact.ParseDecimal("0.2")
	c, _ := exact.ParseDecimal("0.3")
	if a.Add(b).Cmp(c) != 0 {
		t.Fatalf("0.1 + 0.2 != 0.3; got %s", a.Add(b).StringFixed(20))
	}
}

func TestDivIsExactAndReversible(t *testing.T) {
	num, _ := exact.ParseDecimal("88.45")
	den, _ := exact.ParseDecimal("0.495")
	q, err := num.Div(den)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	if q.Mul(den).Cmp(num) != 0 {
		t.Fatalf("(88.45/0.495)*0.495 != 88.45; got %s", q.Mul(den).StringFixed(20))
	}
	// The published two-decimal answer from the spec's worked example.
	if got := q.StringFixed(2); got != "178.69" {
		t.Fatalf("StringFixed(2) = %q, want %q", got, "178.69")
	}
}

func TestDivByZeroIsAnError(t *testing.T) {
	num, _ := exact.ParseDecimal("1")
	zero := exact.FromInt(0)
	if _, err := num.Div(zero); err == nil {
		t.Fatal("Div by zero returned no error")
	}
}

func TestStringFixedRoundsHalfToEven(t *testing.T) {
	cases := map[string]string{
		"2.345":  "2.34", // half, down to even
		"2.355":  "2.36", // half, up to even
		"2.344":  "2.34",
		"2.346":  "2.35",
		"-2.345": "-2.34",
	}
	for in, want := range cases {
		d, err := exact.ParseDecimal(in)
		if err != nil {
			t.Fatalf("ParseDecimal(%q): %v", in, err)
		}
		if got := d.StringFixed(2); got != want {
			t.Fatalf("StringFixed(2) of %q = %q, want %q", in, got, want)
		}
	}
}

func TestZeroValueIsUsable(t *testing.T) {
	var d exact.Decimal
	if !d.IsZero() {
		t.Fatal("zero Decimal.IsZero() = false")
	}
	if got := d.StringFixed(2); got != "0.00" {
		t.Fatalf("zero StringFixed(2) = %q, want %q", got, "0.00")
	}
}
