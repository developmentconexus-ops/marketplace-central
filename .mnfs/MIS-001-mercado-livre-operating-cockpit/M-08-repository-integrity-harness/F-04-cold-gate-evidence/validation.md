# F-04 Validation

```yaml
feature_id: F-04
status: blocked
candidate_sha: 396c9d038f2e2c82f7f76f4e9cc73eb032d20ce5
context_pack: scripts/.runs/f52c2e54ac864dd1b6c7dbfcfbda9ec0/context-pack.json
```

## RED/GREEN evidence

- RED contract commits: `e94d642`, `20b7e09`.
- GREEN implementation commit: `9be30c2`.
- F-03 inventory consumption and cache alignment corrections: `0de0546`, `371d25c`, `396c9d0`.
- Contract suites passed before real execution: cold evidence, cold snapshot, and governance contracts.
- Context pack compiled and current-validated at each fixed candidate SHA; latest artifact is the path above.

## Real cold attempts

- Candidate `9be30c2`: two runs passed provisioning and fake contract stages with image identity `sha256:da788743d2060767375896de4d646f7576f5911461444b372616f19ea61db2ec`; F-03 regression failed in `postgres-go.tests.ps1` because isolated `GOPROXY=off` cache lacked `github.com/godror/godror@v0.49.4`. Run directories: `scripts/.runs/3907f68f7c5e409a93855dd1b80cee4e`, `scripts/.runs/34505b7ab62a4c1482bded6c422fbcf1`.
- Candidate `371d25c`: same F-03 cache failure after provisioning Go modules from `apps/server_core`.
- Candidate `396c9d0`: provisioning remained in isolated workspace Go module download without completing; run was interrupted after no output for more than two minutes. No success manifest was persisted.

## Safety/invariance

- Caller checkout was clean and remained clean after each attempt (`git status --short` empty).
- No provider, Oracle, browser, or dev database write was attempted.
- PostgreSQL tag resolved to the stable identity above and no labelled ephemeral resources were created by the cold gate attempts.
- Real acceptance is blocked: F-03 `godror` dependency provisioning in the isolated snapshot cache did not complete, so `32/0` migration and zero-resource evidence are unavailable.

Milestone Orchestrator must perform final SPEC/SAFETY and QUALITY reviews on a fixed candidate after the provisioning blocker is resolved. F-04 and M-08 are not passed.
