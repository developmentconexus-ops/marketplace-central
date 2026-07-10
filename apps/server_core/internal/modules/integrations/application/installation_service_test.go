package application

import (
	"context"
	"reflect"
	"testing"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

type stubInstallationRepo struct {
	saved            []domain.Installation
	list             []domain.Installation
	updated          []installationStatusUpdate
	appliedSnapshots []appliedConnectionSnapshot
	getByID          map[string]domain.Installation
}

type installationStatusUpdate struct {
	installationID string
	status         domain.InstallationStatus
	health         domain.HealthStatus
}

type appliedConnectionSnapshot struct {
	installationID     string
	snapshot           domain.ConnectionSnapshot
	activeCredentialID string
}

type stubRuntimeCapabilitySource struct {
	items []domain.RuntimeCapability
}

type stubInstallationCapabilityStateResolver struct {
	states []domain.CapabilityState
}

type stubAuthSessionLookup struct {
	session domain.AuthSession
	found   bool
}

func (s stubAuthSessionLookup) GetAuthSession(_ context.Context, installationID string) (domain.AuthSession, bool, error) {
	if !s.found || s.session.InstallationID != installationID {
		return domain.AuthSession{}, false, nil
	}
	return s.session, true, nil
}

type stubCredentialLookup struct {
	credential domain.Credential
	found      bool
}

func (s stubCredentialLookup) GetActiveCredential(_ context.Context, installationID string) (domain.Credential, bool, error) {
	if !s.found || s.credential.InstallationID != installationID {
		return domain.Credential{}, false, nil
	}
	return s.credential, true, nil
}

type stubPayloadDecryptor struct {
	payload map[string]any
}

func (s stubPayloadDecryptor) DecryptJSON(encoded []byte) (map[string]any, string, error) {
	return s.payload, "key-1", nil
}

type stubAccountProber struct {
	snapshot connectorsdomain.AccountSnapshot
}

func (s stubAccountProber) ProbeAccount(_ context.Context, ref connectorsdomain.ProviderAccountRef) (connectorsdomain.AccountSnapshot, error) {
	result := s.snapshot
	if result.ProviderAccountID == "" {
		result.ProviderAccountID = ref.ProviderAccountID
	}
	return result, nil
}

func (s stubRuntimeCapabilitySource) Project(inst domain.Installation) []domain.RuntimeCapability {
	return append([]domain.RuntimeCapability(nil), s.items...)
}

func (s stubInstallationCapabilityStateResolver) Resolve(_ context.Context, _ string, _ ports.MarketplaceCapabilities) ([]domain.CapabilityState, error) {
	return append([]domain.CapabilityState(nil), s.states...), nil
}

func (s *stubInstallationRepo) CreateInstallation(_ context.Context, inst domain.Installation) error {
	s.saved = append(s.saved, inst)
	return nil
}

func (s *stubInstallationRepo) GetInstallation(_ context.Context, installationID string) (domain.Installation, bool, error) {
	if s.getByID != nil {
		inst, ok := s.getByID[installationID]
		if ok {
			return inst, true, nil
		}
	}
	for _, inst := range s.list {
		if inst.InstallationID == installationID {
			return inst, true, nil
		}
	}
	return domain.Installation{}, false, nil
}

func (s *stubInstallationRepo) ListInstallations(_ context.Context) ([]domain.Installation, error) {
	return append([]domain.Installation(nil), s.list...), nil
}

func (s *stubInstallationRepo) UpdateInstallationStatus(_ context.Context, installationID string, status domain.InstallationStatus, health domain.HealthStatus) error {
	s.updated = append(s.updated, installationStatusUpdate{
		installationID: installationID,
		status:         status,
		health:         health,
	})
	return nil
}

func (s *stubInstallationRepo) ApplyConnectionSnapshot(_ context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error {
	s.appliedSnapshots = append(s.appliedSnapshots, appliedConnectionSnapshot{
		installationID:     installationID,
		snapshot:           snapshot,
		activeCredentialID: activeCredentialID,
	})
	return nil
}

func TestCreateDraftInstallation(t *testing.T) {
	t.Parallel()

	repo := &stubInstallationRepo{}
	svc := NewInstallationService(repo, "tenant-default")

	inst, err := svc.CreateDraft(context.Background(), CreateInstallationInput{
		InstallationID: "inst_001",
		ProviderCode:   "mercado_livre",
		DisplayName:    "ML Primary",
		Family:         string(domain.IntegrationFamilyMarketplace),
	})
	if err != nil {
		t.Fatalf("CreateDraft() error = %v", err)
	}

	if inst.Status != domain.InstallationStatusDraft {
		t.Fatalf("installation status = %q, want %q", inst.Status, domain.InstallationStatusDraft)
	}
	if inst.HealthStatus != domain.HealthStatusHealthy {
		t.Fatalf("installation health = %q, want %q", inst.HealthStatus, domain.HealthStatusHealthy)
	}
	if inst.TenantID != "tenant-default" {
		t.Fatalf("tenant_id = %q, want %q", inst.TenantID, "tenant-default")
	}
	if len(repo.saved) != 1 {
		t.Fatalf("saved installations = %d, want 1", len(repo.saved))
	}
	if got := repo.saved[0]; got.InstallationID != "inst_001" || got.ProviderCode != "mercado_livre" || got.DisplayName != "ML Primary" || got.Family != domain.IntegrationFamilyMarketplace {
		t.Fatalf("saved installation = %#v, want installation draft fields to match", got)
	}
}

func TestCreateDraftInstallationRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	svc := NewInstallationService(&stubInstallationRepo{}, "tenant-default")

	_, err := svc.CreateDraft(context.Background(), CreateInstallationInput{})
	if err == nil {
		t.Fatal("CreateDraft() error = nil, want invalid input error")
	}
	if got, want := err.Error(), "INTEGRATIONS_INSTALLATION_INVALID"; got != want {
		t.Fatalf("CreateDraft() error = %q, want %q", got, want)
	}
}

func TestInstallationServiceGetListAndUpdateStatus(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0).UTC()
	repo := &stubInstallationRepo{
		list: []domain.Installation{{
			InstallationID: "inst_001",
			TenantID:       "tenant-default",
			ProviderCode:   "mercado_livre",
			Family:         domain.IntegrationFamilyMarketplace,
			DisplayName:    "ML Primary",
			Status:         domain.InstallationStatusConnected,
			HealthStatus:   domain.HealthStatusHealthy,
			CreatedAt:      now,
			UpdatedAt:      now,
		}},
		getByID: map[string]domain.Installation{
			"inst_001": {
				InstallationID: "inst_001",
				TenantID:       "tenant-default",
				ProviderCode:   "mercado_livre",
				Family:         domain.IntegrationFamilyMarketplace,
				DisplayName:    "ML Primary",
				Status:         domain.InstallationStatusConnected,
				HealthStatus:   domain.HealthStatusHealthy,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
	}
	svc := NewInstallationService(repo, "tenant-default")

	inst, ok, err := svc.Get(context.Background(), "inst_001")
	if err != nil || !ok {
		t.Fatalf("Get() = (%#v, %v, %v), want installation, true, nil", inst, ok, err)
	}
	if inst.InstallationID != "inst_001" {
		t.Fatalf("Get() installation_id = %q, want %q", inst.InstallationID, "inst_001")
	}

	items, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if !reflect.DeepEqual(items, repo.list) {
		t.Fatalf("List() = %#v, want %#v", items, repo.list)
	}

	if err := svc.UpdateStatus(context.Background(), "inst_001", domain.InstallationStatusDegraded, domain.HealthStatusWarning); err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if got, want := repo.updated, []installationStatusUpdate{{installationID: "inst_001", status: domain.InstallationStatusDegraded, health: domain.HealthStatusWarning}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("UpdateStatus() calls = %#v, want %#v", got, want)
	}
}

func TestInstallationServiceRejectsInvalidUpdateStatus(t *testing.T) {
	t.Parallel()

	svc := NewInstallationService(&stubInstallationRepo{}, "tenant-default")

	err := svc.UpdateStatus(context.Background(), "", domain.InstallationStatusDraft, domain.HealthStatusHealthy)
	if err == nil {
		t.Fatal("UpdateStatus() error = nil, want invalid input error")
	}
	if got, want := err.Error(), "INTEGRATIONS_INSTALLATION_INVALID"; got != want {
		t.Fatalf("UpdateStatus() error = %q, want %q", err.Error(), "INTEGRATIONS_INSTALLATION_INVALID")
	}
}

func TestInstallationServiceApplyConnectionSnapshotRejectsEmptyInstallation(t *testing.T) {
	t.Parallel()

	svc := NewInstallationService(&stubInstallationRepo{}, "tenant-default")

	err := svc.ApplyConnectionSnapshot(context.Background(), "", domain.ConnectionSnapshot{}, "cred-1")
	if err == nil {
		t.Fatal("ApplyConnectionSnapshot() error = nil, want invalid input error")
	}
	if got, want := err.Error(), "INTEGRATIONS_INSTALLATION_INVALID"; got != want {
		t.Fatalf("ApplyConnectionSnapshot() error = %q, want %q", got, want)
	}
}

func TestInstallationServiceApplyConnectionSnapshotDelegatesToRepository(t *testing.T) {
	t.Parallel()

	repo := &stubInstallationRepo{}
	svc := NewInstallationService(repo, "tenant-default")
	verifiedAt := time.Unix(100, 0).UTC()
	expiresAt := verifiedAt.Add(time.Hour)
	snapshot := domain.ConnectionSnapshot{
		State:               domain.ConnectionStateConnected,
		Health:              domain.HealthStatusHealthy,
		ProviderCode:        "mercado_livre",
		ExternalAccountID:   "seller-1",
		ExternalAccountName: "Seller One",
		AuthStrategy:        domain.AuthStrategyOAuth2,
		LastVerifiedAt:      &verifiedAt,
		ExpiresAt:           &expiresAt,
		NextAction:          domain.ConnectionNextActionNone,
	}

	err := svc.ApplyConnectionSnapshot(context.Background(), "inst-1", snapshot, "cred-1")
	if err != nil {
		t.Fatalf("ApplyConnectionSnapshot() error = %v", err)
	}

	want := []appliedConnectionSnapshot{{
		installationID:     "inst-1",
		snapshot:           snapshot,
		activeCredentialID: "cred-1",
	}}
	if !reflect.DeepEqual(repo.appliedSnapshots, want) {
		t.Fatalf("ApplyConnectionSnapshot() calls = %#v, want %#v", repo.appliedSnapshots, want)
	}
}

func TestInstallationServiceProjectsRuntimeCapabilitiesOnRead(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0).UTC()
	repo := &stubInstallationRepo{
		getByID: map[string]domain.Installation{
			"inst_001": {
				InstallationID: "inst_001",
				TenantID:       "tenant-default",
				ProviderCode:   "mercado_livre",
				Family:         domain.IntegrationFamilyMarketplace,
				DisplayName:    "ML Primary",
				Status:         domain.InstallationStatusConnected,
				HealthStatus:   domain.HealthStatusHealthy,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
	}
	svc := NewInstallationService(repo, "tenant-default")
	svc.SetRuntimeCapabilitySource(stubRuntimeCapabilitySource{
		items: []domain.RuntimeCapability{{
			Code:       domain.RuntimeCapabilityAccountProbe,
			State:      domain.RuntimeCapabilityStateAvailable,
			Executable: true,
		}},
	})

	inst, found, err := svc.Get(context.Background(), "inst_001")
	if err != nil || !found {
		t.Fatalf("Get() = (%#v, %v, %v), want installation, true, nil", inst, found, err)
	}
	if len(inst.RuntimeCapabilities) != 1 || inst.RuntimeCapabilities[0].Code != domain.RuntimeCapabilityAccountProbe {
		t.Fatalf("runtime capabilities = %#v", inst.RuntimeCapabilities)
	}
}

func TestInstallationServiceOverlaysPersistedCapabilityStateOnRead(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0).UTC()
	validatedAt := now.Add(10 * time.Minute)
	repo := &stubInstallationRepo{
		getByID: map[string]domain.Installation{
			"inst_001": {
				InstallationID:    "inst_001",
				TenantID:          "tenant-default",
				ProviderCode:      "mercado_livre",
				Family:            domain.IntegrationFamilyMarketplace,
				DisplayName:       "ML Primary",
				Status:            domain.InstallationStatusConnected,
				HealthStatus:      domain.HealthStatusHealthy,
				ExternalAccountID: "691607102",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	svc := NewInstallationService(repo, "tenant-default")
	svc.SetRuntimeCapabilitySource(stubRuntimeCapabilitySource{
		items: []domain.RuntimeCapability{{
			Code:       domain.RuntimeCapabilityListingRead,
			State:      domain.RuntimeCapabilityStateAvailable,
			Executable: true,
		}},
	})
	svc.SetCapabilityStateResolver(stubInstallationCapabilityStateResolver{
		states: []domain.CapabilityState{{
			CapabilityStateID: "cap_listing_read",
			InstallationID:    "inst_001",
			CapabilityCode:    "listing_read",
			Status:            domain.CapabilityStatusEnabled,
			ReasonCode:        "INTEGRATIONS_OPERATION_SUCCEEDED",
			LastEvaluatedAt:   &validatedAt,
		}},
	})

	inst, found, err := svc.Get(context.Background(), "inst_001")
	if err != nil || !found {
		t.Fatalf("Get() = (%#v, %v, %v), want installation, true, nil", inst, found, err)
	}
	if len(inst.RuntimeCapabilities) != 1 {
		t.Fatalf("runtime capabilities = %#v, want one item", inst.RuntimeCapabilities)
	}
	if !inst.RuntimeCapabilities[0].LiveValidated {
		t.Fatalf("live_validated = false, want true")
	}
	if inst.RuntimeCapabilities[0].LastValidatedAt == nil || !inst.RuntimeCapabilities[0].LastValidatedAt.Equal(validatedAt) {
		t.Fatalf("last_validated_at = %v, want %v", inst.RuntimeCapabilities[0].LastValidatedAt, validatedAt)
	}
}

func TestInstallationServiceKeepsDegradedCapabilityExecutableWhenBasePathExists(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0).UTC()
	validatedAt := now.Add(10 * time.Minute)
	repo := &stubInstallationRepo{
		getByID: map[string]domain.Installation{
			"inst_001": {
				InstallationID:    "inst_001",
				TenantID:          "tenant-default",
				ProviderCode:      "mercado_livre",
				Family:            domain.IntegrationFamilyMarketplace,
				DisplayName:       "ML Primary",
				Status:            domain.InstallationStatusConnected,
				HealthStatus:      domain.HealthStatusHealthy,
				ExternalAccountID: "691607102",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	svc := NewInstallationService(repo, "tenant-default")
	svc.SetRuntimeCapabilitySource(stubRuntimeCapabilitySource{
		items: []domain.RuntimeCapability{{
			Code:       domain.RuntimeCapabilityFeeQuoteRead,
			State:      domain.RuntimeCapabilityStateAvailable,
			Executable: true,
		}},
	})
	svc.SetCapabilityStateResolver(stubInstallationCapabilityStateResolver{
		states: []domain.CapabilityState{{
			CapabilityStateID: "cap_fee_quote_read",
			InstallationID:    "inst_001",
			CapabilityCode:    "fee_quote_read",
			Status:            domain.CapabilityStatusDegraded,
			ReasonCode:        "CONNECTORS_PROVIDER_PAYLOAD_INVALID",
			LastEvaluatedAt:   &validatedAt,
		}},
	})

	inst, found, err := svc.Get(context.Background(), "inst_001")
	if err != nil || !found {
		t.Fatalf("Get() = (%#v, %v, %v), want installation, true, nil", inst, found, err)
	}
	if len(inst.RuntimeCapabilities) != 1 {
		t.Fatalf("runtime capabilities = %#v, want one item", inst.RuntimeCapabilities)
	}
	if inst.RuntimeCapabilities[0].State != domain.RuntimeCapabilityStateDegraded {
		t.Fatalf("state = %q, want degraded", inst.RuntimeCapabilities[0].State)
	}
	if !inst.RuntimeCapabilities[0].Executable {
		t.Fatal("executable = false, want true for degraded revalidation path")
	}
}

func TestInstallationServiceReconcilesConnectedSnapshotFromAuthSession(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0).UTC()
	verifiedAt := now.Add(5 * time.Minute)
	expiresAt := now.Add(time.Hour)
	repo := &stubInstallationRepo{
		getByID: map[string]domain.Installation{
			"inst_001": {
				InstallationID:     "inst_001",
				TenantID:           "tenant-default",
				ProviderCode:       "mercado_livre",
				Family:             domain.IntegrationFamilyMarketplace,
				DisplayName:        "ML Primary",
				Status:             domain.InstallationStatusConnected,
				HealthStatus:       domain.HealthStatusHealthy,
				ActiveCredentialID: "cred-1",
				ConnectionSnapshot: domain.ConnectionSnapshot{
					State:        domain.ConnectionStateConnected,
					Health:       domain.HealthStatusHealthy,
					ProviderCode: "mercado_livre",
					AuthStrategy: domain.AuthStrategyUnknown,
					NextAction:   domain.ConnectionNextActionNone,
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	svc := NewInstallationService(repo, "tenant-default")
	svc.SetStateReconciler(NewInstallationStateReconciler(InstallationStateReconcilerConfig{
		TenantID:      "tenant-default",
		Installations: svc,
		AuthSessions: stubAuthSessionLookup{
			session: domain.AuthSession{
				InstallationID:       "inst_001",
				ProviderAccountID:    "691607102",
				State:                domain.AuthStateValid,
				AccessTokenExpiresAt: &expiresAt,
				LastVerifiedAt:       &verifiedAt,
			},
			found: true,
		},
		Credentials: stubCredentialLookup{
			credential: domain.Credential{
				CredentialID:   "cred-1",
				InstallationID: "inst_001",
			},
			found: true,
		},
		Decryptor: stubPayloadDecryptor{
			payload: map[string]any{
				"access_token":  "token-1",
				"refresh_token": "refresh-1",
			},
		},
		Capabilities: connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
			ProviderCode: "mercado_livre",
			AccountProbes: stubAccountProber{snapshot: connectorsdomain.AccountSnapshot{
				ProviderCode:        "mercado_livre",
				ProviderAccountID:   "691607102",
				ProviderAccountName: "METALNOBREACABAMENTOS",
				SiteID:              "MLB",
				Status:              "active",
				FetchedAt:           verifiedAt,
			}},
		}}),
	}))
	svc.SetRuntimeCapabilitySource(NewRuntimeCapabilityProjector(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode:  "mercado_livre",
		AccountProbes: stubAccountProber{},
	}})))

	inst, found, err := svc.Get(context.Background(), "inst_001")
	if err != nil || !found {
		t.Fatalf("Get() = (%#v, %v, %v), want reconciled installation, true, nil", inst, found, err)
	}
	if inst.ExternalAccountID != "691607102" {
		t.Fatalf("external_account_id = %q, want reconciled provider account id", inst.ExternalAccountID)
	}
	if inst.ExternalAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("external_account_name = %q, want reconciled provider account name", inst.ExternalAccountName)
	}
	if inst.ConnectionSnapshot.AuthStrategy != domain.AuthStrategyOAuth2 {
		t.Fatalf("auth_strategy = %q, want oauth2", inst.ConnectionSnapshot.AuthStrategy)
	}
	if len(inst.RuntimeCapabilities) != 1 || !inst.RuntimeCapabilities[0].Executable {
		t.Fatalf("runtime capabilities = %#v, want executable capability after reconciliation", inst.RuntimeCapabilities)
	}
	if len(repo.appliedSnapshots) != 1 {
		t.Fatalf("applied snapshots = %d, want 1 persisted repair", len(repo.appliedSnapshots))
	}
	if got := repo.appliedSnapshots[0].snapshot; got.ExternalAccountID != "691607102" || got.ExternalAccountName != "METALNOBREACABAMENTOS" || got.AuthStrategy != domain.AuthStrategyOAuth2 {
		t.Fatalf("applied snapshot = %#v", got)
	}
}
