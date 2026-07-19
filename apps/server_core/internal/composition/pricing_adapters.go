package composition

import (
	"context"
	"strconv"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// pricing_adapters.go wires degrau 3 (live COTAÇÃO) of the pricing tariff
// ladder to the connectors module's read-only mercado_livre capabilities. Like
// market_adapters.go, these adapters live in composition (not modules/pricing)
// because they bridge pricing's ports onto connectors' readers, and only
// composition is allowed to import both. They reuse the raw connectors readers
// directly (never the integrations catalog-match probe wrapper, which records a
// provider_operations audit row on every call — a decompose must not spam the
// audit log).

// firstConnectedMLAccountRef resolves the tenant's first Connected
// mercado_livre installation's account ref. found=false (not an error) when
// none is Connected: degrau 3 is simply unavailable and the composite falls to
// degrau 4 (unlike the market module's 503, degrau 3 degrades silently).
func firstConnectedMLAccountRef(ctx context.Context, installations *integrationsapp.InstallationService, tenantID string) (connectorsdomain.ProviderAccountRef, bool, error) {
	insts, err := installations.List(ctx)
	if err != nil {
		return connectorsdomain.ProviderAccountRef{}, false, err
	}
	for _, inst := range insts {
		if inst.ProviderCode != mercadoLivreProviderCode {
			continue
		}
		if inst.Status == integrationsdomain.InstallationStatusConnected {
			return accountRefForInstallation(tenantID, inst), true, nil
		}
	}
	return connectorsdomain.ProviderAccountRef{}, false, nil
}

// --- pricingports.CategoryResolver ------------------------------------------

// pricingCategoryResolverAdapter backs pricing's CategoryResolver with the
// mercado_livre catalog-match reader (source 1: EAN+title). A future catmap
// source lands behind the same port.
type pricingCategoryResolverAdapter struct {
	catalogMatch  connectorsports.CatalogMatchReader
	installations *integrationsapp.InstallationService
	tenantID      string
}

var _ pricingports.CategoryResolver = (*pricingCategoryResolverAdapter)(nil)

func newPricingCategoryResolverAdapter(catalogMatch connectorsports.CatalogMatchReader, installations *integrationsapp.InstallationService, tenantID string) *pricingCategoryResolverAdapter {
	return &pricingCategoryResolverAdapter{catalogMatch: catalogMatch, installations: installations, tenantID: tenantID}
}

func (r *pricingCategoryResolverAdapter) ResolveCategory(ctx context.Context, q pricingports.CategoryQuery) (pricingports.CategoryResolution, bool, error) {
	ref, found, err := firstConnectedMLAccountRef(ctx, r.installations, r.tenantID)
	if err != nil {
		return pricingports.CategoryResolution{}, false, err
	}
	if !found {
		return pricingports.CategoryResolution{}, false, nil
	}
	// D-83: EAN and title are always sent together.
	snapshot, err := r.catalogMatch.ReadCatalogMatch(ctx, connectorsdomain.CatalogMatchInput{
		AccountRef: ref,
		EAN:        q.EAN,
		Query:      q.Titulo,
	})
	if err != nil {
		return pricingports.CategoryResolution{}, false, err
	}
	categoryID, fonte, ok := selectCatalogMatchCategory(snapshot)
	if !ok {
		return pricingports.CategoryResolution{}, false, nil
	}
	return pricingports.CategoryResolution{
		CategoriaID: categoryID,
		Fonte:       fonte,
		// The provider exposes no confidence score on catalog match (ADR-17:
		// left nil, never a fabricated 0).
		Confianca: nil,
		Data:      snapshot.FetchedAt.UTC().Format(time.RFC3339),
	}, true, nil
}

// selectCatalogMatchCategory picks the fee-quote category, mirroring the
// integrations probe: the buy-box category when present (a resolved catalog
// product, EAN-anchored => EAN-CATALOGO), else the top-ranked domain-discovery
// prediction (title-driven => TITULO). Pure.
func selectCatalogMatchCategory(snapshot connectorsdomain.CatalogMatchSnapshot) (string, pricingports.FonteCategoria, bool) {
	if snapshot.BuyBox != nil {
		if id := strings.TrimSpace(snapshot.BuyBox.CategoryID); id != "" {
			return id, pricingports.FonteCategoriaEANCatalogo, true
		}
	}
	for _, prediction := range snapshot.DomainDiscovery {
		if id := strings.TrimSpace(prediction.CategoryID); id != "" {
			return id, pricingports.FonteCategoriaTitulo, true
		}
	}
	return "", "", false
}

// --- pricingports.CommissionQuoter ------------------------------------------

// pricingCommissionQuoterAdapter backs pricing's CommissionQuoter with the
// mercado_livre fee-quote (listing_prices) reader.
type pricingCommissionQuoterAdapter struct {
	fees          connectorsports.FeeQuoteReader
	installations *integrationsapp.InstallationService
	tenantID      string
}

var _ pricingports.CommissionQuoter = (*pricingCommissionQuoterAdapter)(nil)

func newPricingCommissionQuoterAdapter(fees connectorsports.FeeQuoteReader, installations *integrationsapp.InstallationService, tenantID string) *pricingCommissionQuoterAdapter {
	return &pricingCommissionQuoterAdapter{fees: fees, installations: installations, tenantID: tenantID}
}

func (r *pricingCommissionQuoterAdapter) QuoteCommission(ctx context.Context, q pricingports.CommissionQuery) (pricingports.CommissionQuote, bool, error) {
	ref, found, err := firstConnectedMLAccountRef(ctx, r.installations, r.tenantID)
	if err != nil {
		return pricingports.CommissionQuote{}, false, err
	}
	if !found {
		return pricingports.CommissionQuote{}, false, nil
	}
	price, err := strconv.ParseFloat(strings.TrimSpace(q.Preco.Amount), 64)
	if err != nil || price <= 0 {
		// An unparseable/non-positive basis cannot anchor a real quote (ADR-17).
		return pricingports.CommissionQuote{}, false, nil
	}
	snapshot, err := r.fees.ReadFeeQuote(ctx, connectorsdomain.FeeQuoteInput{
		AccountRef:    ref,
		CategoryID:    q.CategoriaID,
		ListingTypeID: q.ListingTypeID,
		PriceAmount:   price,
		CurrencyID:    q.Preco.Currency,
	})
	if err != nil {
		return pricingports.CommissionQuote{}, false, err
	}
	quote, ok := commissionQuoteFrom(snapshot)
	return quote, ok, nil
}

// commissionQuoteFrom maps a fee snapshot to a degrau-3 commission quote. A nil
// commission percent is an honest miss (ADR-17: never a fabricated 0). The ML
// percentage_fee is already a percentage (e.g. 11 => 11%), matching the
// engine's ComissaoPct unit. Data prefers the provider's source-updated time,
// falling back to the fetch time. Pure.
func commissionQuoteFrom(snapshot connectorsdomain.FeeQuoteSnapshot) (pricingports.CommissionQuote, bool) {
	if snapshot.CommissionPercent == nil {
		return pricingports.CommissionQuote{}, false
	}
	data := snapshot.FetchedAt
	if snapshot.SourceUpdatedAt != nil {
		data = *snapshot.SourceUpdatedAt
	}
	return pricingports.CommissionQuote{
		ComissaoPct: strconv.FormatFloat(*snapshot.CommissionPercent, 'f', -1, 64),
		Data:        data.UTC().Format(time.RFC3339),
	}, true
}
