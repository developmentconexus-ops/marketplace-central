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

// TestResolveListingFeesPrefersLayer2OverAbsent proves the ledger-only
// ladder's first rung: when a layer 2 row exists, ResolveListingFees returns
// it verbatim — value, value_type, currency, layer:2, detail (JSONB
// verbatim), origem, coletado_em (IC-01 Required Outputs) — never the
// config, never synthesized.
func TestResolveListingFeesPrefersLayer2OverAbsent(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_resolve_l2")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM channel_fees WHERE tenant_id = $1`, "tenant_f02_resolve_l2")
	})

	writer := NewWriter(pool)
	currency := "BRL"
	coletadoEm := time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC)
	detail := json.RawMessage(`{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,"price_used":79.90,"listing_type_id":"gold_special"}`)
	fee := domain.Fee{
		TenantID:       "tenant_f02_resolve_l2",
		Provider:       "mercado_livre",
		InstallationID: "inst1",
		Layer:          domain.LayerListing,
		FeeKind:        domain.FeeKindCommission,
		SubjectType:    domain.SubjectTypeListing,
		SubjectID:      "MLB123",
		ValueType:      domain.ValueTypeAmount,
		Value:          "15.99",
		Currency:       &currency,
		Detail:         detail,
		Origem:         domain.OrigemAPIListingPrices,
		ColetadoEm:     coletadoEm,
	}
	if err := writer.UpsertFee(ctx, fee); err != nil {
		t.Fatalf("seed layer 2 commission: %v", err)
	}

	reader := NewReader(pool)
	resolved, err := reader.ResolveListingFees(ctx, "tenant_f02_resolve_l2", "mercado_livre", "inst1", "MLB123")
	if err != nil {
		t.Fatalf("ResolveListingFees: %v", err)
	}

	c := resolved.Commission
	if !c.Found {
		t.Fatal("expected commission to be found at layer 2")
	}
	if c.Layer != domain.LayerListing {
		t.Fatalf("expected layer 2, got %d", c.Layer)
	}
	if c.Value != "15.9900" { // numeric(14,4) round-trips zero-padded to 4 decimals
		t.Fatalf("expected value 15.99, got %q", c.Value)
	}
	if c.ValueType != domain.ValueTypeAmount {
		t.Fatalf("expected value_type amount, got %q", c.ValueType)
	}
	if c.Currency == nil || *c.Currency != "BRL" {
		t.Fatalf("expected currency BRL, got %v", c.Currency)
	}
	if c.Origem != domain.OrigemAPIListingPrices {
		t.Fatalf("expected origem api_listing_prices, got %q", c.Origem)
	}
	if !c.ColetadoEm.Truncate(time.Microsecond).Equal(coletadoEm.Truncate(time.Microsecond)) {
		t.Fatalf("expected coletado_em %v, got %v", coletadoEm, c.ColetadoEm)
	}

	var decodedDetail map[string]any
	if err := json.Unmarshal(c.Detail, &decodedDetail); err != nil {
		t.Fatalf("decode returned detail: %v", err)
	}
	if decodedDetail["percentage_fee"] != 12.5 {
		t.Fatalf("expected detail to be the VERBATIM row jsonb (percentage_fee=12.5), got %v", decodedDetail)
	}

	// Freight was never written for this listing: honest-unknown, never 0.
	if resolved.Freight.Found {
		t.Fatalf("expected freight to be absent (typed-absent), got %+v", resolved.Freight)
	}
}

// TestResolveListingFeesAbsentWhenNoLedgerRowExists proves the ladder's
// terminal rung: with zero rows at any layer, the reader returns a typed
// absent result (Found=false) — never a zero amount, never a default. The
// config fallback is out of scope for this port (M-07 composes it).
func TestResolveListingFeesAbsentWhenNoLedgerRowExists(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_resolve_absent")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	reader := NewReader(pool)
	resolved, err := reader.ResolveListingFees(ctx, "tenant_f02_resolve_absent", "mercado_livre", "inst1", "MLB_NOFEE")
	if err != nil {
		t.Fatalf("ResolveListingFees: %v", err)
	}

	if resolved.Commission.Found {
		t.Fatalf("expected commission absent, got %+v", resolved.Commission)
	}
	if resolved.Commission.Value != "" {
		t.Fatalf("expected zero-value Value on absent result, got %q — must never be a fabricated amount", resolved.Commission.Value)
	}
	if resolved.Freight.Found {
		t.Fatalf("expected freight absent, got %+v", resolved.Freight)
	}
}

