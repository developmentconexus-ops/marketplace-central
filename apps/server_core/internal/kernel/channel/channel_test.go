package channel_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/channel"
)

func TestParseCodeRejectsEmpty(t *testing.T) {
	if _, err := channel.ParseCode(""); err == nil {
		t.Fatal("ParseCode(\"\") returned no error")
	}
}

func TestParseCodeNormalisesCase(t *testing.T) {
	code, err := channel.ParseCode("MercadoLivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	if got := code.String(); got != "mercadolivre" {
		t.Fatalf("String() = %q, want %q", got, "mercadolivre")
	}
}

func TestNewAccountRefRejectsZeroCode(t *testing.T) {
	var zero channel.Code
	if _, err := channel.NewAccountRef(zero, "acc-1"); err == nil {
		t.Fatal("NewAccountRef with zero Code returned no error")
	}
}

func TestNewAccountRefRejectsEmptyExternal(t *testing.T) {
	code, err := channel.ParseCode("mercadolivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	if _, err := channel.NewAccountRef(code, ""); err == nil {
		t.Fatal("NewAccountRef with empty external id returned no error")
	}
}

func TestAccountRefRoundTrips(t *testing.T) {
	code, err := channel.ParseCode("mercadolivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	ref, err := channel.NewAccountRef(code, "123456789")
	if err != nil {
		t.Fatalf("NewAccountRef: %v", err)
	}
	if ref.Channel().String() != "mercadolivre" || ref.External() != "123456789" {
		t.Fatalf("round trip lost data: %q / %q", ref.Channel().String(), ref.External())
	}
}
