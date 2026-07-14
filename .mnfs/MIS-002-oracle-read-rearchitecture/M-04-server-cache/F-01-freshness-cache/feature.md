# F-01-freshness-cache

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. IC-01 cache TTL table, `as_of` semantics, no-cache mapping. Mission ADR-05, MIS-02-C05.

## Milestone

M-04 (single feature).

## Brief

Generic caching decorator wrapping the read ports at composition root: key = port method + canonical args hash (no tenant component — Oracle read ports have no tenant dimension, mission accepted assumption); value = result + snapshot timestamp; TTL from FreshnessPolicy per data class; `golang.org/x/sync/singleflight` for miss collapse; MaxAge=0 (from `Cache-Control: no-cache`, plumbed in M-02 F-02) bypasses and repopulates; LRU entry cap; linkage ports NEVER decorated. Plus evict-on-mutation: the cache exposes `InvalidateClass(factClass)` (classes: `catalog`, `inventory`, `pricecost`) and write paths call it per the IC-01 Invalidation Crosswalk table — linkage confirm → `catalog`; stock-affecting actions → `inventory`; product edits → `catalog`; margin-input imports → `pricecost` (risk R9).

## Inputs

- Ports from M-02/M-03 (page port, batch ports) — decorate these; existing FreshnessPolicy type in domain (honor MaxAge).
- `Cache-Control` → MaxAge=0 mapping already wired in transport (M-02 F-02).
- IC-01 TTL defaults: catalog 300s, stock 45s, price/cost 120s, linkage 0/never. Env knobs `MPC_CACHE_TTL_*`.
- Timing decorator from M-01 (compose: timing outside cache, so hits log ~0ms with `cache=hit`).

## Expected Output

- `internal_read/adapters/cache` (or equivalent adapter-layer package) decorator + composition wiring; dependency `golang.org/x/sync` added intentionally (single go.mod change, no other deps).
- `InvalidateClass(factClass)` on the decorator; existing write application services (linkage confirm, stock actions, product edits, margin-input import) call it post-commit via a small port injected at composition root.
- `as_of` on responses = cached snapshot time, not response time.
- Structured log fields `cache=hit|miss|bypass key_class=<catalog|inventory|pricecost>`.
- One intentional commit.

## Constraints

- Linkage exclusion is structural (linkage ports never pass through decorator), not a TTL=0 configuration that someone can flip.
- Key hashing must not log or expose raw args; log key_class only.
- No external cache (Redis etc.) — in-memory only (mission Non-Scope).
- Cache is per-process; no invalidation bus (accepted assumption).
- `GOCACHE=.gocache`.

## Inputs/Outputs

Decorator is transparent: same port interfaces in/out. Cached entry stores the full result struct including quality flags; `AsOf` field sourced from entry creation time.

## Negative Scenarios

- While an Oracle call errors on miss, when concurrent waiters exist (singleflight), the system shall propagate the error to all waiters and cache NOTHING (no negative caching).
- While the entry cap is reached, when a new key arrives, the system shall evict LRU and store the new entry.
- While MaxAge=0, when the call runs, the system shall skip lookup, query Oracle, store fresh entry.
- While a cache entry for a fact class is warm and inside TTL, when a write path calls `InvalidateClass` for that class, the system shall evict all matching keys so the next read queries Oracle and returns a strictly newer `as_of`.
- While a write fails (rolled back), when the write path errors, the system shall NOT call `InvalidateClass` (cache stays consistent with unchanged data).

## Validation Expectations

- Fake-clock TTL tests per class; 20-goroutine singleflight test; bypass test; linkage exclusion test; cap/LRU test; error-not-cached test; evict-on-mutation test (warm entry inside TTL → InvalidateClass → next read hits Oracle, newer `as_of`; other classes untouched).
- `GOCACHE=.gocache go test ./...` green.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
