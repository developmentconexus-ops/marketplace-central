# ADR-021: Frontend platform seam — one shell, TanStack Query exclusive for server state

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-003's frontend work but no document was ever
written under its original number (ADR-15). It is reconstructed here from the 14 live
citations of `ADR-15`, harvested at
`docs/architecture/decisions/_citations/adr-015-citations.md`. Every clause below is
traceable to the mission's binding contracts or to code that already asserts it.

## Context

Before MIS-003, each page in the operator frontend owned its own data fetching: its own
`useEffect`/`useState` fetch, its own retry and loading logic, its own ad-hoc cache
invalidation after a write. That pattern reinvents the same context/invalidation logic on
every page and gives each page a different, undocumented failure and loading vocabulary —
there is no single place to fix a stale-cache bug or add a freshness indicator.

MIS-003 needed a single frontend platform seam: one shell that owns routing, installation
context, and query invalidation, with one server-state mechanism used everywhere new or
rebuilt screens are built.

## Decision

**A single SPA shell owns routing, installation context, and query-key invalidation.
TanStack Query is the exclusive mechanism for server state on any new or rebuilt screen:
no per-page ad-hoc `fetch`, `localStorage` cache, or `useEffect`-driven fetch loop.**

**§1 — One shell, single writer.** The SPA shell is the single writer of routing and
installation context; screens consume it rather than re-deriving it.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-01-shell-routes-context/feature.md:17` — "ADR-15 (frontend platform seam; SPA shell single writer)."

**§2 — TanStack Query is exclusive for server state.** New and rebuilt screens read and
invalidate server state only through the shared TanStack Query contract (query keys,
crosswalk table, state vocabulary defined in IC-05); no page-local fetch/cache
mechanism is introduced alongside it.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-02-web-query-state-components/feature.md:17` — "ADR-15 (TanStack Query exclusive), ADR-12 (freshness = fetched_at + manual refresh)."

**§3 — A shared six-state vocabulary and pt-BR-only copy bind new surfaces.** Loading,
error, empty and related UI states use one vocabulary across screens, and copy on
new/rebuilt surfaces is pt-BR only.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:211-212` — Q5 Usability / Q6 Maintainability tie pt-BR copy, the 6-state vocabulary, and "TanStack-only server state" to ADR-15.

**§4 — Named carve-out: legacy pages keep direct fetch until rebuilt.** Classifications
and Marketplaces are the two legacy pages that existed before this seam. They are
explicitly exempted: they keep their own direct data-fetch until each is individually
rebuilt onto the platform. This is a recorded, accepted trade-off, not a violation of §2.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:118` — "ADR-15 — legacy pages (Classifications, Marketplaces) keep direct fetch until rebuilt; recorded as migration briefs, not silent debt."
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/PLAN-ADJUDICATION.md:12` — "legacy pages mounting at new paths keep their fetch, recorded debt not defect."

**§5 — The pt-BR-only rule binds new/rebuilt surfaces only.** The carve-out in §4 extends
to copy: old, not-yet-migrated legacy pages are not required to be pt-BR-only.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/PLAN-ADJUDICATION.md:12` — "pt-BR blocking rule binds new/rebuilt surfaces only."

**§6 — Retiring the legacy pages is the seam's own success criterion.** The read
workspaces milestone frames rebuilding the remaining legacy pages onto the platform as
proof the seam holds under real fan-out, not a separate concern.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-04-read-workspaces/milestone.md:25` — "retiring the direct-fetch legacy pages that motivated ADR-15."

## Rationale

Centralizing routing, installation context, and server-state caching in one shell means a
fix to invalidation, freshness, or error handling is made once and applies to every screen
built after it. Naming the legacy carve-out explicitly (§4) rather than leaving it silent
means the two pages that still bypass the seam are a tracked migration item, not an
unnoticed regression discovered later.

## Consequences

- Every new or rebuilt screen must be planned and reviewed against the TanStack Query
  contract (query keys, invalidation, six-state vocabulary) — a page cannot introduce its
  own fetch/cache logic to avoid that discipline.
- Classifications and Marketplaces remain outside the seam by design until rebuilt; any
  audit of "is this page on the platform" must check this carve-out before flagging a
  gap.
- **Verified against the current codebase (2026-08-05):** the carve-out is still in
  force. `packages/feature-classifications/src/ClassificationsPage.tsx:93-112` loads data
  through a `useCallback`/`useEffect` pair calling the client directly (`loadData`,
  `client.listCatalogProductFacts`, `client.listClassifications`) with local
  `useState` for `loading`/`error`/`products`/`classifications` — no TanStack Query hook
  is used anywhere in the file. This page is still mounted at `/classifications` in
  `apps/web/src/app/AppRouter.tsx:73` (`ClassificationsPageWrapper`).
  `/marketplaces` (`apps/web/src/app/AppRouter.tsx:74`) now renders
  `WorkspacePlaceholder` (`apps/web/src/pages/WorkspacePlaceholder.tsx:1-8`), a static
  "em construção" stub with no data fetch of any kind — it no longer bypasses the seam
  because it does not fetch server state at all; it has not yet been rebuilt onto
  TanStack Query either.

## Alternatives Considered

**Migrate every legacy page before starting new platform work.** Rejected: the mission
needed the shell and query seam available for new/rebuilt screens immediately; blocking
on a full legacy rewrite first would have delayed the platform itself. The named
carve-out let the seam ship while tracking the remaining migration explicitly.

**Leave the exemption unrecorded and let each legacy page be judged case by case.**
Rejected: an unrecorded exception reads as an accidental gap in every later audit. Naming
the two pages and the condition for retiring the exemption (rebuild) turns it into a
tracked migration brief instead of silent debt.
