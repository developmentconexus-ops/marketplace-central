package oracle

import (
	"strings"
	"testing"
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
