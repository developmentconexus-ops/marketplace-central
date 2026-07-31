package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func TestMarketErrorEnvelopeWholeBody(t *testing.T) {
	t.Parallel()

	h := NewHandler(&fakeReadService{})
	req := httptest.NewRequest(http.MethodGet, "/market/observations", nil)
	rr := httptest.NewRecorder()

	h.handleObservations(rr, req)

	if got, want := rr.Code, http.StatusBadRequest; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	want := `{"error":{"code":"installation_required","message":"installation_id é obrigatório","details":{"key":"installation_id"}}}`
	if trimJSON(rr.Body.String()) != trimJSON(want) {
		t.Fatalf("body = %s, want %s", trimJSON(rr.Body.String()), trimJSON(want))
	}
}

type fakeReadService struct {
	observations    application.ObservationReadResult
	references      application.ReferenceReadResult
	observationsErr error
	referencesErr   error
	observationIDs  []string
	referenceIDs    []string
}

func (f *fakeReadService) ListObservations(_ context.Context, ids []string) (application.ObservationReadResult, error) {
	f.observationIDs = append([]string(nil), ids...)
	if len(ids) > application.MaxReadIDs {
		return application.ObservationReadResult{}, &application.TooManyIDsError{Count: len(ids), Limit: application.MaxReadIDs}
	}
	return f.observations, f.observationsErr
}

func (f *fakeReadService) ListReferences(_ context.Context, ids []string) (application.ReferenceReadResult, error) {
	f.referenceIDs = append([]string(nil), ids...)
	if len(ids) > application.MaxReadIDs {
		return application.ReferenceReadResult{}, &application.TooManyIDsError{Count: len(ids), Limit: application.MaxReadIDs}
	}
	return f.references, f.referencesErr
}

func TestObservationsRequiresInstallation(t *testing.T) {
	t.Parallel()

	h := NewHandler(&fakeReadService{})
	req := httptest.NewRequest(http.MethodGet, "/market/observations?listing_ids=garbage", nil)
	rr := httptest.NewRecorder()

	h.handleObservations(rr, req)

	if got, want := rr.Code, http.StatusBadRequest; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"installation_required"`)
}

func TestObservationsTwoHundredOneIDsReturns422(t *testing.T) {
	t.Parallel()

	ids := make([]string, 201)
	for i := range ids {
		ids[i] = "listing_" + strconv.Itoa(i)
	}
	h := NewHandler(&fakeReadService{})
	req := httptest.NewRequest(http.MethodGet, "/market/observations?installation_id=inst_1&listing_ids="+strings.Join(ids, ","), nil)
	rr := httptest.NewRecorder()

	h.handleObservations(rr, req)

	if got, want := rr.Code, http.StatusUnprocessableEntity; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"too_many_ids"`)
}

func TestMalformedListingIDRemains200(t *testing.T) {
	t.Parallel()

	service := &fakeReadService{observations: application.ObservationReadResult{
		Items: []domain.MarketObservation{{
			ListingID:     "garbage",
			EvidenceState: domain.EvidenceStateNoPriceEvidence,
		}},
		AsOf: time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC),
	}}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/market/observations?installation_id=inst_1&listing_ids=garbage", nil)
	rr := httptest.NewRecorder()

	h.handleObservations(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	body := rr.Body.Bytes()
	for _, field := range []string{
		`"our_sale_price":null`,
		`"competitive_status":null`,
		`"winner_price":null`,
		`"competitive_target":null`,
		`"catalog_offer_price":null`,
		`"catalog_stats":null`,
		`"captured_at":null`,
		`"source":null`,
	} {
		assertBodyContains(t, body, field)
	}
	assertBodyContains(t, body, `"evidence_state":"no_price_evidence"`)
	assertBodyContains(t, body, `"as_of":"2026-07-16T12:00:00Z"`)
}

func TestObservationsStoredRowSerializesMoneyObjects(t *testing.T) {
	t.Parallel()

	status := domain.CompetitiveStatusWinning
	h := NewHandler(&fakeReadService{observations: application.ObservationReadResult{
		Items: []domain.MarketObservation{{
			ListingID:         "listing_1",
			EvidenceState:     domain.EvidenceStateObserved,
			OurSalePrice:      &domain.Money{Amount: "12.34", Currency: "BRL"},
			CompetitiveStatus: &status,
		}},
		AsOf: time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC),
	}})
	req := httptest.NewRequest(http.MethodGet, "/market/observations?installation_id=inst_1&listing_ids=listing_1", nil)
	rr := httptest.NewRecorder()

	h.handleObservations(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"our_sale_price":{"amount":"12.34","currency":"BRL"}`)
	assertBodyContains(t, rr.Body.Bytes(), `"as_of":"2026-07-16T12:00:00Z"`)
}

func TestReferencesTwoHundredOneIDsReturns422(t *testing.T) {
	t.Parallel()

	ids := make([]string, 201)
	for i := range ids {
		ids[i] = "product_" + strconv.Itoa(i)
	}
	h := NewHandler(&fakeReadService{})
	req := httptest.NewRequest(http.MethodGet, "/market/references?product_ids="+strings.Join(ids, ","), nil)
	rr := httptest.NewRecorder()

	h.handleReferences(rr, req)

	if got, want := rr.Code, http.StatusUnprocessableEntity; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"too_many_ids"`)
}

func assertBodyContains(t *testing.T, body []byte, want string) {
	t.Helper()
	if !strings.Contains(string(body), want) {
		t.Fatalf("body = %s, want substring %s", body, want)
	}
}
