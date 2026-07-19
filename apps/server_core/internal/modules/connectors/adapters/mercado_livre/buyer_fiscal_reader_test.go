package mercadolivre

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// The two-step ML fiscal flow (docs: faturamento):
//  1. GET /orders/{id}          -> buyer.billing_info.id
//  2. GET /orders/billing-info/{SITE_ID}/{billing_info_id} -> name, identification, address
//
// Fixtures mirror the documented payload shapes verbatim (deepmap Q2). For MLB the
// identification.type literal is UNVERIFIED in the docs (examples are MLA DNI/CUIT), so the
// decoder treats it as an opaque string — see TestGetBuyerFiscalInfoTypeIsOpaque.
func TestGetBuyerFiscalInfoMapsTwoStepFlow(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-1":
			_, _ = w.Write([]byte(`{"id":2000010733434062,"site_id":"MLB","buyer":{"id":2212631646,"nickname":"JOAOSILVA","billing_info":{"id":"677487519924852462"}}}`))
		case "/orders/billing-info/MLB/677487519924852462":
			_, _ = w.Write([]byte(`{"site_id":"MLB","buyer":{"cust_id":2212631646,"billing_info":{"name":"João","last_name":"Silva","identification":{"type":"CPF","number":"12345678909"},"address":{"street_name":"Avenida Rio Branco","street_number":"1","city_name":"Rio de Janeiro","state":{"code":"BR-RJ","name":"Rio de Janeiro"},"zip_code":"20040-002","country_id":"BR"}}},"seller":{"cust_id":81387353}}`))
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	info, err := pricingTestAdapter(server.URL, fetchedAt).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-1")
	if err != nil {
		t.Fatalf("GetBuyerFiscalInfo() error = %v", err)
	}
	if !info.HasData() {
		t.Fatal("HasData() = false, want true")
	}
	if !info.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", info.FetchedAt, fetchedAt)
	}
	if info.Name == nil || *info.Name != "João Silva" {
		t.Fatalf("Name = %#v, want João Silva", info.Name)
	}
	if info.DocType == nil || *info.DocType != "CPF" {
		t.Fatalf("DocType = %#v, want CPF", info.DocType)
	}
	if info.DocNumber == nil || *info.DocNumber != "12345678909" {
		t.Fatalf("DocNumber = %#v, want 12345678909", info.DocNumber)
	}
	if info.Address == nil {
		t.Fatal("Address = nil, want populated")
	}
	addr := info.Address
	if addr.StreetName == nil || *addr.StreetName != "Avenida Rio Branco" {
		t.Fatalf("StreetName = %#v", addr.StreetName)
	}
	if addr.StreetNumber == nil || *addr.StreetNumber != "1" {
		t.Fatalf("StreetNumber = %#v", addr.StreetNumber)
	}
	if addr.City == nil || *addr.City != "Rio de Janeiro" {
		t.Fatalf("City = %#v", addr.City)
	}
	if addr.StateCode == nil || *addr.StateCode != "BR-RJ" {
		t.Fatalf("StateCode = %#v, want verbatim BR-RJ", addr.StateCode)
	}
	if addr.StateName == nil || *addr.StateName != "Rio de Janeiro" {
		t.Fatalf("StateName = %#v", addr.StateName)
	}
	if addr.ZipCode == nil || *addr.ZipCode != "20040-002" {
		t.Fatalf("ZipCode = %#v", addr.ZipCode)
	}
	if addr.CountryID == nil || *addr.CountryID != "BR" {
		t.Fatalf("CountryID = %#v", addr.CountryID)
	}
}

// identification.type is rendered as an opaque string: the decoder must never assume the
// MLB literal is "CPF"/"CNPJ" nor branch on it (ADR-17, dispatch constraint). An unexpected
// type string passes straight through unchanged.
func TestGetBuyerFiscalInfoTypeIsOpaque(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-2":
			_, _ = w.Write([]byte(`{"id":42,"site_id":"MLB","buyer":{"id":7,"billing_info":{"id":"BI-2"}}}`))
		case "/orders/billing-info/MLB/BI-2":
			_, _ = w.Write([]byte(`{"buyer":{"billing_info":{"name":"ACME","last_name":"","identification":{"type":"CNPJ_MEI","number":"11222333000181"}}}}`))
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	info, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-2")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if info.DocType == nil || *info.DocType != "CNPJ_MEI" {
		t.Fatalf("DocType = %#v, want opaque CNPJ_MEI passthrough", info.DocType)
	}
	// last_name empty -> Name is just the first name, no trailing space.
	if info.Name == nil || *info.Name != "ACME" {
		t.Fatalf("Name = %#v, want ACME", info.Name)
	}
}

