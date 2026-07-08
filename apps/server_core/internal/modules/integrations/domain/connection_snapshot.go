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
