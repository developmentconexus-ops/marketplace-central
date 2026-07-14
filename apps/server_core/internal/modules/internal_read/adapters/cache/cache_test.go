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

type stagedCatalogReader struct {
	clock        Clock
	calls        atomic.Int64
	firstStart   chan struct{}
	firstRelease chan struct{}
}

type staticCatalogReader struct {
	page internalreadports.CatalogFactPage
}

func (r *staticCatalogReader) ListCatalogProductFacts(context.Context, internalreadports.Cursor, int) (internalreadports.CatalogFactPage, error) {
	return r.page, nil
}

func (r *staticCatalogReader) SearchCatalogProductFacts(context.Context, string, int) (internalreadports.CatalogFactPage, error) {
	return r.page, nil
}

func (r *stagedCatalogReader) ListCatalogProductFacts(context.Context, internalreadports.Cursor, int) (internalreadports.CatalogFactPage, error) {
	call := r.calls.Add(1)
	if call == 1 {
		close(r.firstStart)
		<-r.firstRelease
	}
	return internalreadports.CatalogFactPage{AsOf: r.clock.Now().Add(time.Duration(call) * time.Nanosecond)}, nil
}

func (r *stagedCatalogReader) SearchCatalogProductFacts(context.Context, string, int) (internalreadports.CatalogFactPage, error) {
	return internalreadports.CatalogFactPage{AsOf: r.clock.Now()}, nil
}

type mutableBatchReader struct {
	cost map[int64]*internalreaddomain.CostAsOf
	tax  map[int64]*internalreaddomain.TaxInputs
}

func (r *mutableBatchReader) GetCostFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.CostAsOf, error) {
	return r.cost, nil
}

func (r *mutableBatchReader) GetTaxFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.TaxInputs, error) {
	return r.tax, nil
}

type mutableStockReader struct {
	facts map[int64]*internalreaddomain.StockFact
}

func (r *mutableStockReader) GetStockFactsByIDs(context.Context, []int64) (map[int64]*internalreaddomain.StockFact, error) {
	return r.facts, nil
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

func TestFreshnessCacheFencesInFlightLoadAfterInvalidation(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	c := testCache(clock, 10, nil)
	downstream := &fakeReader{clock: clock, blockStart: make(chan struct{}, 1), blockRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, c)
	firstResult := make(chan internalreadports.CatalogFactPage, 1)
	go func() {
		page, _ := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		firstResult <- page
	}()
	select {
	case <-downstream.blockStart:
	case <-time.After(time.Second):
		t.Fatal("downstream was not entered")
	}

	c.InvalidateClass(ClassCatalog)
	clock.Advance(time.Second)
	close(downstream.blockRelease)
	first := <-firstResult
	clock.Advance(time.Second)
	second, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	if !second.AsOf.After(first.AsOf) || downstream.listCalls.Load() != 2 {
		t.Fatalf("stale in-flight load was cached: calls=%d first=%s second=%s", downstream.listCalls.Load(), first.AsOf, second.AsOf)
	}
}

func TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	c := testCache(clock, 10, nil)
	downstream := &stagedCatalogReader{clock: clock, firstStart: make(chan struct{}), firstRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, c)
	firstResult := make(chan internalreadports.CatalogFactPage, 1)
	go func() {
		page, _ := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		firstResult <- page
	}()
	select {
	case <-downstream.firstStart:
	case <-time.After(time.Second):
		t.Fatal("first downstream load was not entered")
	}

	c.InvalidateClass(ClassCatalog)
	clock.Advance(time.Second)
	secondResult := make(chan internalreadports.CatalogFactPage, 1)
	go func() {
		page, _ := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		secondResult <- page
	}()
	close(downstream.firstRelease)
	first := <-firstResult
	second := <-secondResult
	if downstream.calls.Load() != 2 {
		t.Fatalf("post-invalidation reader joined stale flight: downstream calls=%d, want 2", downstream.calls.Load())
	}
	if !second.AsOf.After(first.AsOf) {
		t.Fatalf("post-invalidation reader got stale data: first=%s second=%s", first.AsOf, second.AsOf)
	}
}

func TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	downstream := &stagedCatalogReader{clock: clock, firstStart: make(chan struct{}), firstRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, testCache(clock, 10, nil))
	normalResult := make(chan internalreadports.CatalogFactPage, 1)
	go func() {
		page, _ := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		normalResult <- page
	}()
	select {
	case <-downstream.firstStart:
	case <-time.After(time.Second):
		t.Fatal("normal leader was not entered")
	}

	clock.Advance(time.Second)
	bypass := internalreaddomain.WithFreshnessPolicy(context.Background(), internalreaddomain.FreshnessPolicy{MaxAge: 0})
	fresh, err := catalog.ListCatalogProductFacts(bypass, internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	close(downstream.firstRelease)
	older := <-normalResult
	if downstream.calls.Load() != 2 {
		t.Fatalf("expected normal and bypass downstream calls, got %d", downstream.calls.Load())
	}
	if !fresh.AsOf.After(older.AsOf) {
		t.Fatalf("bypass snapshot=%s, normal snapshot=%s; want bypass newer", fresh.AsOf, older.AsOf)
	}

	got, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	if !got.AsOf.Equal(fresh.AsOf) {
		t.Fatalf("older normal load overwrote bypass refresh: got=%s, want=%s", got.AsOf, fresh.AsOf)
	}
}

