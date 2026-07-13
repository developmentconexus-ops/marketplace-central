# M-09 Final C01 Correction Checkpoint

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
- final fixed-SHA review:
  - M-09-C02 through M-09-C05 pass.
  - M-09-C01 fails because product-link candidate generation can still promote
    legacy `ProductCandidate.ProductID` instead of requiring positive canonical
    `InternalProductID`.
  - OpenAPI does not require nullable `brand_name` and `product_group_name`
    although Go always emits them and SDK requires them.
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
- next: freeze the correction commit SHA, run `mpc-verifier` fixed-SHA review,
  then proportional QA only if review passes.
