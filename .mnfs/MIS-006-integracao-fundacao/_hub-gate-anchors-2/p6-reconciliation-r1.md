# P6 dual gate — reconciliation, round 1 (CHIP-ANCHORS-2)

**Joint verdict: REFUTED.** Both model-side reviewers refuted independently, from different
evidence, on the same frozen input. The merge at `dbdcdfb1` therefore stands as **merged but
NOT gate-approved**; the correctives below land on `main`.

| | Reviewer | Verdict | Blocking findings |
|---|---|---|---|
| GPT side | `gpt-5.6-sol` / medium, OS-process, read-only sandbox | REFUTED | G1, G2 |
| Claude side | Opus via `harness:gate-reviewer`, physically read-only | REFUTED | B-01 |

Frozen input for both: `p6-input-r1.patch`,
sha256 `0762d05f92aa2c620354f1ea4dcb5c2be99814e60c7b762dcb003da4e0b4dd32` (`git diff e98d8193
dbdcdfb1 -- apps contracts packages`, 20 files, +1220/-124). Both reviewers received the same
brief, ran concurrently, and were blind to each other. Verbatim verdicts:
`p6-sol-gate-r1.md`, `p6-opus-gate-r1.md`.

Why this gate was run at all: the first marker on this merge (`ec70d80c`) declared honestly that
no GPT-side verdict existed, because the chip's dispatch ledger showed GPT as the IMPLEMENTER,
not as a gate. The operator ruled that an honest partial marker is not a substitute for the gate
the harness prescribes, and required the real dual gate — Opus and Sol medium, adversarial, and
specific about what was implemented and what it impacts. That ruling is correct and this file
supersedes the earlier coverage statement.

## What the two sides did NOT share

This is the value the dual gate produced, and it is worth recording precisely.

**Opus alone found B-01**, the only defect that makes the system state something false. Sol
graded the same criterion (C11a) PASS. So did the hub's own pre-merge read. Three independent
readers passed a classifier that decides "does the ERP side have this anchor?" by looking at the
wrong field — because all three checked that the reclassification was placed at the sites the
ruling named, and none checked what the value it classifies actually IS.

**Sol alone escalated G1 to blocking** and built the four-link reachability chain for it
(`IsValidCodprod` accepts leading zeros → import stores the raw string → linking canonicalises
via `ParseInt` → the join compares `'101' = '00101'`). Opus found the same locus as B-04 but
downgraded it to non-blocking for lack of live data. The hub had recorded the same observation
at merge time as "an assumption nobody wrote down" — weaker than either.

## Disagreements, resolved by the hub

**G1 / B-04 — leading-zero CODPROD undercounts `vinculados`. Hub rules BLOCKING** (with Sol,
against Opus). Opus downgraded on the ground that it cannot show today's ERP data contains a
padded codprod. But the value is user-authored: the xlsx upload is a spreadsheet a human fills
in, `IsValidCodprod` explicitly accepts leading zeros
(`internal_read/domain/seller_sku.go:20-34` — digits only, at least one non-zero), and the
failure mode is a silent undercount of an operator-facing number with no error anywhere. Absence
of the shape in today's data is not evidence the path is closed. Verified by the hub: all four
links hold.

**G2 / B-05 — malformed `{id}` answers 500. Hub rules NON-BLOCKING** (with Opus, against Sol),
and folds it into the corrective anyway. Opus's grading is the more precise one: the sibling
route `GET /erp/imports/{id}` has identical behaviour and predates this chip, and the OpenAPI
already declares a 500 for the path. Verified by the hub: `grep 'PathValue("id")'` returns
exactly the two sites, `type ImportID string` has no validator, and no invalid-uuid guard exists
anywhere in `internal/`. This discharges the chip's F-8, which was filed as "not verified by
anyone" — it is now verified, and inherited rather than introduced.

## Criteria the reviewers could not execute — discharged by the hub

Both reviewers correctly refused to grade execution-shaped criteria from pasted output. The hub
ran them at the merged tip `dbdcdfb1`.

| Criterion | Reviewer grade | Hub result |
|---|---|---|
| C5(a) determinism | NOT-PROVEN (both) | **PASS** — the three order guards at `-count=10`, `ok ... 1.474s` |
| C5(b) must-fail | NOT-PROVEN (both) | **PASS** — `SortStableFunc`→`SortFunc` at `generation_service.go:861`, guard FAILS on every one of 10 runs; keys identical, only the display orientation flips (`100MM ≡ 10cm` where `10cm ≡ 100MM` is wanted), which is the correct failure signature. Reverted; `git diff HEAD` on the file is empty |
| C9 same-commit | NOT-PROVEN (Opus) | **PASS by git object** — `9c030154` carries `generation_service.go` + `marketplace-central.openapi.yaml` + `sdk-runtime/src/index.ts`; `37d6b7cc` carries the Go files + the yaml + 3 SDK files |
| C11(d) C4 at final tip | NOT-PROVEN (Opus) | **PASS** — the D-121 policy suite re-run at `dbdcdfb1`, 22/22 |

