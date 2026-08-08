package internalread

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	erpports "marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var _ readports.Reader = (*Reader)(nil)
var _ readports.CatalogPageReader = (*Reader)(nil)

var importedAt = time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
var fetchedAt = time.Date(2026, 7, 11, 9, 0, 0, 0, time.UTC)

func ptr[T any](v T) *T { return &v }

func mirrorRows() []erpdomain.MirrorProduct {
	return []erpdomain.MirrorProduct{
		{CodigoProduto: "1", Descricao: ptr("Blue Widget"), Custo: ptr("12.30"), EstoqueTotal: ptr("7"), EAN: ptr("7894900011517"), Referencia: ptr(" REF-1 "), Marca: ptr("Acme"), NCM: ptr("12345678"), ImportedAt: timePtr(importedAt), UpdatedAt: importedAt.Add(time.Hour)},
		{CodigoProduto: "2", Descricao: ptr("Blue Spare"), Custo: ptr("3.25"), EstoqueTotal: ptr("4"), EAN: ptr("7894900011517"), ImportedAt: timePtr(importedAt.Add(time.Minute)), UpdatedAt: importedAt.Add(time.Hour)},
		{CodigoProduto: "10", Descricao: ptr("Red Widget"), Custo: ptr("7"), EAN: ptr("bad"), ImportedAt: timePtr(importedAt.Add(2 * time.Minute)), UpdatedAt: importedAt.Add(time.Hour)},
	}
}

func timePtr(value time.Time) *time.Time { return &value }

type fakeRepo struct {
	rows            []erpdomain.MirrorProduct
	err             error
	tenant          string
	source          erpdomain.ImportSource
	snapshotCalled  bool
	mirrorRowsCalls int
}

func (f *fakeRepo) PersistSnapshotAtomically(context.Context, string, erpdomain.ImportSnapshot) error {
	return nil
}
func (f *fakeRepo) FindByFileSHA256(context.Context, string, erpdomain.FileSHA256) (*erpdomain.ImportReport, error) {
	return nil, nil
}
func (f *fakeRepo) ListImports(context.Context, string) ([]erpdomain.ImportReport, error) {
	return nil, nil
}
func (f *fakeRepo) GetImport(context.Context, string, erpdomain.ImportID) (erpdomain.ImportReport, error) {
	return erpdomain.ImportReport{}, nil
}
func (f *fakeRepo) LatestCompletedSnapshot(context.Context, string, erpdomain.ImportSource) (erpdomain.ImportSnapshot, error) {
	f.snapshotCalled = true
	panic("reader must not rescan snapshots")
}
func (f *fakeRepo) SyncLatestCompletedSnapshot(context.Context, string, erpdomain.ImportSource) (int, error) {
	return 0, nil
}
func (f *fakeRepo) record(tenant string, source erpdomain.ImportSource) {
	f.tenant, f.source = tenant, source
}
func (f *fakeRepo) MirrorRows(_ context.Context, tenant string, source erpdomain.ImportSource) ([]erpdomain.MirrorProduct, error) {
	f.record(tenant, source)
	f.mirrorRowsCalls++
	return f.rows, f.err
}
func (f *fakeRepo) MirrorProductByCode(_ context.Context, tenant string, source erpdomain.ImportSource, code string) (erpdomain.MirrorProduct, bool, error) {
	f.record(tenant, source)
	if f.err != nil {
		return erpdomain.MirrorProduct{}, false, f.err
	}
	for _, row := range f.rows {
		if strings.TrimSpace(row.CodigoProduto) == code {
			return row, true, nil
		}
	}
	return erpdomain.MirrorProduct{}, false, nil
}
func (f *fakeRepo) MirrorProductsByCodes(_ context.Context, tenant string, source erpdomain.ImportSource, codes []string) ([]erpdomain.MirrorProduct, error) {
	f.record(tenant, source)
	if f.err != nil {
		return nil, f.err
	}
	wanted := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		wanted[code] = struct{}{}
	}
	rows := make([]erpdomain.MirrorProduct, 0, len(codes))
	for _, row := range f.rows {
		if _, ok := wanted[strings.TrimSpace(row.CodigoProduto)]; ok {
			rows = append(rows, row)
		}
	}
	return rows, nil
}
func (f *fakeRepo) MirrorCatalogPage(_ context.Context, tenant string, source erpdomain.ImportSource, query string, after int64, limit int, policy erpports.MirrorAssortmentPolicy) ([]erpdomain.MirrorProduct, error) {
	f.record(tenant, source)
	if f.err != nil {
		return nil, f.err
	}
	rows := make([]erpdomain.MirrorProduct, 0)
	for _, row := range f.rows {
		id, err := strconv.ParseInt(strings.TrimSpace(row.CodigoProduto), 10, 64)
		if err != nil || id <= after || id <= 0 || query != "" && (row.Descricao == nil || !strings.Contains(strings.ToLower(*row.Descricao), strings.ToLower(query))) {
			continue
		}
		eligible, err := matchesSellableAssortment(row, SellableAssortmentPolicy(policy))
		if err != nil || !eligible {
			continue
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		a, _ := strconv.ParseInt(rows[i].CodigoProduto, 10, 64)
		b, _ := strconv.ParseInt(rows[j].CodigoProduto, 10, 64)
		return a < b
	})
	if limit > 0 && len(rows) > limit+1 {
		rows = rows[:limit+1]
	}
	return rows, nil
}
func (f *fakeRepo) MirrorCatalogAssortmentCounts(_ context.Context, tenant string, source erpdomain.ImportSource, policy erpports.MirrorAssortmentPolicy) (int, int, error) {
	f.record(tenant, source)
	sellable := 0
	for _, row := range f.rows {
		eligible, err := matchesSellableAssortment(row, SellableAssortmentPolicy(policy))
		if err == nil && eligible {
			sellable++
		}
	}
	return sellable, len(f.rows), nil
}
func (f *fakeRepo) MirrorEANCollisionCounts(_ context.Context, tenant string, source erpdomain.ImportSource, _ erpports.MirrorAssortmentPolicy) (map[string]int, error) {
	f.record(tenant, source)
	if f.err != nil {
		return nil, f.err
	}
	return map[string]int{"7894900011517": 2}, nil
}

