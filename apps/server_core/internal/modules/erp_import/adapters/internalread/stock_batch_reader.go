package internalread

import (
	"context"
	"strconv"

	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// BatchStockReader answers the dashboard's "how much internal stock does each of
// these products have?" question from the product mirror, in one round-trip per
// call. It is the upload-source counterpart of the oracle BatchStockReader: the
// mirror IS the tenant's ERP data for xlsx/catalogo_cliente sources, and before
// this reader existed those tenants had NO batch stock path at all, so every
// Stock Seguro row rendered "o ERP não informou o estoque" while the mirror held
// the quantity.
//
// Honest-unknown rules (ADR-17): a product with no mirror row is absent from the
// result map; a row whose estoque_total is NULL answers with a nil Quantity plus
// QualityMissingStock — never a zero, which the risk engine would read as "sold
// out" and act on.
type BatchStockReader struct {
	reader *Reader
}

func NewBatchStockReader(reader *Reader) *BatchStockReader {
	return &BatchStockReader{reader: reader}
}

var _ readports.StockBatchReader = (*BatchStockReader)(nil)

func (r *BatchStockReader) GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*readdomain.StockFact, error) {
	if r == nil || r.reader == nil {
		return nil, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "mirror stock batch reader is not configured", nil)
	}
	source, err := r.reader.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return nil, err
	}

	seen := make(map[int64]struct{}, len(ids))
	codes := make([]string, 0, len(ids))
	for _, id := range ids {
		if _, duplicate := seen[id]; duplicate {
			continue
		}
		seen[id] = struct{}{}
		codes = append(codes, strconv.FormatInt(id, 10))
	}

	result := make(map[int64]*readdomain.StockFact, len(codes))
	if len(codes) == 0 {
		return result, nil
	}

	rows, err := r.reader.repo.MirrorProductsByCodes(ctx, r.reader.tenantID, source, codes)
	if err != nil {
		return nil, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}

	for _, row := range rows {
		id, parseErr := strconv.ParseInt(row.CodigoProduto, 10, 64)
		if parseErr != nil {
			// CODPROD is not representable as an internal product id; the caller
			// asked by id, so such a row can never be the answer to this question.
			continue
		}
		dataAt, ok := dataObservationTime(source, row)
		if !ok {
			// The row's data time is unknown, so any quantity we returned would carry
			// a false "observed at" claim. Leave it out: the caller renders the
			// product as unanswered rather than as fresh.
			continue
		}
		fact := &readdomain.StockFact{ProductID: id, Source: r.reader.source(source, dataAt)}
		if row.EstoqueTotal == nil {
			fact.QualityFlags = []readdomain.QualityFlag{readdomain.QualityMissingStock}
			result[id] = fact
			continue
		}
		quantity, err := parseNumber(*row.EstoqueTotal, "estoque_total")
		if err != nil {
			return nil, err
		}
		fact.Quantity = &quantity
		result[id] = fact
	}
	return result, nil
}
