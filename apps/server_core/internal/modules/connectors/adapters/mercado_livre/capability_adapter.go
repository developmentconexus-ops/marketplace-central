package mercadolivre

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

const (
	defaultBaseURL = "https://api.mercadolibre.com"
	defaultSiteID  = "MLB"
)

type AccessTokenResolver func(context.Context, domain.ProviderAccountRef) (string, error)

type CapabilityAdapterConfig struct {
	BaseURL              string
	SiteID               string
	HTTPClient           *http.Client
	AccessTokenResolver  AccessTokenResolver
	Now                  func() time.Time
	CatalogOffersEnabled bool
}

type CapabilityAdapter struct {
	baseURL              string
	siteID               string
	httpClient           *http.Client
	accessTokenResolver  AccessTokenResolver
	now                  func() time.Time
	catalogOffersEnabled bool
}

var _ ports.MarketReader = (*CapabilityAdapter)(nil)
var _ ports.CatalogOffersReader = (*CapabilityAdapter)(nil)
var _ ports.CatalogReader = (*CapabilityAdapter)(nil)
var _ ports.ShipmentReader = (*CapabilityAdapter)(nil)

func NewCapabilityAdapter(cfg CapabilityAdapterConfig) *CapabilityAdapter {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	siteID := strings.TrimSpace(cfg.SiteID)
	if siteID == "" {
		siteID = defaultSiteID
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}

	return &CapabilityAdapter{
		baseURL:              baseURL,
		siteID:               siteID,
		httpClient:           httpClient,
		accessTokenResolver:  cfg.AccessTokenResolver,
		now:                  now,
		catalogOffersEnabled: cfg.CatalogOffersEnabled,
	}
}

func (a *CapabilityAdapter) ProviderCapabilitySet() connectorsapp.ProviderCapabilitySet {
	return connectorsapp.ProviderCapabilitySet{
		ProviderCode:  "mercado_livre",
		AccountProbes: a,
		Listings:      a,
		FeeQuotes:     a,
		StockReads:    a,
		Orders:        a,
		PriceWrites:   a,
		StockWrites:   a,
		ListingWrites: a,
	}
}

func (a *CapabilityAdapter) GetOwnItemPricing(ctx context.Context, accountRef domain.ProviderAccountRef, itemID string) (domain.OwnItemPricing, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.OwnItemPricing{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.OwnItemPricing{}, mapPricingReaderError(err)
	}
	return a.getOwnItemPricing(ctx, accountRef, token, strings.TrimSpace(itemID))
}

func (a *CapabilityAdapter) GetPriceToWin(ctx context.Context, accountRef domain.ProviderAccountRef, itemID string) (domain.PriceToWin, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.PriceToWin{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.PriceToWin{}, mapPricingReaderError(err)
	}
	return a.getPriceToWin(ctx, accountRef, token, strings.TrimSpace(itemID))
}

func (a *CapabilityAdapter) SearchCatalogByEAN(ctx context.Context, accountRef domain.ProviderAccountRef, ean string) (domain.CatalogSearchResult, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.CatalogSearchResult{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.CatalogSearchResult{}, mapPricingReaderError(err)
	}
	return a.searchCatalogByEAN(ctx, accountRef, token, strings.TrimSpace(ean))
}

func (a *CapabilityAdapter) GetCatalogProduct(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) (domain.CatalogProduct, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.CatalogProduct{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.CatalogProduct{}, mapPricingReaderError(err)
	}
	return a.getCatalogProduct(ctx, accountRef, token, strings.TrimSpace(catalogProductID))
}

func (a *CapabilityAdapter) GetShipmentInfo(ctx context.Context, accountRef domain.ProviderAccountRef, shipmentID string) (domain.ShipmentInfo, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.ShipmentInfo{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.ShipmentInfo{}, mapPricingReaderError(err)
	}
	return a.getShipmentInfo(ctx, accountRef, token, strings.TrimSpace(shipmentID))
}

func (a *CapabilityAdapter) GetFreeShippingCost(ctx context.Context, accountRef domain.ProviderAccountRef, query domain.FreeShippingQuery) (domain.FreeShippingCost, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.FreeShippingCost{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.FreeShippingCost{}, mapPricingReaderError(err)
	}
	return a.getFreeShippingCost(ctx, accountRef, token, query)
}

