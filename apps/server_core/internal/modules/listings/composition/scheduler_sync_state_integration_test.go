//go:build integration

package composition_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	syncpg "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// zeroPageEnumerator is a minimal ports.IDEnumerator that reports exhaustion
// (NextCursor == "") on its very first page, with no ids — the "legitimate
// empty terminal tick" backfill.go:130-134 documents. This makes the FIRST
// Scheduler.RunOnce call land the terminal sweepCursor in one tick, which is
// exactly what this test needs to observe (backfill.go:141-171,
// specifically the NextCursor=="" branch at :160-165).
type zeroPageEnumerator struct{}

func (zeroPageEnumerator) ReadIDPage(context.Context, listingsports.InstallationAccount, string, int) (listingsports.BatchIDPage, error) {
	return listingsports.BatchIDPage{NextCursor: ""}, nil
}

// neverCalledHydrator asserts ingestBatch's zero-ids short-circuit
// (backfill.go:28-31: len(ids)==0 returns before ever calling the hydrator)
// actually holds — a real hydrator implementation is unnecessary scope for a
// test that only cares about the sync_state round trip.
type neverCalledHydrator struct{ t *testing.T }

func (h neverCalledHydrator) HydrateBatch(context.Context, listingsports.InstallationAccount, []string) ([]listingsdomain.Listing, error) {
	h.t.Fatal("HydrateBatch must not be called: the enumerator's page has zero ids")
	return nil, errors.New("unreachable")
}

// noopCompletedPullStore is a minimal ports.CompletedPullStore. UpsertPulledRows
// is never expected to run (paired with neverCalledHydrator/zeroPageEnumerator);
// MarkRunComplete just needs to succeed so NewListingsJob's terminal branch can
// proceed to build the sweep cursor. Real listings-row persistence (a real
// Postgres CompletedPullStore, FK/schema setup) is out of scope for this test —
// it is already covered by repository_integration_test.go in this milestone.
type noopCompletedPullStore struct{ markRunCompleteCalls int }

func (s *noopCompletedPullStore) UpsertPulledRows(context.Context, string, []listingsdomain.Listing, time.Time) error {
	return errors.New("UpsertPulledRows must not be called: the enumerator page has zero ids")
}

func (s *noopCompletedPullStore) MarkRunComplete(context.Context, string, time.Time) error {
	s.markRunCompleteCalls++
	return nil
}

// sweepCursorShape mirrors listingsapp's unexported sweepCursor (backfill.go:96-99)
// structurally, for this external test package's own unmarshal-and-compare —
// it must never assert a byte-exact JSONB string (a named trap in this repo's
// own history: "PG string-eq byte-exact em round-trip JSONB").
type sweepCursorShape struct {
	Phase           string    `json:"phase"`
	LastFullSweepAt time.Time `json:"last_full_sweep_at"`
}

// TestRealListingsJobThroughSchedulerPersistsTerminalSweepCursor is the
// substitute observable for M-04/F-04's daily (24h-interval) listings
// scheduler: nobody can watch a real 24h tick land live, so this test proves
// the wiring end-to-end instead — a REAL *syncapp.Scheduler, running the REAL
// listingsapp.NewListingsJob (backfill.go), against a REAL ephemeral Postgres
// sync_state table (syncpg.NewSyncStateRepository, not a fake StateStore).
//
// It asserts three things a fake-store unit test cannot: (a) the sync_state
// row actually exists in Postgres keyed by entity="listings" after RunOnce,
// (b) the persisted cursor round-trips through real JSONB with the
// {"phase":"sweep","last_full_sweep_at":...} structure intact (compared by
// unmarshalled fields, never a byte-exact string), and (c) RunOnce's own
// incremental-timestamp routing (scheduler.go:160, inferIncremental at
// scheduler.go:178-193) resolves last_full_sync_at forward while leaving
// last_incremental_at untouched, because ADR-026's ratified phase vocabulary
// maps "sweep" to incremental=false (scheduler.go:188-193, default case).
func TestRealListingsJobThroughSchedulerPersistsTerminalSweepCursor(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "listings_scheduler_sync_state")

	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenantID := "tenant-listings-sched-" + token
	installationID := "installation-listings-sched-" + token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenantID)
	})

	at := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	clock := func() time.Time { return at }

	repo := syncpg.NewSyncStateRepository(pool, tenantID)
	scheduler := syncapp.NewScheduler(repo, installationID, time.Hour, clock)

	account := listingsports.InstallationAccount{
		TenantID:          tenantID,
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "seller-listings-sched-" + token,
	}
	store := &noopCompletedPullStore{}
	job := listingsapp.NewListingsJob(zeroPageEnumerator{}, neverCalledHydrator{t: t}, store, account, clock)

	if err := scheduler.RegisterJob(domain.EntityListings, job); err != nil {
		t.Fatalf("RegisterJob(EntityListings): %v", err)
	}

	scheduler.RunOnce(ctx)

	if store.markRunCompleteCalls != 1 {
		t.Fatalf("MarkRunComplete calls = %d, want 1 (the terminal empty page must still complete the run)", store.markRunCompleteCalls)
	}

	// (a) the row exists, keyed by entity="listings", via a REAL Postgres read
	// (not the fake StateStore products_regression_test.go uses).
	state, found, err := repo.Read(ctx, installationID, domain.EntityListings)
	if err != nil {
		t.Fatalf("Read after RunOnce: %v", err)
	}
	if !found {
		t.Fatal("found = false, want true: RunOnce must have persisted a sync_state row for entity=listings")
	}

	// (b) structural cursor comparison — never byte-exact JSONB.
	var cursor sweepCursorShape
	if err := json.Unmarshal(state.Cursor, &cursor); err != nil {
		t.Fatalf("persisted cursor %q did not unmarshal: %v", state.Cursor, err)
	}
	if cursor.Phase != "sweep" {
		t.Errorf("cursor.Phase = %q, want %q", cursor.Phase, "sweep")
	}
	if !cursor.LastFullSweepAt.Equal(at) {
		t.Errorf("cursor.LastFullSweepAt = %v, want %v (the fixed clock value, stored verbatim by RecordSuccess)", cursor.LastFullSweepAt, at)
	}

	// (c) last_full_sync_at advanced, last_incremental_at did NOT — proving
	// RunOnce's incremental flag resolved false for a "sweep"-phase cursor,
	// through the REAL scheduler's runJob -> inferIncremental path
	// (scheduler.go:160 calls RecordSuccess with inferIncremental(next);
	// inferIncremental's switch at scheduler.go:188-193 has no "sweep" case,
	// so it falls to the default: false).
	if state.LastFullSyncAt == nil {
		t.Fatal("LastFullSyncAt = nil, want the fixed clock value: sweep must resolve incremental=false")
	}
	if !state.LastFullSyncAt.Equal(at) {
		t.Errorf("LastFullSyncAt = %v, want %v", state.LastFullSyncAt, at)
	}
	if state.LastIncrementalAt != nil {
		t.Errorf("LastIncrementalAt = %v, want nil: a sweep-phase cursor must never advance the incremental timestamp", state.LastIncrementalAt)
	}
}
