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

Order (replanned 2026-07-16, mission Parallel Execution Plan): **F-02 first** (W1, inside
CHIP-SAT — depends only on M-01), F-01 later (W3, rebases on merged F-02). Both write the
shared `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime`, and server
composition root — the original F-01 → F-02 serialization guarded that; it is now handled by
disjoint OpenAPI sections + the additive composition-root lock (see Ownership & Concurrency).

## Dependencies

M-03 (listing_edit envelope), M-04 F-01 (produto/anúncio detail surfaces host the flow entry). F-02 depends only on M-01 module precedent.

## Ownership & Concurrency

Split execution (mission Parallel Execution Plan): **F-02 runs in W1 inside CHIP-SAT**
(depends only on M-01 — contract-only, no adapter, no UI); F-01 runs in W3 as CHIP-M06
after M-03, M-04 F-01 (host surface), and merged F-02.

- F-02 (in CHIP-SAT, W1): migration block **0043–0045** reserved (do not exceed; more →
  `REQUEST`). OpenAPI/SDK sections = market + category-attribute paths (additive; never
  touch CHIP-M03's mutation sections or M-05 F-01's dashboard/orders/sync sections).
  Additive contract-lock: composition-root registration lines for the market module.
- F-02 closes at feature grain (CHIP-SAT reports its `CLOSED`); this milestone stays open
  until W3.
- F-01 (CHIP-M06, W3): consumes merged F-02 contracts; AppRouter/nav route rows disjoint
  from CHIP-M05's but files shared — hub serializes merges (rebase-then-merge).
- Governance base anchor: pinned per chip at dispatch (profile §2).

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