func (a *CapabilityAdapter) ListCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) ([]domain.CatalogOffer, error) {
	started := time.Now()
	status := "error"
	pageCount := 0
	offerCount := 0
	defer func() { logCatalogOffers(status, pageCount, offerCount, started) }()
	if !a.catalogOffersEnabled {
		status = "unavailable"
		return nil, domain.ErrCatalogOffersUnavailable
	}
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return nil, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, mapPricingReaderError(err)
	}
	offers, pages, err := a.listCatalogOffers(ctx, accountRef, token, strings.TrimSpace(catalogProductID))
	pageCount = pages
	if err != nil {
		if errors.Is(err, domain.ErrCatalogOffersUnavailable) {
			status = "unavailable"
		}
		return nil, err
	}
	offerCount = len(offers)
	status = "ok"
	return offers, nil
}

func (a *CapabilityAdapter) ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error) {
	accountRef, err := normalizeAccountRef(ref)
	if err != nil {
		return domain.AccountSnapshot{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.AccountSnapshot{}, err
	}

	var response mlAccountResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/me", nil, &response); err != nil {
		return domain.AccountSnapshot{}, err
	}

	return domain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   firstNonEmpty(normalizeAnyID(response.ID), accountRef.ProviderAccountID),
		ProviderAccountName: firstNonEmpty(response.Nickname, strings.TrimSpace(response.FirstName+" "+response.LastName)),
		SiteID:              strings.TrimSpace(response.SiteID),
		Status:              firstNonEmpty(response.Status.SiteStatus, response.Status.List.Status),
		FetchedAt:           a.now(),
		RawProviderRef:      "/users/me",
	}, nil
}

func (a *CapabilityAdapter) ListListings(ctx context.Context, input domain.ListListingsInput) ([]domain.ListingSnapshot, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return nil, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("limit", strconv.Itoa(limitOrDefault(input.Limit, 50)))
	query.Set("offset", strconv.Itoa(offsetFromCursor(input.Cursor)))
	query.Set("include_filters", "false")

	var response struct {
		Results []string `json:"results"`
	}
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/"+url.PathEscape(accountRef.ProviderAccountID)+"/items/search?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}

	snapshots := make([]domain.ListingSnapshot, 0, len(response.Results))
	for _, itemID := range response.Results {
		snapshot, err := a.ReadListing(ctx, domain.ProviderListingRef{
			AccountRef:     accountRef,
			ProviderItemID: itemID,
		})
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// ListListingsScanPage pages the seller's items with ML's search_type=scan,
// which has no depth cap — plain offset paging on /users/{id}/items/search is
// rejected by the provider past 1000 items, so accounts above that size can
// only be fully ingested via scan. The returned cursor is the provider
// scroll_id to pass back on the next call; the first call passes an empty
// cursor. An empty result page ends the scan.
func (a *CapabilityAdapter) ListListingsScanPage(ctx context.Context, input domain.ListListingsInput) ([]domain.ListingSnapshot, string, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return nil, "", err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, "", err
	}

	query := url.Values{}
	query.Set("search_type", "scan")
	query.Set("limit", strconv.Itoa(limitOrDefault(input.Limit, 50)))
	if cursor := strings.TrimSpace(input.Cursor); cursor != "" {
		query.Set("scroll_id", cursor)
	}

	var response struct {
		ScrollID string   `json:"scroll_id"`
		Results  []string `json:"results"`
	}
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/"+url.PathEscape(accountRef.ProviderAccountID)+"/items/search?"+query.Encode(), nil, &response); err != nil {
		return nil, "", err
	}

	snapshots := make([]domain.ListingSnapshot, 0, len(response.Results))
	for _, itemID := range response.Results {
		snapshot, err := a.ReadListing(ctx, domain.ProviderListingRef{
			AccountRef:     accountRef,
			ProviderItemID: itemID,
		})
		if err != nil {
			return nil, "", err
		}
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, strings.TrimSpace(response.ScrollID), nil
}

