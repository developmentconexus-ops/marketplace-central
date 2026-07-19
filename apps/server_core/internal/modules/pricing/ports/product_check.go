package ports

import "context"

// ProductChecker reports whether a product exists at all, distinct from
// whether it has a cost fact. The IC-04 error matrix separates the two: an
// unknown product id is 404 ITEM_NOT_FOUND, whereas a known product with no
// custo_erp is a 200 decomposition carrying blocking_state SEM_CUSTO (ADR-17
// — absent cost is not an error). Backed by the catalog/internal_read reader.
type ProductChecker interface {
	Exists(ctx context.Context, productID int) (bool, error)
}
