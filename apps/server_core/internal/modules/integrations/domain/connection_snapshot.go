package domain

import (
	"strings"
	"time"
)

type ConnectionState string
type ConnectionNextAction string

const (
	ConnectionStateDraft        ConnectionState = "draft"
	ConnectionStatePending      ConnectionState = "pending_connection"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateDegraded     ConnectionState = "degraded"
	ConnectionStateNeedsReauth  ConnectionState = "needs_reauth"
	ConnectionStateDisconnected ConnectionState = "disconnected"

	ConnectionNextActionNone      ConnectionNextAction = "none"
	ConnectionNextActionAuthorize ConnectionNextAction = "authorize"
	ConnectionNextActionReauth    ConnectionNextAction = "reauth"
	ConnectionNextActionConfigure ConnectionNextAction = "configure"
	ConnectionNextActionRetry     ConnectionNextAction = "retry"
)

type ConnectionSnapshot struct {
	State               ConnectionState      `json:"state"`
	Health              HealthStatus         `json:"health"`
	ProviderCode        string               `json:"provider_code"`
	ExternalAccountID   string               `json:"external_account_id"`
	ExternalAccountName string               `json:"external_account_name"`
	AuthStrategy        AuthStrategy         `json:"auth_strategy"`
	LastVerifiedAt      *time.Time           `json:"last_verified_at,omitempty"`
	ExpiresAt           *time.Time           `json:"expires_at,omitempty"`
	NextAction          ConnectionNextAction `json:"next_action"`
	ReauthReason        string               `json:"reauth_reason,omitempty"`
}

type ConnectedSnapshotInput struct {
	ProviderCode        string
	ExternalAccountID   string
	ExternalAccountName string
	AuthStrategy        AuthStrategy
	VerifiedAt          time.Time
	ExpiresAt           *time.Time
}

func NewConnectedSnapshot(input ConnectedSnapshotInput) ConnectionSnapshot {
	verifiedAt := input.VerifiedAt.UTC()

	var expiresAt *time.Time
	if input.ExpiresAt != nil {
		value := input.ExpiresAt.UTC()
		expiresAt = &value
	}

	authStrategy := input.AuthStrategy
	if authStrategy == "" {
		authStrategy = AuthStrategyNone
	}

	return ConnectionSnapshot{
		State:               ConnectionStateConnected,
		Health:              HealthStatusHealthy,
		ProviderCode:        strings.TrimSpace(input.ProviderCode),
		ExternalAccountID:   strings.TrimSpace(input.ExternalAccountID),
		ExternalAccountName: strings.TrimSpace(input.ExternalAccountName),
		AuthStrategy:        authStrategy,
		LastVerifiedAt:      &verifiedAt,
		ExpiresAt:           expiresAt,
		NextAction:          ConnectionNextActionNone,
	}
}

func NewDisconnectedSnapshot(providerCode string, reason string) ConnectionSnapshot {
	return ConnectionSnapshot{
		State:        ConnectionStateDisconnected,
		Health:       HealthStatusWarning,
		ProviderCode: strings.TrimSpace(providerCode),
		AuthStrategy: AuthStrategyNone,
		NextAction:   ConnectionNextActionAuthorize,
		ReauthReason: strings.TrimSpace(reason),
	}
}

