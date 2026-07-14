package internalread

import (
	"context"
	"fmt"
	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	catalogports "marketplace-central/apps/server_core/internal/modules/catalog/ports"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"strconv"
	"strings"
	"time"
)

type Reader struct{ reader readports.Reader }

type UnavailableReader struct{ Err error }

func (r UnavailableReader) ListCanonicalProducts(context.Context) ([]catalogdomain.CanonicalProduct, error) {
	return nil, r.Err
}
func (r UnavailableReader) GetCanonicalProduct(context.Context, catalogdomain.InternalProductID) (catalogdomain.CanonicalProduct, error) {
	return catalogdomain.CanonicalProduct{}, r.Err
}
func (r UnavailableReader) SearchCanonicalProducts(context.Context, string) ([]catalogdomain.CanonicalProduct, error) {
	return nil, r.Err
}
func (r UnavailableReader) ListProducts(context.Context) ([]catalogdomain.Product, error) {
	return nil, r.Err
}
func (r UnavailableReader) ListProductsByIDs(context.Context, []string) ([]catalogdomain.Product, error) {
	return nil, r.Err
}
func (r UnavailableReader) GetProduct(context.Context, string) (catalogdomain.Product, error) {
	return catalogdomain.Product{}, r.Err
}
func (r UnavailableReader) SearchProducts(context.Context, string) ([]catalogdomain.Product, error) {
	return nil, r.Err
}
func (r UnavailableReader) ListTaxonomyNodes(context.Context) ([]catalogdomain.TaxonomyNode, error) {
	return nil, r.Err
}
func (r UnavailableReader) ListCatalogProductFacts(context.Context, readports.Cursor, int) (readports.CatalogFactPage, error) {
	return readports.CatalogFactPage{}, r.Err
}
func (r UnavailableReader) SearchCatalogProductFacts(context.Context, string, int) (readports.CatalogFactPage, error) {
	return readports.CatalogFactPage{}, r.Err
}

func NewReader(reader readports.Reader) *Reader { return &Reader{reader: reader} }

var _ catalogports.CanonicalProductReader = (*Reader)(nil)
var _ catalogports.CanonicalProductReader = UnavailableReader{}
var _ catalogports.ProductReader = UnavailableReader{}
var _ readports.CatalogPageReader = UnavailableReader{}

func (r *Reader) ListCanonicalProducts(ctx context.Context) ([]catalogdomain.CanonicalProduct, error) {
	var (
		cursor readports.Cursor
		out    []catalogdomain.CanonicalProduct
	)
	for {
		page, err := r.ListCatalogProductFacts(ctx, cursor, 100)
		if err != nil {
			return nil, err
		}
		products, err := canonicalProductsFromPage(page)
		if err != nil {
			return nil, err
		}
		out = append(out, products...)
		if page.NextCursor == nil {
			return out, nil
		}
		cursor = *page.NextCursor
	}
}
func (r *Reader) SearchCanonicalProducts(ctx context.Context, q string) ([]catalogdomain.CanonicalProduct, error) {
	page, err := r.SearchCatalogProductFacts(ctx, q, 50)
	if err != nil {
		return nil, err
	}
	return canonicalProductsFromPage(page)
}

func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor readports.Cursor, limit int) (readports.CatalogFactPage, error) {
	reader, ok := r.reader.(readports.CatalogPageReader)
	if !ok {
		return readports.CatalogFactPage{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "oracle catalog page reader is unavailable", nil)
	}
	return reader.ListCatalogProductFacts(ctx, cursor, limit)
}

func (r *Reader) SearchCatalogProductFacts(ctx context.Context, q string, limit int) (readports.CatalogFactPage, error) {
	reader, ok := r.reader.(readports.CatalogPageReader)
	if !ok {
		return readports.CatalogFactPage{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "oracle catalog page reader is unavailable", nil)
	}
	return reader.SearchCatalogProductFacts(ctx, q, limit)
}

