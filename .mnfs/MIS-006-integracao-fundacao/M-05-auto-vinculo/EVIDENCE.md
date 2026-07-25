# M-05 auto-vínculo produto↔anúncio — Evidence Pack

Branch `claude/elated-albattani-323511`, worktree `.claude/worktrees/trusting-mayer-a5c8f6`.
BASE-SHA `e3c081ae43b72af070185939253b745080acf68b`. Chip tip **`f9dbc868`**.

This pack is updated in place and committed on top of the code it describes, so the tip named
here is always the last CODE commit; the commit carrying this file is the tip of the branch.
(Round 2's cold gate raised the mismatch as a blocker — the pack said `5807d634` while the branch
was at `a4a0ad89`, the docs commit. Naming both is the fix.)

Commits (`git diff e3c081ae..f9dbc868`, zero files under `apps/web`):

- `124cd9e8` F-04 — `seller_sku` anchored on `p.CODPROD`, not `p.REFFORN`, with an `IsValidCodprod` validity guard, mirrored in `erp_import/adapters/internalread/`
- `100f0343` F-03 — E10 decision trail, migration `0082_product_link_decisions.sql`, written in the same transition as the state change
- `b01a2579` F-02 — corroborated auto-approve + CONFIRM queue (D-121-2), OpenAPI + sdk-runtime enum in the same changeset
- `84323d0a` F-01 — generation wired to the products sync, and the xlsx import's existing call changed from fail-hard to logged-and-isolated (the import already called it at BASE-SHA; this commit did not introduce that call)
- `5807d634` correctives — the four defects the round-1 dual gate found
- `f9dbc868` correctives — the three defects the round-2 dual gate found, plus two it raised as non-blocking (see **Round-2 dual gate** below)

Ruling applied throughout: where any other artifact still says "auto-approve is exclusive to the
EAN-exact-unique path" or "CODPROD-único auto-aprova", `milestone.md` at BASE-SHA wins
(ADR-05 AMENDADO D-121 / REVISADO D-121-2, ratified by the operator).

## Ladders

| Ladder | Command | Result |
|--------|---------|--------|
| L1 Go build/vet | `cd apps/server_core; GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go build ./... && go vet ./...` | clean, zero output |
| L1 Go tests | `go test ./...` | **106 packages ok, zero FAIL** |
| L1 product_links | `go test -count=1 ./internal/modules/product_links/... -v` | **all PASS**, zero FAIL (69 tests at round 1, 86 after the round-1 correctives, more after round 2) |
| L1 tsc | from the MAIN checkout root: `npx --no-install tsc -p .claude/worktrees/trusting-mayer-a5c8f6/packages/sdk-runtime/tsconfig.json --noEmit` | exit 0 (run from the main root deliberately — a worktree-local `npx --no-install tsc` is a vacuous pass, D-120 finding) |
| L2 integration (hermetic PG) | `pwsh -NoProfile -File scripts/harness.ps1 -Command integration` | **status=passed**, `migrations_first=69`, `migrations_second=0`, `run_id=847fa06f6784412eb38ba74ea7fc8cf7`, `run_dir=scripts/.runs/847fa06f6784412eb38ba74ea7fc8cf7` — **and proven non-vacuous**, see below |
| gofmt (CRLF-aware) | `tr -d '\r' < FILE \| gofmt -d` over every touched `.go` | zero diff (a bare `gofmt -l` false-alarms on CRLF, profile §3) |

**The lane's `status=passed` is proven to mean the E10 tests RAN.** Both round-2 gates raised the
same doubt independently, and they were right to: `scripts/harness/Postgres.psm1:276` runs
`go test -tags=integration … -count=1` with **no `-v`**, every DB test opens with
`SkipWithoutTarget`, and `summary.txt` records only three lines
(`target`, `status`, `run_id`) — so a run in which all six E10 tests skipped would report
`status=passed` identically. Proof by must-fail: with `got != 1` changed to `got != 999` in
`TestCorroboratedApprovalWritesTheE10Row`, the same command returns

```
failure_token=test=TestCorroboratedApprovalWritesTheE10Row
status=blocked
postgres lifecycle failed reasons=HPG_TEST_FAILED exit_code=1
```

