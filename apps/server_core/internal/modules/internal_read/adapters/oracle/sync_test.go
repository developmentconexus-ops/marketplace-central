package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

func TestSankhyaAdapterKindIsLiveReadThrough(t *testing.T) {
	a := NewSankhyaAdapter(&dispatchQueryer{}, &fakeMirror{}, "tenant_x")
	if got := a.Kind(); got != sourcekind.LiveReadThrough {
		t.Fatalf("Kind() = %q, want %q", got, sourcekind.LiveReadThrough)
	}
}

// TestSankhyaSyncMapsSnapshotHonestNull drives Sync over canned Oracle results and
// asserts the ratified Sankhya→E2.1 mapping: EAN from REFERENCIA only when
// EAN-shaped, referencia from REFFORN, custo/preco/estoque as-of merged by CODPROD,
// per-CODLOCAL stock summed into estoque_total, and honest-NULL (never 0) for every
// fact a product lacks — distinct from a real 0 balance.
func TestSankhyaSyncMapsSnapshotHonestNull(t *testing.T) {
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFPRO": {cols: 8, rows: [][]driver.Value{
			// CODPROD, DESCRPROD, NCM, REFERENCIA(EAN), REFFORN, CODGRUPOPROD, DESCRGRUPOPROD, DESCRICAO(marca)
			{int64(100), "Torneira", "84818090", "7894900011517", "DOCOL-99", int64(5), "Metais", "Docol"},
			{int64(200), "Parafuso", nil, "ABC123", nil, nil, nil, nil},
			{int64(300), "Sifao Orfao", nil, nil, nil, nil, nil, nil},
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
	if fv(p.EstoqueTotal) != 15.0 {
		t.Errorf("100.EstoqueTotal = %v, want 15 (10+5)", fv(p.EstoqueTotal))
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
	if p.EstoqueTotal == nil || fv(p.EstoqueTotal) != 0.0 {
		t.Errorf("200.EstoqueTotal = %v, want 0 (real zero balance, NOT NULL)", p.EstoqueTotal)
	}

	// Product 300 — orphan: present in base, absent everywhere else → honest NULL.
	p = rows["300"]
	if p.Custo != nil || p.PrecoVenda != nil || p.EstoqueTotal != nil {
		t.Errorf("300 should be honest-NULL, got custo=%v preco=%v estoque=%v", p.Custo, p.PrecoVenda, p.EstoqueTotal)
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
}

func (m *fakeMirror) ApplySnapshot(_ context.Context, tenantID string, rows []mirror.Row, locs []mirror.StockLocation) (int, error) {
	m.tenantID = tenantID
	m.rows = rows
	m.locs = locs
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

// --- query-dispatching fake driver -----------------------------------------

type fakeResult struct {
	cols int
	rows [][]driver.Value
}

type dispatchQueryer struct {
	results map[string]fakeResult
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

func (q *dispatchQueryer) QueryContext(ctx context.Context, query string, _ ...any) (*sql.Rows, error) {
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
