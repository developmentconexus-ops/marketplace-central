# M-01-vtex-removal-architecture-reset

```yaml
id: M-01
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

VTEX is removed from active target architecture surfaces and cannot receive new feature work. Provider catalog foundations remain intact for future marketplaces.

## Why This Milestone Exists

ADR-005 makes VTEX a dead path. Removing it first prevents new Mercado Livre work from depending on obsolete routes, SDK types, tests, docs, or UI navigation.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | VTEX surface inventory | Identify active VTEX code, routes, SDK methods, UI pages, tests, docs, env keys, and migrations before removal. |
| F-02 | VTEX active surface removal | Remove VTEX routes, SDK functions/types, frontend pages/nav, adapters, and tests that no longer belong to target architecture. |
| F-03 | Architecture truth alignment | Update OpenAPI, architecture/wiki/brain-facing docs and verification evidence after removal. |

## Dependencies

- ADR-005 accepted.
- Current tests and OpenAPI must be readable before edits.

## Risks

- Removing VTEX can expose hidden test or route dependencies.
- Some docs may mention VTEX historically; historical references must be classified instead of blindly deleted.

## Done Means

- `rg` finds no active VTEX route, SDK method, frontend navigation, adapter registration, or target-architecture docs.
- Tests/builds pass or remaining failures are scoped and documented.
- Historical docs, if retained, are explicitly marked legacy/historical.

## Handoff

- Current status: passed.
- Next owner: Mission Strategist.
- Next action: Start the next Mercado Livre-first mission using M-01 validation evidence as baseline.
- Required files/evidence: F-*/validation.md and M-01/validation-result.md.
- Blockers or open decisions: None.

## Correction Handoff

- QA failure summary: None.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: Mission continuation artifacts for the next Mercado Livre capability area.
- Revalidation evidence required: none for M-01 unless later regression appears.
