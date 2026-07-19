package application

import (
	"context"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

// Listing types quoted by the catalog-match probe. classico = gold_special
// (the standard paid tier), premium = gold_pro (docs/marketplaces/mercado-livre.md).
const (
	catalogMatchListingTypeClassico = "gold_special"
	catalogMatchListingTypePremium  = "gold_pro"
	catalogMatchDefaultCurrencyID   = "BRL"

	providerOperationTypeCatalogMatchProbe = "catalog_match_probe"
)

// CatalogMatchProbeInput drives the read-only tier-3 catalog-match probe.
type CatalogMatchProbeInput struct {
	EAN         string
	Query       string
	PriceAmount float64
}

// CatalogMatchProbeResult is the normalized probe surface: catalog hits, buy-box,
// category predictions, a two-tier fee quote against the resolved/predicted
// category, and honest flags. Every unknown is null + flagged, never zero (ADR-17).
type CatalogMatchProbeResult struct {
	CatalogHits     []connectorsdomain.CatalogHit       `json:"catalog_hits"`
	BuyBox          *connectorsdomain.BuyBoxSnapshot    `json:"buy_box"`
	DomainDiscovery []connectorsdomain.DomainPrediction `json:"domain_discovery"`
	FeeQuote        *CatalogMatchFeeQuote               `json:"fee_quote"`
	Flags           CatalogMatchFlags                   `json:"flags"`
	FetchedAt       time.Time                           `json:"fetched_at"`
}

// CatalogMatchFeeQuote holds the resolved category used for the quote plus the
// per-listing-type fee tiers. Nil (JSON null) when no category could be resolved
// or no price was available to quote against.
type CatalogMatchFeeQuote struct {
	CategoryID  string               `json:"category_id"`
	PriceAmount float64              `json:"price_amount"`
	CurrencyID  string               `json:"currency_id"`
	Classico    *CatalogMatchFeeTier `json:"classico"`
	Premium     *CatalogMatchFeeTier `json:"premium"`
}

// CatalogMatchFeeTier is one listing-type fee tier. Fees are null when the
// provider omits them.
type CatalogMatchFeeTier struct {
	ListingTypeID string   `json:"listing_type_id"`
	PercentageFee *float64 `json:"percentage_fee"`
	FixedFee      *float64 `json:"fixed_fee"`
}

// CatalogMatchFlags carry the honest observation signals the probe exists to surface.
type CatalogMatchFlags struct {
	CategoryPredita bool `json:"category_predita"`
	BuyBoxNull      bool `json:"buy_box_null"`
	NoCatalogHit    bool `json:"no_catalog_hit"`
}

// ProbeCatalogMatch runs the read-only tier-3 COTAÇÃO-MATCH probe end to end:
// EAN → catalog → buy-box/category → domain-discovery prediction → listing_prices
// fee quote. Gated on the fee-quote runtime capability (the probe composes fee
// quotes). No provider writes.
func (s *ProviderOperationService) ProbeCatalogMatch(ctx context.Context, installationID string, input CatalogMatchProbeInput) (CatalogMatchProbeResult, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityFeeQuoteRead)
	if err != nil {
		return CatalogMatchProbeResult{}, err
	}
	if s.catalogMatch == nil {
		return CatalogMatchProbeResult{}, domain.ErrInstallationNotFound
	}

	accountRef := s.accountRef(inst)
	startedAt := s.now()
	snapshot, execErr := s.catalogMatch.ReadCatalogMatch(ctx, connectorsdomain.CatalogMatchInput{
		AccountRef: accountRef,
		EAN:        input.EAN,
		Query:      input.Query,
	})
	if execErr == nil {
		result := buildCatalogMatchResult(snapshot)
		categoryID, predicted := resolveCatalogMatchCategory(result)
		result.Flags.CategoryPredita = predicted
		result.FeeQuote, execErr = s.composeCatalogMatchFeeQuote(ctx, inst.ProviderCode, accountRef, categoryID, catalogMatchQuotePrice(input.PriceAmount, result.BuyBox))
		if execErr == nil {
			if recordErr := s.recordCatalogMatchProbe(ctx, inst.InstallationID, startedAt, nil, result); recordErr != nil {
				return CatalogMatchProbeResult{}, recordErr
			}
			return result, nil
		}
	}

	if recordErr := s.recordCatalogMatchProbe(ctx, inst.InstallationID, startedAt, execErr, CatalogMatchProbeResult{}); recordErr != nil {
		return CatalogMatchProbeResult{}, recordErr
	}
	return CatalogMatchProbeResult{}, execErr
}

