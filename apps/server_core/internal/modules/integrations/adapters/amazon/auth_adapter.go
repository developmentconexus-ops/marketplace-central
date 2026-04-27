package amazon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type Config struct {
	ApplicationID string
	ClientID      string
	ClientSecret  string
	AuthorizeURL  string
	TokenURL      string
	AuthVersion   string
	HTTPClient    *http.Client
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
			ApplicationID: strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_APPLICATION_ID")),
			ClientID:      strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_CLIENT_ID")),
			ClientSecret:  strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_CLIENT_SECRET")),
			AuthorizeURL:  "https://sellercentral.amazon.com.br/apps/authorize/consent",
			TokenURL:      "https://api.amazon.com/auth/o2/token",
			AuthVersion:   strings.TrimSpace(os.Getenv("MPC_PROVIDER_AMAZON_AUTH_VERSION")),
		})
	})
}

func NewAdapter(cfg Config) *Adapter {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 15 * time.Second}
	}
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
	applicationID := strings.TrimSpace(a.cfg.ApplicationID)
	if !strings.HasPrefix(applicationID, "amzn1.sellerapps.app.") {
		return "", domain.ErrAuthConfigurationInvalid
	}

	base, err := url.Parse(strings.TrimSpace(a.cfg.AuthorizeURL))
	if err != nil {
		return "", err
	}

	query := base.Query()
	query.Set("response_type", "code")
	query.Set("application_id", applicationID)
	query.Set("redirect_uri", strings.TrimSpace(redirectURI))
	query.Set("state", strings.TrimSpace(state))

	if version := strings.TrimSpace(a.cfg.AuthVersion); version != "" {
		query.Set("version", version)
	}
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

func (a *Adapter) ExchangeCallback(ctx context.Context, input application.HandleCallbackAdapterInput) (application.CredentialPayload, error) {
	result, err := a.exchangeToken(ctx, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {strings.TrimSpace(input.Code)},
		"redirect_uri":  {strings.TrimSpace(input.RedirectURI)},
		"client_id":     {strings.TrimSpace(a.cfg.ClientID)},
		"client_secret": {strings.TrimSpace(a.cfg.ClientSecret)},
	})
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := time.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}

	providerAccountID := strings.TrimSpace(input.ProviderAccountID)
	extra := result.RawExtras
	if extra == nil {
		extra = map[string]any{}
	}
	if providerAccountID != "" {
		extra["selling_partner_id"] = providerAccountID
	}

	return application.CredentialPayload{
		SecretType:        "lwa",
		AccessToken:       result.AccessToken,
		RefreshToken:      result.RefreshToken,
		ProviderAccountID: providerAccountID,
		ExpiresAt:         expiresAt,
		Extra:             extra,
	}, nil
}

func (a *Adapter) VerifyAPIKey(context.Context, application.SubmitAPIKeyAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) Refresh(ctx context.Context, input application.RefreshCredentialAdapterInput) (application.CredentialPayload, error) {
	result, err := a.exchangeToken(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {strings.TrimSpace(input.RefreshToken)},
		"client_id":     {strings.TrimSpace(a.cfg.ClientID)},
		"client_secret": {strings.TrimSpace(a.cfg.ClientSecret)},
	})
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := time.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}
	refreshToken := strings.TrimSpace(result.RefreshToken)
	if refreshToken == "" {
		refreshToken = input.RefreshToken
	}
	return application.CredentialPayload{
		SecretType:   "lwa",
		AccessToken:  result.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Extra:        result.RawExtras,
	}, nil
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

func (a *Adapter) exchangeToken(ctx context.Context, form url.Values) (*domain.TokenResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimSpace(a.cfg.TokenURL), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status=%d body=%s", domain.ErrAuthCodeExchangeFailed, resp.StatusCode, readProviderErrorBody(resp))
	}

	var payload struct {
		AccessToken  string         `json:"access_token"`
		RefreshToken string         `json:"refresh_token"`
		ExpiresIn    int            `json:"expires_in"`
		TokenType    string         `json:"token_type"`
		Raw          map[string]any `json:"-"`
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	rawExtras := map[string]any{}
	_ = json.Unmarshal(raw, &rawExtras)

	return &domain.TokenResult{
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
		ExpiresIn:    payload.ExpiresIn,
		TokenType:    payload.TokenType,
		RawExtras:    rawExtras,
	}, nil
}

func readProviderErrorBody(resp *http.Response) string {
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "unavailable"
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return "empty"
	}
	return text
}
