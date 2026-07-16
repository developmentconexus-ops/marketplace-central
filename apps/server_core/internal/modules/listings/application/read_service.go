package application

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
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

// maxBelowMarginGroupScanPages separately bounds group-key pages because each
// page expands to complete child sets and one cost batch; clients resume from
// the last scanned group key when the bound is reached.
const maxBelowMarginGroupScanPages = 50

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

func (s ReadService) Summary(ctx context.Context, query ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	if strings.TrimSpace(query.InstallationID) == "" {
		return ports.ListingSummaryRow{}, ErrInstallationIDRequired
	}
	found, err := s.installations.InstallationExists(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingSummaryRow{}, fmt.Errorf("validate installation: %w", err)
	}
	if !found {
		return ports.ListingSummaryRow{}, ErrInstallationNotFound
	}
	policy, policyFound, err := s.policies.GetPricingPolicyForInstallation(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingSummaryRow{}, fmt.Errorf("read pricing policy: %w", err)
	}
	ceilings, ceilingErr := s.facts.GetICMSCeilingByOrigin(ctx, originUFMG)
	row, err := s.repo.GetListingsSummary(ctx, query)
	if err != nil {
		return ports.ListingSummaryRow{}, err
	}
	row.AsOf = s.now().UTC()
	if ceilingErr != nil {
		row.BelowMarginWorstCase, row.MarginUnknown = nil, nil
		return row, nil
	}
	ids := make([]int64, 0, len(row.Linked))
	seen := make(map[int64]bool, len(row.Linked))
	for _, linked := range row.Linked {
		if !seen[linked.CostID] {
			seen[linked.CostID] = true
			ids = append(ids, linked.CostID)
		}
	}
	costs := map[int64]*ports.CostFact{}
	if len(ids) > 0 {
		costs, err = s.facts.GetCostFactsByIDs(ctx, ids)
		if err != nil {
			row.BelowMarginWorstCase, row.MarginUnknown = nil, nil
			return row, nil
		}
	}
	below, unknown := 0, 0
	ceiling := maximumCeiling(ceilings)
	var margin *float64
	if policyFound {
		margin = &policy.MinMarginPercent
	}
	for _, linked := range row.Linked {
		var cost *float64
		if fact := costs[linked.CostID]; fact != nil {
			cost = fact.Amount
		}
		value := belowMargin(linked.Price, cost, ceiling, margin)
		if value == nil {
			unknown++
		} else if *value {
			below++
		}
	}
	row.BelowMarginWorstCase, row.MarginUnknown = &below, &unknown
	return row, nil
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
	ceilings, ceilingErr := s.facts.GetICMSCeilingByOrigin(ctx, originUFMG)
	factsUnavailable := ceilingErr != nil
	maxCeiling := maximumCeiling(ceilings)
	serveTime := s.now().UTC()
	if !needsBelowMarginScan(query.Filter) {
		page, err := s.repo.ListListingRows(ctx, query)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		_, err = s.enrich(ctx, page.Items, maxCeiling, policy, policyFound, factsUnavailable)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		page.AsOf = serveTime
		return page, nil
	}
	return s.scan(ctx, query, maxCeiling, policy, policyFound, factsUnavailable, serveTime)
}

