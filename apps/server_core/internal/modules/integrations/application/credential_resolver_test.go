package application

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type fakeCredentialLookup struct {
	credential domain.Credential
	found      bool
	err        error
}

func (f *fakeCredentialLookup) GetActiveCredential(context.Context, string) (domain.Credential, bool, error) {
	return f.credential, f.found, f.err
}

type fakePayloadDecryptor struct {
	payload map[string]any
	err     error
}

func (f fakePayloadDecryptor) DecryptJSON([]byte) (map[string]any, string, error) {
	return f.payload, "key-1", f.err
}

func TestCredentialResolverRejectsIncompleteReference(t *testing.T) {
	t.Parallel()

	resolver := NewCredentialResolver(nil, nil)

	_, err := resolver.ResolveAccessToken(context.Background(), CredentialResolutionRef{
		TenantID:       "tenant_default",
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
	})

	if err == nil || err.Error() != "INTEGRATIONS_CREDENTIAL_INVALID_REFERENCE" {
		t.Fatalf("err = %v", err)
	}
}

func TestCredentialResolverDoesNotReturnRefreshToken(t *testing.T) {
	t.Parallel()

	store := &fakeCredentialLookup{
		credential: domain.Credential{
			CredentialID:     "cred-1",
			InstallationID:   "inst-1",
			EncryptedPayload: []byte("ciphertext"),
		},
		found: true,
	}
	decryptor := fakePayloadDecryptor{
		payload: map[string]any{
			"access_token":  "access-token",
			"refresh_token": "refresh-token",
		},
	}
	resolver := NewCredentialResolver(store, decryptor)

	resolved, err := resolver.ResolveAccessToken(context.Background(), CredentialResolutionRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	})

	if err != nil {
		t.Fatalf("ResolveAccessToken() error = %v", err)
	}
	if resolved.AccessToken != "access-token" {
		t.Fatalf("access token = %q", resolved.AccessToken)
	}
	if resolved.ProviderAccountID != "691607102" {
		t.Fatalf("provider account id = %q", resolved.ProviderAccountID)
	}
}
