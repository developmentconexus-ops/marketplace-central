package domain

import (
	"testing"
	"time"
)

func TestClassifyStockRiskBlocksUnresolvedAndConflictLinks(t *testing.T) {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)

	cases := []struct {
		name      string
		linkState LinkState
		want      StockRiskState
	}{
		{name: "missing link", linkState: LinkStateNone, want: StockRiskUnresolved},
		{name: "unresolved link", linkState: LinkStateUnresolved, want: StockRiskUnresolved},
		{name: "rejected link", linkState: LinkStateRejected, want: StockRiskUnresolved},
		{name: "conflict link", linkState: LinkStateConflict, want: StockRiskConflict},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			row := ClassifyStockRisk(RiskInput{
				Policy: DefaultStockPolicy(),
				Link:   LinkEvidence{State: tt.linkState},
				Now:    now,
			})

			assertRiskState(t, row, tt.want)
			if row.RecommendedQuantity != nil {
				t.Fatalf("blocked link should not have recommendation, got %d", *row.RecommendedQuantity)
			}
		})
	}
}

func TestClassifyStockRiskBlocksStaleAndUnsupportedSources(t *testing.T) {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)
	old := now.Add(-DefaultFreshnessMaxAge - time.Second)
	fresh := now.Add(-time.Minute)

	cases := []struct {
		name     string
		input    RiskInput
		want     StockRiskState
		wantCode string
	}{
		{
			name:     "missing internal quantity",
			input:    baseRiskInput(now, nil, intPtr(4), &fresh, &fresh),
			want:     StockRiskStale,
			wantCode: "missing_internal_stock",
		},
		{
			name:     "stale internal source",
			input:    baseRiskInput(now, intPtr(5), intPtr(4), &old, &fresh),
			want:     StockRiskStale,
			wantCode: "stale_internal_source",
		},
		{
			name:     "stale provider source",
			input:    baseRiskInput(now, intPtr(5), intPtr(4), &fresh, &old),
			want:     StockRiskStale,
			wantCode: "stale_provider_source",
		},
		{
			name:     "missing internal source timestamp",
			input:    baseRiskInput(now, intPtr(5), intPtr(4), nil, &fresh),
			want:     StockRiskStale,
			wantCode: "stale_internal_source",
		},
		{
			name:     "missing provider source timestamp",
			input:    baseRiskInput(now, intPtr(5), intPtr(4), &fresh, nil),
			want:     StockRiskStale,
			wantCode: "stale_provider_source",
		},
		{
			name:     "unsupported provider quantity",
			input:    baseRiskInput(now, intPtr(5), nil, &fresh, &fresh),
			want:     StockRiskUnsupported,
			wantCode: "missing_provider_quantity",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			row := ClassifyStockRisk(tt.input)
			assertRiskState(t, row, tt.want)
			if row.BlockingReason.Code != tt.wantCode {
				t.Fatalf("blocking code=%q, want %q", row.BlockingReason.Code, tt.wantCode)
			}
			if row.BlockingReason.Message == "" {
				t.Fatalf("expected blocking message for %s", tt.name)
			}
			if row.RecommendedQuantity != nil {
				t.Fatalf("blocked source should not have recommendation, got %d", *row.RecommendedQuantity)
			}
		})
	}
}

func TestClassifyStockRiskPreservesTimestampsOnBlockedRows(t *testing.T) {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)
	old := now.Add(-DefaultFreshnessMaxAge - time.Second)
	fresh := now.Add(-time.Minute)

	row := ClassifyStockRisk(baseRiskInput(now, intPtr(5), intPtr(4), &old, &fresh))

	assertRiskState(t, row, StockRiskStale)
	if row.InternalObservedAt == nil || !row.InternalObservedAt.Equal(old) {
		t.Fatalf("expected stale internal timestamp to stay visible, got %+v", row.InternalObservedAt)
	}
	if row.ProviderObservedAt == nil || !row.ProviderObservedAt.Equal(fresh) {
		t.Fatalf("expected provider timestamp to stay visible, got %+v", row.ProviderObservedAt)
	}
}

func TestClassifyStockRiskBlocksIneligibleProducts(t *testing.T) {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	input := baseRiskInput(now, intPtr(6), intPtr(5), &fresh, &fresh)
	input.Policy.Eligibility = EligibilityPolicy{
		Rules: []EligibilityRule{GroupExclusion(77, "showroom-only group")},
	}
	input.Product = ProductEvidence{ProductID: 20318, GroupID: intPtr(77)}

	row := ClassifyStockRisk(input)

	assertRiskState(t, row, StockRiskIneligible)
	if row.BlockingReason.Code != "ineligible_product" {
		t.Fatalf("blocking code=%q, want ineligible_product", row.BlockingReason.Code)
	}
	if row.BlockingReason.Message != "showroom-only group" {
		t.Fatalf("blocking message=%q, want explicit eligibility reason", row.BlockingReason.Message)
	}
	if row.RecommendedQuantity != nil {
		t.Fatalf("ineligible product should not have recommendation, got %d", *row.RecommendedQuantity)
	}
}

func TestClassifyStockRiskComparesProviderQuantityToRecommendation(t *testing.T) {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)

	cases := []struct {
		name     string
		provider int
		want     StockRiskState
	}{
		{name: "provider above recommendation is oversell", provider: 9, want: StockRiskOversell},
		{name: "provider below recommendation is undersell", provider: 6, want: StockRiskUndersell},
		{name: "provider equals recommendation is healthy", provider: 7, want: StockRiskHealthy},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			row := ClassifyStockRisk(baseRiskInput(now, intPtr(8), intPtr(tt.provider), &fresh, &fresh))

			assertRiskState(t, row, tt.want)
			if row.RecommendedQuantity == nil || *row.RecommendedQuantity != 7 {
				t.Fatalf("recommended quantity=%v, want 7", row.RecommendedQuantity)
			}
			if row.PolicyID != "stock-seguro-default" {
				t.Fatalf("policy id=%q, want stock-seguro-default", row.PolicyID)
			}
			if row.InternalObservedAt == nil || row.ProviderObservedAt == nil {
				t.Fatalf("expected source timestamps in row: %+v", row)
			}
		})
	}
}

func baseRiskInput(now time.Time, internalQty *int, providerQty *int, internalObservedAt *time.Time, providerObservedAt *time.Time) RiskInput {
	return RiskInput{
		Policy: DefaultStockPolicy(),
		Link: LinkEvidence{
			State:             LinkStateResolved,
			InternalProductID: intPtr(20318),
		},
		InternalStock: InternalStockEvidence{
			Quantity: internalQty,
			Source:   SourceEvidence{ObservedAt: internalObservedAt},
		},
		ProviderStock: ProviderStockEvidence{
			AvailableQuantity: providerQty,
			Source:            SourceEvidence{ObservedAt: providerObservedAt},
		},
		Product: ProductEvidence{ProductID: 20318},
		Now:     now,
	}
}

func assertRiskState(t *testing.T, row StockRiskRow, want StockRiskState) {
	t.Helper()
	if row.State != want {
		t.Fatalf("state=%s, want %s; row=%+v", row.State, want, row)
	}
}
