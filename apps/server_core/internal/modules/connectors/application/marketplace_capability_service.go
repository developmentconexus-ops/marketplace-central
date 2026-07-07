package application

import (
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

type ProviderCapabilitySet struct {
	ProviderCode string
	Listings     ports.ListingReader
	StockReads   ports.StockReader
	StockWrites  ports.StockWriter
	Orders       ports.OrderReader
}

type MarketplaceCapabilityService struct {
	byProvider map[string]ProviderCapabilitySet
}

func NewMarketplaceCapabilityService(capabilities []ProviderCapabilitySet) *MarketplaceCapabilityService {
	byProvider := make(map[string]ProviderCapabilitySet, len(capabilities))
	for _, capability := range capabilities {
		code := strings.TrimSpace(capability.ProviderCode)
		if code == "" {
			continue
		}
		capability.ProviderCode = code
		byProvider[code] = capability
	}

	return &MarketplaceCapabilityService{byProvider: byProvider}
}

func (s *MarketplaceCapabilityService) ListingReader(providerCode string) (ports.ListingReader, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.Listings == nil {
		return nil, unsupported(providerCode, ports.CapabilityListingRead)
	}
	return capability.Listings, nil
}

func (s *MarketplaceCapabilityService) StockReader(providerCode string) (ports.StockReader, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.StockReads == nil {
		return nil, unsupported(providerCode, ports.CapabilityStockRead)
	}
	return capability.StockReads, nil
}

func (s *MarketplaceCapabilityService) StockWriter(providerCode string) (ports.StockWriter, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.StockWrites == nil {
		return nil, unsupported(providerCode, ports.CapabilityStockWrite)
	}
	return capability.StockWrites, nil
}

func (s *MarketplaceCapabilityService) OrderReader(providerCode string) (ports.OrderReader, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.Orders == nil {
		return nil, unsupported(providerCode, ports.CapabilityOrderRead)
	}
	return capability.Orders, nil
}

func (s *MarketplaceCapabilityService) provider(providerCode string) (ProviderCapabilitySet, error) {
	code := strings.TrimSpace(providerCode)
	if code == "" {
		return ProviderCapabilitySet{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider code is required")
	}

	capability, ok := s.byProvider[code]
	if !ok {
		return ProviderCapabilitySet{}, unsupported(code, "provider")
	}
	return capability, nil
}

func unsupported(providerCode string, operation string) error {
	return domain.NewCapabilityError(
		domain.ErrCodeProviderUnsupportedShape,
		"provider "+strings.TrimSpace(providerCode)+" does not support "+strings.TrimSpace(operation),
	)
}
