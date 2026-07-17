# P7 Claude readiness — r04 (targeted, operator-authorized)

```yaml
round: r04
mode: targeted (NOT full crew) — operator instruction 2026-07-17 "Corrija isso barato e revise especifico" after 3-round cap
manifest: planning-reviews/p7-input-r04.sha256
digest: 8da3bfb7f567fa8f92eb548c7813d69d47a42ceacdfb8345ad632690fac3cca0
reviewers: 1 cold mission-reviewer (sonnet), scoped
scope: (A) ★3 H1-repair verification; (B) ★2 consistency of all post-r03 edits
```

## Pre-review repairs (advisories r03, applied before manifest freeze)

1. mission.md ADR-06 — enums `match_status` (IC-01) vs `blocking_state` (IC-03) separados com valores nomeados.
2. M-02 VC C02 — mesma separação (ver FAIL abaixo — 1ª versão do repair estava errada).
3. M-04 F-01:38 + M-04 VC C02 — "quem" = ator fixo `operator` (MVP sem auth).
4. M-01 VC C01 — critério log-payload (grep fixture ⇒ zero hits) + 409 `DUPLICATE_FILE` nomeado.
5. mission.md ledger — bullet AL transição de alíquota (R-04, override operador).
6. R-02 (`research/design-screens-2026-07-17.md`) citado nos Inputs de 9 briefs FE: M-03 F-02/F-03, M-04 F-02, M-05 F-02, M-06 F-01, M-07 F-02, M-08 F-01/F-02, M-09 F-01.

## Result

- Manifest check: OK (digest match).
- **Scope A (★3, H1): PASS.** Token `packages/feature-simulator/**` idêntico em mission.md:146, M-07 milestone.md:43, F-02:65, M-07 VC:138; F-01 não reivindica; nenhum outro dono. H1 fechado.
- **Scope B (★2, 6 itens): 5 PASS, 1 FAIL.**
  - Item 2 FAIL — M-02 VC C02 (linha 52): conflacionava `price_evidence_status`/`blocking_state` sob um value-set único de 4 valores. IC-01:44 define `price_evidence_status` ∈ {OK, NO_PRICE_EVIDENCE, INSUFFICIENT_MARKET} (3 valores, com OK, sem NO_CANDIDATE/SEM_CUSTO) — enum DIFERENTE de `blocking_state` (IC-03, 4 valores, sem OK). Yes-if: separar em três enums nomeados e value-matched.
  - Itens 1, 3, 4, 5, 6: PASS com excertos citados (ver dispatch record / transcript).

## Repair disposition

- Item-2 FAIL: **APPLIED pós-review** — M-02 VC C02 reescrito exatamente per yes-if: três enums distintos (`match_status` IC-01 / `price_evidence_status` IC-01 / `blocking_state` IC-03), value-sets fiéis às fontes. Mudança pontual, autoridade local (consolidação de contrato ratificado). NOTA: essa edição invalida o digest r04 — qualquer rodada formal futura exige novo manifest.

## Advisories (não-bloqueantes)

- r03 advisories fechadas nesta rodada: M-01 log-payload, M-04 "quem", AL ledger.
- Ainda abertas: PR FCP/FECOEP wording em `research/difal-interna-rates-2026.md:37` (documentado como discrepancy note no próprio arquivo — skip decidido); M-03 F-02/F-03 citam R-02 na mesma bullet do handoff (cosmético).

## Standing

Revisão targeted concluída; todos os findings blocking de r01–r04 aplicados. NÃO é verdict `Ready` formal (rodada fora do protocolo full-crew, cap atingido em r03). `status: planned` requer dual gate (crew full + Sol HIGH) ou decisão explícita do operador aceitando a rodada lean como suficiente — decisão pendente do operador.
