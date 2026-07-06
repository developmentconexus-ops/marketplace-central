# Feature Specification

```yaml
id: F-03
type: feature-spec
status: executed
owner: Feature Implementer
parent: F-03-architecture-truth-alignment
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-architecture-truth-alignment

## Summary

Update active architecture/wiki/brain truth so Marketplace Central no longer describes VTEX removal as pending work and instead records the completed reset toward a Mercado Livre-first control plane.

## Problem

F-02 removed active VTEX runtime, contract, SDK, and frontend surfaces, but active planning/truth documents still describe VTEX inventory/deletion as future work. That leaves the repo with contradictory execution truth versus planning truth.

## Goals

- Mark M-01 progress in `.brain/roadmap.json` and `.brain/system-pulse.md` based on completed F-01/F-02 evidence.
- Keep wiki environment guidance accurate by classifying VTEX keys as historical/legacy only.
- Remove architecture-facing phrases that imply VTEX connector work is still a forward path.
- Preserve historical evidence docs instead of rewriting them as if the work never happened.

## Non-Goals

- Do not edit forward-only database migrations.
- Do not rewrite historical research/evidence docs under `docs/`.
- Do not introduce new Mercado Livre feature scope beyond documenting the completed reset and next phase.

## Acceptance Criteria

1. Active truth docs no longer present VTEX inventory/removal as upcoming work.
2. `.brain/roadmap.json` reflects phase-4 task `4.2` as completed and points next work toward Mercado Livre modules.
3. `.brain/system-pulse.md` reflects F-01/F-02 completion, updates next steps, and downgrades VTEX residue to migration/history only.
4. Wiki and architecture-facing docs that remain in active use clearly mark VTEX env/config references as legacy-only and not part of current operations.
5. `rg -n "VTEX|vtex" ARCHITECTURE.md wiki .brain` returns only explicitly legacy/historical references or ADR context, not active next-step instructions.

## Files In Scope

- `.brain/system-pulse.md`
- `.brain/roadmap.json`
- `.brain/session-log.md`
- `.brain/decisions/001-metalshopping-direct-read.md`
- `wiki/operations/environment-and-db.md`
- `.mnfs/.../F-03-architecture-truth-alignment/{spec.md,plan.md,validation.md}`

## Risks

- Over-editing ADR history could erase why older decisions were made.
- Under-editing active planning docs could leave the repo claiming VTEX work is still next.

## Done Definition

- Active truth docs align with F-02 removal evidence.
- Remaining VTEX mentions in active truth are clearly legacy/historical.
- Validation evidence records exact `rg` checks and updated planning/doc files.
