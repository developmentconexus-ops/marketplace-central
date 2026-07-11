# F-07-governance-context-compiler

```yaml
id: F-07
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

Create schema-validated governance contracts and compile bounded, source-hashed
worker context from the active MNFS feature and repository truth.

## Inputs

- Accepted F-06 truth order, module/config/runtime code, OpenAPI/SDK boundary,
  M08 criteria, shared seams, and current harness command taxonomy.

## Expected Output

- Canonical JSON registries and JSON Schemas for modules, runtime config,
  execution lanes, invariants, shared seams, context packs, evidence, state,
  and eval cases.
- Semantic drift checker and context compiler available through the stable
  harness entrypoint.
- Ignored context pack containing base SHA, hashes, risk, criteria, paths,
  seams, side effects, commands, evidence types, and stop conditions.

## Constraints

- One fact has one writable owner; narrative sources link rather than duplicate
  machine-owned facts.
- No new package/runtime dependency.
- Compiler reads only declared sources and targeted paths; no broad repository
  dump enters the pack.

## Negative Scenarios

- Schema or semantic drift: exit non-zero before dispatch.
- Missing criterion-to-proof mapping: reject pack.
- Base SHA or source hash changes: reject stale pack.
- Requested path intersects an undeclared shared seam: reject readiness.

## Validation Expectations

- All registries and a positive context fixture validate with exit 0.
- Invalid schema, stale SHA/hash, missing proof mapping, and seam conflict each
  exit non-zero with stable reason codes and no secret values.
- Serialized context pack remains within the declared 2,000-token estimate.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after F-06.
- Next action: Create `spec.md` and `plan.md` with exact registry fields and
  compiler interfaces.
- Required files/evidence: schema fixtures, semantic drift cases, positive and
  stale context packs.
- Blockers or open decisions: F-06 truth-order acceptance.
