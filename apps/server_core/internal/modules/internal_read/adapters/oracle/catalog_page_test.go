package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestCatalogPageUsesOneQueryForEveryListSize(t *testing.T) {
	for _, limit := range []int{1, 50, 100} {
		t.Run("limit-"+itoa(limit), func(t *testing.T) {
			rows := make([][]driver.Value, limit+1)
			for i := range rows {
				rows[i] = catalogDriverRow(int64(i + 1))
			}
			queryer := &fakeCatalogQueryer{rows: rows}
			reader := NewReader(queryer)

			page, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{}, limit)
			if err != nil {
				t.Fatalf("ListCatalogProductFacts() error = %v", err)
			}
			if queryer.queryCount != 1 {
				t.Fatalf("query count = %d, want 1", queryer.queryCount)
			}
			if len(page.Items) != limit || page.NextCursor == nil {
				t.Fatalf("page = items:%d next:%v, want %d items and next cursor", len(page.Items), page.NextCursor, limit)
			}
			if !strings.Contains(queryer.queries[0], "METALPRD.TGFEXC") || strings.Contains(queryer.queries[0], "VW_PRECO_TABELA") {
				t.Fatalf("catalog query did not use base price tables: %s", queryer.queries[0])
			}
		})
	}
}

func TestCatalogPageCursorChainIsGaplessAndNonOverlapping(t *testing.T) {
	ctx := context.Background()
	cases := [][][]driver.Value{
		{catalogDriverRow(1), catalogDriverRow(2), catalogDriverRow(3), catalogDriverRow(4)},
		{catalogDriverRow(4), catalogDriverRow(5), catalogDriverRow(6), catalogDriverRow(7)},
		{catalogDriverRow(7), catalogDriverRow(8)},
	}
	var cursor ports.Cursor
	var got []int64
	for i, rows := range cases {
		queryer := &fakeCatalogQueryer{rows: rows}
		reader := NewReader(queryer)
		page, err := reader.ListCatalogProductFacts(ctx, cursor, 3)
		if err != nil {
			t.Fatalf("page %d error = %v", i+1, err)
		}
		for _, item := range page.Items {
			got = append(got, item.InternalProductID)
		}
		if i < 2 {
			if page.NextCursor == nil || page.NextCursor.InternalProductID != int64((i+1)*3) {
				t.Fatalf("page %d next cursor = %+v", i+1, page.NextCursor)
			}
			cursor = *page.NextCursor
		} else if page.NextCursor != nil {
			t.Fatalf("last page next cursor = %+v, want nil", page.NextCursor)
		}
		if queryer.queryCount != 1 {
			t.Fatalf("page %d query count = %d, want 1", i+1, queryer.queryCount)
		}
		if i > 0 && !strings.Contains(queryer.queries[0], "p.CODPROD > :1") {
			t.Fatalf("page %d query missing keyset predicate: %s", i+1, queryer.queries[0])
		}
	}
	want := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	if !equalInt64s(got, want) {
		t.Fatalf("cursor chain IDs = %v, want %v", got, want)
	}
}

func TestCatalogPageMapsNullableFactsAndAmbiguousPrice(t *testing.T) {
	row := catalogDriverRow(42664)
	row[5] = nil
	row[6] = nil
	row[7] = int64(2)
	row[8] = nil
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{row}}
	reader := NewReader(queryer)

	page, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{}, 1)
	if err != nil {
		t.Fatalf("ListCatalogProductFacts() error = %v", err)
	}
	item := page.Items[0]
	if item.SellableStock.Quantity != nil || !hasCatalogQuality(item.SellableStock.Quality, domain.QualityMissingStock) {
		t.Fatalf("stock = %+v, want nil + missing_stock", item.SellableStock)
	}
	if item.CurrentPrice.Amount != nil || !hasCatalogQuality(item.CurrentPrice.Quality, domain.QualityAmbiguousPrice) {
		t.Fatalf("price = %+v, want nil + ambiguous_price", item.CurrentPrice)
	}
	if item.Cost.Amount != nil || !hasCatalogQuality(item.Cost.Quality, domain.QualityMissingCost) {
		t.Fatalf("cost = %+v, want nil + missing_cost", item.Cost)
	}
}

func hasCatalogQuality(flags []string, target domain.QualityFlag) bool {
	for _, flag := range flags {
		if flag == string(target) {
			return true
		}
	}
	return false
}

