VERDICT: CONFIRMED

## (a) False universal — PASS

[generation_service_test.go:418](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/apps/server_core/internal/modules/product_links/application/generation_service_test.go:418):

> "The ERP side of seller_sku is the CODPROD, and findProducts drops any candidate without one — so with a product present, side=erp cannot arise for seller_sku. It arises on the nil-product (unresolved) path instead"

This is true within its stated production scope:

- `findProducts` rejects candidates without a positive canonical ID at `generation_service.go:277-282`.
- A surviving product therefore supplies the `seller_sku` ERP value through its canonical CODPROD.
- The production unresolved path explicitly passes `nil` at `generation_service.go:215-216`.
- `missingMatchedAnchorReason` assigns listing-empty to `SideProvider` and nil/missing product to `SideERP` at `generation_service.go:645-650`.

The prior universal was deleted and narrowed under R-25, not preserved through annotation.

## (b) Restored coverage — PASS

[generation_service_test.go:616](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/apps/server_core/internal/modules/product_links/application/generation_service_test.go:616):

> "TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide drives the whole generator down the unresolved path"

True. `generateSingle` invokes the public `GenerateLinkCandidates` method at `generation_service_test.go:1308-1324`. Its empty matcher results make SKU, EAN, and title resolution empty, reaching the production unresolved call at `generation_service.go:215-217`.

The test asserts all material properties:

- `NO_CANDIDATE` at `generation_service_test.go:639-640`.
- Presence of the `seller_sku` reason at `:642-644`.
- Exact `INCOMPARABLE`, `SideERP`, and detail at `:646-648`.

[ generation_service_test.go:624](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/apps/server_core/internal/modules/product_links/application/generation_service_test.go:624):

> "classifyProviderIdentityAnchor emits nothing when the product is nil and the listing value is present."

True. Mercado Livre supplies `seller_sku`; with `comparison.product == nil` and a non-empty listing value, `classifyProviderIdentityAnchor` returns `("", "", "", false)` at `generation_service.go:723-727`. Consequently the declared-anchor loop continues at `:686-688`, leaving the seeded reason unchanged. The test is not passing because a promotion happens to produce an equivalent value—the promotion does not fire.

A separate test was justified. The existing table runner passes `nil` declarations and invokes the internal scorer at `generation_service_test.go:455-457`; it therefore cannot demonstrate either production reachability or survival through the declared-anchor pass.

## (c) Six EVIDENCE corrections — PASS

Each correction in [EVIDENCE.md:578](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/.mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/EVIDENCE.md:578) is accurate:

1. **A2:** "`side=erp` ficou INALCANÇÁVEL" was false. The production unresolved path at `generation_service.go:216` makes it reachable.

2. **A3:** "um subteste discrimina" was false. `a3-mustfail-raw.txt:18-19` records two failing subtests and only the `has_a_refforn` case passing.

3. **A7:** "root-caused a skew de ~600ms" was overclaimed. A constant 600 ms skew does not explain intermittent 4/9 failures reported as roughly 1 ms. Retaining REPORT while grading the diagnosis unproven is honest.

4. **A12:** The old `ladder-l0-l1-raw.txt` contains neither build/vet commands nor the `-count=10` guard and includes the corrupted string `[no test fi107`. The correction accurately names `ladder-l0-l1-raw-r2.txt` as the replacement claim source.

5. **A11:** "degrada como as duas irmãs" was false. With a nil product, `buildConcordantCandidate` seeds two `FOR` reasons, cannot produce a hard negative because the internal product name is empty, and therefore assigns `95 / ALTA / ACCEPT`.

6. **Ledger:** The slice reviewer's "verified by hand-trace" statement at `review-adversarial-r1.md:19` was wrong: the raw must-fail shows `listing_has_no_seller_sku` is discriminating. Treating that CONFIRMED as weaker rather than corroborating evidence is accurate.

## Other delta hunks

**Queued identity — PASS.** At [query_repository.go:118](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:118), both CTEs retain:

> `SELECT DISTINCT products.codprod AS codprod`

while `queued_products` now uses:

> `ON ltrim(products.codprod, '0') = ltrim(pending.codprod, '0')`

This matches the canonical membership predicate used by `resolved_products` while preserving the raw key counted by `importados`. The fixture exercises both padding directions and asserts `Importados=3`, `Vinculados=1`, `Enfileirados=2`.

**Nil-degradation comment — PASS.** [generation_service.go:490](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0/apps/server_core/internal/modules/product_links/application/generation_service.go:490):

> "function still emits seller_sku FOR and ean FOR at 95 / ALTA / ACCEPT … The row carries a null internal_product_id and autoApprovals still reads that ACCEPT as auto-approvable."

Every part follows from the code at `:497-517`, `newCandidate`'s canonical-ID guard, `detectHardNegative` returning false for an empty internal name, and `autoApprovals` selecting every `ACCEPT`.

## Findings

None.

## What I could not verify, and why

These belong to the separate execution seat and do not alter this reading verdict:

- Execution of Go tests, integration tests, must-fail mutations, build/vet, policy suites, or repeated flake runs.
- PostgreSQL provisioning and runtime confirmation of the SQL fixtures.
- HEAD/commit identities, commit atomicity, git history, restored empty diffs, file lists, or patch hashes.
- The real-database measurements used to adjudicate the rejected R4 overcount framing.
- Reproduction of the A7 timing flake or determination of its actual root cause.
