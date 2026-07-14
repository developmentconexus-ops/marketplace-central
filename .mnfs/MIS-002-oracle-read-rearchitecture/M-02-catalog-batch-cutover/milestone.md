# M-02-catalog-batch-cutover

```yaml
id: M-02
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

MIS-002 (`../mission.md`). Governing contract: IC-01 (`../research/catalog-read-interface-contract.md`).

## Outcome

Catalog listing goes from 1+3N sequential Oracle calls to exactly ONE set-based keyset-paginated query per page, exposed via the IC-01 cursor envelope on the existing `/catalog/*` route namespace (IC-01: no new prefixes); OpenAPI and `sdk-runtime` updated together; old entity-composition path removed from the catalog listing flow.

## Why This Milestone Exists

This is the hot path causing the user-visible slowness (mission root cause: `catalog/adapters/internalread/reader.go:67-111`). It also establishes the page-port + envelope pattern M-03 replicates.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | catalog-page-port | new port `ListCatalogProductFacts(ctx, cursor, limit)` + search variant `SearchCatalogProductFacts(ctx, q, limit)` + Oracle adapter with the IC-01 single JOIN/subquery SQL, keyset pagination, nullable facts + quality flags |
| F-02 | catalog-api-envelope | transport routes (list + `/catalog/products/search`) serving IC-01 envelope `{items,next_cursor,page_size,as_of}`; cursor/limit validation; error matrix; OpenAPI + sdk-runtime in same change |

## Dependencies

M-01 passed (baseline facts). GATE: consume M-01-C04 plan verdict — if `FULL_SCAN`, F-01 implements the recorded fallback (JOIN base TGF* tables / adjusted predicates) instead of the candidate query as-is.

F-02 does NOT wait for F-01 completion — interface-first handshake: port signatures and page types are already contract-fixed (IC-01 + F-01 brief). F-01's FIRST deliverable step is the port interface + domain types file (ports package); as soon as that lands, F-02 dispatches and builds transport/envelope/OpenAPI/SDK against the interface using a fake port implementation. Integration when F-01's Oracle adapter is accepted.

## Feature Parallelization

| Lane | Content | Starts |
| --- | --- | --- |
| A | F-01: port interface + types (step 1), then Oracle adapter + tests | immediately |
| B | F-02: transport routes, envelope, error matrix, OpenAPI + SDK against fake port | after lane A step 1 (interface handshake) |

Seam ownership: F-01 owns `internal_read/ports` + oracle adapter files; F-02 owns transport/catalog routes + OpenAPI + `sdk-runtime`. Neither edits the other's files; composition wiring integrated by orchestrator at accept time.

## Risks

R1 (view performance — resolved input from M-01), R2 (cursor semantics drift — IC-01 fixes base64 CODPROD keyset), R5 (OpenAPI/SDK drift — single-commit rule).

## Done Means

Fake-queryer test proves 1 Oracle query per page call; envelope matches IC-01 examples byte-shape; invalid cursor → 400 `invalid_cursor`; limit out of range (0 or 101) → 400 `invalid_limit`; missing facts → null + quality flags; OpenAPI + sdk-runtime updated in the same feature commit; old N+1 composition no longer reachable from catalog listing.

## Handoff

- Current status: planned
- Next owner: Milestone Orchestrator
- Next action: verify M-01-C04 facts exist, dispatch F-01; dispatch F-02 the moment F-01's port-interface step lands (two parallel workers, gpt-5.6-luna high)
- Required files/evidence: feature `validation.md` files; `validation-result.md`
- Blockers or open decisions: plan-verdict gate from M-01

## Correction Handoff

Not applicable during initial planning.
