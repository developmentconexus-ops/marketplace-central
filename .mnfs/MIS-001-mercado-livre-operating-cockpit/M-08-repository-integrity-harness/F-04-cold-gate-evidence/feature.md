# F-04-cold-gate-evidence

```yaml
id: F-04
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
Aggregate deterministic Go, JS, build, static, boundary, and integration checks into one cold gate with a redacted evidence manifest.

## Inputs
- F-02 command taxonomy, F-03 ephemeral integration lane, active Go/npm workspaces, architecture boundary rules.

## Expected Output
- One cold-gate command shared by local and future CI.
- Run-ID manifest with SHA, branch, dirty state, tools, commands, targets, exit codes, and artifact paths.

## Constraints
- Live Oracle/provider/browser evidence is separate and never required for deterministic CI.
- Gate cannot skip an active workspace silently.
- Raw logs stay ignored; MNFS stores summaries and paths.

## Negative Scenarios
- Dirty or wrong SHA when clean required: fail preflight.
- Unknown target type or secret-like output: fail evidence validation.
- One lane fails: aggregate result fails while preserving other lane evidence.

## Validation Expectations
- Two runs from same SHA produce same command inventory and result classification.
- Manifest differentiates fake, ephemeral PostgreSQL, live, and browser evidence.

## Execution Artifact Rules
`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff
- Current status: Briefed.
- Next owner: Feature Implementer after F-02/F-03.
- Next action: Create `spec.md` and `plan.md` for gate and manifest schema.
- Required files/evidence: two cold-run manifests and boundary checks.
- Blockers or open decisions: F-02/F-03 accepted interfaces.

