# F-07 Governance Registry and Context Compiler — Validation

```yaml
id: F-07
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-07
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-07-governance-context-compiler

## Summary

Quick validation passed from clean base
`9db79d7d01b161ba824c056dddb8772e7c82f9bd`. Strict governance contracts,
semantic drift fixtures, context compiler fixtures, current-repository
governance, F-07 context compilation, and current-base/source-hash validation
all exited 0. All evidence targets are `fake`; no database, network, provider,
or other external side effect was used.

This result is a Feature Implementer handoff, not feature acceptance, a formal
QA verdict, or an M-08 pass claim.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-11
- Final feature state for handoff: `quick_validation_passed`
- Base SHA: `9db79d7d01b161ba824c056dddb8772e7c82f9bd`
- Target: `fake`
- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: Pass

## Evidence Honesty

- Every passing command below was run in this Phase 4 session and records its
  observable output or artifact path.
- Fixture and repository checks prove deterministic contract behavior only.
  They do not constitute live integration evidence.
- The context pack serializes command evidence types as `assumed` because it is
  a pre-dispatch worker input. This validation artifact separately records the
  Phase 4 executions as `ran`.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None observed in the focused F-07 scope.
- Reason: All five F-07 acceptance criteria have executed proof mappings; the
  current context pack is schema-valid, criterion-complete, hash-current,
  path/seam bounded, fake-targeted, and below the 2,000-token estimate ceiling.

## Commands Run

### Governance contract fixtures

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0 and `PASS governance contract tests`
- Actual: exit 0 and `PASS governance contract tests`
- Artifact: `contracts/governance/` and `scripts/tests/governance-contracts.tests.ps1`
- Blocking condition: None

### Governance drift fixtures

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-drift.tests.ps1`
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0 and `PASS governance drift tests`
- Actual: exit 0 and `PASS governance drift tests`
- Artifact: `scripts/tests/governance-drift.tests.ps1`
- Blocking condition: None

### Context compiler fixtures

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0 and `PASS context compiler tests`
- Actual: exit 0 and `PASS context compiler tests`
- Artifact: `scripts/tests/context-compiler.tests.ps1`
- Blocking condition: None

### Current repository governance

- Command: `npm run harness:governance -- -BaseSha (git rev-parse HEAD)`
- Expanded base: `9db79d7d01b161ba824c056dddb8772e7c82f9bd`
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0, schema validity, semantic drift validity, and no undeclared
  repository growth
- Actual: exit 0, `status=passed`, 35 exact `baseline_exception` records, and
  `artifact_path=contracts/governance`
- Artifact: `contracts/governance/`
- Blocking condition: None

### Current F-07 context compilation

