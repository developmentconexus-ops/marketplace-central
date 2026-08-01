# M-02 (sync-core-seam) — EVIDENCE

milestone: MIS-007/M-02-sync-core-seam · session `local_c0c3c6c4-9f68-4e6d-ade5-50d23046b13c` (hub) ·
this orchestration session on `claude/gallant-banach-2f909b` · BASE-SHA `295e293fdc273ed0fad9c3eb2445b7f2152586ed`
· HEAD `7d02bc75ba2c0ec899927f2aa0164277fa9ab2ed` (post hub-directed fix, see addendum below;
milestone-review tip was `e51c09f6219a0b21d1d51633c6a781ebcdb95eb7`). Not pushed.

Process: operator override "review enxuto" — no dual gate. 1 adversarial sonnet reviewer per
feature at feature-end, 1 cold adversarial reviewer on the full milestone diff at the end.
4 features dispatched to sonnet subagents in worktree isolation per the internal DAG
(F-01 ∥ F-03 ∥ F-04; F-02 after F-01 merge). Opus orchestration only — no production code
written by this session.

## Deliverables (write-set)

| File | Kind | Feature |
|---|---|---|
| `apps/server_core/migrations/0086_channel_fees.sql` (+`channel_fees_test.go`) | NEW channel_fees ledger table (IC-01) | F-01 |
| `apps/server_core/migrations/0087_divergences.sql` (+`divergences_test.go`) | NEW divergences table, partial unique one-open-row (IC-02) | F-01 |
| `apps/server_core/migrations/0088_order_shipments.sql` (+`order_shipments_test.go`) | NEW order_shipments table (IC-03) | F-01 |
| `apps/server_core/migrations/0089_orders_marketplace_orders_sync_fields.sql` (+test) | NEW additive orders buyer-fiscal 9-col (IC-03) | F-01 |
| `apps/server_core/tests/integration/mis007_f01_core_ddl_must_fail_test.go` | NEW permanent must-fail regression (23514/23505) | F-01 |
| `apps/server_core/internal/platform/migrate/runner_test.go` | MOD migration count 72→76 | F-01 |
| `apps/server_core/internal/modules/channelfees/{domain,ports,adapters/postgres}/**` (+tests) | NEW ChannelFeeWriter/Reader ports + Postgres impl | F-02 |
| `apps/server_core/internal/modules/divergences/{domain,ports,adapters/postgres}/**` (+tests) | NEW DivergenceRecorder/Reader ports + Postgres impl | F-02 |
| `contracts/governance/modules.json` | MOD +2 entries (channelfees, divergences) | F-02 |
| `apps/server_core/internal/modules/sync/application/scheduler.go` | MOD `inferIncremental` replaces hardcoded `false` | F-03 |
| `apps/server_core/internal/modules/sync/application/cursor_contract.go` (+test) | NEW `AssertTerminalCursor` reusable contract (M-04/M-06 consume later) | F-03 |
| `apps/server_core/internal/modules/sync/application/products_regression_test.go` | NEW real-job-through-real-scheduler regression | F-03 |
| `apps/server_core/internal/platform/archguard/archguard_test.go` (+3 testdata fixtures) | NEW AST shrinking-allowlist guard, alias-resolving | F-04 |

## Criterion verdicts (M02-C1..C7 + 2 user-drive, per validation-contract.md)

