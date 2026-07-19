package application

import (
	"context"
	"testing"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type fakeCatalogMatchReader struct {
	snapshot connectorsdomain.CatalogMatchSnapshot
	input    connectorsdomain.CatalogMatchInput
}

func (f *fakeCatalogMatchReader) ReadCatalogMatch(_ context.Context, input connectorsdomain.CatalogMatchInput) (connectorsdomain.CatalogMatchSnapshot, error) {
	f.input = input
	return f.snapshot, nil
}

// fakeCatalogMatchFeeReader echoes the requested listing type and category so the
// probe's two-tier composition can be asserted per tier.
type fakeCatalogMatchFeeReader struct {
	inputs     []connectorsdomain.FeeQuoteInput
	commission float64
	fixed      float64
}

func (f *fakeCatalogMatchFeeReader) ReadFeeQuote(_ context.Context, in connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	f.inputs = append(f.inputs, in)
	commission := f.commission
	fixed := f.fixed
	return connectorsdomain.FeeQuoteSnapshot{
		CategoryID:        in.CategoryID,
		ListingTypeID:     in.ListingTypeID,
		PriceAmount:       in.PriceAmount,
		CurrencyID:        in.CurrencyID,
		CommissionPercent: &commission,
		FixedFeeAmount:    &fixed,
	}, nil
}

func newCatalogMatchProbeService(now time.Time, store *providerOperationRunStore, catalog *fakeCatalogMatchReader, fees *fakeCatalogMatchFeeReader) *ProviderOperationService {
	return NewProviderOperationService(ProviderOperationServiceConfig{
		TenantID: "tenant_default",
		Installations: fakeProviderOperationInstallations{
			inst: domain.Installation{
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				Status:            domain.InstallationStatusConnected,
				HealthStatus:      domain.HealthStatusHealthy,
				ExternalAccountID: "691607102",
				RuntimeCapabilities: []domain.RuntimeCapability{
					{Code: domain.RuntimeCapabilityFeeQuoteRead, State: domain.RuntimeCapabilityStateAvailable, Executable: true},
				},
			},
			found: true,
		},
		Capabilities: connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
			ProviderCode: "mercado_livre",
			FeeQuotes:    fees,
		}}),
		CapabilityStates: &providerOperationCapabilityStore{},
		Operations:       NewOperationService(store, "tenant_default"),
		CatalogMatch:     catalog,
		Now:              func() time.Time { return now },
	})
}

func floatPtr(v float64) *float64 { return &v }

func TestProbeCatalogMatchQuotesBuyBoxCategory(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	catalog := &fakeCatalogMatchReader{snapshot: connectorsdomain.CatalogMatchSnapshot{
		CatalogHits: []connectorsdomain.CatalogHit{{ProductID: "MLB20270041", Name: "Torneira Gourmet", DomainID: "MLB-BATHROOM_FAUCETS_AND_MIXERS", Status: "active"}},
		BuyBox:      &connectorsdomain.BuyBoxSnapshot{CategoryID: "MLB1276", Price: floatPtr(199.9), ListingType: "gold_special"},
		DomainDiscovery: []connectorsdomain.DomainPrediction{
			{CategoryID: "MLB9999", CategoryName: "Outra", DomainID: "MLB-OTHER"},
		},
		FetchedAt: now,
	}}
	fees := &fakeCatalogMatchFeeReader{commission: 12.5, fixed: 6}
	service := newCatalogMatchProbeService(now, store, catalog, fees)

	result, err := service.ProbeCatalogMatch(context.Background(), "inst-1", CatalogMatchProbeInput{EAN: "7891234567890", Query: "Torneira Cozinha"})
	if err != nil {
		t.Fatalf("ProbeCatalogMatch() error = %v", err)
	}

	if catalog.input.EAN != "7891234567890" || catalog.input.Query != "Torneira Cozinha" || catalog.input.AccountRef.ProviderAccountID != "691607102" {
		t.Fatalf("catalog input = %+v", catalog.input)
	}
	if result.Flags.CategoryPredita || result.Flags.BuyBoxNull || result.Flags.NoCatalogHit {
		t.Fatalf("flags = %+v, want all false", result.Flags)
	}
	if !result.FetchedAt.Equal(now) {
		t.Fatalf("fetched_at = %v, want %v", result.FetchedAt, now)
	}
	if result.FeeQuote == nil {
		t.Fatal("fee quote = nil, want composed")
	}
	// Buy-box category wins over domain-discovery prediction.
	if result.FeeQuote.CategoryID != "MLB1276" || result.FeeQuote.PriceAmount != 199.9 || result.FeeQuote.CurrencyID != "BRL" {
		t.Fatalf("fee quote = %+v", result.FeeQuote)
	}
	if result.FeeQuote.Classico == nil || result.FeeQuote.Classico.ListingTypeID != "gold_special" {
		t.Fatalf("classico = %+v", result.FeeQuote.Classico)
	}
	if result.FeeQuote.Premium == nil || result.FeeQuote.Premium.ListingTypeID != "gold_pro" {
		t.Fatalf("premium = %+v", result.FeeQuote.Premium)
	}
	if result.FeeQuote.Classico.PercentageFee == nil || *result.FeeQuote.Classico.PercentageFee != 12.5 {
		t.Fatalf("classico fee = %+v", result.FeeQuote.Classico)
	}
	if len(fees.inputs) != 2 {
		t.Fatalf("fee reads = %d, want 2", len(fees.inputs))
	}
	for _, in := range fees.inputs {
		if in.CategoryID != "MLB1276" || in.PriceAmount != 199.9 {
			t.Fatalf("fee input = %+v", in)
		}
	}
	if len(store.runs) != 1 || store.runs[0].OperationType != providerOperationTypeCatalogMatchProbe || store.runs[0].Status != domain.OperationRunStatusSucceeded {
		t.Fatalf("runs = %#v", store.runs)
	}
}