func buildCatalogMatchResult(snapshot connectorsdomain.CatalogMatchSnapshot) CatalogMatchProbeResult {
	result := CatalogMatchProbeResult{
		CatalogHits:     snapshot.CatalogHits,
		BuyBox:          snapshot.BuyBox,
		DomainDiscovery: snapshot.DomainDiscovery,
		FetchedAt:       snapshot.FetchedAt,
	}
	result.Flags = CatalogMatchFlags{
		BuyBoxNull:   snapshot.BuyBox == nil || strings.TrimSpace(snapshot.BuyBox.CategoryID) == "",
		NoCatalogHit: len(snapshot.CatalogHits) == 0,
	}
	return result
}

// resolveCatalogMatchCategory picks the fee-quote category: the buy-box category
// when present, else the top-ranked domain-discovery prediction (predicted=true).
func resolveCatalogMatchCategory(result CatalogMatchProbeResult) (categoryID string, predicted bool) {
	if result.BuyBox != nil {
		if categoryID = strings.TrimSpace(result.BuyBox.CategoryID); categoryID != "" {
			return categoryID, false
		}
	}
	for _, prediction := range result.DomainDiscovery {
		if categoryID = strings.TrimSpace(prediction.CategoryID); categoryID != "" {
			return categoryID, true
		}
	}
	return "", false
}

// catalogMatchQuotePrice picks the price to quote against: the caller override
// when positive, else the buy-box price when the provider exposed one.
func catalogMatchQuotePrice(override float64, buyBox *connectorsdomain.BuyBoxSnapshot) float64 {
	if override > 0 {
		return override
	}
	if buyBox != nil && buyBox.Price != nil {
		return *buyBox.Price
	}
	return 0
}

func (s *ProviderOperationService) composeCatalogMatchFeeQuote(ctx context.Context, providerCode string, accountRef connectorsdomain.ProviderAccountRef, categoryID string, price float64) (*CatalogMatchFeeQuote, error) {
	if categoryID == "" || price <= 0 {
		return nil, nil
	}

	reader, err := s.capabilities.FeeQuoteReader(providerCode)
	if err != nil {
		return nil, err
	}

	classico, err := readCatalogMatchFeeTier(ctx, reader, accountRef, categoryID, catalogMatchListingTypeClassico, price)
	if err != nil {
		return nil, err
	}
	premium, err := readCatalogMatchFeeTier(ctx, reader, accountRef, categoryID, catalogMatchListingTypePremium, price)
	if err != nil {
		return nil, err
	}

	return &CatalogMatchFeeQuote{
		CategoryID:  categoryID,
		PriceAmount: price,
		CurrencyID:  catalogMatchDefaultCurrencyID,
		Classico:    classico,
		Premium:     premium,
	}, nil
}

func readCatalogMatchFeeTier(ctx context.Context, reader connectorsports.FeeQuoteReader, accountRef connectorsdomain.ProviderAccountRef, categoryID, listingTypeID string, price float64) (*CatalogMatchFeeTier, error) {
	quote, err := reader.ReadFeeQuote(ctx, connectorsdomain.FeeQuoteInput{
		AccountRef:    accountRef,
		CategoryID:    categoryID,
		ListingTypeID: listingTypeID,
		PriceAmount:   price,
		CurrencyID:    catalogMatchDefaultCurrencyID,
	})
	if err != nil {
		return nil, err
	}
	return &CatalogMatchFeeTier{
		ListingTypeID: quote.ListingTypeID,
		PercentageFee: quote.CommissionPercent,
		FixedFee:      quote.FixedFeeAmount,
	}, nil
}

func (s *ProviderOperationService) recordCatalogMatchProbe(ctx context.Context, installationID string, startedAt time.Time, execErr error, result CatalogMatchProbeResult) error {
	return s.recordOperation(ctx, installationID, providerOperationTypeCatalogMatchProbe, startedAt, execErr, map[string]any{
		"catalog_hit_count": len(result.CatalogHits),
		"category_predita":  result.Flags.CategoryPredita,
		"buy_box_null":      result.Flags.BuyBoxNull,
		"no_catalog_hit":    result.Flags.NoCatalogHit,
	})
}
