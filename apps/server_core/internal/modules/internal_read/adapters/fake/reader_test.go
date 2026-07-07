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
				Codprod:         42664,
				Quantity:        float64ptr(3),
				Scope:           domain.DefaultSellableStockScope(),
				SourceFetchedAt: time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
			},
		},
	})

	got, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{Codprod: 42664})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity == nil || *got.Quantity != 3 {
		t.Fatalf("expected revenda stock only, got %v", got.Quantity)
	}
	for _, location := range got.Scope.Locations {
		if location == 10108 {
			t.Fatal("expected showroom location 10108 to stay excluded")
		}
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

	got, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{Codprod: 42664})
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

	got, err := reader.GetCostAsOf(context.Background(), ports.CostAsOfInput{Codprod: 42664, Codemp: 1, SaleDate: "2026-07-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CUSSEMICM != nil {
		t.Fatalf("expected nil cost, got %v", *got.CUSSEMICM)
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingCost) {
		t.Fatalf("expected missing_cost flag, got %v", got.QualityFlags)
	}
}

func TestFakeReaderMissingTaxRemainsFlagged(t *testing.T) {
	reader := NewReader(Fixtures{})

	got, err := reader.GetTaxInputs(context.Background(), ports.TaxInput{Codprod: 42664, SaleDate: "2026-07-06"})
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

func float64ptr(v float64) *float64 {
	return &v
}
