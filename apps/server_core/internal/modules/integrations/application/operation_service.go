package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

type RecordOperationInput struct {
	OperationRunID      string
	InstallationID      string
	OperationType       string
	Status              domain.OperationRunStatus
	ResultCode          string
	FailureCode         string
	TranslatedErrorCode string
	AttemptCount        int
	ActorType           string
	ActorID             string
	ProviderEvidence    map[string]any
	DurationMs          int64
	StartedAt           *time.Time
	CompletedAt         *time.Time
}

type OperationService struct {
	store    ports.OperationRunStore
	tenantID string
}

func NewOperationService(store ports.OperationRunStore, tenantID string) *OperationService {
	return &OperationService{store: store, tenantID: tenantID}
}

func (s *OperationService) Record(ctx context.Context, input RecordOperationInput) (domain.OperationRun, error) {
	if input.OperationRunID == "" || input.InstallationID == "" || input.OperationType == "" || !isValidOperationRunStatus(input.Status) || input.AttemptCount < 0 {
		return domain.OperationRun{}, errors.New("INTEGRATIONS_OPERATION_INVALID")
	}

	now := time.Now().UTC()
	run := domain.OperationRun{
		OperationRunID:      input.OperationRunID,
		TenantID:            s.tenantID,
		InstallationID:      input.InstallationID,
		OperationType:       input.OperationType,
		Status:              input.Status,
		ResultCode:          input.ResultCode,
		FailureCode:         input.FailureCode,
		TranslatedErrorCode: input.TranslatedErrorCode,
		AttemptCount:        input.AttemptCount,
		ActorType:           input.ActorType,
		ActorID:             input.ActorID,
		ProviderEvidence:    cloneMap(input.ProviderEvidence),
		DurationMs:          input.DurationMs,
		StartedAt:           cloneTimePtr(input.StartedAt),
		CompletedAt:         cloneTimePtr(input.CompletedAt),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	return run, s.store.SaveOperationRun(ctx, run)
}

func (s *OperationService) ListByInstallation(ctx context.Context, installationID string) ([]domain.OperationRun, error) {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return nil, errors.New("INTEGRATIONS_OPERATION_INVALID")
	}

	return s.store.ListByInstallation(ctx, installationID)
}

func isValidOperationRunStatus(status domain.OperationRunStatus) bool {
	switch status {
	case domain.OperationRunStatusQueued,
		domain.OperationRunStatusRunning,
		domain.OperationRunStatusSucceeded,
		domain.OperationRunStatusFailed,
		domain.OperationRunStatusCancelled:
		return true
	default:
		return false
	}
}

func cloneMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return nil
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
