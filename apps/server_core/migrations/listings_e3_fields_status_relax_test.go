package migrations

import (
	"strings"
	"testing"
)

const listingsE3Migration = "0090_listings_e3_fields_status_relax.sql"

// TestListingsE3FieldsAreAdditiveAndNullable proves every column IC-07 adds to
// listings arrives as a nullable, no-op-on-rerun ADD COLUMN — honest-unknown
// (ADR-17) until M-04's writer (IngestListing) populates it.
func TestListingsE3FieldsAreAdditiveAndNullable(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, listingsE3Migration))
	lowerSQL := strings.ToLower(sql)

	for _, declaration := range []string{
		"alter table listings add column if not exists sold_quantity int",
		"alter table listings add column if not exists category_id text",
		"alter table listings add column if not exists condition text",
		"alter table listings add column if not exists permalink text",
		"alter table listings add column if not exists thumbnail text",
		"alter table listings add column if not exists date_created_ml timestamptz",
		"alter table listings add column if not exists tags text[]",
		"alter table listings add column if not exists catalog_product_id text",
		"alter table listings add column if not exists shipping_mode text",
		"alter table listings add column if not exists free_shipping boolean",
		"alter table listings add column if not exists logistic_type text",
		"alter table listings add column if not exists available_quantity int",
		"alter table listings add column if not exists last_seen_at timestamptz",
		"alter table listings add column if not exists absent_since timestamptz",
		"alter table listings add column if not exists raw jsonb",
		"alter table listings add column if not exists raw_truncated boolean",
	} {
		if !strings.Contains(lowerSQL, declaration) {
			t.Errorf("listings E3 migration missing exact declaration %q", declaration)
		}
	}

	// Every ADD COLUMN in this migration must be bare (nullable, no DEFAULT) —
	// NULL is the honest-unknown state until IngestListing backfills it.
	added := []string{
		"sold_quantity", "category_id", "condition", "permalink", "thumbnail",
		"date_created_ml", "tags", "catalog_product_id", "shipping_mode",
		"free_shipping", "logistic_type", "available_quantity", "last_seen_at",
		"absent_since", "raw", "raw_truncated",
	}
	for _, column := range added {
		for _, suffix := range []string{
			column + " int not null", column + " text not null",
			column + " timestamptz not null", column + " text[] not null",
			column + " boolean not null", column + " jsonb not null",
			column + " default",
		} {
			if strings.Contains(lowerSQL, suffix) {
				t.Errorf("%s must be nullable with no DEFAULT (honest-unknown, ADR-17)", column)
			}
		}
	}

	// No destructive ALTER anywhere except the one sanctioned DROP CONSTRAINT
	// below — additive columns only.
	for _, forbidden := range []string{"alter column", "drop column", "drop table"} {
		if strings.Contains(lowerSQL, forbidden) {
			t.Errorf("listings E3 migration must be additive only, found %q", forbidden)
		}
	}

	// IC-07 decision (planning design §7 item 6): fee data lives only in
	// channel_fees. listings must never grow these columns. Matched as an
	// actual ADD COLUMN, not bare substring — this file's own rationale
	// comment names them in prose on purpose.
	for _, forbidden := range []string{"commission_amount", "commission_pct", "free_shipping_cost"} {
		if strings.Contains(lowerSQL, "add column if not exists "+forbidden) {
			t.Errorf("listings E3 migration must not add fee column %q — fee data lives only in channel_fees (IC-01)", forbidden)
		}
	}
}

// TestListingsStatusCheckIsDropped proves IC-07's "status is verbatim from
// the provider, no CHECK" decision: the CHECK this repo wrote in 0036/0037 no
// longer constrains the column. Idempotent-safe (IF EXISTS) so a re-run is a
// no-op whether or not the constraint is still present.
func TestListingsStatusCheckIsDropped(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, listingsE3Migration))
	if !strings.Contains(strings.ToLower(sql), "alter table listings drop constraint if exists listings_status_check") {
		t.Fatal("listings E3 migration must drop listings_status_check (IF EXISTS, for re-run safety)")
	}
	// Guard against re-adding a narrower list in the same breath — IC-07 is
	// "no CHECK", not "a wider CHECK".
	if strings.Contains(strings.ToLower(sql), "add constraint listings_status_check") {
		t.Error("listings E3 migration must not re-add a status CHECK — IC-07 requires no CHECK at all")
	}
}
