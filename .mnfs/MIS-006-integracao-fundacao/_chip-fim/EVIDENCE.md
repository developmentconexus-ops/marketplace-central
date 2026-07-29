# CHIP-FIM — evidence pack

- BASE-SHA: `9c82804479bd01084c1ee349954694084990e804`
- Branch: `worktree-chip-fim`
- Contract: `CORTE-YAGNI.md` ("Ordem até o fim") + `_hub-live-drive/DRIVE-vinculos-2026-07-28.md` (F-1)

Raw lane output is in `evidence/`. Every must-fail below is a versioned pair: mutation
applied → red run captured → mutation removed → residue count 0 → green run captured.

## Item → files → test that sustains it

| # | Defect | Files | Test that sustains it | Must-fail |
|---|--------|-------|----------------------|-----------|
| 1 | B-08: batch route got the 15s interactive deadline | `platform/httpx/route_deadline.go`, `route_deadline_test.go`, `erp_import/transport/http_handler_test.go` | `TestRouteClassMuxAppliesDeclaredClassToMethodQualifiedPattern`, `TestHandlerPostImportRunsUnderTheDeclaredBatchDeadline` | `mustfail-01-{red,green}.txt` |
| 2 | F-1: product in both anchors emitted twice | `product_links/application/generation_service.go`, `generation_service_test.go` | `TestCollisionCandidatesEmitOneRowPerProductAcrossBothAnchors` (2 subtests) | `mustfail-02-{red,green}.txt` |
| 3 | Chain panel stuck at `Carregando…` on 5xx | none — already fixed at BASE-SHA | `ImportChainPanel.test.tsx` "renders a generic error for server failures without a chain" | `mustfail-03-red.txt` |
| 4 | STALE candidate outside the cap | `product_links/adapters/postgres/link_candidate_repo.go`, `link_candidate_repo_integration_test.go` (new) | `TestReplaceLinkCandidatesClearsPendingCandidatesLeftOutOfTheRun` | `mustfail-04-{red,green}.txt` |
| 5 | ean detail alleged absence when the value exists | `product_links/application/generation_service.go`, `generation_service_test.go` | `TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason`, `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide` | `mustfail-05-{red,green}.txt` |

## Lanes

| Lane | Command | EXIT |
|------|---------|------|
| Go build + vet | `cd apps/server_core && go build ./... && go vet ./...` | 0 / 0 |
| Go test | `go test ./internal/modules/product_links/... ./internal/modules/erp_import/... ./internal/platform/httpx/... -count=1 -v` | 0 — 205 `--- PASS`, 14 packages `ok` |
| Integration | `npm run harness:integration` | 0 — `status=passed` |
| FE vitest | `cd apps/web && npx --no-install vitest run src/pages/importacoes src/pages/vinculos` | 0 — 10 files, 61 tests passed |
| FE tsc | `npx --no-install tsc --noEmit` | 2 — **0 errors in `src/pages/importacoes` or `src/pages/vinculos`** |

The tsc errors are in `anuncios`, `mutations` and `produto`. None can be mine: this chip
changed no FE file at all (`git status --porcelain apps/web` is empty), so every one of them
predates BASE-SHA.

## Item 1 — which side of the keying was wrong

`RegisterRouteClass` is called at 14 sites and **every one of them passes a bare path**
(`composition/root.go:268` registers 8 batch routes, plus `catalog`, `erp_import:39`,
`integrations:31`, `listings:27`/`:139`, `tenant_config:43`). The transport side registers the
ServeMux pattern with its method. So the miss was not specific to `/erp/imports` — **every
batch route in the application was silently running on the 15s interactive deadline.** Fixing
the mux lookup (exact pattern first, then method-stripped path) repairs all of them at once
and leaves the existing `/batch` test on its exact-match branch.

The red run states the defect in the units the contract used:

```
route_deadline_test.go:112: effective budget = 15s, want the 2m0s batch budget (interactive default is 15s)
http_handler_test.go:140: upload budget = 15s, want the 120s batch budget; 15s is the interactive default
```

## Item 2 — no invented confidence band

The contract forbids deciding alone to promote confidence. Nothing was promoted: of the two
rows the anchors already produce for the shared product, the higher-confidence one is kept
(40 from `applyAmbiguousCorroborationScore` vs 20 from `applyCollisionScore`), ties keep the
first anchor. Both values are ones an existing branch already assigns; no merged score, no new
band. The 40 row is also the only one whose reasons name both anchors.

The red run reproduces EXEMPLO-IO verbatim, in both anchor directions:

```
generation_service_test.go:422: codprod 22467 emitted twice for MLB4735378521: seller_sku/40% and ean/20%
generation_service_test.go:422: codprod 22467 emitted twice for MLB4735378521: seller_sku/20% and ean/40%
```

