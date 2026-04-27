package shopee

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestAdapterBuildAuthorizeURLUsesSignedPartnerFlow(t *testing.T) {
	t.Parallel()

	now := time.Unix(1594897040, 0).UTC()
	adapter := NewAdapter(Config{
		PartnerID:  "10090",
		PartnerKey: "test-partner-key",
		BaseURL:    "https://partner.test-stable.shopeemobile.com",
		Now:        func() time.Time { return now },
	})

	got, err := adapter.BuildAuthorizeURL("state-123", "https://example.com/callback", "", nil)
	if err != nil {
		t.Fatalf("BuildAuthorizeURL() error = %v", err)
	}

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got, want := parsed.Path, "/api/v2/shop/auth_partner"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}

	query := parsed.Query()
	if got, want := query.Get("partner_id"), "10090"; got != want {
		t.Fatalf("partner_id = %q, want %q", got, want)
	}
	if got, want := query.Get("timestamp"), "1594897040"; got != want {
		t.Fatalf("timestamp = %q, want %q", got, want)
	}
	if got, want := query.Get("redirect"), "https://example.com/callback"; got != want {
		t.Fatalf("redirect = %q, want %q", got, want)
	}
	if got, want := query.Get("state"), "state-123"; got != want {
		t.Fatalf("state = %q, want %q", got, want)
	}
	if got, want := query.Get("sign"), "5c3357a5d7ede41e954ce12fe4a885c77dd0f4ff59666ddf16fbab373fdfe636"; got != want {
		t.Fatalf("sign = %q, want %q", got, want)
	}
}

func TestAdapterExchangeCallbackPostsShopeeTokenRequest(t *testing.T) {
	t.Parallel()

	var gotMethod, gotPath string
	var gotQuery url.Values
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"request_id":"req-1",
			"response":{
				"access_token":"access-123",
				"refresh_token":"refresh-456",
				"expire_in":14400,
				"shop_id":"9001"
			}
		}`))
	}))
	defer server.Close()

	adapter := NewAdapter(Config{
		PartnerID:  "10090",
		PartnerKey: "test-partner-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1594897040, 0).UTC() },
	})

	credential, err := adapter.ExchangeCallback(context.Background(), application.HandleCallbackAdapterInput{
		Code:              "code-123",
		ProviderAccountID: "9001",
	})
	if err != nil {
		t.Fatalf("ExchangeCallback() error = %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/api/v2/auth/token/get" {
		t.Fatalf("path = %q, want /api/v2/auth/token/get", gotPath)
	}
	if got, want := gotQuery.Get("partner_id"), "10090"; got != want {
		t.Fatalf("query partner_id = %q, want %q", got, want)
	}
	if got, want := gotQuery.Get("timestamp"), "1594897040"; got != want {
		t.Fatalf("query timestamp = %q, want %q", got, want)
	}
	if got, want := gotQuery.Get("sign"), "fdca6bda23a875e5a2b836c825009085487ff88af73c681b6019d5a5b70d3b34"; got != want {
		t.Fatalf("query sign = %q, want %q", got, want)
	}

	if got := strings.TrimSpace(stringField(gotBody, "code")); got != "code-123" {
		t.Fatalf("body code = %q, want code-123", got)
	}
	if got := strings.TrimSpace(stringField(gotBody, "shop_id")); got != "9001" {
		t.Fatalf("body shop_id = %q, want 9001", got)
	}

	if credential.SecretType != "shopee_partner" {
		t.Fatalf("secret type = %q, want shopee_partner", credential.SecretType)
	}
	if credential.AccessToken != "access-123" || credential.RefreshToken != "refresh-456" {
		t.Fatalf("credential tokens = %#v, want upstream values", credential)
	}
	if credential.ProviderAccountID != "9001" {
		t.Fatalf("provider account id = %q, want 9001", credential.ProviderAccountID)
	}
}

func TestAdapterRefreshPostsShopeeRefreshRequest(t *testing.T) {
	t.Parallel()

	var gotMethod, gotPath string
	var gotQuery url.Values
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"request_id":"req-2",
			"response":{
				"access_token":"access-789",
				"refresh_token":"refresh-999",
				"expire_in":14400,
				"shop_id":"9001"
			}
		}`))
	}))
	defer server.Close()

	adapter := NewAdapter(Config{
		PartnerID:  "10090",
		PartnerKey: "test-partner-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1594897040, 0).UTC() },
	})

	credential, err := adapter.Refresh(context.Background(), application.RefreshCredentialAdapterInput{
		RefreshToken:      "refresh-456",
		ProviderAccountID: "9001",
	})
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/api/v2/auth/access_token/get" {
		t.Fatalf("path = %q, want /api/v2/auth/access_token/get", gotPath)
	}
	if got, want := gotQuery.Get("sign"), "095229669dda19f50e82da0c4b2067215c24f11b582282705d507fc117a79838"; got != want {
		t.Fatalf("query sign = %q, want %q", got, want)
	}
	if got := strings.TrimSpace(stringField(gotBody, "refresh_token")); got != "refresh-456" {
		t.Fatalf("body refresh_token = %q, want refresh-456", got)
	}
	if got := strings.TrimSpace(stringField(gotBody, "shop_id")); got != "9001" {
		t.Fatalf("body shop_id = %q, want 9001", got)
	}

	if credential.SecretType != "shopee_partner" {
		t.Fatalf("secret type = %q, want shopee_partner", credential.SecretType)
	}
	if credential.AccessToken != "access-789" || credential.RefreshToken != "refresh-999" {
		t.Fatalf("credential tokens = %#v, want upstream values", credential)
	}
	if credential.ProviderAccountID != "9001" {
		t.Fatalf("provider account id = %q, want 9001", credential.ProviderAccountID)
	}
}

func TestAdapterMapsShopeeSignErrorsToConfigurationError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"request_id":"req-3","error":"error_sign","message":"Wrong sign."}`))
	}))
	defer server.Close()

	adapter := NewAdapter(Config{
		PartnerID:  "10090",
		PartnerKey: "test-partner-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1594897040, 0).UTC() },
	})

	_, err := adapter.ExchangeCallback(context.Background(), application.HandleCallbackAdapterInput{
		Code:              "code-123",
		ProviderAccountID: "9001",
	})
	if !errors.Is(err, domain.ErrAuthConfigurationInvalid) {
		t.Fatalf("error = %v, want auth configuration invalid", err)
	}
}

func TestAdapterRejectsMissingShopeeCallbackData(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{
		PartnerID:  "10090",
		PartnerKey: "test-partner-key",
		BaseURL:    "https://partner.test-stable.shopeemobile.com",
		Now:        func() time.Time { return time.Unix(1594897040, 0).UTC() },
	})

	if _, err := adapter.ExchangeCallback(context.Background(), application.HandleCallbackAdapterInput{Code: "code-123"}); !errors.Is(err, domain.ErrCredentialValidationFailed) {
		t.Fatalf("missing shop id error = %v, want credential validation failure", err)
	}
	if _, err := adapter.ExchangeCallback(context.Background(), application.HandleCallbackAdapterInput{ProviderAccountID: "9001"}); !errors.Is(err, domain.ErrAuthCodeExchangeFailed) {
		t.Fatalf("missing code error = %v, want auth code exchange failure", err)
	}
}

func stringField(body map[string]any, key string) string {
	value, _ := body[key].(string)
	return value
}
