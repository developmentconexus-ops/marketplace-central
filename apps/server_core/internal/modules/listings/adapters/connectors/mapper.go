package connectors

import (
	"fmt"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
)

var listingTypeAllowlist = map[string]struct{}{
	"gold_pro":     {},
	"gold_special": {},
	"gold_premium": {},
	"gold":         {},
	"silver":       {},
	"bronze":       {},
	"free":         {},
}

// MapListingSnapshotToCanonicalRows maps one provider page atomically in memory.
// It returns no partial rows when a required provider fact is invalid.
func MapListingSnapshotToCanonicalRows(accountRef connectorsdomain.ProviderAccountRef, snapshots []connectorsdomain.ListingSnapshot) ([]listingsdomain.Listing, error) {
	tenantID := strings.TrimSpace(accountRef.TenantID)
	installationID := strings.TrimSpace(accountRef.InstallationID)
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if installationID == "" {
		return nil, fmt.Errorf("installation_id is required")
	}

	rows := make([]listingsdomain.Listing, 0, len(snapshots))
	for index, snapshot := range snapshots {
		provider := strings.TrimSpace(snapshot.ProviderCode)
		if provider == "" {
			return nil, fmt.Errorf("snapshot %d: provider is required", index)
		}
		providerListingID := strings.TrimSpace(snapshot.ProviderItemID)
		if providerListingID == "" {
			return nil, fmt.Errorf("snapshot %d: provider item id is required", index)
		}
		title := strings.TrimSpace(snapshot.Title)
		if title == "" {
			return nil, fmt.Errorf("snapshot %d: title is required", index)
		}
		if snapshot.FetchedAt.IsZero() {
			return nil, fmt.Errorf("snapshot %d: fetched_at is required", index)
		}

		fetchedAt := snapshot.FetchedAt
		listingTypeCode := canonicalListingTypeCode(snapshot.ListingTypeCode)
		priceAmount, priceCurrency := canonicalPrice(snapshot.PriceAmount, snapshot.PriceCurrency)
		variationRows := snapshot.Variations
		if len(variationRows) == 0 {
			row, err := newCanonicalListing(tenantID, installationID, provider, providerListingID, listingsdomain.NoVariationID, title, snapshot, listingTypeCode, priceAmount, priceCurrency, snapshot.AvailableQuantity, fetchedAt)
			if err != nil {
				return nil, fmt.Errorf("snapshot %d: %w", index, err)
			}
			rows = append(rows, row)
			continue
		}

		for variationIndex, variation := range variationRows {
			variationID := strings.TrimSpace(variation.ProviderVariationID)
			if variationID == "" {
				return nil, fmt.Errorf("snapshot %d variation %d: provider variation id is required", index, variationIndex)
			}
			row, err := newCanonicalListing(tenantID, installationID, provider, providerListingID, variationID, title, snapshot, listingTypeCode, priceAmount, priceCurrency, variation.AvailableQuantity, fetchedAt)
			if err != nil {
				return nil, fmt.Errorf("snapshot %d variation %d: %w", index, variationIndex, err)
			}
			rows = append(rows, row)
		}
	}

	return rows, nil
}

func newCanonicalListing(tenantID, installationID, provider, providerListingID, variationID, title string, snapshot connectorsdomain.ListingSnapshot, listingTypeCode *listingsdomain.ListingTypeCode, priceAmount *listingsdomain.PriceAmount, priceCurrency *listingsdomain.PriceCurrency, quantity *int, fetchedAt time.Time) (listingsdomain.Listing, error) {
	fetchedAtValue := fetchedAt
	return listingsdomain.NewListing(listingsdomain.ListingInput{
		TenantID:          tenantID,
		InstallationID:    installationID,
		Provider:          provider,
		ProviderListingID: providerListingID,
		VariationID:       variationID,
		Title:             title,
		ListingTypeCode:   listingTypeCode,
		Status:            canonicalListingStatus(snapshot.ProviderStatus),
		PriceAmount:       priceAmount,
		PriceCurrency:     priceCurrency,
		PublishedQuantity: quantity,
		SyncState:         listingsdomain.ListingSyncStateSynced,
		FetchedAt:         &fetchedAtValue,
		CreatedAt:         fetchedAt,
		UpdatedAt:         fetchedAt,
	})
}

func canonicalListingStatus(providerStatus string) listingsdomain.ListingStatus {
	switch strings.ToLower(strings.TrimSpace(providerStatus)) {
	case "active":
		return listingsdomain.ListingStatusActive
	case "paused":
		return listingsdomain.ListingStatusPaused
	case "closed":
		return listingsdomain.ListingStatusClosed
	default:
		return listingsdomain.ListingStatusUnknown
	}
}

func canonicalListingTypeCode(providerCode string) *listingsdomain.ListingTypeCode {
	providerCode = strings.TrimSpace(providerCode)
	if _, ok := listingTypeAllowlist[providerCode]; !ok {
		return nil
	}
	code := listingsdomain.ListingTypeCode(providerCode)
	return &code
}

func canonicalPrice(amount *string, currency string) (*listingsdomain.PriceAmount, *listingsdomain.PriceCurrency) {
	if amount == nil || strings.TrimSpace(*amount) == "" || strings.TrimSpace(currency) == "" {
		return nil, nil
	}
	amountValue := listingsdomain.PriceAmount(strings.TrimSpace(*amount))
	currencyValue := listingsdomain.PriceCurrency(strings.TrimSpace(currency))
	return &amountValue, &currencyValue
}
