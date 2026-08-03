package mercadolivre

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

const scrubbedShipmentFixture = `{
  "id": 44556677,
  "site_id": "MLB",
  "status": "shipped",
  "substatus": "in_hub",
  "date_created": "2026-07-18T10:01:00.000-04:00",
  "tracking_number": "BR123456789",
  "lead_time": {
    "estimated_delivery_limit": {"date": "2026-07-25T23:59:59.000-04:00"}
  },
  "destination": {
    "receiver_name": "COMPRADOR TESTE",
    "shipping_address": {
      "state": {"id": "BR-SP", "name": "Sao Paulo"},
      "city": {"id": "BR-SP-XYZ", "name": "Sao Paulo"},
      "zip_code": "01000-000"
    }
  }
}`

const scrubbedShipmentCostsFixture = `{
  "gross_amount": 24.90,
  "senders": [{"cost": 12.45}]
}`

func newShipmentTestServer(t *testing.T, primaryStatus int, primaryBody string, costsStatus int, costsBody string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-format-new") != "true" {
			t.Fatalf("missing x-format-new header on %s", r.URL.Path)
		}
		switch r.URL.Path {
		case "/shipments/44556677":
			w.WriteHeader(primaryStatus)
			if primaryBody != "" {
				_, _ = w.Write([]byte(primaryBody))
			}
		case "/shipments/44556677/costs":
			w.WriteHeader(costsStatus)
			if costsBody != "" {
				_, _ = w.Write([]byte(costsBody))
			}
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
}

func TestGetShipmentDetailDecodesAllTargetFields(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	server := newShipmentTestServer(t, http.StatusOK, scrubbedShipmentFixture, http.StatusOK, scrubbedShipmentCostsFixture)
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, fetchedAt).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if err != nil {
		t.Fatalf("GetShipmentDetail() error = %v", err)
	}
	if !detail.Found {
		t.Fatal("Found = false, want true")
	}
	if detail.ProviderShipmentID != "44556677" {
		t.Fatalf("ProviderShipmentID = %q", detail.ProviderShipmentID)
	}
	if detail.Status != "shipped" {
		t.Fatalf("Status = %q", detail.Status)
	}
	if detail.Substatus == nil || *detail.Substatus != "in_hub" {
		t.Fatalf("Substatus = %#v", detail.Substatus)
	}
	if detail.TrackingNumber == nil || *detail.TrackingNumber != "BR123456789" {
		t.Fatalf("TrackingNumber = %#v", detail.TrackingNumber)
	}
	if detail.SourceTime == nil {
		t.Fatal("SourceTime = nil, want set from date_created")
	}
	if !detail.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", detail.FetchedAt, fetchedAt)
	}
	if detail.ReceiverName == nil || *detail.ReceiverName != "COMPRADOR TESTE" {
		t.Fatalf("ReceiverName = %#v", detail.ReceiverName)
	}
	if detail.DestState == nil || *detail.DestState != "SP" {
		t.Fatalf("DestState = %#v, want SP (BR- prefix trimmed)", detail.DestState)
	}
	if detail.DestCity == nil || *detail.DestCity != "Sao Paulo" {
		t.Fatalf("DestCity = %#v", detail.DestCity)
	}
	if detail.DestZip == nil || *detail.DestZip != "01000-000" {
		t.Fatalf("DestZip = %#v", detail.DestZip)
	}
	if detail.CostGross == nil || *detail.CostGross != "24.90" {
		t.Fatalf("CostGross = %#v, want 24.90", detail.CostGross)
	}
	if detail.CostSeller == nil || *detail.CostSeller != "12.45" {
		t.Fatalf("CostSeller = %#v, want 12.45", detail.CostSeller)
	}
	if detail.Currency == nil || *detail.Currency == "" {
		t.Fatalf("Currency = %#v, want a resolved currency", detail.Currency)
	}

	// SLALimitAt comes from lead_time.estimated_delivery_limit.date on this SAME primary
	// payload (NOT the separate /sla sub-resource) — same field the existing, still-LIVE
	// mapShipmentInfo already decodes (shipping_reader.go:236-237).
	wantSLALimit := time.Date(2026, 7, 26, 3, 59, 59, 0, time.UTC)
	if detail.SLALimitAt == nil || !detail.SLALimitAt.Equal(wantSLALimit) {
		t.Fatalf("SLALimitAt = %#v, want %v", detail.SLALimitAt, wantSLALimit)
	}

	// scrubbedShipmentFixture above carries neither `logistic` nor `tracking_method` — both
	// are now decoded (mlIngestShipmentResponse), but absence on THIS fixture must still
	// degrade to nil (ADR-17), never fabricated. See TestGetShipmentDetailDecodesLogisticAndTrackingMethod
	// for the positive-value assertion.
	if detail.LogisticType != nil {
		t.Fatalf("LogisticType = %#v, want nil (not present on this fixture)", detail.LogisticType)
	}
	if detail.TrackingMethod != nil {
		t.Fatalf("TrackingMethod = %#v, want nil (not present on this fixture)", detail.TrackingMethod)
	}
	// SLAStatus stays nil: the separate GET /shipments/{id}/sla endpoint is deliberately
	// not called (unconfirmed response shape, tracked honest gap) — unlike SLALimitAt, it
	// has no source on the primary payload.
	if detail.SLAStatus != nil {
		t.Fatalf("SLAStatus = %#v, want nil (sla sub-resource not called, tracked gap)", detail.SLAStatus)
	}
}

