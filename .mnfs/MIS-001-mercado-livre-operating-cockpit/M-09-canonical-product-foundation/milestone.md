# M-09-canonical-product-foundation

```yaml
id: M-09
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit MVP replan.

## Outcome

Every MVP product is anchored by positive Sankhya `CODPROD`; EAN, manufacturer
reference, and provider-observed seller SKU remain separate nullable identifiers;
catalog reads no longer require the legacy MSDB runtime and unknown cost/price/stock
remain null with source quality.

## Why This Milestone Exists

Product 360 cannot join listings, stock, sales, and margin safely while `product_id`
is a legacy string and catalog values are zero-filled from a second database path.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Canonical product identity contract | Align domain/API/SDK identity and nullable source facts with IC-003 |
| F-02 | Oracle catalog cutover | Replace blocking MSDB catalog reads with internal_read/Oracle-backed product reads |
| F-03 | Catalog compatibility and cutover | Preserve enrichment/classification compatibility, remove MSDB composition/config, and prove consumers use CODPROD |

## Dependencies

- M-03 Oracle internal-read contract and real read-only lane.
- M-04 product link state and provider listing identity.
- IC-003 product/listing identity semantics.

## Risks

- Existing enrichments/classifications reference string product IDs.
- Legacy query uses zero defaults for missing cost/price/stock.
- Product links may contain values derived from manufacturer or provider references.

## Done Means

- Catalog `internal_product_id` is a positive integer equal to CODPROD.
- EAN/reference/provider SKU cannot occupy the internal product identity field.
- Catalog API/SDK represents missing numeric source facts as null plus quality/source.
- `MS_DATABASE_URL`, `MS_TENANT_ID`, and `platform/msdb` are absent from active composition.
- Existing classifications/enrichments either migrate by proven CODPROD equality or
  remain explicitly unlinked; no guessed mapping is accepted.

## Handoff

- Current status: Ready for user-started Milestone execution.
- Next owner: Milestone Orchestrator.
- Next action: Follow `execution-guide.md` and dispatch F-01 first.
- Required files/evidence: M-09 feature `validation.md` files and `validation-result.md`.
- Blockers or open decisions: Stop if legacy IDs cannot be deterministically mapped to CODPROD.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: None until QA reports a named failed criterion.
- Attempts used/remaining: 0/2.
- Next artifact: `F-01-canonical-product-identity/feature.md`.
- Revalidation evidence required: M-09 contract criteria.
