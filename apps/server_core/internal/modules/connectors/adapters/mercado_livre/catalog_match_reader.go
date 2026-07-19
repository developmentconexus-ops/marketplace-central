package mercadolivre

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

// ReadCatalogMatch runs the read-only tier-3 catalog-match probe: resolve a
// seller product against the provider catalog by EAN, read its buy-box winner,
// and rank category predictions from domain discovery. All calls are GET.
//
// Category resolution is deliberately buy-box-anchored: when the catalog product
// exposes no buy_box_winner (common — observed live 2026-07-13), the buy-box is
// nil and the caller falls back to the domain-discovery prediction. Unknown
// facts stay null, never zero-defaulted (ADR-17).
func (a *CapabilityAdapter) ReadCatalogMatch(ctx context.Context, input domain.CatalogMatchInput) (domain.CatalogMatchSnapshot, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return domain.CatalogMatchSnapshot{}, err
	}
	ean := strings.TrimSpace(input.EAN)
	query := strings.TrimSpace(input.Query)
	if ean == "" && query == "" {
		return domain.CatalogMatchSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "ean or query is required")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.CatalogMatchSnapshot{}, err
	}

	snapshot := domain.CatalogMatchSnapshot{
		FetchedAt:      a.now(),
		RawProviderRef: "/products/search",
	}

	if ean != "" {
		hits, err := a.searchCatalog(ctx, accountRef, token, ean)
		if err != nil {
			return domain.CatalogMatchSnapshot{}, err
		}
		snapshot.CatalogHits = hits
	}

	var resolvedName string
	if len(snapshot.CatalogHits) > 0 {
		buyBox, name, err := a.readCatalogProduct(ctx, accountRef, token, snapshot.CatalogHits[0].ProductID)
		if err != nil {
			return domain.CatalogMatchSnapshot{}, err
		}
		snapshot.BuyBox = buyBox
		resolvedName = name
	}

	discoveryQuery := firstNonEmpty(query, resolvedName)
	if discoveryQuery != "" {
		predictions, err := a.discoverDomains(ctx, accountRef, token, discoveryQuery)
		if err != nil {
			return domain.CatalogMatchSnapshot{}, err
		}
		snapshot.DomainDiscovery = predictions
	}

	return snapshot, nil
}

func (a *CapabilityAdapter) searchCatalog(ctx context.Context, accountRef domain.ProviderAccountRef, token, ean string) ([]domain.CatalogHit, error) {
	query := url.Values{}
	query.Set("site_id", a.siteID)
	query.Set("product_identifier", ean)

	var response mlProductSearchResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/products/search?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}

	hits := make([]domain.CatalogHit, 0, len(response.Results))
	for _, result := range response.Results {
		productID := strings.TrimSpace(result.ID)
		if productID == "" {
			continue
		}
		hits = append(hits, domain.CatalogHit{
			ProductID: productID,
			Name:      strings.TrimSpace(result.Name),
			DomainID:  strings.TrimSpace(result.DomainID),
			Status:    strings.TrimSpace(result.Status),
		})
	}
	return hits, nil
}

func (a *CapabilityAdapter) readCatalogProduct(ctx context.Context, accountRef domain.ProviderAccountRef, token, productID string) (*domain.BuyBoxSnapshot, string, error) {
	var product mlCatalogProductResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/products/"+url.PathEscape(productID), nil, &product); err != nil {
		return nil, "", err
	}

	name := strings.TrimSpace(product.Name)
	if product.BuyBoxWinner == nil {
		return nil, name, nil
	}

	winner := product.BuyBoxWinner
	return &domain.BuyBoxSnapshot{
		CategoryID:  strings.TrimSpace(winner.CategoryID),
		Price:       winner.Price,
		ListingType: strings.TrimSpace(winner.ListingTypeID),
	}, name, nil
}

func (a *CapabilityAdapter) discoverDomains(ctx context.Context, accountRef domain.ProviderAccountRef, token, discoveryQuery string) ([]domain.DomainPrediction, error) {
	query := url.Values{}
	query.Set("q", discoveryQuery)
	query.Set("limit", strconv.Itoa(catalogMatchDomainDiscoveryLimit))

	var response []mlDomainDiscoveryResult
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/sites/"+url.PathEscape(a.siteID)+"/domain_discovery/search?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}

	predictions := make([]domain.DomainPrediction, 0, len(response))
	for _, result := range response {
		categoryID := strings.TrimSpace(result.CategoryID)
		if categoryID == "" {
			continue
		}
		predictions = append(predictions, domain.DomainPrediction{
			CategoryID:   categoryID,
			CategoryName: strings.TrimSpace(result.CategoryName),
			DomainID:     strings.TrimSpace(result.DomainID),
		})
	}
	return predictions, nil
}

const catalogMatchDomainDiscoveryLimit = 3

type mlProductSearchResponse struct {
	Results []mlProductSearchResult `json:"results"`
}

type mlProductSearchResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	DomainID string `json:"domain_id"`
}

type mlCatalogProductResponse struct {
	Name         string          `json:"name"`
	BuyBoxWinner *mlBuyBoxWinner `json:"buy_box_winner"`
}

type mlBuyBoxWinner struct {
	CategoryID    string   `json:"category_id"`
	Price         *float64 `json:"price"`
	ListingTypeID string   `json:"listing_type_id"`
}

type mlDomainDiscoveryResult struct {
	DomainID     string `json:"domain_id"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

var _ ports.CatalogMatchReader = (*CapabilityAdapter)(nil)
