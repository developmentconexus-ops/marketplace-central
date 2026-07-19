package composition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	marketdomain "marketplace-central/apps/server_core/internal/modules/market/domain"
	marketports "marketplace-central/apps/server_core/internal/modules/market/ports"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

// market_adapters.go wires the M-02 price-intel module's outbound ports to
// concrete composition-owned adapters. These adapters live here (not inside
// modules/market/adapters) because modules/connectors already imports
// modules/market, so a market-side adapter over connectors would be an
// import cycle. modules/market itself stays untouched.

// --- ports.CostReader -------------------------------------------------------

// marketCostReaderAdapter backs ports.CostReader with the same IC-02
// internal_read.GetCostAsOf pattern used by profitability's FactReader
// (apps/server_core/internal/modules/profitability/adapters/internalread/fact_reader.go).
// When no internal_read source is configured (source != "oracle"/"xlsx", or
// oracle connection failed) it honestly returns a nil Money rather than a
// fabricated cost.
type marketCostReaderAdapter struct {
	service   internalreadapp.Service
	available bool
}

var _ marketports.CostReader = (*marketCostReaderAdapter)(nil)

func newMarketCostReaderAdapter(service internalreadapp.Service, available bool) *marketCostReaderAdapter {
	return &marketCostReaderAdapter{service: service, available: available}
}

func (r *marketCostReaderAdapter) GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (*marketdomain.Money, time.Time, error) {
	if !r.available {
		return nil, effectiveAt, nil
	}
	cost, err := r.service.GetCostAsOf(ctx, internalreadports.CostAsOfInput{
		ProductID: productID,
		Policy: internalreaddomain.CostAsOfPolicy{
			CompanyID:   1,
			EffectiveAt: effectiveAt,
			Basis:       internalreaddomain.CostBasisCUSSEMICM,
		},
		Freshness: internalreaddomain.FreshnessPolicy{},
	})
	if err != nil {
		if internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable) {
			return nil, effectiveAt, nil
		}
		return nil, effectiveAt, err
	}
	if cost.Amount == nil {
		return nil, cost.EffectiveAt, nil
	}
	money := &marketdomain.Money{
		Amount:   strconv.FormatFloat(*cost.Amount, 'f', -1, 64),
		Currency: "BRL",
	}
	return money, cost.EffectiveAt, nil
}

// --- ports.ProductIdentityReader --------------------------------------------

// productLinkFanOutLimit bounds the per-installation ListProductLinks scan
// used by LinkedListings below. LinkedListings has no by-productID index at
// any layer (product_links is installation-scoped only), so this adapter
// fans out over every installation for the tenant and filters client-side.
// This is a known scaling caveat: fine for the small number of installations
// per tenant expected today, but not a substitute for a future by-productID
// index if installation counts grow.
const productLinkFanOutLimit = 1000

// marketProductIdentityReaderAdapter backs ports.ProductIdentityReader with
// ERP product identity (internal_read.FindProductsForLinking) and resolved
// product_links rows across all of the tenant's installations.
type marketProductIdentityReaderAdapter struct {
	service       internalreadapp.Service
	available     bool
	installations *integrationsapp.InstallationService
	links         *productlinkspostgres.LinkCandidateRepository
}

var _ marketports.ProductIdentityReader = (*marketProductIdentityReaderAdapter)(nil)

func newMarketProductIdentityReaderAdapter(
	service internalreadapp.Service,
	available bool,
	installations *integrationsapp.InstallationService,
	links *productlinkspostgres.LinkCandidateRepository,
) *marketProductIdentityReaderAdapter {
	return &marketProductIdentityReaderAdapter{
		service:       service,
		available:     available,
		installations: installations,
		links:         links,
	}
}

