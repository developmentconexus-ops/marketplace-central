package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const (
	batchChunkSize = 500
	salesRowCap    = 5000
	salesPeekLimit = salesRowCap + 1
)

// BatchReader owns only batch-shaped Oracle reads. The existing Reader remains
// the compatibility/interactive adapter and is deliberately not modified.
type BatchReader struct {
	db        queryer
	semaphore *oraclebatch.Semaphore
	now       func() time.Time
	logger    *slog.Logger
}

func NewBatchReader(db queryer, semaphore *oraclebatch.Semaphore) *BatchReader {
	return &BatchReader{
		db:        db,
		semaphore: semaphore,
		now:       func() time.Time { return time.Now().UTC() },
		logger:    slog.Default(),
	}
}

func (r *BatchReader) GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.CostAsOf, error) {
	ids = uniquePositiveIDs(ids)
	result := emptyCostFacts(ids)
	if len(ids) == 0 {
		return result, nil
	}
	if err := r.ensureBatchAvailable(ctx); err != nil {
		return nil, err
	}
	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return nil, wrapOracleError("acquire Oracle batch semaphore", err)
	}
	defer release()

	for _, chunk := range oraclebatch.Chunks(ids, batchChunkSize) {
		query, args := buildCostFactsQuery(chunk, r.now().UTC())
		rows, queryErr := r.db.QueryContext(ctx, query, args...)
		if queryErr != nil {
			return nil, wrapOracleError("read cost facts batch", queryErr)
		}
		for rows.Next() {
			var productID sql.NullInt64
			var amount sql.NullFloat64
			var effectiveAt sql.NullTime
			if scanErr := rows.Scan(&productID, &amount, &effectiveAt); scanErr != nil {
				rows.Close()
				return nil, wrapOracleError("scan cost facts batch", scanErr)
			}
			if !productID.Valid {
				continue
			}
			effective := effectiveAtValue(effectiveAt, r.now())
			fact := &domain.CostAsOf{
				ProductID:   int(productID.Int64),
				CompanyID:   1,
				Basis:       domain.CostBasisCUSSEMICM,
				EffectiveAt: effective,
				Amount:      nullableFloat(amount),
				AmountScope: domain.CostAmountScopePerUnit,
				Source:      domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC(), ObservedAt: timePtr(effective)},
			}
			if fact.Amount == nil {
				fact.QualityFlags = []domain.QualityFlag{domain.QualityMissingCost}
			} else {
				fact.QualityFlags = []domain.QualityFlag{domain.QualityComplete}
			}
			result[productID.Int64] = fact
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			rows.Close()
			return nil, wrapOracleError("iterate cost facts batch", rowsErr)
		}
		rows.Close()
	}
	return result, nil
}

func (r *BatchReader) GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.TaxInputs, error) {
	ids = uniquePositiveIDs(ids)
	result := emptyTaxFacts(ids)
	if len(ids) == 0 {
		return result, nil
	}
	if err := r.ensureBatchAvailable(ctx); err != nil {
		return nil, err
	}
	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return nil, wrapOracleError("acquire Oracle batch semaphore", err)
	}
	defer release()

	for _, chunk := range oraclebatch.Chunks(ids, batchChunkSize) {
		query, args := buildTaxFactsQuery(chunk)
		rows, queryErr := r.db.QueryContext(ctx, query, args...)
		if queryErr != nil {
			return nil, wrapOracleError("read tax facts batch", queryErr)
		}
		for rows.Next() {
			var productID sql.NullInt64
			var icms, ipi, pis, cofins sql.NullFloat64
			if scanErr := rows.Scan(&productID, &icms, &ipi, &pis, &cofins); scanErr != nil {
				rows.Close()
				return nil, wrapOracleError("scan tax facts batch", scanErr)
			}
			if !productID.Valid {
				continue
			}
			fact := &domain.TaxInputs{
				ProductID:    int(productID.Int64),
				ICMSAmount:   nullableFloat(icms),
				IPIAmount:    nullableFloat(ipi),
				PISAmount:    nullableFloat(pis),
				COFINSAmount: nullableFloat(cofins),
				Source:       domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
			}
			if fact.ICMSAmount == nil && fact.IPIAmount == nil && fact.PISAmount == nil && fact.COFINSAmount == nil {
				fact.QualityFlags = []domain.QualityFlag{domain.QualityMissingTax}
			} else {
				fact.QualityFlags = []domain.QualityFlag{domain.QualityComplete}
			}
			result[productID.Int64] = fact
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			rows.Close()
			return nil, wrapOracleError("iterate tax facts batch", rowsErr)
		}
		rows.Close()
	}
	return result, nil
}

