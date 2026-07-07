package leroymerlin

import (
	"context"
	"strings"

	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type Config struct{}

type Adapter struct{}

func init() {
	integrationsproviders.RegisterDefinition(domain.ProviderDefinition{
		ProviderCode: "leroy_merlin",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Leroy Merlin",
		AuthStrategy: domain.AuthStrategyAPIKey,
		InstallMode:  domain.InstallModeManual,
		Metadata: map[string]any{
			"country":                     "BR",
			"rollout_stage":               "wave_2",
			"execution_mode":              "available",
			"fee_source":                  "seed",
			"baseline_commission_percent": 0.18,
			"baseline_fixed_fee_amount":   0.0,
			"credential_schema": []map[string]any{
				{"key": "api_key", "label": "API Key", "secret": true},
				{"key": "shop_id", "label": "Shop ID", "secret": false},
			},
		},
		DeclaredCapabilities: []string{},
		IsActive:             true,
	})
	integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{})
	})
}

func NewAdapter(Config) *Adapter {
	return &Adapter{}
}

func (a *Adapter) ProviderCode() string { return "leroy_merlin" }

func (a *Adapter) AuthStrategy() domain.AuthStrategy { return domain.AuthStrategyAPIKey }

func (a *Adapter) StartAuthorize(context.Context, application.StartAuthorizeAdapterInput) (application.AuthorizeStart, error) {
	return application.AuthorizeStart{}, domain.ErrNotSupported
}

func (a *Adapter) ExchangeCallback(context.Context, application.HandleCallbackAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) Refresh(context.Context, application.RefreshCredentialAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) VerifyAPIKey(_ context.Context, input application.SubmitAPIKeyAdapterInput) (application.CredentialPayload, error) {
	if strings.TrimSpace(input.APIKey) == "" {
		return application.CredentialPayload{}, domain.ErrAPIKeyValidationFailed
	}

	return application.CredentialPayload{
		SecretType:          "api_key",
		APIKey:              input.APIKey,
		ProviderAccountID:   strings.TrimSpace(input.Metadata["shop_id"]),
		ProviderAccountName: strings.TrimSpace(input.Metadata["shop_name"]),
	}, nil
}