func (r *marketProductIdentityReaderAdapter) GetLocalIdentity(ctx context.Context, codprod string) (marketdomain.LocalProductIdentity, error) {
	productID, err := strconv.Atoi(strings.TrimSpace(codprod))
	if err != nil {
		return marketdomain.LocalProductIdentity{}, marketports.ErrProductNotFound
	}
	if !r.available {
		return marketdomain.LocalProductIdentity{}, marketports.ErrProductNotFound
	}
	candidates, err := r.service.FindProductsForLinking(ctx, internalreadports.FindProductsInput{ProductID: &productID})
	if err != nil {
		// A codprod absent from the ERP source surfaces as the reader's own
		// not-found error; translate it to the market port sentinel so the
		// collection handler answers 404 instead of a raw 500 (FINDING-M02-COLLECT-4XX, D-86).
		var notFound *erpinternalread.ERPProductNotFoundError
		if errors.As(err, &notFound) {
			return marketdomain.LocalProductIdentity{}, marketports.ErrProductNotFound
		}
		return marketdomain.LocalProductIdentity{}, err
	}
	if len(candidates) == 0 {
		return marketdomain.LocalProductIdentity{}, marketports.ErrProductNotFound
	}
	candidate := candidates[0]
	return marketdomain.LocalProductIdentity{
		Codprod: codprod,
		Descr:   candidate.Name,
		EAN:     candidate.EAN,
		RefForn: candidate.ReferenceCode,
		Marca:   candidate.BrandName,
	}, nil
}

func (r *marketProductIdentityReaderAdapter) LinkedListings(ctx context.Context, codprod string) ([]string, error) {
	productID, err := strconv.Atoi(strings.TrimSpace(codprod))
	if err != nil {
		return nil, nil
	}
	installations, err := r.installations.List(ctx)
	if err != nil {
		return nil, err
	}
	listingIDs := make([]string, 0)
	for _, inst := range installations {
		links, err := r.links.ListProductLinks(ctx, inst.InstallationID, productLinkFanOutLimit)
		if err != nil {
			return nil, err
		}
		for _, link := range links {
			if link.State != productlinksdomain.ProductLinkStateResolved {
				continue
			}
			if link.InternalProductID == nil || *link.InternalProductID != productID {
				continue
			}
			listingID := listingsdomain.ListingID{
				InstallationID:    link.InstallationID,
				ProviderListingID: link.ProviderItemID,
				VariationID:       link.ProviderVariationID,
			}
			listingIDs = append(listingIDs, listingID.String())
		}
	}
	return listingIDs, nil
}

// --- ports.PriceIntelCollector -----------------------------------------------

// marketPriceIntelCollectorAdapter backs ports.PriceIntelCollector with the
// mercado_livre connectors capability adapter.
//
// AMBIGUITY (flagged per task instructions, not silently guessed):
// connectorsapp.MarketplaceCapabilityService (the module's public capability
// facade) exposes only AccountProbes/Listings/FeeQuotes/StockReads/
// StockWrites/PriceWrites/ListingWrites/Orders — it has no catalog-read or
// market-read accessor. The concrete *mercadolivreconnector.CapabilityAdapter
// directly implements ports.MarketReader/ports.CatalogReader/
// ports.CatalogOffersReader (see var _ assertions in capability_adapter.go),
// so this adapter is wired against that concrete adapter instead, which is
// already constructed in root.go before the market wiring point. Editing the
// connectors module to add a facade accessor is out of scope for this task.
//
// AMBIGUITY 2: SearchCatalogByEAN/GetCatalogProduct/ListCatalogOffers carry
// no installation/account context in ports.PriceIntelCollector's signature
// (only an ean or catalogProductID string), yet every connectors capability
// call requires a full ProviderAccountRef. There is no existing "the
// tenant's marketplace installation" concept outside of a specific listing.
// This adapter resolves account context for these three methods by picking
// the tenant's first Connected mercado_livre installation; if none is
// Connected (or none exists), it returns a typed 503 ProviderStatusError
// rather than guessing at credentials.
type marketPriceIntelCollectorAdapter struct {
	capabilities  *mercadolivreconnector.CapabilityAdapter
	installations *integrationsapp.InstallationService
	tenantID      string
	now           func() time.Time
	log           *slog.Logger
	warned        sync.Map // operation|status -> struct{}, warn-once dedup
}

