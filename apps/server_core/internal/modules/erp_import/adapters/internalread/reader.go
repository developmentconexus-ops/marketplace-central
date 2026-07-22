package internalread

import (
	"context"
	"errors"
	"fmt"
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

// WithActiveSource pins the source selected by M-02 routing from the tenant's
// active_source config. Without a pin, the reader fails closed.
func WithActiveSource(ctx context.Context, source erpdomain.ImportSource) context.Context {
	return context.WithValue(ctx, activeSourceKey{}, source)
}

// ActiveSourceFromContext reports the pinned dataset source and whether one was set.
func ActiveSourceFromContext(ctx context.Context) (erpdomain.ImportSource, bool) {
	if source, ok := ctx.Value(activeSourceKey{}).(erpdomain.ImportSource); ok && source != "" {
		return source, true
	}
	return "", false
}

func activeSourceFromContext(ctx context.Context) (erpdomain.ImportSource, error) {
	if source, ok := ActiveSourceFromContext(ctx); ok {
		return source, nil
	}
	return "", ErrUnknownActiveSource
}

// ErrUnknownActiveSource is returned by ParseActiveSource for a non-empty value
// that names no known dataset. Transport maps it to 400 — never a silent
// fallback to the default (which would hide a client bug behind ERP data).
var ErrUnknownActiveSource = errors.New("unknown_erp_source")

// ParseActiveSource maps a transport erp_source selector to an ImportSource.
// Empty ("") means the caller made no selection: present=false. A known value
// returns present=true. Anything else is ErrUnknownActiveSource.
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

func (r *Reader) sourceFor(ctx context.Context, message string) (erpdomain.ImportSource, error) {
	if r == nil || r.repo == nil {
		return "", readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, message, ErrNoErpSnapshot)
	}
	source, err := activeSourceFromContext(ctx)
	if err != nil {
		return "", readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, message, err)
	}
	return source, nil
}

