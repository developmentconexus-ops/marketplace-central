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

func TestCatalogPageSellableAssortmentCutAndScreenEscape(t *testing.T) {
	cursor := ports.Cursor{InternalProductID: 41}
	allOn := ports.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}
	filtered, filteredArgs := buildCatalogPageQueryWithPolicy(cursor, 51, " bolt ", allOn)
	all, allArgs := buildCatalogPageQueryWithPolicy(cursor, 51, " bolt ", ports.AllProductsAssortment())

	ruleClauses := []string{
		"UPPER(TRIM(p.USOPROD)) = 'R'",
		"NVL(stock.sellable_qty, 0) > 0",
		"NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'",
	}
	for _, clause := range ruleClauses {
		if !strings.Contains(filtered, clause) {
			t.Fatalf("filtered catalog query missing %q:\n%s", clause, filtered)
		}
		if strings.Contains(all, clause) {
			t.Fatalf("include_all catalog query retained %q:\n%s", clause, all)
		}
	}
	for _, clause := range []string{
		"est.CODEMP IN (1, 2)",
		"est.CODLOCAL IN (10101, 10102)",
		"p.ATIVO = 'S'",
		"p.CODPROD > 0",
		"p.CODPROD > :1",
		"UPPER(p.DESCRPROD) LIKE :2",
		"FETCH FIRST :3 ROWS ONLY",
	} {
		if !strings.Contains(filtered, clause) || !strings.Contains(all, clause) {
			t.Fatalf("catalog query mode missing shared clause %q:\nFILTERED:\n%s\nINCLUDE_ALL:\n%s", clause, filtered, all)
		}
	}
	wantArgs := []any{int64(41), "%BOLT%", 51}
	if !equalCatalogArgs(filteredArgs, wantArgs) || !equalCatalogArgs(allArgs, wantArgs) {
		t.Fatalf("catalog binds = filtered:%#v include_all:%#v, want %#v", filteredArgs, allArgs, wantArgs)
	}
	t.Logf("PAGE SQL:\n%s\nPAGE BINDS: %#v", filtered, filteredArgs)
	t.Logf("INCLUDE_ALL SQL:\n%s\nINCLUDE_ALL BINDS: %#v", all, allArgs)
}

func TestCatalogAssortmentCountUsesThePagePredicate(t *testing.T) {
	allOn := ports.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}
	page, _ := buildCatalogPageQueryWithPolicy(ports.Cursor{}, 51, "", allOn)
	count, args := buildCatalogAssortmentCountQuery(allOn)

	for _, clause := range []string{
		"FROM METALPRD.TGFEST est",
		"est.CODEMP IN (1, 2)",
		"est.CODLOCAL IN (10101, 10102)",
		"UPPER(TRIM(p.USOPROD)) = 'R'",
		"NVL(stock.sellable_qty, 0) > 0",
		"NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'",
		"p.ATIVO = 'S'",
		"p.CODPROD > 0",
	} {
		if !strings.Contains(page, clause) || !strings.Contains(count, clause) {
			t.Fatalf("page/count predicate drift at %q:\nPAGE:\n%s\nCOUNT:\n%s", clause, page, count)
		}
	}
	if !strings.Contains(count, "sellable_count") || !strings.Contains(count, "total_count") {
		t.Fatalf("count query does not project both exact counts:\n%s", count)
	}
	if len(args) != 0 {
		t.Fatalf("count binds = %#v, want []", args)
	}
	t.Logf("COUNT SQL:\n%s\nCOUNT BINDS: %#v", count, args)
}

func TestCatalogPageAndCountUseTheExactPolicyValue(t *testing.T) {
	policy := ports.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true}
	page, _ := buildCatalogPageQueryWithPolicy(ports.Cursor{}, 51, "", policy)
	count, _ := buildCatalogAssortmentCountQuery(policy)
	stockClause := "NVL(stock.sellable_qty, 0) > 0"
	for name, query := range map[string]string{"page": page, "count": count} {
		if strings.Contains(query, stockClause) {
			t.Fatalf("%s query for policy read %+v retained disabled clause %q:\n%s", name, policy, stockClause, query)
		}
		for _, enabled := range []string{"UPPER(TRIM(p.USOPROD)) = 'R'", "NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'"} {
			if !strings.Contains(query, enabled) {
				t.Fatalf("%s query for policy read %+v lost enabled clause %q:\n%s", name, policy, enabled, query)
			}
		}
	}
}

