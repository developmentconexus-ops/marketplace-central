package postgres

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

const targetKey = "MPC_TEST_DATABASE_URL"

var generatedDatabase = regexp.MustCompile(`^mpc_test_[0-9a-f]{32}$`)

// LoadConfig reads only the harness-owned test target and constructs the
// ordinary pgdb configuration without consulting application ambient state.
func LoadConfig(getenv func(string) string, tenantID, encryptionKey string) (pgdb.Config, error) {
	target := getenv(targetKey)
	parsed, err := url.Parse(target)
	if err != nil || (parsed.Scheme != "postgres" && parsed.Scheme != "postgresql") {
		return pgdb.Config{}, fmt.Errorf("HPG_TARGET_INVALID")
	}
	host := parsed.Hostname()
	ip := net.ParseIP(host)
	if parsed.Port() == "" || ip == nil || !ip.IsLoopback() {
		return pgdb.Config{}, fmt.Errorf("HPG_TARGET_INVALID")
	}
	database := strings.TrimPrefix(parsed.Path, "/")
	if strings.Contains(database, "/") || !generatedDatabase.MatchString(database) {
		return pgdb.Config{}, fmt.Errorf("HPG_TARGET_INVALID")
	}
	return pgdb.Config{
		DatabaseURL:     target,
		DefaultTenantID: tenantID,
		EncryptionKey:   encryptionKey,
	}, nil
}
