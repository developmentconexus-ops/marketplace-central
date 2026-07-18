package catalog

import (
	"context"
	"strconv"

	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// ProductChecker answers pricing/ports.ProductChecker over the catalog service.
// It separates "item does not exist" (→ 404 ITEM_NOT_FOUND) from "item exists
// but has no cost fact" (→ 200 SEM_CUSTO). The latter is the costread adapter's
// concern; this checker only reports presence. Any transport error propagates
// (never collapsed into a false "not found").
type ProductChecker struct {
	svc productLister
}

func NewProductChecker(svc productLister) ProductChecker { return ProductChecker{svc: svc} }

var _ pricingports.ProductChecker = ProductChecker{}

func (c ProductChecker) Exists(ctx context.Context, productID int) (bool, error) {
	products, err := c.svc.ListProductsByIDs(ctx, []string{strconv.Itoa(productID)})
	if err != nil {
		return false, err
	}
	return len(products) > 0, nil
}
