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
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type batchDriverState struct {
	mu        sync.Mutex
	queries   int
	failAt    int
	lastSQL   []string
	salesRows int
}

type batchDriver struct{ state *batchDriverState }

func (d batchDriver) Open(string) (driver.Conn, error) { return batchConn{state: d.state}, nil }

type batchConn struct{ state *batchDriverState }

func (c batchConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c batchConn) Close() error                        { return nil }
func (c batchConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (c batchConn) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.state.mu.Lock()
	defer c.state.mu.Unlock()
	c.state.queries++
	c.state.lastSQL = append(c.state.lastSQL, query)
	if c.state.failAt > 0 && c.state.queries == c.state.failAt {
		return nil, codedBatchError{}
	}
	switch {
	case strings.Contains(query, "FETCH FIRST 5001 ROWS ONLY"):
		return &batchRows{columns: []string{"sale_date", "product_id", "product_group_id", "quantity", "net_amount", "price_table_id", "operation_type_id", "line_type"}, count: c.state.salesRows}, nil
	case strings.Contains(query, "CUSSEMICM"):
		return &batchRows{columns: []string{"product_id", "amount", "effective_at"}}, nil
	default:
		return &batchRows{columns: []string{"product_id", "icms", "ipi", "pis", "cofins"}}, nil
	}
}

type batchRows struct {
	columns []string
	count   int
	index   int
}

func (r *batchRows) Columns() []string { return r.columns }
func (r *batchRows) Close() error      { return nil }
func (r *batchRows) Next(values []driver.Value) error {
	if r.index >= r.count {
		return io.EOF
	}
	r.index++
	if len(r.columns) == 8 {
		values[0] = time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
		values[1] = int64(r.index)
		values[2] = int64(12)
		values[3] = float64(1)
		values[4] = float64(10)
		values[5] = int64(1)
		values[6] = int64(1)
		values[7] = "sale"
	}
	return nil
}

type codedBatchError struct{}

func (codedBatchError) Error() string { return "driver text with DSN secret" }
func (codedBatchError) Code() int     { return 942 }

var batchDriverSequence int32

func openBatchDB(t *testing.T, state *batchDriverState) *sql.DB {
	t.Helper()
	driverName := "f02_batch_driver_" + strconv.Itoa(int(atomic.AddInt32(&batchDriverSequence, 1)))
	sql.Register(driverName, batchDriver{state: state})
	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestBatchReaderUses500ChunksAndNoPartialMapOnFailure(t *testing.T) {
	state := &batchDriverState{}
	reader := NewBatchReader(openBatchDB(t, state), oraclebatch.NewSemaphore(4))
	ids := make([]int64, 0, 1200)
	for index := int64(1); index <= 1200; index++ {
		ids = append(ids, index)
	}
	ids = append(ids, 1, 2)
	if _, err := reader.GetCostFactsByIDs(context.Background(), ids); err != nil {
		t.Fatalf("cost batch error = %v", err)
	}
	if _, err := reader.GetTaxFactsByIDs(context.Background(), ids); err != nil {
		t.Fatalf("tax batch error = %v", err)
	}
	if state.queries != 6 {
		t.Fatalf("query count = %d, want 6 (3 cost + 3 tax)", state.queries)
	}

	state.failAt = state.queries + 2
	result, err := reader.GetCostFactsByIDs(context.Background(), ids)
	if err == nil || !domain.IsReadErrorCode(err, domain.ReadErrorSourceUnavailable) {
		t.Fatalf("failure = %v, want typed source_unavailable", err)
	}
	if result != nil {
		t.Fatalf("partial result = %#v, want nil", result)
	}
}

func TestBatchReaderSalesHistoryPeeksAt5001AndMarksTruncated(t *testing.T) {
	state := &batchDriverState{salesRows: 5001}
	reader := NewBatchReader(openBatchDB(t, state), oraclebatch.NewSemaphore(4))
	productID := 42
	result, err := reader.GetSalesHistory(context.Background(), ports.SalesHistoryInput{ProductID: &productID, Window: domain.SalesHistoryWindow{Start: time.Now().Add(-time.Hour), End: time.Now()}})
	if err != nil {
		t.Fatalf("GetSalesHistory() error = %v", err)
	}
	if len(result.Entries) != 5000 || !result.Truncated {
		t.Fatalf("sales result entries=%d truncated=%v, want 5000/true", len(result.Entries), result.Truncated)
	}
	if len(state.lastSQL) != 1 || !strings.Contains(state.lastSQL[0], "FETCH FIRST 5001 ROWS ONLY") {
		t.Fatalf("sales SQL = %q, want 5001-row peek", state.lastSQL)
	}
}
