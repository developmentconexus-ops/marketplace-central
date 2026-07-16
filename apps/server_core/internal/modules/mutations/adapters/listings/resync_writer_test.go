package listings

import (
	"context"
	"testing"
	"time"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type captureSource struct {
	calls   int
	account listingsports.InstallationAccount
	row     listingsdomain.Listing
}

func (s *captureSource) ReadPage(_ context.Context, account listingsports.InstallationAccount, _ string, _ int) (listingsports.ListingPage, error) {
	s.calls++
	s.account = account
	if s.calls == 1 {
		return listingsports.ListingPage{ProviderItemCount: 1, Rows: []listingsdomain.Listing{s.row}}, nil
	}
	return listingsports.ListingPage{}, nil
}

type captureStore struct {
	installation string
	rows         []listingsdomain.Listing
	completed    time.Time
}

func (s *captureStore) ApplyCompletedPull(_ context.Context, installation string, rows []listingsdomain.Listing, completed time.Time) error {
	s.installation, s.rows, s.completed = installation, rows, completed
	return nil
}

func TestResyncWriterUsesExistingIngestionAndCanonicalUpsert(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	row, err := listingsdomain.NewListing(listingsdomain.ListingInput{TenantID: "tenant-1", InstallationID: "installation-1", Provider: "mercado_livre", ProviderListingID: "MLB-1", Title: "Canonical title", Status: listingsdomain.ListingStatusClosed, SyncState: listingsdomain.ListingSyncStateSynced, FetchedAt: &now, CreatedAt: now, UpdatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	source, store := &captureSource{row: row}, &captureStore{}
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre", ProviderAccountID: "seller-1"}
	writer := NewResyncWriter(source, store, account, 1, func() time.Time { return now.Add(time.Minute) })
	outcome, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "MLB-1"})
	if err != nil || outcome.Failure != nil {
		t.Fatalf("outcome=%+v err=%v", outcome, err)
	}
	if source.calls != 2 || source.account.TenantID != "tenant-1" || source.account.ProviderAccountID != "seller-1" {
		t.Fatalf("read calls=%d account=%+v", source.calls, source.account)
	}
	if store.installation != "installation-1" || len(store.rows) != 1 || store.rows[0].Status != listingsdomain.ListingStatusClosed || store.rows[0].Key.TenantID != "tenant-1" {
		t.Fatalf("persisted installation=%q rows=%+v", store.installation, store.rows)
	}
}

func TestResyncWriterRejectsCrossInstallationBeforeRead(t *testing.T) {
	source := &captureSource{}
	writer := NewResyncWriter(source, &captureStore{}, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre"}, 50, time.Now)
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "other", ListingID: "MLB-1"})
	if err == nil || source.calls != 0 {
		t.Fatalf("err=%v calls=%d", err, source.calls)
	}
}
