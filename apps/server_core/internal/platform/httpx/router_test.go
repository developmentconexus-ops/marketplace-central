package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// The web client sends Cache-Control: no-cache when the operator presses a
// Refresh control, and the read handlers honour it to bypass their caches. That
// header is not CORS-simple, so a preflight that does not allow it makes the
// browser block the actual request — the Refresh button then fails silently in
// a real browser while every mocked-fetch test passes.
func TestPreflightAllowsTheRefreshCacheHeaders(t *testing.T) {
	handler := httpx.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/catalog/products?limit=50", nil)
	req.Header.Set("Origin", "http://localhost:5174")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "cache-control")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight status=%d want %d", rec.Code, http.StatusNoContent)
	}
	allowed := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	for _, header := range []string{"content-type", "authorization", "cache-control", "pragma"} {
		if !strings.Contains(allowed, header) {
			t.Fatalf("Access-Control-Allow-Headers=%q does not allow %q", allowed, header)
		}
	}
}
