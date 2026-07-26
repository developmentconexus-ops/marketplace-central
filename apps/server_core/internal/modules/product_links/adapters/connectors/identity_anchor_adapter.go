package connectors

import (
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	productlinksports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type IdentityAnchorAdapter struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
}

func NewIdentityAnchorAdapter(capabilities *connectorsapp.MarketplaceCapabilityService) IdentityAnchorAdapter {
	return IdentityAnchorAdapter{capabilities: capabilities}
}

func (a IdentityAnchorAdapter) ProviderIdentityAnchors(providerCode string) ([]productlinksports.ProviderIdentityAnchor, error) {
	declared, err := a.capabilities.IdentityAnchors(providerCode)
	if err != nil {
		return nil, err
	}

	declaredSet := make(map[connectorsports.IdentityAnchor]struct{}, len(declared))
	for _, anchor := range declared {
		declaredSet[anchor] = struct{}{}
	}

	anchors := make([]productlinksports.ProviderIdentityAnchor, 0, len(connectorsports.KnownIdentityAnchors()))
	for _, anchor := range connectorsports.KnownIdentityAnchors() {
		_, supplied := declaredSet[anchor]
		anchors = append(anchors, productlinksports.ProviderIdentityAnchor{
			Anchor:   string(anchor),
			Supplied: supplied,
		})
	}
	return anchors, nil
}

var _ productlinksports.ProviderIdentityAnchorReader = IdentityAnchorAdapter{}
