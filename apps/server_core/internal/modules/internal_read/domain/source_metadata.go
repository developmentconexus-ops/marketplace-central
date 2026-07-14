package domain

import (
	"context"
	"time"
)

type SourceMetadata struct {
	System     string
	FetchedAt  time.Time
	ObservedAt *time.Time
}

type FreshnessPolicy struct {
	MaxAge time.Duration
}

func (p FreshnessPolicy) IsZero() bool {
	return p.MaxAge <= 0
}

type freshnessPolicyContextKey struct{}

// WithFreshnessPolicy carries transport freshness intent to read adapters
// without coupling the adapter layer to a transport package.
func WithFreshnessPolicy(ctx context.Context, policy FreshnessPolicy) context.Context {
	return context.WithValue(ctx, freshnessPolicyContextKey{}, policy)
}

func FreshnessPolicyFromContext(ctx context.Context) (FreshnessPolicy, bool) {
	policy, ok := ctx.Value(freshnessPolicyContextKey{}).(FreshnessPolicy)
	return policy, ok
}
