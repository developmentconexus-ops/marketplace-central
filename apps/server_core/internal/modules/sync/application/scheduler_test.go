package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

type successCall struct {
	installationID string
	entity         domain.Entity
	cursor         json.RawMessage
	at             time.Time
	incremental    bool
}

type failureCall struct {
	installationID string
	entity         domain.Entity
	syncErr        domain.SyncError
}

// fakeStore records calls and serves canned reads. It performs NO external I/O —
// proving the scheduler drives only the persistence seam (M01-C5 at the unit
// layer: no ML/Oracle client is reachable from the loop).
type fakeStore struct {
	mu        sync.Mutex
	readState map[domain.Entity]domain.SyncState
	successes []successCall
	failures  []failureCall
	readErr   error
}

func newFakeStore() *fakeStore {
	return &fakeStore{readState: map[domain.Entity]domain.SyncState{}}
}

func (f *fakeStore) Read(_ context.Context, _ string, entity domain.Entity) (domain.SyncState, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.readErr != nil {
		return domain.SyncState{}, false, f.readErr
	}
	st, ok := f.readState[entity]
	return st, ok, nil
}

func (f *fakeStore) RecordSuccess(_ context.Context, installationID string, entity domain.Entity, cursor json.RawMessage, at time.Time, incremental bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.successes = append(f.successes, successCall{installationID, entity, cursor, at, incremental})
	return nil
}

func (f *fakeStore) RecordFailure(_ context.Context, installationID string, entity domain.Entity, syncErr domain.SyncError) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failures = append(f.failures, failureCall{installationID, entity, syncErr})
	return nil
}

func fixedClock(ts time.Time) func() time.Time { return func() time.Time { return ts } }

func TestRunOnceRecordsSuccessWithCursor(t *testing.T) {
	store := newFakeStore()
	store.readState[domain.EntityProducts] = domain.SyncState{Cursor: json.RawMessage(`{"page":1}`)}
	at := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	s := NewScheduler(store, "erp", time.Minute, fixedClock(at))

	var gotCursor json.RawMessage
	if err := s.RegisterJob(domain.EntityProducts, func(_ context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		gotCursor = cursor
		return json.RawMessage(`{"page":2}`), nil
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	s.RunOnce(context.Background())

	if string(gotCursor) != `{"page":1}` {
		t.Errorf("job received cursor %q, want {\"page\":1}", gotCursor)
	}
	if len(store.successes) != 1 || len(store.failures) != 0 {
		t.Fatalf("successes=%d failures=%d, want 1/0", len(store.successes), len(store.failures))
	}
	got := store.successes[0]
	if got.installationID != "erp" || got.entity != domain.EntityProducts {
		t.Errorf("success key = %s/%s", got.installationID, got.entity)
	}
	if string(got.cursor) != `{"page":2}` {
		t.Errorf("persisted cursor = %q, want {\"page\":2}", got.cursor)
	}
	if got.incremental {
		t.Error("skeleton cycle must record a full sync, not incremental")
	}
	if !got.at.Equal(at) {
		t.Errorf("timestamp = %v, want %v", got.at, at)
	}
}

func TestRunOnceFailureIsIsolatedPerEntity(t *testing.T) {
	store := newFakeStore()
	at := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	s := NewScheduler(store, "erp", time.Minute, fixedClock(at))

	if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) {
		return nil, errors.New("boom")
	}); err != nil {
		t.Fatal(err)
	}
	// A second, healthy entity must still run after the first one fails.
	ran := false
	if err := s.RegisterJob(domain.EntityListings, func(context.Context, json.RawMessage) (json.RawMessage, error) {
		ran = true
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}

	s.RunOnce(context.Background())

	if !ran {
		t.Error("second entity did not run — failure was not isolated")
	}
	if len(store.failures) != 1 {
		t.Fatalf("failures = %d, want 1", len(store.failures))
	}
	fc := store.failures[0]
	if fc.entity != domain.EntityProducts || fc.syncErr.Message != "boom" {
		t.Errorf("failure = %+v", fc)
	}
	if !fc.syncErr.At.Equal(at) {
		t.Errorf("failure timestamp = %v, want %v", fc.syncErr.At, at)
	}
	if len(store.successes) != 1 || store.successes[0].entity != domain.EntityListings {
		t.Errorf("healthy entity success missing: %+v", store.successes)
	}
}

func TestRunOnceSkipsEntityOnReadError(t *testing.T) {
	store := newFakeStore()
	store.readErr = errors.New("transient read error")
	s := NewScheduler(store, "erp", time.Minute, fixedClock(time.Now()))

	ran := false
	if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) {
		ran = true
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}

	s.RunOnce(context.Background())

	// Fail honest: a read error is not a first run. The job must NOT run with a
	// fabricated cursor, and neither success nor failure is recorded.
	if ran {
		t.Error("job ran despite a cursor read error — must skip honestly, not fabricate a first run")
	}
	if len(store.successes) != 0 || len(store.failures) != 0 {
		t.Errorf("read error recorded state (successes=%d failures=%d), want 0/0", len(store.successes), len(store.failures))
	}
}

