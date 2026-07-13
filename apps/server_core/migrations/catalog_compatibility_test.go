package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

func TestCatalogCompatibilityMigrationPreservesEvidenceAndReadbackStates(t *testing.T) {
	raw, err := fs.ReadFile(Source(), "0035_catalog_codprod_compatibility.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(raw)
	for _, required := range []string{
		"catalog_codprod_compatibility_candidates",
		"catalog_legacy_product_reference_evidence",
		"source text NOT NULL CHECK (source = 'oracle/sankhya')",
		"read_only boolean NOT NULL CHECK (read_only)",
		"e.legacy_product_id = c.internal_product_id::text",
		"c.internal_product_id > 0",
		"WHEN r.candidate_count = 1 THEN 'mapped'",
		"WHEN r.candidate_count = 0 THEN 'not_found'",
		"ELSE 'identity_conflict'",
		"resolution_status = 'mapped' AND internal_product_id > 0",
	} {
		if !strings.Contains(sql, required) {
			t.Fatalf("compatibility migration missing %q", required)
		}
	}
	for _, consumer := range []string{"product_enrichments", "classification_products", "pricing_simulations"} {
		if !strings.Contains(sql, "UPDATE "+consumer) {
			t.Fatalf("compatibility migration has no %s readback update", consumer)
		}
	}
}
