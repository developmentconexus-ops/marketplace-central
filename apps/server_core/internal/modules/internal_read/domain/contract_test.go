package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultSellableStockPolicyPreservesMissionDefaults(t *testing.T) {
	policy := DefaultSellableStockPolicy()

	if policy.Formula != "SUM(ESTOQUE - RESERVADO)" {
		t.Fatalf("expected stock formula to stay unchanged, got %q", policy.Formula)
	}
	if policy.Scope != SellableStockScopeResale {
		t.Fatalf("expected scope %q, got %q", SellableStockScopeResale, policy.Scope)
	}
	if !reflect.DeepEqual(policy.CompanyIDs, []int{1, 2}) {
		t.Fatalf("expected default companies [1 2], got %v", policy.CompanyIDs)
	}
	if !reflect.DeepEqual(policy.LocationIDs, []int{10101, 10102}) {
		t.Fatalf("expected default locations [10101 10102], got %v", policy.LocationIDs)
	}
	if !reflect.DeepEqual(policy.ExcludedLocationIDs, []int(nil)) {
		t.Fatalf("expected no exclusions outside the whitelist, got %v", policy.ExcludedLocationIDs)
	}
}

func TestRequiredQualityFlagsRemainExplicit(t *testing.T) {
	flags := []QualityFlag{
		QualityComplete,
		QualityMissingProduct,
		QualityMissingStock,
		QualityMissingPrice,
		QualityMissingCost,
		QualityMissingTax,
		QualityAmbiguousProduct,
		QualityStaleSource,
		QualityInvalidEAN,
		QualityEANCollision,
	}

	expected := []QualityFlag{
		"complete",
		"missing_product",
		"missing_stock",
		"missing_price",
		"missing_cost",
		"missing_tax",
		"ambiguous_product",
		"stale_source",
		"invalid_ean",
		"ean_collision",
	}

	if !reflect.DeepEqual(flags, expected) {
		t.Fatalf("expected quality flags %v, got %v", expected, flags)
	}
}

func TestProductCandidateUsesOracleFirstFields(t *testing.T) {
	ean := "7890000000000"
	reference := "SKU-42664"
	candidate := ProductCandidate{
		ProductID:     42664,
		Name:          "Produto teste",
		EAN:           &ean,
		ReferenceCode: &reference,
		IsActive:      true,
		QualityFlags:  []QualityFlag{QualityComplete},
	}

	if candidate.EAN == nil || candidate.ReferenceCode == nil {
		t.Fatal("expected product candidate to preserve EAN/reference pointers")
	}
	if !reflect.DeepEqual(candidate.QualityFlags, []QualityFlag{QualityComplete}) {
		t.Fatalf("expected complete quality flag, got %v", candidate.QualityFlags)
	}
}

func TestProductAndTaxNCMRemainNullable(t *testing.T) {
	candidate := ProductCandidate{}
	tax := TaxInputs{}
	if candidate.NCM != nil || tax.NCM != nil {
		t.Fatalf("absent NCM must remain nil: candidate=%v tax=%v", candidate.NCM, tax.NCM)
	}
}

func TestProductCandidateKeepsCanonicalIdentitySeparateFromLegacyID(t *testing.T) {
	canonical := InternalProductID(1001)
	candidate := ProductCandidate{InternalProductID: &canonical, ProductID: 77}
	if candidate.InternalProductID == nil || *candidate.InternalProductID != 1001 {
		t.Fatalf("canonical identity = %#v", candidate.InternalProductID)
	}
	if candidate.ProductID != 77 {
		t.Fatalf("legacy id = %d", candidate.ProductID)
	}
}

func TestMissingCostStaysNilWithQualityFlag(t *testing.T) {
	now := time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)
	cost := CostAsOf{
		ProductID:    42664,
		CompanyID:    1,
		Basis:        CostBasisCUSSEMICM,
		EffectiveAt:  now,
		Amount:       nil,
		QualityFlags: []QualityFlag{QualityMissingCost},
	}

	if cost.Amount != nil {
		t.Fatalf("expected missing CUSSEMICM cost to stay nil, got %v", *cost.Amount)
	}
	if !reflect.DeepEqual(cost.QualityFlags, []QualityFlag{QualityMissingCost}) {
		t.Fatalf("expected missing cost quality flag, got %v", cost.QualityFlags)
	}
}
