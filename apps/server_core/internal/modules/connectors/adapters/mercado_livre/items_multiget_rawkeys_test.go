package mercadolivre

import (
	"encoding/json"
	"os"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/rawkeys"
)

// multigetBodyIgnoredKeys are keys the multiget body DTO deliberately does NOT
// type. Every entry is a decision, not an oversight: the value is either
// already carried elsewhere, or out of scope for this seam.
//
// ADR-C6: a key may only be added here with a reason. An unexplained entry
// turns the detector back into the silence it exists to break.
var multigetBodyIgnoredKeys = []string{
	"accepts_mercadopago", // não usado em nenhuma tela
	"automatic_relist",    // não usado em nenhuma tela
	"base_price",          // preço de lista antes de promoção; fora do escopo do P1
	"catalog_listing",     // competição de catálogo = PLANEJADO, bloqueada por opt-in
	"deal_ids",            // promoções; fora do escopo do P1
	"domain_id",           // domínio da categoria; fora do escopo do P1
	"health",              // descontinuado pelo ML; a fonte é /item/{id}/performance
	"inventory_id",        // estoque Full; a conta não usa
	"last_updated",        // sincronização incremental é do M-06
	"original_price",      // preço de lista antes de promoção; fora do escopo do P1
	"pictures",            // thumbnail já é tipado
	"sale_terms",          // garantia/termos; fora do escopo do P1
	"seller_address",      // PII do vendedor: nunca sai do adapter
	"seller_id",           // já conhecido pela conta da instalação
	"sub_status",          // motivo de bloqueio = PLANEJADO (P4)
	"user_product_id",     // estoque Full; a conta não usa
	"warranty",            // fora do escopo do P1
}

func TestMultigetItemBodyDeclaresEveryConsumedKey(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_multiget_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	missing, err := rawkeys.Undeclared(json.RawMessage(raw), mlMultigetItemBody{}, multigetBodyIgnoredKeys)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("mlMultigetItemBody nao declara chaves que o ML manda: %v", missing)
	}
}
