# Feature Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
split_decision: single
split_reason: one shared cache core and one composition seam form one coherent vertical unit
```

## Steps

1. Add the neutral freshness context setter/getter, cache invalidation port, and
   direct `golang.org/x/sync` dependency.
2. Implement shared generic cache state, canonical keys, TTL policy/env config,
   singleflight, LRU, structured logs, and three typed decorators.
3. Wire one shared cache at the composition root while keeping timing outside
   the catalog decorator and leaving linkage readers untouched.
4. Add post-persistence invalidation to the four write application services.
5. Add table-driven fake-clock, concurrency, bypass, exclusion, LRU, error, and
   mutation tests; run targeted and full Go validation.
6. Record real output in `validation.md`, review the diff, and create one
   intentional conventional commit.

## Files Expected To Change

- `apps/server_core/internal/modules/internal_read/adapters/cache/*`: cache core, decorators, and tests.
- `apps/server_core/internal/modules/internal_read/domain/source_metadata.go`: neutral context seam.
- `apps/server_core/internal/modules/internal_read/ports/cache_invalidation.go`: application-facing invalidation port.
- `apps/server_core/internal/modules/catalog/transport/http_handler.go`: use neutral freshness context seam.
- `apps/server_core/internal/composition/root.go`: shared cache/config and decorator/invalidation wiring.
- `apps/server_core/internal/modules/catalog/application/service.go`: catalog post-commit invalidation.
- `apps/server_core/internal/modules/inventory/application/stock_action_service.go`: applied-action post-commit invalidation.
- `apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service.go`: confirmation post-commit invalidation.
- `apps/server_core/internal/modules/profitability/application/service.go`: margin-input post-commit invalidation.
- `apps/server_core/go.mod`, `apps/server_core/go.sum`: direct singleflight dependency metadata.
- `.mnfs/.../F-01-freshness-cache/{spec,plan,validation}.md`: execution evidence.

## Verification Commands

- Command: `$env:GOCACHE='<absolute worktree>\\.gocache'; go test ./... -run 'FreshnessCache' -v`
  - Satisfies criterion ID: M-04-C01, M-04-C02, M-04-C03, M-04-C04
  - Expected result: PASS; fake-clock TTL, env config, singleflight, bypass,
    structural exclusion, LRU, error, nil-map, and log assertions pass.
- Command: `$env:GOCACHE='<absolute worktree>\\.gocache'; go test ./... -run 'EvictOnMutation' -v`
  - Satisfies criterion ID: M-04-C05
  - Expected result: PASS; matching class is evicted after successful writes,
    other classes survive, and failed writes do not invalidate.
- Command: `$env:GOCACHE='<absolute worktree>\\.gocache'; go test ./...`
  - Satisfies criterion ID: M-04-C01, M-04-C02, M-04-C03, M-04-C04, M-04-C05
  - Expected result: PASS for the complete server_core module.

## QA Steps

- Inspect `git diff --check` and changed-path list; expected no whitespace
  errors and only scoped implementation/evidence paths.
- Inspect cache logs in the unit test; expected `cache`, `key_class`, and no
  raw IDs, query text, key hash, or data values.
- Inspect root construction; expected linkage source remains direct/timed and
  catalog timing wraps the cache.

## Rollback/Risk Notes

The cache is process-local and intentionally loses warm state on restart. If
tests reveal stale or cross-class invalidation, fix the cache seam before
commit; do not broaden scope or change public port signatures. No destructive
git recovery operations are permitted.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, exact command output, commit SHA
- Blockers or open decisions: None
