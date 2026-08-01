package connectors

import (
	"fmt"
	"strings"
	"time"

	mercadolivre "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
)

// MultigetItemError names a single hydrated item that could not become a
// canonical row (per-item mapping failure — e.g. missing title/variation id).
// IC-06: an error on one item never aborts the batch; the caller records
// these and moves on.
type MultigetItemError struct {
	ProviderItemID string
	Err            error
}

func (e MultigetItemError) Error() string {
	return fmt.Sprintf("item %s: %v", e.ProviderItemID, e.Err)
}

// MapMultigetItemsToListings is Passo 2 of F-03 (backfill-cursor-ingest): it
// converts a hydrated batch of mercadolivre.ItemMultigetDTO (Passo 1's
// GetItemsMultiget, M-01/F-02) into canonical domain.Listing rows.
//
// Unlike MapListingSnapshotToCanonicalRows (mapper.go, used by the existing
// single-item ReadListing/refresh path), this mapper:
//   - keeps status VERBATIM (no canonicalListingStatus remap) — F-01/F-02 of
//     THIS milestone already dropped the listings.status CHECK so the writer
//     persists whatever the provider says (IC-07); a closed remap here would
//     silently discard provider statuses beyond the 8 named constants (e.g.
//     "suspended") that domain.NewListing's isNonEmpty (0a) now accepts.
//   - emits ONE row per item (VariationID = NoVariationID sentinel, ADR-13),
//     with per-variation data attached as child ListingVariation entries on
//     that single row's Variations slice — matching IC-07's canonical
//     example ("dois níveis, mesma transação") and the shape F-02's own
//     repository test already exercises (repository_integration_test.go:
//     TestApplyCompletedPull_VariationsUpsertIdempotently), NOT the
//     per-variation-flattened-row shape mapper.go uses for xlsx/Sankhya
//     snapshots (those providers hand back per-variation snapshots already;
//     multiget hands back one item with a nested variations array).
//
// A per-item error (item.Err != nil, set by GetItemsMultiget for a 404/429/
// decode failure on that slot) is recorded and skipped, never aborts the
// batch — same IC-06 guarantee.
func MapMultigetItemsToListings(tenantID, installationID, provider string, items []mercadolivre.ItemMultigetDTO, fetchedAt time.Time) ([]listingsdomain.Listing, []MultigetItemError) {
	rows := make([]listingsdomain.Listing, 0, len(items))
	var itemErrs []MultigetItemError
	for _, item := range items {
		if item.Err != nil {
			itemErrs = append(itemErrs, MultigetItemError{ProviderItemID: item.ProviderItemID, Err: item.Err})
			continue
		}
		row, err := MapMultigetItemToListing(tenantID, installationID, provider, item, fetchedAt)
		if err != nil {
			itemErrs = append(itemErrs, MultigetItemError{ProviderItemID: item.ProviderItemID, Err: err})
			continue
		}
		rows = append(rows, row)
	}
	return rows, itemErrs
}

// MapMultigetItemToListing maps one hydrated item. Exported so a caller that
// already has a single DTO (e.g. IngestListing's single-item refresh path)
// can reuse it without going through the batch/error-collection wrapper.
func MapMultigetItemToListing(tenantID, installationID, provider string, item mercadolivre.ItemMultigetDTO, fetchedAt time.Time) (listingsdomain.Listing, error) {
	providerListingID := strings.TrimSpace(item.ProviderItemID)
	if providerListingID == "" {
		return listingsdomain.Listing{}, fmt.Errorf("provider item id is required")
	}
	title := strings.TrimSpace(item.Title)
	if title == "" {
		return listingsdomain.Listing{}, fmt.Errorf("item %s: title is required", providerListingID)
	}

	var shippingMode, logisticType *string
	var freeShipping *bool
	if item.Shipping != nil {
		shippingMode = nonEmptyPtr(item.Shipping.Mode)
		logisticType = nonEmptyPtr(item.Shipping.LogisticType)
		freeShipping = item.Shipping.FreeShipping
	}

	// IC-07 available_quantity grain rule: variation grain when variations
	// exist (the sum is NOT materialized on the parent — a reader sums the
	// child rows); listing grain otherwise. item.AvailableQuantity itself may
	// already be nil (provider omitted it) — that stays honest-nil either way.
	availableQuantity := item.AvailableQuantity
	if len(item.Variations) > 0 {
		availableQuantity = nil
	}

	fetchedAtValue := fetchedAt
	row, err := listingsdomain.NewListing(listingsdomain.ListingInput{
		TenantID:          tenantID,
		InstallationID:    installationID,
		Provider:          provider,
		ProviderListingID: providerListingID,
		VariationID:       listingsdomain.NoVariationID,
		Title:             title,
		Status:            listingsdomain.ListingStatus(strings.TrimSpace(item.Status)),
		SyncState:         listingsdomain.ListingSyncStateSynced,
		FetchedAt:         &fetchedAtValue,
		CreatedAt:         fetchedAt,
		UpdatedAt:         fetchedAt,
		SoldQuantity:      item.SoldQuantity,
		CategoryID:        nonEmptyPtr(item.CategoryID),
		Condition:         nonEmptyPtr(item.Condition),
		Permalink:         nonEmptyPtr(item.Permalink),
		Thumbnail:         nonEmptyPtr(item.Thumbnail),
		DateCreatedML:     item.DateCreated,
		Tags:              item.Tags,
		CatalogProductID:  nonEmptyPtr(item.CatalogProductID),
		ShippingMode:      shippingMode,
		FreeShipping:      freeShipping,
		LogisticType:      logisticType,
		AvailableQuantity: availableQuantity,
	})
	if err != nil {
		return listingsdomain.Listing{}, fmt.Errorf("item %s: %w", providerListingID, err)
	}

	if len(item.Variations) == 0 {
		return row, nil
	}

	row.Variations = make([]listingsdomain.ListingVariation, 0, len(item.Variations))
	for _, variation := range item.Variations {
		variationID := strings.TrimSpace(variation.ProviderVariationID)
		if variationID == "" {
			return listingsdomain.Listing{}, fmt.Errorf("item %s: variation id is required", providerListingID)
		}
		row.Variations = append(row.Variations, listingsdomain.ListingVariation{
			VariationID:       variationID,
			Price:             multigetPriceAmount(variation.Price),
			AvailableQuantity: variation.AvailableQuantity,
			SoldQuantity:      variation.SoldQuantity,
			SellerSKU:         nonEmptyPtr(variation.SellerSKU),
			Attributes:        multigetAttributesToMap(variation.Attributes),
		})
	}
	return row, nil
}

