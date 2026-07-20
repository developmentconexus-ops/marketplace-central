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

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
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
	for _, field := range []string{"retorno_liquido", "margem_pct"} {
		if _, present := item[field]; present {
			t.Fatalf("%s must be omitted when unknown; body=%s", field, body)
		}
	}
	decomposicao, ok := item["decomposicao"].(map[string]any)
	if !ok {
		t.Fatalf("decomposicao missing or wrong type: %#v; body=%s", item["decomposicao"], body)
	}
	for _, field := range []string{"comissao", "taxa_fixa", "frete", "imposto", "difal", "tarifa_full", "custo", "margem_valor", "margem_pct"} {
		if _, present := decomposicao[field]; present {
			t.Fatalf("decomposicao.%s must be omitted when unknown; body=%s", field, body)
		}
	}
	componentesDesconhecidos, ok := decomposicao["componentes_desconhecidos"].([]any)
	if !ok {
		t.Fatalf("decomposicao.componentes_desconhecidos missing or wrong type: %#v; body=%s", decomposicao["componentes_desconhecidos"], body)
	}
	wantComponents := []string{"comissao", "taxa_fixa", "frete", "imposto", "difal", "tarifa_full", "custo"}
	if len(componentesDesconhecidos) != len(wantComponents) {
		t.Fatalf("decomposicao.componentes_desconhecidos = %#v, want %v; body=%s", componentesDesconhecidos, wantComponents, body)
	}
	for i, want := range wantComponents {
		if componentesDesconhecidos[i] != want {
			t.Fatalf("decomposicao.componentes_desconhecidos[%d] = %#v, want %q; body=%s", i, componentesDesconhecidos[i], want, body)
		}
	}
	difal, ok := item["difal"].(map[string]any)
	if !ok {
		t.Fatalf("difal missing or wrong type: %#v; body=%s", item["difal"], body)
	}
	for _, field := range []string{"amount", "uf_route", "due_date", "paid"} {
		if _, present := difal[field]; present {
			t.Fatalf("difal.%s must be omitted when unknown; body=%s", field, body)
		}
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

func TestMapEnrichedOrderDerivesBucketFromFaturadoAt(t *testing.T) {
	order := enrichedOrderModels()[1]
	order.Status = "paid"
	order.ShippingID = ""
	order.FaturadoAt = nil

	notFaturado := application.EnrichedOrder{Order: order}
	if dto := mapEnrichedOrder(notFaturado); dto.Bucket != domain.BucketFaturar {
		t.Fatalf("bucket = %q, want %q (paid, not faturado)", dto.Bucket, domain.BucketFaturar)
	}

	// shipping_id column presence alone no longer flips the bucket (goal B):
	// only the persisted faturado_at fact does. The two paths (this enriched
	// mapper and the SQL summary by_status path in order_repo.go) both key off
	// faturado_at, so they must never diverge on the same order (P6 dual-gate
	// finding carried forward).
	shippingOnly := order
	shippingOnly.ShippingID = "ship-col-1"
	withShippingNoFaturado := application.EnrichedOrder{Order: shippingOnly}
	if dto := mapEnrichedOrder(withShippingNoFaturado); dto.Bucket != domain.BucketFaturar {
		t.Fatalf("bucket = %q, want %q (paid, shipping_id present but not faturado)", dto.Bucket, domain.BucketFaturar)
	}

	faturadoAt := time.Date(2026, 7, 20, 9, 0, 0, 0, time.UTC)
	faturado := order
	faturado.FaturadoAt = &faturadoAt
	withFaturado := application.EnrichedOrder{Order: faturado}
	if dto := mapEnrichedOrder(withFaturado); dto.Bucket != domain.BucketEnviar {
		t.Fatalf("bucket = %q, want %q (paid, faturado_at set)", dto.Bucket, domain.BucketEnviar)
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(withFaturado))
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

// TestMapEnrichedOrderSurfacesShipmentDestinoCarrierFrete asserts the round-4
// additive surfacing: the shipment destination (cep/destinatario), the carrier
// (rastreio.transportadora/url_rastreio) and the real freight actuals
// (frete_real.bruto/receiver/sender) all reach the JSON DTO when the shipment
// enrichment carries them. Every field is additive/omitempty (ADR-17).
func TestMapEnrichedOrderSurfacesShipmentDestinoCarrierFrete(t *testing.T) {
	order := enrichedOrderModels()[1]
	order.Status = "shipped"
	order.ShippingID = "ship-col-1"

	zip := "20040-002"
	receiver := "João Silva"
	city := "Rio de Janeiro"
	uf := "RJ"
	carrier := "Total Express"
	tracking := "http://tracking.totalexpress.com.br/poupup_track.php?reid=3"

	e := application.EnrichedOrder{
		Order: order,
		Buyer: domain.MaskedBuyer{City: &city, UF: &uf},
		Shipment: &application.ShipmentEnrichment{
			ShipmentID:      "ship-col-1",
			Status:          "shipped",
			DestinationUF:   &uf,
			DestinationCity: &city,
			DestinationZip:  &zip,
			ReceiverName:    &receiver,
			CarrierName:     &carrier,
			TrackingURL:     &tracking,
			Costs: &connectorsdomain.ShipmentCosts{
				GrossAmount:  &connectorsdomain.Money{Amount: "23.45", Currency: "BRL"},
				ReceiverCost: &connectorsdomain.Money{Amount: "10.00", Currency: "BRL"},
				SenderCost:   &connectorsdomain.Money{Amount: "13.45", Currency: "BRL"},
			},
		},
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}

	if decoded["destino_cep"] != "20040-002" {
		t.Fatalf("destino_cep = %#v, want 20040-002; body=%s", decoded["destino_cep"], dtoJSON)
	}
	if decoded["destinatario"] != "João Silva" {
		t.Fatalf("destinatario = %#v, want João Silva; body=%s", decoded["destinatario"], dtoJSON)
	}
	if decoded["destino_uf"] != "RJ" {
		t.Fatalf("destino_uf = %#v, want RJ; body=%s", decoded["destino_uf"], dtoJSON)
	}

	rastreio, ok := decoded["rastreio"].(map[string]any)
	if !ok {
		t.Fatalf("rastreio missing/not object; body=%s", dtoJSON)
	}
	if rastreio["transportadora"] != "Total Express" {
		t.Fatalf("rastreio.transportadora = %#v, want Total Express; body=%s", rastreio["transportadora"], dtoJSON)
	}
	if rastreio["url_rastreio"] != tracking {
		t.Fatalf("rastreio.url_rastreio = %#v, want %q; body=%s", rastreio["url_rastreio"], tracking, dtoJSON)
	}

	frete, ok := decoded["frete_real"].(map[string]any)
	if !ok {
		t.Fatalf("frete_real missing/not object; body=%s", dtoJSON)
	}
	if frete["bruto"] != 23.45 {
		t.Fatalf("frete_real.bruto = %#v, want 23.45; body=%s", frete["bruto"], dtoJSON)
	}
	if frete["receiver"] != 10.0 {
		t.Fatalf("frete_real.receiver = %#v, want 10.0; body=%s", frete["receiver"], dtoJSON)
	}
	if frete["sender"] != 13.45 {
		t.Fatalf("frete_real.sender = %#v, want 13.45; body=%s", frete["sender"], dtoJSON)
	}
}

// TestMapEnrichedOrderOmitsShipmentFactsWhenAbsent asserts the honest-absence
// contract: a shipment with no carrier/costs/destino must omit frete_real and
// the carrier fields entirely (never a fabricated zero/blank — ADR-17).
func TestMapEnrichedOrderOmitsShipmentFactsWhenAbsent(t *testing.T) {
	order := enrichedOrderModels()[1]
	order.Status = "shipped"
	order.ShippingID = "ship-col-1"

	e := application.EnrichedOrder{
		Order:    order,
		Shipment: &application.ShipmentEnrichment{ShipmentID: "ship-col-1", Status: "ready_to_ship"},
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}

	for _, absent := range []string{"destino_cep", "destinatario", "destino_uf", "frete_real"} {
		if _, present := decoded[absent]; present {
			t.Fatalf("%s must be omitted when the shipment carries no such fact; body=%s", absent, dtoJSON)
		}
	}
	rastreio, ok := decoded["rastreio"].(map[string]any)
	if !ok {
		t.Fatalf("rastreio missing/not object; body=%s", dtoJSON)
	}
	if _, present := rastreio["transportadora"]; present {
		t.Fatalf("rastreio.transportadora must be omitted when carrier is nil; body=%s", dtoJSON)
	}
	if _, present := rastreio["url_rastreio"]; present {
		t.Fatalf("rastreio.url_rastreio must be omitted when tracking is nil; body=%s", dtoJSON)
	}
}

// TestMapEnrichedOrderOmitsFreteRealWhenAllAmountsNil asserts that a non-nil
// Costs whose every amount is unresolvable omits frete_real entirely rather
// than emitting an empty {} object (ADR-17 honest absence).
func TestMapEnrichedOrderOmitsFreteRealWhenAllAmountsNil(t *testing.T) {
	order := enrichedOrderModels()[1]
	order.Status = "shipped"
	order.ShippingID = "ship-col-1"

	e := application.EnrichedOrder{
		Order: order,
		Shipment: &application.ShipmentEnrichment{
			ShipmentID: "ship-col-1",
			Status:     "shipped",
			Costs:      &connectorsdomain.ShipmentCosts{}, // present, but every Money nil
		},
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}
	if _, present := decoded["frete_real"]; present {
		t.Fatalf("frete_real must be omitted when all amounts are nil; body=%s", dtoJSON)
	}
}

// TestMapEnrichedOrderSurfacesCompradorFiscal asserts the buyer fiscal identity
// (name, opaque document type + number, billing address) reaches the JSON DTO
// as an additive comprador_fiscal block when the enrichment carries it. doc_tipo
// is passed through verbatim (a non-CPF value proves it is never mapped to an
// enum — ADR-17).
func TestMapEnrichedOrderSurfacesCompradorFiscal(t *testing.T) {
	order := enrichedOrderModels()[1]

	name := "ACME Comercio LTDA"
	docType := "CNPJ_MEI" // deliberately non-CPF: proves opaque passthrough
	docNumber := "11222333000181"
	street := "Av. Paulista"
	number := "1000"
	city := "São Paulo"
	stateCode := "SP"
	stateName := "São Paulo"
	zip := "01310-100"
	country := "BR"

	e := application.EnrichedOrder{
		Order: order,
		BuyerFiscal: &connectorsdomain.BuyerFiscalInfo{
			Name:      &name,
			DocType:   &docType,
			DocNumber: &docNumber,
			Address: &connectorsdomain.BuyerFiscalAddress{
				StreetName:   &street,
				StreetNumber: &number,
				City:         &city,
				StateCode:    &stateCode,
				StateName:    &stateName,
				ZipCode:      &zip,
				CountryID:    &country,
			},
		},
	}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}

	cf, ok := decoded["comprador_fiscal"].(map[string]any)
	if !ok {
		t.Fatalf("comprador_fiscal missing/not object; body=%s", dtoJSON)
	}
	if cf["nome"] != name {
		t.Fatalf("comprador_fiscal.nome = %#v, want %q; body=%s", cf["nome"], name, dtoJSON)
	}
	if cf["doc_tipo"] != docType {
		t.Fatalf("comprador_fiscal.doc_tipo = %#v, want %q (opaque passthrough); body=%s", cf["doc_tipo"], docType, dtoJSON)
	}
	if cf["doc_numero"] != docNumber {
		t.Fatalf("comprador_fiscal.doc_numero = %#v, want %q; body=%s", cf["doc_numero"], docNumber, dtoJSON)
	}
	end, ok := cf["endereco"].(map[string]any)
	if !ok {
		t.Fatalf("comprador_fiscal.endereco missing/not object; body=%s", dtoJSON)
	}
	for k, want := range map[string]string{
		"logradouro": street, "numero": number, "cidade": city,
		"uf_codigo": stateCode, "uf_nome": stateName, "cep": zip, "pais": country,
	} {
		if end[k] != want {
			t.Fatalf("comprador_fiscal.endereco.%s = %#v, want %q; body=%s", k, end[k], want, dtoJSON)
		}
	}
}

// TestMapEnrichedOrderOmitsCompradorFiscalWhenAbsent asserts the honest-absence
// contract: a nil BuyerFiscal enrichment omits comprador_fiscal entirely (never
// a fabricated empty block — ADR-17).
func TestMapEnrichedOrderOmitsCompradorFiscalWhenAbsent(t *testing.T) {
	order := enrichedOrderModels()[1]
	e := application.EnrichedOrder{Order: order}

	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}
	if _, present := decoded["comprador_fiscal"]; present {
		t.Fatalf("comprador_fiscal must be omitted when BuyerFiscal is nil; body=%s", dtoJSON)
	}
}

