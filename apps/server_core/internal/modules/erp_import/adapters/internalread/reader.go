package internalread

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	erpports "marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var ErrNoErpSnapshot = errors.New("no_erp_snapshot")

// activeSourceKey carries the caller-selected dataset source through the request
// context. The two-source toggle (xlsx = Sankhya ERP, catalogo_cliente = prospect
// catalog) is resolved here so every reader method observes the same dataset.
type activeSourceKey struct{}

// WithActiveSource pins the dataset source for the reader on this request. Absent
// it, the reader defaults to the ERP (xlsx) snapshot so the demo opens on real data.
func WithActiveSource(ctx context.Context, source erpdomain.ImportSource) context.Context {
	return context.WithValue(ctx, activeSourceKey{}, source)
}

// ActiveSourceFromContext reports the pinned dataset source and whether one was
// set. Absent (present=false) the reader defaults to xlsx; the presence flag lets
// transport tests confirm the toggle threaded through without asserting on the
// default.
func ActiveSourceFromContext(ctx context.Context) (erpdomain.ImportSource, bool) {
	if source, ok := ctx.Value(activeSourceKey{}).(erpdomain.ImportSource); ok && source != "" {
		return source, true
	}
	return "", false
}

func activeSourceFromContext(ctx context.Context) erpdomain.ImportSource {
	if source, ok := ActiveSourceFromContext(ctx); ok {
		return source
	}
	return erpdomain.SourceXLSX
}

// ErrUnknownActiveSource is returned by ParseActiveSource for a non-empty value
// that names no known dataset. Transport maps it to 400 — never a silent
// fallback to the default (which would hide a client bug behind ERP data).
var ErrUnknownActiveSource = errors.New("unknown_erp_source")

// ParseActiveSource maps a transport erp_source selector to an ImportSource.
// Empty ("") means the caller made no selection: present=false, and the reader
// keeps its xlsx default (so an absent param is byte-stable with prior clients).
// A known value returns present=true. Anything else is ErrUnknownActiveSource.
func ParseActiveSource(raw string) (erpdomain.ImportSource, bool, error) {
	switch strings.TrimSpace(raw) {
	case "":
		return "", false, nil
	case string(erpdomain.SourceXLSX):
		return erpdomain.SourceXLSX, true, nil
	case string(erpdomain.SourceCatalogoCliente):
		return erpdomain.SourceCatalogoCliente, true, nil
	default:
		return "", false, ErrUnknownActiveSource
	}
}

type ERPProductNotFoundError struct{ ProductID int }

func (e *ERPProductNotFoundError) Error() string {
	return fmt.Sprintf("ERP product %d not found", e.ProductID)
}

type Option func(*Reader)

func WithClock(now func() time.Time) Option {
	return func(r *Reader) { r.now = now }
}

type Reader struct {
	repo     erpports.ImportRepository
	tenantID string
	now      func() time.Time
}

func NewReader(repo erpports.ImportRepository, tenantID string, opts ...Option) *Reader {
	r := &Reader{repo: repo, tenantID: tenantID, now: time.Now}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Reader) snapshot(ctx context.Context, message string) (erpdomain.ImportSnapshot, error) {
	if r == nil || r.repo == nil {
		return erpdomain.ImportSnapshot{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, message, ErrNoErpSnapshot)
	}
	snapshot, err := r.repo.LatestCompletedSnapshot(ctx, r.tenantID, activeSourceFromContext(ctx))
	if err != nil || snapshot.Status != erpdomain.ImportStatusCompleted {
		return erpdomain.ImportSnapshot{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, message, ErrNoErpSnapshot)
	}
	return snapshot, nil
}

