# Milestone Validation Contract

```yaml
id: M-05
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-05

## QA Level

QA-2

## Required Outcome

TanStack Query L1 cache adopted; freshness visible; mutations invalidate correctly.

## Criteria

## Criterion: Infinite catalog pagination over cursor envelope
ID: M-05-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `npm test` (component test with mocked sdk client) in `apps/web`
- Expected: `useInfiniteQuery` consumes IC-01 envelope; `getNextPageParam` returns `next_cursor`; null cursor stops fetching; pages appended without refetch of prior pages
- Actual:
- Artifact: `M-05-web-tanstack/F-01-tanstack-adoption/validation.md`
Blocking failure: offset-style params or page-1 refetch per scroll
Blocking failure observed: No
Owner: QA Validator

## Criterion: staleTime prevents redundant refetch
ID: M-05-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: component test with fetch spy — mount view, unmount, remount within staleTime
- Expected: 1 network call total; staleTime per IC-01 (catalog 300s, stock 45s, price/cost 120s); remount after staleTime → refetch
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: refetch on every mount
Blocking failure observed: No
Owner: QA Validator

## Criterion: as_of indicator and manual refresh
ID: M-05-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: component test — render with envelope fixture; click refresh
- Expected: as_of rendered as local time indicator; refresh triggers request with `Cache-Control: no-cache` header and updates indicator with new as_of
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: refresh serving cached as_of
Blocking failure observed: No
Owner: QA Validator

## Criterion: Mutations invalidate correct namespaces; linkage never stale
ID: M-05-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: component test — perform linkage-confirm mutation with spy on invalidateQueries; inspect linkage query options
- Expected: mutations invalidate exactly the IC-01 crosswalk namespaces (linkage confirm → `['linkage']`+`['catalog']`; margin-input import → `['catalog']`+`['profitability']`); linkage candidate query configured staleTime 0, gcTime 0 (fresh every open)
- Actual:
- Artifact: `M-05-web-tanstack/F-02-mutation-invalidation/validation.md`
Blocking failure: stale linkage candidates reachable after mutation
Blocking failure observed: No
Owner: QA Validator

## Criterion: Build health and no bypassing fetches
ID: M-05-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `npm run build` in `apps/web`; grep migrated views for direct `fetch(`/raw client calls outside hooks
- Expected: build green; migrated Oracle-backed views call sdk only through query/mutation hooks
- Actual:
- Artifact: F-02 `validation.md`
Blocking failure: build red or direct fetch in migrated views
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Component-test transcripts + spy assertions; build output; TanStack v5 API claims re-verified at install (`npm view @tanstack/react-query version` + changelog note in validation.md).

## Blocking Failures

Per criterion.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: execute F-01 → F-02, then QA; then mission QA (MIS-02-C01..C05)
- Required files/evidence: feature validation.md files, `validation-result.md`
- Blockers or open decisions: None
