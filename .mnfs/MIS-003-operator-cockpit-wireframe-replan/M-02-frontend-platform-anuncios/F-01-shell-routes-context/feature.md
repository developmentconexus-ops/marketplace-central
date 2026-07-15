# F-01-shell-routes-context

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contract: IC-05 `../../research/frontend-platform-interface-contract.md` (route table, redirect map, context contract, sidebar order). ADR-15 (frontend platform seam; SPA shell single writer).

## Milestone

M-02 frontend-platform-anuncios.

## Brief

Rebuild the SPA shell to the wireframe deck-2 chrome: sidebar in IC-05 order, route map + `<Navigate replace>` redirects (query preserved), `InstallationContext` (single `listIntegrationInstallations` fetch, URL `?installation=` persistence, first-installation fallback with URL rewrite), top-bar pills `ML: <nome> ▾` (functional) and empresa (static single-option). Routes owned by later milestones mount an "em construção" `EmptyState` stub.

EARS:
- While the app loads any route, when `?installation=` names a known installation, the context shall select it; when it names an unknown one, the context shall fall back to the first installation AND rewrite the URL param.
- While the operator switches installation in the top-bar pill, when selection changes, the URL param shall update on the current route without a full reload.
- While a legacy URL is visited (any IC-05 redirect row), when the router resolves, the app shall land on the new route with the query string preserved.
- While no installations exist, when any workspace renders, the context shall report `status="empty"` and pages shall render `EmptyState` with hint "Conecte uma conta em Integrações".

## Inputs

- IC-05 route table + redirect map + context contract (verbatim), current `apps/web/src/app/` AppRouter/Layout files, R-02 (current per-page `?installation=` pattern to replace), sdk-runtime `listIntegrationInstallations`.

## Expected Output

- `apps/web/src/app/InstallationContext.tsx` + `useInstallation()` hook.
- New Layout/sidebar per IC-05 order; AppRouter with full route table + redirects; stub elements for M-03/M-04/M-05 routes.
- Legacy pages still reachable at NEW paths (rendered unmodified under new shell) until their rebuild milestone.
- Component tests: redirect map (each row), unknown-installation URL rewrite, context single-fetch (mock asserts one call across route changes).

## Constraints

- Do not rebuild legacy page internals; wrap only.
- Context persistence URL-only, no localStorage (IC-05 Persistence).
- Sidebar pt-BR labels exactly per IC-05; Mercado omitted.
- Routes not in the IC-05 table must not be added.
- `/classifications`, `/marketplaces` stay routable, off sidebar.

## Negative Scenarios

- Unknown `?installation=inst_ghost` → first installation selected, URL rewritten (no silent mismatch).
- Zero installations → `status="empty"`, no crash, EmptyState hint rendered.
- Direct deep link to stub route → stub EmptyState, shell intact.

## State Model

Context states: `loading` (installations fetch in flight) → `ready` (selection valid) | `empty` (no installations) | `error` (fetch failed → ErrorState with retry re-running the single fetch). Selection changes never re-enter `loading` (no refetch on switch).

## Validation Expectations

- Vitest run output: redirect table test (all 6 rows), URL-rewrite test, single-fetch test green.
- Browser proof (screenshot or DOM assert): `/inventory/stock-seguro?installation=inst_1` lands `/estoque?installation=inst_1`; sidebar order matches IC-05.
- Reload-survival proof: route + `?installation=` identical after F5 on `/anuncios` stub.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; read IC-05 + named app files only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
