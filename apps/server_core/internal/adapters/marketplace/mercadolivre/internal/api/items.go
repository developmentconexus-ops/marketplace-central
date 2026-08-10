package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const multigetBatchSize = 20

// ScanPage enumerates ids via search_type=scan. An empty results slice means
// the walk is over (measured behavior of the legacy reader).
type ScanPage struct {
	ScrollID string   `json:"scroll_id"`
	Results  []string `json:"results"`
}

func (c *Client) ScanIDs(ctx context.Context, scrollID string, limit int) (ScanPage, error) {
	q := url.Values{}
	q.Set("search_type", "scan")
	q.Set("limit", strconv.Itoa(limit))
	if s := strings.TrimSpace(scrollID); s != "" {
		q.Set("scroll_id", s)
	}
	var page ScanPage
	err := c.getJSON(ctx, "/users/"+url.PathEscape(c.userID)+"/items/search?"+q.Encode(), &page)
	return page, err
}

// Item is the multiget body wire shape — the fields this slice consumes plus
// Raw, which retains the channel's full bytes for the observation record.
type Item struct {
	ID                string          `json:"id"`
	Title             string          `json:"title"`
	Status            string          `json:"status"`
	ListingTypeID     string          `json:"listing_type_id"`
	Price             *json.Number    `json:"price"`
	CurrencyID        string          `json:"currency_id"`
	AvailableQuantity *int            `json:"available_quantity"`
	SellerSKU         string          `json:"seller_sku"`
	Attributes        []Attribute     `json:"attributes"`
	Variations        []Variation     `json:"variations"`
	Raw               json.RawMessage `json:"-"`
}

type Attribute struct {
	ID        string `json:"id"`
	ValueName string `json:"value_name"`
}

type Variation struct {
	ID                json.Number  `json:"id"`
	Price             *json.Number `json:"price"`
	AvailableQuantity *int         `json:"available_quantity"`
	SellerSKU         string       `json:"seller_sku"`
	Attributes        []Attribute  `json:"attributes"`
}

// GTIN resolves the EAN attribute: id GTIN with EAN as fallback, value_name —
// the exact rule the legacy mapper measured (multiget_mapper.go:197).
func GTIN(attrs []Attribute) string {
	for _, want := range []string{"GTIN", "EAN"} {
		for _, a := range attrs {
			if a.ID == want && strings.TrimSpace(a.ValueName) != "" {
				return strings.TrimSpace(a.ValueName)
			}
		}
	}
	return ""
}

type multigetElement struct {
	Code int             `json:"code"`
	Body json.RawMessage `json:"body"`
}

// ItemsMultiget hydrates ids in request order, batching at 20 (measured ML
// cap). A per-item code!=200 fails the WHOLE call naming the id: a silently
// skipped listing would read as "absent from the channel", which is a
// different fact.
func (c *Client) ItemsMultiget(ctx context.Context, ids []string) ([]Item, error) {
	items := make([]Item, 0, len(ids))
	for start := 0; start < len(ids); start += multigetBatchSize {
		end := start + multigetBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[start:end]
		var elements []multigetElement
		if err := c.getJSON(ctx, "/items?ids="+url.QueryEscape(strings.Join(batch, ",")), &elements); err != nil {
			return nil, err
		}
		if len(elements) != len(batch) {
			return nil, fmt.Errorf("mercadolivre api: multiget asked %d ids, got %d elements", len(batch), len(elements))
		}
		for i, el := range elements {
			if el.Code != 200 {
				return nil, fmt.Errorf("mercadolivre api: multiget item %s returned code %d", batch[i], el.Code)
			}
			var item Item
			if err := json.Unmarshal(el.Body, &item); err != nil {
				return nil, fmt.Errorf("mercadolivre api: decode item %s: %w", batch[i], err)
			}
			item.Raw = el.Body
			items = append(items, item)
		}
	}
	return items, nil
}
