package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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

// SkipWithoutTarget skips the test when the harness-owned PostgreSQL target
// is absent, so integration-tagged packages compile and skip clean without a DSN.
func SkipWithoutTarget(t testing.TB) {
	t.Helper()
	if os.Getenv("MPC_TEST_DATABASE_URL") == "" {
		t.Skip("MPC_TEST_DATABASE_URL is unset")
	}
}

// OpenPool fails closed when the harness-owned target is absent or unsafe.
// Integration tests must use this boundary instead of application ambient config.
func OpenPool(t testing.TB, tenantID string) (*pgxpool.Pool, pgdb.Config) {
	t.Helper()
	cfg, err := LoadConfig(os.Getenv, tenantID, "harness-test-key")
	if err != nil {
		t.Fatalf("load harness PostgreSQL target: %v", err)
	}
	pool, err := pgdb.NewPool(context.Background(), cfg)
	if err != nil {
		t.Fatalf("open harness PostgreSQL target: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool, cfg
}
