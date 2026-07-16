package connectors

import (
	"context"
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
