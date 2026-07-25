package routing

import (
	"context"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

// StockBatchReader routes the dashboard's batch stock question to the reader
// that owns the tenant's active source.
//
// Before it existed the batch path was oracle-only, so a tenant on an upload
// source (xlsx, catalogo_cliente) had no stock reader at all and every Stock
// Seguro row rendered "o ERP não informou o estoque disponível deste produto"
// while products_mirror held the quantity.
//
// For a live source the oracle reader is preferred when it is wired, because it
// reads the ERP itself. When it is not, the mirror still holds that same
// source's rows (written by the sankhya sync) with an honest observation time —
// that is the SAME source answering with an older observation, not the silent
// cross-source fallback ADR-17 forbids.
type StockBatchReader struct {
	mirror   internalreadports.StockBatchReader
	live     internalreadports.StockBatchReader
	lookup   tenant_config.ActiveSourceLookup
	tenantID string
}

func NewStockBatchReader(mirror, live internalreadports.StockBatchReader, lookup tenant_config.ActiveSourceLookup, tenantID string) *StockBatchReader {
	return &StockBatchReader{mirror: mirror, live: live, lookup: lookup, tenantID: tenantID}
}

var _ internalreadports.StockBatchReader = (*StockBatchReader)(nil)

func (r *StockBatchReader) GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.StockFact, error) {
	cfg, err := r.lookup.Get(ctx, r.tenantID)
	if err != nil {
		return nil, err
	}
	if !cfg.Source.Valid() {
		return nil, tenant_config.ErrUnknownActiveSource
	}
	ctx = tenant_config.WithActiveSource(ctx, cfg)
	if cfg.Source == tenant_config.SourceSankhya && r.live != nil {
		return r.live.GetStockFactsByIDs(ctx, ids)
	}
	if r.mirror == nil {
		return nil, ErrActiveSourceUnavailable
	}
	ctx = erpinternalread.WithActiveSource(ctx, erpdomain.ImportSource(cfg.Source))
	return r.mirror.GetStockFactsByIDs(ctx, ids)
}
