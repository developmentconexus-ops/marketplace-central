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
	SellerID  string       `json:"seller_id"`
	Price     *json.Number `json:"price"`
	Currency  string       `json:"currency_id"`
	Condition string       `json:"condition"`
	Shipping  struct {
		Mode string `json:"mode"`
	} `json:"shipping"`
}

func (a *CapabilityAdapter) listCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, token, catalogProductID string) ([]domain.CatalogOffer, int, error) {
	offers := make([]domain.CatalogOffer, 0)
	offset := 0
	pageCount := 0
	for {
		query := url.Values{}
		query.Set("offset", strconv.Itoa(offset))
		path := "/products/" + url.PathEscape(catalogProductID) + "/items?" + query.Encode()
		var response mlCatalogOffersResponse
		if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
			mapped := mapPricingReaderError(err)
			if errors.Is(mapped, domain.ErrUnauthorized) || errors.Is(mapped, domain.ErrNotFound) || errors.Is(mapped, domain.ErrRateLimited) {
				return nil, pageCount, mapped
			}
			return nil, pageCount, fmt.Errorf("%w: %w", domain.ErrCatalogOffersUnavailable, mapped)
		}
		pageCount++
		for _, result := range response.Results {
			price, err := providerMoney("price", result.Price, result.Currency)
			if err != nil {
				return nil, pageCount, err
			}
			offers = append(offers, domain.CatalogOffer{
				SellerID:     strings.TrimSpace(result.SellerID),
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
			return nil, pageCount, fmt.Errorf("%w: incomplete pagination", domain.ErrCatalogOffersUnavailable)
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
