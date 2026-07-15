package application

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

const originUFMG int64 = 13

// maxBelowMarginScanPages limits one request to 50 keyset fetches and 50 cost
// batches. This covers up to 50×limit candidates while bounding Oracle load;
// clients resume from the last scanned row when the cap is reached.
const maxBelowMarginScanPages = 50

var ErrInstallationIDRequired = errors.New("installation_id is required")
var ErrInstallationNotFound = errors.New("installation not found")

type ReadService struct {
	repo          ports.ListingReadRepository
	facts         ports.CostReader
	policies      ports.PolicyReader
	installations ports.InstallationReader
	now           func() time.Time
}

func NewReadService(repo ports.ListingReadRepository, facts ports.CostReader, policies ports.PolicyReader, installations ports.InstallationReader, now func() time.Time) ReadService {
	return ReadService{repo: repo, facts: facts, policies: policies, installations: installations, now: now}
}

func (s ReadService) List(ctx context.Context, query ports.ListingQuery) (ports.ListingRowPage, error) {
	if strings.TrimSpace(query.InstallationID) == "" {
		return ports.ListingRowPage{}, ErrInstallationIDRequired
	}
	found, err := s.installations.InstallationExists(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("validate installation: %w", err)
	}
	if !found {
		return ports.ListingRowPage{}, ErrInstallationNotFound
	}
	policy, policyFound, err := s.policies.GetPricingPolicyForInstallation(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("read pricing policy: %w", err)
	}
	ceilings, err := s.facts.GetICMSCeilingByOrigin(ctx, originUFMG)
	if err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("read ICMS ceiling: %w", err)
	}
	maxCeiling := maximumCeiling(ceilings)
	serveTime := s.now().UTC()
	if !needsBelowMarginScan(query.Filter) {
		page, err := s.repo.ListListingRows(ctx, query)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		if err := s.enrich(ctx, page.Items, maxCeiling, policy, policyFound); err != nil {
			return ports.ListingRowPage{}, err
		}
		page.AsOf = serveTime
		return page, nil
	}
	return s.scan(ctx, query, maxCeiling, policy, policyFound, serveTime)
}

func (s ReadService) scan(ctx context.Context, q ports.ListingQuery, ceiling *float64, policy ports.PricingPolicy, policyFound bool, asOf time.Time) (ports.ListingRowPage, error) {
	result := ports.ListingRowPage{Items: make([]domain.ListingReadModel, 0, q.Limit), AsOf: asOf}
	cursor := q.Cursor
	for pageNo := 0; pageNo < maxBelowMarginScanPages; pageNo++ {
		candidate := q
		candidate.Cursor = cursor
		page, err := s.repo.ListListingRows(ctx, candidate)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		if err := s.enrich(ctx, page.Items, ceiling, policy, policyFound); err != nil {
			return ports.ListingRowPage{}, err
		}
		for _, item := range page.Items {
			lastScanned := ports.ListingCursor{LastTitle: item.Title, ListingID: item.ListingID}
			result.NextCursor = &lastScanned
			if matchesDependentFilter(item, q.Filter) {
				result.Items = append(result.Items, item)
				if len(result.Items) == q.Limit {
					if page.NextCursor == nil && item.ListingID == page.Items[len(page.Items)-1].ListingID {
						result.NextCursor = nil
					}
					return result, nil
				}
			}
		}
		if page.NextCursor == nil {
			// Source keyset exhausted before reaching the limit: no more rows to
			// scan, so drop the last-scanned cursor to signal true end-of-data
			// (distinct from a scan-cap resume, which keeps the cursor).
			result.NextCursor = nil
			return result, nil
		}
		cursor = *page.NextCursor
	}
	return result, nil
}

