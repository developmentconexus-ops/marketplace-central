package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

type CreateInstallationInput struct {
	InstallationID string
	ProviderCode   string
	DisplayName    string
	Family         string
}

type InstallationService struct {
	repo     ports.InstallationRepository
	tenantID string
}

func NewInstallationService(repo ports.InstallationRepository, tenantID string) *InstallationService {
	return &InstallationService{repo: repo, tenantID: tenantID}
}

func (s *InstallationService) CreateDraft(ctx context.Context, input CreateInstallationInput) (domain.Installation, error) {
	installationID := strings.TrimSpace(input.InstallationID)
	providerCode := strings.TrimSpace(input.ProviderCode)
	displayName := strings.TrimSpace(input.DisplayName)
	family := strings.TrimSpace(input.Family)
	if installationID == "" || providerCode == "" || displayName == "" || family == "" {
		return domain.Installation{}, errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}

	now := time.Now().UTC()
	inst := domain.Installation{
		InstallationID:      installationID,
		TenantID:            s.tenantID,
		ProviderCode:        providerCode,
		Family:              domain.IntegrationFamily(family),
		DisplayName:         displayName,
		Status:              domain.InstallationStatusDraft,
		HealthStatus:        domain.HealthStatusHealthy,
		CreatedAt:           now,
		UpdatedAt:           now,
		RuntimeCapabilities: []domain.RuntimeCapability{},
	}
	inst.ConnectionSnapshot = domain.ProjectConnectionSnapshot(inst, domain.AuthStrategyUnknown, nil, "")

	return inst, s.repo.CreateInstallation(ctx, inst)
}

func (s *InstallationService) Get(ctx context.Context, installationID string) (domain.Installation, bool, error) {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return domain.Installation{}, false, errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}

	return s.repo.GetInstallation(ctx, installationID)
}

func (s *InstallationService) List(ctx context.Context) ([]domain.Installation, error) {
	return s.repo.ListInstallations(ctx)
}

func (s *InstallationService) UpdateStatus(ctx context.Context, installationID string, status domain.InstallationStatus, health domain.HealthStatus) error {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" || !isValidInstallationStatus(status) || !isValidHealthStatus(health) {
		return errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}

	return s.repo.UpdateInstallationStatus(ctx, installationID, status, health)
}

func (s *InstallationService) ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" || !isValidConnectionState(snapshot.State) || !isValidHealthStatus(snapshot.Health) {
		return errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}

	return s.repo.ApplyConnectionSnapshot(ctx, installationID, snapshot, strings.TrimSpace(activeCredentialID))
}

func isValidInstallationStatus(status domain.InstallationStatus) bool {
	switch status {
	case domain.InstallationStatusDraft,
		domain.InstallationStatusPendingConnection,
		domain.InstallationStatusConnected,
		domain.InstallationStatusDegraded,
		domain.InstallationStatusRequiresReauth,
		domain.InstallationStatusDisconnected,
		domain.InstallationStatusSuspended,
		domain.InstallationStatusFailed:
		return true
	default:
		return false
	}
}

func isValidHealthStatus(status domain.HealthStatus) bool {
	switch status {
	case domain.HealthStatusHealthy,
		domain.HealthStatusWarning,
		domain.HealthStatusCritical:
		return true
	default:
		return false
	}
}

func isValidConnectionState(state domain.ConnectionState) bool {
	switch state {
	case domain.ConnectionStateDraft,
		domain.ConnectionStatePending,
		domain.ConnectionStateConnected,
		domain.ConnectionStateDegraded,
		domain.ConnectionStateNeedsReauth,
		domain.ConnectionStateDisconnected:
		return true
	default:
		return false
	}
}
