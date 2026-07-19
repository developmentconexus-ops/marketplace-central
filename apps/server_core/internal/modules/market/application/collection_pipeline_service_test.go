package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

// --- fakes -------------------------------------------------------------

type fakeIdentityReader struct {
	identities  map[string]domain.LocalProductIdentity
	identityErr map[string]error
	linked      map[string][]string
	linkedErr   map[string]error
}

func (f *fakeIdentityReader) GetLocalIdentity(_ context.Context, codprod string) (domain.LocalProductIdentity, error) {
	if err, ok := f.identityErr[codprod]; ok {
		return domain.LocalProductIdentity{}, err
	}
	identity, ok := f.identities[codprod]
	if !ok {
		return domain.LocalProductIdentity{}, ports.ErrProductNotFound
	}
	return identity, nil
}

func (f *fakeIdentityReader) LinkedListings(_ context.Context, codprod string) ([]string, error) {
	if err, ok := f.linkedErr[codprod]; ok {
		return nil, err
	}
	return f.linked[codprod], nil
}

type fakeCollector struct {
	searchResult ports.CatalogIdentityResult
	searchErr    error

	catalogProduct ports.CatalogProductDetail
	catalogErr     error

	offers    []domain.ValidatedOffer
	offersErr error

	ownPricing    map[string]ports.OwnItemPricing
	ownPricingErr map[string]error

	priceToWin    map[string]ports.PriceToWinSignal
	priceToWinErr map[string]error

	searchCalls     int
	catalogCalls    int
	offersCalls     int
	ownPricingCalls int
	priceToWinCalls int
}

func (f *fakeCollector) SearchCatalogByEAN(_ context.Context, _ string) (ports.CatalogIdentityResult, error) {
	f.searchCalls++
	return f.searchResult, f.searchErr
}

func (f *fakeCollector) GetCatalogProduct(_ context.Context, _ string) (ports.CatalogProductDetail, error) {
	f.catalogCalls++
	return f.catalogProduct, f.catalogErr
}

func (f *fakeCollector) ListCatalogOffers(_ context.Context, _ string) ([]domain.ValidatedOffer, error) {
	f.offersCalls++
	return f.offers, f.offersErr
}

func (f *fakeCollector) GetOwnItemPricing(_ context.Context, listingID string) (ports.OwnItemPricing, error) {
	f.ownPricingCalls++
	if err, ok := f.ownPricingErr[listingID]; ok {
		return ports.OwnItemPricing{}, err
	}
	return f.ownPricing[listingID], nil
}

func (f *fakeCollector) GetPriceToWin(_ context.Context, listingID string) (ports.PriceToWinSignal, error) {
	f.priceToWinCalls++
	if err, ok := f.priceToWinErr[listingID]; ok {
		return ports.PriceToWinSignal{}, err
	}
	return f.priceToWin[listingID], nil
}

func (f *fakeCollector) totalCalls() int {
	return f.searchCalls + f.catalogCalls + f.offersCalls + f.ownPricingCalls + f.priceToWinCalls
}

type fakeMatchDecisionStore struct {
	appended []domain.MatchDecision
}

func (f *fakeMatchDecisionStore) AppendMatchDecision(_ context.Context, decision domain.MatchDecision) error {
	f.appended = append(f.appended, decision)
	return nil
}

func (f *fakeMatchDecisionStore) LatestMatchDecisions(_ context.Context, _ []string) ([]domain.MatchDecision, error) {
	return f.appended, nil
}

type fakeSnapshotStore struct {
	appended []domain.MarketPriceSnapshot
}

func (f *fakeSnapshotStore) AppendSnapshot(_ context.Context, snapshot domain.MarketPriceSnapshot) error {
	f.appended = append(f.appended, snapshot)
	return nil
}

func (f *fakeSnapshotStore) LatestValidSnapshot(_ context.Context, _ domain.MarketPriceScope, _ string, _ domain.MarketPriceSource, _ time.Time) (*domain.MarketPriceSnapshot, error) {
	return nil, nil
}

