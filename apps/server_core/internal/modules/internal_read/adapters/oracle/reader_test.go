package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestBuildSellableStockQueryUsesOracleFirstPolicy(t *testing.T) {
	query, args, err := buildSellableStockQuery(42664, domain.DefaultSellableStockPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(query, "METALPRD.TGFEST") {
		t.Fatalf("expected oracle stock query to stay inside TGFEST, got %q", query)
	}
	if len(args) != 5 {
		t.Fatalf("expected bound args for product + companies + locations + exclusions, got %d", len(args))
	}
}

func TestReaderReturnsTypedUnavailableWhenDBMissing(t *testing.T) {
	reader := NewReader(nil)

	_, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{
		ProductID: 42664,
		Policy:    domain.DefaultSellableStockPolicy(),
	})
	if !domain.IsReadErrorCode(err, domain.ReadErrorSourceUnavailable) {
		t.Fatalf("expected source_unavailable read error, got %v", err)
	}
}

func TestTaxInputsRequireExactSaleLineIdentityBeforeOracleAccess(t *testing.T) {
	reader := NewReader(nil)
	got, err := reader.GetTaxInputs(context.Background(), ports.TaxInput{
		ProductID: 42664,
		Policy:    domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatalf("GetTaxInputs() error = %v, want missing tax without Oracle access", err)
	}
	if got.ICMSAmount != nil || len(got.QualityFlags) != 1 || got.QualityFlags[0] != domain.QualityMissingTax {
		t.Fatalf("GetTaxInputs() = %+v, want nil amounts with missing_tax", got)
	}
}

func TestBuildTaxInputsQueryUsesExactSaleLinePredicates(t *testing.T) {
	policy := domain.DefaultTaxPolicy(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC))
	policy.Source = domain.TaxSourceIdentity{DocumentID: 98765, LineNumber: 3}
	query, args := buildTaxInputsQuery(42664, policy)

	for _, predicate := range []string{"d.NUNOTA = :1", "d.SEQUENCIA = :2", "i.CODPROD = :3", "d.CODINC = :4"} {
		if !strings.Contains(query, predicate) {
			t.Fatalf("tax query missing exact predicate %q: %s", predicate, query)
		}
	}
	if strings.Contains(query, "DTNEG") || len(args) != 4 || args[0] != 98765 || args[1] != 3 || args[2] != 42664 {
		t.Fatalf("tax query args/query = %#v / %s, want exact document-line-product identity", args, query)
	}
}

func TestFindProductsQueryMatchesGTINEANAgainstTGFPROReference(t *testing.T) {
	// REFERENCIA is the governed barcode source: the same query projects it as
	// the candidate's EAN and the mirror sync writes products_mirror.ean from it
	// under this exact GTIN rule.
	ean := "7894900011517"
	query, args, err := buildFindProductsQuery(ports.FindProductsInput{EAN: &ean})
	if err != nil {
		t.Fatalf("EAN-only query error = %v", err)
	}
	if !strings.Contains(query, "TRIM(p.REFERENCIA) = :1") || len(args) != 1 || args[0] != ean {
		t.Fatalf("EAN was not matched against the governed reference: query=%q args=%v", query, args)
	}
}

func TestFindProductsQueryTreatsANonGTINBarcodeAsNoMatch(t *testing.T) {
	// A malformed provider barcode is a no-match, never a widened query and
	// never a failure that aborts the caller's whole linking run.
	barcode := "not-a-gtin"
	query, args, err := buildFindProductsQuery(ports.FindProductsInput{EAN: &barcode})
	if !errors.Is(err, errNoMatchableInput) {
		t.Fatalf("non-GTIN EAN-only query error = %v, want errNoMatchableInput", err)
	}
	if strings.Contains(query, "TRIM(p.REFERENCIA) = :") || len(args) != 0 {
		t.Fatalf("non-GTIN barcode leaked into the query: query=%q args=%v", query, args)
	}

	title := "produto"
	query, _, err = buildFindProductsQuery(ports.FindProductsInput{EAN: &barcode, Title: &title})
	if err != nil || strings.Contains(query, "TRIM(p.REFERENCIA) = :") || !strings.Contains(query, "UPPER(p.DESCRPROD) LIKE") {
		t.Fatalf("non-GTIN barcode must not add a clause beside title: query=%q err=%v", query, err)
	}
}

