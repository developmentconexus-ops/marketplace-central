# F-04 Validation

```yaml
feature_id: F-04
status: blocked
candidate_sha: 396c9d038f2e2c82f7f76f4e9cc73eb032d20ce5
context_pack: scripts/.runs/f52c2e54ac864dd1b6c7dbfcfbda9ec0/context-pack.json
```

## RED/GREEN evidence

## Consolidated correction batch — 2026-07-11

- Status: `blocked` — correction contracts are green; Phase 4 has not run.
- Assigned P0/P1 findings addressed: schema-valid blocked evidence, detached
  snapshot and secret-free cold provisioning, mandatory inventory, F-03
  integration classification, redaction/projection/comparison, and current
  Context evidence semantics.
- RED observed before implementation:
  - `cold-gate-evidence.tests.ps1`: outcome schema/acceptance-link contract was
    not valid for the candidate representation.
  - `cold-gate-snapshot.tests.ps1`: no isolated `cold-provision` child
    environment, stable Docker/image blocks, or mandatory ordered inventory.
  - `cold-gate.integration.tests.ps1`: unconditional fake PASS; no F-03 real
    integration target/class record.
- GREEN targeted contracts:
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1` — exit 0, fake contract evidence.
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-evidence.tests.ps1` — exit 0, fake contract evidence.
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-snapshot.tests.ps1` — exit 0, fake contract evidence.
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate.integration.tests.ps1` — exit 0, F-03 integration-contract evidence (not live database evidence).
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1` — exit 0, fake context-contract evidence.
- Safe negative evidence: a dirty-caller invocation wrote schema-valid blocked
  `outcome.json` and `trace.jsonl` under
  `scripts/.runs/1ff8067dd98040a2933003d30558c918/`; it records
  `COLD_CALLER_DIRTY`, `dirty: true`, a repository-relative run path, and no
  provisioning or provider/Oracle/browser action.
- Remaining blocker: the F-03 regression command launched its isolated child
  and printed `contamination_count=0` plus `PASS postgres contract tests
  target=fake`, but did not produce its terminal F-03 pass marker in this
  session. Its exit/terminal result must be reproduced conclusively before
  compiling a fixed-HEAD current pack or authorizing Phase 4. No two-run cold
  evidence, caller-invariance proof, or Docker resource proof is claimed.

- RED contract commits: `e94d642`, `20b7e09`.
- GREEN implementation commit: `9be30c2`.
- F-03 inventory consumption and cache alignment corrections: `0de0546`, `371d25c`, `396c9d0`.
- Contract suites passed before real execution: cold evidence, cold snapshot, and governance contracts.
- Context pack compiled and current-validated at each fixed candidate SHA; latest artifact is the path above.

## Prior real cold attempts

- Candidate `9be30c2`: two runs passed provisioning and fake contract stages with image identity `sha256:da788743d2060767375896de4d646f7576f5911461444b372616f19ea61db2ec`; F-03 regression failed in `postgres-go.tests.ps1` because isolated `GOPROXY=off` cache lacked `github.com/godror/godror@v0.49.4`. Run directories: `scripts/.runs/3907f68f7c5e409a93855dd1b80cee4e`, `scripts/.runs/34505b7ab62a4c1482bded6c422fbcf1`.
- Candidate `371d25c`: same F-03 cache failure after provisioning Go modules from `apps/server_core`.
- Candidate `396c9d0`: provisioning remained in isolated workspace Go module download without completing; run was interrupted after no output for more than two minutes. No success manifest was persisted.

## Safety/invariance

- Caller checkout was clean and remained clean after each attempt (`git status --short` empty).
- No provider, Oracle, browser, or dev database write was attempted.
- PostgreSQL tag resolved to the stable identity above and no labelled ephemeral resources were created by the cold gate attempts.
- Real acceptance is blocked: F-03 `godror` dependency provisioning in the isolated snapshot cache did not complete, so `32/0` migration and zero-resource evidence are unavailable.

Milestone Orchestrator must perform final SPEC/SAFETY and QUALITY reviews on a fixed candidate after the provisioning blocker is resolved. F-04 and M-08 are not passed.
