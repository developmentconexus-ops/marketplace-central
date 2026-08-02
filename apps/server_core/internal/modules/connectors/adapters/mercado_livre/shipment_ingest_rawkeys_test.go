package mercadolivre

import (
	"encoding/json"
	"os"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/rawkeys"
)

// shipmentBodyIgnoredKeys: cada entrada e uma DECISAO, nao um esquecimento.
var shipmentBodyIgnoredKeys = []string{}

func TestShipmentBodyDeclaresEveryConsumedKey(t *testing.T) {
	raw, err := os.ReadFile("testdata/shipment_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	missing, err := rawkeys.Undeclared(json.RawMessage(raw), mlIngestShipmentResponse{}, shipmentBodyIgnoredKeys)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("DTO nao declara chaves que o ML manda: %v", missing)
	}
}
