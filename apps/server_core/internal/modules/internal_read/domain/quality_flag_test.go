package domain

import "testing"

func TestAllQualityFlagsAreNonEmpty(t *testing.T) {
	flags := []QualityFlag{
		QualityComplete,
		QualityMissingProduct,
		QualityMissingStock,
		QualityMissingCost,
		QualityMissingTax,
		QualityAmbiguousProduct,
		QualityStaleSource,
	}

	for _, flag := range flags {
		if flag == "" {
			t.Fatal("expected every quality flag to have a stable string value")
		}
	}
}
