# F-01-orders-projection-api

```yaml
id: F-01
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

Projeção enriquecida de pedidos: retorno líquido + decomposição por pedido (motor pricing via port IC-04), SLA/atraso/rastreio (shipment via port IC-06), chip DIFAL por UF destino, custo ERP via Reader. KPIs + filtros/tabs. Extensão ADITIVA de `/orders*`.

## Inputs

- OpenAPI atual `/orders*` (list, import, sankhya-linkage — preservar).
- Ports consumidos: `Decompose`/`DifalForUF` (M-07 F-01, IC-04), `GetShipmentInfo`/`GetFreeShippingCost` (M-02 F-01, IC-06), `GetCostAsOf` (IC-02), vínculo item→CODPROD via product_links (API/port público).
- Migrations bloco 0060–0064 (colunas/tabela de projeção se cache necessário). Tabela/coluna nova carrega `tenant_id`; toda query nova escopa `tenant_id` (invariante da missão).
- Design handoff (campos da lista/drawer/KPIs — dita o shape) + leitura R-02 `research/design-screens-2026-07-17.md`.

## Expected Output

- `GET /orders` aditivo por pedido: `retorno_liquido: decimal|null`, `margem_pct: number|null`, `decomposicao: {...IC-04...}|null`, `componentes_desconhecidos: []`, `sla: {due, delayed: bool|null}`, `destino_uf: string|null`, `difal: {aliquota, valor}|null`, `rastreio: {status, tracking}|null`, `vinculo_status: OK|SEM_VINCULO`.
- `GET /orders/summary`: KPIs (pedidos hoje/semana, retorno líquido agregado SÓ dos conhecidos + contagem de desconhecidos, atrasados, sem vínculo) + agregado DIFAL do período `{total_difal, pedidos_com_difal}`.
- Filtros aditivos: `?status_tab=`, `?exception=atrasado|sem_vinculo|margem_negativa`, ordenação por SLA.
- Timeline por pedido (`GET /orders/{id}` aditivo): eventos de status conhecidos (criado/pago/enviado/entregue — do dado existente + shipment) com timestamps; evento desconhecido omitido, nunca inventado.
- Comprador: APENAS forma mascarada na lista/detalhe (`primeiro nome + inicial`, cidade/UF) — dado completo NÃO trafega nos endpoints de lista.
- EARS: While pedido tem item vinculado + custo + payment + shipment, when lista responde, the sistema shall incluir decomposição completa com soma exata. While item sem vínculo CODPROD, when lista responde, the sistema shall enviar `retorno_liquido: null` + `vinculo_status: SEM_VINCULO` (nunca margem sem custo real). While shipment indisponível no ML, when enriquecimento roda, the sistema shall enviar `sla: {due: null, delayed: null}` + rastreio null (desconhecido honesto).

## Inputs/Outputs

Decomposição por pedido usa valores REAIS do pedido (comissão/frete/taxas do payload de order/payment quando existirem) e completa faltantes via motor IC-04; campo `decomposicao.origem` distingue `real|simulado` por componente. Codes: 200 lista; pedido inexistente 404; filtro inválido 422.

## Negative Scenarios

- Pedido multi-item com vínculos parciais ⇒ decomposição por item; total do pedido null se qualquer item sem custo (parcial explícito por item).
- `ErrRateLimited` no shipment ⇒ campos shipment null + telemetria; lista NUNCA falha por enriquecimento.
- UF destino ausente ⇒ `difal: null` (não SC default, não zero).
- Pedido cancelado ⇒ excluído dos KPIs de retorno (filtro de estado cancelado NÃO é servido no MVP — P5-F-19; projeção de cancelados = MIS-005).

## Constraints

- ADITIVO estrito em `/orders*`; consumidores atuais intocados.
- Motor de decomposição NUNCA reimplementado — port M-07 (duplicação = defeito).
- PII: LGPD — mascaramento no backend (não no front); logs sem dados de comprador.
- Enriquecimento ML tolerante a falha (estados null), sem retry-storm.

## Ownership

- Owned paths: `modules/orders/**`, `apps/server_core/migrations/0060*–0062*` (+ fixture `apps/server_core/internal/platform/migrate/runner_test.go` bump), seção `/orders*` OpenAPI (aditiva), `sdk-runtime/src/orders.ts` (novo).
- Forbidden paths: `modules/pricing/**`, `modules/connectors/**`, `modules/product_links/**` (consumo via ports), barrel SDK, `apps/web/**` (F-02).
- Parallel-safe with: none — depends on ports M-02 F-01 + M-07 F-01 (hub gate no merge desses trechos); F-02 depende da OpenAPI deste F.

## Validation Expectations

- Transcript GET /orders com seed cobrindo: decomposição completa (soma conferida manualmente no transcript), SEM_VINCULO, shipment null, difal null — JSON exato por caso.
- Transcript summary: KPIs batendo com a seed (contagens exatas + desconhecidos contados à parte).
- Grep/inspeção payload: nenhum campo de comprador não-mascarado nos responses de lista.
- Golden test: decomposição de pedido real (lane live-provider-read) com componentes `real` vs `simulado` marcados.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-08). Complexidade: candidata a Sol low (projeção multi-port + unknown matrix — hub decide).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + transcripts.
- Blockers or open decisions: none.
