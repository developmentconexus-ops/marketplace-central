package pgdb

import "testing"

func TestLoadConfigRequiresEncryptionKey(t *testing.T) {
	t.Setenv("MC_DATABASE_URL", "postgres://localhost/marketplace")
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_custom")
	t.Setenv("MPC_ENCRYPTION_KEY", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when MPC_ENCRYPTION_KEY is empty")
	}
}

// A command that never decrypts a credential must not be able to fail for the
// lack of a key it would never use, and must not ask an operator to put a
// secret into the environment of a process with no use for it. Splitting the
// two loaders is what makes "this process holds no secret" a property of the
// code rather than a claim in a doc comment.
func TestLoadPoolConfigDoesNotRequireAnEncryptionKey(t *testing.T) {
	t.Setenv("MC_DATABASE_URL", "postgres://localhost/marketplace")
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_custom")
	t.Setenv("MPC_ENCRYPTION_KEY", "")

	cfg, err := LoadPoolConfig()
	if err != nil {
		t.Fatalf("LoadPoolConfig with no encryption key: %v", err)
	}
	if cfg.DatabaseURL != "postgres://localhost/marketplace" {
		t.Fatalf("DatabaseURL = %q, want the configured URL", cfg.DatabaseURL)
	}
	if cfg.DefaultTenantID != "tenant_custom" {
		t.Fatalf("DefaultTenantID = %q, want tenant_custom", cfg.DefaultTenantID)
	}
	if cfg.EncryptionKey != "" {
		t.Fatalf("EncryptionKey = %q, want empty: a pool config must not carry a secret it did not need", cfg.EncryptionKey)
	}
}

// The database URL is the one thing a pool cannot be opened without, so its
// absence stays a failure in both loaders.
func TestLoadPoolConfigStillRequiresTheDatabaseURL(t *testing.T) {
	t.Setenv("MC_DATABASE_URL", "")
	t.Setenv("MC_DEFAULT_TENANT_ID", "tenant_custom")
	t.Setenv("MPC_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")

	if _, err := LoadPoolConfig(); err == nil {
		t.Fatal("LoadPoolConfig error = nil, want an error naming MC_DATABASE_URL")
	}
}

func TestLoadPoolConfigDefaultsTenantID(t *testing.T) {
	t.Setenv("MC_DATABASE_URL", "postgres://localhost/marketplace")
	t.Setenv("MC_DEFAULT_TENANT_ID", "")
	t.Setenv("MPC_ENCRYPTION_KEY", "")

	cfg, err := LoadPoolConfig()
	if err != nil {
		t.Fatalf("LoadPoolConfig: %v", err)
	}
	if cfg.DefaultTenantID != "tenant_default" {
		t.Fatalf("DefaultTenantID = %q, want tenant_default", cfg.DefaultTenantID)
	}
}

func TestLoadConfigDefaultsTenantID(t *testing.T) {
	t.Setenv("MC_DATABASE_URL", "postgres://localhost/marketplace")
	t.Setenv("MC_DEFAULT_TENANT_ID", "")
	t.Setenv("MPC_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultTenantID != "tenant_default" {
		t.Fatalf("expected tenant_default, got %q", cfg.DefaultTenantID)
	}
}
