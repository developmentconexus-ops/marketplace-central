# F-10-pragmatic-harness-cutover

```yaml
id: F-10
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-08 Repository Integrity Harness.

## Brief

Atomically retire the superseded cold-clone/provisioning design and replace it
with a current-checkout impact gate whose commands come from the active context
pack and whose redacted outcome records what actually ran.

## Inputs

- F-04 blocked evidence and operator supersession decision.
- Accepted F-07 context pack, F-08 child execution, F-03 PostgreSQL, governance,
  target labels, and generic evidence/redaction primitives.
- Pragmatic V1 design and M-08 validation contract.

## Expected Output

- No active `cold` command, npm alias, `cold-provision` lane, snapshot/provisioning
  implementation, cold-only schema vocabulary, or cold-only tests.
- Reusable evidence/redaction remains and records base/commit SHA, changed paths,
  risk, ordered task-declared commands, target labels, exit classifications,
  relative artifacts, remaining risks, and aggregate outcome.
- One impact-gate command validates the current pack/base/changed paths, reports
  declared seams for F-05 lease enforcement, and runs only versioned registered
  deterministic commands selected by the pack in the normal checkout.
- Existing unit, ephemeral PostgreSQL, live-read, browser, and provider-write
  commands remain separate and retain their safety gates.

## Inputs/Outputs

- Input: current context-pack path and its base SHA.
- Output: schema-valid redacted outcome with `passed|failed|blocked`, ordered
  commands, target labels, exits/reasons, changed paths, and relative artifacts.
- F-04 artifacts remain historical and are never converted to a passed result.

## Constraints

- Delete wrong active cold surfaces rather than leave aliases or compatibility wrappers.
- Do not rewrite accepted Environment, Execution, Policy, Context, or Postgres behavior without a demonstrated regression.
- No dependency installation, Docker pull, clone, cache reset, Oracle/provider/browser execution, or provider write in cutover validation.
- No product code, OpenAPI/SDK, migration, frontend, or provider adapter changes.
- Keep `scripts/harness.ps1` as the stable façade and reduce cold-specific complexity.
- Command IDs resolve to structured argv through versioned harness commands;
  never evaluate arbitrary plan text as a shell command.

## Negative Scenarios

- `harness:cold`, `-Command cold`, `cold-provision`, snapshot clone, or cold acceptance vocabulary remains active: cutover fails.
- Impact gate runs a command absent from the current pack: fail closed with a stable reason.
- Stale base/source or out-of-scope changed path: aggregate blocked before tests.
- Evidence contains an absolute path, secret-like value, or unknown target/class: persistence fails.
- One command fails: aggregate fails while preserving earlier safe records.

## Validation Expectations

- Boundary search returns zero active cold command/lane/alias/snapshot/provision references outside clearly marked historical artifacts.
- A fixture pack containing two commands runs exactly those two in order; an
  undeclared third command is not executed.
- Stale pack, changed path outside scope, and unknown target fixtures each exit nonzero with pinned reasons.
- Existing governance, context, environment, execution, evidence, alias, and
  F-03 fake-contract suites remain green.

## Criterion Mapping

| Criterion | Ownership | Minimum proof |
| --- | --- | --- |
| M-08-C17 | Primary | Zero active cold runtime/alias/lane/test references plus preserved F-04 history in this feature's `validation.md`. |
| M-08-C15 | Supporting | Fixture proves exactly two registered command IDs execute in order and arbitrary/undeclared text does not execute. |

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed and next executable feature.
- Next owner: Fresh F-10 Feature Implementer.
- Next action: Create spec/plan, identify reusable evidence primitives, then TDD the cutover.
- Required files/evidence: RED/GREEN, changed paths, boundary search, outcome fixtures, fixed-commit review.
- Blockers or open decisions: None; do not reopen cold behavior.
