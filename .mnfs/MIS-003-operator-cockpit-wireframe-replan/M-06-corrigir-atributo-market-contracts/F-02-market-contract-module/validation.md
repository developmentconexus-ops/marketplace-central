# F-02-market-contract-module — validation

```yaml
id: F-02
type: feature-validation
status: complete
owner: CHIP-SAT
parent: M-06
updated: 2026-07-16
branch: chip/sat-m05f01-m06f02
governance_anchor: a49168e641ffd6f61932ca57c29b1d1bdcde2fb0
head: 72551f20
```

## Slice trail (all reviewed before dependent work; ledger rows 51–79)

| Slice | Commit | Review verdict |
|---|---|---|
| M1 migrations 0043/0044 + spec.md | 947bc238 | ACCEPT |
| M2 domain + ports (CollectorPort verbatim IC-04) | a3e70538 | ACCEPT |
| M3 observation repo (corrective: tautological INSERT → plain VALUES) | 9b1e6501 | ACCEPT (delta) |
| M4 reference repo | 5cd992fa | ACCEPT |
| M5 collection service | 0f45d11c | ACCEPT |
| M6 honest-empty read service | c8ba9114 | ACCEPT |
| M7 /market transport (corrective: per-id TrimSpace + lock test) | 8fa6c0fa | ACCEPT (delta) |
| R1 composition registration (reads only) | b5f3ea70 | ACCEPT |
| C1 OpenAPI market paths/schemas + SDK (same commit) | b16fcf80 | ACCEPT, zero findings |
| C2 category-attribute contract reservation | 3d17b906 | ACCEPT-WITH-CONDITIONS (1 non-blocking, see findings) |
| I1 integration tests + absence proofs (+ migrate fixture 37→39) | eb479ddc | ACCEPT |
| Governance registry: market module entry | e1819778 | orchestrator corrective (GOV_MODULE_COVERAGE) |
| I1 corrective: catalog_product_id assertion (test-only) | 72551f20 | DELTA ACCEPT |

## Hard gates (feature-defining)

- NO production CollectorPort implementation, NO seed, NO scheduler, NO UI — enforced twice:
  composition negative test (`root_test.go`, substring `marketapp.NewCollectionService(`)
  AND I1 absence proofs (filepath.Walk scans, each with inline self-test proving the scan
  fails on a synthetic violation): no non-test CollectorPort method impl; `INSERT INTO market_`
  only in the two adapter repositories; forbidden scraping/ML strings absent; root.go free of
  collection/scheduler wiring.
- Migrations exactly 0043+0044; 0045 reserved-unused (no file exists — verified, no empty file).
- Stored `source`/`captured_at` NOT NULL (SQL CHECK + domain ValidateStored + I1 SQLSTATE 23502 test).
- Honest-empty per IC-04: unknown ids → 200 item-level `evidence_state=no_price_evidence` with
  ALL explicit-null signals (raw-body null-key assertions in transport + I1 tests); >200 ids →
  422 `too_many_ids`; missing installation_id → 400 `installation_required`.
- 6-signal separation round-trips independently (I1: all four money pairs + catalog_stats +
  competitive/match state asserted separately).

## Ladder

### L0 — build + static + contract (COMPLETE)

| Check | Result | Evidence |
|---|---|---|
| `go build ./...` (apps/server_core) | exit 0 | run 2026-07-16 |
| Governance (`harness.ps1 governance -BaseSha a49168e6…`) | status=passed | clean temp worktree `mpc-gov-check-f02` @ e1819778; only baseline exceptions; first run failed GOV_MODULE_COVERAGE id=market → registry entry added (e1819778) |
| SDK `tsc --noEmit` | exit 0 | sdk-runtime @ HEAD |
| SDK vitest | 52/52 PASS | includes C1 market tests (URL order, all-null fixture, six-signal) + C2 category-attribute tests |
| apps/web vitest | 11/11 PASS | baseline unchanged — zero web files touched (seam grep clean) |

### L1 — unit + integration lanes

