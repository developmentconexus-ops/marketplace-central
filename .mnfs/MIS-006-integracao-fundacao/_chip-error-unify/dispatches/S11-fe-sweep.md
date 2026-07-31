# S11 — FE: one match pattern, TOTAL sweep (VC-6)

Blocks 1 and 2 are in `_common-blocks.md` in this same directory. **Read that file first, in
full.** Blocks 3-5 are Go-specific and do NOT apply.

No other slice is running. You are the only writer in this checkout.

## Read this first — this slice fixes a regression THIS CHIP CAUSED

This is not a cosmetic migration. Three FE sites match the error shape the backend emitted
BEFORE this chip, and they are broken RIGHT NOW at runtime:

- `ImportChainPanel.tsx:12` — `candidate.error === "import_not_found"`. On the SDK's
  `MarketplaceCentralClientError`, `.error` is the OBJECT `{code,message,details}`. An
  object is never `===` a string, so this is **always false**.
- `IntegracoesPage.tsx:304` — `const { status, error: code } = ...; code === "unknown_erp_source"`.
  Same defect via destructuring. **Always false.**
- `useErpImportUpload.ts:42-44` — reads `maybe?.body`, and the error class has NO `body`
  property. `code` is always `undefined`, the `switch` always falls to `default`, and only
  the status-only fallback below it does anything. The operator loses the `column` name on
  422 and the `protocol` on 409 — an unknown fact silently becoming a default, which is the
  exact ADR-17 violation this chip exists to prevent.

They are masked by 7 FE tests that mock the OLD flat shape (listed below). The tests are
wrong, not the SDK. `IntegracoesPage.test.tsx` even contradicts itself — line 407 mocks the
correct nested shape for `unknown_erp_source` while 293/307/423 mock the flat shape for the
same code on sibling calls.

So: fixing these is REQUIRED for this chip to be correct, not optional polish.

## The SDK surface you migrate onto

From `@marketplace-central/sdk-runtime` (already landed, commit ea5fd9ff):
- `class MarketplaceCentralClientError extends Error` — `status`, `code`, `message`,
  `details` (always an object), plus a legacy `error: {code,message,details}`
- `isApiError(value): value is MarketplaceCentralClientError` — a real `instanceof` check
- `hasCode(value, code)` — narrows AND runtime-validates spec-guaranteed detail fields.
  For `missing_required_column` it guarantees `details.column: string`; for
  `duplicate_file`, `details.import_id` and `details.protocol`; for `invalid_erp_source` /
  `invalid_limit` / `invalid_ids`, `details.allowed_range`.
- `isRefreshInProgressError(value)` — the dedicated 409 guard
- `ApiErrorCode` and the per-domain unions

**The rule (VC-6): zero string-literal error-code comparison anywhere in the FE outside
`hasCode` or a switch over a TYPED union member.** A raw `=== "some_code"` against an
untyped value is what dies. A `switch` over a value typed as `ApiErrorCode` is fine — it is
type-checked, which is the substance.

## The population — measured, classified, and NOT the same as "every hit"

IN SCOPE (HTTP-ERROR, 6 sites):
1. `apps/web/src/pages/AnunciosPage.tsx:136` — `invalid_filter`, currently correct (nested
   read) but uses an inline cast. → `hasCode`.
2. `apps/web/src/pages/ListingsRefreshControl.tsx:47-48` — hand-rolls exactly what
   `isRefreshInProgressError` already does. → call the helper; delete the hand-rolled check.
3. `apps/web/src/pages/mutations/mutationPresentation.ts:18-24` (feeds
   `MutationPreviewModal.tsx:57,96,111,112,141,150,173`) — `preview_stale`,
   `selection_too_large`. Currently correct nested reads. → `isApiError`/`hasCode`.
4. `apps/web/src/pages/integracoes/useErpImportUpload.ts:41-70` — BROKEN, see above.
5. `apps/web/src/pages/integracoes/IntegracoesPage.tsx:302-305` — BROKEN.
6. `apps/web/src/pages/importacoes/ImportChainPanel.tsx:9-13` — BROKEN.

DEAD (1 site): `apps/web/src/pages/vinculos/useVinculosQueue.ts:45-48` matches
`ALREADY_RESOLVED` with `status === 409`. Traced: that string is produced only by
`product_links/application/batch_service.go:274`, as a per-item `cause` inside a **200**
batch response. The single-item routes this hook is wired to (`ApproveCandidate` in
`resolution_service.go`) never emit it. **Verify that trace yourself before acting.** If it
confirms, DELETE the dead branch rather than migrating it — migrating dead code to a typed
guard would dress a falsehood in a type. If your trace CONTRADICTS mine, stop and report;
do not delete on my say-so.

OUT OF SCOPE — do NOT touch these. They are domain data on 2xx bodies, not HTTP errors:
- `packages/feature-inventory/src/StockSeguroPage.tsx:164,166` (`blocking_reason.code`)
- `apps/web/src/pages/mutations/MutationResultSummary.tsx:41-42` and
  `ProtocoloPage.tsx:46` (`item.failure.code`, from `listMutationItems()`)
