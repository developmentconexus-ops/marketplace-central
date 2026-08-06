package mercadolivre

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// itemsMultigetBatchSize is ML's documented multiget page size for
// GET /items?ids= (fact #13, research/external-ml-api-facts.md: "Items scan:
// 1000/batch com scroll_id; multiget 20 — verified live D-120 scan-paging +
// doc"). ids beyond this count must be partitioned into multiple calls.
const itemsMultigetBatchSize = 20

// itemMultigetRawCap bounds the PER-ITEM Raw capture. It is DELIBERATELY
// SMALLER than the adapter's whole-response cap (doRawWithHeaders,
// io.LimitReader(resp.Body, 1<<20), capability_adapter.go:747): that outer
// cap bounds the WHOLE multiget response (shared across every item in the
// batch, up to 20 of them) and is hit FIRST — if a single item's body alone
// exceeded 1MB, the whole response would already be over the outer 1MB
// ceiling and get cut mid-JSON before this reader ever sees a parseable
// array (which is correctly a BATCH-level decode failure, not a per-item
// one — see ErrItemsMultigetBatchDecode). For the per-item truncation
// marker to ever be reachable at all, this inner cap must sit strictly
// below the outer one; 256KiB is generous headroom over any real ML item
// payload (a few KB) while leaving room for a batch of up to 20 such items
// to still fit under the outer 1MB ceiling.
const itemMultigetRawCap = 256 * 1024

// ErrItemsMultigetBatchDecode marks a BATCH-level failure: the multiget
// response body did not parse as the documented `[]{code,body}` shape (e.g.
// truncated/garbled JSON, or an element count that does not match the
// requested id count for that batch). This is distinct from a PER-ITEM
// error — a batch-level failure means we cannot safely attribute any element
// to any requested id, so the whole batch call fails rather than guessing.
var ErrItemsMultigetBatchDecode = errors.New("items multiget: batch response undecodable")

// ItemMultigetDTO is the per-item multiget result. Following ADR-025 (raw
// seletivo): fields actually consumed by callers are typed, and the exact
// per-item response bytes are ALSO captured in Raw for anything not yet
// typed. Persistence of Raw (if any) is a decision for the ingest milestones
// (M-04) — this reader only guarantees capture, never decides storage.
//
// When Err is non-nil the item failed (either the ML multiget marked it with
// a non-200 `code`, or its `body` did not decode into the expected shape);
// every other field is the zero value except ProviderItemID, which always
// reflects the originally-requested id for that slot so a caller can always
// tell WHICH id failed.
type ItemMultigetDTO struct {
	ProviderItemID string
	Title          string
	Status         string
	// SellerSKU/SellerCustomField/Attributes are the item's OWN top-level
	// fields (as opposed to per-variation ones on ItemMultigetVariationDTO) —
	// the same wire fields mlItemResponse (capability_adapter.go, single-item
	// ReadListing) already types, since ML's multiget `body` element is the
	// same full item representation as the single-item GET. Needed so a
	// no-variation item can still derive a listing-level SellerSKU/EAN
	// (multiget_mapper.go's MapMultigetItemToListingSnapshot) instead of
	// leaving that ADR-019 seam honest-empty for lack of a typed field.
	SellerSKU         string
	SellerCustomField string
	Attributes        []ItemMultigetAttributeDTO
	// Price é a forma decimal em string (mesma convenção de
	// ItemMultigetVariationDTO.Price) — nunca float.
	Price         *string
	CurrencyID    string
	ListingTypeID string
	// InitialQuantity é a quantidade PUBLICADA. Não é
	// AvailableQuantity + SoldQuantity: reposição de estoque incrementa a
	// disponível sem tocar a inicial. São dois números distintos.
	InitialQuantity   *int
	SoldQuantity      *int
	CategoryID        string
	Condition         string
	Permalink         string
	Thumbnail         string
	DateCreated       *time.Time
	Tags              []string
	CatalogProductID  string
	Shipping          *ItemMultigetShippingDTO
	AvailableQuantity *int
	Variations        []ItemMultigetVariationDTO
	Raw               json.RawMessage
	RawTruncated      bool
	Err               error
}

// ItemMultigetShippingDTO is the "whatever shape is available" shipping
// sub-object the feature brief asks for (IC-07 E3) — deliberately narrow
// (mode/logistic_type/free_shipping) since the FULL shipping shape already
// has a dedicated reader (shipping_reader.go) for when a shipment_id exists;
// this is only the summary ML embeds inline on the item itself.
type ItemMultigetShippingDTO struct {
	Mode         string
	LogisticType string
	FreeShipping *bool
}

