VERDICT: REFUTED

| Criterion | Grade | Evidence checked |
|---|---|---|
| A1 | PASS | `generation_service.go:736` — `func identityAnchorValues`; seller-SKU path contains `canonicalProductID(*product)` and `strconv.Itoa(productID)` at `:751-752`; `ReferenceCode` occurs zero times inside the function. |
| A2 | PASS | `generation_service.go:279` — `canonicalProductID(product)` filters candidates; `generation_service_test.go:478` — `TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference`; `:519` — `InternalProductID: canonicalIDPtr(101)`; `:522` — `generateSingle`. `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide` was corrected, and the unreachable `"exact EAN ERP seller SKU empty"` case was deleted. |
| A3 | NOT-PROVEN | `dispatches/a3-mustfail-raw.txt:12` claims `INCOMPARABLE`/`side:"erp"` instead of `UNAVAILABLE`/`side:""`; current `git diff HEAD -- generation_service.go` is empty. I could not independently execute the mutation in the read-only gate. |
| A4 | PASS | `generation_service.go:739-753` compares listing `SellerSKU` with canonical CODPROD. `missingMatchedAnchorReason` at `:646-648` sends both-present-and-different to `UNAVAILABLE`; `classifyProviderIdentityAnchor` at `:724-725` emits no replacement. The anchor remains non-emitted, but no longer because `seller_sku` was compared with `refforn`. |
| A5 | NOT-PROVEN | `query_repository.go:98` — `ltrim(links.internal_product_id::text, '0') = ltrim(products.codprod, '0')`; fixture at `chain_query_repository_integration_test.go:266,273,282` checks `'00101'`, `101`, and `Vinculados == 1`. Real-Postgres execution was not independently reproducible. |
| A6 | NOT-PROVEN | `dispatches/impl-integration-a5-a7-a10.md:216` claims exact old-join result `Vinculados:0`. This execution-shaped must-fail exists only as chip-authored output and could not be independently repeated read-only. |
| A7 | REPORT | `chain_query_repository_integration_test.go:51-53,58-60,80` retains the duplicate-link fixture and `Vinculados == 2` assertion. The supplied record says counts passed 9/9, but the same test failed 4/9 at its separate `QueueReadAt` assertion (`:83-84`). |
| A8 | PASS | `http_handler.go:87-94` — `readImportID`; calls at `:97` and `:114`; `http_handler_test.go:365-390` checks both `not-a-uuid` routes, exact 400 body, and untouched `getID`/`chainID`. OpenAPI contains both 400 responses and `invalid_import_id`; `erpImport.ts:80-90` mirrors the union. |
| A9 | PASS | `marketplace_capability.go:26-42` — exact checked text begins `knownIdentityAnchors is the identity vocabulary THIS file governs` and explicitly identifies the wired `market/domain/identity_resolver.go` exception instead of asserting a universal. |
| A10 | NOT-PROVEN | `query_repository.go:110-116` uses `CASE WHEN jsonb_typeof(...) = 'array'`; fixtures at `chain_query_repository_integration_test.go:311-312,341-347` cover object and scalar. The claimed real-Postgres PASS and SQLSTATE 22023 must-fail were not independently executable. |
| A11 | NOT-PROVEN | `generation_service.go:489-495` contains the nil guard. The only production call at `:303` passes `&product`, where `product := skuMatches.Products[0]`; `generation_service_test.go:542-564` reaches nil only by calling the private scorer directly. Production reachability is therefore not proven. |
| A12 | NOT-PROVEN | Diff inspection confirms 11 files, zero `apps/web`, zero migrations, and zero `platform/httpx`. The two governance artifacts contain the same 53 `error_code=` findings, but the required build, vet, full test, order-guard, policy-suite, SDK, and real-Postgres executions could not be independently rerun. |
| A13 | PASS | `EVIDENCE.md:13-14,315-316` points to `_hub-gate-anchors-2/p6-reconciliation-r1.md` and names `B-01, G1/B-04, G2/B-05, B-03, B-07, B-09` without recopying their bodies. |

## TestCase3EANAloneYieldsMediaConfirm

The new assertion encodes the code that was written, not behavior that can be established as correct.

The symmetry precedent is real: `TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason` sends an EAN present-and-different on both sides through the same `missingMatchedAnchorReason` default branch and expects `UNAVAILABLE` with an empty side. After CORR-1, `TestCase3EANAloneYieldsMediaConfirm` reaches that pre-existing branch because `"SKU-NO-MATCH"` and canonical CODPROD `"300"` are both present and different.

That proves the assertion accurately describes the resulting path. It does not make `UNAVAILABLE` semantically correct: that direction means the provider does not supply the anchor, while this listing does supply it. CORR-1 did not create the bucket error—it made the existing branch common by replacing `refforn` with the always-present canonical CODPROD. Per the brief, the direction decision is out of scope and this question is REPORT, not PASS.

## Findings

### Blocking

- `dispatches/impl-integration-a5-a7-a10.md:179-249` — A5, A6, and A10 depend on real-Postgres PASS/must-fail executions supplied by the chip itself. The gate could inspect the SQL and fixtures but could not independently execute them; the read-only environment denied creation of Go’s work directory. Correct would be independently reproducible cold execution, or an immutable external CI artifact establishing the exact tests, database target, no skips, and mutation signatures.

- `dispatches/ladder-l0-l1-raw.txt:1` — A12’s execution obligations are not independently proven. The supplied transcript cannot substitute for a cold gate run, particularly where the brief explicitly says supporting evidence is a claim rather than proof. Correct would be independently verified build, vet, full suite, both `-count=10` guards, D-121 suite, SDK tests, and governance comparison.

### Non-blocking

- `chain_query_repository_integration_test.go:74-84` — the DISTINCT guard shares a test with a host-clock/database-clock bracket that failed 4 of 9 claimed runs. The count assertion appears intact, but the test as a whole is not a reliable green guard. Correct would separate the count invariant from timing or validate time using a database-controlled bound that does not compare unsynchronized clocks.

- `generation_service_test.go:1554-1569` — the assertion correctly captures the newly reached default branch but labels both-present-and-different as `UNAVAILABLE`. Correct behavior requires the operator’s pending A2-R1 direction decision; this chip correctly discloses that it made a pre-existing defect frequent rather than creating it.

- `generation_service.go:490-495` — the nil guard prevents a panic, but its only test calls an unreachable production shape and still constructs a concordant candidate with no internal product ID. Correct would be to retain the honest NOT-PROVEN grade until a real production path can supply nil, or redesign the function contract if such a path is introduced.

## What I could not verify, and why

- I could not run Go tests, builds, or vet: the read-only environment denied creation of Go’s temporary work directory.
- I could not mutate `generation_service.go` or `query_repository.go`, so A3 and A6 could not be independently must-failed.
- I could not provision or write to a PostgreSQL test database, so A5, A7, and A10 execution claims remain unconfirmed.
- I could not independently reproduce the SDK test run.
- I did not inspect the other concurrent P6 reviewer’s output, as the brief forbids it.
- L2/live-drive behavior, production-scale query performance, B-02, B-08, G4, `apps/web` TypeScript failures, and the final A2-R1 direction are explicitly outside this gate’s scope.