func (s ReadService) enrich(ctx context.Context, items []domain.ListingReadModel, ceiling *float64, policy ports.PricingPolicy, policyFound bool) error {
	ids := make([]int64, 0, len(items))
	seen := map[int64]bool{}
	for _, item := range items {
		if item.Link.ProductID != nil {
			if id, err := strconv.ParseInt(*item.Link.ProductID, 10, 64); err == nil && !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
	}
	costs, err := s.facts.GetCostFactsByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("read listing costs: %w", err)
	}
	for i := range items {
		item := &items[i]
		var amount *float64
		if item.Link.ProductID != nil {
			if id, e := strconv.ParseInt(*item.Link.ProductID, 10, 64); e == nil {
				if fact := costs[id]; fact != nil {
					amount = fact.Amount
				}
			}
		}
		if amount != nil {
			item.Cost = &domain.Money{Amount: strconv.FormatFloat(*amount, 'f', -1, 64), Currency: priceCurrency(item.Price)}
		}
		var margin *float64
		if policyFound {
			margin = &policy.MinMarginPercent
		}
		item.BelowMarginWorstCase = belowMargin(item.Price, amount, ceiling, margin)
		item.ICMSWorstCaseByUF = nil
		item.PendingIssue = pendingIssue(*item)
	}
	return nil
}

func priceCurrency(price *domain.Money) domain.PriceCurrency {
	if price != nil {
		return price.Currency
	}
	return "BRL"
}
func maximumCeiling(values map[int64]*ports.ICMSCeiling) *float64 {
	var max *float64
	for _, v := range values {
		if v != nil && v.Percent != nil && (max == nil || *v.Percent > *max) {
			n := *v.Percent
			max = &n
		}
	}
	return max
}
func needsBelowMarginScan(f domain.ListingFilter) bool {
	return f.Exception == domain.ListingExceptionBelowMargin || f.HasException != nil
}
func matchesDependentFilter(item domain.ListingReadModel, f domain.ListingFilter) bool {
	active := sqlException(item) || item.BelowMarginWorstCase != nil && *item.BelowMarginWorstCase
	if f.Exception == domain.ListingExceptionBelowMargin {
		return item.BelowMarginWorstCase != nil && *item.BelowMarginWorstCase
	}
	if f.HasException == nil {
		return true
	}
	if *f.HasException {
		return active
	}
	return item.BelowMarginWorstCase != nil && !active
}
func sqlException(i domain.ListingReadModel) bool {
	return i.SyncState == domain.ListingSyncStateError || i.SyncError != nil || i.SyncState == domain.ListingSyncStateStale || i.Link.State != domain.LinkStateResolved
}
func pendingIssue(i domain.ListingReadModel) *domain.PendingIssue {
	if i.SyncState == domain.ListingSyncStateError || i.SyncError != nil {
		return &domain.PendingIssue{Kind: domain.PendingIssueSyncError, MessagePT: "Erro de sincronização"}
	}
	if i.SyncState == domain.ListingSyncStateStale {
		return &domain.PendingIssue{Kind: domain.PendingIssueStale, MessagePT: "Dados desatualizados"}
	}
	if i.BelowMarginWorstCase != nil && *i.BelowMarginWorstCase {
		return &domain.PendingIssue{Kind: domain.PendingIssueBelowMargin, MessagePT: "Margem abaixo do mínimo no pior caso"}
	}
	if i.Link.State != domain.LinkStateResolved {
		return &domain.PendingIssue{Kind: domain.PendingIssueUnlinked, MessagePT: "Produto não vinculado"}
	}
	return nil
}

func belowMargin(price *domain.Money, cost, ceiling, margin *float64) *bool {
	if price == nil || cost == nil || ceiling == nil || margin == nil {
		return nil
	}
	gross, ok := new(big.Rat).SetString(price.Amount)
	if !ok {
		return nil
	}
	one := big.NewRat(1, 1)
	pct := new(big.Rat).Quo(new(big.Rat).SetFloat64(*ceiling), big.NewRat(100, 1))
	net := new(big.Rat).Mul(gross, new(big.Rat).Sub(one, pct))
	costRat := new(big.Rat).SetFloat64(*cost)
	marginRat := new(big.Rat).Quo(new(big.Rat).SetFloat64(*margin), big.NewRat(100, 1))
	threshold := new(big.Rat).Mul(costRat, new(big.Rat).Add(one, marginRat))
	result := net.Cmp(threshold) < 0
	return &result
}
