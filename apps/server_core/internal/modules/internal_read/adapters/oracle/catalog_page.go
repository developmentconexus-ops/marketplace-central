package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const (
	catalogListMinLimit   = 1
	catalogListMaxLimit   = 100
	catalogSearchMaxLimit = 50
	catalogCurrency       = "BRL"
)

func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor ports.Cursor, limit int) (ports.CatalogFactPage, error) {
	if cursor.InternalProductID < 0 {
		return ports.CatalogFactPage{}, ports.NewInvalidCursorError()
	}
	if limit < catalogListMinLimit || limit > catalogListMaxLimit {
		return ports.CatalogFactPage{}, fmt.Errorf("invalid catalog page limit: %d", limit)
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogFactPage{}, err
	}

	query, args := buildCatalogPageQuery(cursor, limit+1, "")
	return r.readCatalogPage(ctx, query, args, limit, true)
}

func (r *Reader) SearchCatalogProductFacts(ctx context.Context, q string, limit int) (ports.CatalogFactPage, error) {
	if limit < catalogListMinLimit || limit > catalogSearchMaxLimit {
		return ports.CatalogFactPage{}, fmt.Errorf("invalid catalog search limit: %d", limit)
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogFactPage{}, err
	}

	query, args := buildCatalogPageQuery(ports.Cursor{}, limit, q)
	return r.readCatalogPage(ctx, query, args, limit, false)
}

func (r *Reader) readCatalogPage(ctx context.Context, query string, args []any, limit int, paginate bool) (ports.CatalogFactPage, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return ports.CatalogFactPage{}, wrapOracleError("read catalog page", err)
	}
	defer rows.Close()

	asOf := r.now().UTC()
	page := ports.CatalogFactPage{Items: make([]ports.CatalogProductFact, 0, limit), AsOf: asOf}
	for rows.Next() {
		var (
			productID       int64
			reference       sql.NullString
			description     sql.NullString
			ean             sql.NullString
			active          sql.NullString
			stock           sql.NullFloat64
			price           sql.NullString
			activePriceRows sql.NullInt64
			cost            sql.NullString
		)
		if err := rows.Scan(&productID, &reference, &description, &ean, &active, &stock, &price, &activePriceRows, &cost); err != nil {
			return ports.CatalogFactPage{}, wrapOracleError("scan catalog page", err)
		}

		page.Items = append(page.Items, catalogProductFact(productID, reference, description, ean, active, stock, price, activePriceRows, cost))
	}
	if err := rows.Err(); err != nil {
		return ports.CatalogFactPage{}, wrapOracleError("iterate catalog page", err)
	}

	if paginate && len(page.Items) > limit {
		last := page.Items[limit-1].InternalProductID
		page.Items = page.Items[:limit]
		page.NextCursor = &ports.Cursor{InternalProductID: last}
	}
	return page, nil
}

