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

// fakeReader is an internalreadports.Source that records the ctx it was invoked
// with, for assertions on what the routing Reader threaded through, and the
// policy each paged read was asked for.
//
// There used to be a second fake here — a "pager" variant — because the paged
// reads lived in a nominally optional port and a source could honestly implement
// one half. It cannot any more: a source either is one or is not, and the
// compiler decides which. One fake is what the port now describes.
type fakeReader struct {
	called  bool
	gotCtx  context.Context
	product internalreaddomain.ProductCandidate
	err     error

	pageCalled   bool
	searchCalled bool
	byIDsCalled  bool
	pageCtx      context.Context
	policies     []*internalreadports.SellableAssortmentPolicy
}

func (f *fakeReader) ListCatalogProductFacts(ctx context.Context, _ internalreadports.Cursor, _ int, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	f.pageCalled = true
	f.pageCtx = ctx
	f.policies = append(f.policies, policy)
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakeReader) SearchCatalogProductFacts(ctx context.Context, _ string, _ internalreadports.Cursor, _ int, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	f.searchCalled = true
	f.pageCtx = ctx
	f.policies = append(f.policies, policy)
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakeReader) CatalogProductFactsByIDs(ctx context.Context, _ []int64) (internalreadports.CatalogFactPage, error) {
	f.byIDsCalled = true
	f.pageCtx = ctx
	return internalreadports.CatalogFactPage{}, nil
}

func (f *fakeReader) GetCatalogAssortmentCounts(_ context.Context, policy *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogAssortmentCounts, error) {
	f.policies = append(f.policies, policy)
	return internalreadports.CatalogAssortmentCounts{}, nil
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

var _ internalreadports.Source = (*fakeReader)(nil)

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

// Regression guard for the /catalog 503: the routing Reader must forward the
// paged reads to the ACTIVE source and no other, with the same context the
// unpaged reads get.
func TestReaderRoutesCatalogPageToUploadSource(t *testing.T) {
	upload := &fakeReader{}
	live := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SetAt: time.Now()}}
	r := NewReader(upload, live, lookup, "t1")

	if _, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10, nil); err != nil {
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
	if _, err := r.SearchCatalogProductFacts(context.Background(), "q", internalreadports.Cursor{}, 10, nil); err != nil {
		t.Fatalf("SearchCatalogProductFacts() error = %v", err)
	}
	if !upload.searchCalled {
		t.Fatal("expected upload search pager to be invoked")
	}
}

func TestReaderUsesTheTenantPolicyForPageAndCounts(t *testing.T) {
	upload := &fakeReader{}
	readPolicy := tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SellableAssortment: readPolicy}}
	r := NewReader(upload, nil, lookup, "t1")

	if _, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := r.GetCatalogAssortmentCounts(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	want := internalreadports.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true}
	if len(upload.policies) != 2 {
		t.Fatalf("downstream calls = %d, want 2 (page and count)", len(upload.policies))
	}
	for i, got := range upload.policies {
		// A nil below this seam is the invariant breaking, not a default being
		// taken: the adapter would have to invent a rule to keep going.
		if got == nil {
			t.Fatalf("call %d: policy arrived nil below the routing seam, want %+v", i, want)
		}
		if *got != want {
			t.Fatalf("call %d: policy read = %+v, want %+v", i, *got, want)
		}
	}
}

// A caller that names a policy is asking a different question from a caller that
// names none, and routing must not confuse the two. The shape this guards against
// shipped once: the parameter was compared for equality against a constant and
// then thrown away, so the only value a caller could actually communicate was the
// zero one, and every explicit policy was silently replaced by the stored rule.
func TestReaderHonorsAnExplicitPolicyInsteadOfTheStoredOne(t *testing.T) {
	upload := &fakeReader{}
	stored := tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceXLSX, SellableAssortment: stored}}
	r := NewReader(upload, nil, lookup, "t1")

	asked := internalreadports.AllProductsAssortment()
	if _, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10, &asked); err != nil {
		t.Fatal(err)
	}
	if _, err := r.GetCatalogAssortmentCounts(context.Background(), &asked); err != nil {
		t.Fatal(err)
	}
	if len(upload.policies) != 2 {
		t.Fatalf("downstream calls = %d, want 2 (page and count)", len(upload.policies))
	}
	for i, got := range upload.policies {
		if got == nil {
			t.Fatalf("call %d: policy arrived nil below the routing seam, want %+v", i, asked)
		}
		if *got != asked {
			t.Fatalf("call %d: policy read = %+v, want the requested %+v (tenant stores %+v)", i, *got, asked, stored)
		}
	}
}

func TestReaderCatalogPageSankhyaNoLiveFailsHonest(t *testing.T) {
	upload := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}
	r := NewReader(upload, nil, lookup, "t1")

	_, err := r.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 10, nil)
	if !errors.Is(err, ErrActiveSourceUnavailable) {
		t.Fatalf("ListCatalogProductFacts() error = %v, want ErrActiveSourceUnavailable", err)
	}
	if upload.pageCalled {
		t.Fatal("expected upload pager NOT to be invoked (no fallback)")
	}
}

// TestReaderCatalogPageNonPagerFailsHonest stood here. It wired a source that
// implemented only the non-paged half and asserted the routing reader answered
// source_unavailable rather than panicking. Folding the ports deleted the state
// it described: NewReader takes an internalreadports.Source, so a half source is
// not a value this package can be handed. The guarantee survives as a build
// error at the wiring site, which is where it was always supposed to land.

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
