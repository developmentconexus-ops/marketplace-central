package connectors

import (
	"errors"
	"testing"
	"time"

	mercadolivre "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
)

func intPtr(v int) *int    { return &v }
func boolPtr(v bool) *bool { return &v }

// TestMapMultigetItemToListingStatusVerbatim is required test (7): status
// persists verbatim through THIS mapper too (no canonicalListingStatus reuse,
// unlike mapper.go), including a value with no named ListingStatus constant.
func TestMapMultigetItemToListingStatusVerbatim(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	item := mercadolivre.ItemMultigetDTO{ProviderItemID: "MLB1", Title: "Item novel status", Status: "suspended"}

	row, err := MapMultigetItemToListing("tenant-1", "installation-1", "mercado_livre", item, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}
	if row.Status != listingsdomain.ListingStatus("suspended") {
		t.Fatalf("Status = %q, want verbatim %q", row.Status, "suspended")
	}
	if row.Key.VariationID != listingsdomain.NoVariationID {
		t.Fatalf("VariationID = %q, want sentinel %q", row.Key.VariationID, listingsdomain.NoVariationID)
	}
}

// TestMapMultigetItemToListingNoVariationsUsesListingGrain proves the IC-07
// available_quantity grain rule's "no variations" branch: the listing-level
// quantity passes through, and Variations stays empty (no child rows).
func TestMapMultigetItemToListingNoVariationsUsesListingGrain(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	item := mercadolivre.ItemMultigetDTO{
		ProviderItemID:    "MLB1",
		Title:             "Item sem variação",
		Status:            "active",
		AvailableQuantity: intPtr(7),
		SoldQuantity:      intPtr(3),
		CategoryID:        "MLB1055",
		Condition:         "new",
		Permalink:         "https://example.com/MLB1",
		Thumbnail:         "https://example.com/thumb.jpg",
		Tags:              []string{"good_quality_thumbnail"},
		CatalogProductID:  "CAT1",
		Shipping:          &mercadolivre.ItemMultigetShippingDTO{Mode: "me2", LogisticType: "fulfillment", FreeShipping: boolPtr(true)},
	}

	row, err := MapMultigetItemToListing("tenant-1", "installation-1", "mercado_livre", item, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}
	if row.AvailableQuantity == nil || *row.AvailableQuantity != 7 {
		t.Fatalf("AvailableQuantity = %v, want 7 (listing grain, no variations)", row.AvailableQuantity)
	}
	if len(row.Variations) != 0 {
		t.Fatalf("Variations = %#v, want empty", row.Variations)
	}
	if row.SoldQuantity == nil || *row.SoldQuantity != 3 {
		t.Fatalf("SoldQuantity = %v, want 3", row.SoldQuantity)
	}
	if row.CategoryID == nil || *row.CategoryID != "MLB1055" {
		t.Fatalf("CategoryID = %v, want MLB1055", row.CategoryID)
	}
	if row.ShippingMode == nil || *row.ShippingMode != "me2" {
		t.Fatalf("ShippingMode = %v, want me2", row.ShippingMode)
	}
	if row.LogisticType == nil || *row.LogisticType != "fulfillment" {
		t.Fatalf("LogisticType = %v, want fulfillment", row.LogisticType)
	}
	if row.FreeShipping == nil || !*row.FreeShipping {
		t.Fatalf("FreeShipping = %v, want true", row.FreeShipping)
	}
	if len(row.Tags) != 1 || row.Tags[0] != "good_quality_thumbnail" {
		t.Fatalf("Tags = %#v", row.Tags)
	}
}

// TestMapMultigetItemToListingWithVariationsUsesVariationGrain proves the
// OTHER half of the grain rule: when variations exist, the PARENT
// available_quantity is nil (never a materialized sum — reader sums the
// child rows), and each variation becomes ONE ListingVariation entry on the
// SAME parent row's Variations slice (one listings row, N child rows — IC-07
// canonical shape), never a separate flattened `listings` row per variation.
func TestMapMultigetItemToListingWithVariationsUsesVariationGrain(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	price181 := "79.90"
	item := mercadolivre.ItemMultigetDTO{
		ProviderItemID: "MLB2",
		Title:          "Item com variações",
		Status:         "active",
		Variations: []mercadolivre.ItemMultigetVariationDTO{
			{
				ProviderVariationID: "181",
				Price:               &price181,
				AvailableQuantity:   intPtr(5),
				SoldQuantity:        intPtr(1),
				SellerSKU:           "SKU-A",
				Attributes:          []mercadolivre.ItemMultigetAttributeDTO{{ID: "GTIN", ValueName: "7891234567890"}},
			},
			{ProviderVariationID: "182", AvailableQuantity: intPtr(2)},
		},
	}

	row, err := MapMultigetItemToListing("tenant-1", "installation-1", "mercado_livre", item, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}
	if row.Key.VariationID != listingsdomain.NoVariationID {
		t.Fatalf("parent VariationID = %q, want sentinel", row.Key.VariationID)
	}
	if row.AvailableQuantity != nil {
		t.Fatalf("parent AvailableQuantity = %v, want nil (grain moved to variations)", *row.AvailableQuantity)
	}
	if len(row.Variations) != 2 {
		t.Fatalf("Variations count = %d, want 2", len(row.Variations))
	}
	v1 := row.Variations[0]
	if v1.VariationID != "181" || v1.Price == nil || string(*v1.Price) != "79.90" {
		t.Fatalf("variation[0] = %#v", v1)
	}
	if v1.AvailableQuantity == nil || *v1.AvailableQuantity != 5 {
		t.Fatalf("variation[0].AvailableQuantity = %v, want 5", v1.AvailableQuantity)
	}
	if v1.SoldQuantity == nil || *v1.SoldQuantity != 1 {
		t.Fatalf("variation[0].SoldQuantity = %v, want 1", v1.SoldQuantity)
	}
	if v1.SellerSKU == nil || *v1.SellerSKU != "SKU-A" {
		t.Fatalf("variation[0].SellerSKU = %v, want SKU-A", v1.SellerSKU)
	}
	if v1.Attributes["GTIN"] != "7891234567890" {
		t.Fatalf("variation[0].Attributes = %#v", v1.Attributes)
	}
	v2 := row.Variations[1]
	if v2.VariationID != "182" || v2.Price != nil || v2.SellerSKU != nil {
		t.Fatalf("variation[1] = %#v, want price/sku nil (DTO omitted them)", v2)
	}
}

