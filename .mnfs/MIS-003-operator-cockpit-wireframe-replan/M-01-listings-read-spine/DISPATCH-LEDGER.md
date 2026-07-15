# M-01 Dispatch Ledger

Chip session: M-01-listings-read-spine. Branch `mis-003/m-01-listings-read-spine` off base SHA `d0d30d68` (rebased from 1eb99600 after hub committed the harness amendment; docs-only fast-forward).
Isolation: dedicated git worktree `.claude/worktrees/m01-listings` (session re-rooted). Established after the shared-tree hub commit clobbered HEAD onto `main` — chip now isolated, hub owns main.
Harness amendment (operator-ratified 2026-07-15): planning is batch up-front; all codex dispatches via `/codex:rescue --model <m> --effort <e> --wait`; raw `codex exec` only for the precondition probe.

## Precondition
- Codex probe: `codex 0.144.4` present; `CODEX_OK` (gpt-5.6-sol) + `LUNA_OK` (gpt-5.6-luna) returned. Sandbox healthy.

## P2 PLAN

| # | Feature | Role | Model / effort | Path | Result |
|---|---|---|---|---|---|
| D1 | F-01 | Feature planner | gpt-5.6-sol / medium (raw `codex exec` — pre-amendment) | evidence/planner-sol-medium.log → plan.md | DELIVERED — 5 slices, real file:line anchors, 15 open risks; surfaced 2 blocking contradictions |

### D1 findings triage
- **BLOCKING (hub):** migration block 0033–0035 already consumed at base; connectors capability seam lacks price/currency/modality (IC-02-required) → cross-module edit outside listings lane. → BLOCKED event sent to hub.
- **Milestone-session ratified** (documented, non-blocking): price_currency nullable; 409 active-run id at `error.details.operation_run_id`; duplicate-canonical-key → fail-honest before persistence; no cross-module FK (installation existence enforced at app layer); modality allowlist authored here; known-but-not-connected installation → new error case (scoped connected-only for M-01 tests).
- **ESCALATION (out-of-scope):** ADR-12/ADR-17 have no formal record under docs/architecture/decisions (behavior unambiguous; architecture owner repairs).
- **Deferred:** stale queued/running run crash-recovery (manual cleanup for M-01).

| D2 | F-02 | Feature planner | gpt-5.6-sol / medium (streamed to live log; two stuck restarts, see note) | F-02/plan.md | DELIVERED — 32.8KB/219 lines. 7 slices, 7 shared-seam reqs for F-01, real file:line anchors, 5 blockers surfaced (R1–R5) |

