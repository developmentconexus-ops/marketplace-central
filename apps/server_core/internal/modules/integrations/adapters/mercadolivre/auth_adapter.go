package mercadolivre

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	integrationsproviders "marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type Config struct {
	ClientID     string
	ClientSecret string
	AuthorizeURL string
	TokenURL     string
	APIBaseURL   string
	HTTPClient   *http.Client
}

type Adapter struct {
	cfg Config
}

func init() {
	integrationsproviders.RegisterDefinition(domain.ProviderDefinition{
		ProviderCode: "mercado_livre",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Mercado Livre",
		AuthStrategy: domain.AuthStrategyOAuth2,
		InstallMode:  domain.InstallModeInteractive,
		Metadata: map[string]any{
			"country":                     "BR",
			"rollout_stage":               "v1",
			"execution_mode":              "available",
			"fee_source":                  "api_sync",
			"baseline_commission_percent": 0.16,
			"baseline_fixed_fee_amount":   0.0,
			"credential_schema": []map[string]any{
				{"key": "client_id", "label": "Client ID", "secret": false},
				{"key": "client_secret", "label": "Client Secret", "secret": true},
				{"key": "redirect_uri", "label": "Redirect URI", "secret": false},
			},
		},
		DeclaredCapabilities: []string{
			"listing_read",
			"pricing_fee_sync",
			"stock_read",
			"stock_write",
			"order_read",
			"message_read",
			"message_reply",
			"shipment_tracking",
			"webhook_receive",
		},
		IsActive: true,
	})
	integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{
			ClientID:     strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET")),
			AuthorizeURL: "https://auth.mercadolivre.com.br/authorization",
			TokenURL:     "https://api.mercadolibre.com/oauth/token",
			APIBaseURL:   "https://api.mercadolibre.com",
		})
	})
}

