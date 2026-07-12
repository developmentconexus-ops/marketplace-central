# M-06-orders-margin-ml

```yaml
id: M-06
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

Mercado Livre orders are ingested idempotently and each sale can show contribution/margin quality using revenue, sale fee, `CUSSEMICM`, taxes, and manual freight/commission adjustments where needed.

## Why This Milestone Exists

After stock risk is controlled, operators need to understand whether Mercado Livre sales are actually profitable and which missing inputs reduce confidence.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Order ingestion | Poll/read Mercado Livre orders and persist normalized order/item/payment/shipment/cancellation snapshots. |
| F-02 | Margin input model | Define revenue, sale fee, cost, tax, freight, commission/manual adjustment inputs and quality states. |
| F-03 | Profit snapshot calculation | Calculate per-order/per-item margin with `CUSSEMICM` as cost basis and explicit missing flags. |
| F-04 | Orders and margin UI | Build sale detail and margin-quality views for operators. |
| F-05 | Docker live-Oracle validation runner | Run only the registered Oracle SELECT validation target in the canonical Linux/CGO/Instant Client image with explicitly injected lane variables. |

## Dependencies

- M-02 capability framework.
- M-03 MNOS/Sankhya read contract.
- M-04 Product Links ML for margin resolution quality.
- M-05 is recommended before this milestone because stock is the urgent operational risk.

## Risks

- Mercado Livre fee/freight fields may be incomplete for some order shapes.
- Manual freight/commission entry can drift without audit.
- Missing product links can reduce margin quality but must not block order ingestion.

## Done Means

- Orders are idempotent by provider order id and provider update timestamp.
- Missing link/cost/freight/tax creates quality flags.
- Profit snapshots never use zero defaults for unknown inputs.
- UI shows margin and confidence side by side.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: Start after Stock Seguro is usable or when operator prioritizes margin.
- Required files/evidence: F-*/validation.md and M-06/validation-result.md.
- Blockers or open decisions: Manual adjustment UX details can be refined during F-02/F-04.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-06/validation-result.md.
- Revalidation evidence required: order ingestion, margin, quality flag, and UI tests.
