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
			// Real new-format GET /shipments/{id} shape (ML docs gerenciamento-de-envios):
			// the buyer delivery address lives under `destination.shipping_address`
			// ({state,city as {id,name}}, zip_code, …) + `destination.receiver_name` —
			// NOT a flat `receiver_address` field (that key does not exist in the schema).
			_, _ = io.WriteString(w, `{"id":"SHIP-1","site_id":"MLB","status":"ready_to_ship","substatus":"ready_to_print","lead_time":{"estimated_delivery_limit":{"date":"2026-07-22T15:04:05Z"}},"delayed":true,"destination":{"receiver_id":74425755,"receiver_name":"João Silva","shipping_address":{"city":{"id":"BR-RJ-RIO","name":"Rio de Janeiro"},"zip_code":"20040-002","street_name":"Avenida Rio Branco","street_number":"1","state":{"id":"BR-RJ","name":"Rio de Janeiro"}}}}`)
		case "/shipments/SHIP-1/costs":
			// Real new-format /costs shape (ML docs gerenciamento-de-envios): NO
			// currency field — the amount currency is resolved from the shipment site.
			_, _ = io.WriteString(w, `{"gross_amount":19.90,"receiver":{"user_id":74425755,"cost":2.50},"senders":[{"user_id":81387353,"cost":17.40}]}`)
		case "/shipments/SHIP-1/carrier":
			// Dedicated carrier sub-resource (ML docs gerenciamento-de-envios): {url, name}.
			_, _ = io.WriteString(w, `{"url":"http://tracking.totalexpress.com.br/poupup_track.php?reid=3&pedido=14&nfiscal=1","name":"Total Express"}`)
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
	if shipment.DestinationCity == nil || *shipment.DestinationCity != "Rio de Janeiro" {
		t.Fatalf("DestinationCity = %#v, want Rio de Janeiro", shipment.DestinationCity)
	}
	if shipment.DestinationZip == nil || *shipment.DestinationZip != "20040-002" {
		t.Fatalf("DestinationZip = %#v, want 20040-002", shipment.DestinationZip)
	}
	if shipment.ReceiverName == nil || *shipment.ReceiverName != "João Silva" {
		t.Fatalf("ReceiverName = %#v, want João Silva", shipment.ReceiverName)
	}
	if shipment.CarrierName == nil || *shipment.CarrierName != "Total Express" {
		t.Fatalf("CarrierName = %#v, want Total Express", shipment.CarrierName)
	}
	if shipment.TrackingURL == nil || *shipment.TrackingURL != "http://tracking.totalexpress.com.br/poupup_track.php?reid=3&pedido=14&nfiscal=1" {
		t.Fatalf("TrackingURL = %#v", shipment.TrackingURL)
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
		case "/shipments/SHIP-2/carrier":
			w.WriteHeader(http.StatusNotFound)
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
		case "/shipments/SHIP-5/carrier":
			w.WriteHeader(http.StatusNotFound)
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

	var costsHeader, primaryHeader, carrierHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-7":
			primaryHeader = r.Header.Get("x-format-new")
			_, _ = io.WriteString(w, `{"id":"SHIP-7","site_id":"MLB","status":"shipped","substatus":"out_for_delivery"}`)
		case "/shipments/SHIP-7/carrier":
			carrierHeader = r.Header.Get("x-format-new")
			_, _ = io.WriteString(w, `{"url":"http://track.example/OP1","name":"Correios"}`)
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
	if carrierHeader != "true" {
		t.Fatalf("carrier x-format-new = %q, want true", carrierHeader)
	}
	if shipment.Substatus != "out_for_delivery" {
		t.Fatalf("Substatus = %q, want out_for_delivery", shipment.Substatus)
	}
	if shipment.CarrierName == nil || *shipment.CarrierName != "Correios" {
		t.Fatalf("CarrierName = %#v, want Correios", shipment.CarrierName)
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
			_, _ = io.WriteString(w, `{"id":28264263908,"mode":"me2","order_id":2339711980,"order_cost":99.9,"base_cost":22.07,"site_id":"MLB","status":"shipped","substatus":"out_for_delivery","date_created":"2024-04-23T10:48:51.245-04:00","destination":{"shipping_address":{"state":{"id":"BR-SP","name":"São Paulo"}}}}`)
		case "/shipments/28264263908/carrier":
			w.WriteHeader(http.StatusNotFound)
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
		case "/shipments/SHIP-6/carrier":
			w.WriteHeader(http.StatusNotFound)
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
		case "/shipments/SHIP-8/carrier":
			w.WriteHeader(http.StatusNotFound)
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

func TestGetShipmentInfoMaskedDestinationDegradesToNil(t *testing.T) {
	t.Parallel()

	// ML obfuscates the buyer address in `destination` until payment is confirmed
	// (ML docs gerenciamento-de-envios). The masked payload keeps the object shape
	// but blanks the values. Masked-to-empty must degrade to nil honest-absence —
	// never a fabricated blank, never a WARN/error (ADR-17).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-9":
			_, _ = io.WriteString(w, `{"id":"SHIP-9","status":"pending","destination":{"receiver_name":"","shipping_address":{"state":{"id":"","name":""},"city":{"id":"","name":""},"zip_code":""}}}`)
		case "/shipments/SHIP-9/carrier":
			w.WriteHeader(http.StatusNotFound)
		case "/shipments/SHIP-9/costs":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-9")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want nil (masked address is not an error)", err)
	}
	if shipment.Status != "pending" {
		t.Fatalf("Status = %q, want pending (must survive masked destination)", shipment.Status)
	}
	if shipment.DestinationUF != nil || shipment.DestinationCity != nil || shipment.DestinationZip != nil || shipment.ReceiverName != nil {
		t.Fatalf("masked destination fields = uf:%#v city:%#v zip:%#v receiver:%#v, want all nil", shipment.DestinationUF, shipment.DestinationCity, shipment.DestinationZip, shipment.ReceiverName)
	}
}

