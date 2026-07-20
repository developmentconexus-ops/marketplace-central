package internalread

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	erpports "marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var _ readports.Reader = (*Reader)(nil)
var _ readports.CatalogPageReader = (*Reader)(nil)

type fakeRepo struct {
	snapshot erpdomain.ImportSnapshot
	err      error
	tenant   string
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
func (f *fakeRepo) LatestCompletedSnapshot(_ context.Context, tenant string, _ erpdomain.ImportSource) (erpdomain.ImportSnapshot, error) {
	f.tenant = tenant
	return f.snapshot, f.err
}

var _ erpports.ImportRepository = (*fakeRepo)(nil)

var importedAt = time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
var fetchedAt = time.Date(2026, 7, 11, 9, 0, 0, 0, time.UTC)

func ptr[T any](v T) *T { return &v }
func rows() []erpdomain.NormalizedRow {
	return []erpdomain.NormalizedRow{
		{Codprod: "10", Descrprod: "Blue Widget", Custo: "12.30", StockPhysical: "9.5", StockReserved: ptr("2.5"), EAN: ptr("7894900011517"), Refforn: ptr(" REF-10 "), Marca: ptr("Acme"), NCM: ptr("12345678")},
		{Codprod: "20", Descrprod: "Blue Spare", Custo: "3.25", StockPhysical: "4", StockReserved: ptr("0"), EAN: ptr("7894900011517")},
		{Codprod: "30", Descrprod: "Red Widget", Custo: "7", StockPhysical: "5", EAN: ptr("bad")},
	}
}
func readerWith(rs []erpdomain.NormalizedRow) *Reader {
	return NewReader(&fakeRepo{snapshot: erpdomain.ImportSnapshot{Status: erpdomain.ImportStatusCompleted, ImportedAt: importedAt, AcceptedRows: rs}}, "tenant-a", WithClock(func() time.Time { return fetchedAt }))
}

func TestFindProductsForLinking(t *testing.T) {
	tests := []struct {
		name        string
		input       readports.FindProductsInput
		wantIDs     []int
		flags       []readdomain.QualityFlag
		wantMissing bool
	}{
		{"product id", readports.FindProductsInput{ProductID: ptr(10)}, []int{10}, []readdomain.QualityFlag{readdomain.QualityComplete, readdomain.QualityEANCollision}, false},
		{"seller sku", readports.FindProductsInput{SellerSKU: ptr(" 20 ")}, []int{20}, []readdomain.QualityFlag{readdomain.QualityComplete, readdomain.QualityEANCollision}, false},
		{"invalid ean projected nil", readports.FindProductsInput{ProductID: ptr(30)}, []int{30}, []readdomain.QualityFlag{readdomain.QualityInvalidEAN}, false},
		{"multiple title matches ambiguous", readports.FindProductsInput{Title: ptr("blue")}, []int{10, 20}, []readdomain.QualityFlag{readdomain.QualityComplete, readdomain.QualityEANCollision, readdomain.QualityAmbiguousProduct}, false},
		{"requested missing", readports.FindProductsInput{ProductID: ptr(99)}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readerWith(rows()).FindProductsForLinking(context.Background(), tt.input)
			var missing *ERPProductNotFoundError
			if tt.wantMissing {
				if !errors.As(err, &missing) || len(got) != 0 {
					t.Fatalf("got=%v err=%v", got, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			ids := make([]int, len(got))
			for i, c := range got {
				ids[i] = c.ProductID
				if !reflect.DeepEqual(c.QualityFlags, tt.flags) {
					t.Fatalf("flags=%v", c.QualityFlags)
				}
			}
			if !reflect.DeepEqual(ids, tt.wantIDs) {
				t.Fatalf("ids=%v", ids)
			}
			if tt.name == "invalid ean projected nil" && got[0].EAN != nil {
				t.Fatal("invalid EAN retained")
			}
		})
	}
}

func TestFindProductsSkipsNonNumericCodprod(t *testing.T) {
	rs := []erpdomain.NormalizedRow{
		{Codprod: "50", Descrprod: "Green Gadget", Custo: "9", StockPhysical: "3"},
		{Codprod: "SKU-ABC", Descrprod: "Green Gizmo", Custo: "9", StockPhysical: "3"},
	}
	got, err := readerWith(rs).FindProductsForLinking(context.Background(), readports.FindProductsInput{Title: ptr("green")})
	if err != nil {
		t.Fatalf("unrepresentable CODPROD must be skipped, not fail the call: %v", err)
	}
	if len(got) != 1 || got[0].ProductID != 50 {
		t.Fatalf("want only the numeric-CODPROD candidate 50, got %+v", got)
	}
}

func TestSellableStock(t *testing.T) {
	for _, tt := range []struct {
		id       int
		quantity *float64
		flag     readdomain.QualityFlag
		missing  bool
	}{{10, ptr(7.0), "", false}, {20, ptr(4.0), "", false}, {30, nil, readdomain.QualityMissingStock, false}, {99, nil, "", true}} {
		got, err := readerWith(rows()).GetSellableStock(context.Background(), readports.SellableStockInput{ProductID: tt.id})
		var missing *ERPProductNotFoundError
		if tt.missing {
			if !errors.As(err, &missing) {
				t.Fatalf("id=%d err=%v", tt.id, err)
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got.Quantity, tt.quantity) {
			t.Fatalf("id=%d quantity=%v", tt.id, got.Quantity)
		}
		if tt.flag != "" && !readdomain.HasQualityFlag(got.QualityFlags, tt.flag) {
			t.Fatalf("flags=%v", got.QualityFlags)
		}
		if got.Source.ObservedAt == nil || !got.Source.ObservedAt.Equal(importedAt) {
			t.Fatalf("source=%+v", got.Source)
		}
	}
	f := &fakeRepo{err: erpports.ErrImportNotFound}
	_, err := NewReader(f, "tenant-a").GetSellableStock(context.Background(), readports.SellableStockInput{ProductID: 10})
	if !errors.Is(err, ErrNoErpSnapshot) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) || f.tenant != "tenant-a" {
		t.Fatalf("dual err=%v tenant=%q", err, f.tenant)
	}
	_, err = NewReader(&fakeRepo{}, "tenant-a").GetSellableStock(context.Background(), readports.SellableStockInput{ProductID: 10})
	if !errors.Is(err, ErrNoErpSnapshot) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
		t.Fatalf("empty snapshot err=%v", err)
	}
}

func TestUnsupportedQueries(t *testing.T) {
	r := readerWith(rows())
	_, err := r.GetCurrentPrice(context.Background(), readports.CurrentPriceInput{})
	if err == nil || err.Error() != "xlsx snapshot does not contain current price" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
		t.Fatalf("price err=%v", err)
	}
	_, err = r.GetSalesHistory(context.Background(), readports.SalesHistoryInput{})
	if err == nil || err.Error() != "xlsx snapshot does not contain sales history" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
		t.Fatalf("sales err=%v", err)
	}
}