func NewAdapter(cfg Config) *Adapter {
	if strings.TrimSpace(cfg.APIBaseURL) == "" {
		cfg.APIBaseURL = "https://api.mercadolibre.com"
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Adapter{cfg: cfg}
}

func (a *Adapter) ProviderCode() string {
	return "mercado_livre"
}

func (a *Adapter) AuthStrategy() domain.AuthStrategy {
	return domain.AuthStrategyOAuth2
}

func (a *Adapter) StartAuthorize(_ context.Context, input application.StartAuthorizeAdapterInput) (application.AuthorizeStart, error) {
	authURL, err := a.BuildAuthorizeURL(input.State, input.RedirectURI, input.CodeChallenge)
	if err != nil {
		return application.AuthorizeStart{}, err
	}
	return application.AuthorizeStart{AuthURL: authURL}, nil
}

func (a *Adapter) BuildAuthorizeURL(state, redirectURI, codeChallenge string) (string, error) {
	base, err := url.Parse(a.cfg.AuthorizeURL)
	if err != nil {
		return "", err
	}
	q := base.Query()
	q.Set("response_type", "code")
	q.Set("client_id", a.cfg.ClientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("state", state)
	if codeChallenge != "" {
		q.Set("code_challenge", codeChallenge)
		q.Set("code_challenge_method", "S256")
	}
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func (a *Adapter) ExchangeCallback(ctx context.Context, input application.HandleCallbackAdapterInput) (application.CredentialPayload, error) {
	result, err := a.ExchangeCode(ctx, input.Code, input.RedirectURI, input.CodeVerifier)
	if err != nil {
		return application.CredentialPayload{}, err
	}
	profile, err := a.lookupAccount(ctx, result.AccessToken, result.ProviderAccountID)
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := time.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}

	return application.CredentialPayload{
		SecretType:          "oauth2",
		AccessToken:         result.AccessToken,
		RefreshToken:        result.RefreshToken,
		ProviderAccountID:   profile.ProviderAccountID,
		ProviderAccountName: profile.ProviderAccountName,
		ExpiresAt:           expiresAt,
		Extra:               result.RawExtras,
	}, nil
}

func (a *Adapter) ExchangeCode(ctx context.Context, code, redirectURI, codeVerifier string) (*domain.TokenResult, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", a.cfg.ClientID)
	form.Set("client_secret", a.cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	if strings.TrimSpace(codeVerifier) != "" {
		form.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status=%d body=%s", domain.ErrAuthCodeExchangeFailed, resp.StatusCode, readProviderErrorBody(resp))
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		UserID       any    `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &domain.TokenResult{
		AccessToken:       payload.AccessToken,
		RefreshToken:      payload.RefreshToken,
		ExpiresIn:         payload.ExpiresIn,
		TokenType:         "Bearer",
		ProviderAccountID: normalizeAnyString(payload.UserID),
	}, nil
}

func (a *Adapter) Refresh(ctx context.Context, input application.RefreshCredentialAdapterInput) (application.CredentialPayload, error) {
	result, err := a.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return application.CredentialPayload{}, err
	}
	profile, err := a.lookupAccount(ctx, result.AccessToken, firstNonEmpty(result.ProviderAccountID, input.ProviderAccountID))
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := time.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}

	return application.CredentialPayload{
		SecretType:          "oauth2",
		AccessToken:         result.AccessToken,
		RefreshToken:        result.RefreshToken,
		ProviderAccountID:   profile.ProviderAccountID,
		ProviderAccountName: profile.ProviderAccountName,
		ExpiresAt:           expiresAt,
		Extra:               result.RawExtras,
	}, nil
}

func (a *Adapter) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResult, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", a.cfg.ClientID)
	form.Set("client_secret", a.cfg.ClientSecret)
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body := readProviderErrorBody(resp)
		return nil, fmt.Errorf("%w: status=%d body=%s", classifyRefreshHTTPError(resp.StatusCode, body), resp.StatusCode, body)
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		UserID       any    `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &domain.TokenResult{
		AccessToken:       payload.AccessToken,
		RefreshToken:      payload.RefreshToken,
		ExpiresIn:         payload.ExpiresIn,
		TokenType:         "Bearer",
		ProviderAccountID: normalizeAnyString(payload.UserID),
	}, nil
}

func (a *Adapter) RevokeToken(context.Context, string) error {
	return nil
}

func (a *Adapter) VerifyAPIKey(context.Context, application.SubmitAPIKeyAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) ValidateCredentials(context.Context, map[string]string) (*domain.TokenResult, error) {
	return nil, domain.ErrNotSupported
}

func normalizeAnyString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		if value == nil {
			return ""
		}
		return fmt.Sprintf("%v", value)
	}
}

var errInvalidConfig = errors.New("INTEGRATIONS_AUTH_PROVIDER_UNREACHABLE")

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

// classifyRefreshHTTPError traduz a resposta de erro do endpoint de token do ML
// para a sentinela de domínio correspondente. O vocabulário do provider
// ("invalid_grant") vive aqui e só aqui: o domínio não conhece nome de erro de
// marketplace nenhum (ADR-C4).
//
// Sem essa tradução todo erro >= 400 virava ErrRefreshProviderError, que
// ClassifyRefreshError (domain/refresh_policy.go:52) considera TRANSITÓRIO —
// um refresh token revogado seria retentado para sempre e a conta nunca
// chegaria a requires_reauth.
func classifyRefreshHTTPError(status int, body string) error {
	// O ML responde 400 com {"error":"invalid_grant"} para refresh token
	// inválido, revogado ou já usado. Nenhum retry conserta: só reautorização.
	if strings.Contains(body, "invalid_grant") {
		return domain.ErrRefreshTokenInvalid
	}
	if status == http.StatusTooManyRequests {
		return domain.ErrRefreshRateLimited
	}
	return domain.ErrRefreshProviderError
}

type accountProfile struct {
	ProviderAccountID   string
	ProviderAccountName string
}

func (a *Adapter) lookupAccount(ctx context.Context, accessToken string, fallbackAccountID string) (accountProfile, error) {
	accountID := strings.TrimSpace(fallbackAccountID)
	if strings.TrimSpace(accessToken) == "" {
		return accountProfile{ProviderAccountID: accountID}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(a.cfg.APIBaseURL, "/")+"/users/me", nil)
	if err != nil {
		return accountProfile{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.cfg.HTTPClient.Do(req)
	if err != nil {
		return accountProfile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return accountProfile{}, fmt.Errorf("%w: status=%d body=%s", domain.ErrAuthCodeExchangeFailed, resp.StatusCode, readProviderErrorBody(resp))
	}

	var payload struct {
		ID        any    `json:"id"`
		Nickname  string `json:"nickname"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return accountProfile{}, err
	}

	name := strings.TrimSpace(payload.Nickname)
	if name == "" {
		name = strings.TrimSpace(strings.TrimSpace(payload.FirstName) + " " + strings.TrimSpace(payload.LastName))
	}

	return accountProfile{
		ProviderAccountID:   firstNonEmpty(normalizeAnyString(payload.ID), accountID),
		ProviderAccountName: name,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
