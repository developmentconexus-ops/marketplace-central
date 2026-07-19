package mercadolivre

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type mlCatalogSearchResponse struct {
	Results []mlCatalogProduct `json:"results"`
}

type mlCatalogProduct struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Title        string                 `json:"title"`
	Attributes   []mlCatalogAttribute   `json:"attributes"`
	BuyBoxWinner *mlCatalogBuyBoxWinner `json:"buy_box_winner"`
}

type mlCatalogAttribute struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ValueName *string `json:"value_name"`
}

type mlCatalogBuyBoxWinner struct {
	SellerID string       `json:"seller_id"`
	Price    *json.Number `json:"price"`
	Currency string       `json:"currency_id"`
}

func (a *CapabilityAdapter) searchCatalogByEAN(ctx context.Context, accountRef domain.ProviderAccountRef, token, ean string) (domain.CatalogSearchResult, error) {
	// ML /products/search 4XXs without site_id (FINDING-M02-COLLECT-4XX, D-86);
	// same contract catalog_match_reader.searchCatalog honors and D-83 proved live.
	query := url.Values{}
	query.Set("site_id", a.siteID)
	query.Set("product_identifier", ean)
	var response mlCatalogSearchResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/products/search?"+query.Encode(), nil, &response); err != nil {
		return domain.CatalogSearchResult{}, mapPricingReaderError(err)
	}

	products := make([]domain.CatalogCandidate, 0, len(response.Results))
	for _, result := range response.Results {
		products = append(products, domain.CatalogCandidate{
			CatalogProductID: strings.TrimSpace(result.ID),
			Title:            firstNonEmpty(result.Name, result.Title),
			Attrs:            mapCatalogAttributes(result.Attributes),
		})
	}
	return domain.CatalogSearchResult{Products: products, FetchedAt: a.now()}, nil
}

func (a *CapabilityAdapter) getCatalogProduct(ctx context.Context, accountRef domain.ProviderAccountRef, token, catalogProductID string) (domain.CatalogProduct, error) {
	var response mlCatalogProduct
	path := "/products/" + url.PathEscape(catalogProductID)
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
		return domain.CatalogProduct{}, mapPricingReaderError(err)
	}

	var winner *domain.BuyBoxWinner
	if response.BuyBoxWinner != nil {
		price, err := providerMoney("buy_box_winner_price", response.BuyBoxWinner.Price, response.BuyBoxWinner.Currency)
		if err != nil {
			return domain.CatalogProduct{}, err
		}
		winner = &domain.BuyBoxWinner{
			Price:    price,
			SellerID: strings.TrimSpace(response.BuyBoxWinner.SellerID),
		}
	}

	return domain.CatalogProduct{
		ID:           strings.TrimSpace(response.ID),
		Title:        firstNonEmpty(response.Name, response.Title),
		BuyBoxWinner: winner,
		Attrs:        mapCatalogAttributes(response.Attributes),
		FetchedAt:    a.now(),
	}, nil
}

func mapCatalogAttributes(attributes []mlCatalogAttribute) []domain.CatalogAttribute {
	result := make([]domain.CatalogAttribute, 0, len(attributes))
	for _, attribute := range attributes {
		var valueName *string
		if attribute.ValueName != nil {
			value := strings.TrimSpace(*attribute.ValueName)
			valueName = &value
		}
		result = append(result, domain.CatalogAttribute{
			ID:        strings.TrimSpace(attribute.ID),
			Name:      strings.TrimSpace(attribute.Name),
			ValueName: valueName,
		})
	}
	return result
}