// ItemMultigetVariationDTO covers the per-variation E3 fields named in the
// feature brief: price/available_quantity/sold_quantity/seller_sku/attributes.
type ItemMultigetVariationDTO struct {
	ProviderVariationID string
	// Price is the decimal string form (mirrors decimalString/providerMoney's
	// convention elsewhere in this package) — NOT a float, to avoid
	// binary-float precision drift on money-shaped values.
	Price             *string
	AvailableQuantity *int
	SoldQuantity      *int
	SellerSKU         string
	Attributes        []ItemMultigetAttributeDTO
}

// ItemMultigetAttributeDTO mirrors mlAttribute (capability_adapter.go) as a
// package-exported-shape DTO field; kept as its own type here (rather than
// reusing mlAttribute directly on the DTO) so the DTO surface stays decoupled
// from the wire-decode struct's field tags.
type ItemMultigetAttributeDTO struct {
	ID        string
	ValueID   string
	ValueName string
}

// mlMultigetElement is the ML multiget response element shape: an array of
// `{code, body}` pairs, one per requested id, IN REQUEST ORDER per batch
// (feature brief F-02; fact #13). code is the per-item HTTP-equivalent status
// (200 on success); body is kept as json.RawMessage so a per-item decode
// failure never blocks parsing the OTHER items in the same batch.
type mlMultigetElement struct {
	Code int             `json:"code"`
	Body json.RawMessage `json:"body"`
}

// mlMultigetItemBody is the per-item body shape for a successful (code=200)
// multiget element, covering the E3 fields named in the feature brief
// (IC-07). Unknown fields are ignored (Go's default json.Unmarshal
// behavior) — Raw carries the rest.
type mlMultigetItemBody struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// SellerSKU/SellerCustomField/Attributes: the item's OWN top-level fields
	// (distinct from mlMultigetVariation's per-variation fields of the same
	// name below). mlItemResponse (capability_adapter.go) already types these
	// three at top level for the single-item GET; this struct originally
	// omitted them (json.Unmarshal silently drops undeclared fields), which
	// left every no-variation item's listing-level SellerSKU/EAN honest-empty
	// with no way to derive them even though the bytes are present in Raw.
	SellerSKU         string        `json:"seller_sku"`
	SellerCustomField string        `json:"seller_custom_field"`
	Attributes        []mlAttribute `json:"attributes"`
	Status            string        `json:"status"`
	// Price/CurrencyID/ListingTypeID/InitialQuantity: o ML manda os quatro em
	// TODO item, e este struct não os declarava — encoding/json os descartava
	// em silêncio, deixando price_amount/price_currency/listing_type_code/
	// published_quantity NULL em 34 de 34 anúncios. Price é json.Number (e não
	// float64) pela mesma convenção de dinheiro do resto do pacote: evita
	// deriva de precisão binária em valor monetário.
	Price             *json.Number          `json:"price"`
	CurrencyID        string                `json:"currency_id"`
	ListingTypeID     string                `json:"listing_type_id"`
	InitialQuantity   *int                  `json:"initial_quantity"`
	SoldQuantity      *int                  `json:"sold_quantity"`
	CategoryID        string                `json:"category_id"`
	Condition         string                `json:"condition"`
	Permalink         string                `json:"permalink"`
	Thumbnail         string                `json:"thumbnail"`
	DateCreated       string                `json:"date_created"`
	Tags              []string              `json:"tags"`
	CatalogProductID  string                `json:"catalog_product_id"`
	Shipping          *mlMultigetShipping   `json:"shipping"`
	AvailableQuantity *int                  `json:"available_quantity"`
	Variations        []mlMultigetVariation `json:"variations"`
}

type mlMultigetShipping struct {
	Mode         string `json:"mode"`
	LogisticType string `json:"logistic_type"`
	FreeShipping *bool  `json:"free_shipping"`
}

// mlMultigetVariation mirrors mlVariation (capability_adapter.go) but adds
// SoldQuantity, which the single-item ReadListing path does not currently
// consume but IC-07 requires here.
type mlMultigetVariation struct {
	ID                any           `json:"id"`
	Price             *json.Number  `json:"price"`
	AvailableQuantity *int          `json:"available_quantity"`
	SoldQuantity      *int          `json:"sold_quantity"`
	SellerSKU         string        `json:"seller_sku"`
	SellerCustomField string        `json:"seller_custom_field"`
	Attributes        []mlAttribute `json:"attributes"`
}

