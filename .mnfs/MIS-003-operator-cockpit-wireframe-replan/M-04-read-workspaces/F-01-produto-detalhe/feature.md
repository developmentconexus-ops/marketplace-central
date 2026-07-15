# F-01-produto-detalhe

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/catalogo/produtos/:productId`, keys, states), IC-02 (`listListingsByProduct`), IC-01 dormant row (catalog invalidation — ACTIVATES here), IC-03 (listing_edit/link_apply intents), R-01 screen 2b. ADR-16/17.

## Milestone

M-04 read-workspaces.

## Brief

Build Detalhe do produto (wireframe 2b): header (CODPROD, título ERP, ⚑ badges catálogo/estoque/NF: ERP · preço/anúncio: HUB), cadastro section (ERP master data read-only + HUB enrichment fields editable), completude (cadastro) vs qualidade (anúncio) indicators as separate scores, anúncios do produto table via `listingsQueryKeys.byProduct`, estoque + custo panel (null → UnknownValue "custo?"), GTIN state ("sem GTIN" badge when null), and the 2b secondary tabs Pedidos (orders of this product, existing orders read via IC-05 keys) and Histórico (sync/mutation timeline rows for the product's listings). Enrichment save writes through existing catalog enrichment endpoint and MUST invalidate server class `catalog` + client `['catalog']` (IC-01). Listing rows offer per-listing actions reusing the M-03 modal (listing_edit, link_apply).

EARS:
- While a product exists, when `/catalogo/produtos/123` loads (deep link or reload), the page shall render from `productId` alone — no navigation-state dependency.
- While enrichment is saved, when the write succeeds, the client shall invalidate `['catalog']` and the UI shall reflect fresh data without manual reload (IC-01 proof).
- While ERP master fields render, when displayed, they shall be visibly read-only with ⚑ "catálogo: ERP" provenance; only HUB-owned fields editable.
- While the product has no linked listings, when the anúncios table renders, it shall show EmptyState with hint linking to `/vinculos`.
- While cost is unknown, when panels render, margem/simulação shall show "não simulado" hint, never computed-from-zero.

## Inputs

- R-01 §2b element inventory, IC-01 verbatim (R-03 carries it), existing catalog product endpoints + enrichment write (sdk-runtime), IC-02 byProduct operation, M-03 modal component, IC-05 components/keys.

## Expected Output

- Page + route element at IC-05-assigned row; breadcrumb Catálogo → produto.
- Completude/qualidade computed client-side from documented field lists (spec.md fixes the lists from R-01).
- Enrichment form with dirty-state guard ("Alterações não salvas").
- Component tests: deep-link render, IC-01 invalidation (spy), read-only enforcement, unknown-cost rendering.

## Constraints

- ERP fields never writable (⚑ mandate); no ERP write path exists or is added.
- Enrichment uses existing catalog endpoint — no new server endpoint in this feature.
- Listing mutations only through M-03 modal; no bespoke write calls.
- pt-BR; glossary terms exact (completude vs qualidade distinction preserved).

## Negative Scenarios

- Unknown productId → ErrorState 404 copy "Produto não encontrado." with link back to `/catalogo`.
- Enrichment save failure → form stays dirty, ErrorState inline, values not lost.
- byProduct API error → anúncios section ErrorState; rest of page functional (independent queries).

## State Model

Independent queries: product (catalog ns), listings byProduct (listings ns), each with own loading/error. Edit form: pristine→dirty→saving→pristine|error. Navigation with dirty form prompts confirm. IC-01: save success → invalidate catalog → product query refetches (staleTime honored via invalidation, not manual setQueryData).

## Validation Expectations

- Vitest output: IC-01 invalidation spy test green (the formerly dormant assertion now active — name test after IC-01), deep-link, read-only, unknown-cost tests.
- Browser proof: 2b screenshot with badges, completude/qualidade, anúncios table; F5 fidelity on deep link.
- Network proof: enrichment save followed by automatic product refetch (invalidation, not manual).

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; read R-01 §2b + IC-01/IC-02/IC-05 + catalog endpoint docs.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: completude/qualidade field lists fixed in spec.md from R-01 (bounded, not operator-blocking).
