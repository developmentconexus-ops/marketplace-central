# M-02-frontend-platform-anuncios

```yaml
id: M-02
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

MIS-003 — `../mission.md`; contracts: IC-05 (`../research/frontend-platform-interface-contract.md`), IC-02 (read API consumed), ADR-15.

## Outcome

The SPA has the wireframe deck-2 shell: new sidebar + route map with legacy redirects, `InstallationContext` with URL-persisted selection, web-query expanded with new namespaces/staleTimes/key builders/`invalidateAfterMutation`, the six shared state components in `packages/ui`, and the first rebuilt workspace — Anúncios unificado (screen 2a) — read-only against IC-02 endpoints. Observable: `/anuncios?installation=X&tab=pendencia&filter.exception=sync_error` renders the filtered table and survives F5; `/products` redirects to `/catalogo` with query preserved.

## Why This Milestone Exists

IC-05 is the single-writer seam every later UI milestone consumes. Building shell + platform + one real consumer (Anúncios) in one slice proves the seam before M-04/M-05 fan out in parallel.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | shell-routes-context | Sidebar, route map, redirects, InstallationContext, top-bar pills |
| F-02 | web-query-state-components | Query namespaces/keys/staleTimes, invalidateAfterMutation, six state components, failureCopy module |
| F-03 | anuncios-workspace | Anúncios unificado screen (2a) read-only on IC-02 |

Order F-01 → F-02 → F-03. No parallel: all three touch IC-05-owned files (one writer).

## Dependencies

M-01 (IC-02 endpoints live; F-03 consumes them). Mutation actions in 2a render disabled with tooltip "disponível em breve" until M-03 (buttons present per wireframe, no dead handlers).

## Ownership & Concurrency

Wave W1 (mission Parallel Execution Plan) — runs concurrent with CHIP-M03 and CHIP-SAT.

- Owns exclusively: frontend platform seam — `AppRouter`, `Layout`/nav, redirects,
  `InstallationContext`, web-query namespaces, state components, failureCopy. No other W1
  chip touches frontend files.
- Migrations: none. OpenAPI/SDK: none (read-only consumer of M-01 SDK).
- Composition root: not touched.
- Outbound edge: F-03 (Anúncios workspace) unblocks M-03 F-04 FE work — report F-03 in the
  `COMMITTED` event so the hub can trigger CHIP-M03's rebase without waiting for CLOSED.
- Governance base anchor: pinned in chip prompt at dispatch (profile §2).

## Risks

- RK-05 (legacy direct-fetch pages regress under new router): redirects tested; legacy pages keep working unrebuilt under new Layout.
- RK-04 (seam contention on M-02-owned files): mitigated by this milestone's one-writer feature order F-01 → F-02 → F-03 and by later milestones editing AppRouter/nav sequentially (M-04 → M-05 per mission dependency graph).
- Scope creep into other workspaces: only route stubs for M-04/M-05 rows ("em construção" EmptyState), no partial builds.

## Done Means

All IC-05 seam elements exist and Anúncios proves them per `validation-contract.md` (M02-C01..C11); pt-BR-only copy in new surfaces; no direct fetch in new code; governance lanes green.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 with feature context pack (IC-05 + current AppRouter/Layout paths).
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: none.

## Correction Handoff

Not applicable at planning time.