var _ marketports.PriceIntelCollector = (*marketPriceIntelCollectorAdapter)(nil)

const mercadoLivreProviderCode = "mercado_livre"

func newMarketPriceIntelCollectorAdapter(
	capabilities *mercadolivreconnector.CapabilityAdapter,
	installations *integrationsapp.InstallationService,
	tenantID string,
) *marketPriceIntelCollectorAdapter {
	return &marketPriceIntelCollectorAdapter{
		capabilities:  capabilities,
		installations: installations,
		tenantID:      tenantID,
		now:           time.Now,
	}
}

// accountRefForTenant resolves the account ref for the tenant's mercado_livre
// installation, used by the three catalog methods that carry no per-call
// installation context (see AMBIGUITY 2 above).
func (r *marketPriceIntelCollectorAdapter) accountRefForTenant(ctx context.Context) (connectorsdomain.ProviderAccountRef, error) {
	installations, err := r.installations.List(ctx)
	if err != nil {
		return connectorsdomain.ProviderAccountRef{}, err
	}
	for _, inst := range installations {
		if inst.ProviderCode != mercadoLivreProviderCode {
			continue
		}
		if inst.Status == integrationsdomain.InstallationStatusConnected {
			return accountRefForInstallation(r.tenantID, inst), nil
		}
	}
	// No Connected mercado_livre installation: same honest 503 as zero
	// installations — never fall back to a non-Connected account's credentials.
	return connectorsdomain.ProviderAccountRef{}, &marketports.ProviderStatusError{StatusCode: http.StatusServiceUnavailable}
}

// accountRefForListing resolves account context from the listing ID itself
// (installationID~providerListingID~variationID), which is the honest source
// of truth for GetOwnItemPricing/GetPriceToWin: no cross-installation
// heuristic needed.
func (r *marketPriceIntelCollectorAdapter) accountRefForListing(ctx context.Context, listingID string) (connectorsdomain.ProviderAccountRef, string, error) {
	parsed, err := listingsdomain.ParseListingID(listingID)
	if err != nil {
		return connectorsdomain.ProviderAccountRef{}, "", &marketports.ProviderStatusError{StatusCode: http.StatusBadRequest}
	}
	inst, found, err := r.installations.Get(ctx, parsed.InstallationID)
	if err != nil {
		return connectorsdomain.ProviderAccountRef{}, "", err
	}
	if !found {
		return connectorsdomain.ProviderAccountRef{}, "", &marketports.ProviderStatusError{StatusCode: http.StatusNotFound}
	}
	return accountRefForInstallation(r.tenantID, inst), parsed.ProviderListingID, nil
}

func accountRefForInstallation(tenantID string, inst integrationsdomain.Installation) connectorsdomain.ProviderAccountRef {
	return connectorsdomain.ProviderAccountRef{
		TenantID:          tenantID,
		InstallationID:    inst.InstallationID,
		ProviderCode:      inst.ProviderCode,
		ProviderAccountID: inst.ExternalAccountID,
	}
}

func (r *marketPriceIntelCollectorAdapter) SearchCatalogByEAN(ctx context.Context, ean string) (marketports.CatalogIdentityResult, error) {
	ref, err := r.accountRefForTenant(ctx)
	if err != nil {
		return marketports.CatalogIdentityResult{}, err
	}
	result, err := r.capabilities.SearchCatalogByEAN(ctx, ref, ean)
	if err != nil {
		return marketports.CatalogIdentityResult{}, mapMarketProviderError(ctx, err)
	}
	candidates := make([]marketdomain.IdentityCandidate, 0, len(result.Products))
	for _, p := range result.Products {
		candidates = append(candidates, marketdomain.IdentityCandidate{
			CatalogProductID: p.CatalogProductID,
			Title:            p.Title,
			Attrs:            convertCatalogAttrs(p.Attrs),
		})
	}
	return marketports.CatalogIdentityResult{Candidates: candidates, FetchedAt: result.FetchedAt}, nil
}

