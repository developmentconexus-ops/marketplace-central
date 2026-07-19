package catalog

import (
	"context"
	"strconv"
	"strings"

	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// ProductIdentityReader wraps the catalog service and implements
// pricing/ports.ProductIdentityReader: it resolves a pricing product id to the
// barcode + title + optional catalog price that degrau 3 needs to find a
// provider category. The pricing engine keys products by int; the catalog
// service keys them by decimal string.
type ProductIdentityReader struct {
	svc productLister
}

func NewProductIdentityReader(svc productLister) *ProductIdentityReader {
	return &ProductIdentityReader{svc: svc}
}

// identityCurrency is the marketplace's settlement currency; catalog prices are
// stored as a bare amount with no per-row currency.
const identityCurrency = "BRL"

func (r *ProductIdentityReader) IdentityForProduct(ctx context.Context, productID int) (pricingports.ProductIdentity, bool, error) {
	products, err := r.svc.ListProductsByIDs(ctx, []string{strconv.Itoa(productID)})
	if err != nil {
		return pricingports.ProductIdentity{}, false, err
	}
	if len(products) == 0 {
		return pricingports.ProductIdentity{}, false, nil
	}
	p := products[0]
	titulo := strings.TrimSpace(p.Name)
	ean := strings.TrimSpace(p.EAN)
	// Nothing to match on => a miss, so no pointless downstream probe.
	if ean == "" && titulo == "" {
		return pricingports.ProductIdentity{}, false, nil
	}
	// A non-positive price is absent, never a fabricated 0 quote basis (ADR-17).
	var preco *pricingdomain.Money
	if p.PriceAmount > 0 {
		preco = &pricingdomain.Money{
			Amount:   strconv.FormatFloat(p.PriceAmount, 'f', 2, 64),
			Currency: identityCurrency,
		}
	}
	return pricingports.ProductIdentity{EAN: ean, Titulo: titulo, Preco: preco}, true, nil
}