The assertion was restored in place and the lane re-run green. The tests execute against real
Postgres. **FINDING**: the lane should record executed/skipped test counts in `summary.txt` —
as it stands, `status=passed` cannot by itself distinguish a green run from a fully-skipped one,
and no reviewer can close an L2 criterion on that artifact alone.

`scripts/harness.ps1 -Command integration` must run under **pwsh (PS7)**, not Windows
PowerShell 5.1 — under 5.1 the lane dies with
`[System.Security.Cryptography.RandomNumberGenerator] não contém um método denominado 'Fill'`
and reports `status=blocked`. **FINDING**, reported to the hub for profile ratification.

## Must-fail proofs (executed, then restored)

| Criterion | Revert applied | Observed failure |
|-----------|----------------|------------------|
| M05-C17 | `if hardNeg, detail := detectHardNegative(...); false && hardNeg` | `--- FAIL: TestHardNegativeBlocksTheAutomaticPathEvenWithBothAnchors: match_status = ACCEPT, want the kit/unit contradiction to block the automatic path` |
| M05-C21 | `LinkCandidateMatchStatusConfirm = "REVIEW"` | `--- FAIL: TestConfirmationAndReviewAreCountedSeparately: ... want confirmation and review as separate groups` (one REVIEW group of 3) |
| M05-C18 | `p.CODPROD = :%d` → `p.REFFORN = :%d` in `buildFindProductsQuery` | `--- FAIL: TestFindProductsQueryMatchesSellerSKUAgainstCodprod: seller_sku was not matched against CODPROD` (query showed `AND (p.REFFORN = :1)`) |
| M05-C10 / gate defect 1 | terminal-state guard removed from `AutoApproveCandidate` | `--- FAIL: TestAutoApprovalNeverReopensAListingTheOperatorRejected` — transition emitted with `PreviousState:"rejected", NextState:"resolved"`, `Actor{ActorType:"system", ActorID:"auto_linker"}` |
| M05-C22 / gate defect 2 | `decisionRuleForCandidate` `MatchStatus` gate removed | `--- FAIL: TestApprovingAnAmbiguousCandidateIsRecordedAsAManualCall/colliding_EAN: rule_matched = "exact_ean_unique", want manual: no anchor resolved this listing` (and the `conflicting SKU` subtest, `exact_codprod_unique`) |
| M05-C8 / gate defect 3 | undo stopped writing its own decision | `--- FAIL: TestUndoSupersedesTheDecisionItReverts: decisions = 1, want the undo recorded alongside what it reverted` |
| L2 lane genuinely executes | `got != 1` → `got != 999` in `TestCorroboratedApprovalWritesTheE10Row` | lane returned `failure_token=test=TestCorroboratedApprovalWritesTheE10Row`, `status=blocked` — see the L2 section above |

Every revert was undone **in place** and the suites re-run green afterwards. (Reverting via
`git checkout -- <file>` on a file carrying uncommitted work destroys that work — it happened
once here, was caught immediately by a verification print, and every subsequent proof edited
and restored in the same command. Reported as a FINDING below.)

## Criteria ledger

