//go:build integration

// Must-fail proof for MIS-007 M-02 F-01-core-ddl, kept as a permanent
// regression guard (feature.md Validation Expectations): constraint
// rejection must be proven by SQLSTATE with the named constraint, and the
// partial-unique upsert must be proven in both directions — reject while a
// row for the same (tenant, provider, entity_type, entity_id, kind) is still
// open, accept once it is resolved.
package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
	canonical "marketplace-central/apps/server_core/migrations"
)

func TestF01ChannelFeesCurrencyCheckRejectsAmountWithoutCurrency(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_mis007_f01")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO channel_fees
			(tenant_id, provider, installation_id, layer, fee_kind, subject_type, subject_id,
			 value_type, value, currency, origem, coletado_em)
		VALUES
			('tenant_mis007_f01', 'mercado_livre', 'inst1', 2, 'commission', 'listing', 'MLB1',
			 'amount', 15.99, NULL, 'api_listing_prices', now())
	`)
	if err == nil {
		t.Fatal("MUST-FAIL: expected INSERT with value_type='amount' AND currency=NULL to be rejected, got no error")
	}
	t.Logf("RED (rejected, as required): %v", err)

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected a *pgconn.PgError, got %T: %v", err, err)
	}
	if pgErr.Code != "23514" {
		t.Fatalf("expected SQLSTATE 23514 (check_violation), got %s", pgErr.Code)
	}
	if pgErr.ConstraintName != "channel_fees_currency_when_amount_check" {
		t.Fatalf("expected constraint channel_fees_currency_when_amount_check named in the error, got %q", pgErr.ConstraintName)
	}
	t.Logf("GREEN (constraint identified, as required): SQLSTATE=%s constraint=%s", pgErr.Code, pgErr.ConstraintName)
}

func TestF01DivergencesOneOpenRowRejectsSecondOpenThenAcceptsAfterResolve(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_mis007_f01")
	ctx := context.Background()
	if _, err := migrate.Run(ctx, pool, canonical.Source()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	insert := func() error {
		_, err := pool.Exec(ctx, `
			INSERT INTO divergences
				(tenant_id, provider, entity_type, entity_id, kind,
				 expected_value, observed_value, expected_source, observed_source,
				 expected_observed_at, observed_observed_at, detected_at, last_evaluated_at)
			VALUES
				('tenant_mis007_f01', 'mercado_livre', 'listing', 'MLB123mustfail', 'estoque',
				 12, 9, 'mirror:sankhya:vendavel', 'api_listing:available_quantity',
				 now(), now(), now(), now())
		`)
		return err
	}

	if err := insert(); err != nil {
		t.Fatalf("first open row insert should succeed: %v", err)
	}

	err := insert()
	if err == nil {
		t.Fatal("MUST-FAIL: expected a second open row for the same (tenant, provider, entity_type, entity_id, kind) to be rejected, got no error")
	}
	t.Logf("RED (rejected, as required): %v", err)

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected a *pgconn.PgError, got %T: %v", err, err)
	}
	if pgErr.Code != "23505" {
		t.Fatalf("expected SQLSTATE 23505 (unique_violation), got %s", pgErr.Code)
	}
	if pgErr.ConstraintName != "divergences_one_open_row" {
		t.Fatalf("expected constraint divergences_one_open_row named in the error, got %q", pgErr.ConstraintName)
	}
	t.Logf("constraint identified: SQLSTATE=%s constraint=%s", pgErr.Code, pgErr.ConstraintName)

	if _, err := pool.Exec(ctx, `
		UPDATE divergences SET resolved_at = now()
		WHERE tenant_id = 'tenant_mis007_f01' AND provider = 'mercado_livre'
		  AND entity_type = 'listing' AND entity_id = 'MLB123mustfail' AND kind = 'estoque'
		  AND resolved_at IS NULL
	`); err != nil {
		t.Fatalf("resolve first open row: %v", err)
	}

	if err := insert(); err != nil {
		t.Fatalf("GREEN: expected insert to succeed once the prior open row is resolved, got: %v", err)
	}
	t.Log("GREEN (accepted after resolve, as required): new open row created for the same entity+kind")
}
