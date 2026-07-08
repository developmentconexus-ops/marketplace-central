package domain

import (
	"testing"
	"time"
)

func TestNewConnectedSnapshotNormalizesAccountAndNextAction(t *testing.T) {
	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	expires := now.Add(6 * time.Hour)

	snapshot := NewConnectedSnapshot(ConnectedSnapshotInput{
		ProviderCode:        " mercado_livre ",
		ExternalAccountID:   " 691607102 ",
		ExternalAccountName: " METALNOBREACABAMENTOS ",
		AuthStrategy:        AuthStrategyOAuth2,
		VerifiedAt:          now,
		ExpiresAt:           &expires,
	})

	if snapshot.State != ConnectionStateConnected {
		t.Fatalf("state = %q, want %q", snapshot.State, ConnectionStateConnected)
	}
	if snapshot.Health != HealthStatusHealthy {
		t.Fatalf("health = %q, want %q", snapshot.Health, HealthStatusHealthy)
	}
	if snapshot.ProviderCode != "mercado_livre" {
		t.Fatalf("provider = %q", snapshot.ProviderCode)
	}
	if snapshot.ExternalAccountID != "691607102" {
		t.Fatalf("external account id = %q", snapshot.ExternalAccountID)
	}
	if snapshot.ExternalAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("external account name = %q", snapshot.ExternalAccountName)
	}
	if snapshot.NextAction != ConnectionNextActionNone {
		t.Fatalf("next action = %q, want none", snapshot.NextAction)
	}
	if snapshot.ExpiresAt == nil || !snapshot.ExpiresAt.Equal(expires) {
		t.Fatalf("expires at was not preserved")
	}
}

func TestDisconnectedSnapshotRequiresAuthAction(t *testing.T) {
	snapshot := NewDisconnectedSnapshot("mercado_livre", "reauth required")

	if snapshot.State != ConnectionStateDisconnected {
		t.Fatalf("state = %q", snapshot.State)
	}
	if snapshot.Health != HealthStatusWarning {
		t.Fatalf("health = %q", snapshot.Health)
	}
	if snapshot.NextAction != ConnectionNextActionAuthorize {
		t.Fatalf("next action = %q", snapshot.NextAction)
	}
	if snapshot.ReauthReason != "reauth required" {
		t.Fatalf("reauth reason = %q", snapshot.ReauthReason)
	}
}
