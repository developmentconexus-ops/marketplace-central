package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRouteClassDeadlinesAndCancellation(t *testing.T) {
	if got := routeClassDeadline(InteractiveRouteClass); got != 15*time.Second {
		t.Fatalf("interactive deadline = %s, want 15s", got)
	}
	if got := routeClassDeadline(BatchRouteClass); got != 120*time.Second {
		t.Fatalf("batch deadline = %s, want 120s", got)
	}

	const stalledFor = 30 * time.Millisecond
	tests := []struct {
		name       string
		class      RouteClass
		budget     time.Duration
		wantStatus int
		wantBody   string
	}{
		{name: "interactive expires", class: InteractiveRouteClass, budget: 15 * time.Millisecond, wantStatus: http.StatusGatewayTimeout, wantBody: `{"error":"deadline_exceeded"}`},
		{name: "batch keeps its longer budget", class: BatchRouteClass, budget: 120 * time.Millisecond, wantStatus: http.StatusOK, wantBody: `{"status":"completed"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cancelled := make(chan struct{})
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case <-time.After(stalledFor):
					WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
				case <-r.Context().Done():
					close(cancelled)
				}
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			startedAt := time.Now()
			deadlineMiddlewareWithTimeout(tc.budget, handler).ServeHTTP(recorder, request)
			elapsed := time.Since(startedAt)

			if tc.wantStatus == http.StatusGatewayTimeout {
				if elapsed < tc.budget {
					t.Fatalf("returned after %s, before %s budget", elapsed, tc.budget)
				}
				select {
				case <-cancelled:
				case <-time.After(time.Second):
					t.Fatal("handler context was not canceled")
				}
			} else if elapsed < stalledFor {
				t.Fatalf("returned after %s, before stalled handler completed", elapsed)
			}
			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tc.wantStatus)
			}
			if recorder.Body.String() != tc.wantBody+"\n" && recorder.Body.String() != tc.wantBody {
				t.Fatalf("body = %q, want %q", recorder.Body.String(), tc.wantBody)
			}
		})
	}
}

func TestRouteClassMuxDeclaresBatchBeforeTransportRegistration(t *testing.T) {
	mux := NewRouteClassMux()
	mux.RegisterRouteClass("/batch", BatchRouteClass)
	mux.HandleFunc("/batch", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Deadline(); !ok {
			t.Error("batch route has no context deadline")
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
	})

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/batch", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
}

// Transports declare the class under the bare path and register the ServeMux
// pattern with its method, so a lookup keyed on the raw pattern misses every
// time and the route silently falls back to the interactive default. Asserting
// that a deadline EXISTS cannot see it — interactive has one too — so this
// measures the budget the handler actually received.
func TestRouteClassMuxAppliesDeclaredClassToMethodQualifiedPattern(t *testing.T) {
	mux := NewRouteClassMux()
	mux.RegisterRouteClass("/erp/imports", BatchRouteClass)

	var budget time.Duration
	mux.HandleFunc("POST /erp/imports", func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Deadline()
		if !ok {
			t.Error("batch route has no context deadline")
		}
		budget = time.Until(deadline)
		WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
	})

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/erp/imports", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if budget <= interactiveDeadline {
		t.Fatalf("effective budget = %s, want the %s batch budget (interactive default is %s)", budget, batchDeadline, interactiveDeadline)
	}
}
