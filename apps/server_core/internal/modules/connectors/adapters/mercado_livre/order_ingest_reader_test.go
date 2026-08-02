package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// scrubbedOrderFixture is a realistic /orders/{id} payload shaped after the documented ML
// order object (ML docs: gestion-de-pedidos) — status/status_detail/tags/date_last_updated/
// pack_id/shipping.id/buyer/payments/order_items with sale_fee PER UNIT (fact T2). PII in
// buyer.billing_info is masked-but-present (never a real fixture — synthetic per PII rules).
const scrubbedOrderFixture = `{
  "id": 2000003508372839,
  "status": "paid",
  "status_detail": "not_specified",
  "date_created": "2026-07-18T10:00:00.000-04:00",
  "date_closed": "2026-07-18T10:05:00.000-04:00",
  "last_updated": "2026-07-18T10:05:01.000-04:00",
  "date_last_updated": "2026-07-18T10:05:01.000-04:00",
  "pack_id": 2000003508372800,
  "tags": ["paid", "not_delivered"],
  "shipping": {"id": 44556677},
  "buyer": {
    "id": 555444333,
    "nickname": "COMPRADOR_TESTE",
    "billing_info": {"id": "999-billing", "buyer_ssn_or_pii_marker": "SHOULD_NEVER_SURVIVE"}
  },
  "order_items": [
    {
      "item": {
        "id": "MLB123456",
        "variation_id": 987,
        "seller_sku": "SKU-A",
        "title": "Produto A",
        "variation_attributes": [{"id": "GTIN", "value_id": "", "value_name": "7891234567890"}]
      },
      "quantity": 3,
      "unit_price": 49.90,
      "sale_fee": 6.25
    },
    {
      "item": {
        "id": "MLB654321",
        "variation_id": null,
        "seller_sku": "SKU-B",
        "title": "Produto B",
        "variation_attributes": []
      },
      "quantity": 1,
      "unit_price": 120.00,
      "sale_fee": 15.60
    }
  ],
  "payments": [
    {"id": 778899, "status": "approved", "total_paid_amount": 269.70, "transaction_amount": 269.70}
  ]
}`

