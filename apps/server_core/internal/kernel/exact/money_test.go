package exact_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/exact"
)

func TestParseCurrencyRejectsGarbage(t *testing.T) {
	for _, in := range []string{"", "R$", "brlx", "12"} {
		if _, err := exact.ParseCurrency(in); err == nil {
			t.Fatalf("ParseCurrency(%q) returned no error", in)
		}
	}
}

func TestNewMoneyRejectsZeroCurrency(t *testing.T) {
	var zero exact.Currency
	if _, err := exact.NewMoney(exact.FromInt(10), zero); err == nil {
		t.Fatal("NewMoney with zero Currency returned no error")
	}
}

func TestMoneyStringCarriesCurrency(t *testing.T) {
	brl, err := exact.ParseCurrency("brl")
	if err != nil {
		t.Fatalf("ParseCurrency: %v", err)
	}
	m, err := exact.ParseMoney("82.45", brl)
	if err != nil {
		t.Fatalf("ParseMoney: %v", err)
	}
	if got := m.String(); got != "BRL 82.45" {
		t.Fatalf("String() = %q, want %q", got, "BRL 82.45")
	}
}

// Adding two currencies is not a rounding problem, it is a wrong answer.
func TestMoneyAddRejectsMismatchedCurrency(t *testing.T) {
	brl, _ := exact.ParseCurrency("BRL")
	usd, _ := exact.ParseCurrency("USD")
	a, _ := exact.ParseMoney("10.00", brl)
	b, _ := exact.ParseMoney("10.00", usd)
	if _, err := a.Add(b); err == nil {
		t.Fatal("Add across currencies returned no error")
	}
	if _, err := a.Sub(b); err == nil {
		t.Fatal("Sub across currencies returned no error")
	}
}

func TestMoneyMulDecimalKeepsExactness(t *testing.T) {
	brl, _ := exact.ParseCurrency("BRL")
	price, _ := exact.ParseMoney("178.69", brl)
	rate := exact.MustParseDecimal("0.16")
	fee := price.MulDecimal(rate)
	if got := fee.Amount().StringFixed(4); got != "28.5904" {
		t.Fatalf("StringFixed(4) = %q, want %q", got, "28.5904")
	}
}
