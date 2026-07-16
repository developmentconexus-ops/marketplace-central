package integrations

import (
	"context"
	integrationdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	listingports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type Service interface {
	Get(context.Context, string) (integrationdomain.Installation, bool, error)
}
type InstallationReader struct{ service Service }

func NewInstallationReader(service Service) InstallationReader {
	return InstallationReader{service: service}
}
func (a InstallationReader) InstallationExists(ctx context.Context, id string) (bool, error) {
	_, found, err := a.service.Get(ctx, id)
	return found, err
}

var _ listingports.InstallationReader = InstallationReader{}