func (r *Reader) FindProductsForLinking(ctx context.Context, input readports.FindProductsInput) ([]readdomain.ProductCandidate, error) {
	source, err := r.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return nil, err
	}
	rows, err := r.repo.MirrorRows(ctx, r.tenantID, source)
	if err != nil {
		return nil, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}
	collisions := validEANCounts(rows)
	results := make([]readdomain.ProductCandidate, 0)
	for _, row := range rows {
		if !matches(row, input) {
			continue
		}
		id, parseErr := strconv.ParseInt(strings.TrimSpace(row.CodigoProduto), 10, 0)
		if parseErr != nil || id <= 0 {
			// CODPROD is not representable as an internal product id; skip the row
			// (mirrors catalogPage) rather than aborting the whole call or leaking its value.
			continue
		}
		candidate, err := r.candidate(int(id), row, collisions, source)
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
	source, err := r.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	row, found, err := r.repo.MirrorProductByCode(ctx, r.tenantID, source, strconv.Itoa(input.ProductID))
	if err != nil {
		return readdomain.SellableStock{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}
	if !found {
		return readdomain.SellableStock{}, &ERPProductNotFoundError{ProductID: input.ProductID}
	}
	result := readdomain.SellableStock{ProductID: input.ProductID, Policy: input.Policy, Source: r.source(source, row.UpdatedAt)}
	if row.EstoqueTotal == nil {
		result.QualityFlags = []readdomain.QualityFlag{readdomain.QualityMissingStock}
		return result, nil
	}
	quantity, err := parseNumber(*row.EstoqueTotal, "estoque_total")
	if err != nil {
		return readdomain.SellableStock{}, err
	}
	result.Quantity = &quantity
	return result, nil
}

func (r *Reader) GetCurrentPrice(context.Context, readports.CurrentPriceInput) (readdomain.CurrentPrice, error) {
	return readdomain.CurrentPrice{}, readdomain.NewReadError(readdomain.ReadErrorUnsupportedQuery, "xlsx snapshot does not contain current price", nil)
}

func (r *Reader) GetCostAsOf(ctx context.Context, input readports.CostAsOfInput) (readdomain.CostAsOf, error) {
	source, err := r.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return readdomain.CostAsOf{}, err
	}
	row, found, err := r.repo.MirrorProductByCode(ctx, r.tenantID, source, strconv.Itoa(input.ProductID))
	if err != nil {
		return readdomain.CostAsOf{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}
	if !found {
		return readdomain.CostAsOf{}, &ERPProductNotFoundError{ProductID: input.ProductID}
	}
	if row.UpdatedAt.After(input.Policy.EffectiveAt) {
		return readdomain.CostAsOf{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "mirror row is newer than requested cost as-of", ErrNoErpSnapshot)
	}
	if row.Custo == nil {
		return readdomain.CostAsOf{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "custo is unknown for this product", nil)
	}
	amount, err := parseNumber(*row.Custo, "custo")
	if err != nil {
		return readdomain.CostAsOf{}, err
	}
	return readdomain.CostAsOf{ProductID: input.ProductID, CompanyID: input.Policy.CompanyID, Basis: input.Policy.Basis, EffectiveAt: input.Policy.EffectiveAt, Amount: &amount, AmountScope: readdomain.CostAmountScopePerUnit, Source: r.source(source, row.UpdatedAt)}, nil
}

func (r *Reader) GetSalesHistory(context.Context, readports.SalesHistoryInput) (readdomain.SalesHistory, error) {
	return readdomain.SalesHistory{}, readdomain.NewReadError(readdomain.ReadErrorUnsupportedQuery, "xlsx snapshot does not contain sales history", nil)
}

func (r *Reader) GetTaxInputs(ctx context.Context, input readports.TaxInput) (readdomain.TaxInputs, error) {
	source, err := r.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return readdomain.TaxInputs{}, err
	}
	row, found, err := r.repo.MirrorProductByCode(ctx, r.tenantID, source, strconv.Itoa(input.ProductID))
	if err != nil {
		return readdomain.TaxInputs{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}
	if !found {
		return readdomain.TaxInputs{}, &ERPProductNotFoundError{ProductID: input.ProductID}
	}
	result := readdomain.TaxInputs{ProductID: input.ProductID, EffectiveAt: input.Policy.EffectiveAt, IncidenceCode: input.Policy.IncidenceCode, SourceIdentity: input.Policy.Source, Source: r.source(source, row.UpdatedAt)}
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
	source, err := r.sourceFor(ctx, "active ERP source is unavailable")
	if err != nil {
		return readports.CatalogFactPage{}, err
	}
	rows, err := r.repo.MirrorCatalogPage(ctx, r.tenantID, source, query, cursor.InternalProductID, limit)
	if err != nil {
		return readports.CatalogFactPage{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}
	hasMore := limit > 0 && len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items := make([]readports.CatalogProductFact, 0, len(rows))
	collisions, err := r.repo.MirrorEANCollisionCounts(ctx, r.tenantID, source)
	if err != nil {
		return readports.CatalogFactPage{}, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "EAN collision data is unavailable", err)
	}
	var asOf time.Time
	var lastID int64
	for _, row := range rows {
		id, parseErr := strconv.ParseInt(strings.TrimSpace(row.CodigoProduto), 10, 64)
		if parseErr != nil || id <= 0 {
			return readports.CatalogFactPage{}, fmt.Errorf("parse codigo_produto: %w", parseErr)
		}
		fact, buildErr := catalogFact(id, row, collisions)
		if buildErr != nil {
			return readports.CatalogFactPage{}, buildErr
		}
		items = append(items, fact)
		lastID = id
		if row.UpdatedAt.After(asOf) {
			asOf = row.UpdatedAt
		}
	}
	page := readports.CatalogFactPage{Items: items, AsOf: asOf}
	if hasMore {
		page.NextCursor = &readports.Cursor{InternalProductID: lastID}
	}
	return page, nil
}

func (r *Reader) candidate(id int, row erpdomain.MirrorProduct, collisions map[string]int, source erpdomain.ImportSource) (readdomain.ProductCandidate, error) {
	canonical, err := readdomain.NewInternalProductID(id)
	if err != nil {
		return readdomain.ProductCandidate{}, fmt.Errorf("canonical codprod: %w", err)
	}
	ean, flags := identityQuality(row.EAN, collisions)
	name := ""
	if row.Descricao != nil {
		name = strings.TrimSpace(*row.Descricao)
	}
	return readdomain.ProductCandidate{InternalProductID: &canonical, ProductID: id, Name: name, EAN: ean, ReferenceCode: copyTrimmed(row.Referencia), NCM: copyTrimmed(row.NCM), BrandName: copyTrimmed(row.Marca), IsActive: true, Source: r.source(source, row.UpdatedAt), QualityFlags: flags}, nil
}

func catalogFact(id int64, row erpdomain.MirrorProduct, collisions map[string]int) (readports.CatalogProductFact, error) {
	ean, identityFlags := identityQuality(row.EAN, collisions)
	quality := make([]string, len(identityFlags))
	for i, flag := range identityFlags {
		quality[i] = string(flag)
	}
	fact := readports.CatalogProductFact{InternalProductID: id, Reference: copyTrimmed(row.Referencia), ManufacturerReference: copyTrimmed(row.Referencia), Description: copyTrimmed(row.Descricao), EAN: ean, BrandName: copyTrimmed(row.Marca), NCM: copyTrimmed(row.NCM), QualityFlags: quality, Active: true, CurrentPrice: readports.CatalogMoneyFact{Currency: "BRL", Quality: []string{string(readdomain.QualityMissingPrice)}}, Cost: readports.CatalogMoneyFact{Currency: "BRL"}}
	if row.EstoqueTotal == nil {
		fact.SellableStock.Quality = []string{string(readdomain.QualityMissingStock)}
	} else {
		quantity, err := parseNumber(*row.EstoqueTotal, "estoque_total")
		if err != nil {
			return fact, err
		}
		fact.SellableStock.Quantity = &quantity
	}
	if row.Custo == nil {
		fact.Cost.Quality = []string{string(readdomain.QualityMissingCost)}
	} else if _, err := parseNumber(*row.Custo, "custo"); err != nil {
		return fact, err
	} else {
		cost := strings.TrimSpace(*row.Custo)
		fact.Cost.Amount = &cost
	}
	return fact, nil
}

func (r *Reader) source(source erpdomain.ImportSource, observed time.Time) readdomain.SourceMetadata {
	t := observed
	return readdomain.SourceMetadata{System: string(source), FetchedAt: r.now().UTC(), ObservedAt: &t}
}
func matches(row erpdomain.MirrorProduct, input readports.FindProductsInput) bool {
	if input.ProductID != nil && strings.TrimSpace(row.CodigoProduto) == strconv.Itoa(*input.ProductID) {
		return true
	}
	if sku := trimmed(input.SellerSKU); sku != nil && strings.TrimSpace(row.CodigoProduto) == *sku {
		return true
	}
	if ean := trimmed(input.EAN); ean != nil && row.EAN != nil && catalogdomain.IsValidGTIN(*row.EAN) && strings.TrimSpace(*row.EAN) == *ean {
		return true
	}
	if title := trimmed(input.Title); title != nil && row.Descricao != nil && strings.Contains(strings.ToLower(*row.Descricao), strings.ToLower(*title)) {
		return true
	}
	return input.ProductID == nil && trimmed(input.SellerSKU) == nil && trimmed(input.EAN) == nil && trimmed(input.Title) == nil
}
func validEANCounts(rows []erpdomain.MirrorProduct) map[string]int {
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
