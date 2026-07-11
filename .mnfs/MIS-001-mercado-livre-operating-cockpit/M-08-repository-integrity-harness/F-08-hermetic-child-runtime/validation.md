# F-08 Hermetic Child Runtime — Validation

```yaml
id: F-08
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-08-hermetic-child-runtime

## Summary

Focused validation passed from clean base
`68174966e61af5b6474a7f147e9a98fdcd57bdb4`. Environment, typed process,
lane, alias, governance, context, and normal unit behavior all exited `0`.
Normal unit execution passed one Go package plus 20 web test files containing
185 tests. Current governance reported `status=passed` and exactly 35 declared
baseline exceptions. Final context compilation/current-base validation passed.

All evidence targets are `fake`/local contract execution. No database, Oracle,
Mercado Livre, external network, browser, migration, or provider write ran.
This is Feature Implementer quick validation, not feature acceptance, formal
milestone QA, or an M-08 pass.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-11
- Final feature state: `quick_validation_passed`
- Base SHA: `68174966e61af5b6474a7f147e9a98fdcd57bdb4`
- Target: `fake`
- Spec adherence: all `F08-AC01` through `F08-AC05` proof mappings ran and
  passed; no focused-scope deviation observed.

## Evidence Honesty

- Passing commands below ran this Phase 4 session and name observable output
  plus an artifact path.
- Fixture, governance, cache, context, and unit checks prove deterministic local
  contract behavior only. They are not live integration evidence.
- The context pack stores pre-dispatch command evidence as `assumed`; this file
  separately records the actual Phase 4 executions as `ran`.

## Quick Validation State

- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None observed in focused F-08 scope.
- Reason: Every acceptance criterion has at least one executed proof command;
  external systems remain explicit non-goals.

## Commands Run

Every passing entry has evidence type `ran`, owner Feature Implementer, target
`fake`, and blocking condition `None`.

### Focused runtime fixtures

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-environment.tests.ps1`
  - Expected/actual: exit `0`, `PASS harness environment tests`.
  - Artifact: `scripts/tests/harness-environment.tests.ps1`.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-execution.tests.ps1`
  - Expected/actual: exit `0`, `PASS harness execution tests`.
  - Artifact: `scripts/tests/harness-execution.tests.ps1`.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/hermetic-lanes.tests.ps1`
  - Expected/actual: exit `0`, `PASS hermetic lane tests`.
  - Artifact: `scripts/tests/hermetic-lanes.tests.ps1`.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-aliases.tests.ps1`
  - Expected/actual: exit `0`, `PASS harness alias tests`; all eight root npm
    harness aliases exercised behavior from a foreign CWD.
  - Artifact: `scripts/tests/harness-aliases.tests.ps1`.

### Governance and context regressions

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
  - Expected/actual: exit `0`, `PASS governance contract tests`.
  - Artifact: `scripts/tests/governance-contracts.tests.ps1`.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1`
  - Expected/actual: exit `0`, `PASS governance drift tests`.
  - Artifact: `scripts/tests/governance-drift.tests.ps1`.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
  - Expected/actual: exit `0`, `PASS context compiler tests`.
  - Artifact: `scripts/tests/context-compiler.tests.ps1`.
- `npm run harness:governance -- -BaseSha 68174966e61af5b6474a7f147e9a98fdcd57bdb4`
  - Expected/actual: exit `0`, `status=passed`, exactly 35 named
    `baseline_exception` records, `artifact_path=contracts/governance`.
  - Artifact: `contracts/governance/`.

### Normal unit execution

- `npm run harness:unit`
  - Expected: exit `0`, fake target, disabled DB/Oracle/provider/migrations,
    repository-root execution, redacted minimal summary.
  - Actual: exit `0`; Go `marketplace-central/apps/server_core/tests/unit`
    passed; Vitest reported 20 files and 185 tests passed; summary reported
    `target=fake`, `status=passed`.
  - Artifact: `scripts/.runs/57093563a36c420d9639693d8d28adb3/summary.txt`.
  - Non-blocking existing diagnostics: Node `DEP0205` warnings and React
    `act(...)` warnings appeared; neither changed exit/result.

