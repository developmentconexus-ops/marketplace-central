# F-04-stock-seguro-dashboard

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Build the operator-facing Stock Seguro dashboard and action panel using SDK runtime data only.

## Inputs

- Stock risk API.
- Stock action API.
- Existing UI primitives and route patterns.

## Expected Output

- Dashboard supports filtering by risk/link/actionability.
- Row detail shows internal stock, ML stock, recommendation, source timestamps, policy, and block reason.
- Manual apply flow is clearly confirmatory.

## Constraints

- React does not calculate stock safety math.
- No visible instructional marketing copy; UI should be an operational cockpit.

## Interaction Model

- Open dashboard -> scan risk rows -> filter oversell/blockers -> inspect detail -> confirm manual action -> see applied/failed audit result -> row refetches.

## Validation Expectations

- UI tests cover loading, error, empty, healthy, oversell, conflict, unresolved, stale, and action result states.
- Browser validation confirms no text overlap on desktop/mobile.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-04/validation.md.
- Blockers or open decisions: final route path can use `/inventory/stock-seguro` unless router conventions force another path.
