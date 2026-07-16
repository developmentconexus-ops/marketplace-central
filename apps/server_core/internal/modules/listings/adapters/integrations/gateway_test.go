package integrations

import (
	"context"
	"testing"

	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type installationStub struct {
	value integrationsdomain.Installation
}

func (s installationStub) Get(context.Context, string) (integrationsdomain.Installation, bool, error) {
	return s.value, true, nil
}

type operationsStub struct {
	began    integrationsapp.RecordOperationInput
	recorded integrationsapp.RecordOperationInput
}

func (s *operationsStub) BeginExclusive(_ context.Context, in integrationsapp.RecordOperationInput) (integrationsdomain.OperationRun, bool, error) {
	s.began = in
	return integrationsdomain.OperationRun{OperationRunID: in.OperationRunID, InstallationID: in.InstallationID, Status: in.Status}, false, nil
}
func (s *operationsStub) Record(_ context.Context, in integrationsapp.RecordOperationInput) (integrationsdomain.OperationRun, error) {
	s.recorded = in
	return integrationsdomain.OperationRun{}, nil
}

func TestGatewayMapsOnlyPublishedIntegrationBoundary(t *testing.T) {
	ops := &operationsStub{}
	g := NewGateway(installationStub{integrationsdomain.Installation{TenantID: "tenant", InstallationID: "inst", ProviderCode: "provider", ExternalAccountID: "seller"}}, ops)
	account, found, err := g.InstallationAccount(context.Background(), "inst")
	if err != nil || !found || account.ProviderAccountID != "seller" {
		t.Fatalf("account=%+v found=%v err=%v", account, found, err)
	}
	_, _, err = g.BeginRefresh(context.Background(), listingsports.RefreshRun{OperationRunID: "op", InstallationID: "inst"})
	if err != nil || ops.began.OperationType != listingsports.ListingsRefreshOperationType {
		t.Fatalf("input=%+v err=%v", ops.began, err)
	}
}
