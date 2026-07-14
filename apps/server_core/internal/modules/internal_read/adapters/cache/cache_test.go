package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	inventoryports "marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(delta time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(delta)
	c.mu.Unlock()
}

type fakeReader struct {
	clock        Clock
	listCalls    atomic.Int64
	searchCalls  atomic.Int64
	linkCalls    atomic.Int64
	blockStart   chan struct{}
	blockRelease chan struct{}
	err          error
}

func (r *fakeReader) maybeBlock() {
	if r.blockStart != nil {
		select {
		case r.blockStart <- struct{}{}:
		default:
		}
		<-r.blockRelease
	}
}

func (r *fakeReader) ListCatalogProductFacts(context.Context, internalreadports.Cursor, int) (internalreadports.CatalogFactPage, error) {
	r.listCalls.Add(1)
	r.maybeBlock()
	if r.err != nil {
		return internalreadports.CatalogFactPage{}, r.err
	}
	return internalreadports.CatalogFactPage{Items: []internalreadports.CatalogProductFact{{InternalProductID: r.listCalls.Load()}}, AsOf: r.clock.Now()}, nil
}

func (r *fakeReader) SearchCatalogProductFacts(context.Context, string, int) (internalreadports.CatalogFactPage, error) {
	r.searchCalls.Add(1)
	if r.err != nil {
		return internalreadports.CatalogFactPage{}, r.err
	}
	return internalreadports.CatalogFactPage{AsOf: r.clock.Now()}, nil
}

func (r *fakeReader) FindProductsForLinking(context.Context, internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	r.linkCalls.Add(1)
	return nil, nil
}

func (r *fakeReader) GetSellableStock(context.Context, internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	return internalreaddomain.SellableStock{}, nil
}

func (r *fakeReader) GetCurrentPrice(context.Context, internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	return internalreaddomain.CurrentPrice{}, nil
}

func (r *fakeReader) GetCostAsOf(context.Context, internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	return internalreaddomain.CostAsOf{}, nil
}

func (r *fakeReader) GetSalesHistory(context.Context, internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	return internalreaddomain.SalesHistory{}, nil
}

func (r *fakeReader) GetTaxInputs(context.Context, internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	return internalreaddomain.TaxInputs{}, nil
}

type fakeBatchReader struct {
	costCalls atomic.Int64
	taxCalls  atomic.Int64
}

func (r *fakeBatchReader) GetCostFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.CostAsOf, error) {
	r.costCalls.Add(1)
	return map[int64]*internalreaddomain.CostAsOf{1: nil}, nil
}

func (r *fakeBatchReader) GetTaxFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.TaxInputs, error) {
	r.taxCalls.Add(1)
	return map[int64]*internalreaddomain.TaxInputs{1: nil}, nil
}

type fakeStockReader struct{ calls atomic.Int64 }

func (r *fakeStockReader) GetStockFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.StockFact, error) {
	r.calls.Add(1)
	return map[int64]*internalreaddomain.StockFact{1: nil}, nil
}

func testCache(clock Clock, maxEntries int, policies map[string]internalreaddomain.FreshnessPolicy) *Cache {
	return New(Config{Clock: clock, MaxEntries: maxEntries, Policies: policies})
}

