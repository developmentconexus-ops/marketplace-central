# M-06-orders-margin-ml

```yaml
id: M-06
type: milestone
status: correction_needed
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-13
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

- Current status: Historical fixed-SHA review passed at `1eb8831fb1d0d1b84f4d1325978bbc4f76c9ed0f`; proportional QA failed.
- Next owner: Mission Strategist/Portfolio for historical preservation only.
- Next action: Reuse accepted order/margin capabilities in M-13 and M-14; do not dispatch another M-06 correction under the MVP scope.
- Required files/evidence: existing F-*/validation.md, terminal checkpoint, fixed-SHA review, and `validation-result.md`.
- Blockers or open decisions: C03 trusted-principal/authorization remains unproved and is intentionally deferred to post-MVP M-11.

## Correction Handoff

- QA failure summary: C03 remained deliberately failing and the registered Go QA command used invalid Windows/cwd paths, so C01/C02 were not independently reverified.
- Correction scope: None authorized in M-06 after the 2026-07-13 MVP replan.
- Attempts used/remaining: historical ledger remains 2/0; it is not reset.
- Next artifact: `../M-13-integrated-operator-workspaces/milestone.md` after M-09.
- Revalidation evidence required: M-14 vertical real-read/browser gate, not another M-06 production-auth gate.

## MVP Replan Disposition

- M-06 is not relabeled Pass and its QA history is not rewritten.
- Its implemented orders, profitability, Sankhya linkage, and UI evidence are reusable inputs.
- Production authentication/authorization and provider-write durability move to M-11.