| ID | Verdict | Evidence |
|----|---------|----------|
| M05-C1 | **could-not-run at L2 (named) — `ran` at L1** | The internal trigger exists and is proven without HTTP: `erp_import/application/import_service.go:154-177` calls `s.generator.GenerateLinkCandidates` inside `runImport` after the snapshot is persisted, per installation returned by the enqueuer; `import_service_test.go` `TestImportServiceGenerationFailureDoesNotFailTheImport` exercises that path with no transport involved. The contract's L2 form (post-import candidate row on the dev stack) is **could-not-run**: the dev stack is a hub-owned seam a chip may not boot, and `listings` = 0 / `integration_installations` = `pending_connection` means IO-A has no listing to pair with. Hub re-drive required. |
| M05-C2 | ran | `generation_service.go:481-486` — EAN-only ⇒ `MatchStatusConfirm`, detail `sem CODPROD para corroborar o EAN`. `auto_link_policy_test.go` `TestSingleAnchorGoesToConfirmationNeverAutoApproved` asserts the state, the exact warning text, zero approved link and zero `actor=system` decision (golden `ean:7909251260214` → 100000). |
| M05-C3 | ran (L1 + L2) | `resolution_service.go:270-276` stamps `DecisionRule=concordant_codprod_ean`, `Actor{ActorType:"system", ActorID:"auto_linker"}`, `CollisionsAtDecision=&collisions`. `TestConcordantAnchorsAreTheOnlyAutomaticPath` asserts rule + actor + `collisions_at_decision = 1`; the single-anchor tests assert no `actor=system` row exists at all. DB-level: `tests/integration/product_link_decisions_test.go` `TestCorroboratedApprovalWritesTheE10Row` + `TestTheDatabaseRefusesASystemDecisionOnASingleAnchor`, run in the hermetic lane above. |
| M05-C4 | ran | `generation_service.go:307` `buildCollisionCandidates` + `applyCollisionScore` (`REVIEW`, conf 20/BAIXA, detail naming the colliding codprods). `TestCollidingAnchorGoesToReviewAndApprovesNothing` uses the real ERP collision `7896902180697` → 4 products (42535-42538) and asserts zero approvals. |
| M05-C5 | **ran, with a wording reconcile requested** | `applyUnresolvedScore` (`generation_service.go:556`) — no anchor ⇒ confidence 0 and a **non-empty reason**, never a silent blank; `generation_service_test.go` `TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons` asserts both. The status is `NO_CANDIDATE`, not the literal `REVIEW` the criterion's wording asks for: `/vinculos` keys its "sem candidato / Criar produto / Ignorar" affordance **and its batch-select guard** off that exact value (`QueueRow.tsx:77`, `VinculoDrawer.tsx:63`), so emitting `REVIEW` there turns an anchorless listing into an ordinary review row the operator can bulk-approve into a wrong link — the `/catalogo`-503 class of invisible regression. The criterion's *intent* (no silent blank, nothing auto-approved) is met. **Hub reconcile requested**: amend C5's wording, or hand the status flip to M-06 with the screen change. |
| M05-C6 | ran | The decision reads the generator's own counts: `autoApprovals` carries `len(skuMatches.Products)` (`generation_service.go:192-199`), and `buildCollisionCandidates` takes `skuMatches`/`eanMatches` directly. `grep -rn "validEANCounts\|identityQuality" internal/modules/product_links/` returns **zero hits** — AC-04 clear. No new counting function in the write-set. |
| M05-C7 | **ran (signature) / could-not-run (curl)** | `product_links/transport/http_handler.go:88-100` `Register` is byte-unchanged in the write-set: all 11 routes, including `/product-links/link-candidates/generations`, still registered; `git diff e3c081ae..84323d0a -- .../transport/` is empty. The `curl … 200` half is could-not-run for the same reason as C1 (dev stack = hub seam). |
| M05-C8 | ran | `resolution_service.go:252-256` — any decision with `SupersededBy == ""` already in force ⇒ `AutoApproveCandidate` returns `(false, nil)` and writes nothing. `TestRerunningGenerationDoesNotDuplicateTheAutomaticLink` runs generation twice over the same snapshot and asserts one link, one decision. DB-level: `TestListingIdentityKeepsApprovalIdempotentPerVariation` (integration lane). |
| M05-C9 | ran | `migrations/0082_product_link_decisions.sql` contains exactly one `CREATE TABLE IF NOT EXISTS product_link_decisions` and **zero** `ALTER TABLE`; `migrations/product_link_decisions_test.go` `TestE10MigrationDoesNotTouchProductLinks` asserts it statically, and `TestOneProductMayLinkToManyListings` (integration) proves N listings per product survive. AC-09 clear. |
| M05-C10 | ran + must-fail proven | Two guards, because the trail alone was not enough. (a) A decision in force blocks the automatic path — `TestAutoApprovalNeverOverridesAnOperatorDecision`; DB-level `TestAnOperatorOverrideSupersedesTheAutomaticDecision` proves the earlier row receives `superseded_by` instead of being erased. (b) The link's own terminal state blocks it too (`resolution_service.go`, `state == rejected \|\| resolved ⇒ (false, nil)`), which covers the two cases the trail cannot see: **rejection writes no decision row at all**, and links approved before 0082 have none. `TestAutoApprovalNeverReopensAListingTheOperatorRejected`, must-fail proven above. |
| M05-C11 | ran | Import half: `import_service.go:170-175` logs and continues, so the report stays `completed` — `TestImportServiceGenerationFailureDoesNotFailTheImport` asserts no error, completed status and a completed persisted snapshot with a failing generator. Sync half: `sync/composition/products_job.go:115-122` — `TestProductsJobSurvivesALinkCandidateRefreshFailure` (cycle succeeds, cursor keeps the real `processed=42`) and `TestProductsJobDoesNotRefreshWhenTheSyncFailed`. |
| M05-C12 | ran | `migrations/product_link_decisions_test.go` `TestProductLinkDecisionsMatchesE10` asserts the column set (`link_id, rule_matched, actor, collisions_at_decision, created_at, superseded_by`) and the accepted `rule_matched` values; `TestCollisionsAtDecisionIsHonestlyUnknownNeverZero` proves the column is nullable with no default (ADR-17 / AC-03). Both re-verified against a real migrated PG in the hermetic lane (`migrations_first=69`). |
| M05-C13 | ran | Migration number `0082`, applied after M-02's `0076` block; `git diff` of the migration shows `CREATE TABLE` only. Inventory fixture bumped 68 → 69 in `internal/platform/migrate/runner_test.go:25,64` (it was left stale by the F-03 commit and caught by the full `go test ./...`). |
| M05-C14 | ran | `generation_service.go:471-476` — CODPROD-only ⇒ `MatchStatusConfirm`, detail `sem EAN para corroborar o CODPROD`. Second sub-case of `TestSingleAnchorGoesToConfirmationNeverAutoApproved` (golden `sku:100001`), asserting zero approved link and zero `actor=system`. |
| M05-C15 | ran | `autoApprovals` (`generation_service.go:192`) selects **only** `MatchStatus == ACCEPT`, which `buildConcordantCandidate` alone produces; `AutoApproveCandidate` refuses anything else with `PRODUCT_LINKS_AUTO_APPROVE_NOT_CORROBORATED` (`resolution_service.go:213-232`), proven by `TestAutoApproveRefusesACandidateItDidNotCorroborate`. Positive path: `TestConcordantAnchorsAreTheOnlyAutomaticPath` (golden `100002` / `7909251304727`). |
| M05-C16 | ran | `applyConflictScore` (`generation_service.go:515`) sets `REVIEW` on BOTH sides with the detail `… aponta codprod X (conflito, nenhuma âncora vence)` — no precedence rule anywhere. `TestConflictingAnchorsApproveNothingAndElectNoWinner` (golden IO-E `seller_sku="100000"` + `ean=7899656858195`) asserts neither anchor won and nothing was approved. AC-08 clear. |
| M05-C17 | ran + must-fail proven | `detectHardNegative` overrides to `REJECT` even with both anchors concordant; `TestHardNegativeBlocksTheAutomaticPathEvenWithBothAnchors`. Must-fail proof above. |
| M05-C18 | ran + must-fail proven | `internal_read/adapters/oracle/reader.go:451-460` matches `p.CODPROD` and no `p.REFFORN` clause exists; `TestFindProductsQueryMatchesSellerSKUAgainstCodprod`. Mirrored in `erp_import/adapters/internalread/reader.go:469-476`. Must-fail proof above. |
| M05-C19 | ran | `internal_read/domain/seller_sku.go` `IsValidCodprod` guards the clause; a legacy REFFORN (`"ZP1704.1."`, IO-F) produces `errNoMatchableInput` on its own and does not widen a query that also carries a valid EAN — `TestFindProductsQueryDropsANonCodprodSellerSKU`. |
| M05-C20 | ran | `TestOneProductAutoLinksToManyListings` — two distinct listings carrying the same CODPROD+EAN both auto-approve, neither flagged. DB-level `TestOneProductMayLinkToManyListings` (integration lane). AC-09 clear. |
| M05-C21 | ran + must-fail proven | `domain/link_candidate.go:42-49` `LinkCandidateMatchStatusConfirm = "CONFIRM"`, carried verbatim by the postgres repo (no CHECK on `match_status`) and surfaced at the transport boundary — `contracts/api/marketplace-central.openapi.yaml:6459` enum `[ACCEPT, REVIEW, REJECT, NO_CANDIDATE, CONFIRM]` + `packages/sdk-runtime/src/index.ts` `ProductLinkMatchStatus`. `TestConfirmationAndReviewAreCountedSeparately` counts the two groups apart with the warning preserved. Must-fail proof above. AC-10/AC-11 clear at the API layer (see the user-drive section for the screen layer). |
| M05-C22 | ran + must-fail proven | `decision_trail_test.go` `TestOperatorConfirmationRecordsTheAnchorItWasTakenOn` — the operator confirming the C14 candidate writes `actor=operator, rule_matched=exact_codprod_unique`, never `system`/`concordant_codprod_ean`. And the trail only names an anchor when one actually resolved the listing: `decisionRuleForCandidate` records `manual` for anything that reached the operator through the collision or conflict path, since naming `exact_ean_unique` over a real four-way EAN collision asserts a uniqueness nobody established (ADR-17) and reads as though an anchor had won (AC-08). `TestApprovingAnAmbiguousCandidateIsRecordedAsAManualCall`, must-fail proven above. Undo is itself a decision (`TestUndoSupersedesTheDecisionItReverts`), so a reverted row never reads as still in force.. A rejection likewise supersedes what it overrules (`TestRejectingAListingSupersedesTheDecisionItOverrules`, round 2), and a corroborated candidate approved BY HAND records `concordant_codprod_ean` + `actor=operator` rather than being under-claimed as a single anchor (`TestOperatorApprovingACorroboratedCandidateRecordsCorroboration`). |

