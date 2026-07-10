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

type fakeProviderOperationStockReader struct {
	snapshot connectorsdomain.StockSnapshot
}

func (f fakeProviderOperationStockReader) ReadStock(context.Context, connectorsdomain.ProviderListingRef) (connectorsdomain.StockSnapshot, error) {
	return f.snapshot, nil
}

type providerOperationRunStore struct {
	runs []domain.OperationRun
}

type providerOperationCapabilityStore struct {
	states []domain.CapabilityState
}

func (s *providerOperationRunStore) SaveOperationRun(_ context.Context, run domain.OperationRun) error {
	s.runs = append(s.runs, run)
	return nil
}

func (s *providerOperationRunStore) ListByInstallation(_ context.Context, _ string) ([]domain.OperationRun, error) {
	return append([]domain.OperationRun(nil), s.runs...), nil
}

func (s *providerOperationCapabilityStore) Upsert(_ context.Context, states []domain.CapabilityState) error {
	s.states = append(s.states, states...)
	return nil
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
	capabilityStore := &providerOperationCapabilityStore{}
	prober := fakeProviderOperationAccountProber{snapshot: connectorsdomain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   "691607102",
		ProviderAccountName: "METALNOBREACABAMENTOS",
		SiteID:              "MLB",
		Status:              "active",
		FetchedAt:           now,
	}}
	service := newProviderOperationServiceForTest(now, store, capabilityStore, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
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
	if len(capabilityStore.states) != 1 || capabilityStore.states[0].Status != domain.CapabilityStatusEnabled {
		t.Fatalf("capability states = %#v", capabilityStore.states)
	}
	if capabilityStore.states[0].LastEvaluatedAt == nil || !capabilityStore.states[0].LastEvaluatedAt.Equal(now) {
		t.Fatalf("last_evaluated_at = %v, want %v", capabilityStore.states[0].LastEvaluatedAt, now)
	}
}

func TestProviderOperationServiceReadsFeeQuote(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	capabilityStore := &providerOperationCapabilityStore{}
	reader := fakeProviderOperationFeeQuoteReader{snapshot: connectorsdomain.FeeQuoteSnapshot{
		ProviderCode:  "mercado_livre",
		SiteID:        "MLB",
		CategoryID:    "MLB123",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
		FetchedAt:     now,
	}}
	service := newProviderOperationServiceForTest(now, store, capabilityStore, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
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
	if len(capabilityStore.states) != 1 || capabilityStore.states[0].CapabilityCode != "fee_quote_read" {
		t.Fatalf("capability states = %#v", capabilityStore.states)
	}
}

func TestProviderOperationServiceReadsStock(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	store := &providerOperationRunStore{}
	capabilityStore := &providerOperationCapabilityStore{}
	reader := fakeProviderOperationStockReader{snapshot: connectorsdomain.StockSnapshot{
		ProviderCode:      "mercado_livre",
		ProviderItemID:    "MLB123",
		AvailableQuantity: 11,
		ProviderStatus:    "active",
		FetchedAt:         now,
		Scope:             connectorsdomain.StockScopeItem,
	}}
	service := newProviderOperationServiceForTest(now, store, capabilityStore, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode: "mercado_livre",
		StockReads:   reader,
	}}))

	result, err := service.ReadStock(context.Background(), "inst-1", connectorsdomain.ProviderListingRef{
		ProviderItemID: "MLB123",
	})
	if err != nil {
		t.Fatalf("ReadStock() error = %v", err)
	}
	if result.AvailableQuantity != 11 {
		t.Fatalf("available quantity = %d", result.AvailableQuantity)
	}
	if len(store.runs) != 1 || store.runs[0].OperationType != providerOperationTypeStockRead {
		t.Fatalf("runs = %#v", store.runs)
	}
	if len(capabilityStore.states) != 1 || capabilityStore.states[0].CapabilityCode != "stock_read" {
		t.Fatalf("capability states = %#v", capabilityStore.states)
	}
}

func TestProviderOperationServiceMarksCapabilityDegradedOnProviderFailure(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	runStore := &providerOperationRunStore{}
	capabilityStore := &providerOperationCapabilityStore{}
	failingProber := fakeProviderOperationAccountProberErr{
		err: connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, "provider unavailable"),
	}
	service := newProviderOperationServiceForTest(now, runStore, capabilityStore, connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode:  "mercado_livre",
		AccountProbes: failingProber,
	}}))

	_, err := service.ProbeAccount(context.Background(), "inst-1")
	if err == nil {
		t.Fatal("ProbeAccount() error = nil, want provider failure")
	}
	if len(capabilityStore.states) != 1 || capabilityStore.states[0].Status != domain.CapabilityStatusDegraded {
		t.Fatalf("capability states = %#v", capabilityStore.states)
	}
}

func TestProviderOperationServiceDefaultClockAdvancesOperationIDs(t *testing.T) {
	t.Parallel()

	store := &providerOperationRunStore{}
	capabilityStore := &providerOperationCapabilityStore{}
	prober := fakeProviderOperationAccountProber{snapshot: connectorsdomain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   "691607102",
		ProviderAccountName: "METALNOBREACABAMENTOS",
		SiteID:              "MLB",
		Status:              "active",
	}}
	service := NewProviderOperationService(ProviderOperationServiceConfig{
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
				},
			},
			found: true,
		},
		Capabilities: connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
			ProviderCode:  "mercado_livre",
			AccountProbes: prober,
		}}),
		CapabilityStates: capabilityStore,
		Operations:       NewOperationService(store, "tenant_default"),
	})

	if _, err := service.ProbeAccount(context.Background(), "inst-1"); err != nil {
		t.Fatalf("first ProbeAccount() error = %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	if _, err := service.ProbeAccount(context.Background(), "inst-1"); err != nil {
		t.Fatalf("second ProbeAccount() error = %v", err)
	}

	if len(store.runs) != 2 {
		t.Fatalf("runs = %d, want 2", len(store.runs))
	}
	if store.runs[0].OperationRunID == store.runs[1].OperationRunID {
		t.Fatalf("operation ids did not advance: %q == %q", store.runs[0].OperationRunID, store.runs[1].OperationRunID)
	}
	if store.runs[0].StartedAt == nil || store.runs[1].StartedAt == nil || !store.runs[1].StartedAt.After(*store.runs[0].StartedAt) {
		t.Fatalf("started_at did not advance: first=%v second=%v", store.runs[0].StartedAt, store.runs[1].StartedAt)
	}
}

type fakeProviderOperationAccountProberErr struct {
	err error
}

func (f fakeProviderOperationAccountProberErr) ProbeAccount(context.Context, connectorsdomain.ProviderAccountRef) (connectorsdomain.AccountSnapshot, error) {
	return connectorsdomain.AccountSnapshot{}, f.err
}

func newProviderOperationServiceForTest(now time.Time, store *providerOperationRunStore, capabilityStore *providerOperationCapabilityStore, capabilities *connectorsapp.MarketplaceCapabilityService) *ProviderOperationService {
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
					{Code: domain.RuntimeCapabilityStockRead, State: domain.RuntimeCapabilityStateAvailable, Executable: true},
				},
			},
			found: true,
		},
		Capabilities:     capabilities,
		CapabilityStates: capabilityStore,
		Operations:       NewOperationService(store, "tenant_default"),
		Now:              func() time.Time { return now },
	})
}
