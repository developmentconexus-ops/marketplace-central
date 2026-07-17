# F-02-pedidos-ui

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-08-pedidos.

## Brief

Tela `/pedidos` real (substitui placeholder W1 via `routes/pedidos.tsx`): KPIs clicáveis, Fila de atenção rankeada, Lista com tabs por estado, Kanban read-only, drawer com decomposição de margem + timeline + rastreio + comprador mascarado, banner DIFAL agregado.

## Inputs

- Design handoff (tela Pedidos) + leitura R-02 `research/design-screens-2026-07-17.md`.
- `sdk-runtime/src/orders.ts` (F-01).
- IC-05 (tokens/MarginChip/FreshnessIndicator/DataTable/DetailDrawer), IC-04 (labels de decomposição).

## Expected Output

- KPIs topo (summary): pedidos hoje/semana, retorno líquido (+ "N sem margem conhecida"), atrasados, sem vínculo — clique aplica filtro correspondente na lista (URL).
- Fila de atenção: ranking por severidade (atrasado > margem negativa > sem vínculo), cards com ação direta (abrir drawer / ir a /vinculos).
- Lista: tabs por estado (Todos/Pagos/Enviados/Entregues; `Cancelados` = tab DISABLED "em breve" — P5-F-19, sem fetch), colunas do design (pedido, comprador mascarado, itens, retorno líquido com MarginChip, SLA com badge atraso, DIFAL chip quando aplicável), ordenação por SLA.
- Kanban read-only: colunas de estado, cards arrastáveis NÃO (somente leitura no MVP).
- Drawer: decomposição completa (componente + valor + origem real|simulado), timeline vertical, rastreio, itens com vínculo/CODPROD, banner "margem desconhecida" quando aplicável.
- Banner DIFAL agregado do período no topo da lista (total + link simulador).
- EARS: While KPI atrasados clicado, when filtro aplica, the lista shall mostrar só atrasados e URL refletir (`?exception=atrasado`, deep-linkável). While pedido SEM_VINCULO aberto no drawer, when renderiza, the decomposição shall mostrar bloqueio "sem vínculo — resolver em /vinculos" com link. While tab Cancelados (disabled) é clicada, when clique ocorre, the sistema shall não navegar e manter affordance "em breve" (nenhum fetch disparado).

## Negative Scenarios

- Lista vazia (conta demo sem pedidos) ⇒ EmptyState com orientação runbook (pedido de teste).
- Summary falha ⇒ KPIs ocultos, lista funciona.
- Pedido com decomposição parcial ⇒ componentes conhecidos renderizados, faltantes "—" com motivo, chip neutro.

## State / Interaction Model

- Tab/filtros/ordenação = URL params (deep-link + F5); drawer = `?order=` na URL.
- Keys: `['orders','list',params]`, `['orders','summary']`, `['orders','detail',id]`; refetch manual único (botão) invalida os três.
- Kanban deriva da MESMA query da lista (agrupamento client-side) — zero fetch extra, zero drag.

## Constraints

- Read-only total: nenhuma ação de escrita em pedido no MVP.
- Nenhum dado de comprador além do mascarado da API (front não "desmascara").
- Zero cálculo de margem no front (exibe API).

## Ownership

- Owned paths: `apps/web/src/pages/pedidos/**` (novo), `apps/web/src/routes/pedidos.tsx`.
- Forbidden paths: `apps/web/src/app/**`, `packages/ui/**`, outros routes/pages, `sdk-runtime/**`, backend.
- Parallel-safe with: none — depends on F-01 (`orders.ts`) + M-03 seams.

## Validation Expectations

- Screenshot lista com MarginChip nos 4 estados + badge atraso + chip DIFAL; drawer com decomposição origem real|simulado visível.
- Transcript: clique KPI ⇒ URL + lista filtrada; F5 mantém.
- Deep-link `?order=` + F5 ⇒ drawer aberto.
- Kanban espelha contagens das tabs (mesmos números).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-08, após F-01 OpenAPI).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots + transcripts.
- Blockers or open decisions: none.
