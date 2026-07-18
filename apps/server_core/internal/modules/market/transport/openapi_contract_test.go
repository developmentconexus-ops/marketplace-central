package transport

import (
	"os"
	"strings"
	"testing"
)

// TestMarketCollectVerdictOpenAPIContractParity mirrors
// composition/root_test.go's TestRefreshListingsOpenAPIContractParity idiom
// (string-slicing, no yaml/spectral tool in repo) for the F-04-S5 additive
// /market/* paths. It lives inside transport (not composition) per the
// ratified OpenAPI gate (D-31 Option B): composition/root_test.go stays
// untouched.
func TestMarketCollectVerdictOpenAPIContractParity(t *testing.T) {
	raw, err := os.ReadFile("../../../../../../contracts/api/marketplace-central.openapi.yaml")
	if err != nil {
		t.Fatalf("read openapi spec: %v", err)
	}
	spec := string(raw)

	assertPathBlock(t, spec, "  /market/collections:", "\n  /", []string{
		"operationId: collectMarketPriceIntel",
		"$ref: '#/components/schemas/MarketCollectionRequest'",
		"$ref: '#/components/schemas/MarketCollectionResponse'",
		`"200":`,
		`"404":`,
		`"409":`,
	}, nil)

	assertPathBlock(t, spec, "  /market/signals:", "\n  /", []string{
		"operationId: listMarketSignals",
		"$ref: '#/components/schemas/MarketCompetitiveSignal'",
		`"200":`,
		`"422":`,
	}, nil)

	assertPathBlock(t, spec, "  /market/aggregates:", "\n  /", []string{
		"operationId: listMarketAggregates",
		"$ref: '#/components/schemas/MarketAggregate'",
		`"200":`,
		`"422":`,
	}, nil)

	assertPathBlock(t, spec, "  /market/verdicts:", "\n  /", []string{
		"operationId: listMarketVerdicts",
		"$ref: '#/components/schemas/MarketVerdict'",
		`"200":`,
		`"422":`,
	}, nil)

	for _, schema := range []string{
		"    MarketCollectionRequest:",
		"    MarketCollectionResponse:",
		"    MarketCollectionDecision:",
		"    MarketCollectionCause:",
		"    MarketCompetitiveSignal:",
		"    MarketAggregate:",
		"    MarketVerdict:",
		"    MarketVerdictInputRef:",
		"    MarketVerdictMarketRange:",
	} {
		if !strings.Contains(spec, schema) {
			t.Fatalf("spec missing schema %s", schema)
		}
	}

	// Untouched confirmation: the pre-existing /market/observations and
	// /market/references blocks must still be present and unmodified by
	// this additive change.
	assertPathBlock(t, spec, "  /market/observations:", "\n  /", []string{
		"operationId: listMarketObservations",
	}, nil)
	assertPathBlock(t, spec, "  /market/references:", "\n  /", []string{
		"operationId: listMarketReferences",
	}, nil)
}

func assertPathBlock(t *testing.T, spec, startMarker, endMarker string, want, unwanted []string) {
	t.Helper()
	start := strings.Index(spec, startMarker)
	if start < 0 {
		t.Fatalf("spec missing path block %s", startMarker)
	}
	relEnd := strings.Index(spec[start+len(startMarker):], endMarker)
	end := len(spec)
	if relEnd >= 0 {
		end = start + len(startMarker) + relEnd
	}
	block := spec[start:end]
	for _, w := range want {
		if !strings.Contains(block, w) {
			t.Fatalf("path block %s missing %q\nblock:\n%s", startMarker, w, block)
		}
	}
	for _, u := range unwanted {
		if strings.Contains(block, u) {
			t.Fatalf("path block %s unexpectedly contains %q", startMarker, u)
		}
	}
}
