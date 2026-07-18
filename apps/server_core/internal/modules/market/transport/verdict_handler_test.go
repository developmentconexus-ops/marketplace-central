package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

type fakeEvidenceReader struct {
	signals         []domain.CompetitiveSignal
	aggregates      []domain.MarketAggregate
	verdicts        []domain.Verdict
	signalsErr      error
	aggregatesErr   error
	verdictsErr     error
	gotSignalIDs    []string
	gotAggregateIDs []string
	gotVerdictIDs   []string
}

func (f *fakeEvidenceReader) Signals(_ context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error) {
	f.gotSignalIDs = append([]string(nil), listingIDs...)
	if len(listingIDs) > application.MaxReadIDs {
		return nil, &application.TooManyIDsError{Count: len(listingIDs), Limit: application.MaxReadIDs}
	}
	return f.signals, f.signalsErr
}

func (f *fakeEvidenceReader) Aggregates(_ context.Context, codprods []string) ([]domain.MarketAggregate, error) {
	f.gotAggregateIDs = append([]string(nil), codprods...)
	if len(codprods) > application.MaxReadIDs {
		return nil, &application.TooManyIDsError{Count: len(codprods), Limit: application.MaxReadIDs}
	}
	return f.aggregates, f.aggregatesErr
}

func (f *fakeEvidenceReader) Verdicts(_ context.Context, codprods []string) ([]domain.Verdict, error) {
	f.gotVerdictIDs = append([]string(nil), codprods...)
	if len(codprods) > application.MaxReadIDs {
		return nil, &application.TooManyIDsError{Count: len(codprods), Limit: application.MaxReadIDs}
	}
	return f.verdicts, f.verdictsErr
}

func TestVerdictsSerializesB2ShapeWithNullLabelAndMarketRange(t *testing.T) {
	t.Parallel()

	blocking := domain.BlockingStateSemCusto
	evidence := &fakeEvidenceReader{verdicts: []domain.Verdict{{
		ProductID:           "123",
		MatchStatus:         domain.MatchStatusAccept,
		PriceEvidenceStatus: domain.PriceEvidenceStatusOK,
		VerdictLabel:        nil,
		BlockingState:       &blocking,
		InputsUsed: map[string]domain.VerdictInputRef{
			"our_price":    {Source: "ml_sale_price", AsOf: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)},
			"target_price": {Source: "ml_price_to_win", AsOf: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)},
			"market":       {Source: "ml_catalog_offers", AsOf: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)},
		},
		MarketRange: &domain.VerdictMarketRange{
			MinValid: &domain.Money{Amount: "72.33", Currency: "BRL"},
			Median:   &domain.Money{Amount: "88.88", Currency: "BRL"},
			Currency: "BRL",
			NOffers:  6,
			NSellers: 5,
		},
	}}}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, evidence)
	req := httptest.NewRequest(http.MethodGet, "/market/verdicts?codprod=123", nil)
	rr := httptest.NewRecorder()

	h.handleVerdicts(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	body := rr.Body.Bytes()
	for _, field := range []string{
		`"match_status":"ACCEPT"`,
		`"price_evidence_status":"OK"`,
		`"verdict_label":null`,
		`"blocking_state":"SEM_CUSTO"`,
		`"our_price":{"source":"ml_sale_price","as_of":"2026-07-18T12:00:00Z"}`,
		`"market_range":{"min_valid":{"amount":"72.33","currency":"BRL"},"median":{"amount":"88.88","currency":"BRL"},"currency":"BRL","n_offers":6,"n_sellers":5}`,
	} {
		assertBodyContains(t, body, field)
	}
	if strings.Contains(string(body), `"max"`) {
		t.Fatalf("market_range must not invent a max field: %s", body)
	}
}

