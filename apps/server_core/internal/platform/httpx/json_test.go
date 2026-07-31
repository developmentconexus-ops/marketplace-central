package httpx

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// trimJSON normalises a JSON body so that key ordering and insignificant
// whitespace do not affect comparison. Copied from
// internal/modules/catalog/transport/http_handler_test.go:387 — it is 6
// lines and crossing a package boundary for it is not worth an export.
func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func TestWriteJSON_EncodeFailureAnswersEnvelopeNotPlainText(t *testing.T) {
	w := httptest.NewRecorder()

	// map[string]any{"f": func(){}} is not encodable by encoding/json.
	WriteJSON(w, 200, map[string]any{"f": func() {}})

	if w.Code != 500 {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	want := trimJSON(`{"error":{"code":"internal_error","message":"failed to encode response","details":{}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
	for _, forbidden := range []string{"json", "unsupported", "func"} {
		if strings.Contains(w.Body.String(), forbidden) {
			t.Fatalf("body leaked internal encoder text %q: %s", forbidden, w.Body.String())
		}
	}
}
