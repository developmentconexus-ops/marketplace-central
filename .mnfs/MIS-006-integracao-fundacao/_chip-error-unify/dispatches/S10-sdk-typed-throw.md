# S10 — SDK: one typed throw path (VC-5)

Blocks 1 and 2 are in `_common-blocks.md` in this same directory. **Read that file first,
in full.** Blocks 3, 4 and 5 are Go-specific and do NOT apply; your work and validation
are below.

The OpenAPI spec was re-specced in the immediately preceding slice and is now the truth:
`contracts/api/marketplace-central.openapi.yaml`. Read the error schemas there
(`ErrorResponse`, `CatalogPageErrorResponse`, `ErpImportError`, `ErpImportConflict`,
`ActiveSourceError`, `SellableAssortmentError`, `MutationErrorResponse`,
`ConnectorErrorResponse`) — their `code` enums are where your unions come from. Do not
invent codes and do not copy them from stale FE strings.

No other slice is running.

## What is broken today

`packages/sdk-runtime/src/index.ts`:
- `:1720 ErrorResponse` and `:1728 MarketplaceCentralClientError` are **interfaces**, so
  `instanceof` cannot work and a caught error is a bare object literal, not an `Error`.
- Six throw sites blind-cast and throw an object literal:
  `:1751` (getJson), `:1826` (putJson), `:1835` (deleteJson), `:1847` (postJson),
  `:1860` (postVoid), `:1877` (postMultipart, which also attaches `body: data`).
  Every one does `(data as ErrorResponse).error` with no runtime check whatsoever. If the
  server ever answers a shape the cast lies about, `.error` is `undefined` and every
  consumer downstream reads `undefined.code` — a blank screen instead of an error state.
- `packages/sdk-runtime/src/market.ts:149` and `:162` throw a SECOND error type
  (`MarketPriceIntelApiError`). VC-5 says one typed throw path; two is the defect.
- `:1899 getCatalogAssortmentCounts: () => getJson(...)` takes no arguments, so a caller
  cannot pass `erp_source` even though the endpoint accepts it and the spec now declares
  it.

## The work

### 1. A real error class

`MarketplaceCentralClientError` becomes a `class ... extends Error`. It must carry
`status: number`, `code: string`, `message: string`, and `details: Record<string, unknown>`
(always an object, never `undefined` — the backend now always sends it and the spec now
requires it). Set `this.name`, and pass the envelope's human message to `super(...)` so
`e.message` is the human string and a bare `console.error(e)` is readable.

Keep the existing `error` property too, shaped as before (`{ code, message, details }`),
so consumers that already read `e.error.code` do not break. Say in your report that you
checked how many call sites read `e.error` versus `e.code` — grep for it, do not guess.

`ListingsRefreshConflictError` (`:1733`) currently `extends` the interface. Rework it to
whatever the class version requires while keeping its narrowing power (status 409, code
`refresh_in_progress`, `details.operation_run_id: string`). If a type-level narrowing helper
serves it better than a subclass, use that and say why.

### 2. Code unions from the spec

Per-domain unions (e.g. `CatalogErrorCode`, `ErpImportErrorCode`, `MutationErrorCode`,
`ActiveSourceErrorCode`, `ConnectorErrorCode`) plus a global `ApiErrorCode` that is the
union of them. Source every member from the spec's enums.

A union type is only worth building if something CHECKS it. Its job here is that
`hasCode(e, "invalid_erp_soruce")` — note the typo — must fail to compile. Write a test
that proves the union rejects a bogus code (see tests below); a union nothing validates
against is decoration.

### 3. `isApiError` and `hasCode`, both exported, both narrowing

```ts
export function isApiError(value: unknown): value is MarketplaceCentralClientError
export function hasCode<C extends ApiErrorCode>(value: unknown, code: C): value is MarketplaceCentralClientError & { code: C; details: Record<string, unknown> }
```

Exact signatures are yours to finalise, but after `hasCode(e, "invalid_erp_source")` the
expression `e.details.allowed_range` MUST be readable and MUST type as something a
consumer can assign to `const allowed: string` — the day-1 golden test already asserts
exactly this and it is not yours to weaken.

`isApiError` must be a real runtime check (an `instanceof` test), not a duck-type on the
presence of `.code`.

### 4. One parse path with runtime validation — the core of VC-5

Every throw site funnels into ONE function that takes the raw parsed body plus the status
and returns a `MarketplaceCentralClientError`. All six sites in `index.ts` and both in
`market.ts` call it. Delete `MarketPriceIntelApiError`'s separate throw shape (check its
consumers first and report how many there were).

