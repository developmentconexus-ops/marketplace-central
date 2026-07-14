# M-01-foundation-observability

```yaml
id: M-01
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

MIS-002 — Oracle read re-architecture. See `../mission.md` and `../research/catalog-read-interface-contract.md` (IC-01).

## Outcome

Oracle adapter hardened and observable BEFORE any query reshaping: uncommitted refactor integrated and committed; pool at 12; server + per-route-class deadlines active; latency/slow-query/pool-stats logging live; baseline evidence captured (catalog COUNT, network RTT, EXPLAIN PLAN of the M-02 candidate page query) so M-02 cutover is gated on facts, not hope.

## Why This Milestone Exists

Mission rule "measure before changing" (risks R1, R8, R10). The working tree already contains an unreviewed adapter hardening refactor; it must land first as one intentional commit or every later milestone builds on uncommitted state.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | adapter-hardening-pool-deadlines | integrate + commit the pending adapter refactor; pool default 12; http.Server timeouts; route-class context deadlines (interactive 15s / batch 120s per IC-01) |
| F-02 | observability-baseline | per-port-method latency logging, slow-query log >1s, periodic db.Stats(); baseline evidence via live lane: TGFPRO active COUNT, RTT, EXPLAIN PLAN of candidate page query |

## Dependencies

None (first milestone). Split F-02 into two lanes: the **baseline-evidence lane** (live docker lane: TGFPRO COUNT, RTT, EXPLAIN PLAN — zero code dependency) runs IN PARALLEL with F-01 from minute one; only the **instrumentation lane** (timing decorator, slow-query log, pool stats) waits for F-01's committed adapter surface.

## Feature Parallelization

| Lane | Content | Starts |
| --- | --- | --- |
| A | F-01 adapter hardening (code) | immediately |
| B | F-02 baseline evidence (live lane via `scripts/run-live-oracle-docker.ps1`, read-only) | immediately, parallel worker |
| C | F-02 instrumentation (code) | after F-01 accepted |

Seam rule: lane B writes only evidence markdown (no source files) — no collision possible.

## Risks

R1 (plan evidence may force base-table JOIN fallback in M-02 design), R8 (WAN RTT would raise TTLs), R10 (real COUNT ≫ assumption). All three RESOLVE here — that is the point.

## Done Means

`go build ./...` + `GOCACHE=.gocache go test ./...` green; requests to Oracle-backed routes log method+duration; slow-query lines appear when forced; pool stats logged periodically; baseline numbers recorded in `validation-result.md`; EXPLAIN PLAN verdict recorded (index-supported or fallback-to-base-tables decision for M-02).

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: dispatch lanes A and B in parallel (two mpc-implementer workers, gpt-5.6-luna high); lane C after A accepted
- Required files/evidence: feature `validation.md` files; `validation-result.md` with baseline numbers
- Blockers or open decisions: None

## Correction Handoff

Not applicable during initial planning.
