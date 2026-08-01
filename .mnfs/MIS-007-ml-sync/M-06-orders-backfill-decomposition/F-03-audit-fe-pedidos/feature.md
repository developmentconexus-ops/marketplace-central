# F-03-audit-fe-pedidos

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-06 orders-backfill-decomposition.

## Brief

Auditoria 3→2 + contrato/FE: (a) auditor pós-decomposição compara comissão OBSERVADA na
linha (camada 3, TOTAL) vs ESPERADA computada do detail da camada 2 do anúncio vinculado:
esperado_unit = detail.percentage_fee × unit_price/100 + detail.fixed_fee +
detail.financing_add_on_fee; esperado_total = esperado_unit × quantity — TODA componente
de fee do detail canônico IC-01 entra (dropar financing_add_on_fee abriria divergência
falsa PERMANENTE em anúncio com parcelamento — auditoria P5 r05 F-r05-2) —
|delta| > R$0,01 (IC-01 tolerância) → divergência kind=`tarifa`
(entity_type=order_line, IC-02) com os dois lados + timestamps; convergência →
auto-resolve. Sem camada 2 p/ o anúncio → NÃO avalia (auditoria muda, nunca falso alarme).
(b) DTO de /orders (list + get) ganha ADITIVO: `margem_pct`, `liquido`, `decomposicao`
(objeto canônico IC-03), `divergences[]`; par OpenAPI + SDK (`listOrders` `index.ts:2288`,
`getOrder` `:2290`) MESMO commit. (c) FE /pedidos: coluna MARGEM na lista (NULL → `—`),
PedidoDrawer ganha seção "Decomposição financeira" (linha a linha: receita, comissão,
frete, custo, líquido, margem; `incompleto[]` renderizado como avisos nomeados) + badge de
divergência tarifa; fila/bucket INTOCADOS (bucket já persistido M-03).

EARS:
- While pedido com camada 2 disponível e delta > R$0,01, when auditor roda, the sistema
  shall abrir divergência tarifa com expected/observed.
- While /pedidos lista, when pedido tem decomposição, the row shall mostrar margem_pct
  formatada; when incompleto, `—` com tooltip dos campos faltantes.
- While drawer abre, when decomposição existe, the seção shall mostrar soma-das-partes
  legível e os avisos de incompleto[].

## Inputs

F-02 (decomposição + camada 3); M-05 F-01 (camada 2 — dependência SOFT: sem ela auditor
muda); IC-01/IC-02/IC-03 (binding); `PedidoDrawer.tsx` estrutura atual (fato #2 de
`research/p5-prerequisites.md` —
CompradorFiscalSection `:355-368` é o idioma de seção); DESIGN-REFERENCE @8144238.

## Expected Output

Auditor (application orders) + DTO/transport aditivo + par OpenAPI+SDK + FE lista/drawer.

## Constraints

- ADITIVO estrito no contrato — golden do DTO baseline.
- Comparação 3→2 na MESMA unidade (TOTAL da linha dos dois lados): camada 3 value é TOTAL;
  esperado vem do DETAIL da camada 2 (percentage_fee × unit_price/100 + fixed_fee +
  financing_add_on_fee, × qty — F-r05-2)
  — NUNCA do amount da camada 2 (amount é ancorado no price_used da observação, não no
  preço do pedido; lição round5: pinar a fórmula, exato não decide a chave contada).
- Divergência avaliada por LINHA (grão order_line dos dois lados).
- FE: sem refetch novo — margem vem no payload da lista existente.

## Inputs/Outputs

DTO canonical: IC-03 §Required Outputs; chaves da decomposição = IC-03 §Persistence
Expectations (incl. proveniência de fee `*_origem`/`*_coletado_em` — ADR-09, P7 r01 B-3).
Divergência: IC-02 §Required Outputs, kind=tarifa.

## Interaction Model

- Coluna MARGEM entra na lista existente (mesma query, mesmo payload — sem fetch paralelo);
  ordenação por margem NÃO entra nesta missão (sort atual intocado).
- Drawer: seção "Decomposição financeira" segue o idioma CompradorFiscalSection
  (`PedidoDrawer.tsx:355-368`); abre com os dados do GET de detalhe existente (aditivo no
  mesmo payload — zero request novo).
- `incompleto[]` renderiza avisos nomeados DENTRO da seção (não toast, não banner global).
- Badge divergência tarifa na fila usa o slot de badge existente do item; estado stale segue
  a política react-query atual do módulo (sem invalidação nova).

## Negative Scenarios

- Anúncio da linha sem vínculo/camada 2 → zero rows de divergência (contado no log do run).
- Margem negativa → renderiza vermelho (DESIGN-REFERENCE), não esconde.
- Pedido só com incompleto[] → badge neutra, nunca divergência.

## Ownership

- Owned paths: auditor em `orders/application/`, `orders/transport/http_handler.go`
  (aditivo), par OpenAPI+SDK /orders, `apps/web/src/pages/pedidos/**`.
- Forbidden paths: listings; pricing; schema; scheduler.
- Parallel-safe with: none — depends on F-02.

## Validation Expectations

- 2 direções da divergência (planta → abre; corrige → resolve) com valores exatos.
- Golden DTO baseline + novo.
- tsc verde; live-drive hub: pedido real com margem na tela, drawer com decomposição
  conferível contra o SELECT.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-02.
- Required files/evidence: `validation.md`; screenshot-métrica live-drive (hub).
- Blockers or open decisions: none.
