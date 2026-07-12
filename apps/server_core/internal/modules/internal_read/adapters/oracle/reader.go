package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type queryer interface {
	PingContext(ctx context.Context) error
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Reader struct {
	db  queryer
	now func() time.Time
}

func NewReader(db queryer) *Reader {
	return &Reader{
		db:  db,
		now: time.Now,
	}
}

func (r *Reader) FindProductsForLinking(ctx context.Context, input ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return nil, err
	}

	query, args, err := buildFindProductsQuery(input)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, wrapOracleError("find products", err)
	}
	defer rows.Close()

	var candidates []domain.ProductCandidate
	for rows.Next() {
		var (
			productID      int
			name           string
			referenceValue sql.NullString
			productGroupID sql.NullInt64
			brandName      sql.NullString
			brandID        sql.NullInt64
			activeValue    string
			usageType      sql.NullString
		)
		if err := rows.Scan(&productID, &name, &referenceValue, &productGroupID, &brandName, &brandID, &activeValue, &usageType); err != nil {
			return nil, wrapOracleError("scan product candidate", err)
		}

		candidate := domain.ProductCandidate{
			ProductID:      productID,
			Name:           name,
			EAN:            nullableString(referenceValue),
			ReferenceCode:  nullableString(referenceValue),
			BrandID:        nullableInt(brandID),
			BrandName:      nullableString(brandName),
			ProductGroupID: nullableInt(productGroupID),
			IsActive:       strings.EqualFold(activeValue, "S"),
			UsageType:      nullableString(usageType),
			Source: domain.SourceMetadata{
				System:    "oracle",
				FetchedAt: r.now().UTC(),
			},
			QualityFlags: []domain.QualityFlag{domain.QualityComplete},
		}
		candidates = append(candidates, candidate)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapOracleError("iterate product candidate rows", err)
	}

	if len(candidates) == 0 {
		return []domain.ProductCandidate{{
			Source:       domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
			QualityFlags: []domain.QualityFlag{domain.QualityMissingProduct},
		}}, nil
	}
	if len(candidates) > 1 {
		for i := range candidates {
			candidates[i].QualityFlags = []domain.QualityFlag{domain.QualityAmbiguousProduct}
		}
	}
	return candidates, nil
}

func (r *Reader) GetSellableStock(ctx context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.SellableStock{}, err
	}

	policy := input.Policy
	if len(policy.CompanyIDs) == 0 && len(policy.LocationIDs) == 0 && len(policy.ExcludedLocationIDs) == 0 && policy.Formula == "" && policy.Scope == "" {
		policy = domain.DefaultSellableStockPolicy()
	}

	query, args, err := buildSellableStockQuery(input.ProductID, policy)
	if err != nil {
		return domain.SellableStock{}, err
	}

	var quantity sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&quantity); err != nil {
		return domain.SellableStock{}, wrapOracleError("read sellable stock", err)
	}

	result := domain.SellableStock{
		ProductID: input.ProductID,
		Quantity:  nullableFloat(quantity),
		Policy:    policy,
		Source: domain.SourceMetadata{
			System:    "oracle",
			FetchedAt: r.now().UTC(),
		},
	}
	if result.Quantity == nil {
		result.QualityFlags = []domain.QualityFlag{domain.QualityMissingStock}
	}
	return result, nil
}

