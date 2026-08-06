package application

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// fakeEnumerator is a scripted ports.IDEnumerator: pages[i] is returned for
// the i-th call, unless errs[i] is set (then the page is not consulted).
type fakeEnumerator struct {
	pages   []ports.BatchIDPage
	errs    map[int]error
	cursors []string
}

func (e *fakeEnumerator) ReadIDPage(_ context.Context, _ ports.InstallationAccount, cursor string, _ int) (ports.BatchIDPage, error) {
	i := len(e.cursors)
	e.cursors = append(e.cursors, cursor)
	if err := e.errs[i]; err != nil {
		return ports.BatchIDPage{}, err
	}
	if i >= len(e.pages) {
		return ports.BatchIDPage{}, nil
	}
	return e.pages[i], nil
}

// fakeHydrator is a scripted ports.BatchHydrator: it hydrates whatever ids
// have a matching entry in byID and silently drops the rest (mirrors a
// real hydrator's per-item skip, without needing MultigetItemError plumbing
// at this layer — that contract is proven in adapters/connectors instead).
type fakeHydrator struct {
	byID map[string]listingsdomain.Listing
}

func (h *fakeHydrator) HydrateBatch(_ context.Context, _ ports.InstallationAccount, ids []string) ([]listingsdomain.Listing, error) {
	rows := make([]listingsdomain.Listing, 0, len(ids))
	for _, id := range ids {
		if row, ok := h.byID[id]; ok {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

type runStoreCall struct {
	rows   []listingsdomain.Listing
	seenAt time.Time
}

// fakeRunStore records every UpsertPulledRows/MarkRunComplete call verbatim
// (including the exact seenAt/runStartedAt argument each received) so tests
// can assert the run-scoped ADR-027 pin survives across ticks and resumes,
// without needing a real Postgres-backed CompletedPullStore.
type fakeRunStore struct {
	upserts   []runStoreCall
	completes []time.Time
}

func (s *fakeRunStore) UpsertPulledRows(_ context.Context, _ string, rows []listingsdomain.Listing, seenAt time.Time) error {
	s.upserts = append(s.upserts, runStoreCall{rows: append([]listingsdomain.Listing(nil), rows...), seenAt: seenAt})
	return nil
}

func (s *fakeRunStore) MarkRunComplete(_ context.Context, _ string, runStartedAt time.Time) error {
	s.completes = append(s.completes, runStartedAt)
	return nil
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

// TestListingIngestorUsesInjectedClockForSeenAt proves the resync-writer
// fix's clock-injectability requirement: NewListingIngestor's now parameter
// is a real constructor-level dependency, not time.Now hardcoded internally
// — a fixed clock's exact value must reach UpsertPulledRows's seenAt
// argument, the same way root.go's installationResyncWriter passes time.Now
// through to it per Apply call.
func TestListingIngestorUsesInjectedClockForSeenAt(t *testing.T) {
	fixed := time.Date(2026, 8, 1, 12, 30, 0, 0, time.UTC)
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{"item-1": testListing(t, "item-1", "-")}}
	store := &fakeRunStore{}

	ingestor := NewListingIngestor(hydrator, store, testAccount(), fixedClock(fixed))
	if err := ingestor.IngestListing(context.Background(), "item-1"); err != nil {
		t.Fatalf("IngestListing() error = %v", err)
	}

	if len(store.upserts) != 1 {
		t.Fatalf("upserts = %+v, want exactly 1", store.upserts)
	}
	if !store.upserts[0].seenAt.Equal(fixed) {
		t.Fatalf("upserts[0].seenAt = %v, want the injected clock value %v", store.upserts[0].seenAt, fixed)
	}
	if len(store.completes) != 0 {
		t.Fatalf("completes = %+v, want none: a single-item ingest must never call MarkRunComplete", store.completes)
	}
}

func TestListingsJobIngestsPageAndAdvancesCursor(t *testing.T) {
	t0 := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	enumerator := &fakeEnumerator{pages: []ports.BatchIDPage{{IDs: []string{"item-1"}, NextCursor: "page-2"}}}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{"item-1": testListing(t, "item-1", "-")}}
	store := &fakeRunStore{}

	job := NewListingsJob(enumerator, hydrator, store, testAccount(), fixedClock(t0))
	next, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("job() error = %v", err)
	}
	if len(store.upserts) != 1 || len(store.upserts[0].rows) != 1 {
		t.Fatalf("upserts = %+v", store.upserts)
	}
	if len(store.completes) != 0 {
		t.Fatalf("completes = %+v, want none (run not exhausted yet)", store.completes)
	}
	var cur backfillCursor
	if err := json.Unmarshal(next, &cur); err != nil {
		t.Fatalf("cursor unmarshal: %v", err)
	}
	if cur.Phase != "backfill" || cur.ScrollID != "page-2" || cur.IDsCollected != 1 || !cur.RunStartedAt.Equal(t0) {
		t.Fatalf("cursor = %+v", cur)
	}
}

// TestListingsJobResumeFromPersistedCursorKeepsSameRunIdentity is the
// job-level kill-and-resume proof (required test 5): a process restart that
// rebuilds the job from the LAST PERSISTED cursor must keep completing the
// SAME run (same run_started_at pin) rather than starting a fresh one at the
// restart's wall-clock time — otherwise the keep-absent bound at
// MarkRunComplete would only cover the resumed tail, silently un-marking
// rows that were genuinely absent for the whole run.
func TestListingsJobResumeFromPersistedCursorKeepsSameRunIdentity(t *testing.T) {
	t0 := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	t1 := t0.Add(10 * time.Minute) // restart happens later

	enumerator1 := &fakeEnumerator{pages: []ports.BatchIDPage{{IDs: []string{"item-1"}, NextCursor: "page-2"}}}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{
		"item-1": testListing(t, "item-1", "-"),
		"item-2": testListing(t, "item-2", "-"),
	}}
	store := &fakeRunStore{}

	job1 := NewListingsJob(enumerator1, hydrator, store, testAccount(), fixedClock(t0))
	persistedCursor, err := job1(context.Background(), nil)
	if err != nil {
		t.Fatalf("first tick error = %v", err)
	}

	// "Kill": job1's closure is discarded here; a fresh job is built (as the
	// scheduler would on restart) with only the persisted cursor bytes and a
	// LATER wall clock.
	enumerator2 := &fakeEnumerator{pages: []ports.BatchIDPage{{IDs: []string{"item-2"}, NextCursor: ""}}}
	job2 := NewListingsJob(enumerator2, hydrator, store, testAccount(), fixedClock(t1))
	next, err := job2(context.Background(), persistedCursor)
	if err != nil {
		t.Fatalf("resumed tick error = %v", err)
	}

	if len(store.upserts) != 2 {
		t.Fatalf("upserts = %+v, want 2 ticks", store.upserts)
	}
	for i, u := range store.upserts {
		if !u.seenAt.Equal(t0) {
			t.Fatalf("upsert[%d].seenAt = %v, want the ORIGINAL run start %v (not the resume time %v)", i, u.seenAt, t0, t1)
		}
	}
	if len(store.completes) != 1 || !store.completes[0].Equal(t0) {
		t.Fatalf("completes = %+v, want exactly [%v] (the run's own start, not the resume time)", store.completes, t0)
	}

	var sweep sweepCursor
	if err := json.Unmarshal(next, &sweep); err != nil || sweep.Phase != "sweep" {
		t.Fatalf("terminal cursor = %s, err = %v", next, err)
	}
}

