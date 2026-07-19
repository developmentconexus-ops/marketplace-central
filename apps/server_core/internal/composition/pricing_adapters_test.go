package composition

import (
	"context"
	"errors"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// --- fakes -----------------------------------------------------------------

type fakeCatalogMatchReader struct {
	snap connectorsdomain.CatalogMatchSnapshot
	err  error
	got  *connectorsdomain.CatalogMatchInput
}

func (f *fakeCatalogMatchReader) ReadCatalogMatch(_ context.Context, in connectorsdomain.CatalogMatchInput) (connectorsdomain.CatalogMatchSnapshot, error) {
	f.got = &in
	return f.snap, f.err
}

type fakeFeeQuoteReader struct {
	snap connectorsdomain.FeeQuoteSnapshot
	err  error
	got  *connectorsdomain.FeeQuoteInput
}

func (f *fakeFeeQuoteReader) ReadFeeQuote(_ context.Context, in connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	f.got = &in
	return f.snap, f.err
}

func connectedMLSvc() *integrationsapp.InstallationService {
	return integrationsapp.NewInstallationService(fakeInstallationRepo{items: []integrationsdomain.Installation{{
		InstallationID:    "inst-1",
		TenantID:          "tenant_default",
		ProviderCode:      mercadoLivreProviderCode,
		Status:            integrationsdomain.InstallationStatusConnected,
		ExternalAccountID: "acct-1",
	}}}, "tenant_default")
}

func noConnectedMLSvc() *integrationsapp.InstallationService {
	return integrationsapp.NewInstallationService(fakeInstallationRepo{items: []integrationsdomain.Installation{{
		InstallationID: "inst-1",
		ProviderCode:   mercadoLivreProviderCode,
		Status:         integrationsdomain.InstallationStatusRequiresReauth,
	}}}, "tenant_default")
}

// --- CategoryResolver -------------------------------------------------------

func TestCategoryResolver_BuyBoxIsEANCatalogo(t *testing.T) {
	reader := &fakeCatalogMatchReader{snap: connectorsdomain.CatalogMatchSnapshot{
		BuyBox:    &connectorsdomain.BuyBoxSnapshot{CategoryID: "MLB1055"},
		FetchedAt: time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC),
	}}
	r := newPricingCategoryResolverAdapter(reader, connectedMLSvc(), "tenant_default")

	res, ok, err := r.ResolveCategory(context.Background(), pricingports.CategoryQuery{EAN: "789", Titulo: "Tenis X"})
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if res.CategoriaID != "MLB1055" || res.Fonte != pricingports.FonteCategoriaEANCatalogo {
		t.Fatalf("res = %+v, want MLB1055/EAN-CATALOGO", res)
	}
	if res.Data != "2026-07-19T12:00:00Z" {
		t.Fatalf("Data = %q", res.Data)
	}
	if res.Confianca != nil {
		t.Fatalf("Confianca = %v, want nil", res.Confianca)
	}
	// D-83: EAN and title sent together.
	if reader.got == nil || reader.got.EAN != "789" || reader.got.Query != "Tenis X" {
		t.Fatalf("catalog-match input = %+v, want EAN+title together (D-83)", reader.got)
	}
	if reader.got.AccountRef.InstallationID != "inst-1" {
		t.Fatalf("account ref not resolved: %+v", reader.got.AccountRef)
	}
}

func TestCategoryResolver_PredictionIsTitulo(t *testing.T) {
	reader := &fakeCatalogMatchReader{snap: connectorsdomain.CatalogMatchSnapshot{
		BuyBox:          &connectorsdomain.BuyBoxSnapshot{CategoryID: "  "},
		DomainDiscovery: []connectorsdomain.DomainPrediction{{CategoryID: "MLB2000"}, {CategoryID: "MLB9999"}},
		FetchedAt:       time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC),
	}}
	r := newPricingCategoryResolverAdapter(reader, connectedMLSvc(), "tenant_default")

	res, ok, err := r.ResolveCategory(context.Background(), pricingports.CategoryQuery{EAN: "789", Titulo: "Tenis X"})
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if res.CategoriaID != "MLB2000" || res.Fonte != pricingports.FonteCategoriaTitulo {
		t.Fatalf("res = %+v, want top prediction MLB2000/TITULO", res)
	}
}

func TestCategoryResolver_NoCategoryIsMiss(t *testing.T) {
	reader := &fakeCatalogMatchReader{snap: connectorsdomain.CatalogMatchSnapshot{}}
	r := newPricingCategoryResolverAdapter(reader, connectedMLSvc(), "tenant_default")

	_, ok, err := r.ResolveCategory(context.Background(), pricingports.CategoryQuery{EAN: "789", Titulo: "X"})
	if ok || err != nil {
		t.Fatalf("ok=%v err=%v, want miss", ok, err)
	}
}

func TestCategoryResolver_NoConnectedInstallationIsMiss(t *testing.T) {
	reader := &fakeCatalogMatchReader{snap: connectorsdomain.CatalogMatchSnapshot{BuyBox: &connectorsdomain.BuyBoxSnapshot{CategoryID: "MLB1"}}}
	r := newPricingCategoryResolverAdapter(reader, noConnectedMLSvc(), "tenant_default")

	_, ok, err := r.ResolveCategory(context.Background(), pricingports.CategoryQuery{EAN: "789", Titulo: "X"})
	if ok || err != nil {
		t.Fatalf("ok=%v err=%v, want silent miss", ok, err)
	}
	if reader.got != nil {
		t.Fatalf("must not call catalog match with no connected installation")
	}
}

