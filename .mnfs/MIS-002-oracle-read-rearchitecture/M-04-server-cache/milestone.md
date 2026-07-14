# M-04-server-cache

```yaml
id: M-04
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Mission

MIS-002 (`../mission.md`), IC-01 cache TTL table + `as_of` semantics.

## Outcome

L2 server cache live: in-memory TTL cache honoring FreshnessPolicy MaxAge per data class (catalog 5min, stock 45s, price/cost 2min, linkage never) with `golang.org/x/sync/singleflight` collapse of concurrent misses; `Cache-Control: no-cache` → MaxAge=0 bypass; evict-on-mutation interface consumed by write paths (IC-01 cache semantics, risk R9); `as_of` reflects true snapshot time on every Oracle-backed response.

## Why This Milestone Exists

After M-02/M-03, every read is 1 query but still 1 query per request. 2–10 operators refreshing dashboards would hammer Oracle for identical pages; TTL cache + singleflight makes repeated reads free and bounds Oracle load to 1 query per key per TTL window (mission ADR-05, journeys J3/J6).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | freshness-cache | generic cache decorator over the read ports (key = port method + canonical args hash — no tenant component, see mission accepted assumptions), TTL per FreshnessPolicy, singleflight, no-cache bypass, evict-on-mutation interface (`InvalidateClass`) wired to write paths, `as_of` = snapshot time from cache entry; linkage ports explicitly excluded; cache-hit/miss structured log fields |

## Dependencies

M-02 and M-03 passed (ports to decorate exist and are stable).

## Risks

R7 (stale ops decisions — TTLs from mission accepted assumptions; no-cache escape hatch; stock lowest TTL 45s), R9 (memory growth — entry cap + LRU eviction, cap sized for ≤100k products in accepted assumptions).

## Done Means

Repeat request within TTL: 0 Oracle queries, same `as_of`; after TTL or with no-cache: fresh query, newer `as_of`; 20 concurrent cold requests for same key → 1 Oracle query (singleflight); write that affects a fact class evicts matching L2 keys — next read is fresh (newer `as_of`) even inside TTL; linkage repeat → always Oracle; hit/miss visible in logs; `GOCACHE=.gocache go test ./...` green.

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: dispatch F-01 (mpc-implementer, gpt-5.6-luna high)
- Required files/evidence: F-01 `validation.md`; `validation-result.md`
- Blockers or open decisions: None

## Correction Handoff

Not applicable during initial planning.
