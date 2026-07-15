package marketplaces

import (
	"context"

	listingports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	marketdomain "marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
)

type Service interface {
	GetPricingPolicyForInstallation(context.Context, string) (marketdomain.Policy, bool, error)
}
type PolicyReader struct{ service Service }

func NewPolicyReader(service Service) PolicyReader { return PolicyReader{service: service} }
func (a PolicyReader) GetPricingPolicyForInstallation(ctx context.Context, id string) (listingports.PricingPolicy, bool, error) {
	p, ok, err := a.service.GetPricingPolicyForInstallation(ctx, id)
	return listingports.PricingPolicy{MinMarginPercent: p.MinMarginPercent}, ok, err
}

var _ listingports.PolicyReader = PolicyReader{}
