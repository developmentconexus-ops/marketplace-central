# Interface Contract

```yaml
id: IC-04
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

pricing (M-07, dono) ↔ FE Simulador (M-07) ↔ orders/Pedidos (M-08, consumidor read-only) ↔ Produto Detalhe/box veredicto (M-06).

## Why This Contract Exists

Decomposição de margem e tabela DIFAL apareceriam em 3 telas com 3 implementações (defeito do mock: DIFAL sempre SP). Fonte única aqui.

## Resources Or Entities

- **CalcProfile** (tenant): regime (`SIMPLES` default 4% | `PRESUMIDO` default 9,25%), aliquota_pct (editável), limiar_verde_pct (default 18), limiar_amarelo_pct (default 10), tarifa_full (editável, nullable=unknown), difal_enabled (bool), difal_destino_uf (UF|null).
- **DifalRate**: uf (27), interna_pct, interestadual_pct (regra: 12% p/ MG,PR,RJ,RS,SC,SP; 7% resto; origem fixa SC), efetivo_pct = max(interna−interestadual, 0), origem_versao `padrao-2026`, override nullable {interna_pct, updated_at} — override persiste só se Δ>0,049pp. Valores seed de `interna_pct` por UF: dataset citado em `research/difal-interna-rates-2026.md` (R-04; 26/27 verified, MS = verify-at-execution — override do operador esperado). Rótulo obrigatório em toda superfície: "seed padrão 2026 — não é orientação fiscal".
- **Decomposition** (fórmula ÚNICA, Simulador E Pedidos): `retorno = preço − comissão(pct_modalidade×preço) − taxa_fixa(preço<79 ⇒ 6,50; senão 0) − frete(preço≥79 ⇒ frete_produto; senão 0) − imposto(aliquota×preço) − difal(efetivo_uf×preço, se enabled e destino conhecido) − tarifa_full(modalidade `full` ⇒ CalcProfile.tarifa_full, null ⇒ componente UNKNOWN propaga; demais modalidades ⇒ 0 explícito) − custo_erp`; `margem_pct = retorno/preço`. `tarifa_full` é componente explícito da decomposição retornada. Comissão/frete vêm de leituras ML vivas (IC-06/listings) — NUNCA hardcoded. Chip margem: verde ≥ limiar_verde · âmbar ≥ limiar_amarelo · vermelho abaixo.
- **Scenario**: id, codprod, listing_id nullable, modalidade (`classico`|`premium`|`full`), preço, snapshot da decomposição, created_at.

Input desconhecido (custo, frete, comissão, uf) ⇒ retorno/margem UNKNOWN (propaga; UnknownValue na UI; veredicto vira estado bloqueante). Nunca 0.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `GET/PUT /pricing/profile` | drawer Parâmetros | CalcProfile | CalcProfile | PUT valida aliquota 0–35% |
| `GET /pricing/difal` | drawer/config | — | DifalRate[27] ordenada por UF | |
| `PUT /pricing/difal/{uf}` | drawer exceção | interna_pct 0–35 | DifalRate | Δ≤0,049pp ⇒ 200 sem persistir override |
| `POST /pricing/simulations` | painel Simular | codprod, preço OU margem alvo, modalidade, destino | Decomposition + veredicto | bidirecional: margem→preço via busca binária |
| `POST/GET/DELETE /pricing/scenarios` | Salvar cenário | Scenario | Scenario[] newest-first | |
| read port `DifalForUF(uf)` | orders M-08 | uf destino do shipment | efetivo_pct + versão | chip Pedidos = SÓ leitura deste port |

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| aliquota fora de 0–35 | 422 | `INVALID_RATE` | |
| UF inválida | 404 | `UF_NOT_FOUND` | |
| simulação sem custo conhecido | 200 | — | decomposição com campos unknown + blocking_state (não é erro) |
| margem alvo inatingível (busca binária não converge) | 200 | — | resultado `UNREACHABLE_TARGET` |
| simulação com preço ≤ 0 | 422 | `INVALID_PRICE` | |
| simulação de item (codprod) inexistente | 404 | `ITEM_NOT_FOUND` | |
| cenário salvo inexistente | 404 | `SCENARIO_NOT_FOUND` | |

## Must Not Decide In Feature Execution

Fórmula, defaults, limiares, regra 12/7, semântica de override, formato do chip.

## Validation Impact

UF de exceção ajustada reflete em Simulador E chip Pedidos (mesma fonte); teste destino ≠ SP; toggle off ⇒ recálculo sem DIFAL; custo unknown ⇒ margem unknown (não 0).