func TestGetOrderDetailDecodesAllTargetFields(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/orders/2000003508372839" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, scrubbedOrderFixture)
	}))
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, fetchedAt).GetOrderDetail(context.Background(), pricingAccountRef(), "2000003508372839")
	if err != nil {
		t.Fatalf("GetOrderDetail() error = %v", err)
	}

	if detail.ProviderCode != "mercado_livre" {
		t.Fatalf("ProviderCode = %q, want mercado_livre", detail.ProviderCode)
	}
	if detail.ProviderOrderID != "2000003508372839" {
		t.Fatalf("ProviderOrderID = %q", detail.ProviderOrderID)
	}
	if detail.ProviderStatus != "paid" || detail.ProviderStatusDetail != "not_specified" {
		t.Fatalf("status/detail = %q/%q", detail.ProviderStatus, detail.ProviderStatusDetail)
	}
	// This fixture's order is "paid", not cancelled, so the payload carries no
	// cancel_detail block. CancellationDetail must stay empty here — mirroring
	// status_detail ("not_specified") was the P2 bug this DTO change fixes (rawkeys:
	// cancel_detail is the real source; status_detail is an unrelated field).
	if detail.CancellationDetail != "" {
		t.Fatalf("CancellationDetail = %q, want empty (no cancel_detail in a non-cancelled order)", detail.CancellationDetail)
	}
	wantCreated := time.Date(2026, 7, 18, 14, 0, 0, 0, time.UTC)
	if detail.ProviderCreatedAt == nil || !detail.ProviderCreatedAt.Equal(wantCreated) {
		t.Fatalf("ProviderCreatedAt = %#v, want %v", detail.ProviderCreatedAt, wantCreated)
	}
	if detail.ProviderClosedAt == nil {
		t.Fatal("ProviderClosedAt = nil, want set")
	}
	if detail.ProviderUpdatedAt == nil {
		t.Fatal("ProviderUpdatedAt = nil, want set")
	}
	if detail.DateLastUpdatedML == nil || !detail.DateLastUpdatedML.Equal(*detail.ProviderUpdatedAt) {
		t.Fatalf("DateLastUpdatedML = %#v, want equal to ProviderUpdatedAt (date_last_updated verbatim) %#v", detail.DateLastUpdatedML, detail.ProviderUpdatedAt)
	}
	if detail.PackID == nil || *detail.PackID != "2000003508372800" {
		t.Fatalf("PackID = %#v, want 2000003508372800", detail.PackID)
	}
	if detail.ShippingID == nil || *detail.ShippingID != "44556677" {
		t.Fatalf("ShippingID = %#v, want 44556677", detail.ShippingID)
	}
	if detail.BuyerNickname == nil || *detail.BuyerNickname != "COMPRADOR_TESTE" {
		t.Fatalf("BuyerNickname = %#v, want COMPRADOR_TESTE", detail.BuyerNickname)
	}
	if len(detail.Tags) != 2 || detail.Tags[0] != "paid" || detail.Tags[1] != "not_delivered" {
		t.Fatalf("Tags = %#v", detail.Tags)
	}
	if !detail.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", detail.FetchedAt, fetchedAt)
	}

	if len(detail.Items) != 2 {
		t.Fatalf("Items = %#v, want 2", detail.Items)
	}
	item0 := detail.Items[0]
	if item0.ProviderItemID != "MLB123456" || item0.ProviderVariationID != "987" || item0.SellerSKU != "SKU-A" {
		t.Fatalf("item0 identity = %#v", item0)
	}
	if item0.EAN != "7891234567890" {
		t.Fatalf("item0.EAN = %q, want 7891234567890", item0.EAN)
	}
	if item0.Quantity != 3 {
		t.Fatalf("item0.Quantity = %d, want 3 (R-4 must-fail needs qty>1)", item0.Quantity)
	}
	if item0.UnitPrice == nil || *item0.UnitPrice != 49.90 {
		t.Fatalf("item0.UnitPrice = %#v, want 49.90", item0.UnitPrice)
	}
	// SaleFeeUnit must be the PROVIDER'S VERBATIM per-line-item sale_fee — 6.25, NOT
	// multiplied by quantity (3) into a total of 18.75. Naming (SaleFeeUnit, never
	// SaleFeeTotal) plus this value assertion together prove the per-unit contract.
	if item0.SaleFeeUnit == nil || *item0.SaleFeeUnit != 6.25 {
		t.Fatalf("item0.SaleFeeUnit = %#v, want 6.25 (per unit, not qty*fee)", item0.SaleFeeUnit)
	}

	item1 := detail.Items[1]
	if item1.ProviderVariationID != "" {
		t.Fatalf("item1.ProviderVariationID = %q, want empty (null variation_id)", item1.ProviderVariationID)
	}
	if item1.EAN != "" {
		t.Fatalf("item1.EAN = %q, want empty (no GTIN attribute)", item1.EAN)
	}

	if len(detail.Payments) != 1 {
		t.Fatalf("Payments = %#v, want 1", detail.Payments)
	}
	payment := detail.Payments[0]
	if payment.PaymentID != "778899" || payment.Status != "approved" {
		t.Fatalf("payment = %#v", payment)
	}
	if payment.TotalPaidAmount == nil || *payment.TotalPaidAmount != 269.70 {
		t.Fatalf("payment.TotalPaidAmount = %#v", payment.TotalPaidAmount)
	}

	if detail.RawOrder == nil {
		t.Fatal("RawOrder = nil, want a raw capture")
	}
	if strings.Contains(string(detail.RawOrder), "billing_info") {
		t.Fatalf("RawOrder contains billing_info substring: %s", detail.RawOrder)
	}
	if strings.Contains(string(detail.RawOrder), "SHOULD_NEVER_SURVIVE") {
		t.Fatalf("RawOrder leaked the billing_info payload content: %s", detail.RawOrder)
	}
	// The raw capture must still carry the non-fiscal buyer facts (nickname) — only
	// billing_info is stripped, not the whole buyer block.
	if !strings.Contains(string(detail.RawOrder), "COMPRADOR_TESTE") {
		t.Fatalf("RawOrder dropped non-fiscal buyer data it should have kept: %s", detail.RawOrder)
	}
	// order_items must survive the redaction byte-for-byte-equivalent (no float64
	// round-trip precision loss) — spot check a decimal sale_fee value.
	if !strings.Contains(string(detail.RawOrder), "6.25") {
		t.Fatalf("RawOrder lost order_items content during redaction: %s", detail.RawOrder)
	}
}

