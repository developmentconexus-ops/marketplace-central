# D6-R2 P4-R1 — Global IA / Operational Mass Reopen

> **Status:** ACCEPTED / CLOSED — IA-01 operator-ratified; OP-READ-01 RESOLVED; corrected B00 global IA `LOCKED`
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **Method:** [Frontend Product Experience Planning Method v2.1](../../development/frontend-product-experience-planning-method.md)
> **Bounded Product repair:** [D5-R2 Operational Read Projection Repair](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md)
> **GF-02 revalidation:** [D8-R2 Operational Read Revalidation](D8-R2-OPERATIONAL-READ-REVALIDATION.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Trigger

During rendered B10 Preparation review, the operator identified that the accepted `OPERAÇÕES` group mixed two different human intents:

```text
pre-sale offer construction / maintenance
  Preparação · Publicações · Disponibilidade · PriceIntent

post-sale execution
  Vendas · materialização/faturamento · Expedição · Pós-venda
```

This was **material user-model evidence**, not visual preference. It falsified only the global functional grouping that placed both masses under one `OPERAÇÕES` label.

The frontend method requires IA to follow human mental models/tasks rather than backend topology. Therefore the smallest affected authority was reopened and has now been re-ratified.

## 2. IA-01 — accepted replacement

The prior `OPERAÇÕES` mass did not express one coherent user intention. The accepted replacement is:

```text
VISÃO GERAL
  Visão geral

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
  Configurações
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

Accepted direction:

- user-facing `Publicações` becomes **Anúncios**;
- **Preços** is first-class under Oferta and remains Offering-owned `PriceIntent` execution;
- Market / Economics / Performance may inform a price decision but do not gain price-write authority;
- **Vendas** remains the user-facing marketplace-sale concept; it is not renamed `Pedidos`, because `MarketplaceSale != BusinessOrderIntent`;
- **Expedição** and **Pós-venda** remain specialist destinations under Operação;
- **Trabalho** remains under Controle because it is cross-owner coordination, not source truth;
- existing technical route identities remain where possible:

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

The organization-wide operational landing is the only new frontend route identity candidate, e.g. `/org/:organizationId/operacao`; exact route spelling remains frontend navigation authority and is not Product API authority.

## 4. Lock disposition

```text
B00 physical shell / Organization / Installation / responsive laws -> LOCKED
B00 corrected global grouping / nav labels                          -> LOCKED
B01 Overview content/state hierarchy                                -> LOCKED
B10 Preparation internal P6 pattern                                 -> remains valid
B10 global placement                                                -> OFERTA > Preparação
B20+                                                                  -> not opened by this closeout
```

The corrected executable HTML was reviewed and explicitly approved by the operator on 2026-08-22. IA-01 is closed. No React/topology/Permission/owner/runtime decision was reopened by IA-01.

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

A global `Nova -> Faturar -> Separar -> Conferir -> Embalar -> Enviar` Kanban remains rejected because no cross-owner Product lifecycle owns those columns. GF-02 composes MarketplaceSales, BusinessSystemMaterialization, Fulfillment, Shipment, PostSale and Work without a transversal workflow owner.

A Kanban-like representation remains a legitimate later hypothesis **inside Fulfillment/Expedição**, where one owner possesses explicit physical checkpoints.

## 6. OP-READ-01 — RESOLVED

The evidenced operational cockpit initially exposed an owner-local read gap: several collection contracts lacked enough state for efficient, honest triage without N+1 detail fan-out or frontend business projection.

The operator approved the smallest owner-local repair and D5-R2 proved it executably while preserving the Product census and authorities.

Accepted repair summary:

```text
BusinessOrderIntent list
  + convergence
  + external_effect_state / convergence narrowing

InvoicingIntent + list
  + source-qualified Sale correlation
  + convergence
  + optional FulfillmentExecution correlation in list
  + external_effect_state / convergence narrowing

FulfillmentExecution list
  + separation / physical_conference / packing / dispatch_handoff
  + optional provider dispatch deadline
  + owner-native readiness/checkpoint/deadline narrowing

Shipment list
  + optional Sale + dispatch deadline
  + state narrowing
```

No repair was required for MarketplaceSales, PostSale or Work.

### 6.1 Preserved negative controls

```text
no /operational-dashboard Product endpoint
no OperationalWorkflow owner
no cross-owner synthetic lifecycle
no operational_stage / next_action / priority / total_count
no global totals inferred from one paginated page
no N+1 detail fan-out as production baseline
```

### 6.2 Proof and revalidation

[D5-R2](D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) proves:

```text
99 Product operations
30 ordinary Permissions
H/A/S only
owner-local projection proof PASS
owner-local filter proof PASS
negative controls 2/2
```

[D8-R2](D8-R2-OPERATIONAL-READ-REVALIDATION.md) confirms GF-02 choreography, owners, writes/effects, physical authority and source-qualified outcomes remain unchanged. Only read projection/filter expressibility was enriched.

## 7. Global-Maximum / YAGNI decision law

For the remainder of D6-R2, every material decision follows:

```text
accepted Product/system authority
→ evidenced human job / mental model
→ smallest frontend structure that serves it
→ owner/wire trace
→ repair only the smallest owning authority when evidence falsifies it
→ no local screen optimization that creates parallel Product truth
→ no speculative capability/platform work
```

A better screen is not sufficient reason to change Product. A proven user job requiring an already-owned truth that is not consumable may justify the smallest owner-local repair, as OP-READ-01 demonstrated.

## 8. B00-R1 closeout

**Artifact:** [`qualification/d6-r2-wireframes/b00-app-shell.html`](../../../qualification/d6-r2-wireframes/b00-app-shell.html)

```text
rendered artifact reviewed: YES
operator disposition:       LOCKED
material changes requested: NONE
physical/context shell:     LOCKED / unchanged
corrected global IA:        LOCKED
visual-design decisions:    NONE
```

## 9. Exact next action

Frontend progression is intentionally paused for the operator-requested **Notification architecture assessment**. Treat the MetalDocs notification direction only as comparative evidence. Do not import Notification semantics, operations, owner, event contract, realtime mechanism or UI into Marketplace Central unless current Product authority demonstrates a real gap and the resulting bounded reopen is separately operator-approved.
