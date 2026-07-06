# M-03-mnos-sankhya-read-contract

```yaml
id: M-03
type: milestone
status: planned
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

MPC has explicit read contracts for internal product, stock, price, cost, tax, and sales inputs derived from MNOS/Sankhya evidence, with no Sankhya write path and no wholesale ERP mirror.

## Why This Milestone Exists

Stock Seguro and margin quality depend on correct internal facts. MNOS already mapped Sankhya semantics; MPC should import the relevant contract, not rediscover or improvise SQL.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Sankhya read contract import | Bring IC-002 semantics into MPC docs/types/ports for product, stock, price, cost, sales, tax. |
| F-02 | Read adapter seam | Add the application-facing read port and fake/test adapter; real Oracle edge can follow the contract. |
| F-03 | Data quality rules | Implement quality flag semantics for missing/ambiguous/stale product, stock, cost, tax, and source data. |

## Dependencies

- IC-002.
- MNOS source files remain available for evidence.
- M-02 capability framework may run in parallel, but Stock Seguro requires both M-03 and M-04.

## Risks

- Drift from MNOS semantics.
- Secret leakage through Oracle errors.

## Done Means

- Stock and cost rules are fixed in code/tests/docs.
- `CODLOCAL=10108` does not contribute to default sellable stock.
- Missing cost/tax returns quality flags rather than zero.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: Start F-01 after M-01 or in parallel with M-02 if reviewer approves.
- Required files/evidence: F-*/validation.md and M-03/validation-result.md.
- Blockers or open decisions: Real Oracle credential usage requires environment confirmation during implementation.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-03/validation-result.md.
- Revalidation evidence required: read-contract tests and secret-safety tests.