func TestProbeCatalogMatchPredictsCategoryWhenBuyBoxNull(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	catalog := &fakeCatalogMatchReader{snapshot: connectorsdomain.CatalogMatchSnapshot{
		CatalogHits: []connectorsdomain.CatalogHit{{ProductID: "MLB22490763", Name: "Ducha Higienica", DomainID: "MLB-SHOWER_HEADS", Status: "active"}},
		BuyBox:      nil,
		DomainDiscovery: []connectorsdomain.DomainPrediction{
			{CategoryID: "MLB5551", CategoryName: "Duchas", DomainID: "MLB-SHOWER_HEADS"},
		},
		FetchedAt: now,
	}}
	fees := &fakeCatalogMatchFeeReader{commission: 14, fixed: 6}
	service := newCatalogMatchProbeService(now, store, catalog, fees)

	// No buy-box price → caller supplies the price to quote against.
	result, err := service.ProbeCatalogMatch(context.Background(), "inst-1", CatalogMatchProbeInput{EAN: "7899999999999", PriceAmount: 150})
	if err != nil {
		t.Fatalf("ProbeCatalogMatch() error = %v", err)
	}

	if !result.Flags.CategoryPredita || !result.Flags.BuyBoxNull {
		t.Fatalf("flags = %+v, want category_predita & buy_box_null", result.Flags)
	}
	if result.Flags.NoCatalogHit {
		t.Fatalf("no_catalog_hit = true, want false")
	}
	if result.FeeQuote == nil || result.FeeQuote.CategoryID != "MLB5551" || result.FeeQuote.PriceAmount != 150 {
		t.Fatalf("fee quote = %+v", result.FeeQuote)
	}
}

func TestProbeCatalogMatchOmitsFeeQuoteWhenCategoryUnresolved(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	catalog := &fakeCatalogMatchReader{snapshot: connectorsdomain.CatalogMatchSnapshot{
		CatalogHits:     nil,
		BuyBox:          nil,
		DomainDiscovery: nil,
		FetchedAt:       now,
	}}
	fees := &fakeCatalogMatchFeeReader{commission: 14, fixed: 6}
	service := newCatalogMatchProbeService(now, store, catalog, fees)

	result, err := service.ProbeCatalogMatch(context.Background(), "inst-1", CatalogMatchProbeInput{Query: "produto sem catalogo", PriceAmount: 150})
	if err != nil {
		t.Fatalf("ProbeCatalogMatch() error = %v", err)
	}

	if result.FeeQuote != nil {
		t.Fatalf("fee quote = %+v, want nil (no category)", result.FeeQuote)
	}
	if len(fees.inputs) != 0 {
		t.Fatalf("fee reads = %d, want 0", len(fees.inputs))
	}
	if !result.Flags.NoCatalogHit || !result.Flags.BuyBoxNull {
		t.Fatalf("flags = %+v, want no_catalog_hit & buy_box_null", result.Flags)
	}
	if result.Flags.CategoryPredita {
		t.Fatalf("category_predita = true, want false")
	}
}