// TestGetOrderDetailCancelledOrderUsesCancelDetail is the P2 regression: a cancelled order's
// cancel_detail (never status_detail, which the provider sent null on 7/7 measured cancelled
// orders) is the source of CancellationDetail, in canonical "<requested_by>:<code>" form.
// testdata/order_body.json is the real cancelled-order shape (PII redacted) this bug was
// measured against.
func TestGetOrderDetailCancelledOrderUsesCancelDetail(t *testing.T) {
	t.Parallel()

	fixture, err := os.ReadFile("testdata/order_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetOrderDetail(context.Background(), pricingAccountRef(), "2000012659424976")
	if err != nil {
		t.Fatalf("GetOrderDetail() error = %v", err)
	}
	if detail.CancellationDetail != "buyer:cancel_purchase_by_buyer" {
		t.Fatalf("CancellationDetail = %q, want buyer:cancel_purchase_by_buyer", detail.CancellationDetail)
	}
}

func TestGetOrderDetailAbsentShippingIDStaysNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"id":"ORD-1","status":"paid"}`)
	}))
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetOrderDetail(context.Background(), pricingAccountRef(), "ORD-1")
	if err != nil {
		t.Fatalf("GetOrderDetail() error = %v", err)
	}
	if detail.ShippingID != nil {
		t.Fatalf("ShippingID = %#v, want nil (no shipping.id in payload — caller must not call GetShipmentDetail)", detail.ShippingID)
	}
	if detail.PackID != nil {
		t.Fatalf("PackID = %#v, want nil (absent from payload)", detail.PackID)
	}
	if detail.BuyerNickname != nil {
		t.Fatalf("BuyerNickname = %#v, want nil (no buyer block)", detail.BuyerNickname)
	}
	if detail.RawOrder == nil || strings.Contains(string(detail.RawOrder), "billing_info") {
		t.Fatalf("RawOrder = %s, want present and billing_info-free", detail.RawOrder)
	}
}

func TestGetOrderDetailMapsProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{
		// 403 (third-party order, a documented live fact) must map the same as 401 —
		// unauthorized/forbidden both collapse to ErrCodeProviderAuth in the shared
		// status mapping, and MUST be a typed, non-retryable error, never zero-valued.
		{http.StatusForbidden, domain.ErrUnauthorized},
		{http.StatusUnauthorized, domain.ErrUnauthorized},
		{http.StatusNotFound, domain.ErrNotFound},
		{http.StatusTooManyRequests, domain.ErrRateLimited},
		{http.StatusInternalServerError, domain.ErrProviderUnavailable},
	} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()

			detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetOrderDetail(context.Background(), pricingAccountRef(), "ORD-ERR")
			if !errors.Is(err, test.want) {
				t.Fatalf("GetOrderDetail error = %v, want %v", err, test.want)
			}
			if detail.ProviderOrderID != "" || detail.Items != nil || detail.RawOrder != nil {
				t.Fatalf("detail = %#v, want zero value on error", detail)
			}
		})
	}
}

func TestGetOrderDetailEmptyProviderOrderIDIsInvalidReference(t *testing.T) {
	t.Parallel()

	_, err := pricingTestAdapter("http://unused.invalid", time.Now().UTC()).GetOrderDetail(context.Background(), pricingAccountRef(), "   ")
	if domain.ErrorCodeOf(err) != domain.ErrCodeProviderInvalidReference {
		t.Fatalf("error code = %v, want %v", domain.ErrorCodeOf(err), domain.ErrCodeProviderInvalidReference)
	}
}