// TestListingsJobZeroRowTerminalTickStillMarksRunComplete is the job-level
// analog of required test 3: a run's final tick can legitimately enumerate
// zero new ids (the catalog count divides evenly into page-sized chunks) —
// that must still trigger MarkRunComplete for the rows upserted by EARLIER
// ticks of the same run, not be mistaken for "nothing happened".
func TestListingsJobZeroRowTerminalTickStillMarksRunComplete(t *testing.T) {
	t0 := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	enumerator := &fakeEnumerator{pages: []ports.BatchIDPage{
		{IDs: []string{"item-1"}, NextCursor: "page-2"},
		{IDs: nil, NextCursor: ""},
	}}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{"item-1": testListing(t, "item-1", "-")}}
	store := &fakeRunStore{}

	job := NewListingsJob(enumerator, hydrator, store, testAccount(), fixedClock(t0))
	cursor1, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("tick 1 error = %v", err)
	}
	final, err := job(context.Background(), cursor1)
	if err != nil {
		t.Fatalf("tick 2 (zero-row terminal) error = %v", err)
	}

	if len(store.upserts) != 1 {
		t.Fatalf("upserts = %+v, want exactly 1 (the zero-row tick must not call UpsertPulledRows)", store.upserts)
	}
	if len(store.completes) != 1 || !store.completes[0].Equal(t0) {
		t.Fatalf("completes = %+v, want exactly [%v]", store.completes, t0)
	}
	var sweep sweepCursor
	if err := json.Unmarshal(final, &sweep); err != nil || sweep.Phase != "sweep" {
		t.Fatalf("terminal cursor = %s, err = %v", final, err)
	}
}

