package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type icmsCeilingDriverState struct {
	mu       sync.Mutex
	query    string
	args     []driver.NamedValue
	rows     [][]driver.Value
	queryErr error
}

type icmsCeilingDriver struct{ state *icmsCeilingDriverState }

func (d icmsCeilingDriver) Open(string) (driver.Conn, error) {
	return icmsCeilingConn{state: d.state}, nil
}

type icmsCeilingConn struct{ state *icmsCeilingDriverState }

func (c icmsCeilingConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c icmsCeilingConn) Close() error                        { return nil }
func (c icmsCeilingConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (c icmsCeilingConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.state.mu.Lock()
	defer c.state.mu.Unlock()
	c.state.query = query
	c.state.args = append([]driver.NamedValue(nil), args...)
	if c.state.queryErr != nil {
		return nil, c.state.queryErr
	}
	return &icmsCeilingRows{rows: c.state.rows}, nil
}

type icmsCeilingRows struct {
	rows  [][]driver.Value
	index int
}

func (r *icmsCeilingRows) Columns() []string { return []string{"UFDEST", "CEILING"} }
func (r *icmsCeilingRows) Close() error      { return nil }
func (r *icmsCeilingRows) Next(values []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(values, r.rows[r.index])
	r.index++
	return nil
}

var icmsCeilingDriverSequence int32

func openICMSCeilingDB(t *testing.T, state *icmsCeilingDriverState) *sql.DB {
	t.Helper()
	driverName := "f02_icms_ceiling_driver_" + strconv.Itoa(int(atomic.AddInt32(&icmsCeilingDriverSequence, 1)))
	sql.Register(driverName, icmsCeilingDriver{state: state})
	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestBuildICMSCeilingQueryUsesRatifiedCeilingColumnsAndBind(t *testing.T) {
	query, args := buildICMSCeilingQuery(domain.UF(13))
	upperQuery := strings.ToUpper(query)
	for _, want := range []string{"ALIQUFDEST", "PERCICMSFCP", "METALPRD.TGFICM", "GROUP BY UFDEST", ":UFORIG"} {
		if !strings.Contains(upperQuery, want) {
			t.Fatalf("query = %q, want %q", query, want)
		}
	}
	for _, forbidden := range []string{"ALIQUOTA", "REDBASE", "TGFITE.ALIQICMS"} {
		if strings.Contains(upperQuery, forbidden) {
			t.Fatalf("query = %q, must not contain %q", query, forbidden)
		}
	}
	if strings.Contains(query, "13") {
		t.Fatalf("query = %q, origin must not be string-concatenated", query)
	}
	if len(args) != 1 {
		t.Fatalf("args = %#v, want one named bind", args)
	}
	bind, ok := args[0].(sql.NamedArg)
	if !ok || bind.Name != "uforig" || bind.Value != int64(13) {
		t.Fatalf("args = %#v, want sql.Named(\"uforig\", int64(13))", args)
	}
}

func TestGetICMSCeilingByOriginMapsValuesAndPreservesUnknowns(t *testing.T) {
	state := &icmsCeilingDriverState{rows: [][]driver.Value{
		{int64(8), float64(22)},
		{int64(10), nil},
	}}
	reader := NewBatchReader(openICMSCeilingDB(t, state), oraclebatch.NewSemaphore(1))

	got, err := reader.GetICMSCeilingByOrigin(context.Background(), domain.UF(13))
	if err != nil {
		t.Fatalf("GetICMSCeilingByOrigin() error = %v", err)
	}
	if got[domain.UF(8)] == nil || got[domain.UF(8)].CeilingPercent == nil || *got[domain.UF(8)].CeilingPercent != 22 {
		t.Fatalf("value row = %#v, want non-nil 22%% ceiling", got[domain.UF(8)])
	}
	if got[domain.UF(10)] != nil {
		t.Fatalf("NULL row = %#v, want nil map value", got[domain.UF(10)])
	}
	if _, ok := got[domain.UF(11)]; ok {
		t.Fatalf("absent UF map = %#v, want key absent", got)
	}

	state.mu.Lock()
	query, args := state.query, state.args
	state.mu.Unlock()
	if !strings.Contains(strings.ToUpper(query), "WHERE UFORIG=:UFORIG") {
		t.Fatalf("executed query = %q, want named origin predicate", query)
	}
	if len(args) != 1 || args[0].Name != "uforig" || args[0].Value != int64(13) {
		t.Fatalf("executed args = %#v, want named origin bind 13", args)
	}
}
