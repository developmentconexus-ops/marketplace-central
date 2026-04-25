package amazon

import (
	"context"
	"net/url"
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
		ClientID:     "amazon-client-id",
		AuthorizeURL: "https://sellercentral.amazon.com.br/apps/authorize/consent",
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
	if got, want := query.Get("client_id"), "amazon-client-id"; got != want {
		t.Fatalf("client_id = %q, want %q", got, want)
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