func TestVerdictsNeverCollectedReturns200NoPriceEvidence(t *testing.T) {
	t.Parallel()

	evidence := &fakeEvidenceReader{verdicts: []domain.Verdict{{
		ProductID:           "never-collected",
		MatchStatus:         domain.MatchStatusNoCandidate,
		PriceEvidenceStatus: domain.PriceEvidenceStatusNoPriceEvidence,
		BlockingState:       func() *domain.BlockingState { s := domain.BlockingStateNoCandidate; return &s }(),
		InputsUsed:          map[string]domain.VerdictInputRef{},
	}}}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, evidence)
	req := httptest.NewRequest(http.MethodGet, "/market/verdicts?codprod=never-collected", nil)
	rr := httptest.NewRecorder()

	h.handleVerdicts(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	assertBodyContains(t, rr.Body.Bytes(), `"price_evidence_status":"NO_PRICE_EVIDENCE"`)
	assertBodyContains(t, rr.Body.Bytes(), `"market_range":null`)
}

func TestVerdictsTwoHundredOneIDsReturns422(t *testing.T) {
	t.Parallel()

	ids := make([]string, 201)
	for i := range ids {
		ids[i] = "p_" + strconv.Itoa(i)
	}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodGet, "/market/verdicts?codprod="+strings.Join(ids, ","), nil)
	rr := httptest.NewRecorder()

	h.handleVerdicts(rr, req)

	if got, want := rr.Code, http.StatusUnprocessableEntity; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"too_many_ids"`)
}

func TestSignalsPreservesInputOrderAndEmptyArray(t *testing.T) {
	t.Parallel()

	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, &fakeEvidenceReader{signals: []domain.CompetitiveSignal{}})
	req := httptest.NewRequest(http.MethodGet, "/market/signals", nil)
	rr := httptest.NewRecorder()

	h.handleSignals(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	assertBodyContains(t, rr.Body.Bytes(), `[]`)
}

func TestSignalsSerializesNullableFieldsWithoutOmission(t *testing.T) {
	t.Parallel()

	evidence := &fakeEvidenceReader{signals: []domain.CompetitiveSignal{{
		ListingID: "listing_1",
		OurPrice:  &domain.Money{Amount: "10.00", Currency: "BRL"},
		FetchedAt: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
	}}}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, evidence)
	req := httptest.NewRequest(http.MethodGet, "/market/signals?listing_ids=listing_1", nil)
	rr := httptest.NewRecorder()

	h.handleSignals(rr, req)

	body := rr.Body.Bytes()
	for _, field := range []string{
		`"listing_id":"listing_1"`,
		`"our_price":{"amount":"10.00","currency":"BRL"}`,
		`"winner_price":null`,
		`"target_price":null`,
		`"position":null`,
		`"fetched_at":"2026-07-18T12:00:00Z"`,
	} {
		assertBodyContains(t, body, field)
	}
}

func TestAggregatesTwoHundredOneIDsReturns422(t *testing.T) {
	t.Parallel()

	ids := make([]string, 201)
	for i := range ids {
		ids[i] = "p_" + strconv.Itoa(i)
	}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodGet, "/market/aggregates?codprod="+strings.Join(ids, ","), nil)
	rr := httptest.NewRecorder()

	h.handleAggregates(rr, req)

	if got, want := rr.Code, http.StatusUnprocessableEntity; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"too_many_ids"`)
}

func TestAggregatesSerializesNullableMoneyWithoutOmission(t *testing.T) {
	t.Parallel()

	evidence := &fakeEvidenceReader{aggregates: []domain.MarketAggregate{{
		ProductID:  "123",
		Status:     domain.MarketAggregateStatusNoPriceEvidence,
		Source:     domain.MarketPriceSourceMLCatalogOffers,
		FetchedAt:  time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
		ComputedAt: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
	}}}
	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, evidence)
	req := httptest.NewRequest(http.MethodGet, "/market/aggregates?codprod=123", nil)
	rr := httptest.NewRecorder()

	h.handleAggregates(rr, req)

	body := rr.Body.Bytes()
	for _, field := range []string{
		`"product_id":"123"`,
		`"median":null`,
		`"min_valid":null`,
		`"status":"NO_PRICE_EVIDENCE"`,
	} {
		assertBodyContains(t, body, field)
	}
}