| Check | Result | Evidence |
|---|---|---|
| `npm run harness:unit` (canonical) | **status=passed**, exit 0 | run_id 22f693be7c304771ad12c56e06b055fb (`scripts/.runs/…`); go unit ok + web vitest 11/11 |
| `npm run harness:integration` (canonical, run 1) | 3 failures: allowlisted TestPhase1SmokeFlow + known F-2 flake + `TestMarketContractRoundTripAndTenantIsolation` (mine) | tasks/b9289goye.output; migrations 39/0 idempotent, resource_count=0, session pg 51700 |
| Corrective (test-only, 1 line) | market_contract_test.go:285 expected-value string doubled "round-trip-" vs seed convention (`"catalog-"+product_id`); stored value was correct | commit 72551f20; DELTA ACCEPT (review-verdict-f02-i1.md append); isolated repro: per-run `mpc_test_<32hex>` DB, 15/15 PASS × 3 rounds |
| `npm run harness:integration` (canonical, run 2, post-fix) | failure set = ONLY allowlisted `TestPhase1SmokeFlow` + pre-existing F-2 flake (`TestOrderRepositoryDuplicateIdentityGroupPreservesIDSetAndRemainsAmbiguous`, hub landing fix on main) — ALL F-02 market tests PASS live | tasks/bm53levbh.output; migrations 39/0 idempotent, resource_count=0, session pg 51700 |

### L2 — live QA

Not chip-owned (feature grain closes at L1 + evidence; milestone close = hub dual gate + QA).

## Findings / carried notes (for hub)

1. **Migration-count fixture drift (new, resolved)**: `internal/platform/migrate/runner_test.go`
   hardcodes canonical migration count (37→39 bumped in eb479ddc). Any future migration grant
   must bump it; surfaced only in I1's full-suite run because slice testing was module-scoped.
   Profile candidate: "migration grant ⇒ runner_test.go count bump" note.
2. **Governance module coverage (new, resolved)**: creating a module directory fails the
   governance lane until `contracts/governance/modules.json` declares it (GOV_MODULE_COVERAGE).
   Registered market @ e1819778 (precedent: dashboard @ 17cce1a9, F-01).
3. **C2 error-code casing (non-blocking condition from review)**: `category_not_found` /
   `provider_unavailable` are lowercase like the committed market transport codes, but
   `ErrorResponse.error.code` documents `MODULE_ENTITY_REASON` uppercase (openapi.yaml:5330-5340).
   Two casing conventions now coexist in the contract; resolve when the M-06 F-01 handler lands.
4. **M1 (carried)**: migrations `market_test.go` lacks CREATE INDEX/TRIGGER negative asserts.
5. **M2 (carried)**: blank percentile strings pass `validateCatalogStats` when NSellers≥5.
6. **M3 (carried)**: `scanMoney` lone-nil pair panic risk — guarded by DB pair CHECKs.
7. **M4 (carried)**: catalog_stats marshal block duplicated from M3 (factoring forbidden in-slice).
8. **M7 (carried)**: dead `filter.` TrimPrefix in query.go; no handler-level empty-ids test.
9. **I1 (review notes)**: broad `search?` substring in forbidden-string scan; absence-proof
   test lives in the integration-tagged file so it only runs in the DB lane despite needing none.
10. **F-2 ratification (from M-05 F-01, hub-requested)**: pre-existing flake
    `TestOrderRepositoryDuplicateIdentityGroupPreservesIDSetAndRemainsAmbiguous`
    (`orders/adapters/postgres/order_repo_test.go:172`) — crypto/rand `NewMPCLineID` makes the
    positional assertion fail ~2/3 of live runs; fix = set-containment. Hub landing fix on main.
11. **F-4 ratification (from M-05 F-01, hub-requested, profile §3 candidate)**: Windows Go
    `time.Now()` carries 100ns; postgres timestamptz stores µs. Integration fixtures asserting
    `.Equal()` on round-tripped timestamps must `Truncate(time.Microsecond)` (production
    convention `orders/domain/sankhya_linkage.go:178`). Signature: got/want differ only in the
    7th fractional digit.