func TestCategoryResolver_ReaderErrorPropagates(t *testing.T) {
	sentinel := errors.New("boom")
	reader := &fakeCatalogMatchReader{err: sentinel}
	r := newPricingCategoryResolverAdapter(reader, connectedMLSvc(), "tenant_default")

	_, ok, err := r.ResolveCategory(context.Background(), pricingports.CategoryQuery{EAN: "789", Titulo: "X"})
	if ok || !errors.Is(err, sentinel) {
		t.Fatalf("ok=%v err=%v, want sentinel", ok, err)
	}
}

// --- CommissionQuoter -------------------------------------------------------

func TestCommissionQuoter_QuotesPercentPrefersSourceTime(t *testing.T) {
	pct := 11.0
	src := time.Date(2026, 7, 18, 9, 0, 0, 0, time.UTC)
	reader := &fakeFeeQuoteReader{snap: connectorsdomain.FeeQuoteSnapshot{
		CommissionPercent: &pct,
		SourceUpdatedAt:   &src,
		FetchedAt:         time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC),
	}}
	r := newPricingCommissionQuoterAdapter(reader, connectedMLSvc(), "tenant_default")

	q, ok, err := r.QuoteCommission(context.Background(), pricingports.CommissionQuery{
		CategoriaID: "MLB1055", Preco: pricingdomain.Money{Amount: "199.90", Currency: "BRL"}, ListingTypeID: "gold_special",
	})
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if q.ComissaoPct != "11" {
		t.Fatalf("ComissaoPct = %q, want 11", q.ComissaoPct)
	}
	if q.Data != "2026-07-18T09:00:00Z" {
		t.Fatalf("Data = %q, want source-updated time", q.Data)
	}
	// fee-quote input mapping.
	if reader.got == nil || reader.got.CategoryID != "MLB1055" || reader.got.ListingTypeID != "gold_special" || reader.got.PriceAmount != 199.90 || reader.got.CurrencyID != "BRL" {
		t.Fatalf("fee-quote input = %+v", reader.got)
	}
}

func TestCommissionQuoter_FallsBackToFetchedAt(t *testing.T) {
	pct := 13.5
	reader := &fakeFeeQuoteReader{snap: connectorsdomain.FeeQuoteSnapshot{
		CommissionPercent: &pct,
		FetchedAt:         time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC),
	}}
	r := newPricingCommissionQuoterAdapter(reader, connectedMLSvc(), "tenant_default")

	q, ok, err := r.QuoteCommission(context.Background(), pricingports.CommissionQuery{
		CategoriaID: "MLB1", Preco: pricingdomain.Money{Amount: "10.00", Currency: "BRL"}, ListingTypeID: "gold_pro",
	})
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if q.ComissaoPct != "13.5" || q.Data != "2026-07-19T12:00:00Z" {
		t.Fatalf("quote = %+v", q)
	}
}

func TestCommissionQuoter_NilPercentIsMiss(t *testing.T) {
	reader := &fakeFeeQuoteReader{snap: connectorsdomain.FeeQuoteSnapshot{CommissionPercent: nil, FetchedAt: time.Now().UTC()}}
	r := newPricingCommissionQuoterAdapter(reader, connectedMLSvc(), "tenant_default")

	_, ok, err := r.QuoteCommission(context.Background(), pricingports.CommissionQuery{
		CategoriaID: "MLB1", Preco: pricingdomain.Money{Amount: "10.00", Currency: "BRL"}, ListingTypeID: "gold_pro",
	})
	if ok || err != nil {
		t.Fatalf("ok=%v err=%v, want miss (no fabricated commission)", ok, err)
	}
}

func TestCommissionQuoter_NoConnectedInstallationIsMiss(t *testing.T) {
	reader := &fakeFeeQuoteReader{}
	r := newPricingCommissionQuoterAdapter(reader, noConnectedMLSvc(), "tenant_default")

	_, ok, err := r.QuoteCommission(context.Background(), pricingports.CommissionQuery{
		CategoriaID: "MLB1", Preco: pricingdomain.Money{Amount: "10.00", Currency: "BRL"}, ListingTypeID: "gold_pro",
	})
	if ok || err != nil {
		t.Fatalf("ok=%v err=%v, want silent miss", ok, err)
	}
	if reader.got != nil {
		t.Fatalf("must not quote with no connected installation")
	}
}

func TestCommissionQuoter_ReaderErrorPropagates(t *testing.T) {
	sentinel := errors.New("boom")
	reader := &fakeFeeQuoteReader{err: sentinel}
	r := newPricingCommissionQuoterAdapter(reader, connectedMLSvc(), "tenant_default")

	_, ok, err := r.QuoteCommission(context.Background(), pricingports.CommissionQuery{
		CategoriaID: "MLB1", Preco: pricingdomain.Money{Amount: "10.00", Currency: "BRL"}, ListingTypeID: "gold_pro",
	})
	if ok || !errors.Is(err, sentinel) {
		t.Fatalf("ok=%v err=%v, want sentinel", ok, err)
	}
}
