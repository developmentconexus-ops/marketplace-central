# F-01-dashboard-mpc

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-09
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-09-dashboard-demo (CORTÁVEL — decisão do hub na entrada da wave C).

## Brief

Visão geral `/` real: estender `/dashboard/summary` aditivamente com agregações reais e reconstruir DashboardPage nos tokens. Regra de ouro: TODO número do dashboard bate com a tela de origem (paridade verificável) — dado ausente é estado honesto, nunca zero.

## Inputs

- Design handoff (tela Visão geral) + leitura R-02 `research/design-screens-2026-07-17.md` (dashboard MPC, não multichannel).
- `/dashboard/summary` existente + `modules/dashboard/**`.
- Fontes por card (via ports/APIs públicas dos módulos donos): orders summary (vendas/retorno), listings summary (ativos/exceções), product_links (sem vínculo), erp_import (último import + idade).
- IC-05 (tokens/seams); `sdk-runtime/src/dashboard.ts` (novo).

## Expected Output

- Backend: `/dashboard/summary` aditivo — `{vendas: {hoje, semana}|null, retorno_liquido: {...}|null, anuncios: {ativos, excecoes}|null, vinculos_pendentes: n|null, ultimo_import: {at, idade}|null}` — cada bloco null-ável independente (fonte fora ⇒ bloco null, resto responde).
- Front: KPIs (mono font), Fila de atenção cross-módulo (top itens: pedidos atrasados, anúncios abaixo do custo, produtos sem vínculo — cada item deep-linka p/ tela de origem COM filtro aplicado), pedidos recentes (5), atalhos (importar planilha / simular / vínculos).
- EARS: While todas as fontes respondem, when dashboard carrega, the KPIs shall mostrar números idênticos aos das telas de origem (mesma janela temporal explícita no card). While uma fonte falha, when summary responde, the card correspondente shall mostrar ErrorState isolado (demais cards normais). While item da Fila de atenção clicado, when navegação ocorre, the tela destino shall abrir já filtrada (URL com filtro).

## Inputs/Outputs

Shape aditivo acima; janela temporal explícita por bloco (`window: {from, to}`). Codes: 200 sempre; blocos null carregam `reason: SOURCE_UNAVAILABLE|NO_DATA`.

## Negative Scenarios

- Zero pedidos na janela ⇒ `vendas: {hoje: 0, semana: 0}` é legítimo (zero REAL, com janela citada) vs fonte fora ⇒ `null` + reason — os dois estados são visualmente distintos.
- Sem import ERP ⇒ card último-import em call-to-action (não erro).
- Summary inteiro falha ⇒ página com ErrorState global + retry (nunca branco).

## State / Interaction Model

- Uma query `['dashboard','summary']`; refresh botão + refetch on focus off (demo: sem surpresa visual).
- Deep-links da Fila carregam filtro por URL param das telas destino (contratos de URL de M-05/M-08/M-04 — usar os params ratificados nos briefs deles).

## Constraints

- Leituras cross-módulo APENAS via APIs/ports públicos dos donos — dashboard não conhece tabela alheia.
- Paridade de número é critério de aceite (QA compara card vs tela).
- Se milestone cortado: nada deste brief entra; demo abre /anuncios (registro no mission close).

## Ownership

- Owned paths: `modules/dashboard/**`, seção `/dashboard*` OpenAPI (aditiva), `sdk-runtime/src/dashboard.ts`, `apps/web/src/pages/dashboard/**` (rebuild), `apps/web/src/routes/dashboard.tsx`.
- Forbidden paths: módulos-fonte (orders/listings/product_links/erp_import — consumo público), barrel SDK, `apps/web/src/app/**`, `packages/ui/**`.
- Parallel-safe with: M-06 F-01 (disjoint: rotas/páginas/módulos distintos).

## Validation Expectations

- Screenshot dashboard light+dark com KPIs + Fila de atenção.
- Transcript paridade: número do card vendas ⇒ MESMO número em /pedidos com mesma janela (dois screenshots + valores citados).
- Transcript falha isolada: fonte listings derrubada em teste ⇒ card anúncios ErrorState, demais renderizam (JSON do summary com bloco null + reason).
- Clique item da Fila ⇒ URL destino com filtro aplicado.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-09, wave C — ou corte pelo hub).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots + paridade.
- Blockers or open decisions: corte (hub/operator na entrada da wave C).
