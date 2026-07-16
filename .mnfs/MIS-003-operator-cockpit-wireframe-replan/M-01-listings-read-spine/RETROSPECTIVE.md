# M-01 Pilot Retrospective → HUB

**Milestone:** M-01-listings-read-spine (MIS-003). First milestone run under the hub-and-chips harness — the pilot.
**Author:** M-01 milestone session (Opus).
**Status:** F-02 COMPLETE (Slices 2–7 CLOSED); F-01 Slice 4/5 pending. Retrospective written before close so its lessons apply to M-02..M-06 planning, not after.
**Purpose:** The pilot delivered, but slowly and with too much avoidable friction. This is the honest post-mortem: what cost wall-clock, what was old-harness residue, what redundancy accumulated, and the concrete harness changes that must land before the next milestone so execution is fast and low-error.

---

## 1. Cost summary (the pain, quantified)

Pulled from DISPATCH-LEDGER.md + script inventory. These are the numbers that matter:

| Metric | Count | Note |
|---|---|---|
| Hub round-trips (BLOCKING / ESCALATION / REQUEST) | **~8** | Dominant wall-clock cost — each is a synchronous stall on an idle hub session |
| Slices dispatched then BLOCKED with zero code (wasted dispatch) | **2** | Slice 2, Slice 4 — both contract conflicts, both plan-time-detectable |
| Prerequisites discovered at implementation, not planning | **2** | D-21 (`GetPricingPolicyForInstallation` didn't exist); root.go had no listings wiring (plan assumed it did) |
| Planning passes (should be 1) | **2** | D1 (F-01) ran pre-batch-amendment, then D2 (F-02) separately — process deviation, partial sunk cost |
| Hand-rolled one-off scripts (scratchpad) | **4** | 3 PowerShell lane/dispatch runners + 1 Node dashboard, all re-implementing harness capability |
| Distinct test-DB bring-up code paths | **4** | 3 of them redundant hand-rolls of the one canonical lane |
| Recurring infra failures that masqueraded as test failures | **2** | cold modcache → false `HPG_MIGRATION_FAILED`; pg_isready first-boot race → `3D000 database does not exist` |
| Reviewer findings dismissed as false positives | **≥5** | Slice 1 (1/1 FP), Slice 4 (3/3 all FP), Slice 7 (several 🟡) — each cost owner adjudication time |
| Dispatch-channel doctrine reversals | **3** | rescue-only → switch-to-rescue → "actually raw-exec is the panel-visible one" |
| Old-harness (`mpc-goal-harness`) reference sites still on disk | **8+** | Including a test that HARD-ASSERTS the dead skill exists |

**Root read:** the milestone was slow not because the work was hard — the code shipped clean (p95 3.26ms, zero real correctness escapes) — but because **the harness was still half-migrated, the plan wasn't built for the parallel model, and contract conflicts were pushed to implementation time.** Every one of those is fixable at the harness/process layer.

---

## 2. Findings (grouped, most-costly first)

### A. Plan was authored for the OLD (serial, feature-at-a-time) harness — not the parallel model
The plan came out feature-at-a-time (D1 F-01, then D2 F-02), because D1 was dispatched *before* the batch-once amendment ratified. The deeper problem isn't the double pass — it's that **the plan never produced the artifacts the parallel harness needs to go fast:**
- No **file-level write-DAG** up front. Which slices touch `repository.go` / `read_service.go` (→ serial, one-writer) vs which are file-disjoint (→ parallel) had to be re-derived slice by slice. Only ONE parallel wave ever happened (I3/I4/I5); everything after was strictly serial because slices 2–7 all write the same two files.
- The **additive-lock sub-worker pattern** (D-21, D-24 — split a disjoint additive dependency out as its own concurrent worker) was invented mid-flight, not planned. It's the single best parallelism lever we found and it was reactive.
- **No contract-reconciliation pass.** The plan listed decisions but never cross-checked them for mutual satisfiability. That's what caused §B.

> **The plan told us WHAT to build. It never told us what could run in parallel or where the contract contradicted itself. In a parallel harness, that's the plan's whole job.**

### B. Contract conflicts discovered at implementation → 2 wasted dispatches + 2 escalations
- **Slice 2 BLOCKED:** `below_margin` filter over-constrained by D-20 (no PG cost projection) + D-18 (keyset limit+1) + IC-02 (exact filter) — four requirements that can't all hold. Worker honestly wrote zero code, stopped, escalated. Hub ruled Option 2. Redispatch.
- **Slice 4 BLOCKED:** group-key pagination not covered by the Slice-2 row-scoped ruling. Same cycle: dispatch → zero-code stop → escalation → ruling → redispatch.
- **D-21 / root.go:** `GetPricingPolicyForInstallation` was referenced in decisions but didn't exist in code; root.go had no listings wiring though the plan said "inject into existing handler." Both caught only at pre-dispatch audit.

All four were **detectable by reading the decisions against each other and against the code at plan time.** A one-hour contract-satisfiability + prerequisite-existence pass during planning removes two full blocked-dispatch cycles and two hub stalls.

### C. Old-harness residue actively sabotaged workers (`mpc-goal-harness`)
The Slice 3 worker auto-invoked the **superseded** `mpc-goal-harness` skill — because the new harness shipped as a doc with no worker-facing skill, so skill discovery resolved to the old `.agents/skills/mpc-goal-harness`. That triggered a 6-ask migration escalation mid-milestone. Residue still on disk right now:
- `.agents/skills/mpc-goal-harness/` skill + `dispatch_preflight.py` + `openai.yaml` — present.
- `scripts/tests/harness-orchestration.tests.ps1:7,37` — a test that **hard-asserts the dead skill exists** ("RED: mpc-goal-harness skill is missing").
- `contracts/governance/shared-seams.json:10`, `knowledge-routes.json:23` — governance still points at it.
- **Doctrine-path split:** `AGENTS.md:3` binds `docs/superpowers/HARNESS.md` (worktree has this, stale); the M-01 dispatch prompts declare `docs/HARNESS.md` is the only binding doctrine (worktree lacks it). Every worker prompt had to carry an interim "pin HARNESS.md + forbid mpc-goal-harness" workaround clause because the base was ambiguous.

**Until this is deleted atomically and the worker skill ships into the branch base, every future worker rolls the same dice.**

### D. Hermetic test-DB machinery — overkill per-run + fragile ("pra que ser hermetic toda hora")
The operator's instinct is correct. Facts:
- **The code does not require an ephemeral container.** `testsupport/postgres/target.go` requires only `MPC_TEST_DATABASE_URL` pointing at a DB named `^mpc_test_[0-9a-f]{32}$`. Ephemeral-container-per-invocation is *convention*, not a code constraint.
- **4 distinct bring-up paths exist; 3 are redundant hand-rolls** (scratchpad `s6-repro.ps1`, `s7-lane.ps1`, plus the canonical `Postgres.psm1`). Workers hand-wrote runners because the canonical lane didn't fit (see next point).
- **`Postgres.psm1:119` hardcodes the integration package list.** Every new module must be hand-added there. That's *why* the listings lane didn't run out of the box and *why* workers wrote their own runners.
- **Two recurring infra failures, both pure infra, both burned cycles:**
  1. Cold modcache in a fresh worktree → `go run ./cmd/testdb migrate` build fails → harness maps it to `HPG_MIGRATION_FAILED, migrations_first=-1` (a *build* failure disguised as a *migration* failure). Diagnosed, now HARNESS §5 doctrine, but still a manual first-act.
  2. `pg_isready` returns exit 0 during postgres initdb's temp-server phase, *before* the real TCP restart → `CREATE DATABASE` races the restart → `3D000 database does not exist` → every test "fails." Fix was a 20-retry `CREATE DATABASE` loop, hand-written in `s7-lane.ps1`, **not yet in the harness.**

Per-run ephemeral bring-up + a boot race + a hardcoded package list = every single test run risked a 20-retry stall on infra that has nothing to do with the code under test.

### E. Review churn — high false-positive tax
Real 🔴s were caught and mattered (Slice 2 scan-cursor exhaustion, Slice 7 missing ANALYZE, Slice 3 unbounded page loop) — the reviewer earned its place. But the **false-positive rate was high**: Slice 1 (1/1 FP — `provider` col was IC-02-mandated), Slice 4 (**3/3 all FP**, all dismissed with evidence — the reviewer flagged established repo idioms: pgx `**string` nullable scan, `len(page.Items)` page_size semantics). Each FP cost the milestone owner an adjudication + evidence write. The reviewer lacked the repo's own idioms as context, so it re-litigated settled patterns.

### F. Dispatch-channel doctrine churned 3× — never a stable answer
Ledger line 5: rescue-only. Line 196: switch Slice 6+ to `/codex:rescue` because "raw exec isn't panel-visible." Line 208: **reverse** — empirically, OS-process raw-exec-with-tee IS what drives the operator panel, `/codex:rescue` is NOT. The team never had one correct, stable answer on which path surfaces in the panel, and re-decided it per slice.

### G. Pre-existing sibling breakage = permanent gate noise
`TestPhase1SmokeFlow` (`PRICING_INVALID_PRODUCT_ID`) and an orders duplicate-identity flake fail **every full lane run**, unrelated to listings. Proving non-linkage cost 3 separate repro captures and an acceptance-criterion adjustment (RULING B). As long as the gate runs the full lane, it can **never be green** for an M-01-scoped change — the gate is structurally unable to give a clean signal.

---

## 3. Recommendations to the HUB (ranked by wall-clock returned)

1. **Planning must emit a write-DAG + a contract-satisfiability pass, before any dispatch.** The single batched planner pass (already doctrine) must additionally output: (a) file-level shared-seam map → which slices serialize on which files, which are disjoint; (b) pre-identified additive-lock sub-workers (formalize the D-21/D-24 pattern as a planning output, not a mid-flight discovery); (c) a decisions-vs-decisions-vs-code reconciliation that resolves conflicts like D-18/D-20/IC-02 **at plan time.** → Removes findings A + B: 2 wasted dispatches, 2 escalations, 2 late prereqs.

2. **Delete `mpc-goal-harness` residue atomically and ship the worker skill into the branch base.** One corrective chip: remove the skill dir, the hard-assert test, the `shared-seams.json`/`knowledge-routes.json` entries; fix `AGENTS.md` to bind `docs/HARNESS.md`; delete the stale `docs/superpowers/HARNESS.md`; land `harness-worker`/`codex-dispatch` skills so discovery resolves correctly with **zero interim pin clause.** → Removes finding C for every future worker.

3. **Make the integration lane self-discovering and race-proof, once.** (a) `Postgres.psm1` globs `./internal/modules/*/adapters/postgres` + `./tests/integration` instead of a hardcoded list. (b) Replace the `pg_isready` gate with the retry-`CREATE DATABASE`-until-exit-0 loop (promote the `s7-lane.ps1` fix into the harness). (c) Fold modcache-warm into worktree bootstrap as an automatic step, not a manual first-act. → Kills all of finding D's recurring failures and the reason 3 scratchpad runners existed.

4. **Reconsider hermetic-per-run: one long-lived container per milestone session.** Create one ephemeral `postgres:16` at worktree bootstrap, drop it at close; each test run just `CREATE DATABASE mpc_test_<hex>` against it. The code already supports this (only needs the URL). Removes the boot-race from *every* run and answers the operator's "why hermetic every time" directly — keep isolation (fresh DB per run) without paying container-boot per run.

5. **Partition the gate: `--owned <module>` lane flag.** Gate on M-01-owned packages + a pinned known-failure allowlist for siblings. Full-lane-green is already dropped as acceptance (RULING B); make it structural so pre-existing breakage (finding G) never reaches a milestone verdict and never needs re-proving.

6. **Settle the dispatch channel once and document the panel-visibility truth.** Empirically established (ledger 208): OS-process raw-exec-with-tee is the panel-visible path. Either sanction that as the milestone-session dispatch and retire the rescue-only clause, or add the companion-tee to `/codex:rescue` (the deferred "Option A") so both converge. Stop re-deciding per slice.

7. **Give the reviewer the repo's idioms as context to cut the FP tax.** Feed the established patterns (pgx `**string` nullable scan, `len(page.Items)` page_size, the shared scanner idiom) so it stops flagging settled code. Keep it scoped to changed behavior. Do NOT over-correct — the real 🔴s (scan cursor, ANALYZE, unbounded loop) must still surface.

8. **Reduce hub-stall latency.** The 8 round-trips dominated wall-clock mainly because the hub session sits idle and needs the operator to open it to rule. Recs 1 removes the two biggest (Slice 2/4). For the rest: batch related asks into one escalation (the 6-ask migration package was the right shape — do that by default), and/or an async hub-inbox so a chip doesn't hard-stall on an idle hub.

---

## 4. What went RIGHT (keep these)

- **Honest zero-code stops on contract conflicts.** Workers refused to guess through an over-constrained contract (Slice 2, Slice 4). That discipline is exactly right — the fix is to catch the conflict at plan time, not to loosen the discipline.
- **Additive-lock sub-worker pattern (D-21, D-24).** The one real parallelism win. Formalize it as a planning output.
- **ADR-17 held everywhere.** Unknown cost → null, never 0/default, across list/detail/summary/by-product. No fabricated data escaped.
- **The multi-agent live dashboard.** Genuinely useful observability; already adopted as HARNESS §8.
- **Evidence discipline.** Every slice's lane result, EXPLAIN, and repro is written down — this retrospective was possible *because* the ledger is complete.

---

## 5. One-line ask

Land recs **1, 2, 3** before M-02 planning starts. They remove the three friction classes (non-parallel plan, harness residue, fragile per-run DB lane) that caused the bulk of M-01's wasted wall-clock. Recs 4–8 are the next tier.
