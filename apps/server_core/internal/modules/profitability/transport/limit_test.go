package transport

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type limitErrorService struct{}

func (limitErrorService) ImportMarginInputs(context.Context, profitabilityapp.ImportMarginInputsInput) (profitabilitydomain.ImportInputsResult, error) {
	return profitabilitydomain.ImportInputsResult{}, errors.New("limit_exceeded: maximum limit is 200")
}
func (limitErrorService) ListInputs(context.Context, string, int) ([]profitabilitydomain.MarginInput, error) {
	return nil, nil
}
func (limitErrorService) CreateManualAdjustment(context.Context, profitabilityapp.CreateManualAdjustmentInput) (profitabilitydomain.ManualAdjustment, error) {
	return profitabilitydomain.ManualAdjustment{}, nil
}
func (limitErrorService) ListAdjustments(context.Context, string, int) ([]profitabilitydomain.ManualAdjustment, error) {
	return nil, nil
}
func (limitErrorService) CalculateSnapshots(context.Context, string, int) (profitabilitydomain.CalculateSnapshotsResult, error) {
	return profitabilitydomain.CalculateSnapshotsResult{}, nil
}
func (limitErrorService) ListSnapshots(context.Context, string, int) ([]profitabilitydomain.ProfitSnapshot, error) {
	return nil, nil
}

func TestHandleImportLimitExceededReturns422AndNamesCap(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/profitability/margin-inputs/import", bytes.NewBufferString(`{"installation_id":"inst-1","limit":201}`))
	recorder := httptest.NewRecorder()
	NewHandler(limitErrorService{}).handleImport(recorder, req)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", recorder.Code)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, `"error":"limit_exceeded"`) || !strings.Contains(body, `"limit":200`) {
		t.Fatalf("body = %s, want limit_exceeded and limit 200", body)
	}
}
