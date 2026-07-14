# Feature Validation

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Summary

Feature quick validation passed. The implementation is ready for Milestone
Orchestrator review; this is not a milestone or QA verdict.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-14
- Final feature state for handoff: `quick_validation_passed`

## Evidence Honesty

All passing evidence below is `ran` and includes real command output excerpts.
The first cold full-suite attempt timed out during compilation; the exact same
full command was rerun after the targeted build and passed.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None
- Reason: Port signatures and public contracts were unchanged; no OpenAPI or
  sdk-runtime edit was required.

## Changes Made

- `apps/server_core/internal/modules/internal_read/adapters/cache/cache.go`:
  shared typed decorators, SHA-256 canonical keys, FreshnessPolicy TTL table,
  injected clock, singleflight, nil-preserving maps, LRU cap, invalidation, and
  structured redacted logs.
- `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`:
  fake-clock, table-style TTL, concurrency, bypass, linkage, LRU, logging,
  error, nil-map, and invalidation tests.
- `apps/server_core/internal/modules/internal_read/domain/source_metadata.go`,
  `.../catalog/transport/http_handler.go`: neutral freshness context seam.
- `apps/server_core/internal/modules/internal_read/ports/cache_invalidation.go`:
  application-facing post-commit invalidation port.
- `apps/server_core/internal/composition/root.go`: shared cache composition;
  timing remains outside catalog cache and linkage remains direct.
- Catalog, inventory, orders, and profitability application services: invalidate
  only after successful persistence; catalog test covers failed-write behavior.
- `apps/server_core/go.mod`: promoted existing `golang.org/x/sync v0.14.0` to a
  direct dependency; `go mod tidy` made no additional dependency change.

## Commands Run

- Command: `go mod tidy`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: x/sync direct and tidy dependency metadata
  - Actual: exit code 0; x/sync is direct in `go.mod`; no `go.sum` delta
  - Artifact: `apps/server_core/go.mod`

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'; go test ./... -run 'FreshnessCache' -v`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: C01-C04 scenarios pass
  - Actual excerpt:

    ```text
    === RUN   TestFreshnessCacheTTLPerClass
    --- PASS: TestFreshnessCacheTTLPerClass (0.01s)
    === RUN   TestFreshnessCacheSingleflight
    --- PASS: TestFreshnessCacheSingleflight (0.00s)
    === RUN   TestFreshnessCacheErrorNotCached
    --- PASS: TestFreshnessCacheErrorNotCached (0.05s)
    === RUN   TestFreshnessCacheBypassAndLinkageExclusion
    --- PASS: TestFreshnessCacheBypassAndLinkageExclusion (0.00s)
    === RUN   TestFreshnessCacheLRUAndLogs
    --- PASS: TestFreshnessCacheLRUAndLogs (0.00s)
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 3.405s
    ```

  - Artifact: `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'; go test ./... -run 'EvictOnMutation' -v`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: C05 successful class eviction, unrelated-class retention, and
    failed-write no-invalidation
  - Actual excerpt:

    ```text
    === RUN   TestEvictOnMutationFailedWriteDoesNotInvalidate
    --- PASS: TestEvictOnMutationFailedWriteDoesNotInvalidate (0.00s)
    === RUN   TestEvictOnMutation
    --- PASS: TestEvictOnMutation (0.01s)
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 3.467s
    ```

  - Artifact: `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/F-01-freshness-cache/validation.md`

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'; go test ./...`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: complete server_core suite green
  - Actual excerpt:

    ```text
    ok   marketplace-central/apps/server_core/internal/modules/catalog/application (cached)
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 1.463s
    ok   marketplace-central/apps/server_core/internal/modules/inventory/application (cached)
    ok   marketplace-central/apps/server_core/internal/modules/orders/application (cached)
    ok   marketplace-central/apps/server_core/internal/modules/profitability/application (cached)
    ok   marketplace-central/apps/server_core/tests/integration (cached)
    ok   marketplace-central/apps/server_core/tests/unit (cached)
    ```

  - Artifact: `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'; go test ./...` (initial cold run)
  - Status: Timeout; superseded by the passing rerun above
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: complete server_core suite green
  - Actual: command timed out after 124.2s during cold compilation, before test
    failure output; rerun with the same command passed.
  - Blocking condition: None after successful rerun

- Command: `git diff --check`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no whitespace errors
  - Actual: exit code 0; only Git line-ending normalization warnings
  - Artifact: working-tree diff

## Manual QA

- QA level: QA-2
- Flow or step: structural composition inspection
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: `Timing(Cache(Oracle))` on catalog; linkage reader is not decorated;
  batch readers have cache without new timing
- Actual: root wiring confirms `NewTimingReader(NewReader(...))`; linkage uses
  its direct timing wrapper; cost/tax and stock use cache-only decorators.

## Evidence Mapped to M-04 Criteria

| Criterion | Evidence | Result |
| --- | --- | --- |
| M-04-C01 | `TestFreshnessCacheTTLPerClass`: fake clock, catalog/stock/pricecost TTLs, env values, identical/newer `AsOf` | Pass |
| M-04-C02 | `TestFreshnessCacheSingleflight` and `TestFreshnessCacheErrorNotCached`: 20 waiters, one successful call, shared result, retry after error | Pass |
| M-04-C03 | `TestFreshnessCacheBypassAndLinkageExclusion`: MaxAge=0 repopulates; linkage reaches downstream every time | Pass |
| M-04-C04 | `TestFreshnessCacheLRUAndLogs`: cap remains 2, LRU is evicted, logs contain only status/class fields | Pass |
| M-04-C05 | `TestEvictOnMutation` plus `TestEvictOnMutationFailedWriteDoesNotInvalidate`: matching class evicts, pricecost survives, failed persistence does not invalidate | Pass |

## Design-Decision Notes

- Context-based MaxAge: the decorator reads the neutral
  `FreshnessPolicyFromContext`; a present zero policy bypasses lookup and
  repopulates after success. No port signature changed.
- Timing-outside-cache: the root composes `TimingReader(CacheReader(Oracle))`
  so cache hits are timed and logged as hits without timing batch readers.
- No signature change: OpenAPI and sdk-runtime remain untouched because the
  decorators implement the existing ports transparently.
- Structural linkage exclusion: `Reader` only overrides catalog page methods;
  linkage methods remain promoted direct downstream calls, and the separate
  Sankhya linkage source is never passed through cache.
- Freshness values: catalog pages preserve stored `CatalogFactPage.AsOf`; batch
  entries capture the injected clock at insertion. Nil maps and nil map values
  are preserved.

## Risks

- Cache is intentionally per-process and has no cross-process invalidation bus.
- The default cap is 100,000 entries; env configuration is accepted only for
  positive values and positive TTL durations.
- Milestone-level QA and fixed-SHA review remain outstanding.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review spec, plan, changed paths, commit, and evidence; then run
  fixed-SHA review and proportional QA.
- Required files/evidence: feature brief, `spec.md`, `plan.md`, this validation,
  targeted test transcripts, full-suite result
- Blockers or open decisions: None
