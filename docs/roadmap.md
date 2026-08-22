# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 authority | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **IA-01 CLOSED; B00/B01 LOCKED; OP-READ-01 RESOLVED; B10 SUSPENDED** |
| NOTIF-01 design | [Personal Notification Architecture Design](engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md) — **PROPOSED / operator direction approved / written-spec review pending** |
| Product repair | [D5-R2](engineering/rebaseline/D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) — **ACCEPTED / CANONICAL** |
| GF-02 revalidation | [D8-R2](engineering/rebaseline/D8-R2-OPERATIONAL-READ-REVALIDATION.md) — **ACCEPTED / PASS** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator reviews/adjudicates the written NOTIF-01 architecture design. Do not reopen D0/D1, edit the OAD, alter B00 for the bell, or resume B10/B20 until the written design is explicitly approved.** |
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
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 written design pending review** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current D6-R2 result

- IA-01 is closed: **OFERTA / OPERAÇÃO / ESTRATÉGIA E INTELIGÊNCIA / CONTROLE / CONFIGURAÇÕES**; `Publicações` is **Anúncios** and **Preços** is first-class under Oferta.
- `OPERAÇÃO` = **Visão operacional / Vendas / Expedição / Pós-venda**; a global cross-owner Kanban is rejected.
- OP-READ-01 is resolved by D5-R2; Product remains **99/30/H-A-S**. D8-R2 preserves GF-02.
- **B00 App Shell + corrected global IA and B01 content are LOCKED.** B10 remains valid under `OFERTA > Preparação` but is paused for NOTIF-01.
- NOTIF-01 direction: persistent Personal Notifications owner candidate, recoverable PostgreSQL+River delivery, no broker, no Work/Audit/authorization conflation, and an Organization-scoped **bell in the topbar** as Inbox entry point. This remains design, not Product authority, until written-spec approval.
- Global-Maximum/YAGNI remains binding.

## D8 authority carried forward

D8 remains accepted. Implementation-readiness retains **P2** redeferred to the first qualifying real ML Sale/beta drive and **P5** capability narrowing: contact reference is not a supported full alternate street/fiscal destination override.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
