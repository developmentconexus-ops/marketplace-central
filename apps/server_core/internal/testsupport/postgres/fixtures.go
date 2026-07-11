package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderFixture struct {
	Code        string
	DisplayName string
}

type MarketplaceDefinitionFixture struct {
	Code        string
	DisplayName string
	FeeSource   string
}

func SeedMarketplaceDefinition(t testing.TB, pool *pgxpool.Pool, fixture MarketplaceDefinitionFixture) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO marketplace_definitions (
			marketplace_code, display_name, fee_source, capabilities, credential_schema, active, created_at
		) VALUES ($1, $2, $3, '{}'::text[], '[]'::jsonb, true, $4)
		ON CONFLICT (marketplace_code) DO NOTHING
	`, fixture.Code, fixture.DisplayName, fixture.FeeSource, time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("seed marketplace definition fixture %q: %v", fixture.Code, err)
	}
}

type InstallationFixture struct {
	TenantID       string
	InstallationID string
	Provider       ProviderFixture
}

func SeedProvider(t testing.TB, pool *pgxpool.Pool, fixture ProviderFixture) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO integration_provider_definitions (
			tenant_id, provider_code, family, display_name, auth_strategy, install_mode,
			metadata_json, declared_caps_json, is_active, created_at, updated_at
		) VALUES ('system', $1, 'marketplace', $2, 'none', 'manual', '{}'::jsonb, '[]'::jsonb, true, $3, $3)
		ON CONFLICT (provider_code) DO NOTHING
	`, fixture.Code, fixture.DisplayName, time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("seed provider fixture %q: %v", fixture.Code, err)
	}
}

func SeedInstallation(t testing.TB, pool *pgxpool.Pool, fixture InstallationFixture) {
	t.Helper()
	SeedProvider(t, pool, fixture.Provider)
	_, err := pool.Exec(context.Background(), `
		INSERT INTO integration_installations (
			installation_id, tenant_id, provider_code, family, display_name,
			status, health_status, created_at, updated_at
		) VALUES ($1, $2, $3, 'marketplace', $4, 'draft', 'healthy', $5, $5)
	`, fixture.InstallationID, fixture.TenantID, fixture.Provider.Code, fixture.Provider.DisplayName,
		time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("seed installation fixture %q: %v", fixture.InstallationID, err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM integration_installations WHERE tenant_id=$1 AND installation_id=$2`, fixture.TenantID, fixture.InstallationID); err != nil {
			t.Errorf("cleanup installation fixture %q: %v", fixture.InstallationID, err)
		}
	})
}

func CleanupMarketplaceAccounts(t testing.TB, pool *pgxpool.Pool, tenantID string, accountIDs ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, accountID := range accountIDs {
			for _, statement := range []string{
				`DELETE FROM pricing_simulations WHERE tenant_id=$1 AND account_id=$2`,
				`DELETE FROM marketplace_pricing_policies WHERE tenant_id=$1 AND account_id=$2`,
				`DELETE FROM marketplace_accounts WHERE tenant_id=$1 AND account_id=$2`,
			} {
				if _, err := pool.Exec(context.Background(), statement, tenantID, accountID); err != nil {
					t.Errorf("cleanup marketplace account fixture %q: %v", accountID, err)
				}
			}
		}
	})
}
