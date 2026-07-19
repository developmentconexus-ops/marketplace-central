package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

// ErrCollectionInProgress signals a concurrent Collect call for the same
// codprod; IC-03 forbids a job table, so concurrency is guarded in-process.
var ErrCollectionInProgress = errors.New("COLLECTION_IN_PROGRESS")

const snapshotTTL = time.Hour

type CollectionStatus string

const (
	CollectionStatusCompleted CollectionStatus = "COMPLETED"
	CollectionStatusPartial   CollectionStatus = "PARTIAL"
)

type CollectionCause string

const (
	CollectionCauseNoIdentity       CollectionCause = "NO_IDENTITY"
	CollectionCauseFlagDisabled     CollectionCause = "FLAG_DISABLED"
	CollectionCauseProvider4xx      CollectionCause = "PROVIDER_4XX"
	CollectionCauseProvider5xx      CollectionCause = "PROVIDER_5XX"
	CollectionCauseTimeout          CollectionCause = "TIMEOUT"
	CollectionCausePriceUnavailable CollectionCause = "PRICE_UNAVAILABLE"
	// CollectionCauseInvalidListingRef is a LOCAL failure to resolve a listing
	// id into account context — distinct from any provider fault so an
	// unparseable id is never mislabeled PROVIDER_4XX (ADR-17; live-drive
	// reprova, D-86).
	CollectionCauseInvalidListingRef CollectionCause = "INVALID_LISTING_REF"
)

// CollectionDecision is one collection outcome row for a codprod. Blocking
// mirrors domain.Verdict's own nilable BlockingState: nil means nothing
// blocks (price evidence is OK).
type CollectionDecision struct {
	Codprod       string
	MatchStatus   domain.MatchStatus
	PriceEvidence domain.PriceEvidenceStatus
	Blocking      *domain.BlockingState
}

// CollectionCauseDetail pairs a PARTIAL/skip cause with an optional honest,
// MPC-native observability note. Detail names the pipeline operation that
// failed and — when the provider surfaced one — its HTTP status; it never
// carries a raw provider payload, URL, PII, or token. Detail is nil when the
// cause has no extra provider-status context (e.g. a local NO_IDENTITY).
type CollectionCauseDetail struct {
	Cause  CollectionCause
	Detail *string
}

// CollectionSummary is the synchronous result of one Collect call. Contagens
// keys mirror amendment B.1's lower snake_case categories; sem_custo always
// stays 0 here because cost is a verdict-read (S4) concern, never consulted
// at collection time (D-F4-i).
type CollectionSummary struct {
	Codprod   string
	Status    CollectionStatus
	Decisoes  []CollectionDecision
	Contagens map[string]int
	Causas    []CollectionCauseDetail
}

// CollectionPipelineService orchestrates the on-demand collection pipeline
// for one codprod: identity -> catalog -> snapshots -> offers -> aggregate ->
// competitive signals, persisting evidence via the F-02 store ports.
type CollectionPipelineService struct {
	identity   ports.ProductIdentityReader
	collector  ports.PriceIntelCollector
	decisions  ports.MatchDecisionStore
	snapshots  ports.SnapshotStore
	offers     ports.ValidatedOfferStore
	aggregates ports.MarketAggregateStore
	signals    ports.CompetitiveSignalStore
	now        func() time.Time
	inflight   sync.Map
}

func NewCollectionPipelineService(
	identity ports.ProductIdentityReader,
	collector ports.PriceIntelCollector,
	decisions ports.MatchDecisionStore,
	snapshots ports.SnapshotStore,
	offers ports.ValidatedOfferStore,
	aggregates ports.MarketAggregateStore,
	signals ports.CompetitiveSignalStore,
	now func() time.Time,
) *CollectionPipelineService {
	if now == nil {
		now = time.Now
	}
	return &CollectionPipelineService{
		identity:   identity,
		collector:  collector,
		decisions:  decisions,
		snapshots:  snapshots,
		offers:     offers,
		aggregates: aggregates,
		signals:    signals,
		now:        now,
	}
}

