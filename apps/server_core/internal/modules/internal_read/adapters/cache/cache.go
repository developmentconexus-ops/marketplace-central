package cache

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	inventoryports "marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

const (
	ClassCatalog   = "catalog"
	ClassInventory = "inventory"
	ClassPriceCost = "pricecost"

	defaultCatalogTTL = 5 * time.Minute
	defaultStockTTL   = 45 * time.Second
	defaultPriceTTL   = 2 * time.Minute
	defaultMaxEntries = 100_000
)

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now().UTC() }

type Config struct {
	Policies   map[string]domain.FreshnessPolicy
	MaxEntries int
	Clock      Clock
}

func ConfigFromEnv(getenv func(string) string) Config {
	if getenv == nil {
		getenv = os.Getenv
	}
	return Config{
		Policies: map[string]domain.FreshnessPolicy{
			ClassCatalog:   {MaxAge: envDuration(getenv, "MPC_CACHE_TTL_CATALOG", defaultCatalogTTL)},
			ClassInventory: {MaxAge: envDuration(getenv, "MPC_CACHE_TTL_STOCK", defaultStockTTL)},
			ClassPriceCost: {MaxAge: envDuration(getenv, "MPC_CACHE_TTL_PRICECOST", defaultPriceTTL)},
		},
		MaxEntries: envPositiveInt(getenv, "MPC_CACHE_MAX_ENTRIES", defaultMaxEntries),
		Clock:      systemClock{},
	}
}

