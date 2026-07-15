package application

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

func TestIngestionWalksEveryConnectorPageBeforePersistence(t *testing.T) {
	account := testAccount()
	finished := false
	source := &stubPageSource{pages: []ports.ListingPage{
		{ProviderItemCount: 2, Rows: []listingsdomain.Listing{testListing(t, "item-1", "-"), testListing(t, "item-2", "v-1"), testListing(t, "item-2", "v-2")}},
		{ProviderItemCount: 2, Rows: []listingsdomain.Listing{testListing(t, "item-3", "-"), testListing(t, "item-4", "-")}},
		{ProviderItemCount: 1, Rows: []listingsdomain.Listing{testListing(t, "item-5", "-")}},
	}, onRead: func(index int) {
		if index == 2 {
			finished = true
		}
	}}
	store := &stubCompletedStore{beforeApply: func() {
		if !finished {
			t.Fatal("store called before pagination finished")
		}
	}}

	err := NewIngestion(source, store, 2, func() time.Time { return time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC) }).Pull(context.Background(), account)
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if !reflect.DeepEqual(source.cursors, []string{"", "2", "4"}) {
		t.Fatalf("cursors = %#v", source.cursors)
	}
	if store.calls != 1 {
		t.Fatalf("store calls = %d, want 1", store.calls)
	}
	want := []string{"item-1/-", "item-2/v-1", "item-2/v-2", "item-3/-", "item-4/-", "item-5/-"}
	if !reflect.DeepEqual(keys(store.rows), want) {
		t.Fatalf("stored keys = %#v, want %#v", keys(store.rows), want)
	}
}

func TestIngestionCapabilityErrorMidPullDoesNotApplySnapshot(t *testing.T) {
	capabilityErr := connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, "provider unavailable")
	source := &stubPageSource{pages: []ports.ListingPage{{ProviderItemCount: 1, Rows: []listingsdomain.Listing{testListing(t, "item-1", "-")}}}, errs: map[int]error{1: capabilityErr}}
	store := &stubCompletedStore{}
	err := NewIngestion(source, store, 1, time.Now).Pull(context.Background(), testAccount())
	if connectorsdomain.ErrorCodeOf(err) != connectorsdomain.ErrCodeProviderTransient {
		t.Fatalf("ErrorCodeOf(error) = %q; error = %v", connectorsdomain.ErrorCodeOf(err), err)
	}
	if store.calls != 0 {
		t.Fatalf("store calls = %d, want 0", store.calls)
	}
}

func TestIngestionRejectsDuplicateCanonicalKey(t *testing.T) {
	duplicate := testListing(t, "item-1", "variation-1")
	source := &stubPageSource{pages: []ports.ListingPage{{ProviderItemCount: 1, Rows: []listingsdomain.Listing{duplicate}}, {ProviderItemCount: 0, Rows: []listingsdomain.Listing{duplicate}}}}
	store := &stubCompletedStore{}
	err := NewIngestion(source, store, 1, time.Now).Pull(context.Background(), testAccount())
	if !errors.Is(err, ErrDuplicateCanonicalKey) {
		t.Fatalf("Pull() error = %v, want duplicate-key error", err)
	}
	if store.calls != 0 {
		t.Fatalf("store calls = %d, want 0", store.calls)
	}
}

func TestIngestionFailsWhenProviderNeverSignalsFinalPage(t *testing.T) {
	// A connector that always returns a full page (distinct keys, so the dedup guard
	// never fires) must be stopped by the honest page cap, not loop forever.
	source := &unboundedPageSource{}
	store := &stubCompletedStore{}
	err := NewIngestion(source, store, 1, time.Now).Pull(context.Background(), testAccount())
	if err == nil || !strings.Contains(err.Error(), "exceeded") {
		t.Fatalf("Pull() error = %v; want page-cap exceeded error", err)
	}
	if source.n < maxIngestionPages {
		t.Fatalf("read %d pages; want cap at %d", source.n, maxIngestionPages)
	}
	if store.calls != 0 {
		t.Fatalf("store calls = %d, want 0 (no partial snapshot applied)", store.calls)
	}
}

type unboundedPageSource struct{ n int }

func (s *unboundedPageSource) ReadPage(_ context.Context, _ ports.InstallationAccount, _ string, _ int) (ports.ListingPage, error) {
	s.n++
	fetched := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	item := "item-" + strconv.Itoa(s.n)
	row, err := listingsdomain.NewListing(listingsdomain.ListingInput{TenantID: "tenant-a", InstallationID: "installation-a", Provider: "mercado_livre", ProviderListingID: item, VariationID: "-", Title: item, Status: listingsdomain.ListingStatusActive, SyncState: listingsdomain.ListingSyncStateSynced, FetchedAt: &fetched, CreatedAt: fetched, UpdatedAt: fetched})
	if err != nil {
		return ports.ListingPage{}, err
	}
	return ports.ListingPage{ProviderItemCount: 1, Rows: []listingsdomain.Listing{row}}, nil
}

type stubPageSource struct {
	pages   []ports.ListingPage
	errs    map[int]error
	cursors []string
	onRead  func(int)
}

func (s *stubPageSource) ReadPage(_ context.Context, _ ports.InstallationAccount, cursor string, _ int) (ports.ListingPage, error) {
	i := len(s.cursors)
	s.cursors = append(s.cursors, cursor)
	if s.onRead != nil {
		s.onRead(i)
	}
	if err := s.errs[i]; err != nil {
		return ports.ListingPage{}, err
	}
	return s.pages[i], nil
}

type stubCompletedStore struct {
	calls       int
	rows        []listingsdomain.Listing
	beforeApply func()
}

func (s *stubCompletedStore) ApplyCompletedPull(_ context.Context, _ string, rows []listingsdomain.Listing, _ time.Time) error {
	if s.beforeApply != nil {
		s.beforeApply()
	}
	s.calls++
	s.rows = append([]listingsdomain.Listing(nil), rows...)
	return nil
}
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
func keys(rows []listingsdomain.Listing) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.Key.ProviderListingID + "/" + row.Key.VariationID
	}
	return out
}
