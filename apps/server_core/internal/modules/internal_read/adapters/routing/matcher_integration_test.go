//go:build integration

package routing

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	internalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	tenant := "matcher-rule-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM active_source WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `DELETE FROM products_mirror_stock_locations WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, tenant)
	})

	importedAt := time.Now().UTC().Add(-time.Minute).Truncate(time.Microsecond)
	protocolID := fmt.Sprintf("70000000-0000-0000-0000-%012d", time.Now().UTC().UnixNano()%1000000000000)
	snapshot := erpdomain.ImportSnapshot{
		ID:         erpdomain.ImportID(protocolID),
		Protocol:   erpdomain.Protocol("#751-E"),
		FileSHA256: erpdomain.FileSHA256("s5-matcher-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)),
		Source:     erpdomain.SourceXLSX,
		ImportedAt: importedAt,
		Status:     erpdomain.ImportStatusCompleted,
		AcceptedRows: []erpdomain.NormalizedRow{
			{Codprod: "7001", Descrprod: "V only", Custo: "1", StockPhysical: "5", Usoprod: matcherStringPtr("V"), ADEcommerce: matcherStringPtr("S")},
			{Codprod: "7002", Descrprod: "R sellable", Custo: "1", StockPhysical: "5", Usoprod: matcherStringPtr("R"), ADEcommerce: matcherStringPtr("S")},
			{Codprod: "7003", Descrprod: "Unknown USOPROD", Custo: "1", StockPhysical: "5", ADEcommerce: matcherStringPtr("S")},
		},
	}
	importRepo := erppostgres.NewRepository(pool)
	if err := importRepo.PersistSnapshotAtomically(ctx, tenant, snapshot); err != nil {
		t.Fatal(err)
	}
	if _, err := importRepo.SyncLatestCompletedSnapshot(ctx, tenant, erpdomain.SourceXLSX); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		UPDATE products_mirror
		SET usoprod = CASE codigo_produto
			WHEN '7001' THEN 'V'
			WHEN '7002' THEN 'R'
			ELSE NULL
		END,
		ad_ecommerce = 'S'
		WHERE tenant_id=$1 AND source=$2
	`, tenant, erpdomain.SourceXLSX); err != nil {
		t.Fatal(err)
	}
	readBack, found, err := importRepo.MirrorProductByCode(ctx, tenant, erpdomain.SourceXLSX, "7001")
	if err != nil {
		t.Fatal(err)
	}
	if !found || readBack.Usoprod == nil || *readBack.Usoprod != "V" || readBack.ADEcommerce == nil || *readBack.ADEcommerce != "S" {
		t.Fatalf("mirror columns read back = usoprod=%v ad_ecommerce=%v found=%v, want V/S/true", readBack.Usoprod, readBack.ADEcommerce, found)
	}

	configRepo := tenant_config.NewRepository(pool)
	if err := configRepo.Set(ctx, tenant_config.Config{TenantID: tenant, Source: tenant_config.SourceXLSX}); err != nil {
		t.Fatal(err)
	}
	if err := configRepo.SetSellableAssortment(ctx, tenant, tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}); err != nil {
		t.Fatal(err)
	}
	matcher := NewMirrorMatcher(internalread.NewReader(importRepo, tenant), configRepo, tenant)
	input := internalreadports.FindProductsInput{}

	before, err := matcher.FindProductsForLinking(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	beforeIDs := candidateIDs(before)
	t.Logf("candidate IDs before toggle = %v", beforeIDs)
	wantBefore := []int{7002, 7003}
	if !reflect.DeepEqual(beforeIDs, wantBefore) {
		t.Fatalf("candidate IDs before toggle = %v, want %v", beforeIDs, wantBefore)
	}

	if err := configRepo.SetSellableAssortment(ctx, tenant, tenant_config.SellableAssortment{}); err != nil {
		t.Fatal(err)
	}
	after, err := matcher.FindProductsForLinking(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	afterIDs := candidateIDs(after)
	t.Logf("candidate IDs after toggle = %v", afterIDs)
	wantAfter := []int{7001, 7002, 7003}
	if !reflect.DeepEqual(afterIDs, wantAfter) {
		t.Fatalf("candidate IDs after toggle = %v, want %v", afterIDs, wantAfter)
	}
}

func candidateIDs(candidates []internalreaddomain.ProductCandidate) []int {
	ids := make([]int, 0, len(candidates))
	for _, candidate := range candidates {
		ids = append(ids, candidate.ProductID)
	}
	return ids
}

func matcherStringPtr(value string) *string { return &value }