func TestRunOncePanicIsIsolatedAsFailure(t *testing.T) {
	store := newFakeStore()
	// Seed a real cursor for the panicking entity: a panic must record a failure
	// and NEVER reach RecordSuccess(nil), which would erase this stored cursor
	// (the #7×#3 interaction — panic isolation must not clobber progress).
	store.readState[domain.EntityProducts] = domain.SyncState{Cursor: json.RawMessage(`{"page":7}`)}
	at := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	s := NewScheduler(store, "erp", time.Minute, fixedClock(at))

	if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) {
		panic("kaboom")
	}); err != nil {
		t.Fatal(err)
	}
	ran := false
	if err := s.RegisterJob(domain.EntityListings, func(context.Context, json.RawMessage) (json.RawMessage, error) {
		ran = true
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}

	// Must not propagate the panic; must record it as the entity's failure and
	// keep serving the other entity.
	s.RunOnce(context.Background())

	if !ran {
		t.Error("panic in first job aborted the loop — not isolated")
	}
	if len(store.failures) != 1 || store.failures[0].entity != domain.EntityProducts {
		t.Fatalf("panic not recorded as products failure: %+v", store.failures)
	}
	if !strings.Contains(store.failures[0].syncErr.Message, "panicked") {
		t.Errorf("panic failure message = %q, want it to mention panicked", store.failures[0].syncErr.Message)
	}
	// The panicking entity must NOT have recorded a success: RecordSuccess(nil)
	// on the panic path would clobber the seeded {"page":7} cursor. Only the
	// healthy listings entity may appear in successes.
	for _, sc := range store.successes {
		if sc.entity == domain.EntityProducts {
			t.Errorf("panic path recorded a success for products (%q) — would clobber the stored cursor", sc.cursor)
		}
	}
}

func TestRegisterJobGuards(t *testing.T) {
	s := NewScheduler(newFakeStore(), "erp", time.Minute, nil)

	if err := s.RegisterJob(domain.Entity("bogus"), func(context.Context, json.RawMessage) (json.RawMessage, error) { return nil, nil }); !errors.Is(err, domain.ErrUnknownEntity) {
		t.Errorf("unknown entity error = %v, want ErrUnknownEntity", err)
	}
	if err := s.RegisterJob(domain.EntityProducts, nil); err == nil {
		t.Error("nil job func must be rejected")
	}
	if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) { return nil, nil }); err != nil {
		t.Fatalf("first registration: %v", err)
	}
	if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) { return nil, nil }); err == nil {
		t.Error("duplicate entity registration must be rejected")
	}
}

func TestStartStopsOnContextCancel(t *testing.T) {
	s := NewScheduler(newFakeStore(), "erp", time.Millisecond, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { s.Start(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Start did not return after context cancel")
	}
}

func TestStartNoOpOnZeroInterval(t *testing.T) {
	s := NewScheduler(newFakeStore(), "erp", 0, nil)
	// Must return immediately without a ticker panic.
	s.Start(context.Background())
}

// TestInferIncrementalTolerates pins inferIncremental's tolerant-parse
// contract: an absent, empty, unrecognized, or unparseable phase all resolve
// to false with no error — never a crash, never a spurious incremental flag.
// The "legacy products shape" case is the load-bearing one: the real
// products job's ProductsCursor has no phase key at all and must keep
// recording incremental=false with zero behavior change (F-03).
func TestInferIncrementalTolerates(t *testing.T) {
	cases := []struct {
		name   string
		cursor json.RawMessage
		want   bool
	}{
		{"nil cursor (state erased)", nil, false},
		{"empty cursor", json.RawMessage(``), false},
		{"malformed json", json.RawMessage(`{not json`), false},
		{"json array, not an object", json.RawMessage(`[1,2,3]`), false},
		{"phase incremental", json.RawMessage(`{"phase":"incremental"}`), true},
		{"phase sweep", json.RawMessage(`{"phase":"sweep"}`), true},
		{"phase backfill", json.RawMessage(`{"phase":"backfill"}`), false},
		{"phase empty string", json.RawMessage(`{"phase":""}`), false},
		{"phase unrecognized", json.RawMessage(`{"phase":"bogus"}`), false},
		{
			"legacy products shape (no phase key)",
			json.RawMessage(`{"source":"xlsx","processed":1845,"completed_at":"2026-07-24T12:00:00Z"}`),
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := inferIncremental(tc.cursor); got != tc.want {
				t.Errorf("inferIncremental(%s) = %v, want %v", tc.cursor, got, tc.want)
			}
		})
	}
}

// TestRunOnceDerivesIncrementalFromCursorPhase drives the derivation through
// the real RunOnce/runJob path (not just the helper in isolation): whatever
// phase the job's returned cursor carries must be what RecordSuccess is told.
func TestRunOnceDerivesIncrementalFromCursorPhase(t *testing.T) {
	cases := []struct {
		name        string
		nextCursor  json.RawMessage
		incremental bool
	}{
		{"phase incremental", json.RawMessage(`{"phase":"incremental"}`), true},
		{"phase sweep", json.RawMessage(`{"phase":"sweep"}`), true},
		{"phase backfill", json.RawMessage(`{"phase":"backfill"}`), false},
		{"legacy products shape (no phase key)", json.RawMessage(`{"source":"xlsx","processed":3}`), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newFakeStore()
			store.readState[domain.EntityProducts] = domain.SyncState{Cursor: json.RawMessage(`{"phase":"backfill"}`)}
			at := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
			s := NewScheduler(store, "erp", time.Minute, fixedClock(at))

			if err := s.RegisterJob(domain.EntityProducts, func(context.Context, json.RawMessage) (json.RawMessage, error) {
				return tc.nextCursor, nil
			}); err != nil {
				t.Fatal(err)
			}

			s.RunOnce(context.Background())

			if len(store.successes) != 1 {
				t.Fatalf("successes = %d, want 1", len(store.successes))
			}
			if got := store.successes[0].incremental; got != tc.incremental {
				t.Errorf("incremental = %v, want %v (cursor %s)", got, tc.incremental, tc.nextCursor)
			}
		})
	}
}
