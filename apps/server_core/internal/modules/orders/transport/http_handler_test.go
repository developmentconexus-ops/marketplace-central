package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
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

type stubOrderReadService struct {
	page    ports.OrderPage
	listErr error
	model   domain.OrderReadModel
	getErr  error
}

func (s stubOrderReadService) List(context.Context, ports.OrderListQuery) (ports.OrderPage, error) {
	return s.page, s.listErr
}

func (s stubOrderReadService) Get(context.Context, string, string) (domain.OrderReadModel, error) {
	return s.model, s.getErr
}

func TestListOrdersReturnsCursorEnvelopeAndCanonicalProjection(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	service := stubOrderReadService{page: ports.OrderPage{
		Items: []domain.OrderReadModel{{
			ProviderOrderID: "order-1", ProviderCode: "mercado_livre", Status: "paid",
			BuyerNickname: nil, Total: nil, Currency: nil, Fulfillment: nil, NFState: nil,
			CreatedAt: &now, Items: []domain.MarketplaceOrderItem{}, Payments: []domain.MarketplaceOrderPayment{},
		}},
		NextCursor: &ports.OrderCursor{Timestamp: now, ProviderOrderID: "order-1"},
	}}
	handler := NewHandlerWithReader(stubOrderImporter{}, service)
	req := httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1&limit=1&status=paid&date_from=2026-07-01&q=order", nil)
	rr := httptest.NewRecorder()

	handler.handleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	items := body["items"].([]any)
	item := items[0].(map[string]any)
	for _, field := range []string{"buyer_nickname", "total", "currency", "fulfillment", "nf_state"} {
		if value, present := item[field]; !present || value != nil {
			t.Fatalf("field %q = %#v, present=%v; body=%s", field, value, present, rr.Body.String())
		}
	}
	if _, ok := body["next_cursor"].(string); !ok || body["next_cursor"] == "" {
		t.Fatalf("next_cursor=%#v body=%s", body["next_cursor"], rr.Body.String())
	}
}

func TestListOrdersMalformedCursorReturns400InvalidCursor(t *testing.T) {
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{})
	rr := httptest.NewRecorder()
	handler.handleList(rr, httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1&cursor=bad", nil))

	if rr.Code != http.StatusBadRequest || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":"invalid_cursor"`)) {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestListOrdersMalformedDateReturns400InvalidFilter(t *testing.T) {
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{})
	rr := httptest.NewRecorder()
	handler.handleList(rr, httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1&date_from=not-a-date", nil))

	if rr.Code != http.StatusBadRequest || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":"invalid_filter"`)) {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestGetUnknownOrderReturns404OrderNotFound(t *testing.T) {
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{getErr: &ports.OrderNotFoundError{InstallationID: "inst-1", ProviderOrderID: "missing"}})
	req := httptest.NewRequest(http.MethodGet, "/orders/missing?installation_id=inst-1", nil)
	req.SetPathValue("provider_order_id", "missing")
	rr := httptest.NewRecorder()

	handler.handleGet(rr, req)

	if rr.Code != http.StatusNotFound || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":"order_not_found"`)) {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestOrderTransportNeverSerializesRawProviderRef(t *testing.T) {
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{model: domain.OrderReadModel{ProviderOrderID: "order-1"}})
	req := httptest.NewRequest(http.MethodGet, "/orders/order-1?installation_id=inst-1", nil)
	req.SetPathValue("provider_order_id", "order-1")
	rr := httptest.NewRecorder()

	handler.handleGet(rr, req)

	if rr.Code != http.StatusOK || bytes.Contains(rr.Body.Bytes(), []byte("raw_provider_ref")) {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestListOrdersLegacyLimitOnlyStillWorks(t *testing.T) {
	handler := NewHandler(stubOrderImporter{}, stubOrderLister{})
	rr := httptest.NewRecorder()
	handler.handleList(rr, httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1&limit=1", nil))

	if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"items"`)) {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestNewHandlerLegacyDoesNotRegisterDetailRoute(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(stubOrderImporter{}, stubOrderLister{}).Register(mux)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/orders/order-1?installation_id=inst-1", nil))

	if rr.Code != http.StatusNotFound && rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d body=%s, want 404 or 405", rr.Code, rr.Body.String())
	}
}

func TestNewHandlerLegacyKeepsIgnoringEvolvedListFilters(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(stubOrderImporter{}, stubOrderLister{}).Register(mux)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1&status=paid&q=order", nil))

	if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"provider_order_id":"2001"`)) {
		t.Fatalf("status=%d body=%s, want legacy list response", rr.Code, rr.Body.String())
	}
}

func TestRequiredInstallationUsesTypedSentinel(t *testing.T) {
	_, err := requiredInstallation(url.Values{})
	if !errors.Is(err, ErrInstallationRequired) {
		t.Fatalf("error=%v, want errors.Is(..., ErrInstallationRequired)", err)
	}
}
