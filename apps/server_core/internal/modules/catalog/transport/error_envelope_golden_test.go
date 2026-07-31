package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// TestCatalogCountsInvalidErpSourceEmitsErrorEnvelope pins the unified error
// contract on the chip's EXEMPLO-IO route. It was written RED against the flat
// legacy body {"error":"invalid_erp_source","allowed_range":"..."} that
// handleCounts used to write, and went green when writeCatalogPageError moved
// to platform/apierror.
//
// It pins the WHOLE body, not a substring: a substring assertion cannot tell a
// correct envelope from one carrying leftover sibling fields. Re-proven in
// VC-7(a) — reintroducing the flat body makes this test, plus
// TestCatalogPageRoutesValidateBeforePortCall and
// TestCatalogPageRejectsUnknownActiveSource, fail by name.
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
