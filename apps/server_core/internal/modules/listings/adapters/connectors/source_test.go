package connectors

import (
	"context"
	"errors"
	"testing"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

func TestSourceSelectsListingCapabilityAndMapsPage(t *testing.T) {
	reader := &recordingListingReader{snapshots: []connectorsdomain.ListingSnapshot{{ProviderCode: "mercado_livre", ProviderItemID: "MLB1", ProviderStatus: "active", Title: "Item", FetchedAt: time.Now(), Variations: []connectorsdomain.ListingVariationSnapshot{{ProviderVariationID: "V1"}, {ProviderVariationID: "V2"}}}}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	account := listingsports.InstallationAccount{TenantID: "tenant-a", InstallationID: "installation-a", ProviderCode: "mercado_livre", ProviderAccountID: "seller-a"}
	page, err := NewSource(service).ReadPage(context.Background(), account, "50", 25)
	if err != nil {
		t.Fatalf("ReadPage() error = %v", err)
	}
	if page.ProviderItemCount != 1 || len(page.Rows) != 2 {
		t.Fatalf("page = %+v", page)
	}
	if reader.input.Cursor != "50" || reader.input.Limit != 25 || reader.input.AccountRef.ProviderAccountID != "seller-a" {
		t.Fatalf("ListListings input = %+v", reader.input)
	}
}

func TestSourceHandsProviderIdentityFieldsToTheObserver(t *testing.T) {
	// EAN and seller SKU exist ONLY on the provider snapshot — the canonical
	// listing row has no column for either — so the observer is the linker's
	// only chance to see them on a full-catalog sync.
	reader := &recordingListingReader{snapshots: []connectorsdomain.ListingSnapshot{{
		ProviderCode: "mercado_livre", ProviderItemID: "MLB1", ProviderStatus: "active",
		Title: "Item", SellerSKU: "90001", EAN: "7891234567890", FetchedAt: time.Now(),
	}}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	observer := &recordingSnapshotObserver{}
	account := listingsports.InstallationAccount{TenantID: "tenant-a", InstallationID: "installation-a", ProviderCode: "mercado_livre", ProviderAccountID: "seller-a"}

	if _, err := NewSourceWithObserver(service, observer).ReadPage(context.Background(), account, "", 25); err != nil {
		t.Fatalf("ReadPage() error = %v", err)
	}
	if observer.installationID != "installation-a" || len(observer.snapshots) != 1 {
		t.Fatalf("observer = %+v", observer)
	}
	if observer.snapshots[0].EAN != "7891234567890" || observer.snapshots[0].SellerSKU != "90001" {
		t.Fatalf("observed snapshot lost provider identity: %+v", observer.snapshots[0])
	}
}

func TestSourceFailsThePageWhenTheObserverFails(t *testing.T) {
	// A pull that silently skipped the linker would report a complete catalog
	// while the matcher stayed blind to part of it.
	reader := &recordingListingReader{snapshots: []connectorsdomain.ListingSnapshot{{ProviderCode: "mercado_livre", ProviderItemID: "MLB1", ProviderStatus: "active", Title: "Item", FetchedAt: time.Now()}}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	observer := &recordingSnapshotObserver{err: errors.New("snapshot store down")}
	account := listingsports.InstallationAccount{TenantID: "tenant-a", InstallationID: "installation-a", ProviderCode: "mercado_livre", ProviderAccountID: "seller-a"}

	if _, err := NewSourceWithObserver(service, observer).ReadPage(context.Background(), account, "", 25); err == nil {
		t.Fatal("ReadPage() error = nil, want the observer failure")
	}
}

type recordingSnapshotObserver struct {
	installationID string
	snapshots      []connectorsdomain.ListingSnapshot
	err            error
}

func (o *recordingSnapshotObserver) AbsorbProviderSnapshots(_ context.Context, installationID string, snapshots []connectorsdomain.ListingSnapshot) error {
	o.installationID = installationID
	o.snapshots = append(o.snapshots, snapshots...)
	return o.err
}

var _ SnapshotObserver = (*recordingSnapshotObserver)(nil)

type recordingListingReader struct {
	input     connectorsdomain.ListListingsInput
	snapshots []connectorsdomain.ListingSnapshot
}

func (r *recordingListingReader) ListListings(_ context.Context, input connectorsdomain.ListListingsInput) ([]connectorsdomain.ListingSnapshot, error) {
	r.input = input
	return r.snapshots, nil
}
func (r *recordingListingReader) ReadListing(context.Context, connectorsdomain.ProviderListingRef) (connectorsdomain.ListingSnapshot, error) {
	return connectorsdomain.ListingSnapshot{}, nil
}

var _ connectorsports.ListingReader = (*recordingListingReader)(nil)
