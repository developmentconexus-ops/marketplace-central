# F-02-precos-provenance-fe

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-07
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-07 pricing-fee-read.

## Brief

Proveniência na superfície: DTOs de /pricing (simulação single + batch + decompose) ganham
ADITIVO por componente de tarifa `origem` (vocabulário IC-01 com produtor NESTA missão:
api_listing_prices | config — `api_shipping_options` fica no CHECK do schema como reserva
aditiva futura, SEM produtor aqui) e `coletado_em` (null quando config); par OpenAPI + SDK
(métodos `runPricingSimulation` `index.ts:2359`, `runBatchSimulation` `:2361`,
`pricingDecompose` `:2374`) MESMO commit. FE /precos: junto do valor de comissão, chip de
origem (`anúncio`/`config`) + ⚠ com tooltip SÓ quando origem=config ("tarifa padrão — sem
observação do anúncio"); coletado_em exibido como data da observação, SEM regra de
staleness (limiar de stale = decisão futura, fora desta missão); DESIGN-REFERENCE tokens.
Nenhum valor de tarifa aparece sem origem.

EARS:
- While simulação usa camada 2, when tela mostra comissão, the chip shall dizer `anúncio` +
  data da observação (coletado_em).
- While origem=config, when tela mostra comissão, the chip shall dizer `config` com ⚠ e
  tooltip.
- While payload antigo sem campos novos (SDK velho impossibilitado — par same-commit), the
  render shall degradar sem crash (campo optional).

## Inputs

F-01 (proveniência no domain); IC-01 (vocabulário origem); fatos #6/#9 de
`research/p5-prerequisites.md` (SDK métodos
pricing verbatim, NENHUM método tariff-defaults — fora de escopo aqui); DESIGN-REFERENCE
@8144238; memória `pricing-strategy-3-camadas` (decisão por margem — proveniência é o
alicerce da confiança).

## Expected Output

Transport pricing aditivo + par OpenAPI+SDK + FE /precos (componentes de simulação).

## Constraints

- ADITIVO estrito — golden do DTO baseline de simulação.
- ⚠ SÓ para origem=config — nenhuma heurística de staleness nesta missão (auditoria P5
  A-7: limiar não-ratificado não entra).
- SÓ /precos — /anuncios tarifa é M-05, /pedidos margem é M-06 (disjunção FE da lane C).
- tsc verde no worktree com `npm ci` na raiz (profile §3).

## Inputs/Outputs

Shape aditivo por componente: `{"origem":"api_listing_prices","coletado_em":"..."}` —
exemplos canônicos na spec seguem IC-01.

## Interaction Model

- Chip de origem renderiza inline no componente de resultado de simulação existente (single
  e batch) — sem painel novo, sem rota nova.
- Tooltip do ⚠ é hover/focus (idioma de tooltip do DESIGN-REFERENCE); sem estado persistido.
- Dados vêm no MESMO payload da simulação (zero refetch; estado segue o fluxo
  request→resultado atual da tela).

## Negative Scenarios

- Simulação batch com mix de origens → cada linha carrega a SUA origem (nunca a do lote).
- origem ausente num componente não-tarifa (frete config) → chip config coerente, sem ⚠
  duplicado.

## Ownership

- Owned paths: `pricing/transport/` (DTOs aditivos), par OpenAPI+SDK /pricing,
  `apps/web/src/pages/precos/**` (componentes de simulação/resultado).
- Forbidden paths: resolver (F-01); outras rotas FE; defaults write path.
- Parallel-safe with: none — depends on F-01.

## Validation Expectations

- Golden DTO baseline + novo.
- tsc verde; live-drive hub: simulação real com chip de origem visível nos 2 estados
  (anúncio requer M-05 populado — se lane C ainda não tem camada 2, estado `config` é o
  provado + fixture prova o outro).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md`; screenshot-métrica.
- Blockers or open decisions: none.
