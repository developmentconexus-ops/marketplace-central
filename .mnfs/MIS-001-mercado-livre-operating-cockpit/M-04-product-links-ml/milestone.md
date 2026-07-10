# M-04-product-links-ml

```yaml
id: M-04
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

Existing Mercado Livre listings and variations can be linked to internal `CODPROD` with confidence states, conflict handling, and audit.

## Why This Milestone Exists

Stock writes, margin, and repricing are unsafe unless the listing/product relationship is resolved and unambiguous.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Listing snapshot import | Fetch and persist normalized Mercado Livre listing snapshots for existing listings. |
| F-02 | Link candidate engine | Generate candidates by exact EAN and `seller_sku`, then title heuristic with confidence states. |
| F-03 | Link resolution workflow | Add API/SDK/UI workflow to approve, reject, resolve conflict, and audit product links. |

## Dependencies

- M-02 capability framework.
- M-03 product read contract for candidates.

## Risks

- EAN/SKU reuse or dirty listing titles can create false matches.
- Variation-level links need careful identity semantics.

## Done Means

- Resolved links are unique for provider listing/variation where policy requires.
- Conflicts block stock writes.
- Operator can see evidence behind a candidate.

## Handoff

- Current status: passed.
- Next owner: Milestone Orchestrator.
- Next action: start M-05 Stock Seguro using explicit `product_links` states as downstream safety gates.
- Required files/evidence: F-*/validation.md and M-04/validation-result.md.
- Blockers or open decisions: None at milestone scope.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-04/validation-result.md.
- Revalidation evidence required: link conflict and audit tests.