| Crit | Verdict | Evidence | Type |
|---|---|---|---|
| M02-C1 (4 migrations apply clean, idempotent, exact shapes) | PASS | `harness:integration` → `migrations_first=76, migrations_second=0, status=passed`; `runner_test.go` fixture 76 = `ls migrations/*.sql` count | ran |
| M02-C2 (CHECKs reject invalid, proven by INSERT) | PASS | `mis007_f01_core_ddl_must_fail_test.go`: 23514 `channel_fees_currency_when_amount_check`; 23505 `divergences_one_open_row`, accept after resolve | ran (real Postgres) |
| M02-C3 (port contracts: resolution/tolerance/refusal) | PASS | `channelfees`/`divergences` postgres+domain suites green on real Postgres; layer 2→1→absent-typed, layer-3 commission/freight asymmetry, tolerance stock=0/tariff=R$0.01, one-open-row+auto-resolve, detected_at immutable | ran (real Postgres) |
| M02-C4 (nil cursor must-fail, named) | PASS | `cursor_contract_test.go` asserts message contains `"terminal cursor must be non-nil"` | ran |
| M02-C5 (incremental fix, zero products regression) | PASS | `products_regression_test.go` runs real `synccomposition.NewProductsJob` through real `Scheduler.RunOnce`, asserts `incremental=false` + cursor fields unchanged | ran |
| M02-C6 (guard allowlist, must-fail names new site) | PASS | baseline 4 real sites; 5th-site + aliased-site fixtures fail-and-name; shrunk-allowlist fixture proves not hardcoded | ran |
| M02-C7 (zero raw PII in new schema) | PASS | 0089's 9 buyer-fiscal columns map 1:1 onto the 9 fields the FE drawer renders (p5-prerequisites.md); `uf_nome`/`fetched_at` correctly excluded; no raw/jsonb billing column | ran (cross-checked against FE) |
| M02-U1/U2 | NOT RUN — hub-driven browser QA, out of scope for this session (dev stack forbidden to chips) | — | deferred to hub |

## P5 ladder (ran, full milestone diff, from `apps/server_core`)

- `go build ./...` → GREEN, whole module.
- `go vet ./...` (plain + `-tags=integration`) → clean, no output.
- `go test ./...` (unit) → all packages green, zero collisions across the 4 features.
- `npm run harness:integration` (hermetic Postgres, full `//go:build integration` glob) → `status=passed`, `resource_count=0` (no docker leak).
- `gofmt -l` on all new packages → clean.

## Cross-feature checks (milestone-level, not visible per-feature)

- F-01 ↔ F-02 seam: hand-compared `channelfees`/`divergences` Postgres adapter SQL against the actual merged 0086/0087 migrations — exact match on natural-key constraint names and the partial-unique target.
- Ownership disjointness: `git diff 295e293f..e51c09f6 --stat` — no file touched by two features. F-01 (migrations + its integration test), F-02 (`internal/modules/{channelfees,divergences}/**` + governance registry), F-03 (`sync/application/{scheduler,cursor_contract}*`), F-04 (`internal/platform/archguard/**`).
- Governance registry: 2 new entries well-formed (`composition_required:false`, `dependencies:[]`); grepped repo-wide for any consumer import of `channelfees`/`divergences` — zero, confirming nothing is wired yet (correct — M-05/M-06/M-07 wire consumers).
- F-03 deferral (`AssertTerminalCursor` exported but not wired into `runJob`) confirmed intentional and milestone-scope-correct: no M02 criterion assumes live scheduler enforcement; wiring is M-04/M-06's job when they write real job bodies.

## Findings / follow-ups (non-blocking)

- **F-04 round 1 → REFUTE → fixed.** Adversarial review found `scanWiringFile` matched capability args by literal identifier name only — a one-line local-variable alias (`mlAlias := mercadoLivreCapabilities`) bypassed detection silently. Fixed with `resolveCapabilityAliases` (single-hop AST fixpoint over local `:=`/`=` assignments). New fixture `testdata/aliased_site/` + `TestFixture_AliasedSiteIsDetectedAndNamed` proves it, red-before-fix/green-after captured. Re-verified independently by orchestrator before merge.
- **F-02 round 1 → REFUTE (test-only) → fixed.** `TestResolveListingFeesIgnoresLayer3` seeded its layer-3 row under a different `subject_type` than the ladder ever queries — vacuous, would pass with the layer filter deleted. Underlying SQL logic independently verified correct by the reviewer regardless. Fixed by reseeding under the identical natural key; red→green proof captured (temporarily deleted the layer filter, confirmed the rewritten test now fails naming the regression, restored, confirmed green).
- **Worktree provisioning staleness** (process finding, not a code defect): 3 of 4 isolated feature worktrees (F-01, F-03, F-02) forked from a point 20-480+ commits behind the intended base. F-01/F-03 self-healed via `git merge --ff-only` after confirming strict ancestry. F-02's first dispatch was even further stale (predated all of MIS-007) and correctly reported BLOCKED rather than guessing; resumed with an explicit fetch-and-ff-merge onto `claude/gallant-banach-2f909b` (the actual milestone branch, not `main` — `main`'s tip does not carry this milestone's merges). **FINDING → hub**: worktree-isolation pre-provisioning for this session pool is stale; worth a harness-side fix (fresher snapshot or an explicit rebase step baked into the isolation dispatch) so implementers don't have to self-diagnose this per dispatch.
- **Minor duplication** (non-blocking): `channelfees/boundary_test.go` and `divergences/boundary_test.go` are near-identical copy-pasted AST import-scanners. Same feature (F-02), not a cross-feature collision. Could collapse into a shared test-support helper later.

