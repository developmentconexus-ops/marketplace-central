package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

func TestSankhyaAdapterKindIsLiveReadThrough(t *testing.T) {
	a := NewSankhyaAdapter(&dispatchQueryer{}, &fakeMirror{}, "tenant_x")
	if got := a.Kind(); got != sourcekind.LiveReadThrough {
		t.Fatalf("Kind() = %q, want %q", got, sourcekind.LiveReadThrough)
	}
}

func TestSankhyaStockSQLPinsSellableCompaniesAndLocations(t *testing.T) {
	if !strings.Contains(sankhyaStockSQL, "CODEMP IN (1, 2)") {
		t.Error("sankhyaStockSQL missing sellable company predicate CODEMP IN (1, 2)")
	}
	if !strings.Contains(sankhyaStockSQL, "CODLOCAL IN (10101, 10102)") {
		t.Error("sankhyaStockSQL missing sellable location predicate CODLOCAL IN (10101, 10102)")
	}
}

func TestSankhyaSyncMapsSellableAssortmentFields(t *testing.T) {
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFPRO": {cols: 12, rows: [][]driver.Value{
			{int64(100), "Produto R", nil, nil, nil, nil, nil, nil, "R", "S", nil, nil},
			{int64(200), "Produto blank", nil, nil, nil, nil, nil, nil, "   ", nil, nil, nil},
		}},
		"TGFCUS": {cols: 2},
		"TGFEXC": {cols: 2},
		"TGFEST": {cols: 3},
	}}
	mw := &fakeMirror{}

	if _, err := NewSankhyaAdapter(q, mw, "tenant_x").Sync(context.Background()); err != nil {
		t.Fatalf("Sync: %v", err)
	}
	rows := indexRows(mw.rows)
	if got := strv(rows["100"].Usoprod); got != "R" {
		t.Errorf("100.Usoprod = %q, want %q", got, "R")
	}
	if got := strv(rows["100"].ADEcommerce); got != "S" {
		t.Errorf("100.ADEcommerce = %q, want %q", got, "S")
	}
	if rows["200"].Usoprod != nil {
		t.Errorf("200.Usoprod = %q, want nil for Oracle whitespace", strv(rows["200"].Usoprod))
	}
	if rows["200"].ADEcommerce != nil {
		t.Errorf("200.ADEcommerce = %q, want nil", strv(rows["200"].ADEcommerce))
	}
}

func TestSankhyaSyncStoresKnownZeroForProductsOutsideTheSellableCut(t *testing.T) {
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFPRO": {cols: 12, rows: [][]driver.Value{
			{int64(100), "Showroom only", nil, nil, nil, nil, nil, nil, "R", "S", nil, nil},
			{int64(200), "Sellable", nil, nil, nil, nil, nil, nil, "R", "S", nil, nil},
		}},
		"TGFCUS": {cols: 2},
		"TGFEXC": {cols: 2},
		"TGFEST": {cols: 3, rows: [][]driver.Value{
			{int64(200), int64(10101), 7.0},
		}},
	}}
	mw := &fakeMirror{
		stored: map[string]mirror.Row{
			"300": {CodigoProduto: "300", EstoqueTotal: ptrFloat(9)},
		},
		absent: map[string]bool{},
	}
	if _, err := NewSankhyaAdapter(q, mw, "tenant_x").Sync(context.Background()); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	rows := indexRows(mw.rows)
	// Fatal, not Error: the checks below dereference EstoqueTotal, and a nil here is
	// exactly the A-9 regression. Continuing would replace the diagnostic with a panic.
	if rows["100"].EstoqueTotal == nil || *rows["100"].EstoqueTotal != 0 {
		t.Fatalf("100.EstoqueTotal = %v, want known zero 0", rows["100"].EstoqueTotal)
	}
	if *rows["100"].EstoqueTotal > 0 {
		t.Errorf("100 unexpectedly entered sellable set with estoque_total=%v", *rows["100"].EstoqueTotal)
	}
	if rows["200"].EstoqueTotal == nil || *rows["200"].EstoqueTotal != 7 {
		t.Fatalf("200.EstoqueTotal = %v, want 7", rows["200"].EstoqueTotal)
	}
	if *rows["200"].EstoqueTotal <= 0 {
		t.Errorf("200 unexpectedly absent from sellable set with estoque_total=%v", *rows["200"].EstoqueTotal)
	}
	// Product 300 is absent from Q1, so the keep-absent merge preserves its prior
	// facts and only flags it absent.
	prior := mw.stored["300"]
	if !mw.absent["300"] || prior.EstoqueTotal == nil || *prior.EstoqueTotal != 9 {
		t.Errorf("300 absent state/value = %v/%v, want true/9 preserved", mw.absent["300"], prior.EstoqueTotal)
	}
}

