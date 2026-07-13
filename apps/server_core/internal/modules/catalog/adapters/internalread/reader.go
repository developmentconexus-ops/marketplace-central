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

func NewReader(reader readports.Reader) *Reader { return &Reader{reader: reader} }

var _ catalogports.CanonicalProductReader = (*Reader)(nil)
var _ catalogports.CanonicalProductReader = UnavailableReader{}
var _ catalogports.ProductReader = UnavailableReader{}

func (r *Reader) ListCanonicalProducts(ctx context.Context) ([]catalogdomain.CanonicalProduct, error) {
	return r.products(ctx, readports.FindProductsInput{})
}
func (r *Reader) SearchCanonicalProducts(ctx context.Context, q string) ([]catalogdomain.CanonicalProduct, error) {
	return r.products(ctx, readports.FindProductsInput{Title: &q})
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
