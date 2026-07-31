package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

// errorEnvelopeShape decodes the apierror.Write envelope for assertions that
// need the code/details fields.
type errorEnvelopeShape struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	} `json:"error"`
}

func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func TestIntegrationsErrorEnvelopeWholeBody(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	NewRunReadHandler(&stubRunReadService{}).HandleList(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs?installation_id=inst_001&cursor=bad", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
	want := `{"error":{"code":"invalid_cursor","message":"cursor inválido","details":{"key":"cursor"}}}`
	if trimJSON(recorder.Body.String()) != trimJSON(want) {
		t.Fatalf("body = %s, want %s", trimJSON(recorder.Body.String()), trimJSON(want))
	}
}

type stubRunReadService struct {
	page  ports.RunPage
	query ports.RunListQuery
}

func (s *stubRunReadService) ListRuns(_ context.Context, query ports.RunListQuery) (ports.RunPage, error) {
	s.query = query
	return s.page, nil
}

func TestListSyncRunsRunningHasNullFinishedAt(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	service := &stubRunReadService{page: ports.RunPage{Items: []ports.RunReadModel{{
		OperationRunID: "run_running",
		InstallationID: "inst_001",
		Module:         "listing_read",
		Status:         domain.OperationRunStatusRunning,
		StartedAt:      &startedAt,
		FinishedAt:     nil,
	}}}}
	recorder := httptest.NewRecorder()

	NewRunReadHandler(service).HandleList(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs?installation_id=inst_001", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}
	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	finishedAt, present := payload.Items[0]["finished_at"]
	if !present || finishedAt != nil {
		t.Fatalf("finished_at = %#v (present=%v), want explicit null", finishedAt, present)
	}
}

func TestListSyncRunsInvalidCursor(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	NewRunReadHandler(&stubRunReadService{}).HandleList(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs?installation_id=inst_001&cursor=bad", nil))

	assertRunError(t, recorder, http.StatusBadRequest, "invalid_cursor")
}

func TestListSyncRunsInvalidStatusOrModule(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	NewRunReadHandler(&stubRunReadService{}).HandleList(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs?installation_id=inst_001&filter.status=bogus", nil))

	assertRunError(t, recorder, http.StatusBadRequest, "invalid_filter")
}

func TestListSyncRunsRequiresInstallation(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	NewRunReadHandler(&stubRunReadService{}).HandleList(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs", nil))

	assertRunError(t, recorder, http.StatusBadRequest, "installation_required")
}

func assertRunError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
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
