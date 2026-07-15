package connectors

import (
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
)

func TestMapListingSnapshotToCanonicalRows(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 15, 12, 30, 0, 0, time.UTC)
	quantityItem := 5
	quantityA := 2
	quantityB := 9
	accountRef := connectorsdomain.ProviderAccountRef{
		TenantID:          "tenant-1",
		InstallationID:    "installation-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "seller-1",
	}

	withoutVariations := connectorsdomain.ListingSnapshot{
		ProviderCode:      "mercado_livre",
		ProviderItemID:    "MLB-NO-VARIATIONS",
		ProviderStatus:    "active",
		Title:             "Item simples",
		PriceAmount:       stringPtr("89.00"),
		PriceCurrency:     "BRL",
		ListingTypeCode:   "gold_pro",
		AvailableQuantity: &quantityItem,
		FetchedAt:         fetchedAt,
	}
	withVariations := connectorsdomain.ListingSnapshot{
		ProviderCode:    "mercado_livre",
		ProviderItemID:  "MLB-VARIATIONS",
		ProviderStatus:  "paused",
		Title:           "Item com variações",
		PriceAmount:     stringPtr("89.00"),
		PriceCurrency:   "BRL",
		ListingTypeCode: "gold_special",
		FetchedAt:       fetchedAt,
		Variations: []connectorsdomain.ListingVariationSnapshot{
			{ProviderVariationID: "101", AvailableQuantity: &quantityA},
			{ProviderVariationID: "102", AvailableQuantity: &quantityB},
		},
	}

	rows, err := MapListingSnapshotToCanonicalRows(accountRef, []connectorsdomain.ListingSnapshot{withoutVariations, withVariations})
	if err != nil {
		t.Fatalf("MapListingSnapshotToCanonicalRows() error = %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("row count = %d, want 3", len(rows))
	}

	if got := rows[0]; got.Key.VariationID != listingsdomain.NoVariationID || got.PublishedQuantity == nil || *got.PublishedQuantity != quantityItem {
		t.Fatalf("no-variation row = %#v", got)
	}
	if rows[1].Key.VariationID != "101" || rows[2].Key.VariationID != "102" {
		t.Fatalf("variation ids = %q, %q", rows[1].Key.VariationID, rows[2].Key.VariationID)
	}
	if rows[1].PublishedQuantity == nil || *rows[1].PublishedQuantity != quantityA || rows[2].PublishedQuantity == nil || *rows[2].PublishedQuantity != quantityB {
		t.Fatalf("variation quantities = %#v, %#v", rows[1].PublishedQuantity, rows[2].PublishedQuantity)
	}
	for i, row := range rows {
		if row.Title == "" || row.PriceAmount == nil || *row.PriceAmount != listingsdomain.PriceAmount("89.00") || row.PriceCurrency == nil || *row.PriceCurrency != listingsdomain.PriceCurrency("BRL") || row.ListingTypeCode == nil || row.FetchedAt == nil || !row.FetchedAt.Equal(fetchedAt) {
			t.Fatalf("row %d propagated facts = %#v", i, row)
		}
		if row.SyncState != listingsdomain.ListingSyncStateSynced || row.SyncError != nil || row.QualityScore != nil || row.Sales30D != nil {
			t.Fatalf("row %d nullable/sync facts = %#v", i, row)
		}
	}

	statusCases := []struct {
		provider string
		want     listingsdomain.ListingStatus
	}{
		{provider: "active", want: listingsdomain.ListingStatusActive},
		{provider: "paused", want: listingsdomain.ListingStatusPaused},
		{provider: "closed", want: listingsdomain.ListingStatusClosed},
		{provider: "something-new", want: listingsdomain.ListingStatusUnknown},
	}
	for _, tt := range statusCases {
		snapshot := withoutVariations
		snapshot.ProviderItemID = "MLB-STATUS-" + tt.provider
		snapshot.ProviderStatus = tt.provider
		rows, err := MapListingSnapshotToCanonicalRows(accountRef, []connectorsdomain.ListingSnapshot{snapshot})
		if err != nil {
			t.Fatalf("status %q error = %v", tt.provider, err)
		}
		if rows[0].Status != tt.want {
			t.Fatalf("status %q = %q, want %q", tt.provider, rows[0].Status, tt.want)
		}
	}

	for _, modality := range []string{"not-a-modality", ""} {
		snapshot := withoutVariations
		snapshot.ProviderItemID = "MLB-MODALITY-" + modality
		snapshot.ListingTypeCode = modality
		rows, err := MapListingSnapshotToCanonicalRows(accountRef, []connectorsdomain.ListingSnapshot{snapshot})
		if err != nil {
			t.Fatalf("modality %q error = %v", modality, err)
		}
		if rows[0].ListingTypeCode != nil {
			t.Fatalf("modality %q = %v, want nil", modality, rows[0].ListingTypeCode)
		}
	}
}

