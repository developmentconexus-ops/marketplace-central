package domain

import (
	"errors"
	"testing"
)

func TestListingStatusValues(t *testing.T) {
	tests := []struct {
		name  string
		value ListingStatus
	}{
		{"active", ListingStatusActive},
		{"paused", ListingStatusPaused},
		{"closed", ListingStatusClosed},
		{"unknown", ListingStatusUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.name {
				t.Fatalf("status = %q, want %q", tt.value, tt.name)
			}
		})
	}
}

func TestListingSyncStateValues(t *testing.T) {
	tests := []struct {
		name  string
		value ListingSyncState
	}{
		{"synced", ListingSyncStateSynced},
		{"error", ListingSyncStateError},
		{"stale", ListingSyncStateStale},
		{"queued", ListingSyncStateQueued},
		{"syncing", ListingSyncStateSyncing},
		{"paused_sync", ListingSyncStatePausedSync},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.name {
				t.Fatalf("sync state = %q, want %q", tt.value, tt.name)
			}
		})
	}
}

func TestListingKeyNormalizesVariationAndPreservesIdentity(t *testing.T) {
	tests := []struct {
		name      string
		variation string
		want      string
	}{
		{"no provider variation", "", "-"},
		{"first variation", "red", "red"},
		{"second variation", "blue", "blue"},
	}
	keys := make([]ListingKey, 0, len(tests))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewListingKey("tenant-1", "installation-1", "item-1", tt.variation)
			if err != nil {
				t.Fatal(err)
			}
			if key.VariationID != tt.want {
				t.Fatalf("variation_id = %q, want %q", key.VariationID, tt.want)
			}
			keys = append(keys, key)
		})
	}
	if keys[1] == keys[2] {
		t.Fatal("different provider variations must produce distinct listing keys")
	}
}

func TestNewListingRejectsEmptyRequiredFields(t *testing.T) {
	base := ListingInput{
		TenantID:          "tenant-1",
		InstallationID:    "installation-1",
		Provider:          "mercado_livre",
		ProviderListingID: "item-1",
		Title:             "A listing",
		Status:            ListingStatusActive,
		SyncState:         ListingSyncStateSynced,
	}
	tests := []struct {
		name string
		edit func(*ListingInput)
		want error
	}{
		{"tenant", func(input *ListingInput) { input.TenantID = "" }, ErrTenantIDRequired},
		{"installation", func(input *ListingInput) { input.InstallationID = "" }, ErrInstallationIDRequired},
		{"provider listing id", func(input *ListingInput) { input.ProviderListingID = "" }, ErrProviderListingIDRequired},
		{"title", func(input *ListingInput) { input.Title = "" }, ErrTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := base
			tt.edit(&input)
			if _, err := NewListing(input); !errors.Is(err, tt.want) {
				t.Fatalf("NewListing() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewListingPreservesNullableFacts(t *testing.T) {
	listing, err := NewListing(ListingInput{
		TenantID:          "tenant-1",
		InstallationID:    "installation-1",
		Provider:          "mercado_livre",
		ProviderListingID: "item-1",
		Title:             "A listing",
		Status:            ListingStatusUnknown,
		SyncState:         ListingSyncStateStale,
	})
	if err != nil {
		t.Fatal(err)
	}
	if listing.ListingTypeCode != nil || listing.PriceAmount != nil || listing.PriceCurrency != nil ||
		listing.SyncError != nil || listing.QualityScore != nil || listing.Sales30D != nil || listing.FetchedAt != nil {
		t.Fatal("unknown nullable facts must remain nil")
	}
	if listing.PublishedQuantity != nil {
		t.Fatal("nil published quantity must not become a zero value")
	}
}
