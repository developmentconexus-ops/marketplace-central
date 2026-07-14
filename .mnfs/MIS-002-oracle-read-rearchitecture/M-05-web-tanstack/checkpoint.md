# M-05-web-tanstack — Milestone Checkpoint (TERMINAL)

**Verdict: QA PASS** · **Frozen SHA `c2aea877`** · branch `claude/epic-pare-2eaed9` · worktree `.claude/worktrees/epic-pare-2eaed9`
**Diff range:** `d413fd78..c2aea877` · **NOT merged to main** (hub owns main merges)
**Date:** 2026-07-14

## Commits
| SHA | Feature |
|---|---|
| `2c8739a6` | F-01-tanstack-adoption — `feat(web): adopt TanStack Query reads` |
| `f0ab3c5e` | F-02-mutation-invalidation — `feat(web): invalidate query namespaces on Oracle-backed writes` |
| `c2aea877` | merge of hub base `3f48ca91` (dual-review protocol + mpc-verifier on sol) |

## Result
Both dual reviewers passed independently at `c2aea877`: Codex `gpt-5.6-sol` → NO BLOCKING FINDINGS; Claude QA validator → PASS on M-05-C01..C05, no correction scope. Full detail in `validation-result.md` (same directory).

Evidence at frozen SHA, real exit codes: feature suites 32/32 (exit 0), canonical `npm test` 11/11 (exit 0), `apps/web` build 1832 modules (exit 0). Seams clean — no Go, OpenAPI, or `sdk-runtime` edits.

## Decisions and deviations (orchestrator-owned — hub should be aware)

1. **Option B (user-approved): new CatalogPage, legacy ProductsPage DELETED.** Post-M-02 the endpoint returns `CatalogProductFact` while ProductsPage still rendered old `CatalogProduct` fields — and `vite build` is esbuild-only, so it typechecks nothing and the staleness was invisible. Deleting beat patching. **Consequence:** IC-01's crosswalk row `product edit → ['catalog']` has no implementation target and is MOOT. Both reviewers independently judged that acceptable (the row is a conditional; nothing can violate it). **Forward obligation for IC-01's owner: annotate the row as dormant — do not delete it. If a product-edit surface returns, the `['catalog']` obligation binds it.**

2. **Sequential F-01 → F-02, not parallel** (deviation from the milestone plan). `StockSeguroPage` and `OrdersPage` carry BOTH F-01 reads and F-02 writes, so the disjoint-seam precondition for parallel writers failed. Same-file overlap would have meant two writers on one seam.

3. **`--ff-only` to hub base `3f48ca91` was not possible.** The branch had **diverged** (two commits of my own), not merely fallen behind, so the hub's conditional ff-only instruction did not apply. Incoming commits were docs + agent config with **zero seam overlap**, so a `--no-ff` merge was safe and conflict-free → `c2aea877`. Nothing reset, reverted, stashed, or cleaned at any point.

4. **Orchestrator applied F-02's correction directly** (deviation from strict role separation). Codex was killed by host teardown twice mid-run — the second time before it could apply a correction it had been handed. Its code survived uncommitted and was contract-correct; the outstanding work was 6 mechanically-diagnosed lines plus artifacts. Rather than burn a third long run, the orchestrator applied the fix and ran all verification. **Cross-model integrity is preserved regardless** — it rests on the milestone-end Sol review, which ran and passed.

## Defects the orchestrator caught that the implementer did not

Codex's F-02 was contract-correct on crosswalk, namespaces, cache options, and error handling, but shipped two defects; it was killed before self-verifying:

1. **Production defect — `mutationFn` passed by reference at 5 sites.** TanStack v5 calls `mutationFn(variables, context)`, so the second `{client, meta, mutationKey}` argument was forwarded into SDK methods and unbound them from the client. Caught by spy output. Fixed with typed arrow wrappers (`Parameters<Client["method"]>[0]`, no `any`).
2. **Test defect — error assertions raced the production `retry: 1`.** Fixed with a test-local `createTestQueryClient()` deriving from the production factory and disabling only retry. **Production defaults unchanged** (contract-fixed).

## Process lessons (worth propagating)

- **Codex self-reports are not evidence.** Earlier in this milestone Codex printed a plausible correct answer while BOTH its `exec` calls failed — it answered from context injection, not execution. Every Codex claim here was independently re-verified against the repo. The F-02 defects above are exactly what that policy catches.
- **Never judge a test runner through a `| tail` pipe** — the pipe's exit code masks the runner's. This produced one false "exit 0" while 8 tests were failing.
- **Never judge tests from a contended machine.** A run concurrent with a build reported 8 failures (environment 538s, all "timed out in 5000ms"); quiet, it finished in 7.4s. Both reviewers confirmed the discard was legitimate. Real failures survive a quiet re-run — contention artifacts do not.

## Infrastructure finding — machine-wide, affects other sessions

**The Codex CLI 0.144.4 did NOT ship broken.** Proven by running the exact binary from its release path with sandbox ON *before* changing anything (`SANDBOX_OK`, 3374ms). The fault was the **local launcher shim**: `Programs\OpenAI\Codex\` contained only the `bin` junction, missing the `codex-resources` / `codex-path` siblings that `codex-package.json` requires → `program not found` → `CreateProcessWithLogonW failed: 2`. Fixed with two junctions (verified `SHIM_FIXED`, 977ms, sandbox ON; runner `codex-command-runner-0.144.4.exe` now cached).

`.sandbox-bin` held working runners for `0.144.0-alpha.4` (Jul 9) and `0.144.2` (Jul 13), so copies worked until Jul 13 — **something (installer/update) removed the junctions while leaving `bin`. A future Codex update may strip them again; recreate the two junctions if `program not found` returns.**

**Cross-session:** the sibling **M-04 worktree (`focused-borg-9c9811`) appears in the sandbox log failing identically at 12:34** — it was silently getting nothing done. The junction fix is machine-wide and unblocks M-04 too.

**Security posture:** the sandbox bypass the user consented to was used ONLY for the in-flight F-01. Once the shim was fixed, F-02 and all review runs used the real sandbox (`-s workspace-write`, and `-s read-only` for the Sol review). Bypass retired.

## Carried forward — non-blocking, do NOT reopen M-05
1. `packages/web-query/src/index.ts:72` — `noCacheDepth` is transport-wide; a concurrent unrelated GET can pick up `no-cache` during a refresh window. *Found only by the cross-model reviewer.*
2. `CatalogPage.test.tsx:92` — under-asserts that `as_of` actually changes (same regex before and after refresh). **Flagged independently by BOTH reviewers.**
3. C02's behavioral remount proof covers catalog only; stock/pricecost verified by constant + call-site inspection.
4. `CatalogPage.test.tsx:73` double-unmounts `first`.
5. Product enrichment editing is now unreachable in the UI (SDK methods survive → dormant, not destroyed).

## Routed to Mission Strategist (not blocking M-05)
- **Annotate IC-01's `product edit → ['catalog']` row as dormant with forward obligation.**
- **Repo typecheck is broken** (`TS2688: Cannot find type definition file for 'node'`), accepted as pre-existing baseline. **A green build is not a type-safe build** — this is the exact blindness that hid the stale-field ProductsPage. Mission hygiene.
- **`AGENTS.md` drifts from config**: it pins `mpc-verifier` to `gpt-5.6-luna`/high; `.codex/agents/mpc-verifier.toml` is now `gpt-5.6-sol`/medium (user's own change at `6e992d3c`).

## Next
Hub to review and own the merge to main. Nothing further is owned by this milestone session.
