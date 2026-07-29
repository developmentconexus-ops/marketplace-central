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
	err     error
}

func (f *fakeReader) FindProductsForLinking(ctx context.Context, _ internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	f.called = true
	f.gotCtx = ctx
	if f.err != nil {
		return nil, f.err
	}
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

// fakePagerReader is a fakeReader that also implements the optional
// CatalogPageReader capability, mirroring the real upload/oracle chains.
type fakePagerReader struct {
	fakeReader
	pageCalled   bool
	searchCalled bool
	byIDsCalled  bool
	pageCtx      context.Context
	policies     []internalreadports.SellableAssortmentPolicy
}

func (f *fakePagerReader) ListCatalogProductFacts(ctx context.Context, _ internalreadports.Cursor, _ int) (internalreadports.CatalogFactPage, error) {
	f.pageCalled = true
	f.pageCtx = ctx
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakePagerReader) CatalogProductFactsByIDs(ctx context.Context, _ []int64) (internalreadports.CatalogFactPage, error) {
	f.byIDsCalled = true
	f.pageCtx = ctx
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakePagerReader) SearchCatalogProductFacts(ctx context.Context, _ string, _ internalreadports.Cursor, _ int) (internalreadports.CatalogFactPage, error) {
	f.searchCalled = true
	f.pageCtx = ctx
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakePagerReader) ListCatalogProductFactsWithPolicy(ctx context.Context, cursor internalreadports.Cursor, limit int, policy internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	f.policies = append(f.policies, policy)
	return f.ListCatalogProductFacts(ctx, cursor, limit)
}

func (f *fakePagerReader) SearchCatalogProductFactsWithPolicy(ctx context.Context, query string, cursor internalreadports.Cursor, limit int, _ internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	return f.SearchCatalogProductFacts(ctx, query, cursor, limit)
}

func (f *fakePagerReader) GetCatalogAssortmentCounts(_ context.Context, policy internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogAssortmentCounts, error) {
	f.policies = append(f.policies, policy)
	return internalreadports.CatalogAssortmentCounts{}, nil
}

var _ internalreadports.CatalogPageReader = (*fakePagerReader)(nil)
var _ internalreadports.CatalogAssortmentReader = (*fakePagerReader)(nil)

// Regression guard for the /catalog 503: wrapping the upload reader in the
// routing Reader must NOT hide its CatalogPageReader capability from
// internal_read/application.Service's type assertion.
func TestReaderRoutesCatalogPageToUploadSource(t *testing.T) {
	upload := &fakePagerReader{}
	live := &fakePagerReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SetAt: time.Now()}}
	r := NewReader(upload, live, lookup, "t1")

	if _, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10); err != nil {
		t.Fatalf("ListCatalogProductFacts() error = %v", err)
	}
	if !upload.pageCalled {
		t.Fatal("expected upload pager to be invoked")
	}
	if live.pageCalled {
		t.Fatal("expected live pager NOT to be invoked")
	}
	source, ok := erpinternalread.ActiveSourceFromContext(upload.pageCtx)
	if !ok || source != erpdomain.SourceXLSX {
		t.Fatalf("erp ActiveSourceFromContext() = (%v, %v), want (xlsx, true)", source, ok)
	}
	if _, err := r.SearchCatalogProductFacts(context.Background(), "q", internalreadports.Cursor{}, 10); err != nil {
		t.Fatalf("SearchCatalogProductFacts() error = %v", err)
	}
	if !upload.searchCalled {
		t.Fatal("expected upload search pager to be invoked")
	}
}

func TestReaderUsesTheTenantPolicyForPageAndCounts(t *testing.T) {
	upload := &fakePagerReader{}
	readPolicy := tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SellableAssortment: readPolicy}}
	r := NewReader(upload, nil, lookup, "t1")

	if _, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10); err != nil {
		t.Fatal(err)
	}
	if _, err := r.GetCatalogAssortmentCounts(context.Background(), internalreadports.DefaultSellableAssortment()); err != nil {
		t.Fatal(err)
	}
	want := internalreadports.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true}
	if len(upload.policies) != 2 || upload.policies[0] != want || upload.policies[1] != want {
		t.Fatalf("page/count policies read = %+v, want [%+v %+v]", upload.policies, want, want)
	}
}

func TestReaderCatalogPageSankhyaNoLiveFailsHonest(t *testing.T) {
	upload := &fakePagerReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}
	r := NewReader(upload, nil, lookup, "t1")

	_, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10)
	if !errors.Is(err, ErrActiveSourceUnavailable) {
		t.Fatalf("ListCatalogProductFacts() error = %v, want ErrActiveSourceUnavailable", err)
	}
	if upload.pageCalled {
		t.Fatal("expected upload pager NOT to be invoked (no fallback)")
	}
}

func TestReaderCatalogPageNonPagerFailsHonest(t *testing.T) {
	upload := &fakeReader{} // no CatalogPageReader capability
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SetAt: time.Now()}}
	r := NewReader(upload, nil, lookup, "t1")

	_, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10)
	if !internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable) {
		t.Fatalf("ListCatalogProductFacts() error = %v, want ReadErrorSourceUnavailable", err)
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
