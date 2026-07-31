package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

// TestClassificationsErrorEnvelopeWholeBody pins the COMPLETE JSON body of a
// representative classifications error through apierror.Write, proving no
// stray top-level key survives the migration off the module-local envelope.
func TestClassificationsErrorEnvelopeWholeBody(t *testing.T) {
	handler := Handler{}
	req := httptest.NewRequest(http.MethodPost, "/classifications", strings.NewReader("{not-json"))
	rr := httptest.NewRecorder()

	handler.handleCollection(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body=%s", rr.Code, rr.Body.String())
	}
	want := `{"error":{"code":"CLASSIFICATIONS_CREATE_INVALID","message":"malformed request body","details":{}}}`
	if got := trimJSON(rr.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