func TestCostAsOf(t *testing.T) {
	policy := readdomain.CostAsOfPolicy{CompanyID: 1, Basis: readdomain.CostBasisCUSSEMICM, EffectiveAt: importedAt}
	got, err := readerWith(rows()).GetCostAsOf(context.Background(), readports.CostAsOfInput{ProductID: 10, Policy: policy})
	if err != nil || got.Amount == nil || *got.Amount != 12.3 || got.AmountScope != readdomain.CostAmountScopePerUnit || got.Source.ObservedAt == nil || !got.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	policy.EffectiveAt = importedAt.Add(-time.Second)
	_, err = readerWith(rows()).GetCostAsOf(context.Background(), readports.CostAsOfInput{ProductID: 10, Policy: policy})
	if !errors.Is(err, ErrNoErpSnapshot) {
		t.Fatalf("as-of err=%v", err)
	}
	policy.EffectiveAt = importedAt
	_, err = readerWith(rows()).GetCostAsOf(context.Background(), readports.CostAsOfInput{ProductID: 99, Policy: policy})
	var missing *ERPProductNotFoundError
	if !errors.As(err, &missing) {
		t.Fatalf("missing err=%v", err)
	}
}

func TestTaxInputs(t *testing.T) {
	for _, tt := range []struct {
		id      int
		ncm     *string
		missing bool
	}{{10, ptr("12345678"), false}, {20, nil, true}} {
		got, err := readerWith(rows()).GetTaxInputs(context.Background(), readports.TaxInput{ProductID: tt.id})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got.NCM, tt.ncm) || got.ICMSAmount != nil || got.IPIAmount != nil || got.PISAmount != nil || got.COFINSAmount != nil {
			t.Fatalf("tax=%+v", got)
		}
		if tt.missing && !readdomain.HasQualityFlag(got.QualityFlags, readdomain.QualityMissingTax) {
			t.Fatalf("flags=%v", got.QualityFlags)
		}
	}
}