func (a *CapabilityAdapter) ReadListing(ctx context.Context, ref domain.ProviderListingRef) (domain.ListingSnapshot, error) {
	accountRef, err := normalizeAccountRef(ref.AccountRef)
	if err != nil {
		return domain.ListingSnapshot{}, err
	}
	itemID := strings.TrimSpace(ref.ProviderItemID)
	if itemID == "" {
		return domain.ListingSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider item id is required")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.ListingSnapshot{}, err
	}

	item, err := a.readItem(ctx, accountRef, token, itemID)
	if err != nil {
		return domain.ListingSnapshot{}, err
	}
	return a.mapListing(item), nil
}

func (a *CapabilityAdapter) ReadStock(ctx context.Context, ref domain.ProviderListingRef) (domain.StockSnapshot, error) {
	accountRef, err := normalizeAccountRef(ref.AccountRef)
	if err != nil {
		return domain.StockSnapshot{}, err
	}
	itemID := strings.TrimSpace(ref.ProviderItemID)
	if itemID == "" {
		return domain.StockSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider item id is required")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.StockSnapshot{}, err
	}

	item, err := a.readItem(ctx, accountRef, token, itemID)
	if err != nil {
		return domain.StockSnapshot{}, err
	}

	variationID := strings.TrimSpace(ref.ProviderVariationID)
	if variationID == "" {
		if len(item.Variations) > 0 {
			return domain.StockSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderUnsupportedShape, "item requires variation reference for stock reads")
		}
		quantity := intValue(item.AvailableQuantity)
		return domain.StockSnapshot{
			ProviderCode:      "mercado_livre",
			ProviderItemID:    item.ID,
			AvailableQuantity: quantity,
			ProviderStatus:    item.Status,
			SellerSKU:         firstNonEmpty(item.SellerSKU, item.SellerCustomField),
			EAN:               attributeValue(item.Attributes, "GTIN", "EAN"),
			Title:             item.Title,
			SourceUpdatedAt:   parseTimePtr(item.LastUpdated),
			FetchedAt:         a.now(),
			RawProviderRef:    "/items/" + item.ID,
			Scope:             domain.StockScopeItem,
		}, nil
	}

	variation, found := findVariation(item.Variations, variationID)
	if !found {
		return domain.StockSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderUnsupportedShape, "variation not found for stock read")
	}

	return domain.StockSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderItemID:      item.ID,
		ProviderVariationID: normalizeAnyID(variation.ID),
		AvailableQuantity:   intValue(variation.AvailableQuantity),
		ProviderStatus:      item.Status,
		SellerSKU:           firstNonEmpty(variation.SellerSKU, variation.SellerCustomField, item.SellerSKU, item.SellerCustomField),
		EAN:                 attributeValue(variation.Attributes, "GTIN", "EAN"),
		Title:               item.Title,
		SourceUpdatedAt:     parseTimePtr(item.LastUpdated),
		FetchedAt:           a.now(),
		RawProviderRef:      "/items/" + item.ID,
		Scope:               domain.StockScopeVariation,
	}, nil
}