func (r *Reader) FindProductsForLinking(ctx context.Context, input readports.FindProductsInput) ([]readdomain.ProductCandidate, error) {
	snapshot, err := r.snapshot(ctx, "xlsx completed snapshot is unavailable")
	if err != nil {
		return nil, err
	}
	collisions := validEANCounts(snapshot.AcceptedRows)
	results := make([]readdomain.ProductCandidate, 0)
	for _, row := range snapshot.AcceptedRows {
		if !matches(row, input) {
			continue
		}
		id, parseErr := strconv.ParseInt(strings.TrimSpace(row.Codprod), 10, 0)
		if parseErr != nil || id <= 0 {
			// CODPROD is not representable as an internal product id; skip the row
			// (mirrors catalogPage) rather than aborting the whole call or leaking its value.
			continue
		}
		candidate, err := r.candidate(int(id), row, collisions)
		if err != nil {
			return nil, err
		}
		results = append(results, candidate)
	}
	if len(results) == 0 && (input.ProductID != nil || trimmed(input.SellerSKU) != nil) {
		productID := 0
		if input.ProductID != nil {
			productID = *input.ProductID
		}
		return nil, &ERPProductNotFoundError{ProductID: productID}
	}
	if len(results) > 1 {
		for i := range results {
			if !readdomain.HasQualityFlag(results[i].QualityFlags, readdomain.QualityAmbiguousProduct) {
				results[i].QualityFlags = append(results[i].QualityFlags, readdomain.QualityAmbiguousProduct)
			}
		}
	}
	return results, nil
}

func (r *Reader) GetSellableStock(ctx context.Context, input readports.SellableStockInput) (readdomain.SellableStock, error) {
	snapshot, err := r.snapshot(ctx, "xlsx completed snapshot is unavailable")
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	row, err := rowByProductID(snapshot.AcceptedRows, input.ProductID)
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	result := readdomain.SellableStock{ProductID: input.ProductID, Policy: input.Policy, Source: r.source(snapshot.ImportedAt)}
	if row.StockReserved == nil || strings.TrimSpace(row.StockPhysical) == "" {
		result.QualityFlags = []readdomain.QualityFlag{readdomain.QualityMissingStock}
		return result, nil
	}
	physical, err := parseNumber(row.StockPhysical, "stock_physical")
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	reserved, err := parseNumber(*row.StockReserved, "stock_reserved")
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	quantity := physical - reserved
	result.Quantity = &quantity
	return result, nil
}

func (r *Reader) GetCurrentPrice(context.Context, readports.CurrentPriceInput) (readdomain.CurrentPrice, error) {
	return readdomain.CurrentPrice{}, readdomain.NewReadError(readdomain.ReadErrorUnsupportedQuery, "xlsx snapshot does not contain current price", nil)
}

func (r *Reader) GetCostAsOf(ctx context.Context, input readports.CostAsOfInput) (readdomain.CostAsOf, error) {
	snapshot, err := r.snapshot(ctx, "xlsx completed snapshot is unavailable")
	if err != nil {
		return readdomain.CostAsOf{}, err
	}
	row, err := rowByProductID(snapshot.AcceptedRows, input.ProductID)
	if err != nil {
		return readdomain.CostAsOf{}, err
	}
	if snapshot.ImportedAt.After(input.Policy.EffectiveAt) {
		return readdomain.CostAsOf{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "xlsx snapshot is newer than requested cost as-of", ErrNoErpSnapshot)
	}
	custo := strings.TrimSpace(string(row.Custo))
	if custo == "" {
		// Client-catalog (lenient) imports honestly omit custo (ADR-17); report
		// the source as unavailable for this query rather than fabricating a cost.
		return readdomain.CostAsOf{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "custo is unknown for this product", nil)
	}
	amount, err := parseNumber(custo, "custo")
	if err != nil {
		return readdomain.CostAsOf{}, err
	}
	return readdomain.CostAsOf{ProductID: input.ProductID, CompanyID: input.Policy.CompanyID, Basis: input.Policy.Basis, EffectiveAt: input.Policy.EffectiveAt, Amount: &amount, AmountScope: readdomain.CostAmountScopePerUnit, Source: r.source(snapshot.ImportedAt)}, nil
}

func (r *Reader) GetSalesHistory(context.Context, readports.SalesHistoryInput) (readdomain.SalesHistory, error) {
	return readdomain.SalesHistory{}, readdomain.NewReadError(readdomain.ReadErrorUnsupportedQuery, "xlsx snapshot does not contain sales history", nil)
}