// GetItemsMultiget resolves item DTOs for a set of provider item ids via
// ML's GET /items?ids= multiget, batching at 20 ids per call (fact #13) and
// returning DTOs in the SAME order as providerItemIDs regardless of how many
// batches were needed. It never fails the whole call for a single bad id —
// see ItemMultigetDTO.Err. All HTTP traffic goes through the adapter's
// existing choke point (doJSON -> doRawWithHeaders), so it inherits F-01's
// resilience decorator automatically; this reader implements no retry of its
// own (feature brief negative scenario: batch 429 propagates as-is).
func (a *CapabilityAdapter) GetItemsMultiget(ctx context.Context, accountRef domain.ProviderAccountRef, providerItemIDs []string) ([]ItemMultigetDTO, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return nil, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, mapPricingReaderError(err)
	}
	return a.getItemsMultiget(ctx, accountRef, token, providerItemIDs)
}

func (a *CapabilityAdapter) getItemsMultiget(ctx context.Context, accountRef domain.ProviderAccountRef, token string, providerItemIDs []string) ([]ItemMultigetDTO, error) {
	if len(providerItemIDs) == 0 {
		return nil, nil
	}

	results := make([]ItemMultigetDTO, 0, len(providerItemIDs))
	for start := 0; start < len(providerItemIDs); start += itemsMultigetBatchSize {
		end := start + itemsMultigetBatchSize
		if end > len(providerItemIDs) {
			end = len(providerItemIDs)
		}
		batchIDs := providerItemIDs[start:end]
		batchDTOs, err := a.getItemsMultigetBatch(ctx, accountRef, token, batchIDs)
		if err != nil {
			return nil, err
		}
		results = append(results, batchDTOs...)
	}
	return results, nil
}

func (a *CapabilityAdapter) getItemsMultigetBatch(ctx context.Context, accountRef domain.ProviderAccountRef, token string, batchIDs []string) ([]ItemMultigetDTO, error) {
	query := url.Values{}
	query.Set("ids", strings.Join(batchIDs, ","))
	path := "/items?" + query.Encode()

	var elements []mlMultigetElement
	// doJSON is the SAME choke point (doRawWithHeaders underneath) every other
	// reader in this package uses; a whole-batch 429/5xx/401 is classified and
	// returned here with zero retry — that is F-01's decorator's job, wrapping
	// doRawWithHeaders, not this reader's.
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &elements); err != nil {
		if domain.ErrorCodeOf(err) == domain.ErrCodeProviderPayloadInvalid {
			return nil, fmt.Errorf("%w: %w", ErrItemsMultigetBatchDecode, mapPricingReaderError(err))
		}
		return nil, mapPricingReaderError(err)
	}
	if len(elements) != len(batchIDs) {
		return nil, fmt.Errorf("%w: requested %d ids, got %d elements", ErrItemsMultigetBatchDecode, len(batchIDs), len(elements))
	}

	dtos := make([]ItemMultigetDTO, len(batchIDs))
	for i, element := range elements {
		dtos[i] = mapMultigetElement(batchIDs[i], element)
	}
	return dtos, nil
}

