package routing

import (
	"context"
	"errors"
	"testing"
	"time"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

// fakeReader is a bare internalreadports.Reader that records the ctx it was
// invoked with, for assertions on what the routing Reader threaded through.
type fakeReader struct {
	called  bool
	gotCtx  context.Context
	product internalreaddomain.ProductCandidate
}

func (f *fakeReader) FindProductsForLinking(ctx context.Context, _ internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	f.called = true
	f.gotCtx = ctx
	return []internalreaddomain.ProductCandidate{f.product}, nil
}

func (f *fakeReader) GetSellableStock(context.Context, internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	return internalreaddomain.SellableStock{}, nil
}

func (f *fakeReader) GetCurrentPrice(context.Context, internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	return internalreaddomain.CurrentPrice{}, nil
}

func (f *fakeReader) GetCostAsOf(context.Context, internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	return internalreaddomain.CostAsOf{}, nil
}

func (f *fakeReader) GetSalesHistory(context.Context, internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	return internalreaddomain.SalesHistory{}, nil
}

func (f *fakeReader) GetTaxInputs(context.Context, internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	return internalreaddomain.TaxInputs{}, nil
}

var _ internalreadports.Reader = (*fakeReader)(nil)

// fakeLookup is a scripted tenant_config.ActiveSourceLookup.
type fakeLookup struct {
	cfg tenant_config.Config
	err error
}

func (f fakeLookup) Get(context.Context, string) (tenant_config.Config, error) {
	return f.cfg, f.err
}

var _ tenant_config.ActiveSourceLookup = fakeLookup{}

func TestReaderResolvesXLSX(t *testing.T) {
	upload := &fakeReader{}
	live := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SetAt: time.Now()}}
	r := NewReader(upload, live, lookup, "t1")

	if _, err := r.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{}); err != nil {
		t.Fatalf("FindProductsForLinking() error = %v", err)
	}
	if !upload.called {
		t.Fatal("expected upload reader to be invoked")
	}
	if live.called {
		t.Fatal("expected live reader NOT to be invoked")
	}
	source, ok := erpinternalread.ActiveSourceFromContext(upload.gotCtx)
	if !ok || source != erpdomain.SourceXLSX {
		t.Fatalf("erp ActiveSourceFromContext() = (%v, %v), want (xlsx, true)", source, ok)
	}
	cfg, ok := tenant_config.FromContext(upload.gotCtx)
	if !ok || cfg.Source != tenant_config.SourceXLSX {
		t.Fatalf("tenant_config.FromContext() = (%v, %v), want (xlsx, true)", cfg, ok)
	}
}

func TestReaderResolvesCatalogoCliente(t *testing.T) {
	upload := &fakeReader{}
	live := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceCatalogoCliente, SetAt: time.Now()}}
	r := NewReader(upload, live, lookup, "t1")

	if _, err := r.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{}); err != nil {
		t.Fatalf("FindProductsForLinking() error = %v", err)
	}
	if !upload.called {
		t.Fatal("expected upload reader to be invoked")
	}
	if live.called {
		t.Fatal("expected live reader NOT to be invoked")
	}
	source, ok := erpinternalread.ActiveSourceFromContext(upload.gotCtx)
	if !ok || source != erpdomain.SourceCatalogoCliente {
		t.Fatalf("erp ActiveSourceFromContext() = (%v, %v), want (catalogo_cliente, true)", source, ok)
	}
}

func TestReaderResolvesSankhyaWithLivePresent(t *testing.T) {
	upload := &fakeReader{}
	live := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}
	r := NewReader(upload, live, lookup, "t1")

	if _, err := r.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{}); err != nil {
		t.Fatalf("FindProductsForLinking() error = %v", err)
	}
	if !live.called {
		t.Fatal("expected live reader to be invoked")
	}
	if upload.called {
		t.Fatal("expected upload reader NOT to be invoked")
	}
	cfg, ok := tenant_config.FromContext(live.gotCtx)
	if !ok || cfg.Source != tenant_config.SourceSankhya {
		t.Fatalf("tenant_config.FromContext() = (%v, %v), want (sankhya, true)", cfg, ok)
	}
	if _, ok := erpinternalread.ActiveSourceFromContext(live.gotCtx); ok {
		t.Fatal("expected erp active-source sub-toggle NOT to be set for sankhya")
	}
}

func TestReaderSankhyaWithNoLiveReaderFailsHonest(t *testing.T) {
	upload := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}
	r := NewReader(upload, nil, lookup, "t1")

	_, err := r.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{})
	if !errors.Is(err, ErrActiveSourceUnavailable) {
		t.Fatalf("FindProductsForLinking() error = %v, want ErrActiveSourceUnavailable", err)
	}
	if upload.called {
		t.Fatal("expected upload reader NOT to be invoked (no fallback)")
	}
}

func TestReaderUnknownActiveSourceFailsClosed(t *testing.T) {
	upload := &fakeReader{}
	live := &fakeReader{}
	lookup := fakeLookup{err: tenant_config.ErrUnknownActiveSource}
	r := NewReader(upload, live, lookup, "t1")

	_, err := r.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{})
	if !errors.Is(err, tenant_config.ErrUnknownActiveSource) {
		t.Fatalf("FindProductsForLinking() error = %v, want ErrUnknownActiveSource", err)
	}
	if upload.called || live.called {
		t.Fatal("expected neither reader to be invoked")
	}
}