func envDuration(getenv func(string) string, name string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func envPositiveInt(getenv func(string) string, name string, fallback int) int {
	value := strings.TrimSpace(getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

type entry struct {
	class    string
	key      string
	value    any
	created  time.Time
	snapshot time.Time
}

type Cache struct {
	mu         sync.Mutex
	entries    map[string]*list.Element
	lru        *list.List
	group      singleflight.Group
	generations map[string]uint64
	policies   map[string]domain.FreshnessPolicy
	maxEntries int
	clock      Clock
}

var _ internalreadports.CacheInvalidator = (*Cache)(nil)

func New(config Config) *Cache {
	clock := config.Clock
	if clock == nil {
		clock = systemClock{}
	}
	policies := make(map[string]domain.FreshnessPolicy, len(config.Policies))
	for class, policy := range config.Policies {
		policies[class] = policy
	}
	if _, ok := policies[ClassCatalog]; !ok {
		policies[ClassCatalog] = domain.FreshnessPolicy{MaxAge: defaultCatalogTTL}
	}
	if _, ok := policies[ClassInventory]; !ok {
		policies[ClassInventory] = domain.FreshnessPolicy{MaxAge: defaultStockTTL}
	}
	if _, ok := policies[ClassPriceCost]; !ok {
		policies[ClassPriceCost] = domain.FreshnessPolicy{MaxAge: defaultPriceTTL}
	}
	maxEntries := config.MaxEntries
	if maxEntries <= 0 {
		maxEntries = defaultMaxEntries
	}
	return &Cache{
		entries:    make(map[string]*list.Element),
		lru:        list.New(),
		generations: make(map[string]uint64),
		policies:   policies,
		maxEntries: maxEntries,
		clock:      clock,
	}
}

func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

func (c *Cache) InvalidateClass(factClass string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generations[factClass]++
	for key, element := range c.entries {
		if element.Value.(*entry).class == factClass {
			delete(c.entries, key)
			c.lru.Remove(element)
		}
	}
}

func (c *Cache) load(ctx context.Context, class, key string, loader func() (any, time.Time, error)) (any, error) {
	policy := c.policies[class]
	bypass := false
	if requested, ok := domain.FreshnessPolicyFromContext(ctx); ok && requested.IsZero() {
		bypass = true
	}
	if !bypass {
		if value, ok := c.lookup(class, key, policy.MaxAge); ok {
			cacheLog("hit", class)
			return value, nil
		}
		cacheLog("miss", class)
	} else {
		cacheLog("bypass", class)
	}

	groupKey := key
	if bypass {
		groupKey = "bypass:" + key
	}
	done := c.group.DoChan(groupKey, func() (any, error) {
		generation := c.classGeneration(class)
		if !bypass {
			if cached, ok := c.lookup(class, key, policy.MaxAge); ok {
				return cached, nil
			}
		}
		loaded, snapshot, err := loader()
		if err != nil {
			return nil, err
		}
		c.storeIfCurrent(class, key, loaded, snapshot, generation)
		return loaded, nil
	})
	select {
	case result := <-done:
		return result.Val, result.Err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Cache) classGeneration(class string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.generations[class]
}

func (c *Cache) lookup(class, key string, ttl time.Duration) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	item := element.Value.(*entry)
	if item.class != class || ttl <= 0 || c.clock.Now().Sub(item.created) >= ttl {
		delete(c.entries, key)
		c.lru.Remove(element)
		return nil, false
	}
	c.lru.MoveToFront(element)
	return cloneCachedValue(item.value), true
}

func (c *Cache) storeIfCurrent(class, key string, value any, snapshot time.Time, generation uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.generations[class] != generation {
		return
	}
	value = cloneCachedValue(value)
	if existing, ok := c.entries[key]; ok {
		existing.Value.(*entry).value = value
		existing.Value.(*entry).created = c.clock.Now()
		existing.Value.(*entry).snapshot = snapshot
		existing.Value.(*entry).class = class
		c.lru.MoveToFront(existing)
		return
	}
	for len(c.entries) >= c.maxEntries {
		oldest := c.lru.Back()
		if oldest == nil {
			break
		}
		item := oldest.Value.(*entry)
		delete(c.entries, item.key)
		c.lru.Remove(oldest)
	}
	item := &entry{class: class, key: key, value: value, created: c.clock.Now(), snapshot: snapshot}
	c.entries[key] = c.lru.PushFront(item)
}

func cacheLog(status, class string) {
	slog.Default().Info("internal_read.cache", "cache", status, "key_class", class)
}

func canonicalKey(method string, parts ...string) string {
	var builder strings.Builder
	builder.WriteString(method)
	for _, part := range parts {
		builder.WriteByte(0)
		builder.WriteString(strconv.Itoa(len(part)))
		builder.WriteByte(':')
		builder.WriteString(part)
	}
	digest := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(digest[:])
}

func canonicalIDs(ids []int64) string {
	ordered := append([]int64(nil), ids...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i] < ordered[j] })
	parts := make([]string, len(ordered))
	for i, id := range ordered {
		parts[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(parts, ",")
}

type CatalogPageReader struct {
	downstream internalreadports.CatalogPageReader
	cache      *Cache
}

var _ internalreadports.CatalogPageReader = CatalogPageReader{}

func NewCatalogPageReader(downstream internalreadports.CatalogPageReader, cache *Cache) CatalogPageReader {
	return CatalogPageReader{downstream: downstream, cache: cache}
}

func (r CatalogPageReader) ListCatalogProductFacts(ctx context.Context, cursor internalreadports.Cursor, limit int) (internalreadports.CatalogFactPage, error) {
	key := canonicalKey("ListCatalogProductFacts", strconv.FormatInt(cursor.InternalProductID, 10), strconv.Itoa(limit))
	value, err := r.cache.load(ctx, ClassCatalog, key, func() (any, time.Time, error) {
		page, err := r.downstream.ListCatalogProductFacts(ctx, cursor, limit)
		return cloneCatalogPage(page), page.AsOf, err
	})
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return cloneCatalogPage(value.(internalreadports.CatalogFactPage)), nil
}

func (r CatalogPageReader) SearchCatalogProductFacts(ctx context.Context, query string, limit int) (internalreadports.CatalogFactPage, error) {
	key := canonicalKey("SearchCatalogProductFacts", query, strconv.Itoa(limit))
	value, err := r.cache.load(ctx, ClassCatalog, key, func() (any, time.Time, error) {
		page, err := r.downstream.SearchCatalogProductFacts(ctx, query, limit)
		return cloneCatalogPage(page), page.AsOf, err
	})
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return cloneCatalogPage(value.(internalreadports.CatalogFactPage)), nil
}

type Reader struct {
	internalreadports.Reader
	pages internalreadports.CatalogPageReader
}

var _ internalreadports.Reader = Reader{}
var _ internalreadports.CatalogPageReader = Reader{}

func NewReader(downstream internalreadports.Reader, cache *Cache) Reader {
	pages, _ := downstream.(internalreadports.CatalogPageReader)
	return Reader{Reader: downstream, pages: NewCatalogPageReader(pages, cache)}
}

func (r Reader) ListCatalogProductFacts(ctx context.Context, cursor internalreadports.Cursor, limit int) (internalreadports.CatalogFactPage, error) {
	return r.pages.ListCatalogProductFacts(ctx, cursor, limit)
}

func (r Reader) SearchCatalogProductFacts(ctx context.Context, query string, limit int) (internalreadports.CatalogFactPage, error) {
	return r.pages.SearchCatalogProductFacts(ctx, query, limit)
}

type BatchReader struct {
	downstream internalreadports.BatchReader
	cache      *Cache
}

var _ internalreadports.BatchReader = BatchReader{}

func NewBatchReader(downstream internalreadports.BatchReader, cache *Cache) BatchReader {
	return BatchReader{downstream: downstream, cache: cache}
}

func (r BatchReader) GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.CostAsOf, error) {
	key := canonicalKey("GetCostFactsByIDs", canonicalIDs(ids))
	value, err := r.cache.load(ctx, ClassPriceCost, key, func() (any, time.Time, error) {
		facts, err := r.downstream.GetCostFactsByIDs(ctx, ids)
		return cloneCostFacts(facts), r.cache.clock.Now(), err
	})
	if err != nil {
		return nil, err
	}
	return cloneCostFacts(value.(map[int64]*domain.CostAsOf)), nil
}

func (r BatchReader) GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.TaxInputs, error) {
	key := canonicalKey("GetTaxFactsByIDs", canonicalIDs(ids))
	value, err := r.cache.load(ctx, ClassPriceCost, key, func() (any, time.Time, error) {
		facts, err := r.downstream.GetTaxFactsByIDs(ctx, ids)
		return cloneTaxFacts(facts), r.cache.clock.Now(), err
	})
	if err != nil {
		return nil, err
	}
	return cloneTaxFacts(value.(map[int64]*domain.TaxInputs)), nil
}

type StockBatchReader struct {
	downstream inventoryports.InternalStockBatchReader
	cache      *Cache
}

var _ inventoryports.InternalStockBatchReader = StockBatchReader{}

func NewStockBatchReader(downstream inventoryports.InternalStockBatchReader, cache *Cache) StockBatchReader {
	return StockBatchReader{downstream: downstream, cache: cache}
}

func (r StockBatchReader) GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.StockFact, error) {
	key := canonicalKey("GetStockFactsByIDs", canonicalIDs(ids))
	value, err := r.cache.load(ctx, ClassInventory, key, func() (any, time.Time, error) {
		facts, err := r.downstream.GetStockFactsByIDs(ctx, ids)
		return cloneStockFacts(facts), r.cache.clock.Now(), err
	})
	if err != nil {
		return nil, err
	}
	return cloneStockFacts(value.(map[int64]*domain.StockFact)), nil
}

