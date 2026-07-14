package composition

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/platform/httpx"
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
