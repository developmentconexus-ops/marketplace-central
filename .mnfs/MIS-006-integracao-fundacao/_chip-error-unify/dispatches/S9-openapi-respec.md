# S9 — OpenAPI re-spec (VC-4)

Blocks 1 and 2 are in `_common-blocks.md` in this same directory. **Read that file first,
in full.** Blocks 3, 4 and 5 are Go-specific and do NOT apply to this slice; your work and
your validation are specified below instead.

The spec is `contracts/api/marketplace-central.openapi.yaml` (8391 lines). Note the path:
it is under `contracts/api/`, not `contracts/`.

No other slice is running. Nothing else in the tree is mid-edit.

## What is true today (measured by the chip — verify each line before you act on it)

The backend now emits ONE shape on every error, without exception:
`{"error":{"code":"...","message":"...","details":{...}}}`, with `details` ALWAYS present
(as `{}` when empty). The spec has not caught up, and it describes FIVE flat error
schemas, not the one the brief named:

| line | schema | refs | shape today |
|------|--------|------|-------------|
| 3984 | `CatalogPageErrorResponse` | 12 | flat: `error` is a **string** enum, plus a sibling `allowed_range` |
| 8271 | `ErpImportError` | 11 | flat: `error` string enum, plus `detail`, plus `column` |
| 8283 | `ErpImportConflict` | 2 | flat: `error` string enum, plus `import_id`, plus `protocol` |
| 8331 | `ActiveSourceError` | 4 | flat: `error` string enum, plus `detail` |
| 8366 | `SellableAssortmentError` | 5 | flat: `error` string enum, plus `detail` |

Every one of them is now WRONG about the wire. All five die.

The base to specialise from is `ErrorResponse` (`:7553`, 107 refs). The pattern the mission
names is `MutationErrorResponse` (`:8129`): a standalone object with
`required: [error]`, whose `error` object has `required: [code, message, details]` and a
per-domain `code` enum. Follow THAT pattern, not `ConnectorErrorResponse`'s `allOf` form
(`:7584`) — two idioms already exist in this file and the mission picked one. Say in your
report that you noticed both and which you followed.

## The work

### 1. Replace the five flat schemas

Each becomes a `MutationErrorResponse`-shaped specialisation. The code enums are already
written in the flat versions — **carry them over verbatim, do not re-derive them**. Their
sibling ad hoc fields move into `details` as declared properties:

- `CatalogPageErrorResponse` → `allowed_range` becomes a property of `details`. Keep its
  existing description text, which explains exactly when the field appears; it is accurate
  and still applies.
- `ErpImportError` → `column` into `details`; the flat `detail` field is GONE, because the
  backend now writes that human string as `message`.
- `ErpImportConflict` → `import_id` and `protocol` into `details`; `detail` gone.
- `ActiveSourceError` and `SellableAssortmentError` → `detail` gone (it is `message` now).
  These two have no other ad hoc field.

Keep each schema's NAME and every `$ref` to it working. This is a shape change, not a
rename: renaming would touch 34 reference sites for no gain and is out of scope.

### 2. `details` becomes required on the base

`ErrorResponse` (`:7553`) currently has `required: [code, message]` with `details`
optional. The backend now always emits `details`, as `{}` when empty — that was the entire
point of having one producer. Make `details` required so the spec states the guarantee the
code actually provides. A client that must branch on "is details there or not" is exactly
the defect this chip exists to kill.

### 3. Fix the false prose in `ErrorResponse`'s `code` description

`:7566-7576` says codes follow `MODULE_ENTITY_REASON` and lists SCREAMING_SNAKE examples.
That is now only half true: most of the tree uses snake_case (`invalid_erp_source`,
`deadline_exceeded`, `internal_error`), while `integrations` still uses SCREAMING_SNAKE
(`INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED`). A description that states one convention as THE
convention is false, and false prose in a published contract gets corrected, not decorated.

Rewrite it to state what is actually true: the codebase carries both conventions, the
snake_case form is the prevailing one, and the per-domain schemas below are the
authoritative enumerations. **Do not invent a migration promise or a deprecation notice** —
you do not own that decision, and a contract that promises a future rename is a claim
nobody has agreed to.

### 4. Declare `erp_source` on the counts endpoint

`/catalog/products/counts` (`:429`) accepts `erp_source` at runtime and rejects an unknown
value with 400 — the whole EXEMPLO-IO of this chip is that rejection — but the endpoint
declares NO parameters at all. Add the query parameter, mirroring the sibling declaration
at `:370-380` (same `enum: [xlsx, catalogo_cliente]`, `required: false`). Reuse that
block's wording where it applies rather than writing a second description of the same
parameter.

While you are there: the `"400"` response description at `:443-445` says "the error field
names which one (erp_source)". Under the envelope the field is `error.code`, so make the
sentence true.

## Validation

There is no Go in this slice. Prove the YAML three ways, all `ran`:

1. **It parses.** Use a parser already available in the repo's installed toolchain — check
   `packages/sdk-runtime` and the root `node_modules` for something that reads YAML before
   reaching for anything else. Do NOT install a dependency. If nothing in the tree can
   parse YAML, say so and fall back to the tests in step 3 rather than installing.
2. **The flat shape is gone.** This is VC-4's stated evidence. Grep for an `error` property
   declared as `type: string` inside any error schema and show the result is empty, with
   its byte and line count. Then reconcile POPULATION vs EXTRACTION and print BOTH numbers:
   how many error-ish schemas exist in the file in total, and how many you changed. They
   will not be equal — `ListingReadSyncError` (`:3511`) and `ConnectorErrorResponse`
   (`:7584`) exist too. Read both and state explicitly whether each is an HTTP error body
   this chip governs or something else (one of them is a per-item status field inside a
   200 response, not an error body at all — determine which by reading, and do not change
   what this chip does not govern).
3. **The existing contract tests still pass.** `packages/sdk-runtime` has tests that read
   this YAML directly (`activeSource.test.ts` asserts against `operationId:
   getCatalogAssortmentCounts`, among others). Run:
   `cd <checkout>/packages/sdk-runtime && npx --no-install vitest run`
   Strip ANSI before grepping anything (`sed 's/\x1b\[[0-9;]*m//g'`) — vitest's colour
   codes sit BEFORE the leading whitespace and defeat `^`-anchored patterns, which returns
   an empty filter that reads as clean. Report tests passed/failed as COUNTS, never as the
   tail.
   **`src/errorContract.golden.test.ts` is EXPECTED to fail** — it pins the future SDK
   surface, which a later slice builds. Report it as an expected red, by name, and confirm
   nothing ELSE is red. If something else is red, that is yours.

## write_set

- `contracts/api/marketplace-central.openapi.yaml`

Only this file. The SDK is a later slice and touching it here would split a change that
AGENTS.md requires to land together. If you believe the SDK must change for your edit to
be coherent, STOP and report — do not reach into it.

**open_questions**: none. If you find one, stop and report rather than deciding.
