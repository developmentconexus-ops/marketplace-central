package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/catalog/application"
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type Handler struct {
	Service              application.CanonicalService
	CompatibilityService application.Service
	PageReader           CatalogReader
}

const (
	catalogListPath     = "/catalog/products"
	catalogSearchPath   = "/catalog/products/search"
	catalogCountsPath   = "/catalog/products/counts"
	catalogDefaultLimit = 50
	catalogMaxLimit     = 100
	searchMaxLimit      = 50
)

type CatalogReader interface {
	ports.CatalogPageReader
}

// FreshnessPolicyFromContext exposes the transport-to-read freshness seam to
// the composition-owned reader without coupling transport callers to the
// context key implementation.
func FreshnessPolicyFromContext(ctx context.Context) (internalreaddomain.FreshnessPolicy, bool) {
	return internalreaddomain.FreshnessPolicyFromContext(ctx)
}

func requestContext(r *http.Request) (context.Context, error) {
	ctx := r.Context()
	// erp_source selects which imported snapshot the reader serves (xlsx = Sankhya
	// ERP, catalogo_cliente = prospect catalog). Absent = reader default (xlsx),
	// so prior clients stay byte-stable; an unknown value is a 400, never a silent
	// fallback to the default dataset.
	source, present, err := erpinternalread.ParseActiveSource(r.URL.Query().Get("erp_source"))
	if err != nil {
		return nil, &catalogPageQueryError{code: "invalid_erp_source", allowedRange: "xlsx|catalogo_cliente"}
	}
	if present {
		ctx = erpinternalread.WithActiveSource(ctx, source)
	}
	if hasNoCacheDirective(r.Header.Get("Cache-Control")) {
		ctx = internalreaddomain.WithFreshnessPolicy(ctx, internalreaddomain.FreshnessPolicy{MaxAge: 0})
	}
	return ctx, nil
}

func hasNoCacheDirective(header string) bool {
	for _, directive := range strings.Split(header, ",") {
		if strings.EqualFold(strings.TrimSpace(strings.SplitN(directive, "=", 2)[0]), "no-cache") {
			return true
		}
	}
	return false
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, map[string]any{
		"error": map[string]any{"code": code, "message": message, "details": map[string]any{}},
	})
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	// There is one catalog list route and one catalog search route, and both go
	// through the page reader. The registration used to fork: a mux that was not
	// a RouteClassMux, with no page reader, got a pair of legacy handlers that
	// read the WHOLE catalog with no assortment policy at all. Production never
	// took that fork, so the uncut read stayed alive on the strength of two unit
	// tests — and any future wiring that forgot the page reader would have been
	// served the uncut catalog silently instead of failing.
	registerInteractiveRoute(mux, catalogListPath, h.handleProducts)
	registerInteractiveRoute(mux, catalogSearchPath, h.handleSearch)
	registerInteractiveRoute(mux, catalogCountsPath, h.handleCounts)
	mux.HandleFunc("/catalog/taxonomy", h.handleTaxonomy)
	mux.HandleFunc("GET /catalog/products/{id}", h.handleGetProduct)
	mux.HandleFunc("GET /catalog/products/{id}/enrichment", h.handleGetEnrichment)
	mux.HandleFunc("PUT /catalog/products/{id}/enrichment", h.handleUpsertEnrichment)
}

type routeClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func registerInteractiveRoute(mux httpx.RouteRegistrar, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if registrar, ok := mux.(routeClassRegistrar); ok {
		registrar.RegisterRouteClass(pattern, httpx.InteractiveRouteClass)
	}
	mux.HandleFunc("GET "+pattern, handler)
}