func (s ReadService) ByProduct(ctx context.Context, query ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	if strings.TrimSpace(query.InstallationID) == "" {
		return ports.ListingGroupRowPage{}, ErrInstallationIDRequired
	}
	found, err := s.installations.InstallationExists(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingGroupRowPage{}, fmt.Errorf("validate installation: %w", err)
	}
	if !found {
		return ports.ListingGroupRowPage{}, ErrInstallationNotFound
	}
	policy, policyFound, err := s.policies.GetPricingPolicyForInstallation(ctx, query.InstallationID)
	if err != nil {
		return ports.ListingGroupRowPage{}, fmt.Errorf("read pricing policy: %w", err)
	}
	ceilings, ceilingErr := s.facts.GetICMSCeilingByOrigin(ctx, originUFMG)
	factsUnavailable := ceilingErr != nil
	ceiling, asOf := maximumCeiling(ceilings), s.now().UTC()
	if needsBelowMarginScan(query.Filter) {
		return s.scanGroups(ctx, query, ceiling, policy, policyFound, factsUnavailable, asOf)
	}
	page, err := s.repo.ListListingGroupRows(ctx, query)
	if err != nil {
		return ports.ListingGroupRowPage{}, err
	}
	_, err = s.enrichGroups(ctx, page.Groups, ceiling, policy, policyFound, factsUnavailable)
	if err != nil {
		return ports.ListingGroupRowPage{}, err
	}
	finalizeGroups(page.Groups)
	page.AsOf = asOf
	return page, nil
}

func (s ReadService) scanGroups(ctx context.Context, q ports.ListingGroupQuery, ceiling *float64, policy ports.PricingPolicy, policyFound, factsUnavailable bool, asOf time.Time) (ports.ListingGroupRowPage, error) {
	if factsUnavailable {
		return s.passThroughGroups(ctx, q, ceiling, policy, policyFound, asOf)
	}
	result := ports.ListingGroupRowPage{Groups: []domain.ListingGroup{}, AsOf: asOf}
	cursor := q.Cursor
	for pageNo := 0; pageNo < maxBelowMarginGroupScanPages; pageNo++ {
		candidate := q
		candidate.Cursor = cursor
		page, err := s.repo.ListListingGroupRows(ctx, candidate)
		if err != nil {
			return ports.ListingGroupRowPage{}, err
		}
		factsUnavailable, err = s.enrichGroups(ctx, page.Groups, ceiling, policy, policyFound, factsUnavailable)
		if err != nil {
			return ports.ListingGroupRowPage{}, err
		}
		if factsUnavailable {
			return s.passThroughGroups(ctx, q, ceiling, policy, policyFound, asOf)
		}
		for groupIndex := range page.Groups {
			group := page.Groups[groupIndex]
			last := cursorForGroup(group)
			result.NextCursor = &last
			survivors := make([]domain.ListingReadModel, 0, len(group.Listings))
			for _, child := range group.Listings {
				if matchesDependentFilter(child, q.Filter) {
					survivors = append(survivors, child)
				}
			}
			if len(survivors) == 0 {
				continue
			}
			group.Listings = survivors
			group.ListingCount = len(survivors)
			group.GroupState = groupState(survivors)
			result.Groups = append(result.Groups, group)
			if len(result.Groups) == q.Limit {
				if page.NextCursor == nil && groupIndex == len(page.Groups)-1 {
					result.NextCursor = nil
				}
				return result, nil
			}
		}
		if page.NextCursor == nil {
			result.NextCursor = nil
			return result, nil
		}
		cursor = *page.NextCursor
	}
	return result, nil
}

func (s ReadService) passThroughGroups(ctx context.Context, q ports.ListingGroupQuery, ceiling *float64, policy ports.PricingPolicy, policyFound bool, asOf time.Time) (ports.ListingGroupRowPage, error) {
	page, err := s.repo.ListListingGroupRows(ctx, q)
	if err != nil {
		return ports.ListingGroupRowPage{}, err
	}
	if _, err := s.enrichGroups(ctx, page.Groups, ceiling, policy, policyFound, true); err != nil {
		return ports.ListingGroupRowPage{}, err
	}
	finalizeGroups(page.Groups)
	page.AsOf = asOf
	return page, nil
}

func (s ReadService) enrichGroups(ctx context.Context, groups []domain.ListingGroup, ceiling *float64, policy ports.PricingPolicy, policyFound, factsUnavailable bool) (bool, error) {
	items := make([]domain.ListingReadModel, 0)
	for _, group := range groups {
		items = append(items, group.Listings...)
	}
	factsUnavailable, err := s.enrich(ctx, items, ceiling, policy, policyFound, factsUnavailable)
	if err != nil {
		return false, err
	}
	offset := 0
	for i := range groups {
		n := len(groups[i].Listings)
		copy(groups[i].Listings, items[offset:offset+n])
		offset += n
	}
	return factsUnavailable, nil
}

