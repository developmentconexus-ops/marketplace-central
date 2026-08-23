# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) + [B10 P6](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) — **B00/B01/B00-R2/B11/B12 LOCKED; B10 SUSPENDED; NOTIF-01 D2-R6 NEXT / P9 PAUSED** |
| NOTIF-01 | [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) **ACCEPTED** · [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) + [D3-R2](engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md) **PASS** · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) + [D5-R3](engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-RATIFICATION.md) **OPERATOR-RATIFIED** · [D5-R4](engineering/rebaseline/D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md) **PROVED** · [D6-R](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md) + [P8 ratification](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md) **APPROVED/LOCKED** · [P9 supersession](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md) **OPERATOR-RATIFIED; D1-R2 PASS; D2-R6 NEXT** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **104 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only NOTIF-01 D2-R6 `AuthorizationRequest` identity/lifecycle. OAD remains 104/31. Do not begin D3/D5 wire repair, resume P9/B10, D7-R/D8-R or Product implementation first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; NOTIF-01 D1-R ACCEPTED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED baseline; NOTIF-01 D2-R6 TARGETED REOPEN / NEXT |
| D3 — Communication / Events | ACCEPTED / CLOSED; NOTIF-01 D3-R ACCEPTED / D3-R1+R2 PASS |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; NOTIF-01 D5-R3 ACCEPTED / D5-R4 OAD 104/31 PROVED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED; NOTIF-01 P8 LOCKED / P9 PAUSED BY D2-R6** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED; NOTIF-01 D7-R BLOCKED BY D6-R** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; NOTIF-01 D8-R BLOCKED BY D7-R** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — NOTIF-01 D2-R6 NEXT; P9 PAUSED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Canonical OAD remains **104/31**; historical 95/29 + 99/30 proof remains intact.
- P8 B00-R2/B11/B12 remain operator-`LOCKED`; reviewed HTML snapshots are unchanged.
- P9 query-only remedy is superseded: repeated authorization episodes can share one target/revision, so pre-decision `AuthorizationRequest` identity requires targeted D2 repair.
- D1-R2 confirms `Controlled Action Governance` remains the correct owner; no domain/edge change.
- D2-R6 exact identity/lifecycle is the only next adjudication. D3/D5/P9/D7-R/D8-R remain blocked behind it.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
