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
	repo                    ports.InstallationRepository
	tenantID                string
	runtimeCapabilitySource installationRuntimeCapabilitySource
	capabilityStates        installationCapabilityStateResolver
	stateReconciler         installationStateReconciler
}

func NewInstallationService(repo ports.InstallationRepository, tenantID string) *InstallationService {
	return &InstallationService{repo: repo, tenantID: tenantID}
}

type installationRuntimeCapabilitySource interface {
	Project(inst domain.Installation) []domain.RuntimeCapability
}

type installationCapabilityStateResolver interface {
	Resolve(ctx context.Context, installationID string, declared ports.MarketplaceCapabilities) ([]domain.CapabilityState, error)
}

type installationStateReconciler interface {
	Reconcile(ctx context.Context, inst domain.Installation) (domain.Installation, error)
}

func (s *InstallationService) SetRuntimeCapabilitySource(source installationRuntimeCapabilitySource) {
	s.runtimeCapabilitySource = source
}

func (s *InstallationService) SetCapabilityStateResolver(resolver installationCapabilityStateResolver) {
	s.capabilityStates = resolver
}

func (s *InstallationService) SetStateReconciler(reconciler installationStateReconciler) {
	s.stateReconciler = reconciler
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

	inst, found, err := s.repo.GetInstallation(ctx, installationID)
	if err != nil || !found {
		return inst, found, err
	}
	inst, err = s.reconcileState(ctx, inst)
	if err != nil {
		return domain.Installation{}, false, err
	}
	inst, err = s.projectRuntimeCapabilities(ctx, inst)
	if err != nil {
		return domain.Installation{}, false, err
	}
	return inst, true, nil
}

func (s *InstallationService) List(ctx context.Context) ([]domain.Installation, error) {
	items, err := s.repo.ListInstallations(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i], err = s.reconcileState(ctx, items[i])
		if err != nil {
			return nil, err
		}
		items[i], err = s.projectRuntimeCapabilities(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
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

func (s *InstallationService) projectRuntimeCapabilities(ctx context.Context, inst domain.Installation) (domain.Installation, error) {
	if s.runtimeCapabilitySource == nil {
		return inst, nil
	}
	inst.RuntimeCapabilities = s.runtimeCapabilitySource.Project(inst)
	if s.capabilityStates == nil || len(inst.RuntimeCapabilities) == 0 {
		return inst, nil
	}

	declared := make(ports.MarketplaceCapabilities, 0, len(inst.RuntimeCapabilities))
	for _, capability := range inst.RuntimeCapabilities {
		declared = append(declared, string(capability.Code))
	}

	resolved, err := s.capabilityStates.Resolve(ctx, inst.InstallationID, declared)
	if err != nil {
		return domain.Installation{}, err
	}
	inst.RuntimeCapabilities = overlayRuntimeCapabilityStates(inst.RuntimeCapabilities, resolved)
	return inst, nil
}

func (s *InstallationService) reconcileState(ctx context.Context, inst domain.Installation) (domain.Installation, error) {
	if s.stateReconciler == nil {
		return inst, nil
	}
	return s.stateReconciler.Reconcile(ctx, inst)
}