func (r *marketPriceIntelCollectorAdapter) GetCatalogProduct(ctx context.Context, catalogProductID string) (marketports.CatalogProductDetail, error) {
	ref, err := r.accountRefForTenant(ctx)
	if err != nil {
		return marketports.CatalogProductDetail{}, err
	}
	product, err := r.capabilities.GetCatalogProduct(ctx, ref, catalogProductID)
	if err != nil {
		return marketports.CatalogProductDetail{}, r.mapAndObserve(ctx, "GetCatalogProduct", err)
	}
	detail := marketports.CatalogProductDetail{
		CatalogProductID: product.ID,
		FetchedAt:        product.FetchedAt,
	}
	if product.BuyBoxWinner != nil {
		detail.BuyBoxWinnerPrice = convertMoney(product.BuyBoxWinner.Price)
		detail.BuyBoxSellerID = product.BuyBoxWinner.SellerID
	}
	return detail, nil
}

func (r *marketPriceIntelCollectorAdapter) ListCatalogOffers(ctx context.Context, catalogProductID string) ([]marketdomain.ValidatedOffer, error) {
	ref, err := r.accountRefForTenant(ctx)
	if err != nil {
		return nil, err
	}
	offers, err := r.capabilities.ListCatalogOffers(ctx, ref, catalogProductID)
	if err != nil {
		if errors.Is(err, connectorsdomain.ErrCatalogOffersUnavailable) {
			return nil, marketports.ErrCatalogOffersDisabled
		}
		return nil, r.mapAndObserve(ctx, "ListCatalogOffers", err)
	}
	fetchedAt := r.now()
	validated := make([]marketdomain.ValidatedOffer, 0, len(offers))
	for _, offer := range offers {
		vo, err := marketdomain.NewValidatedOffer(marketdomain.ValidatedOfferInput{
			CatalogProductID: catalogProductID,
			SellerID:         offer.SellerID,
			Price:            convertMoney(offer.Price),
			Condition:        offer.Condition,
			MatchStatus:      marketdomain.MatchStatusAccept,
			FetchedAt:        fetchedAt,
		})
		if err != nil {
			return nil, err
		}
		validated = append(validated, vo)
	}
	return validated, nil
}

func (r *marketPriceIntelCollectorAdapter) GetOwnItemPricing(ctx context.Context, listingID string) (marketports.OwnItemPricing, error) {
	ref, providerListingID, err := r.accountRefForListing(ctx, listingID)
	if err != nil {
		return marketports.OwnItemPricing{}, err
	}
	pricing, err := r.capabilities.GetOwnItemPricing(ctx, ref, providerListingID)
	if err != nil {
		return marketports.OwnItemPricing{}, r.mapAndObserve(ctx, "GetOwnItemPricing", err)
	}
	var ourPrice *marketdomain.Money
	if pricing.Amount != nil {
		ourPrice = &marketdomain.Money{Amount: *pricing.Amount, Currency: pricing.Currency}
	}
	// WinnerPrice: connectors' OwnItemPricing (market_read.go) carries only
	// our own Amount/RegularAmount, no competitor/buy-box price. There is no
	// connectors read that supplies a winner price alongside own pricing, so
	// this is left honestly nil rather than guessed from PriceToWin (a
	// separate call with different semantics/timing). Flagged per task
	// instructions.
	return marketports.OwnItemPricing{
		ListingID:   listingID,
		OurPrice:    ourPrice,
		WinnerPrice: nil,
		FetchedAt:   pricing.FetchedAt,
	}, nil
}

