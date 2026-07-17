package domain

type FailureCode string

const (
	FailureCodeProviderValidation    FailureCode = "provider_validation"
	FailureCodeProviderRateLimited   FailureCode = "provider_rate_limited"
	FailureCodeProviderUnavailable   FailureCode = "provider_unavailable"
	FailureCodeProviderAuth          FailureCode = "provider_auth"
	FailureCodeListingPausedRemote   FailureCode = "listing_paused_remote"
	FailureCodeLinkUnresolved        FailureCode = "link_unresolved"
	FailureCodePolicyMissing         FailureCode = "policy_missing"
	FailureCodeSKUInvariantViolation FailureCode = "sku_invariant_violation"
	FailureCodeStaleSource           FailureCode = "stale_source"
	FailureCodeConflictRemoteChanged FailureCode = "conflict_remote_changed"
	FailureCodeTypeNotEnabled        FailureCode = "type_not_enabled"
	FailureCodeInternal              FailureCode = "internal"
)

type Failure struct {
	Code            FailureCode
	MessagePT       string
	MessageProvider string
	Retryable       bool
}

func (c FailureCode) RetryableByDefault() bool {
	switch c {
	case FailureCodeProviderRateLimited, FailureCodeProviderUnavailable, FailureCodeStaleSource:
		return true
	default:
		return false
	}
}
