# M-05-stock-seguro-ml

```yaml
id: M-05
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

Operators can see Mercado Livre stock risk for existing linked listings and apply manual audited stock corrections with explicit policy and source evidence.

## Why This Milestone Exists

This is the first direct business value: prevent overselling and avoid cancellations caused by manual stock drift.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Stock policy model | Define sellable stock policy, buffer, ineligibility rules, and source freshness. |
| F-02 | Stock risk engine | Compare internal sellable stock with Mercado Livre available quantity and classify risk. |
| F-03 | Manual stock action audit | Propose, approve, apply, and audit manual Mercado Livre stock actions. |
| F-04 | Stock Seguro dashboard | Build operator UI for risk scanning, blocked cases, policy evidence, and manual action. |

## Dependencies

- M-02 capability framework.
- M-03 MNOS/Sankhya read contract.
- M-04 Product Links ML.

## Risks

- Provider data can be stale.
- Product exclusion rules may need operator refinement after real examples.
- Manual action UI can accidentally encourage unsafe writes if blocked states are unclear.

## Done Means

- Stock risk list never depends on synchronous provider calls.
- Unresolved/conflict links block actions.
- Manual action stores before/after/policy/source/provider response.
- Default policy excludes reserved and showroom stock.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: Start after M-02, M-03, and M-04 pass.
- Required files/evidence: F-*/validation.md and M-05/validation-result.md.
- Blockers or open decisions: Product exclusion groups/weight/size examples should be reviewed during F-01.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-05/validation-result.md.
- Revalidation evidence required: stock policy, risk engine, API, SDK, UI, and audit tests.