func TestFindProductsReturnsNoCandidatesForANonGTINBarcode(t *testing.T) {
	barcode := "not-a-gtin"
	db := &identityQueryer{}
	candidates, err := NewReader(db).FindProductsForLinking(context.Background(), ports.FindProductsInput{EAN: &barcode})
	if err != nil {
		t.Fatalf("FindProductsForLinking() error = %v, want an honest empty result", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("candidates = %+v, want none", candidates)
	}
}

func TestProductIdentityProjectionUsesGTINReferenceAndNCM(t *testing.T) {
	rows := [][]driver.Value{
		{42664, "Valid product", "7894900011517", "REF-VALID", "82055900", nil, "Doka", 10, "S", nil, int64(1)},
		{42665, "Invalid product", "not-a-gtin", "REF-INVALID", nil, nil, nil, nil, "S", nil, int64(1)},
		{42666, "Blank product", "   ", nil, "", nil, "Menegotti", 11, "S", nil, int64(1)},
	}
	db := &identityQueryer{rows: rows}
	reader := NewReader(db)

	got, err := reader.FindProductsForLinking(context.Background(), ports.FindProductsInput{})
	if err != nil {
		t.Fatalf("FindProductsForLinking() error = %v", err)
	}
	if got[0].EAN == nil || *got[0].EAN != "7894900011517" || domain.HasQualityFlag(got[0].QualityFlags, domain.QualityInvalidEAN) {
		t.Fatalf("valid GTIN projection = %+v", got[0])
	}
	if got[0].ReferenceCode == nil || *got[0].ReferenceCode != "REF-VALID" || got[0].NCM == nil || *got[0].NCM != "82055900" {
		t.Fatalf("independent REFFORN/NCM projection = %+v", got[0])
	}
	for _, index := range []int{1, 2} {
		if got[index].EAN != nil || !domain.HasQualityFlag(got[index].QualityFlags, domain.QualityInvalidEAN) {
			t.Fatalf("invalid or blank REFERENCIA must be discarded with invalid_ean: %+v", got[index])
		}
	}
	if got[1].ReferenceCode == nil || *got[1].ReferenceCode != "REF-INVALID" || got[2].ReferenceCode != nil {
		t.Fatalf("REFFORN must remain independent of REFERENCIA: %+v %+v", got[1], got[2])
	}
	if got[1].BrandName != nil || got[1].NCM != nil || got[2].NCM != nil {
		t.Fatalf("absent brand/NCM must remain nil: %+v %+v", got[1], got[2])
	}
}

func TestProductEANCollisionPreservesDistinctActiveCODPRODs(t *testing.T) {
	rows := [][]driver.Value{
		{42664, "First", "7894900011517", "REF-1", "82055900", nil, "Doka", 10, "S", nil, int64(2)},
		{42665, "Second", "7894900011517", "REF-2", "82055900", nil, "Doka", 10, "S", nil, int64(2)},
	}
	db := &identityQueryer{rows: rows}
	reader := NewReader(db)

	got, err := reader.FindProductsForLinking(context.Background(), ports.FindProductsInput{})
	if err != nil {
		t.Fatalf("FindProductsForLinking() error = %v", err)
	}
	if len(got) != 2 || got[0].ProductID == got[1].ProductID {
		t.Fatalf("collision must preserve both CODPRODs: %+v", got)
	}
	for _, candidate := range got {
		if !domain.HasQualityFlag(candidate.QualityFlags, domain.QualityEANCollision) {
			t.Fatalf("candidate missing ean_collision: %+v", candidate)
		}
	}
	if len(db.queries) != 1 || !strings.Contains(db.queries[0], "COUNT(DISTINCT collision.CODPROD)") || !strings.Contains(db.queries[0], "collision.ATIVO = 'S'") {
		t.Fatalf("collision must use full active-product population query: %v", db.queries)
	}
}

type identityQueryer struct {
	rows    [][]driver.Value
	queries []string
	dbs     []*sql.DB
}

func (q *identityQueryer) PingContext(context.Context) error { return nil }

func (q *identityQueryer) QueryContext(ctx context.Context, query string, _ ...any) (*sql.Rows, error) {
	q.queries = append(q.queries, query)
	db := sql.OpenDB(identityConnector{rows: q.rows})
	q.dbs = append(q.dbs, db)
	return db.QueryContext(ctx, query)
}

func (q *identityQueryer) QueryRowContext(ctx context.Context, query string, _ ...any) *sql.Row {
	db := sql.OpenDB(identityConnector{rows: q.rows})
	q.dbs = append(q.dbs, db)
	return db.QueryRowContext(ctx, query)
}

type identityConnector struct{ rows [][]driver.Value }

func (c identityConnector) Connect(context.Context) (driver.Conn, error) {
	return &identityConn{rows: c.rows}, nil
}
func (identityConnector) Driver() driver.Driver { return identityDriver{} }

type identityDriver struct{}

func (identityDriver) Open(string) (driver.Conn, error) { return &identityConn{}, nil }

type identityConn struct{ rows [][]driver.Value }

func (c *identityConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}
func (c *identityConn) Close() error { return nil }
func (c *identityConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}
func (c *identityConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &identityRows{rows: c.rows}, nil
}
func (c *identityConn) CheckNamedValue(*driver.NamedValue) error { return nil }

type identityRows struct {
	rows  [][]driver.Value
	index int
}

func (r *identityRows) Columns() []string {
	return []string{"CODPROD", "DESCRPROD", "REFERENCIA", "REFFORN", "NCM", "CODGRUPOPROD", "MARCA", "CODMARCA", "ATIVO", "USOPROD", "EAN_ACTIVE_COUNT"}
}
func (r *identityRows) Close() error { return nil }
func (r *identityRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.index])
	r.index++
	return nil
}
