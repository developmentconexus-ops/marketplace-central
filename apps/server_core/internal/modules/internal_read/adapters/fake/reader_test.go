package fake

import (
	"context"
	"slices"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestFakeReaderSellableStockExcludesShowroomByDefault(t *testing.T) {
	reader := NewReader(Fixtures{
		Stocks: map[int]domain.SellableStock{
			42664: {
				ProductID: 42664,
				Quantity:  float64ptr(3),
				Policy:    domain.DefaultSellableStockPolicy(),
				Source: domain.SourceMetadata{
					System:    "fake",
					FetchedAt: time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	})

	got, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{
		ProductID: 42664,
		Policy:    domain.DefaultSellableStockPolicy(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity == nil || *got.Quantity != 3 {
		t.Fatalf("expected revenda stock only, got %v", got.Quantity)
	}
	for _, location := range got.Policy.LocationIDs {
		if location == 10108 {
			t.Fatal("showroom location 10108 entered the selling whitelist")
		}
	}
	if !slices.Equal(got.Policy.LocationIDs, []int{10101, 10102}) {
		t.Fatalf("selling location whitelist = %v, want [10101 10102]", got.Policy.LocationIDs)
	}
}

func TestFakeReaderMissingProductReturnsUnresolvedCandidate(t *testing.T) {
	title := "Produto inexistente"
	reader := NewReader(Fixtures{})

	got, err := reader.FindProductsForLinking(context.Background(), ports.FindProductsInput{Title: &title})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one unresolved candidate, got %d", len(got))
	}
	if !slices.Contains(got[0].QualityFlags, domain.QualityMissingProduct) {
		t.Fatalf("expected missing_product flag, got %v", got[0].QualityFlags)
	}
}

func TestFakeReaderMissingStockStaysNilWithQualityFlag(t *testing.T) {
	reader := NewReader(Fixtures{})

	got, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{ProductID: 42664})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity != nil {
		t.Fatalf("expected missing stock quantity to stay nil, got %v", *got.Quantity)
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingStock) {
		t.Fatalf("expected missing_stock flag, got %v", got.QualityFlags)
	}
}

func TestFakeReaderMissingCostRemainsFlagged(t *testing.T) {
	reader := NewReader(Fixtures{})

	got, err := reader.GetCostAsOf(context.Background(), ports.CostAsOfInput{
		ProductID: 42664,
		Policy: domain.CostAsOfPolicy{
			CompanyID:   1,
			EffectiveAt: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),
			Basis:       domain.CostBasisCUSSEMICM,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != nil {
		t.Fatalf("expected nil cost, got %v", *got.Amount)
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingCost) {
		t.Fatalf("expected missing_cost flag, got %v", got.QualityFlags)
	}
}

func TestFakeReaderMissingTaxRemainsFlagged(t *testing.T) {
	reader := NewReader(Fixtures{})

	got, err := reader.GetTaxInputs(context.Background(), ports.TaxInput{
		ProductID: 42664,
		Policy:    domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ICMSAmount != nil || got.IPIAmount != nil || got.PISAmount != nil || got.COFINSAmount != nil {
		t.Fatal("expected missing tax numeric values to stay nil")
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingTax) {
		t.Fatalf("expected missing_tax flag, got %v", got.QualityFlags)
	}
}

func TestFakeReaderSelectsTaxByExactSourceAndReturnsProvenance(t *testing.T) {
	icms := 7.5
	identity := domain.TaxSourceIdentity{DocumentID: 98765, LineNumber: 3}
	reader := NewReader(Fixtures{Taxes: map[string]domain.TaxInputs{
		"42664:98765:3:0": {ProductID: 42664, SourceIdentity: identity, ICMSAmount: &icms, Source: domain.SourceMetadata{System: "fake"}},
	}})
	policy := domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC))
	policy.Source = identity

	got, err := reader.GetTaxInputs(context.Background(), ports.TaxInput{ProductID: 42664, Policy: policy})
	if err != nil {
		t.Fatalf("GetTaxInputs() error = %v", err)
	}
	if got.SourceIdentity != identity || got.ICMSAmount == nil || *got.ICMSAmount != icms || got.IPIAmount != nil {
		t.Fatalf("GetTaxInputs() = %+v, want exact provenance with partial tax preserved", got)
	}
}

func TestFakeReaderMissingPriceRemainsFlagged(t *testing.T) {
	reader := NewReader(Fixtures{})

	got, err := reader.GetCurrentPrice(context.Background(), ports.CurrentPriceInput{
		ProductID: 42664,
		Policy:    domain.DefaultCurrentPricePolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != nil {
		t.Fatalf("expected nil price, got %v", *got.Amount)
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingPrice) {
		t.Fatalf("expected missing_price flag, got %v", got.QualityFlags)
	}
}

func float64ptr(v float64) *float64 {
	return &v
}
