# M-06-corrigir-atributo-market-contracts

```yaml
id: M-06
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-003 — `../mission.md`; contracts: IC-03 (`listing_edit`), IC-04 (`../research/market-data-interface-contract.md`), R-04 (G1–G7, scraping forbidden), ADR-14 (market contract-only) / ADR-12 (`listings` module hosts the category-attributes read).

## Outcome

Two closing slices. **Corrigir-atributo mini-flow**: category-attribute read endpoint (ML category attributes via connectors capability, cached), and a guided fix flow on listing pendências — operator picks flagged attribute gap, edits value with attribute-schema validation, applies via `listing_edit` envelope. **Market module contract-only**: `market` Go module skeleton with IC-04 tables, `GET /market/observations` + `GET /market/references` returning honest `no_price_evidence` empties, `CollectorPort` defined with test-double proving the contract — NO production adapter, NO seed data (G1 OAuth still FAILED; scraping forbidden per ML ToS 7.6).

## Why This Milestone Exists

Corrigir-atributo is the wireframe's remaining actionable pendência (quality gaps block ML exposure). Market contract-only locks IC-04 seams now so the future collector mission plugs in without reshaping — while honestly shipping zero fake market data.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | corrigir-atributo-flow | Category-attribute read + guided attribute-fix UI via listing_edit envelope |
| F-02 | market-contract-module | market module tables/endpoints/CollectorPort, contract tests, no adapter |

Order F-01 → F-02 sequential: both features write the shared `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime`, and server composition root (one writer per seam; F-02 rebases on accepted F-01).

## Dependencies

M-03 (listing_edit envelope), M-04 F-01 (produto/anúncio detail surfaces host the flow entry). F-02 depends only on M-01 module precedent.

## Risks

- Attribute schema variance across categories: validation is schema-driven from ML attribute metadata; unmappable constraints surface as provider_validation at apply, never client-guessed.
- Market module misuse: endpoints live but empty — UI consumers (M-04 /precos placeholders) must render evidence_state honestly; contract test pins empty behavior.

## Done Means

Mini-flow applies a real attribute fix via protocolo (stub lane; live lane operator-gated) and market contract proven per `validation-contract.md` (M06-C01..C08); governance lanes green; grep proves no market production adapter or seed.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 and F-02 (parallel allowed).
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: none.

## Correction Handoff

Not applicable at planning time.
