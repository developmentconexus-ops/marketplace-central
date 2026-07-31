package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// trimJSON re-serializes body through encoding/json so a map-key-order
// difference cannot fail an otherwise byte-identical comparison. Applied to
// BOTH SIDES of a whole-body pin (see catalog/transport/http_handler_test.go
// for the idiom this is copied from) — applying it to only one side would
// force a hand-alphabetised expected literal and hide a stray top-level key.
func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

// TestMutationsErrorEnvelopeWholeBody pins the complete JSON body of a
// mutations transport error. This module is the family-A reference the rest
// of the chip's wire shape was modelled on: the point of this test is that
// the shape is unchanged after routing through apierror.Write instead of the
// module's own local mutationErrorEnvelope/mutationError structs.
func TestMutationsErrorEnvelopeWholeBody(t *testing.T) {
	h := NewHandler(&commandServiceFake{}, &approvalServiceFake{}, &retryServiceFake{})
	req := httptest.NewRequest(http.MethodPost, "/mutations", strings.NewReader(`{`))
	w := httptest.NewRecorder()
	mux := http.NewServeMux()
	h.Register(mux)
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", w.Code, w.Body.String())
	}
	want := `{"error":{"code":"invalid_body","message":"corpo da requisição inválido","details":{}}}`
	if got := trimJSON(w.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