### D2 findings triage (post-plan)
- **OPERATOR-RATIFIED:** R2 below_margin formula → contract-literal (DECISIONS D-16).
- **Milestone-ratified:** R3 exception precedence/group severity (D-17); R5 keyset stability (D-18); R4 timeline event source `listing_sync_events` folded into F-01 (D-19).
- **ESCALATION → HUB:** R1 Oracle-cost vs C09 single-query summary (D-20). Blocks only Slice 5 + below_margin counter; Slices 1–4/6/7 proceed.
- **Live-visibility infra:** codex now streamed to a tailable log + local dashboard (http://127.0.0.1:7391) rendering codex events as chat cards. Native task panel does NOT stream codex (wrapper hides child proc → shows "Parado"); log+dashboard is the real live view. Two planner restarts: first raw `codex exec` hung 3.7h (no reasoning output); relaunch via /codex:rescue got interrupted; final clean raw-exec-with-live-log run delivered plan.md at 11:59:50.

### Process deviation (recorded, one-time)
Harness §4.1 (amended) requires ONE batched planner pass per milestone covering ALL features + shared seam. M-01 ran TWO passes (D1 F-01, D2 F-02). Cause: D1 (F-01 planner) was dispatched BEFORE the batch-once amendment was ratified (old one-feature-at-a-time model). Post-amendment the correct move was a single combined pass; instead F-02 was planned separately. Sunk cost — F-01 already planned; a combined re-plan would waste D1. Coherence recovered by feeding F-01/plan.md + DECISIONS.md into D2 as inputs (F-02 designs against F-01's ports; F-02's Shared-Seam-Requirements folded back into F-01 Slice 3). M-02..M-06: single batched planning pass, no repeat.

## Hub rulings applied
- Migration block 0036–0037 granted → 0036_listings.sql (see DECISIONS D-1).
- Connectors capability seam contract-lock (option a), additive-only + CONNECTORS_PROVIDER_AUTH (see DECISIONS D-3). Lock ends at CLOSED; connectors diff called out in CLOSED payload.

## Base / isolation
- Rebased 1eb99600 → d0d30d68 (hub harness amendment, docs-only ff).
- Isolated worktree .claude/worktrees/m01-listings after shared-tree clobber onto main.

## P3 IMPLEMENT
- Starts at F-01 Slice 1 once D2 (F-02 plan) lands and Shared-Seam-Requirements are folded into F-01 repository design. Slices: Luna-high (standard) / Sol-low (complex), failing-test-first, commit per green, independent sonnet review per slice before next.

| # | Feature/Slice | Role | Model / effort | Log | Result |
|---|---|---|---|---|---|
| I1 | F-01 Slice 1 (schema+domain+modules.json) | Implementer | gpt-5.6-luna / high (direct `codex exec`, live log) | scratchpad/agent__f01-slice1.log | **GREEN — committed 746a97d4**. 0036_listings.sql (listings + listing_sync_events per D-19), domain listing.go, RED-first migration+domain tests, modules.json register. Tests+build green (re-run by milestone owner in worktree), governance-validate passed, no prefix-exception. Sonnet review clean (1 finding = false positive: `provider` col IS IC-02-mandated, line 31/134). Governance-composition gap resolved via D-23. |

| I2 | F-01 Slice 2 (connectors seam + ML mapper) | Implementer | gpt-5.6-luna / high (direct `codex exec`, live log) | scratchpad/agent__f01-slice2.log | **GREEN — committed 9df16851** (+ review fixup 2167dbb7). Additive connectors ListingSnapshot fields (price/currency/listing_type, D-3) + ErrCodeProviderAuth (401/403 split); ML decode via json.Number (no float); listings mapper ML→canonical (D-5 allowlist, ADR-17 nil-not-zero, '-' expansion). Tests+build green (re-run by milestone owner), governance-validate passed, existing connectors tests untouched, no ML-type leak, no product_links. Sonnet review: 1 finding (test-isolation defect in tenant reject case) FIXED in 2167dbb7 — production mapper confirmed to reject blank tenant_id. |

Note: D-19 folded — `listing_sync_events` table created in the SAME 0036 migration in Slice 1 (schema only; write logic stays Slice 3).

### Live visibility (multi-agent)
Dashboard rebuilt as multi-agent (http://127.0.0.1:7391): sidebar lists every `scratchpad/agent__<id>.log`, operator picks which worker's chat stream to view. Each dispatch writes its own `agent__<id>.log` + `agent__<id>.done` sentinel on exit; server `/agents` reports state live|idle|done. Current agents: f01-slice1 (live), f02-planner (done).

| I3 | F-01 Slice 3 (atomic upsert-and-close ingestion) | Implementer | gpt-5.6-sol / low (complex; direct `codex exec` OS-process, live log) | scratchpad/agent__f01-slice3.log | **BLOCKED — zero code written.** Worker stopped BEFORE implementing on a verification/ownership conflict: registered integration lane `scripts/harness/Postgres.psm1:119` hardcodes its package list and omits `./internal/modules/listings/adapters/postgres`, so the tagged integration test can't execute without a `harness-control` shared-seam edit (outside the 7 permitted files). Disciplined-but-OVER-conservative stop (should have written all 7 files incl. the tagged test + green unit tests, deferring only lane execution). Evidence in slice3-notes.md. Re-dispatch pending hub §3 direction (decoupled prompt: write tagged test inside the 7 files; defer lane execution to milestone owner). |

## HUB RULINGS applied (2026-07-15) — all six asks resolved
Hub ruled the escalation (see proposal). Applied on the chip: (1/2) hub executes skill fixes + un-guards the orchestration test + queues dispatch_preflight migration — nothing on chip plate; permanent rule adopted: **every dispatch prompt pins HARNESS.md + forbids mpc-goal-harness** (done in both prompts below). (3) DELEGATED lane append — `./internal/modules/listings/adapters/postgres` added to `scripts/harness/Postgres.psm1:119` by milestone owner (harness-control touch under delegation, **flag in CLOSED**). Slice 3 UNHELD. (4) RATIFIED: OS-process codex dispatch (stdin-closed, teed log, .done) completes to the dispatching loop → allowed intra-milestone incl. backgrounded; Agent/Task nested children stay SYNC-only. (5) dashboard adopted as HARNESS §8 standard. (6) SPLIT option A: D-24 fanned out as a 2nd OS-process worker INSIDE this chip (disjoint files, one ledger, no sibling worktree). F-01 internals + F-01→F-02 serial ACCEPTED as data-DAG-mandated.

| I3-retry | F-01 Slice 3 (re-dispatch, lane fixed, HARNESS-pinned) | Implementer | gpt-5.6-sol / low (complex, OS-process bg) | scratchpad/agent__f01-slice3.log (task bou52vp3o) | RUNNING — concurrent with I4 per §4.1. Build/test scoped to listings subtree (sibling active); full ./... build deferred to owner at integration. |
| I4 | D-24 internal_read `GetICMSCeilingByOrigin` (additive ceiling read) | Implementer | gpt-5.6-luna / high (standard, OS-process bg) | scratchpad/agent__d24-icms-ceiling.log (task bu7gn8oam) | RUNNING — concurrent with I3. Additive: domain ICMSCeiling + port method + oracle adapter (specialist query verbatim, current-config no as-of) + fake stub + REDBASE-trap NEGATIVE test. Disjoint from listings. |

| I5 | F-02 Slice 1 (filter grammar + read model + opaque cursors) | Implementer | gpt-5.6-sol / low (complex, OS-process bg) | scratchpad/agent__f02-slice1.log (task bnytu93bu) | RUNNING — 3rd concurrent worker. Creates 7 NEW files (domain/read_model, ports/read_repository+cursor, transport/query) disjoint from I3 (listings/adapters+application) and I4 (internal_read). NO compile-assertion in read_repository.go (would bind Slice 2's repository.go). Build scoped to domain+ports+transport packages only (siblings mid-write in adapters/postgres). Banks F-02 foundation. |

### Concurrency (3-way, §4.1 ratified)
Three OS-process codex workers running in ONE worktree, file-disjoint, one ledger, no sibling worktrees: I3 (listings ingestion, adapters/postgres+application+ports/ingestion.go), I4 (internal_read ceiling), I5 (F-02 read foundation: domain/read_model+ports/read_repository+cursor+transport/query). Each build/test SCOPED to its own packages to avoid compiling siblings' half-written files; the integrated `go build ./...` is the milestone owner's at integration. DAG rationale (operator Q): F-02 S2-5 all edit the SAME repository.go+read_service.go → serial-mandated (one-writer-per-seam); only file-disjoint units parallelize. Next parallel wave candidate: F-01 Slice 4 (refresh/wiring, needs integrations lock) ∥ nothing in listings (S2 collides on repository.go). Dashboard shows all three streams.

## Harness-control ESCALATION (2026-07-15) — old→new migration + field standards
Slice 3 worker auto-invoked the SUPERSEDED `mpc-goal-harness` skill (root cause: new harness ships as a DOC + missing `.claude/skills/harness-hub`, NO worker skill → discovery resolves to old `.agents/skills/mpc-goal-harness`). Operator directed full migration + standardization. Package: harness-migration-and-standards-proposal.md. ESCALATION sent to hub `local_efa46c30` with 6 asks: (1) worker-facing skill; (2) migrate `dispatch_preflight.py` collision/one-writer validation concept, retire /goal prose, un-guard the orchestration test, keep shared psm1; (3) authorize Postgres.psm1 listings lane append + auto-discovery; (4) ratify OS-process-dispatch safety clarification (completes to dispatching loop, not the nested-child→hub deadlock — verified task b1umx1rds); (5) adopt multi-agent dashboard as standard; (6) SPLIT-REQUEST for D-24 internal_read ceiling method — intra-chip fan-out (if §4.1 ratified) vs sibling worktree. All paths are `harness-control` seam (shared-seams.json:10) → chip cannot self-edit. **HOLDING for hub ruling per operator.**

### Below_margin / ICMS (parallel, non-blocking F-01)
ICMS consultation SENT to DB specialist `local_ec787804` (3 questions per hub route-a ruling — see cost-read-verification.md). below_margin (F-02 Slice 2 enrichment + Slice 5 counter) parked pending reply. F-01 Slice 1 has no cost dependency → proceeds.

## Wave-1 acceptance (2026-07-15): I3 / I4 / I5 all CLOSED green
Three concurrent workers completed exit 0; each independently reviewed (sonnet cavecrew-reviewer), findings fixed by milestone owner, verified.
- **I5 = F-02 Slice 1** (bnytu93bu, 69.8k tok): 7 files. Review 0🔴 3🟡 — all fixed: (a) present-but-empty filter param (`filter.status=`) now rejected as invalid_filter instead of aliasing to default (query.go filterScalar + 3 regression cases); (b) unknown-filter-key error now names bare key not `filter.x` prefixed; (c) trivial-coverage below_margin test rewritten to honest zero-value=nil assertion. `go test domain/ports/transport` GREEN, build GREEN.
- **I3 = F-01 Slice 3** (bou52vp3o, 150.3k tok): 7 files (ingestion port/app/connector/postgres + D-19 sync-event ledger). Review 0🔴 1🟡 — fixed: unbounded provider page loop → added `maxIngestionPages=10_000` honest ceiling (fails loudly, no silent truncation) + regression test. Reviewer confirmed D-19 kind derivation exact (closed/paused/synced), single-txn rollback proven, tenant+installation scoping on every SQL, ErrDuplicateCanonicalKey pre-persist, provider payload stays at connector. Unit GREEN; **hermetic integration GREEN** (see below).
- **I4 = D-24 internal_read ceiling** (bu7gn8oam, 239k tok): domain/icms_ceiling.go + BatchReader.GetICMSCeilingByOrigin + oracle adapter + REDBASE-trap negative test + cache/legacy forwarding. Review 1🔴 — FIXED: adding the interface method broke a compile-assertion in `tests/unit/cache_composed_test.go:243` (composed double missing the method); added stub. Reviewer confirmed query verbatim (ALIQUFDEST+NVL(MAX(PERCICMSFCP),0), GROUP BY UFDEST, UFORIG bind, no date filter), negative test asserts absence of ALIQUOTA/REDBASE/ALIQICMS, nil discipline correct, additive-only. `go build/test ./internal/modules/internal_read/...` GREEN + `go vet ./tests/unit/...` GREEN.

## Integration lane — root-caused + PRE-EXISTING sibling breakage (2026-07-15)
Slice 3's registered-lane block (`HPG_MIGRATION_FAILED`, `migrations_first=-1`) FULLY diagnosed and it was **NOT a listings/migration defect**:
- **Root cause = cold hermetic module cache.** The integration/unit lanes run Go with `GOPROXY=off GOSUMDB=off GOMODCACHE=apps/server_core/.gomodcache` (Environment.psm1:180-184). A fresh worktree's `.gomodcache` is empty → `go run ./cmd/testdb migrate` can't fetch `godror` etc. → build exit 1 → harness maps ANY nonzero migrate exit to `HPG_MIGRATION_FAILED` with count -1 (the -1 is "no `applied N` line", i.e. build failed before running — NOT a SQL error). Proven: `GOPROXY=off go build ./cmd/testdb` → `godror@v0.51.0: module lookup disabled by GOPROXY=off`.
- **Fix (env prep, NOT a dep change / NOT a purge):** `cd apps/server_core && GOMODCACHE=$(pwd)/.gomodcache go mod download all` (~129M). After warm: hermetic offline build GREEN; registered lane advanced to `migrations_first=36 migrations_second=0` (apply + idempotent PROVEN). Recorded to memory [[hermetic-integration-lane-modcache]] as a field standard for the hub (fresh-worktree onboarding gap).
- **Residual = PRE-EXISTING broken sibling test, out of scope.** With cache warm the lane now fails `HPG_TEST_FAILED` on `tests/integration/TestPhase1SmokeFlow` (`PRICING_INVALID_PRODUCT_ID`, phase1_smoke_test.go:75). Fails **in isolation on a clean DB (0.34s, deterministic)** and **with the original 4-package set (no listings)** — reproduces with zero of my changes; phase1 has no internal_read/ceiling refs, so D-24 is not implicated. This is a branch-level/other-module defect → **hub/one-writer territory, chip must not fix** (orders/pricing/tests-integration seam).
- `orders/.../TestOrderRepositoryDuplicateIdentityGroup...` failed once in the 5-package concurrent run, PASSED in the 4-package run → pre-existing shared-DB test-isolation flakiness surfaced by adding a 5th concurrent package; listings uses disjoint tables (listings, listing_sync_events) so no logical coupling. Flag to hub.
- **Listings deliverable is GREEN in the lane:** `internal/modules/listings/adapters/postgres` `ok` both isolated (7.8s) and in the full 5-package set; profitability + product_links `ok`. Manual hermetic replication used own postgres:16-alpine, DB name `mpc_test_<32hex>` (LoadConfig target.go rule), applied 36 migrations.

**For hub (CLOSED will carry):** (1) adopt cache-warm as standard fresh-worktree step or pre-warm in worktree bootstrap; (2) pre-existing `TestPhase1SmokeFlow` failure blocks any full-integration-lane green — needs hub owner/other-milestone fix; (3) orders integration test isolation weakness; (4) Postgres.psm1 listings lane append (already delegated) now verified correct.

## HUB RULINGS (2026-07-15) — A–D + 6 asks, HOLD LIFTED
Hub (local_efa46c30) ruled the escalation. **A) ADOPTED** — cache-warm is doctrine, ratified HARNESS §5 (hub commit 4c1f22c): standard first act in every new chip worktree before any hermetic lane; `migrations_first=-1` = build died before migrate. **B) ROUTED to hub** — TestPhase1SmokeFlow is hub-owned (post-M-01 corrective chip); M-01 acceptance criterion ADJUSTED: full-lane green NOT required — required = (1) green on all M-01-owned packages in the hermetic lane, (2) ledger evidence that phase1 fails without our changes → see evidence/phase1-preexisting-repro.md (4-package clean-DB repro + isolation repro + D-24 non-linkage). Do not fix. **C) ACK** — orders flake = hub board task, bundled with B's chip. **D) CHIP COMMITS go.work.sum** alone (hub-delegated single write): "chore(deps): normalize go.work.sum via go mod download (checksum prune only, no version changes)". **6 asks:** (1) mpc-goal-harness DELETED on master (fa7c6fcc) + worker skill .agents/skills/harness-worker (1e7df7cc) + .agents/skills/codex-dispatch (b6b02513) — keep interim pin verbatim until merge (worktree predates these); (2)/(3) auto-discovery + dispatch_preflight migration DEFERRED post-M-01 (hub board); (3) e44d74c7 lane append RATIFIED; (4) OS-process dispatch safety in HARNESS §3 + codex-dispatch skill; (5) dashboard ADOPTED standard; (6) MOOT (D-24 intra-chip §4.1 correct). **HOLD LIFTED — proceed F-02 Slice 2 after D.**