func finalizeGroups(groups []domain.ListingGroup) {
	for i := range groups {
		groups[i].ListingCount = len(groups[i].Listings)
		groups[i].GroupState = groupState(groups[i].Listings)
	}
}
func groupState(items []domain.ListingReadModel) domain.GroupState {
	state := domain.GroupStateOK
	for _, item := range items {
		if item.SyncState == domain.ListingSyncStateError || item.SyncError != nil {
			return domain.GroupStateError
		}
		if item.SyncState == domain.ListingSyncStateStale || item.Link.State != domain.LinkStateResolved || item.BelowMarginWorstCase != nil && *item.BelowMarginWorstCase {
			state = domain.GroupStateAttention
		}
	}
	return state
}
func cursorForGroup(group domain.ListingGroup) ports.GroupCursor {
	if group.ProductID == nil {
		return ports.GroupCursor{NullLast: true}
	}
	return ports.GroupCursor{ProductTitle: *group.ProductTitle, ProductID: *group.ProductID}
}

func (s ReadService) Get(ctx context.Context, id domain.ListingID) (domain.ListingReadModel, []domain.TimelineEvent, error) {
	key := domain.ListingKey{InstallationID: id.InstallationID, ProviderListingID: id.ProviderListingID, VariationID: id.VariationID}
	model, found, err := s.repo.GetListingRow(ctx, key)
	if err != nil {
		return domain.ListingReadModel{}, nil, fmt.Errorf("get listing row: %w", err)
	}
	if !found {
		return domain.ListingReadModel{}, nil, &domain.ListingNotFoundError{}
	}
	policy, policyFound, err := s.policies.GetPricingPolicyForInstallation(ctx, id.InstallationID)
	if err != nil {
		return domain.ListingReadModel{}, nil, fmt.Errorf("read pricing policy: %w", err)
	}
	ceilings, ceilingErr := s.facts.GetICMSCeilingByOrigin(ctx, originUFMG)
	factsUnavailable := ceilingErr != nil
	items := []domain.ListingReadModel{model}
	factsUnavailable, err = s.enrich(ctx, items, maximumCeiling(ceilings), policy, policyFound, factsUnavailable)
	if err != nil {
		return domain.ListingReadModel{}, nil, err
	}
	model = items[0]
	if ceilingErr != nil {
		model.ICMSWorstCaseByUF = nil
	} else {
		model.ICMSWorstCaseByUF = icmsWorstCaseByUF(model, ceilings, policy, policyFound)
	}
	timeline, err := s.repo.ListListingTimeline(ctx, key, 10)
	if err != nil {
		return domain.ListingReadModel{}, nil, fmt.Errorf("read listing timeline: %w", err)
	}
	if timeline == nil {
		timeline = []domain.TimelineEvent{}
	}
	return model, timeline, nil
}

