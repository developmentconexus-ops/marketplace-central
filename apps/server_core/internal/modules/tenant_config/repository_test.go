package tenant_config_test

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/sourcekind"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Integration test against the harness-owned Postgres target. Skips clean
// when MPC_TEST_DATABASE_URL is unset (testpostgres.SkipWithoutTarget) —
// no fake pool is substituted.
func TestRepository_Get_NoRow_ReturnsErrUnknownActiveSource(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	pool, cfg := testpostgres.OpenPool(t, "tenant_repo_test_norow")
	repo := tenant_config.NewRepository(pool)

	_, err := repo.Get(context.Background(), cfg.DefaultTenantID)
	if !errors.Is(err, tenant_config.ErrUnknownActiveSource) {
		t.Fatalf("Get() error = %v, want ErrUnknownActiveSource", err)
	}
}

func TestRepository_Set_Get_RoundTrip(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	pool, cfg := testpostgres.OpenPool(t, "tenant_repo_test_roundtrip")
	repo := tenant_config.NewRepository(pool)

	in := tenant_config.Config{
		TenantID: cfg.DefaultTenantID,
		Source:   tenant_config.SourceXLSX,
		SetBy:    "test-operator",
	}
	if err := repo.Set(context.Background(), in); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	got, err := repo.Get(context.Background(), cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Source != tenant_config.SourceXLSX {
		t.Errorf("Source = %q, want %q", got.Source, tenant_config.SourceXLSX)
	}
	if got.Kind != sourcekind.UploadSnapshot {
		t.Errorf("Kind = %q, want %q (derived via DefaultKind)", got.Kind, sourcekind.UploadSnapshot)
	}
	if got.SetBy != "test-operator" {
		t.Errorf("SetBy = %q, want %q", got.SetBy, "test-operator")
	}
	if got.SetAt.IsZero() {
		t.Error("SetAt is zero, want now()")
	}
}

func TestRepository_SetSellableAssortment_RoundTripPersistsPerTenant(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	pool, cfg := testpostgres.OpenPool(t, "tenant_repo_test_sellable_assortment")
	repo := tenant_config.NewRepository(pool)
	ctx := context.Background()
	firstTenantID := cfg.DefaultTenantID + "-sellable-assortment"
	secondTenantID := firstTenantID + "-second"
	missingTenantID := firstTenantID + "-missing"
	if _, err := pool.Exec(ctx, `
		DELETE FROM active_source
		WHERE tenant_id IN ($1, $2, $3)
	`, firstTenantID, secondTenantID, missingTenantID); err != nil {
		t.Fatalf("clear sellable assortment test tenants: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `
			DELETE FROM active_source
			WHERE tenant_id IN ($1, $2, $3)
		`, firstTenantID, secondTenantID, missingTenantID); err != nil {
			t.Errorf("clean sellable assortment test tenants: %v", err)
		}
	})

	if err := repo.Set(ctx, tenant_config.Config{
		TenantID: firstTenantID,
		Source:   tenant_config.SourceXLSX,
	}); err != nil {
		t.Fatalf("Set() first tenant error = %v", err)
	}

	got, err := repo.Get(ctx, firstTenantID)
	if err != nil {
		t.Fatalf("Get() first tenant defaults error = %v", err)
	}
	if got.SellableAssortment.OnlyRevenda != true {
		t.Errorf("first tenant OnlyRevenda = %v, want true", got.SellableAssortment.OnlyRevenda)
	}
	if got.SellableAssortment.OnlyEmEstoque != true {
		t.Errorf("first tenant OnlyEmEstoque = %v, want true", got.SellableAssortment.OnlyEmEstoque)
	}
	if got.SellableAssortment.OnlyEcommerceEligible != false {
		t.Errorf("first tenant OnlyEcommerceEligible = %v, want false", got.SellableAssortment.OnlyEcommerceEligible)
	}

	firstRule := tenant_config.SellableAssortment{
		OnlyRevenda:           false,
		OnlyEmEstoque:         false,
		OnlyEcommerceEligible: true,
	}
	if err := repo.SetSellableAssortment(ctx, firstTenantID, firstRule); err != nil {
		t.Fatalf("SetSellableAssortment() first tenant error = %v", err)
	}

	got, err = repo.Get(ctx, firstTenantID)
	if err != nil {
		t.Fatalf("Get() first tenant rule error = %v", err)
	}
	if got.SellableAssortment.OnlyRevenda != firstRule.OnlyRevenda {
		t.Errorf("first tenant OnlyRevenda = %v, want %v", got.SellableAssortment.OnlyRevenda, firstRule.OnlyRevenda)
	}
	if got.SellableAssortment.OnlyEmEstoque != firstRule.OnlyEmEstoque {
		t.Errorf("first tenant OnlyEmEstoque = %v, want %v", got.SellableAssortment.OnlyEmEstoque, firstRule.OnlyEmEstoque)
	}
	if got.SellableAssortment.OnlyEcommerceEligible != firstRule.OnlyEcommerceEligible {
		t.Errorf("first tenant OnlyEcommerceEligible = %v, want %v", got.SellableAssortment.OnlyEcommerceEligible, firstRule.OnlyEcommerceEligible)
	}

	if err := repo.Set(ctx, tenant_config.Config{
		TenantID: secondTenantID,
		Source:   tenant_config.SourceSankhya,
	}); err != nil {
		t.Fatalf("Set() second tenant error = %v", err)
	}
	secondRule := tenant_config.SellableAssortment{
		OnlyRevenda:           true,
		OnlyEmEstoque:         true,
		OnlyEcommerceEligible: true,
	}
	if err := repo.SetSellableAssortment(ctx, secondTenantID, secondRule); err != nil {
		t.Fatalf("SetSellableAssortment() second tenant error = %v", err)
	}

	got, err = repo.Get(ctx, secondTenantID)
	if err != nil {
		t.Fatalf("Get() second tenant rule error = %v", err)
	}
	if got.SellableAssortment.OnlyRevenda != secondRule.OnlyRevenda {
		t.Errorf("second tenant OnlyRevenda = %v, want %v", got.SellableAssortment.OnlyRevenda, secondRule.OnlyRevenda)
	}
	if got.SellableAssortment.OnlyEmEstoque != secondRule.OnlyEmEstoque {
		t.Errorf("second tenant OnlyEmEstoque = %v, want %v", got.SellableAssortment.OnlyEmEstoque, secondRule.OnlyEmEstoque)
	}
	if got.SellableAssortment.OnlyEcommerceEligible != secondRule.OnlyEcommerceEligible {
		t.Errorf("second tenant OnlyEcommerceEligible = %v, want %v", got.SellableAssortment.OnlyEcommerceEligible, secondRule.OnlyEcommerceEligible)
	}

	got, err = repo.Get(ctx, firstTenantID)
	if err != nil {
		t.Fatalf("Get() first tenant after second write error = %v", err)
	}
	if got.SellableAssortment.OnlyRevenda != firstRule.OnlyRevenda {
		t.Errorf("first tenant after second write OnlyRevenda = %v, want %v", got.SellableAssortment.OnlyRevenda, firstRule.OnlyRevenda)
	}
	if got.SellableAssortment.OnlyEmEstoque != firstRule.OnlyEmEstoque {
		t.Errorf("first tenant after second write OnlyEmEstoque = %v, want %v", got.SellableAssortment.OnlyEmEstoque, firstRule.OnlyEmEstoque)
	}
	if got.SellableAssortment.OnlyEcommerceEligible != firstRule.OnlyEcommerceEligible {
		t.Errorf("first tenant after second write OnlyEcommerceEligible = %v, want %v", got.SellableAssortment.OnlyEcommerceEligible, firstRule.OnlyEcommerceEligible)
	}

	if err := repo.Set(ctx, tenant_config.Config{
		TenantID: firstTenantID,
		Source:   tenant_config.SourceCatalogoCliente,
	}); err != nil {
		t.Fatalf("Set() source switch error = %v", err)
	}
	got, err = repo.Get(ctx, firstTenantID)
	if err != nil {
		t.Fatalf("Get() after source switch error = %v", err)
	}
	if got.SellableAssortment.OnlyRevenda != firstRule.OnlyRevenda {
		t.Errorf("after source switch OnlyRevenda = %v, want %v", got.SellableAssortment.OnlyRevenda, firstRule.OnlyRevenda)
	}
	if got.SellableAssortment.OnlyEmEstoque != firstRule.OnlyEmEstoque {
		t.Errorf("after source switch OnlyEmEstoque = %v, want %v", got.SellableAssortment.OnlyEmEstoque, firstRule.OnlyEmEstoque)
	}
	if got.SellableAssortment.OnlyEcommerceEligible != firstRule.OnlyEcommerceEligible {
		t.Errorf("after source switch OnlyEcommerceEligible = %v, want %v", got.SellableAssortment.OnlyEcommerceEligible, firstRule.OnlyEcommerceEligible)
	}

	if err := repo.SetSellableAssortment(ctx, missingTenantID, firstRule); !errors.Is(err, tenant_config.ErrUnknownActiveSource) {
		t.Fatalf("SetSellableAssortment() missing tenant error = %v, want ErrUnknownActiveSource", err)
	}
	if _, err := repo.Get(ctx, missingTenantID); !errors.Is(err, tenant_config.ErrUnknownActiveSource) {
		t.Fatalf("Get() missing tenant after failed set error = %v, want ErrUnknownActiveSource", err)
	}
}
