package application

import (
	"context"
	"errors"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type CredentialResolutionRef struct {
	TenantID          string
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
}

type ResolvedCredential struct {
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
	AccessToken       string
}

type credentialLookup interface {
	GetActiveCredential(ctx context.Context, installationID string) (domain.Credential, bool, error)
}

type payloadDecryptor interface {
	DecryptJSON(encoded []byte) (map[string]any, string, error)
}

type CredentialResolver struct {
	credentials credentialLookup
	decryptor   payloadDecryptor
}

func NewCredentialResolver(credentials credentialLookup, decryptor payloadDecryptor) *CredentialResolver {
	return &CredentialResolver{
		credentials: credentials,
		decryptor:   decryptor,
	}
}

func (r *CredentialResolver) ResolveAccessToken(ctx context.Context, ref CredentialResolutionRef) (ResolvedCredential, error) {
	ref.TenantID = strings.TrimSpace(ref.TenantID)
	ref.InstallationID = strings.TrimSpace(ref.InstallationID)
	ref.ProviderCode = strings.TrimSpace(ref.ProviderCode)
	ref.ProviderAccountID = strings.TrimSpace(ref.ProviderAccountID)
	if ref.TenantID == "" || ref.InstallationID == "" || ref.ProviderCode == "" || ref.ProviderAccountID == "" {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_INVALID_REFERENCE")
	}
	if r.credentials == nil || r.decryptor == nil {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_RESOLVER_NOT_CONFIGURED")
	}

	credential, found, err := r.credentials.GetActiveCredential(ctx, ref.InstallationID)
	if err != nil {
		return ResolvedCredential{}, err
	}
	if !found {
		return ResolvedCredential{}, domain.ErrCredentialNotFound
	}

	payload, _, err := r.decryptor.DecryptJSON(credential.EncryptedPayload)
	if err != nil {
		return ResolvedCredential{}, domain.ErrCredentialDecryptionFailed
	}
	accessToken, ok := credentialPayloadString(payload, "access_token")
	if !ok {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_ACCESS_TOKEN_MISSING")
	}

	return ResolvedCredential{
		InstallationID:    ref.InstallationID,
		ProviderCode:      ref.ProviderCode,
		ProviderAccountID: ref.ProviderAccountID,
		AccessToken:       accessToken,
	}, nil
}
