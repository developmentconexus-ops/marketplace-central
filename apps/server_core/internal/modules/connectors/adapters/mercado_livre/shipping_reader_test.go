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
			_, _ = io.WriteString(w, `{"id":"SHIP-1","status":"ready_to_ship","lead_time":{"estimated_delivery_limit":{"date":"2026-07-22T15:04:05Z"}},"delayed":true,"receiver_address":{"state":{"id":"BR-RJ"}}}`)
		case "/shipments/SHIP-1/costs":
			_, _ = io.WriteString(w, `{"gross_amount":19.90,"receiver":{"cost":2.50},"senders":[{"cost":17.40}],"currency_id":"BRL"}`)
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

func TestGetShipmentInfoAbsentOptionalsStayNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/shipments/SHIP-3":
			_, _ = io.WriteString(w, `{"id":"SHIP-3","status":"pending","delayed":null}`)
		case "/shipments/SHIP-3/costs":
			_, _ = io.WriteString(w, `{"gross_amount":12.00,"currency_id":"BRL"}`)
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
