# F-05-operations-ux-convergence

```yaml
id: F-05
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-13 Integrated Operator Workspaces.

## Brief

Converge Marketplaces and Integrations into `/operations`, retire duplicate primary
navigation, and finalize cross-workspace states, Portuguese copy, and responsive behavior.

## Inputs

- Existing marketplace catalog/settings and Integrations Hub components/SDK operations.
- IC-003 Operations order, route redirects, quality/source/time vocabulary.
- Accepted F-01–F-04 workspace shells.

## Inputs/Outputs

- Input: optional installation/provider/health/needs-action filters.
- Output: channel setup, installations, OAuth/credential state, capability health,
  probes, operation runs, actionable errors, source, and observed time.

## Interaction Model

- Channel setup and operational health are internal views of one workspace.
- Selecting an installation uses shared context and filters other workspaces.
- Retry refreshes only the failed run/probe and keeps drawer/detail context.

## State Model

Connection/auth/run statuses remain server-owned provider states. Operations renders
them unchanged and separately consumes server-owned IC-003 quality for attention.

## Negative Scenarios

- Reauthorization required: visible action, no false connected state.
- Probe failure: source/time/error/Retry visible; other installations remain usable.
- Secret-bearing credential response: redact and fail test; never render raw values.

## Expected Output

One Operations workspace answers whether ML and Sankhya data are trustworthy and what
requires attention, while all five workspaces share a consistent visual/state language.

## Constraints

- Do not perform M-10 provider runtime consolidation.
- No new provider types or credential storage.
- Owned paths: feature-marketplaces/integrations convergence, feature-local UI
  consistency, and this feature root. AppRouter/Layout/context/redirects are read-only
  because F-01 already owns and wires them.
- `packages/ui` is read-only; cross-workspace consistency uses the accepted F-01
  primitives or returns a proved gap to the Orchestrator.
- Forbidden paths: provider registries/adapters, production auth, domain formulas.

## Criteria IDs

- M-13-C03 Operations convergence.
- M-13-C04 Proportional security and simulation.
- M-13-C05 State and responsive consistency.
- M-13-C06 SDK-only thin client.

## Validation Expectations

- `/marketplaces` and `/integrations` redirect to `/operations` preserving filters.
- One degraded installation shows source, observed time, error, and Retry.
- Desktop and 390x844 drives show no horizontal page overflow and consistent Portuguese copy.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after earlier workspace seams settle.
- Next action: Create spec/plan and finish convergence/polish in one bounded commit.
- Required files/evidence: accepted F-01–F-04 `validation.md` files and this feature's `validation.md`.
- Blockers or open decisions: Runtime duplication that blocks UX requires new M-10 authority.
