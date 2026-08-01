package migrations

import (
	"strings"
	"testing"
)

const (
	listingVariationsMigration      = "0091_listing_variations.sql"
	listingVariationsIndexMigration = "0092_listing_variations_tenant_provider_listing_index.sql"
)

func TestListingVariationsMigrationMatchesIC07(t *testing.T) {
	sql := readMigration(t, listingVariationsMigration)
	normalizedSQL := normalizeMigrationSQL(sql)

	body := createTableBody(t, sql, "listing_variations")
	assertExactColumns(t, body, []string{
		"tenant_id", "installation_id", "provider", "provider_listing_id", "variation_id",
		"price", "available_quantity", "sold_quantity", "seller_sku", "attributes",
		"last_seen_at", "absent_since", "created_at", "updated_at",
	})

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS listing_variations",
		"price numeric(14,2)",
		"attributes jsonb",
		"PRIMARY KEY (tenant_id, installation_id, provider, provider_listing_id, variation_id)",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("listing_variations migration missing %q", required)
		}
	}

	for _, column := range []string{
		"tenant_id", "installation_id", "provider", "provider_listing_id", "variation_id",
		"created_at", "updated_at",
	} {
		assertNotNullableColumn(t, body, column)
	}
	for _, column := range []string{
		"price", "available_quantity", "sold_quantity", "seller_sku", "attributes",
		"last_seen_at", "absent_since",
	} {
		assertNullableColumn(t, body, column)
	}
	assertOnlyLifecycleDefaults(t, body)

	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("listing_variations migration must not add a foreign key")
	}
	// No raw column of its own — the parent listing's raw (0090) covers the
	// full /items payload; a variation has no independent raw fetch.
	if strings.Contains(strings.ToLower(body), "raw") {
		t.Error("listing_variations must not carry its own raw column")
	}
}

func TestListingVariationsIndexMatchesIC07(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, listingVariationsIndexMigration))
	want := "CREATE INDEX IF NOT EXISTS idx_listing_variations_tenant_provider_listing ON listing_variations (tenant_id, provider_listing_id)"
	if !strings.Contains(sql, normalizeMigrationSQL(want)) {
		t.Fatalf("listing_variations index migration missing %q", want)
	}
}