var _ erpports.ImportRepository = (*fakeRepo)(nil)

func readerWith(rows []erpdomain.MirrorProduct) (*Reader, *fakeRepo, context.Context) {
	repo := &fakeRepo{rows: rows}
	reader := NewReader(repo, "tenant-a", WithClock(func() time.Time { return fetchedAt }))
	return reader, repo, WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente)
}

func TestMirrorBackedReads(t *testing.T) {
	r, repo, ctx := readerWith(mirrorRows())
	candidates, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{Title: ptr("blue")})
	if err != nil || len(candidates) != 2 || candidates[0].ProductID != 1 || candidates[0].Source.ObservedAt == nil || !candidates[0].Source.ObservedAt.Equal(importedAt) || !readdomain.HasQualityFlag(candidates[0].QualityFlags, readdomain.QualityEANCollision) || !readdomain.HasQualityFlag(candidates[0].QualityFlags, readdomain.QualityAmbiguousProduct) {
		t.Fatalf("candidates=%+v err=%v", candidates, err)
	}
	stock, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1})
	if err != nil || stock.Quantity == nil || *stock.Quantity != 7 || stock.Source.System != string(erpdomain.SourceCatalogoCliente) || stock.Source.ObservedAt == nil || !stock.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("stock=%+v err=%v", stock, err)
	}
	cost, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt}})
	if err != nil || cost.Amount == nil || *cost.Amount != 12.3 || cost.AmountScope != readdomain.CostAmountScopePerUnit || cost.Source.ObservedAt == nil || !cost.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("cost=%+v err=%v", cost, err)
	}
	tax, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1})
	if err != nil || tax.NCM == nil || *tax.NCM != "12345678" || tax.Source.ObservedAt == nil || !tax.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("tax=%+v err=%v", tax, err)
	}
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 2, noCutPolicy())
	if err != nil || len(page.Items) != 2 || page.Items[0].InternalProductID != 1 || page.Items[1].InternalProductID != 2 || page.NextCursor == nil || page.NextCursor.InternalProductID != 2 || !page.AsOf.Equal(importedAt.Add(time.Minute)) {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	page2, err := r.ListCatalogProductFacts(ctx, *page.NextCursor, 2, noCutPolicy())
	if err != nil || len(page2.Items) != 1 || page2.Items[0].InternalProductID != 10 || page2.NextCursor != nil {
		t.Fatalf("page2=%+v err=%v", page2, err)
	}
	if repo.snapshotCalled || repo.tenant != "tenant-a" || repo.source != erpdomain.SourceCatalogoCliente {
		t.Fatalf("repo calls snapshot=%v tenant=%q source=%q", repo.snapshotCalled, repo.tenant, repo.source)
	}
}