// TestResolveListingFeesIgnoresLayer3 proves the ladder never falls through
// to layer 3, and does so under a key the layer filter is actually forced to
// discriminate on: a layer 2 row AND a layer 3 row are seeded under the
// IDENTICAL (tenant_id, provider, installation_id, subject_type=listing,
// subject_id, fee_kind) natural key — the natural key's UNIQUE constraint
// allows this because `layer` is itself part of that key. If the reader's
// `layer = ANY($7)` filter were ever deleted, `ORDER BY layer DESC LIMIT 1`
// alone would return the layer 3 row (the higher layer) instead of layer 2 —
// so this test can only stay green because the layer filter is doing real
// work, not because the two rows happen to live under different keys.
func TestResolveListingFeesIgnoresLayer3(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_f02_resolve_l3ignore")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM channel_fees WHERE tenant_id = $1`, "tenant_f02_resolve_l3ignore")
	})

	writer := NewWriter(pool)
	currency := "BRL"
	const (
		sameProvider  = "mercado_livre"
		sameInstall   = "inst1"
		sameSubjectID = "MLB_L3IGNORE"
	)

	layer2 := domain.Fee{
		TenantID:       "tenant_f02_resolve_l3ignore",
		Provider:       sameProvider,
		InstallationID: sameInstall,
		Layer:          domain.LayerListing,
		FeeKind:        domain.FeeKindCommission,
		SubjectType:    domain.SubjectTypeListing,
		SubjectID:      sameSubjectID,
		ValueType:      domain.ValueTypeAmount,
		Value:          "15.99",
		Currency:       &currency,
		Detail:         json.RawMessage(`{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,"price_used":79.90,"listing_type_id":"gold_special"}`),
		Origem:         domain.OrigemAPIListingPrices,
		ColetadoEm:     time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC),
	}
	if err := writer.UpsertFee(ctx, layer2); err != nil {
		t.Fatalf("seed layer 2 commission: %v", err)
	}

	// Same tenant/provider/installation/subject_type/subject_id/fee_kind as
	// layer2 above — only `layer` differs, which is exactly why the natural
	// key UNIQUE constraint permits both rows to coexist.
	layer3 := layer2
	layer3.Layer = domain.LayerOrder
	layer3.Value = "99.99" // deliberately far from 15.99 so a wrong pick is unmissable
	layer3.Detail = json.RawMessage(`{"sale_fee_unit":8.11,"quantity":3}`)
	layer3.Origem = domain.OrigemAPIOrder
	if err := writer.UpsertFee(ctx, layer3); err != nil {
		t.Fatalf("seed layer 3 commission (same natural key modulo layer): %v", err)
	}

	reader := NewReader(pool)
	resolved, err := reader.ResolveListingFees(ctx, "tenant_f02_resolve_l3ignore", sameProvider, sameInstall, sameSubjectID)
	if err != nil {
		t.Fatalf("ResolveListingFees: %v", err)
	}
	if !resolved.Commission.Found {
		t.Fatal("expected commission found at layer 2, even with a layer-3 row present at the same key")
	}
	if resolved.Commission.Layer != domain.LayerListing {
		t.Fatalf("expected layer 2 to win over layer 3 at the identical key, got layer %d", resolved.Commission.Layer)
	}
	if resolved.Commission.Value != "15.9900" {
		t.Fatalf("expected the layer-2 value 15.9900 (not the layer-3 value 99.9900), got %q", resolved.Commission.Value)
	}
}