func ptrFloat(v float64) *float64 { return &v }

// TestSankhyaSyncMapsSnapshotHonestNull drives Sync over canned Oracle results and
// asserts the ratified Sankhya→E2.1 mapping: EAN from REFERENCIA only when
// EAN-shaped, referencia from REFFORN, custo/preco/estoque as-of merged by CODPROD,
// per-CODLOCAL stock summed into estoque_total, honest-NULL for unknown descriptive
// facts, and known-zero stock for a Q1 product absent from pinned Q4.
func TestSankhyaSyncMapsSnapshotHonestNull(t *testing.T) {
	fiscalDtRef := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFPRO": {cols: 12, rows: [][]driver.Value{
			// CODPROD, DESCRPROD, NCM, REFERENCIA(EAN), REFFORN, CODGRUPOPROD, DESCRGRUPOPROD, DESCRICAO(marca), USOPROD, AD_ECOMMERCE, GRUPOICMS, ORIGPROD
			// 100's EAN is space-padded (Oracle CHAR) — must still resolve after trimming.
			// 100 also carries a classified fiscal group (GRUPOICMS=7, distinct from CODGRUPOPROD=5)
			// and a classified ORIGPROD=1 (P2.b Task 3).
			{int64(100), "Torneira", "84818090", "  7894900011517 ", "DOCOL-99", int64(5), "Metais", "Docol", "R", "S", int64(7), int64(1)},
			// 200's GRUPOICMS/ORIGPROD are Oracle NULL — the ADR-17 negative control: an
			// unclassified fiscal fact must stay Go nil, never fall back to 0.
			{int64(200), "Parafuso", nil, "ABC123", nil, nil, nil, nil, nil, nil, nil, nil},
			// 300's description is whitespace-only → honest-NULL, not "".
			{int64(300), "   ", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		}},
		"TGFCUS": {cols: 2, rows: [][]driver.Value{
			{int64(100), 12.50},
			{int64(200), 3.00},
			// 300 absent → custo NULL
		}},
		"TGFEXC": {cols: 2, rows: [][]driver.Value{
			{int64(100), 169.90},
			// 200, 300 absent → preco NULL
		}},
		// S_ST (TIPIMPOSTO='S' only) vs R_TOTAL (SUM of 'S'+'I') are DIFFERENT
		// quantities — 100 proves they must not be conflated (2.50 vs 3.60).
		"TGFEFDVMRSTDIA": {cols: 4, rows: [][]driver.Value{
			{int64(100), 2.50, 3.60, fiscalDtRef},
			// 200 absent from the ST block entirely → all three stay NULL.
			{int64(300), nil, 1.10, fiscalDtRef}, // only 'I' present, no 'S' row → S_ST NULL
		}},
		"TGFEST": {cols: 3, rows: [][]driver.Value{
			// CODPROD, CODLOCAL, DISPONIVEL
			{int64(100), int64(1), 10.0},
			{int64(100), int64(2), 5.0},
			{int64(200), int64(1), 0.0}, // real zero balance, present row
			{int64(999), int64(1), 7.0}, // not in base → must be ignored
			// 300 absent → estoque NULL
		}},
	}}
	mw := &fakeMirror{}
	a := NewSankhyaAdapter(q, mw, "tenant_x")
	fixed := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	a.now = func() time.Time { return fixed }

	res, err := a.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if res.Processed != 3 || res.Errors != 0 {
		t.Fatalf("SyncResult = %+v, want {Processed:3 Errors:0}", res)
	}
	if mw.tenantID != "tenant_x" {
		t.Fatalf("tenantID = %q, want tenant_x", mw.tenantID)
	}

	// Bind order/count is load-bearing: a real Oracle rejects or mis-binds a
	// swapped company/dataRef, but the fake ignores binds at execution — so these
	// assertions are the only guard against a Q2/Q3 Go-argument regression.
	if a, ok := q.argsFor("TGFCUS"); !ok || len(a) != 2 || a[0] != sankhyaCompany || a[1] != fixed {
		t.Errorf("cost binds = %v, want [%d %v] (CODEMP then data_ref)", a, sankhyaCompany, fixed)
	}
	if a, ok := q.argsFor("TGFEXC"); !ok || len(a) != 1 || a[0] != fixed {
		t.Errorf("price binds = %v, want [%v] (data_ref only)", a, fixed)
	}
	// The Go-arg assertions above only lock argument ORDER; a swap of the SQL
	// placeholders themselves (e.g. CODEMP=:2 AND DTATUAL<=:1) would keep the same
	// arg slice and slip through. Assert the placeholder-to-column binding in the SQL
	// text so both halves — argument order AND placeholder position — are guarded.
	if s, ok := q.queryFor("TGFCUS"); !ok || !strings.Contains(s, "CODEMP = :1") || !strings.Contains(s, "DTATUAL <= :2") {
		t.Errorf("cost SQL placeholder binding wrong; want CODEMP=:1 and DTATUAL<=:2\n%s", s)
	}
	if s, ok := q.queryFor("TGFEXC"); !ok || !strings.Contains(s, "DTVIGOR <= :1") {
		t.Errorf("price SQL placeholder binding wrong; want DTVIGOR<=:1\n%s", s)
	}
	if a, ok := q.argsFor("TGFEFDVMRSTDIA"); !ok || len(a) != 2 || a[0] != sankhyaCompany || a[1] != fixed {
		t.Errorf("fiscal ST binds = %v, want [%d %v] (CODEMP then data_ref)", a, sankhyaCompany, fixed)
	}
	if s, ok := q.queryFor("TGFEFDVMRSTDIA"); !ok || !strings.Contains(s, "CODEMP = :1") || !strings.Contains(s, "DTMOV <= :2") {
		t.Errorf("fiscal ST SQL placeholder binding wrong; want CODEMP=:1 and DTMOV<=:2\n%s", s)
	}

	rows := indexRows(mw.rows)
	if len(rows) != 3 {
		t.Fatalf("mirror rows = %d, want 3", len(rows))
	}

	// Product 100 — fully resolved.
	p := rows["100"]
	if strv(p.EAN) != "7894900011517" {
		t.Errorf("100.EAN = %v, want 7894900011517 (from REFERENCIA)", strv(p.EAN))
	}
	if strv(p.Referencia) != "DOCOL-99" {
		t.Errorf("100.Referencia = %v, want DOCOL-99 (from REFFORN)", strv(p.Referencia))
	}
	if strv(p.Marca) != "Docol" || strv(p.GrupoCodigo) != "5" || strv(p.GrupoDescricao) != "Metais" || strv(p.NCM) != "84818090" {
		t.Errorf("100 identity mismatch: %+v", p)
	}
	if fv(p.Custo) != 12.50 || fv(p.PrecoVenda) != 169.90 {
		t.Errorf("100 custo/preco = %v/%v, want 12.5/169.9", fv(p.Custo), fv(p.PrecoVenda))
	}
	if p.GrupoICMS == nil || *p.GrupoICMS != 7 {
		t.Errorf("100.GrupoICMS = %v, want 7", intv(p.GrupoICMS))
	}
	if fv(p.EstoqueTotal) != 15.0 {
		t.Errorf("100.EstoqueTotal = %v, want 15 (10+5)", fv(p.EstoqueTotal))
	}
	if p.Origprod == nil || *p.Origprod != 1 {
		t.Errorf("100.Origprod = %v, want 1", intv(p.Origprod))
	}
	if p.StRetidoEntrada == nil || *p.StRetidoEntrada != 2.50 {
		t.Errorf("100.StRetidoEntrada = %v, want 2.50", p.StRetidoEntrada)
	}
	if p.RestituicaoUnit == nil || *p.RestituicaoUnit != 3.60 {
		t.Errorf("100.RestituicaoUnit = %v, want 3.60 (S+I, distinct from S alone)", p.RestituicaoUnit)
	}
	if p.FiscalDtRef == nil || !p.FiscalDtRef.Equal(fiscalDtRef) {
		t.Errorf("100.FiscalDtRef = %v, want %v", p.FiscalDtRef, fiscalDtRef)
	}

	// Product 200 — non-EAN REFERENCIA → NULL; no price → NULL; real 0 stock (not NULL).
	p = rows["200"]
	if p.EAN != nil {
		t.Errorf("200.EAN = %v, want NULL (ABC123 is not EAN-shaped)", strv(p.EAN))
	}
	if p.Referencia != nil {
		t.Errorf("200.Referencia = %v, want NULL", strv(p.Referencia))
	}
	if fv(p.Custo) != 3.00 {
		t.Errorf("200.Custo = %v, want 3", fv(p.Custo))
	}
	if p.PrecoVenda != nil {
		t.Errorf("200.PrecoVenda = %v, want NULL (no CODTAB=0 row)", fv(p.PrecoVenda))
	}
	if p.GrupoICMS != nil {
		t.Errorf("200.GrupoICMS = %v, want NULL (Oracle NULL must never scan as 0 — ADR-17)", intv(p.GrupoICMS))
	}
	if p.EstoqueTotal == nil || fv(p.EstoqueTotal) != 0.0 {
		t.Errorf("200.EstoqueTotal = %v, want 0 (real zero balance, NOT NULL)", p.EstoqueTotal)
	}
	if p.Origprod != nil {
		t.Errorf("200.Origprod = %v, want NULL", intv(p.Origprod))
	}
	if p.StRetidoEntrada != nil || p.RestituicaoUnit != nil || p.FiscalDtRef != nil {
		t.Errorf("200 fiscal ST = st=%v r=%v dt=%v, want all NULL (absent from the ST block entirely)", p.StRetidoEntrada, p.RestituicaoUnit, p.FiscalDtRef)
	}

	// Product 300 — present in the ERP snapshot but absent from pinned Q4 → known zero stock.
	p = rows["300"]
	if p.Custo != nil || p.PrecoVenda != nil || p.EstoqueTotal == nil || *p.EstoqueTotal != 0 {
		t.Errorf("300 should have NULL custo/preco and known-zero estoque, got custo=%v preco=%v estoque=%v", p.Custo, p.PrecoVenda, p.EstoqueTotal)
	}
	if p.GrupoICMS != nil {
		t.Errorf("300.GrupoICMS = %v, want NULL", intv(p.GrupoICMS))
	}
	if p.Descricao != nil {
		t.Errorf("300.Descricao = %q, want NULL (whitespace-only Oracle text is honest-unknown, not \"\")", strv(p.Descricao))
	}
	// 300 has an 'I' row but no 'S' row in the ST block: S_ST must stay NULL
	// while R_TOTAL (S+I, here just I) is a real value — proves the two
	// aggregates are computed independently, not one derived from the other.
	if p.StRetidoEntrada != nil {
		t.Errorf("300.StRetidoEntrada = %v, want NULL (no 'S' row for this product)", p.StRetidoEntrada)
	}
	if p.RestituicaoUnit == nil || *p.RestituicaoUnit != 1.10 {
		t.Errorf("300.RestituicaoUnit = %v, want 1.10 ('I' only)", p.RestituicaoUnit)
	}
	if p.FiscalDtRef == nil || !p.FiscalDtRef.Equal(fiscalDtRef) {
		t.Errorf("300.FiscalDtRef = %v, want %v", p.FiscalDtRef, fiscalDtRef)
	}

	// Stock locations: 100×2 + 200×1 = 3; codprod 999 (not in base) excluded.
	if len(mw.locs) != 3 {
		t.Fatalf("stock locations = %d, want 3 (999 must be excluded)", len(mw.locs))
	}
	for _, l := range mw.locs {
		if l.CodigoProduto == "999" {
			t.Errorf("codprod 999 leaked into stock locations")
		}
	}
}

