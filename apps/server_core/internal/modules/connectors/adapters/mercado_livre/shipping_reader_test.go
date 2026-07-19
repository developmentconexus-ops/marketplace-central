package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func TestGetShipmentInfoMapsShipmentAndCosts(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		switch r.URL.Path {
		case "/shipments/SHIP-1":
			_, _ = io.WriteString(w, `{"id":"SHIP-1","site_id":"MLB","status":"ready_to_ship","substatus":"ready_to_print","lead_time":{"estimated_delivery_limit":{"date":"2026-07-22T15:04:05Z"}},"delayed":true,"receiver_address":{"state":{"id":"BR-RJ"}}}`)
		case "/shipments/SHIP-1/costs":
			// Real new-format /costs shape (ML docs gerenciamento-de-envios): NO
			// currency field — the amount currency is resolved from the shipment site.
			_, _ = io.WriteString(w, `{"gross_amount":19.90,"receiver":{"user_id":74425755,"cost":2.50},"senders":[{"user_id":81387353,"cost":17.40}]}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, fetchedAt).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-1")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v", err)
	}
	if shipment.ID != "SHIP-1" || shipment.Status != "ready_to_ship" || !shipment.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("shipment = %#v", shipment)
	}
	if shipment.Substatus != "ready_to_print" {
		t.Fatalf("Substatus = %q, want ready_to_print", shipment.Substatus)
	}
	wantSLA := time.Date(2026, 7, 22, 15, 4, 5, 0, time.UTC)
	if shipment.SLADue == nil || !shipment.SLADue.Equal(wantSLA) {
		t.Fatalf("SLADue = %#v, want %v", shipment.SLADue, wantSLA)
	}
	if shipment.Delayed == nil || !*shipment.Delayed {
		t.Fatalf("Delayed = %#v, want true", shipment.Delayed)
	}
	if shipment.DestinationUF == nil || *shipment.DestinationUF != "RJ" {
		t.Fatalf("DestinationUF = %#v, want RJ", shipment.DestinationUF)
	}
	if shipment.Costs == nil {
		t.Fatal("Costs = nil, want mapped costs")
	}
	if got := shipment.Costs.GrossAmount; got == nil || *got != (domain.Money{Amount: "19.90", Currency: "BRL"}) {
		t.Fatalf("GrossAmount = %#v", got)
	}
	if got := shipment.Costs.ReceiverCost; got == nil || *got != (domain.Money{Amount: "2.50", Currency: "BRL"}) {
		t.Fatalf("ReceiverCost = %#v", got)
	}
	if got := shipment.Costs.SenderCost; got == nil || *got != (domain.Money{Amount: "17.40", Currency: "BRL"}) {
		t.Fatalf("SenderCost = %#v", got)
	}
}

func TestGetShipmentInfoCostsNotFoundLeavesCostsUnknown(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-2":
			_, _ = io.WriteString(w, `{"id":"SHIP-2","status":"shipped"}`)
		case "/shipments/SHIP-2/costs":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-2")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v", err)
	}
	if shipment.Costs != nil {
		t.Fatalf("Costs = %#v, want nil", shipment.Costs)
	}
}

func TestGetShipmentInfoCostsServerErrorFails(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-5":
			_, _ = io.WriteString(w, `{"id":"SHIP-5","status":"shipped"}`)
		case "/shipments/SHIP-5/costs":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-5")
	if !errors.Is(err, domain.ErrProviderUnavailable) {
		t.Fatalf("GetShipmentInfo error = %v, want ErrProviderUnavailable", err)
	}
	if shipment != (domain.ShipmentInfo{}) {
		t.Fatalf("shipment = %#v, want zero value on costs failure", shipment)
	}
}

func TestGetShipmentInfoRequestsSendXFormatNewHeader(t *testing.T) {
	t.Parallel()

	var costsHeader, primaryHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-7":
			primaryHeader = r.Header.Get("x-format-new")
			_, _ = io.WriteString(w, `{"id":"SHIP-7","site_id":"MLB","status":"shipped","substatus":"out_for_delivery"}`)
		case "/shipments/SHIP-7/costs":
			costsHeader = r.Header.Get("x-format-new")
			_, _ = io.WriteString(w, `{"gross_amount":10.00}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-7")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v", err)
	}
	// BOTH ML shipment GETs require x-format-new:true. Without it the primary
	// GET /shipments/{id} returns the legacy shape and the whole read sinks
	// before the costs degrade is reached (ML docs: gerenciamento-de-envios,
	// "é necessário enviar o header 'x-format-new: true'"); the /costs call
	// likewise needs it or costs sink.
	if primaryHeader != "true" {
		t.Fatalf("primary x-format-new = %q, want true", primaryHeader)
	}
	if costsHeader != "true" {
		t.Fatalf("costs x-format-new = %q, want true", costsHeader)
	}
	if shipment.Substatus != "out_for_delivery" {
		t.Fatalf("Substatus = %q, want out_for_delivery", shipment.Substatus)
	}
}

