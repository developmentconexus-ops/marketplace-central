# M-03 verification ladder — L0/L1 evidence (pre-F-04 state)

Chip: CHIP-M03 · branch `chip/m-03-mutation-envelope-writes` · run window 2026-07-17.
Anchor/BaseSha: `a49168e641ffd6f61932ca57c29b1d1bdcde2fb0` (merge-base vs main).
Evidence head: `8a53105c`. F-04 remains GATED (hub rebase trigger); ladder re-runs after
F-04 lands, this file then extends.

## L0

- `go build ./...` + `go vet ./...` — GREEN @ 1ebb7cc3 (absolute GOCACHE, apps/server_core);
  re-proven implicitly by the full-sweep run below at 5e5c14ee.
- Web typecheck: no standalone script exists (apps/web has no `typecheck`); web compilation
  proven via harness:unit web member (vite/vitest transform) — see L1.
- Governance lane: **status=passed** @ 8a53105c, clean detached worktree
  (`Documents/mpc-govcheck-m03`), full 40-hex BaseSha, zero violations,
  14 baseline_exceptions (all pre-existing registry entries).

### Governance blockage → resolution trail (all adjudicated by hub)

| Finding | Cause | Resolution |
| --- | --- | --- |
| StrictMode crash "property 'dependencies' cannot be found" (Policy.psm1:316) | chip-created `mutations` module dir undeclared in modules.json (12 declared / 13 on disk) | hub grant #1 → entry committed isolado @ 10ed323b; tooling null-guard fixed by hub in main @ 449156f4 |
| GOV_MODULE_DEPENDENCY ×3 (inventory → listings/mutations via envelope.go) | F-02 envelope fold imports not in registry | hub grant #2a → inventory deps updated @ 31f85d3c (domain/ports layers only, no exception) |
| RCFG_UNDECLARED_READ MPC_PROVIDER_WRITES_ENABLED (root.go:640) | F-02 gate key never declared | hub grant #2b → runtime-config entry @ 8a53105c (owner modules/mutations, public, sole typed reader root.go) |
| RCFG_UNAPPROVED_READER MPC_TEST_DATABASE_URL (lifecycle_test.go) | raw os.Getenv skip-guard | chip-side fix @ 2efbe8e6 — guard removed, file already integration-tagged, relies on approved typed reader (testpostgres.OpenPool) |

## L1

- **Full `go test ./...`** (migrations touched → full sweep): GREEN — 79 packages ok,
  0 failures, exit 0 @ 5e5c14ee. Prereq fix: runner_test.go canonical-migration count
  37→39 (chip migrations 0038/0039), commit 5e5c14ee — per profile §3 test-fixture
  convention (2c095f3b: count bump is part of the migration grant).
- **harness:unit**: GREEN (exit 0, status=passed) after cherry-pick 89892b5e (hub fix:
  drop `-Force` in Postgres.psm1:4 nested import — pre-existing repo-wide lane breakage,
  independently root-caused by chip probe, same defect CHIP-SAT reported).
- **harness:integration**: GREEN-with-allowlist @ session pg (39 migrations embedded,
  fresh `mpc_test_*` db) — sole failure `TestPhase1SmokeFlow` = ratified profile §2
  allowlist entry, cited not re-proven. Orders duplicate-identity flake (F-B) eliminated
  by cherry-pick d9a36a0a (hub fix: set containment vs positional assert; allowlist
  request DENIED in favor of the existing fix — correct outcome).
- **Web vitest**: 11/11 (4 files) via workspace run; sdk-runtime vitest 51/51.
- `-race`: unavailable on this machine (accepted; hub decides cross-machine run
  pre-mission-close).

## Final re-run — post-rebase + F-04 (evidence head `a093f0f0`)

Rebase onto main @ `79d6787f18916cf9906fd355fb8eae9b2bc3067a` (hub F-04 trigger): 55 commits
replayed, 3 additive conflicts resolved (openapi.yaml schemas, sdk index.ts, root_test.go),
migration union 41 → fixture bump ea19ac33. New BaseSha for governance =
`79d6787f18916cf9906fd355fb8eae9b2bc3067a` (merge-base vs main).

### L0 (@ a093f0f0)

- `go build ./...` + `go vet ./...` — GREEN (absolute GOCACHE, apps/server_core).
- Governance lane: **status=passed** in clean detached worktree (`Documents/mpc-govcheck-m03`
  @ a093f0f0), full 40-hex BaseSha 79d6787f..., zero violations, 23 baseline_exceptions
  (all pre-existing main registry entries; count grew 14→23 from M-02 upstream additions).

### L1 (@ a093f0f0)

- **Full `go test ./...`**: GREEN — 86 packages ok, 0 failures, exit 0.
- **harness:unit**: GREEN (status=passed, exit 0).
- **harness:integration**: GREEN-with-allowlist @ session pg (migrations_first=41 embedded,
  fresh db, container session-reuse) — sole failure `TestPhase1SmokeFlow` = ratified profile
  §2 allowlist entry, cited not re-proven. Orders flake absent (d9a36a0a fix holding).
