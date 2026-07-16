package oracle

import (
	"context"
	"database/sql"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

func (r *BatchReader) GetICMSCeilingByOrigin(ctx context.Context, originUF domain.UF) (map[domain.UF]*domain.ICMSCeiling, error) {
	if err := r.ensureBatchAvailable(ctx); err != nil {
		return nil, err
	}
	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return nil, wrapOracleError("acquire Oracle batch semaphore", err)
	}
	defer release()

	query, args := buildICMSCeilingQuery(originUF)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, wrapOracleError("read ICMS ceiling", err)
	}
	defer rows.Close()

	result := make(map[domain.UF]*domain.ICMSCeiling)
	for rows.Next() {
		var destinationUF sql.NullInt64
		var ceiling sql.NullFloat64
		if err := rows.Scan(&destinationUF, &ceiling); err != nil {
			return nil, wrapOracleError("scan ICMS ceiling", err)
		}
		if !destinationUF.Valid || !ceiling.Valid {
			if destinationUF.Valid {
				result[domain.UF(destinationUF.Int64)] = nil
			}
			continue
		}
		result[domain.UF(destinationUF.Int64)] = &domain.ICMSCeiling{
			DestinationUF:  domain.UF(destinationUF.Int64),
			CeilingPercent: nullableFloat(ceiling),
			Source:         domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
			QualityFlags:   []domain.QualityFlag{domain.QualityComplete},
		}
	}
	if err := rows.Err(); err != nil {
		return nil, wrapOracleError("iterate ICMS ceiling rows", err)
	}
	return result, nil
}

func buildICMSCeilingQuery(originUF domain.UF) (string, []any) {
	// current-config rate, not effective-dated (D-22.4).
	return `
SELECT UFDEST, MAX(ALIQUFDEST)+NVL(MAX(PERCICMSFCP),0) AS ceiling FROM METALPRD.TGFICM WHERE UFORIG=:uforig GROUP BY UFDEST`, []any{sql.Named("uforig", int64(originUF))}
}