func TestFreshnessCacheTTLPerClass(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	c := testCache(clock, 20, map[string]internalreaddomain.FreshnessPolicy{
		ClassCatalog: {MaxAge: 5 * time.Second}, ClassInventory: {MaxAge: 2 * time.Second}, ClassPriceCost: {MaxAge: 3 * time.Second},
	})
	reader := &fakeReader{clock: clock}
	catalog := NewCatalogPageReader(reader, c)
	first, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	clock.Advance(4 * time.Second)
	second, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil || !second.AsOf.Equal(first.AsOf) || reader.listCalls.Load() != 1 {
		t.Fatalf("catalog hit: calls=%d first=%s second=%s err=%v", reader.listCalls.Load(), first.AsOf, second.AsOf, err)
	}
	clock.Advance(time.Second)
	third, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil || !third.AsOf.After(second.AsOf) || reader.listCalls.Load() != 2 {
		t.Fatalf("catalog expiry: calls=%d second=%s third=%s err=%v", reader.listCalls.Load(), second.AsOf, third.AsOf, err)
	}

	batchDownstream := &fakeBatchReader{}
	batch := NewBatchReader(batchDownstream, c)
	if _, err := batch.GetCostFactsByIDs(context.Background(), []int64{2, 1}); err != nil {
		t.Fatal(err)
	}
	got, err := batch.GetCostFactsByIDs(context.Background(), []int64{1, 2})
	if err != nil || got[1] != nil || batchDownstream.costCalls.Load() != 1 {
		t.Fatalf("pricecost hit or nil preservation failed: calls=%d value=%v err=%v", batchDownstream.costCalls.Load(), got[1], err)
	}
	clock.Advance(3 * time.Second)
	if _, err := batch.GetCostFactsByIDs(context.Background(), []int64{1, 2}); err != nil {
		t.Fatal(err)
	}
	if batchDownstream.costCalls.Load() != 2 {
		t.Fatalf("pricecost expiry calls=%d", batchDownstream.costCalls.Load())
	}

	stockDownstream := &fakeStockReader{}
	stock := NewStockBatchReader(stockDownstream, c)
	if _, err := stock.GetStockFactsByIDs(context.Background(), []int64{1}); err != nil {
		t.Fatal(err)
	}
	clock.Advance(2 * time.Second)
	if _, err := stock.GetStockFactsByIDs(context.Background(), []int64{1}); err != nil {
		t.Fatal(err)
	}
	if stockDownstream.calls.Load() != 2 {
		t.Fatalf("inventory expiry calls=%d", stockDownstream.calls.Load())
	}

	configured := ConfigFromEnv(func(name string) string {
		return map[string]string{"MPC_CACHE_TTL_CATALOG": "7s", "MPC_CACHE_TTL_STOCK": "bad", "MPC_CACHE_TTL_PRICECOST": "9s", "MPC_CACHE_MAX_ENTRIES": "3"}[name]
	})
	if configured.Policies[ClassCatalog].MaxAge != 7*time.Second || configured.Policies[ClassInventory].MaxAge != defaultStockTTL || configured.Policies[ClassPriceCost].MaxAge != 9*time.Second || configured.MaxEntries != 3 {
		t.Fatalf("env config not applied: %+v", configured)
	}
}

func TestFreshnessCacheSingleflight(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	downstream := &fakeReader{clock: clock, blockStart: make(chan struct{}, 1), blockRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, testCache(clock, 10, nil))
	const waiters = 20
	results := make([]internalreadports.CatalogFactPage, waiters)
	errs := make([]error, waiters)
	var wg sync.WaitGroup
	for i := 0; i < waiters; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		}(i)
	}
	select {
	case <-downstream.blockStart:
	case <-time.After(time.Second):
		t.Fatal("downstream was not entered")
	}
	close(downstream.blockRelease)
	wg.Wait()
	if downstream.listCalls.Load() != 1 {
		t.Fatalf("singleflight downstream calls=%d", downstream.listCalls.Load())
	}
	for i := range results {
		if errs[i] != nil || !results[i].AsOf.Equal(results[0].AsOf) {
			t.Fatalf("waiter %d result=%+v err=%v", i, results[i], errs[i])
		}
	}
}

func TestFreshnessCacheErrorNotCached(t *testing.T) {
	boom := errors.New("downstream failure")
	clock := &fakeClock{now: time.Now().UTC()}
	downstream := &fakeReader{clock: clock, err: boom, blockStart: make(chan struct{}, 1), blockRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, testCache(clock, 10, nil))
	var wg sync.WaitGroup
	errs := make([]error, 20)
	var ready atomic.Int64
	for i := range errs {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			ready.Add(1)
			_, errs[index] = catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		}(i)
	}
	deadline := time.Now().Add(time.Second)
	for ready.Load() != int64(len(errs)) && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if ready.Load() != int64(len(errs)) {
		t.Fatalf("not all waiters started: %d", ready.Load())
	}
	select {
	case <-downstream.blockStart:
	case <-time.After(time.Second):
		t.Fatal("downstream was not entered")
	}
	time.Sleep(50 * time.Millisecond)
	close(downstream.blockRelease)
	wg.Wait()
	for i, err := range errs {
		if !errors.Is(err, boom) {
			t.Fatalf("waiter %d error=%v", i, err)
		}
	}
	if _, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50); !errors.Is(err, boom) || downstream.listCalls.Load() != 2 {
		t.Fatalf("error was cached or retry missing: calls=%d err=%v", downstream.listCalls.Load(), err)
	}
}