Both directions are exercised deliberately: a naive "keep the first anchor's row" passes one
direction and fails the other.

## Item 3 — already satisfied at BASE-SHA, and the existing test is not vacuous

`ImportChainPanel.tsx:50-55` already renders `ErrorState` with a manual retry on
`isError || !data`, and `ImportChainPanel.test.tsx` already covers a 500. `git log` on the file
shows `1bdf70b9` / `b97cd9a8` / `b91c7507` — CHIP-ANCHORS-3 landed it.

Reporting "already fixed" is worthless without showing the test can fail, so the error branch
was mutated to fall back to `LoadingState` (the pre-fix behaviour). Exactly the 3 error tests
went red and the other 8 stayed green — clean attribution, no over-broad mutation:

```
× renders a not-found error without a chain
× renders a generic error when a 404 has no not-found body
× renders a generic error for server failures without a chain   (Test timed out in 5000ms)
```

**No FE file was changed by this chip.** After the mutation was reverted, `git status` still
reported the file modified; `git hash-object` and `git rev-parse HEAD:<file>` both returned
`dbcd8416923a579f4627d3cf6ef8350c21c639d8`, proving zero content delta — the flag was an
on-disk LF/CRLF artifact left by `sed`, not an edit.

## Item 4 — no migration needed

"Pending" is expressible with the exact `NOT EXISTS (product_links …)` predicate
`ListLinkCandidates` already reads with, so the sweep needed no schema change and no
migration grant. The per-identity DELETE loop was left untouched; the sweep is additive.

Two boundaries are deliberate:

- a candidate the operator already decided is **not** swept — it is no longer pending;
- an empty identity set clears nothing, because a run that processed nothing asserts nothing
  about the installation.

ADR-17 is why the stale row is deleted rather than kept: the fila renders it as the current
answer, but the generator never asked the matcher about that listing.

## Item 5 — class sweep, reconciled

`sweep-05-detail-class.txt`. Population predicate = every `Detail`-carrying site in
`generation_service.go`; **POPULATION_COUNT=21, EXTRACTION_COUNT=21** — reconciled, no site
silently dropped.

An earlier narrower pattern (`Detail:|Detail = |Detail :?=`) returned 19 and silently missed
the `missingMatchedAnchorReason("seller_sku", …)` sites, where the detail is a positional
argument. That is the exact failure the profile's reconciliation rule exists to catch, and it
is why both counts are printed here.

Of the surviving `"sem …"` details, both are true where they fire: `"sem EAN para corroborar
o CODPROD"` (:559) only when `strings.TrimSpace(snapshot.EAN) == ""`, and `"sem CODPROD para
corroborar o EAN"` (:572) only when the seller_sku is absent. **The class is closed with a
single defective site**, now corrected to mirror the C1 wording ANCHORS-3 used on the
seller_sku side.

`auto_link_policy_test.go` was **not** touched: it pins only the short string (:185, :381),
which stays true. Only `generation_service_test.go:385,:432` pinned the false sentence. Under
the item-5 mutation the auto-link tests stayed green while the two anchor tests went red,
which confirms the split rather than assuming it.

## Field findings for the hub

1. **`npm run harness:pg:up` failed on first invocation with no container created.**
   `status=blocked`, `HPG_READY_TIMEOUT`, EXIT=1, and `docker ps -a --filter name=mpc-pg-session`
   showed no new container at all. The identical command re-run immediately succeeded
   (`container=mpc-pg-session-5ae2f176 port=52798`, ready in 26s). This is a sibling of the
   documented `pg_isready`-lies-on-first-boot race but a **distinct signature** — there the
   container exists and is not ready; here nothing was created. Retry once before debugging.

2. **The integration lane still cannot prove by itself that a test ran.** `summary.txt` holds
   only `target/status/run_id`; a fully-skipped run and a fully-green run are byte-identical.
   What discharged it here was the red arm: with the sweep disabled the lane printed
   `failure_token=test=TestReplaceLinkCandidatesClearsPendingCandidatesLeftOutOfTheRun` —
   named attribution. **A green integration run is only trustworthy as evidence when a red
   arm has named the test.**

3. **A pre-existing test was vacuous and hid B-08.**
   `TestRouteClassMuxDeclaresBatchBeforeTransportRegistration` asserts that a deadline
   *exists* — interactive has one too, so it passes under the defect. It stayed green in the
   item-1 red run. Asserting presence cannot detect a wrong budget; the assertion has to
   measure the value.

4. **`sed -i` on a CRLF file produces a phantom `M` in `git status`.** `git diff` shows no
   hunks and `git hash-object` matches HEAD. Compare blobs before concluding an edit leaked.
