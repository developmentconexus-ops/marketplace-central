package application

import (
	"context"
	"testing"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type fakeProviderOperationInstallations struct {
	inst  domain.Installation
	found bool
}

func (f fakeProviderOperationInstallations) Get(context.Context, string) (domain.Installation, bool, error) {
	return f.inst, f.found, nil
}

type fakeProviderOperationAccountProber struct {
	snapshot connectorsdomain.AccountSnapshot
}

func (f fakeProviderOperationAccountProber) ProbeAccount(context.Context, connectorsdomain.ProviderAccountRef) (connectorsdomain.AccountSnapshot, error) {
	return f.snapshot, nil
}

type fakeProviderOperationFeeQuoteReader struct {
	snapshot connectorsdomain.FeeQuoteSnapshot
}

func (f fakeProviderOperationFeeQuoteReader) ReadFeeQuote(context.Context, connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	return f.snapshot, nil
}

type providerOperationRunStore struct {
	runs []domain.OperationRun
}

func (s *providerOperationRunStore) SaveOperationRun(_ context.Context, run domain.OperationRun) error {
	s.runs = append(s.runs, run)
	return nil
}

func (s *providerOperationRunStore) ListByInstallation(_ context.Context, _ string) ([]domain.OperationRun, error) {
	return append([]domain.OperationRun(nil), s.runs...), nil
}

func TestProviderOperationServiceRejectsUnavailableCapability(t *testing.T) {
	t.Parallel()

	service := ProviderOperationService{
		tenantID: "tenant_default",
		installations: fakeProviderOperationInstallations{
			inst: domain.Installation{
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				Status:            domain.InstallationStatusConnected,
				HealthStatus:      domain.HealthStatusHealthy,
				ExternalAccountID: "691607102",
				RuntimeCapabilities: []domain.RuntimeCapability{{
					Code:       domain.RuntimeCapabilityStockWrite,
					State:      domain.RuntimeCapabilityStateUnavailable,
					Executable: false,
				}},
			},
			found: true,
		},
		capabilities: connectorsapp.NewMarketplaceCapabilityService(nil),
		operations:   NewOperationService(&providerOperationRunStore{}, "tenant_default"),
		now:          func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	}

	_, err := service.ProbeAccount(context.Background(), "inst-1")
	if err == nil || err.Error() != "INTEGRATIONS_CAPABILITY_UNAVAILABLE" {
		t.Fatalf("err = %v", err)
	}
}

func TestProviderOperationServiceProbesAccount(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	prober := fakeProviderOperationAccountProber{snapshot: connectorsdomain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   "691607102",
		ProviderAccountName: "METALNOBREACABAMENTOS",
		SiteID:              "MLB",
		Status:              "active",
		FetchedAt:           now,
	}}
	service := newProviderOperationServiceForTest(now, store, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode:  "mercado_livre",
		AccountProbes: prober,
	}}))

	result, err := service.ProbeAccount(context.Background(), "inst-1")
	if err != nil {
		t.Fatalf("ProbeAccount() error = %v", err)
	}
	if result.ProviderAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("account name = %q", result.ProviderAccountName)
	}
	if len(store.runs) != 1 || store.runs[0].OperationType != providerOperationTypeAccountProbe || store.runs[0].Status != domain.OperationRunStatusSucceeded {
		t.Fatalf("runs = %#v", store.runs)
	}
}

func TestProviderOperationServiceReadsFeeQuote(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	reader := fakeProviderOperationFeeQuoteReader{snapshot: connectorsdomain.FeeQuoteSnapshot{
		ProviderCode:  "mercado_livre",
		SiteID:        "MLB",
		CategoryID:    "MLB123",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
		FetchedAt:     now,
	}}
	service := newProviderOperationServiceForTest(now, store, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode: "mercado_livre",
		FeeQuotes:    reader,
	}}))

	result, err := service.ReadFeeQuote(context.Background(), "inst-1", connectorsdomain.FeeQuoteInput{
		SiteID:        "MLB",
		CategoryID:    "MLB123",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
	})
	if err != nil {
		t.Fatalf("ReadFeeQuote() error = %v", err)
	}
	if result.ListingTypeID != "gold_special" {
		t.Fatalf("listing type = %q", result.ListingTypeID)
	}
	if len(store.runs) != 1 || store.runs[0].OperationType != providerOperationTypeFeeQuoteRead {
		t.Fatalf("runs = %#v", store.runs)
	}
}

func newProviderOperationServiceForTest(now time.Time, store *providerOperationRunStore, capabilities *connectorsapp.MarketplaceCapabilityService) *ProviderOperationService {
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
					{Code: domain.RuntimeCapabilityAccountProbe, State: domain.RuntimeCapabilityStateAvailable, Executable: true},
					{Code: domain.RuntimeCapabilityFeeQuoteRead, State: domain.RuntimeCapabilityStateAvailable, Executable: true},
				},
			},
			found: true,
		},
		Capabilities: capabilities,
		Operations:   NewOperationService(store, "tenant_default"),
		Now:          func() time.Time { return now },
	})
}
