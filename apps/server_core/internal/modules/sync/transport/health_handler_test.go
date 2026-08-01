package transport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/application"
	synctransport "marketplace-central/apps/server_core/internal/modules/sync/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type fakeHealthReader struct {
	rows []application.HealthEntityRow
}

var _ application.HealthReader = (*fakeHealthReader)(nil)

func (f *fakeHealthReader) ReadAll(context.Context, string) ([]application.HealthEntityRow, error) {
	return f.rows, nil
}

type fakeWebhookStatsReader struct {
	stats application.WebhookStats
}

var _ application.WebhookStatsReader = (*fakeWebhookStatsReader)(nil)

func (f *fakeWebhookStatsReader) WebhookStats(context.Context, string) (application.WebhookStats, error) {
	return f.stats, nil
}

// The default webhook block is byte-for-byte the canonical zero state — no
// fabricated timestamp, no error.
func TestHealthHandlerDefaultWebhookBlockIsByteExact(t *testing.T) {
	svc := application.NewHealthService(&fakeHealthReader{}, "tenant-1")
	handler := synctransport.NewHealthHandler(svc)

	mux := httpx.NewRouteClassMux()
	handler.Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"webhook":{"last_notification_at":null,"pending":0,"dropped_24h":0}`) {
		t.Fatalf("body = %s, want the canonical webhook block byte-for-byte", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"entities":[]`) {
		t.Fatalf("body = %s, want entities:[] for a tenant with no rows", recorder.Body.String())
	}
}

// TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive is the
// highest-scrutiny test in F-01 (IC-05): a future milestone will call
// WithWebhookStatsReader on the ALREADY-CONSTRUCTED, ALREADY-Register()'d
// service. The swap must be visible on the next hit of the SAME live route —
// not require re-registering a new handler. A value-receiver builder that
// returned an updated copy (like pricing transport's WithCalc idiom) would
// fail this test: the handler would keep serving the pre-injection default
// forever, because it only ever held the original *HealthService pointer.
func TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive(t *testing.T) {
	svc := application.NewHealthService(&fakeHealthReader{}, "tenant-1")
	handler := synctransport.NewHealthHandler(svc)

	mux := httpx.NewRouteClassMux()
	handler.Register(mux) // route is live before injection, exactly like production root.go

	firstResponse := httptest.NewRecorder()
	mux.ServeHTTP(firstResponse, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if !strings.Contains(firstResponse.Body.String(), `"pending":0`) {
		t.Fatalf("pre-injection body = %s, want default pending=0", firstResponse.Body.String())
	}

	notified := time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC)
	fake := &fakeWebhookStatsReader{stats: application.WebhookStats{
		LastNotificationAt: &notified,
		Pending:            7,
		Dropped24h:         2,
	}}
	// Injected AFTER Register — this is the exact ordering a future
	// milestone's own root.go region will use against the already-live route.
	svc.WithWebhookStatsReader(fake)

	secondResponse := httptest.NewRecorder()
	mux.ServeHTTP(secondResponse, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if secondResponse.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", secondResponse.Code, secondResponse.Body.String())
	}

	var body struct {
		Webhook struct {
			LastNotificationAt *string `json:"last_notification_at"`
			Pending            int     `json:"pending"`
			Dropped24h         int     `json:"dropped_24h"`
		} `json:"webhook"`
	}
	if err := json.Unmarshal(secondResponse.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v (body=%s)", err, secondResponse.Body.String())
	}
	if body.Webhook.Pending != 7 || body.Webhook.Dropped24h != 2 {
		t.Fatalf("post-injection webhook = %+v, want the fake's pending=7 dropped_24h=2 — the SAME live route must reflect the swap", body.Webhook)
	}
	if body.Webhook.LastNotificationAt == nil || *body.Webhook.LastNotificationAt != notified.Format(time.RFC3339) {
		t.Fatalf("post-injection last_notification_at = %v, want %v", body.Webhook.LastNotificationAt, notified.Format(time.RFC3339))
	}
}

// Negative control (IC-05 registration requirement): /sync/health must be
// INTERACTIVE (15s budget), never BATCH (120s) — it is a single tenant-scoped
// read, not a bulk operation, and registerBatchRoutes in composition/root.go
// must never list it. RouteClassMux keeps its class table private, so the
// discriminating observable is the actual context deadline the handler
// receives — mirrors httpx's own
// TestRouteClassMuxAppliesDeclaredClassToMethodQualifiedPattern.
func TestHealthHandlerRegistersInteractiveNotBatch(t *testing.T) {
	var budget time.Duration
	probe := deadlineProbeService{onCall: func(ctx context.Context) {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("handler received no context deadline at all")
		}
		budget = time.Until(deadline)
	}}
	handler := synctransport.NewHealthHandler(probe)

	mux := httpx.NewRouteClassMux()
	handler.Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}

	const batchDeadline = 120 * time.Second
	const interactiveDeadline = 15 * time.Second
	if budget >= batchDeadline {
		t.Fatalf("effective budget = %s, this IS the 120s batch deadline — /sync/health must never be batch-classed", budget)
	}
	if budget > interactiveDeadline || budget < interactiveDeadline/2 {
		t.Fatalf("effective budget = %s, want close to the %s interactive deadline", budget, interactiveDeadline)
	}
}

type deadlineProbeService struct {
	onCall func(ctx context.Context)
}

func (p deadlineProbeService) GetHealth(ctx context.Context) (application.HealthResult, error) {
	p.onCall(ctx)
	return application.HealthResult{Entities: []application.HealthEntity{}, Webhook: application.WebhookStats{}}, nil
}