func canonicalProductsFromPage(page readports.CatalogFactPage) ([]catalogdomain.CanonicalProduct, error) {
	out := make([]catalogdomain.CanonicalProduct, 0, len(page.Items))
	for _, item := range page.Items {
		stock, err := catalogNumericFact(item.SellableStock.Quantity, item.SellableStock.Quality, page.AsOf, readdomain.QualityMissingStock, "stock unavailable")
		if err != nil {
			return nil, fmt.Errorf("catalog stock projection: %w", err)
		}
		priceValue, err := decimalFloat(item.CurrentPrice.Amount)
		if err != nil {
			return nil, fmt.Errorf("catalog price projection: %w", err)
		}
		price, err := catalogNumericFact(priceValue, item.CurrentPrice.Quality, page.AsOf, readdomain.QualityMissingPrice, "price unavailable")
		if err != nil {
			return nil, fmt.Errorf("catalog price projection: %w", err)
		}
		costValue, err := decimalFloat(item.Cost.Amount)
		if err != nil {
			return nil, fmt.Errorf("catalog cost projection: %w", err)
		}
		cost, err := catalogNumericFact(costValue, item.Cost.Quality, page.AsOf, readdomain.QualityMissingCost, "cost unavailable")
		if err != nil {
			return nil, fmt.Errorf("catalog cost projection: %w", err)
		}
		out = append(out, catalogdomain.CanonicalProduct{
			InternalProductID:     catalogdomain.InternalProductID(item.InternalProductID),
			Name:                  ptr(item.Description),
			EAN:                   item.EAN,
			ManufacturerReference: item.Reference,
			PriceAmount:           price,
			CostAmount:            cost,
			StockQuantity:         stock,
		})
	}
	return out, nil
}

func catalogNumericFact(value *float64, flags []string, asOf time.Time, missingFlag readdomain.QualityFlag, fallbackReason string) (catalogdomain.NumericSourceFact, error) {
	domainFlags := make([]readdomain.QualityFlag, 0, len(flags))
	for _, flag := range flags {
		domainFlags = append(domainFlags, readdomain.QualityFlag(flag))
	}
	return fact(value, readdomain.SourceMetadata{System: "oracle", FetchedAt: asOf}, domainFlags, missingFlag, fallbackReason)
}