func buildCatalogPageQuery(cursor ports.Cursor, fetchLimit int, search string) (string, []any) {
	query := `
WITH stock AS (
    SELECT est.CODPROD,
           SUM(CASE WHEN est.CODPARC = 0 THEN NVL(est.ESTOQUE, 0) - NVL(est.RESERVADO, 0) ELSE 0 END) AS sellable_qty
    FROM METALPRD.TGFEST est
    WHERE est.CODEMP IN (1, 2)
      AND est.CODLOCAL = 10101
    GROUP BY est.CODPROD
), price_candidates AS (
    SELECT e.CODPROD,
           e.VLRVENDA,
           ROW_NUMBER() OVER (PARTITION BY e.CODPROD ORDER BY t.DTVIGOR DESC, e.NUTAB DESC) AS price_rank,
           COUNT(*) OVER (PARTITION BY e.CODPROD, e.NUTAB) AS active_price_rows
    FROM METALPRD.TGFEXC e
    JOIN METALPRD.TGFTAB t ON t.NUTAB = e.NUTAB
    WHERE t.CODTAB = 0
      AND e.CODLOCAL = 10101
      AND t.DTVIGOR <= SYSDATE
), price AS (
    SELECT CODPROD,
           MAX(CASE WHEN price_rank = 1 AND active_price_rows = 1 THEN VLRVENDA END) AS amount,
           MAX(CASE WHEN price_rank = 1 THEN active_price_rows ELSE 0 END) AS active_price_rows
    FROM price_candidates
    GROUP BY CODPROD
), cost_candidates AS (
    SELECT c.CODPROD,
           c.CUSSEMICM,
           ROW_NUMBER() OVER (PARTITION BY c.CODPROD ORDER BY c.DTATUAL DESC, c.CODLOCAL DESC) AS cost_rank
    FROM METALPRD.TGFCUS c
    WHERE c.CODEMP = 1
      AND c.DTATUAL <= SYSDATE
), cost AS (
    SELECT CODPROD,
           MAX(CASE WHEN cost_rank = 1 THEN CUSSEMICM END) AS amount
    FROM cost_candidates
    GROUP BY CODPROD
)
SELECT
    p.CODPROD,
    p.REFERENCIA,
    p.DESCRPROD,
    CAST(NULL AS VARCHAR2(1)) AS EAN,
    p.ATIVO,
    stock.sellable_qty,
    price.amount,
    price.active_price_rows,
    cost.amount
FROM METALPRD.TGFPRO p
LEFT JOIN stock ON stock.CODPROD = p.CODPROD
LEFT JOIN price ON price.CODPROD = p.CODPROD
LEFT JOIN cost ON cost.CODPROD = p.CODPROD
WHERE p.ATIVO = 'S'
  AND p.CODPROD > 0`

	args := make([]any, 0, 3)
	if cursor.InternalProductID > 0 {
		args = append(args, cursor.InternalProductID)
		query += fmt.Sprintf("\n  AND p.CODPROD > :%d", len(args))
	}
	if value := strings.TrimSpace(search); value != "" {
		args = append(args, "%"+strings.ToUpper(value)+"%")
		query += fmt.Sprintf("\n  AND UPPER(p.DESCRPROD) LIKE :%d", len(args))
	}
	query += "\nORDER BY p.CODPROD"
	args = append(args, fetchLimit)
	query += fmt.Sprintf("\nFETCH FIRST :%d ROWS ONLY", len(args))
	return query, args
}

func catalogProductFact(productID int64, reference, description, ean, active sql.NullString, stock sql.NullFloat64, price sql.NullString, activePriceRows sql.NullInt64, cost sql.NullString) ports.CatalogProductFact {
	result := ports.CatalogProductFact{
		InternalProductID: productID,
		Reference:         nullableString(reference),
		Description:       nullableString(description),
		EAN:               nullableString(ean),
		Active:            strings.EqualFold(active.String, "S"),
		SellableStock: ports.CatalogQuantityFact{
			Quantity: nullableFloat(stock),
		},
		CurrentPrice: ports.CatalogMoneyFact{Currency: catalogCurrency},
		Cost:         ports.CatalogMoneyFact{Currency: catalogCurrency},
	}
	if result.SellableStock.Quantity == nil {
		result.SellableStock.Quality = []string{string(domain.QualityMissingStock)}
	}
	if activePriceRows.Valid && activePriceRows.Int64 > 1 {
		result.CurrentPrice.Quality = []string{string(domain.QualityAmbiguousPrice)}
	} else if !price.Valid || strings.TrimSpace(price.String) == "" {
		result.CurrentPrice.Quality = []string{string(domain.QualityMissingPrice)}
	} else {
		value := strings.TrimSpace(price.String)
		result.CurrentPrice.Amount = &value
	}
	if !cost.Valid || strings.TrimSpace(cost.String) == "" {
		result.Cost.Quality = []string{string(domain.QualityMissingCost)}
	} else {
		value := strings.TrimSpace(cost.String)
		result.Cost.Amount = &value
	}
	return result
}

var _ ports.CatalogPageReader = (*Reader)(nil)
