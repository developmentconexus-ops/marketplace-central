package integrations

import (
	"context"

	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type InstallationReader struct {
	service *integrationsapp.InstallationService
}

func NewInstallationReader(service *integrationsapp.InstallationService) InstallationReader {
	return InstallationReader{service: service}
}

func (r InstallationReader) GetInstallation(ctx context.Context, installationID string) (inventorydomain.ProviderStockRef, bool, error) {
	inst, found, err := r.service.Get(ctx, installationID)
	if err != nil || !found {
		return inventorydomain.ProviderStockRef{}, found, err
	}
	return inventorydomain.ProviderStockRef{
		TenantID:          inst.TenantID,
		InstallationID:    inst.InstallationID,
		ProviderCode:      inst.ProviderCode,
		ProviderAccountID: inst.ExternalAccountID,
	}, true, nil
}
