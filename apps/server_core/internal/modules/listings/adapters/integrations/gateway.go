package integrations

import (
	"context"

	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type installationService interface {
	Get(context.Context, string) (integrationsdomain.Installation, bool, error)
}
type operationService interface {
	BeginExclusive(context.Context, integrationsapp.RecordOperationInput) (integrationsdomain.OperationRun, bool, error)
	Record(context.Context, integrationsapp.RecordOperationInput) (integrationsdomain.OperationRun, error)
}

type Gateway struct {
	installations installationService
	operations    operationService
}

func NewGateway(installations installationService, operations operationService) Gateway {
	return Gateway{installations: installations, operations: operations}
}

func (g Gateway) InstallationAccount(ctx context.Context, id string) (listingsports.InstallationAccount, bool, error) {
	inst, found, err := g.installations.Get(ctx, id)
	if err != nil || !found {
		return listingsports.InstallationAccount{}, found, err
	}
	return listingsports.InstallationAccount{TenantID: inst.TenantID, InstallationID: inst.InstallationID, ProviderCode: inst.ProviderCode, ProviderAccountID: inst.ExternalAccountID}, true, nil
}

func (g Gateway) BeginRefresh(ctx context.Context, run listingsports.RefreshRun) (listingsports.RefreshRun, bool, error) {
	got, active, err := g.operations.BeginExclusive(ctx, integrationsapp.RecordOperationInput{
		OperationRunID: run.OperationRunID, InstallationID: run.InstallationID, OperationType: listingsports.ListingsRefreshOperationType,
		Status: integrationsdomain.OperationRunStatusQueued, ResultCode: "LISTINGS_REFRESH_QUEUED", ActorType: "system", ActorID: "listings_refresh",
	})
	return mapRun(got), active, err
}

func (g Gateway) RecordRefresh(ctx context.Context, run listingsports.RefreshRun) error {
	_, err := g.operations.Record(ctx, integrationsapp.RecordOperationInput{
		OperationRunID: run.OperationRunID, InstallationID: run.InstallationID, OperationType: listingsports.ListingsRefreshOperationType,
		Status: integrationsdomain.OperationRunStatus(run.Status), ResultCode: resultCode(run.Status), FailureCode: failureCode(run.Status),
		TranslatedErrorCode: run.TranslatedErrorCode, AttemptCount: 1, ActorType: "system", ActorID: "listings_refresh",
		StartedAt: run.StartedAt, CompletedAt: run.CompletedAt,
	})
	return err
}

func mapRun(run integrationsdomain.OperationRun) listingsports.RefreshRun {
	return listingsports.RefreshRun{OperationRunID: run.OperationRunID, InstallationID: run.InstallationID, Status: string(run.Status), TranslatedErrorCode: run.TranslatedErrorCode, StartedAt: run.StartedAt, CompletedAt: run.CompletedAt}
}
func resultCode(status string) string {
	if status == "succeeded" {
		return "LISTINGS_REFRESH_SUCCEEDED"
	}
	if status == "running" {
		return "LISTINGS_REFRESH_RUNNING"
	}
	return ""
}
func failureCode(status string) string {
	if status == "failed" {
		return "LISTINGS_REFRESH_FAILED"
	}
	return ""
}

var _ listingsports.RefreshIntegrationGateway = Gateway{}
