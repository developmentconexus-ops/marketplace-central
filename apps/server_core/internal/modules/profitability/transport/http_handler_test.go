package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type stubService struct{}
func (stubService) ImportMarginInputs(context.Context, profitabilityapp.ImportMarginInputsInput) (profitabilitydomain.ImportInputsResult, error) {
	return profitabilitydomain.ImportInputsResult{InstallationID: "inst-1", ImportedCount: 1}, nil
}
func (stubService) ListInputs(context.Context, string, int) ([]profitabilitydomain.MarginInput, error) {
	return []profitabilitydomain.MarginInput{{InstallationID: "inst-1", ProviderOrderID: "o1", Kind: profitabilitydomain.InputKindRevenue}}, nil
}
func (stubService) CreateManualAdjustment(context.Context, profitabilityapp.CreateManualAdjustmentInput) (profitabilitydomain.ManualAdjustment, error) {
	return profitabilitydomain.ManualAdjustment{AdjustmentID: "adj-1", InstallationID: "inst-1"}, nil
}
func (stubService) ListAdjustments(context.Context, string, int) ([]profitabilitydomain.ManualAdjustment, error) {
	return []profitabilitydomain.ManualAdjustment{{AdjustmentID: "adj-1", InstallationID: "inst-1"}}, nil
}
func (stubService) CalculateSnapshots(context.Context, string, int) (profitabilitydomain.CalculateSnapshotsResult, error) {
	return profitabilitydomain.CalculateSnapshotsResult{InstallationID: "inst-1", CalculatedCount: 1}, nil
}
func (stubService) ListSnapshots(context.Context, string, int) ([]profitabilitydomain.ProfitSnapshot, error) {
	return []profitabilitydomain.ProfitSnapshot{{InstallationID: "inst-1", ProviderOrderID: "o1", Scope: profitabilitydomain.InputScopeOrder, Quality: profitabilitydomain.ProfitSnapshotIncomplete}}, nil
}

func TestHandleImportReturnsJSON(t *testing.T) {
	h := NewHandler(stubService{})
	req := httptest.NewRequest(http.MethodPost, "/profitability/margin-inputs/import", bytes.NewBufferString(`{"installation_id":"inst-1","limit":1}`))
	rr := httptest.NewRecorder()
	h.handleImport(rr, req)
	if rr.Code != http.StatusOK { t.Fatalf("status=%d", rr.Code) }
}

func TestHandleListInputsReturnsItems(t *testing.T) {
	h := NewHandler(stubService{})
	req := httptest.NewRequest(http.MethodGet, "/profitability/margin-inputs?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()
	h.handleListInputs(rr, req)
	var payload struct{ Items []profitabilitydomain.MarginInput `json:"items"` }
	_ = json.NewDecoder(rr.Body).Decode(&payload)
	if len(payload.Items) != 1 { t.Fatalf("items=%d", len(payload.Items)) }
}

func TestHandleCreateAdjustmentReturnsAudit(t *testing.T) {
	h := NewHandler(stubService{})
	req := httptest.NewRequest(http.MethodPost, "/profitability/manual-adjustments", bytes.NewBufferString(`{"installation_id":"inst-1","provider_order_id":"o1","amount":12.5,"reason":"manual freight"}`))
	rr := httptest.NewRecorder()
	h.handleAdjustments(rr, req)
	if rr.Code != http.StatusOK { t.Fatalf("status=%d", rr.Code) }
}

func TestHandleCalculateSnapshotsReturnsResult(t *testing.T) {
	h := NewHandler(stubService{})
	req := httptest.NewRequest(http.MethodPost, "/profitability/profit-snapshots/calculate", bytes.NewBufferString(`{"installation_id":"inst-1","limit":1}`))
	rr := httptest.NewRecorder()
	h.handleCalculateSnapshots(rr, req)
	if rr.Code != http.StatusOK { t.Fatalf("status=%d", rr.Code) }
}

func TestHandleListSnapshotsReturnsItems(t *testing.T) {
	h := NewHandler(stubService{})
	req := httptest.NewRequest(http.MethodGet, "/profitability/profit-snapshots?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()
	h.handleListSnapshots(rr, req)
	var payload struct{ Items []profitabilitydomain.ProfitSnapshot `json:"items"` }
	_ = json.NewDecoder(rr.Body).Decode(&payload)
	if len(payload.Items) != 1 { t.Fatalf("items=%d", len(payload.Items)) }
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

type limitExceededService struct{ stubService }

func (limitExceededService) CalculateSnapshots(context.Context, string, int) (profitabilitydomain.CalculateSnapshotsResult, error) {
	return profitabilitydomain.CalculateSnapshotsResult{}, errors.New("limit_exceeded: maximum limit is 200")
}

// TestProfitabilityErrorEnvelopeWholeBody pins the complete JSON body of the
// limit_exceeded case — the migrating ad hoc field this module carries — to
// prove "limit" landed inside details and the envelope has no stray top-level
// key (a field-by-field read, like limit_test.go's, cannot prove that).
func TestProfitabilityErrorEnvelopeWholeBody(t *testing.T) {
	h := NewHandler(limitExceededService{})
	req := httptest.NewRequest(http.MethodPost, "/profitability/profit-snapshots/calculate", bytes.NewBufferString(`{"installation_id":"inst-1","limit":1}`))
	rr := httptest.NewRecorder()
	h.handleCalculateSnapshots(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422, body=%s", rr.Code, rr.Body.String())
	}
	want := `{"error":{"code":"limit_exceeded","message":"limit_exceeded: maximum limit is 200","details":{"limit":200}}}`
	if got := trimJSON(rr.Body.String()); got != trimJSON(want) {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
