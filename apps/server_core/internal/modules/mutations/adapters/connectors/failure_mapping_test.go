package connectors

import (
	"context"
	"errors"
	"testing"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

func TestFailureMappingCoversProviderSurface(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		rejected, paused bool
		want             domain.FailureCode
		retry            bool
	}{
		{"429", connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderRateLimited, "safe 429"), false, false, domain.FailureCodeProviderRateLimited, true},
		{"401", connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderAuth, "safe 401"), false, false, domain.FailureCodeProviderAuth, false},
		{"403", connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderAuth, "safe 403"), false, false, domain.FailureCodeProviderAuth, false},
		{"validation", connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderValidation, "safe validation"), false, false, domain.FailureCodeProviderValidation, false},
		{"paused", nil, true, true, domain.FailureCodeListingPausedRemote, false},
		{"timeout", context.DeadlineExceeded, false, false, domain.FailureCodeProviderUnavailable, true},
		{"5xx", connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, "safe 503"), false, false, domain.FailureCodeProviderUnavailable, true},
		{"unknown", errors.New("unknown"), false, false, domain.FailureCodeInternal, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got domain.Failure
			if tt.rejected {
				got = MapRejected("safe paused", tt.paused)
			} else {
				got = MapFailure(tt.err)
			}
			if got.Code != tt.want || got.Retryable != tt.retry {
				t.Fatalf("got=%+v", got)
			}
			if tt.err != nil && connectorsdomain.ErrorCodeOf(tt.err) != "" && got.MessageProvider == "" {
				t.Fatal("sanitized provider message not retained")
			}
		})
	}
}