func (f *fakeSnapshotStore) SnapshotSeries(_ context.Context, _ domain.MarketPriceScope, _ string, _ domain.MarketPriceSource, _ time.Time) ([]domain.MarketPriceSnapshot, error) {
	return nil, nil
}

type fakeOfferStore struct {
	appended []domain.ValidatedOffer
}

func (f *fakeOfferStore) AppendValidatedOffers(_ context.Context, offers []domain.ValidatedOffer) error {
	f.appended = append(f.appended, offers...)
	return nil
}

func (f *fakeOfferStore) ValidatedOffersByCatalogProduct(_ context.Context, _ string) ([]domain.ValidatedOffer, error) {
	return f.appended, nil
}

type fakeAggregateStore struct {
	appended []domain.MarketAggregate
}

func (f *fakeAggregateStore) AppendMarketAggregates(_ context.Context, aggregates []domain.MarketAggregate) error {
	f.appended = append(f.appended, aggregates...)
	return nil
}

func (f *fakeAggregateStore) LatestMarketAggregates(_ context.Context, _ []string) ([]domain.MarketAggregate, error) {
	return f.appended, nil
}

func (f *fakeAggregateStore) LatestMarketAggregatesBySource(_ context.Context, _ []string, _ domain.MarketPriceSource) ([]domain.MarketAggregate, error) {
	return f.appended, nil
}

type fakeCompetitiveSignalStore struct {
	appended []domain.CompetitiveSignal
}

func (f *fakeCompetitiveSignalStore) AppendCompetitiveSignals(_ context.Context, signals []domain.CompetitiveSignal) error {
	f.appended = append(f.appended, signals...)
	return nil
}

func (f *fakeCompetitiveSignalStore) LatestCompetitiveSignals(_ context.Context, _ []string) ([]domain.CompetitiveSignal, error) {
	return f.appended, nil
}

var (
	_ ports.ProductIdentityReader  = (*fakeIdentityReader)(nil)
	_ ports.PriceIntelCollector    = (*fakeCollector)(nil)
	_ ports.MatchDecisionStore     = (*fakeMatchDecisionStore)(nil)
	_ ports.SnapshotStore          = (*fakeSnapshotStore)(nil)
	_ ports.ValidatedOfferStore    = (*fakeOfferStore)(nil)
	_ ports.MarketAggregateStore   = (*fakeAggregateStore)(nil)
	_ ports.CompetitiveSignalStore = (*fakeCompetitiveSignalStore)(nil)
)

// --- test scaffolding ----------------------------------------------------

func fixedNow() time.Time {
	return time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
}

func stringPtr(v string) *string { return &v }

func money(amount string) *domain.Money {
	return &domain.Money{Amount: amount, Currency: "BRL"}
}

type pipelineFakes struct {
	identity   *fakeIdentityReader
	collector  *fakeCollector
	decisions  *fakeMatchDecisionStore
	snapshots  *fakeSnapshotStore
	offers     *fakeOfferStore
	aggregates *fakeAggregateStore
	signals    *fakeCompetitiveSignalStore
}

func newPipelineFakes() *pipelineFakes {
	return &pipelineFakes{
		identity:   &fakeIdentityReader{identities: map[string]domain.LocalProductIdentity{}, linked: map[string][]string{}},
		collector:  &fakeCollector{ownPricing: map[string]ports.OwnItemPricing{}, ownPricingErr: map[string]error{}, priceToWin: map[string]ports.PriceToWinSignal{}, priceToWinErr: map[string]error{}},
		decisions:  &fakeMatchDecisionStore{},
		snapshots:  &fakeSnapshotStore{},
		offers:     &fakeOfferStore{},
		aggregates: &fakeAggregateStore{},
		signals:    &fakeCompetitiveSignalStore{},
	}
}

func (fk *pipelineFakes) service() *CollectionPipelineService {
	return NewCollectionPipelineService(fk.identity, fk.collector, fk.decisions, fk.snapshots, fk.offers, fk.aggregates, fk.signals, fixedNow)
}