func (h Handler) handleProducts(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if _, asked := r.URL.Query()["ids"]; asked {
		h.handleProductsByIDs(w, r, start)
		return
	}
	cursor, limit, err := parseCatalogPageQuery(r, catalogMaxLimit)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	if !h.requirePageReader(w) {
		return
	}
	ctx, err := requestContext(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	policy, err := catalogPolicy(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	page, err := h.PageReader.ListCatalogProductFacts(ctx, cursor, limit, policy)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	slog.Info("catalog.products", "action", "list", "result", "200", "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, newCatalogPageResponse(page, false))
}

// handleProductsByIDs serves ?ids=1,2,3. A screen that already knows which
// products it needs — the ones linked to listings, for instance — cannot reach
// them through the keyset cursor without reading the entire catalog in between,
// so it asks for them by id. Ids the active source does not carry come back
// absent, never as blank rows.
func (h Handler) handleProductsByIDs(w http.ResponseWriter, r *http.Request, start time.Time) {
	ids, err := parseCatalogIDs(r.URL.Query()["ids"])
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	if !h.requirePageReader(w) {
		return
	}
	ctx, err := requestContext(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	page, err := h.PageReader.CatalogProductFactsByIDs(ctx, ids)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	slog.Info("catalog.products", "action", "by_ids", "result", "200", "asked", len(ids), "found", len(page.Items), "duration_ms", time.Since(start).Milliseconds())
	// The response is the exact set asked for, not a page of a larger sequence:
	// emitting a next_cursor would invite a caller to page past the end of it.
	httpx.WriteJSON(w, http.StatusOK, newCatalogPageResponse(page, true))
}

func parseCatalogIDs(raw []string) ([]int64, error) {
	invalid := &catalogPageQueryError{code: "invalid_ids", allowedRange: "1.." + strconv.Itoa(catalogMaxLimit) + " positive integers"}
	if len(raw) != 1 {
		return nil, invalid
	}
	fields := strings.Split(raw[0], ",")
	ids := make([]int64, 0, len(fields))
	for _, field := range fields {
		value, convErr := strconv.ParseInt(strings.TrimSpace(field), 10, 64)
		if convErr != nil || value <= 0 {
			return nil, invalid
		}
		ids = append(ids, value)
	}
	if len(ids) == 0 || len(ids) > catalogMaxLimit {
		return nil, invalid
	}
	return ids, nil
}

func (h Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	q := r.URL.Query().Get("q")
	// q is required: true in contracts/api/marketplace-central.openapi.yaml. The
	// paged handler never enforced it — the check lived in the legacy handler
	// deleted with the uncut route — so a blank q reached the source as the empty
	// pattern and came back as the whole catalog, dressed as search results.
	if strings.TrimSpace(q) == "" {
		writeCatalogPageError(w, &catalogPageQueryError{code: "invalid_q"})
		return
	}
	// Search is a page of a longer sequence, exactly like the list read, so it
	// accepts the same cursor. Without it the response could only ever be the
	// first page and a caller had no way to reach the rest of the matches.
	cursor, limit, err := parseCatalogPageQuery(r, searchMaxLimit)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	if !h.requirePageReader(w) {
		return
	}
	ctx, err := requestContext(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	policy, err := catalogPolicy(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	page, err := h.PageReader.SearchCatalogProductFacts(ctx, q, cursor, limit, policy)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	slog.Info("catalog.search", "action", "search", "result", "200", "items", len(page.Items), "has_more", page.NextCursor != nil, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, newCatalogPageResponse(page, false))
}

// catalogPolicy answers with the only two things the wire can say. include_all
// is a boolean and it stays a boolean up to here: absent or false is the request
// "apply whatever rule this tenant configured", which crosses the port as nil
// and is answered ONCE at the routing seam; true names a concrete policy — the
// honest no-cut one — which the chain then applies verbatim. The transport never
// invents a rule of its own, which is why no default lives in this file.
func catalogPolicy(r *http.Request) (*ports.SellableAssortmentPolicy, error) {
	raw, present := r.URL.Query()["include_all"]
	if !present {
		return nil, nil
	}
	if len(raw) != 1 {
		return nil, &catalogPageQueryError{code: "invalid_include_all"}
	}
	value := strings.TrimSpace(raw[0])
	if value != "true" && value != "false" {
		return nil, &catalogPageQueryError{code: "invalid_include_all"}
	}
	if value == "false" {
		return nil, nil
	}
	all := ports.AllProductsAssortment()
	return &all, nil
}

func (h Handler) handleCounts(w http.ResponseWriter, r *http.Request) {
	if !h.requirePageReader(w) {
		return
	}
	ctx, err := requestContext(r)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	// The counter reports the tenant's rule against the whole catalog, so it asks
	// the routing seam for the stored policy rather than naming one here.
	counts, err := h.PageReader.GetCatalogAssortmentCounts(ctx, nil)
	if err != nil {
		writeCatalogPageError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]int{"sellable_count": counts.SellableCount, "total_count": counts.TotalCount})
}

type catalogPageQueryError struct {
	code         string
	allowedRange string
}

func (e *catalogPageQueryError) Error() string { return e.code }

func parseCatalogPageQuery(r *http.Request, maxLimit int) (ports.Cursor, int, error) {
	cursor := ports.Cursor{}
	if raw, ok := r.URL.Query()["cursor"]; ok {
		if len(raw) != 1 {
			return cursor, 0, &catalogPageQueryError{code: "invalid_cursor"}
		}
		decoded, err := ports.DecodeCursor(raw[0])
		if err != nil {
			return cursor, 0, &catalogPageQueryError{code: "invalid_cursor"}
		}
		cursor = decoded
	}
	limit, err := parseLimit(r, maxLimit)
	return cursor, limit, err
}

func parseLimit(r *http.Request, maxLimit int) (int, error) {
	query := r.URL.Query()
	raw, ok := query["limit"]
	if !ok {
		return catalogDefaultLimit, nil
	}
	if len(raw) != 1 {
		return 0, &catalogPageQueryError{code: "invalid_limit", allowedRange: "1.." + strconv.Itoa(maxLimit)}
	}
	limit, err := strconv.Atoi(strings.TrimSpace(raw[0]))
	if err != nil || limit < 1 || limit > maxLimit {
		return 0, &catalogPageQueryError{code: "invalid_limit", allowedRange: "1.." + strconv.Itoa(maxLimit)}
	}
	return limit, nil
}

type catalogPageEnvelope struct {
	Items      []catalogProductFactResponse `json:"items"`
	NextCursor *string                      `json:"next_cursor"`
	PageSize   int                          `json:"page_size"`
	AsOf       time.Time                    `json:"as_of"`
}

type catalogProductFactResponse struct {
	InternalProductID     int                         `json:"internal_product_id"`
	Reference             *string                     `json:"reference"`
	ManufacturerReference *string                     `json:"manufacturer_reference"`
	Description           *string                     `json:"description"`
	EAN                   *string                     `json:"ean"`
	BrandName             *string                     `json:"brand_name"`
	NCM                   *string                     `json:"ncm"`
	QualityFlags          []string                    `json:"quality_flags"`
	Active                bool                        `json:"active"`
	SellableStock         catalogQuantityFactResponse `json:"sellable_stock"`
	CurrentPrice          catalogMoneyFactResponse    `json:"current_price"`
	Cost                  catalogMoneyFactResponse    `json:"cost"`
}

type catalogQuantityFactResponse struct {
	Quantity *float64 `json:"quantity"`
	Quality  []string `json:"quality"`
}

type catalogMoneyFactResponse struct {
	Amount   *string  `json:"amount"`
	Currency string   `json:"currency"`
	Quality  []string `json:"quality"`
}

// suppressCursor is for reads that are NOT a page of a longer sequence — the
// explicit-id read answers exactly the set asked for, so a next_cursor would
// invite paging past the end of it. Paged reads (list and search) must emit
// theirs: withholding it makes a truncated page look complete.
func newCatalogPageResponse(page ports.CatalogFactPage, suppressCursor bool) catalogPageEnvelope {
	items := make([]catalogProductFactResponse, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, catalogProductFactResponse{
			InternalProductID:     int(item.InternalProductID),
			Reference:             item.Reference,
			ManufacturerReference: item.ManufacturerReference,
			Description:           item.Description,
			EAN:                   item.EAN,
			BrandName:             item.BrandName,
			NCM:                   item.NCM,
			QualityFlags:          nonNilStrings(item.QualityFlags),
			Active:                item.Active,
			SellableStock: catalogQuantityFactResponse{
				Quantity: item.SellableStock.Quantity,
				Quality:  nonNilStrings(item.SellableStock.Quality),
			},
			CurrentPrice: catalogMoneyFactResponse{
				Amount:   item.CurrentPrice.Amount,
				Currency: item.CurrentPrice.Currency,
				Quality:  nonNilStrings(item.CurrentPrice.Quality),
			},
			Cost: catalogMoneyFactResponse{
				Amount:   item.Cost.Amount,
				Currency: item.Cost.Currency,
				Quality:  nonNilStrings(item.Cost.Quality),
			},
		})
	}
	response := catalogPageEnvelope{Items: items, PageSize: len(items), AsOf: page.AsOf.UTC()}
	if !suppressCursor && page.NextCursor != nil {
		if encoded, err := page.NextCursor.Encode(); err == nil {
			response.NextCursor = &encoded
		}
	}
	return response
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func writeCatalogPageError(w http.ResponseWriter, err error) {
	if queryErr, ok := err.(*catalogPageQueryError); ok {
		// Every code that knows its accepted domain says it. Gating this on one
		// hardcoded code made allowed_range dead data on the two other sites
		// that fill it (invalid_erp_source, invalid_ids): built, then dropped.
		if queryErr.allowedRange != "" {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": queryErr.code, "allowed_range": queryErr.allowedRange})
			return
		}
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": queryErr.code})
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		httpx.WriteJSON(w, http.StatusGatewayTimeout, map[string]string{"error": "deadline_exceeded"})
		return
	}
	if internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable) {
		httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "source_unavailable"})
		return
	}
	// Anything else is a fault on our side, and it says so. The branch this
	// replaced answered 503 source_unavailable to EVERY unclassified error —
	// including a mis-wired reader and an unresolved assortment policy, both of
	// which are bugs in this process. An operator reading source_unavailable goes
	// and checks Oracle, which is working fine. A read error that genuinely means
	// "the source is down" carries ReadErrorSourceUnavailable and is matched
	// above; it does not need to be recognised a second time by sniffing the
	// error's text, which is how an error merely CONTAINING that word — the bare
	// errors.New the nil-reader guard used to raise — was getting classified as
	// an outage.
	slog.Error("catalog.page", "result", "500", "error", err.Error())
	httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
}

// requirePageReader answers the one question every paged catalog handler has to
// ask first, in one place instead of four. A nil PageReader is composition
// declining to wire a catalog source, which IS an unavailable source — unlike
// the errors below it, which are faults in a source that is wired.
func (h Handler) requirePageReader(w http.ResponseWriter) bool {
	if h.PageReader != nil {
		return true
	}
	httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "source_unavailable"})
	return false
}