Anti-criteria: AC-01 (every E10/`product_links` write tenant-scoped), AC-03 (`collisions_at_decision`
nullable, no default, honest-unknown — and, after round 2, a `*int` the whole way from the
generator, so a caller that read no count cannot be forced to write 0), AC-04 (zero `validEANCounts`/`identityQuality` in
`product_links/*`), AC-05, AC-07, AC-08, AC-09, AC-10, AC-11 — all clear. AC-06: `git status -sb`
shows no upstream; nothing pushed.

## User-drive (D-120 amendment)

| ID | Verdict | Blocker / evidence |
|----|---------|--------------------|
| M05-U1 | **could-not-run** | Conta ML não autorizada: `integration_installations` = `pending_connection`, `listings` = 0. No real anúncio exists to auto-approve. Never passed by code inspection — hub re-drives after the operator completes the OAuth. |
| M05-U3 | **could-not-run** | Same blocker: `/anuncios` cannot reflect a link that no listing can produce. |
| M05-U2 | **could-not-run by the chip — hub drive requested** | The drive needs the dev stack, which is a hub-owned seam a chip may not boot (profile §9). The API-layer equivalent is proven (`TestCollidingAnchorGoesToReviewAndApprovesNothing`, `TestConflictingAnchorsApproveNothingAndElectNoWinner`, `applyUnresolvedScore`): collision, conflict and no-anchor all stay pending, never a false auto-vínculo. |
| M05-U4 | **BLOCKED — screen gap owned by M-06** | `/vinculos` renders one flat queue: `apps/web/src/pages/vinculos/QueueTab.tsx` groups nothing, and `QueueRow.tsx:77` / `VinculoDrawer.tsx:63` special-case only `NO_CANDIDATE`. A `CONFIRM` candidate therefore reaches the screen (it survives to `data-match-status`, and the warning text is in the candidate detail) but is **not visually separated from the review queue**, which is exactly what U4 demands. `apps/web` is explicitly out of this milestone's scope, and the telas are M-06's seam. Escalated to the hub: U4 closes on M-06, not here. |