func TestGetShipmentInfoCarrierNotFoundLeavesCarrierNil(t *testing.T) {
	t.Parallel()

	// The /carrier sub-resource returns 404 as a normal "none" condition (ML docs
	// precedent: /delays 404 = no delay). A carrier 404 must degrade to a nil
	// carrier with the shipment status intact — never a WARN, never a sink.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-10":
			_, _ = io.WriteString(w, `{"id":"SHIP-10","status":"shipped","substatus":"out_for_delivery"}`)
		case "/shipments/SHIP-10/carrier":
			w.WriteHeader(http.StatusNotFound)
		case "/shipments/SHIP-10/costs":
			_, _ = io.WriteString(w, `{"gross_amount":8.00}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-10")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want nil (carrier 404 is honest absence)", err)
	}
	if shipment.Status != "shipped" {
		t.Fatalf("Status = %q, want shipped (must survive carrier 404)", shipment.Status)
	}
	if shipment.CarrierName != nil || shipment.TrackingURL != nil {
		t.Fatalf("carrier = name:%#v url:%#v, want both nil", shipment.CarrierName, shipment.TrackingURL)
	}
	if shipment.Costs == nil {
		t.Fatal("Costs = nil, want mapped (carrier failure must not drop costs)")
	}
}

func TestGetShipmentInfoCarrierSurvivesCostsDegrade(t *testing.T) {
	t.Parallel()

	// Carrier is fetched independent of costs: a costs decode failure degrades to
	// nil Costs, but the carrier fetched before it must still be present.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-11":
			_, _ = io.WriteString(w, `{"id":"SHIP-11","status":"shipped"}`)
		case "/shipments/SHIP-11/carrier":
			_, _ = io.WriteString(w, `{"url":"http://track.example/OP99","name":"Jadlog"}`)
		case "/shipments/SHIP-11/costs":
			// Legacy/undecodable cost shape → costs degrade.
			_, _ = io.WriteString(w, `{"senders":{"cost":5.0}}`)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	shipment, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentInfo(context.Background(), pricingAccountRef(), "SHIP-11")
	if err != nil {
		t.Fatalf("GetShipmentInfo() error = %v, want degrade", err)
	}
	if shipment.Costs != nil {
		t.Fatalf("Costs = %#v, want nil (undecodable costs)", shipment.Costs)
	}
	if shipment.CarrierName == nil || *shipment.CarrierName != "Jadlog" {
		t.Fatalf("CarrierName = %#v, want Jadlog (carrier must survive costs degrade)", shipment.CarrierName)
	}
}

func TestGetShipmentInfoAbsentOptionalsStayNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-3":
			_, _ = io.WriteString(w, `{"id":"SHIP-3","site_id":"MLB","status":"pending","delayed":null}`)
		case "/shipments/SHIP-3/carrier":
			w.WriteHeader(http.StatusNotFound)
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
			// state.id already a bare UF (no "BR-" prefix) — trim must be a no-op.
			_, _ = io.WriteString(w, `{"id":"SHIP-4","destination":{"shipping_address":{"state":{"id":"SP","name":"São Paulo"}}}}`)
		case "/shipments/SHIP-4/carrier":
			w.WriteHeader(http.StatusNotFound)
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