// --- fake mirror writer -----------------------------------------------------

type fakeMirror struct {
	tenantID string
	rows     []mirror.Row
	locs     []mirror.StockLocation
	stored   map[string]mirror.Row
	absent   map[string]bool
}

func (m *fakeMirror) ApplySnapshot(_ context.Context, tenantID string, rows []mirror.Row, locs []mirror.StockLocation) (int, error) {
	m.tenantID = tenantID
	m.rows = rows
	m.locs = locs
	if m.stored != nil {
		for code := range m.stored {
			m.absent[code] = true
		}
		for _, row := range rows {
			m.stored[row.CodigoProduto] = row
			m.absent[row.CodigoProduto] = false
		}
	}
	return len(rows), nil
}

func indexRows(rows []mirror.Row) map[string]mirror.Row {
	out := make(map[string]mirror.Row, len(rows))
	for _, r := range rows {
		out[r.CodigoProduto] = r
	}
	return out
}

func strv(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

func fv(p *float64) float64 {
	if p == nil {
		return -1
	}
	return *p
}

func intv(p *int) string {
	if p == nil {
		return "<nil>"
	}
	return strconv.Itoa(*p)
}

// --- query-dispatching fake driver -----------------------------------------

type fakeResult struct {
	cols int
	rows [][]driver.Value
}

type capturedCall struct {
	query string
	args  []any
}

type dispatchQueryer struct {
	results map[string]fakeResult
	calls   []capturedCall
	dbs     []*sql.DB
}

func (q *dispatchQueryer) PingContext(context.Context) error { return nil }

func (q *dispatchQueryer) match(query string) fakeResult {
	for key, res := range q.results {
		if strings.Contains(query, key) {
			return res
		}
	}
	return fakeResult{cols: 1}
}

// argsFor returns the bound arguments of the first captured query containing substr.
func (q *dispatchQueryer) argsFor(substr string) ([]any, bool) {
	for _, c := range q.calls {
		if strings.Contains(c.query, substr) {
			return c.args, true
		}
	}
	return nil, false
}

// queryFor returns the SQL text of the first captured query containing substr, so a
// test can assert placeholder-to-column binding, not just Go argument order.
func (q *dispatchQueryer) queryFor(substr string) (string, bool) {
	for _, c := range q.calls {
		if strings.Contains(c.query, substr) {
			return c.query, true
		}
	}
	return "", false
}

func (q *dispatchQueryer) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	q.calls = append(q.calls, capturedCall{query: query, args: args})
	db := sql.OpenDB(dispatchConnector{res: q.match(query)})
	q.dbs = append(q.dbs, db)
	return db.QueryContext(ctx, query)
}

