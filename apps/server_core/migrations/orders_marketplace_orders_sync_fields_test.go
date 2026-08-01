package migrations

import (
	"strings"
	"testing"
)

const ordersSyncFieldsMigration = "0089_orders_marketplace_orders_sync_fields.sql"

// TestOrdersMarketplaceOrdersSyncFieldsAreAdditiveAndNullable proves every
// column IC-03 adds to orders_marketplace_orders arrives as a nullable,
// no-op-on-rerun ADD COLUMN — honest-unknown until a later milestone's
// producer populates it (M-03 shipment-enrich, M-06 backfill+decomposition).
func TestOrdersMarketplaceOrdersSyncFieldsAreAdditiveAndNullable(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, ordersSyncFieldsMigration))

	for _, declaration := range []string{
		"alter table orders_marketplace_orders add column if not exists pack_id text",
		"alter table orders_marketplace_orders add column if not exists provider_shipment_id text",
		"alter table orders_marketplace_orders add column if not exists bucket text",
		"alter table orders_marketplace_orders add column if not exists date_last_updated_ml timestamptz",
		"alter table orders_marketplace_orders add column if not exists buyer_name text",
		"alter table orders_marketplace_orders add column if not exists buyer_doc_type text",
		"alter table orders_marketplace_orders add column if not exists buyer_doc_number text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_street text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_number text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_city text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_state_code text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_zip text",
		"alter table orders_marketplace_orders add column if not exists buyer_address_country text",
		"alter table orders_marketplace_orders add column if not exists decomposition jsonb",
		"alter table orders_marketplace_orders add column if not exists net_amount numeric(14,2)",
		"alter table orders_marketplace_orders add column if not exists margin_pct numeric(7,4)",
		"create index if not exists idx_orders_marketplace_orders_bucket on orders_marketplace_orders (tenant_id, bucket)",
	} {
		if !strings.Contains(strings.ToLower(sql), declaration) {
			t.Errorf("orders_marketplace_orders sync-fields migration missing exact declaration %q", declaration)
		}
	}

	// Every ADD COLUMN in this migration must be bare (nullable, no DEFAULT) —
	// NULL is the honest-unknown state until a producer backfills it.
	added := []string{
		"pack_id", "provider_shipment_id", "bucket", "date_last_updated_ml",
		"buyer_name", "buyer_doc_type", "buyer_doc_number", "buyer_address_street",
		"buyer_address_number", "buyer_address_city", "buyer_address_state_code",
		"buyer_address_zip", "buyer_address_country", "decomposition",
		"net_amount", "margin_pct",
	}
	lowerSQL := strings.ToLower(sql)
	for _, column := range added {
		if strings.Contains(lowerSQL, column+" text not null") ||
			strings.Contains(lowerSQL, column+" timestamptz not null") ||
			strings.Contains(lowerSQL, column+" jsonb not null") ||
			strings.Contains(lowerSQL, column+" numeric(14,2) not null") ||
			strings.Contains(lowerSQL, column+" numeric(7,4) not null") ||
			strings.Contains(lowerSQL, column+" default") {
			t.Errorf("%s must be nullable with no DEFAULT (honest-unknown, ADR-17)", column)
		}
	}

	// No destructive ALTER anywhere — additive columns only.
	for _, forbidden := range []string{"alter column", "drop column", "drop table", "alter table listings", "alter table products_mirror"} {
		if strings.Contains(lowerSQL, forbidden) {
			t.Errorf("orders sync-fields migration must be additive only, found %q", forbidden)
		}
	}

	// R-6: buyer fiscal columns are typed only — no raw billing_info payload
	// column anywhere in the ADD COLUMN statements.
	if strings.Contains(lowerSQL, "add column if not exists billing_info") ||
		strings.Contains(lowerSQL, "add column if not exists buyer_raw") {
		t.Error("buyer fiscal fields must stay typed-only, no raw billing_info payload column")
	}
}