// TestMapMultigetItemsToListingsSkipsPerItemErrorsWithoutAborting is part of
// IC-06's "erro por item não aborta run": a batch of DTOs where one item
// already failed multiget (item.Err set) and another fails mapping
// (missing title) must still yield every OTHER valid row, with both
// failures reported back, never a batch-level abort.
func TestMapMultigetItemsToListingsSkipsPerItemErrorsWithoutAborting(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	items := []mercadolivre.ItemMultigetDTO{
		{ProviderItemID: "MLB-OK-1", Title: "Item ok 1", Status: "active"},
		{ProviderItemID: "MLB-404", Err: errors.New("provider resource was not found")},
		{ProviderItemID: "MLB-NO-TITLE", Title: "", Status: "active"},
		{ProviderItemID: "MLB-OK-2", Title: "Item ok 2", Status: "paused"},
	}

	rows, itemErrs := MapMultigetItemsToListings("tenant-1", "installation-1", "mercado_livre", items, fetchedAt)
	if len(rows) != 2 {
		t.Fatalf("rows = %#v, want 2 valid rows", rows)
	}
	if rows[0].Key.ProviderListingID != "MLB-OK-1" || rows[1].Key.ProviderListingID != "MLB-OK-2" {
		t.Fatalf("rows out of order or wrong ids: %#v", rows)
	}
	if len(itemErrs) != 2 {
		t.Fatalf("itemErrs = %#v, want 2", itemErrs)
	}
	if itemErrs[0].ProviderItemID != "MLB-404" || itemErrs[1].ProviderItemID != "MLB-NO-TITLE" {
		t.Fatalf("itemErrs identity wrong: %#v", itemErrs)
	}
}

// TestMapMultigetItemToListingRequiresVariationID proves a variation with a
// blank id fails mapping (never silently coerced to the parent sentinel,
// which would collide two distinct variations onto one child PK).
func TestMapMultigetItemToListingRequiresVariationID(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	item := mercadolivre.ItemMultigetDTO{
		ProviderItemID: "MLB3",
		Title:          "Item variação sem id",
		Status:         "active",
		Variations:     []mercadolivre.ItemMultigetVariationDTO{{ProviderVariationID: "  "}},
	}
	if _, err := MapMultigetItemToListing("tenant-1", "installation-1", "mercado_livre", item, fetchedAt); err == nil {
		t.Fatal("MapMultigetItemToListing() error = nil, want error for blank variation id")
	}
}

// TestMapMultigetItemToListingSnapshotFeedsObserverHonestly proves the
// ADR-13 observer-feed shape: SellerSKU/EAN come through per-variation
// (attribute-derived) but stay honest-empty at the top level for a
// no-variation item, because mlMultigetItemBody carries no top-level
// attributes (see the doc comment on MapMultigetItemToListingSnapshot) —
// this must never fabricate a value, and must never skip producing a
// snapshot at all (that would starve the observer, the exact must-fail
// ADR-13 guards against).
func TestMapMultigetItemToListingSnapshotFeedsObserverHonestly(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)

	noVariation := mercadolivre.ItemMultigetDTO{ProviderItemID: "MLB1", Title: "Sem variação", Status: "active", AvailableQuantity: intPtr(4)}
	snap := MapMultigetItemToListingSnapshot("mercado_livre", noVariation, fetchedAt)
	if snap.ProviderItemID != "MLB1" || snap.SellerSKU != "" || snap.EAN != "" {
		t.Fatalf("no-variation snapshot = %#v, want empty top-level SellerSKU/EAN", snap)
	}
	if snap.AvailableQuantity == nil || *snap.AvailableQuantity != 4 {
		t.Fatalf("snapshot.AvailableQuantity = %v, want 4", snap.AvailableQuantity)
	}

	withVariation := mercadolivre.ItemMultigetDTO{
		ProviderItemID: "MLB2",
		Title:          "Com variação",
		Status:         "active",
		Variations: []mercadolivre.ItemMultigetVariationDTO{{
			ProviderVariationID: "181",
			SellerSKU:           "SKU-A",
			Attributes:          []mercadolivre.ItemMultigetAttributeDTO{{ID: "EAN", ValueName: "7891234567890"}},
			AvailableQuantity:   intPtr(5),
		}},
	}
	snap2 := MapMultigetItemToListingSnapshot("mercado_livre", withVariation, fetchedAt)
	if len(snap2.Variations) != 1 {
		t.Fatalf("snapshot.Variations = %#v, want 1", snap2.Variations)
	}
	v := snap2.Variations[0]
	if v.SellerSKU != "SKU-A" || v.EAN != "7891234567890" {
		t.Fatalf("variation snapshot = %#v, want SellerSKU=SKU-A EAN=7891234567890", v)
	}
}