func (a *CapabilityAdapter) UpdateAvailableQuantity(ctx context.Context, request domain.StockWriteRequest) (domain.StockWriteResult, error) {
	accountRef, err := normalizeAccountRef(domain.ProviderAccountRef{
		TenantID:          request.TenantID,
		InstallationID:    request.InstallationID,
		ProviderCode:      request.ProviderCode,
		ProviderAccountID: request.ProviderAccountID,
	})
	if err != nil {
		return domain.StockWriteResult{}, err
	}
	itemID := strings.TrimSpace(request.ProviderItemID)
	if itemID == "" {
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider item id is required")
	}
	if strings.TrimSpace(request.IdempotencyKey) == "" {
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "idempotency key is required")
	}
	if request.RequestedQuantity < 0 {
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "requested quantity must be non-negative")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.StockWriteResult{}, err
	}

	item, err := a.readItem(ctx, accountRef, token, itemID)
	if err != nil {
		return domain.StockWriteResult{}, err
	}

	body := make(map[string]any)
	variationID := strings.TrimSpace(request.ProviderVariationID)
	if variationID == "" {
		if len(item.Variations) > 0 {
			return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderUnsupportedShape, "item requires variation reference for stock writes")
		}
		body["available_quantity"] = request.RequestedQuantity
	} else {
		if len(item.Variations) == 0 {
			return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderUnsupportedShape, "item does not expose variation stock writes")
		}
		variation, found := findVariation(item.Variations, variationID)
		if !found {
			return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderUnsupportedShape, "variation not found for stock write")
		}
		body["variations"] = []map[string]any{{
			"id":                 variation.ID,
			"available_quantity": request.RequestedQuantity,
		}}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, "failed to serialize stock write payload")
	}

	resp, respBody, err := a.doRaw(ctx, accountRef, token, http.MethodPut, "/items/"+url.PathEscape(item.ID), bytes.NewReader(payload))
	if err != nil {
		return domain.StockWriteResult{}, err
	}

	result := domain.StockWriteResult{
		IDempotencyKey:      request.IdempotencyKey,
		ProviderCode:        "mercado_livre",
		ProviderItemID:      item.ID,
		ProviderVariationID: variationID,
		ProviderResponseRef: "/items/" + item.ID,
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		result.Result = domain.StockWriteResultApplied
		result.Message = "provider write applied"
		return result, nil
	case resp.StatusCode == http.StatusTooManyRequests:
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, "provider rate limit while updating stock")
	case resp.StatusCode >= 500:
		return domain.StockWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider transient failure while updating stock")
	default:
		result.Result = domain.StockWriteResultRejected
		result.Message = strings.TrimSpace(string(respBody))
		if result.Message == "" {
			result.Message = resp.Status
		}
		return result, nil
	}
}

func (a *CapabilityAdapter) ListOrders(ctx context.Context, input domain.ListOrdersInput) ([]domain.OrderSnapshot, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return nil, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("seller", accountRef.ProviderAccountID)
	query.Set("limit", strconv.Itoa(limitOrDefault(input.Limit, 50)))
	query.Set("offset", strconv.Itoa(offsetFromCursor(input.Cursor)))
	if status := strings.TrimSpace(input.Status); status != "" {
		query.Set("order.status", status)
	}
	if input.UpdatedAfter != nil {
		query.Set("order.date_last_updated.from", input.UpdatedAfter.UTC().Format(time.RFC3339))
	}

	var response struct {
		Results []struct {
			ID any `json:"id"`
		} `json:"results"`
	}
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/orders/search?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}

	snapshots := make([]domain.OrderSnapshot, 0, len(response.Results))
	for _, result := range response.Results {
		orderID := normalizeAnyID(result.ID)
		if orderID == "" {
			continue
		}
		snapshot, err := a.ReadOrder(ctx, domain.ProviderOrderRef{
			AccountRef:      accountRef,
			ProviderOrderID: orderID,
		})
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

func (a *CapabilityAdapter) ReadOrder(ctx context.Context, ref domain.ProviderOrderRef) (domain.OrderSnapshot, error) {
	accountRef, err := normalizeAccountRef(ref.AccountRef)
	if err != nil {
		return domain.OrderSnapshot{}, err
	}
	orderID := strings.TrimSpace(ref.ProviderOrderID)
	if orderID == "" {
		return domain.OrderSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider order id is required")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.OrderSnapshot{}, err
	}

	var order mlOrderResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/orders/"+url.PathEscape(orderID), nil, &order); err != nil {
		return domain.OrderSnapshot{}, err
	}
	return a.mapOrder(order), nil
}

