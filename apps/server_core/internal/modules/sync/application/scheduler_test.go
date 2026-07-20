package application

import (
	"context"
	"encoding/json"
	"errors"
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