That function VALIDATES rather than casts:
- Envelope present and well-formed (`error` is an object with a string `code`) → use it,
  with `details` defaulting to `{}` when absent.
- **Anything else** — a bare string body, an HTML error page from a proxy, `null`, a
  legacy flat `{"error":"some_code"}`, a 502 with no body — becomes code
  `internal_error` with the ORIGINAL payload preserved at `details.raw`.
  **Never `undefined`, never a silently dropped body.** This is ADR-17's rule in the SDK:
  an unknown operational fact does not become a default that looks like knowledge. A
  consumer that gets `internal_error` can still see what actually arrived.
- `postMultipart`'s extra `body: data` is superseded by `details.raw`; do not keep two
  ways to reach the payload.

Zero blind `as ErrorResponse` may remain on the error path. That is a graded criterion —
you will grep for it in validation.

### 5. `erp_source` on the counts call

`getCatalogAssortmentCounts` takes an optional options object with `erp_source?: "xlsx" |
"catalogo_cliente"`. Build the query with the `catalogQuery` helper already at `:1756`
rather than concatenating a string. Calling it with no arguments must stay byte-identical
to today (no `?` suffix), because existing callers and `activeSource.test.ts:165` do that.

## Tests

`src/errorContract.golden.test.ts` already exists and is currently RED. **It is the
acceptance test for this slice and you must NOT edit it.** If you find yourself needing to
change it, stop and report — that means the contract moved, which is not your call.

Add `src/errorContract.test.ts` (a separate file from the golden) covering what the golden
does not:
1. A malformed body (e.g. `"<html>502 Bad Gateway</html>"`) produces `internal_error` with
   the original text at `details.raw`, AND `e instanceof MarketplaceCentralClientError`.
   Assert `details.raw` explicitly — a test that only checks the code would pass on an
   implementation that threw the payload away, which is the exact failure being prevented.
2. A legacy FLAT body `{"error":"invalid_q"}` also lands as `internal_error` with the flat
   body preserved at `details.raw` — it is NOT silently reinterpreted as code
   `invalid_q`. Guessing at an old shape would resurrect the ambiguity this chip is
   deleting.
3. An envelope with `details` absent yields `details === {}` (an object), not `undefined`.
4. `isApiError` returns false for a plain object that merely LOOKS like one
   (`{ status: 400, code: "x" }`) — proving it is an `instanceof` check and not a
   duck-type.
5. A type-level test that a bogus code is rejected by the union. Since vitest's esbuild
   transform ERASES types at runtime, a runtime assertion cannot prove this. Use
   `@ts-expect-error` on the bogus-code line: `tsc` then fails if the line does NOT
   error, which makes the tsc lane the instrument. Say explicitly in your report that this
   assertion is discharged by `tsc`, not by vitest.

## Validation

1. `cd <checkout>/packages/sdk-runtime && npx --no-install vitest run`
   ANSI-strip before grepping (`sed 's/\x1b\[[0-9;]*m//g'`) — vitest's colour codes sit
   BEFORE leading whitespace and defeat `^`-anchored patterns, returning an empty filter
   that reads as clean. Report COUNTS, never the tail.
   `errorContract.golden.test.ts` must now be GREEN. If it is not, the slice is not done.
2. `cd <checkout>/apps/web && npx --no-install tsc --noEmit` — this is the FE lane, and
   the `cd` is PART of it. The ceiling is **12 pre-existing errors, compared as a LIST,
   not as a count**. The exact list is at
   `.mnfs/MIS-006-integracao-fundacao/_chip-error-unify/evidence/EU-baseline-tsc.txt` —
   read it and diff against it. New errors are yours; the FE migration is a LATER slice, so
   errors caused by the FE not yet using the new API are EXPECTED. Report them as a
   precise list so the next slice inherits an accurate work item, and do NOT fix FE files:
   they are outside your write_set.
3. `grep -n 'as ErrorResponse' packages/sdk-runtime/src/*.ts` → must be 0 on the error
   path. Print the output with its byte and line count.
4. Prove the new tests are real: capture a NAMED red for each before it passes.

## write_set

- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/market.ts`
- `packages/sdk-runtime/src/errorContract.test.ts` (new)

Other `packages/sdk-runtime/src/*.ts` files may need a type import updated if they
reference the changed types — that is allowed, but report each one and why. Do NOT touch
`apps/web` or any `packages/feature-*`: the FE migration is its own slice.

**open_questions**: none. If you find one, stop and report rather than deciding.
