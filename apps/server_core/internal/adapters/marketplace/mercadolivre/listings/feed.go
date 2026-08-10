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
	sum := sha256.Sum256(item.Raw)
	evidence, err := provenance.NewEvidence(systemName, "item", item.ID, observedAt, hex.EncodeToString(sum[:]))
	if err != nil {
		return contracts.ListingObservation{}, err
	}

	obs := contracts.ListingObservation{Key: key, Evidence: evidence, RawPayload: item.Raw}
	if obs.Title, err = stringFact(item.Title, "title", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.Status, err = stringFact(item.Status, "status", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.ListingType, err = stringFact(item.ListingTypeID, "listing_type_id", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.Price, err = moneyFact(item.Price, item.CurrencyID, evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.AvailableQuantity, err = intFact(item.AvailableQuantity, "available_quantity", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.SellerSKU, err = stringFact(item.SellerSKU, "seller_sku", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.GTIN, err = stringFact(api.GTIN(item.Attributes), "gtin attribute", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	for _, v := range item.Variations {
		mapped := contracts.VariationObservation{VariationID: v.ID.String()}
		if mapped.Price, err = moneyFact(v.Price, item.CurrencyID, evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.AvailableQuantity, err = intFact(v.AvailableQuantity, "variation available_quantity", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.SellerSKU, err = stringFact(v.SellerSKU, "variation seller_sku", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.GTIN, err = stringFact(api.GTIN(v.Attributes), "variation gtin attribute", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		obs.Variations = append(obs.Variations, mapped)
	}
	return obs, nil
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
