# M-09 Inventory Clock QA-Unblock Checkpoint

- status: `correction_complete_pending_fixed_sha_review`
- milestone task: `019f5d00-0b82-7b61-9920-32c7bd490333`
- correction contract: `M-09-CORR-03` (final authorized attempt; no retry)
- accepted dispatch base: `371b2daf1c45c65d437779a1196dcab7fd85a8fd`
- runner correction base: `5da9e774bd3acb0ebc6a72b6741a52c5570c7847`
- prior reviewed SHA: `230dc78306d3775894a00b5424238529382cc9b0`
- correction scope: narrow the governed Oracle runner from the full live smoke
  test to the exact product lookup subtest required by M-09-C05.
- selector: `^TestOracleLiveSmoke$/^product_lookup$`.
- proof: governed runner Pester contract 14/14 PASS.
- governed Oracle evidence:
  `F-02-oracle-catalog-cutover/_fixed-sha-oracle-evidence.md`; exit 0,
  source `oracle/sankhya`, `read_only=true`, and
  `positive_codprod_observed=true` at the frozen SHA.
- limitations: no guessed mapping or provider/Oracle write occurred.
- prior C01 review findings were resolved by the final authorized correction.
- final correction outcome:
  - generation requires non-nil positive canonical `InternalProductID` for
    filtering, conflict comparison, deduplication, candidate IDs, and persisted
    identity; legacy `ProductID` is not promoted;
  - nil/zero/negative canonical IDs with positive legacy metadata produce only
    unresolved candidates without canonical identity;
  - OpenAPI requires nullable `brand_name` and `product_group_name`, matching
    the existing SDK contract.
- correction proof:
  - targeted `product_links` Go suite with absolute repository `.gocache` — PASS;
  - integration-tag generation fixture compilation without execution/database
    writes — PASS;
  - SDK runtime and OpenAPI parity — PASS (40/40).
- side effects: no network, database, provider, or Oracle action; no dependency
  installation.
- retry usage: final Portfolio-authorized exception consumed; no retry remains.
- frozen final SHA: `97fd4b58d55a7d14a2b45f0c3bae15b2e374822a`.
- fixed-SHA review: PASS with no findings; C02-C05 reuse was accepted as not
  invalidated by the final C01 correction.
- proportional QA: FAIL. The targeted Go lane passed, but registered full Go
  command `go test ./... -count=1` failed in
  `TestStockRiskServiceClassifiesOversellAndFilters` with `len(items)=0, want 1`.
- QA stop: SDK, residue scan, runner-contract, and frozen-SHA Oracle lanes were
  not run after the deterministic failure. C03 and C05 therefore remain not
  completed for this QA SHA.
- QA evidence:
  - `_fixed-sha-qa/deterministic-qa.md`
  - `validation-result.md`
- retry: none authorized by `M-09-CORR-03`.
- Portfolio subsequently authorized test-only correction `M-09-QA-CORR-01` at
  dispatch base `ee71d7ab2a66f1a55bcc5dc2e9928c34dab78eb2`.
- inventory clock correction:
  - only `TestStockRiskServiceClassifiesOversellAndFilters` changed;
  - its shared observation now uses fresh runtime UTC time within the existing
    30-minute policy;
  - production freshness behavior, risk classification, quantities, filter,
    and assertions are unchanged.
- correction proof with absolute repository `.gocache` and `-count=1`:
  - exact inventory test — PASS;
  - inventory application package — PASS;
  - full `apps/server_core` `go test ./...` lane, run once — PASS.
- correction evidence: `qa-inventory-clock/validation.md`.
- side effects: no production, database, network, provider, or Oracle changes
  or writes; no dependency installation.
- next: freeze the containing correction commit, request fixed-SHA review, then
  restart proportional QA from the beginning only after review passes.