func TestCatalogPages(t *testing.T) {
	r := readerWith(rows())
	page, err := r.ListCatalogProductFacts(context.Background(), readports.Cursor{}, 1)
	if err != nil || len(page.Items) != 1 || page.Items[0].InternalProductID != 10 || page.Items[0].Cost.Amount == nil || *page.Items[0].Cost.Amount != "12.30" || page.Items[0].CurrentPrice.Amount != nil || !reflect.DeepEqual(page.Items[0].CurrentPrice.Quality, []string{"missing_price"}) || page.NextCursor == nil || page.NextCursor.InternalProductID != 10 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	page, err = r.ListCatalogProductFacts(context.Background(), *page.NextCursor, 10)
	if err != nil || len(page.Items) != 2 || page.Items[0].InternalProductID != 20 {
		t.Fatalf("page2=%+v err=%v", page, err)
	}
	empty, err := r.SearchCatalogProductFacts(context.Background(), "absent", 10)
	if err != nil || empty.Items == nil || len(empty.Items) != 0 {
		t.Fatalf("empty=%#v err=%v", empty.Items, err)
	}
	none := NewReader(&fakeRepo{err: erpports.ErrImportNotFound}, "tenant-a")
	_, err = none.ListCatalogProductFacts(context.Background(), readports.Cursor{}, 10)
	if !errors.Is(err, ErrNoErpSnapshot) {
		t.Fatalf("err=%v", err)
	}
}

// GOAL C — the active-source toggle selects which imported snapshot the reader
// serves, and ParseActiveSource is the sole validated seam into it.

func TestParseActiveSource(t *testing.T) {
	for _, tc := range []struct {
		raw         string
		wantSource  erpdomain.ImportSource
		wantPresent bool
		wantErr     bool
	}{
		{raw: "", wantPresent: false},
		{raw: "   ", wantPresent: false},
		{raw: "xlsx", wantSource: erpdomain.SourceXLSX, wantPresent: true},
		{raw: " catalogo_cliente ", wantSource: erpdomain.SourceCatalogoCliente, wantPresent: true},
		{raw: "sankhya", wantErr: true},
		{raw: "XLSX", wantErr: true}, // case-sensitive: only the exact dataset keys resolve
	} {
		source, present, err := ParseActiveSource(tc.raw)
		if tc.wantErr {
			if !errors.Is(err, ErrUnknownActiveSource) {
				t.Fatalf("ParseActiveSource(%q) err = %v, want ErrUnknownActiveSource", tc.raw, err)
			}
			continue
		}
		if err != nil || present != tc.wantPresent || source != tc.wantSource {
			t.Fatalf("ParseActiveSource(%q) = (%q,%v,%v), want (%q,%v,nil)", tc.raw, source, present, err, tc.wantSource, tc.wantPresent)
		}
	}
}