// acceptCandidate builds a local identity + catalog candidate pair that
// ResolveIdentity accepts (ean + marca anchors, no contradictions).
func acceptCandidate(codprod, catalogProductID string) (domain.LocalProductIdentity, domain.IdentityCandidate) {
	local := domain.LocalProductIdentity{
		Codprod: codprod, Descr: "Furadeira Acme 500W", EAN: stringPtr("7891234567890"), Marca: stringPtr("Acme"),
	}
	candidate := domain.IdentityCandidate{
		CatalogProductID: catalogProductID, Title: "Furadeira Acme 500W",
		Attrs: []domain.IdentityCandidateAttribute{
			{ID: "GTIN", ValueName: stringPtr("7891234567890")},
			{ID: "BRAND", ValueName: stringPtr("Acme")},
		},
	}
	return local, candidate
}

func qualifiedOffer(t *testing.T, catalogProductID, sellerID, amount string) domain.ValidatedOffer {
	t.Helper()
	offer, err := domain.NewValidatedOffer(domain.ValidatedOfferInput{
		CatalogProductID: catalogProductID, SellerID: sellerID, Price: money(amount),
		Condition: "new", MatchStatus: domain.MatchStatusAccept, FetchedAt: fixedNow(),
	})
	if err != nil {
		t.Fatalf("qualifiedOffer() error = %v", err)
	}
	return offer
}

// --- (a) no EAN ------------------------------------------------------------

func TestCollectNoEANSkipsCatalogEntirely(t *testing.T) {
	fk := newPipelineFakes()
	fk.identity.identities["P1"] = domain.LocalProductIdentity{Codprod: "P1", Descr: "no ean product"}
	fk.identity.linked["P1"] = []string{"listing-1"} // present to prove the collector is still never touched

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if got := fk.collector.totalCalls(); got != 0 {
		t.Fatalf("collector calls = %d, want 0", got)
	}
	if summary.Status != CollectionStatusCompleted {
		t.Fatalf("status = %v, want COMPLETED", summary.Status)
	}
	if len(summary.Decisoes) != 1 || summary.Decisoes[0].MatchStatus != domain.MatchStatusNoCandidate ||
		summary.Decisoes[0].PriceEvidence != domain.PriceEvidenceStatusNoPriceEvidence ||
		summary.Decisoes[0].Blocking == nil || *summary.Decisoes[0].Blocking != domain.BlockingStateNoCandidate {
		t.Fatalf("decisoes = %+v, want one NO_CANDIDATE/NO_PRICE_EVIDENCE/NO_CANDIDATE row", summary.Decisoes)
	}
	if summary.Contagens["no_candidate"] != 1 {
		t.Fatalf("contagens = %+v, want no_candidate=1", summary.Contagens)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCauseNoIdentity {
		t.Fatalf("causas = %v, want [NO_IDENTITY]", summary.Causas)
	}
	if summary.Causas[0].Detail != nil {
		t.Fatalf("NO_IDENTITY detail = %v, want nil (local resolution, no provider status)", *summary.Causas[0].Detail)
	}
	if len(fk.decisions.appended) != 1 || fk.decisions.appended[0].Status != domain.MatchStatusNoCandidate {
		t.Fatalf("persisted decisions = %+v, want one NO_CANDIDATE row", fk.decisions.appended)
	}
}

// --- (b) ACCEPT + >=5 sellers -> OK ----------------------------------------

func TestCollectAcceptWithFiveSellersProducesOKAggregate(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offers = []domain.ValidatedOffer{
		qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
		qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
		qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
		qualifiedOffer(t, "CAT1", "seller-4", "103.00"),
		qualifiedOffer(t, "CAT1", "seller-5", "104.00"),
	}
	fk.identity.linked["P1"] = []string{"listing-1"}
	fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: money("99.00"), WinnerPrice: money("98.00"), FetchedAt: fixedNow()}
	fk.collector.priceToWin["listing-1"] = ports.PriceToWinSignal{ListingID: "listing-1", TargetPrice: money("97.50"), Position: &domain.MarketPosition{Rank: 2, Total: 6}, FetchedAt: fixedNow()}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusCompleted {
		t.Fatalf("status = %v, want COMPLETED", summary.Status)
	}
	row := summary.Decisoes[0]
	if row.MatchStatus != domain.MatchStatusAccept || row.PriceEvidence != domain.PriceEvidenceStatusOK || row.Blocking != nil {
		t.Fatalf("decision row = %+v, want ACCEPT/OK/nil-blocking", row)
	}
	if summary.Contagens["ok"] != 1 {
		t.Fatalf("contagens = %+v, want ok=1", summary.Contagens)
	}
	if len(fk.offers.appended) != 5 {
		t.Fatalf("persisted offers = %d, want 5", len(fk.offers.appended))
	}
	if len(fk.aggregates.appended) != 1 || fk.aggregates.appended[0].Status != domain.MarketAggregateStatusOK || fk.aggregates.appended[0].NSellers != 5 {
		t.Fatalf("persisted aggregates = %+v, want one OK/NSellers=5 row", fk.aggregates.appended)
	}
	if len(fk.signals.appended) != 1 || fk.signals.appended[0].ListingID != "listing-1" {
		t.Fatalf("persisted competitive signals = %+v, want one listing-1 row", fk.signals.appended)
	}
	if len(fk.snapshots.appended) != 2 {
		t.Fatalf("persisted snapshots = %d, want 2 (ml_sale_price + ml_price_to_win)", len(fk.snapshots.appended))
	}
	for _, snapshot := range fk.snapshots.appended {
		if snapshot.Status != domain.MarketPriceSnapshotStatusValid {
			t.Fatalf("snapshot %+v, want VALID", snapshot)
		}
	}
}