// TestMapEnrichedOrderCompradorFiscalOmitsBlankAddressFields asserts that within
// a present comprador_fiscal, individual absent fields (nil pointers) are omitted
// rather than emitted as empty strings, and an all-nil address omits endereco.
func TestMapEnrichedOrderCompradorFiscalOmitsBlankAddressFields(t *testing.T) {
	order := enrichedOrderModels()[1]
	docType := "CPF"
	docNumber := "12345678909"
	e := application.EnrichedOrder{
		Order: order,
		BuyerFiscal: &connectorsdomain.BuyerFiscalInfo{
			DocType:   &docType,
			DocNumber: &docNumber,
			// Name nil, Address nil.
		},
	}
	dtoJSON, err := json.Marshal(mapEnrichedOrder(e))
	if err != nil {
		t.Fatalf("marshal dto: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(dtoJSON, &decoded); err != nil {
		t.Fatalf("decode dto: %v", err)
	}
	cf, ok := decoded["comprador_fiscal"].(map[string]any)
	if !ok {
		t.Fatalf("comprador_fiscal missing/not object; body=%s", dtoJSON)
	}
	if _, present := cf["nome"]; present {
		t.Fatalf("comprador_fiscal.nome must be omitted when nil; body=%s", dtoJSON)
	}
	if _, present := cf["endereco"]; present {
		t.Fatalf("comprador_fiscal.endereco must be omitted when address is nil; body=%s", dtoJSON)
	}
	if cf["doc_tipo"] != docType {
		t.Fatalf("comprador_fiscal.doc_tipo = %#v, want %q; body=%s", cf["doc_tipo"], docType, dtoJSON)
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

// fakeOrderFaturador is the transport-level fake for OrderFaturador (goal B),
// mirroring the fake*Store idiom already used in this package for summary.
type fakeOrderFaturador struct {
	err             error
	called          bool
	installationID  string
	providerOrderID string
}

func (f *fakeOrderFaturador) MarkFaturado(_ context.Context, installationID, providerOrderID string) error {
	f.called = true
	f.installationID = installationID
	f.providerOrderID = providerOrderID
	return f.err
}

func TestHandleFaturadoMarksOrderAndReturns204(t *testing.T) {
	fake := &fakeOrderFaturador{}
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{}).WithFaturador(fake)

	body := bytes.NewBufferString(`{"installation_id":" inst-1 "}`)
	req := httptest.NewRequest(http.MethodPost, "/orders/order-1/faturado", body)
	req.SetPathValue("provider_order_id", " order-1 ")
	rr := httptest.NewRecorder()

	handler.handleFaturado(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s, want 204", rr.Code, rr.Body.String())
	}
	if !fake.called {
		t.Fatal("handleFaturado did not call the faturador")
	}
	if fake.installationID != " inst-1 " {
		t.Fatalf("installation ID passed through = %q", fake.installationID)
	}
	if fake.providerOrderID != "order-1" {
		t.Fatalf("provider order ID = %q, want order-1 (trimmed from path)", fake.providerOrderID)
	}
}

func TestHandleFaturadoWrongMethodReturns405(t *testing.T) {
	fake := &fakeOrderFaturador{}
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{}).WithFaturador(fake)

	req := httptest.NewRequest(http.MethodGet, "/orders/order-1/faturado", nil)
	req.SetPathValue("provider_order_id", "order-1")
	rr := httptest.NewRecorder()

	handler.handleFaturado(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d body=%s, want 405", rr.Code, rr.Body.String())
	}
	if fake.called {
		t.Fatal("handleFaturado called the faturador on a GET")
	}
	if allow := rr.Header().Get("Allow"); allow != http.MethodPost {
		t.Fatalf("Allow header = %q, want POST", allow)
	}
}

func TestHandleFaturadoNotFoundReturns404(t *testing.T) {
	fake := &fakeOrderFaturador{err: &ports.OrderNotFoundError{InstallationID: "inst-1", ProviderOrderID: "missing"}}
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{}).WithFaturador(fake)

	body := bytes.NewBufferString(`{"installation_id":"inst-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/orders/missing/faturado", body)
	req.SetPathValue("provider_order_id", "missing")
	rr := httptest.NewRecorder()

	handler.handleFaturado(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s, want 404", rr.Code, rr.Body.String())
	}
}

func TestHandleFaturadoMalformedBodyReturns400(t *testing.T) {
	fake := &fakeOrderFaturador{}
	handler := NewHandlerWithReader(stubOrderImporter{}, stubOrderReadService{}).WithFaturador(fake)

	req := httptest.NewRequest(http.MethodPost, "/orders/order-1/faturado", bytes.NewBufferString(`{not json`))
	req.SetPathValue("provider_order_id", "order-1")
	rr := httptest.NewRecorder()

	handler.handleFaturado(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s, want 400", rr.Code, rr.Body.String())
	}
	if fake.called {
		t.Fatal("handleFaturado called the faturador with a malformed body")
	}
}

func TestRegisterOnlyWiresFaturadoRouteWhenBothReaderAndFaturadorSet(t *testing.T) {
	// No readLister -> no faturado route, even with a faturador set.
	mux := http.NewServeMux()
	NewHandler(stubOrderImporter{}, stubOrderLister{}).WithFaturador(&fakeOrderFaturador{}).Register(mux)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/orders/order-1/faturado", bytes.NewBufferString(`{}`)))
	if rr.Code != http.StatusNotFound && rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d body=%s, want 404 or 405 (no readLister => no faturado route)", rr.Code, rr.Body.String())
	}
}
