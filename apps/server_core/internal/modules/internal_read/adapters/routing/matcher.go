package routing

import (
	"context"
	"errors"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

// MirrorMatcher answers product-linking questions from products_mirror for
// EVERY active source, including the live ones.
//
// Reader routes a single product question to the source that owns it, which is
// correct for one read and wrong for candidate generation: generation asks the
// matcher about every listing in an account, and one live ERP round-trip per
// listing cannot finish inside a request deadline. The mirror holds the same
// rows the live sync writes under the same source, so the answer is unchanged —
// only the number of round-trips is. The tenant's active source still decides
// which rows are visible, and there is still no fallback between sources: an
// unknown source is an error, never another source's catalog.
type MirrorMatcher struct {
	mirror   internalreadports.Reader
	lookup   tenant_config.ActiveSourceLookup
	tenantID string
}

func NewMirrorMatcher(mirror internalreadports.Reader, lookup tenant_config.ActiveSourceLookup, tenantID string) *MirrorMatcher {
	return &MirrorMatcher{mirror: mirror, lookup: lookup, tenantID: tenantID}
}

func (m *MirrorMatcher) FindProductsForLinking(ctx context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	cfg, err := m.lookup.Get(ctx, m.tenantID)
	if err != nil {
		return nil, err
	}
	if !cfg.Source.Valid() {
		return nil, tenant_config.ErrUnknownActiveSource
	}
	ctx = tenant_config.WithActiveSource(ctx, cfg)
	ctx = erpinternalread.WithActiveSource(ctx, erpdomain.ImportSource(cfg.Source))
	ctx = erpinternalread.WithSellableAssortment(ctx,
		cfg.SellableAssortment.OnlyRevenda,
		cfg.SellableAssortment.OnlyEmEstoque,
		cfg.SellableAssortment.OnlyEcommerceEligible,
	)
	candidates, err := m.mirror.FindProductsForLinking(ctx, input)
	// A linking anchor that resolves to no product is "sem correspondência" for
	// THAT listing, not a failed run: the reader reports an absent codprod as a
	// not-found error (right for "show me this product", wrong for "search"), and
	// letting it out would abort generation for the whole account on the first
	// unmatched seller_sku.
	var notFound *erpinternalread.ERPProductNotFoundError
	if errors.As(err, &notFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return candidates, nil
}
