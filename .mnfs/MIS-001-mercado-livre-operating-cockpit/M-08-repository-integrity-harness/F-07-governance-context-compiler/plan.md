# F-07 Governance Registry and Context Compiler — Plan

```yaml
id: F-07
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
split_decision: split
split_reason: governance contracts, repository drift, and context compilation are three reviewable serial deliverables across shared seams
```

## Feature ID

F-07-governance-context-compiler

## Interfaces Frozen For All Phases

- JSON Schema draft: `https://json-schema.org/draft/2020-12/schema`.
- Registry version: `1.0`.
- Schema validation: `Test-Json -LiteralPath <json> -SchemaFile <schema>`.
- Public module functions and reason codes are exactly those in `spec.md`.
- All paths serialize repo-relative with `/`; rooted paths and `..` fail.
- Root dispatcher alone writes stdout and maps result to exit 0/1.
- F-08 owns child-process environment and must not be implemented here.
- Registry records must follow
  `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/m08-governance-runtime-inventory.md`;
  any mismatch with active source blocks the phase instead of inventing a value.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-07",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/research/m08-governance-runtime-inventory.md",
    "docs/superpowers/specs/2026-07-10-repository-native-agent-harness-design.md"
  ],
  "allowed_paths": [
    "AGENTS.md",
    "package.json",
    "contracts/governance/**",
    "scripts/harness.ps1",
    "scripts/harness/**",
    "scripts/tests/governance-contracts.tests.ps1",
    "scripts/tests/governance-drift.tests.ps1",
    "scripts/tests/context-compiler.tests.ps1",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler/**"
  ],
  "forbidden_paths": [
    "apps/server_core/internal/modules/**",
    "apps/server_core/migrations/**",
    "apps/web/src/**",
    "packages/sdk-runtime/**",
    "contracts/api/**",
    "docker/**",
    "docker-compose.yml"
  ],
  "side_effects": {
    "allowed": [],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {
      "id": "governance-contracts",
      "command_template": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1",
      "lane_id": "unit",
      "expected_exit_code": 0
    },
    {
      "id": "governance-drift-fixtures",
      "command_template": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1",
      "lane_id": "unit",
      "expected_exit_code": 0
    },
    {
      "id": "context-fixtures",
      "command_template": "pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1",
      "lane_id": "unit",
      "expected_exit_code": 0
    },
    {
      "id": "governance-current",
      "command_template": "npm run harness:governance -- -BaseSha {base_sha}",
      "lane_id": "unit",
      "expected_exit_code": 0
    },
    {
      "id": "context-current",
      "command_template": "npm run harness:context:validate -- -ContextPath {context_path} -RequireCurrentBase",
      "lane_id": "unit",
      "expected_exit_code": 0
    }
  ],
  "criteria": [
    {"id": "F07-AC01", "milestone_criterion_id": "M-08-C09", "command_ids": ["governance-contracts", "governance-current"]},
    {"id": "F07-AC02", "milestone_criterion_id": "M-08-C09", "command_ids": ["governance-drift-fixtures", "governance-current"]},
    {"id": "F07-AC03", "milestone_criterion_id": "M-08-C09", "command_ids": ["context-fixtures", "context-current"]},
    {"id": "F07-AC04", "milestone_criterion_id": "M-08-C09", "command_ids": ["context-fixtures"]},
    {"id": "F07-AC05", "milestone_criterion_id": "M-08-C09", "command_ids": ["context-fixtures"]}
  ],
  "stop_conditions": [
    {"code": "scope-conflict", "condition": "A requested write path is not allowed or intersects a forbidden path."},
    {"code": "contract-drift", "condition": "A schema, registry, source hash, criterion proof, or lane reference is inconsistent."},
    {"code": "external-side-effect", "condition": "Execution would access a real database, external network, or provider write."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "done", "evidence", "blockers", "next"]
}
```

## Phase 1 — Schemas, Registries, and Contract Tests

### Files

- Modify: `contracts/governance/README.md`
- Create: `contracts/governance/modules.json`
- Create: `contracts/governance/runtime-config.json`
- Create: `contracts/governance/execution-lanes.json`
- Create: `contracts/governance/invariants.json`
- Create: `contracts/governance/shared-seams.json`
- Create: `contracts/governance/schemas/modules.schema.json`
- Create: `contracts/governance/schemas/runtime-config.schema.json`
- Create: `contracts/governance/schemas/execution-lanes.schema.json`
- Create: `contracts/governance/schemas/invariants.schema.json`
- Create: `contracts/governance/schemas/shared-seams.schema.json`
- Create: `contracts/governance/schemas/context-pack.schema.json`
- Create during Phase 3 correction: `contracts/governance/schemas/feature-work-contract.schema.json`
- Create: `scripts/tests/governance-contracts.tests.ps1`

### Required registry content

- `modules.json` covers exactly 11 active module directories and includes exact
  temporary forbidden-layer exceptions for current cross-module imports.
- `runtime-config.json` contains all active canonical keys and declared Oracle
  aliases discovered in active Go, web, scripts, and Docker entrypoints. It
  distinguishes `active`, `legacy_current`, `legacy_alias`,
  `temporary_test_contract`, and `reserved_not_ambient` lifecycle values.
- `execution-lanes.json` contains `unit`, `integration`, `live-oracle`,
  `live-provider-read`, `browser`, and `provider-write`. Unit declares
  `inherit_parent: false`, network/database disabled, fake target, and no
  application runtime keys. Provider-write declares every actor/idempotency/
  execute/link/policy/timestamp/before-after gate.
- `invariants.json` initially defines enforceable module coverage/dependency/
  layer/composition, application import, PostgreSQL driver, production panic,
  API/SDK atomicity, frontend fetch, and migration-prefix checks.
- `shared-seams.json` defines narrow exclusive seams: `api-sdk`,
  `migration-sequence`, `composition-root`, `dependency-graph`,
  `architecture-decisions`, and `provider-capability-contract`.

### RED/GREEN

1. Write contract test that expects all six schemas/registries, strict unknown
   property rejection, unique IDs, reference resolution, lane invariants, and
   exact exception shapes.
2. Run:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1
```

Expected RED: exit 1 because schemas/registries do not exist.

3. Create schemas and registries with no placeholder records.
4. Rerun same command.

Expected GREEN: exit 0 and `PASS governance contract tests`.

5. Commit only Phase 1 paths:

```text
feat(governance): add executable repository contracts
```

## Phase 2 — Policy Loader and Semantic Drift

### Files

- Create: `scripts/harness/Policy.psm1`
- Create: `scripts/tests/governance-drift.tests.ps1`
- Modify: `scripts/harness.ps1`
- Modify: `package.json`
- Modify: `AGENTS.md`

### Semantic checks

- Module directory coverage, roots, declared dependencies, forbidden target
  layers, and composition imports.
- Runtime key/alias ownership, exact reader inventory, temporary direct-reader
  exceptions, sensitivity, lane access, and missing/stale readers.
- Application imports, PostgreSQL adapter driver boundary, production panic,
  OpenAPI/SDK atomicity relative to base SHA, frontend direct fetch, and unique
  migration prefixes.
- Temporary exceptions match exact path/key/edge and cannot use globs.
- Existing known violations produce `baseline_exception` records, not pass
  claims; any new violation exits 1.

### RED/GREEN

1. Write temp-fixture tests for one positive repository and pinned negative
   cases: missing module, undeclared dependency, forbidden application import,
   undeclared env reader, alias collision, secret classification, production
   panic, duplicate migration prefix, OpenAPI/SDK split, and frontend fetch.
2. Run:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1
```

Expected RED: exit 1 because `Policy.psm1` functions do not exist.

3. Implement `Policy.psm1`, then add dispatcher commands and npm alias
   `harness:governance` that runs validate plus drift.
4. Rerun focused tests and current repository command:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1
npm run harness:governance -- -BaseSha (git rev-parse HEAD)
```

Expected GREEN: fixture suite exits 0; current command exits 0 while reporting
exact known `baseline_exception` IDs and no undeclared growth.

5. Update `AGENTS.md` machine-governance line from conditional reservation to
accepted authority only after both commands pass.
6. Commit only Phase 2 paths:

```text
feat(harness): enforce governance drift
```

## Phase 3 — Context Compiler and Currentness Validation

### Files

- Create: `scripts/harness/Context.psm1`
- Create: `scripts/tests/context-compiler.tests.ps1`
- Modify: `scripts/harness.ps1`
- Modify: `package.json`
- Create at runtime only: `scripts/.runs/<run-id>/context-pack.json`

### Context behavior

- Parse the active feature/spec/plan plus parent milestone/validation contract.
- Parse and schema-validate exactly one `## Machine Work Contract` JSON block
  from the plan. Resolve its criteria and commands against spec acceptance IDs,
  parent milestone IDs, and execution-lane IDs; missing/dangling proof fails.
- Always include canonical feature/spec/plan, parent milestone/validation,
  mission, root guidance, and declared extra sources. Validation recompiles all
  derived fields and requires exact equality; it never trusts pack-declared
  sources, objective, commands, labels, criteria, or risk.
- Compute source SHA-256, Git base SHA, paths/seams, deterministic risk/review
  policy, side effects, commands/targets, stop conditions, retry budget, and
  handoff fields.
- Risk rules: external/live/provider-write => L3; exclusive seam or cross-module
  scope => L2; one module => L1; docs/mechanical and no seam => L0.
- Estimate serialized input size excluding the estimate field and reject above
  2,000 using `[Text.Encoding]::UTF8.GetByteCount(...)`.
- `context-validate -RequireCurrentBase` recomputes Git SHA and every source
  hash before dispatch.

### RED/GREEN

1. Write two feature-agnostic positive fixtures and negative tests for missing/
   duplicate/invalid machine work contract, invalid SHA, stale HEAD,
   missing/mutated source, missing criterion proof, dangling command proof,
   pack-field tampering, ancestor scope escalation, rooted/traversal path,
   allowed/forbidden overlap, out-of-scope path, undeclared seam, side-effect
   conflict, fake-to-live inflation, UTF-8 multibyte estimate overflow, and
   2,001 estimate.
2. Run:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1
```

Expected RED: exit 1 because `Context.psm1` functions do not exist.

3. Add `feature-work-contract.schema.json`; implement generic compiler/
   validator, add dispatcher commands and npm aliases
   `harness:context:compile` and `harness:context:validate`.
4. Rerun focused tests, then compile F-07 from current HEAD:

```powershell
$sha = git rev-parse HEAD
$allowed = @(
  'contracts/governance/**'
  'scripts/harness/**'
  'scripts/tests/governance-contracts.tests.ps1'
  'scripts/tests/governance-drift.tests.ps1'
  'scripts/tests/context-compiler.tests.ps1'
  'scripts/harness.ps1'
  'package.json'
  'AGENTS.md'
  '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler/**'
)
& ./scripts/harness.ps1 -Command context-compile -FeaturePath '.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-07-governance-context-compiler' -BaseSha $sha -AllowedPath $allowed
npm run harness:context:validate -- -ContextPath '<reported context-pack path>' -RequireCurrentBase
```

Expected GREEN: both commands exit 0; pack reports risk `L2`, target `fake`,
all F07 criteria/proofs, source hashes, declared seams, and estimate <= 2,000.

5. Commit only Phase 3 paths:

```text
feat(harness): compile bounded feature context
```

## Phase 4 — Feature Validation and Handoff

### Files

- Create: `F-07-governance-context-compiler/validation.md`
- Modify: `F-07-governance-context-compiler/feature.md`

Run complete focused gate from clean current HEAD:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1
npm run harness:governance -- -BaseSha (git rev-parse HEAD)
```

Recompile context using final HEAD, validate it, then record exact commands,
exit codes, target `fake`, known exceptions, changed paths, evidence artifacts,
and risks. Status may become `quick_validation_passed`, never `accepted`.

Commit:

```text
docs(m08): record governance context evidence
```

## Allowed Paths

- All files named in Phases 1-4.
- No product, migration, OpenAPI, SDK, Docker, provider, or `.mnfs` path outside
  F-07 artifacts may change.

## Verification Mapping

- F07-AC01: `governance-contracts.tests.ps1` and `harness:governance`.
- F07-AC02: `governance-drift.tests.ps1` plus current repository drift result.
- F07-AC03: `context-compiler.tests.ps1` and current F-07 compile/validate.
- F07-AC04: context negative path/seam/side-effect/target fixtures.
- F07-AC05: context estimate boundary fixtures.

## Rollback and Risk Notes

- One writer owns `contracts/governance`, `scripts/harness.ps1`, and
  `package.json`; phases remain serial.
- Do not change production code to make drift pass. Add only exact temporary
  exceptions proven by current HEAD and assign removal owner.
- Do not weaken a schema or invariant after a RED fixture; fix the registry,
  parser, or test contract.
- Do not use reset/revert/stash/clean/checkout/restore. Stop on unrelated dirty
  paths.
- Raw context artifacts remain ignored under `scripts/.runs/`.

## Handoff

- Current status: `planned`.
- Next owner: Fresh Phase 1 Feature Implementer in build mode.
- Next action: Execute Phase 1 only, return commit/evidence, then dispatch Phase
  2 from the accepted Phase 1 SHA.
- Required files/evidence: feature, spec, plan, PowerShell 7.6 capability,
  current runtime/module inventories, RED/GREEN output, and phase commit.
- Blockers or open decisions: None.