// --- (c) ACCEPT + <5 sellers -> INSUFFICIENT_MARKET -------------------------

func TestCollectAcceptWithFewerThanFiveSellersProducesInsufficientMarket(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offers = []domain.ValidatedOffer{
		qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
		qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
		qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
	}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	row := summary.Decisoes[0]
	if row.PriceEvidence != domain.PriceEvidenceStatusInsufficientMarket || row.Blocking == nil || *row.Blocking != domain.BlockingStateInsufficientMarket {
		t.Fatalf("decision row = %+v, want INSUFFICIENT_MARKET", row)
	}
	if summary.Contagens["insufficient_market"] != 1 {
		t.Fatalf("contagens = %+v, want insufficient_market=1", summary.Contagens)
	}
}

// --- (d) ACCEPT + catalog offers disabled -> NO_PRICE_EVIDENCE + FLAG_DISABLED

func TestCollectAcceptWithCatalogOffersDisabledProducesNoPriceEvidence(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offersErr = ports.ErrCatalogOffersDisabled

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v, want nil (no 500)", err)
	}
	if summary.Status != CollectionStatusCompleted {
		t.Fatalf("status = %v, want COMPLETED (flag-disabled is not a provider failure)", summary.Status)
	}
	row := summary.Decisoes[0]
	if row.PriceEvidence != domain.PriceEvidenceStatusNoPriceEvidence {
		t.Fatalf("decision row = %+v, want NO_PRICE_EVIDENCE", row)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCauseFlagDisabled {
		t.Fatalf("causas = %v, want [FLAG_DISABLED]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "catalog offers capability disabled for this account" {
		t.Fatalf("FLAG_DISABLED detail = %v, want the capability-disabled note", summary.Causas[0].Detail)
	}
	if len(fk.aggregates.appended) != 1 || fk.aggregates.appended[0].Status != domain.MarketAggregateStatusNoPriceEvidence || fk.aggregates.appended[0].NOffers != 0 {
		t.Fatalf("persisted aggregates = %+v, want one NO_PRICE_EVIDENCE/n=0 row", fk.aggregates.appended)
	}
	if len(fk.offers.appended) != 0 {
		t.Fatalf("persisted offers = %d, want 0", len(fk.offers.appended))
	}
}

// --- (e) REVIEW -> no catalog-product/offers calls, blocking NO_CANDIDATE ---

func TestCollectReviewSkipsCatalogProductAndOffers(t *testing.T) {
	fk := newPipelineFakes()
	local := domain.LocalProductIdentity{Codprod: "P1", Descr: "Furadeira Acme 500W", EAN: stringPtr("7891234567890"), Marca: stringPtr("Acme")}
	candidate := domain.IdentityCandidate{
		CatalogProductID: "CAT1", Title: "Furadeira Acme 500W",
		Attrs: []domain.IdentityCandidateAttribute{{ID: "GTIN", ValueName: stringPtr("7891234567890")}}, // ean anchor only, no brand/model -> REVIEW
	}
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.identity.linked["P1"] = []string{"listing-1"}
	fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: money("50.00"), WinnerPrice: money("48.00"), FetchedAt: fixedNow()}
	fk.collector.priceToWin["listing-1"] = ports.PriceToWinSignal{ListingID: "listing-1", TargetPrice: money("47.50"), FetchedAt: fixedNow()}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if fk.collector.catalogCalls != 0 || fk.collector.offersCalls != 0 {
		t.Fatalf("catalogCalls=%d offersCalls=%d, want 0/0", fk.collector.catalogCalls, fk.collector.offersCalls)
	}
	row := summary.Decisoes[0]
	if row.MatchStatus != domain.MatchStatusReview {
		t.Fatalf("match status = %v, want REVIEW", row.MatchStatus)
	}
	if row.Blocking == nil || *row.Blocking != domain.BlockingStateNoCandidate {
		t.Fatalf("blocking = %v, want NO_CANDIDATE", row.Blocking)
	}
	if summary.Contagens["no_candidate"] != 1 {
		t.Fatalf("contagens = %+v, want no_candidate=1", summary.Contagens)
	}
	if fk.collector.ownPricingCalls != 1 || fk.collector.priceToWinCalls != 1 {
		t.Fatalf("competitive calls ownPricing=%d priceToWin=%d, want 1/1 (D-27 option X: signals fire on non-ACCEPT)", fk.collector.ownPricingCalls, fk.collector.priceToWinCalls)
	}
	if len(fk.signals.appended) != 1 || fk.signals.appended[0].ListingID != "listing-1" {
		t.Fatalf("persisted competitive signals = %+v, want one listing-1 row", fk.signals.appended)
	}
}

