package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const stockBatchChunkSize = 500

type BatchStockReader struct {
	db        queryer
	semaphore *oraclebatch.Semaphore
	now       func() time.Time
}

func NewBatchStockReader(db queryer, semaphore *oraclebatch.Semaphore) *BatchStockReader {
	if semaphore == nil {
		semaphore = oraclebatch.NewSemaphore(4)
	}
	return &BatchStockReader{db: db, semaphore: semaphore, now: time.Now}
}

var _ ports.StockBatchReader = (*BatchStockReader)(nil)

func (r *BatchStockReader) GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.StockFact, error) {
	if r == nil || r.db == nil {
		return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle stock batch reader is not configured", nil)
	}

	uniqueIDs := dedupeIDs(ids)
	result := make(map[int64]*domain.StockFact, len(uniqueIDs))
	if len(uniqueIDs) == 0 {
		return result, nil
	}

	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	fetchedAt := r.now().UTC()
	for chunkNumber, chunk := range oraclebatch.Chunks(uniqueIDs, stockBatchChunkSize) {
		query, args := buildStockBatchQuery(chunk)
		rows, err := r.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, wrapOracleError(fmt.Sprintf("read stock batch chunk %d", chunkNumber+1), err)
		}
		for rows.Next() {
			var (
				productID int64
				quantity  sql.NullFloat64
			)
			if err := rows.Scan(&productID, &quantity); err != nil {
				rows.Close()
				return nil, wrapOracleError(fmt.Sprintf("scan stock batch chunk %d", chunkNumber+1), err)
			}
			fact := &domain.StockFact{
				ProductID: productID,
				Quantity:  nullableFloat(quantity),
				Source: domain.SourceMetadata{
					System:    "oracle",
					FetchedAt: fetchedAt,
				},
			}
			if fact.Quantity == nil {
				fact.QualityFlags = []domain.QualityFlag{domain.QualityMissingStock}
			}
			result[productID] = fact
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, wrapOracleError(fmt.Sprintf("iterate stock batch chunk %d", chunkNumber+1), err)
		}
		if err := rows.Close(); err != nil {
			return nil, wrapOracleError(fmt.Sprintf("close stock batch chunk %d", chunkNumber+1), err)
		}
	}
	return result, nil
}

func dedupeIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	unique := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}

func buildStockBatchQuery(ids []int64) (string, []any) {
	policy := domain.DefaultSellableStockPolicy()
	args := make([]any, 0, len(ids)+len(policy.CompanyIDs)+len(policy.LocationIDs)+len(policy.ExcludedLocationIDs))
	query := `
SELECT
    est.CODPROD,
    CAST(SUM(CASE WHEN est.CODPARC = 0 THEN NVL(est.ESTOQUE, 0) - NVL(est.RESERVADO, 0) ELSE 0 END) AS NUMBER(38,2)) AS SELLABLE_QTY
FROM METALPRD.TGFEST est
WHERE est.CODPROD IN (`
	query += appendInt64Placeholders(ids, &args)
	query += ")"
	query += buildIntListClause("est.CODEMP", policy.CompanyIDs, &args)
	query += buildIntListClause("est.CODLOCAL", policy.LocationIDs, &args)
	query += buildNotIntListClause("est.CODLOCAL", policy.ExcludedLocationIDs, &args)
	query += " GROUP BY est.CODPROD"
	return query, args
}

func appendInt64Placeholders(values []int64, args *[]any) string {
	placeholders := make([]string, 0, len(values))
	for _, value := range values {
		*args = append(*args, value)
		placeholders = append(placeholders, fmt.Sprintf(":%d", len(*args)))
	}
	return strings.Join(placeholders, ", ")
}
