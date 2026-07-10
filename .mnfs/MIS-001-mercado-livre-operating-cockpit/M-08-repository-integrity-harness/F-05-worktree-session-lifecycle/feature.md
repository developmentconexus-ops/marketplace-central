# F-05-worktree-session-lifecycle

```yaml
id: F-05
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Mission
MIS-001 Mercado Livre Operating Cockpit.

## Milestone
M-08 Repository Integrity and Deterministic Harness.

## Brief
Make milestone sessions and worktree runtimes reproducible through project/run namespaces, bounded context packs, compact handoffs, runbooks, and a reusable local Codex skill.

## Inputs
- Accepted baseline SHA, F-02 command surface, F-04 evidence manifest, orchestration design.

## Expected Output
- Parameterized Compose/project/port/database identity with no fixed container assumptions.
- Runbook, local skill, context-pack validator, and fresh-thread/worktree proof.

## Constraints
- Prefer native Codex worktree tools.
- Worktree directory/state cannot enter Git.
- Skill contains workflow, not project secrets or duplicated architecture prose.

## Negative Scenarios
- Dirty/unaccepted starting state: dispatch rejected.
- Shared runtime namespace collision: startup rejected with conflicting run ID.
- Context pack omits objective, constraints, paths, gates, or evidence types: validation fails.

## Validation Expectations
- Fresh Terra thread starts from named SHA, reads only listed files, and reports first action without broad discovery.
- Worktree stack uses independent runtime identity and leaves primary checkout unchanged.

## Execution Artifact Rules
`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff
- Current status: Briefed.
- Next owner: Feature Implementer after F-02, final validation after F-04.
- Next action: Create `spec.md` and `plan.md` for worktree/session lifecycle.
- Required files/evidence: skill validation, runbook review, fresh-thread/worktree proof.
- Blockers or open decisions: Accepted baseline and frozen command surface.

