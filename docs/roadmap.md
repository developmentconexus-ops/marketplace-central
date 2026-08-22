# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED** |
| D6-R2 authority | [D6-R2 Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1 IA Reopen](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8 Ledger](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **IA-01 adjudicated; OP-READ-01 OPEN; B10 SUSPENDED** |
| Execution method | [Frontend Product Experience Planning Method v2.1](development/frontend-product-experience-planning-method.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates OP-READ-01. Do not edit the OAD or re-render B00/Visão operacional before approval.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — IA-01 partial B00 reopen; B01 content LOCKED; B10 SUSPENDED; OP-READ-01 OPEN** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current D6-R2 result

- IA-01 replaces the mixed `OPERAÇÕES` mass with **OFERTA / OPERAÇÃO / ESTRATÉGIA E INTELIGÊNCIA / CONTROLE / CONFIGURAÇÕES**; user-facing `Publicações` becomes **Anúncios**, and **Preços** is first-class under Oferta.
- Candidate `OPERAÇÃO` = **Visão operacional / Vendas / Expedição / Pós-venda**. The landing is a hybrid action cockpit; a global cross-owner Kanban is rejected.
- **OP-READ-01**: owner-local Materialization/Fulfillment/Shipment collections need bounded triage enrichment before that cockpit can be rendered honestly without N+1/client workflow projection.
- B00 physical/context/responsive shell and B01 content remain LOCKED; B00 grouping is reopened. B10's internal pattern remains valid but is SUSPENDED under `OFERTA > Preparação`.
- D5 remains CLOSED and the OAD remains unchanged until explicit OP-READ-01 approval. If approved, repair only D5-B2 W2/W3 + OAD, preserve **99/30/H-A-S**, and revalidate affected GF-02 read/composition properties.

## D8 authority carried forward

D8 remains accepted. [Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md), [Proof Closure](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md), [Probe Protocol](engineering/rebaseline/D8-LIVE-PROBE-PROTOCOL.md) and [Probe Evidence](engineering/rebaseline/D8-LIVE-PROBE-EVIDENCE.md) remain reachable.

Implementation-readiness retains **P2** redeferred to the first qualifying real ML Sale/beta drive and **P5** capability narrowing: contact reference is not a supported full alternate street/fiscal destination override.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
