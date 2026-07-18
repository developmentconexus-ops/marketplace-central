package application_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

type fakeEvidenceDecisionStore struct {
	rows      []domain.MatchDecision
	called    int
	requested []string
}

func (f *fakeEvidenceDecisionStore) AppendMatchDecision(context.Context, domain.MatchDecision) error {
	return nil
}

func (f *fakeEvidenceDecisionStore) LatestMatchDecisions(_ context.Context, productIDs []string) ([]domain.MatchDecision, error) {
	f.called++
	f.requested = append([]string(nil), productIDs...)
	return f.rows, nil
}

type fakeEvidenceSignalStore struct {
	rows      []domain.CompetitiveSignal
	called    int
	requested []string
}

func (f *fakeEvidenceSignalStore) AppendCompetitiveSignals(context.Context, []domain.CompetitiveSignal) error {
	return nil
}

func (f *fakeEvidenceSignalStore) LatestCompetitiveSignals(_ context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error) {
	f.called++
	f.requested = append([]string(nil), listingIDs...)
	return f.rows, nil
}

type fakeEvidenceAggregateStore struct {
	byProduct       []domain.MarketAggregate
	bySource        []domain.MarketAggregate
	calledByProduct int
	calledBySource  int
}

func (f *fakeEvidenceAggregateStore) AppendMarketAggregates(context.Context, []domain.MarketAggregate) error {
	return nil
}

func (f *fakeEvidenceAggregateStore) LatestMarketAggregates(_ context.Context, productIDs []string) ([]domain.MarketAggregate, error) {
	f.calledByProduct++
	return f.byProduct, nil
}

func (f *fakeEvidenceAggregateStore) LatestMarketAggregatesBySource(_ context.Context, productIDs []string, source domain.MarketPriceSource) ([]domain.MarketAggregate, error) {
	f.calledBySource++
	return f.bySource, nil
}

type fakeEvidenceCostReader struct {
	byProductID map[int]*domain.Money
	asOf        time.Time
	err         error
}

func (f *fakeEvidenceCostReader) GetCostAsOf(_ context.Context, productID int, _ time.Time) (*domain.Money, time.Time, error) {
	if f.err != nil {
		return nil, time.Time{}, f.err
	}
	if money, ok := f.byProductID[productID]; ok {
		return money, f.asOf, nil
	}
	return nil, time.Time{}, nil
}

func fixedNow() time.Time {
	return time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
}

func newEvidenceReadService(
	decisions *fakeEvidenceDecisionStore,
	aggregates *fakeEvidenceAggregateStore,
	signals *fakeEvidenceSignalStore,
	cost *fakeEvidenceCostReader,
) *application.EvidenceReadService {
	return application.NewEvidenceReadService(decisions, aggregates, signals, cost, fixedNow)
}

func okAggregate(productID string, source domain.MarketPriceSource, status domain.MarketAggregateStatus) domain.MarketAggregate {
	v, err := domain.NewMarketAggregate(domain.MarketAggregateInput{
		ProductID: productID, Median: &domain.Money{Amount: "50", Currency: "BRL"}, NOffers: 6, NSellers: 5,
		Source: source, FetchedAt: fixedNow(), ComputedAt: fixedNow(), Status: status,
	})
	if err != nil {
		panic(err)
	}
	return v
}

