package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// TestCatalogCountsInvalidErpSourceEmitsErrorEnvelope pins the FUTURE unified
// error contract: today handleCounts writes the flat legacy body
// {"error":"invalid_erp_source","allowed_range":"xlsx|catalogo_cliente"} via
// writeCatalogPageError; after the chip it must write the universal envelope
// with a human message and a details object. This test is expected to be RED
// until that migration lands — red is the deliverable for this slice.
func TestCatalogCountsInvalidErpSourceEmitsErrorEnvelope(t *testing.T) {
	fake := &fakeCatalogPageReader{}
	mux := httpx.NewRouteClassMux()
	(Handler{PageReader: fake}).Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/catalog/products/counts?erp_source=banana", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", recorder.Code, recorder.Body.String())
	}
	want := `{"error":{"code":"invalid_erp_source","message":"erp_source inválido: use xlsx ou catalogo_cliente","details":{"allowed_range":"xlsx|catalogo_cliente"}}}`
	if trimJSON(recorder.Body.String()) != trimJSON(want) {
		t.Fatalf("body = %s, want %s", trimJSON(recorder.Body.String()), trimJSON(want))
	}
	if len(fake.listCursors) != 0 || len(fake.searchQueries) != 0 {
		t.Fatalf("a validation error must not reach the reader: listCursors=%+v searchQueries=%+v", fake.listCursors, fake.searchQueries)
	}
}
