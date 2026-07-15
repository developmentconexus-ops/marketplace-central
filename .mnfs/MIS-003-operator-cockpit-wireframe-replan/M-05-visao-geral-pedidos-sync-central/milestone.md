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
