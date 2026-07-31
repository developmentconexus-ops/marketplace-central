package apierror

import (
	"encoding/json"
	"net/http/httptest"
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

func TestWrite_FullBodyWithDetails(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, 400, "invalid_erp_source", "erp_source inválido: use xlsx ou catalogo_cliente", map[string]any{
		"allowed_range": "xlsx|catalogo_cliente",
	})

	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	want := trimJSON(`{"error":{"code":"invalid_erp_source","message":"erp_source inválido: use xlsx ou catalogo_cliente","details":{"allowed_range":"xlsx|catalogo_cliente"}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
}

func TestWrite_NilDetailsSerialisesAsEmptyObject(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, 500, "internal_error", "boom", nil)

	want := trimJSON(`{"error":{"code":"internal_error","message":"boom","details":{}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
}

func TestWrite_EmptyNonNilDetailsMatchesNilDetails(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, 500, "internal_error", "boom", map[string]any{})

	want := trimJSON(`{"error":{"code":"internal_error","message":"boom","details":{}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
}

func TestWrite_ContentTypeAndStatus(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, 409, "conflict", "conflito", map[string]any{})

	if w.Code != 409 {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
}

func TestWrite_MigratingAdHocFieldsRoundTripInsideDetailsOnly(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, 400, "bad_request", "requisição inválida", map[string]any{
		"allowed_range": "xlsx|catalogo_cliente",
		"limit":         100,
		"column":        "codprod",
		"import_id":     "abc-123",
		"protocol":      "sankhya",
	})

	want := trimJSON(`{"error":{"code":"bad_request","message":"requisição inválida","details":{"allowed_range":"xlsx|catalogo_cliente","column":"codprod","import_id":"abc-123","limit":100,"protocol":"sankhya"}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
}