- Command:
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
  ```
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0 using only Machine Work Contract allowed paths
- Actual: exit 0 and `status=passed`
- Artifact: `scripts/.runs/57b45bced59b475091b4c0f4f0a5345c/context-pack.json`
- Blocking condition: None

### Current-base context validation

- Command: `npm run harness:context:validate -- -ContextPath 'scripts/.runs/57b45bced59b475091b4c0f4f0a5345c/context-pack.json' -RequireCurrentBase`
- Status: Pass
- Evidence type: `ran`
- Owner: Feature Implementer
- Target: `fake`
- Expected: exit 0 only if Git base, source hashes, derived fields, scope,
  seams, criteria, targets, side effects, and estimate remain current
- Actual: exit 0 and `status=passed`
- Artifact: `scripts/.runs/57b45bced59b475091b4c0f4f0a5345c/context-pack.json`
- Blocking condition: None

## Known Baseline Exceptions

Current governance emitted exactly 35 declared baseline exceptions. The
authoritative records and removal owners live in `contracts/governance/`; this
feature does not treat them as repaired behavior.

- Direct readers (29): `direct-reader-amazon-application-id`,
  `direct-reader-amazon-auth-version`, `direct-reader-amazon-client-id`,
  `direct-reader-amazon-client-secret`,
  `direct-reader-integrations-test-database`,
  `direct-reader-magalu-client-id`, `direct-reader-magalu-client-secret`,
  `direct-reader-manual-adjustment-test-database`,
  `direct-reader-marketplaces-test-database`, `direct-reader-me-client-id`,
  `direct-reader-me-client-secret`, `direct-reader-me-redirect`,
  `direct-reader-migrate-runner-database`,
  `direct-reader-migrate-runner-directory`, `direct-reader-migrations-dir`,
  `direct-reader-ml-client-id`, `direct-reader-ml-client-secret`,
  `direct-reader-oauth-redirect`, `direct-reader-oracle-live-gate`,
  `direct-reader-orders-test-database`,
  `direct-reader-phase1-smoke-database`,
  `direct-reader-product-links-installation`,
  `direct-reader-product-links-live-gate`,
  `direct-reader-product-links-postgres`,
  `direct-reader-profit-snapshot-test-database`,
  `direct-reader-shopee-base-url`, `direct-reader-shopee-partner-id`,
  `direct-reader-shopee-partner-key`, and `direct-reader-web-origin`.
- Repository structure (6): `migration-prefix-0021-duplicate`,
  `module-edge-connectors-magalu-integrations-adapter`,
  `module-edge-connectors-mercado-livre-integrations-adapter`,
  `module-edge-connectors-shopee-integrations-adapter`,
  `module-edge-integrations-feesync-marketplaces-registry`, and
  `production-panic-product-link-transition`.

## Context Pack Evidence

- Artifact: `scripts/.runs/57b45bced59b475091b4c0f4f0a5345c/context-pack.json`
- Evidence type: `ran`
- Status: Pass
- Base SHA: `9db79d7d01b161ba824c056dddb8772e7c82f9bd`
- Feature: `F-07-governance-context-compiler`
- Risk: `L2`
- Review policy: `independent-review`
- Advisory model: `strong-reasoning-model`
- Target labels: five commands, all `fake`
- Sources: nine ordered paths with SHA-256 hashes
- Shared seams: `dependency-graph`
- Side effects: none allowed
- Estimate: 1,272 input tokens, below the 2,000-token ceiling
- Criteria: `F07-AC01` through `F07-AC05`, all mapped to `M-08-C09` and
  declared proof command IDs
- Stop conditions: `scope-conflict`, `contract-drift`, and
  `external-side-effect`
- Retry budget: 1 correction attempt
- Handoff: Milestone Orchestrator

## Changed Paths Across Phases 1-3

- Feature contracts: `F-07-governance-context-compiler/spec.md`,
  `F-07-governance-context-compiler/plan.md`
- Governance guidance: `AGENTS.md`, `contracts/governance/README.md`
- Registries: `contracts/governance/modules.json`,
  `contracts/governance/runtime-config.json`,
  `contracts/governance/execution-lanes.json`,
  `contracts/governance/invariants.json`,
  `contracts/governance/shared-seams.json`
- Schemas: `contracts/governance/schemas/context-pack.schema.json`,
  `execution-lanes.schema.json`, `feature-work-contract.schema.json`,
  `invariants.schema.json`, `modules.schema.json`,
  `runtime-config.schema.json`, and `shared-seams.schema.json`
- Harness: `package.json`, `scripts/harness.ps1`,
  `scripts/harness/Context.psm1`, and `scripts/harness/Policy.psm1`
- Fixtures: `scripts/tests/context-compiler.tests.ps1`,
  `scripts/tests/governance-contracts.tests.ps1`, and
  `scripts/tests/governance-drift.tests.ps1`

No product code, migrations, OpenAPI, SDK, Docker, provider behavior, or MNFS
artifact outside F-07 changed in the accepted F-07 sequence.

## RED/GREEN and Correction History

- Phase 1 RED was the absent governance contract surface; GREEN was recorded in
  `5fdb88e` (`feat(governance): add executable repository contracts`). Contract
  review corrections were consolidated in `a4a063c`
  (`fix(governance): tighten executable contracts`).
- Phase 2 RED was the absent policy loader/drift enforcement; GREEN was recorded
  in `c2a288b` (`feat(harness): enforce governance drift`). Scanner-bypass
  corrections were consolidated in `e5afa44`
  (`fix(harness): close governance scanner bypasses`).
- Phase 3 RED was the absent context compiler/currentness validation; GREEN was
  recorded in `6d149c7` (`feat(harness): compile bounded feature context`). The
  work-contract structure was frozen in `2e1a4dd`, and contract-binding
  corrections were consolidated in `6de0af1`
  (`fix(harness): bind context to work contract`).
- The first Phase 4 compile attempt at base `6de0af1` exited 1 before the
  compiler ran because the plan example passed multiple scalar arguments to
  the PowerShell array parameter. This was an evidence-command shape gap, not
  an implementation defect. Commit `9db79d7`
  (`docs(m08): fix context path array invocation`) corrected the plan to bind an
  explicit `$allowed` array. Direct reproduction and this full refreshed gate
  then passed from the corrected base.
- Phase 4 made no same-session implementation fixup. Independent L2 acceptance
  review remains required from the Milestone Orchestrator; no review verdict is
  claimed here.

## Manual QA

- QA level: QA-3 deterministic feature gate
- Flow or step: No browser or live manual flow applies to governance/context
  compilation.
- Status: Not run
- Evidence type: `assumed`
- Owner: Feature Implementer
- Expected: Not applicable to this non-UI, fake-targeted feature
- Actual: Automated deterministic proof completed as recorded above
- Blocking condition: None; formal QA and feature acceptance remain downstream

## Risks

- The 35 exact baseline exceptions remain debt with assigned authoritative
  records; they prove no undeclared growth, not full conformance.
- The pack's `L2` classification requires the planned independent acceptance
  review because F-07 owns the `dependency-graph` shared seam.
- Pack artifacts are intentionally ignored runtime evidence. Recompilation is
  required if HEAD or any of the nine source hashes changes.
- Unit/fake evidence does not validate live databases, Oracle, providers, or
  provider writes; none are part of F-07 scope.
- F-08 must still implement the declared non-inheriting child environment;
  F-07 only defines and validates that contract.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review the fixed F-07 commit sequence, spec, plan, changed paths,
  and this validation evidence; then accept, reject, or block F-07.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, this
  `validation.md`, commits `5fdb88e` through the Phase 4 evidence commit,
  `contracts/governance/`, focused fixture suites, and the ignored pack path
  recorded above.
- Blockers or open decisions: No quick-validation blocker. Feature acceptance,
  independent L2 review, and M-08 QA remain outstanding.

## Milestone Orchestrator Acceptance

- Decision: `accepted`.
- Accepted implementation/evidence HEAD:
  `8a57678d1670a32e1aec0e10f210d51ea4f5155a`.
- Fresh acceptance commands:
  - governance contract fixtures: exit `0`, `PASS`;
  - governance drift fixtures: exit `0`, `PASS`;
  - context compiler/adversarial fixtures: exit `0`, `PASS`;
  - current `harness:governance`: exit `0`, `status=passed`, exactly 35
    declared baseline exceptions and no unknown drift code;
  - current F-07 context compile: exit `0`;
  - current-base context validation: exit `0`;
  - `git diff --check` and visible status: clean.
- Fresh acceptance context:
  `scripts/.runs/d13ec27f22bf4d95a3044bc256622c2a/context-pack.json`.
- Review history: each Phase 1-3 contract gate used one fixed-commit review;
  findings were reproduced, corrected in one batch, and focused revalidation
  passed. Phase 3 architecture was explicitly replanned to use a strict machine
  work contract instead of preserving an F-07-specific Markdown heuristic.
- Evidence class: `fake`; no database, Oracle, provider, browser, or external
  network target was used or claimed.
- Next owner: F-08 Feature Implementer.