func (s ReadService) scan(ctx context.Context, q ports.ListingQuery, ceiling *float64, policy ports.PricingPolicy, policyFound, factsUnavailable bool, asOf time.Time) (ports.ListingRowPage, error) {
	if factsUnavailable {
		return s.passThrough(ctx, q, ceiling, policy, policyFound, asOf)
	}
	result := ports.ListingRowPage{Items: make([]domain.ListingReadModel, 0, q.Limit), AsOf: asOf}
	cursor := q.Cursor
	for pageNo := 0; pageNo < maxBelowMarginScanPages; pageNo++ {
		candidate := q
		candidate.Cursor = cursor
		page, err := s.repo.ListListingRows(ctx, candidate)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		factsUnavailable, err = s.enrich(ctx, page.Items, ceiling, policy, policyFound, factsUnavailable)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		if factsUnavailable {
			return s.passThrough(ctx, q, ceiling, policy, policyFound, asOf)
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

func (s ReadService) passThrough(ctx context.Context, q ports.ListingQuery, ceiling *float64, policy ports.PricingPolicy, policyFound bool, asOf time.Time) (ports.ListingRowPage, error) {
	page, err := s.repo.ListListingRows(ctx, q)
	if err != nil {
		return ports.ListingRowPage{}, err
	}
	if _, err := s.enrich(ctx, page.Items, ceiling, policy, policyFound, true); err != nil {
		return ports.ListingRowPage{}, err
	}
	page.AsOf = asOf
	return page, nil
}

func (s ReadService) enrich(ctx context.Context, items []domain.ListingReadModel, ceiling *float64, policy ports.PricingPolicy, policyFound, factsUnavailable bool) (bool, error) {
	ids := make([]int64, 0, len(items))
	seen := map[int64]bool{}
	if !factsUnavailable {
		for _, item := range items {
			if item.Link.ProductID != nil {
				if id, err := strconv.ParseInt(*item.Link.ProductID, 10, 64); err == nil && !seen[id] {
					seen[id] = true
					ids = append(ids, id)
				}
			}
		}
	}
	costs := map[int64]*ports.CostFact{}
	if len(ids) > 0 {
		var err error
		costs, err = s.facts.GetCostFactsByIDs(ctx, ids)
		if err != nil {
			factsUnavailable = true
		}
	}
	for i := range items {
		item := &items[i]
		var amount *float64
		if factsUnavailable {
			item.Cost = nil
			item.BelowMarginWorstCase = nil
			item.ICMSWorstCaseByUF = nil
		} else {
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
		}
		item.PendingIssue = pendingIssue(*item)
	}
	return factsUnavailable, nil
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

func icmsWorstCaseByUF(item domain.ListingReadModel, ceilings map[int64]*ports.ICMSCeiling, policy ports.PricingPolicy, policyFound bool) *[]domain.ICMWorstCaseByUF {
	destinations := make([]int64, 0, len(ceilings))
	for destination := range ceilings {
		destinations = append(destinations, destination)
	}
	sort.Slice(destinations, func(i, j int) bool { return destinations[i] < destinations[j] })

	var cost *float64
	if item.Cost != nil {
		if parsed, err := strconv.ParseFloat(item.Cost.Amount, 64); err == nil {
			cost = &parsed
		}
	}
	var margin *float64
	if policyFound {
		margin = &policy.MinMarginPercent
	}
	rows := make([]domain.ICMWorstCaseByUF, 0, len(destinations))
	for _, destination := range destinations {
		row := domain.ICMWorstCaseByUF{DestinationUF: strconv.FormatInt(destination, 10)}
		ceiling := ceilings[destination]
		var percent *float64
		if ceiling != nil {
			percent = ceiling.Percent
		}
		if percent != nil {
			value := strconv.FormatFloat(*percent, 'f', -1, 64)
			row.WorstCaseICMSPct = &value
			row.PriceNetBasis = priceNetBasis(item.Price, percent)
			row.BelowMarginAtUF = belowMargin(item.Price, cost, percent, margin)
		}
		rows = append(rows, row)
	}
	return &rows
}

func priceNetBasis(price *domain.Money, ceiling *float64) *string {
	if price == nil || ceiling == nil {
		return nil
	}
	gross, ok := new(big.Rat).SetString(price.Amount)
	if !ok {
		return nil
	}
	pct := new(big.Rat).Quo(new(big.Rat).SetFloat64(*ceiling), big.NewRat(100, 1))
	net := new(big.Rat).Mul(gross, new(big.Rat).Sub(big.NewRat(1, 1), pct))
	value := net.FloatString(12)
	value = strings.TrimRight(strings.TrimRight(value, "0"), ".")
	return &value
}