func (r *marketPriceIntelCollectorAdapter) GetPriceToWin(ctx context.Context, listingID string) (marketports.PriceToWinSignal, error) {
	ref, providerListingID, err := r.accountRefForListing(ctx, listingID)
	if err != nil {
		return marketports.PriceToWinSignal{}, err
	}
	signal, err := r.capabilities.GetPriceToWin(ctx, ref, providerListingID)
	if err != nil {
		return marketports.PriceToWinSignal{}, r.mapAndObserve(ctx, "GetPriceToWin", err)
	}
	var position *marketdomain.MarketPosition
	if signal.Position != nil {
		position = &marketdomain.MarketPosition{Rank: signal.Position.Rank, Total: signal.Position.Total}
	}
	return marketports.PriceToWinSignal{
		ListingID:   listingID,
		TargetPrice: convertMoney(signal.TargetPrice),
		Position:    position,
		FetchedAt:   signal.FetchedAt,
	}, nil
}

func convertMoney(m *connectorsdomain.Money) *marketdomain.Money {
	if m == nil {
		return nil
	}
	return &marketdomain.Money{Amount: m.Amount, Currency: m.Currency}
}

func convertCatalogAttrs(attrs []connectorsdomain.CatalogAttribute) []marketdomain.IdentityCandidateAttribute {
	out := make([]marketdomain.IdentityCandidateAttribute, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, marketdomain.IdentityCandidateAttribute{ID: a.ID, Name: a.Name, ValueName: a.ValueName})
	}
	return out
}

// mapMarketProviderError is the per-call error translator used by all five
// PriceIntelCollector methods. It first checks the caller's context for a
// deadline: if the collection pipeline's ctx expired, the failure is a
// TIMEOUT regardless of how the connector dressed the transport error, so an
// error wrapping context.DeadlineExceeded is returned for
// classifyProviderFailure to recognize (IC-03 Amendment A1).
func mapMarketProviderError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if ctxErr := ctx.Err(); errors.Is(ctxErr, context.DeadlineExceeded) {
		return fmt.Errorf("provider call timed out: %w", context.DeadlineExceeded)
	}
	return mapMarketConnectorError(err)
}

// mapMarketConnectorError translates connectors' error vocabulary
// (sentinel-wrapped by mapPricingReaderError, or a raw *CapabilityError) into
// the market module's connector-agnostic ports vocabulary, per the CONTRACT
// mapping: provider 4xx -> ProviderStatusError, 5xx -> ProviderStatusError,
// rate-limit -> ErrRateLimited, timeout -> passthrough wrapping
// context.DeadlineExceeded (the connectors CapabilityAdapter now preserves
// the transport cause on its error chain).
func mapMarketConnectorError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("provider call timed out: %w", context.DeadlineExceeded)
	}
	switch {
	case errors.Is(err, connectorsdomain.ErrCatalogOffersUnavailable):
		return marketports.ErrCatalogOffersDisabled
	case errors.Is(err, connectorsdomain.ErrRateLimited):
		return marketports.ErrRateLimited
	case errors.Is(err, connectorsdomain.ErrUnauthorized):
		return &marketports.ProviderStatusError{StatusCode: http.StatusUnauthorized}
	case errors.Is(err, connectorsdomain.ErrNotFound):
		return &marketports.ProviderStatusError{StatusCode: http.StatusNotFound}
	case errors.Is(err, connectorsdomain.ErrProviderUnavailable):
		return &marketports.ProviderStatusError{StatusCode: http.StatusBadGateway}
	}
	var capErr *connectorsdomain.CapabilityError
	if errors.As(err, &capErr) && capErr != nil {
		switch capErr.Code {
		case connectorsdomain.ErrCodeProviderRateLimited:
			return marketports.ErrRateLimited
		case connectorsdomain.ErrCodeProviderAuth:
			return &marketports.ProviderStatusError{StatusCode: http.StatusUnauthorized}
		case connectorsdomain.ErrCodeProviderInvalidReference:
			return &marketports.ProviderStatusError{StatusCode: http.StatusNotFound}
		case connectorsdomain.ErrCodeProviderTransient:
			return &marketports.ProviderStatusError{StatusCode: http.StatusBadGateway}
		case connectorsdomain.ErrCodeProviderValidation, connectorsdomain.ErrCodeProviderPayloadInvalid, connectorsdomain.ErrCodeProviderUnsupportedShape:
			return &marketports.ProviderStatusError{StatusCode: http.StatusBadRequest}
		}
	}
	return &marketports.ProviderStatusError{StatusCode: http.StatusBadGateway}
}

