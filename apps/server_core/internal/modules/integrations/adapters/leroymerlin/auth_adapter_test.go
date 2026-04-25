package leroymerlin

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestAdapterAcceptsAPIKeyCredential(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{})

	credential, err := adapter.VerifyAPIKey(context.Background(), application.SubmitAPIKeyAdapterInput{
		InstallationID: "inst-leroy",
		APIKey:         "leroy-api-key",
		Metadata: map[string]string{
			"shop_id":   "shop-leroy-1",
			"shop_name": "Leroy Loja",
		},
	})
	if err != nil {
		t.Fatalf("VerifyAPIKey() error = %v", err)
	}
	if got, want := credential.SecretType, "api_key"; got != want {
		t.Fatalf("SecretType = %q, want %q", got, want)
	}
	if got, want := credential.ProviderAccountID, "shop-leroy-1"; got != want {
		t.Fatalf("ProviderAccountID = %q, want %q", got, want)
	}
	if got, want := credential.ProviderAccountName, "Leroy Loja"; got != want {
		t.Fatalf("ProviderAccountName = %q, want %q", got, want)
	}
}

func TestAdapterRejectsBlankAPIKey(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(Config{})

	_, err := adapter.VerifyAPIKey(context.Background(), application.SubmitAPIKeyAdapterInput{
		InstallationID: "inst-leroy",
		APIKey:         " ",
	})
	if err == nil {
		t.Fatal("VerifyAPIKey() error = nil, want validation error")
	}
	if err != domain.ErrAPIKeyValidationFailed {
		t.Fatalf("VerifyAPIKey() error = %v, want %v", err, domain.ErrAPIKeyValidationFailed)
	}
}
