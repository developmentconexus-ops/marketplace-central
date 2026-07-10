package application

import (
	"context"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubListingSource struct {
	installationID string
	limit          int
	listings       []connectorsdomain.ListingSnapshot
}

func (s *stubListingSource) ListListings(_ context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error) {
	s.installationID = installationID
	s.limit = limit
	return s.listings, nil
}

type stubListingSnapshotStore struct {
	snapshots []productlinksdomain.ListingSnapshot
}

func (s *stubListingSnapshotStore) UpsertListingSnapshots(_ context.Context, snapshots []productlinksdomain.ListingSnapshot) error {
	s.snapshots = append([]productlinksdomain.ListingSnapshot(nil), snapshots...)
	return nil
}

func TestImportListingSnapshotsPersistsItemWithoutVariation(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 19, 0, 0, 0, time.UTC)
	fetchedAt := now.Add(-time.Minute)
	sourceUpdatedAt := now.Add(-2 * time.Minute)
	quantity := 11

	source := &stubListingSource{
		listings: []connectorsdomain.ListingSnapshot{{
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB123",
			ProviderStatus:    "paused",
			SellerSKU:         "SKU-1",
			EAN:               "7891234567890",
			Title:             "Produto base",
			AvailableQuantity: &quantity,
			SourceUpdatedAt:   &sourceUpdatedAt,
			FetchedAt:         fetchedAt,
		}},
	}
	store := &stubListingSnapshotStore{}
	svc := NewImportService(ImportServiceConfig{Source: source, Store: store, Now: func() time.Time { return now }})

	result, err := svc.ImportListingSnapshots(context.Background(), ImportListingSnapshotsInput{
		InstallationID: "inst-1",
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("ImportListingSnapshots() error = %v", err)
	}

	if source.installationID != "inst-1" || source.limit != 25 {
		t.Fatalf("source called with installation=%q limit=%d", source.installationID, source.limit)
	}
	if result.ImportedCount != 1 || len(store.snapshots) != 1 {
		t.Fatalf("imported=%d stored=%d, want 1", result.ImportedCount, len(store.snapshots))
	}
	if store.snapshots[0].ProviderVariationID != "" {
		t.Fatalf("variation id = %q, want empty", store.snapshots[0].ProviderVariationID)
	}
	if store.snapshots[0].Title != "Produto base" || store.snapshots[0].SellerSKU != "SKU-1" {
		t.Fatalf("snapshot=%#v", store.snapshots[0])
	}
}

func TestImportListingSnapshotsFlattensVariationsIntoLinkableRows(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 19, 5, 0, 0, time.UTC)
	fetchedAt := now.Add(-time.Minute)
	qtyA := 3
	qtyB := 7

	source := &stubListingSource{
		listings: []connectorsdomain.ListingSnapshot{{
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB999",
			ProviderStatus: "active",
			Title:          "Produto com variações",
			FetchedAt:      fetchedAt,
			Variations: []connectorsdomain.ListingVariationSnapshot{
				{ProviderVariationID: "VAR-A", SellerSKU: "SKU-A", EAN: "111", AvailableQuantity: &qtyA},
				{ProviderVariationID: "VAR-B", SellerSKU: "SKU-B", EAN: "222", AvailableQuantity: &qtyB},
			},
		}},
	}
	store := &stubListingSnapshotStore{}
	svc := NewImportService(ImportServiceConfig{Source: source, Store: store, Now: func() time.Time { return now }})

	result, err := svc.ImportListingSnapshots(context.Background(), ImportListingSnapshotsInput{
		InstallationID: "inst-2",
		Limit:          10,
	})
	if err != nil {
		t.Fatalf("ImportListingSnapshots() error = %v", err)
	}

	if result.ImportedCount != 2 || len(store.snapshots) != 2 {
		t.Fatalf("imported=%d stored=%d, want 2", result.ImportedCount, len(store.snapshots))
	}
	if store.snapshots[0].ProviderItemID != "MLB999" || store.snapshots[1].ProviderItemID != "MLB999" {
		t.Fatalf("snapshots=%#v", store.snapshots)
	}
	if store.snapshots[0].ProviderVariationID == "" || store.snapshots[1].ProviderVariationID == "" {
		t.Fatalf("snapshots=%#v, want per-variation rows", store.snapshots)
	}
	if store.snapshots[0].Title != "Produto com variações" || store.snapshots[1].ProviderStatus != "active" {
		t.Fatalf("snapshots=%#v", store.snapshots)
	}
}
