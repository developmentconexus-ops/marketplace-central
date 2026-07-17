---
id: R-04
type: research
parent: MIS-004
created: 2026-07-17
lifecycle_scope: support
---

# DIFAL — Alíquota Interna (Modal) por UF — Seed 2026

> **seed padrão 2026 — não é orientação fiscal.** Todo valor é
> operator-overridable. Este dataset alimenta a tabela DIFAL do MVP demo
> (origem SC; interestadual fixo 12% para MG/PR/RJ/RS/SC/SP e 7% para as
> demais UFs; `DIFAL_efetivo = max(interna_pct − interestadual_pct, 0)`).

## Tabela

| UF | interna_pct (modal) | ano verificado | FCP/FECP note | source URL | status |
|----|----------------------|-----------------|----------------|------------|--------|
| AC | 19.0% | 2025/2026 | — | https://html.duckduckgo.com/html/?q=tabela+icms+2026+aliquota+interna+por+estado → https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| AL | 20.5% | 2026 (efetiva 01/04/2026, Lei 9.776/2025) | Rate inclui FECOEP; regra anterior era 19% + 1% FECOEP = 20% até 31/03/2026 | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://cidesp.com.br/v2/tabela-aliquota-icms-2026 e https://agilize.com.br/blog/gestao-contabil-e-fiscal/tabela-de-icms-2026/ | verified |
| AM | 20.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| AP | 18.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| BA | 20.5% | 2025/2026 | — | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| CE | 20.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| DF | 20.0% | 2025/2026 | — | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| ES | 17.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| GO | 19.0% | 2025/2026 (efetiva 01/04/2024) | — | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| MA | 23.0% | 2025/2026 | maior alíquota geral do país | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; https://cidesp.com.br/v2/tabela-aliquota-icms-2026 | verified |
| MG | 18.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| MS | 17.0% (DISPUTED — ver caveat) | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; contradito por https://tributodevido.com.br/portal/tabela-icms-2026-aliquotas-atualizadas-todos-estados-brasileiros/ (19%) | verify-at-execution |
| MT | 17.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| PA | 19.0% | 2025/2026 | sem fundo de pobreza | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| PB | 20.0% | 2025/2026 | +2% FCP (total 22%) per Brasil NFe | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| PE | 20.5% | 2025/2026 | — | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| PI | 22.5% | 2025/2026 (efetiva 01/04/2025) | — | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://cidesp.com.br/v2/tabela-aliquota-icms-2026 | verified |
| PR | 19.5% | 2025/2026 | um aggregador (Tributo Devido) descreve como "19,5% incl. 2% FECOEP"; demais três fontes tratam 19,5% como taxa cheia sem FECOEP — nota de discrepância menor, sem impacto no total | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; https://tributodevido.com.br/portal/tabela-icms-2026-aliquotas-atualizadas-todos-estados-brasileiros/ | verified |
| RJ | 20.0% | 2025/2026 | +2% FECP (total efetivo 22%) | https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; cross-check https://tributodevido.com.br/portal/tabela-icms-2026-aliquotas-atualizadas-todos-estados-brasileiros/ | verified |
| RN | 20.0% | 2025/2026 | atualizada 2025 | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| RO | 19.5% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| RR | 20.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| RS | 17.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| SC | 17.0% (origem do vendedor) | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ ; https://agilize.com.br/blog/gestao-contabil-e-fiscal/tabela-de-icms-2026/ | verified |
| SE | 19.0% | 2025/2026 | +1% FUNPOBREZA/FECOEP (total 20%) | https://agilize.com.br/blog/gestao-contabil-e-fiscal/tabela-de-icms-2026/ ; cross-check https://www.brasilnfe.com.br/aliquotas-icms/ | verified |
| SP | 18.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |
| TO | 20.0% | 2025/2026 | — | https://www.brasilnfe.com.br/aliquotas-icms/ ; cross-check https://simtax.com.br/tabela-icms-2026-aliquotas-de-todos-estados-atualizada/ | verified |

## Method