// MapMultigetItemToListingSnapshot converts one hydrated item into the
// provider-agnostic connectorsdomain.ListingSnapshot shape the product-links
// matcher consumes via SnapshotObserver.AbsorbProviderSnapshots (source.go),
// so the NEW multiget-based hydration path keeps feeding the SAME re-vínculo
// seam the old single-item ReadListing path feeds (ADR-13, IC-07 "hidratação
// nova CONTINUA alimentando AbsorbProviderSnapshots").
//
// KNOWN, DOCUMENTED GAP vs the old path: mlItemResponse (capability_adapter.go,
// the single-item ReadListing shape) carries a top-level `attributes` array,
// letting mapListing derive a listing-level EAN via attributeValue(item.
// Attributes, "GTIN","EAN"); mlMultigetItemBody (items_multiget_reader.go,
// M-01/F-02 — frozen, out of this feature's ownership) has NO top-level
// attributes field, only per-variation ones. So for an item WITHOUT
// variations, this snapshot's top-level SellerSKU/EAN stay "" (honest-empty,
// the struct's own zero value for an unset string field — never fabricated).
// For an item WITH variations, each ListingVariationSnapshot's SellerSKU/EAN
// ARE populated from that variation's own fields, matching the old path
// exactly. This does not starve the observer (every item and every
// variation still produces a snapshot, preserving count) — see backfill.go's
// non-regression note and the ADR-13 test for the file:line evidence.
func MapMultigetItemToListingSnapshot(providerCode string, item mercadolivre.ItemMultigetDTO, fetchedAt time.Time) connectorsdomain.ListingSnapshot {
	snapshot := connectorsdomain.ListingSnapshot{
		ProviderCode:   providerCode,
		ProviderItemID: strings.TrimSpace(item.ProviderItemID),
		ProviderStatus: strings.TrimSpace(item.Status),
		Title:          strings.TrimSpace(item.Title),
		FetchedAt:      fetchedAt,
	}
	if len(item.Variations) == 0 {
		snapshot.AvailableQuantity = item.AvailableQuantity
		return snapshot
	}
	snapshot.Variations = make([]connectorsdomain.ListingVariationSnapshot, 0, len(item.Variations))
	for _, variation := range item.Variations {
		snapshot.Variations = append(snapshot.Variations, connectorsdomain.ListingVariationSnapshot{
			ProviderVariationID: strings.TrimSpace(variation.ProviderVariationID),
			SellerSKU:           strings.TrimSpace(variation.SellerSKU),
			EAN:                 multigetAttributeValue(variation.Attributes, "GTIN", "EAN"),
			AvailableQuantity:   variation.AvailableQuantity,
		})
	}
	return snapshot
}

// multigetAttributeValue mirrors capability_adapter.go's unexported
// attributeValue helper (same lookup-by-id, ValueName-then-ValueID
// preference) but operates on the exported ItemMultigetAttributeDTO type
// from THIS package, since capability_adapter.go is frozen for this feature
// and its helper is unexported (unreachable from here regardless).
func multigetAttributeValue(attributes []mercadolivre.ItemMultigetAttributeDTO, keys ...string) string {
	keyset := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keyset[strings.ToUpper(strings.TrimSpace(key))] = struct{}{}
	}
	for _, attribute := range attributes {
		if _, ok := keyset[strings.ToUpper(strings.TrimSpace(attribute.ID))]; !ok {
			continue
		}
		if value := strings.TrimSpace(attribute.ValueName); value != "" {
			return value
		}
		if value := strings.TrimSpace(attribute.ValueID); value != "" {
			return value
		}
	}
	return ""
}

func nonEmptyPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func multigetPriceAmount(price *string) *listingsdomain.PriceAmount {
	if price == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*price)
	if trimmed == "" {
		return nil
	}
	value := listingsdomain.PriceAmount(trimmed)
	return &value
}

func multigetAttributesToMap(attributes []mercadolivre.ItemMultigetAttributeDTO) map[string]any {
	if len(attributes) == 0 {
		return nil
	}
	out := make(map[string]any, len(attributes))
	for _, attribute := range attributes {
		id := strings.TrimSpace(attribute.ID)
		if id == "" {
			continue
		}
		if value := strings.TrimSpace(attribute.ValueName); value != "" {
			out[id] = value
			continue
		}
		if value := strings.TrimSpace(attribute.ValueID); value != "" {
			out[id] = value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