func TestGetShipmentInfoDecodesNewFormatNumericID(t *testing.T) {
	t.Parallel()

	// The x-format-new shipment JSON carries a NUMERIC id (ML docs:
	// status-de-pedidos-rastreamento — "id": 28264263908) plus many extra
	// top-level fields. The read must decode it (numeric id, unknown fields
	// ignored) and surface status+substatus — not sink on the id type.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/28264263908":
			_, _ = io.WriteString(w, `{"id":28264263908,"mode":"me2","order_id":2339711980,"order_cost":99.9,"base_cost":22.07,"site_id":"MLB","status":"shipped","substatus":"out_for_delivery","date_created":"2024-04-23T10:48:51.245-04:00","receiver_address":{"state":{"id":"BR-SP"}}}`)
		case "/shipments/28264263908/costs":
			_, _ = io.WriteString(w, `{"gross_amount":19.90}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "28264263908")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want nil (numeric id must decode)", err)
	}
	if shipment.ID != "28264263908" {
		t.Fatalf("ID = %q, want 28264263908", shipment.ID)
	}
	if shipment.Status != "shipped" || shipment.Substatus != "out_for_delivery" {
		t.Fatalf("status/substatus = %q/%q, want shipped/out_for_delivery", shipment.Status, shipment.Substatus)
	}
	if shipment.DestinationUF == nil || *shipment.DestinationUF != "SP" {
		t.Fatalf("DestinationUF = %#v, want SP", shipment.DestinationUF)
	}
}

func TestGetShipmentInfoCostsLegacyShapeDegradesToStatusOnly(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-6":
			_, _ = io.WriteString(w, `{"id":"SHIP-6","status":"shipped","substatus":"receiver_absent"}`)
		case "/shipments/SHIP-6/costs":
			// Legacy cost shape (what ML returns when x-format-new is absent): `senders` is an
			// object, not the new-format array → decode into mlShipmentCostsResponse fails with
			// ProviderPayloadInvalid. That must degrade to shipment-with-status, not sink the read.
			_, _ = io.WriteString(w, `{"senders":{"cost":17.40},"currency_id":"BRL"}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-6")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want degrade (nil error)", err)
	}
	if shipment.Status != "shipped" {
		t.Fatalf("Status = %q, want shipped (must survive costs decode failure)", shipment.Status)
	}
	if shipment.Substatus != "receiver_absent" {
		t.Fatalf("Substatus = %q, want receiver_absent (must survive costs decode failure)", shipment.Substatus)
	}
	if shipment.Costs != nil {
		t.Fatalf("Costs = %#v, want nil (undecodable legacy body → costs unknown)", shipment.Costs)
	}
}

