package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/dashboard/domain"
)

type summaryServiceStub struct {
	summary        domain.Summary
	err            error
	installationID string
	calls          int
}

func (s *summaryServiceStub) Summary(_ context.Context, installationID string) (domain.Summary, error) {
	s.calls++
	s.installationID = installationID
	return s.summary, s.err
}

func TestDashboardSummaryPartialFailureIs200(t *testing.T) {
	t.Parallel()

	service := &summaryServiceStub{
		summary: domain.Summary{
			PendingLinks: nil,
			Degraded:     []string{"linkage"},
			LastSyncAt:   map[string]*time.Time{"listings_refresh": nil},
			AsOf:         time.Date(2026, 7, 16, 15, 4, 5, 0, time.UTC),
		},
	}
	recorder := httptest.NewRecorder()

	NewHandler(service).HandleSummary(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary?installation_id=inst_001", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"pending_links":null`) {
		t.Fatalf("response = %s, want explicit pending_links null", recorder.Body.String())
	}
	var payload domain.Summary
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Degraded) != 1 || payload.Degraded[0] != "linkage" {
		t.Fatalf("degraded = %v, want [linkage]", payload.Degraded)
	}
}

func TestDashboardSummaryUnknownInstallationIs404(t *testing.T) {
	t.Parallel()

	service := &summaryServiceStub{err: &domain.InstallationNotFoundError{InstallationID: "missing"}}
	recorder := httptest.NewRecorder()

	NewHandler(service).HandleSummary(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary?installation_id=missing", nil))

	assertDashboardError(t, recorder, http.StatusNotFound, "installation_not_found")
}

func TestDashboardSummaryRequiresInstallation(t *testing.T) {
	t.Parallel()

	service := &summaryServiceStub{}
	recorder := httptest.NewRecorder()

	NewHandler(service).HandleSummary(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil))

	assertDashboardError(t, recorder, http.StatusBadRequest, "installation_required")
	if service.calls != 0 {
		t.Fatalf("service calls = %d, want 0", service.calls)
	}
}

func TestDashboardSummaryEmitsExplicitNullKeys(t *testing.T) {
	t.Parallel()

	service := &summaryServiceStub{summary: domain.Summary{
		LastSyncAt: map[string]*time.Time{"listings_refresh": nil},
		Degraded:   []string{},
		AsOf:       time.Date(2026, 7, 16, 15, 4, 5, 0, time.UTC),
	}}
	recorder := httptest.NewRecorder()

	NewHandler(service).HandleSummary(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary?installation_id=inst_001", nil))

	for _, key := range []string{
		`"sync_errors":null`, `"pending_links":null`, `"below_margin":null`,
		`"missing_gtin":null`, `"orders_today":null`, `"orders_7d":null`,
		`"listings_refresh":null`,
	} {
		if !strings.Contains(recorder.Body.String(), key) {
			t.Errorf("response = %s, want %s", recorder.Body.String(), key)
		}
	}
}

// trimJSON canonicalises JSON for whole-body comparisons (key order is
// otherwise a false failure), from catalog/transport/http_handler_test.go:387.
func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

// TestDashboardErrorEnvelopeWholeBody pins the complete JSON body of the
// installation_required case — the migrating ad hoc field this module
// carries ("key") — to prove it landed inside details and nowhere else,
// which a field-by-field .error.code read cannot prove.
func TestDashboardErrorEnvelopeWholeBody(t *testing.T) {
	t.Parallel()

	service := &summaryServiceStub{}
	recorder := httptest.NewRecorder()

	NewHandler(service).HandleSummary(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body=%s", recorder.Code, recorder.Body.String())
	}
	want := `{"error":{"code":"installation_required","message":"installation_id é obrigatório","details":{"key":"installation_id"}}}`
	if got := trimJSON(recorder.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func assertDashboardError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, status, recorder.Body.String())
	}
	var payload struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if payload.Error.Code != code {
		t.Fatalf("error code = %q, want %q", payload.Error.Code, code)
	}
}