// A buyer with no billing_info.id (order payload carries no billing block) is honest
// absence, not an error: the billing-info call is skipped and HasData() is false. The
// server fails the test if the billing-info endpoint is hit.
func TestGetBuyerFiscalInfoNoBillingInfoIDIsAbsent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-3":
			_, _ = w.Write([]byte(`{"id":99,"site_id":"MLB","buyer":{"id":7,"nickname":"NB"}}`))
		default:
			t.Fatalf("billing-info must NOT be called when billing_info.id is absent; got %s", r.URL.Path)
		}
	}))
	defer server.Close()

	info, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-3")
	if err != nil {
		t.Fatalf("error = %v, want nil (honest absence)", err)
	}
	if info.HasData() {
		t.Fatalf("HasData() = true, want false for a buyer with no billing_info")
	}
}

// A 404 on the billing-info call is honest absence (a buyer without billing data), NOT an
// error: the order still renders. Degrade to empty + nil error, never a WARN-worthy failure.
func TestGetBuyerFiscalInfoBillingInfoNotFoundIsAbsent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-4":
			_, _ = w.Write([]byte(`{"id":100,"site_id":"MLB","buyer":{"id":7,"billing_info":{"id":"BI-4"}}}`))
		case "/orders/billing-info/MLB/BI-4":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	info, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-4")
	if err != nil {
		t.Fatalf("error = %v, want nil (404 = honest absence)", err)
	}
	if info.HasData() {
		t.Fatalf("HasData() = true, want false on billing-info 404")
	}
}

// A non-degradable provider failure (500) on billing-info propagates as an error so the
// caller can degrade+warn ONCE — distinct from the silent 404 absence path.
func TestGetBuyerFiscalInfoBillingInfoServerErrorFails(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-5":
			_, _ = w.Write([]byte(`{"id":101,"site_id":"MLB","buyer":{"id":7,"billing_info":{"id":"BI-5"}}}`))
		case "/orders/billing-info/MLB/BI-5":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	_, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-5")
	if !errors.Is(err, domain.ErrProviderUnavailable) {
		t.Fatalf("error = %v, want provider unavailable", err)
	}
}

// The billing-info call is Bearer-only: no x-format-new header (that is a shipments-only
// requirement). Masked/blank provider fields degrade to nil (ADR-17), never a fabricated blank.
func TestGetBuyerFiscalInfoBillingInfoIsBearerOnlyAndMasksDegrade(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orders/ORDER-6":
			_, _ = w.Write([]byte(`{"id":102,"site_id":"MLB","buyer":{"id":7,"billing_info":{"id":"BI-6"}}}`))
		case "/orders/billing-info/MLB/BI-6":
			if got := r.Header.Get("x-format-new"); got != "" {
				t.Fatalf("x-format-new = %q, want absent on billing-info", got)
			}
			if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "Bearer ") {
				t.Fatalf("Authorization = %q, want Bearer token", auth)
			}
			// name present but address fields all blank -> Address degrades to nil.
			_, _ = w.Write([]byte(`{"buyer":{"billing_info":{"name":"Ana","last_name":"","identification":{"type":"CPF","number":"1"},"address":{"street_name":"  ","street_number":"","city_name":"","state":{"code":"","name":""},"zip_code":"","country_id":""}}}}`))
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	info, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetBuyerFiscalInfo(context.Background(), pricingAccountRef(), "ORDER-6")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if info.Address != nil {
		t.Fatalf("Address = %#v, want nil (all fields masked/blank)", info.Address)
	}
	if info.Name == nil || *info.Name != "Ana" {
		t.Fatalf("Name = %#v, want Ana", info.Name)
	}
}
