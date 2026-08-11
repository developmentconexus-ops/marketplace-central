package application_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func storedKey(t *testing.T, listingID string) contracts.SourceListingKey {
	t.Helper()
	tid, err := tenant.Parse("tenant_default")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	code, err := channel.ParseCode("mercado_livre")
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	account, err := channel.NewAccountRef(code, "179571326")
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	key, err := contracts.NewSourceListingKey(tid, account, listingID)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	return key
}

func seedStored(t *testing.T, store *memStore, ids ...string) {
	t.Helper()
	for _, id := range ids {
		store.stored = append(store.stored, contracts.StoredObservation{
			Key:         storedKey(t, id),
			RawPayload:  []byte(`{"id":"` + id + `"}`),
			PayloadHash: "hash-" + id,
			ObservedAt:  time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
		})
	}
}

// TestStoredObservationsWalksEveryRowExactlyOnce is the property a reprocess
// depends on and the one a cursor bug destroys silently: a walk that skips a
// row leaves it holding the old mapping's facts forever, and a walk that
// repeats one does redundant work while reporting a count nobody can check.
func TestStoredObservationsWalksEveryRowExactlyOnce(t *testing.T) {
	store := newMemStore()
	seedStored(t, store, "MLB1", "MLB2", "MLB3", "MLB4", "MLB5")
	svc := application.NewService(store)
	tid, _ := tenant.Parse("tenant_default")

	seen := []string{}
	cursor := contracts.StoredCursor{}
	pages := 0
	for {
		page, err := svc.StoredObservations(context.Background(), tid, cursor, 2)
		if err != nil {
			t.Fatalf("page %d: %v", pages+1, err)
		}
		pages++
		for _, o := range page.Observations {
			seen = append(seen, o.Key.ListingID())
		}
		if page.Done {
			break
		}
		if page.Next.IsStart() {
			t.Fatal("a page that is not Done must advance the cursor; this one handed back the start")
		}
		cursor = page.Next
		if pages > 10 {
			t.Fatal("walk did not terminate")
		}
	}

	want := []string{"MLB1", "MLB2", "MLB3", "MLB4", "MLB5"}
	if len(seen) != len(want) {
		t.Fatalf("saw %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("saw %v, want %v", seen, want)
		}
	}
	// 5 rows at 2 per page: two full pages, one short page that ends the walk.
	if pages != 3 {
		t.Fatalf("pages = %d, want 3", pages)
	}
}

// TestStoredObservationsReportsDoneOnAnEmptyStore: nothing stored is a legal
// answer and must end the walk, not loop forever on a start cursor.
func TestStoredObservationsReportsDoneOnAnEmptyStore(t *testing.T) {
	svc := application.NewService(newMemStore())
	tid, _ := tenant.Parse("tenant_default")

	page, err := svc.StoredObservations(context.Background(), tid, contracts.StoredCursor{}, 50)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	if !page.Done || len(page.Observations) != 0 {
		t.Fatalf("page = %+v, want Done with no observations", page)
	}
}

// TestStoredObservationsRejectsANonPositiveLimit: a limit of zero would ask
// the store for nothing and read as "the walk is over", turning a caller's bug
// into a reprocess that silently reprocessed nothing.
func TestStoredObservationsRejectsANonPositiveLimit(t *testing.T) {
	svc := application.NewService(newMemStore())
	tid, _ := tenant.Parse("tenant_default")

	for _, limit := range []int{0, -1} {
		if _, err := svc.StoredObservations(context.Background(), tid, contracts.StoredCursor{}, limit); err == nil {
			t.Fatalf("limit %d was accepted", limit)
		}
	}
}

// TestStoredObservationsScopesTheTenant: the store is asked for one tenant's
// rows, and the service must pass the tenant it was given rather than let a
// walk cross tenants.
func TestStoredObservationsScopesTheTenant(t *testing.T) {
	store := newMemStore()
	seedStored(t, store, "MLB1")
	svc := application.NewService(store)
	tid, _ := tenant.Parse("tenant_default")

	if _, err := svc.StoredObservations(context.Background(), tid, contracts.StoredCursor{}, 50); err != nil {
		t.Fatalf("page: %v", err)
	}
	if store.storedTenant != tid {
		t.Fatalf("store was asked for tenant %q, want %q", store.storedTenant, tid)
	}
}