func cloneCatalogPage(page internalreadports.CatalogFactPage) internalreadports.CatalogFactPage {
	if page.Items != nil {
		items := make([]internalreadports.CatalogProductFact, len(page.Items))
		copy(items, page.Items)
		page.Items = items
		for i := range page.Items {
			item := &page.Items[i]
			item.Reference = cloneString(item.Reference)
			item.Description = cloneString(item.Description)
			item.EAN = cloneString(item.EAN)
			item.SellableStock.Quantity = cloneFloat64(item.SellableStock.Quantity)
			item.SellableStock.Quality = cloneStrings(item.SellableStock.Quality)
			item.CurrentPrice.Amount = cloneString(item.CurrentPrice.Amount)
			item.CurrentPrice.Quality = cloneStrings(item.CurrentPrice.Quality)
			item.Cost.Amount = cloneString(item.Cost.Amount)
			item.Cost.Quality = cloneStrings(item.Cost.Quality)
		}
	}
	if page.NextCursor != nil {
		cursor := *page.NextCursor
		page.NextCursor = &cursor
	}
	return page
}

func cloneCostFacts(source map[int64]*domain.CostAsOf) map[int64]*domain.CostAsOf {
	if source == nil {
		return nil
	}
	target := make(map[int64]*domain.CostAsOf, len(source))
	for id, fact := range source {
		if fact == nil {
			target[id] = nil
			continue
		}
		copy := *fact
		copy.Amount = cloneFloat64(fact.Amount)
		copy.Source = cloneSourceMetadata(fact.Source)
		copy.QualityFlags = cloneQualityFlags(fact.QualityFlags)
		target[id] = &copy
	}
	return target
}

