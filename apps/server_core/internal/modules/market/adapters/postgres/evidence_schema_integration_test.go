//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestEvidenceSchemaShape(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "market_evidence_schema")

	for table, columns := range map[string]map[string]string{
		"market_price_snapshots": {
			"tenant_id": "NO", "scope": "NO", "ref_id": "NO", "source": "NO",
			"price_amount": "YES", "price_currency": "YES", "observed_at": "NO",
			"fetched_at": "NO", "expires_at": "NO", "status": "NO",
			"failure_reason": "YES", "request_id": "NO",
		},
		"market_validated_offers": {
			"tenant_id": "NO", "catalog_product_id": "NO", "seller_id": "NO",
			"price_amount": "YES", "price_currency": "YES", "condition": "NO",
			"match_status": "NO", "fetched_at": "NO",
		},
		"market_match_decisions": {
			"tenant_id": "NO", "product_id": "NO", "catalog_product_id": "YES",
			"match_status": "NO", "anchors": "NO", "contradictions": "NO",
			"confidence_band": "NO", "decided_at": "NO",
		},
		"market_competitive_signals": {
			"tenant_id": "NO", "listing_id": "NO", "our_price_amount": "NO",
			"our_price_currency": "NO", "winner_price_amount": "YES",
			"winner_price_currency": "YES", "target_price_amount": "YES",
			"target_price_currency": "YES", "position_rank": "YES",
			"position_total": "YES", "fetched_at": "NO",
		},
		"market_aggregates": {
			"tenant_id": "NO", "product_id": "NO", "median_amount": "YES",
			"median_currency": "YES", "min_valid_amount": "YES",
			"min_valid_currency": "YES", "n_offers": "NO", "n_sellers": "NO",
			"source": "NO", "fetched_at": "NO", "computed_at": "NO", "status": "NO",
		},
	} {
		assertColumns(t, ctx, pool, table, columns)
	}

	for _, constraint := range []string{
		"market_price_snapshots_pkey",
		"market_price_snapshots_scope_check", "market_price_snapshots_source_check",
		"market_price_snapshots_price_pair_check", "market_price_snapshots_status_check",
		"market_price_snapshots_evidence_check",
		"market_validated_offers_pkey", "market_validated_offers_price_pair_check",
		"market_validated_offers_match_status_check", "market_match_decisions_pkey",
		"market_match_decisions_match_status_check", "market_match_decisions_confidence_band_check",
		"market_competitive_signals_pkey", "market_competitive_signals_our_price_check",
		"market_competitive_signals_winner_price_pair_check", "market_competitive_signals_target_price_pair_check",
		"market_competitive_signals_position_pair_check", "market_aggregates_pkey",
		"market_aggregates_median_pair_check", "market_aggregates_min_valid_pair_check",
		"market_aggregates_counts_check", "market_aggregates_source_check", "market_aggregates_status_check",
	} {
		assertConstraint(t, ctx, pool, constraint)
	}

	for _, index := range []string{
		"market_price_snapshots_latest_valid_idx", "market_validated_offers_catalog_idx",
		"market_match_decisions_product_idx", "market_competitive_signals_listing_idx",
		"market_aggregates_product_idx",
	} {
		assertIndexTenantFirst(t, ctx, pool, index)
	}

	assertPrimaryKeyColumns(t, ctx, pool, "market_price_snapshots", []string{"tenant_id", "scope", "ref_id", "source", "fetched_at"})
}

func assertColumns(t *testing.T, ctx context.Context, pool *pgxpool.Pool, table string, want map[string]string) {
	t.Helper()
	rows, err := pool.Query(ctx, `SELECT column_name, is_nullable FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1`, table)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var name, nullable string
		if err := rows.Scan(&name, &nullable); err != nil {
			t.Fatal(err)
		}
		got[name] = nullable
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for name, nullable := range want {
		if got[name] != nullable {
			t.Errorf("%s.%s nullability = %q, want %q", table, name, got[name], nullable)
		}
	}
}

func assertConstraint(t *testing.T, ctx context.Context, pool *pgxpool.Pool, name string) {
	t.Helper()
	var exists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = $1)`, name).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("constraint %s does not exist", name)
	}
}

func assertPrimaryKeyColumns(t *testing.T, ctx context.Context, pool *pgxpool.Pool, table string, want []string) {
	t.Helper()
	rows, err := pool.Query(ctx, `
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass AND i.indisprimary
		ORDER BY array_position(i.indkey, a.attnum)
	`, table)
	if err != nil {
		t.Fatalf("read primary key columns for %s: %v", table, err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		got = append(got, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Errorf("%s primary key columns = %v, want %v", table, got, want)
		return
	}
	for i, name := range want {
		if got[i] != name {
			t.Errorf("%s primary key columns = %v, want %v", table, got, want)
			return
		}
	}
}

func assertIndexTenantFirst(t *testing.T, ctx context.Context, pool *pgxpool.Pool, name string) {
	t.Helper()
	var firstColumn string
	if err := pool.QueryRow(ctx, `
		SELECT attribute.attname
		FROM pg_class index_class
		JOIN pg_index index_info ON index_info.indexrelid = index_class.oid
		JOIN pg_attribute attribute ON attribute.attrelid = index_info.indrelid AND attribute.attnum = index_info.indkey[0]
		WHERE index_class.relname = $1
	`, name).Scan(&firstColumn); err != nil {
		t.Fatalf("read index %s: %v", name, err)
	}
	if firstColumn != "tenant_id" {
		t.Errorf("index %s first column = %q, want tenant_id", name, firstColumn)
	}
}