func TestCatalogProductFactsByIDsRemainsUnfiltered(t *testing.T) {
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{catalogDriverRow(7)}}
	reader := NewReader(queryer)
	page, err := reader.CatalogProductFactsByIDs(context.Background(), []int64{7})
	if err != nil {
		t.Fatalf("CatalogProductFactsByIDs() error = %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].InternalProductID != 7 {
		t.Fatalf("explicit-id page = %+v, want product 7", page.Items)
	}
	query := queryer.queries[0]
	for _, clause := range []string{
		"UPPER(TRIM(p.USOPROD)) = 'R'",
		"NVL(stock.sellable_qty, 0) > 0",
		"NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'",
	} {
		if strings.Contains(query, clause) {
			t.Fatalf("explicit-id query retained assortment clause %q:\n%s", clause, query)
		}
	}
	if !strings.Contains(query, "p.CODPROD IN (:2)") {
		t.Fatalf("explicit-id query missing bound ID predicate:\n%s", query)
	}
}

func TestCatalogPageMapsNullableFactsAndAmbiguousPrice(t *testing.T) {
	row := catalogDriverRow(42664)
	row[8] = nil
	row[9] = nil
	row[10] = int64(2)
	row[11] = nil
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{row}}
	reader := NewReader(queryer)

	page, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{}, 1)
	if err != nil {
		t.Fatalf("ListCatalogProductFacts() error = %v", err)
	}
	item := page.Items[0]
	if item.SellableStock.Quantity == nil || *item.SellableStock.Quantity != 0 || len(item.SellableStock.Quality) != 0 {
		t.Fatalf("stock = %+v, want quantity 0 with no quality flag", item.SellableStock)
	}
	if item.CurrentPrice.Amount != nil || !hasCatalogQuality(item.CurrentPrice.Quality, domain.QualityAmbiguousPrice) {
		t.Fatalf("price = %+v, want nil + ambiguous_price", item.CurrentPrice)
	}
	if item.Cost.Amount != nil || !hasCatalogQuality(item.Cost.Quality, domain.QualityMissingCost) {
		t.Fatalf("cost = %+v, want nil + missing_cost", item.Cost)
	}
}

