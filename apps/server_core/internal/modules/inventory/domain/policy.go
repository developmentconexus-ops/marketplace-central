package domain

import (
	"time"

	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

const (
	DefaultStockFormula               = "SUM(ESTOQUE - RESERVADO)"
	StockScopeResale       StockScope = "resale"
	StockScopeCustom       StockScope = "custom"
	DefaultFreshnessMaxAge            = 30 * time.Minute
	// DefaultSnapshotFreshnessMaxAge bounds a snapshot-shaped source (an xlsx or
	// catalogo_cliente upload). Its data is only as new as the last import, so the
	// live 30-minute bar would mark EVERY row stale minutes after a successful
	// import and Stock Seguro could never act — the operator's cadence is the
	// upload, and a snapshot older than a day no longer describes today's stock.
	DefaultSnapshotFreshnessMaxAge = 24 * time.Hour
)

type StockScope string

type StockSourceScope struct {
	CompanyIDs          []int
	LocationIDs         []int
	ExcludedLocationIDs []int
	Formula             string
	Scope               StockScope
}

type StockBuffer struct {
	Quantity int
}

type StockPolicy struct {
	PolicyID    string
	Scope       StockSourceScope
	Buffer      StockBuffer
	Freshness   FreshnessPolicy
	Eligibility EligibilityPolicy
}

func DefaultStockPolicy() StockPolicy {
	return StockPolicy{
		PolicyID: "stock-seguro-default",
		Scope: StockSourceScope{
			CompanyIDs:          []int{1, 2},
			LocationIDs:         []int{10101},
			ExcludedLocationIDs: []int{10108},
			Formula:             DefaultStockFormula,
			Scope:               StockScopeResale,
		},
		Buffer:    StockBuffer{Quantity: 1},
		Freshness: FreshnessPolicy{MaxAge: DefaultFreshnessMaxAge, SnapshotMaxAge: DefaultSnapshotFreshnessMaxAge},
	}
}

func (p StockPolicy) RecommendedProviderQuantity(internalSellableQuantity int) int {
	buffer := p.Buffer.Quantity
	if buffer < 0 {
		buffer = 0
	}
	recommended := internalSellableQuantity - buffer
	if recommended < 0 {
		return 0
	}
	return recommended
}

type SourceEvidence struct {
	ObservedAt *time.Time
	// Kind says whether the observation came from a live read or from an uploaded
	// snapshot; the freshness bar differs because the cadences differ. An empty
	// kind is read as live (the stricter bar), never as the lenient one.
	Kind sourcekind.SourceKind
}

type SourceFreshnessState string

const (
	SourceFreshnessFresh SourceFreshnessState = "fresh"
	SourceFreshnessStale SourceFreshnessState = "stale"
)

type FreshnessPolicy struct {
	// MaxAge bounds a live source (the ERP read through, the marketplace API).
	MaxAge time.Duration
	// SnapshotMaxAge bounds an upload-snapshot source. Zero falls back to MaxAge,
	// so a policy that does not distinguish the two keeps its old behaviour.
	SnapshotMaxAge time.Duration
}

// maxAgeFor picks the bar that matches how the observation was produced.
func (p FreshnessPolicy) maxAgeFor(kind sourcekind.SourceKind) time.Duration {
	if kind == sourcekind.UploadSnapshot && p.SnapshotMaxAge > 0 {
		return p.SnapshotMaxAge
	}
	return p.MaxAge
}

type FreshnessResult struct {
	State  SourceFreshnessState
	Reason string
}

func (p FreshnessPolicy) Evaluate(source SourceEvidence, now time.Time) FreshnessResult {
	if source.ObservedAt == nil {
		return FreshnessResult{State: SourceFreshnessStale, Reason: "missing_source_timestamp"}
	}
	if maxAge := p.maxAgeFor(source.Kind); maxAge > 0 && now.Sub(*source.ObservedAt) > maxAge {
		return FreshnessResult{State: SourceFreshnessStale, Reason: "source_older_than_policy"}
	}
	return FreshnessResult{State: SourceFreshnessFresh}
}

type ProductEvidence struct {
	ProductID   int
	GroupID     *int
	WeightGrams *int
}

type ProductEligibilityState string

const (
	ProductEligibilityEligible   ProductEligibilityState = "eligible"
	ProductEligibilityIneligible ProductEligibilityState = "ineligible"
)

type ProductEligibilityResult struct {
	State    ProductEligibilityState
	Reason   string
	RuleType EligibilityRuleType
}

type EligibilityRuleType string

const (
	EligibilityRuleTypeGroup  EligibilityRuleType = "group"
	EligibilityRuleTypeWeight EligibilityRuleType = "weight"
	EligibilityRuleTypeSize   EligibilityRuleType = "size"
	EligibilityRuleTypeMargin EligibilityRuleType = "margin"
)

type EligibilityRule struct {
	Type    EligibilityRuleType
	Reason  string
	Matches func(ProductEvidence) bool
}

type EligibilityPolicy struct {
	Rules []EligibilityRule
}

func GroupExclusion(groupID int, reason string) EligibilityRule {
	return EligibilityRule{
		Type:   EligibilityRuleTypeGroup,
		Reason: reason,
		Matches: func(product ProductEvidence) bool {
			return product.GroupID != nil && *product.GroupID == groupID
		},
	}
}

func (p EligibilityPolicy) Evaluate(product ProductEvidence) ProductEligibilityResult {
	for _, rule := range p.Rules {
		if rule.Matches != nil && rule.Matches(product) {
			return ProductEligibilityResult{
				State:    ProductEligibilityIneligible,
				Reason:   rule.Reason,
				RuleType: rule.Type,
			}
		}
	}
	return ProductEligibilityResult{State: ProductEligibilityEligible}
}