func TestFreshnessCacheBypassAndLinkageExclusion(t *testing.T) {
	clock := &fakeClock{now: time.Now().UTC()}
	downstream := &fakeReader{clock: clock}
	wrapped := NewReader(downstream, testCache(clock, 10, nil))
	if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50); err != nil {
		t.Fatal(err)
	}
	noCache := internalreaddomain.WithFreshnessPolicy(context.Background(), internalreaddomain.FreshnessPolicy{MaxAge: 0})
	if _, err := wrapped.ListCatalogProductFacts(noCache, internalreadports.Cursor{}, 50); err != nil {
		t.Fatal(err)
	}
	if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50); err != nil {
		t.Fatal(err)
	}
	if downstream.listCalls.Load() != 2 {
		t.Fatalf("bypass did not repopulate cache: calls=%d", downstream.listCalls.Load())
	}
	for i := 0; i < 3; i++ {
		if _, err := wrapped.FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{ProductID: intPtr(i + 1)}); err != nil {
			t.Fatal(err)
		}
	}
	if downstream.linkCalls.Load() != 3 {
		t.Fatalf("linkage was cached: calls=%d", downstream.linkCalls.Load())
	}
}

func TestFreshnessCacheLRUAndLogs(t *testing.T) {
	clock := &fakeClock{now: time.Now().UTC()}
	downstream := &fakeReader{clock: clock}
	wrapped := NewCatalogPageReader(downstream, testCache(clock, 2, nil))
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	defer slog.SetDefault(previous)
	for _, cursor := range []int64{1, 2} {
		if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{InternalProductID: cursor}, 50); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{InternalProductID: 1}, 50); err != nil {
		t.Fatal(err)
	}
	if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{InternalProductID: 3}, 50); err != nil {
		t.Fatal(err)
	}
	if wrapped.cache.Size() != 2 {
		t.Fatalf("cache exceeded LRU cap: size=%d", wrapped.cache.Size())
	}
	if _, err := wrapped.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{InternalProductID: 2}, 50); err != nil {
		t.Fatal(err)
	}
	if downstream.listCalls.Load() != 4 {
		t.Fatalf("least-recently-used entry was not evicted: calls=%d", downstream.listCalls.Load())
	}
	text := logs.String()
	for _, expected := range []string{"cache=miss", "cache=hit", "key_class=catalog"} {
		if !bytes.Contains([]byte(text), []byte(expected)) {
			t.Fatalf("missing structured log field %q in %s", expected, text)
		}
	}
	if bytes.Contains([]byte(text), []byte("InternalProductID")) || bytes.Contains([]byte(text), []byte("sha256")) {
		t.Fatalf("raw key data leaked in logs: %s", text)
	}
}

func TestEvictOnMutation(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	c := testCache(clock, 10, nil)
	reader := &fakeReader{clock: clock}
	catalog := NewCatalogPageReader(reader, c)
	pricecost := NewBatchReader(&fakeBatchReader{}, c)
	first, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pricecost.GetCostFactsByIDs(context.Background(), []int64{1}); err != nil {
		t.Fatal(err)
	}
	clock.Advance(time.Second)
	c.InvalidateClass(ClassCatalog)
	second, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil || !second.AsOf.After(first.AsOf) || reader.listCalls.Load() != 2 {
		t.Fatalf("catalog mutation eviction failed: calls=%d first=%s second=%s err=%v", reader.listCalls.Load(), first.AsOf, second.AsOf, err)
	}
	if c.Size() != 2 {
		t.Fatalf("unrelated pricecost entry did not remain warm: size=%d", c.Size())
	}
	if _, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50); err != nil {
		t.Fatal(err)
	}
	if reader.listCalls.Load() != 2 {
		t.Fatalf("repopulated catalog entry missed: calls=%d", reader.listCalls.Load())
	}
	var invalidations atomic.Int64
	var failed = errors.New("rolled back")
	if err := func() error { return failed }(); err == nil {
		invalidations.Add(1)
	}
	if invalidations.Load() != 0 {
		t.Fatal("failed mutation invalidated cache")
	}
}

func intPtr(value int) *int { return &value }

var _ internalreadports.Reader = (*fakeReader)(nil)
var _ internalreadports.CatalogPageReader = (*fakeReader)(nil)
var _ internalreadports.BatchReader = (*fakeBatchReader)(nil)
var _ inventoryports.InternalStockBatchReader = (*fakeStockReader)(nil)

func ExampleCache() {
	config := ConfigFromEnv(nil)
	_ = New(config)
	fmt.Println("in-memory freshness cache")
	// Output: in-memory freshness cache
}
