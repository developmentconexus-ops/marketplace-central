package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// errNoMatchableInput means the caller's input carried nothing this query can
// filter on, so the honest answer is an empty candidate list. It is internal:
// it never escapes FindProductsForLinking, and it exists only so the query
// builder can refuse to run an unfiltered table scan.
var errNoMatchableInput = errors.New("no matchable product lookup input")

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
	if errors.Is(err, errNoMatchableInput) {
		return nil, nil
	}
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
			eanValue       sql.NullString
			referenceValue sql.NullString
			ncmValue       sql.NullString
			productGroupID sql.NullInt64
			brandName      sql.NullString
			brandID        sql.NullInt64
			activeValue    string
			usageType      sql.NullString
			activeEANCount sql.NullInt64
		)
		if err := rows.Scan(&productID, &name, &eanValue, &referenceValue, &ncmValue, &productGroupID, &brandName, &brandID, &activeValue, &usageType, &activeEANCount); err != nil {
			return nil, wrapOracleError("scan product candidate", err)
		}

		ean := nullableString(eanValue)
		qualityFlags := []domain.QualityFlag{domain.QualityComplete}
		if ean == nil || !catalogdomain.IsValidGTIN(*ean) {
			ean = nil
			qualityFlags = []domain.QualityFlag{domain.QualityInvalidEAN}
		} else if activeEANCount.Valid && activeEANCount.Int64 >= 2 {
			qualityFlags = append(qualityFlags, domain.QualityEANCollision)
		}

		candidate := domain.ProductCandidate{
			InternalProductID: canonicalProductID(productID),
			ProductID:         productID,
			Name:              name,
			EAN:               ean,
			ReferenceCode:     nullableString(referenceValue),
			NCM:               nullableNCM(ncmValue),
			BrandID:           nullableInt(brandID),
			BrandName:         nullableString(brandName),
			ProductGroupID:    nullableInt(productGroupID),
			IsActive:          strings.EqualFold(activeValue, "S"),
			UsageType:         nullableString(usageType),
			Source: domain.SourceMetadata{
				System:    "oracle",
				FetchedAt: r.now().UTC(),
			},
			QualityFlags: qualityFlags,
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
			if !domain.HasQualityFlag(candidates[i].QualityFlags, domain.QualityAmbiguousProduct) {
				candidates[i].QualityFlags = append(candidates[i].QualityFlags, domain.QualityAmbiguousProduct)
			}
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
		Quantity:  knownOracleStock(quantity),
		Policy:    policy,
		Source: domain.SourceMetadata{
			System:    "oracle",
			FetchedAt: r.now().UTC(),
		},
	}
	return result, nil
}

func knownOracleStock(quantity sql.NullFloat64) *float64 {
	if value := nullableFloat(quantity); value != nil {
		return value
	}
	zero := 0.0
	return &zero
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
	policy := input.Policy
	if !policy.Source.Verified() {
		return domain.TaxInputs{
			ProductID: input.ProductID, EffectiveAt: policy.EffectiveAt,
			IncidenceCode: policy.IncidenceCode,
			QualityFlags:  []domain.QualityFlag{domain.QualityMissingTax},
		}, nil
	}
	if err := r.ensureAvailable(ctx); err != nil {
		return domain.TaxInputs{}, err
	}

	query, args := buildTaxInputsQuery(input.ProductID, policy)

	var (
		icms   sql.NullFloat64
		ipi    sql.NullFloat64
		pis    sql.NullFloat64
		cofins sql.NullFloat64
	)
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&icms, &ipi, &pis, &cofins); err != nil {
		return domain.TaxInputs{}, wrapOracleError("read tax inputs", err)
	}

	result := domain.TaxInputs{
		ProductID:      input.ProductID,
		EffectiveAt:    policy.EffectiveAt,
		IncidenceCode:  policy.IncidenceCode,
		SourceIdentity: policy.Source,
		ICMSAmount:     nullableFloat(icms),
		IPIAmount:      nullableFloat(ipi),
		PISAmount:      nullableFloat(pis),
		COFINSAmount:   nullableFloat(cofins),
		Source:         domain.SourceMetadata{System: "oracle", FetchedAt: r.now().UTC()},
	}
	if result.ICMSAmount == nil && result.IPIAmount == nil && result.PISAmount == nil && result.COFINSAmount == nil {
		result.QualityFlags = []domain.QualityFlag{domain.QualityMissingTax}
	}
	return result, nil
}