func decimalFloat(value *string) (*float64, error) {
	if value == nil {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(strings.TrimSpace(*value), 64)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
func (r *Reader) GetCanonicalProduct(ctx context.Context, id catalogdomain.InternalProductID) (catalogdomain.CanonicalProduct, error) {
	n := int(id)
	products, err := r.products(ctx, readports.FindProductsInput{ProductID: &n})
	if err != nil {
		return catalogdomain.CanonicalProduct{}, err
	}
	if len(products) != 1 {
		return catalogdomain.CanonicalProduct{}, fmt.Errorf("CATALOG_PRODUCT_NOT_FOUND")
	}
	return products[0], nil
}
func (r *Reader) products(ctx context.Context, input readports.FindProductsInput) ([]catalogdomain.CanonicalProduct, error) {
	candidates, err := r.reader.FindProductsForLinking(ctx, input)
	if err != nil {
		return nil, err
	}
	out := make([]catalogdomain.CanonicalProduct, 0, len(candidates))
	for _, c := range candidates {
		if c.InternalProductID == nil {
			continue
		}
		p, err := r.product(ctx, c)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}
func (r *Reader) product(ctx context.Context, c readdomain.ProductCandidate) (catalogdomain.CanonicalProduct, error) {
	id := int(*c.InternalProductID)
	stock, err := r.reader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: id, Policy: readdomain.DefaultSellableStockPolicy()})
	if err != nil {
		return catalogdomain.CanonicalProduct{}, err
	}
	price, err := r.reader.GetCurrentPrice(ctx, readports.CurrentPriceInput{ProductID: id, Policy: readdomain.DefaultCurrentPricePolicy(time.Now().UTC())})
	if err != nil {
		return catalogdomain.CanonicalProduct{}, err
	}
	cost, err := r.reader.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: id, Policy: readdomain.CostAsOfPolicy{CompanyID: 1, EffectiveAt: time.Now().UTC(), Basis: readdomain.CostBasisCUSSEMICM}})
	if err != nil {
		return catalogdomain.CanonicalProduct{}, err
	}
	costFact, err := fact(cost.Amount, cost.Source, cost.QualityFlags, readdomain.QualityMissingCost, "cost unavailable")
	if err != nil {
		return catalogdomain.CanonicalProduct{}, fmt.Errorf("catalog cost projection: %w", err)
	}
	priceFact, err := fact(price.Amount, price.Source, price.QualityFlags, readdomain.QualityMissingPrice, "price unavailable")
	if err != nil {
		return catalogdomain.CanonicalProduct{}, fmt.Errorf("catalog price projection: %w", err)
	}
	stockFact, err := fact(stock.Quantity, stock.Source, stock.QualityFlags, readdomain.QualityMissingStock, "stock unavailable")
	if err != nil {
		return catalogdomain.CanonicalProduct{}, fmt.Errorf("catalog stock projection: %w", err)
	}
	return catalogdomain.CanonicalProduct{InternalProductID: *c.InternalProductID, Name: c.Name, EAN: c.EAN, ManufacturerReference: c.ReferenceCode, BrandName: c.BrandName, ProductGroupName: c.ProductGroupName, CostAmount: costFact, PriceAmount: priceFact, StockQuantity: stockFact}, nil
}
func fact(value *float64, source readdomain.SourceMetadata, flags []readdomain.QualityFlag, missingFlag readdomain.QualityFlag, fallbackReason string) (catalogdomain.NumericSourceFact, error) {
	q := catalogdomain.SourceFactQualityCurrent
	reason := (*string)(nil)
	var observed *time.Time
	if source.ObservedAt != nil {
		observed = source.ObservedAt
	} else if !source.FetchedAt.IsZero() {
		t := source.FetchedAt
		observed = &t
	}
	if containsFlag(flags, readdomain.QualityAmbiguousProduct) {
		q = catalogdomain.SourceFactQualityConflict
		value = nil
	} else if containsFlag(flags, readdomain.QualityStaleSource) && value != nil && observed != nil {
		q = catalogdomain.SourceFactQualityStale
		r := "stale_source"
		reason = &r
	} else if value == nil || containsFlag(flags, missingFlag) || containsFlag(flags, readdomain.QualityMissingProduct) || containsFlag(flags, readdomain.QualityStaleSource) {
		q = catalogdomain.SourceFactQualityUnknown
		value = nil
		observed = nil
		r := fallbackReason
		for _, flag := range flags {
			if flag != readdomain.QualityComplete {
				r = string(flag)
				break
			}
		}
		reason = &r
	}
	system := strings.TrimSpace(source.System)
	if system == "" {
		system = "oracle"
	}
	return catalogdomain.NewNumericSourceFact(system, value, q, observed, reason)
}
func containsFlag(flags []readdomain.QualityFlag, target readdomain.QualityFlag) bool {
	for _, flag := range flags {
		if flag == target {
			return true
		}
	}
	return false
}
func (r *Reader) ListProducts(ctx context.Context) ([]catalogdomain.Product, error) {
	ps, e := r.ListCanonicalProducts(ctx)
	return legacy(ps), e
}
func (r *Reader) ListProductsByIDs(ctx context.Context, ids []string) ([]catalogdomain.Product, error) {
	out := make([]catalogdomain.Product, 0, len(ids))
	for _, v := range ids {
		n, e := strconv.Atoi(v)
		if e != nil || n <= 0 {
			continue
		}
		p, e := r.GetCanonicalProduct(ctx, catalogdomain.InternalProductID(n))
		if e != nil {
			return nil, e
		}
		out = append(out, legacy([]catalogdomain.CanonicalProduct{p})...)
	}
	return out, nil
}
func (r *Reader) GetProduct(ctx context.Context, id string) (catalogdomain.Product, error) {
	n, e := strconv.Atoi(id)
	if e != nil {
		return catalogdomain.Product{}, fmt.Errorf("CATALOG_PRODUCT_NOT_FOUND")
	}
	p, e := r.GetCanonicalProduct(ctx, catalogdomain.InternalProductID(n))
	if e != nil {
		return catalogdomain.Product{}, e
	}
	return legacy([]catalogdomain.CanonicalProduct{p})[0], nil
}
func (r *Reader) SearchProducts(ctx context.Context, q string) ([]catalogdomain.Product, error) {
	ps, e := r.SearchCanonicalProducts(ctx, q)
	return legacy(ps), e
}
func (r *Reader) ListTaxonomyNodes(context.Context) ([]catalogdomain.TaxonomyNode, error) {
	return []catalogdomain.TaxonomyNode{}, nil
}
func legacy(ps []catalogdomain.CanonicalProduct) []catalogdomain.Product {
	out := make([]catalogdomain.Product, 0, len(ps))
	for _, p := range ps {
		out = append(out, catalogdomain.Product{ProductID: strconv.Itoa(int(p.InternalProductID)), SKU: strings.TrimSpace(ptr(p.ManufacturerReference)), Name: p.Name, Status: "active", EAN: ptr(p.EAN), Reference: ptr(p.ManufacturerReference)})
	}
	return out
}
func ptr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
