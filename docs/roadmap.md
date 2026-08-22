# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED** |
| D6-R2 authority | [D6-R2 Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1 IA Reopen](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P5 Inventory](engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md) + [P8 Ledger](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **IA-01 operator-adjudicated direction; OP-READ-01 OPEN; B10 SUSPENDED** |
| Execution method | [Frontend Product Experience Planning Method v2.1](development/frontend-product-experience-planning-method.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates the derived OP-READ-01 repair candidate. Do not edit canonical OAD before explicit approval; do not re-render B00 or Visão operacional first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | **ACCEPTED / CLOSED — global IA grouping partially reopened only through D6-R2 finding IA-01** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; GF-02 must be revalidated if OP-READ-01 changes Product contract** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00 global IA partial reopen; B01 content LOCKED; B10 SUSPENDED; OP-READ-01 OPEN** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current D6-R2 result

- P0–P3 remain valid: **4 actor contexts → 16 needs → 16 complete flows → 99/99 Product coverage** under the current Product surface.
- **IA-01 is a material mental-model falsifier**: prior `OPERAÇÕES` mixed pre-sale offer work with post-sale execution.
- Operator-approved P4-R1 direction is now **OFERTA / OPERAÇÃO / ESTRATÉGIA E INTELIGÊNCIA / CONTROLE / CONFIGURAÇÕES**.
- User-facing `Publicações` becomes **Anúncios**; **Preços** becomes first-class under Oferta while remaining Offering-owned PriceIntent execution; existing technical routes should remain where possible.
- `OPERAÇÃO` gains one candidate landing, **Visão operacional**; Vendas, Expedição and Pós-venda remain specialist destinations.
- Global operational Kanban is rejected: GF-02 crosses Sales, Materialization, Fulfillment, Shipment, PostSale and Work without one transversal workflow owner.
- Operational landing candidate is a **hybrid action cockpit**: attention/exceptions → normal actionable work → monitoring → specialist entry points.
- **OP-READ-01** blocks that wireframe: several owner-local collections lack enough triage state/filtering to avoid N+1 detail fan-out or frontend-authored workflow projection.
- B00 physical/context/responsive shell remains LOCKED; only global grouping/labels are reopened. B01 content/state hierarchy remains LOCKED.
- B10 internal `search/list → selected exact-subject detail` evidence remains valid but B10 is **SUSPENDED** under `OFERTA > Preparação` until corrected B00 global IA is re-rendered.
- D5 remains ACCEPTED / CLOSED while OP-READ-01 is only a repair candidate; no Product operation, Permission, semantic owner or OAD has been changed by the finding yet.

## OP-READ-01 candidate guardrail

Leading repair shape, awaiting operator adjudication:

```text
BusinessOrderIntent list
  expose convergence + typed state/convergence filters

InvoicingIntent list
  expose convergence + necessary owner-local correlation + typed state/convergence filters

FulfillmentExecution list
  expose checkpoint states + provider dispatch deadline
  add typed readiness/checkpoint/deadline filters

Shipment list
  expose sale + dispatch deadline
  add typed state/deadline filters

PostSale / Work
  no repair currently justified
```

Binding negative controls:

```text
no new operation
no new Permission
no new semantic owner
no screen-shaped /operational-dashboard
no generic sort/filter DSL
no total_count
no cross-owner synthetic lifecycle
```

If the operator approves this direction, reopen only the smallest D5-B2 W2/W3 + canonical OAD authority needed, preserve the **99-operation / 30-Permission / H-A-S** census, prove generated projections and revalidate affected GF-02 read/composition properties. Live external-effect probes are not repeated unless a write/effect invariant is actually changed.

## D8 authority carried forward

D8 remains accepted. [Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md), [Proof Closure](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md), [Probe Protocol](engineering/rebaseline/D8-LIVE-PROBE-PROTOCOL.md) and [Probe Evidence](engineering/rebaseline/D8-LIVE-PROBE-EVIDENCE.md) remain reachable.

Implementation-readiness retains **P2** redeferred to the first qualifying real ML Sale/beta drive and **P5** capability narrowing: contact reference is not a supported full alternate street/fiscal destination override.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
