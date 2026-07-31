// Package routing selects, per request, which internal_read.Reader serves a
// tenant's product data: the upload reader (xlsx/catalogo_cliente) or the
// live reader (sankhya/oracle). The tenant's active-source Config drives the
// choice; there is no fallback between sources (ADR-17, fail honest).
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

// ErrActiveSourceUnavailable is returned when the tenant's configured active
// source is sankhya but no live reader is wired (oracle unconfigured or
// unreachable). It is never a fallback to the upload reader: serving the
// wrong source's data would misrepresent the tenant's catalog (ADR-17).
var ErrActiveSourceUnavailable = errors.New("active_source_unavailable")

// Reader routes internal_read.Reader calls to the upload or live reader
// based on the tenant's active-source Config, and pins that Config (plus the
// erp upload sub-source toggle, where applicable) onto the request context
// for downstream consumers (e.g. M-04's cache-key derivation).
type Reader struct {
	upload   internalreadports.Source
	live     internalreadports.Source
	lookup   tenant_config.ActiveSourceLookup
	tenantID string
}

// NewReader builds a routing Reader. live may be nil when no live source is
// wired (e.g. oracle unconfigured); the tenant then fails honest if its
// active source is sankhya rather than silently reading the upload source.
func NewReader(upload, live internalreadports.Source, lookup tenant_config.ActiveSourceLookup, tenantID string) *Reader {
	return &Reader{upload: upload, live: live, lookup: lookup, tenantID: tenantID}
}

// resolve looks up the tenant's active source, pins it (and, for upload
// sources, the erp sub-source toggle) onto ctx, and returns the reader that
// should serve this call.
func (r *Reader) resolve(ctx context.Context) (internalreadports.Source, context.Context, error) {
	rd, ctx, _, err := r.resolveWithConfig(ctx)
	return rd, ctx, err
}

func (r *Reader) resolveWithConfig(ctx context.Context) (internalreadports.Source, context.Context, tenant_config.Config, error) {
	cfg, err := r.lookup.Get(ctx, r.tenantID)
	if err != nil {
		return nil, ctx, tenant_config.Config{}, err
	}
	ctx = tenant_config.WithActiveSource(ctx, cfg)
	switch cfg.Source {
	case tenant_config.SourceXLSX, tenant_config.SourceCatalogoCliente:
		ctx = erpinternalread.WithActiveSource(ctx, erpdomain.ImportSource(cfg.Source))
		return r.upload, ctx, cfg, nil
	case tenant_config.SourceSankhya:
		if r.live == nil {
			return nil, ctx, tenant_config.Config{}, ErrActiveSourceUnavailable
		}
		return r.live, ctx, cfg, nil
	default:
		return nil, ctx, tenant_config.Config{}, tenant_config.ErrUnknownActiveSource
	}
}

func (r *Reader) FindProductsForLinking(ctx context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return nil, err
	}
	return rd.FindProductsForLinking(ctx, input)
}

func (r *Reader) GetSellableStock(ctx context.Context, input internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreaddomain.SellableStock{}, err
	}
	return rd.GetSellableStock(ctx, input)
}

func (r *Reader) GetCurrentPrice(ctx context.Context, input internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreaddomain.CurrentPrice{}, err
	}
	return rd.GetCurrentPrice(ctx, input)
}

func (r *Reader) GetCostAsOf(ctx context.Context, input internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreaddomain.CostAsOf{}, err
	}
	return rd.GetCostAsOf(ctx, input)
}

func (r *Reader) GetSalesHistory(ctx context.Context, input internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreaddomain.SalesHistory{}, err
	}
	return rd.GetSalesHistory(ctx, input)
}

func (r *Reader) GetTaxInputs(ctx context.Context, input internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreaddomain.TaxInputs{}, err
	}
	return rd.GetTaxInputs(ctx, input)
}

// ListCatalogProductFacts routes catalog paging to the resolved source's
// reader, applying the assortment rule that resolveCatalogAssortment picked.
func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor internalreadports.Cursor, limit int, requested *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	pager, ctx, policy, err := r.resolveCatalogAssortment(ctx, requested)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return pager.ListCatalogProductFacts(ctx, cursor, limit, policy)
}

// SearchCatalogProductFacts routes catalog search to the resolved source's
// reader, with the same honest-failure contract as ListCatalogProductFacts.
func (r *Reader) SearchCatalogProductFacts(ctx context.Context, query string, cursor internalreadports.Cursor, limit int, requested *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	pager, ctx, policy, err := r.resolveCatalogAssortment(ctx, requested)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return pager.SearchCatalogProductFacts(ctx, query, cursor, limit, policy)
}

// CatalogProductFactsByIDs routes an explicit-id catalog read to the resolved
// source's reader. It carries no assortment rule on purpose: the caller named
// the exact products it needs — the ones linked to listings, typically — and
// those have to come back whether or not the tenant's cut would have included
// them, or a linked listing would render blank the moment the cut moved.
func (r *Reader) CatalogProductFactsByIDs(ctx context.Context, ids []int64) (internalreadports.CatalogFactPage, error) {
	pager, ctx, err := r.resolve(ctx)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return pager.CatalogProductFactsByIDs(ctx, ids)
}

func (r *Reader) GetCatalogAssortmentCounts(ctx context.Context, requested *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogAssortmentCounts, error) {
	reader, ctx, policy, err := r.resolveCatalogAssortment(ctx, requested)
	if err != nil {
		return internalreadports.CatalogAssortmentCounts{}, err
	}
	return reader.GetCatalogAssortmentCounts(ctx, policy)
}

// resolveCatalogAssortment is the ONE seat that answers "which assortment rule
// applies", for the page and the count alike. A caller that named a policy gets
// that policy applied; a caller that named none (nil) is asking for the tenant's
// stored rule, which is read here from the same active-source row the matcher
// reads for linking. Past this function the policy is always concrete.
func (r *Reader) resolveCatalogAssortment(ctx context.Context, requested *internalreadports.SellableAssortmentPolicy) (internalreadports.Source, context.Context, *internalreadports.SellableAssortmentPolicy, error) {
	rd, ctx, cfg, err := r.resolveWithConfig(ctx)
	if err != nil {
		return nil, ctx, nil, err
	}
	if requested != nil {
		named := *requested
		return rd, ctx, &named, nil
	}
	policy := internalreadports.SellableAssortmentPolicy{
		OnlyRevenda:           cfg.SellableAssortment.OnlyRevenda,
		OnlyEmEstoque:         cfg.SellableAssortment.OnlyEmEstoque,
		OnlyEcommerceEligible: cfg.SellableAssortment.OnlyEcommerceEligible,
	}
	return rd, ctx, &policy, nil
}

var _ internalreadports.Source = (*Reader)(nil)