// --- (f) provider 5xx on offers -> PARTIAL + PROVIDER_5XX, latest-VALID untouched

func TestCollectProvider5xxOnOffersPreservesLatestValidAggregate(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offersErr = &ports.ProviderStatusError{StatusCode: 503}

	existingValid, err := domain.NewMarketAggregate(domain.MarketAggregateInput{
		ProductID: "P1", MinValid: money("90.00"), Median: money("95.00"), NOffers: 6, NSellers: 5,
		Source: domain.MarketPriceSourceMLCatalogOffers, FetchedAt: fixedNow().Add(-time.Hour), ComputedAt: fixedNow().Add(-time.Hour),
		Status: domain.MarketAggregateStatusOK,
	})
	if err != nil {
		t.Fatalf("seed aggregate error = %v", err)
	}
	fk.aggregates.appended = []domain.MarketAggregate{existingValid}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusPartial {
		t.Fatalf("status = %v, want PARTIAL", summary.Status)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCauseProvider5xx {
		t.Fatalf("causas = %v, want [PROVIDER_5XX]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "catalog offers fetch: provider responded HTTP 503" {
		t.Fatalf("PROVIDER_5XX detail = %v, want honest operation+status note", summary.Causas[0].Detail)
	}
	if len(fk.aggregates.appended) != 2 {
		t.Fatalf("persisted aggregates = %d, want 2 (untouched valid + new no-price-evidence)", len(fk.aggregates.appended))
	}
	if !reflect.DeepEqual(fk.aggregates.appended[0], existingValid) {
		t.Fatalf("existing valid aggregate mutated: got %+v, want %+v", fk.aggregates.appended[0], existingValid)
	}
	if fk.aggregates.appended[1].Status != domain.MarketAggregateStatusNoPriceEvidence {
		t.Fatalf("new aggregate = %+v, want NO_PRICE_EVIDENCE", fk.aggregates.appended[1])
	}
}

// --- (g) ErrRateLimited -> PARTIAL, remaining provider calls skipped -------

func TestCollectRateLimitedStopsRemainingProviderWork(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offersErr = ports.ErrRateLimited
	fk.identity.linked["P1"] = []string{"listing-1", "listing-2"}
	fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: money("50.00"), FetchedAt: fixedNow()}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusPartial {
		t.Fatalf("status = %v, want PARTIAL", summary.Status)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCauseProvider4xx {
		t.Fatalf("causas = %v, want [PROVIDER_4XX]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "catalog offers fetch: provider rate-limited (HTTP 429)" {
		t.Fatalf("PROVIDER_4XX detail = %v, want the rate-limited note", summary.Causas[0].Detail)
	}
	if fk.collector.ownPricingCalls != 0 || fk.collector.priceToWinCalls != 0 {
		t.Fatalf("ownPricingCalls=%d priceToWinCalls=%d, want 0/0 (stopped after rate limit)", fk.collector.ownPricingCalls, fk.collector.priceToWinCalls)
	}
	if len(fk.signals.appended) != 0 {
		t.Fatalf("persisted competitive signals = %d, want 0", len(fk.signals.appended))
	}
}

// --- (g2) GetOwnItemPricing nil-price -> FAILED PRICE_UNAVAILABLE snapshot --

func TestCollectOwnItemPricingNullPriceWritesFailedSnapshot(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offers = []domain.ValidatedOffer{
		qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
		qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
		qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
		qualifiedOffer(t, "CAT1", "seller-4", "103.00"),
		qualifiedOffer(t, "CAT1", "seller-5", "104.00"),
	}
	fk.identity.linked["P1"] = []string{"listing-1"}
	fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: nil, FetchedAt: fixedNow()}
	fk.collector.priceToWin["listing-1"] = ports.PriceToWinSignal{ListingID: "listing-1", TargetPrice: money("97.50"), FetchedAt: fixedNow()}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusPartial {
		t.Fatalf("status = %v, want PARTIAL", summary.Status)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCausePriceUnavailable {
		t.Fatalf("causas = %v, want [PRICE_UNAVAILABLE]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "own listing pricing: provider omitted our price" {
		t.Fatalf("PRICE_UNAVAILABLE detail = %v, want the omitted-our-price note", summary.Causas[0].Detail)
	}
	if fk.collector.priceToWinCalls != 1 {
		t.Fatalf("priceToWinCalls = %d, want 1 (null OurPrice must not suppress the price-to-win probe)", fk.collector.priceToWinCalls)
	}

	var saleFailed, saleValid int
	for _, snap := range fk.snapshots.appended {
		if snap.Source != domain.MarketPriceSourceMLSalePrice {
			continue
		}
		switch snap.Status {
		case domain.MarketPriceSnapshotStatusFailed:
			saleFailed++
			if snap.FailureReason == nil || *snap.FailureReason != string(CollectionCausePriceUnavailable) {
				t.Fatalf("failed sale-price snapshot reason = %v, want PRICE_UNAVAILABLE", snap.FailureReason)
			}
		case domain.MarketPriceSnapshotStatusValid:
			saleValid++
		}
	}
	if saleFailed != 1 {
		t.Fatalf("failed ml_sale_price snapshots = %d, want 1", saleFailed)
	}
	if saleValid != 0 {
		t.Fatalf("valid ml_sale_price snapshots = %d, want 0", saleValid)
	}
	if len(fk.signals.appended) != 0 {
		t.Fatalf("persisted competitive signals = %d, want 0 (no OurPrice, no signal)", len(fk.signals.appended))
	}
}

// --- (g3) GetPriceToWin nil-price -> FAILED PRICE_UNAVAILABLE snapshot -----

func TestCollectPriceToWinNullPriceWritesFailedSnapshot(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offers = []domain.ValidatedOffer{
		qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
		qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
		qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
		qualifiedOffer(t, "CAT1", "seller-4", "103.00"),
		qualifiedOffer(t, "CAT1", "seller-5", "104.00"),
	}
	fk.identity.linked["P1"] = []string{"listing-1"}
	fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: money("99.00"), WinnerPrice: money("98.00"), FetchedAt: fixedNow()}
	fk.collector.priceToWin["listing-1"] = ports.PriceToWinSignal{ListingID: "listing-1", TargetPrice: nil, FetchedAt: fixedNow()}

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusPartial {
		t.Fatalf("status = %v, want PARTIAL", summary.Status)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCausePriceUnavailable {
		t.Fatalf("causas = %v, want [PRICE_UNAVAILABLE]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "price-to-win: provider omitted the target price" {
		t.Fatalf("PRICE_UNAVAILABLE detail = %v, want the omitted-target-price note", summary.Causas[0].Detail)
	}

	var winFailed, winValid, saleValid int
	for _, snap := range fk.snapshots.appended {
		switch {
		case snap.Source == domain.MarketPriceSourceMLPriceToWin && snap.Status == domain.MarketPriceSnapshotStatusFailed:
			winFailed++
			if snap.FailureReason == nil || *snap.FailureReason != string(CollectionCausePriceUnavailable) {
				t.Fatalf("failed price-to-win snapshot reason = %v, want PRICE_UNAVAILABLE", snap.FailureReason)
			}
		case snap.Source == domain.MarketPriceSourceMLPriceToWin && snap.Status == domain.MarketPriceSnapshotStatusValid:
			winValid++
		case snap.Source == domain.MarketPriceSourceMLSalePrice && snap.Status == domain.MarketPriceSnapshotStatusValid:
			saleValid++
		}
	}
	if winFailed != 1 {
		t.Fatalf("failed ml_price_to_win snapshots = %d, want 1", winFailed)
	}
	if winValid != 0 {
		t.Fatalf("valid ml_price_to_win snapshots = %d, want 0", winValid)
	}
	if saleValid != 1 {
		t.Fatalf("valid ml_sale_price snapshots = %d, want 1 (sale-price path stayed green)", saleValid)
	}
	if len(fk.signals.appended) != 1 || fk.signals.appended[0].ListingID != "listing-1" {
		t.Fatalf("persisted competitive signals = %+v, want one listing-1 row (OurPrice present)", fk.signals.appended)
	}
}

// --- (g4) context.DeadlineExceeded -> TIMEOUT, run stops -------------------

func TestCollectDeadlineExceededStopsRunWithTimeoutCause(t *testing.T) {
	fk := newPipelineFakes()
	local, candidate := acceptCandidate("P1", "CAT1")
	fk.identity.identities["P1"] = local
	fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
	fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
	fk.collector.offers = []domain.ValidatedOffer{
		qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
		qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
		qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
		qualifiedOffer(t, "CAT1", "seller-4", "103.00"),
		qualifiedOffer(t, "CAT1", "seller-5", "104.00"),
	}
	fk.identity.linked["P1"] = []string{"listing-1", "listing-2"}
	fk.collector.ownPricingErr["listing-1"] = context.DeadlineExceeded

	summary, err := fk.service().Collect(context.Background(), "P1")
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if summary.Status != CollectionStatusPartial {
		t.Fatalf("status = %v, want PARTIAL", summary.Status)
	}
	if len(summary.Causas) != 1 || summary.Causas[0].Cause != CollectionCauseTimeout {
		t.Fatalf("causas = %v, want [TIMEOUT]", summary.Causas)
	}
	if summary.Causas[0].Detail == nil || *summary.Causas[0].Detail != "own listing pricing: provider request timed out" {
		t.Fatalf("TIMEOUT detail = %v, want the timeout note", summary.Causas[0].Detail)
	}
	if fk.collector.ownPricingCalls != 1 {
		t.Fatalf("ownPricingCalls = %d, want 1 (stopped after timeout, listing-2 not probed)", fk.collector.ownPricingCalls)
	}
	if fk.collector.priceToWinCalls != 0 {
		t.Fatalf("priceToWinCalls = %d, want 0 (stopped after timeout)", fk.collector.priceToWinCalls)
	}

	var failed int
	for _, snap := range fk.snapshots.appended {
		if snap.Source == domain.MarketPriceSourceMLSalePrice && snap.Status == domain.MarketPriceSnapshotStatusFailed {
			failed++
			if snap.FailureReason == nil || *snap.FailureReason != string(CollectionCauseTimeout) {
				t.Fatalf("failed snapshot reason = %v, want TIMEOUT", snap.FailureReason)
			}
		}
	}
	if failed != 1 {
		t.Fatalf("failed ml_sale_price snapshots = %d, want 1", failed)
	}
}

// --- (h) concurrent Collect on same codprod -> ErrCollectionInProgress ------

func TestCollectConcurrentSameCodprodReturnsInProgress(t *testing.T) {
	fk := newPipelineFakes()
	fk.identity.identities["P1"] = domain.LocalProductIdentity{Codprod: "P1", Descr: "no ean product"}
	service := fk.service()
	service.inflight.Store("P1", struct{}{})

	_, err := service.Collect(context.Background(), "P1")
	if !errors.Is(err, ErrCollectionInProgress) {
		t.Fatalf("Collect() error = %v, want ErrCollectionInProgress", err)
	}
	if fk.collector.totalCalls() != 0 {
		t.Fatalf("collector calls = %d, want 0", fk.collector.totalCalls())
	}
}

// --- (i) ErrProductNotFound propagates --------------------------------------

func TestCollectProductNotFoundPropagates(t *testing.T) {
	fk := newPipelineFakes()

	_, err := fk.service().Collect(context.Background(), "UNKNOWN")
	if !errors.Is(err, ports.ErrProductNotFound) {
		t.Fatalf("Collect() error = %v, want ErrProductNotFound", err)
	}
}

// --- (j) determinism ---------------------------------------------------------

func TestCollectIsDeterministicAcrossRuns(t *testing.T) {
	run := func() (CollectionSummary, *pipelineFakes) {
		fk := newPipelineFakes()
		local, candidate := acceptCandidate("P1", "CAT1")
		fk.identity.identities["P1"] = local
		fk.collector.searchResult = ports.CatalogIdentityResult{Candidates: []domain.IdentityCandidate{candidate}, FetchedAt: fixedNow()}
		fk.collector.catalogProduct = ports.CatalogProductDetail{CatalogProductID: "CAT1", FetchedAt: fixedNow()}
		fk.collector.offers = []domain.ValidatedOffer{
			qualifiedOffer(t, "CAT1", "seller-1", "100.00"),
			qualifiedOffer(t, "CAT1", "seller-2", "101.00"),
			qualifiedOffer(t, "CAT1", "seller-3", "102.00"),
			qualifiedOffer(t, "CAT1", "seller-4", "103.00"),
			qualifiedOffer(t, "CAT1", "seller-5", "104.00"),
		}
		fk.identity.linked["P1"] = []string{"listing-1"}
		fk.collector.ownPricing["listing-1"] = ports.OwnItemPricing{ListingID: "listing-1", OurPrice: money("99.00"), FetchedAt: fixedNow()}
		fk.collector.priceToWin["listing-1"] = ports.PriceToWinSignal{ListingID: "listing-1", TargetPrice: money("97.50"), FetchedAt: fixedNow()}

		summary, err := fk.service().Collect(context.Background(), "P1")
		if err != nil {
			t.Fatalf("Collect() error = %v", err)
		}
		return summary, fk
	}

	summary1, fk1 := run()
	summary2, fk2 := run()

	if !reflect.DeepEqual(summary1, summary2) {
		t.Fatalf("summaries differ:\n%+v\n%+v", summary1, summary2)
	}
	if !reflect.DeepEqual(fk1.decisions.appended, fk2.decisions.appended) {
		t.Fatalf("persisted decisions differ")
	}
	if !reflect.DeepEqual(fk1.offers.appended, fk2.offers.appended) {
		t.Fatalf("persisted offers differ")
	}
	if !reflect.DeepEqual(fk1.aggregates.appended, fk2.aggregates.appended) {
		t.Fatalf("persisted aggregates differ")
	}
	if !reflect.DeepEqual(fk1.snapshots.appended, fk2.snapshots.appended) {
		t.Fatalf("persisted snapshots differ")
	}
	if !reflect.DeepEqual(fk1.signals.appended, fk2.signals.appended) {
		t.Fatalf("persisted signals differ")
	}
}
