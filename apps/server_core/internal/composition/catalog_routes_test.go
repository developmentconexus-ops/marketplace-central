package composition

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

// The cut that keeps unsellable products off the catalog screen lives behind the
// route registration, not inside it. catalog/transport's Register picks between
// two registrations:
//
//	if h.PageReader == nil && !hasRouteClasses { legacy, UNCUT page } else { cut page + counts }
//
// Nothing asserted which branch production takes, and a property of the
// composition is only observable by mounting the composition — so this file
// mounts NewRootRuntime rather than a stand-in mux.
//
// The first version of this test probed for 404 on the counts path and for 405 on
// a POST to the list path. Both must-fails passed GREEN against it, because both
// probes answer identically in both worlds:
//
//   - "/catalog/products/counts" is matched by "GET /catalog/products/{id}" with
//     id="counts" when the counts route is absent, so the path is never a 404 —
//     it is answered by the wrong handler.
//   - handleLegacyProducts performs its OWN method check and replies 405, so a
//     405 on POST proves nothing about which pattern is mounted.
//
// What actually separates the two worlds is measured below and is never the
// status alone: the cut page degrades honestly through the unavailable reader
// (503 source_unavailable), the legacy page does not, and the router's own 405 is
// plain text where the legacy handler's is JSON.
func TestRootRuntimeDoesNotRegisterLegacyUncutCatalogRoutes(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	serve := func(method, path string) (int, string) {
		recorder := httptest.NewRecorder()
		runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
		return recorder.Code, recorder.Body.String()
	}

	// The legacy branch registers no counts route, and "counts" then reads as a
	// product id. The tell is the identity error, not a 404.
	if code, body := serve(http.MethodGet, "/catalog/products/counts"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /catalog/products/counts = %d, want 503: the counts route is not registered, so "+
			"\"GET /catalog/products/{id}\" answered with id=\"counts\" — the composition took catalog/transport's "+
			"legacy branch and the catalog page is served UNCUT (body=%s)", code, body)
	}

	// The cut list route reaches the canonical reader, which is unavailable
	// without a pool and says so. The legacy list route fails as an internal error
	// instead, which is the same page without the cut and without the honesty.
	if code, body := serve(http.MethodGet, "/catalog/products"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /catalog/products = %d, want 503 source_unavailable: this is catalog/transport's legacy "+
			"UNCUT list handler, which reports an unavailable source as an internal error (body=%s)", code, body)
	}

	// A method-less legacy pattern mounted alongside the cut route is invisible to
	// GET, because "GET /catalog/products" is the more specific pattern. It is only
	// visible on a method the cut route does not serve — and then the status is 405
	// either way. Only the body says who answered: the router writes plain text,
	// handleLegacyProducts writes the module's JSON envelope.
	code, body := serve(http.MethodPost, "/catalog/products")
	if code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /catalog/products = %d, want 405 (body=%s)", code, body)
	}
	if strings.Contains(body, "CATALOG_METHOD_NOT_ALLOWED") {
		t.Fatalf("POST /catalog/products was answered by handleLegacyProducts, not by the router: a method-less "+
			"legacy pattern is mounted on the catalog list path (body=%s)", body)
	}
}
