package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type instrumentedStockDB struct {
	query func(string, []driver.NamedValue) ([][]driver.Value, error)
}

type instrumentedStockConnector struct{ db *instrumentedStockDB }

func (c instrumentedStockConnector) Connect(context.Context) (driver.Conn, error) {
	return instrumentedStockConn{db: c.db}, nil
}
func (c instrumentedStockConnector) Driver() driver.Driver { return instrumentedStockDriver{db: c.db} }

type instrumentedStockDriver struct{ db *instrumentedStockDB }

func (d instrumentedStockDriver) Open(string) (driver.Conn, error) {
	return instrumentedStockConn{db: d.db}, nil
}

type instrumentedStockConn struct{ db *instrumentedStockDB }

func (instrumentedStockConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not used")
}
func (instrumentedStockConn) Close() error { return nil }
func (instrumentedStockConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not used")
}
func (instrumentedStockConn) Ping(context.Context) error { return nil }
func (c instrumentedStockConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	values, err := c.db.query(query, args)
	if err != nil {
		return nil, err
	}
	return &instrumentedStockRows{values: values}, nil
}

type instrumentedStockRows struct {
	values [][]driver.Value
	index  int
}

func (r *instrumentedStockRows) Columns() []string { return []string{"CODPROD", "SELLABLE_QTY"} }
func (r *instrumentedStockRows) Close() error      { return nil }
func (r *instrumentedStockRows) Next(dest []driver.Value) error {
	if r.index >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}

func openInstrumentedStockDB(query func(string, []driver.NamedValue) ([][]driver.Value, error)) *sql.DB {
	db := sql.OpenDB(instrumentedStockConnector{db: &instrumentedStockDB{query: query}})
	db.SetMaxOpenConns(8)
	return db
}

func stockIDsInQuery(query string, args []driver.NamedValue) []int64 {
	marker := "est.CODPROD IN ("
	start := strings.Index(query, marker)
	if start < 0 {
		return nil
	}
	start += len(marker)
	end := strings.Index(query[start:], ")")
	if end < 0 {
		return nil
	}
	count := strings.Count(query[start:start+end], ":")
	ids := make([]int64, 0, count)
	for _, arg := range args[:count] {
		if id, ok := arg.Value.(int64); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func TestStockBatchQueriesExpectedChunkCount(t *testing.T) {
	for _, test := range []struct {
		name        string
		count       int
		wantQueries int32
	}{
		{name: "one", count: 1, wantQueries: 1},
		{name: "exact chunk", count: 500, wantQueries: 1},
		{name: "one over chunk", count: 501, wantQueries: 2},
		{name: "three chunks", count: 1200, wantQueries: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			var queries atomic.Int32
			db := openInstrumentedStockDB(func(query string, args []driver.NamedValue) ([][]driver.Value, error) {
				queries.Add(1)
				rows := make([][]driver.Value, 0)
				for _, id := range stockIDsInQuery(query, args) {
					rows = append(rows, []driver.Value{id, float64(id)})
				}
				return rows, nil
			})
			defer db.Close()
			ids := make([]int64, test.count)
			for i := range ids {
				ids[i] = int64(i + 1)
			}
			facts, err := NewBatchStockReader(db, oraclebatch.NewSemaphore(4)).GetStockFactsByIDs(context.Background(), ids)
			if err != nil {
				t.Fatalf("GetStockFactsByIDs() error = %v", err)
			}
			if queries.Load() != test.wantQueries || len(facts) != test.count {
				t.Fatalf("queries=%d facts=%d, want queries=%d facts=%d", queries.Load(), len(facts), test.wantQueries, test.count)
			}
		})
	}
}

func TestStockBatchDeduplicatesAndEmptyInputDoesNotQuery(t *testing.T) {
	var queries atomic.Int32
	db := openInstrumentedStockDB(func(string, []driver.NamedValue) ([][]driver.Value, error) {
		queries.Add(1)
		return [][]driver.Value{{int64(7), float64(3)}}, nil
	})
	defer db.Close()
	reader := NewBatchStockReader(db, oraclebatch.NewSemaphore(4))
	got, err := reader.GetStockFactsByIDs(context.Background(), []int64{7, 7, 7})
	if err != nil || len(got) != 1 || queries.Load() != 1 {
		t.Fatalf("deduplicated read: facts=%v err=%v queries=%d", got, err, queries.Load())
	}
	got, err = reader.GetStockFactsByIDs(context.Background(), nil)
	if err != nil || len(got) != 0 || queries.Load() != 1 {
		t.Fatalf("empty read: facts=%v err=%v queries=%d", got, err, queries.Load())
	}
}

func TestStockBatchChunkFailureReturnsNoPartialMap(t *testing.T) {
	var queries atomic.Int32
	db := openInstrumentedStockDB(func(string, []driver.NamedValue) ([][]driver.Value, error) {
		if queries.Add(1) == 2 {
			return nil, errors.New("driver leaked host=oracle.example")
		}
		return [][]driver.Value{{int64(1), float64(3)}}, nil
	})
	defer db.Close()
	ids := make([]int64, 501)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	facts, err := NewBatchStockReader(db, oraclebatch.NewSemaphore(4)).GetStockFactsByIDs(context.Background(), ids)
	if err == nil || facts != nil || !domain.IsReadErrorCode(err, domain.ReadErrorSourceUnavailable) {
		t.Fatalf("failure result: facts=%v err=%v, want typed error and nil map", facts, err)
	}
	if strings.Contains(err.Error(), "oracle.example") {
		t.Fatalf("raw driver cause leaked: %v", err)
	}
}

func TestStockBatchAdapterSemaphoreBoundsInstrumentedFake(t *testing.T) {
	var inFlight, maxInFlight atomic.Int32
	db := openInstrumentedStockDB(func(string, []driver.NamedValue) ([][]driver.Value, error) {
		current := inFlight.Add(1)
		for {
			previous := maxInFlight.Load()
			if current <= previous || maxInFlight.CompareAndSwap(previous, current) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		inFlight.Add(-1)
		return [][]driver.Value{{int64(1), float64(1)}}, nil
	})
	defer db.Close()
	reader := NewBatchStockReader(db, oraclebatch.NewSemaphore(4))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := reader.GetStockFactsByIDs(context.Background(), []int64{1}); err != nil {
				t.Errorf("GetStockFactsByIDs() error = %v", err)
			}
		}()
	}
	wg.Wait()
	if maxInFlight.Load() != 4 {
		t.Fatalf("max in-flight = %d, want 4", maxInFlight.Load())
	}
}
