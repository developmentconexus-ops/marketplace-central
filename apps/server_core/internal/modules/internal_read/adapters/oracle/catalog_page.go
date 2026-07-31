package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const (
	catalogListMinLimit   = 1
	catalogListMaxLimit   = 100
	catalogSearchMaxLimit = 50
	catalogCurrency       = "BRL"
)

func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	if cursor.InternalProductID < 0 {
		return ports.CatalogFactPage{}, ports.NewInvalidCursorError()
	}
	if limit < catalogListMinLimit || limit > catalogListMaxLimit {
		return ports.CatalogFactPage{}, fmt.Errorf("invalid catalog page limit: %d", limit)
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogFactPage{}, err
	}

	resolved, err := ports.RequireAssortmentPolicy(policy)
	if err != nil {
		return ports.CatalogFactPage{}, err
	}
	query, args := buildCatalogPageQuery(cursor, limit+1, "", resolved)
	return r.readCatalogPage(ctx, query, args, limit, true)
}

func (r *Reader) SearchCatalogProductFacts(ctx context.Context, q string, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	if cursor.InternalProductID < 0 {
		return ports.CatalogFactPage{}, ports.NewInvalidCursorError()
	}
	if limit < catalogListMinLimit || limit > catalogSearchMaxLimit {
		return ports.CatalogFactPage{}, fmt.Errorf("invalid catalog search limit: %d", limit)
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogFactPage{}, err
	}

	// limit+1 with paginate: the extra row is what proves there are more matches.
	// Fetching exactly `limit` made a truncated result indistinguishable from a
	// complete one.
	resolved, err := ports.RequireAssortmentPolicy(policy)
	if err != nil {
		return ports.CatalogFactPage{}, err
	}
	query, args := buildCatalogPageQuery(cursor, limit+1, q, resolved)
	return r.readCatalogPage(ctx, query, args, limit, true)
}

func (r *Reader) CatalogProductFactsByIDs(ctx context.Context, ids []int64) (ports.CatalogFactPage, error) {
	if len(ids) == 0 {
		return ports.CatalogFactPage{Items: []ports.CatalogProductFact{}, AsOf: r.now().UTC()}, nil
	}
	if len(ids) > catalogListMaxLimit {
		return ports.CatalogFactPage{}, fmt.Errorf("invalid catalog id count: %d", len(ids))
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogFactPage{}, err
	}

	query, args := buildCatalogPageQuery(ports.Cursor{}, len(ids), "", ports.AllProductsAssortment())
	query, args = withCatalogIDFilter(query, args, ids)
	return r.readCatalogPage(ctx, query, args, len(ids), false)
}

// withCatalogIDFilter splices an explicit id list into the page query. The
// FETCH FIRST clause is already the last line, so the predicate goes in front
// of the ORDER BY that precedes it.
func withCatalogIDFilter(query string, args []any, ids []int64) (string, []any) {
	placeholders := make([]string, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf(":%d", len(args)))
	}
	filter := "\n  AND p.CODPROD IN (" + strings.Join(placeholders, ",") + ")"
	orderIndex := strings.LastIndex(query, "\nORDER BY p.CODPROD")
	if orderIndex < 0 {
		return query + filter, args
	}
	return query[:orderIndex] + filter + query[orderIndex:], args
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
			ean             sql.NullString
			reference       sql.NullString
			description     sql.NullString
			brandName       sql.NullString
			ncm             sql.NullString
			activeEANCount  sql.NullInt64
			active          sql.NullString
			stock           sql.NullFloat64
			price           sql.NullString
			activePriceRows sql.NullInt64
			cost            sql.NullString
		)
		if err := rows.Scan(&productID, &ean, &reference, &description, &brandName, &ncm, &activeEANCount, &active, &stock, &price, &activePriceRows, &cost); err != nil {
			return ports.CatalogFactPage{}, wrapOracleError("scan catalog page", err)
		}

		page.Items = append(page.Items, catalogProductFact(productID, ean, reference, description, brandName, ncm, activeEANCount, active, stock, price, activePriceRows, cost))
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

func catalogStockCTE() string {
	return `WITH stock AS (
    SELECT est.CODPROD,
           SUM(CASE WHEN est.CODPARC = 0 THEN NVL(est.ESTOQUE, 0) - NVL(est.RESERVADO, 0) ELSE 0 END) AS sellable_qty
    FROM METALPRD.TGFEST est
    WHERE est.CODEMP IN (1, 2)
      -- Selling-location whitelist: 10101 = 1_REVENDA; 10102 = 2_OUTLET.
      AND est.CODLOCAL IN (10101, 10102)
    GROUP BY est.CODPROD
)`
}