func buildTaxInputsQuery(productID int, policy domain.TaxPolicy) (string, []any) {
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
WHERE d.NUNOTA = :1
  AND d.SEQUENCIA = :2
  AND i.CODPROD = :3
  AND d.CODINC = :4`
	return query, []any{policy.Source.DocumentID, policy.Source.LineNumber, productID, policy.IncidenceCode}
}

func (r *Reader) ensureAvailable(ctx context.Context) error {
	if r == nil || r.db == nil {
		return domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle reader is not configured", nil)
	}
	return nil
}

func buildFindProductsQuery(input ports.FindProductsInput) (string, []any, error) {
	query := `
SELECT
	p.CODPROD,
	p.DESCRPROD,
	p.REFERENCIA,
	p.REFFORN,
	p.NCM,
	p.CODGRUPOPROD,
	p.MARCA,
	p.CODMARCA,
	p.ATIVO,
	p.USOPROD,
	(
		SELECT COUNT(DISTINCT collision.CODPROD)
		FROM METALPRD.TGFPRO collision
		WHERE collision.ATIVO = 'S'
		  AND TRIM(collision.REFERENCIA) = TRIM(p.REFERENCIA)
	) AS EAN_ACTIVE_COUNT
FROM METALPRD.TGFPRO p
WHERE 1 = 1`

	args := make([]any, 0, 3)
	if !input.IncludeInactive {
		query += " AND p.ATIVO = 'S'"
	}

	var matchClauses []string
	if input.ProductID != nil && *input.ProductID > 0 {
		args = append(args, *input.ProductID)
		matchClauses = append(matchClauses, fmt.Sprintf("p.CODPROD = :%d", len(args)))
	}
	// TGFPRO.REFERENCIA IS the governed barcode source: this same query already
	// projects it as the candidate's EAN and counts its active collisions, and
	// the mirror sync (sankhyaBaseSQL) writes products_mirror.ean from it under
	// exactly the GTIN shape rule applied here. Matching on it is therefore the
	// same fact the reader already publishes, not a new inference. A value that
	// is not GTIN-shaped is not an EAN and never becomes a match clause, so a
	// junk barcode can never widen the query.
	if value := trimmedPointer(input.EAN); value != nil && catalogdomain.IsValidGTIN(*value) {
		args = append(args, *value)
		matchClauses = append(matchClauses, fmt.Sprintf("TRIM(p.REFERENCIA) = :%d", len(args)))
	}
	// The provider's seller_sku carries the ERP CODPROD, not REFFORN (D-121,
	// operator-ratified): on Mercado Livre the operator registers the CODPROD in
	// the listing's SKU field, so matching REFFORN — the manufacturer reference —
	// meant the strong anchor never fired and every link fell back to the EAN.
	// A seller_sku that is not CODPROD-shaped is not a CODPROD and never becomes
	// a match clause, exactly as a non-GTIN barcode never becomes one above.
	if value := trimmedPointer(input.SellerSKU); value != nil && domain.IsValidCodprod(*value) {
		args = append(args, *value)
		matchClauses = append(matchClauses, fmt.Sprintf("p.CODPROD = :%d", len(args)))
	}
	if value := trimmedPointer(input.Title); value != nil {
		args = append(args, "%"+strings.ToUpper(*value)+"%")
		matchClauses = append(matchClauses, fmt.Sprintf("UPPER(p.DESCRPROD) LIKE :%d", len(args)))
	}
	if len(matchClauses) == 0 && (trimmedPointer(input.EAN) != nil || trimmedPointer(input.SellerSKU) != nil) {
		// Every anchor the caller supplied was malformed (a barcode that is not
		// GTIN-shaped, a seller_sku that is not CODPROD-shaped, or both). That is
		// a no-match, not a failed lookup: refusing here aborted the caller's whole
		// linking run over one malformed provider identifier. Falling through
		// instead would run the query with NO match clause at all and return the
		// entire catalog as candidates.
		return "", nil, errNoMatchableInput
	}
	if len(matchClauses) > 0 {
		query += " AND (" + strings.Join(matchClauses, " OR ") + ")"
	}
	query += " ORDER BY p.DESCRPROD, p.CODPROD"
	return query, args, nil
}

func canonicalProductID(productID int) *domain.InternalProductID {
	id, err := domain.NewInternalProductID(productID)
	if err != nil {
		return nil
	}
	return &id
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
	// Exclusions may be empty, so absence must emit no SQL. Inclusion lists are
	// deliberately asymmetric: callers reject empty whitelists. A future
	// caller-supplied batch policy must preserve that edge check.
	if len(values) == 0 {
		return ""
	}
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
	return domain.NewReadError(domain.ReadErrorSourceUnavailable, fmt.Sprintf("oracle %s failed", operation), SafeOracleCause(err))
}

// SafeOracleCause strips raw Oracle/OCI driver error text before it reaches a
// log or caller. It is exported because it is also the composition root's
// mitigation of record for the same class of leak outside this package (see
// apps/server_core/cmd/catalogingest/oracle_cgo.go, which mirrors this
// package's cgo/no-cgo split and must not let a raw driver error past its own
// ping check either).
func SafeOracleCause(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	var coded interface{ Code() int }
	if errors.As(err, &coded) {
		return fmt.Errorf("oracle error code=%d", coded.Code())
	}
	return errors.New("oracle driver error")
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

func nullableNCM(value sql.NullString) *string {
	ncm := nullableString(value)
	if ncm == nil || len(*ncm) != 8 {
		return nil
	}
	for _, digit := range *ncm {
		if digit < '0' || digit > '9' {
			return nil
		}
	}
	return ncm
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