// GetSalesHistory applies the IC-01 5000-row cap with one peek row. It is a
// separate batch adapter method because reader.go is an owned read-only seam.
func (r *BatchReader) GetSalesHistory(ctx context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	if err := r.ensureBatchAvailable(ctx); err != nil {
		return domain.SalesHistory{}, err
	}
	query, args, err := buildSalesHistoryQuery(input)
	if err != nil {
		return domain.SalesHistory{}, err
	}
	query += fmt.Sprintf(" FETCH FIRST %d ROWS ONLY", salesPeekLimit)
	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return domain.SalesHistory{}, wrapOracleError("acquire Oracle batch semaphore", err)
	}
	defer release()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.SalesHistory{}, wrapOracleError("read capped sales history", err)
	}
	defer rows.Close()
	result := domain.SalesHistory{Source: domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()}}
	for rows.Next() {
		if len(result.Entries) == salesRowCap {
			result.Truncated = true
			break
		}
		var saleDate time.Time
		var productID int
		var productGroupID, priceTableID, operationTypeID sql.NullInt64
		var quantity, netAmount float64
		var lineType string
		if err := rows.Scan(&saleDate, &productID, &productGroupID, &quantity, &netAmount, &priceTableID, &operationTypeID, &lineType); err != nil {
			return domain.SalesHistory{}, wrapOracleError("scan capped sales history", err)
		}
		result.Entries = append(result.Entries, domain.SalesHistoryEntry{
			ProductID: productID, ProductGroupID: nullableInt(productGroupID), SaleDate: saleDate,
			Quantity: quantity, NetAmount: netAmount, IsConfirmed: true,
			PriceTableID: nullableInt(priceTableID), OperationTypeID: nullableInt(operationTypeID),
		})
		_ = lineType
	}
	if err := rows.Err(); err != nil {
		return domain.SalesHistory{}, wrapOracleError("iterate capped sales history", err)
	}
	if result.Truncated {
		r.logger.Warn("oracle_read", "method", "GetSalesHistory", "slow_query", true, "truncated", true, "row_cap", salesRowCap)
	}
	return result, nil
}

func (r *BatchReader) ensureBatchAvailable(ctx context.Context) error {
	if r == nil || r.db == nil || r.semaphore == nil {
		return domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle batch reader is not configured", nil)
	}
	if err := ctx.Err(); err != nil {
		return wrapOracleError("prepare Oracle batch read", err)
	}
	return nil
}

func uniquePositiveIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func emptyCostFacts(ids []int64) map[int64]*domain.CostAsOf {
	result := make(map[int64]*domain.CostAsOf, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	return result
}

func emptyTaxFacts(ids []int64) map[int64]*domain.TaxInputs {
	result := make(map[int64]*domain.TaxInputs, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	return result
}

func buildCostFactsQuery(ids []int64, effectiveAt time.Time) (string, []any) {
	args := make([]any, 0, len(ids)+3)
	placeholders := batchPlaceholders(ids, &args)
	args = append(args, 1, effectiveAt, effectiveAt)
	query := `
SELECT c.CODPROD, c.CUSSEMICM, c.DTATUAL
FROM METALPRD.TGFCUS c
WHERE c.CODPROD IN (` + placeholders + `)
  AND c.CODEMP = :` + fmt.Sprint(len(ids)+1) + `
  AND c.DTATUAL <= :` + fmt.Sprint(len(ids)+2) + `
  AND c.DTATUAL = (
    SELECT MAX(latest.DTATUAL)
    FROM METALPRD.TGFCUS latest
    WHERE latest.CODPROD = c.CODPROD
      AND latest.CODEMP = c.CODEMP
      AND latest.DTATUAL <= :` + fmt.Sprint(len(ids)+3) + `
  )`
	return query, args
}

func buildTaxFactsQuery(ids []int64) (string, []any) {
	args := make([]any, 0, len(ids))
	placeholders := batchPlaceholders(ids, &args)
	return `
SELECT i.CODPROD,
       SUM(CASE WHEN d.IMPOSTO = 'ICMS' THEN d.VALOR END),
       SUM(CASE WHEN d.IMPOSTO = 'IPI' THEN d.VALOR END),
       SUM(CASE WHEN d.IMPOSTO = 'PIS' THEN d.VALOR END),
       SUM(CASE WHEN d.IMPOSTO = 'COFINS' THEN d.VALOR END)
FROM LEANDRO.VW_IMPOSTO_ITEM d
JOIN METALPRD.TGFITE i ON i.NUNOTA = d.NUNOTA AND i.SEQUENCIA = d.SEQUENCIA
WHERE i.CODPROD IN (` + placeholders + `)
GROUP BY i.CODPROD`, args
}

func batchPlaceholders(ids []int64, args *[]any) string {
	placeholders := make([]string, len(ids))
	for index, id := range ids {
		*args = append(*args, id)
		placeholders[index] = fmt.Sprintf(":%d", len(*args))
	}
	return strings.Join(placeholders, ", ")
}

func effectiveAtValue(value sql.NullTime, fallback time.Time) time.Time {
	if value.Valid {
		return value.Time
	}
	return fallback
}

func timePtr(value time.Time) *time.Time {
	return &value
}
