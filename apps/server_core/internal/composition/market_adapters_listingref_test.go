package composition

import (
	"testing"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
)

// A resolved product link for a no-variation listing carries an EMPTY
// provider_variation_id. Composing it into a ListingID must yield an id that
// ParseListingID accepts (all three segments non-empty) — otherwise own-item
// pricing resolution fabricated a 400 on every variation-less listing, which is
// exactly the demo-core failure the live-drive reprova surfaced (D-86).
func TestListingVariationOrSentinelRoundTripsThroughParse(t *testing.T) {
	for _, providerVariationID := range []string{"", "   "} {
		id := listingsdomain.ListingID{
			InstallationID:    "inst-1",
			ProviderListingID: "MLB4735328201",
			VariationID:       listingVariationOrSentinel(providerVariationID),
		}
		parsed, err := listingsdomain.ParseListingID(id.String())
		if err != nil {
			t.Fatalf("ParseListingID(%q) error = %v; want round-trip success", id.String(), err)
		}
		if parsed.ProviderListingID != "MLB4735328201" || parsed.VariationID != listingsdomain.NoVariationID {
			t.Fatalf("parsed = %#v; want provider MLB4735328201 + NoVariationID sentinel", parsed)
		}
	}
	// A real variation passes through untouched.
	if got := listingVariationOrSentinel("v1"); got != "v1" {
		t.Fatalf("listingVariationOrSentinel(v1) = %q, want v1", got)
	}
}
