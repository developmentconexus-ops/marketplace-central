package domain

import "time"

const (
	DefaultStockFormula               = "SUM(ESTOQUE - RESERVADO)"
	StockScopeResale       StockScope = "resale"
	StockScopeCustom       StockScope = "custom"
	DefaultFreshnessMaxAge            = 30 * time.Minute
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
		Freshness: FreshnessPolicy{MaxAge: DefaultFreshnessMaxAge},
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
}

type SourceFreshnessState string

const (
	SourceFreshnessFresh SourceFreshnessState = "fresh"
	SourceFreshnessStale SourceFreshnessState = "stale"
)

type FreshnessPolicy struct {
	MaxAge time.Duration
}

type FreshnessResult struct {
	State  SourceFreshnessState
	Reason string
}

func (p FreshnessPolicy) Evaluate(source SourceEvidence, now time.Time) FreshnessResult {
	if source.ObservedAt == nil {
		return FreshnessResult{State: SourceFreshnessStale, Reason: "missing_source_timestamp"}
	}
	if p.MaxAge > 0 && now.Sub(*source.ObservedAt) > p.MaxAge {
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
