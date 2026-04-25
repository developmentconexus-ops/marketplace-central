package madeiramadeira

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
		ProviderCode: "madeira_madeira",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Madeira Madeira",
		AuthStrategy: domain.AuthStrategyToken,
		InstallMode:  domain.InstallModeManual,
		Metadata: map[string]any{
			"country":                     "BR",
			"rollout_stage":               "wave_2",
			"execution_mode":              "available",
			"fee_source":                  "seed",
			"baseline_commission_percent": 0.15,
			"baseline_fixed_fee_amount":   0.0,
			"credential_schema": []map[string]any{
				{"key": "api_token", "label": "API Token", "secret": true},
			},
		},
		DeclaredCapabilities: []string{
			"pricing_fee_sync",
			"order_read",
			"shipment_tracking",
		},
		IsActive: true,
	})
	integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{})
	})
}

func NewAdapter(Config) *Adapter {
	return &Adapter{}
}

func (a *Adapter) ProviderCode() string { return "madeira_madeira" }

func (a *Adapter) AuthStrategy() domain.AuthStrategy { return domain.AuthStrategyToken }

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
	token := strings.TrimSpace(input.APIKey)
	if token == "" {
		token = strings.TrimSpace(input.Metadata["api_token"])
	}
	if token == "" {
		return application.CredentialPayload{}, domain.ErrAPIKeyValidationFailed
	}

	return application.CredentialPayload{
		SecretType:          "token",
		APIKey:              token,
		ProviderAccountID:   strings.TrimSpace(input.Metadata["seller_id"]),
		ProviderAccountName: strings.TrimSpace(input.Metadata["seller_name"]),
	}, nil
}
