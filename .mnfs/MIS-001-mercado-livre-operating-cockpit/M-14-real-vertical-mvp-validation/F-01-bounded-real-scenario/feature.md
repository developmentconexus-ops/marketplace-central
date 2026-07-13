# F-01-bounded-real-scenario

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-14
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-14 Real Vertical MVP Validation.

## Brief

Define a bounded, values-minimized local scenario selection that joins one
installation, product, listing, order, source facts, and simulation without
persisting credentials/PII or authorizing a provider write.

## Inputs

- Passed M-09/M-13 contracts and real local integration configuration.
- IC-003 entity IDs/source facts, IC-004 proportional lane, and M-14-C01/C02/C05/C06.

## Inputs/Outputs

- Input identities are supplied through a Git-ignored local file or process-local
  operator selection, never committed values.
- Output is a concise sanitized note with source, evidence type, observed time,
  read-only status, command/interaction, reviewed SHA, and artifact path.
- Reset deletes/reseeds only deterministic fixture-owned rows, never real imported records.

## Negative Scenarios

- No sample joins all required entities: persist IC-004 terminal status
  `externally_blocked` with an allowed `missing_edge`; never return it from an API.
- Selected record contains buyer PII: exclude/redact before evidence.
- Any lane needs provider mutation: fail preflight and stop.
- Ambiguous link: retain as a negative sample, never select it as the successful path.

## Expected Output

A repeatable local selector lets QA address the same scenario while durable evidence
contains no sensitive business/customer values.

## Constraints

- No credential values, order payloads, buyer fields, or full Oracle rows in MNFS.
- No provider/Oracle writes and no synthetic claim of real data.
- Owned paths: the smallest local scenario-selection/testability support proved
  necessary, fixture support outside production rows, the exact `/.mvp-local/`
  entry in `.gitignore` when a file selector is used, and this feature root. QA
  alone writes final `_fixed-sha-qa` artifacts.
- Forbidden paths: provider mutation, production data deletion, auth/RBAC.

## Criteria IDs

- M-14-C01 Real source provenance.
- M-14-C05 Evidence security and honesty.
- M-14-C06 No provider mutation.

## Validation Expectations

- Tests or preflight reject writable live lanes, secrets, and PII fields.
- If a selector file is used, it is Git-ignored and untracked.
- Missing join returns the compact IC-004 checkpoint naming one allowed missing edge.
- Reset test proves real-import rows remain untouched.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create spec/plan and implement only the bounded selector/preflight support proved necessary.
- Required files/evidence: M-14 contract, IC-003, IC-004, and this feature's `validation.md`.
- Blockers or open decisions: External sample availability is not an implementation defect.
