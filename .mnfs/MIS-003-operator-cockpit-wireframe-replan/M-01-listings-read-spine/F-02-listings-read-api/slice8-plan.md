# Slice 8 corrective plan — optional margin-reader degradation

Scope: one behavioral change in `read_service.go` plus test-first coverage in `read_service_test.go`. Keep repository, policy, installation, and timeline failures hard; only `GetICMSCeilingByOrigin` and `GetCostFactsByIDs` degrade. `Summary` remains the reference and needs no edit.

## 1. Add failing tests first

- Replace the old hard-failure expectations at `read_service_test.go:194-201` and `:267-275` with table-driven outage tests covering `ceilingErr` and `costErr` separately for `List`, `ByProduct`, and `Get`. For all six cases assert: nil service error, row/group/detail still present, `Cost == nil`, and `BelowMarginWorstCase == nil`; for `Get`, also assert the timeline is still returned. On ceiling failure, assert the cost reader is not called, matching `Summary:59-67`; on cost failure, assert one attempted batch, matching `Summary:78-83`.
- Add focused `below_margin` filter cases for both `List` and `ByProduct`, with ceiling outage and cost outage. Assert nil error, candidate rows/groups are returned rather than silently excluded, every child has null cost/margin, and the repository cursor/as-of behavior remains valid. These tests pin the `scan`/`scanGroups` fallback and the “no reader-unavailable 503 anywhere” rule.
- Keep/add a compact hard-error guard proving installation validation and pricing-policy read failures still return errors before optional enrichment.

## 2. Implement the smallest service change

- `List` (`read_service.go:107-139`): retain installation and policy validation unchanged. At `:122-126`, convert only the ceiling-reader error into an internal `factsUnavailable` state and a nil ceiling; do not return an error. Pass that state through fast enrichment or `scan`, then return the repository page with UTC `AsOf` as today.
- `ByProduct` (`:142-174`): make the same ceiling handling at `:157-163`; pass degradation state through `enrichGroups`/`scanGroups`, then preserve `finalizeGroups`, group ordering/counts, cursor, and `AsOf`.
- `Get` (`:265-295`): at `:278-287`, tolerate ceiling unavailability, skip cost lookup in that state, and enrich the found row as unknown. If cost alone is unavailable, keep the known ICMS rows but leave their margin results null; if ceilings are unavailable, leave `ICMSWorstCaseByUF` unknown rather than fabricating data. Continue reading/returning timeline normally. Repository, not-found, policy, and timeline errors remain hard.
- `enrich` (`:336-376`): change the private helper contract to accept/return the degradation state. If ceilings already failed, skip `GetCostFactsByIDs`; if the cost call fails, set degraded instead of returning `read listing costs`. In either degraded case explicitly clear `Cost`, `BelowMarginWorstCase`, and the derived ICMS detail on every item, then recompute `PendingIssue` from the honest unknown state. Available-reader behavior stays byte-for-byte equivalent.
- `enrichGroups` (`:223-237`): propagate the boolean degradation result from flattened child enrichment, copy children back, and otherwise preserve existing group finalization.
- `scan` (`:298-333`) and `scanGroups` (`:177-220`): carry the degradation state. While facts are available, retain the current dependent-filter scan exactly. If degradation is known on entry or discovered from a cost batch, stop evaluating `matchesDependentFilter` (unknown is neither true nor false) and return a pass-through candidate page/group page from the original request cursor, with null cost/margin and normal cursor/limit semantics. If failure appears after earlier scan pages, restart once from the original cursor so the response does not mix filtered survivors with omitted unknown candidates.

Filter decision: pass-through degraded rows is the honest 200 response. Returning empty would silently treat unknown as “not below margin”; returning 503 would make an optional Oracle fact a read-spine dependency. The nullable row fields visibly communicate that the exception could not be evaluated. There is no list/detail `margin_unknown` counter to populate; that counter remains Summary-only.

## 3. Regression and validation

- Run the listings unit package first, then `GOCACHE=.gocache go test ./...`/build-vet per the harness.
- Keep M01-C04 green: list ordering, cursor walk, search/filter behavior, and nullable response shape.
- Keep M01-C05 green: by-product grouping, null-product-last behavior, counts, state, and cursor.
- Keep M01-C07 green: unknown cost/margin remains JSON null, never zero/false; Summary counters retain existing semantics.
- Keep M01-C09 green: fast path remains one repository page plus at most one cost batch; Summary remains one SQL aggregate query.

No OpenAPI, SDK, migration, adapter, transport, repository, or other-module changes are needed: all affected fields are already nullable and the HTTP/JSON shape is unchanged.
