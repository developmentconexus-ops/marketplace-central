package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"
)

// fakeStaleReader is the local fake for staleAggregateReader. It records the
// (olderThan, limit) it was called with so the ceiling test can assert on it,
// and returns either a canned ID list or a canned error.
type fakeStaleReader struct {
	ids        []string
	err        error
	gotOlder   time.Duration
	gotLimit   int
	callsCount int
}

func (f *fakeStaleReader) StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error) {
	f.callsCount++
	f.gotOlder = olderThan
	f.gotLimit = limit
	if f.err != nil {
		return nil, f.err
	}
	return f.ids, nil
}

// fakeProductCollector is the local fake for productCollector. failFor names
// codprods that must return an error from Collect; every other codprod
// succeeds. Named distinctly from the pre-existing fakeCollector in
// collection_pipeline_service_test.go (which fakes ports.PriceIntelCollector,
// a different interface entirely) to avoid a same-package redeclaration.
type fakeProductCollector struct {
	failFor map[string]error
	calls   []string
}

func (f *fakeProductCollector) Collect(ctx context.Context, codprod string) (CollectionSummary, error) {
	f.calls = append(f.calls, codprod)
	if err, ok := f.failFor[codprod]; ok {
		return CollectionSummary{}, err
	}
	return CollectionSummary{Codprod: codprod}, nil
}

func testLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, nil))
}

// TestCollectionJobRespectsThePerCycleCeiling proves the job hands its
// configured ceiling straight to the reader's limit argument. This is what
// keeps one cycle from sweeping the whole ~2,923-product catalog and starving
// the shared 900/min ML rate-limit bucket that product/listing/order sync also
// draw from.
func TestCollectionJobRespectsThePerCycleCeiling(t *testing.T) {
	reader := &fakeStaleReader{ids: []string{"P1", "P2"}}
	collector := &fakeProductCollector{}
	ttl := 6 * time.Hour
	ceiling := 50

	job := NewCollectionJob(reader, collector, ttl, ceiling, testLogger(&bytes.Buffer{}))

	if _, err := job(context.Background(), nil); err != nil {
		t.Fatalf("job() error = %v, want nil", err)
	}
	if reader.gotOlder != ttl {
		t.Fatalf("StaleAggregateProductIDs olderThan = %v, want %v", reader.gotOlder, ttl)
	}
	if reader.gotLimit != ceiling {
		t.Fatalf("StaleAggregateProductIDs limit = %d, want %d (the per-cycle ceiling)", reader.gotLimit, ceiling)
	}
}

// TestCollectionJobContinuesAfterOneProductFails guards the same defect class
// integrations/background/refresh_ticker.go:37-51 used to have: RunOnce
// returning on the first item error and silently never trying the rest of the
// batch. One bad codprod must not stop the cycle from reaching the others.
func TestCollectionJobContinuesAfterOneProductFails(t *testing.T) {
	reader := &fakeStaleReader{ids: []string{"P1", "P2", "P3"}}
	collector := &fakeProductCollector{failFor: map[string]error{"P2": errors.New("ml api 500")}}

	job := NewCollectionJob(reader, collector, time.Hour, 100, testLogger(&bytes.Buffer{}))

	if _, err := job(context.Background(), nil); err != nil {
		t.Fatalf("job() error = %v, want nil (a per-product failure is not a cycle failure)", err)
	}
	if len(collector.calls) != 3 {
		t.Fatalf("collector calls = %v, want all 3 products attempted despite P2 failing", collector.calls)
	}
	if collector.calls[0] != "P1" || collector.calls[1] != "P2" || collector.calls[2] != "P3" {
		t.Fatalf("collector calls = %v, want [P1 P2 P3] in order", collector.calls)
	}
}

// TestCollectionJobLogsEveryFailure is the negative control for the exact
// F-A1 defect: `case <-ticker.C: _ = t.RunOnce(ctx)` discarded the error with
// no log at all, so a collection failure was invisible at every layer. A
// per-product failure here must land in the structured log.
func TestCollectionJobLogsEveryFailure(t *testing.T) {
	reader := &fakeStaleReader{ids: []string{"P1", "P2"}}
	boom := errors.New("ml api 500")
	collector := &fakeProductCollector{failFor: map[string]error{"P2": boom}}
	var buf bytes.Buffer

	job := NewCollectionJob(reader, collector, time.Hour, 100, testLogger(&buf))

	// Whether the batch itself returns an error is asserted by
	// TestCollectionJobContinuesAfterOneProductFails; this test's only concern
	// is that the failure left a trace in the log, so it does not gate on the
	// return value here.
	_, _ = job(context.Background(), nil)
	if !bytes.Contains(buf.Bytes(), []byte("P2")) {
		t.Fatalf("log does not mention the failing codprod P2: %s", buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("ml api 500")) {
		t.Fatalf("log does not mention the underlying error: %s", buf.String())
	}
}

// TestCollectionJobReturnsTheEnumerationErrorSoSyncStateRecordsIt proves an
// enumeration failure is returned, not swallowed to nil. Returning the error
// is what makes syncapp.Scheduler.runJob call RecordFailure (scheduler.go:152)
// so the failure surfaces on SyncHealthCard. Swallowing it here would leave
// collection broken with a green card -- exactly the F-A1 defect this feature
// fixed for token refresh, now checked for this new job too.
func TestCollectionJobReturnsTheEnumerationErrorSoSyncStateRecordsIt(t *testing.T) {
	boom := errors.New("db down")
	reader := &fakeStaleReader{err: boom}
	collector := &fakeProductCollector{}

	job := NewCollectionJob(reader, collector, time.Hour, 100, testLogger(&bytes.Buffer{}))

	cursor := json.RawMessage(`{"phase":"incremental"}`)
	got, err := job(context.Background(), cursor)
	if !errors.Is(err, boom) {
		t.Fatalf("job() error = %v, want the enumeration error surfaced (not swallowed)", err)
	}
	if len(collector.calls) != 0 {
		t.Fatalf("collector was called despite enumeration failing: %v", collector.calls)
	}
	if string(got) != string(cursor) {
		t.Fatalf("cursor = %s, want the input cursor preserved unchanged on enumeration failure", got)
	}
}
