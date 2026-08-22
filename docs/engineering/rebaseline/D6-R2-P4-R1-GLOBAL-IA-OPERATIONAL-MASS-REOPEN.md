# D6-R2 P4-R1 — Global IA / Operational Mass Reopen

> **Status:** OPEN / OPERATOR-ADJUDICATED DIRECTION — IA-01 confirmed; OP-READ-01 OPEN
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **Method:** [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Trigger

During rendered B10 Preparation review, the operator identified that the accepted `OPERAÇÕES` group mixed two different human intents:

```text
pre-sale offer construction / maintenance
  Preparação · Publicações · Disponibilidade · PriceIntent

post-sale execution
  Vendas · materialização/faturamento · Expedição · Pós-venda
```

This is **material user-model evidence**, not visual preference. It falsifies only the global functional grouping that placed both masses under one `OPERAÇÕES` label.

The frontend method requires IA to follow human mental models/tasks rather than backend topology. Therefore the smallest affected authority is reopened.

## 2. IA-01 — material falsifier

**Finding:** the prior `OPERAÇÕES` mass does not express one coherent user intention.

Operator-adjudicated replacement masses:

```text
VISÃO GERAL

OFERTA
  Preparação
  Anúncios
  Preços
  Disponibilidade

OPERAÇÃO
  Visão operacional
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Anúncios
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
```

Mental-model questions:

```text
OFERTA                    o que e como estamos colocando para vender?
OPERAÇÃO                  o que aconteceu e o que precisa ser executado agora?
ESTRATÉGIA E INTELIGÊNCIA como estamos performando e que decisão comercial devemos tomar?
CONTROLE                  o que precisa ser coordenado ou aprovado?
CONFIGURAÇÕES             como o funcionamento permitido está configurado?
```

## 3. Terminology / route disposition

Operator-approved direction:

- user-facing `Publicações` becomes **Anúncios**;
- **Preços** is a first-class destination under Oferta and remains Offering-owned `PriceIntent` execution;
- Market / Economics / Performance may inform a price decision but do not gain price-write authority;
- **Vendas** remains the user-facing marketplace-sale concept; it is not renamed `Pedidos`, because `MarketplaceSale != BusinessOrderIntent`;
- **Expedição** and **Pós-venda** remain specialist destinations under Operação;
- **Trabalho** remains under Controle because it is cross-owner coordination, not source truth;
- existing route identities should remain where possible to avoid non-value technical churn:

```text
Anúncios               -> /publicacoes/*
Preços                 -> /publicacoes/precos
Performance / Anúncios -> /performance/publicacoes
Preparação             -> /preparacao
Disponibilidade        -> /disponibilidade
Vendas                 -> /vendas/*
Expedição              -> /expedicao/*
Pós-venda              -> /pos-venda/*
```

The only candidate new route identity is the operational landing, e.g. `/org/:organizationId/operacao`; exact route spelling remains P4/P9 navigation authority and is not Product API authority.

## 4. Lock impact

The reopen is deliberately bounded:

```text
B00 physical shell / Organization / Installation / responsive laws -> remain LOCKED
B00 global grouping / nav labels                                 -> REOPENED by IA-01
B01 Overview content/state hierarchy                             -> remains LOCKED
B01 will later inherit the corrected global nav                  -> structural shell delta only
B10 Preparation internal P6 pattern                              -> remains valid
B10 placement in global IA                                       -> SUSPENDED until corrected B00 nav is re-rendered
B20+                                                               -> BLOCKED
```

No React/topology/API/Permission/owner/runtime decision is reopened by IA-01 itself.

## 5. Operational landing structural hypothesis

The operator approved **cockpit hybrid oriented to action**, not one global Kanban.

```text
VISÃO OPERACIONAL
1. PRECISA DE ATENÇÃO
   explicit exceptions / Work / ambiguity / post-sale attention

2. TRABALHO OPERACIONAL NORMAL
   actionable work that is not itself an exception

3. ACOMPANHAMENTO
   in-progress / dispatched / delivered states without immediate action

4. SPECIALIST ENTRY POINTS
   Vendas · Expedição · Pós-venda · Trabalho
```

A global `Nova -> Faturar -> Separar -> Conferir -> Embalar -> Enviar` Kanban is rejected because no cross-owner Product lifecycle owns those columns. GF-02 intentionally composes MarketplaceSales, BusinessSystemMaterialization, Fulfillment, Shipment, PostSale and Work without a transversal workflow owner.

A Kanban-like representation remains a legitimate later hypothesis **inside Fulfillment/Expedição**, where one owner possesses explicit physical checkpoints.

## 6. OP-READ-01 — owner-local operational read gap

**Status:** OPEN / operator-approved finding for bounded analysis before operational landing wireframe.

The user need is now evidenced, but several current collection reads do not carry enough owner-local state for efficient, honest triage without N+1 detail fan-out or frontend business projection.

Observed gaps:

### Fulfillment

`GetFulfillmentExecution` owns:

```text
separation
physical_conference
packing
dispatch_handoff
physical_readiness
provider_dispatch_deadline
```

but `FulfillmentExecutionListItem` exposes only identity/sale/scope/node, `physical_readiness` and `created_at`; `ListFulfillmentExecutions` currently filters by sale and node only. A UI cannot safely derive paginated `A separar / A conferir / A embalar / A despachar` queues from the current collection contract without detail fan-out and client classification.

### Business-System Materialization

`BusinessOrderIntent` owns `external_effect_state + convergence`; its list item currently omits convergence.

`InvoicingIntent` owns `external_effect_state + convergence` and correlation to BusinessOrder/Fulfillment; its list item omits convergence and richer sale/action correlation.

### Shipment

`Shipment` owns `sale`, `state`, `dispatch_deadline`, `observed_at`; `ShipmentListItem` currently exposes only shipment identity, state and observation time, and the collection has no state/deadline triage filter.

### Post-Sale / Work

These are comparatively queue-capable already: PostSale admits lifecycle narrowing; Work admits lifecycle, responsibility, assignment and origin filters.

## 7. Repair law

OP-READ-01 does **not** authorize:

```text
GET /operational-dashboard
new OperationalWorkflow owner
screen-shaped aggregate API
frontend-computed authoritative workflow state
first-page counts presented as global counts
N+1 detail fan-out as the implementation baseline
```

The next action is to derive the smallest owner-local read-contract enrichment needed for the evidenced cockpit/queues. Prefer enriching existing list-item projections and filters; preserve existing semantic owners and, if possible, the 99-operation census. Any accepted Product-contract change must revalidate the affected D5 proof and GF-02 before D6-R2 continues.

## 8. Exact next action

Derive and operator-review the bounded **OP-READ-01 repair candidate**. Do not edit the canonical OAD until that candidate is explicitly approved. Do not re-render B00 or the Visão operacional before the read gap is dispositioned.
