package shopee

import (
	"bytes"
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
	PartnerID  string
	PartnerKey string
	BaseURL    string
	HTTPClient *http.Client
	Now        func() time.Time
}

type Adapter struct {
	cfg Config
}

func init() {
	integrationsproviders.RegisterDefinition(domain.ProviderDefinition{
		ProviderCode: "shopee",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Shopee",
		AuthStrategy: domain.AuthStrategyShopeePartner,
		InstallMode:  domain.InstallModeInteractive,
		Metadata: map[string]any{
			"country":                     "BR",
			"rollout_stage":               "v1",
			"execution_mode":              "available",
			"fee_source":                  "seed",
			"baseline_commission_percent": 0.14,
			"baseline_fixed_fee_amount":   0.0,
			"credential_schema": []map[string]any{
				{"key": "partner_id", "label": "Partner ID", "secret": false},
				{"key": "partner_key", "label": "Partner Key", "secret": true},
				{"key": "shop_id", "label": "Shop ID", "secret": false},
			},
		},
		DeclaredCapabilities: []string{"pricing_fee_sync"},
		IsActive:             true,
	})
	integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{
			PartnerID:  strings.TrimSpace(os.Getenv("MPC_PROVIDER_SHOPEE_PARTNER_ID")),
			PartnerKey: strings.TrimSpace(os.Getenv("MPC_PROVIDER_SHOPEE_PARTNER_KEY")),
			BaseURL:    strings.TrimSpace(os.Getenv("MPC_PROVIDER_SHOPEE_BASE_URL")),
		})
	})
}

func NewAdapter(cfg Config) *Adapter {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 15 * time.Second}
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://partner.shopeemobile.com"
	}
	if cfg.Now == nil {
		cfg.Now = func() time.Time { return time.Now().UTC() }
	}
	return &Adapter{cfg: cfg}
}

func (a *Adapter) ProviderCode() string { return "shopee" }

func (a *Adapter) AuthStrategy() domain.AuthStrategy { return domain.AuthStrategyShopeePartner }

func (a *Adapter) StartAuthorize(_ context.Context, input application.StartAuthorizeAdapterInput) (application.AuthorizeStart, error) {
	authURL, err := a.BuildAuthorizeURL(input.State, input.RedirectURI, input.CodeChallenge, input.Scopes)
	if err != nil {
		return application.AuthorizeStart{}, err
	}
	return application.AuthorizeStart{AuthURL: authURL}, nil
}

func (a *Adapter) BuildAuthorizeURL(state, redirectURI, _ string, _ []string) (string, error) {
	if err := a.validateAuthConfig(); err != nil {
		return "", err
	}

	base, err := url.Parse(strings.TrimSpace(a.cfg.BaseURL))
	if err != nil {
		return "", err
	}
	base.Path = "/api/v2/shop/auth_partner"

	timestamp := a.cfg.Now().UTC().Unix()
	sign, err := signShopeeRequest(a.cfg.PartnerID, a.cfg.PartnerKey, base.Path, timestamp, "")
	if err != nil {
		return "", err
	}

	query := base.Query()
	query.Set("partner_id", strings.TrimSpace(a.cfg.PartnerID))
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)
	query.Set("redirect", strings.TrimSpace(redirectURI))
	query.Set("state", strings.TrimSpace(state))

	base.RawQuery = query.Encode()
	return base.String(), nil
}

func (a *Adapter) ExchangeCallback(ctx context.Context, input application.HandleCallbackAdapterInput) (application.CredentialPayload, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" {
		return application.CredentialPayload{}, domain.ErrAuthCodeExchangeFailed
	}
	shopID := strings.TrimSpace(input.ProviderAccountID)
	if shopID == "" {
		return application.CredentialPayload{}, domain.ErrCredentialValidationFailed
	}

	result, err := a.exchangeToken(ctx, tokenExchangeRequest{
		Path:   "/api/v2/auth/token/get",
		Code:   code,
		ShopID: shopID,
	})
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := a.cfg.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}

	extra := result.RawExtras
	if extra == nil {
		extra = map[string]any{}
	}
	extra["shop_id"] = shopID

	return application.CredentialPayload{
		SecretType:        "shopee_partner",
		AccessToken:       result.AccessToken,
		RefreshToken:      result.RefreshToken,
		ProviderAccountID: shopID,
		ExpiresAt:         expiresAt,
		Extra:             extra,
	}, nil
}

func (a *Adapter) Refresh(ctx context.Context, input application.RefreshCredentialAdapterInput) (application.CredentialPayload, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return application.CredentialPayload{}, domain.ErrRefreshTokenInvalid
	}
	shopID := strings.TrimSpace(input.ProviderAccountID)
	if shopID == "" {
		return application.CredentialPayload{}, domain.ErrCredentialValidationFailed
	}

	result, err := a.exchangeToken(ctx, tokenExchangeRequest{
		Path:         "/api/v2/auth/access_token/get",
		RefreshToken: refreshToken,
		ShopID:       shopID,
	})
	if err != nil {
		return application.CredentialPayload{}, err
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		ts := a.cfg.Now().UTC().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &ts
	}

	rotatedRefreshToken := strings.TrimSpace(result.RefreshToken)
	if rotatedRefreshToken == "" {
		rotatedRefreshToken = refreshToken
	}

	extra := result.RawExtras
	if extra == nil {
		extra = map[string]any{}
	}
	extra["shop_id"] = shopID

	return application.CredentialPayload{
		SecretType:        "shopee_partner",
		AccessToken:       result.AccessToken,
		RefreshToken:      rotatedRefreshToken,
		ProviderAccountID: shopID,
		ExpiresAt:         expiresAt,
		Extra:             extra,
	}, nil
}

