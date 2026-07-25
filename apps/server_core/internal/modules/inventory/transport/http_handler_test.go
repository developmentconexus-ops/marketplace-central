package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/application"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type stubRiskLister struct {
	items []domain.StockRiskListItem
}

func (s stubRiskLister) List(context.Context, application.ListStockRisksInput) ([]domain.StockRiskListItem, error) {
	return s.items, nil
}

type stubManualActionApplier struct{}

func (s stubManualActionApplier) Apply(context.Context, application.ApplyManualActionInput) (application.ApplyManualActionResult, error) {
	return application.ApplyManualActionResult{
		Action: domain.StockAction{
			ActionID: "act-1",
			State:    domain.StockActionStateApplied,
		},
		Risk: domain.StockRiskListItem{
			Identity: domain.ListingIdentity{
				InstallationID: "inst-1",
				ProviderItemID: "MLB123",
			},
			State: domain.StockRiskHealthy,
		},
	}, nil
}

func TestHandleStockRisksReturnsJSONList(t *testing.T) {
	handler := NewHandler(stubRiskLister{items: []domain.StockRiskListItem{{
		Identity: domain.ListingIdentity{
			InstallationID: "inst-1",
			ProviderItemID: "MLB123",
		},
		ProviderCode: "mercado_livre",
		State:        domain.StockRiskOversell,
		PolicyID:     "stock-seguro-default",
	}}}, stubManualActionApplier{})

	req := httptest.NewRequest(http.MethodGet, "/inventory/stock-risks?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleStockRisks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	var payload struct {
		Items []domain.StockRiskListItem `json:"items"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Items) != 1 || payload.Items[0].State != domain.StockRiskOversell {
		t.Fatalf("payload=%+v", payload.Items)
	}
}

func TestHandleStockRisksReportsNewestObservationAsOf(t *testing.T) {
	older := time.Date(2026, 7, 16, 17, 7, 2, 0, time.UTC)
	newer := time.Date(2026, 7, 18, 9, 30, 0, 0, time.UTC)
	handler := NewHandler(stubRiskLister{items: []domain.StockRiskListItem{
		{Identity: domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB1"}, ProviderObservedAt: &newer},
		{Identity: domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB2"}, ProviderObservedAt: &older, InternalObservedAt: &older},
	}}, stubManualActionApplier{})

	rr := httptest.NewRecorder()
	handler.handleStockRisks(rr, httptest.NewRequest(http.MethodGet, "/inventory/stock-risks?installation_id=inst-1", nil))

	var payload struct {
		AsOf string `json:"as_of"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got, want := payload.AsOf, newer.Format(time.RFC3339Nano); got != want {
		t.Fatalf("as_of=%q, want %q", got, want)
	}
}

func TestHandleStockRisksOmitsAsOfWithoutObservations(t *testing.T) {
	handler := NewHandler(stubRiskLister{items: []domain.StockRiskListItem{
		{Identity: domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB1"}},
	}}, stubManualActionApplier{})

	rr := httptest.NewRecorder()
	handler.handleStockRisks(rr, httptest.NewRequest(http.MethodGet, "/inventory/stock-risks?installation_id=inst-1", nil))

	var payload map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, present := payload["as_of"]; present {
		t.Fatalf("as_of present without any observation: %+v", payload)
	}
}

func TestHandleManualApplyReturnsActionResult(t *testing.T) {
	handler := NewHandler(stubRiskLister{}, stubManualActionApplier{})
	body := bytes.NewBufferString(`{"stock_action_id":"act-1","installation_id":"inst-1","provider_item_id":"MLB123","requested_quantity":7,"approval":{"approved":true,"approved_at":"2026-07-09T12:00:00Z","operator":{"actor_type":"operator","actor_id":"leandro"}}}`)
	req := httptest.NewRequest(http.MethodPost, "/inventory/stock-actions/manual-apply", body)
	rr := httptest.NewRecorder()

	handler.handleManualApply(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	var payload application.ApplyManualActionResult
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Action.ActionID != "act-1" || payload.Action.State != domain.StockActionStateApplied {
		t.Fatalf("action=%+v", payload.Action)
	}
	if payload.Risk.State != domain.StockRiskHealthy {
		t.Fatalf("risk=%+v", payload.Risk)
	}
}