func catalogAssortmentPredicate(policy ports.SellableAssortmentPolicy) string {
	var predicate string
	if policy.OnlyRevenda {
		predicate += "\n  AND UPPER(TRIM(p.USOPROD)) = 'R'"
	}
	if policy.OnlyEmEstoque {
		predicate += "\n  AND NVL(stock.sellable_qty, 0) > 0"
	}
	if policy.OnlyEcommerceEligible {
		predicate += "\n  AND NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'"
	}
	return predicate
}

func buildCatalogPageQuery(cursor ports.Cursor, fetchLimit int, search string, policy ports.SellableAssortmentPolicy) (string, []any) {
	query := "\n" + catalogStockCTE() + `, price_candidates AS (
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
	 p.REFERENCIA AS EAN,
	 p.REFFORN,
    p.DESCRPROD,
	 p.MARCA,
	 p.NCM,
	 (
		 SELECT COUNT(DISTINCT collision.CODPROD)
		 FROM METALPRD.TGFPRO collision
		 WHERE collision.ATIVO = 'S'
		   AND TRIM(collision.REFERENCIA) = TRIM(p.REFERENCIA)
	 ) AS EAN_ACTIVE_COUNT,
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
	query += catalogAssortmentPredicate(policy)

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

func (r *Reader) GetCatalogAssortmentCounts(ctx context.Context, policy *ports.SellableAssortmentPolicy) (ports.CatalogAssortmentCounts, error) {
	resolved, err := ports.RequireAssortmentPolicy(policy)
	if err != nil {
		return ports.CatalogAssortmentCounts{}, err
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return ports.CatalogAssortmentCounts{}, err
	}
	query, args := buildCatalogAssortmentCountQuery(resolved)
	var result ports.CatalogAssortmentCounts
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&result.SellableCount, &result.TotalCount); err != nil {
		return ports.CatalogAssortmentCounts{}, wrapOracleError("read catalog assortment counts", err)
	}
	return result, nil
}

func buildCatalogAssortmentCountQuery(policy ports.SellableAssortmentPolicy) (string, []any) {
	query := "\n" + catalogStockCTE() + `
SELECT
    COUNT(*) AS sellable_count,
    (
        SELECT COUNT(*)
        FROM METALPRD.TGFPRO total
        WHERE total.ATIVO = 'S'
          AND total.CODPROD > 0
    ) AS total_count
FROM METALPRD.TGFPRO p
LEFT JOIN stock ON stock.CODPROD = p.CODPROD
WHERE p.ATIVO = 'S'
  AND p.CODPROD > 0`
	query += catalogAssortmentPredicate(policy)
	return query, nil
}

func catalogProductFact(productID int64, eanValue, reference, description, brandName, ncmValue sql.NullString, activeEANCount sql.NullInt64, active sql.NullString, stock sql.NullFloat64, price sql.NullString, activePriceRows sql.NullInt64, cost sql.NullString) ports.CatalogProductFact {
	ean := nullableString(eanValue)
	qualityFlags := []string{string(domain.QualityComplete)}
	if ean == nil || !catalogdomain.IsValidGTIN(*ean) {
		ean = nil
		qualityFlags = []string{string(domain.QualityInvalidEAN)}
	} else if activeEANCount.Valid && activeEANCount.Int64 >= 2 {
		qualityFlags = append(qualityFlags, string(domain.QualityEANCollision))
	}
	manufacturerReference := nullableString(reference)
	result := ports.CatalogProductFact{
		InternalProductID:     productID,
		Reference:             manufacturerReference,
		ManufacturerReference: manufacturerReference,
		Description:           nullableString(description),
		EAN:                   ean,
		BrandName:             nullableString(brandName),
		NCM:                   nullableNCM(ncmValue),
		QualityFlags:          qualityFlags,
		Active:                strings.EqualFold(active.String, "S"),
		SellableStock: ports.CatalogQuantityFact{
			Quantity: nullableFloat(stock),
		},
		CurrentPrice: ports.CatalogMoneyFact{Currency: catalogCurrency},
		Cost:         ports.CatalogMoneyFact{Currency: catalogCurrency},
	}
	if result.SellableStock.Quantity == nil {
		zero := 0.0
		result.SellableStock.Quantity = &zero
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
