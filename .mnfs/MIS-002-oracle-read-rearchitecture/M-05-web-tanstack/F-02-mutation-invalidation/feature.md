# F-02-mutation-invalidation

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. IC-01: mutation → queryKey invalidation mapping; linkage freshness rule.

## Milestone

M-05. Depends on F-01 (QueryClient + hooks exist).

## Brief

Wrap write flows in `useMutation` with correct `invalidateQueries` per the IC-01 Invalidation Crosswalk table (linkage confirm → `['linkage']` + `['catalog']`; product edits → `['catalog']`; any stock-affecting action → `['inventory']`; margin-input import → `['catalog']` + `['profitability']` — IC-01 reserved namespaces only); configure linkage candidate queries staleTime 0 / gcTime 0 so candidates are always fetched fresh (mirrors server never-cache rule, MIS-02-C05).

## Inputs

- F-01 hooks + queryKey constants module.
- Existing mutation call sites in `apps/web` (linkage confirm, product edit flows) — locate via sdk-runtime write methods.
- IC-01 invalidation table.

## Expected Output

- `useMutation` wrappers with `onSuccess` invalidation per table; linkage query options staleTime 0 gcTime 0; optimistic updates NOT used (lean; refetch is cheap post M-04).
- One intentional commit.

## Constraints

- Invalidation keys come from the shared constants module (no inline string keys — drift risk).
- Do not change server write endpoints.
- No new UI beyond feedback already present.

## Interaction Model

- Confirm linkage → mutation → success → linkage list + catalog views show refreshed data (spinnerless if within staleTime windows invalidated → refetch fires).
- Mutation error → existing error feedback; no invalidation on error.

## Negative Scenarios

- While a mutation fails, when onError runs, the system shall NOT invalidate queries (data unchanged server-side).
- While a linkage candidate panel reopens, when the query mounts, the system shall always hit the network (staleTime 0, gcTime 0).

## Validation Expectations

- Component tests: invalidateQueries spy called with exact namespaces per mutation type; linkage query options asserted; failure path no-invalidation test.
- `npm run build` green; grep proves no inline queryKey strings in migrated code.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
