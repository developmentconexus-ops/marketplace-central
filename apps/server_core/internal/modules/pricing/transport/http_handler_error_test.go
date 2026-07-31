package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/application"
)

func TestMapPricingErrorReturnsServerErrorForBatchLoadFailures(t *testing.T) {
	status, code := mapPricingError("PRICING_BATCH_LOAD_PRODUCTS: database unavailable")
	if status != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", status)
	}
	if code != "PRICING_INTERNAL_ERROR" {
		t.Fatalf("expected PRICING_INTERNAL_ERROR, got %q", code)
	}
}

// trimJSON normalises a JSON body for comparison (key order becomes
// insensitive) — must be applied to BOTH sides of a comparison, never just
// the actual body, or the expected literal has to be hand-alphabetised.
func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

// TestPricingErrorEnvelopeWholeBody pins the COMPLETE JSON body of a
// representative pricing error through apierror.Write, proving no stray
// top-level key survives the migration off the module-local envelope.
func TestPricingErrorEnvelopeWholeBody(t *testing.T) {
	handler := NewHandler(application.Service{}, nil)
	req := httptest.NewRequest(http.MethodPost, "/pricing/simulations", strings.NewReader("{not-json"))
	rr := httptest.NewRecorder()

	handler.handleSimulations(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body=%s", rr.Code, rr.Body.String())
	}
	want := `{"error":{"code":"PRICING_REQUEST_INVALID","message":"malformed request body","details":{}}}`
	if got := trimJSON(rr.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
