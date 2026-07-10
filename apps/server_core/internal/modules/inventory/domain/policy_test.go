package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultStockPolicyPreservesMissionDefaults(t *testing.T) {
	policy := DefaultStockPolicy()

	if !reflect.DeepEqual(policy.Scope.CompanyIDs, []int{1, 2}) {
		t.Fatalf("expected companies [1 2], got %v", policy.Scope.CompanyIDs)
	}
	if !reflect.DeepEqual(policy.Scope.LocationIDs, []int{10101}) {
		t.Fatalf("expected locations [10101], got %v", policy.Scope.LocationIDs)
	}
	if !reflect.DeepEqual(policy.Scope.ExcludedLocationIDs, []int{10108}) {
		t.Fatalf("expected excluded showroom [10108], got %v", policy.Scope.ExcludedLocationIDs)
	}
	if policy.Scope.Formula != "SUM(ESTOQUE - RESERVADO)" {
		t.Fatalf("expected sellable formula, got %q", policy.Scope.Formula)
	}
	if policy.Buffer.Quantity != 1 {
		t.Fatalf("expected default buffer 1, got %d", policy.Buffer.Quantity)
	}
	if policy.Freshness.MaxAge <= 0 {
		t.Fatalf("expected default freshness threshold, got %s", policy.Freshness.MaxAge)
	}
}

func TestRecommendedQuantityAppliesBufferAndNeverGoesNegative(t *testing.T) {
	policy := DefaultStockPolicy()

	cases := []struct {
		name     string
		internal int
		want     int
	}{
		{name: "positive stock subtracts buffer", internal: 8, want: 7},
		{name: "stock equal to buffer becomes zero", internal: 1, want: 0},
		{name: "stock below buffer stays zero", internal: 0, want: 0},
		{name: "negative stock stays zero", internal: -4, want: 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := policy.RecommendedProviderQuantity(tt.internal); got != tt.want {
				t.Fatalf("RecommendedProviderQuantity(%d)=%d, want %d", tt.internal, got, tt.want)
			}
		})
	}
}

func TestNegativeBufferCannotIncreaseRecommendedQuantity(t *testing.T) {
	policy := DefaultStockPolicy()
	policy.Buffer.Quantity = -3

	if got := policy.RecommendedProviderQuantity(8); got != 8 {
		t.Fatalf("negative buffer should not increase recommendation, got %d", got)
	}
}

func TestFreshnessPolicyClassifiesMissingAndOldSources(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	policy := FreshnessPolicy{MaxAge: 30 * time.Minute}
	freshObservedAt := now.Add(-30 * time.Minute)
	staleObservedAt := now.Add(-30*time.Minute - time.Second)

	if got := policy.Evaluate(SourceEvidence{ObservedAt: nil}, now); got.State != SourceFreshnessStale {
		t.Fatalf("missing observed timestamp should be stale, got %s", got.State)
	}
	if got := policy.Evaluate(SourceEvidence{ObservedAt: &freshObservedAt}, now); got.State != SourceFreshnessFresh {
		t.Fatalf("threshold age should still be fresh, got %s", got.State)
	}
	if got := policy.Evaluate(SourceEvidence{ObservedAt: &staleObservedAt}, now); got.State != SourceFreshnessStale {
		t.Fatalf("older than threshold should be stale, got %s", got.State)
	}
}

func TestEligibilityPolicyReturnsVisibleIneligibleReason(t *testing.T) {
	policy := EligibilityPolicy{
		Rules: []EligibilityRule{
			GroupExclusion(77, "showroom-only product group"),
		},
	}

	result := policy.Evaluate(ProductEvidence{ProductID: 20318, GroupID: intPtr(77)})
	if result.State != ProductEligibilityIneligible {
		t.Fatalf("expected ineligible, got %s", result.State)
	}
	if result.Reason != "showroom-only product group" {
		t.Fatalf("expected visible reason, got %q", result.Reason)
	}
}

func TestEligibilityPolicyCanModelFutureRuleTypes(t *testing.T) {
	policy := EligibilityPolicy{
		Rules: []EligibilityRule{
			{
				Type:   EligibilityRuleTypeWeight,
				Reason: "too heavy for Stock Seguro automation",
				Matches: func(product ProductEvidence) bool {
					return product.WeightGrams != nil && *product.WeightGrams > 30000
				},
			},
		},
	}

	result := policy.Evaluate(ProductEvidence{ProductID: 20318, WeightGrams: intPtr(32000)})
	if result.State != ProductEligibilityIneligible {
		t.Fatalf("expected ineligible, got %s", result.State)
	}
	if result.RuleType != EligibilityRuleTypeWeight {
		t.Fatalf("expected weight rule type, got %s", result.RuleType)
	}
}

func intPtr(value int) *int {
	return &value
}