// collectionRun accumulates the mutable outcome of one Collect invocation:
// causes observed so far, whether any provider failure occurred (PARTIAL),
// and whether remaining provider work must stop (rate-limited).
type collectionRun struct {
	codprod string
	causas  []CollectionCauseDetail
	partial bool
	stopped bool
}

func (r *collectionRun) fail(cause CollectionCause, stop bool, detail *string) {
	r.causas = append(r.causas, CollectionCauseDetail{Cause: cause, Detail: detail})
	r.partial = true
	if stop {
		r.stopped = true
	}
}

func (s *CollectionPipelineService) Collect(ctx context.Context, codprod string) (CollectionSummary, error) {
	if _, loaded := s.inflight.LoadOrStore(codprod, struct{}{}); loaded {
		return CollectionSummary{}, ErrCollectionInProgress
	}
	defer s.inflight.Delete(codprod)

	local, err := s.identity.GetLocalIdentity(ctx, codprod)
	if err != nil {
		return CollectionSummary{}, err
	}

	run := &collectionRun{codprod: codprod}

	if local.EAN == nil || strings.TrimSpace(*local.EAN) == "" {
		decision, err := domain.ResolveIdentity(local, nil, s.now())
		if err != nil {
			return CollectionSummary{}, err
		}
		if err := s.decisions.AppendMatchDecision(ctx, decision); err != nil {
			return CollectionSummary{}, err
		}
		run.causas = append(run.causas, CollectionCauseDetail{Cause: CollectionCauseNoIdentity})
		return buildSummary(run, decision.Status, domain.PriceEvidenceStatusNoPriceEvidence, blockingPtr(domain.BlockingStateNoCandidate)), nil
	}

	matchStatus := domain.MatchStatusNoCandidate
	priceEvidence := domain.PriceEvidenceStatusNoPriceEvidence
	blocking := blockingPtr(domain.BlockingStateNoCandidate)

	catalogResult, err := s.collector.SearchCatalogByEAN(ctx, *local.EAN)
	if err != nil {
		cause, stop := classifyProviderFailure(err)
		run.fail(cause, stop, providerCauseDetail("catalog EAN search", err))
	} else {
		decision, err := domain.ResolveIdentity(local, catalogResult.Candidates, s.now())
		if err != nil {
			return CollectionSummary{}, err
		}
		if err := s.decisions.AppendMatchDecision(ctx, decision); err != nil {
			return CollectionSummary{}, err
		}
		matchStatus = decision.Status
		if decision.Status == domain.MatchStatusAccept {
			priceEvidence, blocking, err = s.collectCatalogEvidence(ctx, run, codprod, *decision.CatalogProductID)
			if err != nil {
				return CollectionSummary{}, err
			}
		}
	}

	if !run.stopped {
		if err := s.collectCompetitiveSignals(ctx, run, codprod); err != nil {
			return CollectionSummary{}, err
		}
	}

	return buildSummary(run, matchStatus, priceEvidence, blocking), nil
}