func TestFindProductsForLinkingAppliesAssortmentAfterCachedIndexLoads(t *testing.T) {
	rows := []erpdomain.MirrorProduct{
		{CodigoProduto: "101", Usoprod: ptr("V"), EstoqueTotal: ptr("5"), ADEcommerce: ptr("S"), ImportedAt: timePtr(importedAt)},
		{CodigoProduto: "102", Usoprod: ptr("R"), EstoqueTotal: ptr("0"), ADEcommerce: ptr("S"), ImportedAt: timePtr(importedAt)},
		{CodigoProduto: "103", Usoprod: ptr("R"), EstoqueTotal: ptr("5"), ADEcommerce: ptr("N"), ImportedAt: timePtr(importedAt)},
		{CodigoProduto: "104", ImportedAt: timePtr(importedAt)},
		{CodigoProduto: "105", Usoprod: ptr("R"), EstoqueTotal: ptr("5"), ADEcommerce: ptr("S"), ImportedAt: timePtr(importedAt)},
	}
	r, repo, _ := readerWith(rows)
	strict := WithSellableAssortment(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), true, true, true)
	got, err := r.FindProductsForLinking(strict, readports.FindProductsInput{})
	if err != nil {
		t.Fatal(err)
	}
	if ids := candidateIDsForReaderTest(got); !reflect.DeepEqual(ids, []int{104, 105}) {
		t.Fatalf("strict candidate IDs = %v, want [104 105]", ids)
	}

	relaxed := WithSellableAssortment(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), false, false, false)
	got, err = r.FindProductsForLinking(relaxed, readports.FindProductsInput{})
	if err != nil {
		t.Fatal(err)
	}
	if ids := candidateIDsForReaderTest(got); !reflect.DeepEqual(ids, []int{101, 102, 103, 104, 105}) {
		t.Fatalf("relaxed candidate IDs = %v, want [101 102 103 104 105]", ids)
	}
	if repo.mirrorRowsCalls != 1 {
		t.Fatalf("cached index mirror row calls = %d, want 1", repo.mirrorRowsCalls)
	}
}

func TestFindProductsForLinkingComputesEANCollisionAfterAssortmentCut(t *testing.T) {
	ean := ptr("7894900011517")
	rows := []erpdomain.MirrorProduct{
		{CodigoProduto: "201", EAN: ean, Usoprod: ptr("R"), EstoqueTotal: ptr("5"), ADEcommerce: ptr("S"), ImportedAt: timePtr(importedAt)},
		{CodigoProduto: "202", EAN: ean, Usoprod: ptr("V"), EstoqueTotal: ptr("5"), ADEcommerce: ptr("S"), ImportedAt: timePtr(importedAt)},
	}
	r, _, _ := readerWith(rows)
	strict := WithSellableAssortment(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), true, true, true)
	got, err := r.FindProductsForLinking(strict, readports.FindProductsInput{EAN: ean})
	if err != nil || len(got) != 1 {
		t.Fatalf("strict collision candidates = %+v err=%v, want one survivor", got, err)
	}
	if readdomain.HasQualityFlag(got[0].QualityFlags, readdomain.QualityEANCollision) || readdomain.HasQualityFlag(got[0].QualityFlags, readdomain.QualityAmbiguousProduct) {
		t.Fatalf("strict survivor quality = %v, want non-ambiguous auto-approval candidate", got[0].QualityFlags)
	}

	relaxed := WithSellableAssortment(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), false, false, false)
	got, err = r.FindProductsForLinking(relaxed, readports.FindProductsInput{EAN: ean})
	if err != nil || len(got) != 2 {
		t.Fatalf("relaxed collision candidates = %+v err=%v, want two twins", got, err)
	}
	for _, candidate := range got {
		if !readdomain.HasQualityFlag(candidate.QualityFlags, readdomain.QualityEANCollision) || !readdomain.HasQualityFlag(candidate.QualityFlags, readdomain.QualityAmbiguousProduct) {
			t.Fatalf("relaxed twin %d quality = %v, want EAN collision and ambiguity", candidate.ProductID, candidate.QualityFlags)
		}
	}
}

func candidateIDsForReaderTest(candidates []readdomain.ProductCandidate) []int {
	ids := make([]int, 0, len(candidates))
	for _, candidate := range candidates {
		ids = append(ids, candidate.ProductID)
	}
	return ids
}

func TestMirrorUnknownValuesRemainUnknown(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	stock, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 10})
	if err != nil || stock.Quantity != nil || !readdomain.HasQualityFlag(stock.QualityFlags, readdomain.QualityMissingStock) {
		t.Fatalf("stock=%+v err=%v", stock, err)
	}
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{InternalProductID: 2}, 10, noCutPolicy())
	if err != nil || page.Items[0].Cost.Amount == nil || *page.Items[0].Cost.Amount != "7" || page.Items[0].SellableStock.Quantity != nil {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	rows := mirrorRows()
	rows[0].Custo = nil
	r, _, ctx = readerWith(rows)
	_, err = r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt}})
	if !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
		t.Fatalf("unknown cost err=%v", err)
	}
}

