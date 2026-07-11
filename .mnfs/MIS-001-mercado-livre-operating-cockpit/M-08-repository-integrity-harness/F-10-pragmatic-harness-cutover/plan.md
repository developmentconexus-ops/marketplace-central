# F-10 — Pragmatic Harness Cutover Plan

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-10",
  "required_sources": [
    "scripts/harness.ps1",
    "scripts/harness/Context.psm1",
    "scripts/harness/Evidence.psm1",
    "contracts/governance/execution-lanes.json"
  ],
  "allowed_paths": [
    "package.json",
    "AGENTS.md",
    "scripts/harness.ps1",
    "scripts/harness/**",
    "scripts/tests/**",
    "contracts/governance/**",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-10-pragmatic-harness-cutover/**"
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
    {"id": "impact-probe-one", "command_id": "impact-probe-one", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "impact-probe-two", "command_id": "impact-probe-two", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "governance-contracts", "command_id": "governance-contracts", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F10-AC01", "milestone_criterion_id": "M-08-C17", "command_ids": ["impact-probe-one", "governance-contracts"]},
    {"id": "F10-AC02", "milestone_criterion_id": "M-08-C15", "command_ids": ["impact-probe-one", "impact-probe-two"]},
    {"id": "F10-AC03", "milestone_criterion_id": "M-08-C15", "command_ids": ["impact-probe-one"]},
    {"id": "F10-AC04", "milestone_criterion_id": "M-08-C15", "command_ids": ["impact-probe-one"]}
  ],
  "stop_conditions": [
    {"code": "stale-pack", "condition": "Context base or source hash is not current."},
    {"code": "out-of-scope", "condition": "Changed path is outside the pack allowed paths."},
    {"code": "unregistered-command", "condition": "A selected ID has no registered typed argv."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "changed-paths", "evidence", "review", "next"]
}
```

## Ordered work

1. Make feature contracts encode registered command IDs rather than arbitrary
   command templates, preserving lane-derived target classification.
2. TDD the current-checkout impact gate with fixtures for ordered selection,
   fail-closed invalid inputs, and redacted outcome persistence.
3. Add the closed registry and impact orchestrator, then connect the stable
   façade and alias.
4. Remove cold dispatcher code, npm alias, lane, cold evidence vocabulary, and
   cold-only tests.
5. Update governance/context tests and run focused contracts plus the impact
   test, aliases, and boundary search.
6. Run one fixed-commit proportional review, make at most one consolidated
   correction batch, and record quick-validation evidence.

## Expected changed paths

- `package.json`
- `scripts/harness.ps1`
- `scripts/harness/Impact.psm1`
- `scripts/harness/Context.psm1`
- `scripts/harness/Evidence.psm1`
- `scripts/tests/harness-impact.tests.ps1`
- `scripts/tests/governance-contracts.tests.ps1`
- `contracts/governance/execution-lanes.json`
- `contracts/governance/schemas/{feature-work-contract,context-pack,harness-outcome}.schema.json`
- retired `scripts/tests/cold-gate-*.tests.ps1`

## Verification mapping

| Criterion | Commands |
| --- | --- |
| F10-AC01 | `pwsh scripts/tests/harness-impact.tests.ps1`; boundary `rg` search |
| F10-AC02 | `pwsh scripts/tests/harness-impact.tests.ps1` |
| F10-AC03 | `pwsh scripts/tests/harness-impact.tests.ps1` |
| F10-AC04 | `pwsh scripts/tests/harness-impact.tests.ps1` |

## Execution decision

`split_decision: single` — the ordered removal and replacement spans more than
four paths, but it is an atomic single-writer seam: a split would leave either
a live cold entrypoint or an uncallable impact gate between sessions.

## Risks and rollback

The gate is allowed to block only on contract/current-checkout mismatches. If a
regression appears, restore the affected current-checkout implementation via a
new targeted change; do not restore cold compatibility or historical behavior.
