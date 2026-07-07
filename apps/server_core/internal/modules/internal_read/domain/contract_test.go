package domain

import (
	"reflect"
	"testing"
)

func TestDefaultSellableStockScopePreservesMissionDefaults(t *testing.T) {
	scope := DefaultSellableStockScope()

	if scope.Formula != "SUM(ESTOQUE - RESERVADO)" {
		t.Fatalf("expected stock formula to stay unchanged, got %q", scope.Formula)
	}
	if scope.ScopeCode != StockScopeCodeRevenda {
		t.Fatalf("expected scope code %q, got %q", StockScopeCodeRevenda, scope.ScopeCode)
	}
	if !reflect.DeepEqual(scope.Companies, []int{1, 2}) {
		t.Fatalf("expected default companies [1 2], got %v", scope.Companies)
	}
	if !reflect.DeepEqual(scope.Locations, []int{10101}) {
		t.Fatalf("expected default locations [10101], got %v", scope.Locations)
	}
	for _, location := range scope.Locations {
		if location == 10108 {
			t.Fatalf("expected showroom location 10108 to stay excluded from default scope")
		}
	}
}

func TestRequiredQualityFlagsRemainExplicit(t *testing.T) {
	flags := []QualityFlag{
		QualityComplete,
		QualityMissingProduct,
		QualityMissingStock,
		QualityMissingCost,
		QualityMissingTax,
		QualityAmbiguousProduct,
		QualityStaleSource,
	}

	expected := []QualityFlag{
		"complete",
		"missing_product",
		"missing_stock",
		"missing_cost",
		"missing_tax",
		"ambiguous_product",
		"stale_source",
	}

	if !reflect.DeepEqual(flags, expected) {
		t.Fatalf("expected quality flags %v, got %v", expected, flags)
	}
}

func TestProductCandidateUsesTaskContractFields(t *testing.T) {
	candidate := ProductCandidate{
		Codprod:      42664,
		Produto:      "Produto teste",
		EAN:          "7890000000000",
		Reference:    "SKU-42664",
		QualityFlags: []QualityFlag{QualityComplete},
	}

	if candidate.EAN == "" || candidate.Reference == "" {
		t.Fatalf("expected product candidate EAN and reference to use explicit task contract strings")
	}
	if !reflect.DeepEqual(candidate.QualityFlags, []QualityFlag{QualityComplete}) {
		t.Fatalf("expected complete quality flag, got %v", candidate.QualityFlags)
	}
}

func TestMissingCostStaysNilWithQualityFlag(t *testing.T) {
	cost := CostAsOf{
		Codprod:      42664,
		Codemp:       1,
		SaleDate:     "2026-07-06",
		CUSSEMICM:    nil,
		QualityFlags: []QualityFlag{QualityMissingCost},
	}

	if cost.CUSSEMICM != nil {
		t.Fatalf("expected missing CUSSEMICM to stay nil, got %v", *cost.CUSSEMICM)
	}
	if !reflect.DeepEqual(cost.QualityFlags, []QualityFlag{QualityMissingCost}) {
		t.Fatalf("expected missing cost quality flag, got %v", cost.QualityFlags)
	}
}
