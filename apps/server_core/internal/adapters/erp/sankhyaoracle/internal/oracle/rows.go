// Package oracle holds the wire shape of Sankhya's tables. Nothing outside
// adapters/erp/sankhyaoracle can import it — the compiler says so, not a
// reviewer — which is why CODPROD may exist here and nowhere else.
package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ProductRow is TGFPRO as it comes off the wire, nulls and all.
//
// Every optional column is a sql.Null*, because "the column was NULL" and "the
// column was zero" are different facts and collapsing them here would make the
// lie unrecoverable downstream.
type ProductRow struct {
	Codprod    int64
	Descrprod  sql.NullString
	Referencia sql.NullString
	Marca      sql.NullString
	NCM        sql.NullString
	Ativo      sql.NullString
	ReadAt     time.Time
}

// SelectActiveProducts is the only statement this package issues.
const SelectActiveProducts = `
SELECT p.CODPROD, p.DESCRPROD, p.REFERENCIA, p.MARCA, p.NCM, p.ATIVO
  FROM METALPRD.TGFPRO p
 WHERE p.ATIVO = 'S'
   AND p.CODPROD > :1
 ORDER BY p.CODPROD
 FETCH FIRST :2 ROWS ONLY`

// FetchActiveProducts pages TGFPRO by CODPROD.
//
// The cursor is the last CODPROD seen, not an offset: OFFSET re-reads the whole
// prefix and silently skips rows when the table changes under the paging.
func FetchActiveProducts(ctx context.Context, db *sql.DB, afterCodprod int64, limit int, now time.Time) ([]ProductRow, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("sankhyaoracle: limit must be positive, got %d", limit)
	}
	rows, err := db.QueryContext(ctx, SelectActiveProducts, afterCodprod, limit)
	if err != nil {
		return nil, fmt.Errorf("sankhyaoracle: query TGFPRO: %w", err)
	}
	defer rows.Close()

	out := make([]ProductRow, 0, limit)
	for rows.Next() {
		var r ProductRow
		if err := rows.Scan(&r.Codprod, &r.Descrprod, &r.Referencia, &r.Marca, &r.NCM, &r.Ativo); err != nil {
			return nil, fmt.Errorf("sankhyaoracle: scan TGFPRO: %w", err)
		}
		r.ReadAt = now.UTC()
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sankhyaoracle: iterate TGFPRO: %w", err)
	}
	return out, nil
}
