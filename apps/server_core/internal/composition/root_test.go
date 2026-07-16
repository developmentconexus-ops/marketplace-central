package composition

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/platform/httpx"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

func TestFeeScheduleRoutesUseBatchDeadline(t *testing.T) {
	mux := httpx.NewRouteClassMux()
	registerBatchRoutes(mux)

	for _, path := range []string{"/admin/fee-schedules/sync", "/admin/fee-schedules/seed"} {
		path := path
		t.Run(path, func(t *testing.T) {
			var deadline time.Time
			var hasDeadline bool
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				deadline, hasDeadline = r.Context().Deadline()
				w.WriteHeader(http.StatusNoContent)
			})

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, nil))
			if recorder.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
			}
			if !hasDeadline {
				t.Fatal("batch route has no context deadline")
			}
			if remaining := time.Until(deadline); remaining < 119*time.Second || remaining > 120*time.Second {
				t.Fatalf("batch deadline remaining = %s, want approximately 120s", remaining)
			}
		})
	}
}

func TestRootRuntimeRegistersListingsReadRoutes(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	for _, tc := range []struct {
		path   string
		method string
		status int
	}{
		{path: "/listings", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/by-product", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/summary", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/not-a-listing-id", method: http.MethodGet, status: http.StatusNotFound},
	} {
		t.Run(tc.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(tc.method, tc.path, nil))
			if recorder.Code != tc.status {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tc.status, recorder.Body.String())
			}
		})
	}
}
