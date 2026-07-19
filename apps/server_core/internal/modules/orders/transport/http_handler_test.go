package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

func testTransportLogger() *slog.Logger { return slog.Default() }

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

func enrichedOrderModels() []domain.OrderReadModel {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	return []domain.OrderReadModel{
		{
			ProviderOrderID: "order-1", ProviderCode: "mercado_livre", Status: "paid",
			BuyerNickname: nil, CreatedAt: &now,
			Items: []domain.MarketplaceOrderItem{
				{ProviderItemID: "item-1", SellerSKU: "sku-1", LinkQuality: domain.LinkQualityUnresolved},
			},
			Payments: []domain.MarketplaceOrderPayment{},
		},
		{
			ProviderOrderID: "order-2", ProviderCode: "mercado_livre", Status: "paid",
			BuyerNickname: nil, CreatedAt: &now,
			Items:    []domain.MarketplaceOrderItem{{ProviderItemID: "item-2"}},
			Payments: []domain.MarketplaceOrderPayment{},
		},
	}
}

func TestHandleReadListWithEnricherEmitsHonestNullEnrichment(t *testing.T) {
	service := stubOrderReadService{page: ports.OrderPage{Items: enrichedOrderModels()}}
	enricher := application.NewEnrichService(nil, nil, testTransportLogger())
	handler := NewHandlerWithEnricher(stubOrderImporter{}, service, &enricher)
	req := httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, forbidden := range []string{"raw_provider_ref"} {
		if bytes.Contains([]byte(body), []byte(forbidden)) {
			t.Fatalf("body must not contain %q: %s", forbidden, body)
		}
	}
	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, body)
	}
	if len(payload.Items) != 2 {
		t.Fatalf("items = %d, want 2; body=%s", len(payload.Items), body)
	}
	item := payload.Items[0]
	if item["vinculo_status"] != "SEM_VINCULO" {
		t.Fatalf("vinculo_status = %#v, want SEM_VINCULO; body=%s", item["vinculo_status"], body)
	}
	buyer, ok := item["buyer"].(map[string]any)
	if !ok {
		t.Fatalf("buyer missing or wrong type: %#v; body=%s", item["buyer"], body)
	}
	if display, present := buyer["display"]; !present || display != "" {
		t.Fatalf("buyer.display = %#v, want empty string; body=%s", display, body)
	}
	for _, raw := range []string{"nickname", "cpf", "email", "phone", "city", "uf"} {
		if _, present := buyer[raw]; present && raw != "city" && raw != "uf" {
			t.Fatalf("buyer must not carry raw %q; body=%s", raw, body)
		}
	}
	if _, present := buyer["city"]; present {
		t.Fatalf("buyer.city must be omitted when unknown; body=%s", body)
	}
	if _, present := buyer["uf"]; present {
		t.Fatalf("buyer.uf must be omitted when unknown; body=%s", body)
	}
	componentes, ok := item["componentes_desconhecidos"].([]any)
	if !ok || len(componentes) != 1 || componentes[0] != "sku-1" {
		t.Fatalf("componentes_desconhecidos = %#v, want [sku-1]; body=%s", item["componentes_desconhecidos"], body)
	}
	if _, present := item["sla"]; present {
		t.Fatalf("sla must be omitted when there is no shipment; body=%s", body)
	}
	if _, present := item["rastreio"]; present {
		t.Fatalf("rastreio must be omitted when there is no shipment; body=%s", body)
	}
	if _, present := item["destino_uf"]; present {
		t.Fatalf("destino_uf must be omitted when there is no shipment; body=%s", body)
	}
	items, ok := item["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items[0].items = %#v, want 1 entry; body=%s", item["items"], body)
	}
	firstItem := items[0].(map[string]any)
	if _, present := firstItem["custo_unitario"]; present {
		t.Fatalf("custo_unitario must be omitted when cost is unknown; body=%s", body)
	}
	if firstItem["seller_sku"] != "sku-1" {
		t.Fatalf("items[0].seller_sku = %#v, want sku-1; body=%s", firstItem["seller_sku"], body)
	}
}

func TestHandleReadListEmitsDerivedBucket(t *testing.T) {
	service := stubOrderReadService{page: ports.OrderPage{Items: enrichedOrderModels()}}
	enricher := application.NewEnrichService(nil, nil, testTransportLogger())
	handler := NewHandlerWithEnricher(stubOrderImporter{}, service, &enricher)
	req := httptest.NewRequest(http.MethodGet, "/orders?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rr.Body.String())
	}
	if len(payload.Items) != 2 {
		t.Fatalf("items = %d, want 2; body=%s", len(payload.Items), rr.Body.String())
	}
	for _, item := range payload.Items {
		// enrichedOrderModels() has no shipment for either order, and both
		// have Status "paid", so DeriveOrderBucket must yield "faturar".
		if item["bucket"] != "faturar" {
			t.Fatalf("bucket = %#v, want faturar; body=%s", item["bucket"], rr.Body.String())
		}
	}
}