func cloneTaxFacts(source map[int64]*domain.TaxInputs) map[int64]*domain.TaxInputs {
	if source == nil {
		return nil
	}
	target := make(map[int64]*domain.TaxInputs, len(source))
	for id, fact := range source {
		if fact == nil {
			target[id] = nil
			continue
		}
		copy := *fact
		copy.ICMSAmount = cloneFloat64(fact.ICMSAmount)
		copy.IPIAmount = cloneFloat64(fact.IPIAmount)
		copy.PISAmount = cloneFloat64(fact.PISAmount)
		copy.COFINSAmount = cloneFloat64(fact.COFINSAmount)
		copy.Source = cloneSourceMetadata(fact.Source)
		copy.QualityFlags = cloneQualityFlags(fact.QualityFlags)
		target[id] = &copy
	}
	return target
}

func cloneStockFacts(source map[int64]*domain.StockFact) map[int64]*domain.StockFact {
	if source == nil {
		return nil
	}
	target := make(map[int64]*domain.StockFact, len(source))
	for id, fact := range source {
		if fact == nil {
			target[id] = nil
			continue
		}
		copy := *fact
		copy.Quantity = cloneFloat64(fact.Quantity)
		copy.Source = cloneSourceMetadata(fact.Source)
		copy.QualityFlags = cloneQualityFlags(fact.QualityFlags)
		target[id] = &copy
	}
	return target
}

func cloneCachedValue(value any) any {
	switch value := value.(type) {
	case internalreadports.CatalogFactPage:
		return cloneCatalogPage(value)
	case map[int64]*domain.CostAsOf:
		return cloneCostFacts(value)
	case map[int64]*domain.TaxInputs:
		return cloneTaxFacts(value)
	case map[int64]*domain.StockFact:
		return cloneStockFacts(value)
	default:
		return value
	}
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneFloat64(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneStrings(source []string) []string {
	if source == nil {
		return nil
	}
	target := make([]string, len(source))
	copy(target, source)
	return target
}

func cloneQualityFlags(source []domain.QualityFlag) []domain.QualityFlag {
	if source == nil {
		return nil
	}
	target := make([]domain.QualityFlag, len(source))
	copy(target, source)
	return target
}

func cloneSourceMetadata(source domain.SourceMetadata) domain.SourceMetadata {
	source.ObservedAt = cloneTime(source.ObservedAt)
	return source
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
