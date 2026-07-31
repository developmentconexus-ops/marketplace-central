package composition

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

// catalog/transport's Register is unconditional now — the fork that could mount
// a legacy UNCUT page (and once did so silently whenever the wiring forgot the
// page reader) was deleted with the legacy handlers. What is left to guard at
// this seam is the composition itself: NewRootRuntime must actually mount the
// paged catalog routes and wire them to a reader that degrades honestly. A
// property of the composition is only observable by mounting the composition,
// so this file mounts NewRootRuntime rather than a stand-in mux.
func TestRootRuntimeMountsPagedCatalogRoutes(t *testing.T) {
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

	// If the counts route were not registered, "GET /catalog/products/{id}"
	// would match with id="counts" and answer with an identity error, never a
	// 503 — so the status separates a mounted counts route from its absence.
	if code, body := serve(http.MethodGet, "/catalog/products/counts"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /catalog/products/counts = %d, want 503: the counts route is not registered, so "+
			"\"GET /catalog/products/{id}\" answered with id=\"counts\" (body=%s)", code, body)
	}

	// Without a pool the composition wires the unavailable reader, and the paged
	// list route reports that as 503 source_unavailable instead of dressing it
	// as an internal error. This is the honesty the cut route owes its callers.
	if code, body := serve(http.MethodGet, "/catalog/products"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /catalog/products = %d, want 503 source_unavailable (body=%s)", code, body)
	}

	// Same wiring, search side: the route must exist (a 404/405 here means the
	// search registration fell out of the composition) and must degrade through
	// the same unavailable reader once the q guard is satisfied.
	if code, body := serve(http.MethodGet, "/catalog/products/search?q=cuba"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /catalog/products/search?q=cuba = %d, want 503 source_unavailable (body=%s)", code, body)
	}
}
