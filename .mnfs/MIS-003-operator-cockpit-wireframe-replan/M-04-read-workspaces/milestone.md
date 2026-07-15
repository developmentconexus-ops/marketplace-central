# M-04-read-workspaces

```yaml
id: M-04
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

MIS-003 — `../mission.md`; contracts: IC-05 (routes/keys/states), IC-02 (by-product), IC-01 dormant row (ACTIVATES here), ADR-15/17.

## Outcome

Four wireframe workspaces rebuilt on the M-02 platform: Detalhe do produto (2b) at `/catalogo/produtos/:productId`, Catálogo (1f) at `/catalogo`, Estoque (1h) at `/estoque`, Vínculos & Importação (1g) at `/vinculos`, plus `/precos` (1i/2d minus @mercado) rebuilt from Simulator and deprecated legacy aliases retired. Product enrichment editing activates the IC-01 dormant row (server class `catalog` + client `['catalog']` invalidation). Observable: each screen deep-links + reloads faithfully; product detail shows anúncios do produto via IC-02 by-product; edits produce envelope protocols (link_apply/listing_edit) or catalog enrichment writes with correct invalidation.

## Why This Milestone Exists

These four screens are the operator's daily read surface. They only consume seams built in M-01/M-02/M-03 — proving the platform holds under fan-out and retiring the direct-fetch legacy pages that motivated ADR-15.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | produto-detalhe | Detalhe do produto 2b: cadastro, completude vs qualidade, anúncios do produto, enrichment edit (IC-01 activation) |
| F-02 | catalogo-workspace | Catálogo 1f rebuild on TanStack Query; replaces /products |
| F-03 | estoque-workspace | Estoque 1h: físico/reservado/disponível/seguro, correções via envelope |
| F-04 | vinculos-precos | Vínculos 1g (link_apply via envelope) + /precos rebuild + legacy alias retirement |

F-01 first (defines produto-page patterns others link into). F-02/F-03/F-04 then sequential — each touches AppRouter route rows + shared nav links; keep one writer.

## Dependencies

M-02 (platform + shell), M-03 (envelope writes for corrections/links). IC-02 by-product endpoint from M-01.

## Risks

- IC-01 activation regression: catalog invalidation proven by dormant-row test finally running non-dormant.
- Legacy parity gaps (milestone-local risk, no mission RK ID): each rebuilt screen carries a parity checklist from R-01/R-02 in its spec; missing legacy behavior is a defect unless mission Non-Scope names it.

## Done Means

Four workspaces live per `validation-contract.md` (M04-C01..C10); legacy routes redirect; zero direct-fetch in rebuilt pages; IC-01 invalidation proven; governance lanes green.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 with IC-02/IC-05 + R-01 §2b.
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: none.

## Correction Handoff

Not applicable at planning time.
