# M-07-commercial-intelligence

```yaml
id: M-07
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

MPC turns linked listings, safe stock, sales history, and profit snapshots into commercial recommendations: margin guardrails, repricing candidates, stock aging, and kit/bundle opportunities.

## Why This Milestone Exists

After operational safety and margin visibility, the platform can help sell better instead of only avoiding errors.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Margin guardrail policy | Define minimum margin policy and blocked recommendation states. |
| F-02 | Price recommendation candidates | Recommend review candidates from price, cost, fee, stock, and margin bands without automatic price writes. |
| F-03 | Stock aging and promotion candidates | Identify stock with low turnover, high quantity, or old movement for operator review. |
| F-04 | Kit and bundle opportunity engine | Suggest possible bundles from sales history, complementary groups, stock, and margin quality. |

## Dependencies

- M-05 Stock Seguro.
- M-06 Orders + Margin.
- Sales history and product grouping from M-03.

## Risks

- Recommendations can look authoritative without enough data quality.
- Kits/bundles may need business review and product compatibility knowledge.

## Done Means

- Recommendations always show evidence and quality.
- No automatic price/listing write occurs in this milestone.
- Low-confidence suggestions are labeled or blocked.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator after M-06.
- Next action: Reconfirm business priority before execution.
- Required files/evidence: F-*/validation.md and M-07/validation-result.md.
- Blockers or open decisions: Margin minimum thresholds and kit rules require operator decision before F-01/F-04.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-07/validation-result.md.
- Revalidation evidence required: recommendation tests and UI evidence.