func TestCatalogPageInvalidCursorDoesNotTouchOracle(t *testing.T) {
	queryer := &fakeCatalogQueryer{}
	reader := NewReader(queryer)
	_, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{InternalProductID: -1}, 50)
	if !errors.Is(err, ports.ErrInvalidCursor) {
		t.Fatalf("error = %v, want ErrInvalidCursor", err)
	}
	var typedErr *ports.InvalidCursorError
	if !errors.As(err, &typedErr) {
		t.Fatalf("error = %T, want *InvalidCursorError", err)
	}
	if queryer.queryCount != 0 {
		t.Fatalf("query count = %d, want 0", queryer.queryCount)
	}
}

func TestCatalogSearchUsesOneBoundedQueryAndNoCursor(t *testing.T) {
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{catalogDriverRow(4)}}
	reader := NewReader(queryer)
	page, err := reader.SearchCatalogProductFacts(context.Background(), "bolt", 50)
	if err != nil {
		t.Fatalf("SearchCatalogProductFacts() error = %v", err)
	}
	if queryer.queryCount != 1 || page.NextCursor != nil || len(page.Items) != 1 {
		t.Fatalf("query/page = %d/%+v, want one query and no cursor", queryer.queryCount, page)
	}
	if !strings.Contains(queryer.queries[0], "UPPER(p.DESCRPROD) LIKE") {
		t.Fatalf("search predicate missing: %s", queryer.queries[0])
	}
}

func catalogDriverRow(id int64) []driver.Value {
	return []driver.Value{id, "REF-" + itoa(int(id)), "Description", nil, "S", float64(id), "12.90", int64(1), "4.50"}
}

func itoa(value int) string {
	const digits = "0123456789"
	if value == 0 {
		return "0"
	}
	var out [20]byte
	i := len(out)
	for value > 0 {
		i--
		out[i] = digits[value%10]
		value /= 10
	}
	return string(out[i:])
}

func equalInt64s(left, right []int64) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

type fakeCatalogQueryer struct {
	queryCount int
	queries    []string
	rows       [][]driver.Value
	dbs        []*sql.DB
}

func (f *fakeCatalogQueryer) PingContext(context.Context) error { return nil }

func (f *fakeCatalogQueryer) QueryContext(ctx context.Context, query string, _ ...any) (*sql.Rows, error) {
	f.queryCount++
	f.queries = append(f.queries, query)
	db := sql.OpenDB(fakeCatalogConnector{rows: f.rows})
	f.dbs = append(f.dbs, db)
	return db.QueryContext(ctx, query)
}

func (f *fakeCatalogQueryer) QueryRowContext(ctx context.Context, query string, _ ...any) *sql.Row {
	f.queryCount++
	f.queries = append(f.queries, query)
	db := sql.OpenDB(fakeCatalogConnector{rows: f.rows})
	f.dbs = append(f.dbs, db)
	return db.QueryRowContext(ctx, query)
}

type fakeCatalogConnector struct{ rows [][]driver.Value }

func (c fakeCatalogConnector) Connect(context.Context) (driver.Conn, error) {
	return &fakeCatalogConn{rows: c.rows}, nil
}
func (fakeCatalogConnector) Driver() driver.Driver { return fakeCatalogDriver{} }

type fakeCatalogDriver struct{}

func (fakeCatalogDriver) Open(string) (driver.Conn, error) { return &fakeCatalogConn{}, nil }

type fakeCatalogConn struct {
	rows [][]driver.Value
}

func (c *fakeCatalogConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}
func (c *fakeCatalogConn) Close() error { return nil }
func (c *fakeCatalogConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}
func (c *fakeCatalogConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &fakeCatalogRows{rows: c.rows}, nil
}
func (c *fakeCatalogConn) CheckNamedValue(*driver.NamedValue) error { return nil }

type fakeCatalogRows struct {
	rows  [][]driver.Value
	index int
}

func (r *fakeCatalogRows) Columns() []string {
	return []string{"CODPROD", "REFERENCIA", "DESCRPROD", "EAN", "ATIVO", "SELLABLE_QTY", "PRICE_AMOUNT", "ACTIVE_PRICE_ROWS", "COST_AMOUNT"}
}
func (r *fakeCatalogRows) Close() error { return nil }
func (r *fakeCatalogRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.index])
	r.index++
	return nil
}
