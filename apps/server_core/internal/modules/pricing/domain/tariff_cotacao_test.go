package domain

import "testing"

// FonteCotacao is degrau 3: a live commission quoted from the provider catalog
// match + listing_prices. Its wire value is the design §11 authority ("COTACAO").
func TestFonteCotacaoValue(t *testing.T) {
	if FonteCotacao != "COTACAO" {
		t.Fatalf("FonteCotacao = %q, want COTACAO", FonteCotacao)
	}
}

// A degrau-3 ComponentResolution carries a real source timestamp in Data
// (unlike degrau 4, which leaves it nil per ADR-17). This proves the field exists.
func TestComponentResolutionCarriesData(t *testing.T) {
	ts := "2026-07-19T12:00:00Z"
	c := ComponentResolution{Valor: strptr("13.00"), Fonte: FonteCotacao, Degrau: 3, Data: &ts}
	if c.Data == nil || *c.Data != ts {
		t.Fatalf("ComponentResolution.Data = %v, want %q", c.Data, ts)
	}
}
