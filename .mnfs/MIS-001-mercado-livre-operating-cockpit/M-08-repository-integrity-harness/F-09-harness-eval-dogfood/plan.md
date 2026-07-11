# F-09 Harness Eval and Dogfood — Plan

```yaml
id: F-09
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-09-harness-eval-dogfood

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-09",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/execution-guide.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-contract.md",
    "contracts/governance/harness-evals.json",
    "scripts/harness/Evals.psm1"
  ],
  "knowledge_route_ids": ["harness-control-plane", "root-bootstrap"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-09-harness-eval-dogfood/**",
    "contracts/governance/harness-evals.json",
    "scripts/harness/Evals.psm1",
    "scripts/harness/Impact.psm1",
    "scripts/tests/harness-eval.tests.ps1",
    "scripts/tests/fixtures/harness-eval/**"
  ],
  "forbidden_paths": [
    "apps/**",
    "packages/**",
    "contracts/api/**",
    "docker/**",
    "scripts/harness.ps1",
    ".agents/**"
  ],
  "side_effects": {
    "allowed": ["repository-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "harness-evals", "command_id": "harness-evals", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "orchestration-tests", "command_id": "orchestration-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "context-compiler-tests", "command_id": "context-compiler-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "governance-contracts", "command_id": "governance-contracts", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "harness-aliases", "command_id": "harness-aliases", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F09-AC01", "milestone_criterion_id": "M-08-C16", "command_ids": ["harness-evals", "context-compiler-tests"]},
    {"id": "F09-AC02", "milestone_criterion_id": "M-08-C04", "command_ids": ["harness-evals"]},
    {"id": "F09-AC03", "milestone_criterion_id": "M-08-C14", "command_ids": ["harness-evals", "orchestration-tests"]},
    {"id": "F09-AC04", "milestone_criterion_id": "M-08-C15", "command_ids": ["harness-evals", "harness-aliases", "governance-contracts"]},
    {"id": "F09-AC05", "milestone_criterion_id": "M-08-C17", "command_ids": ["harness-evals", "harness-aliases"]}
  ],
  "stop_conditions": [
    {"code": "stale-pack", "condition": "The accepted base, source hash, or required selector is stale."},
    {"code": "unregistered-gate", "condition": "A selected impact command has no structured registered argv."},
    {"code": "dogfood-scope", "condition": "The fresh bounded worker needs a path outside its three pinned dogfood paths."},
    {"code": "resume-gap", "condition": "A fresh continuation cannot reconstruct state from MNFS, Git, pack, and checkpoint artifacts."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "checkpoint-id", "commit", "changed-paths", "evidence", "review", "blockers", "next"]
}
```

## Steps

1. Add RED assertions for the closed corpus/result shape and the registered
   eval command ID.
2. Implement the pure deterministic grader and declarative corpus; it writes
   only a relative ignored result artifact and never executes manifest text.
3. Register the eval test through the existing impact command registry.
4. Dispatch one fresh bounded worker after the runner exists. It changes only
   the pinned out-of-scope fixture, manifest row, and test assertion; the
   grader is run by this owner after return.
5. Compile/validate the F-09 pack, run focused tests and its current-checkout
   impact gate, then record checkpoint and artifact-only resume evidence.
6. Review one fixed commit once; allow at most one consolidated correction and
   rerun the complete focused validation set if used.

## Files Expected To Change

- `contracts/governance/harness-evals.json` — closed versioned corpus.
- `scripts/harness/Evals.psm1` — deterministic evaluator and result writer.
- `scripts/harness/Impact.psm1` — registered structured eval-test argv.
- `scripts/tests/harness-eval.tests.ps1` — deterministic RED/GREEN proof.
- `scripts/tests/fixtures/harness-eval/out-of-scope-write/**` — fresh-worker
  dogfood input only.
- `F-09/{spec,plan,validation}.md` — feature execution truth.

## Verification Commands

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-eval.tests.ps1`
  - Satisfies: F09-AC01, F09-AC02, F09-AC03, F09-AC04, F09-AC05.
  - Expected: every pinned corpus case emits its deterministic rejection
    verdict/reason and a complete fake/contract result manifest.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-orchestration.tests.ps1`
  - Satisfies: F09-AC03.
  - Expected: checkpoint/resume and depth-one capability contract still pass.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
  - Satisfies: F09-AC01.
  - Expected: route-selected current pack validation remains fail-closed.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
  - Satisfies: F09-AC04.
  - Expected: the corpus governance document has valid deterministic shape.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-aliases.tests.ps1`
  - Satisfies: F09-AC04, F09-AC05.
  - Expected: only current-checkout harness commands remain active.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command impact -ContextPath <validated-F09-pack>`
  - Satisfies: F09-AC04.
  - Expected: only the contract's registered unit command IDs execute.

## QA Steps

- QA level: QA-3. The feature records fake/contract dogfood evidence only;
  independent milestone QA later decides any real-target criteria.

## Rollback/Risk Notes

- Keep corpus data declarative and the evaluator closed; do not add shell text
  or a generic execution API.
- If the fresh worker exceeds its pinned paths or lacks task controls, stop the
  dogfood as blocked or use the labelled fresh-session fallback; do not broaden
  scope or claim native-control parity.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Add RED evaluator tests and implementation.
- Required files/evidence: spec, pack, corpus, test results, checkpoint, and
  impact outcome.
- Blockers or open decisions: None.
