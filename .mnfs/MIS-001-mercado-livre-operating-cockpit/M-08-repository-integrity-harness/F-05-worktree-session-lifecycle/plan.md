# F-05 — Goal Orchestration Control Plane Plan

```yaml
id: F-05
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-05
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-05-worktree-session-lifecycle

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-05",
  "required_sources": [
    "contracts/governance/knowledge-routes.json",
    "scripts/harness/Context.psm1",
    "scripts/harness/State.psm1",
    "scripts/harness/Impact.psm1",
    ".agents/skills/mpc-goal-harness/SKILL.md"
  ],
  "knowledge_route_ids": ["harness-control-plane", "root-bootstrap"],
  "allowed_paths": [
    "AGENTS.md",
    ".agents/skills/mpc-goal-harness/**",
    "contracts/governance/**",
    "scripts/harness.ps1",
    "scripts/harness/**",
    "scripts/tests/**",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/execution-guide.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-05-worktree-session-lifecycle/**"
  ],
  "forbidden_paths": [
    "apps/**",
    "packages/**",
    "contracts/api/**",
    "docker/**"
  ],
  "side_effects": {
    "allowed": ["repository-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "orchestration-tests", "command_id": "orchestration-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "context-compiler-tests", "command_id": "context-compiler-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "impact-probe-one", "command_id": "impact-probe-one", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "governance-contracts", "command_id": "governance-contracts", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F05-AC01", "milestone_criterion_id": "M-08-C09", "command_ids": ["orchestration-tests", "context-compiler-tests", "governance-contracts"]},
    {"id": "F05-AC02", "milestone_criterion_id": "M-08-C10", "command_ids": ["orchestration-tests"]},
    {"id": "F05-AC03", "milestone_criterion_id": "M-08-C11", "command_ids": ["orchestration-tests"]},
    {"id": "F05-AC04", "milestone_criterion_id": "M-08-C13", "command_ids": ["orchestration-tests"]},
    {"id": "F05-AC05", "milestone_criterion_id": "M-08-C14", "command_ids": ["orchestration-tests"]},
    {"id": "F05-AC06", "milestone_criterion_id": "M-08-C15", "command_ids": ["orchestration-tests", "impact-probe-one"]},
    {"id": "F05-AC07", "milestone_criterion_id": "M-08-C06", "command_ids": ["orchestration-tests"]}
  ],
  "stop_conditions": [
    {"code": "stale-pack", "condition": "A route source, context pack, or accepted base is stale."},
    {"code": "lease-conflict", "condition": "A checkout or shared seam already has an active lease."},
    {"code": "unregistered-gate", "condition": "A risk route selects an unregistered impact command."},
    {"code": "unsupported-control", "condition": "The requested native task control is not operator-observed; use the labelled fresh-session fallback."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "checkpoint-id", "changed-paths", "evidence", "review", "next"]
}
```

## Steps

1. Add RED fixtures for route selection, compact pack limits, leases,
   checkpoints/resume, risk routing, and repository-skill capability bounds.
2. Add governance route data/schema and extend the context compiler to emit
   ordered selectors, source hashes, and the measured bootstrap set.
3. Add `State.psm1` and a stable harness façade for lease, checkpoint, recovery,
   resume, handoff, and risk-policy operations.
4. Add the progressive repository skill; shorten root bootstrap and the active
   execution guide atomically, retaining durable policy through pointers.
5. Run focused test suites, governance validation/drift, compile and validate
   an F-05 context pack, and run its registered current-checkout impact gate.
6. Review the fixed commit once, apply at most one consolidated correction, then
   rerun the full focused validation set and record `validation.md`.

## Files Expected To Change

- `AGENTS.md` — compact durable bootstrap and skill/pack routing.
- `.agents/skills/mpc-goal-harness/**` — progressive orchestration procedure.
- `contracts/governance/{knowledge-routes.json,schemas/knowledge-routes.schema.json}` — canonical routes.
- `contracts/governance/schemas/{context-pack,feature-work-contract}.schema.json` — selector and route contract.
- `scripts/harness/{Context,State,Impact}.psm1`, `scripts/harness.ps1` — compiled context, state, and registered router façade.
- `scripts/tests/{context-compiler,harness-orchestration,governance-contracts,harness-impact}.tests.ps1` — deterministic proof.
- `M-08/execution-guide.md` and `F-05/{spec,plan,validation}.md` — execution truth.

## Verification Commands

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/harness-orchestration.tests.ps1`
  - Satisfies: F05-AC01, F05-AC02, F05-AC03, F05-AC04, F05-AC05, F05-AC06, F05-AC07.
  - Expected result: all route/state/risk/skill fixtures pass; native worktree
    starts from the named SHA and leaves the primary checkout unchanged.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/context-compiler.tests.ps1`
  - Satisfies: F05-AC01.
  - Expected result: schema-valid compact F-05 pack validates and tampering fails closed.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`
  - Satisfies: F05-AC01.
  - Expected result: route governance contracts reject invalid shape.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command impact -ContextPath <validated-F05-pack>`
  - Satisfies: F05-AC06.
  - Expected result: only registered unit commands execute in the current checkout.

## QA Steps

- QA-3 milestone dogfood and operator-observed native task/subagent controls
  remain for F-09/independent QA. This feature proves the deterministic contract,
  not live task-control parity.

## Rollback/Risk Notes

- This is one atomic shared control-plane seam despite more than four paths;
  splitting would leave the bootstrap, schema, or state procedure incompatible.
- State files are ignored and recovery is non-destructive. A regression is
  repaired in a new targeted commit; never restore cold-run behavior.

## Execution decision

`split_decision: single` — one writer is required because the plan contract,
context compiler, route registry, state protocol, and bootstrap pointers must
remain compatible at every committed boundary.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement the atomic control-plane slice and record evidence.
- Required files/evidence: spec, route/state tests, current-checkout impact outcome.
- Blockers or open decisions: None.
