# Lane: cicd

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.
>
> Calibration: solid professional level, not Google-tier. Clear dependencies, clear modules, clear
> consumable surfaces, no hand-maintained redundancy, rules that have existed for decades, optimised
> for maintainability and future scaling. The operative sentence: **"this way it gets so much harder
> to send bad PRs."** Success = a mechanism that makes a bad change hard to land, not a cleaner
> codebase.

## Bottom line

**Nothing blocks a merge to `main` today. Nothing ever has.** There is no branch protection
reachable, no `pull_request`-triggered workflow in the current tree, no git hook, no npm
`pre-*`/husky/lint-staged wiring, no `.claude/settings.json` hook. The single GitHub Actions
workflow that exists (`release-images.yml`) triggers on `push` to `main` — it publishes images,
it does not gate anything — and its one recorded execution on the reachable `legacy` remote
**failed on billing**, not on code (see Gate inventory). Every verification mechanism in this repo
(`arch-gate.sh`, `harness.ps1`, `Policy.psm1` governance) is opt-in: a human has to type the command.
Measured this session: run cold with a trivial no-op `-BaseSha` (current `HEAD`), the governance
gate **fails with 58 violations** against code already on `main`, and `arch-gate.sh` fails from (at
least) five independent causes. The instruments exist; none of them are load-bearing.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| F1 | gap | No `pull_request`-triggered workflow exists anywhere in the current tree; the only workflow triggers on `push`/`tag`/`workflow_dispatch` and publishes Docker images | `.github/workflows/release-images.yml:8-11`; `grep -rl "pull_request" .github/` → 0 hits | 1 workflow file, 0 PR gates |
| F2 | gap | No branch protection reachable/observed on either remote; `origin` (`developmentconexus-ops/...`) doesn't resolve from this machine, `legacy` resolves and returns no protection object because the workflow itself never ran clean | `gh api repos/developmentconexus-ops/marketplace-central/branches/main/protection` → 404; `gh repo view` → `Could not resolve to a Repository` (PHASE-0:56-58) | 2 remotes, 0 confirmed protections |
| F3 | gap | The one workflow that has ever executed on a reachable remote (`legacy`) failed **on billing**, not code: "recent account payments have failed or your spending limit needs to be increased" | `gh run view 29712873135 --repo leandrotcawork/marketplace-central` → both `build-push` jobs show that annotation; run triggered 2026-07-20 | 1 of 1 recorded runs on `main` push, 0 succeeded |
| F4 | gap | No pre-commit/git-hook/husky/lint-staged machinery anywhere. `.git/hooks/` holds only `.sample` files; no `.husky/`; no `husky`/`lint-staged`/`pre-commit` string in either `package.json` | `ls .git/hooks/` (all `*.sample`); `grep -riE "husky\|lint-staged\|pre-commit" package.json apps/web/package.json` → 0 hits | 0 hooks |
| F5 | gap | No `.claude/settings.json` hook is configured to run any gate | `grep -n "hooks" .claude/settings.json .claude/settings.local.json` → 0 hits | 0 hooks |
| F6 | drift | Measured this session: `npm run harness:governance -- -BaseSha <HEAD>` (i.e. governance run against the tip of `main`, zero diff) **fails with 58 violations** across 6 error codes, against real production code already merged | full output captured; `grep -c "^error_code="` → 58; breakdown below | 58 violations, 12 baseline-exception categories already in use |
| F7 | gap | `scripts/arch-gate.sh` step 3 (architecture detectors, `archscan`) also fails at `HEAD`, independent of the gofmt/ADR-023 causes already in the established-facts list: 44 findings across 3 of 4 scanned roots | `go run ./internal/arch/cmd/archscan -root internal/contexts` → 1 finding; `-root internal/adapters` → 1; `-root internal/composition` → 42; `-root internal/kernel` → 0 | 44 findings, 3/4 roots dirty |
| F8 | gap | `internal/contexts/` (the target tree of the module migration) has **zero governance drift coverage**. `Policy.psm1` walks only `apps/server_core/internal/modules` for module registration/dependency/layer checks; there is no code path, and no error code, that inspects `internal/contexts` at all | `grep -n "contexts" scripts/harness/Policy.psm1` → 0 hits; module-walk anchored at `Policy.psm1:302` (`$moduleRoot = ... 'apps/server_core/internal/modules'`) | 10 non-test `.go` files under `internal/contexts` today, 0 of them governed |
| F9 | gap | A written negative-fixture test (`scripts/tests/governance-drift.tests.ps1`) already asserts a `GOV_CONTEXT_UNREGISTERED` check should exist for F8's gap — but no such error code exists in `Policy.psm1`, and the test file itself does not pass its own earlier positive-fixture assertion, so the assertion for F8 is never reached | see "Unverified" — file was under concurrent uncommitted edit during this session; independently confirmed `Policy.psm1` (unmodified during the session) has zero `GOV_CONTEXT_UNREGISTERED` handling | 1 test file, non-executing |
| F10 | gap | None of the 11 `scripts/tests/*.tests.ps1` negative-fixture scripts are invoked from any automated place. They are plain PowerShell scripts (not Pester `Describe`/`It`, except `live-oracle-docker-runner.tests.ps1`) that must be run with `pwsh -File <path>` by a human; `Invoke-Pester` appears in zero non-vendor files | `grep -rln "Invoke-Pester" . --include="*.ps1" --include="*.psm1"` (excluding `.runs`) → 0 hits | 11 scripts, 0 wired |
| F11 | gap | Of 11 negative-fixture scripts run cold this session, 6 did not reach a clean pass; 2 of those 6 fail for reasons unrelated to Docker/Oracle availability (pure PowerShell + repo-file logic) | measured, see Gate inventory / commands run | 6/11 non-green |
| F12 | drift | `contracts/governance/knowledge-routes.json` cites a test file, `scripts/tests/harness-orchestration.tests.ps1`, that does not exist in the current tree | `grep -n "harness-orchestration" contracts/governance/knowledge-routes.json:12`; `ls scripts/tests/harness-orchestration.tests.ps1` → No such file | 1 dangling reference |
| F13 | hazard | The governance `-BaseSha` drift check only fires the API/SDK-atomicity rule (`GOV_API_SDK_SPLIT`) and other diff-scoped checks when a caller supplies a valid 40-hex commit SHA by hand; there is no default and no automatic derivation of a base, so a governance run with no `-BaseSha` silently skips all diff-scoped rules | `scripts/harness/Policy.psm1:447` (`if (-not [string]::IsNullOrWhiteSpace($BaseSha))`) gates the entire diff-scoped block; `scripts/harness.ps1:113` only validates the format of `-BaseSha` if one is passed, `Invoke-Governance -Mode validate` (the default alias target `governance-validate`) never asks for one at all | whole diff-scoped rule family (1 rule: `GOV_API_SDK_SPLIT`) is opt-in per invocation |
| F14 | idiom | `scripts/release.sh` exists specifically to bypass CI/GHCR-Actions entirely ("Local release path (no CI minutes needed)") — deploy does not depend on any gate, automated or manual, passing | `scripts/release.sh:1-9` comment block | 1 escape hatch, unconditional |
| F15 | drift | A `pull_request`-triggered workflow (`wiki-lint`, warn-only policy) existed once on a feature branch on the `legacy` remote and failed all 3 times it ran; it was never merged to `main` and is absent from the current tree | `gh api repos/leandrotcawork/marketplace-central/actions/workflows` lists it `active` on GitHub's side; `git branch --all --contains 39934922` → only `remotes/legacy/feat/llm-wiki-m1`; `git merge-base --is-ancestor 39934922 HEAD` → not an ancestor | 3 runs, 3 failures, 0 on `main` |

