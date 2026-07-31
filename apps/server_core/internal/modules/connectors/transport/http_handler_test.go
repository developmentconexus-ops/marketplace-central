package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// TestConnectorsErrorEnvelopeWholeBody pins the complete JSON body written
// when ME_CLIENT_ID is not configured. Before this chip's migration,
// writeConnectorsError wrote {"error":{"code":...,"message":...}} with NO
// details key at all — the site the chip's context pack names as the reason
// `.error.details` was `undefined` on some routes and an object on others.
// This pins that it now always carries details:{}.
func TestConnectorsErrorEnvelopeWholeBody(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/connectors/melhor-envio/auth/start", nil)
	rec := httptest.NewRecorder()

	h.handleMEAuthStart(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s, want 503", rec.Code, rec.Body.String())
	}
	want := `{"error":{"code":"CONNECTORS_ME_NOT_CONFIGURED","message":"Melhor Envio is not configured (ME_CLIENT_ID missing)","details":{}}}`
	if got := trimJSON(rec.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