func TestCatalogPageIdentitySemantics(t *testing.T) {
	validEAN := "4006381333931"
	tests := []struct {
		name      string
		row       []driver.Value
		wantEAN   *string
		wantRef   *string
		wantNCM   *string
		wantFlags []domain.QualityFlag
	}{
		{name: "valid EAN and reference alias", row: catalogIdentityDriverRow(101, validEAN, "MF-101", "Marca", "12345678", 1), wantEAN: &validEAN, wantRef: stringPointer("MF-101"), wantNCM: stringPointer("12345678"), wantFlags: []domain.QualityFlag{domain.QualityComplete}},
		{name: "invalid_ean is nil not empty", row: catalogIdentityDriverRow(102, "not-a-gtin", "MF-102", "Marca", "1234", 1), wantRef: stringPointer("MF-102"), wantFlags: []domain.QualityFlag{domain.QualityInvalidEAN}},
		{name: "blank EAN is nil not empty", row: catalogIdentityDriverRow(103, "   ", nil, nil, nil, 1), wantFlags: []domain.QualityFlag{domain.QualityInvalidEAN}},
		{name: "ean_collision preserves product", row: catalogIdentityDriverRow(104, validEAN, "MF-104", "Marca", "12345678", 2), wantEAN: &validEAN, wantRef: stringPointer("MF-104"), wantNCM: stringPointer("12345678"), wantFlags: []domain.QualityFlag{domain.QualityComplete, domain.QualityEANCollision}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(&fakeCatalogQueryer{rows: [][]driver.Value{tt.row}})
			page, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{}, 1)
			if err != nil {
				t.Fatal(err)
			}
			item := page.Items[0]
			if !equalStringPointers(item.EAN, tt.wantEAN) || !equalStringPointers(item.Reference, tt.wantRef) || !equalStringPointers(item.ManufacturerReference, tt.wantRef) || !equalStringPointers(item.NCM, tt.wantNCM) {
				t.Fatalf("identity = ean:%v reference:%v manufacturer_reference:%v ncm:%v", item.EAN, item.Reference, item.ManufacturerReference, item.NCM)
			}
			if len(item.QualityFlags) != len(tt.wantFlags) {
				t.Fatalf("quality_flags = %v, want %v", item.QualityFlags, tt.wantFlags)
			}
			for i := range tt.wantFlags {
				if item.QualityFlags[i] != string(tt.wantFlags[i]) {
					t.Fatalf("quality_flags = %v, want %v", item.QualityFlags, tt.wantFlags)
				}
			}
		})
	}

	t.Run("collision preserves both CODPRODs", func(t *testing.T) {
		rows := [][]driver.Value{catalogIdentityDriverRow(201, validEAN, "A", "Marca", "12345678", 2), catalogIdentityDriverRow(202, validEAN, "B", "Marca", "12345678", 2)}
		reader := NewReader(&fakeCatalogQueryer{rows: rows})
		page, err := reader.ListCatalogProductFacts(context.Background(), ports.Cursor{}, 2)
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 2 || page.Items[0].InternalProductID != 201 || page.Items[1].InternalProductID != 202 {
			t.Fatalf("items = %+v", page.Items)
		}
	})

	query, _ := buildCatalogPageQueryWithPolicy(ports.Cursor{}, 2, "", ports.AllProductsAssortment())
	if strings.Contains(query, "p.REFERENCIA,\n    p.DESCRPROD") || !strings.Contains(query, "p.REFFORN") || !strings.Contains(query, "EAN_ACTIVE_COUNT") {
		t.Fatalf("query does not isolate REFERENCIA as EAN and REFFORN as reference: %s", query)
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

func TestCatalogSearchUsesOneBoundedQuery(t *testing.T) {
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{catalogDriverRow(4)}}
	reader := NewReader(queryer)
	page, err := reader.SearchCatalogProductFacts(context.Background(), "bolt", ports.Cursor{}, 50)
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

// Search used to fetch exactly `limit` rows and never set NextCursor, so a query
// matching more products than the page held returned a full page with
// next_cursor: null — byte-identical to "that is all there is". The operator saw
// 50 results and no way to ask for the 51st.
func TestCatalogSearchReportsMoreMatchesAndWalksThem(t *testing.T) {
	queryer := &fakeCatalogQueryer{rows: [][]driver.Value{catalogDriverRow(4), catalogDriverRow(9)}}
	reader := NewReader(queryer)

	first, err := reader.SearchCatalogProductFacts(context.Background(), "bolt", ports.Cursor{}, 1)
	if err != nil {
		t.Fatalf("SearchCatalogProductFacts() error = %v", err)
	}
	if len(first.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(first.Items))
	}
	if first.NextCursor == nil || first.NextCursor.InternalProductID != 4 {
		t.Fatalf("next_cursor = %+v, want cursor at 4", first.NextCursor)
	}

	if _, err := reader.SearchCatalogProductFacts(context.Background(), "bolt", *first.NextCursor, 1); err != nil {
		t.Fatalf("second page error = %v", err)
	}
	second := queryer.queries[1]
	if !strings.Contains(second, "p.CODPROD > :1") || !strings.Contains(second, "UPPER(p.DESCRPROD) LIKE :2") {
		t.Fatalf("second page must carry both the cursor and the search predicate: %s", second)
	}
}

func TestCatalogSearchRejectsNegativeCursorWithoutTouchingOracle(t *testing.T) {
	queryer := &fakeCatalogQueryer{}
	reader := NewReader(queryer)
	_, err := reader.SearchCatalogProductFacts(context.Background(), "bolt", ports.Cursor{InternalProductID: -1}, 50)
	if !errors.Is(err, ports.ErrInvalidCursor) {
		t.Fatalf("error = %v, want ErrInvalidCursor", err)
	}
	if queryer.queryCount != 0 {
		t.Fatalf("query count = %d, want 0", queryer.queryCount)
	}
}

func catalogDriverRow(id int64) []driver.Value {
	return catalogIdentityDriverRow(id, nil, "REF-"+itoa(int(id)), nil, nil, nil)
}

func catalogIdentityDriverRow(id int64, ean, reference, brand, ncm any, collisionCount any) []driver.Value {
	return []driver.Value{id, ean, reference, "Description", brand, ncm, collisionCount, "S", float64(id), "12.90", int64(1), "4.50"}
}

func stringPointer(value string) *string { return &value }
func equalStringPointers(left, right *string) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
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

func equalCatalogArgs(left, right []any) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
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
	return []string{"CODPROD", "EAN", "REFFORN", "DESCRPROD", "MARCA", "NCM", "EAN_ACTIVE_COUNT", "ATIVO", "SELLABLE_QTY", "PRICE_AMOUNT", "ACTIVE_PRICE_ROWS", "COST_AMOUNT"}
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