func (r *Reader) GetCurrentPrice(ctx context.Context, input ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.CurrentPrice{}, err
	}

	policy := input.Policy
	if policy.PriceTableID == 0 && policy.LocationID == 0 && policy.EffectiveAt.IsZero() {
		policy = domain.DefaultCurrentPricePolicy(r.now().UTC())
	}

	query := `
SELECT p.CODTAB, p.CODLOCAL, p.PRECO, p.NUTAB, p.VIGENTE_DE, p.VIGENTE_ATE
FROM LEANDRO.VW_PRECO_TABELA p
WHERE p.CODPROD = :1
  AND p.CODTAB = :2
  AND p.CODLOCAL = :3
  AND p.VIGENTE_DE <= :4
  AND (p.VIGENTE_ATE IS NULL OR p.VIGENTE_ATE > :5)
ORDER BY p.VIGENTE_DE DESC, p.NUTAB DESC
FETCH FIRST 2 ROWS ONLY`

	rows, err := r.db.QueryContext(ctx, query, input.ProductID, policy.PriceTableID, policy.LocationID, policy.EffectiveAt, policy.EffectiveAt)
	if err != nil {
		return domain.CurrentPrice{}, wrapOracleError("read current price", err)
	}
	defer rows.Close()

	result := domain.CurrentPrice{
		ProductID:    input.ProductID,
		PriceTableID: policy.PriceTableID,
		LocationID:   policy.LocationID,
		Source:       domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
	}

	rowCount := 0
	for rows.Next() {
		rowCount++
		if rowCount > 1 {
			return domain.CurrentPrice{}, domain.NewReadError(domain.ReadErrorUnsupportedQuery, "oracle current price query returned multiple active rows", nil)
		}
		var (
			priceTableID  int
			locationID    int
			amount        sql.NullFloat64
			versionID     sql.NullInt64
			effectiveFrom time.Time
			effectiveTo   sql.NullTime
		)
		if err := rows.Scan(&priceTableID, &locationID, &amount, &versionID, &effectiveFrom, &effectiveTo); err != nil {
			return domain.CurrentPrice{}, wrapOracleError("scan current price", err)
		}
		result.PriceTableID = priceTableID
		result.LocationID = locationID
		result.Amount = nullableFloat(amount)
		result.VersionID = nullableInt64(versionID)
		result.EffectiveFrom = &effectiveFrom
		result.EffectiveTo = nullableTime(effectiveTo)
	}
	if err := rows.Err(); err != nil {
		return domain.CurrentPrice{}, wrapOracleError("iterate current price rows", err)
	}
	if rowCount == 0 {
		result.QualityFlags = []domain.QualityFlag{domain.QualityMissingPrice}
	}
	return result, nil
}

func (r *Reader) GetCostAsOf(ctx context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.CostAsOf{}, err
	}

	policy := input.Policy
	if policy.Basis == "" {
		policy.Basis = domain.CostBasisCUSSEMICM
	}

	query := `
SELECT c.CUSSEMICM, c.DTATUAL
FROM METALPRD.TGFCUS c
WHERE c.CODPROD = :1
  AND c.CODEMP = :2
  AND c.DTATUAL <= :3
ORDER BY c.DTATUAL DESC, c.CODLOCAL DESC
FETCH FIRST 1 ROW ONLY`

	var (
		amount      sql.NullFloat64
		effectiveAt sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, query, input.ProductID, policy.CompanyID, policy.EffectiveAt).Scan(&amount, &effectiveAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CostAsOf{
				ProductID:   input.ProductID,
				CompanyID:   policy.CompanyID,
				Basis:       policy.Basis,
				EffectiveAt: policy.EffectiveAt,
				Amount:      nil,
				AmountScope: domain.CostAmountScopePerUnit,
				Source:      domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
				QualityFlags: []domain.QualityFlag{
					domain.QualityMissingCost,
				},
			}, nil
		}
		return domain.CostAsOf{}, wrapOracleError("read cost as-of", err)
	}

	result := domain.CostAsOf{
		ProductID:   input.ProductID,
		CompanyID:   policy.CompanyID,
		Basis:       policy.Basis,
		EffectiveAt: policy.EffectiveAt,
		Amount:      nullableFloat(amount),
		AmountScope: domain.CostAmountScopePerUnit,
		Source: domain.SourceMetadata{
			System:     "oracle",
			FetchedAt:  r.now().UTC(),
			ObservedAt: nullableTime(effectiveAt),
		},
	}
	if result.Amount == nil {
		result.QualityFlags = []domain.QualityFlag{domain.QualityMissingCost}
	}
	return result, nil
}

func (r *Reader) GetSalesHistory(ctx context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.SalesHistory{}, err
	}

	query, args, err := buildSalesHistoryQuery(input)
	if err != nil {
		return domain.SalesHistory{}, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.SalesHistory{}, wrapOracleError("read sales history", err)
	}
	defer rows.Close()

	result := domain.SalesHistory{
		Source: domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
	}
	for rows.Next() {
		var (
			saleDate        time.Time
			productID       int
			productGroupID  sql.NullInt64
			quantity        float64
			netAmount       float64
			priceTableID    sql.NullInt64
			operationTypeID sql.NullInt64
			lineType        string
		)
		if err := rows.Scan(&saleDate, &productID, &productGroupID, &quantity, &netAmount, &priceTableID, &operationTypeID, &lineType); err != nil {
			return domain.SalesHistory{}, wrapOracleError("scan sales history", err)
		}
		result.Entries = append(result.Entries, domain.SalesHistoryEntry{
			ProductID:       productID,
			ProductGroupID:  nullableInt(productGroupID),
			SaleDate:        saleDate,
			Quantity:        quantity,
			NetAmount:       netAmount,
			IsConfirmed:     true,
			PriceTableID:    nullableInt(priceTableID),
			OperationTypeID: nullableInt(operationTypeID),
		})
		_ = lineType
	}
	if err := rows.Err(); err != nil {
		return domain.SalesHistory{}, wrapOracleError("iterate sales history rows", err)
	}
	return result, nil
}

