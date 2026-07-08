package application

import (
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type RuntimeCapabilityProjector struct {
	marketplaceCapabilities *connectorsapp.MarketplaceCapabilityService
}

func NewRuntimeCapabilityProjector(marketplaceCapabilities *connectorsapp.MarketplaceCapabilityService) *RuntimeCapabilityProjector {
	return &RuntimeCapabilityProjector{marketplaceCapabilities: marketplaceCapabilities}
}

func (p *RuntimeCapabilityProjector) Project(inst domain.Installation) []domain.RuntimeCapability {
	if p == nil || p.marketplaceCapabilities == nil {
		return []domain.RuntimeCapability{}
	}

	capabilities := make([]domain.RuntimeCapability, 0, 5)
	capabilities = appendCapability(capabilities, domain.RuntimeCapabilityAccountProbe, p.hasAccountProbe(inst.ProviderCode), inst)
	capabilities = appendCapability(capabilities, domain.RuntimeCapabilityListingRead, p.hasListingRead(inst.ProviderCode), inst)
	capabilities = appendCapability(capabilities, domain.RuntimeCapabilityOrderRead, p.hasOrderRead(inst.ProviderCode), inst)
	capabilities = appendCapability(capabilities, domain.RuntimeCapabilityFeeQuoteRead, p.hasFeeQuoteRead(inst.ProviderCode), inst)
	capabilities = appendCapability(capabilities, domain.RuntimeCapabilityStockRead, p.hasStockRead(inst.ProviderCode), inst)
	return capabilities
}

func appendCapability(items []domain.RuntimeCapability, code domain.RuntimeCapabilityCode, supported bool, inst domain.Installation) []domain.RuntimeCapability {
	if !supported {
		return items
	}

	capability := domain.RuntimeCapability{
		Code:           code,
		LocalValidated: false,
		LiveValidated:  false,
	}

	switch inst.Status {
	case domain.InstallationStatusConnected:
		capability.State = domain.RuntimeCapabilityStateAvailable
		capability.Executable = inst.ExternalAccountID != ""
		if !capability.Executable {
			capability.State = domain.RuntimeCapabilityStateNeedsAuth
			capability.UnavailableReason = "provider_account_unverified"
		}
	case domain.InstallationStatusRequiresReauth, domain.InstallationStatusPendingConnection, domain.InstallationStatusDisconnected, domain.InstallationStatusDraft:
		capability.State = domain.RuntimeCapabilityStateNeedsAuth
		capability.Executable = false
		capability.UnavailableReason = "connection_not_ready"
	default:
		capability.State = domain.RuntimeCapabilityStateDegraded
		capability.Executable = false
		capability.UnavailableReason = "connection_degraded"
	}

	capability.LastValidatedAt = inst.LastVerifiedAt
	return append(items, capability)
}

func (p *RuntimeCapabilityProjector) hasAccountProbe(providerCode string) bool {
	_, err := p.marketplaceCapabilities.AccountProber(providerCode)
	return err == nil
}

func (p *RuntimeCapabilityProjector) hasListingRead(providerCode string) bool {
	_, err := p.marketplaceCapabilities.ListingReader(providerCode)
	return err == nil
}

func (p *RuntimeCapabilityProjector) hasOrderRead(providerCode string) bool {
	_, err := p.marketplaceCapabilities.OrderReader(providerCode)
	return err == nil
}

func (p *RuntimeCapabilityProjector) hasFeeQuoteRead(providerCode string) bool {
	_, err := p.marketplaceCapabilities.FeeQuoteReader(providerCode)
	return err == nil
}

func (p *RuntimeCapabilityProjector) hasStockRead(providerCode string) bool {
	_, err := p.marketplaceCapabilities.StockReader(providerCode)
	return err == nil
}