func (r *Reader) GetTaxInputs(ctx context.Context, input readports.TaxInput) (readdomain.TaxInputs, error) {
	snapshot, err := r.snapshot(ctx, "xlsx completed snapshot is unavailable")
	if err != nil {
		return readdomain.TaxInputs{}, err
	}
	row, err := rowByProductID(snapshot.AcceptedRows, input.ProductID)
	if err != nil {
		return readdomain.TaxInputs{}, err
	}
	result := readdomain.TaxInputs{ProductID: input.ProductID, EffectiveAt: input.Policy.EffectiveAt, IncidenceCode: input.Policy.IncidenceCode, SourceIdentity: input.Policy.Source, Source: r.source(snapshot.ImportedAt)}
	if validNCM(row.NCM) {
		result.NCM = copyTrimmed(row.NCM)
	} else {
		result.QualityFlags = []readdomain.QualityFlag{readdomain.QualityMissingTax}
	}
	return result, nil
}

func (r *Reader) ListCatalogProductFacts(ctx context.Context, cursor readports.Cursor, limit int) (readports.CatalogFactPage, error) {
	return r.catalogPage(ctx, cursor, "", limit)
}

func (r *Reader) SearchCatalogProductFacts(ctx context.Context, query string, limit int) (readports.CatalogFactPage, error) {
	return r.catalogPage(ctx, readports.Cursor{}, query, limit)
}

func (r *Reader) catalogPage(ctx context.Context, cursor readports.Cursor, query string, limit int) (readports.CatalogFactPage, error) {
	snapshot, err := r.snapshot(ctx, "xlsx completed snapshot is unavailable")
	if err != nil {
		return readports.CatalogFactPage{}, err
	}
	type indexedRow struct {
		id  int64
		row erpdomain.NormalizedRow
	}
	rows := make([]indexedRow, 0, len(snapshot.AcceptedRows))
	for _, row := range snapshot.AcceptedRows {
		id, parseErr := strconv.ParseInt(strings.TrimSpace(row.Codprod), 10, 64)
		if parseErr != nil || id <= cursor.InternalProductID || id <= 0 || (query != "" && !strings.Contains(strings.ToLower(row.Descrprod), strings.ToLower(query))) {
			continue
		}
		rows = append(rows, indexedRow{id: id, row: row})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].id < rows[j].id })
	hasMore := limit > 0 && len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]readports.CatalogProductFact, 0, len(rows))
	collisions := validEANCounts(snapshot.AcceptedRows)
	for _, indexed := range rows {
		fact, buildErr := catalogFact(indexed.id, indexed.row, collisions)
		if buildErr != nil {
			return readports.CatalogFactPage{}, buildErr
		}
		items = append(items, fact)
	}
	page := readports.CatalogFactPage{Items: items, AsOf: snapshot.ImportedAt}
	if hasMore {
		page.NextCursor = &readports.Cursor{InternalProductID: rows[len(rows)-1].id}
	}
	return page, nil
}

func (r *Reader) candidate(id int, row erpdomain.NormalizedRow, collisions map[string]int) (readdomain.ProductCandidate, error) {
	canonical, err := readdomain.NewInternalProductID(id)
	if err != nil {
		return readdomain.ProductCandidate{}, fmt.Errorf("canonical codprod: %w", err)
	}
	ean, flags := identityQuality(row.EAN, collisions)
	return readdomain.ProductCandidate{InternalProductID: &canonical, ProductID: id, Name: row.Descrprod, EAN: ean, ReferenceCode: copyTrimmed(row.Refforn), NCM: copyTrimmed(row.NCM), BrandName: copyTrimmed(row.Marca), IsActive: true, Source: readdomain.SourceMetadata{System: "xlsx", FetchedAt: r.now().UTC()}, QualityFlags: flags}, nil
}