func (a *CapabilityAdapter) ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	if input.PriceAmount <= 0 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price amount must be positive")
	}
	listingTypeID := strings.TrimSpace(input.ListingTypeID)
	if listingTypeID == "" {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "listing type id is required")
	}
	siteID := firstNonEmpty(input.SiteID, a.siteID)

	query := url.Values{}
	query.Set("price", strconv.FormatFloat(input.PriceAmount, 'f', -1, 64))
	query.Set("listing_type_id", listingTypeID)
	if categoryID := strings.TrimSpace(input.CategoryID); categoryID != "" {
		query.Set("category_id", categoryID)
	}

	resp, rawBody, err := a.doRaw(ctx, accountRef, token, http.MethodGet, "/sites/"+url.PathEscape(siteID)+"/listing_prices?"+query.Encode(), nil)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, "provider rate limit reached")
	}
	if resp.StatusCode >= 500 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	}
	if resp.StatusCode >= 400 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, strings.TrimSpace(string(rawBody)))
	}

	response, err := decodeListingPriceResponses(rawBody)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	if len(response) == 0 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider returned no fee quote")
	}

	price := response[0]
	return domain.FeeQuoteSnapshot{
		ProviderCode:      "mercado_livre",
		SiteID:            siteID,
		CategoryID:        strings.TrimSpace(input.CategoryID),
		ListingTypeID:     firstNonEmpty(price.ListingTypeID, listingTypeID),
		PriceAmount:       input.PriceAmount,
		CurrencyID:        firstNonEmpty(strings.TrimSpace(input.CurrencyID), "BRL"),
		CommissionPercent: price.SaleFeeDetails.PercentageFee,
		FixedFeeAmount:    price.SaleFeeDetails.FixedFee,
		FetchedAt:         a.now(),
		RawProviderRef:    "/sites/" + siteID + "/listing_prices",
	}, nil
}

func (a *CapabilityAdapter) accessToken(ctx context.Context, accountRef domain.ProviderAccountRef) (string, error) {
	if a.accessTokenResolver == nil {
		return "", domain.NewCapabilityError(domain.ErrCodeProviderAuth, "access token resolver is not configured")
	}
	token, err := a.accessTokenResolver(ctx, accountRef)
	if err != nil {
		return "", domain.NewCapabilityError(domain.ErrCodeProviderAuth, "failed to resolve access token")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", domain.NewCapabilityError(domain.ErrCodeProviderAuth, "access token is empty")
	}
	return token, nil
}

func (a *CapabilityAdapter) readItem(ctx context.Context, accountRef domain.ProviderAccountRef, token string, itemID string) (mlItemResponse, error) {
	var item mlItemResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/items/"+url.PathEscape(itemID), nil, &item); err != nil {
		return mlItemResponse{}, err
	}
	return item, nil
}

func (a *CapabilityAdapter) doJSON(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, body io.Reader, target any) error {
	return a.doJSONWithHeaders(ctx, accountRef, token, method, path, body, nil, target)
}

// doJSONWithHeaders is doJSON plus a narrow per-request header map, for the few
// ML endpoints that require an extra header (e.g. the shipments /costs endpoint
// needs `x-format-new: true` or ML returns the legacy cost shape that fails to
// decode). Passing nil headers is identical to doJSON — it does NOT alter the
// default headers of any other call.
func (a *CapabilityAdapter) doJSONWithHeaders(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, body io.Reader, headers map[string]string, target any) error {
	resp, rawBody, err := a.doRawWithHeaders(ctx, accountRef, token, method, path, body, "", headers)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.NewCapabilityError(domain.ErrCodeProviderAuth, "provider credentials were rejected")
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, "provider rate limit reached")
	}
	if resp.StatusCode >= 500 {
		return domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	}
	if resp.StatusCode >= 400 {
		return domain.NewCapabilityError(domain.ErrCodeProviderValidation, providerDiag(method, path, resp.StatusCode, rawBody))
	}
	if target == nil {
		return nil
	}
	if err := json.Unmarshal(rawBody, target); err != nil {
		// A 2xx whose body does not fit the expected struct is a contract drift,
		// not a client bug. Carry the raw (bounded) body + route so the failure is
		// diagnosable instead of a blind "decode failed" (FINDING-M02-LIVE-2 /
		// live-drive reprova, D-86). The market-collection boundary sanitizes
		// (token-redacts + truncates) this message before it reaches any log.
		return domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, providerDiag(method, path, resp.StatusCode, rawBody))
	}
	return nil
}

// providerDiag builds a bounded, single-line diagnostic string from a provider
// response: HTTP route + status + the raw body clipped to a sane length. Token
// redaction is applied downstream at the market-collection boundary; this keeps
// only the size bound so a large body cannot bloat the error chain.
func providerDiag(method, path string, status int, rawBody []byte) string {
	body := strings.Join(strings.Fields(string(rawBody)), " ")
	const maxDiagBody = 512
	if runes := []rune(body); len(runes) > maxDiagBody {
		body = string(runes[:maxDiagBody]) + "…"
	}
	return fmt.Sprintf("%s %s -> HTTP %d: %s", method, path, status, body)
}