## P6 review (lean/enxuto — per-feature + final cold pass, no dual gate)

- F-01: adversarial sonnet — AGREE (verified via live mutation testing: neutered each CHECK/unique constraint independently, confirmed must-fail regression test goes red, reverted clean).
- F-03: adversarial sonnet — AGREE (verified `detected_at`-analog i.e. `inferIncremental` branch coverage + real must-fail flip; confirmed `AssertTerminalCursor` deferral is in-scope per brief).
- F-04: adversarial sonnet — REFUTE (alias bypass, HIGH) → implementer fix → orchestrator independently re-ran full suite, confirmed fix, then merged.
- F-02: adversarial sonnet — AGREE + 1 REFUTE (test theater in layer-3 isolation test) → implementer fix with red→green proof → merged.
- **Milestone-final cold review** (independent sonnet, full diff `295e293f..e51c09f6`, all 7 criteria re-run for real, cross-feature seams checked): **AGREE**. No blocking defects. Recommends proceeding to hub merge + M02-U1/U2 browser QA.

**Verdict: milestone-level AGREE.** Lean review process complete (4× per-feature + 1× final cold pass, 2 real defects found and closed via fix-and-reverify, not waived). Ready for hub merge; M02-U1/U2 (browser QA) remain hub-driven per dev-stack ownership.

## Addendum — hub-directed post-review fix (HOLD-MERGE)

Hub reviewed the merged diff post-CLOSED and held the merge on one finding, `sync/application/scheduler.go`:

> `inferIncremental`'s switch matched `case "incremental", "sweep": return true`. `"sweep"` is not a
> ratified phase — ADR-07 (`mission.md:183-187`) defines only `backfill → incremental`. The function's
> own doc comment already declared "unrecognized phase resolves to false", directly contradicted by the
> `"sweep"` case. Live risk: M-09 (already merged) computes `last_success_at = GREATEST(last_full_sync_at,
> last_incremental_at)`; a future milestone naming a repair pass `"sweep"` would silently record it as
> incremental, `last_full_sync_at` would never advance, and the M-09 health card would show a stuck-old
> full-sync timestamp — root cause here, symptom in another milestone months later.

Verified independently before dispatching the fix: read `mission.md:183-187` directly, confirmed ADR-07
ratifies only `backfill → incremental`, no sweep.

Fix (sonnet subagent, orchestrator-verified, no scope beyond the 2 named files):
- `scheduler.go`: `inferIncremental` now matches only `"incremental"` → `true`; every other value
  (including `"sweep"`) falls to `false`. Doc comment corrected to cite ADR-07 and warn against
  re-adding `sweep` as a special case.
- `scheduler_test.go`: both table-driven tests' `sweep→true` rows renamed to
  `"unrecognized phase sweep falls to tolerant default (ADR-07 has no sweep phase)"` → `false`.

Verification (re-run independently by orchestrator, not just trusting the implementer):
```
--- PASS: TestInferIncrementalTolerates/unrecognized_phase_sweep_falls_to_tolerant_default_(ADR-07_has_no_sweep_phase)
--- PASS: TestRunOnceDerivesIncrementalFromCursorPhase/unrecognized_phase_sweep_falls_to_tolerant_default_(ADR-07_has_no_sweep_phase)
```
`go build ./...` clean. `go vet ./internal/modules/sync/...` clean. `git diff --stat HEAD~1..HEAD` —
only `scheduler.go` (+`scheduler_test.go`) touched, 9 insertions/6 deletions, nothing else.

No REQUEST filed — hub's instruction was to remove the unratified case, not invent repair-phase
semantics, so no new decision was made at this level.

**Tip after fix: `7d02bc75ba2c0ec899927f2aa0164277fa9ab2ed`.** CLOSED reported to hub with this evidence.
