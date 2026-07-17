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