**C11(d) passing is not the same as D-121 being safe.** Among the 22 green tests is
`TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`, which is the test B-01 shows
encodes the defect: its fixture is `&internalreaddomain.ProductCandidate{}`, a candidate with no
canonical id, which `generation_service.go:279` filters out before any candidate exists. A green
suite over a production-impossible fixture proves the code does what the test says, not what the
system needs.

## Criteria the hub REGRADES

- **C11(a) — FAIL.** Both reviewers graded the placement; Opus graded the content and found it
  defective. The reclassification puts the right shape at the sites A2-R1 named, and then
  computes the `side` from `product.ReferenceCode`. The emitted claim is false. This is the R-24
  shape one more time: total in the wording, partial in the code.
- **C2 — PASS except `seller_sku`, which is NOT-PROVEN.** The four-case table is verified for
  `title` and `ean`. The `seller_sku` row is verified against a fixture that cannot occur.

## Correctives, on `main`

1. **CORR-1 (blocking, B-01)** — `identityAnchorValues` case `"seller_sku"` derives the ERP-side
   value from the canonical CODPROD (`product.InternalProductID` / `ProductID`), not from
   `product.ReferenceCode`. The pinning test must carry a candidate WITH a canonical id, so the
   assertion is reachable. Note the alternative the hub rejects: keeping `ReferenceCode` and
   rewriting the details that say "CODPROD" — that would reinstate `refforn` as the ERP side of
   a cross-side comparison in the same mission that removed it as F-01.
2. **CORR-2 (blocking, G1/B-04)** — the `vinculados` join stops comparing a cast integer to a raw
   string. Regression fixture: `erp_import_products.codprod = '00101'`, resolved link
   `internal_product_id = 101`, expected IN `vinculados`.
3. **CORR-3 (non-blocking, G2/B-05)** — validate `{id}` before the query on BOTH `{id}` routes,
   with a test for the malformed value.
4. **CORR-4 (non-blocking, B-03)** — narrow the `marketplace_capability.go` comment. It asserts
   `refforn` "answers no for every provider present and future", and
   `market/domain/identity_resolver.go:90-92` appends a `refforn` anchor from a marketplace
   `MODEL` attribute today. R-25: the false half is deleted or narrowed, never annotated.
5. **CORR-5 (non-blocking, B-06)** — the EVIDENCE sentence claiming the default branch IS
   A2-R1's carve-out. The default fires on both-present, which is strictly wider than
   both-present-and-different.
6. **CORR-6 (non-blocking, B-07 / B-09)** — `jsonb_typeof(...) = 'array'` guard on the queue CTE;
   the unconditional `*comparison.product` deref at `generation_service.go:490`.

## Carried out of this gate, NOT correctives

- **B-02 → contract constraint for CHIP-VINC-NEUTRO.** `QueueRow.tsx:159` enumerates directions
  by string literal, is type-correct, and therefore invisible to `tsc`. A row whose reasons are
  all `INCOMPARABLE` renders ZERO motivo chips, violating the ADR-17 invariant the same file
  documents at `:154-156`, and the glyph/class maps render the literal `undefined`. The wave-2
  pack must name this line, not only the two `Record<Direction, …>` maps the compiler reports:
  **fixing what `tsc` shows leaves this defect silent.** The D-B mechanism made forgetting the
  new case impossible only where the type system reaches.
- **B-08 → its own chip, pre-existing.** `RouteClassMux.RegisterRouteClass` stores the class
  under `"/erp/imports"` while `Handle` receives `"POST /erp/imports"`, so the map never hits and
  `routeClassDeadline("")` falls to `interactiveDeadline`. Verified by the hub at
  `platform/httpx/route_deadline.go:46,50` against
  `erp_import/transport/http_handler.go:39,41,55`. Consequence: the xlsx upload route declares a
  120s batch deadline and receives 15s. Not this chip's write-set; it is a defect on the exact
  path this mission exists to serve.
- **G4 — no index serves `(tenant_id, state, internal_product_id)`.** Needs a production-scale
  `EXPLAIN (ANALYZE, BUFFERS)`; neither reviewer could measure it and neither claimed to.
- **The live-drive obligation is still open**, and Opus named it: the contract assigns live to
  the hub, the operator waived it to CHIP-IMPORT-CHAIN's P7, and no `LIVE-VERIFIED:` exists yet.

## Lesson for the doctrine

An implementer is not a gate over itself — that part the first marker got right. What it got
wrong is the remedy: the hub substituted its own second read and declared the coverage honestly,
and an honest declaration of a missing gate is still a missing gate. **Declaring a gap is not
discharging it.** The gate the operator required found, in one round, a false operator-facing
claim that three prior readers had passed.