func (r *Reader) GetTaxInputs(ctx context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.TaxInputs{}, err
	}

	policy := input.Policy
	query := `
SELECT
	SUM(CASE WHEN d.IMPOSTO = 'ICMS' THEN d.VALOR END) AS ICMS_AMOUNT,
	SUM(CASE WHEN d.IMPOSTO = 'IPI' THEN d.VALOR END) AS IPI_AMOUNT,
	SUM(CASE WHEN d.IMPOSTO = 'PIS' THEN d.VALOR END) AS PIS_AMOUNT,
	SUM(CASE WHEN d.IMPOSTO = 'COFINS' THEN d.VALOR END) AS COFINS_AMOUNT
FROM LEANDRO.VW_IMPOSTO_ITEM d
JOIN METALPRD.TGFITE i
  ON i.NUNOTA = d.NUNOTA
 AND i.SEQUENCIA = d.SEQUENCIA
WHERE i.CODPROD = :1
  AND d.CODINC = :2
  AND TRUNC(d.DTNEG) = TRUNC(:3)`

	var (
		icms   sql.NullFloat64
		ipi    sql.NullFloat64
		pis    sql.NullFloat64
		cofins sql.NullFloat64
	)
	if err := r.db.QueryRowContext(ctx, query, input.ProductID, policy.IncidenceCode, policy.EffectiveAt).Scan(&icms, &ipi, &pis, &cofins); err != nil {
		return domain.TaxInputs{}, wrapOracleError("read tax inputs", err)
	}

	result := domain.TaxInputs{
		ProductID:     input.ProductID,
		EffectiveAt:   policy.EffectiveAt,
		IncidenceCode: policy.IncidenceCode,
		ICMSAmount:    nullableFloat(icms),
		IPIAmount:     nullableFloat(ipi),
		PISAmount:     nullableFloat(pis),
		COFINSAmount:  nullableFloat(cofins),
		Source:        domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
	}
	if result.ICMSAmount == nil && result.IPIAmount == nil && result.PISAmount == nil && result.COFINSAmount == nil {
		result.QualityFlags = []domain.QualityFlag{domain.QualityMissingTax}
	}
	return result, nil
}

func (r *Reader) ensureAvailable(ctx context.Context) error {
	if r == nil || r.db == nil {
		return domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle reader is not configured", nil)
	}
	if err := r.db.PingContext(ctx); err != nil {
		return wrapOracleError("ping oracle", err)
	}
	return nil
}

func buildDSN(cfg Config) string {
	parts := []string{
		fmt.Sprintf("user=%q", cfg.Username),
		fmt.Sprintf("password=%q", cfg.Password),
		fmt.Sprintf("connectString=%q", cfg.ConnectString),
		fmt.Sprintf("poolMinSessions=%d", cfg.PoolMinSessions),
		fmt.Sprintf("poolMaxSessions=%d", cfg.PoolMaxSessions),
		fmt.Sprintf("poolSessionTimeout=%s", cfg.SessionTimeout),
	}
	if cfg.LibDir != "" {
		parts = append(parts, fmt.Sprintf("libDir=%q", cfg.LibDir))
	}
	return strings.Join(parts, " ")
}

func buildFindProductsQuery(input ports.FindProductsInput) (string, []any, error) {
	query := `
SELECT
	p.CODPROD,
	p.DESCRPROD,
	p.REFERENCIA,
	p.CODGRUPOPROD,
	p.MARCA,
	p.CODMARCA,
	p.ATIVO,
	p.USOPROD
FROM METALPRD.TGFPRO p
WHERE 1 = 1`

	args := make([]any, 0, 3)
	if !input.IncludeInactive {
		query += " AND p.ATIVO = 'S'"
	}

	var matchClauses []string
	if value := trimmedPointer(input.EAN); value != nil {
		args = append(args, *value)
		matchClauses = append(matchClauses, fmt.Sprintf("p.REFERENCIA = :%d", len(args)))
	}
	if value := trimmedPointer(input.SellerSKU); value != nil {
		args = append(args, *value)
		matchClauses = append(matchClauses, fmt.Sprintf("p.REFFORN = :%d", len(args)))
	}
	if value := trimmedPointer(input.Title); value != nil {
		args = append(args, "%"+strings.ToUpper(*value)+"%")
		matchClauses = append(matchClauses, fmt.Sprintf("UPPER(p.DESCRPROD) LIKE :%d", len(args)))
	}
	if len(matchClauses) == 0 {
		return "", nil, domain.NewReadError(domain.ReadErrorUnsupportedQuery, "oracle product lookup requires at least one linking key", nil)
	}

	query += " AND (" + strings.Join(matchClauses, " OR ") + ") ORDER BY p.ATIVO DESC, p.CODPROD"
	return query, args, nil
}

