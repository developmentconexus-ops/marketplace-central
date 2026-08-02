package mercadolivre

import (
	"encoding/json"
	"os"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/rawkeys"
)

// orderBodyIgnoredKeys: cada entrada e uma DECISAO, nao um esquecimento.
// Entrada sem motivo escrito devolve exatamente o silencio que o detector
// existe para quebrar.
var orderBodyIgnoredKeys = []string{
	"context",         // canal/site; ja conhecido pela instalacao
	"coupon",          // cupom do ML; entra na decomposicao so se o operador ratificar
	"expiration_date", // prazo de pagamento; sem consumidor
	"feedback",        // reputacao; fora de escopo
	"paid_amount",     // ja coberto por payments[].total_paid_amount
	"seller",          // e a nossa propria conta
	"taxes",           // imposto do ML (medido null); o nosso e P2.b
	"total_amount",    // derivado de order_items; nao e fonte
}

func TestOrderBodyDeclaresEveryConsumedKey(t *testing.T) {
	raw, err := os.ReadFile("testdata/order_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	missing, err := rawkeys.Undeclared(json.RawMessage(raw), mlIngestOrderResponse{}, orderBodyIgnoredKeys)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("DTO nao declara chaves que o ML manda: %v", missing)
	}
}
