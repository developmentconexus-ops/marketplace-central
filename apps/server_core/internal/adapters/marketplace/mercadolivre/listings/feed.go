// Package listings implements contexts/listings/port against Mercado Livre.
// It is the ONLY translator: facts come out with explicit knowledge states,
// and an absent wire field becomes Unknown with a reason, never a zero.
package listings

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

const systemName = "mercado_livre"

type Feed struct {
	client  *api.Client
	account channel.AccountRef
	now     func() time.Time
}

func NewFeed(client *api.Client, account channel.AccountRef, now func() time.Time) *Feed {
	return &Feed{client: client, account: account, now: now}
}

// NextPage: one scan page of ids, hydrated by multiget, mapped to
// observations. The cursor token IS ML's scroll_id — opaque to the context.
func (f *Feed) NextPage(ctx context.Context, t tenant.ID, after port.Cursor, limit int) (port.Page, error) {
	scan, err := f.client.ScanIDs(ctx, after.Token(), limit)
	if err != nil {
		return port.Page{}, err
	}
	if len(scan.Results) == 0 {
		return port.Page{Done: true}, nil
	}
	// A nonempty page with a blank scroll_id is malformed: port.NewCursor("")
	// IS the feed-start cursor (port/feed.go:20), so handing it back here would
	// tell a future caller "go back to page 1" instead of "there is no next
	// page." Today's only caller (composition/listings_ingest.go) happens to
	// catch this one hop later via page.Next.IsStart(), but that is a property
	// of that caller, not of this contract — the adapter is the one place that
	// knows what an empty scroll_id from ML actually means, so it fails here.
	if strings.TrimSpace(scan.ScrollID) == "" {
		return port.Page{}, fmt.Errorf("mercadolivre listings: scan returned %d results with an empty scroll_id", len(scan.Results))
	}
	items, err := f.client.ItemsMultiget(ctx, scan.Results)
	if err != nil {
		return port.Page{}, err
	}
	observations := make([]contracts.ListingObservation, 0, len(items))
	observedAt := f.now().UTC()
	for _, item := range items {
		obs, err := f.mapItem(t, item, observedAt)
		if err != nil {
			return port.Page{}, fmt.Errorf("mercadolivre listings: map %s: %w", item.ID, err)
		}
		observations = append(observations, obs)
	}
	return port.Page{Observations: observations, Next: port.NewCursor(scan.ScrollID)}, nil
}

func (f *Feed) mapItem(t tenant.ID, item api.Item, observedAt time.Time) (contracts.ListingObservation, error) {
	key, err := contracts.NewSourceListingKey(t, f.account, item.ID)
	if err != nil {
		return contracts.ListingObservation{}, err
	}
	evidence, err := provenance.NewEvidence(systemName, "item", item.ID, observedAt, payloadHash(item.Raw))
	if err != nil {
		return contracts.ListingObservation{}, err
	}

	state, err := mapItemState(item, evidence)
	if err != nil {
		return contracts.ListingObservation{}, err
	}
	return contracts.ListingObservation{Key: key, State: state, RawPayload: item.Raw, Evidence: evidence}, nil
}

// mapItemState is THE translation from Mercado Livre's item shape to Listings'
// facts, and the only one. Both entry points go through it: NextPage, which
// maps what it just fetched, and Mapper.MapStored (replay.go), which maps what
// Listings recorded earlier. Two copies of these rules would be two answers to
// "what does this listing say", and the second one would be discovered by a
// reprocess quietly rewriting rows the feed had got right.
func mapItemState(item api.Item, evidence provenance.Evidence) (contracts.ListingState, error) {
	var (
		state contracts.ListingState
		err   error
	)
	if state.Title, err = stringFact(item.Title, "title", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.Status, err = stringFact(item.Status, "status", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.ListingType, err = stringFact(item.ListingTypeID, "listing_type_id", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.Price, err = moneyFact(item.Price, item.CurrencyID, evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.AvailableQuantity, err = intFact(item.AvailableQuantity, "available_quantity", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.SellerSKU, err = stringFact(api.ItemSellerSKU(item), "seller_sku (SELLER_SKU attribute or seller_custom_field)", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	if state.GTIN, err = stringFact(api.GTIN(item.Attributes), "gtin attribute", evidence); err != nil {
		return contracts.ListingState{}, err
	}
	for _, v := range item.Variations {
		mapped := contracts.VariationObservation{VariationID: v.ID.String()}
		if mapped.Price, err = moneyFact(v.Price, item.CurrencyID, evidence); err != nil {
			return contracts.ListingState{}, err
		}
		if mapped.AvailableQuantity, err = intFact(v.AvailableQuantity, "variation available_quantity", evidence); err != nil {
			return contracts.ListingState{}, err
		}
		if mapped.SellerSKU, err = stringFact(api.VariationSellerSKU(item, v), "variation seller_sku (SELLER_SKU attribute or seller_custom_field, item as fallback)", evidence); err != nil {
			return contracts.ListingState{}, err
		}
		if mapped.GTIN, err = stringFact(api.GTIN(v.Attributes), "variation gtin attribute", evidence); err != nil {
			return contracts.ListingState{}, err
		}
		state.Variations = append(state.Variations, mapped)
	}
	return state, nil
}

// payloadHash is the one definition of a listing payload's identity. Both the
// feed (hashing what it fetched) and the reprocess mapper (checking what was
// stored) call it, so the check can never be comparing two different digests
// of the same bytes.
func payloadHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func stringFact(value, field string, e provenance.Evidence) (fact.Fact[string], error) {
	if value == "" {
		return fact.NewUnknown[string]("ml omitted "+field, e)
	}
	return fact.NewKnown(value, e)
}

func intFact(value *int, field string, e provenance.Evidence) (fact.Fact[int], error) {
	if value == nil {
		return fact.NewUnknown[int]("ml omitted "+field, e)
	}
	return fact.NewKnown(*value, e)
}

func moneyFact(price *json.Number, currencyID string, e provenance.Evidence) (fact.Fact[exact.Money], error) {
	if price == nil || currencyID == "" {
		return fact.NewUnknown[exact.Money]("ml omitted price or currency", e)
	}
	currency, err := exact.ParseCurrency(currencyID)
	if err != nil {
		return fact.Fact[exact.Money]{}, err
	}
	money, err := exact.ParseMoney(price.String(), currency)
	if err != nil {
		return fact.Fact[exact.Money]{}, err
	}
	return fact.NewKnown(money, e)
}