func (s *CollectionPipelineService) collectCatalogEvidence(ctx context.Context, run *collectionRun, codprod, catalogProductID string) (domain.PriceEvidenceStatus, *domain.BlockingState, error) {
	catalogProduct, err := s.collector.GetCatalogProduct(ctx, catalogProductID)
	if err != nil {
		cause, stop := classifyProviderFailure(err)
		run.fail(cause, stop, providerCauseDetail("catalog product fetch", err))
		if err := s.persistNoPriceEvidenceAggregate(ctx, codprod, s.now()); err != nil {
			return "", nil, err
		}
		return domain.PriceEvidenceStatusNoPriceEvidence, blockingPtr(domain.BlockingStateNoPriceEvidence), nil
	}

	offers, err := s.collector.ListCatalogOffers(ctx, catalogProductID)
	switch {
	case errors.Is(err, ports.ErrCatalogOffersDisabled):
		run.causas = append(run.causas, CollectionCauseDetail{Cause: CollectionCauseFlagDisabled, Detail: detailPtr("catalog offers capability disabled for this account")})
		if err := s.persistNoPriceEvidenceAggregate(ctx, codprod, catalogProduct.FetchedAt); err != nil {
			return "", nil, err
		}
		return domain.PriceEvidenceStatusNoPriceEvidence, blockingPtr(domain.BlockingStateNoPriceEvidence), nil
	case err != nil:
		cause, stop := classifyProviderFailure(err)
		run.fail(cause, stop, providerCauseDetail("catalog offers fetch", err))
		if err := s.persistNoPriceEvidenceAggregate(ctx, codprod, catalogProduct.FetchedAt); err != nil {
			return "", nil, err
		}
		return domain.PriceEvidenceStatusNoPriceEvidence, blockingPtr(domain.BlockingStateNoPriceEvidence), nil
	}

	if err := s.offers.AppendValidatedOffers(ctx, offers); err != nil {
		return "", nil, err
	}

	// Competitor aggregate EXCLUDES our own offer: ML returns our own listing
	// among the catalog offers, and counting it would corrupt the competitor
	// min/median/max/n_offers/n_sellers the operator prices against. The raw
	// offers above keep the own offer (honest evidence); only the aggregate is
	// filtered. Own seller id is resolved dynamically per installation.
	ownSellerID, err := s.collector.OwnSellerID(ctx)
	if err != nil {
		return "", nil, err
	}
	competitorOffers := excludeOwnSeller(offers, ownSellerID)

	aggregate, err := domain.ComputeMarketAggregate(codprod, competitorOffers, domain.MarketPriceSourceMLCatalogOffers, catalogProduct.FetchedAt, s.now())
	if err != nil {
		return "", nil, err
	}
	if err := s.aggregates.AppendMarketAggregates(ctx, []domain.MarketAggregate{aggregate}); err != nil {
		return "", nil, err
	}
	return aggregateToPriceEvidence(aggregate.Status), aggregateToBlocking(aggregate.Status), nil
}

// excludeOwnSeller drops offers whose seller id matches the installation's own
// seller id so the competitor aggregate never counts our own listing. An empty
// ownSellerID (not resolvable) excludes nothing — we never guess which offer is
// ours (ADR-17: an unknown own-seller id must not silently reshape the market).
func excludeOwnSeller(offers []domain.ValidatedOffer, ownSellerID string) []domain.ValidatedOffer {
	own := strings.TrimSpace(ownSellerID)
	if own == "" {
		return offers
	}
	filtered := make([]domain.ValidatedOffer, 0, len(offers))
	for _, offer := range offers {
		// Trim both sides so a stray-whitespace id (a seed/fixture that skips
		// the repo's write-time trim) can never silently defeat the exclusion
		// and leak our own offer back into the competitor stats.
		if strings.TrimSpace(offer.SellerID) == own {
			continue
		}
		filtered = append(filtered, offer)
	}
	return filtered
}