// The mirror is one observation, not a cost history: an as-of older than the
// snapshot still gets the observed amount, flagged stale so the caller can
// render it as an approximation instead of an exact as-of answer.
func TestCostAsOfOlderThanTheMirrorRowIsFlaggedStale(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	cost, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt.Add(-time.Second)}})
	if err != nil || cost.Amount == nil || *cost.Amount != 12.3 {
		t.Fatalf("cost=%+v err=%v", cost, err)
	}
	if !readdomain.HasQualityFlag(cost.QualityFlags, readdomain.QualityStaleSource) {
		t.Fatalf("cost older than the snapshot must carry stale_source: %+v", cost.QualityFlags)
	}
	if cost.Source.ObservedAt == nil || !cost.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("the observation instant must travel with the amount: %+v", cost.Source)
	}
	exact, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt}})
	if err != nil || readdomain.HasQualityFlag(exact.QualityFlags, readdomain.QualityStaleSource) {
		t.Fatalf("an as-of at the observation instant is not stale: %+v err=%v", exact.QualityFlags, err)
	}
}

func TestLiveReadThroughUsesUpdatedAtAsObservationTime(t *testing.T) {
	// live_read_through source (sankhya): the row carries protocol_id=NULL by
	// design (ImportedAt nil) and updated_at IS the observation instant. sankhya
	// is not yet a ParseActiveSource value — the wiring that routes it through
	// this reader is F1 — so pin it directly to exercise the kind branch.
	syncedAt := time.Date(2026, 7, 12, 8, 0, 0, 0, time.UTC)
	row := erpdomain.MirrorProduct{CodigoProduto: "1", Descricao: ptr("Live Widget"), Custo: ptr("9.90"), EstoqueTotal: ptr("3"), NCM: ptr("12345678"), ImportedAt: nil, UpdatedAt: syncedAt}
	repo := &fakeRepo{rows: []erpdomain.MirrorProduct{row}}
	reader := NewReader(repo, "tenant-a", WithClock(func() time.Time { return fetchedAt }))
	ctx := WithActiveSource(context.Background(), erpdomain.ImportSource("sankhya"))

	cost, err := reader.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: syncedAt}})
	if err != nil || cost.Amount == nil || *cost.Amount != 9.9 || cost.Source.ObservedAt == nil || !cost.Source.ObservedAt.Equal(syncedAt) {
		t.Fatalf("live cost must observe updated_at, not fail on nil imported_at: cost=%+v err=%v", cost, err)
	}
	stale, err := reader.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: syncedAt.Add(-time.Second)}})
	if err != nil || stale.Amount == nil || !readdomain.HasQualityFlag(stale.QualityFlags, readdomain.QualityStaleSource) {
		t.Fatalf("live as-of before the sync instant must answer flagged stale: cost=%+v err=%v", stale, err)
	}
	stock, err := reader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1})
	if err != nil || stock.Source.ObservedAt == nil || !stock.Source.ObservedAt.Equal(syncedAt) {
		t.Fatalf("live stock=%+v err=%v", stock, err)
	}
	page, err := reader.ListCatalogProductFacts(ctx, readports.Cursor{}, 10, noCutPolicy())
	if err != nil || len(page.Items) != 1 || !page.AsOf.Equal(syncedAt) {
		t.Fatalf("live page must set AsOf from updated_at: page=%+v err=%v", page, err)
	}
}

func TestUnknownImportTimeFailsHonestly(t *testing.T) {
	rows := mirrorRows()
	rows[0].ImportedAt = nil
	r, _, ctx := readerWith(rows)
	calls := map[string]func() error{
		"find": func() error {
			_, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: ptr(1)})
			return err
		},
		"stock": func() error {
			_, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1})
			return err
		},
		"cost": func() error {
			_, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt}})
			return err
		},
		"tax": func() error {
			_, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1})
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err := call()
			if !errors.Is(err, ErrNoErpSnapshot) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestCatalogPageUnknownImportTimesFailHonestly(t *testing.T) {
	rows := mirrorRows()
	for i := range rows {
		rows[i].ImportedAt = nil
	}
	r, _, ctx := readerWith(rows)
	_, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 10, noCutPolicy())
	if !errors.Is(err, ErrNoErpSnapshot) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
		t.Fatalf("err=%v", err)
	}
}

