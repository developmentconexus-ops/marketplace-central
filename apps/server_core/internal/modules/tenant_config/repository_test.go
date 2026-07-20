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