// sourceKeyedRepo returns a distinct snapshot per source so a test can prove the
// context source actually drives which dataset the reader reads.
type sourceKeyedRepo struct {
	bySource  map[erpdomain.ImportSource]erpdomain.ImportSnapshot
	gotSource erpdomain.ImportSource
}

func (r *sourceKeyedRepo) PersistSnapshotAtomically(context.Context, string, erpdomain.ImportSnapshot) error {
	return nil
}
func (r *sourceKeyedRepo) FindByFileSHA256(context.Context, string, erpdomain.FileSHA256) (*erpdomain.ImportReport, error) {
	return nil, nil
}
func (r *sourceKeyedRepo) ListImports(context.Context, string) ([]erpdomain.ImportReport, error) {
	return nil, nil
}
func (r *sourceKeyedRepo) GetImport(context.Context, string, erpdomain.ImportID) (erpdomain.ImportReport, error) {
	return erpdomain.ImportReport{}, nil
}
func (r *sourceKeyedRepo) LatestCompletedSnapshot(_ context.Context, _ string, source erpdomain.ImportSource) (erpdomain.ImportSnapshot, error) {
	r.gotSource = source
	snapshot, ok := r.bySource[source]
	if !ok {
		return erpdomain.ImportSnapshot{}, erpports.ErrImportNotFound
	}
	return snapshot, nil
}

var _ erpports.ImportRepository = (*sourceKeyedRepo)(nil)

func TestReaderSelectsSnapshotByActiveSource(t *testing.T) {
	xlsxRows := []erpdomain.NormalizedRow{{Codprod: "10", Descrprod: "Sankhya Widget", Custo: "12.30", StockPhysical: "9", StockReserved: ptr("1"), EAN: ptr("7894900011517"), NCM: ptr("12345678")}}
	clientRows := []erpdomain.NormalizedRow{{Codprod: "20", Descrprod: "Prospect Widget", EAN: ptr("7891000315019")}}
	repo := &sourceKeyedRepo{bySource: map[erpdomain.ImportSource]erpdomain.ImportSnapshot{
		erpdomain.SourceXLSX:            {Status: erpdomain.ImportStatusCompleted, ImportedAt: importedAt, AcceptedRows: xlsxRows},
		erpdomain.SourceCatalogoCliente: {Status: erpdomain.ImportStatusCompleted, ImportedAt: importedAt, AcceptedRows: clientRows},
	}}
	reader := NewReader(repo, "tenant-a", WithClock(func() time.Time { return fetchedAt }))

	// Absent selection => xlsx default (byte-stable with prior callers).
	defPage, err := reader.ListCatalogProductFacts(context.Background(), readports.Cursor{}, 10)
	if err != nil || repo.gotSource != erpdomain.SourceXLSX || len(defPage.Items) != 1 || defPage.Items[0].InternalProductID != 10 {
		t.Fatalf("default page = %+v gotSource=%q err=%v, want the xlsx snapshot", defPage.Items, repo.gotSource, err)
	}

	// catalogo_cliente selection => the prospect snapshot (honest-unknown cost/stock, ADR-17).
	ctx := WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente)
	clientPage, err := reader.ListCatalogProductFacts(ctx, readports.Cursor{}, 10)
	if err != nil || repo.gotSource != erpdomain.SourceCatalogoCliente || len(clientPage.Items) != 1 || clientPage.Items[0].InternalProductID != 20 {
		t.Fatalf("client page = %+v gotSource=%q err=%v, want the catalogo_cliente snapshot", clientPage.Items, repo.gotSource, err)
	}
	fact := clientPage.Items[0]
	if fact.Cost.Amount != nil || len(fact.Cost.Quality) == 0 || fact.Cost.Quality[0] != string(readdomain.QualityMissingCost) {
		t.Fatalf("prospect fact cost = %+v, want honest-unknown missing_cost (never fabricated)", fact.Cost)
	}
	if fact.SellableStock.Quantity != nil || len(fact.SellableStock.Quality) == 0 {
		t.Fatalf("prospect fact stock = %+v, want honest-unknown missing_stock", fact.SellableStock)
	}
}
