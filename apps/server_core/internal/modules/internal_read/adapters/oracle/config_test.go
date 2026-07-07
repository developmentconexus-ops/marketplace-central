package oracle

import (
	"strings"
	"testing"
)

func TestLoadConfigFromEnvDoesNotLeakSecretValues(t *testing.T) {
	_, err := LoadConfigFromEnv(func(key string) string {
		switch key {
		case "SANKHYA_DSN":
			return ""
		case "SANKHYA_PASSWORD":
			return "super-secret"
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
