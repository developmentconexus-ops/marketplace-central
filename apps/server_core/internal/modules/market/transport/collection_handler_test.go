package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

type fakeCollectionService struct {
	summary    application.CollectionSummary
	err        error
	gotCodprod string
}

func (f *fakeCollectionService) Collect(_ context.Context, codprod string) (application.CollectionSummary, error) {
	f.gotCodprod = codprod
	return f.summary, f.err
}

func TestCollectHappyPathReturnsB1Envelope(t *testing.T) {
	t.Parallel()

	blocking := domain.BlockingStateSemCusto
	service := &fakeCollectionService{summary: application.CollectionSummary{
		Codprod: "123",
		Status:  application.CollectionStatusCompleted,
		Decisoes: []application.CollectionDecision{{
			Codprod:       "123",
			MatchStatus:   domain.MatchStatusAccept,
			PriceEvidence: domain.PriceEvidenceStatusOK,
			Blocking:      &blocking,
		}},
		Contagens: map[string]int{"ok": 0, "no_price_evidence": 0, "insufficient_market": 0, "no_candidate": 0, "sem_custo": 1},
		Causas:    []application.CollectionCauseDetail{{Cause: application.CollectionCauseProvider4xx}},
	}}
	h := NewHandlerWithCollections(&fakeReadService{}, service, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":"123"}`))
	rr := httptest.NewRecorder()

	h.handleCollect(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	if service.gotCodprod != "123" {
		t.Fatalf("gotCodprod = %q, want 123", service.gotCodprod)
	}
	body := rr.Body.Bytes()
	for _, field := range []string{
		`"status":"COMPLETED"`,
		`"decisões":[{`,
		`"codprod":"123"`,
		`"match_status":"ACCEPT"`,
		`"price_evidence_status":"OK"`,
		`"blocking_state":"SEM_CUSTO"`,
		`"contagens":{`,
		`"sem_custo":1`,
		`"causas":[{`,
		`"reason":"PROVIDER_4XX"`,
		`"detail":null`,
	} {
		assertBodyContains(t, body, field)
	}
}

func TestCollectSerializesCauseDetailWhenPresent(t *testing.T) {
	t.Parallel()

	detail := "catalog offers fetch: provider responded HTTP 503"
	service := &fakeCollectionService{summary: application.CollectionSummary{
		Codprod:   "123",
		Status:    application.CollectionStatusPartial,
		Contagens: map[string]int{"ok": 0, "no_price_evidence": 0, "insufficient_market": 0, "no_candidate": 1, "sem_custo": 0},
		Causas:    []application.CollectionCauseDetail{{Cause: application.CollectionCauseProvider5xx, Detail: &detail}},
	}}
	h := NewHandlerWithCollections(&fakeReadService{}, service, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":"123"}`))
	rr := httptest.NewRecorder()

	h.handleCollect(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	body := rr.Body.Bytes()
	for _, field := range []string{
		`"reason":"PROVIDER_5XX"`,
		`"detail":"catalog offers fetch: provider responded HTTP 503"`,
	} {
		assertBodyContains(t, body, field)
	}
}

func TestCollectConcurrentReturns409(t *testing.T) {
	t.Parallel()

	service := &fakeCollectionService{err: application.ErrCollectionInProgress}
	h := NewHandlerWithCollections(&fakeReadService{}, service, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":"123"}`))
	rr := httptest.NewRecorder()

	h.handleCollect(rr, req)

	if got, want := rr.Code, http.StatusConflict; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"COLLECTION_IN_PROGRESS"`)
}

func TestCollectUnknownProductReturns404(t *testing.T) {
	t.Parallel()

	service := &fakeCollectionService{err: ports.ErrProductNotFound}
	h := NewHandlerWithCollections(&fakeReadService{}, service, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":"does-not-exist"}`))
	rr := httptest.NewRecorder()

	h.handleCollect(rr, req)

	if got, want := rr.Code, http.StatusNotFound; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
	assertBodyContains(t, rr.Body.Bytes(), `"code":"PRODUCT_NOT_FOUND"`)
}

func TestCollectMissingCodprodReturns400(t *testing.T) {
	t.Parallel()

	h := NewHandlerWithCollections(&fakeReadService{}, &fakeCollectionService{}, &fakeEvidenceReader{})
	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{}`))
	rr := httptest.NewRecorder()

	h.handleCollect(rr, req)

	if got, want := rr.Code, http.StatusBadRequest; got != want {
		t.Fatalf("status = %d, want %d, body=%s", got, want, rr.Body.String())
	}
}
