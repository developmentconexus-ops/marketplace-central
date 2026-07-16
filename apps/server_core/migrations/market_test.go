package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

func TestMarketMigrationsMatchIC04(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		table    string
		columns  []string
		nullable []string
		notNull  []string
		required []string
	}{
		{
			name:  "observations",
			file:  "0043_market_observations.sql",
			table: "market_observations",
			columns: []string{
				"tenant_id", "listing_id", "our_sale_price_amount", "our_sale_price_currency",
				"competitive_status", "winner_price_amount", "winner_price_currency",
				"competitive_target_amount", "competitive_target_currency",
				"catalog_offer_price_amount", "catalog_offer_price_currency", "catalog_stats",
				"evidence_state", "captured_at", "source", "created_at",
			},
			nullable: []string{
				"our_sale_price_amount", "our_sale_price_currency", "competitive_status",
				"winner_price_amount", "winner_price_currency", "competitive_target_amount",
				"competitive_target_currency", "catalog_offer_price_amount",
				"catalog_offer_price_currency", "catalog_stats",
			},
			notNull: []string{
				"tenant_id", "listing_id", "evidence_state", "captured_at", "source", "created_at",
			},
			required: []string{
				"PRIMARY KEY (tenant_id, listing_id, captured_at)",
				"CONSTRAINT market_observations_our_sale_price_pair_check CHECK ((our_sale_price_amount IS NULL AND our_sale_price_currency IS NULL) OR (our_sale_price_amount IS NOT NULL AND our_sale_price_currency IS NOT NULL AND our_sale_price_currency = 'BRL'))",
				"CONSTRAINT market_observations_winner_price_pair_check CHECK ((winner_price_amount IS NULL AND winner_price_currency IS NULL) OR (winner_price_amount IS NOT NULL AND winner_price_currency IS NOT NULL AND winner_price_currency = 'BRL'))",
				"CONSTRAINT market_observations_competitive_target_pair_check CHECK ((competitive_target_amount IS NULL AND competitive_target_currency IS NULL) OR (competitive_target_amount IS NOT NULL AND competitive_target_currency IS NOT NULL AND competitive_target_currency = 'BRL'))",
				"CONSTRAINT market_observations_catalog_offer_price_pair_check CHECK ((catalog_offer_price_amount IS NULL AND catalog_offer_price_currency IS NULL) OR (catalog_offer_price_amount IS NOT NULL AND catalog_offer_price_currency IS NOT NULL AND catalog_offer_price_currency = 'BRL'))",
				"CONSTRAINT market_observations_competitive_status_check CHECK (competitive_status IN ('winning', 'competing', 'not_listed'))",
				"CONSTRAINT market_observations_evidence_state_check CHECK (evidence_state IN ('observed', 'insufficient_market', 'no_price_evidence'))",
				"CONSTRAINT market_observations_source_check CHECK (source IN ('official_api', 'vendor', 'manual'))",
			},
		},
		{
			name:  "references",
			file:  "0044_market_references.sql",
			table: "market_references",
			columns: []string{
				"tenant_id", "product_id", "catalog_product_id", "match_state", "match_method",
				"catalog_stats", "evidence_state", "captured_at", "source", "created_at",
			},
			nullable: []string{"catalog_product_id", "match_method", "catalog_stats"},
			notNull: []string{
				"tenant_id", "product_id", "match_state", "evidence_state", "captured_at", "source", "created_at",
			},
			required: []string{
				"PRIMARY KEY (tenant_id, product_id, captured_at)",
				"CONSTRAINT market_references_match_state_check CHECK (match_state IN ('accept', 'review', 'reject', 'no_candidate'))",
				"CONSTRAINT market_references_evidence_state_check CHECK (evidence_state IN ('observed', 'insufficient_market', 'no_price_evidence'))",
				"CONSTRAINT market_references_source_check CHECK (source IN ('official_api', 'vendor', 'manual'))",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := fs.ReadFile(Source(), tt.file)
			if err != nil {
				t.Fatal(err)
			}
			sql := string(raw)
			normalizedSQL := normalizeMigrationSQL(sql)
			body := createTableBody(t, sql, tt.table)

			assertExactColumns(t, body, tt.columns)
			for _, required := range tt.required {
				if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
					t.Errorf("%s migration missing %q", tt.table, required)
				}
			}
			for _, column := range tt.nullable {
				assertNullableColumn(t, body, column)
			}
			for _, column := range tt.notNull {
				assertNotNullableColumn(t, body, column)
			}
			assertOnlyLifecycleDefaults(t, body)

			if !strings.Contains(sql, "IC-04") || !strings.Contains(strings.ToLower(sql), "append-only") || !strings.Contains(strings.ToLower(sql), "contract-only") {
				t.Error("market migration must identify the IC-04 append-only contract-only intent")
			}
			if strings.Contains(strings.ToUpper(sql), "INSERT") {
				t.Error("market migration must not contain INSERT statements")
			}
			if strings.Contains(strings.ToLower(sql), "market_price") {
				t.Error("market migration must not contain generic market_price")
			}
			if strings.Contains(strings.ToLower(sql), "foreign key") {
				t.Error("market migration must not add a foreign key")
			}
		})
	}
}
