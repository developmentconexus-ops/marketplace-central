# F-04-vinculos-precos

```yaml
id: F-04
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

MIS-003. Binding contracts: IC-05 (routes `/vinculos` `/precos`, linkage/pricecost ns, states), IC-03 (`link_apply`, `price_update`), ADR-008 link states (unresolved|conflict|resolved|rejected), R-01 screens 1g + 1i/2d (minus @mercado), R-02 (ProductLinks + Simulator facts).

## Milestone

M-04 read-workspaces. Depends on M-03 (link_apply/price_update) + F-03 precedent. Last M-04 feature; retires remaining legacy aliases.

## Brief

Two rebuilds, one seam pass. **Vínculos (1g)** at `/vinculos`: link table by ADR-008 state (candidate suggestions, conflitos, resolvidos, rejeitados), existing product_links endpoints via TanStack Query (linkage ns), apply/reject actions through M-03 modal `link_apply` intents, import panel parity from legacy page. **Preços & Simulador (2d minus @mercado)** at `/precos`: simulation table (pricecost ns 120s) from existing profitability/pricing endpoints, política de preço display, margem with UnknownValue when cost null ("não simulado"), bulk "Aplicar preço" via M-03 modal `price_update` (competitive/market columns render only IC-04 evidence_state placeholders — no market adapter this mission). Then retire legacy aliases: legacy page components deleted, only IC-05 redirects remain; audit nav links repo-wide point to new routes.

EARS:
- While links are in `conflict`, when the Vínculos conflitos tab renders, each row shall show both sides (ERP produto vs anúncio) and offer resolver/rejeitar through the modal.
- While a link_apply protocol reaches terminal, when observed, `invalidateAfterMutation('link_apply')` shall invalidate listings+linkage+catalog+mutations (crosswalk row proof in situ).
- While cost is null, when simulation renders, margem shall show "não simulado" — never a number; bulk price selection shall EXCLUDE such rows with visible count "n sem custo excluídos".
- While market evidence is absent (always, this mission), when competitive columns render, they shall show `no_price_evidence` copy per IC-04 — no fake competitor prices.
- While any legacy alias (`/products`, `/product-links`, `/inventory/stock-seguro`, `/orders`, `/integrations`, `/simulator`) is visited, redirect shall preserve query; no legacy component code remains in bundle.

## Inputs

- R-01 §1g + §2d, R-02 ProductLinks + Simulator facts (parity checklists), product_links + profitability/pricing endpoints (sdk-runtime), IC-03 intents, IC-04 evidence_state enum, M-03 modal, IC-05.

## Expected Output

- `/vinculos` + `/precos` pages, TanStack Query only; legacy ProductLinks + Simulator components deleted.
- Bulk price flow excluding null-cost rows with visible exclusion count.
- Grep-audited nav-link sweep (no internal link targets legacy paths).
- Component tests: conflict render, crosswalk in-situ, null-cost exclusion, evidence placeholder, redirect suite (all 6 rows re-proven post-deletion).

## Constraints

- No market data endpoints called beyond IC-04 read stubs (contract-only; M-06 wires module).
- Simulator computation parity: same numbers as legacy for same inputs (golden test from R-02 sample).
- Writes only via M-03 modal; pt-BR.

## Negative Scenarios

- link_apply on already-resolved link → per IC-03 Error Matrix: same outcome → item `skipped` (no-op); conflicting re-resolution → item fails `conflict_remote_changed`. UI shows honest per-item result.
- Simulation with zero-priced policy → validation stops intent pre-preview ("política sem valor").
- Redirect regression: visiting each of the 6 legacy aliases post-deletion → correct new route, no 404.

## Interaction Model

Two pages, one feature: shared modal + crosswalk; selection semantics identical to Anúncios (composite id). URL encodes tab/filter per page. Freshness: linkage default, pricecost 120s.

## Validation Expectations

- Vitest output: all listed component tests + golden simulation parity test green.
- Browser proof: 1g conflitos tab screenshot; 1i table with "não simulado" rows + exclusion count; redirect walk of all 6 aliases.
- Grep proof: legacy components absent; no nav link to legacy paths.
- Protocolo proof: link_apply + price_update stub flows terminal with correct invalidations.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-03 accepted).
- Next action: compile context pack; read R-01 §1g/§2d + R-02 sections + IC-03/IC-04/IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: link_apply-on-resolved semantics pinned in spec.md from product_links API (bounded).