## The five heaviest, with detail

**1. F1+F2+F3 — there is no merge gate, configured or executed, and the one CI mechanism that exists cannot even run.**
`.github/workflows/release-images.yml` is the only workflow in the tree. It triggers on `push` to
`main`, `v*` tags, or manual dispatch — never on `pull_request` — and its job is to build and push
Docker images to GHCR, not to validate anything. Branch protection cannot be confirmed on `origin`
(`developmentconexus-ops/marketplace-central`, unresolvable with this machine's `gh` credentials —
matches PHASE-0's noted fact) and returns nothing on `legacy` because there is nothing to return. The
only recorded execution of this workflow, on `legacy`, from 2026-07-20, failed both matrix jobs
before a single build step ran, with GitHub's own annotation: *"The job was not started because
recent account payments have failed or your spending limit needs to be increased."* This is not a
code failure — it is proof that even the weakest possible gate (a post-merge image build) has never
successfully executed once, for a reason the "zero-spend tooling" constraint (PHASE-0:35-38)
predicts exactly. Whatever remediation is designed, it cannot assume any paid or previously-working
GitHub Actions capacity exists.

**2. F6 — the governance gate, run against the tip of `main` with zero diff, fails with 58 violations.**
`npm run harness:governance -- -BaseSha <HEAD-sha>` runs `Test-GovernanceContracts` +
`Test-GovernanceDrift` against the registries and code exactly as committed. Breakdown of the 58:

```
20 RCFG_UNAPPROVED_READER   (env var read from a file/path the registry doesn't list as a reader)
12 GOV_MODULE_DEPENDENCY    (module imports another module without declaring the dependency)
10 RCFG_READER_MISSING      (registry claims a reader that isn't there / doesn't match)
 8 GOV_MODULE_LAYER         (adapters/transport/registry cross-module import, no exception)
 7 RCFG_UNDECLARED_READ     (env var read that the registry has never heard of)
 1 GOV_MIGRATION_PREFIX     (duplicate numeric prefix: 0093 used by two migration files)
```
19 pre-existing `temporary_exception` baselines already absorb other violations (panic sites,
direct env readers, one adapter-layer exception) — those are not counted above; the 58 are *un*-
absorbed. This means: if this check were wired to CI today in blocking mode, it would fail on `main`
itself, on unrelated code, before a single new PR landed. The check is real and its detectors work
(confirmed against negative fixtures below) — but it has never been run as a gate, so violations
have been accumulating un-flagged since whenever each was introduced. Sizing note: `GOV_MIGRATION_PREFIX`
(`0093_orders_status_details_nullable.sql` / `0093_sync_state_market_queue_entity_split.sql`) is a
correctness hazard independent of governance — two migrations sharing a numeric prefix is an
ordering hazard for any tool that sorts by filename.

**3. F7 — `archscan`, the fourth of `arch-gate.sh`'s five steps, also fails at `HEAD`, on top of the causes already in PHASE-0's established facts.**
PHASE-0 fact 2 names three causes for `arch-gate.sh` FAIL (CRLF gofmt noise, 22 genuinely
unformatted files, 234 ADR-023 violations) but does not mention the architecture-detector step
(`internal/arch/cmd/archscan`, run once per root: `kernel`, `contexts`, `adapters`, `composition`).
Run cold this session: `kernel` is clean, `contexts` and `adapters` each report 1 finding
(`facts/value-discarded` in test files), and `composition` reports 42 (`adapters/vendor-token-
outside-adapters` — `mercado_livre`/`mercadolivre` literals in `internal/composition/root.go`,
matching the D-35 debt's already-recorded count). This is additive detail to the established fact,
not a contradiction: `arch-gate.sh` fails from at least five distinguishable causes, not three, and
step 3 is one more independently-red step sitting between the gofmt failure (step 1) and the unit
tests (step 4) that D-51 already proved run anyway.

**4. F8+F9 — the module-boundary governance check has a hole shaped exactly like the current migration: `internal/contexts` is invisible to it.**
`Policy.psm1`'s drift check (`Test-GovernanceDrift`) walks `apps/server_core/internal/modules` by
directory listing to catch undeclared/orphaned modules (`GOV_MODULE_COVERAGE`,
`Policy.psm1:302-311`) and checks import strings against the registry for dependency/layer
violations — but the walk root is hardcoded to `internal/modules`. There is no equivalent walk, and
no error code, for `internal/contexts` — the tree ADR-023 says is the migration target (PHASE-0 fact
10). A file dropped under `internal/contexts/<anything>/` with no registry entry produces zero
findings from this checker, silently. This is exactly the class the repo's own architecture doctrine
worries about (ADR-023's module protocol). A negative-fixture test already exists for this — 
`scripts/tests/governance-drift.tests.ps1` (diff observed this session, see Unverified section)
asserts a `GOV_CONTEXT_UNREGISTERED` code that must fire on an unregistered `internal/contexts/orphan`
directory — but `Policy.psm1` (confirmed unmodified during this session via `git status`) contains
zero code for that error code, and the test file does not currently reach that assertion because it
fails earlier, on its own positive-fixture setup (`GOV_COMPOSITION_MISSING` x17). Whether the
positive-fixture failure is pre-existing or an artifact of an in-flight, uncommitted edit to that
same file made by another process during this session is `unverified` (see below) — what's not in
doubt is that `Policy.psm1` itself, unmodified, has no `internal/contexts` coverage today.

**5. F10+F11 — the negative-fixture layer is real and detailed where it exists, but disconnected from execution and partly broken.**
`scripts/tests/governance-drift.tests.ps1` alone declares 15+ distinct `Assert-FailureCode` negative
fixtures (one per governance error code: `GOV_MODULE_COVERAGE`, `GOV_MODULE_DEPENDENCY`,
`GOV_APPLICATION_IMPORT`, `RCFG_UNDECLARED_READ`, `RCFG_DYNAMIC_READER_UNBOUNDED` x3,
`RCFG_ALIAS_COLLISION`, `RCFG_SECRET_CLASS_MISMATCH`, `GOV_PRODUCTION_PANIC` x2,
`GOV_MIGRATION_PREFIX`, `GOV_FRONTEND_FETCH`, `GOV_API_SDK_SPLIT` x2). `internal/arch/scan_test.go`
has 10 fires-on-fixture/silent-on-clean-fixture pairs for the `archscan` detectors — that one is
wired to `go test ./internal/...` and does run (as part of `arch-gate.sh` step 4 / any manual
`go test ./internal/...`). But none of the `scripts/tests/*.tests.ps1` files are invoked by anything
automated (F10), and measured cold this session, over half don't finish clean: `governance-drift`
fails on its own positive fixture before reaching any negative assertion; `harness-aliases` fails
because the live `harness:governance` alias itself is red at `HEAD` (consistent with F6);
`governance-contracts` errors on a missing `rg` binary with no declared prerequisite;
`postgres-lifecycle` and `live-oracle-docker-runner` fail in ways plausibly tied to local
Docker/Oracle availability (`unverified`, not chased further — out of lane). Net: a guard nobody has
automated is also, in measured practice today, a guard that frequently doesn't run to green even by
hand.

## Gate inventory

| gate | invoked from | blocks what | fires? | negative fixture? |
|---|---|---|---|---|
| `release-images.yml` (GHA) | `push` to `main`/`v*` tag, `workflow_dispatch` | nothing — publishes images post-hoc | Configured, but its one recorded remote execution failed on billing before any step ran (F3) | No |
| Branch protection (`origin`/`legacy`) | GitHub, pre-merge | would block direct pushes / require checks, if configured | Unverified on `origin` (unresolvable); confirmed absent in practice on `legacy` (workflow allowed to run and fail) | N/A |
| `scripts/arch-gate.sh` | human typing `bash scripts/arch-gate.sh`; **zero** npm script references it (`grep -n "arch-gate" package.json` → 0, established fact 3) | nothing automatically; would gate gofmt/vet/archscan/unit-tests/clean-tree if wired | Runs when invoked; FAILS at `HEAD` from ≥5 causes (CRLF noise, 22 real gofmt files, 234 ADR-023 violations, 44 archscan findings, and whatever else touches the working tree) | Partial — step 3 (archscan) has real fixtures (`internal/arch/scan_test.go`); steps 1/2/5 (gofmt, vet, clean-tree) do not |
| `npm run harness:unit` | human / could be CI | Go `tests/unit` + FE vitest only — **not** `internal/...` (D-51) | Runs, passes (605 FE tests + Go `tests/unit`) | No — by construction it can't prove `internal/...` regressions |
| `npm run harness:integration` | human / could be CI | `tests/integration` against an ephemeral `mpc_test_<hex>` Postgres | Runs when Docker is available; not exercised this session (out of lane, needs Docker) | Some (`postgres-*.tests.ps1`), mixed pass/fail this session |
| `npm run harness:governance` (validate/drift/all) | human only, must pass `-BaseSha` for drift/all | Module dependency/layer, runtime-config reader, migration-prefix, frontend-fetch, API/SDK-split rules — **only if someone runs it and treats a failure as blocking** | Runs; **FAILS with 58 violations at `HEAD`, zero diff** (F6) | Extensive, but the runner script itself is red before reaching most assertions (F9/F11) |
| Git hooks (`.git/hooks/*`) | would be pre-commit/pre-push | nothing | Absent — only `.sample` files | N/A |
| Husky / lint-staged | would be `git commit` | nothing | Absent — no config in either `package.json` | N/A |
| `.claude/settings.json` hooks | would be Claude Code tool events | nothing | Absent | N/A |
| `scripts/tests/*.tests.ps1` (11 files) | human running `pwsh -File <path>` | nothing automatically; would prove the guards above work | 5/11 clean this session, 6/11 not (F11) | These *are* the negative fixtures for the other rows |
| `scripts/release.sh` | human, deliberately bypasses GHA | N/A — it's the escape hatch, not a gate | Runs (not exercised this session — would push images) | N/A |

## What is actually fine

- **The architecture-detector negative fixtures are genuinely good engineering.** `internal/arch/scan_test.go`'s 10 fires-on/silent-on pairs are exactly the right shape (a minimal fixture that must trip the detector, a clean one that must not), they run under plain `go test`, and they are reachable via `go test ./internal/...` without any special setup. If a CI gate gets built, this file is the reference pattern to keep, not touch.
- **The integration lane's ephemeral-database design is sound and already proven** (established fact 4, re-confirmed by nothing contradicting it this session): a fresh `mpc_test_<hex>` database per run means no cross-run state contamination — a common source of CI flakiness elsewhere is already engineered out here.
- **`go vet ./...` is clean at `HEAD`** — zero findings, confirmed this session. Whatever gate design gets chosen, `go vet` is free, fast, and already passing; wiring it costs nothing.
- **`Policy.psm1`'s rule *logic*, where it has coverage, is detailed and defensible** — the runtime-config reader/alias/secret-classification rules, the module dependency/layer rules, and the migration-prefix check are not naive greps; they cross-reference multiple registries and (mostly) have matching negative fixtures written for them. The gap is wiring and reach (`internal/contexts`), not the quality of what's there for `internal/modules`.
- **The ephemeral/no-cross-branch discipline already learned the hard way is holding**: `Policy.psm1:198-206`'s file-scanning exclusion list (`.worktrees`, `.claude`, `node_modules`, etc.) carries a comment explaining a prior incident (cross-branch checkout contamination) and the fix is still in place, unregressed.
- **`scripts/release.sh` and the Docker dev stack are explicitly out of this lane's remit per PHASE-0's do-not-touch candidates** — nothing observed this session contradicts treating them as fine; not exercised further here.

## Unverified / needs judgment

- **Concurrent contamination of governance files during this session.** Mid-session, `git status` showed uncommitted modifications to `contracts/governance/schemas/modules.schema.json` and `scripts/tests/governance-drift.tests.ps1` (later joined by `scripts/harness/Policy.psm1` itself, observed still in-progress at report-write time), plus three untracked scratch files (`scratch_analyze_tables.py`, `scratch_migrations_list.txt`, `scratch_sql_scan.py`) — none created by this lane. The diff on the first two files is an in-flight, uncommitted fix for exactly the `GOV_CONTEXT_UNREGISTERED`/context-registration gap this report describes as F8/F9 (adding a `kind: module|context` schema property and a matching negative-fixture block), consistent with a different concurrent lane or the hub session actively working the same seam live. This lane did not create, revert, or otherwise touch any of these files, per the read-only constraint. Practical effect on this report's evidence: the F6 governance run (58 violations) and the F8 "`Policy.psm1` has zero `internal/contexts` coverage" observation were both captured **before** `Policy.psm1` entered the modified set — they reflect the committed engine as it stood at commit `11b9e494`, not the in-flight edit. By the time this report was written, `Policy.psm1` was itself mid-edit toward closing F8/F9; that work's completeness is `unverified` and un-reviewed by this lane, since it did not exist as a fact at measurement time and reviewing someone else's uncommitted work is outside this lane's discovery mandate. Whether `governance-drift.tests.ps1`'s `GOV_COMPOSITION_MISSING` failure (used as supporting color in finding 4) reflects the committed test file or an intermediate state of someone else's in-progress edit is likewise `unverified` — flagged rather than asserted.
- **Repo visibility of `origin`** — per PHASE-0, `unverified`; this lane did not obtain new evidence either way. Branch-protection reasoning above is scoped to what's reachable (`legacy`) plus the absence of any protection-dependent workflow in the current tree, which holds regardless of `origin`'s visibility.
- **`postgres-lifecycle.tests.ps1` and `live-oracle-docker-runner.tests.ps1` failures** — plausibly Docker/Oracle-availability artifacts of this environment rather than genuine regressions; not chased further, as Docker/Oracle-lane health is explicitly another lane's / out-of-scope territory per PHASE-0's do-not-touch candidates.
- **`governance-contracts.tests.ps1`'s dependency on an unrigged `rg` binary** — this session's `pwsh` had no `rg` on `PATH` (Claude Code's own Grep tool is not the system `rg`), so the script errored before validating anything. Whether the operator's normal interactive shell has `rg` on `PATH` is `unverified`; either way, the dependency is undeclared anywhere in the repo (no README, no check-for-tool guard), which is itself worth noting even if the operator's shell happens to have it.
- **Whether `wiki-lint.yml` (F15) reflects an abandoned direction or a template for what "PR gate" should become here** — out of scope for this lane's discovery mandate (recommending direction is the synthesis step's job); reported as a fact only.

## Commands run

```
git log -1 --format="%h %ad %s" -- scripts/tests/governance-drift.tests.ps1
git log -1 --format="%h %ad %s" -- scripts/harness/Policy.psm1
find .github -type f
ls .git/hooks/ | grep -v sample
find .claude -maxdepth 3 -type f
grep -n "hooks" .claude/settings.json .claude/settings.local.json
cat .gitattributes
gh api repos/developmentconexus-ops/marketplace-central/branches/main/protection
gh repo view --json visibility,nameWithOwner,defaultBranchRef
gh run list --repo leandrotcawork/marketplace-central --limit 50
gh api repos/leandrotcawork/marketplace-central/actions/workflows --jq '.workflows[] | "\(.name) | \(.path) | \(.state)"'
gh run view 29712873135 --repo leandrotcawork/marketplace-central
git branch --all --contains 39934922
git merge-base --is-ancestor 39934922 HEAD
git rev-list --count origin/main..HEAD   # 97 (established fact says 96 as of 1473e863; this session is 1 commit later, at 11b9e494)
git rev-list --count HEAD..origin/main   # 0
grep -riE "husky|lint-staged|pre-commit" package.json apps/web/package.json
find . -name package.json -not -path "*/node_modules/*"
bash scripts/arch-gate.sh   (referenced from prior session evidence embedded in .mnfs/HARNESS-DEBTS.md; not re-run in full this session to avoid re-deriving established fact 2 per budget instructions)
cd apps/server_core && go vet ./...
cd apps/server_core && go run ./internal/arch/cmd/archscan -root internal/kernel|internal/contexts|internal/adapters|internal/composition
SHA=$(git rev-parse HEAD); npm run harness:governance -- -BaseSha "$SHA"
grep -c "^error_code=" / grep "^error_code=" ... | sort | uniq -c   (on captured governance output)
pwsh -NoProfile -File scripts/tests/{governance-contracts,governance-drift,harness-aliases,harness-environment,harness-execution,hermetic-lanes,live-oracle-docker-runner,postgres-contract,postgres-go,postgres-lifecycle,dev-local-runtime}.tests.ps1
grep -rln "Invoke-Pester" . --include="*.ps1" --include="*.psm1"
grep -n "GOV_CONTEXT_UNREGISTERED|contexts" scripts/harness/Policy.psm1
grep -n "^func Test" apps/server_core/internal/arch/scan_test.go
grep -rln "TestModuleBoundaryADR023" apps/server_core --include="*.go"
find apps/server_core/internal/contexts -name "*.go" -not -name "*_test.go" | wc -l   # 10
find apps/server_core/internal/modules -name "*.go" -not -name "*_test.go" | wc -l    # 462
git diff --stat contracts/governance/schemas/modules.schema.json scripts/tests/governance-drift.tests.ps1
git status --porcelain --untracked-files=all
```
