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
		{CodigoProduto: "1", Descricao: ptr("Blue Widget"), Custo: ptr("12.30"), EstoqueTotal: ptr("7"), EAN: ptr("7894900011517"), Referencia: ptr(" REF-1 "), Marca: ptr("Acme"), NCM: ptr("12345678"), UpdatedAt: importedAt},
		{CodigoProduto: "2", Descricao: ptr("Blue Spare"), Custo: ptr("3.25"), EstoqueTotal: ptr("4"), EAN: ptr("7894900011517"), UpdatedAt: importedAt.Add(time.Minute)},
		{CodigoProduto: "10", Descricao: ptr("Red Widget"), Custo: ptr("7"), EAN: ptr("bad"), UpdatedAt: importedAt.Add(2 * time.Minute)},
	}
}

type fakeRepo struct {
	rows           []erpdomain.MirrorProduct
	err            error
	tenant         string
	source         erpdomain.ImportSource
	snapshotCalled bool
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
func (f *fakeRepo) MirrorCatalogPage(_ context.Context, tenant string, source erpdomain.ImportSource, query string, after int64, limit int) ([]erpdomain.MirrorProduct, error) {
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
func (f *fakeRepo) MirrorEANCollisionCounts(_ context.Context, tenant string, source erpdomain.ImportSource) (map[string]int, error) {
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
	if err != nil || len(candidates) != 2 || candidates[0].ProductID != 1 || !readdomain.HasQualityFlag(candidates[0].QualityFlags, readdomain.QualityEANCollision) || !readdomain.HasQualityFlag(candidates[0].QualityFlags, readdomain.QualityAmbiguousProduct) {
		t.Fatalf("candidates=%+v err=%v", candidates, err)
	}
	stock, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1})
	if err != nil || stock.Quantity == nil || *stock.Quantity != 7 || stock.Source.System != string(erpdomain.SourceCatalogoCliente) || stock.Source.ObservedAt == nil || !stock.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("stock=%+v err=%v", stock, err)
	}
	cost, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt}})
	if err != nil || cost.Amount == nil || *cost.Amount != 12.3 || cost.AmountScope != readdomain.CostAmountScopePerUnit {
		t.Fatalf("cost=%+v err=%v", cost, err)
	}
	tax, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1})
	if err != nil || tax.NCM == nil || *tax.NCM != "12345678" {
		t.Fatalf("tax=%+v err=%v", tax, err)
	}
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 2)
	if err != nil || len(page.Items) != 2 || page.Items[0].InternalProductID != 1 || page.Items[1].InternalProductID != 2 || page.NextCursor == nil || page.NextCursor.InternalProductID != 2 || !page.AsOf.Equal(importedAt.Add(time.Minute)) {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	page2, err := r.ListCatalogProductFacts(ctx, *page.NextCursor, 2)
	if err != nil || len(page2.Items) != 1 || page2.Items[0].InternalProductID != 10 || page2.NextCursor != nil {
		t.Fatalf("page2=%+v err=%v", page2, err)
	}
	if repo.snapshotCalled || repo.tenant != "tenant-a" || repo.source != erpdomain.SourceCatalogoCliente {
		t.Fatalf("repo calls snapshot=%v tenant=%q source=%q", repo.snapshotCalled, repo.tenant, repo.source)
	}
}

func TestMirrorUnknownValuesRemainUnknown(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	stock, err := r.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 10})
	if err != nil || stock.Quantity != nil || !readdomain.HasQualityFlag(stock.QualityFlags, readdomain.QualityMissingStock) {
		t.Fatalf("stock=%+v err=%v", stock, err)
	}
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{InternalProductID: 2}, 10)
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

func TestCostAsOfRejectsNewerMirrorRow(t *testing.T) {
	r, _, ctx := readerWith(mirrorRows())
	_, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1, Policy: readdomain.CostAsOfPolicy{EffectiveAt: importedAt.Add(-time.Second)}})
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
		"cost":   func() error { _, err := r.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1}); return err },
		"tax":    func() error { _, err := r.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1}); return err },
		"list":   func() error { _, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 10); return err },
		"search": func() error { _, err := r.SearchCatalogProductFacts(ctx, "blue", 10); return err },
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
	page, err := r.ListCatalogProductFacts(ctx, readports.Cursor{}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(page.Items[0].CurrentPrice.Quality, []string{string(readdomain.QualityMissingPrice)}) || page.Items[2].EAN != nil || !reflect.DeepEqual(page.Items[2].QualityFlags, []string{string(readdomain.QualityInvalidEAN)}) {
		t.Fatalf("items=%+v", page.Items)
	}
}
