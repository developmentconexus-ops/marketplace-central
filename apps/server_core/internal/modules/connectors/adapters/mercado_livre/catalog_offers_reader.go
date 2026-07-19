package mercadolivre

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type mlCatalogOffersResponse struct {
	Paging  mlCatalogOffersPaging `json:"paging"`
	Results []mlCatalogOffer      `json:"results"`
}

type mlCatalogOffersPaging struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type mlCatalogOffer struct {
	// Live ML sends seller_id as a bare NUMBER, not a string (live-drive reprova,
	// D-86); flexString decodes both shapes so a numeric id never sinks the read.
	SellerID  flexString   `json:"seller_id"`
	Price     *json.Number `json:"price"`
	Currency  string       `json:"currency_id"`
	Condition string       `json:"condition"`
	Shipping  struct {
		Mode string `json:"mode"`
	} `json:"shipping"`
}

type mlCatalogProductChildren struct {
	ChildrenIDs []string `json:"children_ids"`
}

// listCatalogOffers resolves offers for a catalog product. Offers/items exist
// only on LEAF catalog products; a PARENT product (children_ids non-empty) is
// not purchasable and its /items 4xxs by design, so read /products/{id} first
// and fan out to each child leaf when present (FINDING-M02-LIVE-2 D2, D-86).
func (a *CapabilityAdapter) listCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, token, catalogProductID string) ([]domain.CatalogOffer, int, error) {
	var product mlCatalogProductChildren
	productPath := "/products/" + url.PathEscape(catalogProductID)
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, productPath, nil, &product); err != nil {
		return nil, 0, mapPricingReaderError(err)
	}
	if len(product.ChildrenIDs) == 0 {
		return a.listLeafCatalogOffers(ctx, accountRef, token, catalogProductID)
	}
	offers := make([]domain.CatalogOffer, 0)
	pageCount := 0
	for _, childID := range product.ChildrenIDs {
		childID = strings.TrimSpace(childID)
		if childID == "" {
			continue // provider noise: a blank ref is not a resolvable leaf
		}
		childOffers, childPages, err := a.listLeafCatalogOffers(ctx, accountRef, token, childID)
		pageCount += childPages
		if err != nil {
			// A leaf whose /items resource is absent (404) legitimately contributes
			// zero offers; skip it. Any other fault (5xx/timeout/401/429) is an
			// UNKNOWN and must abort — returning the surviving siblings would
			// silently undercount competitors, i.e. treat unknown as zero.
			if errors.Is(err, domain.ErrNotFound) {
				continue
			}
			return nil, pageCount, err
		}
		offers = append(offers, childOffers...)
	}
	return offers, pageCount, nil
}

func (a *CapabilityAdapter) listLeafCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, token, catalogProductID string) ([]domain.CatalogOffer, int, error) {
	offers := make([]domain.CatalogOffer, 0)
	offset := 0
	pageCount := 0
	for {
		query := url.Values{}
		query.Set("offset", strconv.Itoa(offset))
		path := "/products/" + url.PathEscape(catalogProductID) + "/items?" + query.Encode()
		var response mlCatalogOffersResponse
		if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
			// A provider fault is a provider fault; ErrCatalogOffersUnavailable is
			// reserved for the flag-off path so composition never reports "disabled"
			// while the capability is on (FINDING-M02-LIVE-2 D1, D-86).
			return nil, pageCount, mapPricingReaderError(err)
		}
		pageCount++
		for _, result := range response.Results {
			price, err := providerMoney("price", result.Price, result.Currency)
			if err != nil {
				return nil, pageCount, err
			}
			offers = append(offers, domain.CatalogOffer{
				SellerID:     strings.TrimSpace(string(result.SellerID)),
				Price:        price,
				Condition:    strings.TrimSpace(result.Condition),
				ShippingMode: strings.TrimSpace(result.Shipping.Mode),
			})
		}
		nextOffset := response.Paging.Offset + len(response.Results)
		if nextOffset >= response.Paging.Total {
			return offers, pageCount, nil
		}
		if len(response.Results) == 0 || nextOffset <= offset {
			return nil, pageCount, fmt.Errorf("%w: incomplete pagination", domain.ErrProviderUnavailable)
		}
		offset = nextOffset
	}
}

func logCatalogOffers(status string, pageCount, offerCount int, started time.Time) {
	slog.Default().Info("catalog offers read",
		"route", "/products/{id}/items",
		"status", status,
		"page_count", pageCount,
		"offer_count", offerCount,
		"duration_ms", time.Since(started).Milliseconds(),
	)
}