func buildSellableStockQuery(productID int, policy domain.SellableStockPolicy) (string, []any, error) {
	if len(policy.CompanyIDs) == 0 || len(policy.LocationIDs) == 0 {
		return "", nil, domain.NewReadError(domain.ReadErrorUnsupportedQuery, "oracle sellable stock policy requires explicit company and location scopes", nil)
	}

	args := []any{productID}
	query := `
SELECT
	CAST(SUM(CASE WHEN est.CODPARC = 0 THEN NVL(est.ESTOQUE, 0) - NVL(est.RESERVADO, 0) ELSE 0 END) AS NUMBER(38,2)) AS SELLABLE_QTY
FROM METALPRD.TGFEST est
WHERE est.CODPROD = :1`

	query += buildIntListClause("est.CODEMP", policy.CompanyIDs, &args)
	query += buildIntListClause("est.CODLOCAL", policy.LocationIDs, &args)
	if len(policy.ExcludedLocationIDs) > 0 {
		query += buildNotIntListClause("est.CODLOCAL", policy.ExcludedLocationIDs, &args)
	}
	return query, args, nil
}

func buildSalesHistoryQuery(input ports.SalesHistoryInput) (string, []any, error) {
	if input.ProductID == nil && input.ProductGroupID == nil {
		return "", nil, domain.NewReadError(domain.ReadErrorUnsupportedQuery, "oracle sales history requires product_id or product_group_id", nil)
	}

	query := `
SELECT
	v.DTNEG,
	v.CODPROD,
	v.CODGRUPOPROD,
	v.QTD_FATURADA,
	v.VLR_FATURADO,
	v.CODTAB,
	v.CODTIPOPER,
	v.TIPO_LINHA
FROM LEANDRO.VW_FAT_VENDA_ITEM v
WHERE v.DTNEG >= :1
  AND v.DTNEG < :2`
	args := []any{input.Window.Start, input.Window.End}

	if input.ProductID != nil {
		args = append(args, *input.ProductID)
		query += fmt.Sprintf(" AND v.CODPROD = :%d", len(args))
	}
	if input.ProductGroupID != nil {
		args = append(args, *input.ProductGroupID)
		query += fmt.Sprintf(" AND v.CODGRUPOPROD = :%d", len(args))
	}
	query += " ORDER BY v.DTNEG, v.NUNOTA, v.SEQUENCIA"
	return query, args, nil
}

func buildIntListClause(column string, values []int, args *[]any) string {
	placeholders := make([]string, 0, len(values))
	for _, value := range values {
		*args = append(*args, value)
		placeholders = append(placeholders, fmt.Sprintf(":%d", len(*args)))
	}
	return fmt.Sprintf(" AND %s IN (%s)", column, strings.Join(placeholders, ", "))
}

func buildNotIntListClause(column string, values []int, args *[]any) string {
	placeholders := make([]string, 0, len(values))
	for _, value := range values {
		*args = append(*args, value)
		placeholders = append(placeholders, fmt.Sprintf(":%d", len(*args)))
	}
	return fmt.Sprintf(" AND %s NOT IN (%s)", column, strings.Join(placeholders, ", "))
}

func wrapOracleError(operation string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.NewReadError(domain.ReadErrorUnsupportedQuery, fmt.Sprintf("oracle %s returned no rows", operation), nil)
	}
	return domain.NewReadError(domain.ReadErrorSourceUnavailable, fmt.Sprintf("oracle %s failed", operation), err)
}

func trimmedPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	v := strings.TrimSpace(value.String)
	if v == "" {
		return nil
	}
	return &v
}

func nullableInt(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int64)
	return &v
}

func nullableInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	v := value.Int64
	return &v
}

func nullableFloat(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	v := value.Float64
	return &v
}

func nullableTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	v := value.Time
	return &v
}
