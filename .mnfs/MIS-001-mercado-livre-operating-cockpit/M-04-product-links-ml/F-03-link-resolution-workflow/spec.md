# Feature Spec

```yaml
id: F-03
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-03
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-link-resolution-workflow

## Problem

`product_links` can already import listing snapshots and generate exact-first candidates, but operators still cannot turn those candidates into audited resolved or rejected link truth. Without an operator workflow, downstream stock, pricing, and order features still lack a safe, explicit linkage boundary.

## Requirements

- Persist operator-managed product link truth separately from generated candidates.
- Support explicit operator actions to approve a candidate, reject a listing identity, and manually resolve a listing identity to an internal product.
- Record audit evidence for every state change, including before/after state, actor metadata, reason, and the source candidate when applicable.
- Keep generated candidate states visible; conflict and unresolved states must remain visible until an operator resolves or rejects them.
- Expose the workflow through API, OpenAPI, SDK, and a dedicated UI surface.
- UI must use SDK runtime only.
- The workflow must not apply stock writes or any other marketplace write beyond link-state persistence.

## Non-Goals

- Automatic resolution of conflict or title-match candidates without operator action.
- Any Mercado Livre stock, price, or listing write.
- Hiding candidate-generation uncertainty behind optimistic UI defaults.

## Acceptance Criteria

- Backend tests prove approve, reject, and manual resolve flows persist link truth and audit entries.
- API/SDK contract exposes candidate list, current link state, audit evidence, and resolution commands consistently.
- UI tests prove loading, error, empty, conflict, unresolved, resolved, and rejected states plus successful operator actions.
- The resulting link truth is suitable for future downstream modules to block writes when the link is unresolved or rejected.

## Execution Slices

1. Backend workflow and persistence.
2. OpenAPI and SDK parity.
3. Dedicated product-links UI surface and route.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: write the execution plan and implement the backend workflow slice first
- Required files/evidence: feature brief, spec, milestone contract, focused validation commands
- Blockers or open decisions: decide whether the UI lives in a new `feature-product-links` package or is folded into an existing feature package; default is a new package
