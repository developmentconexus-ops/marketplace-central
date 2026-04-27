package amazon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestAdapterStrategyIsLWA(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{})
	if got, want := adapter.ProviderCode(), "amazon"; got != want {
		t.Fatalf("ProviderCode() = %q, want %q", got, want)
	}
	if got, want := adapter.AuthStrategy(), domain.AuthStrategyLWA; got != want {
		t.Fatalf("AuthStrategy() = %q, want %q", got, want)
	}
}

func TestAdapterBuildsAuthorizeURL(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{
		ApplicationID: "amazon-application-id",
		ClientID:      "amazon-client-id",
		AuthorizeURL:  "https://sellercentral.amazon.com.br/apps/authorize/consent",
	})

	start, err := adapter.StartAuthorize(context.Background(), application.StartAuthorizeAdapterInput{
		InstallationID: "inst-amazon",
		State:          "state-amazon",
		RedirectURI:    "https://app.test/integrations/callback",
		CodeChallenge:  "challenge-amazon",
		Scopes:         []string{"sellingpartnerapi::notifications", "sellingpartnerapi::orders"},
	})
	if err != nil {
		t.Fatalf("StartAuthorize() error = %v", err)
	}

	parsed, err := url.Parse(start.AuthURL)
	if err != nil {
		t.Fatalf("parse auth URL: %v", err)
	}
	query := parsed.Query()
	if got, want := query.Get("application_id"), "amazon-application-id"; got != want {
		t.Fatalf("application_id = %q, want %q", got, want)
	}
	if got := query.Get("client_id"); got != "" {
		t.Fatalf("client_id = %q, want empty for Amazon consent URL", got)
	}
	if got, want := query.Get("state"), "state-amazon"; got != want {
		t.Fatalf("state = %q, want %q", got, want)
	}
	if got, want := query.Get("code_challenge"), "challenge-amazon"; got != want {
		t.Fatalf("code_challenge = %q, want %q", got, want)
	}
	if got, want := query.Get("code_challenge_method"), "S256"; got != want {
		t.Fatalf("code_challenge_method = %q, want %q", got, want)
	}
	if got, want := query.Get("scope"), "sellingpartnerapi::notifications sellingpartnerapi::orders"; got != want {
		t.Fatalf("scope = %q, want %q", got, want)
	}
}

func TestAdapterExchangeCallbackUsesLWATokenEndpoint(t *testing.T) {
	t.Parallel()

	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; got != want {
			t.Fatalf("Content-Type = %q, want %q", got, want)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		receivedBody = r.Form.Encode()
		assertFormValue(t, r.Form, "grant_type", "authorization_code")
		assertFormValue(t, r.Form, "code", "spapi-code")
		assertFormValue(t, r.Form, "redirect_uri", "https://app.test/integrations/auth/callback")
		assertFormValue(t, r.Form, "client_id", "lwa-client-id")
		assertFormValue(t, r.Form, "client_secret", "lwa-client-secret")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "access-token",
			"refresh_token": "refresh-token",
			"expires_in":    3600,
		})
	}))
	defer server.Close()

	adapter := NewAdapter(Config{
		ClientID:     "lwa-client-id",
		ClientSecret: "lwa-client-secret",
		TokenURL:     server.URL,
	})

	payload, err := adapter.ExchangeCallback(context.Background(), application.HandleCallbackAdapterInput{
		Code:              "spapi-code",
		RedirectURI:       "https://app.test/integrations/auth/callback",
		ProviderAccountID: "A1SELLER",
	})
	if err != nil {
		t.Fatalf("ExchangeCallback() error = %v; body=%s", err, receivedBody)
	}
	if got, want := payload.SecretType, "lwa"; got != want {
		t.Fatalf("SecretType = %q, want %q", got, want)
	}
	if got, want := payload.AccessToken, "access-token"; got != want {
		t.Fatalf("AccessToken = %q, want %q", got, want)
	}
	if got, want := payload.RefreshToken, "refresh-token"; got != want {
		t.Fatalf("RefreshToken = %q, want %q", got, want)
	}
	if got, want := payload.ProviderAccountID, "A1SELLER"; got != want {
		t.Fatalf("ProviderAccountID = %q, want %q", got, want)
	}
	if payload.ExpiresAt == nil {
		t.Fatal("ExpiresAt is nil, want token expiry")
	}
}

func assertFormValue(t *testing.T, form url.Values, key, want string) {
	t.Helper()
	if got := strings.TrimSpace(form.Get(key)); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