## Round-1 dual gate (against `84323d0a`)

| Pass | Verdict | What it found |
|------|---------|---------------|
| Cold gate (independent, contract in hand) | **FAIL** | Blocker: "auto-approve silently reverses an operator REJECTION". Plus a stale blocker — "no evidence pack" — which is an artifact of timing: this file was written after the gate read the tree. |
| Adversarial refuter | **REFUTED** | Five defects, listed below. |

Four were code defects and are fixed in `5807d634`; each guard is proven load-bearing by a
must-fail revert (table above). The fifth is a screen gap outside this milestone.

1. **Auto-approve reopened settled listings.** The precedence guard read only the E10 trail, but `RejectListing` writes no decision row — there is no rule to record for "this anúncio is not ours" — and links approved before 0082 have none either. Both read as undecided, so the next sync flipped them to `actor=system`. Fixed by reading the link's own terminal state, which covers both cases. (Cold-gate blocker and refuter defect 1 are the same hole from two angles.)
2. **The trail named an anchor that had not resolved anything.** Approving a collision or conflict candidate recorded `exact_ean_unique` / `exact_codprod_unique`. Fixed: only ACCEPT/CONFIRM may name their anchor; everything else is `manual`.
3. **Undo left no trace.** The reverted decision still read as in force, so the listing could never be auto-linked again and nothing said why. Fixed: undo writes its own decision, superseding what it took back.
4. **`NO_CANDIDATE` had become unreachable**, breaking `/vinculos`' honest "sem candidato" row and flipping its batch-select guard to enabled. Reverted; see M05-C5 for the reconcile request.
5. **`/vinculos` does not group CONFIRM apart from REVIEW** — real, but `apps/web` is out of scope and the telas are M-06's seam. Closes as M05-U4 BLOCKED, escalated.

