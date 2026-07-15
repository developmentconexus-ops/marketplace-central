# F-02-visao-geral

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/`, sync ns 30s, states), R-01 screen 1e, F-01 summary endpoint. ADR-17.

## Milestone

M-05. Depends on F-01.

## Brief

Build Visão geral (wireframe 1e) at `/`, replacing legacy Dashboard (its direct-fetch migration brief lands here): counter cards (anúncios ativos, pendências, pedidos hoje/7d, vínculos pendentes, abaixo da margem, sem GTIN) from `getDashboardSummary`, pendências feed with each row deep-linking to the owning workspace pre-filtered (e.g. sync errors → `/anuncios?tab=pendencia&filter.exception=sync_error`), últimas sincronizações strip (per-module last-run + status), quick actions per 1e. Null counters render "—" with degraded-source hint. Poll summary at sync staleTime 30s.

EARS:
- While summary returns, when a counter is null, the card shall render "—" with tooltip "fonte indisponível: {source}" — never 0.
- While the operator clicks a pendência row, when navigating, the target workspace shall open pre-filtered via URL (deep-link contract).
- While `/` loads, when legacy Dashboard is gone, no direct fetch shall run (network proof) and the page shall render entirely from `getDashboardSummary` + context.
- While summary is degraded, when rendered, a non-blocking banner shall list degraded sources; healthy cards stay live.

## Inputs

- R-01 §1e card/feed inventory, F-01 summary shape, IC-05 keys (`syncQueryKeys`/dashboard key in sync ns) + components, legacy Dashboard component (deletion target + parity checklist from R-02).

## Expected Output

- `/` page; legacy Dashboard deleted; route table row updated per IC-05 ("until M-05" note discharged).
- Deep-link map card→URL fixed in spec.md (each card names target route + params).
- Component tests: null-render, deep-link URLs, degraded banner, no-direct-fetch.

## Constraints

- No client-side aggregation across endpoints — one summary call (F-01 owns composition).
- pt-BR; card labels per glossary; no invented metrics beyond 1e.

## Negative Scenarios

- Summary 404 (no installation) → context empty-state flow (M-02 contract).
- Summary total failure → ErrorState with retry; sidebar/shell unaffected.
- Pendência deep link to workspace not yet rebuilt (if M-04 pending) → stub EmptyState, no crash.

## Validation Expectations

- Vitest output: null-render ("—" not "0"), deep-link map, degraded banner tests green.
- Browser proof: 1e screenshot with real seeded counters; click-through of one pendência row landing pre-filtered.
- Grep/network proof: legacy Dashboard gone; zero non-TanStack fetches.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-01 accepted).
- Next action: compile context pack; read R-01 §1e + F-01 output + IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
