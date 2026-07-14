# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Feature ID

F-01-freshness-cache

## Problem

Repeated installation-global Oracle reads currently execute downstream for every
request. Add a bounded per-process L2 cache that preserves source freshness
semantics, collapses concurrent misses, and is invalidated by successful writes.

## Requirements

- Decorate catalog pages, cost/tax batch reads, and inventory stock batch reads
  without changing port signatures or caching linkage readers.
- Use deterministic method-and-arguments keys, class-specific FreshnessPolicy
  TTLs, injected time, singleflight, no-cache bypass, and LRU eviction.
- Preserve complete results, including nil map values and catalog `AsOf`.
- Emit structured hit/miss/bypass logs without raw arguments or key contents.
- Inject post-commit invalidation into the four specified write services.

## Non-Goals

- Redis, an invalidation bus, cross-process coordination, or tenant keying.
- Caching interactive linkage or any other internal-read method.
- OpenAPI or sdk-runtime changes; public method signatures remain unchanged.
- Adding timing wrappers to batch readers.

## Design

`internal_read/adapters/cache` owns a shared mutex-protected map plus LRU list,
with one `singleflight.Group`. Entries carry class, value, entry time, and
snapshot time. Per-port decorators canonicalize arguments before hashing and
use the neutral `FreshnessPolicyFromContext` seam. Catalog pages retain the
downstream page `AsOf`; batch entries capture the injected clock at insertion.

The root creates one cache from environment-backed class policies and max-entry
cap. The interactive path is `TimingReader(CacheReader(OracleReader))`; batch
cost/tax and stock paths receive their own decorators over the shared cache.
`internal_read/ports.CacheInvalidator` is passed into catalog, inventory,
orders-linkage, and profitability application services, each invoked only
after its successful persistence point.

## Edge Cases

- Empty or invalid TTL/cap environment values use safe defaults.
- A zero policy present in context bypasses lookup but repopulates after a
  successful downstream call.
- Downstream errors are returned to all singleflight waiters and never stored.
- Nil maps and nil map values remain distinguishable from empty or absent data.
- Invalidation removes only the requested fact class.
- Linkage remains structurally outside the decorator graph.

## Acceptance Criteria

- Criterion: TTL hit/expiry semantics per data class, including configured
  defaults/env values and catalog `AsOf`.
  - Traces to milestone criterion ID: M-04-C01
  - Proven by: `go test ./... -run 'FreshnessCache' -v`
- Criterion: 20 concurrent cold identical calls collapse to one downstream
  call and share result and `AsOf`; errors are not cached.
  - Traces to milestone criterion ID: M-04-C02
  - Proven by: `go test ./... -run 'FreshnessCache' -v`
- Criterion: MaxAge=0 bypasses lookup and linkage is structurally excluded.
  - Traces to milestone criterion ID: M-04-C03
  - Proven by: `go test ./... -run 'FreshnessCache' -v`
- Criterion: LRU cap bounds entries and logs expose only class/status fields.
  - Traces to milestone criterion ID: M-04-C04
  - Proven by: `go test ./... -run 'FreshnessCache' -v`
- Criterion: successful class mutations evict matching entries while failed
  writes do not invalidate and unrelated classes stay warm.
  - Traces to milestone criterion ID: M-04-C05
  - Proven by: `go test ./... -run 'EvictOnMutation' -v`

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md` and implement the scoped feature.
- Required files/evidence: spec, plan, validation, targeted and full test output
- Blockers or open decisions: None
