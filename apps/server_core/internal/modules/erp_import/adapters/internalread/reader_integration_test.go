//go:build integration

package internalread

import (
	"context"
	"errors"
	"testing"
	"time"

	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestReaderRealRepositoryCostReservedAsOfAndRejectedIgnored(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := WithActiveSource(context.Background(), erpdomain.SourceXLSX)
	tenant := "erp-internalread-" + time.Now().UTC().Format("150405.000000000")
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, tenant)
	})
	repo := erppostgres.NewRepository(pool)
	imported := time.Now().UTC().Add(-time.Hour).Truncate(time.Microsecond)
	reserved := "2"
	completed := erpdomain.ImportSnapshot{ID: "a1111111-1111-1111-1111-111111111111", Protocol: "#711-E", FileSHA256: "internalread-completed", Source: erpdomain.SourceXLSX, ImportedAt: imported, Status: erpdomain.ImportStatusCompleted, AcceptedRows: []erpdomain.NormalizedRow{{Codprod: "10", Descrprod: "Widget", Custo: "12.30", StockPhysical: "9", StockReserved: &reserved}}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, completed); err != nil {
		t.Fatal(err)
	}
	rejected := erpdomain.ImportSnapshot{ID: "a2222222-2222-2222-2222-222222222222", Protocol: "#722-E", FileSHA256: "internalread-rejected", Source: erpdomain.SourceXLSX, ImportedAt: imported.Add(time.Minute), Status: erpdomain.ImportStatusRejected}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, rejected); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SyncLatestCompletedSnapshot(ctx, tenant, erpdomain.SourceXLSX); err != nil {
		t.Fatal(err)
	}
	r := NewReader(repo, tenant)
	stock, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 10})
	if err != nil || stock.Quantity == nil || *stock.Quantity != 7 {
		t.Fatalf("stock=%+v err=%v", stock, err)
	}
	cost, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 10, Policy: readdomain.CostAsOfPolicy{EffectiveAt: imported}})
	if err != nil || cost.Amount == nil || *cost.Amount != 12.3 || cost.Source.ObservedAt == nil || !cost.Source.ObservedAt.Equal(imported) {
		t.Fatalf("cost=%+v err=%v", cost, err)
	}
	stale, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 10, Policy: readdomain.CostAsOfPolicy{EffectiveAt: imported.Add(-time.Second)}})
	if err != nil || stale.Amount == nil || !readdomain.HasQualityFlag(stale.QualityFlags, readdomain.QualityStaleSource) {
		t.Fatalf("as-of before the snapshot must answer flagged stale: cost=%+v err=%v", stale, err)
	}
}
