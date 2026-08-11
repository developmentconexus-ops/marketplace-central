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
	SellerCustomField string          `json:"seller_custom_field"`
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
	SellerCustomField string       `json:"seller_custom_field"`
	Attributes        []Attribute  `json:"attributes"`
}

// GTIN resolves the EAN attribute: id GTIN with EAN as fallback, value_name —
// the exact rule the legacy mapper measured (multiget_mapper.go:197).
func GTIN(attrs []Attribute) string {
	return attributeValue(attrs, "GTIN", "EAN")
}

// ItemSellerSKU resolves an item's SKU the way Mercado Livre documents it: the
// SELLER_SKU attribute first, then the legacy seller_custom_field.
//
// There is no top-level seller_sku field on a Mercado Livre item. The full
// documented item response has seller_custom_field and attributes and nothing
// else that carries a SKU (developers.mercadolivre.com.br/pt_br/variacoes), the
// PUT that writes a SKU writes it into attributes, and none of the 34 real
// payloads from the first live drive contained the key. Reading a field the
// channel never sends is not a missing feature — it produced 34 rows claiming
// "ml omitted seller_sku" while the payload stored beside each one carried the
// SKU twice.
//
// The attribute wins over seller_custom_field because ML says the SKU "deve ser
// carregado no atributo SELLER_SKU, e não em campos personalizados do vendedor"
// (developers.mercadolivre.com.br/pt_br/publicacao-de-produtos): the custom
// field is the superseded location, so preferring it would let a stale value
// shadow the one the seller maintains.
func ItemSellerSKU(item Item) string {
	if sku := attributeValue(item.Attributes, "SELLER_SKU"); sku != "" {
		return sku
	}
	return strings.TrimSpace(item.SellerCustomField)
}

// VariationSellerSKU resolves a variation's SKU by ML's published priority:
// the variation's SELLER_SKU attribute, its seller_custom_field, then the
// item's two sources in the same order (mercado-envios-2, "hierarquia de
// prioridade para SKUs"). The fallback to the item is what makes a
// single-variation listing whose SKU is recorded once, at item level, resolve
// for the variation too.
func VariationSellerSKU(item Item, v Variation) string {
	if sku := attributeValue(v.Attributes, "SELLER_SKU"); sku != "" {
		return sku
	}
	if sku := strings.TrimSpace(v.SellerCustomField); sku != "" {
		return sku
	}
	return ItemSellerSKU(item)
}

// attributeValue returns the first nonblank value_name among the wanted
// attribute ids, in the order given — the ids are a priority list, not a set.
func attributeValue(attrs []Attribute, want ...string) string {
	for _, id := range want {
		for _, a := range attrs {
			if a.ID == id && strings.TrimSpace(a.ValueName) != "" {
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
