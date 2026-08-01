package connectors

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mercadolivre "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// scanMultigetReader embeds recordingListingReader (source_test.go, same
// package) so it satisfies connectorsports.ListingReader for free, and adds
// the two optional capabilities F-03 type-asserts for: ListListingsScanIDs
// (idScanReader) and GetItemsMultiget (multigetReader).
type scanMultigetReader struct {
	recordingListingReader

	scanIDs  []string
	scanNext string
	scanErr  error
	scanCall connectorsdomain.ListListingsInput

	items       []mercadolivre.ItemMultigetDTO
	multigetErr error
	multigetRef connectorsdomain.ProviderAccountRef
	multigetIDs []string
}

func (r *scanMultigetReader) ListListingsScanIDs(_ context.Context, input connectorsdomain.ListListingsInput) ([]string, string, error) {
	r.scanCall = input
	return r.scanIDs, r.scanNext, r.scanErr
}

func (r *scanMultigetReader) GetItemsMultiget(_ context.Context, ref connectorsdomain.ProviderAccountRef, ids []string) ([]mercadolivre.ItemMultigetDTO, error) {
	r.multigetRef = ref
	r.multigetIDs = ids
	return r.items, r.multigetErr
}

func testInstallationAccount() listingsports.InstallationAccount {
	return listingsports.InstallationAccount{TenantID: "tenant-a", InstallationID: "installation-a", ProviderCode: "mercado_livre", ProviderAccountID: "seller-a"}
}

func TestBackfillSourceReadIDPageDelegatesToScanCapability(t *testing.T) {
	reader := &scanMultigetReader{scanIDs: []string{"MLB1", "MLB2"}, scanNext: "next-token"}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})

	page, err := NewBackfillSource(service).ReadIDPage(context.Background(), testInstallationAccount(), "cursor-in", 1000)
	if err != nil {
		t.Fatalf("ReadIDPage() error = %v", err)
	}
	if len(page.IDs) != 2 || page.NextCursor != "next-token" {
		t.Fatalf("page = %+v", page)
	}
	if reader.scanCall.Cursor != "cursor-in" || reader.scanCall.Limit != 1000 || reader.scanCall.AccountRef.ProviderAccountID != "seller-a" {
		t.Fatalf("ListListingsScanIDs input = %+v", reader.scanCall)
	}
}

func TestBackfillSourceReadIDPageFailsWhenProviderLacksScanCapability(t *testing.T) {
	reader := &recordingListingReader{} // only ListListings/ReadListing, no ListListingsScanIDs
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})

	_, err := NewBackfillSource(service).ReadIDPage(context.Background(), testInstallationAccount(), "", 1000)
	if err == nil || !strings.Contains(err.Error(), "does not support ids-only backfill enumeration") {
		t.Fatalf("ReadIDPage() error = %v, want the missing-capability error", err)
	}
}

// TestMultigetHydratorFeedsObserverBeforeMappingRows is the ADR-13 must-fail
// proof (required test 6): every hydrated item must ALSO reach the
// SnapshotObserver, exactly once, with the per-variation EAN/SKU intact — if
// HydrateBatch ever stopped calling AbsorbProviderSnapshots (a regression
// that would starve the product-links matcher on the new backfill/sweep
// path while the old ReadListing path kept working), this assertion is the
// one that catches it.
func TestMultigetHydratorFeedsObserverBeforeMappingRows(t *testing.T) {
	fetched := time.Date(2026, 7, 31, 3, 0, 0, 0, time.UTC)
	reader := &scanMultigetReader{items: []mercadolivre.ItemMultigetDTO{
		{ProviderItemID: "MLB1", Title: "Item 1", Status: "active"},
		{ProviderItemID: "MLB2", Title: "Item 2", Status: "active", Variations: []mercadolivre.ItemMultigetVariationDTO{
			{ProviderVariationID: "V1", SellerSKU: "SKU-1", Attributes: []mercadolivre.ItemMultigetAttributeDTO{{ID: "GTIN", ValueName: "7891234567890"}}},
		}},
	}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	observer := &recordingSnapshotObserver{}

	rows, err := NewMultigetHydrator(service, observer, func() time.Time { return fetched }).HydrateBatch(context.Background(), testInstallationAccount(), []string{"MLB1", "MLB2"})
	if err != nil {
		t.Fatalf("HydrateBatch() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	if observer.installationID != "installation-a" || len(observer.snapshots) != 2 {
		t.Fatalf("observer = %+v", observer)
	}
	if len(observer.snapshots[1].Variations) != 1 || observer.snapshots[1].Variations[0].EAN != "7891234567890" || observer.snapshots[1].Variations[0].SellerSKU != "SKU-1" {
		t.Fatalf("observed variation snapshot lost provider identity: %+v", observer.snapshots[1])
	}
}

// TestMultigetHydratorFailsWhenObserverFails mirrors source.go's own
// fail-closed contract (TestSourceFailsThePageWhenTheObserverFails): a pull
// that silently skipped the linker would report a complete catalog while the
// matcher stayed blind to part of it, so the whole batch fails instead.
func TestMultigetHydratorFailsWhenObserverFails(t *testing.T) {
	reader := &scanMultigetReader{items: []mercadolivre.ItemMultigetDTO{{ProviderItemID: "MLB1", Title: "Item 1", Status: "active"}}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	observer := &recordingSnapshotObserver{err: errors.New("snapshot store down")}

	rows, err := NewMultigetHydrator(service, observer, time.Now).HydrateBatch(context.Background(), testInstallationAccount(), []string{"MLB1"})
	if err == nil {
		t.Fatal("HydrateBatch() error = nil, want the observer failure")
	}
	if rows != nil {
		t.Fatalf("rows = %+v, want nil on observer failure", rows)
	}
}

func TestMultigetHydratorFailsWhenProviderLacksMultigetCapability(t *testing.T) {
	reader := &recordingListingReader{} // only ListListings/ReadListing, no GetItemsMultiget
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})

	_, err := NewMultigetHydrator(service, nil, time.Now).HydrateBatch(context.Background(), testInstallationAccount(), []string{"MLB1"})
	if err == nil || !strings.Contains(err.Error(), "does not support multiget hydration") {
		t.Fatalf("HydrateBatch() error = %v, want the missing-capability error", err)
	}
}

// TestMultigetHydratorSkipsPerItemErrorsWithoutAbortingBatch is IC-06's
// "erro por item nunca aborta o batch": one bad slot in a multiget response
// must not lose the rest of the batch, and must not reach the observer
// (an errored item never produced identity fields to absorb).
func TestMultigetHydratorSkipsPerItemErrorsWithoutAbortingBatch(t *testing.T) {
	reader := &scanMultigetReader{items: []mercadolivre.ItemMultigetDTO{
		{ProviderItemID: "MLB1", Err: errors.New("404")},
		{ProviderItemID: "MLB2", Title: "Item 2", Status: "active"},
	}}
	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", Listings: reader}})
	observer := &recordingSnapshotObserver{}

	rows, err := NewMultigetHydrator(service, observer, time.Now).HydrateBatch(context.Background(), testInstallationAccount(), []string{"MLB1", "MLB2"})
	if err != nil {
		t.Fatalf("HydrateBatch() error = %v, want nil (per-item error must not abort the batch)", err)
	}
	if len(rows) != 1 || rows[0].Key.ProviderListingID != "MLB2" {
		t.Fatalf("rows = %+v", rows)
	}
	if len(observer.snapshots) != 1 {
		t.Fatalf("observer snapshots = %d, want 1 (the errored item must never reach the observer)", len(observer.snapshots))
	}
}
