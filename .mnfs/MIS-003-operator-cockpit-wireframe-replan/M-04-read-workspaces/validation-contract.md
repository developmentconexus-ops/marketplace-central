# Milestone Validation Contract

```yaml
id: M-04
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-04-read-workspaces

## QA Level

QA-2 — vitest + browser walkthrough with screenshots; envelope writes on stub lane.

## Required Outcome

Detalhe do produto, Catálogo, Estoque, Vínculos, Preços rebuilt on the platform; IC-01 activated; legacy pages deleted with redirects intact.

## Criteria

## Criterion: Product detail deep link + provenance
ID: M04-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser: direct open `/catalogo/produtos/<seeded CODPROD>`; F5
- Expected: full 2b render from productId alone (header ⚑ badges "catálogo: ERP" + "preço/anúncio: HUB", completude and qualidade as two separate indicators, anúncios table via by-product); identical after F5; ERP master fields visibly read-only
- Actual:
- Artifact: `F-01-produto-detalhe/validation.md` screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: seeded product CODPROD=100001 with 2 linked listings, installation inst_1
- Steps:
  - open http://localhost:5174/catalogo/produtos/100001?installation=inst_1
  - assert text "⚑ catálogo: ERP"
  - assert text "Anúncios do produto"
- Expected: detail renders without prior navigation
Blocking failure: render depends on navigation state or ERP field editable
Blocking failure observed: No
Owner: QA Validator

## Criterion: IC-01 dormant row activated
ID: M04-C02
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: vitest invalidation-spy test named for IC-01; network trace of enrichment save
- Expected: enrichment save success → client invalidates `['catalog']` → product query auto-refetches (observed refetch request, no manual setQueryData); server L2 class `catalog` invalidated (server log/header proof)
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: save without both invalidations
Blocking failure observed: No
Owner: QA Validator

## Criterion: Catálogo rebuild parity
ID: M04-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest URL round-trip + browser walkthrough; spec.md parity checklist review
- Expected: 1f columns render; null cost shows UnknownValue "custo?" (not "R$ 0,00"); filters encode in URL and restore on F5; every legacy Products feature kept or dropped-with-reason in spec.md checklist; legacy component deleted (grep zero references)
- Actual:
- Artifact: `F-02-catalogo-workspace/validation.md`
Blocking failure: parity item silently missing or zero-rendered cost
Blocking failure observed: No
Owner: QA Validator

## Criterion: Estoque divergence + unknown-not-zero
ID: M04-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest render tests on seeded divergent + null-stock rows
- Expected: divergent row shows ConflictTag "divergente: ERP={n}" with both values in detail; null quantity renders "—" never 0; editing a stock-seguro rule for a seeded product persists the new value (subsequent GET returns the updated rule); parity anchored to the R-02 stock-seguro checklist in F-03 spec.md — each legacy behavior kept or dropped-with-reason
- Actual:
- Artifact: `F-03-estoque-workspace/validation.md`
Blocking failure: fake comparison when publish-side unknown, or 0 for null
Blocking failure observed: No
Owner: QA Validator

## Criterion: Stock correction via envelope
ID: M04-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser stub-lane: select rows → "Corrigir estoque" → confirm → terminal
- Expected: protocolo created with `stock_correct` items carrying source_timestamp from inventory fetch; terminal fires crosswalk invalidation {listings, inventory, mutations}; unresolved-link row fails `link_unresolved` with pt-BR failureCopy shown
- Actual:
- Artifact: `F-03.../validation.md` screenshots + protocolo JSON
Blocking failure: write path bypassing envelope
Blocking failure observed: No
Owner: QA Validator

## Criterion: Vínculos states + link_apply
ID: M04-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser: conflitos tab on seeded conflict links; apply flow on stub
- Expected: tabs map ADR-008 states (unresolved|conflict|resolved|rejected); conflict row shows both sides (ERP produto vs anúncio); link_apply terminal invalidates {listings, linkage, catalog, mutations}; import panel parity anchored to the F-04 spec.md checklist built from R-02 ProductLinks facts — each legacy import behavior kept or dropped-with-reason
- Actual:
- Artifact: `F-04-vinculos-precos/validation.md`
Blocking failure: state mislabeled or crosswalk row wrong in situ
Blocking failure observed: No
Owner: QA Validator

## Criterion: Preços simulation honesty + golden parity
ID: M04-C07
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: golden test (R-02 sample inputs) legacy vs rebuilt; null-cost render + bulk exclusion test
- Expected: same numeric outputs as legacy Simulator for identical inputs; null-cost rows show "não simulado" and are EXCLUDED from bulk price selection with visible count "n sem custo excluídos"; competitive columns show `no_price_evidence` copy (no fabricated market values)
- Actual:
- Artifact: `F-04.../validation.md`
Blocking failure: numeric divergence from legacy or null-cost row included in bulk apply
Blocking failure observed: No
Owner: QA Validator

## Criterion: Legacy retirement complete
ID: M04-C08
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: browser walk of all 6 legacy aliases post-deletion; grep for legacy components + deprecated SDK aliases + internal nav links
- Expected: each alias redirects with query preserved, no 404; Products/ProductLinks/Simulator legacy components absent from bundle; no internal link targets legacy paths; deprecated aliases `listCatalogProducts`/`searchCatalogProducts` removed if last consumer gone (else deferred to M-05 with note)
- Actual:
- Artifact: `F-04.../validation.md` grep transcripts + redirect walk
Blocking failure: orphaned legacy code or broken redirect
Blocking failure observed: No
Owner: QA Validator

## Criterion: TanStack-only in rebuilt pages (Q6)
ID: M04-C09
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: grep for direct fetch patterns in M-04-added files; governance lanes at milestone SHA
- Expected: zero direct fetches; all keys from IC-05 builders; lanes green
- Actual:
- Artifact: milestone `validation-result.md`
Blocking failure: direct fetch in rebuilt page
Blocking failure observed: No
Owner: QA Validator

## Criterion: pt-BR copy audit (Q5)
ID: M04-C10
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: copy sweep of the four rebuilt workspaces (grep for English UI literals in new files + screenshot review)
- Expected: zero English UI strings; glossary terms exact (físico/reservado/disponível/seguro; completude vs qualidade); provider raw strings only behind "▸ técnico"
- Actual:
- Artifact: milestone `validation-result.md` copy-audit section
Blocking failure: English string or glossary violation in new surface
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature proofs `F-0*/validation.md`; rollup `validation-result.md` with fixed SHA + dual-gate records.

## Blocking Failures

Any criterion blocking failure; IC-01 regression (enrichment save without catalog invalidation) blocks regardless.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01..F-04 accepted).
- Next action: execute criteria at fixed SHA.
- Required files/evidence: as above.
- Blockers or open decisions: none.
