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

// MapMultigetItemsToListings é o Passo 2 do F-03 (backfill-cursor-ingest):
// converte um lote hidratado de mercadolivre.ItemMultigetDTO (Passo 1,
// GetItemsMultiget) em linhas canônicas domain.Listing.
//
// Duas regras deste mapper, ambas exigidas pelo IC-07:
//   - status VERBATIM, sem remapeamento para um conjunto fechado — o
//     vocabulário é do provider, e um remap silenciosamente descartaria
//     status legítimos que a tabela já aceita.
//   - UMA linha por item (VariationID = NoVariationID, ADR-019), com os dados
//     de variação como filhos em Variations, e não uma linha achatada por
//     variação.
//
// Erro em um item (item.Err != nil) é registrado e pulado, nunca aborta o
// lote — mesma garantia do IC-06.
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

	// ADR-17: ausência do provider vira nil, nunca zero — quem não tem preço
	// não tem preço zero. multigetPriceAmount (below) does the same nil/
	// empty-safe conversion for item.Price, reused rather than duplicated.
	var priceCurrency *listingsdomain.PriceCurrency
	if currency := strings.TrimSpace(item.CurrencyID); currency != "" {
		value := listingsdomain.PriceCurrency(currency)
		priceCurrency = &value
	}
	var listingTypeCode *listingsdomain.ListingTypeCode
	if code := strings.TrimSpace(item.ListingTypeID); code != "" {
		value := listingsdomain.ListingTypeCode(code)
		listingTypeCode = &value
	}

	fetchedAtValue := fetchedAt
	row, err := listingsdomain.NewListing(listingsdomain.ListingInput{
		TenantID:          tenantID,
		InstallationID:    installationID,
		Provider:          provider,
		ProviderListingID: providerListingID,
		VariationID:       listingsdomain.NoVariationID,
		Title:             title,
		ListingTypeCode:   listingTypeCode,
		Status:            listingsdomain.ListingStatus(strings.TrimSpace(item.Status)),
		PriceAmount:       multigetPriceAmount(item.Price),
		PriceCurrency:     priceCurrency,
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
		// PublishedQuantity é a quantidade PUBLICADA (initial_quantity), que o
		// ML mantém distinta de available_quantity: reposição sobe a
		// disponível sem tocar a inicial.
		PublishedQuantity: item.InitialQuantity,
		Raw:               item.Raw,
		RawTruncated:      item.RawTruncated,
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
// seam the old single-item ReadListing path feeds (ADR-019, IC-07 "hidratação
// nova CONTINUA alimentando AbsorbProviderSnapshots").
//
// Top-level SellerSKU/EAN are derived from the item's OWN seller_sku/
// seller_custom_field/attributes — mirroring mapListing's oracle behavior
// exactly (capability_adapter.go:832-846, the single-item ReadListing shape
// this ADR-019 seam must match): those three fields are set UNCONDITIONALLY,
// regardless of whether the item has variations, because mapListing does the
// same (it never gates the top-level SellerSKU/EAN on len(item.Variations)).
// This was previously a KNOWN GAP: items_multiget_reader.go's ItemMultigetDTO
// carried no top-level SellerSKU/SellerCustomField/Attributes fields at all
// (json.Unmarshal silently drops undeclared fields), even though ML's
// multiget `body` element is the same full item representation as the
// single-item GET and so genuinely carries those bytes on the wire (visible
// in Raw). That reader now types them (same package, not a frozen file), so
// this mapper can derive the same values mapListing would for equivalent
// input instead of leaving them honest-empty for lack of a typed source.
//
// Per-variation SellerSKU/EAN come from that variation's own fields only —
// deliberately NOT falling back to the parent item's SellerSKU/
// SellerCustomField, because mapListing's own per-variation loop doesn't
// either (capability_adapter.go:852-860; only ReadStock's single-variation
// StockSnapshot path chains through the parent, and StockSnapshot is a
// different shape/consumer than ListingSnapshot).
func MapMultigetItemToListingSnapshot(providerCode string, item mercadolivre.ItemMultigetDTO, fetchedAt time.Time) connectorsdomain.ListingSnapshot {
	snapshot := connectorsdomain.ListingSnapshot{
		ProviderCode:   providerCode,
		ProviderItemID: strings.TrimSpace(item.ProviderItemID),
		ProviderStatus: strings.TrimSpace(item.Status),
		SellerSKU:      firstNonEmpty(item.SellerSKU, item.SellerCustomField),
		EAN:            multigetAttributeValue(item.Attributes, "GTIN", "EAN"),
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

// firstNonEmpty mirrors capability_adapter.go's unexported same-named helper
// (first argument that is non-blank after trimming, "" if none) — duplicated
// here rather than reached into directly because that helper is unexported
// and capability_adapter.go is frozen for this feature; this package needs
// its own copy for the same seller_sku-then-seller_custom_field fallback.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
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