### Final current context

- F-08 context compile used current HEAD, the exact `allowed_paths` array from
  the Machine Work Contract, and the F-08 feature path.
  - Expected/actual: exit `0`, `status=passed`.
  - Artifact: `scripts/.runs/895794e49f6e4654848caf862b088f84/context-pack.json`.
- `npm run harness:context:validate -- -ContextPath 'scripts/.runs/895794e49f6e4654848caf862b088f84/context-pack.json' -RequireCurrentBase`
  - Expected/actual: exit `0`, `status=passed`.
  - Artifact: same context pack.
- Pack observables: base SHA `68174966e61af5b6474a7f147e9a98fdcd57bdb4`,
  24 hashed sources, 8 commands, all target `fake`, risk `L2`, 1,935 estimated
  input tokens, and retry budget 1.

### Local Go cache risk inventory

- Command: set child-only `GOMODCACHE` to
  `apps/server_core/.gomodcache`, set `GOPROXY=off`/`GOSUMDB=off`, then run
  `go list -m all` from `apps/server_core`; enumerate local
  `cache/download/**/*.zip` files.
  - Evidence type: `ran`; target: `fake`/offline local cache.
  - Expected/actual: exit `0`; module graph contained 32 dependencies while
    local cache contained 24 dependency ZIPs, leaving 8 absent.
  - Absent nonfocused ZIPs: `github.com/golang/protobuf`, `github.com/kr/pretty`,
    `github.com/planetscale/vtprotobuf`, `github.com/stretchr/objx`,
    `golang.org/x/mod`, `golang.org/x/net`, `golang.org/x/tools`, and
    `golang.org/x/xerrors` at graph-selected versions.
  - Artifact: ignored `apps/server_core/.gomodcache/` inventory; no network or
    download attempted.

## Acceptance Mapping

- `F08-AC01` / `M-08-C02`: environment and lane fixtures passed structural
  non-inheritance, safe-set, repository-local cache, EnvFile-ignore, and parent
  immutability cases.
- `F08-AC02` / `M-08-C02`: execution fixtures passed literal argument,
  explicit-CWD, concurrent >2 MiB stream, exit `17`, timeout/tree-kill, start
  error, and parent-state cases.
- `F08-AC03` / `M-08-C04`: execution/lane fixtures passed raw, assigned,
  embedded, URI, stdout, stderr, dispatcher, and summary redaction cases.
- `F08-AC04` / `M-08-C02`,`M-08-C04`: alias/lane fixtures passed all eight
  aliases and fail-before-process/run-directory behavior.
- `F08-AC05` / `M-08-C09`: contract, drift, context, current governance, and
  current F-08 context checks passed without scanner exemption or evidence
  inflation.

## RED/GREEN, Replans, and Reviews

- Planning/spec contract: `4beb9d7` (`docs(m08): plan hermetic child runtime`).
- Phase 1 RED: `979f971` (`test(harness): define hermetic child contract`)
  committed adversarial fixtures before production modules. Failures were the
  intended absent fresh-child/typed-runtime behavior, not accepted as GREEN.
- Phase 2 GREEN: `1bce008` (`feat(harness): add hermetic child runtime`). One
  fixed-commit independent contract/quality review found bounded process/env
  edge cases; one correction batch `f141b62` (`fix(harness): close child runtime
  edge cases`) passed focused re-review. Verdict: Pass.
- Phase 3 stopped before shared-seam implementation when source ownership,
  registry reader shape, runtime precedence, and offline module-cache isolation
  needed explicit plan decisions. Replans were recorded in `af2a720`,
  `d3b5461`, `1c8852c`, and `558fb9f`; these were architecture stops, not
  implementation correction loops.
- Dispatcher GREEN: `647ee9a` (`refactor(harness): route lanes through typed
  runtime`). Review scope/source bindings were clarified in `6b2b7ae` and
  `745e7af`. One independent fixed-commit review found typed dispatcher,
  generic-reader enforcement, redaction, and negative-fixture gaps; one
  correction batch `6817496` (`fix(harness): harden typed dispatcher boundary`)
  passed focused independent re-review. Verdict: Pass.
