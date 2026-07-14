# Milestone Validation Contract

```yaml
id: M-04
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

M-04

## QA Level

QA-2

## Required Outcome

TTL + singleflight cache with correct freshness semantics, bypass, exclusions, evict-on-mutation, and observability.

## Criteria

## Criterion: TTL hit/expiry semantics per data class
ID: M-04-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./... -run 'FreshnessCache' -v` (fake clock)
- Expected: within TTL → 0 Oracle calls, identical `as_of`; after TTL → 1 new call, strictly newer `as_of`; TTLs read from FreshnessPolicy config (catalog 300s, stock 45s, price/cost 120s defaults, env-tunable)
- Actual:
- Artifact: `M-04-server-cache/F-01-freshness-cache/validation.md`
Blocking failure: stale serve past TTL or wrong as_of propagation
Blocking failure observed: No
Owner: QA Validator

## Criterion: Singleflight collapses concurrent misses
ID: M-04-C02
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: concurrency test — 20 goroutines, same key, cold cache, blocking fake reader
- Expected: exactly 1 Oracle call; all 20 receive same page + same `as_of`
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: >1 call for identical concurrent key
Blocking failure observed: No
Owner: QA Validator

## Criterion: no-cache bypass and linkage exclusion
ID: M-04-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest: warm cache → GET with `Cache-Control: no-cache`; repeated linkage candidates calls
- Expected: no-cache → fresh Oracle call + newer `as_of` + cache repopulated; linkage: every call reaches Oracle (log shows no hit ever), regardless of frequency
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: linkage served from cache (blocks provider-write safety) or bypass ignored
Blocking failure observed: No
Owner: QA Validator

## Criterion: Bounded memory and observability
ID: M-04-C04
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: unit test filling cache past entry cap; log capture of hit/miss lines
- Expected: entries never exceed configured cap (LRU eviction, default sized per mission assumption ≤100k products); log lines carry `cache=hit|miss|bypass key_class=<class>` with no key contents leaking data values
- Actual:
- Artifact: F-01 `validation.md`
Blocking failure: unbounded growth or key values (ids/args) logged raw
Blocking failure observed: No
Owner: QA Validator

## Criterion: Evict-on-mutation invalidates matching fact class
ID: M-04-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./... -run 'EvictOnMutation' -v` + httptest: warm catalog cache key → linkage-confirm write → repeat GET inside TTL
- Expected: post-write GET issues a fresh Oracle query and returns strictly newer `as_of` (key evicted, not expired); unrelated fact classes (e.g. `pricecost`) keep their warm entries; failed/rolled-back write triggers NO eviction
- Actual:
- Artifact: `M-04-server-cache/F-01-freshness-cache/validation.md`
Blocking failure: stale L2 entry served after a successful mutation of that fact class
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Fake-clock + concurrency test transcripts; httptest header transcripts; log excerpts.

## Blocking Failures

Per criterion.

## Retry Policy

- correction_attempts: 6
- max_correction_attempts: 2
- last_validation_result: qa-validator PASS at df2cac6a31f68c8980ccfa47da51cb40dfbcacee; Codex gpt-5.6-sol FAIL at the same frozen SHA with 3 blocking evidence-level findings; validation at the new SHA is pending
- overrun: correction rounds 3-6 exceed max_correction_attempts and were user-authorized in the visible Milestone session as deliberate, disclosed test/evidence-only corrections; this is not a silent overrun

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: execute F-01, then QA
- Required files/evidence: F-01 validation.md, `validation-result.md`
- Blockers or open decisions: None