func (a *Adapter) VerifyAPIKey(context.Context, application.SubmitAPIKeyAdapterInput) (application.CredentialPayload, error) {
	return application.CredentialPayload{}, domain.ErrNotSupported
}

func (a *Adapter) ValidateCredentials(context.Context, map[string]string) (*domain.TokenResult, error) {
	return nil, domain.ErrNotSupported
}

func (a *Adapter) ExchangeCode(context.Context, string, string, string) (*domain.TokenResult, error) {
	return nil, domain.ErrNotSupported
}

func (a *Adapter) RefreshToken(context.Context, string) (*domain.TokenResult, error) {
	return nil, domain.ErrNotSupported
}

func (a *Adapter) RevokeToken(context.Context, string) error {
	return nil
}

type tokenExchangeRequest struct {
	Path         string
	Code         string
	RefreshToken string
	ShopID       string
}

func (a *Adapter) exchangeToken(ctx context.Context, req tokenExchangeRequest) (*domain.TokenResult, error) {
	if err := a.validateAuthConfig(); err != nil {
		return nil, err
	}

	timestamp := a.cfg.Now().UTC().Unix()
	sign, err := signShopeeRequest(a.cfg.PartnerID, a.cfg.PartnerKey, req.Path, timestamp, req.ShopID)
	if err != nil {
		return nil, err
	}

	payload := map[string]string{
		"partner_id": strings.TrimSpace(a.cfg.PartnerID),
		"shop_id":    strings.TrimSpace(req.ShopID),
	}
	switch req.Path {
	case "/api/v2/auth/token/get":
		payload["code"] = strings.TrimSpace(req.Code)
	case "/api/v2/auth/access_token/get":
		payload["refresh_token"] = strings.TrimSpace(req.RefreshToken)
	default:
		return nil, domain.ErrAuthConfigurationInvalid
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(strings.TrimSpace(a.cfg.BaseURL))
	if err != nil {
		return nil, err
	}
	endpoint.Path = req.Path
	query := endpoint.Query()
	query.Set("partner_id", strings.TrimSpace(a.cfg.PartnerID))
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)
	endpoint.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := a.cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAuthProviderUnreachable, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, mapShopeeHTTPError(req.Path, resp.StatusCode, raw)
	}

	var payloadResponse shopeeTokenResponse
	if err := json.Unmarshal(raw, &payloadResponse); err != nil {
		return nil, err
	}
	if payloadResponse.Error != "" {
		return nil, mapShopeeAPIError(req.Path, payloadResponse, raw)
	}

	return &domain.TokenResult{
		AccessToken:  strings.TrimSpace(payloadResponse.Response.AccessToken),
		RefreshToken: strings.TrimSpace(payloadResponse.Response.RefreshToken),
		ExpiresIn:    payloadResponse.Response.ExpireIn,
		TokenType:    "Bearer",
		ProviderAccountID: func() string {
			if value := strings.TrimSpace(payloadResponse.Response.ShopID); value != "" {
				return value
			}
			return req.ShopID
		}(),
		RawExtras: map[string]any{
			"request_id": payloadResponse.RequestID,
		},
	}, nil
}

func (a *Adapter) validateAuthConfig() error {
	if strings.TrimSpace(a.cfg.PartnerID) == "" || strings.TrimSpace(a.cfg.PartnerKey) == "" {
		return domain.ErrAuthConfigurationInvalid
	}
	if strings.TrimSpace(a.cfg.BaseURL) == "" {
		return domain.ErrAuthConfigurationInvalid
	}
	return nil
}

func mapShopeeHTTPError(path string, status int, body []byte) error {
	apiErr, _ := decodeShopeeError(body)
	return mapShopeeError(path, status, apiErr, body)
}

func mapShopeeAPIError(path string, status shopeeTokenResponse, body []byte) error {
	return mapShopeeError(path, 400, shopeeAPIError{
		Error:   status.Error,
		Message: status.Message,
	}, body)
}

func mapShopeeError(path string, status int, apiErr shopeeAPIError, body []byte) error {
	code := strings.ToLower(strings.TrimSpace(apiErr.Error))
	message := strings.ToLower(strings.TrimSpace(apiErr.Message))
	bodyText := strings.TrimSpace(string(body))
	switch {
	case strings.Contains(code, "sign"), strings.Contains(message, "sign"), strings.Contains(message, "partner"), strings.Contains(message, "timestamp"):
		return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrAuthConfigurationInvalid, status, apiErr.Error, apiErr.Message, bodyText)
	case path == "/api/v2/auth/token/get" && (strings.Contains(code, "code") || strings.Contains(message, "code")):
		return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrAuthCodeExchangeFailed, status, apiErr.Error, apiErr.Message, bodyText)
	case path == "/api/v2/auth/access_token/get" && (strings.Contains(code, "refresh") || strings.Contains(message, "refresh")):
		return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrRefreshTokenInvalid, status, apiErr.Error, apiErr.Message, bodyText)
	case status == http.StatusTooManyRequests:
		return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrRefreshRateLimited, status, apiErr.Error, apiErr.Message, bodyText)
	default:
		if path == "/api/v2/auth/access_token/get" {
			return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrRefreshProviderError, status, apiErr.Error, apiErr.Message, bodyText)
		}
		return fmt.Errorf("%w: status=%d error=%s message=%s body=%s", domain.ErrAuthCodeExchangeFailed, status, apiErr.Error, apiErr.Message, bodyText)
	}
}

func decodeShopeeError(body []byte) (shopeeAPIError, error) {
	var apiErr shopeeAPIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return shopeeAPIError{}, err
	}
	return apiErr, nil
}