// mapAndObserve maps a connector error to the market ports vocabulary and, at
// this market-collection boundary, emits a warn-once observability signal for
// provider-body faults (e.g. the ML 400 that otherwise dies silently). Dedup is
// per operation+status so a looping pipeline logs the fault once, not per call
// (FINDING-M02-LIVE-2 observability, D-86).
func (r *marketPriceIntelCollectorAdapter) mapAndObserve(ctx context.Context, operation string, err error) error {
	mapped := mapMarketProviderError(ctx, err)
	r.observeProviderFailure(operation, err, mapped)
	return mapped
}

func (r *marketPriceIntelCollectorAdapter) observeProviderFailure(operation string, err, mapped error) {
	// Only provider-body faults carry a *CapabilityError with a payload worth
	// surfacing. Timeouts and bare sentinels have no body — skip them.
	var capErr *connectorsdomain.CapabilityError
	if !errors.As(err, &capErr) || capErr == nil {
		return
	}
	status := 0
	var statusErr *marketports.ProviderStatusError
	if errors.As(mapped, &statusErr) {
		status = statusErr.StatusCode
	}
	key := operation + "|" + strconv.Itoa(status)
	if _, loaded := r.warned.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	logger := r.log
	if logger == nil {
		logger = slog.Default()
	}
	logger.Warn("market collection provider failure",
		"operation", operation,
		"provider_code", mercadoLivreProviderCode,
		"status", status,
		"error_code", string(capErr.Code),
		"provider_body", sanitizeProviderBody(capErr.Message),
	)
}

const maxProviderBodyLog = 256

// The provider body is attacker/provider-controlled (echoed straight from the
// HTTP response), so redaction must cover more than the literal "Bearer X"
// header shape: bare ML credential tokens and token-bearing JSON/query fields
// can appear anywhere in a 4xx/401 body.
var (
	// Authorization header echoed back + bare ML credential token shapes.
	credentialTokenPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/-]+=*`),
		regexp.MustCompile(`(?i)\b(?:APP_USR|APP_ID|TG)-[A-Za-z0-9._~+/-]+`),
	}
	// token-bearing JSON fields or query params: keep the field name, drop value.
	credentialFieldPattern = regexp.MustCompile(`(?i)("?(?:access_token|refresh_token|client_secret|token|secret|authorization)"?\s*[:=]\s*"?)[A-Za-z0-9._~+/-]+`)
)

// sanitizeProviderBody makes a provider payload safe to log: whitespace is
// collapsed, credential tokens are redacted BEFORE truncation (so a token can
// never survive by straddling the length bound), and the result is rune-bounded
// so an oversized body cannot flood the log or split a multibyte rune.
func sanitizeProviderBody(body string) string {
	body = strings.Join(strings.Fields(body), " ")
	body = credentialFieldPattern.ReplaceAllString(body, "$1[REDACTED]")
	for _, pattern := range credentialTokenPatterns {
		body = pattern.ReplaceAllString(body, "[REDACTED]")
	}
	if runes := []rune(body); len(runes) > maxProviderBodyLog {
		body = string(runes[:maxProviderBodyLog]) + "…"
	}
	return body
}

