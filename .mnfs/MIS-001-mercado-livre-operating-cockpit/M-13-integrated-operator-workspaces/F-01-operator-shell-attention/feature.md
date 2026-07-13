# F-01-operator-shell-attention

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-13 Integrated Operator Workspaces.

## Brief

Replace module-shaped primary navigation with the five IC-003 workspaces, add one
shared installation context, and make `/` an actionable attention queue.

## Inputs

- AppRouter, Layout, ClientContext, existing installation/list/read SDK methods.
- IC-003 routes, AttentionItem ordering, quality/severity vocabulary, and seed.

## Inputs/Outputs

- Input: nullable installation filter plus existing domain lists.
- Output: ordered attention items with stable key, entity ref, source/quality/time,
  and exact target URL; no independent attention persistence.

## Interaction Model

- Selecting an installation updates shared URL query state and refetches applicable lists.
- Selecting attention navigates to its exact object and retains `attention=<kind>`.
- Reload reconstructs context from URL; browser back returns to the same filtered queue.

## State Model

`loading → ready`; ready contains empty results or rows whose server-owned source
quality is `current|stale|unknown|conflict`. Fetch failure renders `error` with Retry
and preserves installation. Stale retains prior value/time; unknown is null + reason.

## Negative Scenarios

- Malformed target identity: visible `invalid_identity`; no navigation.
- Source failure with a prior trustworthy value keeps that value/time as `stale`;
  without one it returns null + reason as `unknown`.
- Duplicate conditions for one entity/kind: one AttentionItem stable key.

## Expected Output

The operator starts from work requiring attention and reaches an exact entity without
knowing technical modules.

## Constraints

- Attention is derived from existing SDK results; no new business calculation/table.
- This feature exclusively owns AppRouter/Layout/shared installation context and
  additions to `packages/ui` during execution; the Orchestrator releases the seams
  only after accepting its commit.
- It wires every IC-003 reserved route and legacy redirect to stable workspace
  outlets; later M-13 features supply outlet components without router edits.
- Owned paths: apps/web router/layout/context, shared UI primitives needed by shell,
  feature-dashboard/attention, and this feature root.
- Forbidden paths: backend formulas, provider writes, Product 360/detail implementations.

## Criteria IDs

- M-13-C01 Navigation and attention closure.
- M-13-C05 State and responsive consistency.
- M-13-C06 SDK-only thin client.

## Validation Expectations

- Router tests assert five labels, legacy redirect targets, query preservation, and 404.
- Attention tests assert severity/time/key ordering and exact stock target URL.
- Browser shows a retryable source error and a reloadable stock attention deep link.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer.
- Next action: Create spec/plan and implement the shared shell before other M-13 features.
- Required files/evidence: IC-003, research inventory, and this feature's `validation.md`.
- Blockers or open decisions: Stop if attention requires inventing a domain classification.
