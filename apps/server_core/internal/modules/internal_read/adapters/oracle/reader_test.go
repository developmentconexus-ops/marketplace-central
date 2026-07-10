package oracle

import (
	"context"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestBuildSellableStockQueryUsesOracleFirstPolicy(t *testing.T) {
	query, args, err := buildSellableStockQuery(42664, domain.DefaultSellableStockPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(query, "METALPRD.TGFEST") {
		t.Fatalf("expected oracle stock query to stay inside TGFEST, got %q", query)
	}
	if len(args) != 5 {
		t.Fatalf("expected bound args for product + companies + locations + exclusions, got %d", len(args))
	}
}

func TestReaderReturnsTypedUnavailableWhenDBMissing(t *testing.T) {
	reader := NewReader(nil)

	_, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{
		ProductID: 42664,
		Policy:    domain.DefaultSellableStockPolicy(),
	})
	if !domain.IsReadErrorCode(err, domain.ReadErrorSourceUnavailable) {
		t.Fatalf("expected source_unavailable read error, got %v", err)
	}
}

func TestBuildDSNIncludesOnlyConfigKeys(t *testing.T) {
	dsn := buildDSN(Config{
		Username:        "user",
		Password:        "secret",
		ConnectString:   "host:1521/service",
		PoolMinSessions: 1,
		PoolMaxSessions: 4,
		SessionTimeout:  5 * time.Minute,
	})
	if !strings.Contains(dsn, `connectString="host:1521/service"`) {
		t.Fatalf("expected connectString in dsn, got %q", dsn)
	}
	legacyDatabaseURLKey := "MS" + "_DATABASE_URL"
	if strings.Contains(dsn, legacyDatabaseURLKey) {
		t.Fatalf("did not expect legacy msdb keys in dsn %q", dsn)
	}
}
