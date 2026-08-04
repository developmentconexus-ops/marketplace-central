package domain

import "testing"

func TestEntityValid(t *testing.T) {
	valid := []Entity{
		EntityProducts,
		EntityListings,
		EntityOrders,
		EntityMarket,
		EntityMarketQueue,
		EntityTariffs,
		EntityICMSMatrix,
	}
	for _, e := range valid {
		if !e.Valid() {
			t.Errorf("Entity(%q).Valid() = false, quero true", e)
		}
	}

	invalid := []Entity{"", "icms", "ICMS_MATRIX", "matrix", "produtos"}
	for _, e := range invalid {
		if e.Valid() {
			t.Errorf("Entity(%q).Valid() = true, quero false", e)
		}
	}
}

func TestEntityICMSMatrixWireValue(t *testing.T) {
	// O valor de fio é a chave de sync_state e o rótulo que a tela /integracoes
	// gera por title-case. Mudá-lo abandona a linha já gravada no banco em vez
	// de continuá-la, então o literal é asserção, não detalhe.
	if got := string(EntityICMSMatrix); got != "icms_matrix" {
		t.Fatalf("EntityICMSMatrix = %q, quero \"icms_matrix\"", got)
	}
}
