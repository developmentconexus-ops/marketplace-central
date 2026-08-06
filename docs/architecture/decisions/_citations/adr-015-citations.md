# ADR-015 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** ~14
**Spellings found:** ADR-15 only.

All citations of ADR-15 live in a single mission (`MIS-003-operator-cockpit-wireframe-replan`); unlike ADR-13/14 there is no cross-mission collision — the assertion is consistent everywhere it appears.

## Assertion A1 — Frontend platform seam: one shell owns routing/context/query-invalidation; TanStack Query is the exclusive server-state mechanism (no per-page ad-hoc fetch/localStorage/useEffect-fetch); this seam is what retires the legacy direct-fetch pages
- Citations: 14
- Verbatim: "Frontend platform seam | decided | per-page context/invalidation reinvention | IC-05 route map, redirects, query contract, 6-state vocabulary, pt-BR-only new copy"
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:114` (ADR table row)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:211-212` (Q5 Usability / Q6 Maintainability tie pt-BR copy, 6-state vocabulary, and "TanStack-only server state" to ADR-15)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/frontend-state.md:50` ("Grounds ADR-15 and IC-05")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-01-shell-routes-context/feature.md:17` ("ADR-15 (frontend platform seam; SPA shell single writer)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-02-web-query-state-components/feature.md:17` ("ADR-15 (TanStack Query exclusive), ADR-12 (freshness = fetched_at + manual refresh)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/milestone.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-04-read-workspaces/milestone.md:25` ("retiring the direct-fetch legacy pages that motivated ADR-15")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-04-read-workspaces/F-02-catalogo-workspace/feature.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-05-visao-geral-pedidos-sync-central/milestone.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/PLAN-ADJUDICATION.md:12` (adjudication ruling that a specific legacy-page fetch is a recorded ADR-15 trade-off/debt, not a defect)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-01-shell-routes-context/validation.md:26` (same trade-off cited as validation evidence)

## Contradictions
None found. All 14 citations restate the same rule (frontend platform seam / TanStack-only server state / pt-BR-only new copy / 6-state vocabulary) without divergence.

## Exceptions / carve-outs
- `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md` Migration Briefs (referenced at `M-02-frontend-platform-anuncios/PLAN-ADJUDICATION.md:12`): ADR-15 explicitly accepts that **legacy pages keep their own direct data-fetch** (not migrated onto the TanStack-only seam) until each page is individually rebuilt; this is a named, accepted trade-off ("recorded debt not defect"), not a violation of the ADR-15 rule.
- The pt-BR-only-copy obligation under ADR-15 is scoped to **new/rebuilt surfaces only** — old, not-yet-migrated legacy pages are exempt (`M-02-frontend-platform-anuncios/PLAN-ADJUDICATION.md:12`, "pt-BR blocking rule binds new/rebuilt surfaces only").
