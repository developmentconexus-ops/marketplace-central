package application

import (
	"context"
	"strconv"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

// EvidenceReadService implements ports.EvidenceReader: the batch GET side of
// price-intel evidence. Every method honors the shared read contract — batch
// cap (MaxReadIDs), input order preserved, and honest-empty synthesis for
// IDs without stored evidence (ADR-17: unknown never becomes zero/default).
type EvidenceReadService struct {
	decisions  ports.MatchDecisionStore
	aggregates ports.MarketAggregateStore
	signals    ports.CompetitiveSignalStore
	cost       ports.CostReader
	now        func() time.Time
}

var _ ports.EvidenceReader = (*EvidenceReadService)(nil)

func NewEvidenceReadService(
	decisions ports.MatchDecisionStore,
	aggregates ports.MarketAggregateStore,
	signals ports.CompetitiveSignalStore,
	cost ports.CostReader,
	now func() time.Time,
) *EvidenceReadService {
	if now == nil {
		now = time.Now
	}
	return &EvidenceReadService{
		decisions:  decisions,
		aggregates: aggregates,
		signals:    signals,
		cost:       cost,
		now:        now,
	}
}

// Signals returns the latest competitive signal per requested listing. A
// listing with no stored row gets an honest empty CompetitiveSignal (nil
// Money fields, zero FetchedAt) rather than a fabricated price.
func (s *EvidenceReadService) Signals(ctx context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error) {
	if len(listingIDs) > MaxReadIDs {
		return nil, &TooManyIDsError{Count: len(listingIDs), Limit: MaxReadIDs}
	}
	items := make([]domain.CompetitiveSignal, 0, len(listingIDs))
	if len(listingIDs) == 0 {
		return items, nil
	}

	rows, err := s.signals.LatestCompetitiveSignals(ctx, listingIDs)
	if err != nil {
		return nil, err
	}
	byListingID := make(map[string]domain.CompetitiveSignal, len(rows))
	for _, row := range rows {
		byListingID[row.ListingID] = row
	}
	for _, listingID := range listingIDs {
		if row, ok := byListingID[listingID]; ok {
			items = append(items, row)
			continue
		}
		items = append(items, domain.CompetitiveSignal{ListingID: listingID})
	}
	return items, nil
}

// Aggregates returns the latest product-keyed market aggregate per requested
// codprod. A codprod with no stored row gets an honest empty
// MarketAggregate (NO_PRICE_EVIDENCE, n=0) rather than an invented sample.
func (s *EvidenceReadService) Aggregates(ctx context.Context, codprods []string) ([]domain.MarketAggregate, error) {
	if len(codprods) > MaxReadIDs {
		return nil, &TooManyIDsError{Count: len(codprods), Limit: MaxReadIDs}
	}
	items := make([]domain.MarketAggregate, 0, len(codprods))
	if len(codprods) == 0 {
		return items, nil
	}

	rows, err := s.aggregates.LatestMarketAggregates(ctx, codprods)
	if err != nil {
		return nil, err
	}
	byProductID := make(map[string]domain.MarketAggregate, len(rows))
	for _, row := range rows {
		byProductID[row.ProductID] = row
	}
	for _, codprod := range codprods {
		if row, ok := byProductID[codprod]; ok {
			items = append(items, row)
			continue
		}
		items = append(items, domain.MarketAggregate{ProductID: codprod, Status: domain.MarketAggregateStatusNoPriceEvidence})
	}
	return items, nil
}

// Verdicts assembles a domain.Verdict per requested codprod via
// domain.AssembleVerdict. A never-collected codprod still yields a 200-shaped
// honest verdict (NO_CANDIDATE / NO_PRICE_EVIDENCE), never an error.
func (s *EvidenceReadService) Verdicts(ctx context.Context, codprods []string) ([]domain.Verdict, error) {
	if len(codprods) > MaxReadIDs {
		return nil, &TooManyIDsError{Count: len(codprods), Limit: MaxReadIDs}
	}
	verdicts := make([]domain.Verdict, 0, len(codprods))
	if len(codprods) == 0 {
		return verdicts, nil
	}

	decisionRows, err := s.decisions.LatestMatchDecisions(ctx, codprods)
	if err != nil {
		return nil, err
	}
	byDecisionProductID := make(map[string]domain.MatchDecision, len(decisionRows))
	for _, row := range decisionRows {
		byDecisionProductID[row.ProductID] = row
	}

	aggregateRows, err := s.aggregates.LatestMarketAggregatesBySource(ctx, codprods, domain.MarketPriceSourceMLCatalogOffers)
	if err != nil {
		return nil, err
	}
	byAggregateProductID := make(map[string]domain.MarketAggregate, len(aggregateRows))
	for _, row := range aggregateRows {
		byAggregateProductID[row.ProductID] = row
	}

	for _, codprod := range codprods {
		decision, ok := byDecisionProductID[codprod]
		if !ok {
			decision = domain.MatchDecision{ProductID: codprod, Status: domain.MatchStatusNoCandidate}
		}

		var aggregate *domain.MarketAggregate
		if row, ok := byAggregateProductID[codprod]; ok && row.Status != domain.MarketAggregateStatusNoPriceEvidence {
			aggregateCopy := row
			aggregate = &aggregateCopy
		}

		cost, err := s.costRef(ctx, codprod)
		if err != nil {
			return nil, err
		}

		verdict, err := domain.AssembleVerdict(domain.VerdictInput{
			ProductID: codprod,
			Decision:  decision,
			Aggregate: aggregate,
			Cost:      cost,
		})
		if err != nil {
			return nil, err
		}
		verdicts = append(verdicts, verdict)
	}
	return verdicts, nil
}

// costRef resolves the ERP cost reference for one codprod. A non-numeric
// codprod cannot address the cost port's int key, so it is treated as an
// absent cost (nil), not a batch failure (D-F4-m). The port carries no
// not-found sentinel, so any non-nil error from the reader is a hard error
// and propagates; only a nil error with a nil Money means "absent".
func (s *EvidenceReadService) costRef(ctx context.Context, codprod string) (*domain.CostRef, error) {
	productID, err := strconv.Atoi(codprod)
	if err != nil {
		return nil, nil
	}
	money, asOf, err := s.cost.GetCostAsOf(ctx, productID, s.now())
	if err != nil {
		return nil, err
	}
	if money == nil {
		return nil, nil
	}
	return &domain.CostRef{AsOf: asOf}, nil
}