func TestGetShipmentInfoCostsMappingFailureDegradesToStatusOnly(t *testing.T) {
	t.Parallel()

	// The /costs body decodes cleanly into mlShipmentCostsResponse but a value
	// fails money mapping (exponent amount is not a valid decimal string per
	// domain.ValidateMoney). A mapping failure must degrade to the costs-less
	// shipment — status/substatus survive, Costs nil — NEVER sink the whole read
	// (round-1 invariant: costs failure degrades, status survives).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-8":
			_, _ = io.WriteString(w, `{"id":"SHIP-8","site_id":"MLB","status":"shipped","substatus":"out_for_delivery"}`)
		case "/shipments/SHIP-8/costs":
			_, _ = io.WriteString(w, `{"gross_amount":2.4e1}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-8")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want degrade (nil error)", err)
	}
	if shipment.Status != "shipped" || shipment.Substatus != "out_for_delivery" {
		t.Fatalf("status/substatus = %q/%q, want shipped/out_for_delivery (must survive costs mapping failure)", shipment.Status, shipment.Substatus)
	}
	if shipment.Costs != nil {
		t.Fatalf("Costs = %#v, want nil (unmappable amount → costs unknown)", shipment.Costs)
	}
}

func TestGetShipmentInfoAbsentOptionalsStayNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-3":
			_, _ = io.WriteString(w, `{"id":"SHIP-3","site_id":"MLB","status":"pending","delayed":null}`)
		case "/shipments/SHIP-3/costs":
			_, _ = io.WriteString(w, `{"gross_amount":12.00}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-3")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v", err)
	}
	if shipment.Delayed != nil || shipment.SLADue != nil || shipment.DestinationUF != nil {
		t.Fatalf("optional shipment fields = delayed:%#v sla:%#v uf:%#v", shipment.Delayed, shipment.SLADue, shipment.DestinationUF)
	}
	if shipment.Costs == nil || shipment.Costs.ReceiverCost != nil || shipment.Costs.SenderCost != nil {
		t.Fatalf("optional cost fields = %#v", shipment.Costs)
	}
	if got := shipment.Costs.GrossAmount; got == nil || *got != (domain.Money{Amount: "12.00", Currency: "BRL"}) {
		t.Fatalf("GrossAmount = %#v, want 12.00 BRL", got)
	}
}

func TestGetShipmentInfoDestinationUFPreservesBareUF(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-4":
			_, _ = io.WriteString(w, `{"id":"SHIP-4","receiver_address":{"state":{"id":"SP"}}}`)
		case "/shipments/SHIP-4/costs":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-4")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v", err)
	}
	if shipment.DestinationUF == nil || *shipment.DestinationUF != "SP" {
		t.Fatalf("DestinationUF = %#v, want SP", shipment.DestinationUF)
	}
}

func TestGetShipmentInfoMapsPrimaryProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{{http.StatusUnauthorized, domain.ErrUnauthorized}, {http.StatusNotFound, domain.ErrNotFound}, {http.StatusTooManyRequests, domain.ErrRateLimited}, {http.StatusInternalServerError, domain.ErrProviderUnavailable}} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()

			_, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-ERR")
			if !errors.Is(err, test.want) {
				t.Fatalf("GetShipmentInfo error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestGetFreeShippingCostMapsCostAndRequest(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 18, 13, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/users/seller-1/shipping_options/free" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("item_id"); got != "MLB-ITEM-1" {
			t.Fatalf("item_id = %q", got)
		}
		_, _ = io.WriteString(w, `{"coverage":{"all_country":{"list_cost":25.75}},"currency_id":"BRL"}`)
	}))
	defer server.Close()

	result, err := pricingTestAdapter(server.URL, fetchedAt).GetFreeShippingCost(context.Background(), pricingAccountRef(), domain.FreeShippingQuery{ItemID: "MLB-ITEM-1"})
	if err != nil {
		t.Fatalf("GetFreeShippingCost() error = %v", err)
	}
	if !result.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", result.FetchedAt, fetchedAt)
	}
	if result.Cost == nil || *result.Cost != (domain.Money{Amount: "25.75", Currency: "BRL"}) {
		t.Fatalf("Cost = %#v", result.Cost)
	}
}

func TestGetFreeShippingCostAbsentCostStaysNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"coverage":{},"currency_id":"BRL"}`)
	}))
	defer server.Close()

	result, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetFreeShippingCost(context.Background(), pricingAccountRef(), domain.FreeShippingQuery{ItemID: "MLB-ITEM-2"})
	if err != nil {
		t.Fatalf("GetFreeShippingCost() error = %v", err)
	}
	if result.Cost != nil {
		t.Fatalf("Cost = %#v, want nil", result.Cost)
	}
}

func TestGetFreeShippingCostMapsProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{{http.StatusUnauthorized, domain.ErrUnauthorized}, {http.StatusNotFound, domain.ErrNotFound}, {http.StatusTooManyRequests, domain.ErrRateLimited}, {http.StatusInternalServerError, domain.ErrProviderUnavailable}} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()

			_, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetFreeShippingCost(context.Background(), pricingAccountRef(), domain.FreeShippingQuery{ItemID: "MLB-ITEM-ERR"})
			if !errors.Is(err, test.want) {
				t.Fatalf("GetFreeShippingCost error = %v, want %v", err, test.want)
			}
		})
	}
}