- `apps/web/src/pages/precos/SolverPanel.tsx:76+` (`PricingSolveResponse`)
- `packages/web-query/src/failureCopy.ts` — a presentation lookup keyed by a plain string
  fed from BOTH domain and HTTP paths. It cannot narrow, and forcing it to would be wrong.
- `apps/landing` — measured, zero error-matching code.

If you believe any OUT-OF-SCOPE item is misclassified, report it. Do not migrate it.

## The second error vocabulary — decide, do not silently keep both

`useErpImportUpload.ts:89` does `throw classifyUploadError(err)`, discarding the real
`MarketplaceCentralClientError` and re-throwing a locally-invented
`ErpImportUploadError {kind, column?, protocol?}`. Consumers of `mutation.error` therefore
see a shape `isApiError` cannot recognise — a second producer.

Keep `ErpImportUploadError` as the hook's UI-facing vocabulary if that serves the page (the
page renders a distinct PT message per kind, which is legitimate presentation). But
`classifyUploadError` must derive it via `isApiError`/`hasCode`, reading `details.column`
and `details.protocol` — NOT via `body`, which does not exist.

Preserve the two behaviours the current code documents in comments, because both are real:
the 504/`deadline_exceeded` case must stay distinguishable from `internal_error`, and a
transport failure with no status must stay `network_error`. Do not flatten either into a
generic branch.

Where a code no longer carries the detail it used to, the value is UNKNOWN — it must not
become `""` or a plausible-looking default. Leave it `undefined` and let the UI render its
unknown state.

## The 7 lying test mocks

These mock the pre-chip flat shape and must be corrected to what the SDK actually throws
(construct a real `MarketplaceCentralClientError`, or mock the nested envelope):
- `apps/web/src/pages/importacoes/ImportChainPanel.test.tsx:159` — `{status:404, error:"import_not_found"}`
- `apps/web/src/pages/importacoes/ImportChainPanel.test.tsx:178` — `{status:500, error:"internal_error"}`
- `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx:293, 307, 423` — `{status:400, error:"unknown_erp_source"}`
- `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx:354` — `{status:409, body:{error:"duplicate_file", protocol:"#003-E"}}`
- `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx:365` — `{status:422, body:{error:"missing_required_column", column:"CUSTO"}}`

`IntegracoesPage.test.tsx:407` already mocks the CORRECT nested shape — use it as the
reference, and note in your report that the file previously disagreed with itself.

**Fixing a mock to match the implementation is normally how tests get faked green, so
handle it in the honest order**: for each of the three broken sites, FIRST correct the mock
to the true SDK shape and RUN it, showing the test goes RED against the CURRENT (broken)
implementation. That red is the proof the site was broken. THEN fix the implementation and
show it green. Capture both, by test name. A mock corrected only after the implementation
is fixed proves nothing.

## Validation

Report COUNTS, ANSI-stripped (`sed 's/\x1b\[[0-9;]*m//g'`) — vitest's colour codes precede
leading whitespace and defeat `^`-anchored greps. Never a tail.

1. `cd <checkout>/apps/web && npx --no-install tsc --noEmit` — ceiling is the 12
   pre-existing errors as a LIST, in
   `.mnfs/MIS-006-integracao-fundacao/_chip-error-unify/evidence/EU-baseline-tsc.txt`.
   Diff against it; zero new. Note that two of the 12 (`MutationPreviewModal.tsx:197`,
   `MutationResultSummary.tsx:22`, both `ErrorStateProps.onRetry`) are in files you touch —
   they must still be there, unchanged, and you must not "fix" them.
2. FE vitest, per package. Run them explicitly and name each lane; `apps/web` at minimum,
   plus any `packages/feature-*` whose files you touched. State the command per lane.
3. `cd <checkout>/packages/sdk-runtime && npx --no-install vitest run` and
   `npx --no-install tsc --noEmit -p tsconfig.test.json` — you should not affect these, and
   the second must stay at its 7-error named baseline. Prove you did not.
4. **The VC-6 census.** After your edits, re-run the sweep and print POPULATION vs
   EXTRACTION as two numbers with the arithmetic. The census must show zero
   string-literal error-code comparison outside `hasCode`/typed-union in the HTTP-ERROR
   population, while correctly still counting the OUT-OF-SCOPE domain-data sites as
   present-and-untouched. Give the exact grep commands so they can be re-run.

## write_set

`apps/web/src/pages/` — the six in-scope files, `useVinculosQueue.ts`, and the two test
files named above. Do NOT touch: `packages/sdk-runtime`, any Go file, the OpenAPI yaml,
`packages/feature-inventory`, `packages/web-query`, `apps/landing`, or the out-of-scope
sites.

**open_questions**: none. If you find one, stop and report rather than deciding.