// TestListingsJobEnumeratorErrorPreservesCursorUnchanged proves the
// error-path contract NewListingsJob documents: an enumerator failure
// returns the SAME cursor bytes the job was invoked with (byte-identical),
// matching JobFunc's own convention that the scheduler discards whatever a
// failing tick returns and keeps the stored cursor untouched either way.
func TestListingsJobEnumeratorErrorPreservesCursorUnchanged(t *testing.T) {
	raw, err := json.Marshal(backfillCursor{Phase: "backfill", ScrollID: "resume-here", RunStartedAt: time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("marshal seed cursor: %v", err)
	}
	enumerator := &fakeEnumerator{errs: map[int]error{0: context.DeadlineExceeded}}
	store := &fakeRunStore{}
	job := NewListingsJob(enumerator, &fakeHydrator{}, store, testAccount(), time.Now)

	next, jobErr := job(context.Background(), raw)
	if jobErr == nil {
		t.Fatal("job() error = nil, want the enumerator failure")
	}
	if !bytes.Equal(next, raw) {
		t.Fatalf("returned cursor = %s, want byte-identical to input %s", next, raw)
	}
	if len(store.upserts) != 0 || len(store.completes) != 0 {
		t.Fatalf("store must not be touched on enumerator error: upserts=%+v completes=%+v", store.upserts, store.completes)
	}
}

// TestListingsJobTerminalSweepCursorRestartsFreshBackfillRun proves the
// sweep -> backfill cursor round-trip: a cadence-driven re-invocation that
// hands the job its OWN prior terminal cursor must start a brand new run
// (fresh scroll_id "", fresh run_started_at at the CURRENT clock) rather
// than being rejected or looping forever on the terminal state.
func TestListingsJobTerminalSweepCursorRestartsFreshBackfillRun(t *testing.T) {
	prior, err := json.Marshal(sweepCursor{Phase: "sweep", LastFullSweepAt: time.Date(2026, 7, 30, 3, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("marshal seed cursor: %v", err)
	}
	t1 := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	enumerator := &fakeEnumerator{pages: []ports.BatchIDPage{{IDs: []string{"item-1"}, NextCursor: "page-2"}}}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{"item-1": testListing(t, "item-1", "-")}}
	store := &fakeRunStore{}

	job := NewListingsJob(enumerator, hydrator, store, testAccount(), fixedClock(t1))
	next, err := job(context.Background(), prior)
	if err != nil {
		t.Fatalf("job() error = %v", err)
	}
	if len(enumerator.cursors) != 1 || enumerator.cursors[0] != "" {
		t.Fatalf("enumerator called with cursor = %#v, want a fresh empty cursor", enumerator.cursors)
	}
	var cur backfillCursor
	if err := json.Unmarshal(next, &cur); err != nil {
		t.Fatalf("cursor unmarshal: %v", err)
	}
	if cur.Phase != "backfill" || cur.ScrollID != "page-2" || !cur.RunStartedAt.Equal(t1) {
		t.Fatalf("cursor = %+v, want a fresh run started at %v", cur, t1)
	}
}