func TestMapEnrichedOrderSetsBucketFromShipmentPresence(t *testing.T) {
	order := enrichedOrderModels()[1]
	order.Status = "paid"

	withoutShipment := application.EnrichedOrder{Order: order}
	if dto := mapEnrichedOrder(withoutShipment); dto.Bucket != domain.BucketFaturar {
		t.Fatalf("bucket = %q, want %q (paid, no shipment)", dto.Bucket, domain.BucketFaturar)
	}

	withShipment := application.EnrichedOrder{
		Order:    order,
		Shipment: &application.ShipmentEnrichment{ShipmentID: "ship-1", Status: "shipped"},
	}
	if dto := mapEnrichedOrder(withShipment); dto.Bucket != domain.BucketEnviar {
		t.Fatalf("bucket = %q, want %q (paid, with shipment)", dto.Bucket, domain.BucketEnviar)
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(withShipment))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}
	if decoded["bucket"] != "enviar" {
		t.Fatalf("json field bucket = %#v, want \"enviar\"; body=%s", decoded["bucket"], dtoJSON)
	}
}

func TestHandleGetWithEnricherEmitsHonestNullEnrichment(t *testing.T) {
	model := enrichedOrderModels()[1]
	service := stubOrderReadService{model: model}
	enricher := application.NewEnrichService(nil, nil, testTransportLogger())
	handler := NewHandlerWithEnricher(stubOrderImporter{}, service, &enricher)
	req := httptest.NewRequest(http.MethodGet, "/orders/order-2?installation_id=inst-1", nil)
	req.SetPathValue("provider_order_id", "order-2")
	rr := httptest.NewRecorder()

	handler.handleGet(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rr.Body.String())
	}
	if payload["vinculo_status"] != "SEM_VINCULO" {
		t.Fatalf("vinculo_status = %#v, want SEM_VINCULO; body=%s", payload["vinculo_status"], rr.Body.String())
	}
	if _, present := payload["sla"]; present {
		t.Fatalf("sla must be omitted when there is no shipment; body=%s", rr.Body.String())
	}
}

// fakeSummaryStore mirrors application/summary_service_test.go's
// fakeOrderSummaryStore idiom (unexported there, so redeclared for the
// transport layer) to drive handleSummary through a real SummaryService
// without a database.
type fakeSummaryStore struct {
	summary ports.OrderSummary
	err     error
}

func (f *fakeSummaryStore) GetOrderSummary(context.Context, string, time.Time) (ports.OrderSummary, error) {
	return f.summary, f.err
}

// fakeBucketStore mirrors fakeSummaryStore's idiom for ports.OrderBucketStore
// (F01-A), driving handleSummary's ?by=status branch through a real
// SummaryService without a database.
type fakeBucketStore struct {
	counts ports.OrderBucketCounts
	err    error
}

func (f *fakeBucketStore) GetOrderBucketCounts(context.Context, string) (ports.OrderBucketCounts, error) {
	return f.counts, f.err
}

func TestHandleSummaryReturnsCounts(t *testing.T) {
	summarySvc := application.NewSummaryService(&fakeSummaryStore{summary: ports.OrderSummary{Today: 3, SevenDays: 11}})
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rr.Body.String())
	}
	if len(payload) != 2 {
		t.Fatalf("payload = %#v, want exactly today/seven_days keys", payload)
	}
	if payload["today"] != float64(3) || payload["seven_days"] != float64(11) {
		t.Fatalf("payload = %#v, want today=3 seven_days=11", payload)
	}
}

func TestHandleSummaryStoreErrorReturns500NoFabricatedZeros(t *testing.T) {
	summarySvc := application.NewSummaryService(&fakeSummaryStore{err: errors.New("summary source unavailable")})
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d body=%s, want 500", rr.Code, rr.Body.String())
	}
	if bytes.Contains(rr.Body.Bytes(), []byte(`"today"`)) {
		t.Fatalf("error body must not fabricate a zero summary; body=%s", rr.Body.String())
	}
}

func TestHandleSummaryMissingInstallationReturns400(t *testing.T) {
	summarySvc := application.NewSummaryService(&fakeSummaryStore{summary: ports.OrderSummary{Today: 1, SevenDays: 2}})
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", rr.Code, rr.Body.String())
	}
}