func mapMultigetElement(requestedID string, element mlMultigetElement) ItemMultigetDTO {
	if element.Code != http.StatusOK {
		return ItemMultigetDTO{
			ProviderItemID: requestedID,
			Err:            mapMultigetItemStatus(element.Code, element.Body),
		}
	}

	var body mlMultigetItemBody
	if err := json.Unmarshal(element.Body, &body); err != nil {
		return ItemMultigetDTO{
			ProviderItemID: requestedID,
			// Named per-item error (negative scenario: "Body de item não-JSON"):
			// carries the diagnosable route+body, never sinks the batch.
			Err: mapPricingReaderError(domain.NewCapabilityError(
				domain.ErrCodeProviderPayloadInvalid,
				fmt.Sprintf("items multiget: item %s body decode failed: %s", requestedID, providerDiag(http.MethodGet, "/items", http.StatusOK, element.Body)),
			)),
		}
	}

	dto := ItemMultigetDTO{
		ProviderItemID:    firstNonEmpty(strings.TrimSpace(body.ID), requestedID),
		Title:             body.Title,
		Status:            body.Status,
		Price:             multigetPriceString(body.Price),
		CurrencyID:        strings.TrimSpace(body.CurrencyID),
		ListingTypeID:     strings.TrimSpace(body.ListingTypeID),
		InitialQuantity:   body.InitialQuantity,
		SellerSKU:         body.SellerSKU,
		SellerCustomField: body.SellerCustomField,
		Attributes:        mapMultigetAttributes(body.Attributes),
		SoldQuantity:      body.SoldQuantity,
		CategoryID:        body.CategoryID,
		Condition:         body.Condition,
		Permalink:         body.Permalink,
		Thumbnail:         body.Thumbnail,
		DateCreated:       parseTimePtr(body.DateCreated),
		Tags:              body.Tags,
		CatalogProductID:  body.CatalogProductID,
		AvailableQuantity: body.AvailableQuantity,
	}
	if body.Shipping != nil {
		dto.Shipping = &ItemMultigetShippingDTO{
			Mode:         body.Shipping.Mode,
			LogisticType: body.Shipping.LogisticType,
			FreeShipping: body.Shipping.FreeShipping,
		}
	}
	for _, variation := range body.Variations {
		dto.Variations = append(dto.Variations, ItemMultigetVariationDTO{
			ProviderVariationID: normalizeAnyID(variation.ID),
			Price:               decimalString(variation.Price),
			AvailableQuantity:   variation.AvailableQuantity,
			SoldQuantity:        variation.SoldQuantity,
			SellerSKU:           firstNonEmpty(variation.SellerSKU, variation.SellerCustomField),
			Attributes:          mapMultigetAttributes(variation.Attributes),
		})
	}

	raw, truncated := capRawMessage(element.Body)
	dto.Raw = raw
	dto.RawTruncated = truncated
	return dto
}

func mapMultigetAttributes(attributes []mlAttribute) []ItemMultigetAttributeDTO {
	if len(attributes) == 0 {
		return nil
	}
	mapped := make([]ItemMultigetAttributeDTO, 0, len(attributes))
	for _, attribute := range attributes {
		mapped = append(mapped, ItemMultigetAttributeDTO{
			ID:        attribute.ID,
			ValueID:   attribute.ValueID,
			ValueName: attribute.ValueName,
		})
	}
	return mapped
}

// mapMultigetItemStatus maps the multiget element's embedded `code` (an
// HTTP-equivalent status ML puts INSIDE the array element, not the outer HTTP
// response) onto the same domain error taxonomy doJSONWithHeaders uses for
// the outer response status (capability_adapter.go:643-675) — kept as its own
// function here (rather than calling into capability_adapter.go, which is
// off-limits for this feature) so a per-item 404/401/429/5xx degrades to the
// SAME sentinel a caller would see from a single-item read.
func mapMultigetItemStatus(code int, body json.RawMessage) error {
	var capErr error
	switch {
	case code == http.StatusUnauthorized || code == http.StatusForbidden:
		capErr = domain.NewCapabilityError(domain.ErrCodeProviderAuth, "provider credentials were rejected")
	case code == http.StatusNotFound:
		capErr = domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	case code == http.StatusTooManyRequests:
		capErr = domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, "provider rate limit reached")
	case code >= 500:
		capErr = domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	default:
		capErr = domain.NewCapabilityError(domain.ErrCodeProviderValidation, providerDiag(http.MethodGet, "/items", code, body))
	}
	return mapPricingReaderError(capErr)
}

// capRawMessage bounds a per-item raw body to itemMultigetRawCap (256KiB),
// mirroring the adapter's whole-response LimitReader convention
// (capability_adapter.go:747) but applied per item rather than per response.
// Returns a copy (never aliases the caller's backing array) and an explicit
// truncated flag — an oversize body is NEVER silently prefixed without the
// flag (ADR-025).
func capRawMessage(raw json.RawMessage) (json.RawMessage, bool) {
	if len(raw) <= itemMultigetRawCap {
		out := make(json.RawMessage, len(raw))
		copy(out, raw)
		return out, false
	}
	out := make(json.RawMessage, itemMultigetRawCap)
	copy(out, raw[:itemMultigetRawCap])
	return out, true
}

// multigetPriceString converte o número de fio para a forma decimal em string.
// json.Number já preserva os dígitos exatos recebidos; converter para float
// aqui reintroduziria a deriva binária que a convenção de dinheiro do pacote
// existe para evitar.
func multigetPriceString(n *json.Number) *string {
	if n == nil {
		return nil
	}
	s := strings.TrimSpace(n.String())
	if s == "" {
		return nil
	}
	return &s
}
