# F-08 Hermetic Child Runtime — Plan

```yaml
id: F-08
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Scope and Ownership

F-08 is a serial L2 shared-seam refactor. One writer owns
`scripts/harness.ps1`, `scripts/harness/**`, and root harness aliases. F-03,
F-04, and F-05 remain blocked from these seams until F-08 acceptance.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-08",
  "required_sources": [
    "contracts/governance/execution-lanes.json",
    "contracts/governance/runtime-config.json",
    "contracts/governance/schemas/execution-lanes.schema.json",
    "contracts/governance/schemas/runtime-config.schema.json",
    "docs/superpowers/specs/2026-07-10-repository-native-agent-harness-design.md",
    "scripts/harness.ps1",
    "scripts/harness/Environment.psm1",
    "scripts/harness/Execution.psm1",
    "scripts/harness/Policy.psm1",
    "scripts/tests/harness-aliases.tests.ps1",
    "scripts/tests/harness-environment.tests.ps1",
    "scripts/tests/hermetic-lanes.tests.ps1",
    "scripts/tests/governance-contracts.tests.ps1",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-02-hermetic-execution-lanes/validation.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler/validation.md"
  ],
  "allowed_paths": [
    "scripts/harness.ps1",
    "scripts/harness/**",
    "scripts/tests/hermetic-lanes.tests.ps1",
    "scripts/tests/harness-environment.tests.ps1",
    "scripts/tests/harness-execution.tests.ps1",
    "scripts/tests/harness-aliases.tests.ps1",
    "scripts/tests/fixtures/harness/**",
    "contracts/governance/runtime-config.json",
    "contracts/governance/schemas/runtime-config.schema.json",
    "scripts/tests/governance-contracts.tests.ps1",
    "package.json",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-08-hermetic-child-runtime/**"
  ],
  "forbidden_paths": [
    "apps/**",
    "packages/**",
    "contracts/api/**",
    "migrations/**",
    "docker-compose.yml"
  ],
  "side_effects": {
    "allowed": ["repository-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id":"environment-fixtures","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-environment.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"execution-fixtures","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-execution.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"lane-fixtures","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/hermetic-lanes.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"alias-fixtures","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-aliases.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"governance-contracts","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"governance-drift","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"governance-current","command_template":"npm run harness:governance -- -BaseSha {base_sha}","lane_id":"unit","expected_exit_code":0},
    {"id":"context-regression","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1","lane_id":"unit","expected_exit_code":0}
  ],
  "criteria": [
    {"id":"F08-AC01","milestone_criterion_id":"M-08-C02","command_ids":["environment-fixtures","lane-fixtures"]},
    {"id":"F08-AC02","milestone_criterion_id":"M-08-C02","command_ids":["execution-fixtures"]},
    {"id":"F08-AC03","milestone_criterion_id":"M-08-C04","command_ids":["execution-fixtures","lane-fixtures"]},
    {"id":"F08-AC04","milestone_criterion_id":"M-08-C02","command_ids":["alias-fixtures","lane-fixtures"]},
    {"id":"F08-AC05","milestone_criterion_id":"M-08-C09","command_ids":["governance-contracts","governance-drift","governance-current","context-regression"]}
  ],
  "stop_conditions": [
    {"code":"contract-drift","condition":"Implementation requires a runtime key or lane behavior not authorized by both governance registries."},
    {"code":"external-side-effect","condition":"A fixture or implementation attempts database mutation, external network access, or provider write."},
    {"code":"scope-conflict","condition":"An unrelated dirty path or competing shared-seam owner appears."}
  ],
  "retry_budget": {"max_correction_attempts":1},
  "handoff_fields": ["commit-sha","changed-paths","red-evidence","green-evidence","target-labels","remaining-risks"]
}
```

## Phase 1 — Adversarial RED Contract

### Files

- Modify `scripts/tests/hermetic-lanes.tests.ps1`.
- Create `scripts/tests/harness-environment.tests.ps1`.
- Create `scripts/tests/harness-execution.tests.ps1`.
- Create `scripts/tests/harness-aliases.tests.ps1`.
- Create minimal fixtures under `scripts/tests/fixtures/harness/`.

Write tests before production code. Use a Node child probe for exact environment,
CWD, arguments, streams, exit, volume, and timeout behavior; PowerShell itself
injects `PSModulePath` and cannot prove an exact environment set.

Required RED cases:

1. Parent contains all prior F-02 keys, `MPC_PRODUCT_LINKS_POSTGRES_URL`,
   `ME_CLIENT_ID`, `ME_CLIENT_SECRET`, `ME_REDIRECT_URI`, foreign `GOCACHE`,
   and a GUID-suffixed unknown sentinel. Unit must eventually succeed while the
   child sees none of them.
2. Unit ignores an EnvFile containing those keys and legacy Oracle aliases.
3. Child sees only fixed safe host/tool keys plus canonical `GOCACHE`; parent
   environment and CWD remain unchanged on success and forced failure.
4. Arguments containing spaces, quotes, `&`, `$()`, and semicolons round-trip
   literally; no shell evaluation occurs.
5. Concurrent stdout/stderr greater than 2 MiB completes; exit `17` remains
   `17`; timeout is stable and kills the process tree.
6. Exact credentials appear in raw, assigned, embedded, URI, stdout, and stderr
   forms; every propagated representation is redacted.
7. Absolute direct invocation and all eight npm aliases execute from a foreign
   CWD with expected ready/blocked classifications.
8. Unsupported lane/target fails before a marker process and run directory.

Run all four suites and capture exit/output. Intended RED must be caused by the
missing fresh-child/module behavior, not syntax, fixture, or assertion defects.
Commit tests only as:

```text
test(harness): define hermetic child contract
```

## Phase 2 — Environment and Execution GREEN

### Files

- Create `scripts/harness/Environment.psm1`.
- Create `scripts/harness/Execution.psm1`.
- Modify only the Phase 1 fixtures needed to correct proven test defects.

Implement typed `New-HarnessChildEnvironment`, process request/result creation,
and `Invoke-HarnessProcess` exactly as specified. Resolve tools before clearing
the child environment. Start both asynchronous drains before wait. Redact before
return. Do not mutate global environment or location and do not invoke via `&`,
`Start-Process`, `.Arguments`, or constructed shell strings.

Run environment and execution suites to GREEN. Commit:

```text
feat(harness): add hermetic child runtime
```

Use one fixed-commit combined contract/quality review. Consolidate actionable
findings into at most one correction batch and rerun only affected suites.

## Phase 3 — Dispatcher and Alias Migration

### Files

- Modify `scripts/harness.ps1`.
- Modify `scripts/tests/hermetic-lanes.tests.ps1`.
- Modify `scripts/tests/harness-aliases.tests.ps1`.
- Modify `package.json` only if a delegate defect is proven; alias names and
  stable parameter surfaces do not change.

Route unit subprocesses through the typed runtime. Move env-file parsing and
safe alias resolution into `Environment.psm1`. Keep Policy/Context behavior,
minimal run summaries, F-03 integration blocker, live/browser preflight, and
provider-write pre-network guards. Validate inputs before starting a child or
creating a run directory. Every path derives from `$PSScriptRoot` or explicit
repository root.

If migration exposes the generic registry-driven environment boundary, amend
the runtime-config schema, registry, Policy semantic checks, and governance
fixtures atomically. Represent that boundary with an explicit reader kind;
Policy must validate its declared path and bounded registry use without key
literals. Do not remove reader traceability, duplicate registry keys in code,
or add a scanner exemption. Process environment precedence is allowed only for
the selected lane's `allowed_runtime_keys`; unit still reads none.

Run all F-08 suites plus governance/context regressions. Commit:

```text
refactor(harness): route lanes through typed runtime
```

Use one fixed-commit combined review, one correction batch maximum, and focused
revalidation. A new architectural finding triggers replan instead of review
loops.

## Phase 4 — Evidence and Feature Handoff

Create `validation.md`; update `feature.md` only after a clean complete focused
gate. Recompile an F-08 context pack from final HEAD and validate it with
`-RequireCurrentBase`. Record exact commands/exits, target `fake`, RED/GREEN
history, commits, safe artifact paths, review verdict, and remaining risks.
No fixture is live evidence and no F-08 result can pass M-08 by itself.

Commit:

```text
docs(m08): record hermetic runtime evidence
```

## Verification Mapping

- F08-AC01: environment and hermetic-lane fixtures.
- F08-AC02: execution fixtures.
- F08-AC03: execution redaction plus lane artifact checks.
- F08-AC04: alias behavior and foreign-CWD fixtures.
- F08-AC05: governance/context suites, current governance, and current F-08
  context compile/validate.

## Handoff

- Current status: `planned`.
- Next owner: Fresh Phase 1 Feature Implementer.
- Next action: Create RED fixtures only and return exact failing evidence.
- Review policy: L2, one combined review per fixed implementation phase.
- Evidence class: fake/contract only.
- Blockers: None.