### B evidence line (hub-required)
`go test -tags=integration ./tests/integration ./internal/modules/orders/adapters/postgres ./internal/modules/profitability/adapters/postgres ./internal/modules/product_links/application -count=1` on a clean migrated ephemeral PG (NO listings in set) → `--- FAIL: TestPhase1SmokeFlow (0.44s) phase1_smoke_test.go:75: PRICING_INVALID_PRODUCT_ID`; orders/profitability/product_links `ok`. Isolation `-run TestPhase1SmokeFlow` → FAIL 0.34s deterministic. Full capture: evidence/phase1-preexisting-repro.md. Proves the failure is pre-existing and independent of M-01.

## Wave-2 sequencing (2026-07-15) — HOLD lifted, D-21 prerequisite surfaced
Proceeding to F-02 Slice 2 (tenant-scoped GET /listings keyset + below_margin worst-case enrichment). Pre-dispatch seam audit (milestone owner, read-only) resolved every enrichment input EXCEPT one, and surfaced a hard prerequisite:
- **Cost basis**: internal_read `BatchReader.GetCostFactsByIDs(ctx, []int64) map[int64]*domain.CostAsOf`; `CostAsOf.Amount *float64` = CUSSEMICM (nil = unknown → below_margin unevaluable). Confirmed present.
- **ICMS ceiling**: internal_read `GetICMSCeilingByOrigin(ctx, domain.UF) map[domain.UF]*ICMSCeiling` (D-24, landed 176d8082). Worst-case scalar = MAX ceiling over the returned UFDEST map.
- **Origin UF (:uforig)**: PINNED — D-22 line 104 fixes `UFORIG=13` (MG, the CODEMP=1 company state). Read-service passes it as a documented constant; empty ceiling map → unevaluable/nil (ADR-17).
- **min_margin**: `marketplaces.domain.Policy.MinMarginPercent` (float64) — but the published read method **D-21 `GetPricingPolicyForInstallation` DOES NOT EXIST YET** (git grep: referenced only in DECISIONS.md/plan.md, absent from marketplaces code). Slice 2's below_margin enrichment consumes it → **D-21 is a hard prerequisite**.
- **root.go**: plan's "inject into existing F-01 handler" is STALE — root.go has NO listings wiring yet (F-01 built domain/adapters/app/ports but never composed; no handler exists). **Reconciled deviation:** Slice 2 will NOT touch root.go; all composition wiring deferred to Slice 6 (routes+OpenAPI+SDK, one commit). Read stack is exercised via constructor injection in unit+integration tests. No route, no benefit to half-wiring root twice.

