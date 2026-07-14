package oracle

import (
	"strings"
	"testing"
	"time"
)

func TestLoadConfigFromEnvDoesNotLeakSecretValues(t *testing.T) {
	_, err := LoadConfigFromEnv(func(key string) string {
		switch key {
		case "MPC_ORACLE_USERNAME":
			return ""
		case "MPC_ORACLE_PASSWORD":
			return "super-secret"
		case "MPC_ORACLE_CONNECT_STRING":
			return ""
		default:
			return ""
		}
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "super-secret") {
		t.Fatal("error leaked secret value")
	}
}

func TestLoadConfigFromEnvUsesBoundedDefaults(t *testing.T) {
	cfg, err := LoadConfigFromEnv(func(key string) string {
		return map[string]string{
			"MPC_ORACLE_USERNAME":       "reader",
			"MPC_ORACLE_PASSWORD":       "secret",
			"MPC_ORACLE_CONNECT_STRING": "db:1521/service",
		}[key]
	})
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() error = %v", err)
	}
	if cfg.PoolMinSessions != 1 || cfg.PoolMaxSessions != 12 || cfg.PoolIncrement != 1 {
		t.Fatalf("unexpected pool defaults: %+v", cfg)
	}
	if cfg.PoolWaitTimeout != 5*time.Second || cfg.ConnectTimeout != 10*time.Second || cfg.CallTimeout != 30*time.Second {
		t.Fatalf("unexpected timeout defaults: %+v", cfg)
	}
}

func TestLoadConfigFromEnvAcceptsGovernedOverrides(t *testing.T) {
	values := map[string]string{
		"MPC_ORACLE_USERNAME":             "reader",
		"MPC_ORACLE_PASSWORD":             "secret",
		"MPC_ORACLE_CONNECT_STRING":       "db:1521/service",
		"MPC_ORACLE_POOL_MIN_SESSIONS":    "2",
		"MPC_ORACLE_POOL_MAX_SESSIONS":    "8",
		"MPC_ORACLE_POOL_INCREMENT":       "2",
		"MPC_ORACLE_POOL_WAIT_TIMEOUT":    "7s",
		"MPC_ORACLE_SESSION_TIMEOUT":      "9m",
		"MPC_ORACLE_SESSION_MAX_LIFETIME": "45m",
		"MPC_ORACLE_CONNECT_TIMEOUT":      "12s",
		"MPC_ORACLE_CALL_TIMEOUT":         "40s",
		"MPC_ORACLE_BOOTSTRAP_TIMEOUT":    "50s",
	}
	cfg, err := LoadConfigFromEnv(func(key string) string { return values[key] })
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() error = %v", err)
	}
	if cfg.PoolMinSessions != 2 || cfg.PoolMaxSessions != 8 || cfg.PoolIncrement != 2 || cfg.CallTimeout != 40*time.Second {
		t.Fatalf("overrides not applied: %+v", cfg)
	}
}

func TestLoadConfigFromEnvRejectsUnsafePoolAndTimeoutValues(t *testing.T) {
	base := map[string]string{
		"MPC_ORACLE_USERNAME":       "reader",
		"MPC_ORACLE_PASSWORD":       "secret",
		"MPC_ORACLE_CONNECT_STRING": "db:1521/service",
	}
	for _, tc := range []struct{ key, value string }{
		{"MPC_ORACLE_POOL_MAX_SESSIONS", "0"},
		{"MPC_ORACLE_POOL_INCREMENT", "13"},
		{"MPC_ORACLE_POOL_WAIT_TIMEOUT", "0s"},
		{"MPC_ORACLE_CONNECT_TIMEOUT", "1500ms"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			values := map[string]string{}
			for key, value := range base {
				values[key] = value
			}
			values[tc.key] = tc.value
			if _, err := LoadConfigFromEnv(func(key string) string { return values[key] }); err == nil {
				t.Fatalf("expected %s=%s to fail", tc.key, tc.value)
			}
		})
	}
}

func TestWithConnectTimeoutPreservesDescriptorsAndExistingSetting(t *testing.T) {
	if got := withConnectTimeout("db:1521/service", 12*time.Second); got != "db:1521/service?connect_timeout=12" {
		t.Fatalf("withConnectTimeout() = %q", got)
	}
	if got := withConnectTimeout("db:1521/service?expire_time=2", 12*time.Second); got != "db:1521/service?expire_time=2&connect_timeout=12" {
		t.Fatalf("withConnectTimeout() with options = %q", got)
	}
	for _, value := range []string{
		"db:1521/service?connect_timeout=7",
		"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=db)(PORT=1521)))",
	} {
		if got := withConnectTimeout(value, 12*time.Second); got != value {
			t.Fatalf("withConnectTimeout(%q) changed governed value to %q", value, got)
		}
	}
}