func decodeListingPriceResponses(rawBody []byte) ([]mlListingPriceResponse, error) {
	var batch []mlListingPriceResponse
	if err := json.Unmarshal(rawBody, &batch); err == nil {
		return batch, nil
	}

	var single mlListingPriceResponse
	if err := json.Unmarshal(rawBody, &single); err == nil {
		return []mlListingPriceResponse{single}, nil
	}

	return nil, domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, "provider payload decode failed")
}

func (a *CapabilityAdapter) doRaw(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, body io.Reader) (*http.Response, []byte, error) {
	return a.doRawWithIdempotency(ctx, accountRef, token, method, path, body, "")
}

func (a *CapabilityAdapter) doRawWithIdempotency(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, body io.Reader, idempotencyKey string) (*http.Response, []byte, error) {
	return a.doRawWithHeaders(ctx, accountRef, token, method, path, body, idempotencyKey, nil)
}

func (a *CapabilityAdapter) doRawWithHeaders(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, body io.Reader, idempotencyKey string, headers map[string]string) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, body)
	if err != nil {
		return nil, nil, domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider request build failed")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-Tenant-ID", accountRef.TenantID)
	req.Header.Set("X-Installation-ID", accountRef.InstallationID)
	if strings.TrimSpace(idempotencyKey) != "" {
		req.Header.Set("X-Idempotency-Key", strings.TrimSpace(idempotencyKey))
	}
	// Per-request headers apply last (narrow, endpoint-specific — e.g.
	// x-format-new for the shipments /costs endpoint). They never touch the
	// default headers of other calls.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		// Wrap the transport error as Cause so timeout classification
		// (errors.Is(err, context.DeadlineExceeded)) survives the chain
		// through CapabilityError.Unwrap (IC-03 Amendment A1).
		return nil, nil, &domain.CapabilityError{
			Code:    domain.ErrCodeProviderTransient,
			Message: "provider request failed",
			Cause:   err,
		}
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, nil, domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider response read failed")
	}
	return resp, rawBody, nil
}

func (a *CapabilityAdapter) mapListing(item mlItemResponse) domain.ListingSnapshot {
	priceAmount := strings.TrimSpace(item.Price.String())
	priceCurrency := strings.TrimSpace(item.CurrencyID)
	if priceAmount == "" || priceCurrency == "" {
		priceAmount = ""
		priceCurrency = ""
	}
	var publishedPrice *string
	if priceAmount != "" {
		publishedPrice = &priceAmount
	}

	snapshot := domain.ListingSnapshot{
		ProviderCode:      "mercado_livre",
		ProviderItemID:    item.ID,
		ProviderStatus:    item.Status,
		PriceAmount:       publishedPrice,
		PriceCurrency:     priceCurrency,
		ListingTypeCode:   strings.TrimSpace(item.ListingTypeID),
		SellerSKU:         firstNonEmpty(item.SellerSKU, item.SellerCustomField),
		EAN:               attributeValue(item.Attributes, "GTIN", "EAN"),
		Title:             item.Title,
		AvailableQuantity: item.AvailableQuantity,
		SourceUpdatedAt:   parseTimePtr(item.LastUpdated),
		FetchedAt:         a.now(),
		RawProviderRef:    "/items/" + item.ID,
	}

	if len(item.Variations) == 0 {
		return snapshot
	}

	snapshot.Variations = make([]domain.ListingVariationSnapshot, 0, len(item.Variations))
	for _, variation := range item.Variations {
		snapshot.Variations = append(snapshot.Variations, domain.ListingVariationSnapshot{
			ProviderVariationID: normalizeAnyID(variation.ID),
			SellerSKU:           firstNonEmpty(variation.SellerSKU, variation.SellerCustomField),
			EAN:                 attributeValue(variation.Attributes, "GTIN", "EAN"),
			AvailableQuantity:   variation.AvailableQuantity,
		})
	}

	return snapshot
}

