package mercadolivre

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// Ownership finding (feature brief F-02 point 5 / planning audit P5 F-9,
// cited by M-05/F-01-camada2-fee-ingest/feature.md:47-49 and
// M-05/validation-contract.md:16-17): does GET /items?ids= (or a per-item
// multiget variant) expose sale_price with ?context=channel_marketplace, so
// M-05's fee-layer ingest could read it straight off the multiget hydration
// instead of a second call?
//
// Evidence gathered from THIS repo (no live ML call made — none is available
// from this worktree):
//   - pricing_reader.go:50-58 (getOwnItemPricing) already hits the DEDICATED
//     single-item endpoint GET /items/{id}/sale_price?context=channel_marketplace
//     for the published price (amount/regular_amount) — a call the codebase
//     needed to add SEPARATELY from the plain item read, because ML 400s that
//     endpoint without the context param (comment there cites
//     FINDING-M02-LIVE-2/D-86). If sale_price data were already present on the
//     multiget/item payload, that dedicated call would be redundant — it isn't.
//   - capability_adapter.go:546-609 (ReadFeeQuote) shows the actual COMMISSION
//     fields (percentage_fee, fixed_fee — domain.FeeQuoteSnapshot) come from a
//     THIRD, still different endpoint: GET /sites/{site}/listing_prices,
//     queried by price+category+listing_type, not by item id at all.
//   - docs/design/MIS-007-ML-SYNC-DESIGN.md:65-66,196,198 states explicitly:
//     commission comes from `listing_prices`; "Comissão dentro de
//     /items/{id}/sale_price: não confirmado — usar listing_prices" (i.e. even
//     the single-item sale_price endpoint is NOT trusted for commission,
//     let alone the multiget item shape, which is documented — see
//     mlMultigetItemBody above — as catalog/listing fields: title, status,
//     sold_quantity, category_id, condition, permalink, thumbnail,
//     date_created, tags, catalog_product_id, shipping, available_quantity,
//     variations. None of ML's documented multiget fields are a computed
//     sale_price/commission breakdown.)
//
// CONCLUSION: the multiget endpoint does NOT expose sale_price/commission —
// there are already THREE separate, purpose-built ML endpoints in this
// codebase for price-adjacent facts (published price via /sale_price,
// commission via /listing_prices, and now item-level concurrent prices via
// /prices below), and none of them overlap with the plain /items or
// /items?ids= item shape. Per the feature brief, this file therefore
// delivers the dedicated GET /items/{id}/prices reader so M-05/F-01 has one
// ready-made source instead of writing its own.
//
// CAVEAT (honest-unknown, not fabricated-verified): the exact response shape
// of GET /items/{id}/prices has no research artifact or live dump in this
// repo (`research/external-ml-api-facts.md` does not cover it) — M-05's own
// validation-contract.md:16-17 already flags that "camada 2 real é
// live-driven". The DTO below types the fields ML's Prices API documents
// (id/type/amount/regular_amount/currency_id) tolerantly (unknown fields
// ignored) and ALSO captures Raw, so a live-drive that finds additional
// fields degrades to Raw rather than sinking the read.
type ItemPricesDTO struct {
	ProviderItemID string
	Prices         []ItemPriceEntryDTO
	Raw            json.RawMessage
	RawTruncated   bool
	FetchedAt      time.Time
}

// ItemPriceEntryDTO is one entry of the `prices` array ML's Prices API
// documents (e.g. type "standard"/"promotion"/"campaign"). Amount/
// RegularAmount are decimal strings (mirrors decimalString's convention
// elsewhere in this package) to avoid float precision drift on money.
type ItemPriceEntryDTO struct {
	Type          string
	Amount        *string
	RegularAmount *string
	Currency      string
}

type mlItemPricesResponse struct {
	ItemID string             `json:"id"`
	Prices []mlItemPriceEntry `json:"prices"`
}

type mlItemPriceEntry struct {
	Type          string       `json:"type"`
	Amount        *json.Number `json:"amount"`
	RegularAmount *json.Number `json:"regular_amount"`
	CurrencyID    string       `json:"currency_id"`
}

// GetItemPrices resolves GET /items/{id}/prices for a single item. See the
// ownership finding above for why this dedicated reader exists alongside
// GetItemsMultiget and the pre-existing sale_price/listing_prices readers.
func (a *CapabilityAdapter) GetItemPrices(ctx context.Context, accountRef domain.ProviderAccountRef, itemID string) (ItemPricesDTO, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return ItemPricesDTO{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return ItemPricesDTO{}, mapPricingReaderError(err)
	}
	return a.getItemPrices(ctx, accountRef, token, strings.TrimSpace(itemID))
}

func (a *CapabilityAdapter) getItemPrices(ctx context.Context, accountRef domain.ProviderAccountRef, token, itemID string) (ItemPricesDTO, error) {
	path := "/items/" + url.PathEscape(itemID) + "/prices"
	// doRaw is the SAME choke point (doRawWithHeaders) every other reader in
	// this package uses, inheriting F-01's resilience decorator for free.
	resp, rawBody, err := a.doRaw(ctx, accountRef, token, http.MethodGet, path, nil)
	if err != nil {
		return ItemPricesDTO{}, mapPricingReaderError(err)
	}
	if statusErr := classifyProviderStatus(resp.StatusCode, http.MethodGet, path, rawBody); statusErr != nil {
		return ItemPricesDTO{}, mapPricingReaderError(statusErr)
	}

	var response mlItemPricesResponse
	if err := json.Unmarshal(rawBody, &response); err != nil {
		return ItemPricesDTO{}, mapPricingReaderError(domain.NewCapabilityError(
			domain.ErrCodeProviderPayloadInvalid,
			providerDiag(http.MethodGet, path, resp.StatusCode, rawBody),
		))
	}

	raw, truncated := capRawMessage(rawBody)
	return ItemPricesDTO{
		ProviderItemID: firstNonEmpty(strings.TrimSpace(response.ItemID), itemID),
		Prices:         mapItemPriceEntries(response.Prices),
		Raw:            raw,
		RawTruncated:   truncated,
		FetchedAt:      a.now(),
	}, nil
}

func mapItemPriceEntries(entries []mlItemPriceEntry) []ItemPriceEntryDTO {
	if len(entries) == 0 {
		return nil
	}
	mapped := make([]ItemPriceEntryDTO, 0, len(entries))
	for _, entry := range entries {
		mapped = append(mapped, ItemPriceEntryDTO{
			Type:          strings.TrimSpace(entry.Type),
			Amount:        decimalString(entry.Amount),
			RegularAmount: decimalString(entry.RegularAmount),
			Currency:      strings.TrimSpace(entry.CurrencyID),
		})
	}
	return mapped
}

// classifyProviderStatus maps an HTTP response status to the same domain
// error taxonomy doJSONWithHeaders uses (capability_adapter.go:643-675),
// duplicated here (rather than calling into capability_adapter.go, off-limits
// for this feature) because this reader needs the RAW body bytes for its Raw
// field, which doJSON does not return to its caller.
func classifyProviderStatus(statusCode int, method, path string, rawBody []byte) error {
	switch {
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		return domain.NewCapabilityError(domain.ErrCodeProviderAuth, "provider credentials were rejected")
	case statusCode == http.StatusNotFound:
		return domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	case statusCode == http.StatusTooManyRequests:
		return domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, "provider rate limit reached")
	case statusCode >= 500:
		return domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	case statusCode >= 400:
		return domain.NewCapabilityError(domain.ErrCodeProviderValidation, providerDiag(method, path, statusCode, rawBody))
	default:
		return nil
	}
}
