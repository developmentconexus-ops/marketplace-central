package postgres_test

import (
	"strings"
	"testing"

	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

const validTarget = "postgresql://fixture@127.0.0.1:49152/mpc_test_0123456789abcdef0123456789abcdef?sslmode=disable"

func TestLoadConfigReadsOnlyGeneratedTestTarget(t *testing.T) {
	requested := []string{}
	getenv := func(key string) string {
		requested = append(requested, key)
		if key == "MPC_TEST_DATABASE_URL" {
			return validTarget
		}
		if key == "MC_DATABASE_URL" {
			return "postgresql://fixture@127.0.0.1:5435/marketplace_central"
		}
		return ""
	}

	cfg, err := testpostgres.LoadConfig(getenv, "tenant_fixture", "fixture-key")
	if err != nil {
		t.Fatalf("load safe target: %v", err)
	}
	if cfg.DatabaseURL != validTarget {
		t.Fatalf("unexpected target URL")
	}
	if strings.Join(requested, ",") != "MPC_TEST_DATABASE_URL" {
		t.Fatalf("ambient keys read: %v", requested)
	}
}

func TestLoadConfigRejectsUnsafeTargetsBeforeBuildingConfig(t *testing.T) {
	tests := map[string]string{
		"missing":      "",
		"malformed":    "://bad",
		"remote":       "postgresql://fixture@db.example.test:5432/mpc_test_0123456789abcdef0123456789abcdef",
		"no-port":      "postgresql://fixture@127.0.0.1/mpc_test_0123456789abcdef0123456789abcdef",
		"dev-database": "postgresql://fixture@127.0.0.1:5435/marketplace_central",
		"wrong-prefix": "postgresql://fixture@127.0.0.1:49152/test_0123456789abcdef0123456789abcdef",
		"wrong-run-id": "postgresql://fixture@127.0.0.1:49152/mpc_test_short",
		"uppercase-id": "postgresql://fixture@127.0.0.1:49152/mpc_test_0123456789ABCDEF0123456789ABCDEF",
		"wrong-scheme": "http://127.0.0.1:49152/mpc_test_0123456789abcdef0123456789abcdef",
	}
	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := testpostgres.LoadConfig(func(key string) string {
				if key == "MPC_TEST_DATABASE_URL" {
					return value
				}
				return ""
			}, "tenant_fixture", "fixture-key")
			if err == nil {
				t.Fatalf("unsafe target accepted: %+v", cfg)
			}
			if !strings.Contains(err.Error(), "HPG_TARGET_INVALID") {
				t.Fatalf("unstable target rejection: %v", err)
			}
		})
	}
}