func (q *dispatchQueryer) QueryRowContext(ctx context.Context, query string, _ ...any) *sql.Row {
	db := sql.OpenDB(dispatchConnector{res: q.match(query)})
	q.dbs = append(q.dbs, db)
	return db.QueryRowContext(ctx, query)
}

type dispatchConnector struct{ res fakeResult }

func (c dispatchConnector) Connect(context.Context) (driver.Conn, error) {
	return &dispatchConn{res: c.res}, nil
}
func (dispatchConnector) Driver() driver.Driver { return dispatchDriver{} }

type dispatchDriver struct{}

func (dispatchDriver) Open(string) (driver.Conn, error) { return &dispatchConn{}, nil }

type dispatchConn struct{ res fakeResult }

func (c *dispatchConn) Prepare(string) (driver.Stmt, error) { return nil, io.EOF }
func (c *dispatchConn) Close() error                        { return nil }
func (c *dispatchConn) Begin() (driver.Tx, error)           { return nil, io.EOF }
func (c *dispatchConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &dispatchRows{res: c.res}, nil
}
func (c *dispatchConn) CheckNamedValue(*driver.NamedValue) error { return nil }

type dispatchRows struct {
	res   fakeResult
	index int
}

func (r *dispatchRows) Columns() []string {
	cols := make([]string, r.res.cols)
	for i := range cols {
		cols[i] = "c"
	}
	return cols
}
func (r *dispatchRows) Close() error { return nil }
func (r *dispatchRows) Next(dest []driver.Value) error {
	if r.index >= len(r.res.rows) {
		return io.EOF
	}
	copy(dest, r.res.rows[r.index])
	r.index++
	return nil
}