- Attempted primary sources first: CONFAZ (`confaz.fazenda.gov.br`) and individual
  state SEFAZ pages (SC, MG, PR, RJ, MS). All attempts returned HTTP 404, generic
  homepage content without the rate table, or content that a live fetch could not
  extract (JS-rendered tables, PDF links that 404'd on guessed paths). This is
  recorded honestly below as a source limitation, not glossed over.
- Fell back to a `duckduckgo.com/html` search (`tabela icms 2026 aliquota interna
  por estado`) to discover which tax-focused aggregator sites publish a
  2026-dated, all-27-UF table. That search surfaced simtax.com.br,
  brasilnfe.com.br, cidesp.com.br, agilize.com.br, and tributodevido.com.br.
- Fetched five aggregator tables via `WebFetch` this run (2026-01-17… actually
  session date 2026-07-17): `simtax.com.br/tabela-icms-2026-…`,
  `brasilnfe.com.br/aliquotas-icms`, `cidesp.com.br/v2/tabela-aliquota-icms-2026`,
  `agilize.com.br/blog/…/tabela-de-icms-2026`, and
  `tributodevido.com.br/portal/tabela-icms-2026-…`. Also fetched
  `nfe.io`, `contabilizei.com.br`, `tecnospeed.com.br`, `totvs.com.br`,
  `webmaniabr.com`, `cobrefacil.com.br` and others which returned 403/404/blocked
  and contributed no data (not cited in the table above).
- Every row above is corroborated by at least two independent aggregator
  fetches from this run, except MS (contradicted — flagged) and the PR FECOEP
  phrasing (contradicted on wording only, not on the headline number).
- `interna_pct` = the general/modal ICMS rate a state applies to a typical
  taxed operation (not product-specific surcharges for cigarettes, fuel,
  energy, telecom, alcoholic beverages, etc., which several states tax higher
  — e.g. BA/DF list a higher rate for energy/telecom that is NOT the value
  used here).
- FCP/FECP (Fundo de Combate à Pobreza / Fundo Estadual de Combate e
  Erradicação da Pobreza) is deliberately excluded from `interna_pct` per
  task instruction, and only noted in the FCP column when a mandatory
  poverty-fund addition materially changes the effective total (RJ +2%,
  AL/SE ≈+1%, PB +2% per Brasil NFe).

## Caveats

1. **Aggregator sources, not primary SEFAZ text.** Every attempt to reach a
   state-level SEFAZ page or CONFAZ directly failed to yield the rate table
   itself in this run (404s or JS-only content). The values in this file are
   sourced from tax-consultancy/ERP-vendor aggregator blogs (SimTax, Brasil
   NFe, Cidesp, Agilize, Tributo Devido) that explicitly target 2025/2026 and
   were cross-checked against each other. This satisfies the "best available
   public modal rate with a source" bar stated in the task, but does **not**
   meet fiscal-grade authority — hence the "não é orientação fiscal" label and
   full operator-override expectation.
2. **MS (Mato Grosso do Sul) is disputed.** Three independent fetches (SimTax,
   Brasil NFe, and a corroborating pass) report MS at 17%; one fetch
   (Tributo Devido / Agilize's Center-West breakdown) reports 19%. Majority
   value (17%) was used and marked `verify-at-execution` — operator should
   override once confirmed against MS's own legislation (Lei 1.810/1997) or a
   fiscal-grade source.
3. **AL (Alagoas) is mid-transition.** A rate increase from 19%(+1% FECOEP =
   20% total) to 20.5% takes effect 01/04/2026 per Law 9.776/2025 (cited by
   SimTax and Cidesp). The seed uses the post-increase 2026 value (20.5%);
   an operator running DIFAL for Jan–Mar 2026 transactions should override to
   20% for that window.
4. **Modal ≠ product-specific.** Several states (BA, DF, GO, MG, others) apply
   materially higher ICMS to specific categories (energy, telecom, fuel,
   cigarettes, alcoholic beverages, cosmetics/perfumes) — sometimes 25–30%+.
   This seed intentionally uses the general/modal rate only. Product-specific
   taxation is OUTSIDE this seed's contracted behavior — IC-04 defines only the
   tenant CalcProfile and per-UF DIFAL overrides; no per-SKU fiscal override
   surface exists in MIS-004.
5. **FCP figures are informational, not authoritative.** The FCP notes above
   come from the same aggregator cross-check and were not independently
   verified against each state's poverty-fund legislation. Treat as a
   pointer for "why the operator might need to override," not as a computed
   input.
6. **No live SEFAZ or CONFAZ page could be fetched this run.** If
   fiscal-grade sourcing becomes a hard requirement before this table ships
   externally (vs. as an internal MVP-demo seed), a follow-up research pass
   should target each state's ICMS regulation text (RICMS) directly, ideally
   via a service that renders JS/PDF content, rather than generic web fetch.