func (a *CapabilityAdapter) mapOrder(order mlOrderResponse) domain.OrderSnapshot {
	snapshot := domain.OrderSnapshot{
		ProviderCode:         "mercado_livre",
		ProviderOrderID:      normalizeAnyID(order.ID),
		ProviderStatus:       order.Status,
		ProviderStatusDetail: order.StatusDetail,
		ProviderCreatedAt:    parseTimePtr(order.DateCreated),
		ProviderClosedAt:     parseTimePtr(order.DateClosed),
		ProviderUpdatedAt:    parseTimePtr(firstNonEmpty(order.LastUpdated, order.DateLastUpdated)),
		FetchedAt:            a.now(),
		RawProviderRef:       "/orders/" + normalizeAnyID(order.ID),
		ShippingID:           normalizeAnyID(order.Shipping.ID),
		CancellationDetail:   order.StatusDetail,
		Tags:                 trimStrings(order.Tags),
	}

	var saleFeeSum float64
	var hasSaleFee bool
	for _, item := range order.OrderItems {
		if item.SaleFee != nil {
			saleFeeSum += *item.SaleFee
			hasSaleFee = true
		}
		snapshot.Items = append(snapshot.Items, domain.OrderItemSnapshot{
			ProviderItemID:      item.Item.ID,
			ProviderVariationID: normalizeAnyID(item.Item.VariationID),
			SellerSKU:           item.Item.SellerSKU,
			EAN:                 attributeValue(item.Item.VariationAttributes, "GTIN", "EAN"),
			Title:               item.Item.Title,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			SaleFeeAmount:       item.SaleFee,
		})
	}
	if hasSaleFee {
		snapshot.SaleFeeAmount = &saleFeeSum
	}

	for _, payment := range order.Payments {
		snapshot.Payments = append(snapshot.Payments, domain.OrderPaymentSnapshot{
			PaymentID:         normalizeAnyID(payment.ID),
			Status:            payment.Status,
			Amount:            firstFloatPtr(payment.TotalPaidAmount, payment.TransactionAmount),
			TransactionAmount: payment.TransactionAmount,
			TotalPaidAmount:   payment.TotalPaidAmount,
		})
	}

	return snapshot
}

