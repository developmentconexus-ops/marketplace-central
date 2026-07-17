package marketplaces

import (
	"context"
	"errors"

	marketplacesdomain "marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type pricingPolicyService interface {
	GetPricingPolicyForInstallation(context.Context, string) (marketplacesdomain.Policy, bool, error)
}

type PolicyReader struct {
	service pricingPolicyService
}

func NewPolicyReader(service pricingPolicyService) *PolicyReader {
	return &PolicyReader{service: service}
}

var _ mutationsports.PolicyReader = (*PolicyReader)(nil)

func (r *PolicyReader) GetPricingPolicyForInstallation(ctx context.Context, installationID string) (marketplacesdomain.Policy, bool, error) {
	if r == nil || r.service == nil {
		return marketplacesdomain.Policy{}, false, errors.New("pricing policy reader is not configured")
	}
	policy, found, err := r.service.GetPricingPolicyForInstallation(ctx, installationID)
	if err != nil {
		return marketplacesdomain.Policy{}, false, err
	}
	return policy, found, nil
}

func (r *PolicyReader) HasPricingPolicy(ctx context.Context, installationID string) (bool, error) {
	_, found, err := r.GetPricingPolicyForInstallation(ctx, installationID)
	return found, err
}