- Phase 4 used no implementation correction. This document records fresh
  execution, not inherited review claims.

## Changes Made

- Contracts/docs: `.gitignore`, F-08 `feature.md`/`spec.md`/`plan.md`,
  `contracts/governance/runtime-config.json`, and its schema.
- Runtime: `scripts/harness.ps1`, `Environment.psm1`, `Execution.psm1`, and
  `Policy.psm1`.
- Fixtures: child probe; environment, execution, lane, alias, governance
  contract, and governance drift suites.
- No product code, OpenAPI, SDK, migration, Docker, provider adapter, or UI
  implementation changed.

## Evidence Honesty and Remaining Risks

- Current 35 governance exceptions remain named debt with registry owners.
  Passing governance proves no unknown growth, not full repository conformance.
- Repository-local Go cache proves focused `tests/unit` offline execution only.
  Eight module ZIPs needed by nonfocused server packages remain absent from the
  local cache. F-04 owns deterministic cold provisioning and full cold-gate
  evidence; F-08 does not claim cold-repository completeness.
- Minimal `summary.txt` is intentional. F-04 owns rich evidence manifests and
  aggregation; F-05 owns leases/run state/worktrees.
- No live integration, PostgreSQL lifecycle, Oracle, Mercado Livre, browser,
  OAuth, network, migration, or provider-write behavior was validated.
- `quick_validation_passed` still requires Milestone Orchestrator independent
  acceptance. M-08 remains not passed.

## Manual QA

- QA level: QA-3 deterministic feature gate.
- Flow or step: browser/live manual flow.
- Status: Not run.
- Evidence type: `assumed`.
- Owner: Feature Implementer.
- Expected: Not applicable to F-08 non-UI fake/contract scope.
- Actual: Automated deterministic evidence completed as recorded above.
- Blocking condition: None for feature quick validation; browser/live evidence
  remains outside F-08 and cannot be inferred from fixtures.

## Handoff

- Current status: `quick_validation_passed`.
- Next owner: Milestone Orchestrator.
- Next action: review fixed commits, feature/spec/plan, changed paths, current
  focused evidence, and cold-cache risk; accept, reject, or block F-08 before
  releasing shared seams to F-03/F-04/F-05.
- Required evidence: commits `4beb9d7` through the Phase 4 evidence commit,
  this file, eight focused suites, current governance output, normal unit
  summary, and final context pack.
- Blockers: no F-08 blocker. F-04 cold provisioning and M-08 formal validation
  remain outstanding.

## Milestone Orchestrator Acceptance

- Decision: `accepted`.
- Accepted implementation/evidence HEAD:
  `43f59173009c3c908232129e029ef057158f159d`.
- Fresh acceptance commands from that clean HEAD:
  - environment, execution, lane, alias, governance-contract, governance-drift,
    and context-compiler suites: exit `0`, `PASS`;
  - current `harness:governance`: exit `0`, `status=passed`, exactly 35
    declared baseline exceptions and no unknown drift;
  - post-evidence F-08 context compile/current-base validation: exit `0`,
    artifact
    `scripts/.runs/10f604a8e016408b921454e94b79d50d/context-pack.json`;
  - normal offline `harness:unit`: exit `0`, Go unit passed, Vitest 20 files / 185
    tests passed, summary
    `scripts/.runs/ceb994d242594df7a1cb5894b22feb2a/summary.txt`;
  - `git diff --check` and visible status: clean before this acceptance update.
- Independent review: Phase 2 and Phase 3 each used one fixed-commit combined
  SPEC/QUALITY review, one consolidated correction batch, and one focused
  revalidation. Both final verdicts were `Pass`.
- Evidence class: `fake`/local contract only. No PostgreSQL, Oracle, Mercado
  Livre, browser, network, migration, or provider write was run or claimed.
- Accepted limitation: the repo-local Go cache proves the focused unit graph;
  F-04 must still prove cold dependency provisioning and the complete cold gate.
- Next owner: F-03 Feature Implementer.
