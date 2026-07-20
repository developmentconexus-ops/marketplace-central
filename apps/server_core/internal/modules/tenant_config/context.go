package tenant_config

import "context"

// ctxKey is the unexported key type for the active-source Config carried on a
// request context.
type ctxKey struct{}

// WithActiveSource pins the tenant's resolved active-source Config on ctx.
func WithActiveSource(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, ctxKey{}, cfg)
}

// FromContext reports the Config pinned on ctx and whether one was set.
// Exported so downstream consumers (e.g. M-04's cache-key derivation, C13)
// can read the active source without re-resolving it.
func FromContext(ctx context.Context) (Config, bool) {
	cfg, ok := ctx.Value(ctxKey{}).(Config)
	return cfg, ok
}