func TestMapListingSnapshotToCanonicalRowsPreservesUnknownsAsNil(t *testing.T) {
	t.Parallel()

	snapshot := connectorsdomain.ListingSnapshot{
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB-NULLS",
		ProviderStatus: "active",
		Title:          "Sem fatos",
		FetchedAt:      time.Date(2026, 7, 15, 12, 30, 0, 0, time.UTC),
	}
	rows, err := MapListingSnapshotToCanonicalRows(accountRefForMapperTests(), []connectorsdomain.ListingSnapshot{snapshot})
	if err != nil {
		t.Fatalf("MapListingSnapshotToCanonicalRows() error = %v", err)
	}
	if rows[0].PriceAmount != nil || rows[0].PriceCurrency != nil || rows[0].PublishedQuantity != nil || rows[0].QualityScore != nil || rows[0].Sales30D != nil {
		t.Fatalf("unknown facts were defaulted: %#v", rows[0])
	}
}

func TestMapListingSnapshotToCanonicalRowsRejectsMissingRequiredFacts(t *testing.T) {
	t.Parallel()

	base := connectorsdomain.ListingSnapshot{
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB-REQUIRED",
		ProviderStatus: "active",
		Title:          "Identidade válida",
		FetchedAt:      time.Date(2026, 7, 15, 12, 30, 0, 0, time.UTC),
	}
	tests := []struct {
		name   string
		mutate func(*connectorsdomain.ListingSnapshot)
	}{
		{name: "tenant", mutate: func(snapshot *connectorsdomain.ListingSnapshot) { /* snapshot stays valid; only account.TenantID is blanked below to isolate the tenant-required path */ }},
		{name: "provider item", mutate: func(snapshot *connectorsdomain.ListingSnapshot) { snapshot.ProviderItemID = "" }},
		{name: "provider", mutate: func(snapshot *connectorsdomain.ListingSnapshot) { snapshot.ProviderCode = "" }},
		{name: "title", mutate: func(snapshot *connectorsdomain.ListingSnapshot) { snapshot.Title = "" }},
		{name: "fetched at", mutate: func(snapshot *connectorsdomain.ListingSnapshot) { snapshot.FetchedAt = time.Time{} }},
		{name: "variation identity", mutate: func(snapshot *connectorsdomain.ListingSnapshot) {
			snapshot.Variations = []connectorsdomain.ListingVariationSnapshot{{ProviderVariationID: ""}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := base
			tt.mutate(&snapshot)
			account := accountRefForMapperTests()
			if tt.name == "tenant" {
				account.TenantID = ""
			}
			if _, err := MapListingSnapshotToCanonicalRows(account, []connectorsdomain.ListingSnapshot{snapshot}); err == nil {
				t.Fatalf("expected missing %s to fail honestly", tt.name)
			}
		})
	}
}

func accountRefForMapperTests() connectorsdomain.ProviderAccountRef {
	return connectorsdomain.ProviderAccountRef{
		TenantID:          "tenant-1",
		InstallationID:    "installation-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "seller-1",
	}
}

func stringPtr(value string) *string {
	return &value
}
