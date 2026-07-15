# F-02-web-query-state-components

```yaml
id: F-02
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

MIS-003. Binding contract: IC-05 (query contract, crosswalk table, state vocabulary, glossary), IC-03 (failure codes → pt-BR copy). ADR-15 (TanStack Query exclusive), ADR-12 (freshness = fetched_at + manual refresh).

## Milestone

M-02 frontend-platform-anuncios. Depends on F-01 (shell exists; components render inside it).

## Brief

Expand `packages/web-query`: new `QUERY_STALE_TIME` entries (listings 45_000, mutations 5_000, orders 120_000, sync 30_000, market 300_000), new namespaces, exact key builders per IC-05, `invalidateAfterMutation(queryClient, type)` implementing the crosswalk table, and `failureCopy.ts` mapping every IC-03 failure code to fixed pt-BR strings. Build the six state components in `packages/ui` (`LoadingState`, `ErrorState`, `EmptyState`, `FreshnessIndicator` [exists — verify contract], `UnknownValue`, `ConflictTag`) with IC-05 fixed copy. Fix the two passing build defects that block this seam: Tailwind `@source` for new package paths, `/listings` `/mutations` `/market` `/orders` `/profitability` dev-proxy rows.

EARS:
- While a mutation of IC-03 type T settles, when `invalidateAfterMutation(qc, T)` is called, the client shall invalidate exactly the namespaces in the IC-05 crosswalk row for T (proof: unit test per row).
- While a server value is null/unknown, when a table cell renders through `UnknownValue`, the UI shall show "—" with the contextual hint tooltip, never "0" or "R$ 0,00".
- While an IC-03 failure code arrives, when rendered, the UI shall show the failureCopy pt-BR string; unknown codes fall back to "Falha desconhecida ({code})".

## Inputs

- IC-05 Query contract + State vocabulary + crosswalk (verbatim), IC-03 failure-code enum, current `packages/web-query/src/` (QUERY_STALE_TIME, namespaces, createRefreshableFetch), `packages/ui` conventions, R-02 defects list (Tailwind @source, proxy).

## Expected Output

- web-query additions; existing keys/values untouched (IC-05 Compatibility).
- `invalidateAfterMutation` + unit test asserting each crosswalk row (spy on `queryClient.invalidateQueries`).
- `failureCopy.ts` covering all 12 IC-03 codes + fallback.
- Six components exported from `packages/ui` with fixed copy; stories or render tests per state.
- vite proxy + Tailwind @source fixes.

## Constraints

- No page rebuilds here (F-03 does Anúncios).
- Copy strings byte-exact per IC-05 table; pt-BR only.
- Existing `QUERY_STALE_TIME` values unchanged (L2 mirror).
- No new invalidation paths outside `invalidateAfterMutation` for envelope writes.

## Negative Scenarios

- `invalidateAfterMutation(qc, "bogus")` → throws typed error (crosswalk exhaustive, no silent no-op).
- ErrorState with no `error.message` → generic copy "Erro ao carregar." only.
- `UnknownValue` without hint → renders "—", no empty tooltip.

## Validation Expectations

- Vitest output: crosswalk test (one assert per IC-05 row incl. product-enrichment row), failureCopy exhaustiveness test (iterates IC-03 enum), component render tests.
- Diff proof: vite.config proxy rows + Tailwind @source in commit.
- Type-check proof: `tsc` green across packages.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-01 accepted).
- Next action: compile context pack; read IC-05/IC-03 + web-query/ui sources only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
