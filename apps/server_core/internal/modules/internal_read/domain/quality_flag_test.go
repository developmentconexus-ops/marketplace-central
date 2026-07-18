package domain

import (
	"reflect"
	"testing"
)

func TestQualityFlagsAreStable(t *testing.T) {
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

	want := []QualityFlag{
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

	if !reflect.DeepEqual(flags, want) {
		t.Fatalf("expected stable quality flags %v, got %v", want, flags)
	}
}

func TestReadErrorsRemainTyped(t *testing.T) {
	err := NewReadError(ReadErrorUnsupportedQuery, "unsupported query shape", nil)
	if !IsReadErrorCode(err, ReadErrorUnsupportedQuery) {
		t.Fatalf("expected read error code %q", ReadErrorUnsupportedQuery)
	}
	if IsReadErrorCode(err, ReadErrorSourceUnavailable) {
		t.Fatalf("did not expect read error code %q", ReadErrorSourceUnavailable)
	}
}
