# Marketplace Central — Agent Bootstrap

## How work lands

Branch, commit, open a PR, get it green, get it reviewed, and the operator merges.
There is no other channel. Concretely:

- One branch per change. Conventional-commit PR title (`pr-title.yml` enforces the
  type list: feat fix docs style refactor perf test build ci chore revert).
- **Never `git push` without explicit operator permission**, and never merge — merge
  is the operator's event.
- CI (`.github/workflows/ci.yml`) must be green: `lint-go`, `lint-frontend`,
  `build-go`, `test-go-unit`, `test-frontend`, `verify-full`, aggregated by
  `required`. CodeRabbit reviews the PR. Red is a stop, not a note.
- Work is tracked in GitHub issues and PRs. Nothing else is a queue.

## Verification

`scripts/gate.ps1` is the single implementation of every lane — CI invokes the same
file a developer machine does, and no lane logic lives in a workflow. Two checks live
in `.github/workflows/ci.yml` by nature, because their subject is the workflow itself:
the `required` aggregator with its topology assertion over the job/`needs`/`EXPECTED`
enumerations, and the preflight run of the gate's own counter fixtures
(`scripts/tests/gate-measure.tests.ps1`) ahead of the first lane. A local run is
evidence about the lanes, not about those two.

```bash
npm run gate        # -Lane all: the 13 everyday lanes, seconds
npm run gate:full   # adds selftest, guards, integration, edge — matches CI verify-full
```

Counted checks are ratchets in `contracts/gate/baselines.json`: a number may shrink,
never grow, and an unknown key fails. **Do not raise a baseline to make a lane green** —
that is the defect the ratchet exists to catch. Declared rules live in
`contracts/governance/invariants.json`; a rule that needs an exception gets one there
with a reason, or gets fixed.

Presence is not execution. A test that is skipped and a test that passed are
byte-identical in a log unless the lane prints what it ran — assert the count or the
token, never the exit code alone.

## Repository truth, in order

`ARCHITECTURE.md` and the ADRs under `docs/architecture/decisions/`, then the OpenAPI
spec plus the SDK, then `contracts/governance/`, then the wiki, then tests, builds and
commits. Stop and classify architecture, contract, runtime, ownership, or verification
conflicts instead of picking a side silently.

`.mnfs/` is a frozen archive of the retired hub-and-chips process. It is history, not
authority — read it for context, never cite it as a rule, and do not add to it.

## Architecture rules

Keep domain/application/ports/adapters/transport boundaries; tenant queries scope
`tenant_id`; provider payloads remain at adapters; unknown operational facts never
become zero or a default. API changes update OpenAPI and `sdk-runtime` in the same
commit. Provider writes need resolved linkage, explicit policy/source time, duplicate
protection, and audit. Mocks prove contract behavior, never live integration. Live ML
writes require explicit operator authorization.

## Operational rules

One writer owns a checkout or a shared seam. Do not reset, revert, stash, clean, delete
unknown state, use WSL, expose secrets or PII, cold-clone, purge caches, or install
dependencies as a ritual — a dependency change is its own PR, decided by the operator.

Go commands run from `apps/server_core`, with `GOCACHE` and `GOMODCACHE` bound to
**absolute** paths. A relative value makes `go env GOCACHE` report `off` and every Go
build die with `build cache is required, but could not be located: GOCACHE is not an
absolute path` — 83 bytes and zero `=== RUN`, which reads like a suite with nothing in
it. `scripts/gate.ps1` binds them for every lane; a hand-run needs the same thing,
resolved rather than transcribed:

```powershell
cd apps/server_core
$env:GOCACHE = (Resolve-Path .).Path + '/.gocache'
$env:GOMODCACHE = (Resolve-Path .).Path + '/.gomodcache'
```

Evidence is the PR: the diff, the CI run, and what the lanes printed. A claim with no
run behind it did not happen.