## Round-2 dual gate (against `a4a0ad89` = code `5807d634` + this pack)

| Pass | Verdict | What it found |
|------|---------|---------------|
| Cold gate | **FAIL** | 4 blockers. Explicitly confirmed the round-1 blockers are closed and the guards load-bearing; found no new integrity defect in the code. |
| Adversarial refuter | **REFUTED** | 3 defects + NITs, all in or left standing by the corrective slice. |

Both gates independently reached the same doubt about the L2 rung — that is the finding of this
round, and it was correct. Disposition:

| Raised | Disposition |
|--------|-------------|
| Cold **B3** / refuter ladder: the integration lane's `status=passed` does not show the E10 tests ran | **Proven, not argued.** Must-fail above; the lane names the broken test. Lane artifact itself is a FINDING for the profile. |
| Cold **B1**: pack names a tip that is not the branch tip | **Fixed** — header now names the code tip and the docs commit separately. |
| Refuter **D1**: a rejection writes no decision, so the `actor=system` row it overrules stays in force | **Fixed** in `f9dbc868`. Round 1 fixed this for undo and left it open for reject; the reasoning that excused it ("no rule to record") was wrong — `superseded_by` is a fact, not a rule. `TestRejectingAListingSupersedesTheDecisionItOverrules`. |
| Refuter **D2**: an ACCEPT candidate approved by hand is filed as a single-anchor rule | **Fixed.** Under-claiming the evidence is the same defect as over-claiming it. `TestOperatorApprovingACorroboratedCandidateRecordsCorroboration`. |
| Refuter **D3**: `CollisionsAtDecision int` forces an unread count to 0 (AC-03) | **Fixed** — `*int` from the generator through to the row. `TestAnAutomaticApprovalWithoutACountRecordsNoCount`. |
| Cold **NB-1**: a CONFIRM candidate appears in no aggregate counter | **Fixed** — `summary_reader.go` now counts the confirmation queue as pending. New L2 test `TestPendingLinksCountsTheConfirmationQueue`. A *separate* CONFIRM counter would be a new API field and is M-06's to render; reported, not built. |
| Cold **NB-2**: one failing auto-approval abandons the rest of the batch | **Fixed** — `errors.Join`, all attempted, failures still reported. |
| Refuter NIT: corrective #3's commit message overstates its effect | **Accepted, recorded here.** The undo's own row still blocks re-auto-linking (correctly — an operator taking a link back should not have it re-created by the next sync). Only the trail improved; the message claims more. |
| Cold **B2** (M05-C5 wording) and **B4** (M05-U4) | **Hub adjudication.** Both gates agree the chip's argument is factually right and that a chip may not ratify an amendment to its own contract. Unchanged; see the ledger. |
| Refuter NIT: TOCTOU — the precedence guards read outside the transition's transaction | **Recorded as a FINDING**, not fixed. The scheduler and a live generation call can interleave and both pass the guards, yielding two `actor=system` rows for one link (the second superseding the first). Both rows are corroborated and agree, so the trail stays honest, but the fix belongs in the repository layer. Hub call. |
| Refuter: EVIDENCE citation drift (4 of 8 sampled imprecise) | **Fixed below**, and recorded as a finding about how this pack was written. |