Decision: dispatch **D-21 first** (marketplaces-only, file-disjoint, additive lock already GRANTED by hub — DECISIONS D-21 — so standing authorization, no new REQUEST), then F-02 Slice 2 after D-21 lands+reviews. Mirrors the D-24 additive-lock sub-worker pattern. Slices 2–5 remain serial (one writer on repository.go/read_service.go) — no parallel listings dispatch.

| # | Feature/Slice | Role | Model / effort | Log | Result |
|---|---|---|---|---|---|
| I6 | D-21 marketplaces `GetPricingPolicyForInstallation` (additive published policy read, tenant-scoped via marketplace_accounts.integration_installation_id) | Implementer | gpt-5.6-luna / high (standard, OS-process bg) | scratchpad/agent__d21-policy-read.log (task b6qtngabk) | RUNNING — additive: Repository port method + postgres SQL + Service method + fakeRepository stub + RED-first tests. Existing marketplaces tests untouched-green. Feeds F-02 Slice 2 below_margin (min_margin input). |

## HUB DOCTRINE UPDATE (2026-07-15) — P7 gate replaced (HARNESS.md master commit 320d4a2)
Applies at P7 (not yet reached; currently P3). Read docs/superpowers/HARNESS.md P7 + §5 L4 from MASTER at close (worktree copy predates 320d4a2).
- **P7 = MNFS milestone gate**, replacing ad-hoc fresh-browser QA persona: run `/milestone-validate <milestone-path> --apply` → dispatches independent cold `milestone-reviewer` crew + `qa-validator` live browser drive against validation-contract.md → writes `<milestone-root>/validation-result.md`. ONLY that verdict passes the milestone.
- On Fail: `/correction-create <milestone-path> <report> --apply` scopes; milestone owner (me) dispatches corrective worker (codex per codex-dispatch matrix); full re-gate, never-downgrade.
- **Unchanged:** P6 dual gate (Opus + Sol medium, fixed SHA) STILL required BEFORE P7. Evidence paths unchanged. Rulings A–D unchanged.
- **Role binding** (mnfs-plugin/docs/shared-standards.md §Role Binding): I am Milestone Orchestrator; workers = Feature Implementer / Correction Worker. Deleted upstream — NEVER invoke: milestone-orchestrator/feature-implementer/correction-worker agents, /milestone-start, /feature-context, /feature-accept.

