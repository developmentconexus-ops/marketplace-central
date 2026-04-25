package madeiramadeira

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestAdapterAcceptsTokenFromMetadata(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{})

	credential, err := adapter.VerifyAPIKey(context.Background(), application.SubmitAPIKeyAdapterInput{
		InstallationID: "inst-madeira",
		Metadata: map[string]string{
			"api_token":   "madeira-token",
			"seller_id":   "seller-1",
			"seller_name": "Madeira Loja",
		},
	})
	if err != nil {
		t.Fatalf("VerifyAPIKey() error = %v", err)
	}
	if got, want := credential.SecretType, "token"; got != want {
		t.Fatalf("SecretType = %q, want %q", got, want)
	}
	if got, want := credential.APIKey, "madeira-token"; got != want {
		t.Fatalf("APIKey = %q, want %q", got, want)
	}
	if got, want := credential.ProviderAccountID, "seller-1"; got != want {
		t.Fatalf("ProviderAccountID = %q, want %q", got, want)
	}
	if got, want := credential.ProviderAccountName, "Madeira Loja"; got != want {
		t.Fatalf("ProviderAccountName = %q, want %q", got, want)
	}
}

func TestAdapterRejectsMissingToken(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{})

	_, err := adapter.VerifyAPIKey(context.Background(), application.SubmitAPIKeyAdapterInput{
		InstallationID: "inst-madeira",
		APIKey:         " ",
		Metadata:       map[string]string{},
	})
	if err == nil {
		t.Fatal("VerifyAPIKey() error = nil, want validation error")
	}
	if err != domain.ErrAPIKeyValidationFailed {
		t.Fatalf("VerifyAPIKey() error = %v, want %v", err, domain.ErrAPIKeyValidationFailed)
	}
}