- **Web vitest**: 25 files / 164 tests green (workspace run @ code-identical tree 1d7a5bad;
  a093f0f0 adds docs only) + `npm run build --workspace @marketplace-central/web` green.
- **sdk-runtime vitest**: 60/60 (grew 51→60: mutation query builders + M-02 upstream).
- `-race`: unavailable on this machine (accepted; hub decides cross-machine run
  pre-mission-close).

### F-04 slice evidence

S1 fe268655, S2 a84f0b44, S3 19075cd7, S4 42574869, S5 1d7a5bad — all cold-review ACCEPT
(D-66, D-69, D-71, D-73). Full trail in DISPATCH-LEDGER.md rows D-63..D-73.

## L2

Hub-owned seam (dev stack via docker compose). Chip does not execute. Pending hub
scheduling at milestone close.

## Commits added in this window (post-F-03-closeout 1ebb7cc3)

| Commit | What |
| --- | --- |
| 5e5c14ee | test(migrate): canonical migration inventory 37→39 |
| 10ed323b | chore(governance): register mutations module (hub grant #1) |
| 2efbe8e6 | test(mutations): drop raw MPC_TEST_DATABASE_URL skip-guard |
| 13c56b54 | cherry-pick 89892b5e — fix(harness): drop -Force nested import |
| 50df9d29 | cherry-pick d9a36a0a — fix(orders): de-flake duplicate-identity test |
| 31f85d3c | chore(governance): inventory deps + listings/mutations (hub grant #2a) |
| 8a53105c | chore(governance): declare MPC_PROVIDER_WRITES_ENABLED (hub grant #2b) |

## Corrective round 1 (M03-COR-1) — 2026-07-17

Anchors: before = bd87d15a (dual gate round 1 target, Sol FAIL) → after = 939ea614
(docs SHA; ladder run at this tree). Commits in range:

| Commit | What |
| --- | --- |
| 2ba61fa3 | fix(web): COR-1d failure code literal + totals-based canRetry |
| 94f8efb5 | fix(mutations): COR-1a/1b/1c envelope gate, retry clone scope, msg fallback |
| 75c623cb | docs(mnfs): COR-1f VC ruling (C04 422, C09 409) + PLAN-ADJUDICATION |
| d0ecee1c | fix(mutations): COR-1a review fix — ProviderObservedAt as source time |
| 654721c2 | test(mutations): distinct CreatedAt fixture per review note |
| 939ea614 | docs(mnfs): ledger D-74..D-77 + F-04 validation.md (COR-1e) |

Review: cold sonnet round 1 REJECT on COR-1a (gate validated approval-instant timestamp,
could never fire) → fixed d0ecee1c → re-review ACCEPT; COR-1b/1c/1d ACCEPT round 1
(ledger D-76).

### Ladder re-run @ 939ea614

- **L0**: `go build ./...` + `go vet ./...` green. **governance**: status=passed @
  939ea614 in clean detached worktree mpc-govcheck-m03, BaseSha
  79d6787f18916cf9906fd355fb8eae9b2bc3067a (`scripts\harness.ps1 governance`);
  baseline exceptions only (7, unchanged set).
- **Go sweep**: `go test -count=1 ./...` (server_core, GOCACHE=.gocache) — 86 pkgs ok,
  0 FAIL.
- **harness:unit**: status=passed (run b6277fa360e2427fa030a589c960bf6d).
- **harness:integration**: GREEN-with-allowlist (run 39cec6931f37477c97f2aaf4b9627441,
  ephemeral container port 54738, migrations_first=41) — sole failure
  `TestPhase1SmokeFlow` = ratified profile §2 allowlist entry, cited not re-proven.
  Earlier same-day session-pg run (df7bdc21… + re-run) identical allowlist-only result.
- **Web vitest**: 25 files / 165 tests green (grew 164→165: COR-1d out-of-page retry
  test) + `npm run build --workspace @marketplace-central/web` green.
- **sdk-runtime vitest**: 60/60.
- `-race`: unavailable on this machine (accepted; hub decides pre-mission-close).

L2 unchanged: hub-owned seam.

## Corrective round 2 (M03-COR-2) — 2026-07-17

Anchors: 688496c7 → 32df1d4a. Single item: poller head-of-line livelock on gate
failure (hub delta review @ 688496c7; COR-1b/1c/1d/1e/1f ACCEPT, COR-1a structure
accepted, error path defect).

| Commit | What |
| --- | --- |
| d5a85337 | fix: poller terminalizes gate-failed protocols (failed_preserved) |
| 32df1d4a | fix: drain paginated pending items on gate failure (review REJECT→fix) |

Review: cold sonnet REJECT round 1 (pagination stranding) → 32df1d4a → ACCEPT
(ledger D-79). Affected lanes re-run @ 32df1d4a: go vet + mutations/inventory
sweep 0 FAIL; build green; harness:unit status=passed; harness:integration
GREEN-with-allowlist (sole TestPhase1SmokeFlow, cited). Reviewer finding
(out of scope, reported to hub): totals never updated at finish → COR-1d
canRetry gating never shows retry button for failed protocols.
