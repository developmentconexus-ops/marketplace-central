package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type stubOrderImporter struct{}

func (stubOrderImporter) Import(context.Context, application.ImportOrdersInput) (domain.ImportResult, error) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	return domain.ImportResult{
		InstallationID: "inst-1",
		ImportedCount:  1,
		Items: []domain.MarketplaceOrder{{
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderOrderID:   "2001",
			ProviderStatus:    "paid",
			ProviderUpdatedAt: &now,
			FetchedAt:         now,
		}},
	}, nil
}

type stubOrderLister struct{}

func (stubOrderLister) List(context.Context, application.ListOrdersInput) ([]domain.MarketplaceOrder, error) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	return []domain.MarketplaceOrder{{
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   "2001",
		ProviderStatus:    "paid",
		ProviderUpdatedAt: &now,
		FetchedAt:         now,
	}}, nil
}

func TestHandleImportReturnsResult(t *testing.T) {
	handler := NewHandler(stubOrderImporter{}, stubOrderLister{})
	req := httptest.NewRequest(http.MethodPost, "/orders/import", bytes.NewBufferString(`{"installation_id":"inst-1","limit":1}`))
	rr := httptest.NewRecorder()

	handler.handleImport(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	var payload domain.ImportResult
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.ImportedCount != 1 || payload.Items[0].ProviderOrderID != "2001" {
		t.Fatalf("payload=%+v", payload)
	}
}

func TestHandleListReturnsOrders(t *testing.T) {
	handler := NewHandler(stubOrderImporter{}, stubOrderLister{})
	req := httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	var payload struct {
		Items []domain.MarketplaceOrder `json:"items"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Items) != 1 || payload.Items[0].ProviderOrderID != "2001" {
		t.Fatalf("payload=%+v", payload)
	}
}