func TestFreshnessCacheBypassUsesSeparateSingleflightNamespace(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	downstream := &stagedCatalogReader{clock: clock, firstStart: make(chan struct{}), firstRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, testCache(clock, 10, nil))
	normalResult := make(chan internalreadports.CatalogFactPage, 1)
	go func() {
		page, _ := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		normalResult <- page
	}()
	select {
	case <-downstream.firstStart:
	case <-time.After(time.Second):
		t.Fatal("normal leader was not entered")
	}

	bypass := internalreaddomain.WithFreshnessPolicy(context.Background(), internalreaddomain.FreshnessPolicy{MaxAge: 0})
	fresh, err := catalog.ListCatalogProductFacts(bypass, internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	close(downstream.firstRelease)
	<-normalResult
	if downstream.calls.Load() != 2 || !fresh.AsOf.After(clock.Now()) {
		t.Fatalf("bypass joined cached leader: calls=%d fresh=%s", downstream.calls.Load(), fresh.AsOf)
	}
}

func TestFreshnessCacheWaiterContextCancellation(t *testing.T) {
	clock := &fakeClock{now: time.Now().UTC()}
	downstream := &fakeReader{clock: clock, blockStart: make(chan struct{}, 1), blockRelease: make(chan struct{})}
	catalog := NewCatalogPageReader(downstream, testCache(clock, 10, nil))
	leaderDone := make(chan struct{})
	go func() {
		_, _ = catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
		close(leaderDone)
	}()
	select {
	case <-downstream.blockStart:
	case <-time.After(time.Second):
		t.Fatal("leader was not entered")
	}
	ctx, cancel := context.WithCancel(context.Background())
	canceled := make(chan error, 1)
	go func() {
		_, err := catalog.ListCatalogProductFacts(ctx, internalreadports.Cursor{}, 50)
		canceled <- err
	}()
	cancel()
	select {
	case err := <-canceled:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("waiter error=%v, want context canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("canceled waiter did not return")
	}
	close(downstream.blockRelease)
	<-leaderDone
	if downstream.listCalls.Load() != 1 {
		t.Fatalf("waiter cancellation changed singleflight: calls=%d", downstream.listCalls.Load())
	}
}

func TestFreshnessCacheDeepCopiesCachedFacts(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)}
	reference, description, ean := "REF-1", "Product", "EAN-1"
	quantity, amount := 4.0, "12.50"
	page := internalreadports.CatalogFactPage{
		Items: []internalreadports.CatalogProductFact{{
			InternalProductID: 1, Reference: &reference, Description: &description, EAN: &ean, Active: true,
			SellableStock: internalreadports.CatalogQuantityFact{Quantity: &quantity, Quality: make([]string, 1, 2)},
			CurrentPrice:  internalreadports.CatalogMoneyFact{Amount: &amount, Currency: "BRL", Quality: make([]string, 1, 2)},
			Cost:          internalreadports.CatalogMoneyFact{Amount: &amount, Currency: "BRL", Quality: make([]string, 1, 2)},
		}},
		NextCursor: &internalreadports.Cursor{InternalProductID: 7},
		AsOf:       clock.Now(),
	}
	page.Items[0].SellableStock.Quality[0] = "known"
	page.Items[0].CurrentPrice.Quality[0] = "known"
	page.Items[0].Cost.Quality[0] = "known"
	catalog := NewCatalogPageReader(&staticCatalogReader{page: page}, testCache(clock, 10, nil))
	got, err := catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	*got.Items[0].Reference = "mutated"
	*got.Items[0].Description = "mutated"
	*got.Items[0].EAN = "mutated"
	*got.Items[0].SellableStock.Quantity = 99
	got.Items[0].SellableStock.Quality[0] = "mutated"
	got.Items[0].SellableStock.Quality = append(got.Items[0].SellableStock.Quality, "mutated")
	*got.Items[0].CurrentPrice.Amount = "mutated"
	got.Items[0].CurrentPrice.Quality[0] = "mutated"
	got.Items[0].CurrentPrice.Quality = append(got.Items[0].CurrentPrice.Quality, "mutated")
	*got.Items[0].Cost.Amount = "mutated"
	got.Items[0].Cost.Quality[0] = "mutated"
	got.Items[0].Cost.Quality = append(got.Items[0].Cost.Quality, "mutated")
	got.NextCursor.InternalProductID = 99
	got, err = catalog.ListCatalogProductFacts(context.Background(), internalreadports.Cursor{}, 50)
	if err != nil {
		t.Fatal(err)
	}
	item := got.Items[0]
	if *item.Reference != reference || *item.Description != description || *item.EAN != ean || *item.SellableStock.Quantity != quantity || item.SellableStock.Quality[0] != "known" || len(item.SellableStock.Quality) != 1 || *item.CurrentPrice.Amount != amount || item.CurrentPrice.Quality[0] != "known" || len(item.CurrentPrice.Quality) != 1 || *item.Cost.Amount != amount || item.Cost.Quality[0] != "known" || len(item.Cost.Quality) != 1 || got.NextCursor.InternalProductID != 7 {
		t.Fatalf("catalog cache was aliased: %+v", got)
	}

	observed := clock.Now()
	costAmount, taxICMS, taxIPI, taxPIS, taxCOFINS, stockQuantity := 1.0, 2.0, 3.0, 4.0, 5.0, 6.0
	cost := &internalreaddomain.CostAsOf{Amount: &costAmount, Source: internalreaddomain.SourceMetadata{System: "oracle", ObservedAt: &observed}, QualityFlags: make([]internalreaddomain.QualityFlag, 1, 2)}
	tax := &internalreaddomain.TaxInputs{ICMSAmount: &taxICMS, IPIAmount: &taxIPI, PISAmount: &taxPIS, COFINSAmount: &taxCOFINS, Source: internalreaddomain.SourceMetadata{System: "oracle", ObservedAt: &observed}, QualityFlags: make([]internalreaddomain.QualityFlag, 1, 2)}
	stock := &internalreaddomain.StockFact{Quantity: &stockQuantity, Source: internalreaddomain.SourceMetadata{System: "oracle", ObservedAt: &observed}, QualityFlags: make([]internalreaddomain.QualityFlag, 1, 2)}
	cost.QualityFlags[0], tax.QualityFlags[0], stock.QualityFlags[0] = internalreaddomain.QualityComplete, internalreaddomain.QualityComplete, internalreaddomain.QualityComplete
	batch := NewBatchReader(&mutableBatchReader{cost: map[int64]*internalreaddomain.CostAsOf{1: cost}, tax: map[int64]*internalreaddomain.TaxInputs{1: tax}}, testCache(clock, 10, nil))
	stockReader := NewStockBatchReader(&mutableStockReader{facts: map[int64]*internalreaddomain.StockFact{1: stock}}, batch.cache)
	costResult, _ := batch.GetCostFactsByIDs(context.Background(), []int64{1})
	taxResult, _ := batch.GetTaxFactsByIDs(context.Background(), []int64{1})
	stockResult, _ := stockReader.GetStockFactsByIDs(context.Background(), []int64{1})
	*costResult[1].Amount, *taxResult[1].ICMSAmount, *taxResult[1].IPIAmount, *taxResult[1].PISAmount, *taxResult[1].COFINSAmount, *stockResult[1].Quantity = 99, 99, 99, 99, 99, 99
	costResult[1].Source.System, taxResult[1].Source.System, stockResult[1].Source.System = "mutated", "mutated", "mutated"
	*costResult[1].Source.ObservedAt, *taxResult[1].Source.ObservedAt, *stockResult[1].Source.ObservedAt = time.Time{}, time.Time{}, time.Time{}
	costResult[1].QualityFlags[0], taxResult[1].QualityFlags[0], stockResult[1].QualityFlags[0] = "mutated", "mutated", "mutated"
	costResult[1].QualityFlags = append(costResult[1].QualityFlags, "mutated")
	taxResult[1].QualityFlags = append(taxResult[1].QualityFlags, "mutated")
	stockResult[1].QualityFlags = append(stockResult[1].QualityFlags, "mutated")
	costResult, _ = batch.GetCostFactsByIDs(context.Background(), []int64{1})
	taxResult, _ = batch.GetTaxFactsByIDs(context.Background(), []int64{1})
	stockResult, _ = stockReader.GetStockFactsByIDs(context.Background(), []int64{1})
	if *costResult[1].Amount != costAmount || costResult[1].Source.System != "oracle" || !costResult[1].Source.ObservedAt.Equal(observed) || costResult[1].QualityFlags[0] != internalreaddomain.QualityComplete || len(costResult[1].QualityFlags) != 1 || *taxResult[1].ICMSAmount != taxICMS || *taxResult[1].IPIAmount != taxIPI || *taxResult[1].PISAmount != taxPIS || *taxResult[1].COFINSAmount != taxCOFINS || taxResult[1].Source.System != "oracle" || taxResult[1].QualityFlags[0] != internalreaddomain.QualityComplete || len(taxResult[1].QualityFlags) != 1 || *stockResult[1].Quantity != stockQuantity || stockResult[1].Source.System != "oracle" || stockResult[1].QualityFlags[0] != internalreaddomain.QualityComplete || len(stockResult[1].QualityFlags) != 1 {
		t.Fatalf("batch cache was aliased: cost=%+v tax=%+v stock=%+v", costResult[1], taxResult[1], stockResult[1])
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