func (h Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	productID := r.PathValue("id")
	id, parseErr := strconv.Atoi(productID)
	if parseErr != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_identity", "productId must be a positive integer")
		return
	}
	product, err := h.Service.GetProduct(r.Context(), domain.InternalProductID(id))
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("catalog.product", "action", "get", "result", "404", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
			writeError(w, http.StatusNotFound, "CATALOG_PRODUCT_NOT_FOUND", "product not found")
			return
		}
		slog.Error("catalog.product", "action", "get", "result", "500", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
		writeCatalogReadError(w, err)
		return
	}
	slog.Info("catalog.product", "action", "get", "result", "200", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, product)
}

func writeCatalogReadError(w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "source_unavailable") {
		writeError(w, http.StatusServiceUnavailable, "source_unavailable", "catalog source unavailable")
		return
	}
	writeError(w, http.StatusInternalServerError, "CATALOG_INTERNAL_ERROR", "internal error")
}

func (h Handler) handleTaxonomy(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		slog.Error("catalog.taxonomy", "action", "list", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeError(w, http.StatusMethodNotAllowed, "CATALOG_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	nodes, err := h.CompatibilityService.ListTaxonomyNodes(r.Context())
	if err != nil {
		slog.Error("catalog.taxonomy", "action", "list", "result", "500", "duration_ms", time.Since(start).Milliseconds())
		writeError(w, http.StatusInternalServerError, "CATALOG_INTERNAL_ERROR", "internal error")
		return
	}
	slog.Info("catalog.taxonomy", "action", "list", "result", "200", "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": nodes})
}

func (h Handler) handleGetEnrichment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	productID := r.PathValue("id")
	enrichment, err := h.CompatibilityService.GetEnrichment(r.Context(), productID)
	if err != nil {
		slog.Error("catalog.enrichment", "action", "get", "result", "500", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
		writeError(w, http.StatusInternalServerError, "CATALOG_INTERNAL_ERROR", "internal error")
		return
	}
	slog.Info("catalog.enrichment", "action", "get", "result", "200", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, enrichment)
}

func (h Handler) handleUpsertEnrichment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	productID := r.PathValue("id")
	var req struct {
		HeightCM             *float64 `json:"height_cm"`
		WidthCM              *float64 `json:"width_cm"`
		LengthCM             *float64 `json:"length_cm"`
		WeightG              *float64 `json:"weight_g"`
		SuggestedPriceAmount *float64 `json:"suggested_price_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("catalog.enrichment", "action", "upsert", "result", "400", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
		writeError(w, http.StatusBadRequest, "CATALOG_ENRICHMENT_INVALID", "malformed request body")
		return
	}
	enrichment := domain.ProductEnrichment{
		ProductID:            productID,
		HeightCM:             req.HeightCM,
		WidthCM:              req.WidthCM,
		LengthCM:             req.LengthCM,
		WeightG:              req.WeightG,
		SuggestedPriceAmount: req.SuggestedPriceAmount,
	}
	if err := h.CompatibilityService.UpsertEnrichment(r.Context(), enrichment); err != nil {
		slog.Error("catalog.enrichment", "action", "upsert", "result", "500", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
		writeError(w, http.StatusInternalServerError, "CATALOG_INTERNAL_ERROR", "internal error")
		return
	}
	slog.Info("catalog.enrichment", "action", "upsert", "result", "200", "product_id", productID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, enrichment)
}
