//go:build integration

package transport_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	syncpostgres "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	synctransport "marketplace-central/apps/server_core/internal/modules/sync/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func newHealthTestToken(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type healthEntityWire struct {
	Entity              string  `json:"entity"`
	LastSuccessAt       *string `json:"last_success_at"`
	LastIncrementalAt   *string `json:"last_incremental_at"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	Phase               *string `json:"phase"`
	LastError           *struct {
		Message string `json:"message"`
		At      string `json:"at"`
	} `json:"last_error"`
}

type healthResponseWire struct {
	Entities []healthEntityWire `json:"entities"`
	Webhook  struct {
		LastNotificationAt *string `json:"last_notification_at"`
		Pending            int     `json:"pending"`
		Dropped24h         int     `json:"dropped_24h"`
	} `json:"webhook"`
}

// Golden JSON fixture: one entity with a real success where last_full_sync_at
// is a stale one-time backfill and last_incremental_at is recent (F-r04-1,
// GREATEST semantics end-to-end through the real DB read), one entity with
// consecutive failures + last_error, and one entity name entirely ABSENT from
// sync_state — it must not appear in entities[] at all.
func TestSyncHealthHandlerGoldenFixture(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_health_golden")
	tenant := "tenant-health-golden-" + newHealthTestToken(t)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	staleBackfill := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	recentIncremental := time.Date(2026, 7, 31, 11, 0, 0, 0, time.UTC)
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor, last_full_sync_at, last_incremental_at, consecutive_failures)
		VALUES ($1, 'erp', 'products', '{"phase":"incremental"}'::jsonb, $2, $3, 0)
	`, tenant, staleBackfill, recentIncremental); err != nil {
		t.Fatal(err)
	}

	failAt := time.Date(2026, 7, 30, 8, 0, 0, 0, time.UTC)
	errPayload := fmt.Sprintf(`{"message":"provider timeout","at":%q}`, failAt.Format(time.RFC3339))
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, last_error, consecutive_failures)
		VALUES ($1, 'ml-installation-1', 'listings', $2::jsonb, 3)
	`, tenant, errPayload); err != nil {
		t.Fatal(err)
	}
	// 'orders' is intentionally absent from sync_state for this tenant.

	reader := syncpostgres.NewHealthReader(pool)
	svc := syncapp.NewHealthService(reader, tenant)
	handler := synctransport.NewHealthHandler(svc)
	mux := httpx.NewRouteClassMux()
	handler.Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var body healthResponseWire
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v (body=%s)", err, recorder.Body.String())
	}

	if len(body.Entities) != 2 {
		t.Fatalf("entities = %d, want 2 (orders is absent from sync_state and must not appear): %+v", len(body.Entities), body.Entities)
	}
	// entity ascending: listings, products
	listings, products := body.Entities[0], body.Entities[1]
	if listings.Entity != "listings" || products.Entity != "products" {
		t.Fatalf("entity order = [%s, %s], want [listings, products]", listings.Entity, products.Entity)
	}

	wantSuccess := recentIncremental.Format(time.RFC3339)
	if products.LastSuccessAt == nil || *products.LastSuccessAt != wantSuccess {
		t.Errorf("products.last_success_at = %v, want %s (GREATEST must pick the recent incremental, not the stale last_full_sync_at %s)",
			products.LastSuccessAt, wantSuccess, staleBackfill.Format(time.RFC3339))
	}
	if products.Phase == nil || *products.Phase != "incremental" {
		t.Errorf("products.phase = %v, want incremental", products.Phase)
	}
	if products.LastError != nil {
		t.Errorf("products.last_error = %+v, want nil", products.LastError)
	}

	if listings.ConsecutiveFailures != 3 {
		t.Errorf("listings.consecutive_failures = %d, want 3", listings.ConsecutiveFailures)
	}
	if listings.LastError == nil || listings.LastError.Message != "provider timeout" {
		t.Errorf("listings.last_error = %+v, want message=provider timeout", listings.LastError)
	}
	if listings.Phase != nil {
		t.Errorf("listings.phase = %v, want nil (no cursor was ever written)", listings.Phase)
	}
	if listings.LastSuccessAt != nil {
		t.Errorf("listings.last_success_at = %v, want nil (this entity never succeeded)", listings.LastSuccessAt)
	}

	if body.Webhook.LastNotificationAt != nil || body.Webhook.Pending != 0 || body.Webhook.Dropped24h != 0 {
		t.Errorf("webhook = %+v, want the canonical zero state (no M-08 reader wired in F-01)", body.Webhook)
	}
}

// A tenant with zero sync_state rows answers entities: [] (empty array),
// never null and never a 404.
func TestSyncHealthHandlerEmptyTenant(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "sync_health_empty")
	tenant := "tenant-health-empty-" + newHealthTestToken(t)

	reader := syncpostgres.NewHealthReader(pool)
	svc := syncapp.NewHealthService(reader, tenant)
	handler := synctransport.NewHealthHandler(svc)
	mux := httpx.NewRouteClassMux()
	handler.Register(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sync/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"entities":[]`) {
		t.Fatalf("body = %s, want entities:[] for an empty tenant, not null and not a 404", recorder.Body.String())
	}
}
