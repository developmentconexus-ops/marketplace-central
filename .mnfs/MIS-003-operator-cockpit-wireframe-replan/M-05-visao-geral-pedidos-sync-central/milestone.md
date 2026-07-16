# M-05-visao-geral-pedidos-sync-central

```yaml
id: M-05
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-003 — `../mission.md`; contracts: IC-05 (routes `/` `/pedidos` `/integracoes`, orders/sync ns), IC-02 (summary), IC-03 (protocolos list), ADR-15/17.

## Outcome

The cockpit's aggregate surfaces exist: Visão geral (1e) at `/`, Pedidos read (1j) at `/pedidos`, Integrações & Sync (1k) at `/integracoes` with protocolos list + operation runs. Backed by small new server aggregates: dashboard summary endpoint (composing existing module reads + listings summary), orders list/detail read API over existing orders tables (0027/0033), sync/run listing over `integration_operation_runs`. Observable: `/` renders counter cards + pendências feed with real numbers (null counters render "—"); `/pedidos` lists orders with NF/fulfillment states; `/integracoes` shows installation health, run history, protocolos list linking to `/protocolos/:id`.

## Why This Milestone Exists

Wireframe deck-2's entry screen (Visão geral) and the ops surfaces aggregate every other milestone's data. Building them after M-01/M-03 ensures counters read real seams instead of inventing parallel queries; Dashboard's direct-fetch migration brief lands here.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | aggregate-sync-endpoints | Dashboard summary, orders read API, sync-runs API + OpenAPI/SDK |
| F-02 | visao-geral | Visão geral 1e at `/`, replacing legacy Dashboard |
| F-03 | pedidos-workspace | Pedidos 1j read workspace |
| F-04 | integracoes-sync | Integrações & Sync 1k + protocolos list |

F-01 first (server); F-02/F-03/F-04 sequential after (shared nav/route seam, one writer).

## Dependencies

M-01 (listings summary feeds counters), M-02 (platform), M-03 (protocolos list data), M-04 (sequential before M-05: route rows are disjoint but both milestones edit the same M-02-owned AppRouter/nav files — one writer per seam, no concurrent milestone edits to those files).

## Ownership & Concurrency

Split execution (mission Parallel Execution Plan): **F-01 runs in W1 inside CHIP-SAT**
(depends only on M-01 — backend-only, zero frontend); F-02/F-03/F-04 run in W3 as CHIP-M05
after M-04 CLOSED and F-01 merged.

- F-01 (in CHIP-SAT, W1): OpenAPI/SDK sections = dashboard-summary, orders, sync-runs
  (additive; never touch CHIP-M03's mutation/protocolo sections or M-06's market section).
  Additive contract-lock: composition-root registration lines for these modules. No
  migration block — reads existing tables; new table needed → `REQUEST`, never self-assign.
  **`GET /orders` directive is binding**: evolve existing `listMarketplaceOrders` in place,
  additive-only; duplicate path/operationId forbidden; cursor semantics impossible without
  breakage → `ESCALATION`.
- F-01 closes at feature grain (CHIP-SAT reports its `CLOSED`); this milestone stays open
  until W3.
- F-02..F-04 (CHIP-M05, W3): AppRouter/nav route rows are disjoint from M-06 F-01's but the
  files are shared — hub serializes merges (rebase-then-merge, one at a time). Internal
  order F-02 → F-03 → F-04 sequential stands.
- Governance base anchor: pinned per chip at dispatch (profile §2).

## Risks

- Counter semantics drift: every Visão geral number must name its source query in spec.md; unknown-source counters render "—" (ADR-17), never 0.
- Orders read scope creep: read-only per mission Non-Scope (no faturar/NF emission).

## Done Means

Three workspaces + aggregates per `validation-contract.md` (M05-C01..C09); Dashboard/Integrations legacy pages replaced; governance lanes green.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 with R-03 orders/runs facts + IC-02 summary.
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: none.

## Correction Handoff

Not applicable at planning time.
