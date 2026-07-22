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
	upload   internalreadports.Reader
	live     internalreadports.Reader
	lookup   tenant_config.ActiveSourceLookup
	tenantID string
}

// NewReader builds a routing Reader. live may be nil when no live source is
// wired (e.g. oracle unconfigured); the tenant then fails honest if its
// active source is sankhya rather than silently reading the upload source.
func NewReader(upload, live internalreadports.Reader, lookup tenant_config.ActiveSourceLookup, tenantID string) *Reader {
	return &Reader{upload: upload, live: live, lookup: lookup, tenantID: tenantID}
}

// resolve looks up the tenant's active source, pins it (and, for upload
// sources, the erp sub-source toggle) onto ctx, and returns the reader that
// should serve this call.
func (r *Reader) resolve(ctx context.Context) (internalreadports.Reader, context.Context, error) {
	cfg, err := r.lookup.Get(ctx, r.tenantID)
	if err != nil {
		return nil, ctx, err
	}
	ctx = tenant_config.WithActiveSource(ctx, cfg)
	switch cfg.Source {
	case tenant_config.SourceXLSX, tenant_config.SourceCatalogoCliente:
		ctx = erpinternalread.WithActiveSource(ctx, erpdomain.ImportSource(cfg.Source))
		return r.upload, ctx, nil
	case tenant_config.SourceSankhya:
		if r.live == nil {
			return nil, ctx, ErrActiveSourceUnavailable
		}
		return r.live, ctx, nil
	default:
		return nil, ctx, tenant_config.ErrUnknownActiveSource
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
// reader. The resolved reader must implement the optional CatalogPageReader
// capability (both the upload chain and the oracle chain do); a reader
// without it fails honest with source_unavailable rather than serving
// another source's pages (ADR-17).
func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor internalreadports.Cursor, limit int) (internalreadports.CatalogFactPage, error) {
	pager, ctx, err := r.resolveCatalogPager(ctx)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return pager.ListCatalogProductFacts(ctx, cursor, limit)
}

// SearchCatalogProductFacts routes catalog search to the resolved source's
// reader, with the same honest-failure contract as ListCatalogProductFacts.
func (r *Reader) SearchCatalogProductFacts(ctx context.Context, query string, limit int) (internalreadports.CatalogFactPage, error) {
	pager, ctx, err := r.resolveCatalogPager(ctx)
	if err != nil {
		return internalreadports.CatalogFactPage{}, err
	}
	return pager.SearchCatalogProductFacts(ctx, query, limit)
}

func (r *Reader) resolveCatalogPager(ctx context.Context) (internalreadports.CatalogPageReader, context.Context, error) {
	rd, ctx, err := r.resolve(ctx)
	if err != nil {
		return nil, ctx, err
	}
	pager, ok := rd.(internalreadports.CatalogPageReader)
	if !ok {
		return nil, ctx, internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "active source's reader does not support catalog paging", nil)
	}
	return pager, ctx, nil
}

var _ internalreadports.Reader = (*Reader)(nil)
var _ internalreadports.CatalogPageReader = (*Reader)(nil)
