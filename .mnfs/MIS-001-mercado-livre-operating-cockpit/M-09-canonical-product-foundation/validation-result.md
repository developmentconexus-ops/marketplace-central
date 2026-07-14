# M-09 Validation Result

```yaml
milestone: M-09
mode: proportional_qa
verdict: passed
frozen_sha: 2eabecbc806d652543b60583e036d2e3e686e10e
observed_at: 2026-07-14T22:09:52Z
sha_drift: false
supersedes: historical FAIL rollup at 32b32f6de00875589468c71eb70c6eb3e5d49278 (preserved in git history)
contract: validation-contract.md (M-09-C01..C05, QA-2)
owner: QA Validator
```

## Verdict

**PASS.** All five required criteria are satisfied at the frozen SHA
`2eabecbc`. The prior FAIL (C03 active residue at 32b32f6d) and the prior
externally_blocked C05 (runner stdout/stderr deadlock at a1d4aedd) are both
resolved and re-proved: the exact active-residue scan is clean at this SHA,
and the governed read-only Oracle lane completed with a positive CODPROD at
this exact SHA.

## Frozen-SHA integrity

- `git rev-parse HEAD` = `2eabecbc806d652543b60583e036d2e3e686e10e` before and
  after every QA lane in this session (no drift).
- `a1d4aedd` (M-09 normalization/acceptance SHA) and `314b1ef3` (MIS-002
  mission QA SHA) both verified as ancestors of the frozen SHA.
- Working tree clean on all active paths (`apps/`, `packages/`, `docker/`,
  `scripts/`, `contracts/`) at validation time.

## Criteria

| Criterion | Assessment | Evidence |
| --- | --- | --- |
| M-09-C01 | PASS | Accepted at a1d4aedd (fixed-SHA review PASS, no findings); re-proved by MIS-002 mission QA at 314b1ef3 (full Go ladder exit 0; `.mnfs/MIS-002-oracle-read-rearchitecture/validation-result.md` PASS); `git diff --name-only 314b1ef3..2eabecbc -- apps/ packages/` verified EMPTY this session; corroboration re-run this session at 2eabecbc: targeted Go lane (catalog, internal_read, product_links) all `ok`, SDK lane 41/41 PASS. |
| M-09-C02 | PASS | Same reuse chain (empty executable diff verified); corroboration re-run at 2eabecbc: `internal_read` domain/adapters/fake/oracle packages all `ok`, including nullable-fact and quality-flag coverage; MIS-002 C02 re-ran `Nullable|Quality` tests green at 314b1ef3. |
| M-09-C03 | PASS | Exact active-residue scan re-executed by QA at 2eabecbc: `git grep -E "MetalShopping\|MS_DATABASE_URL\|MS_TENANT_ID\|platform/msdb"` over `apps/`, `packages/`, `docker/`, `contracts/governance/`, `scripts/`, config/compose paths — the three previously failing locations (`catalog/ports/repository.go`, `catalog/application/service.go`, `docker/dev/README.md`) have ZERO matches; `internal/platform/msdb` package absent from the tree. Only remaining hits are contamination-GUARD fixtures in `scripts/tests/harness-environment.tests.ps1` (denylist asserting MS_* vars are scrubbed from child environments — anti-residue enforcement, not a dependency) and non-active historical docs/ADR/plan wording. |
| M-09-C04 | PASS | Same reuse chain; corroboration re-run at 2eabecbc: `internal/platform/migrate`, `migrations`, `product_links`, `classifications`, `pricing`, `composition` packages all `ok` (mapped / `not_found` / `identity_conflict` coverage per feature evidence; no guessing path). |
| M-09-C05 | PASS | Fresh governed read-only lane evidence AT the frozen SHA: `_fixed-sha-qa/c05-oracle-evidence.md` — `status=passed`, `exit_code=0`, `read_only=true`, `positive_codprod_observed=true`, source `oracle/sankhya`, observed 2026-07-14T22:00:25Z, selector `^TestOracleLiveSmoke$/^product_lookup$` (SELECT-only); companion `-EmitBaseline` proof at f099a52c (docs/governance-only diff to 2eabecbc). Runner deadlock fix proven: runner contract suite re-run this session at 2eabecbc — 18/18 assertions PASS, exit 0, including "drains full stdout and stderr pipes concurrently without deadlocking" and timeout termination. Per constraint, the live lane was not re-executed by QA; the governed evidence file at the exact frozen SHA is the accepted proof. |

## QA lanes executed this session (at 2eabecbc)

1. Active-residue scan (`git grep`, exact terms, active paths) — PASS, zero
   dependency matches.
2. Targeted Go lane, `apps/server_core`, absolute `GOCACHE`, `-count=1`:
   catalog, internal_read (incl. cache/fake/oracle/oraclebatch/observability),
   product_links — all `ok`.
3. Second Go lane: platform/migrate, migrations, classifications, pricing,
   composition, cmd — all `ok`, exit 0.
4. SDK lane `npm test --workspace @marketplace-central/sdk-runtime -- --run` —
   PASS 41/41.
5. Governed runner contract `scripts/tests/live-oracle-docker-runner.tests.ps1`
   — PASS 18/18, exit 0 (suite grew from 14 with the MIS-002 deadlock-fix
   contract tests).
6. Diff-emptiness verification `314b1ef3..2eabecbc -- apps/ packages/` — empty;
   non-docs changes since 314b1ef3 confined to governance contracts and the
   Oracle runner seam, each covered by lanes 5 and C05 evidence plus the
   closeout dual reviews (Codex gpt-5.6-sol + independent Claude, PASS) on
   1a7a9389 / f099a52c / 2eabecbc.

## Residual notes (non-blocking)

- `scripts/tests/governance-contracts.tests.ps1` could not be re-run by QA:
  it shells to `rg`, which is not on PATH on this machine (known environment
  limitation). Accepted evidence: the closeout dual reviews covering the
  governance-envelope commit 2eabecbc passed with no blocking findings.
- Live UI drive is not applicable to this milestone's contract scope; the
  contract's runnable surface (governed Oracle read) is proven by the C05
  governed-lane evidence at the frozen SHA.
- No Oracle/provider/database write, credential or raw-row exposure, source
  edit, or destructive Git action occurred in this session.

## Terminal disposition

- status: `passed`
- blocker: none
- next owner: Milestone Orchestrator (M-09 closeout), then MIS-001 rollup.
