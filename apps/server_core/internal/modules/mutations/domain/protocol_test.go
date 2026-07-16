package domain

import (
	"errors"
	"testing"
)

func TestProtocolStateTransitionsExhaustive(t *testing.T) {
	states := []ProtocolState{
		ProtocolStateDraft,
		ProtocolStatePreviewed,
		ProtocolStateApproved,
		ProtocolStateApplying,
		ProtocolStateApplied,
		ProtocolStatePartiallyFailed,
		ProtocolStateFailedPreserved,
		ProtocolStateCancelled,
	}
	legal := map[[2]ProtocolState]bool{
		{ProtocolStateDraft, ProtocolStatePreviewed}:          true,
		{ProtocolStateDraft, ProtocolStateCancelled}:          true,
		{ProtocolStatePreviewed, ProtocolStatePreviewed}:      true,
		{ProtocolStatePreviewed, ProtocolStateApproved}:       true,
		{ProtocolStatePreviewed, ProtocolStateCancelled}:      true,
		{ProtocolStateApproved, ProtocolStateApplying}:        true,
		{ProtocolStateApplying, ProtocolStateApplied}:         true,
		{ProtocolStateApplying, ProtocolStatePartiallyFailed}: true,
		{ProtocolStateApplying, ProtocolStateFailedPreserved}: true,
	}

	assertions := 0
	for _, from := range states {
		for _, to := range states {
			t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
				err := TransitionProtocolState(from, to)
				if legal[[2]ProtocolState{from, to}] {
					if err != nil {
						t.Fatalf("TransitionProtocolState(%q, %q) error = %v", from, to, err)
					}
					return
				}
				var stateErr *InvalidStateTransitionError
				if !errors.As(err, &stateErr) {
					t.Fatalf("TransitionProtocolState(%q, %q) error = %T, want *InvalidStateTransitionError", from, to, err)
				}
				if stateErr.Code != InvalidStateCode {
					t.Fatalf("error code = %q, want %q", stateErr.Code, InvalidStateCode)
				}
			})
			assertions++
		}
	}
	if assertions != 64 {
		t.Fatalf("assertions = %d, want 64", assertions)
	}
}

func TestProtocolAndItemStateValues(t *testing.T) {
	protocolStates := []struct {
		want string
		got  ProtocolState
	}{
		{"draft", ProtocolStateDraft},
		{"previewed", ProtocolStatePreviewed},
		{"approved", ProtocolStateApproved},
		{"applying", ProtocolStateApplying},
		{"applied", ProtocolStateApplied},
		{"partially_failed", ProtocolStatePartiallyFailed},
		{"failed_preserved", ProtocolStateFailedPreserved},
		{"cancelled", ProtocolStateCancelled},
	}
	for _, tt := range protocolStates {
		if string(tt.got) != tt.want {
			t.Errorf("protocol state = %q, want %q", tt.got, tt.want)
		}
	}

	itemStates := []struct {
		want string
		got  ItemState
	}{
		{"previewed", ItemStatePreviewed},
		{"approved", ItemStateApproved},
		{"applying", ItemStateApplying},
		{"applied", ItemStateApplied},
		{"failed", ItemStateFailed},
		{"skipped", ItemStateSkipped},
	}
	for _, tt := range itemStates {
		if string(tt.got) != tt.want {
			t.Errorf("item state = %q, want %q", tt.got, tt.want)
		}
	}
}

func TestTerminalProtocolState(t *testing.T) {
	tests := []struct {
		name    string
		items   []ItemState
		want    ProtocolState
		wantErr bool
	}{
		{"all ok", []ItemState{ItemStateApplied, ItemStateApplied}, ProtocolStateApplied, false},
		{"with skipped", []ItemState{ItemStateApplied, ItemStateSkipped}, ProtocolStateApplied, false},
		{"all skipped", []ItemState{ItemStateSkipped, ItemStateSkipped}, ProtocolStateApplied, false},
		{"mixed", []ItemState{ItemStateApplied, ItemStateSkipped, ItemStateFailed}, ProtocolStatePartiallyFailed, false},
		{"all fail", []ItemState{ItemStateFailed, ItemStateFailed}, ProtocolStateFailedPreserved, false},
		{"empty", nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TerminalProtocolState(tt.items)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TerminalProtocolState() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("TerminalProtocolState() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFailureCodeRetryability(t *testing.T) {
	tests := []struct {
		value     string
		code      FailureCode
		retryable bool
	}{
		{"provider_validation", FailureCodeProviderValidation, false},
		{"provider_rate_limited", FailureCodeProviderRateLimited, true},
		{"provider_unavailable", FailureCodeProviderUnavailable, true},
		{"provider_auth", FailureCodeProviderAuth, false},
		{"listing_paused_remote", FailureCodeListingPausedRemote, false},
		{"link_unresolved", FailureCodeLinkUnresolved, false},
		{"policy_missing", FailureCodePolicyMissing, false},
		{"sku_invariant_violation", FailureCodeSKUInvariantViolation, false},
		{"stale_source", FailureCodeStaleSource, true},
		{"conflict_remote_changed", FailureCodeConflictRemoteChanged, false},
		{"type_not_enabled", FailureCodeTypeNotEnabled, false},
		{"internal", FailureCodeInternal, false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if string(tt.code) != tt.value {
				t.Fatalf("failure code = %q, want %q", tt.code, tt.value)
			}
			if got := tt.code.RetryableByDefault(); got != tt.retryable {
				t.Fatalf("RetryableByDefault() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestProtocolTypeValues(t *testing.T) {
	tests := []struct {
		want string
		got  ProtocolType
	}{
		{"price_update", ProtocolTypePriceUpdate},
		{"stock_correct", ProtocolTypeStockCorrect},
		{"link_apply", ProtocolTypeLinkApply},
		{"listing_pause", ProtocolTypeListingPause},
		{"listing_resync", ProtocolTypeListingResync},
		{"listing_edit", ProtocolTypeListingEdit},
		{"listing_create", ProtocolTypeListingCreate},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Fatalf("protocol type = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
