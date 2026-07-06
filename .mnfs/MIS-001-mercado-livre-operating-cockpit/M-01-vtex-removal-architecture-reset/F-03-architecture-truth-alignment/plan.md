# Feature Plan

```yaml
id: F-03
type: feature-plan
status: executed
owner: Feature Implementer
parent: F-03-architecture-truth-alignment
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: Active truth updates are tightly coupled; roadmap, system pulse, wiki, and ADR wording must move together to avoid contradictory architecture signals.
```

## Feature ID

F-03-architecture-truth-alignment

## Steps

1. Update `.brain/system-pulse.md` to record F-01/F-02 completion, revise current phase next steps, and reclassify remaining VTEX residue as migration/history only.
2. Update `.brain/roadmap.json` to mark phase-4 task `4.2` done and point subsequent work toward Mercado Livre missions/modules.
3. Update `.brain/session-log.md` so the latest session handoff no longer points to VTEX connector depth as the next track.
4. Update active wiki/environment guidance and architecture-facing ADR text where VTEX was still described as an operational target.
5. Run targeted `rg` checks over `ARCHITECTURE.md`, `wiki`, and `.brain`; record remaining legacy/historical matches in `validation.md`.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-03-architecture-truth-alignment/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-03-architecture-truth-alignment/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/F-03-architecture-truth-alignment/validation.md`
- `.brain/system-pulse.md`
- `.brain/roadmap.json`
- `.brain/session-log.md`
- `.brain/decisions/001-metalshopping-direct-read.md`
- `wiki/operations/environment-and-db.md`

## Verification Commands

- Command: `rg -n "VTEX|vtex|connectors/vtex|VTEX_APP_|VTEX_ACCOUNT" ARCHITECTURE.md wiki .brain`
  - Satisfies criterion ID: 1, 4, 5
  - Expected result: only explicit legacy/historical mentions remain.

- Command: `git diff -- .brain/system-pulse.md .brain/roadmap.json .brain/session-log.md .brain/decisions/001-metalshopping-direct-read.md wiki/operations/environment-and-db.md`
  - Satisfies criterion ID: 1, 2, 3, 4
  - Expected result: diffs show active-truth alignment only, with no runtime/code churn.

## QA Steps

- Step: Review roadmap and system pulse together.
  - Expected result: both say VTEX inventory/removal is complete and the next work targets Mercado Livre modules.

- Step: Review wiki environment page.
  - Expected result: VTEX keys are clearly legacy-only, with no suggestion of operational usage.

## Rollback/Risk Notes

- Historical `docs/**` and migration files remain untouched in F-03.
- ADR wording should preserve original context while clarifying that VTEX references were part of the old state, not the target architecture.
