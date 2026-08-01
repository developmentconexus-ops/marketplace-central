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
		{"under_review", ListingStatusUnderReview},
		{"inactive", ListingStatusInactive},
		{"payment_required", ListingStatusPaymentRequired},
		{"not_yet_active", ListingStatusNotYetActive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.name {
				t.Fatalf("status = %q, want %q", tt.value, tt.name)
			}
			if !tt.value.IsValid() {
				t.Fatalf("status %q should be valid", tt.value)
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

func TestNewListingAcceptsUnlistedProviderStatusVerbatim(t *testing.T) {
	// IC-07: status vocabulary is open (F-01/F-02 dropped the DB CHECK, the
	// writer persists provider status verbatim). "suspended" isn't one of
	// the named ListingStatus* constants, but a real ML item can report it.
	// Domain validation must not be the last silent-rejection point.
	const novel ListingStatus = "suspended"
	listing, err := NewListing(ListingInput{
		TenantID:          "tenant-1",
		InstallationID:    "installation-1",
		Provider:          "mercado_livre",
		ProviderListingID: "item-1",
		Title:             "A listing",
		Status:            novel,
		SyncState:         ListingSyncStateSynced,
	})
	if err != nil {
		t.Fatalf("NewListing() error = %v, want nil for unlisted status %q", err, novel)
	}
	if listing.Status != novel {
		t.Fatalf("listing.Status = %q, want verbatim %q", listing.Status, novel)
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
