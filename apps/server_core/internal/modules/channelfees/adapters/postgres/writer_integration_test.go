//go:build integration

package postgres

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/channelfees/domain"
	"marketplace-central/apps/server_core/internal/platform/migrate"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
	canonical "marketplace-central/apps/server_core/migrations"
)

// TestUpsertFeeRoundTripCollapsesToOneRowWithColetadoEmAdvanced is the
// feature.md-required round-trip proof: UpsertFee twice for the same
// subject must collapse to exactly 1 row, with coletado_em advanced —
// proven by a real query, not by trusting the upsert SQL.
func TestUpsertFeeRoundTripCollapsesToOneRowWithColetadoEmAdvanced(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_roundtrip")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM channel_fees WHERE tenant_id = $1`, "tenant_f02_roundtrip")
	})

	writer := NewWriter(pool)
	first := time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC)
	second := time.Date(2026, 7, 31, 15, 30, 0, 0, time.UTC)
	currency := "BRL"

	fee := domain.Fee{
		TenantID:       "tenant_f02_roundtrip",
		Provider:       "mercado_livre",
		InstallationID: "inst1",
		Layer:          domain.LayerListing,
		FeeKind:        domain.FeeKindCommission,
		SubjectType:    domain.SubjectTypeListing,
		SubjectID:      "MLB123",
		ValueType:      domain.ValueTypeAmount,
		Value:          "15.99",
		Currency:       &currency,
		Detail:         json.RawMessage(`{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,"price_used":79.90,"listing_type_id":"gold_special"}`),
		Origem:         domain.OrigemAPIListingPrices,
		ColetadoEm:     first,
	}
	if err := writer.UpsertFee(ctx, fee); err != nil {
		t.Fatalf("first UpsertFee: %v", err)
	}

	fee.Value = "16.50"
	fee.ColetadoEm = second
	if err := writer.UpsertFee(ctx, fee); err != nil {
		t.Fatalf("second UpsertFee: %v", err)
	}

	var count int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM channel_fees
		WHERE tenant_id=$1 AND provider=$2 AND installation_id=$3
		  AND subject_type=$4 AND subject_id=$5 AND layer=$6 AND fee_kind=$7
	`, fee.TenantID, fee.Provider, fee.InstallationID, fee.SubjectType, fee.SubjectID, fee.Layer, fee.FeeKind).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 row after 2 upserts on the same natural key, got %d", count)
	}

	var value string
	var coletadoEm time.Time
	if err := pool.QueryRow(ctx, `
		SELECT value, coletado_em FROM channel_fees
		WHERE tenant_id=$1 AND provider=$2 AND installation_id=$3
		  AND subject_type=$4 AND subject_id=$5 AND layer=$6 AND fee_kind=$7
	`, fee.TenantID, fee.Provider, fee.InstallationID, fee.SubjectType, fee.SubjectID, fee.Layer, fee.FeeKind).Scan(&value, &coletadoEm); err != nil {
		t.Fatalf("select row: %v", err)
	}
	if value != "16.5000" { // numeric(14,4) round-trips zero-padded to 4 decimals
		t.Fatalf("expected value to have been updated to the 2nd upsert's value 16.50, got %q", value)
	}
	if !coletadoEm.Truncate(time.Microsecond).Equal(second.Truncate(time.Microsecond)) {
		t.Fatalf("expected coletado_em to have advanced to %v, got %v", second, coletadoEm)
	}
}

// TestUpsertFeeLayer3FreightWithoutDetailPersists proves (real DB round
// trip, not just validation-in-isolation) that a layer 3 fee_kind=freight
// row with NULL detail is accepted end to end — the asymmetric other half
// of the commission mandate (IC-01; auditoria P5 r06 F-r06-1).
func TestUpsertFeeLayer3FreightWithoutDetailPersists(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_freight")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM channel_fees WHERE tenant_id = $1`, "tenant_f02_freight")
	})

	writer := NewWriter(pool)
	currency := "BRL"
	fee := domain.Fee{
		TenantID:       "tenant_f02_freight",
		Provider:       "mercado_livre",
		InstallationID: "inst1",
		Layer:          domain.LayerOrder,
		FeeKind:        domain.FeeKindFreight,
		SubjectType:    domain.SubjectTypeOrder,
		SubjectID:      "200012345",
		ValueType:      domain.ValueTypeAmount,
		Value:          "9.90",
		Currency:       &currency,
		Detail:         nil,
		Origem:         domain.OrigemAPIShipment,
		ColetadoEm:     time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC),
	}
	if err := writer.UpsertFee(ctx, fee); err != nil {
		t.Fatalf("expected layer 3 freight without detail to persist, got: %v", err)
	}

	var detail *string
	if err := pool.QueryRow(ctx, `
		SELECT detail::text FROM channel_fees
		WHERE tenant_id=$1 AND subject_type=$2 AND subject_id=$3 AND layer=$4 AND fee_kind=$5
	`, fee.TenantID, fee.SubjectType, fee.SubjectID, fee.Layer, fee.FeeKind).Scan(&detail); err != nil {
		t.Fatalf("select row: %v", err)
	}
	if detail != nil {
		t.Fatalf("expected detail to be NULL for freight, got %q", *detail)
	}
}

// TestUpsertFeeLayer3CommissionWithoutDetailRejectedBeforeAnyWrite proves the
// must-fail path end to end: the row must never reach the table.
func TestUpsertFeeLayer3CommissionWithoutDetailRejectedBeforeAnyWrite(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_reject")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM channel_fees WHERE tenant_id = $1`, "tenant_f02_reject")
	})

	writer := NewWriter(pool)
	currency := "BRL"
	fee := domain.Fee{
		TenantID:       "tenant_f02_reject",
		Provider:       "mercado_livre",
		InstallationID: "inst1",
		Layer:          domain.LayerOrder,
		FeeKind:        domain.FeeKindCommission,
		SubjectType:    domain.SubjectTypeOrderLine,
		SubjectID:      "200099999:MLB999",
		ValueType:      domain.ValueTypeAmount,
		Value:          "24.33",
		Currency:       &currency,
		Detail:         nil,
		Origem:         domain.OrigemAPIOrder,
		ColetadoEm:     time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC),
	}
	err := writer.UpsertFee(ctx, fee)
	if err == nil {
		t.Fatal("MUST-FAIL: expected layer 3 commission without detail to be rejected, got nil error")
	}
	t.Logf("RED (rejected, as required): %v", err)

	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM channel_fees WHERE tenant_id=$1`, fee.TenantID).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected zero rows written after rejection, got %d", count)
	}
	t.Log("GREEN (no row written, as required)")
}