### I6 result — D-21 GREEN, committed 1b644ed7 (106.2k tok)
Additive `GetPricingPolicyForInstallation` on marketplaces: Repository port method + Service method + fakeRepository stub + RED-first service tests (blank-guard / found / not-found / repo-error via errors.Is) + postgres SQL. SQL tenant-scoped on BOTH sides (`a.tenant_id=$1` AND join `p.tenant_id=a.tenant_id`), bound params, `LIMIT 2` + `ORDER BY policy_id` duplicate detection, blank installationID short-circuits pre-query. Review (sonnet cavecrew-reviewer) 0🔴 1🟡 1🔵 — 🟡 FIXED by milestone owner: duplicate-policy bare `errors.New` → `ErrMultiplePoliciesForInstallation` sentinel (errors.Is-able, code string preserved, fail-honest ADR-17). 🔵 (inner JOIN semantic) intentional, no change. `go build`+`go vet -tags=integration`+`go test ./internal/modules/marketplaces/...` GREEN; governance passed. Existing marketplaces methods/tests untouched. Unblocks F-02 Slice 2 below_margin (min_margin input now available).

## F-02 Slice 2 — BLOCKED (contract conflict), ESCALATED to hub (2026-07-15)
| # | Feature/Slice | Role | Model / effort | Log | Result |
|---|---|---|---|---|---|
| I7 | F-02 Slice 2 (keyset GET /listings + below_margin) | Implementer | gpt-5.6-sol / low (complex, OS-process bg) | scratchpad/agent__f02-slice2.log (task bs8mvgq0a, 82.8k tok) | **BLOCKED — zero source changed (honest stop).** Real over-constrained conflict: `exception=below_margin`/exact `has_exception` need below_margin (Oracle cost per row), but D-20 forbids a PG cost projection + D-18 mandates keyset limit+1 + IC-02 needs the exact filter — can't hold all four. |

Milestone-owner classification (slice2-escalation.md): block is SMALL — only the below_margin-**dependent filter predicates**. Everything else (keyset list, direct filters, sync_error/stale/unlinked exception filters, q, below_margin per-row DISPLAY via per-page Oracle batch, cursor, IC-02 response) is implementable with ZERO contract change. Options: (1) defer below_margin list filter to later milestone; (2) **iterative bounded keyset scan** for below_margin-filtered queries only — preserves D-18/D-20/IC-02, costs multi-batch on sparse filters [OWNER RECOMMENDATION]; (3) authorize PG cost projection (reverses D-20). ESCALATION sent to hub local_efa46c30; HOLDING for ruling — no redispatch until ruled. D-21 (min_margin seam) already landed 1b644ed7.