func TestHandleSummaryStoreNotConfiguredReturns503(t *testing.T) {
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, application.SummaryService{})
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s, want 503", rr.Code, rr.Body.String())
	}
}

func TestHandleSummaryByStatusReturnsBucketCounts(t *testing.T) {
	summarySvc := application.NewSummaryServiceWithBuckets(
		&fakeSummaryStore{summary: ports.OrderSummary{Today: 3, SevenDays: 11}},
		&fakeBucketStore{counts: ports.OrderBucketCounts{Novo: 4, Faturar: 2, Enviar: 1, Enviado: 7}},
	)
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1&by=status", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var payload struct {
		ByStatus map[string]int64 `json:"by_status"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rr.Body.String())
	}
	want := map[string]int64{"novo": 4, "faturar": 2, "enviar": 1, "enviado": 7}
	if len(payload.ByStatus) != len(want) {
		t.Fatalf("by_status = %#v, want exactly novo/faturar/enviar/enviado keys", payload.ByStatus)
	}
	for key, wantValue := range want {
		if payload.ByStatus[key] != wantValue {
			t.Fatalf("by_status[%q] = %d, want %d; body=%s", key, payload.ByStatus[key], wantValue, rr.Body.String())
		}
	}
}

func TestHandleSummaryByStatusStoreErrorReturns500NoFabricatedZeros(t *testing.T) {
	summarySvc := application.NewSummaryServiceWithBuckets(
		&fakeSummaryStore{summary: ports.OrderSummary{Today: 3, SevenDays: 11}},
		&fakeBucketStore{err: errors.New("bucket source unavailable")},
	)
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1&by=status", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d body=%s, want 500", rr.Code, rr.Body.String())
	}
	if bytes.Contains(rr.Body.Bytes(), []byte(`"by_status"`)) {
		t.Fatalf("error body must not fabricate zero bucket counts; body=%s", rr.Body.String())
	}
}

func TestHandleSummaryByStatusMissingInstallationReturns400(t *testing.T) {
	summarySvc := application.NewSummaryServiceWithBuckets(
		&fakeSummaryStore{summary: ports.OrderSummary{Today: 1, SevenDays: 2}},
		&fakeBucketStore{counts: ports.OrderBucketCounts{Novo: 1}},
	)
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?by=status", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", rr.Code, rr.Body.String())
	}
}

func TestHandleSummaryUnsupportedByValueReturns400(t *testing.T) {
	summarySvc := application.NewSummaryService(&fakeSummaryStore{summary: ports.OrderSummary{Today: 1, SevenDays: 2}})
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1&by=bogus", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", rr.Code, rr.Body.String())
	}
}

func TestHandleSummaryNoParamStillReturnsDefaultShape(t *testing.T) {
	summarySvc := application.NewSummaryServiceWithBuckets(
		&fakeSummaryStore{summary: ports.OrderSummary{Today: 3, SevenDays: 11}},
		&fakeBucketStore{counts: ports.OrderBucketCounts{Novo: 9}},
	)
	handler := NewHandlerWithSummary(stubOrderImporter{}, stubOrderReadService{}, nil, summarySvc)
	req := httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1", nil)
	rr := httptest.NewRecorder()

	handler.handleSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rr.Body.String())
	}
	if len(payload) != 2 {
		t.Fatalf("payload = %#v, want exactly today/seven_days keys (no-param call unchanged)", payload)
	}
	if payload["today"] != float64(3) || payload["seven_days"] != float64(11) {
		t.Fatalf("payload = %#v, want today=3 seven_days=11", payload)
	}
}

func TestRegisterWiresSummaryRouteWithoutShadowingDetailRoute(t *testing.T) {
	summarySvc := application.NewSummaryService(&fakeSummaryStore{summary: ports.OrderSummary{Today: 5, SevenDays: 9}})
	service := stubOrderReadService{model: enrichedOrderModels()[0]}
	mux := http.NewServeMux()
	NewHandlerWithSummary(stubOrderImporter{}, service, nil, summarySvc).Register(mux)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/orders/summary?installation_id=inst-1", nil))
	if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"today":5`)) {
		t.Fatalf("status=%d body=%s, want 200 with summary counts", rr.Code, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/orders/order-1?installation_id=inst-1", nil)
	mux.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"provider_order_id":"order-1"`)) {
		t.Fatalf("status=%d body=%s, want /orders/{id} route unaffected by /orders/summary", rr2.Code, rr2.Body.String())
	}
}
