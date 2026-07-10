package application

import (
	"context"
	"strings"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type installationAuthSessionLookup interface {
	GetAuthSession(ctx context.Context, installationID string) (domain.AuthSession, bool, error)
}

type InstallationStateReconcilerConfig struct {
	TenantID      string
	AuthSessions  installationAuthSessionLookup
	Credentials   credentialLookup
	Decryptor     payloadDecryptor
	Capabilities  *connectorsapp.MarketplaceCapabilityService
	Installations interface {
		ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error
	}
}

type InstallationStateReconciler struct {
	tenantID      string
	authSessions  installationAuthSessionLookup
	credentials   credentialLookup
	decryptor     payloadDecryptor
	capabilities  *connectorsapp.MarketplaceCapabilityService
	installations interface {
		ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error
	}
}

func NewInstallationStateReconciler(cfg InstallationStateReconcilerConfig) *InstallationStateReconciler {
	tenantID := strings.TrimSpace(cfg.TenantID)
	if tenantID == "" {
		tenantID = "tenant_default"
	}
	return &InstallationStateReconciler{
		tenantID:      tenantID,
		authSessions:  cfg.AuthSessions,
		credentials:   cfg.Credentials,
		decryptor:     cfg.Decryptor,
		capabilities:  cfg.Capabilities,
		installations: cfg.Installations,
	}
}

func (r *InstallationStateReconciler) Reconcile(ctx context.Context, inst domain.Installation) (domain.Installation, error) {
	if r == nil || r.authSessions == nil || r.credentials == nil || r.decryptor == nil || r.installations == nil {
		return inst, nil
	}
	if !needsInstallationStateRepair(inst) {
		return inst, nil
	}

	session, found, err := r.authSessions.GetAuthSession(ctx, inst.InstallationID)
	if err != nil || !found {
		return inst, err
	}
	if strings.TrimSpace(session.ProviderAccountID) == "" {
		return inst, nil
	}

	credential, found, err := r.credentials.GetActiveCredential(ctx, inst.InstallationID)
	if err != nil || !found {
		return inst, err
	}

	payload, _, err := r.decryptor.DecryptJSON(credential.EncryptedPayload)
	if err != nil {
		return inst, nil
	}

	repaired := inst
	repaired.ExternalAccountID = strings.TrimSpace(session.ProviderAccountID)
	repaired.ActiveCredentialID = firstNonEmpty(inst.ActiveCredentialID, credential.CredentialID)
	repaired.LastVerifiedAt = firstNonNilTime(session.LastVerifiedAt, inst.LastVerifiedAt)
	repaired.ExternalAccountName = firstNonEmpty(inst.ExternalAccountName, credentialPayloadAccountName(payload))

	if prober, probeErr := r.accountProber(inst.ProviderCode); probeErr == nil {
		snapshot, err := prober.ProbeAccount(ctx, connectorsdomain.ProviderAccountRef{
			TenantID:          r.tenantID,
			InstallationID:    inst.InstallationID,
			ProviderCode:      inst.ProviderCode,
			ProviderAccountID: repaired.ExternalAccountID,
		})
		if err == nil {
			repaired.ExternalAccountID = firstNonEmpty(snapshot.ProviderAccountID, repaired.ExternalAccountID)
			repaired.ExternalAccountName = firstNonEmpty(snapshot.ProviderAccountName, repaired.ExternalAccountName)
			repaired.LastVerifiedAt = firstNonNilTime(&snapshot.FetchedAt, repaired.LastVerifiedAt)
		}
	}

	repaired.ConnectionSnapshot = domain.ProjectConnectionSnapshot(
		repaired,
		inferAuthStrategyFromPayload(payload),
		session.AccessTokenExpiresAt,
		inst.ConnectionSnapshot.ReauthReason,
	)
	if err := r.installations.ApplyConnectionSnapshot(ctx, inst.InstallationID, repaired.ConnectionSnapshot, repaired.ActiveCredentialID); err != nil {
		return inst, err
	}
	return repaired, nil
}

func needsInstallationStateRepair(inst domain.Installation) bool {
	if inst.Status != domain.InstallationStatusConnected {
		return false
	}
	if strings.TrimSpace(inst.ExternalAccountID) == "" {
		return true
	}
	return inst.ConnectionSnapshot.AuthStrategy == "" || inst.ConnectionSnapshot.AuthStrategy == domain.AuthStrategyUnknown
}

func (r *InstallationStateReconciler) accountProber(providerCode string) (connectorsports.AccountProber, error) {
	if r.capabilities == nil {
		return nil, connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderUnsupportedShape, "account probe unavailable")
	}
	return r.capabilities.AccountProber(providerCode)
}

func inferAuthStrategyFromPayload(payload map[string]any) domain.AuthStrategy {
	switch {
	case hasCredentialValue(payload, "api_key"):
		return domain.AuthStrategyAPIKey
	case hasCredentialValue(payload, "refresh_token"), hasCredentialValue(payload, "access_token"):
		return domain.AuthStrategyOAuth2
	default:
		return domain.AuthStrategyUnknown
	}
}

func hasCredentialValue(payload map[string]any, key string) bool {
	_, ok := credentialPayloadString(payload, key)
	return ok
}

func credentialPayloadAccountName(payload map[string]any) string {
	if value, ok := credentialPayloadString(payload, "provider_account_name"); ok {
		return value
	}
	return ""
}

func firstNonNilTime(values ...*time.Time) *time.Time {
	for _, value := range values {
		if value != nil {
			copy := value.UTC()
			return &copy
		}
	}
	return nil
}