func TestVerdictsNeverCollectedCodprodReturnsHonestNoCandidate(t *testing.T) {
	service := newEvidenceReadService(&fakeEvidenceDecisionStore{}, &fakeEvidenceAggregateStore{}, &fakeEvidenceSignalStore{}, &fakeEvidenceCostReader{})

	got, err := service.Verdicts(context.Background(), []string{"never-collected"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("verdicts = %d, want 1", len(got))
	}
	v := got[0]
	if v.MatchStatus != domain.MatchStatusNoCandidate || v.PriceEvidenceStatus != domain.PriceEvidenceStatusNoPriceEvidence {
		t.Fatalf("verdict = %+v, want NO_CANDIDATE/NO_PRICE_EVIDENCE", v)
	}
	if v.BlockingState == nil || *v.BlockingState != domain.BlockingStateNoCandidate {
		t.Fatalf("blocking = %v, want NO_CANDIDATE", v.BlockingState)
	}
}

func TestVerdictsAcceptOKAggregateCostPresentYieldsOK(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	cost := &fakeEvidenceCostReader{byProductID: map[int]*domain.Money{10: {Amount: "20", Currency: "BRL"}}, asOf: fixedNow()}
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		cost,
	)

	got, err := service.Verdicts(context.Background(), []string{"10"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	v := got[0]
	if v.PriceEvidenceStatus != domain.PriceEvidenceStatusOK || v.BlockingState != nil {
		t.Fatalf("verdict = %+v, want OK/nil blocking", v)
	}
	if v.MarketRange == nil || v.MarketRange.Median.Amount != "50" {
		t.Fatalf("market_range = %+v, want to mirror aggregate", v.MarketRange)
	}
	if _, ok := v.InputsUsed["market"]; !ok {
		t.Fatalf("inputs_used = %+v, want market key", v.InputsUsed)
	}
	if _, ok := v.InputsUsed["erp_cost"]; !ok {
		t.Fatalf("inputs_used = %+v, want erp_cost key", v.InputsUsed)
	}
}

func TestVerdictsAcceptOKAggregateCostAbsentYieldsSemCusto(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		&fakeEvidenceCostReader{},
	)

	got, err := service.Verdicts(context.Background(), []string{"10"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	v := got[0]
	if v.PriceEvidenceStatus != domain.PriceEvidenceStatusOK {
		t.Fatalf("price_evidence_status = %q, want OK", v.PriceEvidenceStatus)
	}
	if v.BlockingState == nil || *v.BlockingState != domain.BlockingStateSemCusto {
		t.Fatalf("blocking = %v, want SEM_CUSTO", v.BlockingState)
	}
	if _, ok := v.InputsUsed["erp_cost"]; ok {
		t.Fatalf("inputs_used = %+v, want no erp_cost key when cost absent", v.InputsUsed)
	}
}

func TestVerdictsAcceptInsufficientAggregateYieldsInsufficientMarket(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusInsufficientMarket)
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		&fakeEvidenceCostReader{},
	)

	got, err := service.Verdicts(context.Background(), []string{"10"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	v := got[0]
	if v.PriceEvidenceStatus != domain.PriceEvidenceStatusInsufficientMarket {
		t.Fatalf("price_evidence_status = %q, want INSUFFICIENT_MARKET", v.PriceEvidenceStatus)
	}
	if v.BlockingState == nil || *v.BlockingState != domain.BlockingStateInsufficientMarket {
		t.Fatalf("blocking = %v, want INSUFFICIENT_MARKET", v.BlockingState)
	}
}

func TestVerdictsAggregateNoPriceEvidenceStatusTreatedAsMissing(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusNoPriceEvidence)
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		&fakeEvidenceCostReader{},
	)

	got, err := service.Verdicts(context.Background(), []string{"10"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	v := got[0]
	if v.PriceEvidenceStatus != domain.PriceEvidenceStatusNoPriceEvidence || v.MarketRange != nil {
		t.Fatalf("verdict = %+v, want NO_PRICE_EVIDENCE with nil range", v)
	}
}

func TestSignalsHonestEmptyForMissingListing(t *testing.T) {
	stored := domain.CompetitiveSignal{ListingID: "listing-known", OurPrice: &domain.Money{Amount: "99", Currency: "BRL"}, FetchedAt: fixedNow()}
	service := newEvidenceReadService(&fakeEvidenceDecisionStore{}, &fakeEvidenceAggregateStore{}, &fakeEvidenceSignalStore{rows: []domain.CompetitiveSignal{stored}}, &fakeEvidenceCostReader{})

	got, err := service.Signals(context.Background(), []string{"listing-missing", "listing-known"})
	if err != nil {
		t.Fatalf("Signals() error = %v, want nil", err)
	}
	if len(got) != 2 {
		t.Fatalf("signals = %d, want 2", len(got))
	}
	missing := got[0]
	if missing.ListingID != "listing-missing" || !missing.FetchedAt.IsZero() {
		t.Fatalf("missing signal = %+v, want zero FetchedAt", missing)
	}
	if missing.OurPrice != nil || missing.WinnerPrice != nil || missing.TargetPrice != nil || missing.Position != nil {
		t.Fatalf("missing signal = %+v, want all money/position nil", missing)
	}
	if !reflect.DeepEqual(got[1], stored) {
		t.Fatalf("known signal = %+v, want stored unchanged", got[1])
	}
}

func TestAggregatesHonestEmptyForMissingCodprod(t *testing.T) {
	stored := okAggregate("codprod-known", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	service := newEvidenceReadService(&fakeEvidenceDecisionStore{}, &fakeEvidenceAggregateStore{byProduct: []domain.MarketAggregate{stored}}, &fakeEvidenceSignalStore{}, &fakeEvidenceCostReader{})

	got, err := service.Aggregates(context.Background(), []string{"codprod-missing", "codprod-known"})
	if err != nil {
		t.Fatalf("Aggregates() error = %v, want nil", err)
	}
	if len(got) != 2 {
		t.Fatalf("aggregates = %d, want 2", len(got))
	}
	missing := got[0]
	if missing.ProductID != "codprod-missing" || missing.Status != domain.MarketAggregateStatusNoPriceEvidence || missing.NOffers != 0 || missing.NSellers != 0 {
		t.Fatalf("missing aggregate = %+v, want honest NO_PRICE_EVIDENCE n=0", missing)
	}
	if missing.Median != nil || missing.MinValid != nil {
		t.Fatalf("missing aggregate = %+v, want no invented price", missing)
	}
	if !reflect.DeepEqual(got[1], stored) {
		t.Fatalf("known aggregate = %+v, want stored unchanged", got[1])
	}
}

func TestVerdictsBatchOrderPreserved(t *testing.T) {
	decisions := []domain.MatchDecision{
		{ProductID: "b", Status: domain.MatchStatusNoCandidate},
		{ProductID: "a", Status: domain.MatchStatusNoCandidate},
	}
	service := newEvidenceReadService(&fakeEvidenceDecisionStore{rows: decisions}, &fakeEvidenceAggregateStore{}, &fakeEvidenceSignalStore{}, &fakeEvidenceCostReader{})

	got, err := service.Verdicts(context.Background(), []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	if len(got) != 3 || got[0].ProductID != "a" || got[1].ProductID != "b" || got[2].ProductID != "c" {
		t.Fatalf("verdicts = %+v, want input order a,b,c", got)
	}
}

func TestEvidenceReadServiceTooManyIDs(t *testing.T) {
	decisions := &fakeEvidenceDecisionStore{}
	aggregates := &fakeEvidenceAggregateStore{}
	signals := &fakeEvidenceSignalStore{}
	service := newEvidenceReadService(decisions, aggregates, signals, &fakeEvidenceCostReader{})

	ids := make([]string, 200)
	for i := range ids {
		ids[i] = "id"
	}
	if _, err := service.Verdicts(context.Background(), ids); err != nil {
		t.Fatalf("Verdicts(200 ids) error = %v, want nil", err)
	}
	if decisions.called != 1 {
		t.Fatalf("decision store calls after 200 ids = %d, want 1", decisions.called)
	}

	tooMany := append(append([]string(nil), ids...), "id-201")
	_, err := service.Verdicts(context.Background(), tooMany)
	var tooManyErr *application.TooManyIDsError
	if !errors.As(err, &tooManyErr) || tooManyErr.Code() != "too_many_ids" {
		t.Fatalf("error = %v, want typed too_many_ids", err)
	}
	if decisions.called != 1 {
		t.Fatalf("decision store calls after 201 ids = %d, want unchanged at 1", decisions.called)
	}

	if _, err := service.Signals(context.Background(), tooMany); !errors.As(err, &tooManyErr) {
		t.Fatalf("Signals(201 ids) error = %v, want typed too_many_ids", err)
	}
	if signals.called != 0 {
		t.Fatalf("signal store calls after 201 ids = %d, want 0", signals.called)
	}

	if _, err := service.Aggregates(context.Background(), tooMany); !errors.As(err, &tooManyErr) {
		t.Fatalf("Aggregates(201 ids) error = %v, want typed too_many_ids", err)
	}
	if aggregates.calledByProduct != 0 {
		t.Fatalf("aggregate store calls after 201 ids = %d, want 0", aggregates.calledByProduct)
	}
}

func TestVerdictsNonNumericCodprodTreatsCostAsAbsentWithoutBatchFailure(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "not-a-number", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("not-a-number", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	// A cost reader that would error for any call proves the service never
	// dials out for a non-numeric codprod.
	cost := &fakeEvidenceCostReader{err: errors.New("must not be called")}
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		cost,
	)

	got, err := service.Verdicts(context.Background(), []string{"not-a-number"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil (non-numeric codprod must not fail the batch)", err)
	}
	v := got[0]
	if v.PriceEvidenceStatus != domain.PriceEvidenceStatusOK {
		t.Fatalf("price_evidence_status = %q, want OK", v.PriceEvidenceStatus)
	}
	if v.BlockingState == nil || *v.BlockingState != domain.BlockingStateSemCusto {
		t.Fatalf("blocking = %v, want SEM_CUSTO (cost absent for non-numeric codprod)", v.BlockingState)
	}
}

func TestVerdictsCostReaderHardErrorPropagates(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	cost := &fakeEvidenceCostReader{err: errors.New("erp connection failed")}
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		cost,
	)

	if _, err := service.Verdicts(context.Background(), []string{"10"}); err == nil {
		t.Fatal("Verdicts() error = nil, want cost reader hard error to propagate")
	}
}

func TestVerdictsDeterministic(t *testing.T) {
	decision := domain.MatchDecision{ProductID: "10", Status: domain.MatchStatusAccept}
	aggregate := okAggregate("10", domain.MarketPriceSourceMLCatalogOffers, domain.MarketAggregateStatusOK)
	cost := &fakeEvidenceCostReader{byProductID: map[int]*domain.Money{10: {Amount: "20", Currency: "BRL"}}, asOf: fixedNow()}
	service := newEvidenceReadService(
		&fakeEvidenceDecisionStore{rows: []domain.MatchDecision{decision}},
		&fakeEvidenceAggregateStore{bySource: []domain.MarketAggregate{aggregate}},
		&fakeEvidenceSignalStore{},
		cost,
	)

	got1, err := service.Verdicts(context.Background(), []string{"10", "never-collected"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	got2, err := service.Verdicts(context.Background(), []string{"10", "never-collected"})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("non-deterministic: %#v != %#v", got1, got2)
	}
}

func TestVerdictsEmptyInputReturnsEmptyNonNilSlice(t *testing.T) {
	service := newEvidenceReadService(&fakeEvidenceDecisionStore{}, &fakeEvidenceAggregateStore{}, &fakeEvidenceSignalStore{}, &fakeEvidenceCostReader{})

	got, err := service.Verdicts(context.Background(), []string{})
	if err != nil {
		t.Fatalf("Verdicts() error = %v, want nil", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("verdicts = %#v, want empty non-nil slice", got)
	}
}
