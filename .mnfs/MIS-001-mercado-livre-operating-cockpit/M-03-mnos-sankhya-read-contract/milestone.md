# M-03-mnos-sankhya-read-contract

```yaml
id: M-03
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

MPC must own explicit Oracle read contracts and adapters for internal product, stock, price, cost, tax, and sales inputs, with no ERP write path and no wholesale ERP mirror.

## Why This Milestone Exists

Stock Seguro and margin quality depend on correct internal facts. Legacy MNOS evidence is still useful, but MPC now needs to own the Oracle boundary directly inside `apps/server_core` instead of treating old contracts or `MS_DATABASE_URL` assumptions as execution truth.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Oracle read contract redesign | Re-own MPC docs/types/ports for product, stock, price, cost, sales, and tax against the Oracle source. |
| F-02 | Oracle adapter implementation | Add the real application-facing Oracle adapter seam and supporting query/mapping structure. |
| F-03 | Data quality rules | Preserve quality-flag semantics for missing/ambiguous/stale product, stock, cost, tax, and source data. |

## Dependencies

- IC-002 replacement or rewrite.
- Legacy MNOS source files remain available as reference evidence only.
- M-02 capability framework may run in parallel, but Stock Seguro requires both M-03 and M-04.

## Risks

- Drift between legacy evidence and MPC-owned Oracle semantics.
- Secret leakage through Oracle errors.

## Done Means

- Oracle read contracts and adapter structure are owned by MPC and verified against real source semantics.
- Sellable-stock policy is explicit and auditable in code/tests/docs.
- Missing cost/tax returns quality flags rather than zero.

## Handoff

- Current status: passed.
- Next owner: Mission Strategist.
- Next action: roll the validated Oracle-first runtime path into downstream milestones and preserve the live smoke harness for future regression checks.
- Required files/evidence: current code/test evidence plus retained live Oracle output for product, stock, price, cost, sales, and tax reads.
- Blockers or open decisions: none at milestone scope.

## Correction Handoff

- QA failure summary: previous pass was valid only for the superseded contract-seam shape.
- Correction scope: replace old M-03 execution truth with Oracle-first replanning.
- Attempts used/remaining: 0/2.
- Next artifact: M-03/validation-result.md.
- Revalidation evidence required: read-contract tests and secret-safety tests.