func (s *CollectionPipelineService) collectCompetitiveSignals(ctx context.Context, run *collectionRun, codprod string) error {
	listingIDs, err := s.identity.LinkedListings(ctx, codprod)
	if err != nil {
		return err
	}

	signals := make([]domain.CompetitiveSignal, 0, len(listingIDs))
	for _, listingID := range listingIDs {
		if run.stopped {
			break
		}

		pricing, err := s.collector.GetOwnItemPricing(ctx, listingID)
		if err != nil {
			cause, stop := classifyProviderFailure(err)
			run.fail(cause, stop, providerCauseDetail("own listing pricing", err))
			if err := s.appendFailedSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLSalePrice, cause); err != nil {
				return err
			}
			continue
		}
		if pricing.OurPrice == nil {
			run.fail(CollectionCausePriceUnavailable, false, detailPtr("own listing pricing: provider omitted our price"))
			if err := s.appendFailedSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLSalePrice, CollectionCausePriceUnavailable); err != nil {
				return err
			}
		} else if err := s.appendValidSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLSalePrice, pricing.OurPrice, pricing.FetchedAt); err != nil {
			return err
		}

		var targetPrice *domain.Money
		var position *domain.MarketPosition
		if !run.stopped {
			signal, err := s.collector.GetPriceToWin(ctx, listingID)
			switch {
			case err != nil:
				cause, stop := classifyProviderFailure(err)
				run.fail(cause, stop, providerCauseDetail("price-to-win", err))
				if err := s.appendFailedSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLPriceToWin, cause); err != nil {
					return err
				}
			case signal.TargetPrice == nil:
				run.fail(CollectionCausePriceUnavailable, false, detailPtr("price-to-win: provider omitted the target price"))
				if err := s.appendFailedSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLPriceToWin, CollectionCausePriceUnavailable); err != nil {
					return err
				}
			default:
				targetPrice = signal.TargetPrice
				position = signal.Position
				if err := s.appendValidSnapshot(ctx, codprod, listingID, domain.MarketPriceSourceMLPriceToWin, signal.TargetPrice, signal.FetchedAt); err != nil {
					return err
				}
			}
		}

		if pricing.OurPrice == nil {
			continue
		}
		competitive, err := domain.NewCompetitiveSignal(domain.CompetitiveSignalInput{
			ListingID:   listingID,
			OurPrice:    pricing.OurPrice,
			WinnerPrice: pricing.WinnerPrice,
			TargetPrice: targetPrice,
			Position:    position,
			FetchedAt:   s.now(),
		})
		if err != nil {
			return err
		}
		signals = append(signals, competitive)
	}

	if len(signals) == 0 {
		return nil
	}
	return s.signals.AppendCompetitiveSignals(ctx, signals)
}

func (s *CollectionPipelineService) persistNoPriceEvidenceAggregate(ctx context.Context, codprod string, fetchedAt time.Time) error {
	aggregate, err := domain.NewMarketAggregate(domain.MarketAggregateInput{
		ProductID:  codprod,
		NOffers:    0,
		NSellers:   0,
		Source:     domain.MarketPriceSourceMLCatalogOffers,
		FetchedAt:  fetchedAt,
		ComputedAt: s.now(),
		Status:     domain.MarketAggregateStatusNoPriceEvidence,
	})
	if err != nil {
		return err
	}
	return s.aggregates.AppendMarketAggregates(ctx, []domain.MarketAggregate{aggregate})
}

func (s *CollectionPipelineService) appendValidSnapshot(ctx context.Context, codprod, refID string, source domain.MarketPriceSource, price *domain.Money, fetchedAt time.Time) error {
	snapshot, err := domain.NewMarketPriceSnapshot(domain.MarketPriceSnapshotInput{
		Scope:      domain.MarketPriceScopeListing,
		RefID:      refID,
		Source:     source,
		Price:      price,
		ObservedAt: fetchedAt,
		FetchedAt:  fetchedAt,
		ExpiresAt:  fetchedAt.Add(snapshotTTL),
		Status:     domain.MarketPriceSnapshotStatusValid,
		RequestID:  requestID(codprod, refID, source),
	})
	if err != nil {
		return err
	}
	return s.snapshots.AppendSnapshot(ctx, snapshot)
}

func (s *CollectionPipelineService) appendFailedSnapshot(ctx context.Context, codprod, refID string, source domain.MarketPriceSource, cause CollectionCause) error {
	reason := string(cause)
	observedAt := s.now()
	snapshot, err := domain.NewMarketPriceSnapshot(domain.MarketPriceSnapshotInput{
		Scope:         domain.MarketPriceScopeListing,
		RefID:         refID,
		Source:        source,
		ObservedAt:    observedAt,
		FetchedAt:     observedAt,
		ExpiresAt:     observedAt.Add(snapshotTTL),
		Status:        domain.MarketPriceSnapshotStatusFailed,
		FailureReason: &reason,
		RequestID:     requestID(codprod, refID, source),
	})
	if err != nil {
		return err
	}
	return s.snapshots.AppendSnapshot(ctx, snapshot)
}

func detailPtr(s string) *string { return &s }