func catalogFact(id int64, row erpdomain.NormalizedRow, collisions map[string]int) (readports.CatalogProductFact, error) {
	ean, identityFlags := identityQuality(row.EAN, collisions)
	quality := make([]string, len(identityFlags))
	for i, flag := range identityFlags {
		quality[i] = string(flag)
	}
	fact := readports.CatalogProductFact{InternalProductID: id, Reference: copyTrimmed(row.Refforn), ManufacturerReference: copyTrimmed(row.Refforn), Description: copyTrimmed(&row.Descrprod), EAN: ean, BrandName: copyTrimmed(row.Marca), NCM: copyTrimmed(row.NCM), QualityFlags: quality, Active: true, CurrentPrice: readports.CatalogMoneyFact{Currency: "BRL", Quality: []string{string(readdomain.QualityMissingPrice)}}, Cost: readports.CatalogMoneyFact{Currency: "BRL"}}
	if row.StockReserved == nil || strings.TrimSpace(row.StockPhysical) == "" {
		fact.SellableStock.Quality = []string{string(readdomain.QualityMissingStock)}
	} else {
		physical, err := parseNumber(row.StockPhysical, "stock_physical")
		if err != nil {
			return fact, err
		}
		reserved, err := parseNumber(*row.StockReserved, "stock_reserved")
		if err != nil {
			return fact, err
		}
		quantity := physical - reserved
		fact.SellableStock.Quantity = &quantity
	}
	cost := strings.TrimSpace(string(row.Custo))
	if cost == "" {
		// Client-catalog (lenient) imports honestly omit custo (ADR-17); never
		// fabricate a cost value, just flag it as unknown.
		fact.Cost.Quality = []string{string(readdomain.QualityMissingCost)}
	} else if _, err := parseNumber(cost, "custo"); err != nil {
		return fact, err
	} else {
		fact.Cost.Amount = &cost
	}
	return fact, nil
}

func (r *Reader) source(observed time.Time) readdomain.SourceMetadata {
	t := observed
	return readdomain.SourceMetadata{System: "xlsx", FetchedAt: r.now().UTC(), ObservedAt: &t}
}
func rowByProductID(rows []erpdomain.NormalizedRow, id int) (erpdomain.NormalizedRow, error) {
	target := strconv.Itoa(id)
	for _, row := range rows {
		if strings.TrimSpace(row.Codprod) == target {
			return row, nil
		}
	}
	return erpdomain.NormalizedRow{}, &ERPProductNotFoundError{ProductID: id}
}
func matches(row erpdomain.NormalizedRow, input readports.FindProductsInput) bool {
	if input.ProductID != nil && strings.TrimSpace(row.Codprod) == strconv.Itoa(*input.ProductID) {
		return true
	}
	if sku := trimmed(input.SellerSKU); sku != nil && strings.TrimSpace(row.Codprod) == *sku {
		return true
	}
	if ean := trimmed(input.EAN); ean != nil && row.EAN != nil && catalogdomain.IsValidGTIN(*row.EAN) && strings.TrimSpace(*row.EAN) == *ean {
		return true
	}
	if title := trimmed(input.Title); title != nil && strings.Contains(strings.ToLower(row.Descrprod), strings.ToLower(*title)) {
		return true
	}
	return input.ProductID == nil && trimmed(input.SellerSKU) == nil && trimmed(input.EAN) == nil && trimmed(input.Title) == nil
}
func validEANCounts(rows []erpdomain.NormalizedRow) map[string]int {
	out := map[string]int{}
	for _, row := range rows {
		if row.EAN != nil {
			value := strings.TrimSpace(*row.EAN)
			if catalogdomain.IsValidGTIN(value) {
				out[value]++
			}
		}
	}
	return out
}
func identityQuality(value *string, collisions map[string]int) (*string, []readdomain.QualityFlag) {
	ean := copyTrimmed(value)
	if ean == nil || !catalogdomain.IsValidGTIN(*ean) {
		return nil, []readdomain.QualityFlag{readdomain.QualityInvalidEAN}
	}
	flags := []readdomain.QualityFlag{readdomain.QualityComplete}
	if collisions[*ean] >= 2 {
		flags = append(flags, readdomain.QualityEANCollision)
	}
	return ean, flags
}
func trimmed(value *string) *string { return copyTrimmed(value) }
func copyTrimmed(value *string) *string {
	if value == nil {
		return nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil
	}
	return &v
}
func validNCM(value *string) bool {
	if value == nil {
		return false
	}
	v := strings.TrimSpace(*value)
	if len(v) != 8 {
		return false
	}
	for _, c := range v {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func parseNumber(value, field string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", field, err)
	}
	return parsed, nil
}

var _ readports.Reader = (*Reader)(nil)
var _ readports.CatalogPageReader = (*Reader)(nil)