// --- listingsports.EvidenceReader (F-01 M-05) -------------------------------

// listingsEvidenceAdapter backs listings' consumer-owned
// listingsports.EvidenceReader (internal/modules/listings/ports/evidence.go)
// over the market module's public marketports.EvidenceReader
// (internal/modules/market/ports/evidence_reader.go), per the D-21
// composition-adapter pattern above: this is the ONE place allowed to import
// both modules' public ports, and it does the anti-corruption mapping from
// market's domain types (marketdomain.CompetitiveSignal/MarketAggregate/
// Verdict) to listings' own mirror types (listingsports.EvidenceSignal/
// EvidenceAggregate/EvidenceVerdict) — listings never imports modules/market.
type listingsEvidenceAdapter struct {
	reader marketports.EvidenceReader
}

var _ listingsports.EvidenceReader = (*listingsEvidenceAdapter)(nil)

func newListingsEvidenceAdapter(reader marketports.EvidenceReader) *listingsEvidenceAdapter {
	return &listingsEvidenceAdapter{reader: reader}
}

func (a *listingsEvidenceAdapter) Signals(ctx context.Context, listingIDs []string) ([]listingsports.EvidenceSignal, error) {
	signals, err := a.reader.Signals(ctx, listingIDs)
	if err != nil {
		return nil, err
	}
	out := make([]listingsports.EvidenceSignal, 0, len(signals))
	for _, sig := range signals {
		out = append(out, listingsports.EvidenceSignal{
			ListingID:   sig.ListingID,
			OurPrice:    convertListingsMoney(sig.OurPrice),
			TargetPrice: convertListingsMoney(sig.TargetPrice),
			Position:    convertListingsPosition(sig.Position),
			FetchedAt:   sig.FetchedAt,
		})
	}
	return out, nil
}

func (a *listingsEvidenceAdapter) Aggregates(ctx context.Context, codprods []string) ([]listingsports.EvidenceAggregate, error) {
	aggregates, err := a.reader.Aggregates(ctx, codprods)
	if err != nil {
		return nil, err
	}
	out := make([]listingsports.EvidenceAggregate, 0, len(aggregates))
	for _, agg := range aggregates {
		out = append(out, listingsports.EvidenceAggregate{
			ProductID: agg.ProductID,
			NOffers:   agg.NOffers,
			NSellers:  agg.NSellers,
		})
	}
	return out, nil
}

func (a *listingsEvidenceAdapter) Verdicts(ctx context.Context, codprods []string) ([]listingsports.EvidenceVerdict, error) {
	verdicts, err := a.reader.Verdicts(ctx, codprods)
	if err != nil {
		return nil, err
	}
	out := make([]listingsports.EvidenceVerdict, 0, len(verdicts))
	for _, v := range verdicts {
		out = append(out, listingsports.EvidenceVerdict{
			ProductID:   v.ProductID,
			MatchStatus: string(v.MatchStatus),
		})
	}
	return out, nil
}

// convertListingsMoney maps market's *marketdomain.Money to listings' own
// *listingsdomain.Money mirror (same Amount/Currency shape, listings'
// Currency is its own PriceCurrency type per read_model.go:85).
func convertListingsMoney(m *marketdomain.Money) *listingsdomain.Money {
	if m == nil {
		return nil
	}
	return &listingsdomain.Money{Amount: m.Amount, Currency: listingsdomain.PriceCurrency(m.Currency)}
}

// convertListingsPosition maps market's *marketdomain.MarketPosition to
// listings' own *listingsdomain.SignalPosition mirror (IC-03 position
// {rank, total}).
func convertListingsPosition(p *marketdomain.MarketPosition) *listingsdomain.SignalPosition {
	if p == nil {
		return nil
	}
	return &listingsdomain.SignalPosition{Rank: p.Rank, Total: p.Total}
}