func TestMirrorProductNotFound(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	for name, call := range map[string]func() error{
		"find": func() error {
			_, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: ptr(99)})
			return err
		},
		"stock": func() error {
			_, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 99})
			return err
		},
		"cost": func() error { _, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 99}); return err },
		"tax":  func() error { _, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 99}); return err },
	} {
		t.Run(name, func(t *testing.T) {
			var target *ERPProductNotFoundError
			if !errors.As(call(), &target) {
				t.Fatalf("err=%v", call())
			}
		})
	}
}

func TestUnpinnedReadsFailClosed(t *testing.T) {
	r, repo, _ := readerWith(mirrorRows())
	ctx := context.Background()
	calls := map[string]func() error{
		"find": func() error { _, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{}); return err },
		"stock": func() error {
			_, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1})
			return err
		},
		"cost": func() error { _, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1}); return err },
		"tax":  func() error { _, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1}); return err },
		"list": func() error {
			_, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 10, noCutPolicy())
			return err
		},
		"search": func() error {
			_, err := r.SearchCatalogProductFacts(ctx, "blue", readports.Cursor{}, 10, noCutPolicy())
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err := call()
			if !errors.Is(err, ErrUnknownActiveSource) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
				t.Fatalf("err=%v", err)
			}
		})
	}
	if repo.snapshotCalled || repo.tenant != "" {
		t.Fatalf("repository called on unresolved source")
	}
}

func TestSellerSKUMatchesTheProductCode(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	candidates, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{SellerSKU: ptr("10")})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(candidates) != 1 || candidates[0].ProductID != 10 {
		t.Fatalf("candidates=%+v", candidates)
	}
}

// A seller_sku that is not CODPROD-shaped is a legacy REFFORN, not an anchor.
// It must narrow to nothing — and, the trap, must not be read as "no anchor was
// supplied", which would answer with the entire mirror as candidates.
func TestNonCodprodSellerSKUAnchorsNothingAndDoesNotReturnTheWholeMirror(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	candidates, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{SellerSKU: ptr("ZP1704.1.")})
	if err != nil {
		t.Fatalf("an unmatchable seller_sku is an empty result, not an error: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("candidates=%+v", candidates)
	}
}

func TestNonCodprodSellerSKULeavesTheEANAnchorStanding(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	candidates, err := r.FindProductsForLinking(ctx, readports.FindProductsInput{SellerSKU: ptr("ZP1704.1."), EAN: ptr("7894900011517")})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%+v", candidates)
	}
	for _, candidate := range candidates {
		if candidate.ProductID != 1 && candidate.ProductID != 2 {
			t.Fatalf("the dropped seller_sku widened the match: %+v", candidates)
		}
	}
}

func TestUnsupportedQueries(t *testing.T) {
	r, _, _ := readerWith(mirrorRows())
	_, err := r.GetCurrentPrice(context.Background(), readports.CurrentPriceInput{})
	if err == nil || err.Error() != "xlsx snapshot does not contain current price" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
		t.Fatalf("price err=%v", err)
	}
	_, err = r.GetSalesHistory(context.Background(), readports.SalesHistoryInput{})
	if err == nil || err.Error() != "xlsx snapshot does not contain sales history" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
		t.Fatalf("sales err=%v", err)
	}
}

func TestParseActiveSource(t *testing.T) {
	for _, tc := range []struct {
		raw              string
		want             erpdomain.ImportSource
		present, invalid bool
	}{{raw: ""}, {raw: "   "}, {raw: "xlsx", want: erpdomain.SourceXLSX, present: true}, {raw: " catalogo_cliente ", want: erpdomain.SourceCatalogoCliente, present: true}, {raw: "sankhya", invalid: true}} {
		source, present, err := ParseActiveSource(tc.raw)
		if tc.invalid {
			if !errors.Is(err, ErrUnknownActiveSource) {
				t.Fatalf("raw=%q err=%v", tc.raw, err)
			}
			continue
		}
		if err != nil || source != tc.want || present != tc.present {
			t.Fatalf("raw=%q got=(%q,%v,%v)", tc.raw, source, present, err)
		}
	}
}

func TestCatalogFactQuality(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 10, noCutPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(page.Items[0].CurrentPrice.Quality, []string{string(readdomain.QualityMissingPrice)}) || page.Items[2].EAN != nil || !reflect.DeepEqual(page.Items[2].QualityFlags, []string{string(readdomain.QualityInvalidEAN)}) {
		t.Fatalf("items=%+v", page.Items)
	}
}
