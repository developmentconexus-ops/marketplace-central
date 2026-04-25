package amazon

import (
	"context"
	"net/url"
	"os"
	"strings"

	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type Config struct {
	ClientID     string
	ClientSecret string
	AuthorizeURL string
	TokenURL     string
}

type Adapter struct {
	cfg Config
}

func init() {
	integrationsproviders.RegisterDefinition(domain.ProviderDefinition{
		ProviderCode: "amazon",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Amazon Brasil",
		AuthStrategy: domain.AuthStrategyLWA,
		InstallMode:  domain.InstallModeInteractive,
		Metadata: map[string]any{
			"country":                     "BR",
			"rollout_stage":               "v1",
			"execution_mode":              "available",
			"fee_source":                  "seed",
			"baseline_commission_percent": 0.12,
			"baseline_fixed_fee_amount":   0.0,
			"credential_schema": []map[string]any{
				{"key": "seller_id", "label": "Seller ID", "secret": false},
				{"key": "lwa_app_id", "label": "LWA App ID", "secret": false},
				{"key": "lwa_client_secret", "label": "LWA Client Secret", "secret": true},
				{"key": "refresh_token", "label": "Refresh Token", "secret": true},
			},
		},
		DeclaredCapabilities: []string{
			"catalog_publish",
			"pricing_fee_sync",
			"inventory_sync",
			"order_read",
			"message_read",
			"shipment_tracking",
			"webhook_receive",
		},
		IsActive: true,
	})
	integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{
			ClientID:     strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_CLIENT_SECRET")),
			AuthorizeURL: "https://sellercentral.amazon.com.br/apps/authorize/consent",
			TokenURL:     "https://api.amazon.com/auth/o2/token",
		})
	})
}

func NewAdapter(cfg Config) *Adapter {
	return &Adapter{cfg: cfg}
}

func (a *Adapter) ProviderCode() string { return "amazon" }

func (a *Adapter) AuthStrategy() domain.AuthStrategy { return domain.AuthStrategyLWA }

func (a *Adapter) StartAuthorize(_ context.Context, input application.StartAuthorizeAdapterInput) (application.AuthorizeStart, error) {
	authURL, err := a.BuildAuthorizeURL(input.State, input.RedirectURI, input.CodeChallenge, input.Scopes)
	if err != nil {
		return application.AuthorizeStart{}, err
	}
	return application.AuthorizeStart{AuthURL: authURL}, nil
}

func (a *Adapter) BuildAuthorizeURL(state, redirectURI, codeChallenge string, scopes []string) (string, error) {
	base, err := url.Parse(strings.TrimSpace(a.cfg.AuthorizeURL))
	if err != nil {
		return "", err
	}

	query := base.Query()
	query.Set("response_type", "code")
	query.Set("client_id", strings.TrimSpace(a.cfg.ClientID))
	query.Set("redirect_uri", strings.TrimSpace(redirectURI))
	query.Set("state", strings.TrimSpace(state))

	if scope := joinScopes(scopes); scope != "" {
		query.Set("scope", scope)
	}
	if strings.TrimSpace(codeChallenge) != "" {
		query.Set("code_challenge", strings.TrimSpace(codeChallenge))
		query.Set("code_challenge_method", "S256")
	}

	base.RawQuery = query.Encode()
	return base.String(), nil
}

func (a *Adapter) ExchangeCallback(context.Context, application.HandleCallbackAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) VerifyAPIKey(context.Context, application.SubmitAPIKeyAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) Refresh(context.Context, application.RefreshCredentialAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func joinScopes(scopes []string) string {
	filtered := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		filtered = append(filtered, scope)
	}
	return strings.Join(filtered, " ")
}