func TestGetShipmentDetailDecodesLogisticAndTrackingMethod(t *testing.T) {
	t.Parallel()

	// logistic.type/tracking_method were previously undeclared on mlIngestShipmentResponse
	// (encoding/json silently dropped both — 0/38 rows in order_shipments). This proves the
	// end-to-end wire-to-domain path, not just DTO declaration (rawkeys only proves the
	// latter).
	const fixtureWithLogistic = `{
	  "id": 44556677,
	  "site_id": "MLB",
	  "status": "delivered",
	  "tracking_number": "BR123456789",
	  "tracking_method": "Normal",
	  "logistic": {"type": "drop_off", "mode": "me2", "direction": "forward"}
	}`
	server := newShipmentTestServer(t, http.StatusOK, fixtureWithLogistic, http.StatusOK, scrubbedShipmentCostsFixture)
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if err != nil {
		t.Fatalf("GetShipmentDetail() error = %v", err)
	}
	if detail.LogisticType == nil || *detail.LogisticType != "drop_off" {
		t.Fatalf("LogisticType = %#v, want \"drop_off\"", detail.LogisticType)
	}
	if detail.TrackingMethod == nil || *detail.TrackingMethod != "Normal" {
		t.Fatalf("TrackingMethod = %#v, want \"Normal\"", detail.TrackingMethod)
	}
}

func TestGetShipmentDetailSLALimitAtAbsentStaysNil(t *testing.T) {
	t.Parallel()

	// No lead_time block at all in this payload — SLALimitAt must degrade to nil, not be
	// fabricated, exactly like the existing mapShipmentInfo behaves when lead_time is absent.
	const fixtureWithoutLeadTime = `{
	  "id": 44556677,
	  "site_id": "MLB",
	  "status": "shipped",
	  "tracking_number": "BR123456789"
	}`
	server := newShipmentTestServer(t, http.StatusOK, fixtureWithoutLeadTime, http.StatusOK, scrubbedShipmentCostsFixture)
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if err != nil {
		t.Fatalf("GetShipmentDetail() error = %v", err)
	}
	if detail.SLALimitAt != nil {
		t.Fatalf("SLALimitAt = %#v, want nil (no lead_time in payload)", detail.SLALimitAt)
	}
}

func TestGetShipmentDetailNotFoundIsHonestAbsence(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusNotFound, http.StatusGone} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := newShipmentTestServer(t, status, "", http.StatusOK, scrubbedShipmentCostsFixture)
			defer server.Close()
			// Costs handler won't be reached in the 404/410 case, but the server tolerates it.

			detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
			if err != nil {
				t.Fatalf("GetShipmentDetail() error = %v, want nil (honest-absence)", err)
			}
			if detail.Found {
				t.Fatal("Found = true, want false")
			}
			if detail != (domain.ShipmentDetail{}) {
				t.Fatalf("detail = %#v, want zero value", detail)
			}
		})
	}
}

func TestGetShipmentDetailCostsDegradeOn404(t *testing.T) {
	t.Parallel()

	server := newShipmentTestServer(t, http.StatusOK, scrubbedShipmentFixture, http.StatusNotFound, "")
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if err != nil {
		t.Fatalf("GetShipmentDetail() error = %v, want nil (costs degrade silently)", err)
	}
	if !detail.Found {
		t.Fatal("Found = false, want true (primary succeeded)")
	}
	if detail.CostGross != nil || detail.CostSeller != nil || detail.Currency != nil {
		t.Fatalf("cost fields = %#v/%#v/%#v, want all nil on costs 404", detail.CostGross, detail.CostSeller, detail.Currency)
	}
	if detail.TrackingNumber == nil {
		t.Fatal("TrackingNumber = nil, want the rest of the shipment preserved")
	}
}

func TestGetShipmentDetailCostsRealErrorPropagates(t *testing.T) {
	t.Parallel()

	server := newShipmentTestServer(t, http.StatusOK, scrubbedShipmentFixture, http.StatusInternalServerError, "")
	defer server.Close()

	detail, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if !errors.Is(err, domain.ErrProviderUnavailable) {
		t.Fatalf("GetShipmentDetail error = %v, want ErrProviderUnavailable", err)
	}
	if detail != (domain.ShipmentDetail{}) {
		t.Fatalf("detail = %#v, want zero value on real error", detail)
	}
}

func TestGetShipmentDetailPrimaryErrorPropagates(t *testing.T) {
	t.Parallel()

	server := newShipmentTestServer(t, http.StatusForbidden, "", http.StatusOK, scrubbedShipmentCostsFixture)
	defer server.Close()

	_, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "44556677")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("GetShipmentDetail error = %v, want ErrUnauthorized", err)
	}
}

func TestGetShipmentDetailEmptyShipmentIDIsInvalidReference(t *testing.T) {
	t.Parallel()

	_, err := pricingTestAdapter("http://unused.invalid", time.Now().UTC()).GetShipmentDetail(context.Background(), pricingAccountRef(), "  ")
	if domain.ErrorCodeOf(err) != domain.ErrCodeProviderInvalidReference {
		t.Fatalf("error code = %v, want %v", domain.ErrorCodeOf(err), domain.ErrCodeProviderInvalidReference)
	}
}