**Round 3** runs a focused refuter over `git diff 5807d634..f9dbc868` — the corrective slice is
where new defects enter, and per the D-120 lesson correctives never close on self-verification.

## Deviations reported to the hub

1. **OpenAPI + SDK enum widened** (`marketplace-central.openapi.yaml:6459`, `sdk-runtime/src/index.ts`). M05-C21 requires CONFIRM and REVIEW to be distinct *at the transport read*, so the state must be in the contract; the repo rule that spec and generated SDK land together made this a same-changeset edit. No new endpoint, no FE change, no other enum touched (the market-domain enums at 3447/7618/7777 are deliberately untouched).
2. **`root.go` reordered + refresher wired** — the ownership matrix predicted "no new composition wiring" for M-05. F-01 cannot exist without it: the resolution service is now built before generation so it can be the `AutoApprover`, and the products scheduler receives `WithLinkCandidateRefresh(...)`. Both are additive; the scheduler option is optional and absent it the job behaves exactly as before.
3. **Migration inventory fixture 68 → 69** (`internal/platform/migrate/runner_test.go`) — owed by the F-03 commit, caught here by the full test ladder.
4. **FINDING (pre-existing, fixed in F-04)** — a map-iteration-order flake in the `candidateRows` fixture helper: the test asserted on the first row of a map-ordered slice, so it passed or failed by hash seed. Now deterministic.
5. **FINDING (tooling)** — the integration lane requires pwsh (PS7); Windows PowerShell 5.1 reports `status=blocked` with a `RandomNumberGenerator.Fill` error that reads like a lane failure and is not one. Candidate for profile §3's false-alarm list.
6. **FINDING (method)** — a must-fail proof executed as `git checkout -- <file>` on a file carrying uncommitted work silently destroys that work; the revert and the restore look symmetric but are not. Must-fail proofs belong in one edit-and-restore command, or the slice gets committed first. Candidate for the core §5 verification-ladder notes.
7. **FINDING (method)** — write the evidence pack **before** dispatching the gate. Round 1's cold gate raised "no evidence pack" as a blocker purely because this file did not exist yet when it read the tree; that costs a whole gate round.
8. **FINDING (tooling, high value)** — the integration lane cannot be closed on `summary.txt` alone: it runs `go test -tags=integration` without `-v`, every DB test opens with `SkipWithoutTarget`, and the artifact records only `target`/`status`/`run_id`. A fully-skipped run and a fully-green run are byte-identical in the evidence. Both round-2 gates flagged it independently. Until the lane records executed/skipped counts, an L2 claim needs the must-fail proof shown above. Candidate for the profile's ladder bindings.
9. **DEVIATION (round 2)** — `summary_reader.go` now counts the CONFIRMAÇÃO queue as pending work. Not named in the write-set plan, but the queue is M-05's own creation and it was invisible to the only number an operator reads to know there is work to do. A *separate* CONFIRM counter would be a new API field and belongs to M-06 with the screen that renders it; reported, not built.
10. **FINDING (not fixed — hub call)** — TOCTOU on the precedence guards: `GetProductLink`/`ListDecisionsForLink` read outside the transaction `ApplyProductLinkTransition` opens, so the 15-minute scheduler and a live `POST …/generations` can interleave and both pass, writing two `actor=system` rows for one link (the second superseding the first). Both rows are corroborated and agree, so the trail stays honest and no wrong link is created — but the fix belongs in the repository layer, not in this milestone's write-set.
11. **OBSERVATION** — `product_links/domain/product_link.go` is not gofmt-clean at BASE-SHA (struct alignment in `ProductLinkAuditEntry`). Pre-existing, verified against `e3c081ae`; left alone rather than reformatting a file this chip only touched a doc comment in.
12. **RETRO (mine)** — the round-2 refuter's D1 was a hole I had reasoned myself into: round 1's fix wrote a decision for undo and explicitly declined to for reject, on the ground that "there is no rule to record". `superseded_by` is a fact, not a rule, and I had written a test (`TestRejectingAListingWritesNoDecision`) that encoded the wrong belief — so it could never have caught it. A test written to confirm a decision is not evidence about that decision.