func normalizeAccountRef(ref domain.ProviderAccountRef) (domain.ProviderAccountRef, error) {
	ref.TenantID = strings.TrimSpace(ref.TenantID)
	ref.InstallationID = strings.TrimSpace(ref.InstallationID)
	ref.ProviderCode = strings.TrimSpace(ref.ProviderCode)
	ref.ProviderAccountID = strings.TrimSpace(ref.ProviderAccountID)
	if ref.TenantID == "" || ref.InstallationID == "" || ref.ProviderCode == "" || ref.ProviderAccountID == "" {
		return domain.ProviderAccountRef{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider account reference is incomplete")
	}
	return ref, nil
}

func parseTimePtr(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &ts
}

func attributeValue(attributes []mlAttribute, keys ...string) string {
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

func findVariation(variations []mlVariation, variationID string) (mlVariation, bool) {
	for _, variation := range variations {
		if normalizeAnyID(variation.ID) == variationID {
			return variation, true
		}
	}
	return mlVariation{}, false
}

func normalizeAnyID(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case json.Number:
		return typed.String()
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", value))
	}
}

func intValue(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func firstFloatPtr(values ...*float64) *float64 {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func trimStrings(values []string) []string {
	trimmed := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			trimmed = append(trimmed, value)
		}
	}
	return trimmed
}

func limitOrDefault(limit int, fallback int) int {
	if limit > 0 {
		return limit
	}
	return fallback
}

func offsetFromCursor(cursor string) int {
	offset, err := strconv.Atoi(strings.TrimSpace(cursor))
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

type mlItemResponse struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Status            string        `json:"status"`
	AvailableQuantity *int          `json:"available_quantity"`
	Price             json.Number   `json:"price"`
	CurrencyID        string        `json:"currency_id"`
	ListingTypeID     string        `json:"listing_type_id"`
	LastUpdated       string        `json:"last_updated"`
	SellerSKU         string        `json:"seller_sku"`
	SellerCustomField string        `json:"seller_custom_field"`
	Attributes        []mlAttribute `json:"attributes"`
	Variations        []mlVariation `json:"variations"`
}

type mlAccountResponse struct {
	ID        any    `json:"id"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	SiteID    string `json:"site_id"`
	Status    struct {
		SiteStatus string `json:"site_status"`
		List       struct {
			Status string `json:"status"`
		} `json:"list"`
	} `json:"status"`
}

type mlListingPriceResponse struct {
	ListingTypeID  string `json:"listing_type_id"`
	SaleFeeDetails struct {
		PercentageFee *float64 `json:"percentage_fee"`
		FixedFee      *float64 `json:"fixed_fee"`
	} `json:"sale_fee_details"`
}

type mlVariation struct {
	ID                any           `json:"id"`
	AvailableQuantity *int          `json:"available_quantity"`
	SellerSKU         string        `json:"seller_sku"`
	SellerCustomField string        `json:"seller_custom_field"`
	Attributes        []mlAttribute `json:"attributes"`
}

type mlAttribute struct {
	ID        string `json:"id"`
	ValueID   string `json:"value_id"`
	ValueName string `json:"value_name"`
}

type mlOrderResponse struct {
	ID              any             `json:"id"`
	SiteID          string          `json:"site_id"`
	Status          string          `json:"status"`
	StatusDetail    string          `json:"status_detail"`
	DateCreated     string          `json:"date_created"`
	DateClosed      string          `json:"date_closed"`
	LastUpdated     string          `json:"last_updated"`
	DateLastUpdated string          `json:"date_last_updated"`
	OrderItems      []mlOrderItem   `json:"order_items"`
	Payments        []mlPayment     `json:"payments"`
	Shipping        mlOrderShipping `json:"shipping"`
	Tags            []string        `json:"tags"`
	// Buyer carries the fiscal linkage: buyer.billing_info.id is the key for the
	// step-2 billing-info call (docs: faturamento). Additive/pointer — an order
	// payload without a buyer block leaves it nil (honest absence), and mapOrder
	// ignores it (the buyer fiscal fact is a read-time enrichment, not persisted).
	Buyer *mlOrderBuyer `json:"buyer"`
}

// mlOrderBuyer is the subset of the ML order `buyer` object the fiscal flow needs.
// nickname is NOT surfaced past the adapter (the masked buyer identity is derived
// elsewhere from the read model); only billing_info.id feeds the step-2 call.
type mlOrderBuyer struct {
	ID          any                  `json:"id"`
	Nickname    string               `json:"nickname"`
	BillingInfo *mlOrderBillingInfid `json:"billing_info"`
}

type mlOrderBillingInfid struct {
	ID string `json:"id"`
}

type mlOrderShipping struct {
	ID any `json:"id"`
}

type mlOrderItem struct {
	Item      mlOrderItemDetails `json:"item"`
	Quantity  int                `json:"quantity"`
	UnitPrice *float64           `json:"unit_price"`
	SaleFee   *float64           `json:"sale_fee"`
}

type mlOrderItemDetails struct {
	ID                  string        `json:"id"`
	VariationID         any           `json:"variation_id"`
	SellerSKU           string        `json:"seller_sku"`
	Title               string        `json:"title"`
	VariationAttributes []mlAttribute `json:"variation_attributes"`
}

type mlPayment struct {
	ID                any      `json:"id"`
	Status            string   `json:"status"`
	TotalPaidAmount   *float64 `json:"total_paid_amount"`
	TransactionAmount *float64 `json:"transaction_amount"`
}

var (
	_ ports.AccountProber     = (*CapabilityAdapter)(nil)
	_ ports.FeeQuoteReader    = (*CapabilityAdapter)(nil)
	_ ports.ListingReader     = (*CapabilityAdapter)(nil)
	_ ports.StockReader       = (*CapabilityAdapter)(nil)
	_ ports.StockWriter       = (*CapabilityAdapter)(nil)
	_ ports.PriceWriter       = (*CapabilityAdapter)(nil)
	_ ports.OrderReader       = (*CapabilityAdapter)(nil)
	_ ports.BuyerFiscalReader = (*CapabilityAdapter)(nil)
)