// providerCauseDetail renders an honest, MPC-native observability note for a
// provider failure: the pipeline operation that failed plus the provider's
// HTTP status when it surfaced one. It deliberately excludes response bodies,
// URLs, PII, and tokens — only the operation label and status class leave the
// adapter boundary.
func providerCauseDetail(operation string, err error) *string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return detailPtr(operation + ": provider request timed out")
	case errors.Is(err, ports.ErrRateLimited):
		return detailPtr(operation + ": provider rate-limited (HTTP 429)")
	case errors.Is(err, ports.ErrInvalidListingRef):
		return detailPtr(operation + ": unresolvable listing reference (local)")
	default:
		var statusErr *ports.ProviderStatusError
		if errors.As(err, &statusErr) {
			return detailPtr(fmt.Sprintf("%s: provider responded HTTP %d", operation, statusErr.StatusCode))
		}
		return detailPtr(operation + ": provider request failed")
	}
}

func classifyProviderFailure(err error) (CollectionCause, bool) {
	if errors.Is(err, context.DeadlineExceeded) {
		return CollectionCauseTimeout, true
	}
	if errors.Is(err, ports.ErrRateLimited) {
		return CollectionCauseProvider4xx, true
	}
	if errors.Is(err, ports.ErrInvalidListingRef) {
		// Local resolution failure for this listing — not a provider fault, and
		// not fatal to the other listings under the codprod (stop=false).
		return CollectionCauseInvalidListingRef, false
	}
	var statusErr *ports.ProviderStatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode >= 500 {
		return CollectionCauseProvider5xx, false
	}
	return CollectionCauseProvider4xx, false
}

func aggregateToPriceEvidence(status domain.MarketAggregateStatus) domain.PriceEvidenceStatus {
	switch status {
	case domain.MarketAggregateStatusOK:
		return domain.PriceEvidenceStatusOK
	case domain.MarketAggregateStatusInsufficientMarket:
		return domain.PriceEvidenceStatusInsufficientMarket
	default:
		return domain.PriceEvidenceStatusNoPriceEvidence
	}
}

func aggregateToBlocking(status domain.MarketAggregateStatus) *domain.BlockingState {
	switch status {
	case domain.MarketAggregateStatusOK:
		return nil
	case domain.MarketAggregateStatusInsufficientMarket:
		return blockingPtr(domain.BlockingStateInsufficientMarket)
	default:
		return blockingPtr(domain.BlockingStateNoPriceEvidence)
	}
}

func blockingPtr(state domain.BlockingState) *domain.BlockingState { return &state }

func requestID(codprod, refID string, source domain.MarketPriceSource) string {
	return fmt.Sprintf("collect-%s-%s-%s", codprod, refID, source)
}

func buildSummary(run *collectionRun, matchStatus domain.MatchStatus, priceEvidence domain.PriceEvidenceStatus, blocking *domain.BlockingState) CollectionSummary {
	status := CollectionStatusCompleted
	if run.partial {
		status = CollectionStatusPartial
	}
	decision := CollectionDecision{
		Codprod:       run.codprod,
		MatchStatus:   matchStatus,
		PriceEvidence: priceEvidence,
		Blocking:      blocking,
	}
	return CollectionSummary{
		Codprod:   run.codprod,
		Status:    status,
		Decisoes:  []CollectionDecision{decision},
		Contagens: contagens(decision),
		Causas:    run.causas,
	}
}

func contagens(decision CollectionDecision) map[string]int {
	counts := map[string]int{
		"ok":                  0,
		"no_price_evidence":   0,
		"insufficient_market": 0,
		"no_candidate":        0,
		"sem_custo":           0,
	}
	counts[contagemKey(decision)]++
	return counts
}

func contagemKey(decision CollectionDecision) string {
	if decision.Blocking == nil {
		return "ok"
	}
	switch *decision.Blocking {
	case domain.BlockingStateNoCandidate:
		return "no_candidate"
	case domain.BlockingStateInsufficientMarket:
		return "insufficient_market"
	case domain.BlockingStateSemCusto:
		return "sem_custo"
	default:
		return "no_price_evidence"
	}
}