func ProjectConnectionSnapshot(inst Installation, authStrategy AuthStrategy, expiresAt *time.Time, reauthReason string) ConnectionSnapshot {
	providerCode := strings.TrimSpace(inst.ProviderCode)
	if authStrategy == "" {
		authStrategy = AuthStrategyUnknown
	}

	switch inst.Status {
	case InstallationStatusConnected:
		return NewConnectedSnapshot(ConnectedSnapshotInput{
			ProviderCode:        providerCode,
			ExternalAccountID:   inst.ExternalAccountID,
			ExternalAccountName: inst.ExternalAccountName,
			AuthStrategy:        authStrategy,
			VerifiedAt:          timeOrNow(inst.LastVerifiedAt),
			ExpiresAt:           expiresAt,
		})
	case InstallationStatusRequiresReauth:
		return ConnectionSnapshot{
			State:               ConnectionStateNeedsReauth,
			Health:              normalizeConnectionHealth(inst.HealthStatus),
			ProviderCode:        providerCode,
			ExternalAccountID:   strings.TrimSpace(inst.ExternalAccountID),
			ExternalAccountName: strings.TrimSpace(inst.ExternalAccountName),
			AuthStrategy:        authStrategy,
			LastVerifiedAt:      cloneTimePtr(inst.LastVerifiedAt),
			ExpiresAt:           cloneTimePtr(expiresAt),
			NextAction:          ConnectionNextActionReauth,
			ReauthReason:        strings.TrimSpace(reauthReason),
		}
	case InstallationStatusDegraded:
		return ConnectionSnapshot{
			State:               ConnectionStateDegraded,
			Health:              normalizeConnectionHealth(inst.HealthStatus),
			ProviderCode:        providerCode,
			ExternalAccountID:   strings.TrimSpace(inst.ExternalAccountID),
			ExternalAccountName: strings.TrimSpace(inst.ExternalAccountName),
			AuthStrategy:        authStrategy,
			LastVerifiedAt:      cloneTimePtr(inst.LastVerifiedAt),
			ExpiresAt:           cloneTimePtr(expiresAt),
			NextAction:          ConnectionNextActionRetry,
			ReauthReason:        strings.TrimSpace(reauthReason),
		}
	case InstallationStatusPendingConnection:
		return ConnectionSnapshot{
			State:               ConnectionStatePending,
			Health:              normalizeConnectionHealth(inst.HealthStatus),
			ProviderCode:        providerCode,
			ExternalAccountID:   strings.TrimSpace(inst.ExternalAccountID),
			ExternalAccountName: strings.TrimSpace(inst.ExternalAccountName),
			AuthStrategy:        authStrategy,
			LastVerifiedAt:      cloneTimePtr(inst.LastVerifiedAt),
			ExpiresAt:           cloneTimePtr(expiresAt),
			NextAction:          ConnectionNextActionNone,
		}
	case InstallationStatusDraft:
		return ConnectionSnapshot{
			State:               ConnectionStateDraft,
			Health:              normalizeConnectionHealth(inst.HealthStatus),
			ProviderCode:        providerCode,
			ExternalAccountID:   strings.TrimSpace(inst.ExternalAccountID),
			ExternalAccountName: strings.TrimSpace(inst.ExternalAccountName),
			AuthStrategy:        authStrategy,
			LastVerifiedAt:      cloneTimePtr(inst.LastVerifiedAt),
			ExpiresAt:           cloneTimePtr(expiresAt),
			NextAction:          ConnectionNextActionConfigure,
		}
	default:
		return ConnectionSnapshot{
			State:               ConnectionStateDisconnected,
			Health:              normalizeConnectionHealth(inst.HealthStatus),
			ProviderCode:        providerCode,
			ExternalAccountID:   strings.TrimSpace(inst.ExternalAccountID),
			ExternalAccountName: strings.TrimSpace(inst.ExternalAccountName),
			AuthStrategy:        authStrategy,
			LastVerifiedAt:      cloneTimePtr(inst.LastVerifiedAt),
			ExpiresAt:           cloneTimePtr(expiresAt),
			NextAction:          ConnectionNextActionAuthorize,
			ReauthReason:        strings.TrimSpace(reauthReason),
		}
	}
}

func normalizeConnectionHealth(health HealthStatus) HealthStatus {
	switch health {
	case HealthStatusHealthy, HealthStatusWarning, HealthStatusCritical:
		return health
	default:
		return HealthStatusWarning
	}
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := value.UTC()
	return &copy
}

func timeOrNow(value *time.Time) time.Time {
	if value == nil {
		return time.Now().UTC()
	}
	return value.UTC()
}
