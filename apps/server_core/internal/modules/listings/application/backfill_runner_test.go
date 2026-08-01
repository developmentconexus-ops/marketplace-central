package application

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// testAccount and testListing are shared across this package's tests
// (backfill_test.go's NewListingsJob tests as well as the BackfillRunner
// tests below) — originally defined alongside the now-deleted Ingestion
// tests.
func testAccount() ports.InstallationAccount {
	return ports.InstallationAccount{TenantID: "tenant-a", InstallationID: "installation-a", ProviderCode: "mercado_livre", ProviderAccountID: "seller-a"}
}

func testListing(t *testing.T, item, variation string) listingsdomain.Listing {
	t.Helper()
	fetched := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	row, err := listingsdomain.NewListing(listingsdomain.ListingInput{TenantID: "tenant-a", InstallationID: "installation-a", Provider: "mercado_livre", ProviderListingID: item, VariationID: variation, Title: item, Status: listingsdomain.ListingStatusActive, SyncState: listingsdomain.ListingSyncStateSynced, FetchedAt: &fetched, CreatedAt: fetched, UpdatedAt: fetched})
	if err != nil {
		t.Fatal(err)
	}
	return row
}

// unboundedIDEnumerator never signals exhaustion (NextCursor is always
// non-empty) and always reports a distinct id, so the dedup/duplicate-key
// guard the old Ingestion had can never mask an infinite loop here — only
// the honest page cap can stop it.
type unboundedIDEnumerator struct{ n int }

func (e *unboundedIDEnumerator) ReadIDPage(_ context.Context, _ ports.InstallationAccount, _ string, _ int) (ports.BatchIDPage, error) {
	e.n++
	return ports.BatchIDPage{IDs: []string{"item-" + strconv.Itoa(e.n)}, NextCursor: "next"}, nil
}

func TestBackfillRunnerWalksEveryPageAndMarksRunComplete(t *testing.T) {
	t0 := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	enumerator := &fakeEnumerator{pages: []ports.BatchIDPage{
		{IDs: []string{"item-1"}, NextCursor: "page-2"},
		{IDs: []string{"item-2"}, NextCursor: ""},
	}}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{
		"item-1": testListing(t, "item-1", "-"),
		"item-2": testListing(t, "item-2", "-"),
	}}
	store := &fakeRunStore{}

	runner := NewBackfillRunner(enumerator, hydrator, store, fixedClock(t0))
	if err := runner.Pull(context.Background(), testAccount()); err != nil {
		t.Fatalf("Pull() error = %v", err)
	}

	if !reflectCursorsEqual(enumerator.cursors, []string{"", "page-2"}) {
		t.Fatalf("cursors = %#v", enumerator.cursors)
	}
	if len(store.upserts) != 2 {
		t.Fatalf("upserts = %+v, want 2 pages", store.upserts)
	}
	for i, u := range store.upserts {
		if !u.seenAt.Equal(t0) {
			t.Fatalf("upsert[%d].seenAt = %v, want the run start %v", i, u.seenAt, t0)
		}
	}
	if len(store.completes) != 1 || !store.completes[0].Equal(t0) {
		t.Fatalf("completes = %+v, want exactly [%v]", store.completes, t0)
	}
}

func TestBackfillRunnerEnumeratorErrorMidPullDoesNotMarkComplete(t *testing.T) {
	enumerator := &fakeEnumerator{
		pages: []ports.BatchIDPage{{IDs: []string{"item-1"}, NextCursor: "page-2"}},
		errs:  map[int]error{1: errors.New("scroll expired")},
	}
	hydrator := &fakeHydrator{byID: map[string]listingsdomain.Listing{"item-1": testListing(t, "item-1", "-")}}
	store := &fakeRunStore{}

	runner := NewBackfillRunner(enumerator, hydrator, store, time.Now)
	err := runner.Pull(context.Background(), testAccount())
	if err == nil || !strings.Contains(err.Error(), "enumerate page") {
		t.Fatalf("Pull() error = %v, want an enumerate-page error", err)
	}
	// The first page's rows persisted (per-tick durability, same as
	// NewListingsJob) but the run never completed, so nothing was marked
	// absent.
	if len(store.upserts) != 1 {
		t.Fatalf("upserts = %+v, want the first page's row to have persisted", store.upserts)
	}
	if len(store.completes) != 0 {
		t.Fatalf("completes = %+v, want none (run not exhausted)", store.completes)
	}
}

func TestBackfillRunnerHydratorErrorStopsWithoutMarkingComplete(t *testing.T) {
	enumerator := &fakeEnumerator{pages: []ports.BatchIDPage{{IDs: []string{"item-1"}, NextCursor: ""}}}
	hydrator := &failingHydrator{err: errors.New("provider unavailable")}
	store := &fakeRunStore{}

	runner := NewBackfillRunner(enumerator, hydrator, store, time.Now)
	err := runner.Pull(context.Background(), testAccount())
	if err == nil {
		t.Fatal("Pull() error = nil, want the hydrator failure")
	}
	if len(store.upserts) != 0 || len(store.completes) != 0 {
		t.Fatalf("store must not be touched on hydrate failure: upserts=%+v completes=%+v", store.upserts, store.completes)
	}
}

func TestBackfillRunnerFailsWhenProviderNeverSignalsFinalPage(t *testing.T) {
	enumerator := &unboundedIDEnumerator{}
	store := &fakeRunStore{}
	runner := NewBackfillRunner(enumerator, &fakeHydrator{}, store, time.Now)

	err := runner.Pull(context.Background(), testAccount())
	if err == nil || !strings.Contains(err.Error(), "exceeded") {
		t.Fatalf("Pull() error = %v, want page-cap exceeded error", err)
	}
	if enumerator.n < maxBackfillPages {
		t.Fatalf("read %d pages; want cap at %d", enumerator.n, maxBackfillPages)
	}
	if len(store.completes) != 0 {
		t.Fatalf("completes = %+v, want none", store.completes)
	}
}

func TestBackfillRunnerNotConfiguredGuard(t *testing.T) {
	runner := NewBackfillRunner(nil, nil, nil, time.Now)
	if err := runner.Pull(context.Background(), testAccount()); err == nil {
		t.Fatal("Pull() error = nil, want configuration error")
	}
}

type failingHydrator struct{ err error }

func (h *failingHydrator) HydrateBatch(context.Context, ports.InstallationAccount, []string) ([]listingsdomain.Listing, error) {
	return nil, h.err
}

func reflectCursorsEqual(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
