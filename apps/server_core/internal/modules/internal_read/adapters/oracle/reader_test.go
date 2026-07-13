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

func TestTaxInputsRequireExactSaleLineIdentityBeforeOracleAccess(t *testing.T) {
	reader := NewReader(nil)
	got, err := reader.GetTaxInputs(context.Background(), ports.TaxInput{
		ProductID: 42664,
		Policy:    domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatalf("GetTaxInputs() error = %v, want missing tax without Oracle access", err)
	}
	if got.ICMSAmount != nil || len(got.QualityFlags) != 1 || got.QualityFlags[0] != domain.QualityMissingTax {
		t.Fatalf("GetTaxInputs() = %+v, want nil amounts with missing_tax", got)
	}
}

func TestBuildTaxInputsQueryUsesExactSaleLinePredicates(t *testing.T) {
	policy := domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC))
	policy.Source = domain.TaxSourceIdentity{DocumentID: 98765, LineNumber: 3}
	query, args := buildTaxInputsQuery(42664, policy)

	for _, predicate := range []string{"d.NUNOTA = :1", "d.SEQUENCIA = :2", "i.CODPROD = :3", "d.CODINC = :4"} {
		if !strings.Contains(query, predicate) {
			t.Fatalf("tax query missing exact predicate %q: %s", predicate, query)
		}
	}
	if strings.Contains(query, "DTNEG") || len(args) != 4 || args[0] != 98765 || args[1] != 3 || args[2] != 42664 {
		t.Fatalf("tax query args/query = %#v / %s, want exact document-line-product identity", args, query)
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
