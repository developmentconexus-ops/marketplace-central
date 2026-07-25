package routing

import (
	"context"
	"testing"
	"time"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

type fakeStockBatchReader struct {
	called bool
	gotCtx context.Context
}

func (f *fakeStockBatchReader) GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.StockFact, error) {
	f.called = true
	f.gotCtx = ctx
	facts := make(map[int64]*internalreaddomain.StockFact, len(ids))
	for _, id := range ids {
		facts[id] = &internalreaddomain.StockFact{ProductID: id}
	}
	return facts, nil
}

var _ internalreadports.StockBatchReader = (*fakeStockBatchReader)(nil)

func TestStockBatchReaderServesAnUploadSourceFromTheMirror(t *testing.T) {
	mirror, live := &fakeStockBatchReader{}, &fakeStockBatchReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SetAt: time.Now()}}

	facts, err := NewStockBatchReader(mirror, live, lookup, "t1").GetStockFactsByIDs(context.Background(), []int64{90033})
	if err != nil {
		t.Fatalf("GetStockFactsByIDs() error = %v", err)
	}
	if len(facts) != 1 {
		t.Fatalf("expected one fact, got %d", len(facts))
	}
	if !mirror.called || live.called {
		t.Fatalf("expected the mirror reader only: mirror=%v live=%v", mirror.called, live.called)
	}
	source, ok := erpinternalread.ActiveSourceFromContext(mirror.gotCtx)
	if !ok || source != erpdomain.SourceXLSX {
		t.Fatalf("erp ActiveSourceFromContext() = (%v, %v), want (xlsx, true)", source, ok)
	}
}

func TestStockBatchReaderPrefersTheLiveReaderForSankhya(t *testing.T) {
	mirror, live := &fakeStockBatchReader{}, &fakeStockBatchReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}

	if _, err := NewStockBatchReader(mirror, live, lookup, "t1").GetStockFactsByIDs(context.Background(), []int64{90033}); err != nil {
		t.Fatalf("GetStockFactsByIDs() error = %v", err)
	}
	if !live.called || mirror.called {
		t.Fatalf("expected the live reader only: mirror=%v live=%v", mirror.called, live.called)
	}
}

func TestStockBatchReaderFallsBackToTheSameSourcesMirrorRowsWhenNoLiveReaderIsWired(t *testing.T) {
	mirror := &fakeStockBatchReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}

	if _, err := NewStockBatchReader(mirror, nil, lookup, "t1").GetStockFactsByIDs(context.Background(), []int64{90033}); err != nil {
		t.Fatalf("GetStockFactsByIDs() error = %v", err)
	}
	if !mirror.called {
		t.Fatal("expected the mirror reader to answer with the same source's synced rows")
	}
	// Same source, older observation — never another source's catalog (ADR-17).
	source, ok := erpinternalread.ActiveSourceFromContext(mirror.gotCtx)
	if !ok || string(source) != string(tenant_config.SourceSankhya) {
		t.Fatalf("erp ActiveSourceFromContext() = (%v, %v), want (sankhya, true)", source, ok)
	}
}
